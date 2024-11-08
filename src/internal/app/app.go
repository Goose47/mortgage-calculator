// Package app defines application model.
package app

import (
	"log/slog"
	serverapp "mortgage-calculator/src/internal/app/server"
	"mortgage-calculator/src/internal/server"
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
	router := server.NewRouter(log, env)
	serverApp := serverapp.New(log, port, router)

	return &App{
		Server: serverApp,
	}
}
