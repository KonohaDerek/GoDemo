package handlers

import (
	"context"
	"log"
	"math/rand"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/pkg/errors"
	"k-derek.dev/demo/cqrs/internal/cqrs/commands"
	"k-derek.dev/demo/cqrs/internal/cqrs/events"
)

// OrderBeerHandler is a command handler, which handles OrderBeer command and emits BeerOrdered.
// BeerOrdered is not handled by any event handler, but we may use persistent Pub/Sub to handle it in the future.
type OrderBeerHandler struct {
	EventBus *cqrs.EventBus
}

func (o OrderBeerHandler) HandlerName() string {
	return "OrderBeerHandler"
}

func (o OrderBeerHandler) NewCommand() interface{} {
	return &commands.OrderBeer{}
}

func (o OrderBeerHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*commands.OrderBeer)

	if rand.Int63n(10) == 0 {
		// sometimes there is no beer left, command will be retried
		return errors.Errorf("no beer left for room %s, please try later", cmd.RoomId)
	}

	if err := o.EventBus.Publish(ctx, &events.BeerOrdered{
		RoomId: cmd.RoomId,
		Count:  cmd.Count,
	}); err != nil {
		return err
	}

	log.Printf("%d beers ordered to room %s", cmd.Count, cmd.RoomId)
	return nil
}
