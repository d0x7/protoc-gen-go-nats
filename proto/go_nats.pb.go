// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.27.3
// source: go_nats.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NATSServiceOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version     string  `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Description *string `protobuf:"bytes,3,opt,name=description,proto3,oneof" json:"description,omitempty"`
}

func (x *NATSServiceOptions) Reset() {
	*x = NATSServiceOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_go_nats_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NATSServiceOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NATSServiceOptions) ProtoMessage() {}

func (x *NATSServiceOptions) ProtoReflect() protoreflect.Message {
	mi := &file_go_nats_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NATSServiceOptions.ProtoReflect.Descriptor instead.
func (*NATSServiceOptions) Descriptor() ([]byte, []int) {
	return file_go_nats_proto_rawDescGZIP(), []int{0}
}

func (x *NATSServiceOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *NATSServiceOptions) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *NATSServiceOptions) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

var file_go_nats_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*NATSServiceOptions)(nil),
		Field:         9001,
		Name:          "protoc.gen.go_nats.nats",
		Tag:           "bytes,9001,opt,name=nats",
		Filename:      "go_nats.proto",
	},
}

// Extension fields to descriptorpb.ServiceOptions.
var (
	// optional protoc.gen.go_nats.NATSServiceOptions nats = 9001;
	E_Nats = &file_go_nats_proto_extTypes[0]
)

var File_go_nats_proto protoreflect.FileDescriptor

var file_go_nats_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x67, 0x6f, 0x5f, 0x6e, 0x61, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x67, 0x6f, 0x5f, 0x6e,
	0x61, 0x74, 0x73, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x79, 0x0a, 0x12, 0x4e, 0x41, 0x54, 0x53, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01,
	0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x3a, 0x5f, 0x0a, 0x04, 0x6e, 0x61, 0x74, 0x73, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa9, 0x46, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x67, 0x6f,
	0x5f, 0x6e, 0x61, 0x74, 0x73, 0x2e, 0x4e, 0x41, 0x54, 0x53, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x04, 0x6e, 0x61, 0x74, 0x73, 0x88, 0x01,
	0x01, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_go_nats_proto_rawDescOnce sync.Once
	file_go_nats_proto_rawDescData = file_go_nats_proto_rawDesc
)

func file_go_nats_proto_rawDescGZIP() []byte {
	file_go_nats_proto_rawDescOnce.Do(func() {
		file_go_nats_proto_rawDescData = protoimpl.X.CompressGZIP(file_go_nats_proto_rawDescData)
	})
	return file_go_nats_proto_rawDescData
}

var file_go_nats_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_go_nats_proto_goTypes = []interface{}{
	(*NATSServiceOptions)(nil),          // 0: protoc.gen.go_nats.NATSServiceOptions
	(*descriptorpb.ServiceOptions)(nil), // 1: google.protobuf.ServiceOptions
}
var file_go_nats_proto_depIdxs = []int32{
	1, // 0: protoc.gen.go_nats.nats:extendee -> google.protobuf.ServiceOptions
	0, // 1: protoc.gen.go_nats.nats:type_name -> protoc.gen.go_nats.NATSServiceOptions
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	1, // [1:2] is the sub-list for extension type_name
	0, // [0:1] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_go_nats_proto_init() }
func file_go_nats_proto_init() {
	if File_go_nats_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_go_nats_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NATSServiceOptions); i {
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
	file_go_nats_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_go_nats_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_go_nats_proto_goTypes,
		DependencyIndexes: file_go_nats_proto_depIdxs,
		MessageInfos:      file_go_nats_proto_msgTypes,
		ExtensionInfos:    file_go_nats_proto_extTypes,
	}.Build()
	File_go_nats_proto = out.File
	file_go_nats_proto_rawDesc = nil
	file_go_nats_proto_goTypes = nil
	file_go_nats_proto_depIdxs = nil
}