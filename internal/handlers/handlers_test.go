package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/goethesum/-go-musthave-devops-tpl/internal/env"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"GetMetrics", "/testMetric", "GET", []postData{}, http.StatusOK},
	{"PostHandlerMetrics", "/update", "POST", []postData{
		{key: "id", value: "testMetric"},
		{key: "type", value: "counter"},
		{key: "value", value: "3742"},
	}, http.StatusOK},
}

func getRouter(e *env.Env) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/update", Repo.PostHandlerMetrics)
	mux.HandleFunc("/", Repo.GetMetrics)

	repo := NewRepo(e)
	NewHandlers(repo)

	return mux
}

func TestHandlers(t *testing.T) {
	eTest := &env.Env{
		Data: make(map[string]env.MetricServer),
	}
	routes := getRouter(eTest)
	ts := httptest.NewServer(routes)
	defer ts.Close()

	for _, tt := range theTests {
		fmt.Println(ts.URL + tt.url)
		if tt.method == "GET" {
			request := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(Repo.GetMetrics)
			h.ServeHTTP(w, request)
			resp := w.Result()
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", tt.name, tt.expectedStatusCode, resp.StatusCode)
			}

		} else {
			values := url.Values{}
			for _, a := range tt.params {
				values.Add(a.key, a.value)
			}
			buf := new(bytes.Buffer)
			buf.WriteString(values.Encode())
			request := httptest.NewRequest("POST", tt.url, buf)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(Repo.PostHandlerMetrics)
			h.ServeHTTP(w, request)
			resp := w.Result()

			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", tt.name, tt.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
