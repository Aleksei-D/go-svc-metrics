package handlers_test

import (
	"encoding/json"
	"fmt"
	"go-svc-metrics/internal/domain/mocks"
	"go-svc-metrics/internal/utils/helpers"
	"go-svc-metrics/models"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateBatchMetricslHandler(t *testing.T) {
	tests := []struct {
		name    string
		metrics []models.Metrics
		code    int
	}{
		{
			name:    "positive test update batch model handler",
			metrics: []models.Metrics{helpers.GaugeMetric, helpers.CounterMetric},
			code:    http.StatusOK,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricRepo := mocks.NewMockMetricRepo(ctrl)
	mockMetricRepo.EXPECT().UpdateMetrics(gomock.Any(), []models.Metrics{helpers.GaugeMetric, helpers.CounterMetric}).Return([]models.Metrics{helpers.GaugeMetric, helpers.CounterMetric}, nil).AnyTimes()

	ts := NewTestServer(mockMetricRepo)
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			metricJSON, err := json.Marshal(v.metrics)
			require.NoError(t, err)

			path := "/updates"
			resp, _ := helpers.TestRequest(t, ts, http.MethodPost, path, metricJSON)
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}

func TestUpdateModelHandler(t *testing.T) {
	tests := []struct {
		name       string
		metric     models.Metrics
		code       int
		mockExpect func(mockRepo *mocks.MockMetricRepo)
	}{
		{
			name:   "positive test update model handler  or gauge metric #1",
			metric: helpers.GaugeMetric,
			code:   http.StatusOK,
			mockExpect: func(mockRepo *mocks.MockMetricRepo) {
				mockRepo.EXPECT().UpdateMetrics(gomock.Any(), []models.Metrics{helpers.GaugeMetric}).Return([]models.Metrics{helpers.GaugeMetric}, nil)
			},
		},
		{
			name:   "positive test update model handler for counter metric #3",
			metric: helpers.CounterMetric,
			code:   http.StatusOK,
			mockExpect: func(mockRepo *mocks.MockMetricRepo) {
				mockRepo.EXPECT().UpdateMetrics(gomock.Any(), []models.Metrics{helpers.CounterMetric}).Return([]models.Metrics{helpers.CounterMetric}, nil)
			},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricRepo := mocks.NewMockMetricRepo(ctrl)

	ts := NewTestServer(mockMetricRepo)
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			metricJSON, err := json.Marshal(v.metric)
			require.NoError(t, err)

			path := "/update"
			v.mockExpect(mockMetricRepo)
			resp, _ := helpers.TestRequest(t, ts, http.MethodPost, path, metricJSON)
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name                              string
		metricType, metricID, metricValue string
		code                              int
		mockExpect                        func(mockRepo *mocks.MockMetricRepo)
	}{
		{
			name:        "positive test update handler  for gauge metric #1",
			metricType:  models.Gauge,
			metricID:    helpers.ValidGaugeID,
			metricValue: helpers.ValidGaugeSTRING,
			code:        http.StatusOK,
			mockExpect: func(mockRepo *mocks.MockMetricRepo) {
				mockRepo.EXPECT().UpdateMetrics(gomock.Any(), []models.Metrics{helpers.GaugeMetric}).Return([]models.Metrics{helpers.GaugeMetric}, nil)
			},
		},
		{
			name:        "Negative test update handler  for gauge metric - Invalid value #2",
			metricType:  models.Gauge,
			metricID:    helpers.ValidGaugeID,
			metricValue: helpers.InvalidValueSTRING,
			code:        http.StatusBadRequest,
			mockExpect:  func(mockRepo *mocks.MockMetricRepo) {},
		},
		{
			name:        "positive test update handler for counter metric #3",
			metricType:  models.Counter,
			metricID:    helpers.ValidCounterID,
			metricValue: helpers.ValidCounterSTRING,
			code:        http.StatusOK,
			mockExpect: func(mockRepo *mocks.MockMetricRepo) {
				mockRepo.EXPECT().UpdateMetrics(gomock.Any(), []models.Metrics{helpers.CounterMetric}).Return([]models.Metrics{helpers.CounterMetric}, nil)
			},
		},
		{
			name:        "Negative test update handler  for ounter metric - Invalid value #4",
			metricType:  models.Counter,
			metricID:    helpers.ValidCounterID,
			metricValue: helpers.InvalidValueSTRING,
			code:        http.StatusBadRequest,
			mockExpect:  func(mockRepo *mocks.MockMetricRepo) {},
		},
		{
			name:        "Negative test update handler  for invalid type metric #5",
			metricType:  "creator",
			metricID:    "someMetric",
			metricValue: "xs2",
			code:        http.StatusBadRequest,
			mockExpect:  func(mockRepo *mocks.MockMetricRepo) {},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricRepo := mocks.NewMockMetricRepo(ctrl)
	ts := NewTestServer(mockMetricRepo)
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			path := fmt.Sprintf("/update/%s/%s/%s", v.metricType, v.metricID, v.metricValue)
			v.mockExpect(mockMetricRepo)
			resp, _ := helpers.TestRequest(t, ts, http.MethodPost, path, []byte{})
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}
