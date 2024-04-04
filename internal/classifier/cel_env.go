package classifier

import (
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protobuf"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/ext"
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
	// `flags` is masquerading as a map of strings to regexes, but it's actually individual string constants defined with a dot in the name,
	// along with a placeholder map of strings to nulls. This achieves correct compile-time checking with acceptable error messages.
	for name, tp := range src.FlagTypes {
		rawVal := src.Flags[name]
		val, err := tp.celVal(rawVal)
		if err != nil {
			return err
		}
		options = append(
			options,
			cel.Constant("flags."+name, tp.celType(), val),
		)
	}
	options = append(
		options,
		cel.Constant("flags", cel.MapType(cel.StringType, cel.NullType), types.NullValue),
	)
	// `keywords`, `extensions` etc use a similar trick.
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
		cel.Constant("fileType", cel.MapType(cel.StringType, cel.NullType), types.NullValue),
	)
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
		cel.Constant("contentType", cel.MapType(cel.StringType, cel.NullType), types.NullValue),
	)
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
