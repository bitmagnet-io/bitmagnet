package spec

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Spec struct {
	ProtocVersion string
	Files         []FileInfo
}

type FileInfo struct {
	ProtocVersion          string
	Generate               bool
	SyntaxCommentLocation  protoreflect.SourceLocation
	PackageCommentLocation protoreflect.SourceLocation

	Path string

	RawDescVarName          string
	GeneratedFilenamePrefix string
	GoImportPath            protogen.GoImportPath
	GoPackageName           protogen.GoPackageName

	Deprecated bool

	Enums          []EnumInfo
	Messages       []MessageInfo
	Services       []ServiceInfo
	PluginServices []ServiceInfo // For plugin interfaces
	HostService    *ServiceInfo  // For host functions
}

func New(plugin *protogen.Plugin) (Spec, error) {
	var spec Spec
	protocVersion := "(unknown)"
	if v := plugin.Request.GetCompilerVersion(); v != nil {
		protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
	}
	spec.ProtocVersion = protocVersion
	for _, file := range plugin.Files {
		f, err := newFile(file)
		if err != nil {
			return Spec{}, err
		}
		f.ProtocVersion = protocVersion
		spec.Files = append(spec.Files, f)
	}
	return spec, nil
}

const (
	FileDescriptorProto_Syntax_field_number  protoreflect.FieldNumber = 12
	FileDescriptorProto_Package_field_number protoreflect.FieldNumber = 2
)

func newFile(file *protogen.File) (FileInfo, error) {
	f := FileInfo{
		Generate:                file.Generate,
		SyntaxCommentLocation:   file.Desc.SourceLocations().ByPath(protoreflect.SourcePath{int32(FileDescriptorProto_Syntax_field_number)}),
		PackageCommentLocation:  file.Desc.SourceLocations().ByPath(protoreflect.SourcePath{int32(FileDescriptorProto_Package_field_number)}),
		RawDescVarName:          fileVarName(file, "rawDesc"),
		Path:                    file.Desc.Path(),
		GeneratedFilenamePrefix: file.GeneratedFilenamePrefix,
		GoImportPath:            file.GoImportPath,
		GoPackageName:           file.GoPackageName,
		Deprecated:              file.Proto.GetOptions().GetDeprecated(),
	}

	// Collect all enums, messages, and services in "flattened ordering".
	for _, enum := range file.Enums {
		f.Enums = append(f.Enums, newEnumInfo(enum))
	}

	for _, message := range file.Messages {
		f.Messages = append(f.Messages, newMessageInfo(message))
	}

	for _, service := range file.Services {
		param, err := parseParam(service.Comments.Leading.String())
		if err != nil {
			return FileInfo{}, err
		}

		si := newServiceInfo(service, param)
		switch param.Type {
		case ServicePlugin:
			f.PluginServices = append(f.PluginServices, si)
		case ServiceHost:
			if f.HostService != nil {
				return FileInfo{}, errors.New("type=host can be defined only once")
			}
			f.HostService = &si
		case ServiceUnknown:
			return FileInfo{}, errors.New("unknown go-plugin type")
		}
		f.Services = append(f.Services, si)
	}

	if len(file.Extensions) != 0 {
		return FileInfo{}, errors.New("extensions not supported")
	}

	walkMessages(file.Messages, func(m *protogen.Message) {
		for _, enum := range m.Enums {
			f.Enums = append(f.Enums, newEnumInfo(enum))
		}
		for _, message := range m.Messages {
			f.Messages = append(f.Messages, newMessageInfo(message))
		}
	})

	return f, nil
}

func walkMessages(messages []*protogen.Message, f func(*protogen.Message)) {
	for _, m := range messages {
		f(m)
		walkMessages(m.Messages, f)
	}
}

type EnumInfo struct {
	GoIdent         protogen.GoIdent
	Location        protogen.Location
	Deprecated      bool
	LeadingComments protogen.Comments
	Values          []EnumValueInfo
}

type EnumValueInfo struct {
	GoIdent          protogen.GoIdent
	Location         protogen.Location
	Deprecated       bool
	LeadingComments  protogen.Comments
	TrailingComments protogen.Comments
	Number           protoreflect.EnumNumber
	Name             protoreflect.Name
}

