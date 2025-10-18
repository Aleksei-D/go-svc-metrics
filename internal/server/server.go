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
	"net/http"
	"os/signal"
	"syscall"

	_ "net/http/pprof"
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

	serverCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal(err.Error())
		}
	}()

	<-serverCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), app.cfg.Wait.Duration)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Fatal(err.Error())
	}

	serviceApp.DumpMetricsByInterval(shutdownCtx)
	return nil
}
