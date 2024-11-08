// Package app defines application model.
package app

import (
	"log/slog"
	"mortgage-calculator/src/internal/api/controllers"
	serverapp "mortgage-calculator/src/internal/app/server"
	"mortgage-calculator/src/internal/server"
	"mortgage-calculator/src/internal/services"
)

// App represents application.
type App struct {
	Server *serverapp.Server
}

// New creates all dependencies for App and returns new App instance.
func New(
	log *slog.Logger,
	env string,
	port int,
) *App {
	calcService := services.NewCalculatorService(log)

	calcCon := controllers.NewCalcController(calcService)
	cacheCon := controllers.NewCacheController()

	router := server.NewRouter(log, env, calcCon, cacheCon)
	serverApp := serverapp.New(log, port, router)

	return &App{
		Server: serverApp,
	}
}
