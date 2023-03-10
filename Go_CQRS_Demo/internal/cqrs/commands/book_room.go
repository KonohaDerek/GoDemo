package commands

import (
	"time"
)

type BookRoomCmd struct {
	RoomId    string    `json:"room_id,omitempty"`
	GuestName string    `json:"guest_name,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}
