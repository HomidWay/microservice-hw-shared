package cachinginterceptor

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/HomidWay/microservice-hw-shared/caching"
	"github.com/HomidWay/microservice-hw-shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type ResponseFactory func() proto.Message

type CachingInterceprtor struct {
	log   logger.Logger
	cache []caching.Cache[proto.Message]
}

func NewCachingInterceptor(cache []caching.Cache[proto.Message], log logger.Logger) CachingInterceprtor {
	return CachingInterceprtor{cache: cache, log: log}
}

func (ci CachingInterceprtor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

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

	var result proto.Message

	for i, cache := range ci.cache {
		result = factory()

		err := cache.Get(cacheKey, &result)

		if errors.Is(err, caching.ErrKeyNotFound) {
			continue
		}

		if err != nil {
			ci.log.Info(
				"cache contains malformed response for: ",
				zap.String("key", cacheKey),
				zap.Error(err),
			)
			continue
		}

		ci.log.Info("cache hit on layer ", i+1, zap.String("key", cacheKey))

		for j := range i {
			err = ci.cache[j].Set(cacheKey, result)
			if err != nil {
				ci.log.Error("Failed to save request cache ", err.Error())
			}
		}

		return result, nil
	}

	handlerResult, err := handler(ctx, req)

	if err == nil {
		if resultMsg, ok := handlerResult.(proto.Message); ok {
			ci.log.Info("Saving cached result", zap.String("key", cacheKey))

			for _, cache := range ci.cache {
				err = cache.Set(cacheKey, resultMsg)
				if err != nil {
					ci.log.Error("Failed to save request cache", err.Error())
				}
			}
		}
	}

	return handlerResult, err
}

func generateKey(data interface{}) (string, error) {
	hasher := sha256.New()
	encoder := json.NewEncoder(hasher)
	if err := encoder.Encode(data); err != nil {
		return "", fmt.Errorf("failed to encode data: %w", err)
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
