package grpcutil

import (
	"context"
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceInfo struct {
	ServiceDesc *desc.ServiceDescriptor
	Methods     map[string]*desc.MethodDescriptor // methodName -> methodDesc
}

type DynamicGrpcHelper struct {
	address    string
	conn       *grpc.ClientConn
	client     *grpcreflect.Client
	serviceMap map[string]*ServiceInfo // serviceName -> info
}

func NewGrpcHelper(address string) (*DynamicGrpcHelper, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("new grpc conn: %w", err)
	}
	client := grpcreflect.NewClientAuto(context.Background(), conn)
	return &DynamicGrpcHelper{
		address:    address,
		conn:       conn,
		client:     client,
		serviceMap: make(map[string]*ServiceInfo),
	}, nil
}

func (gh *DynamicGrpcHelper) Close() {
	gh.client.Reset()
	_ = gh.conn.Close()
}

func (gh *DynamicGrpcHelper) RefreshService(_ context.Context) error {
	services, err := gh.client.ListServices()
	if err != nil {
		return err
	}
	for _, service := range services {
		serviceDesc, err := gh.client.ResolveService(service)
		if err != nil {
			return err
		}
		info := &ServiceInfo{
			ServiceDesc: serviceDesc,
			Methods:     make(map[string]*desc.MethodDescriptor, len(serviceDesc.GetMethods())),
		}
		for _, methodDescriptor := range serviceDesc.GetMethods() {
			info.Methods[methodDescriptor.GetName()] = methodDescriptor
		}
		gh.serviceMap[service] = info
	}
	return nil
}

func (gh *DynamicGrpcHelper) ListService(_ context.Context) []*desc.ServiceDescriptor {
	var s = make([]*desc.ServiceDescriptor, 0, len(gh.serviceMap))
	for _, service := range gh.serviceMap {
		s = append(s, service.ServiceDesc)
	}
	return s
}

func (gh *DynamicGrpcHelper) jsonToProtoMessage(descriptor *desc.MessageDescriptor, data []byte) (*dynamic.Message, error) {
	msg := dynamic.NewMessage(descriptor)
	if len(data) == 0 {
		return msg, nil
	}
	err := msg.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (gh *DynamicGrpcHelper) Invoke(ctx context.Context, serviceName, methodName string, jsonInput []byte, jsonOutput *[]byte) error {
	info, ok := gh.serviceMap[serviceName]
	if !ok {
		return fmt.Errorf("service not found: %s", serviceName)
	}
	method, ok := info.Methods[methodName]
	if !ok {
		return fmt.Errorf("method not found: %s", methodName)
	}
	req, err := gh.jsonToProtoMessage(method.GetInputType(), jsonInput)
	if err != nil {
		return err
	}
	res, err := gh.jsonToProtoMessage(method.GetOutputType(), nil)
	if err != nil {
		return err
	}
	err = gh.conn.Invoke(ctx, fmt.Sprintf("/%s/%s", info.ServiceDesc.GetFullyQualifiedName(), method.GetName()), req, res)
	if err != nil {
		return err
	}
	*jsonOutput, err = res.MarshalJSON()
	if err != nil {
		return err
	}
	return nil
}
