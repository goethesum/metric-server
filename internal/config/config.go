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

	"github.com/go-chi/chi/v5"

	"github.com/goethesum/-go-musthave-devops-tpl/internal/history"
	metric "github.com/goethesum/-go-musthave-devops-tpl/internal/metrics"
)

type ConfigAgent struct {
	Server         string        `env:"ADDRESS" envDefault:"http://localhost:8080"`
	URLMetricPush  string        `env:"URL_PATH" envDefault:"/update"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
}

type ConfigServer struct {
	PortNumber    string `env:"ADDRESS" envDefault:"0.0.0.0:8080"`
	Storage       map[string]metric.Metric
	FileStorage   string        `env:"FILE_STORAGE_PATH" envDefault:"./history"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
	*sync.Mutex
}

func (cs *ConfigServer) PostHandlerMetricsJSON(w http.ResponseWriter, r *http.Request) {
	m := metric.Metric{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		log.Printf("unable to decode params in PostHandlerMetricsJSON, %s", err)
		return
	}
	cs.Storage[m.ID] = m

	if cs.StoreInterval == 0 {
		s, _ := history.NewSaver(cs.FileStorage)
		s.WriteMetric(m)
		defer s.Close()
	}

}

// Validate and save metrics via POST URI
func (cs *ConfigServer) PostHandlerMetricByURL(w http.ResponseWriter, r *http.Request) {
	m, err := metric.ParseMetricEntityFromURL(r)
	if err != nil {
		switch {
		case err == metric.ErrMissmatchedType:
			log.Println(err)
			http.Error(w, "Wrong type", http.StatusNotImplemented)
			return

		case err == metric.ErrDeltaAssign:
			log.Println(err)
			http.Error(w, "Wrong delta", http.StatusBadRequest)
			return
		case err == metric.ErrValueAssign:
			log.Println(err)
			http.Error(w, "Wrong value", http.StatusBadRequest)
			return
		default:
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	ID := chi.URLParam(r, "id")
	if m.MType == metric.MetricTypeCounter {
		id1, ok := cs.Storage[ID]
		if ok {
			newDelta := id1.Delta + m.Delta
			m.Delta = newDelta
		}
	}
	cs.Storage[ID] = m

}

// GetMetricsByValue return metrics via GET /value/{type}/{id}
func (cs *ConfigServer) GetMetricsByValueURI(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	if metricType != string(metric.MetricTypeGauge) && metricType != string(metric.MetricTypeCounter) {
		log.Println("missmatched type")
		http.Error(w, "missmatched type", http.StatusBadRequest)
		return
	}
	ID := chi.URLParam(r, "id")
	met, ok := cs.Storage[ID]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
	}
	switch {
	case met.MType == metric.MetricTypeGauge:
		fmt.Fprintf(w, "%v", met.Value)
	case met.MType == metric.MetricTypeCounter:
		fmt.Fprintf(w, "%v", met.Delta)
	}

}

// POSTMetricsByValueJSON return metrics via JSON
func (cs *ConfigServer) POSTMetricsByValueJSON(w http.ResponseWriter, r *http.Request) {
	m := metric.Metric{}
	enc := json.NewDecoder(r.Body)
	if err := enc.Decode(&m); err != nil {
		log.Println(err)
		http.Error(w, "wrong format", http.StatusBadRequest)
		return
	}
	metric, ok := cs.Storage[m.ID]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
	}
	if err := json.NewEncoder(w).Encode(&metric); err != nil {
		http.Error(w, "unable to marshal the struct", http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")

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

// Return metrics data in html
func (cs *ConfigServer) GetMetricsAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, v := range cs.Storage {
		fmt.Fprintf(w, "%v<br>", v)
	}
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
