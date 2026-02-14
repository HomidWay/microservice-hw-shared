package panicinterceptor

import (
	"github.com/HomidWay/microservice-hw-shared/logger"
	"google.golang.org/grpc"
)

type PanicInterceptor struct {
	log logger.Logger
}

func NewPanicInterceptor(log logger.Logger) PanicInterceptor {
	return PanicInterceptor{log: log}
}

func (pi *PanicInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return pi.unaryInterceptor
}

func (pi *PanicInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return pi.streamInterceptor
}
