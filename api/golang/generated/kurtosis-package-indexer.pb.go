// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.4
// source: kurtosis-package-indexer.proto

package generated

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

type IndexerPing struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *IndexerPing) Reset() {
	*x = IndexerPing{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kurtosis_package_indexer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexerPing) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexerPing) ProtoMessage() {}

func (x *IndexerPing) ProtoReflect() protoreflect.Message {
	mi := &file_kurtosis_package_indexer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexerPing.ProtoReflect.Descriptor instead.
func (*IndexerPing) Descriptor() ([]byte, []int) {
	return file_kurtosis_package_indexer_proto_rawDescGZIP(), []int{0}
}

type IndexerPong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *IndexerPong) Reset() {
	*x = IndexerPong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kurtosis_package_indexer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexerPong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexerPong) ProtoMessage() {}

func (x *IndexerPong) ProtoReflect() protoreflect.Message {
	mi := &file_kurtosis_package_indexer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexerPong.ProtoReflect.Descriptor instead.
func (*IndexerPong) Descriptor() ([]byte, []int) {
	return file_kurtosis_package_indexer_proto_rawDescGZIP(), []int{1}
}

var File_kurtosis_package_indexer_proto protoreflect.FileDescriptor

var file_kurtosis_package_indexer_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x6b, 0x75, 0x72, 0x74, 0x6f, 0x73, 0x69, 0x73, 0x2d, 0x70, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x2d, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x18, 0x6b, 0x75, 0x72, 0x74, 0x6f, 0x73, 0x69, 0x73, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x22, 0x0d, 0x0a, 0x0b, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x65, 0x72, 0x50, 0x69, 0x6e, 0x67, 0x22, 0x0d, 0x0a, 0x0b, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x65, 0x72, 0x50, 0x6f, 0x6e, 0x67, 0x32, 0x6e, 0x0a, 0x16, 0x4b, 0x75, 0x72, 0x74,
	0x6f, 0x73, 0x69, 0x73, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x65, 0x72, 0x12, 0x54, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x25, 0x2e, 0x6b, 0x75, 0x72,
	0x74, 0x6f, 0x73, 0x69, 0x73, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x50, 0x69, 0x6e,
	0x67, 0x1a, 0x25, 0x2e, 0x6b, 0x75, 0x72, 0x74, 0x6f, 0x73, 0x69, 0x73, 0x5f, 0x70, 0x61, 0x63,
	0x6b, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x65, 0x72, 0x50, 0x6f, 0x6e, 0x67, 0x42, 0x48, 0x5a, 0x46, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x75, 0x72, 0x74, 0x6f, 0x73, 0x69, 0x73, 0x2d,
	0x74, 0x65, 0x63, 0x68, 0x2f, 0x6b, 0x75, 0x72, 0x74, 0x6f, 0x73, 0x69, 0x73, 0x2d, 0x70, 0x61,
	0x63, 0x6b, 0x61, 0x67, 0x65, 0x2d, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x65, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kurtosis_package_indexer_proto_rawDescOnce sync.Once
	file_kurtosis_package_indexer_proto_rawDescData = file_kurtosis_package_indexer_proto_rawDesc
)

func file_kurtosis_package_indexer_proto_rawDescGZIP() []byte {
	file_kurtosis_package_indexer_proto_rawDescOnce.Do(func() {
		file_kurtosis_package_indexer_proto_rawDescData = protoimpl.X.CompressGZIP(file_kurtosis_package_indexer_proto_rawDescData)
	})
	return file_kurtosis_package_indexer_proto_rawDescData
}

var file_kurtosis_package_indexer_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_kurtosis_package_indexer_proto_goTypes = []interface{}{
	(*IndexerPing)(nil), // 0: kurtosis_package_indexer.IndexerPing
	(*IndexerPong)(nil), // 1: kurtosis_package_indexer.IndexerPong
}
var file_kurtosis_package_indexer_proto_depIdxs = []int32{
	0, // 0: kurtosis_package_indexer.KurtosisPackageIndexer.Ping:input_type -> kurtosis_package_indexer.IndexerPing
	1, // 1: kurtosis_package_indexer.KurtosisPackageIndexer.Ping:output_type -> kurtosis_package_indexer.IndexerPong
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_kurtosis_package_indexer_proto_init() }
func file_kurtosis_package_indexer_proto_init() {
	if File_kurtosis_package_indexer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kurtosis_package_indexer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexerPing); i {
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
		file_kurtosis_package_indexer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexerPong); i {
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
			RawDescriptor: file_kurtosis_package_indexer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_kurtosis_package_indexer_proto_goTypes,
		DependencyIndexes: file_kurtosis_package_indexer_proto_depIdxs,
		MessageInfos:      file_kurtosis_package_indexer_proto_msgTypes,
	}.Build()
	File_kurtosis_package_indexer_proto = out.File
	file_kurtosis_package_indexer_proto_rawDesc = nil
	file_kurtosis_package_indexer_proto_goTypes = nil
	file_kurtosis_package_indexer_proto_depIdxs = nil
}
