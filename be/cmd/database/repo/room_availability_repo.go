package repo

import (
	"context"
	"homestay-be/cmd/database/model"
)

// RoomAvailabilityRepository định nghĩa các phương thức thao tác với bảng room_availability
type RoomAvailabilityRepository interface {
	// Create tạo room availability mới
	Create(ctx context.Context, availability *model.RoomAvailabilityCreateRequest) (*model.RoomAvailability, error)
	
	// GetByID lấy room availability theo ID
	GetByID(ctx context.Context, id int) (*model.RoomAvailability, error)
	
	// Update cập nhật thông tin room availability
	Update(ctx context.Context, id int, availability *model.RoomAvailabilityUpdateRequest) (*model.RoomAvailability, error)
	
	// Delete xóa room availability
	Delete(ctx context.Context, id int) error
	
	// List lấy danh sách room availability với phân trang
	List(ctx context.Context, page, pageSize int) ([]*model.RoomAvailability, int, error)
	
	// Search tìm kiếm room availability
	Search(ctx context.Context, req *model.RoomAvailabilitySearchRequest) ([]*model.RoomAvailability, int, error)
	
	// GetByRoomID lấy danh sách room availability theo room
	GetByRoomID(ctx context.Context, roomID int, page, pageSize int) ([]*model.RoomAvailability, int, error)
	
	// GetByDateRange lấy room availability trong khoảng thời gian
	GetByDateRange(ctx context.Context, roomID int, startDate, endDate string) ([]*model.RoomAvailability, error)
	
	// CreateBatch tạo nhiều room availability cùng lúc
	CreateBatch(ctx context.Context, req *model.RoomAvailabilityBatchRequest) error
	
	// CheckAvailability kiểm tra tính khả dụng của room trong khoảng thời gian
	CheckAvailability(ctx context.Context, roomID int, checkIn, checkOut string) (bool, error)
} 