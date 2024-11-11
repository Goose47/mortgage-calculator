// Package app defines application model.
package app

import (
	"context"
	"log/slog"
	serverapp "mortgage-calculator/src/internal/app/server"
	"mortgage-calculator/src/internal/cache/memory"
	cacherepos "mortgage-calculator/src/internal/cache/repos"
	"mortgage-calculator/src/internal/controllers"
	"mortgage-calculator/src/internal/server"
	"mortgage-calculator/src/internal/services"
)

type clearer interface {
	Clear(ctx context.Context)
}

// App represents application.
type App struct {
	Server *serverapp.Server
	Cache  clearer
}

// New creates all dependencies for App and returns new App instance.
func New(
	log *slog.Logger,
	env string,
	port int,
	cacheTTL int64,
) *App {
	cache := memory.New(log, cacheTTL)
	repo := cacherepos.NewCalcRepository(log, cache)

	calcService := services.NewCalculatorService(log)

	calcCon := controllers.NewCalcController(calcService, repo)
	cacheCon := controllers.NewCacheController(log, repo)

	router := server.NewRouter(log, env, calcCon, cacheCon)
	serverApp := serverapp.New(log, port, router)

	return &App{
		Server: serverApp,
		Cache:  repo,
	}
}
