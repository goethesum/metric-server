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
	Address        string        `env:"ADDRESS"`
	URLMetricPush  string        `env:"URL_PATH" envDefault:"/update"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}

type ConfigServer struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

type Service struct {
	Storage map[string]metric.Metric
	Server  ConfigServer
	*sync.Mutex
}

func NewService(srv *ConfigServer) *Service {
	return &Service{
		Storage: make(map[string]metric.Metric),
		Server:  *srv,
		Mutex:   &sync.Mutex{},
	}
}

func NewConfigServer() *ConfigServer {
	return &ConfigServer{}
}

func NewConfigAgent() *ConfigAgent {
	return &ConfigAgent{}
}

func (s *Service) PostHandlerMetricsJSON(w http.ResponseWriter, r *http.Request) {
	m := metric.Metric{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		log.Printf("unable to decode params in PostHandlerMetricsJSON, %s", err)
		return
	}

	if m.MType == metric.MetricTypeCounter {
		id1, ok := s.Storage[m.ID]
		if ok {
			newDelta := id1.Delta + m.Delta
			m.Delta = newDelta
		}
	}

	s.Storage[m.ID] = m

	if s.Server.StoreInterval == 0 {
		s, _ := history.NewSaver(s.Server.StoreFile)
		s.WriteMetric(m)
		defer s.Close()
	}

}

// Validate and save metrics via POST URI
func (s *Service) PostHandlerMetricByURL(w http.ResponseWriter, r *http.Request) {
	m, err := metric.ParseMetricEntityFromURL(r)
	if err != nil {
		switch err {
		case metric.ErrMissmatchedType:
			log.Println(err)
			http.Error(w, "Wrong type", http.StatusNotImplemented)
			return
		case metric.ErrDeltaAssign:
			log.Println(err)
			http.Error(w, "Wrong delta", http.StatusBadRequest)
			return
		case metric.ErrValueAssign:
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
		id1, ok := s.Storage[ID]
		if ok {
			newDelta := id1.Delta + m.Delta
			m.Delta = newDelta
		}
	}
	s.Storage[ID] = m

}

// GetMetricsByValue return metrics via GET /value/{type}/{id}
func (s *Service) GetMetricsByValueURI(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	if metricType != string(metric.MetricTypeGauge) && metricType != string(metric.MetricTypeCounter) {
		log.Println("missmatched type")
		http.Error(w, "missmatched type", http.StatusBadRequest)
		return
	}
	ID := chi.URLParam(r, "id")
	met, ok := s.Storage[ID]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	switch met.MType {
	case metric.MetricTypeGauge:
		fmt.Fprintf(w, "%v", met.Value)
	case metric.MetricTypeCounter:
		fmt.Fprintf(w, "%v", met.Delta)
	}

}

// POSTMetricsByValueJSON return metrics via JSON
func (s *Service) POSTMetricsByValueJSON(w http.ResponseWriter, r *http.Request) {
	m := metric.Metric{}
	enc := json.NewDecoder(r.Body)
	if err := enc.Decode(&m); err != nil {
		log.Println(err)
		http.Error(w, "wrong format", http.StatusBadRequest)
		return
	}
	if m.MType != metric.MetricTypeGauge && m.MType != metric.MetricTypeCounter {
		log.Println("missmatched type")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	metric, ok := s.Storage[m.ID]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
	}
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&metric); err != nil {
		http.Error(w, "unable to marshal the struct", http.StatusBadRequest)
		return
	}

}

// Return metric data in JSON by Requested URI
func (s *Service) GetMetrics(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get("id")

	value, err := s.GetMetricsByKey(context.Background(), key)
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
func (s *Service) GetMetricsAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, v := range s.Storage {
		fmt.Fprintf(w, "%v<br>", v)
	}
}

func (s *Service) GetMetricsByKey(ctx context.Context, key string) (metric.Metric, error) {
	s.Lock()
	defer s.Unlock()

	m, ok := s.Storage[key]
	if !ok {
		return metric.Metric{}, errors.New("metric not found")
	}
	return m, nil

}

// curl -X POST http://localhost:8080/value -H 'Content-Type: application/json' -d '{"id":"BuckHashSys","type":"gauge"}'
