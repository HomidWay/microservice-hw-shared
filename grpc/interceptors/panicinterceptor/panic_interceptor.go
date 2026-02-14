package panicinterceptor

import (
	"context"
	"runtime/debug"

	"github.com/HomidWay/microservice-hw-shared/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PanicInterceptor struct {
	log logger.Logger
}

func NewPanicInterceptor(log logger.Logger) PanicInterceptor {
	return PanicInterceptor{log: log}
}

func (pi PanicInterceptor) Intercept(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var resp interface{}
	var err error

	defer func() {
		if r := recover(); r != nil {
			pi.log.Error("Panic recovered: ", r)
			pi.log.Error("Stack trace: ", debug.Stack())
			err = status.Errorf(codes.Internal, "panic: %v", r)
		}
	}()

	resp, err = handler(ctx, req)

	return resp, err
}
