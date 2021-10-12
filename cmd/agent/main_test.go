package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	metric "github.com/goethesum/-go-musthave-devops-tpl/internal/metrics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MetricMock struct {
	mock.Mock
}

func (m *MetricMock) MetricSend(ctx context.Context, metrics metric.Metric) error {
	return nil
}

func TestMetricSend(t *testing.T) {
	testClient := &clientHTTP{
		client: *resty.New(),
	}
	testMetric := &metric.Metric{
		ID:    "test",
		MType: metric.MetricTypeCounter,
		Value: 4321,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, http.StatusOK)
	}))
	defer ts.Close()

	resp, err := testClient.MetricSend(ts.URL, *testMetric, &http.Transport{})
	if err != nil {
		t.Errorf("throw an error during the test %s", err)
	}
	assert.Equal(t, 200, resp.StatusCode())
}
