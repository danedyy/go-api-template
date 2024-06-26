package handlers

import (
	"ndewo-mobile-backend/config"
	"ndewo-mobile-backend/db"
	"ndewo-mobile-backend/src/api/controllers"
	"ndewo-mobile-backend/src/common/middleware"
	"ndewo-mobile-backend/src/models"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	controller controllers.Operations
	log        *zerolog.Logger
	config     *config.ConfigType
}

const API_END_POINT_STRING = "api/v1"

type Operations interface {
	// middlewares
	AuthenticatedUserMiddleware() gin.HandlerFunc
	StateTokenMiddleware() gin.HandlerFunc
	JSONLogMiddleware(serName string, exclude []string) gin.HandlerFunc
}

func NewHandler(db *db.Database, config *config.ConfigType) Operations {
	newMiddleware, err := middleware.NewMiddleware(db, config)
	if err != nil {
		log.Fatal().Msgf("middleware error: %v", err)
	}

	h := Handler{
		controller: controllers.New(db, config, newMiddleware),
		config:     config,
	}
	op := Operations(&h)
	return op
}

func getPagingInfo(c *gin.Context) models.APIPagingDto {
	var paging models.APIPagingDto

	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	paging.Filter = c.Query("filter")

	paging.Limit = limit
	paging.Page = page
	return paging
}
