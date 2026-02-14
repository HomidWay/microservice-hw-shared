package panicinterceptor

import (
	"runtime/debug"

	"github.com/HomidWay/microservice-hw-shared/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (pi *PanicInterceptor) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			pi.log.Error("Panic recovered: ", r)
			pi.log.Error("Stack trace: ", debug.Stack())
			err = status.Errorf(codes.Internal, "panic: %v", r)
		}
	}()

	wrappedStream := &PanicWrappedStream{
		ServerStream: ss,
		log:          pi.log,
	}

	err = handler(srv, wrappedStream)

	return err
}

type PanicWrappedStream struct {
	grpc.ServerStream
	log logger.Logger
}

func (pw *PanicWrappedStream) SendMsg(m interface{}) error {

	var err error

	defer func() {
		if r := recover(); r != nil {
			pw.log.Error("Panic recovered: ", r)
			pw.log.Error("Stack trace: ", debug.Stack())
			err = status.Errorf(codes.Internal, "panic: %v", r)
		}
	}()

	err = pw.ServerStream.SendMsg(m)

	return err
}

func (pw *PanicWrappedStream) RecvMsg(m interface{}) error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			pw.log.Error("Panic recovered: ", r)
			pw.log.Error("Stack trace: ", debug.Stack())
			err = status.Errorf(codes.Internal, "panic: %v", r)
		}
	}()

	err = pw.ServerStream.RecvMsg(m)

	return err
}
