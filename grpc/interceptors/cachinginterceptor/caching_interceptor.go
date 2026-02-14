package cachinginterceptor

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/HomidWay/microservice-hw-shared/caching"
	"github.com/HomidWay/microservice-hw-shared/logger"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type ResponseFactory func() proto.Message

type CachingInterceprtor struct {
	log   logger.Logger
	cache caching.Cache[proto.Message]
}

func NewCachingInterceptor(cache caching.Cache[proto.Message], log logger.Logger) CachingInterceprtor {
	return CachingInterceprtor{cache: cache, log: log}
}

func (ci CachingInterceprtor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if factory, ok := cachableRequests[info.FullMethod]; ok {
		return ci.handleCaching(ctx, req, handler, factory)
	}

	return handler(ctx, req)
}

func generateKey(data interface{}) (string, error) {
	hasher := sha256.New()
	encoder := json.NewEncoder(hasher)
	if err := encoder.Encode(data); err != nil {
		return "", fmt.Errorf("failed to encode data: %w", err)
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
