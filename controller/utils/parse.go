package utils

import (
	"log/slog"
	"os"
	"strconv"
)

func ParseBoolFromEnv(envVar string) bool {
	value, err := strconv.ParseBool(os.Getenv(envVar))
	if err != nil {
		slog.Error("Error parsing from environment variables", envVar, err)
		panic(err)
	}
	return value
}
