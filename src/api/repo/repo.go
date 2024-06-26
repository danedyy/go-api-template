package repo

import (
	"ndewo-mobile-backend/db"
)

const (
	DEFAULTPAGE                  = 1
	DEFAULTLIMIT                 = 10
	PageDefaultSortBy            = "created_at"
	PageDefaultSortDirectionDesc = "desc"
)

type Repo struct {
	Repo *db.Database
}

func NewRepo(db *db.Database) *Repo {
	return &Repo{
		Repo: db,
	}
}
