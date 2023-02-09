package events

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"k-derek.dev/demo/cqrs/internal/cqrs/commands"
)

type RoomBooked struct {
	ReservationId string    `json:"reservation_id,omitempty"`
	RoomId        string    `json:"room_id,omitempty"`
	GuestName     string    `json:"guest_name,omitempty"`
	Price         int64     `json:"price,omitempty"`
	StartDate     time.Time `json:"start_date,omitempty"`
	EndDate       time.Time `json:"end_date,omitempty"`
}

// OrderBeerOnRoomBooked is a event handler, which handles RoomBooked event and emits OrderBeer command.
type OrderBeerOnRoomBooked struct {
	CommandBus *cqrs.CommandBus
}

func (o OrderBeerOnRoomBooked) HandlerName() string {
	// this name is passed to EventsSubscriberConstructor and used to generate queue name
	return "OrderBeerOnRoomBooked"
}

func (OrderBeerOnRoomBooked) NewEvent() interface{} {
	return &RoomBooked{}
}

func (o OrderBeerOnRoomBooked) Handle(ctx context.Context, e interface{}) error {
	event := e.(*RoomBooked)

	orderBeerCmd := &commands.OrderBeer{
		RoomId: event.RoomId,
		Count:  rand.Int63n(10) + 1,
	}
	log.Printf("[OrderBeerOnRoomBooked]")
	return o.CommandBus.Send(ctx, orderBeerCmd)
}
