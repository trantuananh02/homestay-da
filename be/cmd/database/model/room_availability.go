package model

import (
	"time"
)

// RoomAvailability đại diện cho bảng room_availability
type RoomAvailability struct {
	ID        int       `db:"id" json:"id"`
	RoomID    int       `db:"room_id" json:"roomId"`
	Date      time.Time `db:"date" json:"date"`
	Status    string    `db:"status" json:"status"`
	Price     *float64  `db:"price" json:"price"` // Có thể null
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// RoomAvailabilityCreateRequest request tạo room availability mới
type RoomAvailabilityCreateRequest struct {
	RoomID int       `json:"roomId" binding:"required"`
	Date   time.Time `json:"date" binding:"required"`
	Status string    `json:"status" binding:"required"`
	Price  *float64  `json:"price" binding:"omitempty,min=0"`
}

// RoomAvailabilityUpdateRequest request cập nhật room availability
type RoomAvailabilityUpdateRequest struct {
	Status *string  `json:"status"`
	Price  *float64 `json:"price" binding:"omitempty,min=0"`
}

// RoomAvailabilitySearchRequest request tìm kiếm room availability
type RoomAvailabilitySearchRequest struct {
	RoomID    *int       `json:"roomId"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	Status    *string    `json:"status"`
}

// RoomAvailabilityBatchRequest request tạo nhiều room availability cùng lúc
type RoomAvailabilityBatchRequest struct {
	RoomID    int       `json:"roomId" binding:"required"`
	StartDate time.Time `json:"startDate" binding:"required"`
	EndDate   time.Time `json:"endDate" binding:"required"`
	Status    string    `json:"status" binding:"required"`
	Price     *float64  `json:"price" binding:"omitempty,min=0"`
}
