package jen

import (
	"github.com/bitmagnet-io/bitmagnet/internal/wasm/gen/spec"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

func GenHostFile(s spec.FileInfo) *jen.File {
	if len(s.PluginServices) == 0 && s.HostService == nil {
		return nil
	}

	f := newFile(s)

	f.HeaderComment("//go:build !wasip1")

	genHeader(f, s)

	if s.HostService != nil {
		f.Add(genHostFunctions(s))
	}

	for _, service := range s.PluginServices {
		f.Add(genHostService(service))
	}

	return f
}

func genHostFunctions(s spec.FileInfo) *jen.Statement {
	structName := "_" + strcase.ToLowerCamel(s.HostService.GoName)
	st := jen.Const().Defs(
		jen.Id("i32").Op("=").Qual(pkgWazeroApi, "ValueTypeI32"),
		jen.Id("i64").Op("=").Qual(pkgWazeroApi, "ValueTypeI64"),
	).Line().
		Type().Id(structName).Struct(
		jen.Id(s.HostService.GoName),
	).Line()

	var funcSt jen.Statement

	for _, method := range s.HostService.Methods {
		funcSt.Line().Id("envBuilder").Dot("NewFunctionBuilder").Call().
			Dot("WithGoModuleFunction").ParamsFunc(func(g *jen.Group) {
			g.Line().Qual(pkgWazeroApi, "GoModuleFunc").Call(
				jen.Id("h").Dot("_" + method.GoName),
			)
			g.Line().Index().Qual(pkgWazeroApi, "ValueType").Values(
				jen.Id("i32"),
				jen.Id("i32"),
			)
			g.Line().Index().Qual(pkgWazeroApi, "ValueType").Values(
				jen.Id("i64"),
			)
			g.Line()
		}).
			Dot("").Line().Id("WithParameterNames").Call(
			jen.Lit("offset"),
			jen.Lit("size"),
		).
			Dot("").Line().Id("Export").Call(
			jen.Lit(strcase.ToSnake(method.GoName)),
		)
	}

	funcSt.Line().Line().Add(jen.Id("_").Op(",").Err().Op(":=").Id("envBuilder").Dot("Instantiate").Call(
		jen.Id("ctx"),
	).Line().
		Return(
			jen.Err(),
		))

	if len(s.PluginServices) == 0 {
		st.Func().Id("Instantiate").Params(
			jen.Id("ctx").Qual(pkgContext, "Context"),
			jen.Id("r").Qual(pkgWazero, "Runtime"),
			jen.Id("hostFunctions").Id(s.HostService.GoName),
		).Error().Block(
			jen.Id("envBuilder").Op(":=").Id("r").Dot("NewHostModuleBuilder").Call(
				jen.Lit(s.HostService.Module),
			).Line().
				Id("h").Op(":=").Id(structName).Values(
				jen.Id("hostFunctions"),
			).Line().Add(&funcSt),
		)
	} else {
		st.Func().Params(
			jen.Id("h").Id(structName),
		).Id("Instantiate").Params(
			jen.Id("ctx").Qual(pkgContext, "Context"),
			jen.Id("r").Qual(pkgWazero, "Runtime"),
			jen.Id("hostFunctions").Id(s.HostService.GoName),
		).Error().Block(
			jen.Id("envBuilder").Op(":=").Id("r").Dot("NewHostModuleBuilder").Call(
				jen.Lit(s.HostService.Module),
			).Line().Add(&funcSt),
		)
	}

	st.Line()

	for _, method := range s.HostService.Methods {
		st.Line().Func().Params(
			jen.Id("h").Id(structName),
		).Id("_"+method.GoName).Params(
			jen.Id("ctx").Qual(pkgContext, "Context"),
			jen.Id("m").Qual(pkgWazeroApi, "Module"),
			jen.Id("stack").Index().Uint64(),
		).Block(
			jen.List(jen.Id("offset"), jen.Id("size")).Op(":=").List(
				jen.Uint32().Call(jen.Id("stack").Index(jen.Lit(0))),
				jen.Uint32().Call(jen.Id("stack").Index(jen.Lit(1))),
			),
			jen.Id("buf").Op(",").Err().Op(":=").Qual(pkgWasm, "ReadMemory").Call(
				jen.Id("m").Dot("Memory").Call(),
				jen.Id("offset"),
				jen.Id("size"),
			),
			jen.If(
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Panic(jen.Id("err")),
			),
			jen.Id("request").Op(":=").New(
				jen.Qual(string(method.InputGoIdent.GoImportPath), method.InputGoIdent.GoName),
			),
			jen.Err().Op("=").Id("request").Dot("UnmarshalVT").Call(
				jen.Id("buf"),
			),
			jen.If(
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Panic(jen.Id("err")),
			),
			jen.List(jen.Id("resp"), jen.Id("err")).Op(":=").Id("h").Dot(method.GoName).Call(
				jen.Id("ctx"),
				jen.Id("request"),
			),
			jen.If(
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Panic(jen.Id("err")),
			),
			jen.Id("buf").Op(",").Err().Op("=").Id("resp").Dot("MarshalVT").Call(),
			jen.If(
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Panic(jen.Id("err")),
			),
			jen.Id("ptr").Op(",").Err().Op(":=").Qual(pkgWasm, "WriteMemory").Call(
				jen.Id("ctx"),
				jen.Id("m"),
				jen.Id("buf"),
			),
			jen.If(
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Panic(jen.Id("err")),
			),
			jen.Id("ptrLen").Op(":=").Parens(
				jen.Id("ptr").Op("<<").Uint64().Call(jen.Lit(32)).Op("|").Uint64().Call(
					jen.Len(jen.Id("buf")),
				),
			),
			jen.Id("stack").Index(jen.Lit(0)).Op("=").Id("ptrLen"),
		)
	}

	return st
}

func genHostService(s spec.ServiceInfo) *jen.Statement {
	pluginName := s.GoName + "Plugin"

	pluginStructName := strcase.ToLowerCamel(s.GoName) + "Plugin"

	loadModuleParams := []jen.Code{
		jen.Id("ctx").Qual(pkgContext, "Context"),
		jen.Id("module").Qual(pkgWazeroApi, "Module"),
	}

	if s.Type == spec.ServiceHost {
		loadModuleParams = append(loadModuleParams,
			jen.Id("hostFunctions").Id(s.GoName),
		)
	}

	st := jen.Const().Id(pluginName+"APIVersion").Op("=").Lit(s.Version).Line().
		Type().Id(pluginName).Struct().Line().
		Func().Params(jen.Id("p").Op("*").Id(pluginName)).Id("LoadModule").Params(
		loadModuleParams...,
	).Params(
		jen.Id(s.GoName),
		jen.Error(),
	).BlockFunc(func(g *jen.Group) {
		// Compare API versions with the loading plugin
		g.Id("apiVersion").Op(":=").Id("module").Dot("ExportedFunction").Call(
			jen.Lit(strcase.ToSnake(s.GoName) + "_api_version"),
		).Line().
			If(jen.Id("apiVersion").Op("==").Nil()).Block(
			jen.Return(
				jen.Nil(),
				jen.Qual("errors", "New").Call(
					jen.Lit(strcase.ToSnake(s.GoName)+"_api_version is not exported"),
				),
			),
		).Line().
			Id("results").Op(",").Err().Op(":=").Id("apiVersion").Dot("Call").Call(
			jen.Id("ctx"),
		).Line().
			If(
				jen.Err().Op("!=").Nil(),
			).Block(
			jen.Return(
				jen.Nil(),
				jen.Err(),
			),
		).Else().If(
			jen.Len(jen.Id("results")).Op("!=").Lit(1),
		).Block(
			jen.Return(
				jen.Nil(),
				jen.Qual("errors", "New").Call(
					jen.Lit("invalid "+strcase.ToSnake(s.GoName)+"_api_version signature"),
				),
			),
		).Line().
			If(
				jen.Id("results").Index(jen.Lit(0)).Op("!=").Id(pluginName + "APIVersion"),
			).Block(
			jen.Return(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("API version mismatch, host: %d, plugin: %d"),
					jen.Id(pluginName+"APIVersion"),
					jen.Id("results").Index(jen.Lit(0)),
				),
			),
		)

		for _, method := range s.Methods {
			varName := strcase.ToLowerCamel(method.GoName)
			funcName := strcase.ToSnake(s.GoName + method.GoName)
			g.Id(varName).Op(":=").Id("module").Dot("ExportedFunction").Call(
				jen.Lit(funcName),
			)
			g.If(
				jen.Id(varName).Op("==").Nil(),
			).Block(
				jen.Return(
					jen.Nil(),
					jen.Qual("errors", "New").Call(
						jen.Lit(funcName+" is not exported"),
					),
				),
			)
		}

		g.Id("malloc").Op(":=").Id("module").Dot("ExportedFunction").Call(
			jen.Lit("malloc"),
		).Line().
			If(
				jen.Id("malloc").Op("==").Nil(),
			).Block(
			jen.Return(
				jen.Nil(),
				jen.Qual("errors", "New").Call(
					jen.Lit("malloc is not exported"),
				),
			),
		).Line().
			Id("free").Op(":=").Id("module").Dot("ExportedFunction").Call(
			jen.Lit("free"),
		).Line().
			If(
				jen.Id("free").Op("==").Nil(),
			).Block(
			jen.Return(
				jen.Nil(),
				jen.Qual("errors", "New").Call(
					jen.Lit("free is not exported"),
				),
			),
		).Line().
			Return(
				jen.Op("&").Id(pluginStructName).ValuesFunc(func(g *jen.Group) {
					// g.Line().Id("sem").Op(":").Make(jen.Chan().Struct(), jen.Lit(1))
					g.Line().Id("module").Op(":").Id("module")
					g.Line().Id("malloc").Op(":").Id("malloc")
					g.Line().Id("free").Op(":").Id("free")
					for _, method := range s.Methods {
						varName := strcase.ToLowerCamel(method.GoName)
						g.Line().Id(varName).Op(":").Id(varName)
					}
					g.Line()
				}),
				jen.Nil(),
			)
	})

	st.Line().Type().Id(pluginStructName).StructFunc(func(g *jen.Group) {
		// g.Id("sem").Chan().Struct()
		g.Id("module").Qual(pkgWazeroApi, "Module")
		g.Id("malloc").Qual(pkgWazeroApi, "Function")
		g.Id("free").Qual(pkgWazeroApi, "Function")
		for _, method := range s.Methods {
			varName := strcase.ToLowerCamel(method.GoName)
			g.Id(varName).Qual(pkgWazeroApi, "Function")
		}
	})
	for _, method := range s.Methods {
		st.Line().Add(genHostServiceMethod(pluginStructName, method))
	}

	return st
}

