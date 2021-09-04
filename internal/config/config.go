package config

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	metric "github.com/goethesum/-go-musthave-devops-tpl/internal/metrics"
)

type ConfigAgent struct {
	Server        string
	URLMetricPush string
	TimeInterval  time.Duration
}

type ConfigServer struct {
	PortNumber string
	Storage    map[string]metric.Metric
	*sync.Mutex
}

func (cs *ConfigServer) PostHandlerMetrics(w http.ResponseWriter, r *http.Request) {

	switch {
	case r.URL.Query().Get("id") != "":
		cs.Storage[r.URL.Query().Get("id")] = metric.Metric{ID: r.URL.Query().Get("id")}
	case r.URL.Query().Get("type") == "gauge":
		cs.Storage[r.URL.Query().Get("type")] = metric.Metric{Type: metric.MetricType(r.URL.Query().Get("id"))}
	case r.URL.Query().Get("type") == "conter":
		cs.Storage[r.URL.Query().Get("type")] = metric.Metric{Type: metric.MetricType(r.URL.Query().Get("id"))}
	case r.URL.Query().Get("value") != "":
		cs.Storage[r.URL.Query().Get("value")] = metric.Metric{ID: r.URL.Query().Get("value")}
	default:
		http.Error(w, "There are some empty keys", http.StatusBadRequest)
	}

}

// Return metric data in JSON by Requested URI
func (cs *ConfigServer) GetMetrics(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get("id")

	value, err := cs.GetMetricsByKey(context.Background(), key)
	if err != nil {
		http.Error(w, "metric not found", http.StatusBadRequest)
	}

	jsonMetric, err := json.Marshal(value)
	if err != nil {
		http.Error(w, "unable to marshal the struct", http.StatusBadRequest)
	}
	w.Write(jsonMetric)

}

// Return metric data in JSON
func (cs *ConfigServer) GetMetricsAll(w http.ResponseWriter, r *http.Request) {

	jsonMetricAll, err := json.MarshalIndent(cs.Storage, "", "    ")
	if err != nil {
		http.Error(w, "unable to marshal the struct", http.StatusBadRequest)
	}
	w.Write(jsonMetricAll)

}

func (cs *ConfigServer) GetMetricsByKey(ctx context.Context, key string) (metric.Metric, error) {
	cs.Lock()
	defer cs.Unlock()

	m, ok := cs.Storage[key]
	if !ok {
		return metric.Metric{}, errors.New("metric not found")
	}
	return m, nil

}
