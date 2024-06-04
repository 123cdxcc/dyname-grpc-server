package grpcutil

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorNotFound 不存在
func ErrorNotFound(errs ...error) error {
	var msg = "not found"
	if len(errs) > 0 {
		msg = errors.Join(errs...).Error()
	}
	return status.Error(codes.NotFound, msg)
}

// ErrorInternalError 内部服务错误
func ErrorInternalError(errs ...error) error {
	var msg = "internal err"
	if len(errs) > 0 {
		msg = errors.Join(errs...).Error()
	}
	return status.Error(codes.Internal, msg)
}
