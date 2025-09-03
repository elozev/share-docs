package util

import "os"

func GetEnv(key string, defaultValue string) string {
	v, found := os.LookupEnv(key)
	if found {
		return v
	}
	return defaultValue
}
