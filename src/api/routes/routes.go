package routes

import (
	"ndewo-mobile-backend/src/api/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	handler handlers.Operations
}

func New(h handlers.Operations) Routes {
	return Routes{handler: h}
}

// Routes registers all system routes
func (ro Routes) Routes(r *gin.Engine) {
	CheckRoutes(r)
}

func CheckRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Ndewo API",
		})
	})
}
