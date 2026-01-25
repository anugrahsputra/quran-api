package helper

import (
	"errors"
	"os"
	"strings"
)

func SanitizeError(err error) string {
	if err == nil {
		return "An unexpected error occurred"
	}

	env := os.Getenv("ENV")
	isProduction := env == "production" || env == "prod"

	if isProduction {
		errMsg := err.Error()

		if strings.Contains(errMsg, "failed to fetch") || strings.Contains(errMsg, "Kemenag") {
			return "Unable to fetch data from external service. Please try again later."
		}
		if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "context deadline exceeded") {
			return "Request timeout. Please try again."
		}
		if strings.Contains(errMsg, "not found") {
			return "Resource not found"
		}
		if strings.Contains(errMsg, "invalid") {
			return "Invalid request parameters"
		}

		return "An internal error occurred. Please try again later."
	}

	return err.Error()
}

func IsProduction() bool {
	env := os.Getenv("ENV")
	return env == "production" || env == "prod"
}

func GetSafeErrorMessage(err error, defaultMsg string) string {
	if err == nil {
		return defaultMsg
	}

	var safeMsg string
	errStr := err.Error()

	switch {
	case strings.Contains(errStr, "not found"):
		safeMsg = "Resource not found"
	case strings.Contains(errStr, "invalid") || strings.Contains(errStr, "bad request"):
		safeMsg = "Invalid request parameters"
	case strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline"):
		safeMsg = "Request timeout. Please try again."
	case strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "forbidden"):
		safeMsg = "Access denied"
	case strings.Contains(errStr, "rate limit") || strings.Contains(errStr, "too many"):
		safeMsg = "Too many requests. Please try again later."
	default:
		if IsProduction() {
			safeMsg = defaultMsg
		} else {
			safeMsg = err.Error()
		}
	}

	return safeMsg
}

func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return errors.New(context + ": " + err.Error())
}
