package resource

import (
	"connectrpc.com/connect"
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated/generatedconnect"
	"google.golang.org/protobuf/types/known/emptypb"
)

// KurtosisPackageIndexerHandlerImpl is the middle layer required by connect-go
// Its only purpose is to wrap and unwrap request / response and forward them to the underlying
// generated.KurtosisPackageIndexerServer
type KurtosisPackageIndexerHandlerImpl struct {
	resource generated.KurtosisPackageIndexerServer
}

func NewKurtosisPackageIndexerHandlerImpl(resource generated.KurtosisPackageIndexerServer) generatedconnect.KurtosisPackageIndexerHandler {
	return &KurtosisPackageIndexerHandlerImpl{
		resource: resource,
	}
}

func (handler KurtosisPackageIndexerHandlerImpl) IsAvailable(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	resp, err := handler.resource.IsAvailable(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (handler KurtosisPackageIndexerHandlerImpl) GetPackages(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[generated.GetPackagesResponse], error) {
	resp, err := handler.resource.GetPackages(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (handler KurtosisPackageIndexerHandlerImpl) Reindex(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	resp, err := handler.resource.Reindex(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (handler KurtosisPackageIndexerHandlerImpl) ReadPackage(ctx context.Context, req *connect.Request[generated.ReadPackageRequest]) (*connect.Response[generated.ReadPackageResponse], error) {
	resp, err := handler.resource.ReadPackage(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (handler KurtosisPackageIndexerHandlerImpl) RereadLoggerLevel(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	resp, err := handler.resource.RereadLoggerLevel(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
