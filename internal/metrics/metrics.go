package metric

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MetricType string

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

const (
	queryKeyMetricID    = "id"
	queryKeyMetricValue = "value"
	queryKeyMetricType  = "type"
)

type Metric struct {
	ID    string     `json:"id"`
	Type  MetricType `json:"type"`
	Value string     `json:"value"`
}

func (m Metric) MarshalJSON() (data []byte, err error) {
	switch {
	case m.Type == MetricTypeCounter:
		digit, _ := strconv.Atoi(m.Value)
		MetricValue := &struct {
			ID    string     `json:"ID"`
			Type  MetricType `json:"type"`
			Value int        `json:"delta"`
		}{
			ID:    m.ID,
			Type:  m.Type,
			Value: digit,
		}
		return json.Marshal(MetricValue)
	case m.Type == MetricTypeGauge:
		digit, _ := strconv.ParseFloat(m.Value, 64)
		MetricDelta := &struct {
			ID    string     `json:"ID"`
			Type  MetricType `json:"type"`
			Value float64    `json:"value"`
		}{
			ID:    m.ID,
			Type:  m.Type,
			Value: digit,
		}
		return json.Marshal(MetricDelta)
	default:
		return nil, errors.New("missmatched type in MarshalJSON")
	}

}

func (m *Metric) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}

	if err := json.Unmarshal(data, &v); err != nil {
		log.Printf("error during UnamarshalJSON %s", err)
		return err
	}
	fmt.Println(v)
	// counter struct
	aliasDelta := &struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Value int    `json:"delta"`
	}{}
	// gauge struct
	aliasValue := &struct {
		ID    string  `json:"id"`
		Type  string  `json:"type"`
		Value float64 `json:"value"`
	}{}

	switch {
	case v["type"].(string) == string(MetricTypeCounter):

		if err := json.Unmarshal(data, &aliasDelta); err != nil {
			return err
		}
		fmt.Println(aliasDelta)
		m.ID = aliasDelta.ID
		m.Type = MetricType(aliasDelta.Type)
		m.Value = strconv.Itoa(aliasDelta.Value)
		fmt.Println(m)
	case v["type"].(string) == string(MetricTypeGauge):
		if err := json.Unmarshal(data, &aliasValue); err != nil {
			return err
		}
		fmt.Println(aliasValue)
		m.ID = aliasValue.ID
		m.Type = MetricType(aliasValue.Type)
		m.Value = fmt.Sprint(aliasValue.Value)
		fmt.Println(m)
	}
	return nil
}

type AgentStorage struct {
	Stats runtime.MemStats
	Data  map[string]Metric
}

func ParseMetricEntityFromRequest(r *http.Request) (*Metric, error) {
	m := new(Metric)

	if m.ID = r.URL.Query().Get(queryKeyMetricID); m.ID == "" {
		return nil, errors.New("empty \"id\" query param")
	}
	if m.Type = MetricType(r.URL.Query().Get(queryKeyMetricType)); m.Type == "" {
		return nil, errors.New("empty \"type\" query param")
	}
	if m.Value = r.URL.Query().Get(queryKeyMetricValue); m.Value == "" {
		return nil, errors.New("empty \"value\" query param")
	}

	if m.Type != MetricTypeGauge && m.Type != MetricTypeCounter {
		return nil, errors.New("missmatched type")
	}

	return m, nil
}

func ParseMetricEntityFromURL(r *http.Request) (*Metric, error) {
	m := new(Metric)

	if m.ID = chi.URLParam(r, queryKeyMetricID); m.ID == "" {
		return nil, errors.New("empty \"id\" query param")
	}
	if m.Type = MetricType(chi.URLParam(r, queryKeyMetricType)); m.Type == "" {
		return nil, errors.New("empty \"type\" query param")
	}
	if m.Value = chi.URLParam(r, queryKeyMetricValue); m.Value == "" {
		return nil, errors.New("empty \"value\" query param")
	}

	if m.Type != MetricTypeGauge && m.Type != MetricTypeCounter {
		return nil, errors.New("missmatched type")
	}

	return m, nil
}

func (as *AgentStorage) PopulateMetricStruct() {
	runtime.ReadMemStats(&as.Stats)
	as.Data["Alloc"] = Metric{
		ID:    "Alloc",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.Alloc, 10),
	}
	as.Data["BuckHashSys"] = Metric{
		ID:    "BuckHashSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.BuckHashSys, 10),
	}
	as.Data["GCSys"] = Metric{
		ID:    "GCSys",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(as.Stats.GCSys, 10),
	}
	as.Data["GCCPUFraction"] = Metric{
		ID:    "GCCPUFraction",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(uint64(as.Stats.GCCPUFraction), 10),
	}
	as.Data["Frees"] = Metric{
		ID:    "Frees",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(as.Stats.Frees, 10),
	}
	as.Data["HeapAlloc"] = Metric{
		ID:    "HeapAlloc",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.HeapAlloc, 10),
	}
	as.Data["HeapIdle"] = Metric{
		ID:    "HeapIdle",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.HeapIdle, 10),
	}
	as.Data["HeapInuse"] = Metric{
		ID:    "HeapInuse",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.HeapInuse, 10),
	}
	as.Data["HeapObjects"] = Metric{
		ID:    "HeapObjects",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(as.Stats.HeapObjects, 10),
	}
	as.Data["HeapReleased"] = Metric{
		ID:    "HeapReleased",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.HeapReleased, 10),
	}
	as.Data["HeapSys"] = Metric{
		ID:    "HeapSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.HeapSys, 10),
	}
	as.Data["LastGC"] = Metric{
		ID:    "LastGC",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(as.Stats.LastGC, 10),
	}
	as.Data["Lookups"] = Metric{ //
		ID:    "Lookups",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.Lookups, 10),
	}
	as.Data["MCacheInuse"] = Metric{
		ID:    "MCacheInuse",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.MCacheInuse, 10),
	}
	as.Data["MCacheSys"] = Metric{
		ID:    "MCacheSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.MCacheSys, 10),
	}
	as.Data["MSpanInuse"] = Metric{
		ID:    "MSpanInuse",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.MSpanInuse, 10),
	}
	as.Data["MSpanSys"] = Metric{
		ID:    "MSpanSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.MSpanSys, 10),
	}
	as.Data["Mallocs"] = Metric{ //
		ID:    "Mallocs",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.Mallocs, 10),
	}
	as.Data["NextGC"] = Metric{
		ID:    "NextGC",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.NextGC, 10),
	}
	as.Data["NumForcedGC"] = Metric{
		ID:    "NumForcedGC",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(uint64(as.Stats.NumForcedGC), 10),
	}
	as.Data["NumGC"] = Metric{
		ID:    "NumGC",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(uint64(as.Stats.NumGC), 10),
	}
	as.Data["OtherSys"] = Metric{
		ID:    "OtherSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.OtherSys, 10),
	}
	as.Data["PauseTotalNs"] = Metric{
		ID:    "PauseTotalNs",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.PauseTotalNs, 10),
	}
	as.Data["StackInuse"] = Metric{
		ID:    "StackInuse",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.StackInuse, 10),
	}
	as.Data["StackSys"] = Metric{
		ID:    "StackSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.StackSys, 10),
	}
	as.Data["Sys"] = Metric{ //
		ID:    "Sys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.Sys, 10),
	}
	as.Data["PollCount"] = Metric{ //
		ID:    "PollCount",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(as.Stats.Sys, 10),
	}
	as.Data["RandomValue"] = Metric{ //
		ID:    "RandomValue",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.Sys, 10),
	}

}
