package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	Storage    map[string]*metric.Metric
	*sync.Mutex
}

func (cs *ConfigServer) PostHandlerMetrics(w http.ResponseWriter, r *http.Request) {

	m, err := metric.ParseMetricEntityFromRequest(r)

	if err != nil {
		if err.Error() == "missmatched type" {
			log.Println(err)
			http.Error(w, fmt.Sprint(err), http.StatusNotImplemented)
			return
		} else {
			log.Println(err)
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
	}
	ID := r.URL.Query().Get("id")
	cs.Storage[ID] = m

}

// Return metric data in JSON by Requested URI
func (cs *ConfigServer) GetMetrics(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get("id")

	value, err := cs.GetMetricsByKey(context.Background(), key)
	if err != nil {
		http.Error(w, "metric not found", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err := json.NewEncoder(w).Encode(&value); err != nil {
		http.Error(w, "unable to marshal the struct", http.StatusBadRequest)
		return
	}

}

// Return metric data in JSON
func (cs *ConfigServer) GetMetricsAll(w http.ResponseWriter, r *http.Request) {

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cs.Storage); err != nil {
		http.Error(w, "unable to marshal the struct", http.StatusBadRequest)
	}

}

func (cs *ConfigServer) GetMetricsByKey(ctx context.Context, key string) (*metric.Metric, error) {
	cs.Lock()
	defer cs.Unlock()

	m, ok := cs.Storage[key]
	if !ok {
		return nil, errors.New("metric not found")
	}
	return m, nil

}
