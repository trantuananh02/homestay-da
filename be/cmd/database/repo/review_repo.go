package repo

import (
	"context"
	"homestay-be/cmd/database/model"
)

// ReviewRepository định nghĩa các phương thức thao tác với bảng review
type ReviewRepository interface {
	// Create tạo review mới
	Create(ctx context.Context, review *model.ReviewCreateRequest) (*model.Review, error)
	
	// GetByID lấy review theo ID
	GetByID(ctx context.Context, id int) (*model.Review, error)
	
	// Update cập nhật thông tin review
	Update(ctx context.Context, id int, review *model.ReviewUpdateRequest) (*model.Review, error)
	
	// Delete xóa review
	Delete(ctx context.Context, id int) error
	
	// List lấy danh sách review với phân trang
	List(ctx context.Context, page, pageSize int) ([]*model.Review, int, error)
	
	// Search tìm kiếm review
	Search(ctx context.Context, req *model.ReviewSearchRequest) ([]*model.Review, int, error)
	
	// GetByUserID lấy danh sách review theo user
	GetByUserID(ctx context.Context, userID int, page, pageSize int) ([]*model.Review, int, error)
	
	// GetByHomestayID lấy danh sách review theo homestay
	GetByHomestayID(ctx context.Context, homestayID int, page, pageSize int) ([]*model.Review, int, error)
	
	// GetByRating lấy danh sách review theo rating
	GetByRating(ctx context.Context, rating int, page, pageSize int) ([]*model.Review, int, error)
	
	// GetAverageRating lấy rating trung bình của homestay
	GetAverageRating(ctx context.Context, homestayID int) (float64, error)
} 