package handlers

import (
	"net/http"
	"time"

	"ndewo-mobile-backend/src/common/message"
	"ndewo-mobile-backend/src/common/response"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ginHands struct {
	SerName    string
	Path       string
	Latency    time.Duration
	Method     string
	StatusCode int
	ClientIP   string
	MsgStr     string
}

func (h *Handler) AuthenticatedUserMiddleware() gin.HandlerFunc {
	// add the middleware function
	return func(c *gin.Context) {
		user, err := h.controller.Middleware().JwtUserAuth(c)
		if err != nil {
			response.Failure(c, http.StatusBadRequest, message.ErrInvalidInput.Error(), err.Error())
			c.Abort()
		} else {
			c.Set("authUser", *user)
		}
		c.Next()
	}
}

func (h *Handler) StateTokenMiddleware() gin.HandlerFunc {
	// add the middleware function
	return func(c *gin.Context) {
		user, err := h.controller.Middleware().StateTokenAuth(c) // state token guard
		if err != nil {
			response.Failure(c, http.StatusBadRequest, message.ErrInvalidToken.Error(), err.Error())
			c.Abort()
		} else {
			c.Set("stateTokenUser", user)
		}
		c.Next()
	}

}

func (h *Handler) JSONLogMiddleware(serName string, exclude []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		// after request
		var exists bool
		for i := 0; i < len(exclude); i++ {
			if exclude[i] == path {
				exists = true
				break
			}
		}
		if !exists {
			if raw != "" {
				path = path + "?" + raw
			}
			msg := c.Errors.String()
			if msg == "" {
				msg = "Request"
			}
			cData := &ginHands{
				SerName:    serName,
				Path:       path,
				Latency:    time.Since(t),
				Method:     c.Request.Method,
				StatusCode: c.Writer.Status(),
				ClientIP:   c.ClientIP(),
				MsgStr:     msg,
			}

			logSwitch(cData)
		}

	}
}

func logSwitch(data *ginHands) {
	switch {
	case data.StatusCode >= 400 && data.StatusCode < 500:
		{
			log.Warn().Str("ser_name", data.SerName).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
		}
	case data.StatusCode >= 500:
		{
			log.Error().Str("ser_name", data.SerName).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
		}
	default:
		log.Info().Str("ser_name", data.SerName).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
	}
}
