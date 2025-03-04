// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.3
// source: internal/protobuf/bitmagnet.proto

package protobuf

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Torrent_File_FileType int32

const (
	Torrent_File_unknown   Torrent_File_FileType = 0
	Torrent_File_archive   Torrent_File_FileType = 1
	Torrent_File_audio     Torrent_File_FileType = 2
	Torrent_File_data      Torrent_File_FileType = 3
	Torrent_File_document  Torrent_File_FileType = 4
	Torrent_File_image     Torrent_File_FileType = 5
	Torrent_File_software  Torrent_File_FileType = 6
	Torrent_File_subtitles Torrent_File_FileType = 7
	Torrent_File_video     Torrent_File_FileType = 8
)

// Enum value maps for Torrent_File_FileType.
var (
	Torrent_File_FileType_name = map[int32]string{
		0: "unknown",
		1: "archive",
		2: "audio",
		3: "data",
		4: "document",
		5: "image",
		6: "software",
		7: "subtitles",
		8: "video",
	}
	Torrent_File_FileType_value = map[string]int32{
		"unknown":   0,
		"archive":   1,
		"audio":     2,
		"data":      3,
		"document":  4,
		"image":     5,
		"software":  6,
		"subtitles": 7,
		"video":     8,
	}
)

func (x Torrent_File_FileType) Enum() *Torrent_File_FileType {
	p := new(Torrent_File_FileType)
	*p = x
	return p
}

func (x Torrent_File_FileType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Torrent_File_FileType) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_protobuf_bitmagnet_proto_enumTypes[0].Descriptor()
}

func (Torrent_File_FileType) Type() protoreflect.EnumType {
	return &file_internal_protobuf_bitmagnet_proto_enumTypes[0]
}

func (x Torrent_File_FileType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Torrent_File_FileType.Descriptor instead.
func (Torrent_File_FileType) EnumDescriptor() ([]byte, []int) {
	return file_internal_protobuf_bitmagnet_proto_rawDescGZIP(), []int{0, 0, 0}
}

type Classification_ContentType int32

const (
	Classification_unknown   Classification_ContentType = 0
	Classification_movie     Classification_ContentType = 1
	Classification_tv_show   Classification_ContentType = 2
	Classification_music     Classification_ContentType = 3
	Classification_ebook     Classification_ContentType = 4
	Classification_comic     Classification_ContentType = 5
	Classification_audiobook Classification_ContentType = 6
	Classification_game      Classification_ContentType = 7
	Classification_software  Classification_ContentType = 8
	Classification_xxx       Classification_ContentType = 9
)

// Enum value maps for Classification_ContentType.
var (
	Classification_ContentType_name = map[int32]string{
		0: "unknown",
		1: "movie",
		2: "tv_show",
		3: "music",
		4: "ebook",
		5: "comic",
		6: "audiobook",
		7: "game",
		8: "software",
		9: "xxx",
	}
	Classification_ContentType_value = map[string]int32{
		"unknown":   0,
		"movie":     1,
		"tv_show":   2,
		"music":     3,
		"ebook":     4,
		"comic":     5,
		"audiobook": 6,
		"game":      7,
		"software":  8,
		"xxx":       9,
	}
)

func (x Classification_ContentType) Enum() *Classification_ContentType {
	p := new(Classification_ContentType)
	*p = x
	return p
}

func (x Classification_ContentType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Classification_ContentType) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_protobuf_bitmagnet_proto_enumTypes[1].Descriptor()
}

func (Classification_ContentType) Type() protoreflect.EnumType {
	return &file_internal_protobuf_bitmagnet_proto_enumTypes[1]
}

