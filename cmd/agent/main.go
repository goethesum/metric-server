package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/goethesum/-go-musthave-devops-tpl/internal/config"
	metric "github.com/goethesum/-go-musthave-devops-tpl/internal/metrics"
)

var conf = config.ConfigAgent{
	Server:        getEnv("SERVER_ADDR", "http://localhost:8080/"),
	URLMetricPush: getEnv("URL_PATH", "update"),
	TimeInterval:  5,
}

// "http://localhost:8080/"

type HTTPClient interface {
	MetricSend(ctx context.Context, metrics metric.Metric) error
}

type clientHTTP struct {
	client resty.Client
}

// Looks up for ENV variables, returns default if it doesn't exist
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// MetricSend takes Server address and relative path from config struct
// Calls NewSendUrl to construct encoded URL
func (client *clientHTTP) MetricSend(ctx context.Context, endpoint string, metrics metric.Metric) (*resty.Response, error) {
	url, err := metrics.NewSendURL()
	if err != nil {
		return nil, fmt.Errorf("unable to parse url:%w", err)
	}

	resp, err := client.client.SetTimeout(1 * time.Second).
		R().
		Post(endpoint + "?" + url)
	if err != nil {
		return nil, fmt.Errorf("unable to send POST request:%w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status cod [%d]: %s", resp.StatusCode(), string(resp.Body()))
	}

	return resp, nil
}

func main() {

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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()

	// make endpoint
	endpoint := conf.Server + conf.URLMetricPush
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
					log.Println("Failed to send", v.ID)

				}
				if resp != nil {
					log.Println(resp.StatusCode(), v.ID)
				}
			}

		}

	}
}
