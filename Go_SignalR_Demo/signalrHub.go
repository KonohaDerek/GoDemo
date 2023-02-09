package main

import (
	"github.com/philippseith/signalr"
	"github.com/rs/zerolog/log"
)

type AppHub struct {
	signalr.Hub
}

func (h *AppHub) SendChatMessage(message string) {
	h.Clients().All().Send("chatMessageReceived", message)
}

func (h *AppHub) OnConnected(connectionID string) {
	// fmt.Printf("%s connected\n", connectionID)
	log.Info().Msgf("user connection : %s", connectionID)
}

func (h *AppHub) OnDisconnected(connectionID string) {
	// fmt.Printf("%s disconnected\n", connectionID)
	log.Info().Msgf("%s disconnected", connectionID)
}
