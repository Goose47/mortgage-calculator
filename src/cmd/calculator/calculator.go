// application entrypoint
package main

import (
	"context"
	"log/slog"
	apppkg "mortgage-calculator/src/internal/app"
	"mortgage-calculator/src/internal/config"
	"mortgage-calculator/src/internal/logger"
	"time"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)
	app := apppkg.New(log, cfg.Env, cfg.Port, int64(cfg.Cache.TTL))

	//todo graceful stop
	go func() {
		tick := time.NewTicker(time.Duration(cfg.Cache.Clear) * time.Second)
		ctx := context.Background()
		for {
			<-tick.C
			app.Cache.Clear(ctx)
		}
	}()

	err := app.Server.Serve()
	log.Error("application has stopped: %s", slog.Any("error", err))
}
