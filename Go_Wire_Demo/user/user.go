package user

import (
	"k-derek.dev/demo/wire/cache"
	"k-derek.dev/demo/wire/user/crowd"
	"k-derek.dev/demo/wire/user/ldap"
)

type Service struct {
	ldap  *ldap.Service
	crowd *crowd.Service
	cache *cache.Service
}

func New(l *ldap.Service, c *crowd.Service, cache *cache.Service) (*Service, error) {
	return &Service{
		ldap:  l,
		crowd: c,
		cache: cache,
	}, nil
}

func (s *Service) Login(username string, password string) bool {
	return true
}
