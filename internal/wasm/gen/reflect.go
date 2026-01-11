// Copyright 2018 The Go Authors. All rights reserved.
// Copyright 2022 Teppei Fukuda. All rights reserved.

package gen

// func genMessageReflectMethods(g *protogen.GeneratedFile, m spec.MessageInfo) {
// 	// ProtoReflect method.
// 	// A dummy method is defined so that it implements proto.Message,
// 	// but it is not supposed to be called.
// 	g.P("func (x *", m.GoIdent, ") ProtoReflect() ", protoreflectPackage.Ident("Message"), " {")
// 	g.P("panic(`not implemented`)")
// 	g.P("}")
// 	g.P()
// }
