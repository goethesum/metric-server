package metric

import (
	"errors"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
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
	ID    string
	Type  MetricType
	Value string
}

func (m Metric) NewSendURL() (string, error) {
	params := url.Values{}
	params.Add("id", m.ID)
	params.Add("type", string(m.Type))
	params.Add("value", m.Value)
	return params.Encode(), nil
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

func (as *AgentStorage) PopulateMetricStruct() {
	runtime.ReadMemStats(&as.Stats)
	as.Data["Alloc"] = Metric{
		ID:    "Alloc",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.Alloc, 10),
	}
	as.Data["TotalAlloc"] = Metric{
		ID:    "TotalAlloc",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(as.Stats.TotalAlloc, 10),
	}
	as.Data["Sys"] = Metric{
		ID:    "Sys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.Sys, 10),
	}
	as.Data["Lookups"] = Metric{
		ID:    "Lookups",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.Lookups, 10),
	}
	as.Data["Mallocs"] = Metric{
		ID:    "Mallocs",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.Mallocs, 10),
	}
	as.Data["Frees"] = Metric{
		ID:    "Frees",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(as.Stats.Mallocs, 10),
	}
	as.Data["HeapAlloc"] = Metric{
		ID:    "HeapAlloc",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.HeapAlloc, 10),
	}
	as.Data["HeapSys"] = Metric{
		ID:    "HeapSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.HeapSys, 10),
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
	as.Data["HeapReleased"] = Metric{
		ID:    "HeapReleased",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.HeapReleased, 10),
	}
	as.Data["HeapObjects"] = Metric{
		ID:    "HeapObjects",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(as.Stats.HeapObjects, 10),
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
	as.Data["MCacheInuse"] = Metric{
		ID:    "MCacheInuse",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.MCacheInuse, 10),
	}
	as.Data["MCacheSys"] = Metric{
		ID:    "MCacheSys",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(as.Stats.MCacheSys, 10),
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
	as.Data["OtherSys"] = Metric{
		ID:    "OtherSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.OtherSys, 10),
	}
	as.Data["NextGC"] = Metric{
		ID:    "NextGC",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(as.Stats.NextGC, 10),
	}
	as.Data["LastGC"] = Metric{
		ID:    "LastGC",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(as.Stats.LastGC, 10),
	}
	as.Data["PauseTotalNs"] = Metric{
		ID:    "PauseTotalNs",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(as.Stats.PauseTotalNs, 10),
	}

}
