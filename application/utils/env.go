package utils

import "os"

type env struct {
	REDIS_URL string
}

var ENV = env{
	REDIS_URL: os.Getenv("REDIS_URL"),
}

func EnvironmentInitialize() {
	ENV = env{
		REDIS_URL: os.Getenv("REDIS_URL"),
	}
	if ENV.REDIS_URL == "" {
		panic("REDIS_URL is empty")
	}
}
