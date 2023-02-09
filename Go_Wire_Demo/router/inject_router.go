package router

import (
	"net/http"

	"k-derek.dev/demo/wire/config"
	"k-derek.dev/demo/wire/user"

	"github.com/google/wire"
)

var RouterSet = wire.NewSet( //nolint:deadcode,unused,varcheck
	ProvideRouter,
)

func ProvideRouter(
	cfg config.Config,
	user *user.Service,
) http.Handler {
	return New(cfg, user)
}
