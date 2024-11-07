// Package server defines router settings and application endpoints.
package server

import (
	"github.com/gin-gonic/gin"
	envpkg "mortgage-calculator/src/internal/lib/env"
)

// NewRouter sets router mode based on env, defines handlers and options and creates new gin router.
func NewRouter(env string) *gin.Engine {
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
