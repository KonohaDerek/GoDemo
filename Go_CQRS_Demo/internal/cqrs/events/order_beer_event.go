package events

type BeerOrdered struct {
	RoomId string `json:"room_id,omitempty"`
	Count  int64  `json:"count,omitempty"`
}
