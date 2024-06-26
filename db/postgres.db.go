package db

import (
	"fmt"
	"ndewo-mobile-backend/config"
	"ndewo-mobile-backend/src/models"

	"github.com/rs/zerolog/log"
	zlog "github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectPostgresDB(config config.ConfigType) *gorm.DB {
	dialect := postgres.Open(fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s TimeZone=Africa/Lagos",
		config.PGHost,
		config.PGPort,
		config.PGUser,
		config.PGDatabase,
		config.PGPassword,
	))

	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err)
		panic(err)
	}

	return db
}

// MigrateModels performs database migration for the models.
//
// It takes a *gorm.DB as a parameter and returns no values.
func MigrateModels(db *gorm.DB) {
	err := db.AutoMigrate(models.AllModels...)
	if err != nil {
		zlog.Fatal().Msgf("Could not initialize database migration: %v", err)

	}
}
