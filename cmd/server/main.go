package main

import (
	"context"
	"errors"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/logger"
	chiRouter "go-svc-metrics/internal/server"
	"go-svc-metrics/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	r := chiRouter.GetMetricRouter(metricStorage)

	server := &http.Server{Addr: configServe.GetServeAddress(), Handler: r}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancel()
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	go func() {
		err := storeProducer.DumpMetricsByInterval(serverCtx)
		if err != nil {
			logger.Log.Fatal(err.Error())
		}
	}()

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Fatal(err.Error())
	}
}
