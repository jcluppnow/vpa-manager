package utils

import (
	"fmt"
	"os"
	"strconv"
)

func ParseBoolFromEnv(envVar string) bool {
	value, err := strconv.ParseBool(os.Getenv(envVar))
	if err != nil {
		fmt.Printf("Error parsing %s from environment variables: %v\n", envVar, err)
		panic(err)
	}
	return value
}
