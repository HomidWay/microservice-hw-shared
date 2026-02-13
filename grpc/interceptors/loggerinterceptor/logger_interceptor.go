package loggerinterceptor

import (
	"context"
	"time"

	"github.com/HomidWay/microservice-hw-shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoggerInterceptor struct {
	logger logger.Logger
}

func NewLoggerInterceptor(logger logger.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{logger: logger}
}

func (li LoggerInterceptor) Intercept(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	startTime := time.Now()
	requestID, ok := ctx.Value("x-request-id").(string)
	if !ok {
		requestID = "UNDEFINED"
	}

	li.logger.Info("gRPC request started",
		zap.String("request_id", requestID),
		zap.String("method", info.FullMethod),
		zap.String("timestamp", startTime.Format(time.RFC3339Nano)),
		zap.Any("request", req),
	)

	result, err := handler(ctx, req)

	duration := time.Since(startTime)

	if err != nil {
		statusValue, ok := status.FromError(err)
		code := codes.Unknown
		message := err.Error()
		if ok {
			code = statusValue.Code()
			message = statusValue.Message()
		}

		li.logger.Error("gRPC request failed",
			zap.String("request_id", requestID),
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.Int32("code", int32(code)),
			zap.String("error_message", message),
			zap.String("timestamp", time.Now().Format(time.RFC3339Nano)),
		)
	} else {
		li.logger.Info("gRPC request completed",
			zap.String("request_id", requestID),
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.Int32("code", int32(codes.OK)),
			zap.String("timestamp", time.Now().Format(time.RFC3339Nano)),
		)
	}

	return result, err
}
