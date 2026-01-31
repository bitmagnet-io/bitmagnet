package jen

import (
	"github.com/bitmagnet-io/bitmagnet/internal/wasm/gen/spec"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

func GenPluginFile(s spec.FileInfo) *jen.File {
	f := newFile(s)

	f.HeaderComment("//go:build wasip1")

	genHeader(f, s)

	for _, svc := range s.PluginServices {
		f.Add(genPluginService(svc))
	}

	if s.HostService != nil {
		f.Anon("unsafe")
		f.Add(genPluginHostFunctions(s))
	}

	return f
}

func genPluginService(s spec.ServiceInfo) *jen.Statement {
	serviceVar := strcase.ToLowerCamel(s.GoName) + "Service"

	st := jen.Const().Id(s.GoName + "PluginAPIVersion").Op("=").Lit(s.Version).Line().
		Comment("//go:wasmexport " + strcase.ToSnake(s.GoName) + "_api_version").Line().
		Func().Id("_" + strcase.ToSnake(s.GoName) + "_api_version").Params().Uint64().Block(
		jen.Return().Id(s.GoName + "PluginAPIVersion"),
	).Line().
		Var().Id(serviceVar).Id(s.GoName).Line().
		Func().Id("Register" + s.GoName).Params(jen.Id("p").Id(s.GoName)).Block(
		jen.Id(serviceVar).Op("=").Id("p"),
	).Line()

	for _, method := range s.Methods {
		exportedName := strcase.ToSnake(s.GoName + method.GoName)
		st = st.Line().Comment("//go:wasmexport "+exportedName).Line().
			Func().
			Id("_"+exportedName).Params(jen.Id("ptr").Uint32(), jen.Id("size").Uint32()).Uint64().Block(
			jen.Id("b").Op(":=").Qual(pkgWasm, "PtrToByte").Call(jen.Id("ptr"), jen.Id("size")),
			jen.Id("req").
				Op(":=").
				New(jen.Qual(string(method.InputGoIdent.GoImportPath), method.InputGoIdent.GoName)),
			jen.Err().Op(":=").Id("req").Dot("UnmarshalVT").Call(jen.Id("b")),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Panic(jen.Id("err")),
			),
			jen.Id("resp").Op(",").Err().Op(":=").Id(serviceVar).Dot(method.GoName).Call(
				jen.Qual(pkgContext, "Background").Call(),
				jen.Id("req"),
			),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.List(jen.Id("ptr"), jen.Id("size")).
					Op("=").
					Qual(pkgWasm, "ByteToPtr").
					Call(jen.Index().Byte().Call(jen.Id("err").Dot("Error").Call())),
				jen.Return().
					Parens(jen.Id("uint64").Call(jen.Id("ptr")).Op("<<").Uint64().Call(jen.Lit(32))).
					Op("|").
					Uint64().
					Call(jen.Id("size")).
					Op("|").
					Parens(jen.Lit(1).Op("<<").Lit(31)),
			),
			jen.Id("b").Op(",").Err().Op("=").Id("resp").Dot("MarshalVT").Call(),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Panic(jen.Id("err")),
			),
			jen.List(jen.Id("ptr"), jen.Id("size")).Op("=").Qual(pkgWasm, "ByteToPtr").Call(jen.Id("b")),
			jen.Return().
				Parens(jen.Id("uint64").Call(jen.Id("ptr")).Op("<<").Uint64().Call(jen.Lit(32))).
				Op("|").
				Uint64().
				Call(jen.Id("size")),
		)
	}

	return st
}

func genPluginHostFunctions(s spec.FileInfo) *jen.Statement {
	structName := strcase.ToLowerCamel(s.HostService.GoName)

	st := jen.Type().Id(structName).Struct().Line().
		Func().Id("New" + s.HostService.GoName).Params().Id(s.HostService.GoName).Block(
		jen.Return().Id(structName).Values(),
	).Line()

	for _, method := range s.HostService.Methods {
		importedName := strcase.ToSnake(method.GoName)
		st.Comment("//go:wasmimport "+s.HostService.Module+" "+importedName).Line().
			Func().Id("_"+importedName).Params(
			jen.Id("ptr").Uint32(),
			jen.Id("size").Uint32(),
		).Uint64().Line().Line()

		st.Func().Params(jen.Id("h").Id(structName)).Id(method.GoName).Params(
			jen.Id("ctx").Qual(pkgContext, "Context"),
			jen.Id("request").
				Op("*").
				Qual(string(method.InputGoIdent.GoImportPath), method.InputGoIdent.GoName),
		).
			Op("(*").
			Qual(string(method.OutputGoIdent.GoImportPath), method.OutputGoIdent.GoName).
			Op(", error)").Block(
			jen.Id("buf").Op(",").Err().Op(":=").Id("request").Dot("MarshalVT").Call(),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return().List(jen.Nil(), jen.Id("err")),
			),
			jen.Id("ptr").Op(",").Id("size").Op(":=").Qual(pkgWasm, "ByteToPtr").Call(jen.Id("buf")),
			jen.Id("ptrSize").Op(":=").Id("_"+importedName).Call(jen.Id("ptr"), jen.Id("size")),
			jen.Qual(pkgWasm, "Free").Call(jen.Id("ptr")),
			jen.Id("ptr").Op("=").Uint32().Parens(jen.Id("ptrSize").Op(">>").Lit(32)),
			jen.Id("size").Op("=").Uint32().Call(jen.Id("ptrSize")),
			jen.Id("buf").Op("=").Qual(pkgWasm, "PtrToByte").Call(jen.Id("ptr"), jen.Id("size")),
			jen.Id("response").
				Op(":=").
				New(jen.Qual(string(method.OutputGoIdent.GoImportPath), method.OutputGoIdent.GoName)),
			jen.Err().Op("=").Id("response").Dot("UnmarshalVT").Call(jen.Id("buf")),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return().List(jen.Nil(), jen.Id("err")),
			),
			jen.Return().List(jen.Id("response"), jen.Nil()),
		).
			Line()
	}

	return st
}
