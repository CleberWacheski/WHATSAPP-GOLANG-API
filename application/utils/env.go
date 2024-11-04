package utils

import "os"

type env struct {
	BASIC_AUTH_REALM    string
	BASIC_AUTH_USERNAME string
	BASIC_AUTH_PASSWORD string
	POSTGRES_URL        string
}

var ENV = env{
	BASIC_AUTH_REALM:    os.Getenv("BASIC_AUTH_REALM"),
	BASIC_AUTH_USERNAME: os.Getenv("BASIC_AUTH_USERNAME"),
	BASIC_AUTH_PASSWORD: os.Getenv("BASIC_AUTH_PASSWORD"),
	POSTGRES_URL:        os.Getenv("POSTGRES_URL"),
}

func EnvironmentInitialize() {
	ENV = env{
		BASIC_AUTH_REALM:    os.Getenv("BASIC_AUTH_REALM"),
		BASIC_AUTH_USERNAME: os.Getenv("BASIC_AUTH_USERNAME"),
		BASIC_AUTH_PASSWORD: os.Getenv("BASIC_AUTH_PASSWORD"),
		POSTGRES_URL:        os.Getenv("POSTGRES_URL"),
	}
	if ENV.POSTGRES_URL == "" {
		panic("POSTGRES_URL is empty")
	}
}