func (x Classification_ContentType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Classification_ContentType.Descriptor instead.
func (Classification_ContentType) EnumDescriptor() ([]byte, []int) {
	return file_internal_protobuf_bitmagnet_proto_rawDescGZIP(), []int{1, 0}
}

type Torrent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InfoHash           string          `protobuf:"bytes,1,opt,name=infoHash,proto3" json:"infoHash,omitempty"`
	Name               string          `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	BaseName           string          `protobuf:"bytes,3,opt,name=baseName,proto3" json:"baseName,omitempty"`
	Size               int64           `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	Extension          *string         `protobuf:"bytes,5,opt,name=extension,proto3,oneof" json:"extension,omitempty"`
	Files              []*Torrent_File `protobuf:"bytes,6,rep,name=files,proto3" json:"files,omitempty"`
	FilesCount         *int32          `protobuf:"varint,7,opt,name=filesCount,proto3,oneof" json:"filesCount,omitempty"`
	FilesSize          *int64          `protobuf:"varint,8,opt,name=filesSize,proto3,oneof" json:"filesSize,omitempty"`
	FileExtensions     []string        `protobuf:"bytes,9,rep,name=fileExtensions,proto3" json:"fileExtensions,omitempty"`
	Seeders            *int32          `protobuf:"varint,10,opt,name=seeders,proto3,oneof" json:"seeders,omitempty"`
	Leechers           *int32          `protobuf:"varint,11,opt,name=leechers,proto3,oneof" json:"leechers,omitempty"`
	HasHint            bool            `protobuf:"varint,12,opt,name=hasHint,proto3" json:"hasHint,omitempty"`
	HasHintedContentId bool            `protobuf:"varint,13,opt,name=hasHintedContentId,proto3" json:"hasHintedContentId,omitempty"`
}

func (x *Torrent) Reset() {
	*x = Torrent{}
	mi := &file_internal_protobuf_bitmagnet_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Torrent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Torrent) ProtoMessage() {}

func (x *Torrent) ProtoReflect() protoreflect.Message {
	mi := &file_internal_protobuf_bitmagnet_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Torrent.ProtoReflect.Descriptor instead.
func (*Torrent) Descriptor() ([]byte, []int) {
	return file_internal_protobuf_bitmagnet_proto_rawDescGZIP(), []int{0}
}

func (x *Torrent) GetInfoHash() string {
	if x != nil {
		return x.InfoHash
	}
	return ""
}

func (x *Torrent) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Torrent) GetBaseName() string {
	if x != nil {
		return x.BaseName
	}
	return ""
}

func (x *Torrent) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Torrent) GetExtension() string {
	if x != nil && x.Extension != nil {
		return *x.Extension
	}
	return ""
}

func (x *Torrent) GetFiles() []*Torrent_File {
	if x != nil {
		return x.Files
	}
	return nil
}

func (x *Torrent) GetFilesCount() int32 {
	if x != nil && x.FilesCount != nil {
		return *x.FilesCount
	}
	return 0
}

func (x *Torrent) GetFilesSize() int64 {
	if x != nil && x.FilesSize != nil {
		return *x.FilesSize
	}
	return 0
}

func (x *Torrent) GetFileExtensions() []string {
	if x != nil {
		return x.FileExtensions
	}
	return nil
}

func (x *Torrent) GetSeeders() int32 {
	if x != nil && x.Seeders != nil {
		return *x.Seeders
	}
	return 0
}

func (x *Torrent) GetLeechers() int32 {
	if x != nil && x.Leechers != nil {
		return *x.Leechers
	}
	return 0
}

func (x *Torrent) GetHasHint() bool {
	if x != nil {
		return x.HasHint
	}
	return false
}

func (x *Torrent) GetHasHintedContentId() bool {
	if x != nil {
		return x.HasHintedContentId
	}
	return false
}

