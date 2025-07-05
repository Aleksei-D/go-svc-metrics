package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/models"
	"net/http"
)

const (
	updatePath              = "http://%s/update/"
	updatesBatchMetricsPath = "http://%s/updates/"
)

type ClientAgent struct {
	httpClient *http.Client
	config     *config.Config
}

func (c *ClientAgent) SendOneMetric(metric models.Metrics) error {
	metricJSON, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	response, err := c.sendMetric(metricJSON, c.getUpdatePath())
	if err != nil {
		return err
	}

	defer response.Body.Close()
	return nil
}

func (c *ClientAgent) SendBatchMetrics(metrics []models.Metrics) error {
	metricJSON, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	response, err := c.sendMetric(metricJSON, c.getUpdateBatchPath())
	if err != nil {
		return err
	}

	defer response.Body.Close()
	return nil
}

func (c *ClientAgent) sendMetric(metricData []byte, updatePath string) (*http.Response, error) {
	gzipData, err := Compress(metricData)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, updatePath, bytes.NewBuffer(gzipData))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "gzip")
	return c.httpClient.Do(request)
}

func (c *ClientAgent) getUpdatePath() string {
	return fmt.Sprintf(updatePath, c.config.GetServeAddress())
}

func (c *ClientAgent) getUpdateBatchPath() string {
	return fmt.Sprintf(updatesBatchMetricsPath, c.config.GetServeAddress())
}

func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(data)
	if err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
