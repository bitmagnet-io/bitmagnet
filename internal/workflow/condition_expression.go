package workflow

import (
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protobuf"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/ext"
	"reflect"
)

func celEnvOption(src workflowSource, ctx *compilerContext) error {
	options := []cel.EnvOption{
		cel.StdLib(),
		Lists(),
		cel.EagerlyValidateDeclarations(true),
		cel.ExtendedValidations(),
		ext.Strings(ext.StringsValidateFormatCalls(true)),
		cel.Types(&protobuf.Torrent{}, &protobuf.Classification{}),
		cel.Variable("torrent", cel.ObjectType("bitmagnet.Torrent")),
		cel.Variable("result", cel.ObjectType("bitmagnet.Classification")),
	}
	// `keywords` is masquerading as a map of strings to regexes, but it's actually individual string constants defined with a dot in the name,
	// along with a placeholder map of strings to nulls. This achieves correct compile-time checking with acceptable error messages.
	for group, keywords := range src.Keywords {
		r := regex.NewRegexFromNames(keywords...)
		options = append(
			options,
			cel.Constant("keywords."+group, cel.StringType, types.String(r.String())),
		)
	}
	options = append(
		options,
		cel.Constant("keywords", cel.MapType(cel.StringType, cel.NullType), types.NullValue),
	)
	// `extensions` uses a similar trick.
	for group, extensions := range src.Extensions {
		options = append(
			options,
			cel.Constant("extensions."+group, cel.ListType(cel.StringType), types.NewStringList(types.DefaultTypeAdapter, extensions)),
		)
	}
	options = append(
		options,
		cel.Constant("extensions", cel.MapType(cel.StringType, cel.NullType), types.NullValue),
	)
	options = append(
		options,
		cel.Constant("fileType.unknown", cel.IntType, types.Int(protobuf.Torrent_File_unknown)),
	)
	for _, ft := range model.FileTypeValues() {
		options = append(
			options,
			cel.Constant(fmt.Sprintf("fileType.%s", ft.String()), cel.IntType, types.Int(protobuf.NewFileType(model.NullFileType{Valid: true, FileType: ft}))),
		)
	}
	options = append(
		options,
		cel.Constant("contentType.unknown", cel.IntType, types.Int(protobuf.Classification_unknown)),
	)
	for _, ct := range model.ContentTypeValues() {
		options = append(
			options,
			cel.Constant(fmt.Sprintf("contentType.%s", ct.String()), cel.IntType, types.Int(protobuf.NewContentType(model.NullContentType{Valid: true, ContentType: ct}))),
		)
	}
	options = append(
		options,
		cel.Constant("mb", cel.IntType, types.Int(1_000_000)),
	)
	options = append(
		options,
		cel.Constant("gb", cel.IntType, types.Int(1_000_000_000)),
	)
	env, err := cel.NewCustomEnv(options...)
	if err != nil {
		return err
	}
	ctx.celEnv = env
	return nil
}

const expressionName = "expression"

type expressionCondition struct{}

var celProgramPayload = payloadTransformer[string, cel.Program]{
	spec: payloadGeneric[string]{},
	transform: func(s string, ctx compilerContext) (cel.Program, error) {
		ast, issues := ctx.celEnv.Compile(s)
		if issues != nil && issues.Err() != nil {
			return nil, ctx.error(fmt.Errorf("type-check error: %w", issues.Err()))
		}
		if !reflect.DeepEqual(ast.OutputType(), cel.BoolType) {
			return nil, ctx.error(fmt.Errorf("got %v, wanted %v output type", ast.OutputType(), cel.BoolType))
		}
		prg, prgErr := ctx.celEnv.Program(ast,
			cel.EvalOptions(cel.OptOptimize),
		)
		if prgErr != nil {
			return nil, ctx.error(fmt.Errorf("program construction error: %w", prgErr))
		}
		return prg, nil
	},
}

var expressionConditionPayload = payloadUnion[cel.Program]{
	oneOf: []TypedPayload[cel.Program]{
		payloadSingleKeyValue[cel.Program]{
			key:       expressionName,
			valueSpec: payloadMustSucceed[cel.Program]{celProgramPayload},
		},
		payloadMustSucceed[cel.Program]{celProgramPayload},
	},
}

func (c expressionCondition) name() string {
	return expressionName
}

func (c expressionCondition) compileCondition(ctx compilerContext) (condition, error) {
	prg, err := expressionConditionPayload.Unmarshal(ctx)
	if err != nil {
		return condition{}, ctx.error(err)
	}
	return condition{
		check: func(ctx executionContext) (bool, error) {
			result, _, err := prg.Eval(map[string]any{
				"torrent": ctx.torrentPb,
				"result":  ctx.resultPb,
			})
			if err != nil {
				return false, err
			}
			bl, ok := result.Value().(bool)
			if !ok {
				return false, errors.New("not bool")
			}
			return bl, nil
		},
	}, nil
}
