package cachinginterceptor

import (
	"context"
	"errors"

	"github.com/HomidWay/microservice-hw-shared/caching"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func (ci CachingInterceprtor) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if factory, ok := cachableRequests[info.FullMethod]; ok {
		return ci.handleCaching(ctx, req, handler, factory)
	}

	return handler(ctx, req)
}

func (ci CachingInterceprtor) handleCaching(ctx context.Context, req interface{}, handler grpc.UnaryHandler, factory ResponseFactory) (interface{}, error) {
	cacheKey, err := generateKey(req)
	if err != nil {
		ci.log.Error("Failed to generate key", err.Error())
		return handler(ctx, req)
	}

	result := factory()

	err = ci.cache.Get(cacheKey, &result)
	if err == nil {
		return handler(ctx, req)
	}

	if !errors.Is(err, caching.ErrKeyNotFound) {
		ci.log.Info(
			"cache contains malformed response for: ",
			zap.String("key", cacheKey),
			zap.Error(err),
		)

		err := ci.cache.Delete(cacheKey)
		if err != nil {
			ci.log.Error("Failed to delete cached response", zap.String("key", cacheKey), zap.Error(err))
		}
	}

	handlerResult, err := handler(ctx, req)
	if err == nil {
		if resultMsg, ok := handlerResult.(proto.Message); ok {
			ci.log.Info("Saving cached result", zap.String("key", cacheKey))

			err = ci.cache.Set(cacheKey, resultMsg)
			if err != nil {
				ci.log.Error("Failed to save cached result", zap.String("key", cacheKey), zap.Error(err))
			}
		}
	}

	return handlerResult, err
}
