package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/goethesum/-go-musthave-devops-tpl/internal/config"
	"github.com/goethesum/-go-musthave-devops-tpl/internal/history"
	metric "github.com/goethesum/-go-musthave-devops-tpl/internal/metrics"
)

var srv *config.Service
var confServ *config.ConfigServer

func main() {

	confServ = config.NewConfigServer()

	flag.StringVar(&confServ.Address, "a", "localhost:8080", "server address")
	flag.BoolVar(&confServ.Restore, "r", false, "true or false for restore from file")
	flag.DurationVar(&confServ.StoreInterval, "i", 300*time.Second, "Interval in seconds between file savings")
	flag.StringVar(&confServ.StoreFile, "f", "/tmp/devops-metrics-db.json", "file path")

	// flag parsing
	flag.Parse()
	confServ.StoreFile = flag.Arg(0)

	// read env variable
	if err := env.Parse(confServ); err != nil {
		log.Fatal(err)
	}

	log.Printf("Address: %s, Path %s, Interval %s, Restore %t", confServ.Address, confServ.StoreFile, confServ.StoreInterval, confServ.Restore)

	srv = config.NewService(confServ)

	// Setup service
	srv = &config.Service{
		Storage: make(map[string]metric.Metric),
		Mutex:   &sync.Mutex{},
		Server:  *confServ,
	}

	server := &http.Server{
		Addr:    confServ.Address,
		Handler: router(srv),
	}

	// Handling signal, waiting for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		for sig := range sigCh {
			log.Println("Recieved sig:", sig)
			fmt.Println("Dying...")
			server.Shutdown(context.Background())
		}

	}()
	// Restore metrics from STOREFILE
	if confServ.Restore {
		r, err := history.NewRestorer(confServ.StoreFile)
		if err != nil {
			log.Println("nothing to restore", err)
		}

		srv.Storage, err = r.RestoreMetrics()
		if err != nil {
			log.Fatalf("dying by...:%s", err)
		}
		r.Close()
	}

	// if int64(confServ.StoreInterval) > 0 {
	// 	go func() {
	// 		tck := time.NewTicker(confServ.StoreInterval)

	// 		select {
	// 		case <-sigCh:
	// 			tck.Stop()
	// 			return
	// 		case <-tck.C:
	// 			s, _ := history.NewSaver(srv.Server.StoreFile)
	// 			s.StoreMetrics(srv.Storage)
	// 			defer s.Close()
	// 		}

	// 	}()
	// }

	log.Println("Starting on port:", confServ.Address)
	log.Fatal(server.ListenAndServe())

}

func router(s *config.Service) http.Handler {
	mux := chi.NewRouter()

	mux.Use(
		middleware.Recoverer,
		middleware.Logger,
	)

	mux.Route("/", func(mux chi.Router) {
		mux.Get("/", s.GetMetricsAll)
	})
	mux.Route("/update", func(mux chi.Router) {
		mux.Get("/", s.GetMetricsAll)
		mux.Post("/", s.PostHandlerMetricsJSON)
		mux.Post("/{type}/{id}/{value}", s.PostHandlerMetricByURL)
	})
	mux.Route("/value", func(mux chi.Router) {
		mux.Post("/", s.POSTMetricsByValueJSON)
		mux.Get("/{type}/{id}", s.GetMetricsByValueURI)
	})

	return mux
}

// curl -X POST http://localhost:8080/value -H 'Content-Type: application/json' -d '{"id":"Sys","type":"gauge"}'
// Post "http://localhost:8080/update/gauge/githubActionGauge/100"
// Get "http://localhost:8080/value/gauge/BuckHashSys"
