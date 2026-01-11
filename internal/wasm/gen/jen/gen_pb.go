package jen

import (
	"github.com/bitmagnet-io/bitmagnet/internal/wasm/gen/spec"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Go names of implementation-specific struct fields in generated messages.
const (
	State_goname = "state"

	SizeCache_goname  = "sizeCache"
	SizeCacheA_goname = "XXX_sizecache"

	WeakFields_goname  = "weakFields"
	WeakFieldsA_goname = "XXX_weak"

	UnknownFields_goname  = "unknownFields"
	UnknownFieldsA_goname = "XXX_unrecognized"

	ExtensionFields_goname  = "extensionFields"
	ExtensionFieldsA_goname = "XXX_InternalExtensions"
	ExtensionFieldsB_goname = "XXX_extensions"

	WeakFieldPrefix_goname = "XXX_weak_"
)

func GenTypes(s spec.FileInfo) *jen.File {
	f := newFile(s)

	genHeader(f, s)

	for _, enum := range s.Enums {
		f.Add(genEnum(enum))
	}

	for _, msg := range s.Messages {
		if msg.IsMapEntry {
			continue
		}

		f.Add(genMessage(msg))
	}

	for _, svc := range s.Services {
		if svc.Type == spec.ServiceUnknown || svc.Type == spec.ServiceNone {
			continue
		}
		f.Add(genService(svc))
	}

	return f
}

func genEnum(s spec.EnumInfo) *jen.Statement {
	return jen.Type().Id(s.GoIdent.GoName).Int32().Line().
		Const().DefsFunc(func(g *jen.Group) {
		for _, v := range s.Values {
			g.Id(v.GoIdent.GoName).Id(s.GoIdent.GoName).Op("=").Lit(int(v.Number))
		}
	}).Line().
		Var().DefsFunc(func(g *jen.Group) {
		g.Id(s.GoIdent.GoName + "_name").Op("=").Map(jen.Int32()).String().ValuesFunc(func(g *jen.Group) {
			for _, v := range s.Values {
				g.Line().Lit(int32(v.Number)).Op(":").Lit(string(v.Name))
			}
			g.Line()
		}).Line()
		g.Id(s.GoIdent.GoName + "_value").Op("=").Map(jen.String()).Int32().ValuesFunc(func(g *jen.Group) {
			for _, v := range s.Values {
				g.Line().Lit(string(v.Name)).Op(":").Lit(int32(v.Number))
			}
			g.Line()
		}).Line()
	}).Line().
		Func().Params(jen.Id("x").Id(s.GoIdent.GoName)).Id("Enum").Params().Op("*").Id(s.GoIdent.GoName).Block(
		jen.Id("p").Op(":=").New(jen.Id(s.GoIdent.GoName)),
		jen.Op("*").Id("p").Op("=").Id("x"),
		jen.Return(jen.Id("p")),
	).Line()
}

func genMessage(s spec.MessageInfo) *jen.Statement {
	st := jen.Type().Id(s.GoIdent.GoName).StructFunc(func(g *jen.Group) {
		g.Id(State_goname).Qual(pkgProtoImpl, "MessageState")
		g.Id(SizeCache_goname).Qual(pkgProtoImpl, "SizeCache")
		g.Id(UnknownFields_goname).Qual(pkgProtoImpl, "UnknownFields")

		for _, f := range s.Fields {
			g.Add(genMessageField(f))
		}
	}).Line().
		Func().Params(jen.Id("x").Op("*").Id(s.GoIdent.GoName)).Id("ProtoReflect").Params().Qual(pkgProtoReflect, "Message").Block(
		jen.Panic(jen.Lit("not implemented")),
	).Line().Line()

	for _, field := range s.Fields {
		// todo: oneof
		if field.LeadingComments != "" {
			st.Comment(string(field.LeadingComments)).Line()
		}

		st.Func().
			Params(jen.Id("x").Op("*").Id(s.GoIdent.GoName)).
			Id("Get" + field.GoName).
			Params().Add(genGoType(field)).Block(
			jen.If(jen.Id("x").Op("!=").Nil()).Block(
				jen.Return().Id("x").Dot(field.GoName),
			).Line().Line().
				Var().Id("zero").Add(genGoType(field)).Line().Line().
				Return().Id("zero"),
		).Line().Line()
	}

	return st
}

func genMessageField(s spec.MessageFieldInfo) *jen.Statement {
	st := &jen.Statement{}

	if s.LeadingComments != "" {
		st = jen.Comment(string(s.LeadingComments)).Line()
	}
	if s.IsDeprecated {
		st = st.Comment("Deprecated: do not use.").Line()
	}
	tag := map[string]string{
		"protobuf": s.ProtobufTagValue,
		"json":     string(s.Name) + ",omitempty",
	}

	if s.IsMap {
		tag["protobuf_key"] = s.MapKeyValue[0].ProtobufTagValue
		tag["protobuf_value"] = s.MapKeyValue[1].ProtobufTagValue
	}
	st.Add(jen.Id(s.GoName).Add(genGoType(s)).Tag(tag))

	if s.TrailingComments != "" {
		st = st.Line().Comment(string(s.TrailingComments))
	}

	return st
}

func genGoType(s spec.MessageFieldInfo) *jen.Statement {
	if s.IsMap {
		return jen.Map(genGoType(s.MapKeyValue[0])).
			Add(genGoType(s.MapKeyValue[1]))
	}

	goType := jen.Interface()

	switch s.Kind {
	case protoreflect.BoolKind:
		goType = jen.Bool()
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		goType = jen.Int32()
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		goType = jen.Uint32()
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		goType = jen.Int64()
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		goType = jen.Uint64()
	case protoreflect.FloatKind:
		goType = jen.Float32()
	case protoreflect.DoubleKind:
		goType = jen.Float64()
	case protoreflect.StringKind:
		goType = jen.String()
	case protoreflect.BytesKind:
		goType = jen.Index().Byte()
	case protoreflect.EnumKind:
		goType = jen.Qual(string(s.EnumGoIdent.GoImportPath), s.EnumGoIdent.GoName)
	case protoreflect.MessageKind, protoreflect.GroupKind:
		goType = jen.Qual(string(s.MessageGoIdent.GoImportPath), s.MessageGoIdent.GoName)
	}

	if s.Kind == protoreflect.MessageKind || s.Kind == protoreflect.GroupKind || s.HasPresence {
		goType = jen.Op("*").Add(goType)
	}

	if s.IsList {
		goType = jen.Index().Add(goType)
	}

	return goType
}

func genService(s spec.ServiceInfo) *jen.Statement {
	return jen.Type().Id(s.GoName).InterfaceFunc(func(g *jen.Group) {
		for _, method := range s.Methods {
			g.Id(method.GoName).Params(
				jen.Id("ctx").Qual(pkgContext, "Context"),
				jen.Id(strcase.ToLowerCamel(method.InputGoIdent.GoName)).Op("*").Qual(string(method.InputGoIdent.GoImportPath), method.InputGoIdent.GoName),
			).Params(
				jen.Op("*").Qual(string(method.OutputGoIdent.GoImportPath), method.OutputGoIdent.GoName),
				jen.Error(),
			)
		}
	}).Line()
}
