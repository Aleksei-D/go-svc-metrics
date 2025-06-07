package agent

import (
	"bytes"
	"fmt"
	"go-svc-metrics/internal/storage"
	"go-svc-metrics/models"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

const pollInterval = 2
const reportInterval = 10
const pathTemplate = "http://localhost:8080/update/%s/%s/%s"

type Metrics struct {
	*storage.MemStorage
}

func (m *Metrics) SendReport() {
	for metricName, metricValue := range m.Metrics {
		switch metricName {
		case "PollCount":
			go func() {
				err := SendMetric(models.Counter, metricName, metricValue)
				if err != nil {
					fmt.Println(err)
					return
				}
			}()
		default:
			go func() {
				err := SendMetric(models.Gauge, metricName, metricValue)
				if err != nil {
					fmt.Println(err)
					return
				}
			}()
		}
	}
}

func (m *Metrics) UpdateCounterMetric() {
	value, ok := m.Metrics["PollCount"]
	if !ok {
		m.Metrics["PollCount"] = "1"
	} else {
		valueInt, _ := strconv.ParseInt(value, 10, 64)
		m.Metrics["PollCount"] = strconv.FormatInt(valueInt+1, 10)
	}
}

func (m *Metrics) UpdateGaugeMetric(memStats *runtime.MemStats) {
	m.Metrics["Alloc"] = strconv.FormatUint(memStats.Alloc, 10)
	m.Metrics["BuckHashSys"] = strconv.FormatUint(memStats.BuckHashSys, 10)
	m.Metrics["Frees"] = strconv.FormatUint(memStats.Frees, 10)
	m.Metrics["GCCPUFraction"] = strconv.FormatFloat(memStats.GCCPUFraction, 'f', -1, 64)
	m.Metrics["GCSys"] = strconv.FormatUint(memStats.GCSys, 10)
	m.Metrics["HeapAlloc"] = strconv.FormatUint(memStats.HeapAlloc, 10)
	m.Metrics["HeapIdle"] = strconv.FormatUint(memStats.HeapIdle, 10)
	m.Metrics["HeapInuse"] = strconv.FormatUint(memStats.HeapInuse, 10)
	m.Metrics["HeapObjects"] = strconv.FormatUint(memStats.HeapObjects, 10)
	m.Metrics["HeapReleased"] = strconv.FormatUint(memStats.HeapReleased, 10)
	m.Metrics["HeapSys"] = strconv.FormatUint(memStats.HeapSys, 10)
	m.Metrics["LastGC"] = strconv.FormatUint(memStats.LastGC, 10)
	m.Metrics["Lookups"] = strconv.FormatUint(memStats.Lookups, 10)
	m.Metrics["MCacheInuse"] = strconv.FormatUint(memStats.MCacheInuse, 10)
	m.Metrics["MCacheSys"] = strconv.FormatUint(memStats.MCacheSys, 10)
	m.Metrics["MSpanInuse"] = strconv.FormatUint(memStats.MSpanInuse, 10)
	m.Metrics["MCacheInuse"] = strconv.FormatUint(memStats.MCacheInuse, 10)
	m.Metrics["MSpanSys"] = strconv.FormatUint(memStats.MSpanSys, 10)
	m.Metrics["Mallocs"] = strconv.FormatUint(memStats.Mallocs, 10)
	m.Metrics["NextGC"] = strconv.FormatUint(memStats.NextGC, 10)
	m.Metrics["NumGC"] = strconv.Itoa(int(memStats.NumGC))
	m.Metrics["NumForcedGC"] = strconv.Itoa(int(memStats.NumForcedGC))
	m.Metrics["OtherSys"] = strconv.FormatUint(memStats.OtherSys, 10)
	m.Metrics["PauseTotalNs"] = strconv.FormatUint(memStats.PauseTotalNs, 10)
	m.Metrics["StackInuse"] = strconv.FormatUint(memStats.StackInuse, 10)
	m.Metrics["StackSys"] = strconv.FormatUint(memStats.StackSys, 10)
	m.Metrics["Sys"] = strconv.FormatUint(memStats.Sys, 10)
	m.Metrics["TotalAlloc"] = strconv.FormatUint(memStats.TotalAlloc, 10)
}

func SendMetric(metricType, metricName, value string) error {
	response, err := http.Post(fmt.Sprintf(pathTemplate, metricType, metricName, value), "text/plain", bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func MetricProcessing() {
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
