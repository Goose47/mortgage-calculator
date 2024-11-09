package controllers

import "github.com/gin-gonic/gin"

// CacheController deals with cache endpoints.
type CacheController struct{}

// NewCacheController is a constructor for CacheController.
func NewCacheController() *CacheController {
	return &CacheController{}
}

// List lists all active cache entries.
func (con *CacheController) List(c *gin.Context) {
	c.JSON(200, gin.H{"message": "hello"})
}
