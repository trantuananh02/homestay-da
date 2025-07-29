package repo

import (
	"context"
	"homestay-be/cmd/database/model"
)

// RoomRepository định nghĩa các phương thức thao tác với bảng room
type RoomRepository interface {
	// Create tạo room mới
	Create(ctx context.Context, room *model.RoomCreateRequest) (*model.Room, error)
	
	// GetByID lấy room theo ID
	GetByID(ctx context.Context, id int) (*model.Room, error)
	
	// Update cập nhật thông tin room
	Update(ctx context.Context, id int, room *model.RoomUpdateRequest) (*model.Room, error)
	
	// Delete xóa room
	Delete(ctx context.Context, id int) error
	
	// List lấy danh sách room với phân trang
	List(ctx context.Context, page, pageSize int) ([]*model.Room, int, error)
	
	// Search tìm kiếm room
	Search(ctx context.Context, req *model.RoomSearchRequest) ([]*model.Room, int, error)
	
	// GetByHomestayID lấy danh sách room theo homestay
	GetByHomestayID(ctx context.Context, homestayID int, page, pageSize int) ([]*model.Room, int, error)
	
	// GetAvailableRooms lấy danh sách room có sẵn trong khoảng thời gian
	GetAvailableRooms(ctx context.Context, homestayID int, checkIn, checkOut string, numGuests int) ([]*model.Room, error)
} 