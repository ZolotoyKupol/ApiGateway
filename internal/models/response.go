package models

import "time"

type GuestResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoomID    int    `json:"room_id"`
}

type RoomResponse struct {
	ID         int       `json:"id"`
	Number     string    `json:"number"`
	Floor      string    `json:"floor"`
	RoomSize   float64   `json:"room_size"`
	Status     string    `json:"status"`
	OccupiedBy string    `json:"occupied_by"`
	CheckIn    time.Time `json:"check_in"`
	CheckOut   time.Time `json:"check_out"`
}

func (g *GuestDB) ConvertToGuestResponse() GuestResponse {
	return GuestResponse{
		ID:        g.ID,
		FirstName: g.FirstName,
		LastName:  g.LastName,
		RoomID:    g.RoomID,
	}
}

func (r *RoomDB) ConvertToRoomResponse() RoomResponse {
	return RoomResponse{
		ID:         r.ID,
		Number:     r.Number,
		Floor:      r.Floor,
		RoomSize:   r.RoomSize,
		Status:     r.Status,
		OccupiedBy: r.OccupiedBy,
		CheckIn:    r.CheckIn,
		CheckOut:   r.CheckOut,
	}
}

func ConvertToGuestResponseList(guests []GuestDB) []GuestResponse {
	guestResponses := make([]GuestResponse, len(guests))
	for i, guest := range guests {
		guestResponses[i] = guest.ConvertToGuestResponse()
	}
	return guestResponses
}

func ConvertToRoomResponseList(rooms []RoomDB) []RoomResponse {
	roomResponses := make([]RoomResponse, len(rooms))
	for i, room := range rooms {
		roomResponses[i] = room.ConvertToRoomResponse()
	}
	return roomResponses
}
