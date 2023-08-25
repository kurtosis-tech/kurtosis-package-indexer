// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: kurtosis_package_indexer.proto

package generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	KurtosisPackageIndexer_Ping_FullMethodName        = "/kurtosis_package_indexer.KurtosisPackageIndexer/Ping"
	KurtosisPackageIndexer_GetPackages_FullMethodName = "/kurtosis_package_indexer.KurtosisPackageIndexer/GetPackages"
)

// KurtosisPackageIndexerClient is the client API for KurtosisPackageIndexer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KurtosisPackageIndexerClient interface {
	Ping(ctx context.Context, in *IndexerPing, opts ...grpc.CallOption) (*IndexerPong, error)
	GetPackages(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetPackagesResponse, error)
}

type kurtosisPackageIndexerClient struct {
	cc grpc.ClientConnInterface
}

func NewKurtosisPackageIndexerClient(cc grpc.ClientConnInterface) KurtosisPackageIndexerClient {
	return &kurtosisPackageIndexerClient{cc}
}

func (c *kurtosisPackageIndexerClient) Ping(ctx context.Context, in *IndexerPing, opts ...grpc.CallOption) (*IndexerPong, error) {
	out := new(IndexerPong)
	err := c.cc.Invoke(ctx, KurtosisPackageIndexer_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kurtosisPackageIndexerClient) GetPackages(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetPackagesResponse, error) {
	out := new(GetPackagesResponse)
	err := c.cc.Invoke(ctx, KurtosisPackageIndexer_GetPackages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KurtosisPackageIndexerServer is the server API for KurtosisPackageIndexer service.
// All implementations should embed UnimplementedKurtosisPackageIndexerServer
// for forward compatibility
type KurtosisPackageIndexerServer interface {
	Ping(context.Context, *IndexerPing) (*IndexerPong, error)
	GetPackages(context.Context, *emptypb.Empty) (*GetPackagesResponse, error)
}

// UnimplementedKurtosisPackageIndexerServer should be embedded to have forward compatible implementations.
type UnimplementedKurtosisPackageIndexerServer struct {
}

func (UnimplementedKurtosisPackageIndexerServer) Ping(context.Context, *IndexerPing) (*IndexerPong, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedKurtosisPackageIndexerServer) GetPackages(context.Context, *emptypb.Empty) (*GetPackagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPackages not implemented")
}

// UnsafeKurtosisPackageIndexerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KurtosisPackageIndexerServer will
// result in compilation errors.
type UnsafeKurtosisPackageIndexerServer interface {
	mustEmbedUnimplementedKurtosisPackageIndexerServer()
}

func RegisterKurtosisPackageIndexerServer(s grpc.ServiceRegistrar, srv KurtosisPackageIndexerServer) {
	s.RegisterService(&KurtosisPackageIndexer_ServiceDesc, srv)
}

func _KurtosisPackageIndexer_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexerPing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KurtosisPackageIndexerServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KurtosisPackageIndexer_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KurtosisPackageIndexerServer).Ping(ctx, req.(*IndexerPing))
	}
	return interceptor(ctx, in, info, handler)
}

func _KurtosisPackageIndexer_GetPackages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KurtosisPackageIndexerServer).GetPackages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KurtosisPackageIndexer_GetPackages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KurtosisPackageIndexerServer).GetPackages(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// KurtosisPackageIndexer_ServiceDesc is the grpc.ServiceDesc for KurtosisPackageIndexer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KurtosisPackageIndexer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kurtosis_package_indexer.KurtosisPackageIndexer",
	HandlerType: (*KurtosisPackageIndexerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _KurtosisPackageIndexer_Ping_Handler,
		},
		{
			MethodName: "GetPackages",
			Handler:    _KurtosisPackageIndexer_GetPackages_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kurtosis_package_indexer.proto",
}