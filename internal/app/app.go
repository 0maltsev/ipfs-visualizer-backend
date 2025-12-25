package app

import (
	"context"
	"ipfs-visualizer/config"
	"log/slog"
	"net/http"
	"time"
	"database/sql"
)

type App struct {
	router                    http.Handler
	serverCfg                 *config.ServerConfig
	sqlDBCfg                  *config.PostgreSqlConfig
	sqlDBPool                 *sql.DB
}

func NewApp(cfg *config.Config) *App {
	app := &App{}

	app.loadServerCfg(cfg)
	app.loadRoutes()
	app.createStorageConnections(cfg)

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":" + a.serverCfg.ServerAddressPort,
		Handler: a.router,
	}

	slog.Info("Server is running on port " + server.Addr)

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- NewServerError("Start", "failed to start server", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := a.CloseStorageConnections()
		if err != nil {
			slog.Error("failed to close storage connections", "error", err)
		}

		return server.Shutdown(timeout)
	}
}
