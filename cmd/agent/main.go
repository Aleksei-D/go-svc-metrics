package main

import (
	"bytes"
	"fmt"
	"go-svc-metrics/internal/storage"
	"net/http"
	"runtime"
	"time"
)

const pollInterval = 2
const reportInterval = 10
const pathGaugeTemplate = "http://localhost:8080/update/gauge/%s/%f"
const pathCounterTemplate = "http://localhost:8080/update/counter/%s/%d"

type Metrics struct {
	*storage.MemStorage
}

func (m *Metrics) SendReport() {
	for metricName, metricValue := range m.GaugeMetrics {
		go func() {
			err := SendGaugeMetric(metricName, metricValue)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
	for metricName, metricValue := range m.CounterMetrics {
		go func() {
			err := SendCounterMetric(metricName, metricValue)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
}

func (m *Metrics) UpdateCounterMetric() {
	_, ok := m.CounterMetrics["PollCount"]
	if !ok {
		m.CounterMetrics["PollCount"] = 1
	} else {
		m.CounterMetrics["PollCount"] += 1
	}
}

func (m *Metrics) UpdateGaugeMetric(memStats *runtime.MemStats) {
	m.GaugeMetrics["Alloc"] = float64(memStats.Alloc)
	m.GaugeMetrics["BuckHashSys"] = float64(memStats.BuckHashSys)
	m.GaugeMetrics["Frees"] = float64(memStats.Frees)
	m.GaugeMetrics["GCCPUFraction"] = memStats.GCCPUFraction
	m.GaugeMetrics["GCSys"] = float64(memStats.GCSys)
	m.GaugeMetrics["HeapAlloc"] = float64(memStats.HeapAlloc)
	m.GaugeMetrics["HeapIdle"] = float64(memStats.HeapIdle)
	m.GaugeMetrics["HeapInuse"] = float64(memStats.HeapInuse)
	m.GaugeMetrics["HeapObjects"] = float64(memStats.HeapObjects)
	m.GaugeMetrics["HeapReleased"] = float64(memStats.HeapReleased)
	m.GaugeMetrics["HeapSys"] = float64(memStats.HeapSys)
	m.GaugeMetrics["LastGC"] = float64(memStats.LastGC)
	m.GaugeMetrics["Lookups"] = float64(memStats.Lookups)
	m.GaugeMetrics["MCacheInuse"] = float64(memStats.MCacheInuse)
	m.GaugeMetrics["MCacheSys"] = float64(memStats.MCacheSys)
	m.GaugeMetrics["MSpanInuse"] = float64(memStats.MSpanInuse)
	m.GaugeMetrics["MCacheInuse"] = float64(memStats.MCacheInuse)
	m.GaugeMetrics["MSpanSys"] = float64(memStats.MSpanSys)
	m.GaugeMetrics["Mallocs"] = float64(memStats.Mallocs)
	m.GaugeMetrics["NextGC"] = float64(memStats.NextGC)
	m.GaugeMetrics["NumGC"] = float64(memStats.NumGC)
	m.GaugeMetrics["NumForcedGC"] = float64(memStats.NumForcedGC)
	m.GaugeMetrics["OtherSys"] = float64(memStats.OtherSys)
	m.GaugeMetrics["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	m.GaugeMetrics["StackInuse"] = float64(memStats.StackInuse)
	m.GaugeMetrics["StackSys"] = float64(memStats.StackSys)
	m.GaugeMetrics["Sys"] = float64(memStats.Sys)
	m.GaugeMetrics["TotalAlloc"] = float64(memStats.TotalAlloc)
}

func SendGaugeMetric(metricName string, value float64) error {
	response, err := http.Post(fmt.Sprintf(pathGaugeTemplate, metricName, value), "text/plain", bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func SendCounterMetric(metricName string, value int64) error {
	response, err := http.Post(fmt.Sprintf(pathCounterTemplate, metricName, value), "text/plain", bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func main() {
	memStats := &runtime.MemStats{}
	metricStorage := &Metrics{MemStorage: storage.InitMemStorage()}
	metricsPollTicker := time.NewTicker(pollInterval * time.Second)
	metricReportTicker := time.NewTicker(reportInterval * time.Second)

	for {
		select {
		case <-metricsPollTicker.C:
			runtime.ReadMemStats(memStats)
			metricStorage.UpdateGaugeMetric(memStats)
			metricStorage.UpdateCounterMetric()
		case <-metricReportTicker.C:
			metricStorage.SendReport()
		}
	}
}
