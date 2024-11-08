package controllers

import "github.com/gin-gonic/gin"

type CacheController struct{}

func NewCacheController() *CacheController {
	return &CacheController{}
}

func (con *CacheController) List(c *gin.Context) {}
