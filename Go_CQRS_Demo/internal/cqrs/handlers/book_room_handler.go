package handlers

import (
	"context"
	"log"
	"math/rand"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"k-derek.dev/demo/cqrs/internal/cqrs/commands"
	"k-derek.dev/demo/cqrs/internal/cqrs/events"
)

// ...
// BookRoomHandler is a command handler, which handles BookRoom command and emits RoomBooked.
//
// In CQRS, one command must be handled by only one handler.
// When another handler with this command is added to command processor, error will be retuerned.
type BookRoomHandler struct {
	EventBus *cqrs.EventBus
}

func (b BookRoomHandler) HandlerName() string {
	return "BookRoomHandler"
}

// NewCommand returns type of command which this handle should handle. It must be a pointer.
func (b BookRoomHandler) NewCommand() interface{} {
	return &commands.BookRoomCmd{}
}

func (b BookRoomHandler) Handle(ctx context.Context, c interface{}) error {
	// c is always the type returned by `NewCommand`, so casting is always safe
	cmd := c.(*commands.BookRoomCmd)

	// some random price, in production you probably will calculate in wiser way
	price := (rand.Int63n(40) + 1) * 10

	log.Printf(
		"Booked %s for %s from %s to %s",
		cmd.RoomId,
		cmd.GuestName,
		cmd.StartDate,
		cmd.EndDate,
	)

	// RoomBooked will be handled by OrderBeerOnRoomBooked event handler,
	// in future RoomBooked may be handled by multiple event handler
	if err := b.EventBus.Publish(ctx, &events.RoomBooked{
		ReservationId: watermill.NewUUID(),
		RoomId:        cmd.RoomId,
		GuestName:     cmd.GuestName,
		Price:         price,
		StartDate:     cmd.StartDate,
		EndDate:       cmd.EndDate,
	}); err != nil {
		return err
	}

	return nil
}
