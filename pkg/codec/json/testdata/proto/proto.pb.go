// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.2
// source: pkg/codec/json/testdata/proto/proto.proto

package proto

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

type Ping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg   string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	Count int32  `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *Ping) Reset() {
	*x = Ping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_codec_json_testdata_proto_proto_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ping) ProtoMessage() {}

func (x *Ping) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_codec_json_testdata_proto_proto_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ping.ProtoReflect.Descriptor instead.
func (*Ping) Descriptor() ([]byte, []int) {
	return file_pkg_codec_json_testdata_proto_proto_proto_rawDescGZIP(), []int{0}
}

func (x *Ping) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *Ping) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type Pong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *Pong) Reset() {
	*x = Pong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_codec_json_testdata_proto_proto_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pong) ProtoMessage() {}

func (x *Pong) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_codec_json_testdata_proto_proto_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pong.ProtoReflect.Descriptor instead.
func (*Pong) Descriptor() ([]byte, []int) {
	return file_pkg_codec_json_testdata_proto_proto_proto_rawDescGZIP(), []int{1}
}

func (x *Pong) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_pkg_codec_json_testdata_proto_proto_proto protoreflect.FileDescriptor

var file_pkg_codec_json_testdata_proto_proto_proto_rawDesc = []byte{
	0x0a, 0x29, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x2f, 0x6a, 0x73, 0x6f, 0x6e,
	0x2f, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x2e, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x14, 0x0a, 0x05,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x22, 0x18, 0x0a, 0x04, 0x50, 0x6f, 0x6e, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x32, 0x5b, 0x0a, 0x04,
	0x45, 0x63, 0x68, 0x6f, 0x12, 0x27, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x69, 0x6e,
	0x67, 0x12, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x1a, 0x0b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6f, 0x6e, 0x67, 0x22, 0x00, 0x12, 0x2a, 0x0a,
	0x0a, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x0b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x1a, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x50, 0x6f, 0x6e, 0x67, 0x22, 0x00, 0x30, 0x01, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x6f, 0x78, 0x67, 0x6f, 0x2f, 0x62, 0x6f,
	0x78, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x2f, 0x6a, 0x73, 0x6f, 0x6e,
	0x2f, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_codec_json_testdata_proto_proto_proto_rawDescOnce sync.Once
	file_pkg_codec_json_testdata_proto_proto_proto_rawDescData = file_pkg_codec_json_testdata_proto_proto_proto_rawDesc
)

func file_pkg_codec_json_testdata_proto_proto_proto_rawDescGZIP() []byte {
	file_pkg_codec_json_testdata_proto_proto_proto_rawDescOnce.Do(func() {
		file_pkg_codec_json_testdata_proto_proto_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_codec_json_testdata_proto_proto_proto_rawDescData)
	})
	return file_pkg_codec_json_testdata_proto_proto_proto_rawDescData
}

var file_pkg_codec_json_testdata_proto_proto_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_codec_json_testdata_proto_proto_proto_goTypes = []interface{}{
	(*Ping)(nil), // 0: proto.Ping
	(*Pong)(nil), // 1: proto.Pong
}
var file_pkg_codec_json_testdata_proto_proto_proto_depIdxs = []int32{
	0, // 0: proto.Echo.StartPing:input_type -> proto.Ping
	0, // 1: proto.Echo.StreamPing:input_type -> proto.Ping
	1, // 2: proto.Echo.StartPing:output_type -> proto.Pong
	1, // 3: proto.Echo.StreamPing:output_type -> proto.Pong
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_codec_json_testdata_proto_proto_proto_init() }
func file_pkg_codec_json_testdata_proto_proto_proto_init() {
	if File_pkg_codec_json_testdata_proto_proto_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_codec_json_testdata_proto_proto_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ping); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_codec_json_testdata_proto_proto_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pong); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_codec_json_testdata_proto_proto_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_codec_json_testdata_proto_proto_proto_goTypes,
		DependencyIndexes: file_pkg_codec_json_testdata_proto_proto_proto_depIdxs,
		MessageInfos:      file_pkg_codec_json_testdata_proto_proto_proto_msgTypes,
	}.Build()
	File_pkg_codec_json_testdata_proto_proto_proto = out.File
	file_pkg_codec_json_testdata_proto_proto_proto_rawDesc = nil
	file_pkg_codec_json_testdata_proto_proto_proto_goTypes = nil
	file_pkg_codec_json_testdata_proto_proto_proto_depIdxs = nil
}
