package utils

import "os"

type env struct {
	REDIS_URL           string
	BASIC_AUTH_REALM    string
	BASIC_AUTH_USERNAME string
	BASIC_AUTH_PASSWORD string
	POSTGRES_URL        string
}

var ENV = env{
	REDIS_URL:           os.Getenv("REDIS_URL"),
	BASIC_AUTH_REALM:    os.Getenv("BASIC_AUTH_REALM"),
	BASIC_AUTH_USERNAME: os.Getenv("BASIC_AUTH_USERNAME"),
	BASIC_AUTH_PASSWORD: os.Getenv("BASIC_AUTH_PASSWORD"),
	POSTGRES_URL:        os.Getenv("POSTGRES_URL"),
}

func EnvironmentInitialize() {
	ENV = env{
		REDIS_URL:           os.Getenv("REDIS_URL"),
		BASIC_AUTH_REALM:    os.Getenv("BASIC_AUTH_REALM"),
		BASIC_AUTH_USERNAME: os.Getenv("BASIC_AUTH_USERNAME"),
		BASIC_AUTH_PASSWORD: os.Getenv("BASIC_AUTH_PASSWORD"),
		POSTGRES_URL:        os.Getenv("POSTGRES_URL"),
	}
	if ENV.REDIS_URL == "" {
		panic("REDIS_URL is empty")
	}
}
