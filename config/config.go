package config

import (
	"ndewo-mobile-backend/src/helpers"
	"ndewo-mobile-backend/src/models"

	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
	"gopkg.in/go-playground/validator.v9"
)

type ConfigType struct {
	AppHost         string
	AppEnv          string `validate:"required"`
	PGHost          string `validate:"required"`
	PGPort          string `validate:"required"`
	PGUser          string `validate:"required"`
	PGPassword      string `validate:"required"`
	PGDatabase      string `validate:"required"`
	RedisUri        string `validate:"required"`
	MonoSecretKey   string `validate:"required"`
	RedisPort       string
	Port            string
	LogLevel        string
	EnableSwagger   string
	JwtSecret       string `validate:"required"`
	JwtSecretExpiry string `validate:"required"`
}

func GetConfig() *ConfigType {
	if os.Getenv("APP_ENV") != "prod" && os.Getenv("APP_ENV") != "stg" && os.Getenv("APP_ENV") != "beta" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal().Err(err).Msgf("env file error: %s", err.Error())
		}
	}
	// if err := godotenv.Load(".env"); err != nil {
	// 	log.Fatal().Err(err).Msgf("env file error: %s", err.Error())
	// }

	ConfigVariables := ConfigType{
		AppHost:         helpers.Getenv("APP_HOST", "0.0.0.0"),
		AppEnv:          helpers.Getenv("APP_ENV", "local"),
		LogLevel:        helpers.Getenv("LOG_LEVEL", "debug"),
		PGHost:          os.Getenv("POSTGRES_HOST"),
		PGPort:          os.Getenv("POSTGRES_PORT"),
		PGUser:          os.Getenv("POSTGRES_USER"),
		PGPassword:      os.Getenv("POSTGRES_PASSWORD"),
		PGDatabase:      os.Getenv("POSTGRES_DATABASE"),
		RedisUri:        helpers.Getenv("REDIS_URL"),
		RedisPort:       helpers.Getenv("REDISPORT"),
		Port:            helpers.Getenv("PORT", "7000"),
		EnableSwagger:   helpers.Getenv("ENABLE_SWAGGER", "true"),
		JwtSecret:       os.Getenv("JWT_SECRET"),
		JwtSecretExpiry: os.Getenv("JWT_SECRET_EXPIRY"),
		MonoSecretKey:   os.Getenv("MONO_SECRET_KEY"),
	}
	fmt.Printf("App environemt is: %s%s%s\n", models.ColorOrange, ConfigVariables.AppEnv, models.ColorReset)

	validate := validator.New()
	err := validate.Struct(ConfigVariables)

	if err != nil {
		log.Fatal().Err(err).Msgf("env validation error: %s", err.Error())
	}
	return &ConfigVariables
}