func genHostServiceMethod(structName string, m spec.ServiceMethodInfo) *jen.Statement {
	// outputType := jen.Id(m.OutputGoIdent.GoName)
	// if m.OutputGoIdent.GoImportPath != "" {
	// 	outputType = jen.Qual(string(m.OutputGoIdent.GoImportPath), m.OutputGoIdent.GoName)
	// }
	return jen.Func().Params(jen.Id("p").Op("*").Id(structName)).Id(m.GoName).Params(
		jen.Id("ctx").Qual(pkgContext, "Context"),
		jen.Id("request").Op("*").Qual(string(m.InputGoIdent.GoImportPath), m.InputGoIdent.GoName),
	).Params(
		jen.Op("*").Qual(string(m.OutputGoIdent.GoImportPath), m.OutputGoIdent.GoName),
		jen.Error(),
	).Block(
		// jen.Select().Block(
		// 	jen.Case(jen.Id("<-ctx").Dot("Done").Call()).Block(
		// 		jen.Return(
		// 			jen.Nil(),
		// 			jen.Id("ctx").Dot("Err").Call(),
		// 		),
		// 	),
		// 	jen.Case(jen.Id("p").Dot("sem").Op("<-").Struct().Values()).Block(),
		// ).Line(),
		// jen.Defer().Func().Params().Block(
		// 	jen.Op("<-").Id("p").Dot("sem"),
		// ).Call().Line(),
		jen.Id("data").Op(",").Err().Op(":=").Id("request").Dot("MarshalVT").Call(),
		jen.If(
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(
				jen.Nil(),
				jen.Err(),
			),
		),
		jen.Var().Id("ptr").Uint64(),
		jen.Id("dataSize").Op(":=").Len(jen.Id("data")),
		jen.If(
			jen.Id("dataSize").Op("!=").Lit(0),
		).Block(
			jen.Id("results").Op(",").Err().Op(":=").Id("p").Dot("malloc").Dot("Call").Call(
				jen.Id("ctx"),
				jen.Uint64().Call(jen.Id("dataSize")),
			),
			jen.If(
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(
					jen.Nil(),
					jen.Err(),
				),
			),
			jen.Id("ptr").Op("=").Uint64().Call(jen.Id("results").Index(jen.Lit(0))),
			jen.Defer().Id("p").Dot("free").Dot("Call").Call(
				jen.Id("ctx"),
				jen.Id("ptr"),
			),
			jen.If(jen.Op("!").Id("p").Dot("module").Dot("Memory").Call().Dot("Write").Call(
				jen.Uint32().Call(jen.Id("ptr")),
				jen.Id("data"),
			).Block(
				jen.Return(
					jen.Nil(),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("out of range memory size"),
					),
				),
			)),
		),
		jen.List(jen.Id("ptrSize"), jen.Err()).Op(":=").Id("p").Dot(strcase.ToLowerCamel(m.GoName)).Dot("Call").Call(
			jen.Id("ctx"),
			jen.Id("ptr"),
			jen.Uint64().Call(jen.Id("dataSize")),
		),
		jen.If(
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(
				jen.Nil(),
				jen.Err(),
			),
		),
		jen.Id("resPtr").Op(":=").Uint32().Call(jen.Id("ptrSize").Index(jen.Lit(0)).Op(">>").Lit(32)),
		jen.Id("resSize").Op(":=").Uint32().Call(jen.Id("ptrSize").Index(jen.Lit(0))),
		jen.Var().Id("isErrResponse").Bool(),
		jen.If(
			jen.Id("resSize").Op("&").Parens(jen.Lit(1).Op("<<").Lit(31)).Op(">").Lit(0),
		).Block(
			jen.Id("isErrResponse").Op("=").Lit(true),
			jen.Id("resSize").Op("&^=").Parens(jen.Lit(1).Op("<<").Lit(31)),
		),
		jen.If(
			jen.Id("resPtr").Op("!=").Lit(0),
		).Block(
			jen.Defer().Id("p").Dot("free").Dot("Call").Call(
				jen.Id("ctx"),
				jen.Uint64().Call(jen.Id("resPtr")),
			),
		),
		jen.List(jen.Id("bytes"), jen.Id("ok")).Op(":=").Id("p").Dot("module").Dot("Memory").Call().Dot("Read").Call(
			jen.Id("resPtr"),
			jen.Id("resSize"),
		),
		jen.If(
			jen.Op("!").Id("ok"),
		).Block(
			jen.Return(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("out of range memory size"),
				),
			),
		),
		jen.If(
			jen.Id("isErrResponse"),
		).Block(
			jen.Return(
				jen.Nil(),
				jen.Qual("errors", "New").Call(
					jen.String().Call(jen.Id("bytes")),
				),
			),
		),
		jen.Id("response").Op(":=").New(
			jen.Qual(string(m.OutputGoIdent.GoImportPath), m.OutputGoIdent.GoName),
		),
		jen.If(
			jen.Err().Op("=").Id("response").Dot("UnmarshalVT").Call(
				jen.Id("bytes"),
			),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(
				jen.Nil(),
				jen.Err(),
			),
		),
		jen.Return(
			jen.Id("response"),
			jen.Nil(),
		),
	).Line()
}
