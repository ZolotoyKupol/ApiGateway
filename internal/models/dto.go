package models

import "time"

type GuestDB struct {
	ID string `gorm:"primaryKey;autoIncrement"`
	FirstName string `gorm:"column:first_name"`
	LastName string `gorm:"column:last_name"`
	RoomID int `gorm:"column:room_id;foreignKey:RoomID;references:ID"`
}

func (GuestDB) TableName() string {
	return "guests"
}

type RoomDB struct {
	ID string `gorm:"primaryKey;colum:id"`
	Number string `gorm:"column:number"`
	Floor string `gorm:"column:floor"`
	RoomSize float64 `gorm:"column:room_size"`
	Status string `gorm:"column:status"`
	OccupiedBy string `gorm:"column:occupied_by"`
	CheckIn time.Time `gorm:"column:check_in"`
	CheckOut time.Time `gorm:"column:check_out"`
}

func (RoomDB) TableName() string {
	return "rooms"
}
