package metrics

import (
	"net/url"
	"runtime"
	"strconv"
)

type MetricType string

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
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

type MetricsStorage struct {
	Stats runtime.MemStats
	Data  map[string]Metric
}

func (ms *MetricsStorage) PopulateMetricStruct() {
	runtime.ReadMemStats(&ms.Stats)
	ms.Data["Alloc"] = Metric{
		ID:    "Alloc",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.Alloc, 10),
	}
	ms.Data["Sys"] = Metric{
		ID:    "Sys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.Sys, 10),
	}

}