type Classification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContentType        Classification_ContentType `protobuf:"varint,1,opt,name=contentType,proto3,enum=bitmagnet.Classification_ContentType" json:"contentType,omitempty"`
	HasAttachedContent bool                       `protobuf:"varint,2,opt,name=hasAttachedContent,proto3" json:"hasAttachedContent,omitempty"`
	HasBaseTitle       bool                       `protobuf:"varint,3,opt,name=hasBaseTitle,proto3" json:"hasBaseTitle,omitempty"`
	Year               *int32                     `protobuf:"varint,4,opt,name=year,proto3,oneof" json:"year,omitempty"`
	Languages          []string                   `protobuf:"bytes,5,rep,name=languages,proto3" json:"languages,omitempty"`
	Episodes           []string                   `protobuf:"bytes,6,rep,name=episodes,proto3" json:"episodes,omitempty"`
	VideoResolution    *string                    `protobuf:"bytes,7,opt,name=videoResolution,proto3,oneof" json:"videoResolution,omitempty"`
	VideoSource        *string                    `protobuf:"bytes,8,opt,name=videoSource,proto3,oneof" json:"videoSource,omitempty"`
	VideoCodec         *string                    `protobuf:"bytes,9,opt,name=videoCodec,proto3,oneof" json:"videoCodec,omitempty"`
	Video3D            *string                    `protobuf:"bytes,10,opt,name=video3d,proto3,oneof" json:"video3d,omitempty"`
	ReleaseGroup       *string                    `protobuf:"bytes,11,opt,name=releaseGroup,proto3,oneof" json:"releaseGroup,omitempty"`
	ContentId          *string                    `protobuf:"bytes,12,opt,name=contentId,proto3,oneof" json:"contentId,omitempty"`
	ContentSource      *string                    `protobuf:"bytes,13,opt,name=contentSource,proto3,oneof" json:"contentSource,omitempty"`
}

func (x *Classification) Reset() {
	*x = Classification{}
	mi := &file_internal_protobuf_bitmagnet_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Classification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Classification) ProtoMessage() {}

