package agent

import (
	"fmt"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/models"
	"math/rand"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

const counterMetricName = "PollCount"

type MetricUpdater struct {
	clientAgent ClientAgent
	*config.Config
	CounterMetric *int64
}

func NewMetricUpdater() (*MetricUpdater, error) {
	agentConfig, err := config.NewAgentConfig()
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
		clientAgent:   agentClient,
		Config:        agentConfig,
		CounterMetric: new(int64),
	}, nil
}

func (m *MetricUpdater) Run() error {
	errors := make(chan error)
	doneCh := make(chan struct{})
	defer close(doneCh)
	pollTicker := time.NewTicker(m.GetPollInterval())
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(m.GetReportInterval())
	defer reportTicker.Stop()

	metricsCh := m.metricGenerator(doneCh, errors, pollTicker)
	go m.sendMetrics(doneCh, metricsCh, reportTicker)

	err := <-errors
	return err
}

func (m *MetricUpdater) metricGenerator(doneCh chan struct{}, errorCh chan<- error, pollTicker *time.Ticker) <-chan []models.Metrics {
	metricSizeCh := *m.ReportInterval / *m.PollInterval + 1
	metricCh := make(chan []models.Metrics, metricSizeCh)

	go func() {
		defer close(metricCh)
		for {
			select {
			case <-doneCh:
				return
			case <-pollTicker.C:
				metrics, err := m.GetMetrics()
				if err != nil {
					errorCh <- err
				}
				metricCh <- metrics
			}
		}
	}()
	return metricCh
}

func (m *MetricUpdater) sendMetrics(doneCh chan struct{}, metricCh <-chan []models.Metrics, reportTicker *time.Ticker) {
	for {
		select {
		case <-doneCh:
			return
		case <-reportTicker.C:
			for w := 1; w <= int(*m.RateLimit); w++ {
				go m.clientAgent.MetricSenderWorker(doneCh, metricCh)
			}
		}
	}
}

func (m *MetricUpdater) GetMetrics() ([]models.Metrics, error) {
	var memStats runtime.MemStats
	metrics := make([]models.Metrics, 0)

	runtime.ReadMemStats(&memStats)
	metrics = append(metrics, m.getGaugeMetric("Alloc", float64(memStats.Alloc)))
	metrics = append(metrics, m.getGaugeMetric("BuckHashSys", float64(memStats.BuckHashSys)))
	metrics = append(metrics, m.getGaugeMetric("Frees", float64(memStats.Frees)))
	metrics = append(metrics, m.getGaugeMetric("GCCPUFraction", memStats.GCCPUFraction))
	metrics = append(metrics, m.getGaugeMetric("GCSys", float64(memStats.GCSys)))
	metrics = append(metrics, m.getGaugeMetric("HeapAlloc", float64(memStats.HeapAlloc)))
	metrics = append(metrics, m.getGaugeMetric("HeapIdle", float64(memStats.HeapIdle)))
	metrics = append(metrics, m.getGaugeMetric("HeapInuse", float64(memStats.HeapInuse)))
	metrics = append(metrics, m.getGaugeMetric("HeapObjects", float64(memStats.HeapObjects)))
	metrics = append(metrics, m.getGaugeMetric("Frees", float64(memStats.Frees)))
	metrics = append(metrics, m.getGaugeMetric("HeapReleased", float64(memStats.HeapReleased)))
	metrics = append(metrics, m.getGaugeMetric("HeapSys", float64(memStats.HeapSys)))
	metrics = append(metrics, m.getGaugeMetric("LastGC", float64(memStats.LastGC)))
	metrics = append(metrics, m.getGaugeMetric("Lookups", float64(memStats.Lookups)))
	metrics = append(metrics, m.getGaugeMetric("MCacheInuse", float64(memStats.MCacheInuse)))
	metrics = append(metrics, m.getGaugeMetric("MCacheSys", float64(memStats.MCacheSys)))
	metrics = append(metrics, m.getGaugeMetric("MSpanInuse", float64(memStats.MSpanInuse)))
	metrics = append(metrics, m.getGaugeMetric("MSpanSys", float64(memStats.MSpanSys)))
	metrics = append(metrics, m.getGaugeMetric("Mallocs", float64(memStats.Mallocs)))
	metrics = append(metrics, m.getGaugeMetric("NextGC", float64(memStats.NextGC)))
	metrics = append(metrics, m.getGaugeMetric("NumGC", float64(memStats.NumGC)))
	metrics = append(metrics, m.getGaugeMetric("NumForcedGC", float64(memStats.NumForcedGC)))
	metrics = append(metrics, m.getGaugeMetric("OtherSys", float64(memStats.OtherSys)))
	metrics = append(metrics, m.getGaugeMetric("PauseTotalNs", float64(memStats.PauseTotalNs)))
	metrics = append(metrics, m.getGaugeMetric("StackInuse", float64(memStats.StackInuse)))
	metrics = append(metrics, m.getGaugeMetric("StackSys", float64(memStats.StackSys)))
	metrics = append(metrics, m.getGaugeMetric("Sys", float64(memStats.Sys)))
	metrics = append(metrics, m.getGaugeMetric("TotalAlloc", float64(memStats.TotalAlloc)))
	metrics = append(metrics, m.getGaugeMetric("RandomValue", rand.Float64()))

	v, err := mem.VirtualMemory()
	if err != nil {
		return metrics, err
	}

	metrics = append(metrics, m.getGaugeMetric("TotalMemory", float64(v.Total)))
	metrics = append(metrics, m.getGaugeMetric("FreeMemory", float64(v.Free)))

	c, err := cpu.Percent(0, true)
	if err != nil {
		return metrics, err
	}
	for i, percent := range c {
		metrics = append(metrics, m.getGaugeMetric(fmt.Sprintf("CPUutilization%d", i+1), percent))
	}

	atomic.AddInt64(m.CounterMetric, 1)
	return append(metrics, m.getCounterMetric()), nil
}

func (m *MetricUpdater) getCounterMetric() models.Metrics {
	return models.Metrics{
		ID:    counterMetricName,
		MType: models.Counter,
		Delta: m.CounterMetric,
	}
}

func (m *MetricUpdater) getGaugeMetric(name string, value float64) models.Metrics {
	return models.Metrics{
		ID:    name,
		MType: models.Gauge,
		Value: &value,
	}
}
