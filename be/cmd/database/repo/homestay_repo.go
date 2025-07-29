package repo

import (
	"context"
	"homestay-be/cmd/database/model"
)

// HomestayRepository định nghĩa các phương thức thao tác với bảng homestay
type HomestayRepository interface {
	// Create tạo homestay mới
	Create(ctx context.Context, homestay *model.HomestayCreateRequest) (*model.Homestay, error)
	
	// GetByID lấy homestay theo ID
	GetByID(ctx context.Context, id int) (*model.Homestay, error)
	
	// Update cập nhật thông tin homestay
	Update(ctx context.Context, id int, homestay *model.HomestayUpdateRequest) (*model.Homestay, error)
	
	// Delete xóa homestay
	Delete(ctx context.Context, id int) error
	
	// List lấy danh sách homestay với phân trang
	List(ctx context.Context, page, pageSize int) ([]*model.Homestay, int, error)
	
	// Search tìm kiếm homestay
	Search(ctx context.Context, req *model.HomestaySearchRequest) ([]*model.Homestay, int, error)
	
	// GetByOwnerID lấy danh sách homestay theo owner
	GetByOwnerID(ctx context.Context, ownerID int, page, pageSize int) ([]*model.Homestay, int, error)

	// GetTopHomestays lấy danh sách homestay nổi bật
	GetTopHomestays(ctx context.Context, limit int) ([]*model.Homestay, error)

	// SearchAvailable tìm kiếm homestay có sẵn theo yêu cầu
	SearchAvailable(ctx context.Context, req *model.HomestaySearchRequest) ([]*model.Homestay, int, error)
} 