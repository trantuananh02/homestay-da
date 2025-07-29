package repo

import (
	"context"
	"homestay-be/cmd/database/model"
)

// BookingRequestRepository định nghĩa các phương thức thao tác với bảng booking_request
type BookingRequestRepository interface {
	// Create tạo booking request mới
	Create(ctx context.Context, req *model.BookingRequestCreateRequest) (*model.BookingRequest, error)
	
	// GetByID lấy booking request theo ID
	GetByID(ctx context.Context, id int) (*model.BookingRequest, error)
	
	// Update cập nhật thông tin booking request
	Update(ctx context.Context, id int, req *model.BookingRequestUpdateRequest) (*model.BookingRequest, error)
	
	// Delete xóa booking request
	Delete(ctx context.Context, id int) error
	
	// List lấy danh sách booking request với phân trang
	List(ctx context.Context, page, pageSize int) ([]*model.BookingRequest, int, error)
	
	// Search tìm kiếm booking request
	Search(ctx context.Context, req *model.BookingRequestSearchRequest) ([]*model.BookingRequest, int, error)
	
	// GetByUserID lấy danh sách booking request theo user
	GetByUserID(ctx context.Context, userID int, page, pageSize int) ([]*model.BookingRequest, int, error)
	
	// GetByRoomID lấy danh sách booking request theo room
	GetByRoomID(ctx context.Context, roomID int, page, pageSize int) ([]*model.BookingRequest, int, error)
	
	// GetByStatus lấy danh sách booking request theo status
	GetByStatus(ctx context.Context, status string, page, pageSize int) ([]*model.BookingRequest, int, error)
} 