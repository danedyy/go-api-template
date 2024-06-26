package tokenservice

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"ndewo-mobile-backend/src/api/repo"
	"ndewo-mobile-backend/src/common/message"
	"ndewo-mobile-backend/src/common/redisservice"
	"ndewo-mobile-backend/src/helpers"
	"ndewo-mobile-backend/src/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

var defaultTtl = "1h"

type TokenService struct {
	redis *redisservice.Redis
	repo  *repo.Repo
}

func NewTokenService(redis *redisservice.Redis, repo *repo.Repo) TokenService {
	tokenService := TokenService{
		redis: redis,
		repo:  repo,
	}
	return tokenService
}

func (t *TokenService) SetToken(ctx context.Context, key string, ttl *time.Duration) string {
	if ttl == nil {
		d := helpers.GetDurationFromTimeString(defaultTtl)
		ttl = &d
	}
	code := helpers.GenerateRandomNumber(6)
	// hash code
	codeHash := helpers.HashString(code)

	// send to redis
	err := t.redis.Set(ctx, key, codeHash, *ttl)
	fmt.Println("TOKEN ERR", err)
	return code
}

func (t *TokenService) ValidateToken(ctx context.Context, key, token string) bool {
	return t.redis.GetValue(ctx, key) == helpers.HashString(token)
}

func (t *TokenService) SetStateToken(ctx context.Context, userID string, ttl *time.Duration, isMobileFriendly bool) string {
	if ttl == nil {
		d := helpers.GetDurationFromTimeString(defaultTtl)
		ttl = &d
	}
	code := helpers.GenerateRandomNumber(6)
	if !isMobileFriendly {
		code = helpers.GenerateRandomByte(16)
	}
	// hash code
	codeHash := helpers.HashString(code)

	// send to redis
	t.redis.Set(ctx, fmt.Sprintf("%s:%s:%s", models.RedisKeys.DataAuthStateTokens, userID, codeHash), "", *ttl)

	strToReturn := fmt.Sprintf("%s:%s", userID, code)
	return base64.RawURLEncoding.EncodeToString([]byte(strToReturn))
}

func (t *TokenService) DecodeStateToken(stateToken string) models.DecodedStateToken {
	var decoded models.DecodedStateToken
	byteStr, _ := base64.RawURLEncoding.DecodeString(stateToken)
	splitted := strings.Split(string(byteStr), ":")
	if len(splitted) > 1 {
		decoded.UserID = splitted[0]
		decoded.Code = splitted[1]
	}
	return decoded
}

func (t *TokenService) ValidateStateToken(ctx context.Context, decoded models.DecodedStateToken) (*models.User, error) {
	value := t.redis.KeyExists(ctx, fmt.Sprintf("%s:%s:%s", models.RedisKeys.DataAuthStateTokens, decoded.UserID, helpers.HashString(decoded.Code)))
	if value < 1 {
		return nil, message.ErrInvalidToken
	}
	id, err := uuid.Parse(decoded.UserID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	user, err := t.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
