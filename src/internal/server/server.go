// Package server defines router settings and application endpoints.
package server

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"mortgage-calculator/src/internal/controllers"
	envpkg "mortgage-calculator/src/internal/lib/env"
	"mortgage-calculator/src/internal/lib/server/middleware"
)

// NewRouter sets router mode based on env, registers middleware, defines handlers and options and creates new gin router.
func NewRouter(
	log *slog.Logger,
	env string,
	calcCon *controllers.CalcController,
	cacheCon *controllers.CacheController,
) *gin.Engine {
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

	r.POST("execute", calcCon.Calculate)
	r.GET("cache", cacheCon.List)

	return r
}
