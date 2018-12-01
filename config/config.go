package config

import (
	"os"
)

var RedisURL = "127.0.0.1:6379"

func init() {
	if os.Getenv("ENV_NAME") == "local" {
		RedisURL = "127.0.0.1:6379"
	}
}
