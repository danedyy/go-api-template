package mono

import (
	"ndewo-mobile-backend/config"
	"ndewo-mobile-backend/src/common/rest"
)

type Anchor struct {
	config  *config.ConfigType
	rest    *rest.RestClient
	headers map[string]string
}

func NewMono(c *config.ConfigType) *Anchor {
	headers := map[string]string{
		"content-type": "application/json",
		"accept":       "application/json",
		"mono-sec-key": c.MonoSecretKey,
	}
	return &Anchor{
		config:  c,
		headers: headers,
		rest:    rest.NewRestClient("https://api.withmono.com/v2"),
	}
}