func newEnumInfo(enum *protogen.Enum) EnumInfo {
	return EnumInfo{
		GoIdent:         enum.GoIdent,
		Location:        enum.Location,
		Deprecated:      enum.Desc.Options().(*descriptorpb.EnumOptions).GetDeprecated(),
		LeadingComments: enum.Comments.Leading,
		Values: slice.Map(enum.Values, func(ev *protogen.EnumValue) EnumValueInfo {
			return EnumValueInfo{
				GoIdent:          ev.GoIdent,
				Location:         ev.Location,
				Deprecated:       ev.Desc.Options().(*descriptorpb.EnumValueOptions).GetDeprecated(),
				LeadingComments:  ev.Comments.Leading,
				TrailingComments: ev.Comments.Trailing,
				Number:           ev.Desc.Number(),
				Name:             ev.Desc.Name(),
			}
		}),
	}
}

type MessageInfo struct {
	IsMapEntry      bool
	GoIdent         protogen.GoIdent
	LeadingComments protogen.Comments
	Location        protogen.Location
	IsDeprecated    bool
	Unknown         protoreflect.RawFields

	Fields []MessageFieldInfo
	Oneofs []OneOfInfo
}

type MessageFieldInfo struct {
	Name                  protoreflect.Name
	GoName                string
	GoIdent               protogen.GoIdent
	MessageGoIdent        protogen.GoIdent
	EnumGoIdent           protogen.GoIdent
	EnumFirstValueGoIdent protogen.GoIdent
	Location              protogen.Location
	LeadingComments       protogen.Comments
	TrailingComments      protogen.Comments
	IsDeprecated          bool
	ProtobufTagValue      string
	IsList                bool
	IsMap                 bool
	HasPresence           bool
	MapKeyValue           *[2]MessageFieldInfo
	OneOf                 *OneOfInfo
	Number                protoreflect.FieldNumber
	Kind                  protoreflect.Kind
}

type OneOfInfo struct {
	Name             protoreflect.Name
	GoName           string
	GoIdent          protogen.GoIdent
	Location         protogen.Location
	IsSynthetic      bool
	LeadingComments  protogen.Comments
	TrailingComments protogen.Comments
	Fields           []MessageFieldInfo
}

func newMessageInfo(message *protogen.Message) MessageInfo {
	m := MessageInfo{
		IsMapEntry:      message.Desc.IsMapEntry(),
		GoIdent:         message.GoIdent,
		LeadingComments: message.Comments.Leading,
		Location:        message.Location,
		IsDeprecated:    message.Desc.Options().(*descriptorpb.MessageOptions).GetDeprecated(),
		Unknown:         message.Desc.Options().(*descriptorpb.MessageOptions).ProtoReflect().GetUnknown(),
	}
	for _, field := range message.Fields {
		m.Fields = append(m.Fields, newMessageFieldInfo(field))
	}
	return m
}

func newMessageFieldInfo(field *protogen.Field) MessageFieldInfo {
	info := MessageFieldInfo{
		Name:             field.Desc.Name(),
		GoName:           field.GoName,
		GoIdent:          field.GoIdent,
		Location:         field.Location,
		LeadingComments:  field.Comments.Leading,
		TrailingComments: field.Comments.Trailing,
		IsDeprecated:     field.Desc.Options().(*descriptorpb.FieldOptions).GetDeprecated(),
		IsList:           field.Desc.IsList(),
		IsMap:            field.Desc.IsMap(),
		HasPresence:      field.Desc.HasPresence(),
		Kind:             field.Desc.Kind(),
		Number:           field.Desc.Number(),
	}
	if info.IsMap {
		info.MapKeyValue = &[2]MessageFieldInfo{
			newMessageFieldInfo(field.Message.Fields[0]),
			newMessageFieldInfo(field.Message.Fields[1]),
		}
	}

	var enumName string
	if field.Desc.Kind() == protoreflect.EnumKind {
		enumName = protoimpl.X.LegacyEnumName(field.Enum.Desc)
		info.EnumGoIdent = field.Enum.GoIdent
		if len(field.Enum.Values) > 0 {
			info.EnumFirstValueGoIdent = field.Enum.Values[0].GoIdent
		}
	}
	info.ProtobufTagValue = marshalTag(field.Desc, enumName)

	// if o := field.Oneof; o != nil {
	// 	oneOf := newOneOfInfo(o)
	// 	info.OneOf = &oneOf
	// }

	if field.Message != nil {
		info.MessageGoIdent = field.Message.GoIdent
	}

	return info
}

func newOneOfInfo(oneof *protogen.Oneof) OneOfInfo {
	return OneOfInfo{
		Name:             oneof.Desc.Name(),
		GoName:           oneof.GoName,
		GoIdent:          oneof.GoIdent,
		Location:         oneof.Location,
		IsSynthetic:      oneof.Desc.IsSynthetic(),
		LeadingComments:  oneof.Comments.Leading,
		TrailingComments: oneof.Comments.Trailing,
		Fields: slice.Map(oneof.Fields, func(fld *protogen.Field) MessageFieldInfo {
			return newMessageFieldInfo(fld)
		}),
	}
}

