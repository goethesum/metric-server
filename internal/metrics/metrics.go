package metric

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MetricType string

var (
	ErrMissmatchedType = errors.New("missmatched type")
	ErrDeltaAssign     = errors.New("type insn't a int64")
	ErrValueAssign     = errors.New("type insn't a float64")
)

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

const (
	queryKeyMetricID    = "id"
	queryKeyMetricValue = "value"
	queryKeyMetricDelta = "delta"
	queryKeyMetricType  = "type"
)

type Metric struct {
	ID    string     `json:"id"`
	MType MetricType `json:"type"`
	Delta int64      `json:"delta,omitempty"`
	Value float64    `json:"value,omitempty"`
}

func (m Metric) MarshalJSON() (data []byte, err error) {

	MetricJSON := &struct {
		ID    string     `json:"id"`
		Mtype MetricType `json:"type"`
		Delta *int64     `json:"delta,omitempty"`
		Value *float64   `json:"value,omitempty"`
	}{}

	switch {

	case m.MType == MetricTypeCounter:

		MetricJSON.ID = m.ID
		MetricJSON.Mtype = m.MType
		MetricJSON.Delta = &m.Delta
		MetricJSON.Value = nil

		marshMetric, err := json.Marshal(MetricJSON)

		return marshMetric, err

	case m.MType == MetricTypeGauge:

		MetricJSON.ID = m.ID
		MetricJSON.Mtype = m.MType
		MetricJSON.Delta = nil
		MetricJSON.Value = &m.Value

		marshMetric, err := json.Marshal(MetricJSON)

		return marshMetric, err

	}

	return nil, errors.New("missmatched type in 123 MarshalJSON")
}

func (m *Metric) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}

	if err := json.Unmarshal(data, &v); err != nil {
		log.Printf("error during UnamarshalJSON %s", err)
		return err
	}

	MetricJSON := &struct {
		ID    string     `json:"id"`
		Mtype MetricType `json:"type"`
		Delta *int64     `json:"delta,omitempty"`
		Value *float64   `json:"value,omitempty"`
	}{}

	switch {
	case v["type"].(string) == string(MetricTypeCounter):

		if err := json.Unmarshal(data, &MetricJSON); err != nil {
			return err
		}

		m.ID = MetricJSON.ID
		m.MType = MetricJSON.Mtype
		if MetricJSON.Delta != nil {
			m.Delta = *MetricJSON.Delta
		}
	case v["type"].(string) == string(MetricTypeGauge):
		if err := json.Unmarshal(data, &MetricJSON); err != nil {
			return err
		}
		m.ID = MetricJSON.ID
		m.MType = MetricJSON.Mtype
		if MetricJSON.Value != nil {
			m.Value = *MetricJSON.Value
		}

	}

	return nil
}

type AgentStorage struct {
	Stats runtime.MemStats
	Data  map[string]Metric
}

func ParseMetricEntityFromURL(r *http.Request) (Metric, error) {
	m := Metric{}
	if m.ID = chi.URLParam(r, queryKeyMetricID); m.ID == "" {

		return Metric{}, errors.New("empty \"id\" query param")
	}
	if m.MType = MetricType(chi.URLParam(r, queryKeyMetricType)); m.MType == "" {
		return Metric{}, errors.New("empty \"type\" query param")
	}
	if m.MType != MetricTypeGauge && m.MType != MetricTypeCounter {
		return Metric{}, ErrMissmatchedType
	}

	switch {
	case m.MType == MetricTypeCounter:
		var err error
		v := chi.URLParam(r, queryKeyMetricValue)
		m.Delta, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return Metric{}, ErrDeltaAssign
		}
	case m.MType == MetricTypeGauge:
		var err error
		v := chi.URLParam(r, queryKeyMetricValue)
		m.Value, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return Metric{}, ErrValueAssign
		}
	}

	return m, nil
}

func (as *AgentStorage) PopulateMetricStruct() {

	runtime.ReadMemStats(&as.Stats)
	as.Data["Alloc"] = Metric{
		ID:    "Alloc",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.Alloc),
	}
	as.Data["BuckHashSys"] = Metric{
		ID:    "BuckHashSys",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.BuckHashSys),
	}
	as.Data["GCSys"] = Metric{
		ID:    "GCSys",
		MType: MetricTypeCounter,
		Value: float64(as.Stats.GCSys),
	}
	as.Data["GCCPUFraction"] = Metric{
		ID:    "GCCPUFraction",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.GCCPUFraction),
	}
	as.Data["Frees"] = Metric{
		ID:    "Frees",
		MType: MetricTypeCounter,
		Value: float64(as.Stats.Frees),
	}
	as.Data["HeapAlloc"] = Metric{
		ID:    "HeapAlloc",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.HeapAlloc),
	}
	as.Data["HeapIdle"] = Metric{
		ID:    "HeapIdle",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.HeapIdle),
	}
	as.Data["HeapInuse"] = Metric{
		ID:    "HeapInuse",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.HeapInuse),
	}
	as.Data["HeapObjects"] = Metric{
		ID:    "HeapObjects",
		MType: MetricTypeCounter,
		Value: float64(as.Stats.HeapObjects),
	}
	as.Data["HeapReleased"] = Metric{
		ID:    "HeapReleased",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.HeapReleased),
	}
	as.Data["HeapSys"] = Metric{
		ID:    "HeapSys",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.HeapSys),
	}
	as.Data["LastGC"] = Metric{
		ID:    "LastGC",
		MType: MetricTypeCounter,
		Value: float64(as.Stats.LastGC),
	}
	as.Data["Lookups"] = Metric{
		ID:    "Lookups",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.Lookups),
	}
	as.Data["MCacheInuse"] = Metric{
		ID:    "MCacheInuse",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.MCacheInuse),
	}
	as.Data["MCacheSys"] = Metric{
		ID:    "MCacheSys",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.MCacheSys),
	}
	as.Data["MSpanInuse"] = Metric{
		ID:    "MSpanInuse",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.MSpanInuse),
	}
	as.Data["MSpanSys"] = Metric{
		ID:    "MSpanSys",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.MSpanSys),
	}
	as.Data["Mallocs"] = Metric{
		ID:    "Mallocs",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.Mallocs),
	}
	as.Data["NextGC"] = Metric{
		ID:    "NextGC",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.NextGC),
	}
	as.Data["NumForcedGC"] = Metric{
		ID:    "NumForcedGC",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.NumForcedGC),
	}
	as.Data["NumGC"] = Metric{
		ID:    "NumGC",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.NumGC),
	}
	as.Data["OtherSys"] = Metric{
		ID:    "OtherSys",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.OtherSys),
	}
	as.Data["PauseTotalNs"] = Metric{
		ID:    "PauseTotalNs",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.PauseTotalNs),
	}
	as.Data["StackInuse"] = Metric{
		ID:    "StackInuse",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.StackInuse),
	}
	as.Data["StackSys"] = Metric{
		ID:    "StackSys",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.StackSys),
	}
	as.Data["Sys"] = Metric{
		ID:    "Sys",
		MType: MetricTypeGauge,
		Value: float64(as.Stats.Sys),
	}
	as.Data["PollCount"] = Metric{
		ID:    "PollCount",
		MType: MetricTypeCounter,
		Delta: 0,
	}
	as.Data["RandomValue"] = Metric{
		ID:    "RandomValue",
		MType: MetricTypeGauge,
		Value: rand.Float64(),
	}

}
