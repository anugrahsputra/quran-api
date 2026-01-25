package helper

import "os"

func GetEnv(key, defaultValue string) string {
	getEnv := os.Getenv(key)

	if len(getEnv) == 0 || getEnv == "" {
		return defaultValue

	}
	return getEnv
}
