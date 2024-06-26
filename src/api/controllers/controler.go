package controllers

import (
	"ndewo-mobile-backend/config"
	"ndewo-mobile-backend/db"
	"ndewo-mobile-backend/src/api/repo"
	"ndewo-mobile-backend/src/common/middleware"
	"ndewo-mobile-backend/src/common/redisservice"
	"ndewo-mobile-backend/src/common/tokenservice"
)

type Controller struct {
	repo   *repo.Repo
	Config *config.ConfigType

	redisService redisservice.Redis

	middleware   *middleware.Middleware
	tokenService tokenservice.TokenService
}

type Operations interface {
	Middleware() *middleware.Middleware
	// allocation
}

func New(db *db.Database, config *config.ConfigType, middleware *middleware.Middleware) Operations {
	repo := repo.NewRepo(db)
	redis := redisservice.Redis{Client: db.RedisClient}

	c := Controller{
		repo:         repo,
		redisService: redis,
		Config:       config,
		middleware:   middleware,

		tokenService: tokenservice.NewTokenService(&redis, repo),
	}
	op := Operations(&c)

	/// TODO start cron jobs
	// Cron(&c)

	return op

}

func (c *Controller) Middleware() *middleware.Middleware {
	return c.middleware
}
