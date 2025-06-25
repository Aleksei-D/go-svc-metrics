package agent

import (
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/models"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type MetricUpdater struct {
	metrics     map[string]float64
	memStats    runtime.MemStats
	PoolCount   int64
	clientAgent ClientAgent
	*config.Config
}

func (m *MetricUpdater) SendReport() {
	var wg sync.WaitGroup
	workerPoolSize := 10

	dataCh := make(chan *models.Metrics, workerPoolSize)
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for metric := range dataCh {
				err := m.clientAgent.SendMetric(metric)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}()
	}

	for id, metricValue := range m.metrics {
		metric := &models.Metrics{
			ID:    id,
			MType: models.Gauge,
			Value: &metricValue,
		}

		dataCh <- metric
	}
	dataCh <- &models.Metrics{
		ID:    "PollCount",
		MType: models.Counter,
		Delta: &m.PoolCount,
	}

	close(dataCh)
	wg.Wait()
}

func (m *MetricUpdater) UpdateGaugeMetric() {
	m.metrics["Alloc"] = float64(m.memStats.Alloc)
	m.metrics["BuckHashSys"] = float64(m.memStats.BuckHashSys)
	m.metrics["Frees"] = float64(m.memStats.Frees)
	m.metrics["GCCPUFraction"] = m.memStats.GCCPUFraction
	m.metrics["GCSys"] = float64(m.memStats.GCSys)
	m.metrics["HeapAlloc"] = float64(m.memStats.HeapAlloc)
	m.metrics["HeapIdle"] = float64(m.memStats.HeapIdle)
	m.metrics["HeapInuse"] = float64(m.memStats.HeapInuse)
	m.metrics["HeapObjects"] = float64(m.memStats.HeapObjects)
	m.metrics["HeapReleased"] = float64(m.memStats.HeapReleased)
	m.metrics["HeapSys"] = float64(m.memStats.HeapSys)
	m.metrics["LastGC"] = float64(m.memStats.LastGC)
	m.metrics["Lookups"] = float64(m.memStats.Lookups)
	m.metrics["MCacheInuse"] = float64(m.memStats.MCacheInuse)
	m.metrics["MCacheSys"] = float64(m.memStats.MCacheSys)
	m.metrics["MSpanInuse"] = float64(m.memStats.MSpanInuse)
	m.metrics["MSpanSys"] = float64(m.memStats.MSpanSys)
	m.metrics["Mallocs"] = float64(m.memStats.Mallocs)
	m.metrics["NextGC"] = float64(m.memStats.NextGC)
	m.metrics["NumGC"] = float64(m.memStats.NumGC)
	m.metrics["NumForcedGC"] = float64(m.memStats.NumForcedGC)
	m.metrics["OtherSys"] = float64(m.memStats.OtherSys)
	m.metrics["PauseTotalNs"] = float64(m.memStats.PauseTotalNs)
	m.metrics["StackInuse"] = float64(m.memStats.StackInuse)
	m.metrics["StackSys"] = float64(m.memStats.StackSys)
	m.metrics["Sys"] = float64(m.memStats.Sys)
	m.metrics["TotalAlloc"] = float64(m.memStats.TotalAlloc)
	m.metrics["RandomValue"] = rand.Float64()
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
			m.UpdateGaugeMetric()
			m.PoolCount++
		case <-metricReportTicker.C:
			m.SendReport()
		}
	}
}

func GetNewMetricUpdater() (*MetricUpdater, error) {
	agentConfig, err := config.GetAgentConfig()
	if err != nil {
		return nil, err
	}

	agentClient := ClientAgent{
		updatePath: fmt.Sprintf(updatePath, agentConfig.GetServeAddress()),
		httpClient: &http.Client{},
	}
	return &MetricUpdater{
		metrics:     make(map[string]float64),
		clientAgent: agentClient,
		Config:      agentConfig,
	}, nil
}
