package main

import (
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/logger"
	"go-svc-metrics/internal/server"

	"go.uber.org/zap"
)

func main() {
	err := logger.Initialize("INFO")
	if err != nil {
		logger.Log.Fatal("cannot initialize zap", zap.Error(err))
	}

	configServe, err := config.NewServerConfig()
	if err != nil {
		logger.Log.Fatal("cannot initialize config", zap.Error(err))
	}

	server, err := server.NewApp(configServe)
	if err != nil {
		logger.Log.Fatal("cannot initialize server", zap.Error(err))
	}

	err = server.Run()
	if err != nil {
		logger.Log.Fatal("cannot start server", zap.Error(err))
	}
}
