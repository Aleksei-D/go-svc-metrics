package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/logger"
	"go-svc-metrics/internal/utils/crypto"
	"go-svc-metrics/internal/utils/delay"
	"go-svc-metrics/models"
	"net/http"
	"time"
)

const (
	updatePath              = "http://%s/update/"
	updatesBatchMetricsPath = "http://%s/updates/"
)

type retryRoundTripper struct {
	next       http.RoundTripper
	maxRetries uint
}

func (rr retryRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	var res *http.Response
	var err error
	delay := delay.NewDelay()
	for attempts := 0; attempts < int(rr.maxRetries); attempts++ {
		res, err = rr.next.RoundTrip(r)
		if err == nil && res.StatusCode < http.StatusInternalServerError {
			break
		}

		select {
		case <-r.Context().Done():
			return res, r.Context().Err()
		case <-time.After(delay()):
		}
	}

	return res, err
}

type ClientAgent struct {
	httpClient *http.Client
	config     *config.Config
}

func (c *ClientAgent) MetricSenderWorker(doneCh chan struct{}, metricCh <-chan []models.Metrics) {
	select {
	case <-doneCh:
		return
	case metrics := <-metricCh:
		err := c.SendBatchMetrics(metrics)
		if err != nil {
			logger.Log.Warn(err.Error())
		}
	}
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

	request, err := c.getRequest(updatePath, gzipData)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "text/html")
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

func (c *ClientAgent) getRequest(url string, data []byte) (*http.Request, error) {
	if c.config.Key != nil {
		hash := crypto.GetHash(*c.config.Key, data)
		cryptData, err := crypto.EncryptData(*c.config.Key, data)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(cryptData))
		if err != nil {
			return nil, err
		}
		request.Header.Set("HashSHA256", hex.EncodeToString(hash))
		return request, nil
	}
	return http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
}
