package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/goethesum/-go-musthave-devops-tpl/internal/config"
	"github.com/goethesum/-go-musthave-devops-tpl/internal/history"
	metric "github.com/goethesum/-go-musthave-devops-tpl/internal/metrics"
)

var confServ *config.ConfigServer

func main() {

	// Setup environmet
	confServ = &config.ConfigServer{
		Storage: make(map[string]*metric.Metric),
		Mutex:   &sync.Mutex{},
	}
	// read env variable
	if err := env.Parse(confServ); err != nil {
		fmt.Printf("%+v\n", err)
	}

	server := &http.Server{
		Addr:    confServ.PortNumber,
		Handler: router(confServ),
	}

	// Handling signal, waiting for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		for sig := range sigCh {
			log.Println("Recieved sig:", sig)
			fmt.Println("Dying...")
			server.Shutdown(context.Background())
		}

	}()
	// Restore metrics from file FILE_STORAGE_PATH
	r, err := history.NewRestorer(confServ.FileStorage)
	if err != nil {
		log.Printf("error during restore from file %s", err)
		return
	}

	confServ.Storage, err = r.RestoreMetrics()
	if err != nil {
		log.Fatalf("dying by...:%s", err)
	}
	r.Close()

	log.Println("Starting on port:", confServ.PortNumber)
	log.Fatal(server.ListenAndServe())

}

func router(cs *config.ConfigServer) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)

	mux.Route("/", func(mux chi.Router) {
		mux.Get("/", cs.GetMetricsAll)
		mux.Get("/{id}", cs.GetMetricsByID)
		mux.Post("/", cs.PostHandlerMetrics)
	})
	mux.Route("/update", func(mux chi.Router) {
		mux.Get("/", cs.GetMetricsAll)
		mux.Post("/", cs.PostHandlerMetricsJSON)
		mux.Post("/{type}/{id}/{value}", cs.PostHandlerMetricByURL)
	})
	mux.Route("/value", func(mux chi.Router) {
		mux.Get("/", cs.GetCheck)
		mux.Get("/{type}/{id}", cs.GetMetricsByValue)
	})

	mux.Get("/metric", cs.GetMetrics)
	return mux
}

// Post "http://localhost:8080/update/gauge/githubActionGauge/100"
// Get "http://localhost:8080/value/gauge/BuckHashSys"
