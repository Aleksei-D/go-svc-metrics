package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/router"
	"go-svc-metrics/internal/storage/local"
	"go-svc-metrics/internal/usecase"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestStatusHandler(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		code   int
		method string
	}{
		{
			name:   "positive test for gauge metric #1",
			path:   "/update/gauge/GaugeMetric/1.0",
			method: http.MethodPost,
			code:   http.StatusOK,
		},
		{
			name:   "Negative test for gauge metric - Invalid value #2",
			path:   "/update/gauge/GaugeMetric/1.a",
			code:   http.StatusBadRequest,
			method: http.MethodPost,
		},
		{
			name:   "positive test for counter metric #3",
			path:   "/update/counter/CounterMetric/4",
			code:   http.StatusOK,
			method: http.MethodPost,
		},
		{
			name:   "Negative test for counter metric - Invalid value #4",
			path:   "/update/counter/CounterMetric/xs2",
			code:   http.StatusBadRequest,
			method: http.MethodPost,
		},
		{
			name:   "Negative test for invalid type metric #5",
			path:   "/update/creator/someMetric/xs2",
			code:   http.StatusBadRequest,
			method: http.MethodPost,
		},
	}

	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	memStorage, _ := local.NewRetryWrapperLocalStorage(configServe, 3)
	metricUseCase := usecase.NewMetricUseCase(memStorage)
	defer memStorage.Close()
	ts := httptest.NewServer(router.GetMetricRouter(metricUseCase))
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, v.method, v.path)
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}
