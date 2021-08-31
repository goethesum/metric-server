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
	ms.Data["TotalAlloc"] = Metric{
		ID:    "TotalAlloc",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(ms.Stats.TotalAlloc, 10),
	}
	ms.Data["Sys"] = Metric{
		ID:    "Sys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.Sys, 10),
	}
	ms.Data["Lookups"] = Metric{
		ID:    "Lookups",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.Lookups, 10),
	}
	ms.Data["Mallocs"] = Metric{
		ID:    "Mallocs",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.Mallocs, 10),
	}
	ms.Data["Frees"] = Metric{
		ID:    "Frees",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(ms.Stats.Mallocs, 10),
	}
	ms.Data["HeapAlloc"] = Metric{
		ID:    "HeapAlloc",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.HeapAlloc, 10),
	}
	ms.Data["HeapSys"] = Metric{
		ID:    "HeapSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.HeapSys, 10),
	}
	ms.Data["HeapIdle"] = Metric{
		ID:    "HeapIdle",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.HeapIdle, 10),
	}
	ms.Data["HeapInuse"] = Metric{
		ID:    "HeapInuse",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.HeapInuse, 10),
	}
	ms.Data["HeapReleased"] = Metric{
		ID:    "HeapReleased",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.HeapReleased, 10),
	}
	ms.Data["HeapObjects"] = Metric{
		ID:    "HeapObjects",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(ms.Stats.HeapObjects, 10),
	}
	ms.Data["StackInuse"] = Metric{
		ID:    "StackInuse",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.StackInuse, 10),
	}
	ms.Data["StackSys"] = Metric{
		ID:    "StackSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.StackSys, 10),
	}
	ms.Data["MSpanInuse"] = Metric{
		ID:    "MSpanInuse",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.MSpanInuse, 10),
	}
	ms.Data["MSpanSys"] = Metric{
		ID:    "MSpanSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.MSpanSys, 10),
	}
	ms.Data["MCacheInuse"] = Metric{
		ID:    "MCacheInuse",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.MCacheInuse, 10),
	}
	ms.Data["MCacheSys"] = Metric{
		ID:    "MCacheSys",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(ms.Stats.MCacheSys, 10),
	}
	ms.Data["BuckHashSys"] = Metric{
		ID:    "BuckHashSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.BuckHashSys, 10),
	}
	ms.Data["GCSys"] = Metric{
		ID:    "GCSys",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(ms.Stats.GCSys, 10),
	}
	ms.Data["OtherSys"] = Metric{
		ID:    "OtherSys",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.OtherSys, 10),
	}
	ms.Data["NextGC"] = Metric{
		ID:    "NextGC",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(ms.Stats.NextGC, 10),
	}
	ms.Data["LastGC"] = Metric{
		ID:    "LastGC",
		Type:  MetricTypeCounter,
		Value: strconv.FormatUint(ms.Stats.LastGC, 10),
	}
	ms.Data["PauseTotalNs"] = Metric{
		ID:    "PauseTotalNs",
		Type:  MetricTypeGauge,
		Value: strconv.FormatUint(ms.Stats.PauseTotalNs, 10),
	}

}