type ServiceInfo struct {
	GoName          string
	Location        protogen.Location
	LeadingComments protogen.Comments
	Version         int
	Type            ServiceType
	Module          string
	Methods         []ServiceMethodInfo
}

type ServiceMethodInfo struct {
	GoName          string
	LeadingComments protogen.Comments
	InputGoIdent    protogen.GoIdent
	OutputGoIdent   protogen.GoIdent
}

func newServiceInfo(service *protogen.Service, param Parameter) ServiceInfo {
	return ServiceInfo{
		Type:            param.Type,
		Version:         param.APIVersion,
		Module:          param.Module,
		GoName:          service.GoName,
		Location:        service.Location,
		LeadingComments: service.Comments.Leading,
		Methods: slice.Map(service.Methods, func(m *protogen.Method) ServiceMethodInfo {
			return newServiceMethodInfo(m)
		}),
	}
}

func newServiceMethodInfo(method *protogen.Method) ServiceMethodInfo {
	return ServiceMethodInfo{
		GoName:          method.GoName,
		LeadingComments: method.Comments.Leading,
		InputGoIdent:    method.Input.GoIdent,
		OutputGoIdent:   method.Output.GoIdent,
	}
}

// marshalTag encodes the protoreflect.FieldDescriptor as a tag.
//
// The enumName must be provided if the kind is an enum.
// Historically, the formulation of the enum "name" was the proto package
// dot-concatenated with the generated Go identifier for the enum type.
// Depending on the context on how Marshal is called, there are different ways
// through which that information is determined. As such it is the caller's
// responsibility to provide a function to obtain that information.
func marshalTag(fd protoreflect.FieldDescriptor, enumName string) string {
	var tag []string
	switch fd.Kind() {
	case protoreflect.BoolKind, protoreflect.EnumKind, protoreflect.Int32Kind, protoreflect.Uint32Kind, protoreflect.Int64Kind, protoreflect.Uint64Kind:
		tag = append(tag, "varint")
	case protoreflect.Sint32Kind:
		tag = append(tag, "zigzag32")
	case protoreflect.Sint64Kind:
		tag = append(tag, "zigzag64")
	case protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind, protoreflect.FloatKind:
		tag = append(tag, "fixed32")
	case protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind, protoreflect.DoubleKind:
		tag = append(tag, "fixed64")
	case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.MessageKind:
		tag = append(tag, "bytes")
	case protoreflect.GroupKind:
		tag = append(tag, "group")
	}
	tag = append(tag, strconv.Itoa(int(fd.Number())))
	switch fd.Cardinality() {
	case protoreflect.Optional:
		tag = append(tag, "opt")
	case protoreflect.Required:
		tag = append(tag, "req")
	case protoreflect.Repeated:
		tag = append(tag, "rep")
	}
	if fd.IsPacked() {
		tag = append(tag, "packed")
	}
	name := string(fd.Name())
	if fd.Kind() == protoreflect.GroupKind {
		// The name of the FieldDescriptor for a group field is
		// lowercased. To find the original capitalization, we
		// look in the field's MessageType.
		name = string(fd.Message().Name())
	}
	tag = append(tag, "name="+name)
	if jsonName := fd.JSONName(); jsonName != "" && jsonName != name && !fd.IsExtension() {
		// NOTE: The jsonName != name condition is suspect, but it preserve
		// the exact same semantics from the previous generator.
		tag = append(tag, "json="+jsonName)
	}
	// The previous implementation does not tag extension fields as proto3,
	// even when the field is defined in a proto3 file. Match that behavior
	// for consistency.
	if fd.Syntax() == protoreflect.Proto3 && !fd.IsExtension() {
		tag = append(tag, "proto3")
	}
	if fd.Kind() == protoreflect.EnumKind && enumName != "" {
		tag = append(tag, "enum="+enumName)
	}
	if fd.ContainingOneof() != nil {
		tag = append(tag, "oneof")
	}
	// This must appear last in the tag, since commas in strings aren't escaped.
	// if fd.HasDefault() {
	// 	def, _ := marshalDefVal(fd.Default(), fd.DefaultEnumValue(), fd.Kind(), GoTag)
	// 	tag = append(tag, "def="+def)
	// }
	return strings.Join(tag, ",")
}
