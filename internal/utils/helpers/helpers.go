package helpers

import (
	"bytes"
	"go-svc-metrics/models"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	ValidCounterSTRING = "4"
	ValidGaugeSTRING   = "1.0"
	InvalidValueSTRING = "1.a"
	ValidCounterID     = "CounterMetric"
	ValidGaugeID       = "GaugeMetric"
	InvalidGaugeID     = "NotExistGaugeMetric"
	InvalidCounterID   = "NotExistCounterMetric"
)

var (
	ValidCounterValue           int64          = 4
	ValidGaugeValue             float64        = 1.0
	GaugeMetric                 models.Metrics = models.Metrics{ID: ValidGaugeID, MType: models.Gauge, Value: &ValidGaugeValue}
	CounterMetric               models.Metrics = models.Metrics{ID: ValidCounterID, MType: models.Counter, Delta: &ValidCounterValue}
	GaugeMetricRequest          models.Metrics = models.Metrics{ID: ValidGaugeID, MType: models.Gauge}
	CounterMetricRequest        models.Metrics = models.Metrics{ID: ValidCounterID, MType: models.Counter}
	InvalidGaugeMetricRequest   models.Metrics = models.Metrics{ID: InvalidGaugeID, MType: models.Gauge}
	InvalidCounterMetricRequest models.Metrics = models.Metrics{ID: InvalidCounterID, MType: models.Counter}
)

func TestRequest(t *testing.T, ts *httptest.Server, method, path string, body []byte) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(body))
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
