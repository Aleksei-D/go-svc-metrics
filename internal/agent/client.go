package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"go-svc-metrics/models"
	"net/http"
)

const (
	updatePath = "http://%s/update/"
)

type ClientAgent struct {
	httpClient *http.Client
	updatePath string
}

func (c *ClientAgent) SendMetric(metric *models.Metrics) error {
	metricJSON, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	gzipData, err := Compress(metricJSON)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, c.updatePath, bytes.NewBuffer(gzipData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "gzip")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
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
