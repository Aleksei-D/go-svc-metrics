package main

import (
	"context"
	"errors"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/domain"
	"go-svc-metrics/internal/logger"
	"go-svc-metrics/internal/router"
	"go-svc-metrics/internal/server"
	"go-svc-metrics/internal/service"
	"go-svc-metrics/internal/utils/helpers"
	"net/http"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var buildVersion, buildDate, buildCommit string

func main() {
	helpers.PrintBuildVersion(buildVersion, buildDate, buildCommit)
	err := logger.Initialize("INFO")
	if err != nil {
		logger.Log.Fatal("cannot initialize zap", zap.Error(err))
	}

	configServe, err := config.NewServerConfig()
	if err != nil {
		logger.Log.Fatal("cannot initialize config", zap.Error(err))
	}

	app, err := server.NewApp(configServe)
	if err != nil {
		logger.Log.Fatal("cannot initialize server", zap.Error(err))
	}

	if err != nil {
		logger.Log.Fatal("cannot start server", zap.Error(err))
	}

	repo, err := domain.NewRepo(configServe)
	if err != nil {
		logger.Log.Fatal("cannot init repo", zap.Error(err))
	}

	serviceApp := service.NewMetricService(repo)

	r, err := router.NewRouter(serviceApp, configServe, app.PrivateKey)
	if err != nil {
		logger.Log.Fatal("cannot init middleware", zap.Error(err))
	}

	server := &http.Server{Addr: *configServe.ServerAddr, Handler: r}

	serverCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal(err.Error())
		}
	}()

	<-serverCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), configServe.Wait.Duration)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Fatal(err.Error())
	}

	serviceApp.DumpMetricsByInterval(shutdownCtx)
}
