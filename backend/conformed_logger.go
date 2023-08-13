package main

import "github.com/rs/zerolog/log"

type conformedLogger struct{}

func (c *conformedLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v)
}

func (c *conformedLogger) Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v)
}
