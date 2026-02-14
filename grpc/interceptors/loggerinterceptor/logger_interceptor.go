package loggerinterceptor

import (
	"github.com/HomidWay/microservice-hw-shared/logger"
	"google.golang.org/grpc"
)

type LoggerInterceptor struct {
	log logger.Logger
}

func NewLoggerInterceptor(log logger.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{log: log}
}

func (li *LoggerInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return li.unaryInterceptor
}

func (li *LoggerInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return li.streamInterceptor
}
