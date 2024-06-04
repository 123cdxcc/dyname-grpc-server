package main

import (
	"dyname-grpc-server/api"
	"dyname-grpc-server/service"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
)

var address string
var _defaultAddress = "dynamic-grpc.socket"

func init() {
	flag.StringVar(&address, "address", _defaultAddress, "监听地址")
	flag.Parse()
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	var listener net.Listener
	var err error
	if address == _defaultAddress {
		{
			err := os.Remove(_defaultAddress)
			if err != nil && !os.IsNotExist(err) {
				panic(err)
			}
		}
		listener, err = net.Listen("unix", _defaultAddress)
	} else {
		listener, err = net.Listen("tcp", address)
	}
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("server listening %s", address))
	server := grpc.NewServer()
	api.RegisterGrpcHelperServer(server, service.NewGrpcHelperService())
	reflection.Register(server)
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
