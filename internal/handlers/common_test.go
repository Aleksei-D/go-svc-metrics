package handlers_test

import (
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/domain"
	"go-svc-metrics/internal/domain/mocks"
	"go-svc-metrics/internal/router"
	"go-svc-metrics/internal/service"
	"go-svc-metrics/internal/utils/helpers"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-svc-metrics/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func NewTestServer(repo domain.MetricRepo) *httptest.Server {
	_ = config.InitDefaultEnv()
	configServe, _ := config.InitConfig()
	metricService := service.NewMetricService(repo)
	return httptest.NewServer(router.NewRouter(metricService, configServe, nil))
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
	mockMetricRepo.EXPECT().GetAllMetrics(gomock.Any()).Return([]models.Metrics{helpers.GaugeMetric, helpers.CounterMetric}, nil).AnyTimes()

	ts := NewTestServer(mockMetricRepo)
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			path := "/"

			resp, _ := helpers.TestRequest(t, ts, http.MethodGet, path, []byte{})
			defer resp.Body.Close()
			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}
