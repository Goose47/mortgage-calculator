// Package middleware provides custom middleware for http server.
package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

// Logger logs time, status code and duration of every request.
func Logger(
	log *slog.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		log.Info(
			fmt.Sprintf(
				"%s status_code: %d, duration: %d ns",
				t.Format("2006/01/02 03:04:05"),
				c.Writer.Status(),
				time.Since(t).Nanoseconds(),
			),
		)
	}
}
