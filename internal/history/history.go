package history

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	metric "github.com/goethesum/-go-musthave-devops-tpl/internal/metrics"
)

type Saver interface {
	WriteMetric(savedMetric *metric.Metric)
	Close() error
}

type Restorer interface {
	RestoreMetric() (*metric.Metric, error)
	Close() error
}

type saver struct {
	file   *os.File
	writer *bufio.Writer
}

type restorer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewSaver(fileName string) (*saver, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &saver{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil

}

func NewRestorer(fileName string) (*restorer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &restorer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (s *saver) Close() error {
	return s.file.Close()
}

func (r *restorer) Close() error {
	return r.file.Close()
}

func (s *saver) WriteMetric(savedMetric metric.Metric) error {
	data, err := json.Marshal(&savedMetric)
	if err != nil {
		return err
	}
	if _, err := s.writer.Write(data); err != nil {
		return err
	}
	if err := s.writer.WriteByte('\n'); err != nil {
		return err
	}
	return s.writer.Flush()
}
func (s *saver) StoreMetrics(storage map[string]metric.Metric) error {
	data, err := json.Marshal(&storage)
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err := s.writer.Write(data); err != nil {
		log.Println(err)
		return err
	}

	return s.writer.Flush()
}

func (r *restorer) RestoreMetrics() (map[string]metric.Metric, error) {
	store := make(map[string]metric.Metric)
	for r.scanner.Scan() {
		item := metric.Metric{}
		data := r.scanner.Bytes()

		err := json.Unmarshal(data, &item)
		if err != nil {
			return nil, err
		}
		store[item.ID] = item

	}

	return store, nil
}
