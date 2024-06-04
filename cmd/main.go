package main

import (
	"dyname-grpc-server/api"
	"dyname-grpc-server/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	{
		err := os.Remove("./dynamic-grpc.socket")
		if err != nil && !os.IsNotExist(err) {
			panic(err)
		}
	}
	listener, err := net.Listen("unix", "dynamic-grpc.socket")
	if err != nil {
		panic(err)
	}
	slog.Info("server listening dynamic-grpc.socket")
	server := grpc.NewServer()
	api.RegisterGrpcHelperServer(server, service.NewGrpcHelperService())
	reflection.Register(server)
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
