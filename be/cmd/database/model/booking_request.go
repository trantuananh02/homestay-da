package model

import (
	"time"
)

// BookingRequest đại diện cho bảng booking_request
type BookingRequest struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	RoomID      int       `db:"room_id" json:"room_id"`
	CheckIn     time.Time `db:"check_in" json:"check_in"`
	CheckOut    time.Time `db:"check_out" json:"check_out"`
	NumGuests   int       `db:"num_guests" json:"num_guests"`
	TotalAmount float64   `db:"total_amount" json:"total_amount"`
	Status      string    `db:"status" json:"status"`
	HostNote    string    `db:"host_note" json:"host_note"`
	GuestNote   string    `db:"guest_note" json:"guest_note"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	// Thông tin bổ sung từ join
	UserName     string `db:"user_name" json:"user_name,omitempty"`
	RoomName     string `db:"room_name" json:"room_name,omitempty"`
	HomestayName string `db:"homestay_name" json:"homestay_name,omitempty"`
}

// BookingRequestCreateRequest request tạo booking request mới
type BookingRequestCreateRequest struct {
	UserID      int       `json:"user_id" binding:"required"`
	RoomID      int       `json:"room_id" binding:"required"`
	CheckIn     time.Time `json:"check_in" binding:"required"`
	CheckOut    time.Time `json:"check_out" binding:"required"`
	NumGuests   int       `json:"num_guests" binding:"required,min=1"`
	TotalAmount float64   `json:"total_amount" binding:"required,min=0"`
	GuestNote   string    `json:"guest_note"`
}

// BookingRequestUpdateRequest request cập nhật booking request
type BookingRequestUpdateRequest struct {
	Status    *string `json:"status" binding:"omitempty,oneof=pending approved rejected cancelled"`
	HostNote  *string `json:"host_note"`
	GuestNote *string `json:"guest_note"`
}

// BookingRequestSearchRequest request tìm kiếm booking request
type BookingRequestSearchRequest struct {
	UserID    *int       `json:"user_id"`
	RoomID    *int       `json:"room_id"`
	Status    *string    `json:"status"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Page      int        `json:"page" binding:"min=1"`
	PageSize  int        `json:"page_size" binding:"min=1,max=100"`
} 