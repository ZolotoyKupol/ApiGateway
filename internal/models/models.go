package models

import "time"


type Guest struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoomID    int    `json:"room_id"`
}

type Room struct {
	ID string `json:"id"` 
	Number string `json:"number"`
	Floor string `json:"floor"`
	RoomSize float64 `json:"room_size"`
	Status string `json:"status"`
	OccupiedBy string `json:"occupied_by"` 
	CheckIn time.Time `json:"check_in"`
	CheckOut time.Time `json:"check_out"`
}