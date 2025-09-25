package util

import (
	"fmt"
	"os"
)

func GetEnv(key string, defaultValue string) string {
	v, found := os.LookupEnv(key)
	if found {
		return v
	}
	return defaultValue
}

func MustGetEnv(key string) string {
	v, found := os.LookupEnv(key)
	if !found {
		panic(fmt.Sprintf("Environment variable with name %s is required\n", key))
	}

	return v
}
