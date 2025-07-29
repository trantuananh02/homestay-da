package model

import (
	"time"
)

// Homestay đại diện cho bảng homestay
type Homestay struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Address     string    `db:"address" json:"address"`
	City        string    `db:"city" json:"city"`
	District    string    `db:"district" json:"district"`
	Ward        string    `db:"ward" json:"ward"`
	Latitude    float64   `db:"latitude" json:"latitude"`
	Longitude   float64   `db:"longitude" json:"longitude"`
	OwnerID     int       `db:"owner_id" json:"ownerId"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	// Thông tin bổ sung từ join
	OwnerName string `db:"owner_name" json:"ownerName,omitempty"`
}

// HomestayCreateRequest request tạo homestay mới
type HomestayCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Address     string  `json:"address" binding:"required"`
	City        string  `json:"city" binding:"required"`
	District    string  `json:"district" binding:"required"`
	Ward        string  `json:"ward" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	OwnerID     int     `json:"ownerId" binding:"required"`
}

// HomestayUpdateRequest request cập nhật homestay
type HomestayUpdateRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Address     *string  `json:"address"`
	City        *string  `json:"city"`
	District    *string  `json:"district"`
	Ward        *string  `json:"ward"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Status      *string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

// HomestaySearchRequest request tìm kiếm homestay
type HomestaySearchRequest struct {
	Name       *string `json:"name"`
	Address    *string `json:"address"`
	City       *string `json:"city"`
	District   *string `json:"district"`
	Status     *string `json:"status" binding:"omitempty,oneof=active inactive"`
	OwnerID    *int    `json:"ownerId"`
	CheckIn    *string `json:"checkIn" binding:"omitempty,datetime=2006-01-02"`
	CheckOut   *string `json:"checkOut" binding:"omitempty,datetime=2006-01-02"`
	GuestCount *int    `json:"guestCount"`
	Page       int     `json:"page" binding:"min=1"`
	PageSize   int     `json:"pageSize" binding:"min=1,max=100"`
}

// HomestayStats thống kê homestay
type HomestayStats struct {
	TotalHomestays    int `json:"totalHomestays"`
	ActiveHomestays   int `json:"activeHomestays"`
	PendingHomestays  int `json:"pendingHomestays"`
	InactiveHomestays int `json:"inactiveHomestays"`
}
