// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.4
// source: kurtosis_package_indexer.proto

package generated

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PackageArgType int32

const (
	PackageArgType_BOOL    PackageArgType = 0
	PackageArgType_STRING  PackageArgType = 1
	PackageArgType_INTEGER PackageArgType = 2
	PackageArgType_FLOAT   PackageArgType = 3
)

// Enum value maps for PackageArgType.
var (
	PackageArgType_name = map[int32]string{
		0: "BOOL",
		1: "STRING",
		2: "INTEGER",
		3: "FLOAT",
	}
	PackageArgType_value = map[string]int32{
		"BOOL":    0,
		"STRING":  1,
		"INTEGER": 2,
		"FLOAT":   3,
	}
)

func (x PackageArgType) Enum() *PackageArgType {
	p := new(PackageArgType)
	*p = x
	return p
}

func (x PackageArgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PackageArgType) Descriptor() protoreflect.EnumDescriptor {
	return file_kurtosis_package_indexer_proto_enumTypes[0].Descriptor()
}

func (PackageArgType) Type() protoreflect.EnumType {
	return &file_kurtosis_package_indexer_proto_enumTypes[0]
}

func (x PackageArgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PackageArgType.Descriptor instead.
func (PackageArgType) EnumDescriptor() ([]byte, []int) {
	return file_kurtosis_package_indexer_proto_rawDescGZIP(), []int{0}
}

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

type GetPackagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Packages []*KurtosisPackage `protobuf:"bytes,1,rep,name=packages,proto3" json:"packages,omitempty"`
}

func (x *GetPackagesResponse) Reset() {
	*x = GetPackagesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kurtosis_package_indexer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPackagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPackagesResponse) ProtoMessage() {}

func (x *GetPackagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kurtosis_package_indexer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPackagesResponse.ProtoReflect.Descriptor instead.
func (*GetPackagesResponse) Descriptor() ([]byte, []int) {
	return file_kurtosis_package_indexer_proto_rawDescGZIP(), []int{2}
}

func (x *GetPackagesResponse) GetPackages() []*KurtosisPackage {
	if x != nil {
		return x.Packages
	}
	return nil
}

type KurtosisPackage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string        `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Args []*PackageArg `protobuf:"bytes,2,rep,name=args,proto3" json:"args,omitempty"`
}

func (x *KurtosisPackage) Reset() {
	*x = KurtosisPackage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kurtosis_package_indexer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KurtosisPackage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KurtosisPackage) ProtoMessage() {}

func (x *KurtosisPackage) ProtoReflect() protoreflect.Message {
	mi := &file_kurtosis_package_indexer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KurtosisPackage.ProtoReflect.Descriptor instead.
func (*KurtosisPackage) Descriptor() ([]byte, []int) {
	return file_kurtosis_package_indexer_proto_rawDescGZIP(), []int{3}
}

func (x *KurtosisPackage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *KurtosisPackage) GetArgs() []*PackageArg {
	if x != nil {
		return x.Args
	}
	return nil
}

type PackageArg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string          `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	IsRequired bool            `protobuf:"varint,2,opt,name=is_required,json=isRequired,proto3" json:"is_required,omitempty"`
	Type       *PackageArgType `protobuf:"varint,3,opt,name=type,proto3,enum=kurtosis_package_indexer.PackageArgType,oneof" json:"type,omitempty"`
}

func (x *PackageArg) Reset() {
	*x = PackageArg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kurtosis_package_indexer_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PackageArg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PackageArg) ProtoMessage() {}

func (x *PackageArg) ProtoReflect() protoreflect.Message {
	mi := &file_kurtosis_package_indexer_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PackageArg.ProtoReflect.Descriptor instead.
func (*PackageArg) Descriptor() ([]byte, []int) {
	return file_kurtosis_package_indexer_proto_rawDescGZIP(), []int{4}
}

func (x *PackageArg) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PackageArg) GetIsRequired() bool {
	if x != nil {
		return x.IsRequired
	}
	return false
}

func (x *PackageArg) GetType() PackageArgType {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return PackageArgType_BOOL
}

var File_kurtosis_package_indexer_proto protoreflect.FileDescriptor

