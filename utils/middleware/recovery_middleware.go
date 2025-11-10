package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/anugrahsputra/go-quran-api/domain/dto"
	"github.com/anugrahsputra/go-quran-api/utils/helper"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("middleware")

// Recovery middleware recovers from panics and returns a proper error response
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				logger.Errorf("Panic recovered - Error: %v, Path: %s, Method: %s, Stack: %s",
					err,
					c.Request.URL.Path,
					c.Request.Method,
					string(debug.Stack()),
				)

				// Return safe error message
				errorMsg := "An internal error occurred"
				if !helper.IsProduction() {
					// In development, include more details
					errorMsg = fmt.Sprintf("Panic: %v", err)
				}

				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: errorMsg,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
