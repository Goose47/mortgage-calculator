// Package server defines router settings and application endpoints.
package server

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	envpkg "mortgage-calculator/src/internal/lib/env"
	"mortgage-calculator/src/internal/lib/server/middleware"
)

// NewRouter sets router mode based on env, registers middleware, defines handlers and options and creates new gin router.
func NewRouter(log *slog.Logger, env string) *gin.Engine {
	var mode string
	switch env {
	case envpkg.Local:
	case envpkg.Dev:
		mode = gin.DebugMode
	case envpkg.Prod:
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	r := gin.New()

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	r.Use(middleware.Logger(log))
	r.Use(gin.Recovery())

	r.POST("execute", execute)
	r.GET("cache", cache)

	return r
}

func execute(c *gin.Context) {
	c.JSON(200, gin.H{
		"execute": true,
	})
}

func cache(c *gin.Context) {
	c.JSON(200, gin.H{
		"cache": true,
	})
}
