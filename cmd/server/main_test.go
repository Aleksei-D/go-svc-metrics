package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/domain/mocks"
	"go-svc-metrics/internal/router"
	"go-svc-metrics/internal/service"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-svc-metrics/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	validCounter     = "4"
	validGauge       = "1.0"
	validCounterID   = "CounterMetric"
	validGaugeID     = "GaugeMetric"
	invalidGaugeID   = "NotExistGaugeMetric"
	invalidCounterID = "NotExistCounterMetric"
)

var (
	validCounterINT64     int64          = 4
	validGaugeFLOAT64     float64        = 1.0
	expectedGaugeMetric   models.Metrics = models.Metrics{ID: validGaugeID, MType: models.Gauge, Value: &validGaugeFLOAT64}
	expectedCounterMetric models.Metrics = models.Metrics{ID: validCounterID, MType: models.Counter, Delta: &validCounterINT64}
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body []byte) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(body))
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name                              string
		metricType, metricID, metricValue string
		code                              int
	}{
		{
			name:        "positive test update handler  for gauge metric #1",
			metricType:  models.Gauge,
			metricID:    validGaugeID,
			metricValue: validGauge,
			code:        http.StatusOK,
		},
		{
			name:        "Negative test update handler  for gauge metric - Invalid value #2",
			metricType:  models.Gauge,
			metricID:    validGaugeID,
			metricValue: "1.a",
			code:        http.StatusBadRequest,
		},
		{
			name:        "positive test update handler for counter metric #3",
			metricType:  models.Counter,
			metricID:    validCounterID,
			metricValue: validCounter,
			code:        http.StatusOK,
		},
		{
			name:        "Negative test update handler  for ounter metric - Invalid value #4",
			metricType:  models.Counter,
			metricID:    validCounterID,
			metricValue: "xs2",
			code:        http.StatusBadRequest,
		},
		{
			name:        "Negative test update handler  for invalid type metric #5",
			metricType:  "creator",
			metricID:    "someMetric",
			metricValue: "xs2",
			code:        http.StatusBadRequest,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricRepo := mocks.NewMockMetricRepo(ctrl)
	mockMetricRepo.EXPECT().UpdateMetrics(gomock.Any(), []models.Metrics{expectedGaugeMetric}).Return([]models.Metrics{expectedGaugeMetric}, nil).AnyTimes()
	mockMetricRepo.EXPECT().UpdateMetrics(gomock.Any(), []models.Metrics{expectedCounterMetric}).Return([]models.Metrics{expectedCounterMetric}, nil).AnyTimes()

	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricService := service.NewMetricService(mockMetricRepo)
	ts := httptest.NewServer(router.NewRouter(metricService, configServe))
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			path := fmt.Sprintf("/update/%s/%s/%s", v.metricType, v.metricID, v.metricValue)
			resp, _ := testRequest(t, ts, http.MethodPost, path, []byte{})
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}

func TestValueHandler(t *testing.T) {
	tests := []struct {
		name                 string
		metricType, metricID string
		code                 int
	}{
		{
			name:       "positive test value handler for gauge metric #1",
			metricType: models.Gauge,
			metricID:   validGaugeID,
			code:       http.StatusOK,
		},
		{
			name:       "Negative test value handler for gauge metric - Invalid value #2",
			metricType: models.Gauge,
			metricID:   invalidGaugeID,
			code:       http.StatusNotFound,
		},
		{
			name:       "positive test value handler for counter metric #3",
			metricType: models.Counter,
			metricID:   validCounterID,
			code:       http.StatusOK,
		},
		{
			name:       "Negative test value handler for counter metric - Invalid value #4",
			metricType: models.Counter,
			metricID:   invalidCounterID,
			code:       http.StatusNotFound,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricRepo := mocks.NewMockMetricRepo(ctrl)
	mockMetricRepo.EXPECT().GetMetric(gomock.Any(), models.Metrics{ID: validGaugeID, MType: models.Gauge}).Return(expectedGaugeMetric, nil).AnyTimes()
	mockMetricRepo.EXPECT().GetMetric(gomock.Any(), models.Metrics{ID: validCounterID, MType: models.Counter}).Return(expectedCounterMetric, nil).AnyTimes()
	mockMetricRepo.EXPECT().GetMetric(gomock.Any(), models.Metrics{ID: invalidGaugeID, MType: models.Gauge}).Return(models.Metrics{}, sql.ErrNoRows).AnyTimes()
	mockMetricRepo.EXPECT().GetMetric(gomock.Any(), models.Metrics{ID: invalidCounterID, MType: models.Counter}).Return(models.Metrics{}, sql.ErrNoRows).AnyTimes()

	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricService := service.NewMetricService(mockMetricRepo)
	ts := httptest.NewServer(router.NewRouter(metricService, configServe))
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			path := fmt.Sprintf("/value/%s/%s", v.metricType, v.metricID)

			resp, _ := testRequest(t, ts, http.MethodGet, path, []byte{})
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}

func TestAllMetricsHandler(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{
			name: "positive test all metrics handler",
			code: http.StatusOK,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricRepo := mocks.NewMockMetricRepo(ctrl)
	mockMetricRepo.EXPECT().GetAllMetrics(gomock.Any()).Return([]models.Metrics{expectedGaugeMetric, expectedCounterMetric}, nil).AnyTimes()

	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricService := service.NewMetricService(mockMetricRepo)
	ts := httptest.NewServer(router.NewRouter(metricService, configServe))
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			path := "/"

			resp, _ := testRequest(t, ts, http.MethodGet, path, []byte{})
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}

func TestUpdateModelHandler(t *testing.T) {
	tests := []struct {
		name   string
		metric models.Metrics
		code   int
	}{
		{
			name: "positive test update model handler  or gauge metric #1",
			metric: models.Metrics{
				ID:    validGaugeID,
				MType: models.Gauge,
				Value: &validGaugeFLOAT64,
			},
			code: http.StatusOK,
		},
		{
			name: "positive test update model handler for counter metric #3",
			metric: models.Metrics{
				ID:    validCounterID,
				MType: models.Counter,
				Delta: &validCounterINT64,
			},
			code: http.StatusOK,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricRepo := mocks.NewMockMetricRepo(ctrl)
	mockMetricRepo.EXPECT().UpdateMetrics(gomock.Any(), []models.Metrics{expectedGaugeMetric}).Return([]models.Metrics{expectedGaugeMetric}, nil).AnyTimes()
	mockMetricRepo.EXPECT().UpdateMetrics(gomock.Any(), []models.Metrics{expectedCounterMetric}).Return([]models.Metrics{expectedCounterMetric}, nil).AnyTimes()

	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricService := service.NewMetricService(mockMetricRepo)
	ts := httptest.NewServer(router.NewRouter(metricService, configServe))
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			metricJSON, err := json.Marshal(v.metric)
			require.NoError(t, err)

			path := "/update"
			resp, _ := testRequest(t, ts, http.MethodPost, path, metricJSON)
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}

func TestGetModelHandler(t *testing.T) {
	tests := []struct {
		name   string
		metric models.Metrics
		code   int
	}{
		{
			name: "positive test update model handler  or gauge metric #1",
			metric: models.Metrics{
				ID:    validGaugeID,
				MType: models.Gauge,
			},
			code: http.StatusOK,
		},
		{
			name: "positive test update model handler for counter metric #3",
			metric: models.Metrics{
				ID:    validCounterID,
				MType: models.Counter,
			},
			code: http.StatusOK,
		},
		{
			name: "Negative test value handler for counter metric - Invalid value #4",
			metric: models.Metrics{
				ID:    invalidCounterID,
				MType: models.Counter,
			},
			code: http.StatusNotFound,
		},
		{
			name: "Negative test value handler for gauge metric - Invalid value #2",
			metric: models.Metrics{
				ID:    invalidGaugeID,
				MType: models.Gauge,
			},
			code: http.StatusNotFound,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricRepo := mocks.NewMockMetricRepo(ctrl)
	mockMetricRepo.EXPECT().GetMetric(gomock.Any(), models.Metrics{ID: validGaugeID, MType: models.Gauge}).Return(expectedGaugeMetric, nil).AnyTimes()
	mockMetricRepo.EXPECT().GetMetric(gomock.Any(), models.Metrics{ID: validCounterID, MType: models.Counter}).Return(expectedCounterMetric, nil).AnyTimes()
	mockMetricRepo.EXPECT().GetMetric(gomock.Any(), models.Metrics{ID: invalidGaugeID, MType: models.Gauge}).Return(models.Metrics{}, sql.ErrNoRows).AnyTimes()
	mockMetricRepo.EXPECT().GetMetric(gomock.Any(), models.Metrics{ID: invalidCounterID, MType: models.Counter}).Return(models.Metrics{}, sql.ErrNoRows).AnyTimes()

	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricService := service.NewMetricService(mockMetricRepo)
	ts := httptest.NewServer(router.NewRouter(metricService, configServe))
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			metricJSON, err := json.Marshal(v.metric)
			require.NoError(t, err)

			path := "/value"
			resp, _ := testRequest(t, ts, http.MethodPost, path, metricJSON)
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}

func TestUpdateBatchMetricslHandler(t *testing.T) {
	tests := []struct {
		name    string
		metrics []models.Metrics
		code    int
	}{
		{
			name: "positive test update batch model handler",
			metrics: []models.Metrics{
				{ID: validGaugeID, MType: models.Gauge, Value: &validGaugeFLOAT64},
				{ID: validCounterID, MType: models.Counter, Delta: &validCounterINT64},
			},
			code: http.StatusOK,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricRepo := mocks.NewMockMetricRepo(ctrl)
	mockMetricRepo.EXPECT().UpdateMetrics(gomock.Any(), []models.Metrics{expectedGaugeMetric, expectedCounterMetric}).Return([]models.Metrics{expectedGaugeMetric, expectedCounterMetric}, nil).AnyTimes()

	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricService := service.NewMetricService(mockMetricRepo)
	ts := httptest.NewServer(router.NewRouter(metricService, configServe))
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			metricJSON, err := json.Marshal(v.metrics)
			require.NoError(t, err)

			path := "/updates"
			resp, _ := testRequest(t, ts, http.MethodPost, path, metricJSON)
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}
