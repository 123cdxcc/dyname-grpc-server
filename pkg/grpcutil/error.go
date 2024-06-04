package grpcutil

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorNotFound 不存在
func ErrorNotFound() error {
	return status.Error(codes.NotFound, "not found")
}

// ErrorInternalError 内部服务错误
func ErrorInternalError() error {
	return status.Error(codes.NotFound, "not found")
}
