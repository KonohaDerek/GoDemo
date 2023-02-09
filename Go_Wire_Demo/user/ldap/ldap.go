package ldap

import (
	"k-derek.dev/demo/wire/cache"
	"k-derek.dev/demo/wire/config"
)

type Service struct {
	bindUsername string
	bindPassword string
	cache        *cache.Service
}

func New(cfg config.Config, cache *cache.Service) (*Service, error) {
	return &Service{
		bindUsername: cfg.Ldap.BindUsername,
		bindPassword: cfg.Ldap.BindPassword,
		cache:        cache,
	}, nil
}
