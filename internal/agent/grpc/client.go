package client

import (
	"context"
	"go-svc-metrics/internal/config"
	pb "go-svc-metrics/internal/pb/metric"
	"go-svc-metrics/models"

	"go-svc-metrics/internal/interceptors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ClientAgent struct {
	conn         *grpc.ClientConn
	metricClient pb.MetrcicClient
}

func NewClientAgent(agentConfig *config.Config) (*ClientAgent, error) {
	opts := []grpc.DialOption{}

	if agentConfig.CryptoKey != nil {
		tlsCreds, err := generateTLSCreds(*agentConfig.CryptoKey)
		if err != nil {
			return nil, err
		}

		opts = append(opts, grpc.WithTransportCredentials(tlsCreds))
	}

	if agentConfig.RealIP != nil {
		opts = append(opts, grpc.WithChainUnaryInterceptor(interceptors.NewRealIPClientInterceptor(*agentConfig.RealIP)))
	}

	conn, err := grpc.NewClient(*agentConfig.AddrGRPC, opts...)
	if err != nil {
		return nil, err
	}

	return &ClientAgent{
		conn:         conn,
		metricClient: pb.NewMetrcicClient(conn),
	}, nil
}

func (c *ClientAgent) ConnClose() {
	c.conn.Close()
}

// SendBatchMetrics отправляет батч метрик на сервер.
func (c *ClientAgent) SendBatchMetrics(ctx context.Context, metrics []models.Metrics) error {
	var batcnMetrics models.BatchMetrics = metrics

	_, err := c.metricClient.V1UpdateManyMetrics(ctx, batcnMetrics.ToProto())
	if err != nil {
		return err
	}
	return nil
}

func generateTLSCreds(certFile string) (credentials.TransportCredentials, error) {
	return credentials.NewClientTLSFromFile(certFile, "")
}
