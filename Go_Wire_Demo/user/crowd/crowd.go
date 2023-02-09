package crowd

import (
	"k-derek.dev/demo/wire/cache"
	"k-derek.dev/demo/wire/config"
)

type Service struct {
	basicUsername string
	basicPassword string
	cache         *cache.Service
}

func New(cfg config.Config, cache *cache.Service) (*Service, error) {
	return &Service{
		basicUsername: cfg.Crowd.BasicUsername,
		basicPassword: cfg.Crowd.BasicPassword,
		cache:         cache,
	}, nil
}
