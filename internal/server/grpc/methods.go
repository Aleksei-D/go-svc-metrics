package server

import (
	"context"
	pb "go-svc-metrics/internal/pb/metric"
	"go-svc-metrics/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (m *MetricServer) GetAll(ctx context.Context, _ *emptypb.Empty) (*pb.BatchMetricsMessage, error) {
	var metricsResponse pb.BatchMetricsMessage

	metrics, err := m.metricService.GetAllMetrics(ctx)
	if err != nil {
		return &metricsResponse, status.Error(codes.Unknown, err.Error())
	}

	for _, metric := range metrics {
		metricsResponse.Metrics = append(
			metricsResponse.Metrics,
			metric.ToProto(),
		)
	}

	return &metricsResponse, nil
}

func (m *MetricServer) V1Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := m.metricService.Ping()
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (m *MetricServer) V1UpdateMetric(ctx context.Context, in *pb.MetricMessage) (*pb.MetricMessage, error) {
	var metric models.Metrics
	metric = metric.FromProto(in)

	updatedMetrics, err := m.metricService.UpdateMetrics(ctx, []models.Metrics{metric})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	return updatedMetrics[0].ToProto(), nil
}

func (m *MetricServer) V1UpdateManyMetrics(ctx context.Context, in *pb.BatchMetricsMessage) (*pb.BatchMetricsMessage, error) {
	var metrics models.BatchMetrics

	updatedMetrics, err := m.metricService.UpdateMetrics(ctx, metrics.FromProto(in))
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	return updatedMetrics.ToProto(), nil
}

func (m *MetricServer) V1GetMetric(ctx context.Context, in *pb.MetricMessage) (*pb.MetricMessage, error) {
	var metricReq models.Metrics

	metric, err := m.metricService.GetMetric(ctx, metricReq.FromProto(in))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return metric.ToProto(), nil
}