var file_kurtosis_package_indexer_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x6b, 0x75, 0x72, 0x74, 0x6f, 0x73, 0x69, 0x73, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x18, 0x6b, 0x75, 0x72, 0x74, 0x6f, 0x73, 0x69, 0x73, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0d, 0x0a, 0x0b, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x65, 0x72, 0x50, 0x69, 0x6e, 0x67, 0x22, 0x0d, 0x0a, 0x0b, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65,
	0x72, 0x50, 0x6f, 0x6e, 0x67, 0x22, 0x5c, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x50, 0x61, 0x63, 0x6b,
	0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a, 0x08,
	0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29,
	0x2e, 0x6b, 0x75, 0x72, 0x74, 0x6f, 0x73, 0x69, 0x73, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67,
	0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2e, 0x4b, 0x75, 0x72, 0x74, 0x6f, 0x73,
	0x69, 0x73, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x52, 0x08, 0x70, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x73, 0x22, 0x5f, 0x0a, 0x0f, 0x4b, 0x75, 0x72, 0x74, 0x6f, 0x73, 0x69, 0x73, 0x50,
	0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x38, 0x0a, 0x04, 0x61, 0x72,
	0x67, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6b, 0x75, 0x72, 0x74, 0x6f,
	0x73, 0x69, 0x73, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x65, 0x72, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x41, 0x72, 0x67, 0x52, 0x04,
	0x61, 0x72, 0x67, 0x73, 0x22, 0x8d, 0x01, 0x0a, 0x0a, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65,
	0x41, 0x72, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x72, 0x65,
	0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x12, 0x41, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x28, 0x2e, 0x6b, 0x75, 0x72, 0x74, 0x6f, 0x73, 0x69,
	0x73, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65,
	0x72, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x41, 0x72, 0x67, 0x54, 0x79, 0x70, 0x65,
	0x48, 0x00, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x2a, 0x3e, 0x0a, 0x0e, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x41,
	0x72, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x42, 0x4f, 0x4f, 0x4c, 0x10, 0x00,
	0x12, 0x0a, 0x0a, 0x06, 0x53, 0x54, 0x52, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07,
	0x49, 0x4e, 0x54, 0x45, 0x47, 0x45, 0x52, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x46, 0x4c, 0x4f,
	0x41, 0x54, 0x10, 0x03, 0x32, 0xc6, 0x01, 0x0a, 0x16, 0x4b, 0x75, 0x72, 0x74, 0x6f, 0x73, 0x69,
	0x73, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x12,
	0x54, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x25, 0x2e, 0x6b, 0x75, 0x72, 0x74, 0x6f, 0x73,
	0x69, 0x73, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x65, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x50, 0x69, 0x6e, 0x67, 0x1a, 0x25,
	0x2e, 0x6b, 0x75, 0x72, 0x74, 0x6f, 0x73, 0x69, 0x73, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67,
	0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65,
	0x72, 0x50, 0x6f, 0x6e, 0x67, 0x12, 0x56, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x50, 0x61, 0x63, 0x6b,
	0x61, 0x67, 0x65, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2d, 0x2e, 0x6b,
	0x75, 0x72, 0x74, 0x6f, 0x73, 0x69, 0x73, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x48, 0x5a,
	0x46, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x75, 0x72, 0x74,
	0x6f, 0x73, 0x69, 0x73, 0x2d, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x6b, 0x75, 0x72, 0x74, 0x6f, 0x73,
	0x69, 0x73, 0x2d, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x2d, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x67, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_kurtosis_package_indexer_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_kurtosis_package_indexer_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_kurtosis_package_indexer_proto_goTypes = []interface{}{
	(PackageArgType)(0),         // 0: kurtosis_package_indexer.PackageArgType
	(*IndexerPing)(nil),         // 1: kurtosis_package_indexer.IndexerPing
	(*IndexerPong)(nil),         // 2: kurtosis_package_indexer.IndexerPong
	(*GetPackagesResponse)(nil), // 3: kurtosis_package_indexer.GetPackagesResponse
	(*KurtosisPackage)(nil),     // 4: kurtosis_package_indexer.KurtosisPackage
	(*PackageArg)(nil),          // 5: kurtosis_package_indexer.PackageArg
	(*emptypb.Empty)(nil),       // 6: google.protobuf.Empty
}
var file_kurtosis_package_indexer_proto_depIdxs = []int32{
	4, // 0: kurtosis_package_indexer.GetPackagesResponse.packages:type_name -> kurtosis_package_indexer.KurtosisPackage
	5, // 1: kurtosis_package_indexer.KurtosisPackage.args:type_name -> kurtosis_package_indexer.PackageArg
	0, // 2: kurtosis_package_indexer.PackageArg.type:type_name -> kurtosis_package_indexer.PackageArgType
	1, // 3: kurtosis_package_indexer.KurtosisPackageIndexer.Ping:input_type -> kurtosis_package_indexer.IndexerPing
	6, // 4: kurtosis_package_indexer.KurtosisPackageIndexer.GetPackages:input_type -> google.protobuf.Empty
	2, // 5: kurtosis_package_indexer.KurtosisPackageIndexer.Ping:output_type -> kurtosis_package_indexer.IndexerPong
	3, // 6: kurtosis_package_indexer.KurtosisPackageIndexer.GetPackages:output_type -> kurtosis_package_indexer.GetPackagesResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
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
		file_kurtosis_package_indexer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPackagesResponse); i {
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
		file_kurtosis_package_indexer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KurtosisPackage); i {
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
		file_kurtosis_package_indexer_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PackageArg); i {
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
	file_kurtosis_package_indexer_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_kurtosis_package_indexer_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_kurtosis_package_indexer_proto_goTypes,
		DependencyIndexes: file_kurtosis_package_indexer_proto_depIdxs,
		EnumInfos:         file_kurtosis_package_indexer_proto_enumTypes,
		MessageInfos:      file_kurtosis_package_indexer_proto_msgTypes,
	}.Build()
	File_kurtosis_package_indexer_proto = out.File
	file_kurtosis_package_indexer_proto_rawDesc = nil
	file_kurtosis_package_indexer_proto_goTypes = nil
	file_kurtosis_package_indexer_proto_depIdxs = nil
}