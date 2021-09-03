package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/goethesum/-go-musthave-devops-tpl/internal/config"
	metric "github.com/goethesum/-go-musthave-devops-tpl/internal/metrics"
)

var confServ *config.ConfigServer

func main() {

	// Setup environmet
	confServ = &config.ConfigServer{
		PortNumber: ":8080",
		Storage:    make(map[string]metric.Metric),
		Mutex:      &sync.Mutex{},
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

	log.Println("Starting on port:", confServ.PortNumber)
	log.Fatal(server.ListenAndServe())

}

func router(cs *config.ConfigServer) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Get("/pull", cs.GetMetricsAll)
	mux.Get("/metric", cs.GetMetrics)
	mux.Post("/update", cs.PostHandlerMetrics)

	return mux
}
