package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"mortgage-calculator/src/internal/domain/dto"
	reposmock "mortgage-calculator/src/internal/mocks/repos"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup2() (*CacheController, *reposmock.MockCacheGetSaver) {
	repo := new(reposmock.MockCacheGetSaver)
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	con := NewCacheController(log, repo)

	return con, repo
}

func TestNewCacheController(t *testing.T) {
	repo := new(reposmock.MockCacheGetSaver)
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	con := NewCacheController(log, repo)

	require.NotEmpty(t, con)
}

func TestCacheController_List_EmptyCache(t *testing.T) {
	con, r := setup2()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/cache", nil)
	c.Request = req

	r.On("List", mock.Anything).Return([]*dto.CacheEntry{}, nil)

	con.List(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(
		t,
		w.Body.String(),
		"empty cache",
	)
}