func (x *Classification) ProtoReflect() protoreflect.Message {
	mi := &file_internal_protobuf_bitmagnet_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Classification.ProtoReflect.Descriptor instead.
func (*Classification) Descriptor() ([]byte, []int) {
	return file_internal_protobuf_bitmagnet_proto_rawDescGZIP(), []int{1}
}

func (x *Classification) GetContentType() Classification_ContentType {
	if x != nil {
		return x.ContentType
	}
	return Classification_unknown
}

func (x *Classification) GetHasAttachedContent() bool {
	if x != nil {
		return x.HasAttachedContent
	}
	return false
}

func (x *Classification) GetHasBaseTitle() bool {
	if x != nil {
		return x.HasBaseTitle
	}
	return false
}

func (x *Classification) GetYear() int32 {
	if x != nil && x.Year != nil {
		return *x.Year
	}
	return 0
}

func (x *Classification) GetLanguages() []string {
	if x != nil {
		return x.Languages
	}
	return nil
}

func (x *Classification) GetEpisodes() []string {
	if x != nil {
		return x.Episodes
	}
	return nil
}

func (x *Classification) GetVideoResolution() string {
	if x != nil && x.VideoResolution != nil {
		return *x.VideoResolution
	}
	return ""
}

func (x *Classification) GetVideoSource() string {
	if x != nil && x.VideoSource != nil {
		return *x.VideoSource
	}
	return ""
}

func (x *Classification) GetVideoCodec() string {
	if x != nil && x.VideoCodec != nil {
		return *x.VideoCodec
	}
	return ""
}

func (x *Classification) GetVideo3D() string {
	if x != nil && x.Video3D != nil {
		return *x.Video3D
	}
	return ""
}

func (x *Classification) GetReleaseGroup() string {
	if x != nil && x.ReleaseGroup != nil {
		return *x.ReleaseGroup
	}
	return ""
}

func (x *Classification) GetContentId() string {
	if x != nil && x.ContentId != nil {
		return *x.ContentId
	}
	return ""
}

func (x *Classification) GetContentSource() string {
	if x != nil && x.ContentSource != nil {
		return *x.ContentSource
	}
	return ""
}

type Torrent_File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index     int32                 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Path      string                `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	BasePath  string                `protobuf:"bytes,3,opt,name=basePath,proto3" json:"basePath,omitempty"`
	BaseName  string                `protobuf:"bytes,4,opt,name=baseName,proto3" json:"baseName,omitempty"`
	Size      int64                 `protobuf:"varint,5,opt,name=size,proto3" json:"size,omitempty"`
	Extension *string               `protobuf:"bytes,6,opt,name=extension,proto3,oneof" json:"extension,omitempty"`
	FileType  Torrent_File_FileType `protobuf:"varint,7,opt,name=fileType,proto3,enum=bitmagnet.Torrent_File_FileType" json:"fileType,omitempty"`
}

func (x *Torrent_File) Reset() {
	*x = Torrent_File{}
	mi := &file_internal_protobuf_bitmagnet_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Torrent_File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Torrent_File) ProtoMessage() {}

func (x *Torrent_File) ProtoReflect() protoreflect.Message {
	mi := &file_internal_protobuf_bitmagnet_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Torrent_File.ProtoReflect.Descriptor instead.
func (*Torrent_File) Descriptor() ([]byte, []int) {
	return file_internal_protobuf_bitmagnet_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Torrent_File) GetIndex() int32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *Torrent_File) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Torrent_File) GetBasePath() string {
	if x != nil {
		return x.BasePath
	}
	return ""
}

func (x *Torrent_File) GetBaseName() string {
	if x != nil {
		return x.BaseName
	}
	return ""
}

func (x *Torrent_File) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Torrent_File) GetExtension() string {
	if x != nil && x.Extension != nil {
		return *x.Extension
	}
	return ""
}

func (x *Torrent_File) GetFileType() Torrent_File_FileType {
	if x != nil {
		return x.FileType
	}
	return Torrent_File_unknown
}

var File_internal_protobuf_bitmagnet_proto protoreflect.FileDescriptor

var file_internal_protobuf_bitmagnet_proto_rawDesc = []byte{
	0x0a, 0x21, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x62, 0x69, 0x74, 0x6d, 0x61, 0x67, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x09, 0x62, 0x69, 0x74, 0x6d, 0x61, 0x67, 0x6e, 0x65, 0x74, 0x22, 0xe3,
	0x06, 0x0a, 0x07, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e,
	0x66, 0x6f, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6e,
	0x66, 0x6f, 0x48, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x61,
	0x73, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x61,
	0x73, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x21, 0x0a, 0x09, 0x65, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52,
	0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a,
	0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x62,
	0x69, 0x74, 0x6d, 0x61, 0x67, 0x6e, 0x65, 0x74, 0x2e, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x23, 0x0a, 0x0a,
	0x66, 0x69, 0x6c, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05,
	0x48, 0x01, 0x52, 0x0a, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x88, 0x01,
	0x01, 0x12, 0x21, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x03, 0x48, 0x02, 0x52, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x53, 0x69, 0x7a,
	0x65, 0x88, 0x01, 0x01, 0x12, 0x26, 0x0a, 0x0e, 0x66, 0x69, 0x6c, 0x65, 0x45, 0x78, 0x74, 0x65,
	0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x66, 0x69,
	0x6c, 0x65, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1d, 0x0a, 0x07,
	0x73, 0x65, 0x65, 0x64, 0x65, 0x72, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x48, 0x03, 0x52,
	0x07, 0x73, 0x65, 0x65, 0x64, 0x65, 0x72, 0x73, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x6c,
	0x65, 0x65, 0x63, 0x68, 0x65, 0x72, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x48, 0x04, 0x52,
	0x08, 0x6c, 0x65, 0x65, 0x63, 0x68, 0x65, 0x72, 0x73, 0x88, 0x01, 0x01, 0x12, 0x18, 0x0a, 0x07,
	0x68, 0x61, 0x73, 0x48, 0x69, 0x6e, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x68,
	0x61, 0x73, 0x48, 0x69, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x12, 0x68, 0x61, 0x73, 0x48, 0x69, 0x6e,
	0x74, 0x65, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x0d, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x12, 0x68, 0x61, 0x73, 0x48, 0x69, 0x6e, 0x74, 0x65, 0x64, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x1a, 0xe7, 0x02, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x61, 0x73,
	0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x61, 0x73,
	0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x61, 0x73, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x61, 0x73, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x21, 0x0a, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x09, 0x65, 0x78, 0x74, 0x65,
	0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x3c, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x62, 0x69, 0x74,
	0x6d, 0x61, 0x67, 0x6e, 0x65, 0x74, 0x2e, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x66, 0x69,
	0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x22, 0x7a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12,
	0x0b, 0x0a, 0x07, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05,
	0x61, 0x75, 0x64, 0x69, 0x6f, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x10,
	0x03, 0x12, 0x0c, 0x0a, 0x08, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x10, 0x04, 0x12,
	0x09, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x10, 0x05, 0x12, 0x0c, 0x0a, 0x08, 0x73, 0x6f,
	0x66, 0x74, 0x77, 0x61, 0x72, 0x65, 0x10, 0x06, 0x12, 0x0d, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x73, 0x10, 0x07, 0x12, 0x09, 0x0a, 0x05, 0x76, 0x69, 0x64, 0x65, 0x6f,
	0x10, 0x08, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e,
	0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x0d,
	0x0a, 0x0b, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x0c, 0x0a,
	0x0a, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x53, 0x69, 0x7a, 0x65, 0x42, 0x0a, 0x0a, 0x08, 0x5f,
	0x73, 0x65, 0x65, 0x64, 0x65, 0x72, 0x73, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6c, 0x65, 0x65, 0x63,
	0x68, 0x65, 0x72, 0x73, 0x22, 0x90, 0x06, 0x0a, 0x0e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x47, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e, 0x62,
	0x69, 0x74, 0x6d, 0x61, 0x67, 0x6e, 0x65, 0x74, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x2e, 0x0a, 0x12, 0x68, 0x61, 0x73, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x65, 0x64, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x12, 0x68, 0x61,
	0x73, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x65, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x12, 0x22, 0x0a, 0x0c, 0x68, 0x61, 0x73, 0x42, 0x61, 0x73, 0x65, 0x54, 0x69, 0x74, 0x6c, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x68, 0x61, 0x73, 0x42, 0x61, 0x73, 0x65, 0x54,
	0x69, 0x74, 0x6c, 0x65, 0x12, 0x17, 0x0a, 0x04, 0x79, 0x65, 0x61, 0x72, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x48, 0x00, 0x52, 0x04, 0x79, 0x65, 0x61, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a,
	0x09, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x09, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x65,
	0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x65,
	0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x2d, 0x0a, 0x0f, 0x76, 0x69, 0x64, 0x65, 0x6f,
	0x52, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x01, 0x52, 0x0f, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74,
	0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x25, 0x0a, 0x0b, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x0b, 0x76,
	0x69, 0x64, 0x65, 0x6f, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x88, 0x01, 0x01, 0x12, 0x23, 0x0a,
	0x0a, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x43, 0x6f, 0x64, 0x65, 0x63, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x03, 0x52, 0x0a, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x43, 0x6f, 0x64, 0x65, 0x63, 0x88,
	0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x33, 0x64, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x04, 0x52, 0x07, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x33, 0x64, 0x88, 0x01,
	0x01, 0x12, 0x27, 0x0a, 0x0c, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x48, 0x05, 0x52, 0x0c, 0x72, 0x65, 0x6c, 0x65, 0x61,
	0x73, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x48, 0x06, 0x52,
	0x09, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x29, 0x0a,
	0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x07, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x88, 0x01, 0x01, 0x22, 0x83, 0x01, 0x0a, 0x0b, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x75, 0x6e, 0x6b, 0x6e,
	0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x10, 0x01,
	0x12, 0x0b, 0x0a, 0x07, 0x74, 0x76, 0x5f, 0x73, 0x68, 0x6f, 0x77, 0x10, 0x02, 0x12, 0x09, 0x0a,
	0x05, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x65, 0x62, 0x6f, 0x6f,
	0x6b, 0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x63, 0x6f, 0x6d, 0x69, 0x63, 0x10, 0x05, 0x12, 0x0d,
	0x0a, 0x09, 0x61, 0x75, 0x64, 0x69, 0x6f, 0x62, 0x6f, 0x6f, 0x6b, 0x10, 0x06, 0x12, 0x08, 0x0a,
	0x04, 0x67, 0x61, 0x6d, 0x65, 0x10, 0x07, 0x12, 0x0c, 0x0a, 0x08, 0x73, 0x6f, 0x66, 0x74, 0x77,
	0x61, 0x72, 0x65, 0x10, 0x08, 0x12, 0x07, 0x0a, 0x03, 0x78, 0x78, 0x78, 0x10, 0x09, 0x42, 0x07,
	0x0a, 0x05, 0x5f, 0x79, 0x65, 0x61, 0x72, 0x42, 0x12, 0x0a, 0x10, 0x5f, 0x76, 0x69, 0x64, 0x65,
	0x6f, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0e, 0x0a, 0x0c, 0x5f,
	0x76, 0x69, 0x64, 0x65, 0x6f, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42, 0x0d, 0x0a, 0x0b, 0x5f,
	0x76, 0x69, 0x64, 0x65, 0x6f, 0x43, 0x6f, 0x64, 0x65, 0x63, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x76,
	0x69, 0x64, 0x65, 0x6f, 0x33, 0x64, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x72, 0x65, 0x6c, 0x65, 0x61,
	0x73, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42, 0x14, 0x5a, 0x12, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_protobuf_bitmagnet_proto_rawDescOnce sync.Once
	file_internal_protobuf_bitmagnet_proto_rawDescData = file_internal_protobuf_bitmagnet_proto_rawDesc
)

func file_internal_protobuf_bitmagnet_proto_rawDescGZIP() []byte {
	file_internal_protobuf_bitmagnet_proto_rawDescOnce.Do(func() {
		file_internal_protobuf_bitmagnet_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_protobuf_bitmagnet_proto_rawDescData)
	})
	return file_internal_protobuf_bitmagnet_proto_rawDescData
}

var file_internal_protobuf_bitmagnet_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_internal_protobuf_bitmagnet_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_protobuf_bitmagnet_proto_goTypes = []any{
	(Torrent_File_FileType)(0),      // 0: bitmagnet.Torrent.File.FileType
	(Classification_ContentType)(0), // 1: bitmagnet.Classification.ContentType
	(*Torrent)(nil),                 // 2: bitmagnet.Torrent
	(*Classification)(nil),          // 3: bitmagnet.Classification
	(*Torrent_File)(nil),            // 4: bitmagnet.Torrent.File
}
var file_internal_protobuf_bitmagnet_proto_depIdxs = []int32{
	4, // 0: bitmagnet.Torrent.files:type_name -> bitmagnet.Torrent.File
	1, // 1: bitmagnet.Classification.contentType:type_name -> bitmagnet.Classification.ContentType
	0, // 2: bitmagnet.Torrent.File.fileType:type_name -> bitmagnet.Torrent.File.FileType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_internal_protobuf_bitmagnet_proto_init() }
func file_internal_protobuf_bitmagnet_proto_init() {
	if File_internal_protobuf_bitmagnet_proto != nil {
		return
	}
	file_internal_protobuf_bitmagnet_proto_msgTypes[0].OneofWrappers = []any{}
	file_internal_protobuf_bitmagnet_proto_msgTypes[1].OneofWrappers = []any{}
	file_internal_protobuf_bitmagnet_proto_msgTypes[2].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_protobuf_bitmagnet_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_protobuf_bitmagnet_proto_goTypes,
		DependencyIndexes: file_internal_protobuf_bitmagnet_proto_depIdxs,
		EnumInfos:         file_internal_protobuf_bitmagnet_proto_enumTypes,
		MessageInfos:      file_internal_protobuf_bitmagnet_proto_msgTypes,
	}.Build()
	File_internal_protobuf_bitmagnet_proto = out.File
	file_internal_protobuf_bitmagnet_proto_rawDesc = nil
	file_internal_protobuf_bitmagnet_proto_goTypes = nil
	file_internal_protobuf_bitmagnet_proto_depIdxs = nil
}
