package model

import "time"

// BookingRoom mapping với bảng booking_room
// Dùng cho booking nhiều phòng

type BookingRoom struct {
	ID        int       `db:"id" json:"id"`
	BookingID int       `db:"booking_id" json:"booking_id"`
	RoomID    int       `db:"room_id" json:"room_id"`
	RoomName  string    `db:"room_name" json:"room_name"`
	RoomType  string    `db:"room_type" json:"room_type"`
	Capacity  int       `db:"capacity" json:"capacity"`
	Price     float64   `db:"price" json:"price"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
