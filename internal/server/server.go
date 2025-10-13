package server

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/datasource"
	"go-svc-metrics/internal/domain"
	"go-svc-metrics/internal/domain/local"
	"go-svc-metrics/internal/domain/postgres"
	"go-svc-metrics/internal/logger"
	"go-svc-metrics/internal/router"
	"go-svc-metrics/internal/service"
	"go-svc-metrics/internal/utils/crypto"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"go.uber.org/zap"
)

type App struct {
	db         *sql.DB
	cfg        *config.Config
	privateKey *rsa.PrivateKey
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := datasource.NewDatabase(*cfg.DatabaseDsn)
	if err != nil {
		logger.Log.Warn("not connected db")
	}

	privatekey, err := crypto.GetPrivateKey(*cfg.CryptoKey)
	if err != nil {
		return nil, err
	}

	return &App{
		db:         db,
		cfg:        cfg,
		privateKey: privatekey,
	}, nil
}

func (app *App) Run() error {
	var metricRepo domain.MetricRepo
	switch app.db {
	case nil:
		localRepo, err := local.NewMetricLocalRepository(app.cfg)
		if err != nil {
			return err
		}
		metricRepo = localRepo
	default:
		postgresRepo := postgres.NewMetricRepository(app.db)
		metricRepo = postgresRepo
	}
	serviceApp := service.NewMetricService(metricRepo)

	r := router.NewRouter(serviceApp, app.cfg, app.privateKey)

	server := &http.Server{Addr: *app.cfg.ServerAddr, Handler: r}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, time.Duration(*app.cfg.Wait)*time.Second)
		defer shutdownCancel()
		defer serviceApp.Close()

		go serviceApp.DumpMetricsByInterval(shutdownCtx)

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

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Fatal("cannot start server", zap.Error(err))
		return err
	}

	<-serverCtx.Done()
	return nil
}
