package main

import (
	"fmt"
	"net/http"
	"os"

	"ndewo-mobile-backend/config"
	"ndewo-mobile-backend/db"
	"ndewo-mobile-backend/docs"
	"ndewo-mobile-backend/src/api/handlers"
	"ndewo-mobile-backend/src/api/routes"
	"ndewo-mobile-backend/src/helpers"
	"ndewo-mobile-backend/src/models"

	"github.com/gin-contrib/cors"
	zlog "github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Configures system wide Logger object

	zlog.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	// make it human-readable, only locally
	if IsLocal() {
		zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	}

	configVariables := config.GetConfig()
	// if configVariables.AppEnv != "dev" {
	// 	runMigrations(configVariables)
	// }

	dbConn := db.ConnectDB(*configVariables)
	db.MigrateModels(dbConn.PostgresDb)

	defer dbConn.RedisClient.Close()
	// defer dbConn.PostgresDb.DB.Close()

	close, err := dbConn.PostgresDb.DB()
	if err != nil {
		zlog.Fatal().Msgf("listen: %s", err)

	}
	close.Close()

	handler := handlers.NewHandler(&dbConn, configVariables)

	server := gin.New()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowAllOrigins = true
	server.Use(handler.JSONLogMiddleware(configVariables.LogLevel, []string{"/health"}))
	server.Use(cors.New(corsConfig), gin.Recovery())
	//server.Use(httpLogger.CustomLog())

	// register routes
	r := routes.New(handler)

	if helpers.StringToBoolean(configVariables.EnableSwagger) {
		// programmatically set swagger info
		docs.SwaggerInfo.Title = "Ndewo Mobile App APIs"
		docs.SwaggerInfo.Description = "This is the API docs for Ndewo Mobile App"
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.BasePath = "/"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}

		url := ginSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", configVariables.AppHost))
		//server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
		server.GET("/swagger/*any", func(c *gin.Context) {
			if c.Param("any") == "/" || c.Param("any") == "" {
				c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
			} else {
				ginSwagger.WrapHandler(swaggerFiles.Handler, url)(c)
			}
		})
	}
	r.Routes(server)
	//run server
	if err := server.Run(fmt.Sprintf("%s:%s", configVariables.AppHost, configVariables.Port)); err != nil && err != http.ErrServerClosed {
		zlog.Fatal().Msgf("listen: %s", err)
	}

}

func IsLocal() bool {
	return os.Getenv("APP_ENV") == "" || os.Getenv("APP_ENV") == models.Local
}
