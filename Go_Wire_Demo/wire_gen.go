// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"k-derek.dev/demo/wire/config"
	"k-derek.dev/demo/wire/router"
	"k-derek.dev/demo/wire/user"
)

import (
	_ "github.com/joho/godotenv/autoload"
)

// Injectors from wire.go:

func InitializeApplication(cfg config.Config) (*application, error) {
	service, err := user.ProvideCache(cfg)
	if err != nil {
		return nil, err
	}
	ldapService, err := user.ProvideLDAP(cfg, service)
	if err != nil {
		return nil, err
	}
	crowdService, err := user.ProvideCROWD(cfg, service)
	if err != nil {
		return nil, err
	}
	userService, err := user.ProvideUser(ldapService, crowdService, service)
	if err != nil {
		return nil, err
	}
	handler := router.ProvideRouter(cfg, userService)
	mainApplication := newApplication(handler, userService)
	return mainApplication, nil
}