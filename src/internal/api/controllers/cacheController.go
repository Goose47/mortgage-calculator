package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"mortgage-calculator/src/internal/domain/dto"
	"net/http"
)

// CacheProvider interacts with cache.
type CacheProvider interface {
	List(ctx context.Context) ([]*dto.CacheEntry, error)
}

// CacheController deals with cache endpoints.
type CacheController struct {
	log   *slog.Logger
	cache CacheProvider
}

// NewCacheController is a constructor for CacheController.
func NewCacheController(
	log *slog.Logger,
	cache CacheProvider,
) *CacheController {
	return &CacheController{
		log:   log,
		cache: cache,
	}
}

// List lists all active cache entries.
func (con *CacheController) List(c *gin.Context) {
	ctx := c.Request.Context()
	entries, err := con.cache.List(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve cache entries",
		})
		return
	}
	if len(entries) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "empty cache",
		})
		return
	}

	c.JSON(http.StatusOK, entries)
}
