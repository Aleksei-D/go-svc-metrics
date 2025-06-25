package main

import (
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/logger"
	chiRouter "go-svc-metrics/internal/server"
	"go-svc-metrics/internal/storage"
	"net/http"
)

func main() {
	configServe, err := config.GetServerConfig()
	if err != nil {
		logger.Log.Fatal("cannot initialize config")
	}

	err = logger.Initialize(*configServe.LogLevel)
	if err != nil {
		logger.Log.Fatal("cannot initialize zap")
	}

	metricStorage := storage.InitMemStorage()
	storeConsumer, err := storage.NewConsumer(*configServe.FileStoragePath)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	if *configServe.Restore {
		err = storeConsumer.RestoreMetrics(metricStorage)
		if err != nil {
			logger.Log.Fatal(err.Error())
		}
		defer storeConsumer.Close()
	}

	storeProducer, err := storage.NewProducer(*configServe.FileStoragePath, metricStorage, configServe.GetStoreInterval())
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	defer storeProducer.Close()

	go func() {
		err := storeProducer.DumpMetricsByInterval()
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}()

	r := chiRouter.GetMetricRouter(metricStorage)
	err = http.ListenAndServe(configServe.GetServeAddress(), r)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}
