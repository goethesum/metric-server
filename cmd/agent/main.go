package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/goethesum/-go-musthave-devops-tpl/internal/metrics"
)

type Config struct {
	Server        string
	URLMetricPush string
	timeInterval  time.Duration
}

var config = Config{
	Server:        "http://localhost:8080/",
	URLMetricPush: "update",
	timeInterval:  5,
}

type clientHTTP struct {
	client resty.Client
}

// MetricSend takes Server address and relative path from config struct
// Calls NewSendUrl to construct encoded URL
func (client *clientHTTP) MetricSend(ctx context.Context, metrics metrics.Metric) error {
	endpoint := config.Server + config.URLMetricPush
	url, err := metrics.NewSendURL()
	if err != nil {
		return fmt.Errorf("unable to parse url:%w", err)
	}
	resp, err := client.client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Content-Length", strconv.Itoa(len(url))).
		SetBody(url).
		Post(endpoint)
	if err != nil {
		return fmt.Errorf("unable to send POST request:%w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("status cod [%d]: %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

func main() {

	// init client
	client := &clientHTTP{
		client: *resty.New(),
	}
	//init metricStorage
	mStorage := &metrics.MetricsStorage{
		Stats: runtime.MemStats{},
		Data:  make(map[string]metrics.Metric),
	}

	// Create Ticker
	tick := time.NewTicker(config.timeInterval * time.Second)
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
				err := client.MetricSend(context.Background(), v)
				if err != nil {
					fmt.Println(fmt.Errorf("unable to send POST request(ticker):%w", err))
				}
			}

		}

	}
}
