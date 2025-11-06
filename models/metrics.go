package models

import (
	pb "go-svc-metrics/internal/pb/metric"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

func (m *Metrics) ToProto() *pb.MetricMessage {
	return &pb.MetricMessage{
		Id:    m.ID,
		Type:  m.MType,
		Delta: m.Delta,
		Value: m.Value,
	}
}

func (m *Metrics) FromProto(in *pb.MetricMessage) Metrics {
	m.ID = in.Id
	m.MType = in.Type
	m.Delta = in.Delta
	m.Value = in.Value
	return *m
}

type BatchMetrics []Metrics

func (b BatchMetrics) FromProto(in *pb.BatchMetricsMessage) BatchMetrics {
	for _, message := range in.Metrics {
		var metric Metrics
		b = append(b, metric.FromProto(message))
	}
	return b
}

func (b BatchMetrics) ToProto() *pb.BatchMetricsMessage {
	var protoMessage pb.BatchMetricsMessage

	for _, metric := range b {
		protoMessage.Metrics = append(protoMessage.Metrics, metric.ToProto())
	}
	return &protoMessage
}
