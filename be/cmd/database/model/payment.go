package model

import (
	"time"
)

// Payment đại diện cho bảng payment
type Payment struct {
	ID            int       `db:"id" json:"id"`
	BookingID     int       `db:"booking_id" json:"booking_id"`
	BookingCode   string    `db:"booking_code" json:"booking_code"`
	Amount        float64   `db:"amount" json:"amount"`
	PaymentMethod string    `db:"payment_method" json:"payment_method"`
	PaymentStatus string    `db:"payment_status" json:"payment_status"`
	TransactionID string    `db:"transaction_id" json:"transaction_id"`
	PaymentDate   time.Time `db:"payment_date" json:"payment_date"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	// Thông tin bổ sung từ join
	UserName     string `db:"user_name" json:"user_name,omitempty"`
	RoomName     string `db:"room_name" json:"room_name,omitempty"`
	HomestayName string `db:"homestay_name" json:"homestay_name,omitempty"`
}

// PaymentCreateRequest request tạo payment mới
type PaymentCreateRequest struct {
	BookingID     int       `json:"booking_id" binding:"required"`
	Amount        float64   `json:"amount" binding:"required,min=0"`
	PaymentMethod string    `json:"payment_method" binding:"required,oneof=cash bank_transfer credit_card"`
	PaymentStatus string    `json:"payment_status" binding:"required,oneof=pending completed failed refunded"`
	TransactionID string    `json:"transaction_id"`
	PaymentDate   time.Time `json:"payment_date"`
}

// PaymentUpdateRequest request cập nhật payment
type PaymentUpdateRequest struct {
	PaymentStatus *string    `json:"payment_status" binding:"omitempty,oneof=pending completed failed refunded"`
	TransactionID *string    `json:"transaction_id"`
	PaymentDate   *time.Time `json:"payment_date"`
}

// PaymentSearchRequest request tìm kiếm payment
type PaymentSearchRequest struct {
	BookingIds    []int      `json:"booking_ids"`
	PaymentStatus *string    `json:"payment_status"`
	PaymentMethod *string    `json:"payment_method"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	Page          int        `json:"page" binding:"min=1"`
	PageSize      int        `json:"page_size" binding:"min=1,max=100"`
}
