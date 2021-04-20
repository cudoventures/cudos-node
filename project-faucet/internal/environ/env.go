package environ

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func GetInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}

	return fallback
}

func GetUint64(key string, fallback uint64) uint64 {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.ParseUint(value, 10, 64); err == nil {
			return i
		}
	}

	return fallback
}

func GetBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		return value == "true"
	}

	return fallback
}
