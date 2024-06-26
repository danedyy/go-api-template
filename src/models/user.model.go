package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		Id        uuid.UUID `json:"id" gorm:"column:id;PRIMARY_KEY;type:uuid;default:gen_random_uuid()"`
		Email     string    `json:"email,omitempty"`
		Status    string    `json:"status,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
	}
)
