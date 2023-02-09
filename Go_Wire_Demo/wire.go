//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"k-derek.dev/demo/wire/config"
	"k-derek.dev/demo/wire/router"
	"k-derek.dev/demo/wire/user"
)

func InitializeApplication(cfg config.Config) (*application, error) {
	wire.Build(
		router.RouterSet,
		user.UserSet,
		newApplication,
	)
	return &application{}, nil
}
