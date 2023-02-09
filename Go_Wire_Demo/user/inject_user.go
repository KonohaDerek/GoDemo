package user

import (
	"k-derek.dev/demo/wire/cache"
	"k-derek.dev/demo/wire/config"
	"k-derek.dev/demo/wire/user/crowd"
	"k-derek.dev/demo/wire/user/ldap"

	"github.com/google/wire"
)

var UserSet = wire.NewSet( //nolint:deadcode,unused,varcheck
	ProvideUser,
	ProvideLDAP,
	ProvideCROWD,
	ProvideCache,
)

func ProvideUser(
	l *ldap.Service,
	c *crowd.Service,
	cache *cache.Service,
) (*Service, error) {
	return New(l, c, cache)
}

func ProvideLDAP(
	cfg config.Config,
	cache *cache.Service,
) (*ldap.Service, error) {
	return ldap.New(cfg, cache)
}

func ProvideCROWD(
	cfg config.Config,
	cache *cache.Service,
) (*crowd.Service, error) {
	return crowd.New(cfg, cache)
}

func ProvideCache(
	cfg config.Config,
) (*cache.Service, error) {
	return cache.New(cfg)
}
