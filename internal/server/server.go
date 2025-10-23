package server

import (
	"crypto/rsa"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/utils/crypto"

	_ "net/http/pprof"
)

type App struct {
	cfg        *config.Config
	PrivateKey *rsa.PrivateKey
}

func NewApp(cfg *config.Config) (*App, error) {
	privatekey, err := crypto.GetPrivateKey(*cfg.CryptoKey)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:        cfg,
		PrivateKey: privatekey,
	}, nil
}
