package agent

import (
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/models"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

const counterMetricName = "PollCount"

type MetricUpdater struct {
	metrics     map[string]models.Metrics
	memStats    runtime.MemStats
	clientAgent ClientAgent
	*config.Config
}

func (m *MetricUpdater) SendReport() {
	for _, metric := range m.metrics {
		go func() {
			err := m.clientAgent.SendOneMetric(metric)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
}

func (m *MetricUpdater) SendAllReports() {
	if len(m.metrics) == 0 {
		return
	}

	metricsToSend := make([]models.Metrics, 0)
	for _, metric := range m.metrics {
		metricsToSend = append(metricsToSend, metric)
	}
	err := m.clientAgent.SendBatchMetrics(metricsToSend)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (m *MetricUpdater) UpdateMetricFromStats() {
	m.updateGaugeMetric("Alloc", float64(m.memStats.Alloc))
	m.updateGaugeMetric("BuckHashSys", float64(m.memStats.BuckHashSys))
	m.updateGaugeMetric("Frees", float64(m.memStats.Frees))
	m.updateGaugeMetric("GCCPUFraction", m.memStats.GCCPUFraction)
	m.updateGaugeMetric("GCSys", float64(m.memStats.GCSys))
	m.updateGaugeMetric("HeapAlloc", float64(m.memStats.HeapAlloc))
	m.updateGaugeMetric("HeapIdle", float64(m.memStats.HeapIdle))
	m.updateGaugeMetric("HeapInuse", float64(m.memStats.HeapInuse))
	m.updateGaugeMetric("HeapObjects", float64(m.memStats.HeapObjects))
	m.updateGaugeMetric("HeapReleased", float64(m.memStats.HeapReleased))
	m.updateGaugeMetric("HeapSys", float64(m.memStats.HeapSys))
	m.updateGaugeMetric("LastGC", float64(m.memStats.LastGC))
	m.updateGaugeMetric("Lookups", float64(m.memStats.Lookups))
	m.updateGaugeMetric("MCacheInuse", float64(m.memStats.MCacheInuse))
	m.updateGaugeMetric("MCacheSys", float64(m.memStats.MCacheSys))
	m.updateGaugeMetric("MSpanInuse", float64(m.memStats.MSpanInuse))
	m.updateGaugeMetric("MSpanSys", float64(m.memStats.MSpanSys))
	m.updateGaugeMetric("Mallocs", float64(m.memStats.Mallocs))
	m.updateGaugeMetric("NextGC", float64(m.memStats.NextGC))
	m.updateGaugeMetric("NumGC", float64(m.memStats.NumGC))
	m.updateGaugeMetric("NumForcedGC", float64(m.memStats.NumForcedGC))
	m.updateGaugeMetric("OtherSys", float64(m.memStats.OtherSys))
	m.updateGaugeMetric("PauseTotalNs", float64(m.memStats.PauseTotalNs))
	m.updateGaugeMetric("StackInuse", float64(m.memStats.StackInuse))
	m.updateGaugeMetric("StackSys", float64(m.memStats.StackSys))
	m.updateGaugeMetric("Sys", float64(m.memStats.Sys))
	m.updateGaugeMetric("TotalAlloc", float64(m.memStats.TotalAlloc))
	m.updateGaugeMetric("RandomValue", rand.Float64())
}

func (m *MetricUpdater) updateGaugeMetric(nameID string, value float64) {
	metric, ok := m.metrics[nameID]
	switch ok {
	case true:
		*metric.Value = value
	case false:
		m.metrics[nameID] = models.Metrics{
			ID:    nameID,
			MType: models.Gauge,
			Value: &value,
		}
	}
}

func (m *MetricUpdater) updateCounterMetric() {
	metric, ok := m.metrics[counterMetricName]
	switch ok {
	case true:
		*metric.Delta++
	case false:
		m.metrics[counterMetricName] = models.Metrics{
			ID:    counterMetricName,
			MType: models.Counter,
			Delta: new(int64),
		}
	}
}

func (m *MetricUpdater) MetricProcessing() {
	metricsPollTicker := time.NewTicker(m.GetPollInterval())
	metricReportTicker := time.NewTicker(m.GetReportInterval())
	defer metricsPollTicker.Stop()
	defer metricReportTicker.Stop()
	for {
		select {
		case <-metricsPollTicker.C:
			runtime.ReadMemStats(&m.memStats)
			m.UpdateMetricFromStats()
			m.updateCounterMetric()
		case <-metricReportTicker.C:
			m.SendAllReports()
		}
	}
}

func GetNewMetricUpdater() (*MetricUpdater, error) {
	agentConfig, err := config.GetAgentConfig()
	if err != nil {
		return nil, err
	}

	agentClient := ClientAgent{
		config: agentConfig,
		httpClient: &http.Client{
			Transport: &retryRoundTripper{
				maxRetries: 3,
				next:       http.DefaultTransport,
			},
		},
	}
	return &MetricUpdater{
		metrics:     make(map[string]models.Metrics),
		clientAgent: agentClient,
		Config:      agentConfig,
	}, nil
}
