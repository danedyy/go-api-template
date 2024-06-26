package db

import (
	"ndewo-mobile-backend/config"
	"ndewo-mobile-backend/src/common/redisservice"

	"github.com/go-redis/redis/v8"
)

func connectRedis(config config.ConfigType) *redis.Client {
	newRedis := redisservice.NewConnection(config)

	return newRedis.Client
}
