package repo

import (
	"context"
	"homestay-be/cmd/database/model"
)

// UserRepository định nghĩa các phương thức thao tác với bảng user
type UserRepository interface {
	// Create tạo user mới
	Create(ctx context.Context, user *model.UserCreateRequest) (*model.User, error)
	
	// GetByID lấy user theo ID
	GetByID(ctx context.Context, id int) (*model.User, error)
	
	// GetByEmail lấy user theo email
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	
	// Update cập nhật thông tin user
	Update(ctx context.Context, id int, user *model.UserUpdateRequest) (*model.User, error)
	
	// Delete xóa user
	Delete(ctx context.Context, id int) error
	
	// List lấy danh sách user với phân trang
	List(ctx context.Context, page, pageSize int) ([]*model.User, int, error)
	
	// Search tìm kiếm user
	Search(ctx context.Context, name, email, role string, page, pageSize int) ([]*model.User, int, error)
} 