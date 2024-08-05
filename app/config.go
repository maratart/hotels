package app

import (
	"time"
)

type Config struct {
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewConfig() Config {
	return Config{
		port:         "8080",
		readTimeout:  5 * time.Second,
		writeTimeout: 10 * time.Second,
	}
}
