package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/goethesum/-go-musthave-devops-tpl/internal/env"
	"github.com/goethesum/-go-musthave-devops-tpl/internal/handlers"
)

func main() {

	// Setup environmet
	e := &env.Env{
		PortNumber: ":8080",
		Data:       make(map[string]env.MetricServer),
	}

	repo := handlers.NewRepo(e)
	handlers.NewHandlers(repo)

	// mux := http.NewServeMux()
	// mux.HandleFunc("/push", handlers.Repo.PostHandlerMetrics)
	// mux.HandleFunc("/", handlers.Repo.GetMetrics)

	server := &http.Server{
		Addr:    e.PortNumber,
		Handler: router(e),
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

	fmt.Println("Starting on port:", e.PortNumber)
	log.Fatal(server.ListenAndServe())

}

func router(e *env.Env) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/push", handlers.Repo.PostHandlerMetrics)
	mux.HandleFunc("/", handlers.Repo.GetMetrics)

	return mux
}
