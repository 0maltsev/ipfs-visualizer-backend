package main

import (
	"context"
	"ipfs-visualizer/config"
	"ipfs-visualizer/internal/app"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("cannot load config", "error", err)
		os.Exit(1)
	}
	newApplication := app.NewApp(cfg)

	slog.Info("config", "config", *cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = newApplication.Start(ctx)
	if err != nil {
		slog.Error("", "error", err)
		os.Exit(1)
	}
}

