package env

import (
	"context"
	"errors"
	"sync"
)

//Holds environment parametrs for the Server

//Env holds the application environment
type Env struct {
	PortNumber string
	Data       map[string]MetricServer
	*sync.Mutex
}

// Holds metrics on the server, the field Type has type string,
// because r.Form.Get() returns string
type MetricServer struct {
	ID    string
	Type  string
	Value string
}

func (e *Env) GetMetricsByKey(ctx context.Context, key string) (MetricServer, error) {
	e.Lock()
	defer e.Unlock()

	m, ok := e.Data[key]
	if !ok {
		return MetricServer{}, errors.New("metric not found")
	}
	return m, nil

}
