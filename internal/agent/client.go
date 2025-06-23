package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/models"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

const (
	updatePath = "http://%s/update/"
)

type Agent struct {
	metrics   map[string]float64
	memStats  runtime.MemStats
	PoolCount int64
	*config.Config
}

func (a *Agent) SendReport() {
	for id, metricValue := range a.metrics {
		metric := &models.Metrics{
			ID:    id,
			MType: models.Gauge,
			Value: &metricValue,
		}
		go func() {
			err := a.SendMetric(metric)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
	go func() {
		err := a.SendMetric(&models.Metrics{
			ID:    "PollCount",
			MType: models.Counter,
			Delta: &a.PoolCount,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
}

func (a *Agent) UpdateGaugeMetric() {
	a.metrics["Alloc"] = float64(a.memStats.Alloc)
	a.metrics["BuckHashSys"] = float64(a.memStats.BuckHashSys)
	a.metrics["Frees"] = float64(a.memStats.Frees)
	a.metrics["GCCPUFraction"] = a.memStats.GCCPUFraction
	a.metrics["GCSys"] = float64(a.memStats.GCSys)
	a.metrics["HeapAlloc"] = float64(a.memStats.HeapAlloc)
	a.metrics["HeapIdle"] = float64(a.memStats.HeapIdle)
	a.metrics["HeapInuse"] = float64(a.memStats.HeapInuse)
	a.metrics["HeapObjects"] = float64(a.memStats.HeapObjects)
	a.metrics["HeapReleased"] = float64(a.memStats.HeapReleased)
	a.metrics["HeapSys"] = float64(a.memStats.HeapSys)
	a.metrics["LastGC"] = float64(a.memStats.LastGC)
	a.metrics["Lookups"] = float64(a.memStats.Lookups)
	a.metrics["MCacheInuse"] = float64(a.memStats.MCacheInuse)
	a.metrics["MCacheSys"] = float64(a.memStats.MCacheSys)
	a.metrics["MSpanInuse"] = float64(a.memStats.MSpanInuse)
	a.metrics["MSpanSys"] = float64(a.memStats.MSpanSys)
	a.metrics["Mallocs"] = float64(a.memStats.Mallocs)
	a.metrics["NextGC"] = float64(a.memStats.NextGC)
	a.metrics["NumGC"] = float64(a.memStats.NumGC)
	a.metrics["NumForcedGC"] = float64(a.memStats.NumForcedGC)
	a.metrics["OtherSys"] = float64(a.memStats.OtherSys)
	a.metrics["PauseTotalNs"] = float64(a.memStats.PauseTotalNs)
	a.metrics["StackInuse"] = float64(a.memStats.StackInuse)
	a.metrics["StackSys"] = float64(a.memStats.StackSys)
	a.metrics["Sys"] = float64(a.memStats.Sys)
	a.metrics["TotalAlloc"] = float64(a.memStats.TotalAlloc)
	a.metrics["RandomValue"] = rand.Float64()
}

func (a *Agent) SendMetric(metric *models.Metrics) error {
	metricJSON, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	response, err := http.Post(fmt.Sprintf(updatePath, a.GetServeAddress()), "application/json", bytes.NewBuffer(metricJSON))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func GetNewAgent() *Agent {
	agentConfig := config.GetAgentConfig()
	return &Agent{
		metrics: make(map[string]float64),
		Config:  agentConfig,
	}
}

func (a *Agent) MetricProcessing() {
	metricsPollTicker := time.NewTicker(a.GetPollInterval())
	metricReportTicker := time.NewTicker(a.GetReportInterval())
	for {
		select {
		case <-metricsPollTicker.C:
			runtime.ReadMemStats(&a.memStats)
			a.UpdateGaugeMetric()
			a.PoolCount++
		case <-metricReportTicker.C:
			a.SendReport()
		}
	}
}
