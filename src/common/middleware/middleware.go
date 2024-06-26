package middleware

import (
	"errors"
	"ndewo-mobile-backend/config"
	"ndewo-mobile-backend/db"
	"ndewo-mobile-backend/src/api/repo"
	"ndewo-mobile-backend/src/common/message"
	"ndewo-mobile-backend/src/common/redisservice"
	"ndewo-mobile-backend/src/common/tokenservice"

	"context"
	"fmt"
	"ndewo-mobile-backend/src/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	authorizationHeader = "Authorization"
	authorizationBearer = "Bearer"
)

type TokenMaker interface {
	CreateAuthRefreshTokens(ctx context.Context, user models.User) (*models.AuthTokens, error)
	VerifyToken(token string) (*Payload, error)
}

type Middleware struct {
	Jwt          TokenMaker
	logger       *zerolog.Logger
	db           *db.Database
	repo         *repo.Repo
	config       *config.ConfigType
	tokenService tokenservice.TokenService
	redis        redisservice.Redis
}

func NewMiddleware(db *db.Database, config *config.ConfigType) (*Middleware, error) {
	l := log.With().Str("middleware", "api").Logger()
	redis := redisservice.Redis{Client: db.RedisClient}
	jwt, err := NewJwtMaker(config, &redis)
	if err != nil {
		return nil, err
	}
	repo := repo.NewRepo(db)
	tokenService := tokenservice.NewTokenService(&redis, repo)
	m := &Middleware{
		Jwt:          jwt,
		logger:       &l,
		config:       config,
		db:           db,
		repo:         repo,
		tokenService: tokenService,
		redis:        redis,
	}

	return m, nil
}

func fetchKey(key string) string {
	splitted := strings.Split(key, ".")
	return splitted[len(splitted)-1]
}

// JwtUserAuth hybrid middleware returns an authorized user
func (m *Middleware) JwtUserAuth(c *gin.Context) (*models.User, error) {
	authorization := c.GetHeader(authorizationHeader)
	if len(authorization) < 0 {
		return nil, message.ErrInvalidToken
	}

	fields := strings.Fields(authorization)
	if len(fields) != 2 {
		return nil, message.ErrInvalidToken
	}
	fmt.Println(fields)
	return m.getUserFromToken(c, fields[1], models.RedisKeys.AccessToken)
}

func (m *Middleware) JwtRefreshTokenAuth(c *gin.Context, token, redisKey string) (*models.User, error) {
	return m.getUserFromToken(c, token, redisKey)
}

func (m *Middleware) getUserFromToken(ctx context.Context, token, redisKey string) (*models.User, error) {
	verified, err := m.Jwt.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(verified.UserID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := m.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.Status != models.Status.Active {
		return nil, message.ErrInactiveUser
	}
	// ensure token is valid on redis too
	value := m.redis.KeyExists(ctx, fmt.Sprintf("%s:%s:%s", redisKey, verified.UserID, verified.ID))
	if value < 1 {
		return nil, message.ErrNoActiveSession
	}
	// get user  permissions

	return user, nil
}

func (m *Middleware) StateTokenAuth(c *gin.Context) (*models.User, error) {
	stateToken := c.Query("stateToken")
	if len(stateToken) < 1 {
		return nil, message.ErrInvalidToken
	}
	decoded := m.tokenService.DecodeStateToken(stateToken)
	user, err := m.tokenService.ValidateStateToken(c, decoded)
	if err != nil {
		return nil, err
	}
	return user, nil
}
