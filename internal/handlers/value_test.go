package handlers_test

import (
	"database/sql"
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

func TestValueHandler(t *testing.T) {
	tests := []struct {
		name                 string
		metricType, metricID string
		code                 int
		mockExpect           func(mockRepo *mocks.MockMetricRepo)
	}{
		{
			name:       "positive test value handler for gauge metric #1",
			metricType: models.Gauge,
			metricID:   helpers.ValidGaugeID,
			code:       http.StatusOK,
			mockExpect: func(mockRepo *mocks.MockMetricRepo) {
				mockRepo.EXPECT().GetMetric(gomock.Any(), helpers.GaugeMetricRequest).Return(helpers.GaugeMetric, nil)
			},
		},
		{
			name:       "Negative test value handler for gauge metric - Invalid value #2",
			metricType: models.Gauge,
			metricID:   helpers.InvalidGaugeID,
			code:       http.StatusNotFound,
			mockExpect: func(mockRepo *mocks.MockMetricRepo) {
				mockRepo.EXPECT().GetMetric(gomock.Any(), helpers.InvalidGaugeMetricRequest).Return(models.Metrics{}, sql.ErrNoRows)
			},
		},
		{
			name:       "positive test value handler for counter metric #3",
			metricType: models.Counter,
			metricID:   helpers.ValidCounterID,
			code:       http.StatusOK,
			mockExpect: func(mockRepo *mocks.MockMetricRepo) {
				mockRepo.EXPECT().GetMetric(gomock.Any(), helpers.CounterMetricRequest).Return(helpers.CounterMetric, nil)
			},
		},
		{
			name:       "Negative test value handler for counter metric - Invalid value #4",
			metricType: models.Counter,
			metricID:   helpers.InvalidCounterID,
			code:       http.StatusNotFound,
			mockExpect: func(mockRepo *mocks.MockMetricRepo) {
				mockRepo.EXPECT().GetMetric(gomock.Any(), helpers.InvalidCounterMetricRequest).Return(models.Metrics{}, sql.ErrNoRows)
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
			path := fmt.Sprintf("/value/%s/%s", v.metricType, v.metricID)
			v.mockExpect(mockMetricRepo)
			resp, _ := helpers.TestRequest(t, ts, http.MethodGet, path, []byte{})
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}

func TestGetModelHandler(t *testing.T) {
	tests := []struct {
		name       string
		metric     models.Metrics
		code       int
		mockExpect func(mockRepo *mocks.MockMetricRepo)
	}{
		{
			name:   "positive test update model handler  or gauge metric #1",
			metric: helpers.GaugeMetricRequest,
			code:   http.StatusOK,
			mockExpect: func(mockRepo *mocks.MockMetricRepo) {
				mockRepo.EXPECT().GetMetric(gomock.Any(), helpers.GaugeMetricRequest).Return(helpers.GaugeMetric, nil)
			},
		},
		{
			name:   "positive test update model handler for counter metric #2",
			metric: helpers.CounterMetricRequest,
			code:   http.StatusOK,
			mockExpect: func(mockRepo *mocks.MockMetricRepo) {
				mockRepo.EXPECT().GetMetric(gomock.Any(), helpers.CounterMetricRequest).Return(helpers.CounterMetric, nil)
			},
		},
		{
			name:   "Negative test value handler for counter metric - Invalid value #3",
			metric: helpers.InvalidCounterMetricRequest,
			code:   http.StatusNotFound,
			mockExpect: func(mockRepo *mocks.MockMetricRepo) {
				mockRepo.EXPECT().GetMetric(gomock.Any(), helpers.InvalidCounterMetricRequest).Return(models.Metrics{}, sql.ErrNoRows)
			},
		},
		{
			name:   "Negative test value handler for gauge metric - Invalid value #4",
			metric: helpers.InvalidGaugeMetricRequest,
			code:   http.StatusNotFound,
			mockExpect: func(mockRepo *mocks.MockMetricRepo) {
				mockRepo.EXPECT().GetMetric(gomock.Any(), helpers.InvalidGaugeMetricRequest).Return(models.Metrics{}, sql.ErrNoRows)
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

			path := "/value"
			v.mockExpect(mockMetricRepo)
			resp, _ := helpers.TestRequest(t, ts, http.MethodPost, path, metricJSON)
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}
