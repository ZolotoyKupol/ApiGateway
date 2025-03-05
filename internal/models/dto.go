package models

import "time"

type GuestDB struct {
    ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
    FirstName string `gorm:"column:first_name;not null" json:"first_name"`
    LastName  string `gorm:"column:last_name;not null" json:"last_name"`
    RoomID    int    `gorm:"column:room_id;foreignKey:ID;references:rooms;not null" json:"room_id"`
}


func (GuestDB) TableName() string {
	return "guests"
}

type RoomDB struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`
	Number string `gorm:"column:number" json:"number"`
	Floor string `gorm:"column:floor" json:"floor"`
	RoomSize float64 `gorm:"column:room_size" json:"room_size"`
	Status string `gorm:"column:status" json:"status"`
	OccupiedBy string `gorm:"column:occupied_by" json:"occupied_by"`
	CheckIn time.Time `gorm:"column:check_in" json:"check_in"`
	CheckOut time.Time `gorm:"column:check_out" json:"check_out"`
}

func (RoomDB) TableName() string {
	return "rooms"
}
