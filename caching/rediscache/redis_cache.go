package rediscache

import (
	"context"
	"fmt"
	"time"

	"github.com/HomidWay/microservice-hw-shared/caching"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

type RedisDB struct {
	rdb *redis.Client

	ttl time.Duration
	ctx context.Context
}

func NewRedis(ctx context.Context, host string, port int, password string, db int, ttl time.Duration) (*RedisDB, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	go func() {
		<-ctx.Done()
		rdb.Close()
	}()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis with error: %w", err)
	}

	return &RedisDB{
		rdb: rdb,
		ttl: ttl,
		ctx: ctx,
	}, nil
}

func (r *RedisDB) Set(key string, value proto.Message) error {

	data, err := proto.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	exists, err := r.rdb.Exists(r.ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check if key exists: %w", err)
	}

	if exists > 0 {
		return nil
	}

	err = r.rdb.Set(r.ctx, key, data, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to store result: %w", err)
	}

	return nil
}

func (r *RedisDB) Get(key string, out *proto.Message) error {
	cmd := r.rdb.Get(r.ctx, key)

	if cmd.Err() != nil {
		return caching.ErrKeyNotFound
	}

	data, err := cmd.Bytes()
	if err != nil {
		return fmt.Errorf("failed to get key %s from Redis with error: %w", key, err)
	}

	err = proto.Unmarshal(data, *out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

func (r *RedisDB) Delete(key string) error {
	cmd := r.rdb.Del(r.ctx, key)

	if cmd.Err() != nil {
		return fmt.Errorf("failed to delete key %s from Redis with error: %w", key, cmd.Err())
	}

	return nil
}
