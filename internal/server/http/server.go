package server

import (
	"context"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/router"
	"go-svc-metrics/internal/service"

	"net/http"
	_ "net/http/pprof"
)

type App struct {
	cfg           *config.Config
	metricService *service.MetricService
	server        *http.Server
}

func NewApp(cfg *config.Config, metricService *service.MetricService) (*App, error) {
	r, err := router.NewRouter(metricService, cfg)
	if err != nil {
		return nil, err
	}

	server := &http.Server{Addr: *cfg.ServerAddr, Handler: r}

	return &App{cfg: cfg, metricService: metricService, server: server}, nil
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
