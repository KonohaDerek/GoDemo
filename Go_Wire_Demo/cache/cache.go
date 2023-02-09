package cache

import (
	"time"

	"k-derek.dev/demo/wire/config"
)

type Service struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

func New(cfg config.Config) (*Service, error) {
	return &Service{
		DefaultExpiration: cfg.Cache.DefaultExpiration,
		CleanupInterval:   cfg.Cache.CleanupInterval,
	}, nil
}
