package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-svc-metrics/internal/handlers"
	"go-svc-metrics/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const pathTemplate = "/update/%s/%s/%s"

func TestStatusHandler(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name        string
		metricType  string
		metricName  string
		metricValue string
		want        want
	}{
		{
			name:        "positive test for gauge metric #1",
			metricType:  "gauge",
			metricName:  "someMetric",
			metricValue: "1.0",
			want: want{
				code: 200,
			},
		},
		{
			name:        "Negative test for gauge metric - Invalid value #2",
			metricType:  "gauge",
			metricName:  "someMetric",
			metricValue: "1.a",
			want: want{
				code: 400,
			},
		},
		{
			name:        "positive test for counter metric #3",
			metricType:  "counter",
			metricName:  "someMetric",
			metricValue: "1",
			want: want{
				code: 200,
			},
		},
		{
			name:        "Negative test for counter metric - Invalid value #4",
			metricType:  "counter",
			metricName:  "someMetric",
			metricValue: "xs2",
			want: want{
				code: 400,
			},
		},
		{
			name:        "Negative test for invalid type metric #5",
			metricType:  "creator",
			metricName:  "someMetric",
			metricValue: "xs2",
			want: want{
				code: 400,
			},
		},
	}
	metricHandler := handlers.MetricHandler{Storage: storage.InitMemStorage()}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := fmt.Sprintf(pathTemplate, test.metricType, test.metricName, test.metricValue)
			request := httptest.NewRequest(http.MethodPost, path, nil)
			request.SetPathValue(handlers.MetricTypePath, test.metricType)
			request.SetPathValue(handlers.MetricNamePath, test.metricName)
			request.SetPathValue(handlers.MetricValuePath, test.metricValue)

			w := httptest.NewRecorder()
			metricHandler.Serve(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}
