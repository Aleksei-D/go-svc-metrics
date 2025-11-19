package server

import (
	"crypto/tls"
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/interceptors"
	pb "go-svc-metrics/internal/pb/metric"
	"go-svc-metrics/internal/service"
	"net"
	"os"

	_ "google.golang.org/grpc/encoding/gzip"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type App struct {
	gRPCServer *grpc.Server
	addrGRPC   string
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", a.addrGRPC)
	if err != nil {
		return err
	}

	if err := a.gRPCServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
}

func NewApp(metricService *service.MetricService, cfg *config.Config) (*App, error) {
	var serverOpts []grpc.ServerOption
	var interceptorsOpts []grpc.UnaryServerInterceptor

	interceptorsOpts = append(interceptorsOpts, interceptors.LoggingInterceptor)

	if cfg.Cert != nil {
		tlsCredentials, err := loadTLSCredentials(*cfg.Cert, *cfg.CryptoKey)
		if err != nil {
			return nil, err
		}

		serverOpts = append(serverOpts, grpc.Creds(tlsCredentials))
	}

	if cfg.TrustedSubnet != nil {
		_, network, err := net.ParseCIDR(*cfg.TrustedSubnet)
		if err != nil {
			return nil, err
		}
		interceptorsOpts = append(interceptorsOpts, interceptors.NewRealIPInterceptor(network))
	}

	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(interceptorsOpts...))
	gRPCServer := grpc.NewServer(serverOpts...)
	registerMetrcicServer(gRPCServer, metricService)

	return &App{
		gRPCServer: gRPCServer,
		addrGRPC:   *cfg.AddrGRPC,
	}, nil
}

func loadTLSCredentials(certFilePath, privateKeyFilePath string) (credentials.TransportCredentials, error) {
	privateKeyPEM, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %v", err)
	}

	certPem, err := os.ReadFile(certFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading certificate key file: %v", err)
	}

	tlsCert, err := tls.X509KeyPair(certPem, privateKeyPEM)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func registerMetrcicServer(gRPCServer *grpc.Server, metricService *service.MetricService) {
	pb.RegisterMetrcicServer(gRPCServer, &MetricServer{metricService: metricService})
}

type MetricServer struct {
	pb.UnimplementedMetrcicServer
	metricService *service.MetricService
}
