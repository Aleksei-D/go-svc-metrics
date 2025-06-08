package agent

import (
	"bytes"
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/models"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

const pathTemplate = "http://%s/update/%s/%s/%s"

type Agent struct {
	metrics  map[string]string
	memStats runtime.MemStats
	*config.Config
}

func (a *Agent) SendReport() {
	for metricName, metricValue := range a.metrics {
		switch metricName {
		case "PollCount":
			go func() {
				err := a.SendMetric(models.Counter, metricName, metricValue)
				if err != nil {
					fmt.Println(err)
					return
				}
			}()
		default:
			go func() {
				err := a.SendMetric(models.Gauge, metricName, metricValue)
				if err != nil {
					fmt.Println(err)
					return
				}
			}()
		}
	}
}

func (a *Agent) UpdateCounterMetric() {
	value, ok := a.metrics["PollCount"]
	if !ok {
		a.metrics["PollCount"] = "1"
	} else {
		valueInt, _ := strconv.ParseInt(value, 10, 64)
		a.metrics["PollCount"] = strconv.FormatInt(valueInt+1, 10)
	}
}

func (a *Agent) UpdateGaugeMetric() {
	a.metrics["Alloc"] = strconv.FormatUint(a.memStats.Alloc, 10)
	a.metrics["BuckHashSys"] = strconv.FormatUint(a.memStats.BuckHashSys, 10)
	a.metrics["Frees"] = strconv.FormatUint(a.memStats.Frees, 10)
	a.metrics["GCCPUFraction"] = strconv.FormatFloat(a.memStats.GCCPUFraction, 'f', -1, 64)
	a.metrics["GCSys"] = strconv.FormatUint(a.memStats.GCSys, 10)
	a.metrics["HeapAlloc"] = strconv.FormatUint(a.memStats.HeapAlloc, 10)
	a.metrics["HeapIdle"] = strconv.FormatUint(a.memStats.HeapIdle, 10)
	a.metrics["HeapInuse"] = strconv.FormatUint(a.memStats.HeapInuse, 10)
	a.metrics["HeapObjects"] = strconv.FormatUint(a.memStats.HeapObjects, 10)
	a.metrics["HeapReleased"] = strconv.FormatUint(a.memStats.HeapReleased, 10)
	a.metrics["HeapSys"] = strconv.FormatUint(a.memStats.HeapSys, 10)
	a.metrics["LastGC"] = strconv.FormatUint(a.memStats.LastGC, 10)
	a.metrics["Lookups"] = strconv.FormatUint(a.memStats.Lookups, 10)
	a.metrics["MCacheInuse"] = strconv.FormatUint(a.memStats.MCacheInuse, 10)
	a.metrics["MCacheSys"] = strconv.FormatUint(a.memStats.MCacheSys, 10)
	a.metrics["MSpanInuse"] = strconv.FormatUint(a.memStats.MSpanInuse, 10)
	a.metrics["MCacheInuse"] = strconv.FormatUint(a.memStats.MCacheInuse, 10)
	a.metrics["MSpanSys"] = strconv.FormatUint(a.memStats.MSpanSys, 10)
	a.metrics["Mallocs"] = strconv.FormatUint(a.memStats.Mallocs, 10)
	a.metrics["NextGC"] = strconv.FormatUint(a.memStats.NextGC, 10)
	a.metrics["NumGC"] = strconv.Itoa(int(a.memStats.NumGC))
	a.metrics["NumForcedGC"] = strconv.Itoa(int(a.memStats.NumForcedGC))
	a.metrics["OtherSys"] = strconv.FormatUint(a.memStats.OtherSys, 10)
	a.metrics["PauseTotalNs"] = strconv.FormatUint(a.memStats.PauseTotalNs, 10)
	a.metrics["StackInuse"] = strconv.FormatUint(a.memStats.StackInuse, 10)
	a.metrics["StackSys"] = strconv.FormatUint(a.memStats.StackSys, 10)
	a.metrics["Sys"] = strconv.FormatUint(a.memStats.Sys, 10)
	a.metrics["TotalAlloc"] = strconv.FormatUint(a.memStats.TotalAlloc, 10)
}

func (a *Agent) SendMetric(metricType, metricName, value string) error {
	response, err := http.Post(fmt.Sprintf(pathTemplate, a.GetServeAddress(), metricType, metricName, value), "text/plain", bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func GetNewAgent() *Agent {
	agentConfig := config.GetAgentConfig()
	return &Agent{
		metrics: make(map[string]string),
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
			a.UpdateCounterMetric()
		case <-metricReportTicker.C:
			a.SendReport()
		}
	}
}
