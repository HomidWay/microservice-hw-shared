package loggerinterceptor

import (
	"time"

	"github.com/HomidWay/microservice-hw-shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (li *LoggerInterceptor) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := time.Now()
	requestID, ok := ss.Context().Value("x-request-id").(string)
	if !ok {
		requestID = "UNDEFINED"
	}

	li.log.Info("gRPC stream started",
		zap.String("request_id", requestID),
		zap.String("method", info.FullMethod),
		zap.String("timestamp", startTime.Format(time.RFC3339Nano)),
		zap.Any("request", srv),
	)

	wrappedStream := &LoggedStream{
		ServerStream: ss,
		log:          li.log,
		requestID:    requestID,
		method:       info.FullMethod,
	}

	err := handler(srv, wrappedStream)

	duration := time.Since(startTime)

	if err != nil {
		statusValue, ok := status.FromError(err)
		code := codes.Unknown
		message := err.Error()
		if ok {
			code = statusValue.Code()
			message = statusValue.Message()
		}

		li.log.Error("gRPC stream failed",
			zap.String("request_id", requestID),
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.Int32("code", int32(code)),
			zap.String("error_message", message),
			zap.String("timestamp", time.Now().Format(time.RFC3339Nano)),
		)

		return err
	}

	li.log.Info("gRPC stream completed",
		zap.String("request_id", requestID),
		zap.String("method", info.FullMethod),
		zap.Duration("duration", duration),
		zap.Int32("code", int32(codes.OK)),
		zap.String("timestamp", time.Now().Format(time.RFC3339Nano)),
	)

	return nil
}

type LoggedStream struct {
	grpc.ServerStream
	log       logger.Logger
	requestID string
	method    string
}

func (ls *LoggedStream) SendMsg(m interface{}) error {

	startTime := time.Now()

	err := ls.ServerStream.SendMsg(m)

	if err != nil {
		statusValue, ok := status.FromError(err)
		code := codes.Unknown
		message := err.Error()
		if ok {
			code = statusValue.Code()
			message = statusValue.Message()
		}

		ls.log.Error("gRPC stream failed to send chunk",
			zap.String("request_id", ls.requestID),
			zap.String("method", ls.method),
			zap.Duration("duration", time.Since(startTime)),
			zap.Int32("code", int32(code)),
			zap.String("error_message", message),
			zap.String("timestamp", time.Now().Format(time.RFC3339Nano)),
		)

		return err
	}

	ls.log.Info("gRPC stream chunk sent",
		zap.String("request_id", ls.requestID),
		zap.String("method", ls.method),
		zap.Duration("duration", time.Since(startTime)),
		zap.Int32("code", int32(codes.OK)),
		zap.String("timestamp", time.Now().Format(time.RFC3339Nano)),
	)

	return nil
}

func (ls *LoggedStream) RecvMsg(m interface{}) error {

	startTime := time.Now()

	err := ls.ServerStream.RecvMsg(m)

	if err != nil {
		statusValue, ok := status.FromError(err)
		code := codes.Unknown
		message := err.Error()
		if ok {
			code = statusValue.Code()
			message = statusValue.Message()
		}

		ls.log.Error("gRPC stream failed to recieve chunk",
			zap.String("request_id", ls.requestID),
			zap.String("method", ls.method),
			zap.Duration("duration", time.Since(startTime)),
			zap.Int32("code", int32(code)),
			zap.String("error_message", message),
			zap.String("timestamp", time.Now().Format(time.RFC3339Nano)),
		)

		return err
	}

	ls.log.Info("gRPC stream chunk recieved",
		zap.String("request_id", ls.requestID),
		zap.String("method", ls.method),
		zap.Duration("duration", time.Since(startTime)),
		zap.Int32("code", int32(codes.OK)),
		zap.String("timestamp", time.Now().Format(time.RFC3339Nano)),
	)

	return nil
}
