package models


type Guest struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoomID    int    `json:"room_id"`
}
