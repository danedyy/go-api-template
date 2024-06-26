package repo

import (
	"context"
	"errors"
	"ndewo-mobile-backend/src/common/message"
	"ndewo-mobile-backend/src/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (r *Repo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	db := r.Repo.PostgresDb.WithContext(ctx).Where("id = ?", id).Find(&user)
	if db.Error != nil {
		log.Err(db.Error).Msgf("Basic::GetUserByID error: %v, (%v)", "record not found", db.Error)
		return &user, errors.New("record not found")
	}
	if user.Id == uuid.Nil {
		return nil, message.ErrUserNotFound
	}
	return &user, nil
}
