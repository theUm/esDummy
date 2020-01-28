package elastic

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
)

type Storage struct {
	client *elasticsearch.Client
}

type Config struct {
	ConnectionString string `env:"ELASTIC_URI"`
}

func New(_ context.Context, c Config) (*Storage, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{c.ConnectionString},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating the client")
	}

	storage := &Storage{
		client: es,
	}
	err = storage.Check()
	if err != nil {
		return nil, errors.Wrapf(err, "error of initial es cluster ping")
	}
	return storage, nil
}

func (s *Storage) Check() error {
	ping, err := s.client.Ping()
	if err != nil {
		return errors.Wrapf(err, "elastic ping error: %s")
	}
	if ping.IsError() {
		return fmt.Errorf("elastic ping error: %s", ping.String())
	}
	return nil
}
