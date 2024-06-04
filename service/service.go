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

func (service *GrpcHelperService) ListService(ctx context.Context, req *api.ListServiceRequest) (*api.ListServiceReply, error) {
	helper, ok := service.helperMap[req.Address]
	if !ok {
		h, err := grpcutil.NewGrpcHelper(req.Address)
		if err != nil {
			return nil, grpcutil.ErrorInternalError()
		}
		service.helperMap[req.Address] = h
		helper = h
		err = helper.RefreshService(ctx)
		if err != nil {
			return nil, grpcutil.ErrorInternalError()
		}
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
