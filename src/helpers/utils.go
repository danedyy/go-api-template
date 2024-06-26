package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	alphabet               = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	MODEL_PREFIX_SEPARATOR = "."
)

func Getenv(variable string, defaultValue ...string) string {
	env := os.Getenv(variable)
	if env == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return env
}
func StringToBoolean(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func GenerateRandomUppercase(length int) string {
	s, _ := uuid.NewRandom()
	return strings.ToUpper(strings.Replace(s.String(), "-", "", -1))
}

type RestUriCredentials struct {
	BaseUrl string
	Id      string
	Secret  string
}

func ExtractURICredentials(uri string, protocol ...string) *RestUriCredentials {
	var url []string

	if len(protocol) > 0 {
		url = strings.Split(uri, protocol[0])
	} else {
		url = strings.Split(uri, "rest://")
	}

	url = strings.Split(url[1], "@")

	if len(url) == 1 {
		return &RestUriCredentials{
			BaseUrl: url[0],
		}
	}

	credentials := strings.Split((url[0]), ":")
	baseUrl := url[1]

	return &RestUriCredentials{
		BaseUrl: baseUrl,
		Id:      credentials[0],
		Secret:  credentials[1],
	}
}

func GetDurationFromTimeString(s string) time.Duration {
	duration, _ := time.ParseDuration(s)
	return duration
}

func GenerateRandomNumber(length int) string {
	return generateRandom(length, 10)
}

func generateRandom(length, randomRange int) string {
	var random string
	for i := 0; i < length; i++ {
		random += string(alphabet[rand.Intn(randomRange)])
	}
	return random
}

func GenerateRandomByte(length int) string {
	return generateRandom(length, len(alphabet))
}

func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
