package redisservice

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"ndewo-mobile-backend/config"
	"ndewo-mobile-backend/src/helpers"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

func NewConnection(config config.ConfigType) Redis {
	var r Redis
	credentials := helpers.ExtractURICredentials(config.RedisUri, "redis://")
	conn := redis.NewClient(&redis.Options{
		Addr:      credentials.BaseUrl,
		Password:  credentials.Secret,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		DB:        0, // use default DB
	})

	r.Client = conn

	return r
}

func (r Redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func IsOpen(ctx context.Context, r Redis) bool {
	status := r.Client.Ping(ctx)
	fmt.Println(status)
	return strings.EqualFold(status.Val(), "PONG")
}

func (r Redis) AddToStream(ctx context.Context, args redis.XAddArgs) string {
	return r.Client.XAdd(ctx, &args).String()
}

func (r Redis) GetValue(ctx context.Context, key string) string {
	value := r.Client.Get(ctx, key).Val()
	return value
}

func (r Redis) GetIntValue(ctx context.Context, key string) int {
	value, _ := r.Client.Get(ctx, key).Int()
	return value
}

func (r Redis) KeyExists(ctx context.Context, key string) int64 {
	return r.Client.Exists(ctx, key).Val()
}

func (r Redis) JsonSet(ctx context.Context, key string, value map[string]interface{}) error {
	m, err := json.Marshal(value)
	if err != nil {
		return err
	}
	val := r.Client.Do(ctx, "JSON.SET", key, "$", m)
	return val.Err()
}

func (r Redis) JsonSetArray(ctx context.Context, key string, value interface{}) error {
	m, err := json.Marshal(value)
	if err != nil {
		return err
	}
	val := r.Client.Do(ctx, "JSON.SET", key, "$", m)
	return val.Err()
}

func (r Redis) JsonGet(ctx context.Context, key string) (interface{}, error) {
	return r.Client.Do(ctx, "JSON.GET", key).Result()
}

func (r Redis) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r Redis) DeleteByPattern(ctx context.Context, pattern string) error {
	result, err := r.Client.Do(ctx, "KEYS", pattern).StringSlice()
	if err != nil {
		return err
	}
	for _, v := range result {
		r.Delete(ctx, v)
	}
	return nil
}
