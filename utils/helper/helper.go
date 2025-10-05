package helper

import "os"

func GetEnv(key, defaultValue string) string {
	getEnv := os.Getenv(key)

	// return default value if env is not set
	if len(getEnv) == 0 || getEnv == "" {
		return defaultValue

	}
	return getEnv
}
