package panicinterceptor

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (pi *PanicInterceptor) unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var resp interface{}
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v", r)
			log.Printf("Stack trace: %s", debug.Stack())
			err = status.Errorf(codes.Internal, "panic: %v", r)
		}
	}()

	resp, err = handler(ctx, req)

	return resp, err
}
