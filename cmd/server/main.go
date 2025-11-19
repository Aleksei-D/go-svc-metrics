package main

import (
	"context"
	"errors"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/domain"
	"go-svc-metrics/internal/logger"
	grpc_server "go-svc-metrics/internal/server/grpc"
	http_server "go-svc-metrics/internal/server/http"
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
	var serverGRPC *grpc_server.App

	err := logger.Initialize("INFO")
	if err != nil {
		logger.Log.Fatal("cannot initialize zap", zap.Error(err))
	}

	configServe, err := config.NewServerConfig()
	if err != nil {
		logger.Log.Fatal("cannot initialize config", zap.Error(err))
	}

	repo, err := domain.NewRepo(configServe)
	if err != nil {
		logger.Log.Fatal("cannot init repo", zap.Error(err))
	}

	serviceApp := service.NewMetricService(repo)

	serverHTTP, err := http_server.NewApp(configServe, serviceApp)
	if err != nil {
		logger.Log.Fatal("cannot initialize server", zap.Error(err))
	}

	if configServe.AddrGRPC != nil {
		app, err := grpc_server.NewApp(serviceApp, configServe)
		if err != nil {
			logger.Log.Fatal(err.Error())
		}
		serverGRPC = app
	}

	serverCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		if err := serverHTTP.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal(err.Error())
		}
	}()

	if serverGRPC != nil {
		go func() {
			if err := serverGRPC.Run(); err != nil {
				logger.Log.Fatal("can not start grpc server", zap.Error(err))
			}
		}()
	}

	<-serverCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), configServe.Wait.Duration)
	defer cancel()

	if err := serverHTTP.Stop(shutdownCtx); err != nil {
		logger.Log.Fatal(err.Error())
	}

	if serverGRPC != nil {
		serverGRPC.Stop()
	}

	serviceApp.DumpMetricsByInterval(shutdownCtx)
}
