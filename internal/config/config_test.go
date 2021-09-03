package config

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	metric "github.com/goethesum/-go-musthave-devops-tpl/internal/metrics"
)

var testCs = &ConfigServer{
	Storage: map[string]metric.Metric{
		"test": {
			ID:    "test",
			Type:  metric.MetricTypeCounter,
			Value: "4343",
		},
	},
	Mutex: &sync.Mutex{},
}

type postData struct {
	key   string
	value string
}

var theTestsGet = struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{"GetMetrics", "/", "GET", http.StatusOK}

var theTestsPost = struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{"PostHandlerMetrics", "/update", "POST", []postData{
	{key: "id", value: "testMetric"},
	{key: "type", value: "counter"},
	{key: "value", value: "3742"},
}, http.StatusOK}

func getRouter(cs *ConfigServer) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/update", cs.PostHandlerMetrics)
	mux.HandleFunc("/pull", cs.GetMetrics)

	return mux
}

func TestHandlerGet(t *testing.T) {

	routes := getRouter(testCs)
	ts := httptest.NewServer(routes)
	defer ts.Close()

	request := httptest.NewRequest("GET", theTestsGet.url, nil)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(testCs.GetMetricsAll)
	h.ServeHTTP(w, request)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != theTestsGet.expectedStatusCode {
		t.Errorf("for %s, expected %d but got %d", theTestsGet.name, theTestsGet.expectedStatusCode, resp.StatusCode)
	}
}

func TestHandlerPost(t *testing.T) {

	values := url.Values{}
	for _, a := range theTestsPost.params {
		values.Add(a.key, a.value)
	}

	buf := new(bytes.Buffer)
	buf.WriteString(values.Encode())
	request := httptest.NewRequest("POST", theTestsPost.url, buf)

	w := httptest.NewRecorder()
	h := http.HandlerFunc(testCs.PostHandlerMetrics)
	h.ServeHTTP(w, request)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != theTestsPost.expectedStatusCode {
		t.Errorf("for %s, expected %d but got %d", theTestsPost.name, theTestsPost.expectedStatusCode, resp.StatusCode)
	}

}
