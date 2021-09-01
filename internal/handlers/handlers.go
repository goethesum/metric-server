package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goethesum/-go-musthave-devops-tpl/internal/env"
)

type Repositry interface {
	PostHandlerMetrics(w http.ResponseWriter, r *http.Request)
	GetMetrics(w http.ResponseWriter, r *http.Request)
}

// Repository is the repository type
type Repository struct {
	E *env.Env
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepois returns reference to Repository struct
func NewRepo(e *env.Env) *Repository {
	return &Repository{
		E: e,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// PostHandlerMetrics parses form from POST request and assigns data to cache
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

}

// Return metric data in JSON by Requested URI
func (e *Repository) GetMetrics(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get("id")

	jsonMetric, err := json.Marshal(e.E.Data[key])
	if err != nil {
		log.Fatalf("unable to marshal the struct %s", err)
	}
	w.Write(jsonMetric)

}
