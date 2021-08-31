package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/goethesum/-go-musthave-devops-tpl/internal/env"
)

// Repository is the repository type
type Repository struct {
	E *env.Env
}

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
func NewRepo(e *env.Env) *Repository {
	return &Repository{
		E: e,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (e *Repository) PostHandlerMetrics(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("unable to parse the form %s", err)
	}
	e.E.Data[r.Form.Get("id")] = env.MetricServer{
		ID:    r.Form.Get("id"),
		Type:  r.Form.Get("type"),
		Value: r.Form.Get("value"),
	}
	//fmt.Println(e.E.Data[r.Form.Get("id")])

}

// Return metric data in JSON by Requested URI
func (e *Repository) GetMetrics(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.RequestURI, "/")
	jsonMetric, err := json.Marshal(e.E.Data[key])
	if err != nil {
		log.Fatalf("unable to marshal the struct %s", err)
	}
	w.Write(jsonMetric)
}
