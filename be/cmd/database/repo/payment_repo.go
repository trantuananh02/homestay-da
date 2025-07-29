package repo

import (
	"context"
	"homestay-be/cmd/database/model"
)

// PaymentRepository định nghĩa các phương thức thao tác với bảng payment
type PaymentRepository interface {
	// Create tạo payment mới
	Create(ctx context.Context, payment *model.PaymentCreateRequest) (*model.Payment, error)

	// GetByID lấy payment theo ID
	GetByID(ctx context.Context, id int) (*model.Payment, error)

	// Update cập nhật thông tin payment
	Update(ctx context.Context, id int, payment *model.PaymentUpdateRequest) (*model.Payment, error)

	// Delete xóa payment
	Delete(ctx context.Context, id int) error

	// List lấy danh sách payment với phân trang
	List(ctx context.Context, page, pageSize int) ([]*model.Payment, int, error)

	// Search tìm kiếm payment
	Search(ctx context.Context, req *model.PaymentSearchRequest) ([]*model.Payment, int, error)

	// GetByBookingID lấy danh sách payment theo booking
	GetByBookingID(ctx context.Context, bookingID int, page, pageSize int) ([]*model.Payment, int, error)

	// GetByStatus lấy danh sách payment theo status
	GetByStatus(ctx context.Context, status string, page, pageSize int) ([]*model.Payment, int, error)

	// GetByTransactionID lấy payment theo transaction ID
	GetByTransactionID(ctx context.Context, transactionID string) (*model.Payment, error)
}
