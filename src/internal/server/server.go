package server

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
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
