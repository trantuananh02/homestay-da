package model

import (
	"time"
)

// Review đại diện cho bảng review
type Review struct {
	ID         int       `db:"id" json:"id"`
	UserID     int       `db:"user_id" json:"userId"`
	HomestayID int       `db:"homestay_id" json:"homestayId"`
	BookingID  int       `db:"booking_id" json:"bookingId"`
	Rating     int       `db:"rating" json:"rating"`
	Comment    string    `db:"comment" json:"comment"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	// Thông tin bổ sung từ join
	UserName     string `db:"user_name" json:"userName,omitempty"`
	HomestayName string `db:"homestay_name" json:"homestayName,omitempty"`
}

// ReviewCreateRequest request tạo review mới
type ReviewCreateRequest struct {
	UserID     int    `json:"userId" binding:"required"`
	HomestayID int    `json:"homestayId" binding:"required"`
	BookingID  int   `json:"bookingId"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Comment    string `json:"comment"`
}

// ReviewUpdateRequest request cập nhật review
type ReviewUpdateRequest struct {
	Rating  *int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Comment *string `json:"comment"`
}

// ReviewSearchRequest request tìm kiếm review
type ReviewSearchRequest struct {
	UserID     *int `json:"userId"`
	HomestayID *int `json:"homestayId"`
	Rating     *int `json:"rating"`
	Page       int  `json:"page" binding:"min=1"`
	PageSize   int  `json:"pageSize" binding:"min=1,max=100"`
}
