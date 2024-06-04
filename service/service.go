package service

import (
	"context"
	"dyname-grpc-server/api"
	"dyname-grpc-server/pkg/grpcutil"
)

type GrpcHelperService struct {
	api.UnimplementedGrpcHelperServer
	helperMap map[string]*grpcutil.DynamicGrpcHelper // address -> helper
}

func NewGrpcHelperService() *GrpcHelperService {
	return &GrpcHelperService{
		helperMap: make(map[string]*grpcutil.DynamicGrpcHelper),
	}
}

func (service *GrpcHelperService) getHelper(ctx context.Context, address string) (*grpcutil.DynamicGrpcHelper, error) {
	helper, ok := service.helperMap[address]
	if !ok {
		_, err := service.RefreshService(ctx, &api.RefreshServiceRequest{
			Address: address,
		})
		if err != nil {
			return nil, grpcutil.ErrorInternalError(err)
		}
		helper = service.helperMap[address]
	}
	return helper, nil
}

func (service *GrpcHelperService) ListService(ctx context.Context, req *api.ListServiceRequest) (*api.ListServiceReply, error) {
	helper, err := service.getHelper(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	services := helper.ListService(ctx)
	respServices := make([]*api.ListServiceReply_ServiceItem, 0, len(services))
	for _, service := range services {
		methodNames := make([]string, 0, len(service.GetMethods()))
		for _, method := range service.GetMethods() {
			methodNames = append(methodNames, method.GetName())
		}
		respServices = append(respServices, &api.ListServiceReply_ServiceItem{
			ServiceName: service.GetFullyQualifiedName(),
			MethodName:  methodNames,
		})
	}
	return &api.ListServiceReply{
		Data: respServices,
	}, nil
}

func (service *GrpcHelperService) RefreshService(ctx context.Context, req *api.RefreshServiceRequest) (*api.RefreshServiceReply, error) {
	h, err := grpcutil.NewGrpcHelper(req.Address)
	if err != nil {
		return nil, grpcutil.ErrorInternalError(err)
	}
	err = h.RefreshService(ctx)
	if err != nil {
		return nil, grpcutil.ErrorInternalError(err)
	}
	service.helperMap[req.Address] = h
	return nil, nil
}

func (service *GrpcHelperService) Invoke(ctx context.Context, req *api.InvokeRequest) (*api.InvokeReply, error) {
	helper, err := service.getHelper(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	var output []byte
	err = helper.Invoke(ctx, req.Service, req.Method, []byte(req.JsonParams), &output)
	if err != nil {
		return nil, grpcutil.ErrorInternalError(err)
	}
	return &api.InvokeReply{
		JsonResponse: string(output),
	}, nil
}
