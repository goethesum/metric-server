package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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
func (client *clientHTTP) MetricSend(endpoint string, metrics metric.Metric, tr *http.Transport) (*resty.Response, error) {
	jsonMetric, err := json.Marshal(metrics)
	if err != nil {
		return nil, fmt.Errorf("error during marshaling in MetricSend %w", err)
	}

	resp, err := client.client.SetCloseConnection(true).
		SetTransport(tr).
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

	conf := config.NewConfigAgent()

	flag.StringVar(&conf.Address, "a", "http://localhost:8080", "server address")
	flag.DurationVar(&conf.ReportInterval, "r", 10*time.Second, "duration of Report Interval")
	flag.DurationVar(&conf.PollInterval, "i", 2*time.Second, "duration of Poll Interval")

	// read env variable
	if err := env.Parse(conf); err != nil {
		log.Fatal(err)
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
	transport := &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 20,
	}
	// make endpoint
	endpoint := conf.Address + conf.URLMetricPush
	log.Println(endpoint)
	// Create Ticker for populating
	tickPoll := time.NewTicker(conf.PollInterval)
	defer tickPoll.Stop()
	// Create Ticker for report
	tickReport := time.NewTicker(conf.ReportInterval)
	defer tickReport.Stop()

	// Handling signal, waiting for graceful shutdown
	done := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for sig := range sigCh {
			log.Println("Recieved sig:", sig)
			tickPoll.Stop()
			tickReport.Stop()
			close(done)
			return
		}

	}()
	log.Println("check")

	// Poll every 2s
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Stopped")
				return
			case <-tickPoll.C:
				mStorage.PopulateMetricStruct()
				log.Println("Polled")
				inc, ok := mStorage.Data["PollCount"]
				if !ok {
					log.Println("Value PollCount doesn't exist")
				}
				inc.Delta += 1
				mStorage.Data["PollCount"] = inc

			}
		}
	}()

	// Report every 10s
	for {
		select {
		case <-done:
			fmt.Println("Stopped")
			return
		case <-tickReport.C:
			for _, v := range mStorage.Data {
				select {
				case <-done:
					return
				default:
					time.Sleep(100 * time.Microsecond)
					resp, err := client.MetricSend(endpoint, v, transport)

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

}
