package config

import (
	"os"
)

var RedisURL = "redis:6379"

func init() {
	if os.Getenv("ENV_NAME") == "local" {
		RedisURL = "localhost:6379"
	}
}
