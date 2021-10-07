package config

import (
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	metric "github.com/goethesum/-go-musthave-devops-tpl/internal/metrics"
)

var testCs = &ConfigServer{
	Storage: map[string]metric.Metric{
		"test": {
			ID:    "test",
			MType: metric.MetricTypeCounter,
			Value: 456,
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
}{" PostHandlerMetricByURL", "/update", "POST", []postData{
	{key: "type", value: "counter"},
	{key: "id", value: "testMetric"},
	{key: "value", value: "3742"},
}, http.StatusOK}

func getRouter(cs *ConfigServer) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/update", cs.PostHandlerMetricByURL)
	mux.HandleFunc("/", cs.GetMetricsAll)

	return mux
}

func getRouterChi(cs *ConfigServer) http.Handler {
	mux := chi.NewRouter()
	mux.Route("/update", func(mux chi.Router) {
		mux.Get("/", cs.GetMetricsAll)
		mux.Post("/", cs.PostHandlerMetricsJSON)
		mux.Post("/{type}/{id}/{value}", cs.PostHandlerMetricByURL)
	})
	return mux
}

func TestHandlerGet(t *testing.T) {

	router := getRouter(testCs)

	ts := httptest.NewServer(router)
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

	router := getRouterChi(testCs)
	ts := httptest.NewServer(router)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/update/counter/id/4324", nil)
	if err != nil {
		log.Println(err)
	}
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != theTestsPost.expectedStatusCode {
		t.Errorf("for %s, expected %d but got %d", theTestsPost.name, theTestsPost.expectedStatusCode, resp.StatusCode)
	}

}
