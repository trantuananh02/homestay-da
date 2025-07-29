package model

import (
	"time"
)

// Booking đại diện cho bảng booking
// Không còn trường RoomID, thay vào đó dùng bảng booking_room
// Thêm trường Rooms để join sang các phòng
// BookingRequestID, UserID giữ lại nếu cần
// Nếu không còn dùng BookingRequestID/UserID thì có thể bỏ

type Booking struct {
	ID            int           `db:"id" json:"id"`
	HomestayID    int           `db:"homestay_id" json:"homestayId"`
	BookingCode   string        `db:"booking_code" json:"bookingCode"`
	Name          string        `db:"name" json:"name"`
	Email         string        `db:"email" json:"email"`
	Phone         string        `db:"phone" json:"phone"`
	CheckIn       time.Time     `db:"check_in" json:"check_in"`
	CheckOut      time.Time     `db:"check_out" json:"check_out"`
	NumGuests     int           `db:"num_guests" json:"num_guests"`
	TotalAmount   float64       `db:"total_amount" json:"total_amount"`
	PaidAmount    float64       `db:"paid_amount" json:"paid_amount"`
	PaymentMethod string        `db:"payment_method" json:"paymentMethod"`
	Status        string        `db:"status" json:"status"`
	Rooms         []BookingRoom `json:"rooms"`
	CreatedAt     time.Time     `db:"created_at" json:"created_at"`
}

// BookingCreateRequest request tạo booking mới
// Bỏ RoomID, thêm Rooms
// Nếu cần truyền nhiều thông tin cho từng phòng, dùng []BookingRoom

type BookingCreateRequest struct {
	BookingCode   string    `db:"booking_code" json:"bookingCode"`
	HomestayID    int       `db:"homestay_id" json:"homestayId" binding:"required"`
	Name          string    `json:"name" binding:"required"`
	Email         string    `json:"email" binding:"required"`
	Phone         string    `json:"phone" binding:"required"`
	CheckIn       time.Time `json:"check_in" binding:"required"`
	CheckOut      time.Time `json:"check_out" binding:"required"`
	NumGuests     int       `json:"num_guests" binding:"required,min=1"`
	TotalAmount   float64   `db:"total_amount" json:"total_amount"`
	PaidAmount    float64   `db:"paid_amount" json:"paid_amount"`
	Status        string    `db:"status" json:"status"`
	PaymentMethod string    `json:"payment_method"`
}

type BookingRoomCreateRequest struct {
	RoomID    int     `json:"room_id" binding:"required"`
	Capacity  int     `json:"capacity" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required,min=0"`
	PriceType string  `json:"price_type" binding:"required"`
}

// BookingUpdateRequest giữ nguyên nếu chỉ update status
type BookingUpdateRequest struct {
	Status     *string  `json:"status"`
	PaidAmount *float64 `json:"paid_amount"`
}

// BookingSearchRequest request tìm kiếm booking
type BookingSearchRequest struct {
	HomestayId    *int       `form:"homestayId"`
	CustomerName  *string    `form:"customerName"`
	CustomerEmail *string    `form:"customerEmail"`
	CustomerPhone *string    `form:"customerPhone"`
	Status        *string    `form:"status"`
	StartDate     *time.Time `form:"startDate"`
	EndDate       *time.Time `form:"endDate"`
	Page          int        `form:"page"`
	PageSize      int        `form:"pageSize"`
}
