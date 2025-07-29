package repo

import (
	"context"
	"homestay-be/cmd/database/model"
	"time"
)

// BookingRepository định nghĩa các phương thức thao tác với bảng booking
type BookingRepository interface {
	// Create tạo booking mới
	Create(ctx context.Context, booking *model.BookingCreateRequest) (*model.Booking, error)

	// GetByID lấy booking theo ID
	GetByID(ctx context.Context, id int) (*model.Booking, error)

	// Update status cập nhật thông tin booking
	UpdateStatus(ctx context.Context, id int, booking *model.BookingUpdateRequest) (*model.Booking, error)

	// Update cập nhật thông tin booking
	Update(ctx context.Context, id int, bookings *model.Booking) (*model.Booking, error)

	// Delete xóa booking
	Delete(ctx context.Context, id int) error

	// List lấy danh sách booking với phân trang
	List(ctx context.Context, page, pageSize int) ([]*model.Booking, int, error)

	// Search tìm kiếm booking
	Search(ctx context.Context, req *model.BookingSearchRequest) ([]*model.Booking, int, error)

	// GetByRoom thông qua bảng booking_room
	GetByRoom(ctx context.Context, roomID int, page, pageSize int) ([]*model.Booking, int, error)

	// Lấy danh sách BookingRoom theo booking_id
	GetRoomsByBookingID(ctx context.Context, bookingID int) ([]*model.BookingRoom, error)

	// Lưu 1 bản ghi booking_room
	InsertBookingRoom(ctx context.Context, bookingRoom *model.BookingRoom) (*model.BookingRoom, error)

	// GetByUserID lấy danh sách booking theo user
	GetByUserID(ctx context.Context, userID int, page, pageSize int) ([]*model.Booking, int, error)

	// GetByStatus lấy danh sách booking theo status
	GetByStatus(ctx context.Context, status string, page, pageSize int) ([]*model.Booking, int, error)

	// GetByBookingRequestID lấy booking theo booking request ID
	GetByBookingRequestID(ctx context.Context, bookingRequestID int) (*model.Booking, error)

	// GetBookingsByHomestayID lấy danh sách booking theo homestay ID
	GetByHomestayID(ctx context.Context, homestayID int, page, pageSize int) ([]*model.Booking, int, error)

	// CheckRoomExists kiểm tra xem phòng đã tồn tại trong booking hay chưa
	CheckRoomExists(ctx context.Context, roomID int, checkIn, checkOut time.Time) (bool, error)

	// CreateReview tạo review cho booking
	CreateReview(ctx context.Context, review *model.ReviewCreateRequest) (*model.Review, error)

	// GetReviewByBookingID lấy review theo booking ID
	GetReviewByBookingID(ctx context.Context, bookingID int) (*model.Review, error)

	// UpdateReview cập nhật review cho booking
	FilterByBookingCode(ctx context.Context, userId int, bookingCode string, page, pageSize int) ([]*model.Booking, int, error)
}
