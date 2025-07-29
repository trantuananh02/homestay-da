package model

import (
	"time"
)

// Room đại diện cho bảng room
type Room struct {
	ID          int       `db:"id" json:"id"`
	HomestayID  int       `db:"homestay_id" json:"homestayId"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Type        string    `db:"type" json:"type"`
	Capacity    int       `db:"capacity" json:"capacity"`
	Price       float64   `db:"price" json:"price"`
	PriceType   string    `db:"price_type" json:"priceType"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	Images      string  `db:"image_urls" json:"images"`
	Amenities   string  `db:"amenities" json:"amenities"`
	// Thông tin bổ sung từ join
	HomestayName string `db:"homestay_name" json:"homestayName,omitempty"`
}

// RoomCreateRequest request tạo room mới
type RoomCreateRequest struct {
	HomestayID  int      `json:"homestayId" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Type        string   `json:"type" binding:"required"`
	Capacity    int      `json:"capacity" binding:"required,min=1"`
	Price       float64  `json:"price" binding:"required,min=0"`
	PriceType   string   `json:"priceType" binding:"required"`
	Images      []string `json:"images"`
	Amenities   []string `json:"amenities"`
}

// RoomUpdateRequest request cập nhật room
type RoomUpdateRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Type        *string  `json:"type"`
	Capacity    *int     `json:"capacity" binding:"omitempty,min=1"`
	Price       *float64 `json:"price" binding:"omitempty,min=0"`
	PriceType   *string  `json:"priceType"`
	Status      *string  `json:"status"`
	Images      []string   `json:"images,omitempty"`
	Amenities   []string   `json:"amenities,omitempty"`
}

// RoomSearchRequest request tìm kiếm room
type RoomSearchRequest struct {
	HomestayID *int     `json:"homestayId"`
	Name       *string  `json:"name"`
	Type       *string  `json:"type"`
	MinPrice   *float64 `json:"minPrice"`
	MaxPrice   *float64 `json:"maxPrice"`
	Capacity   *int     `json:"capacity"`
	Status     *string  `json:"status"`
	Page       int      `json:"page" binding:"min=1"`
	PageSize   int      `json:"pageSize" binding:"min=1,max=100"`
}
