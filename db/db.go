package db

import (
	"ndewo-mobile-backend/config"

	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
)

type Database struct {
	PostgresDb  *gorm.DB
	RedisClient *redis.Client
}

func ConnectDB(config config.ConfigType) Database {
	redis := connectRedis(config)
	db := Database{
		RedisClient: redis,
		PostgresDb:  connectPostgresDB(config),
	}
	return db
}
