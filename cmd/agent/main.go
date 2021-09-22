package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	"github.com/goethesum/-go-musthave-devops-tpl/internal/config"
	metric "github.com/goethesum/-go-musthave-devops-tpl/internal/metrics"
)

type clientHTTP struct {
	client resty.Client
}

// MetricSend takes Server address and relative path from config struct
// Calls NewSendUrl to construct encoded URL
func (client *clientHTTP) MetricSend(ctx context.Context, endpoint string, metrics metric.Metric) (*resty.Response, error) {
	jsonMetric, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("error during marshal in MetricSend %s", err)
	}
	resp, err := client.client.SetTimeout(1*time.Second).
		SetRetryCount(2).
		SetRetryWaitTime(1*time.Second).
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonMetric).
		Post(endpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to send POST request:%w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status cod [%d]: %s", resp.StatusCode(), string(resp.Body()))
	}

	return resp, nil
}

func main() {
	conf := &config.ConfigAgent{
		TimeInterval: 5,
	}

	// read env variable
	if err := env.Parse(conf); err != nil {
		fmt.Printf("%+v\n", err)
	}

	// init client
	client := &clientHTTP{
		client: *resty.New(),
	}

	// Stores agent data
	mStorage := &metric.AgentStorage{
		Stats: runtime.MemStats{},
		Data:  make(map[string]metric.Metric),
	}

	// make ctx
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// make endpoint
	endpoint := conf.Server + conf.URLMetricPush
	fmt.Println(endpoint)
	// Create Ticker
	tick := time.NewTicker(conf.TimeInterval * time.Second)
	defer tick.Stop()
	done := make(chan bool)

	// Handling signal, waiting for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		for sig := range sigCh {
			log.Println("Recieved sig:", sig)
			done <- true
		}

	}()
	// Start handling some logic on each tick
	for {
		select {
		case <-done:
			fmt.Println("Stopped")
			return
		case <-tick.C:
			mStorage.PopulateMetricStruct()
			for _, v := range mStorage.Data {
				resp, err := client.MetricSend(ctx, endpoint, v)

				if err != nil {
					log.Println(err)
					log.Println("Failed to send", v.ID)

				}
				if resp != nil {
					log.Println(resp.StatusCode(), v.ID)
				}
			}

		}

	}
}
