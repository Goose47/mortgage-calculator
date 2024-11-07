// application entrypoint
package main

import (
	"log/slog"
	apppkg "mortgage-calculator/src/internal/app"
	"mortgage-calculator/src/internal/config"
	"mortgage-calculator/src/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)
	app := apppkg.New(log, cfg.Env, cfg.Port)

	err := app.Server.Serve()
	log.Error("application has stopped: %s", slog.Any("error", err))
}
