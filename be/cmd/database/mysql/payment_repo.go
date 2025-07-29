package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"homestay-be/cmd/database/model"
	"homestay-be/cmd/database/repo"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
)

type paymentRepository struct {
	db *sqlx.DB
}

// NewPaymentRepository tạo instance mới của PaymentRepository
func NewPaymentRepository(db *sqlx.DB) repo.PaymentRepository {
	return &paymentRepository{db: db}
}

// Create tạo payment mới
func (r *paymentRepository) Create(ctx context.Context, req *model.PaymentCreateRequest) (*model.Payment, error) {
	query := `
		INSERT INTO payment (booking_id, amount, payment_method, payment_status, transaction_id, payment_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, booking_id, amount, payment_method, payment_status, transaction_id, payment_date, created_at
	`

	var payment model.Payment
	err := r.db.GetContext(ctx, &payment, query, req.BookingID, req.Amount, req.PaymentMethod, req.PaymentStatus, req.TransactionID, req.PaymentDate)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return &payment, nil
}

// GetByID lấy payment theo ID
func (r *paymentRepository) GetByID(ctx context.Context, id int) (*model.Payment, error) {
	query := `
		SELECT p.id, p.booking_id, p.amount, p.payment_method, p.payment_status, p.transaction_id, 
		       p.payment_date, p.created_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM payment p
		LEFT JOIN booking b ON p.booking_id = b.id
		LEFT JOIN "user" u ON b.user_id = u.id
		LEFT JOIN room r ON b.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE p.id = $1
	`

	var payment model.Payment
	err := r.db.GetContext(ctx, &payment, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &payment, nil
}

// Update cập nhật thông tin payment
func (r *paymentRepository) Update(ctx context.Context, id int, req *model.PaymentUpdateRequest) (*model.Payment, error) {
	// Xây dựng query động
	query := `UPDATE payment SET `
	var args []interface{}
	var setClauses []string
	argIndex := 1

	if req.PaymentStatus != nil {
		setClauses = append(setClauses, fmt.Sprintf("payment_status = $%d", argIndex))
		args = append(args, *req.PaymentStatus)
		argIndex++
	}

	if req.TransactionID != nil {
		setClauses = append(setClauses, fmt.Sprintf("transaction_id = $%d", argIndex))
		args = append(args, *req.TransactionID)
		argIndex++
	}

	if req.PaymentDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("payment_date = $%d", argIndex))
		args = append(args, *req.PaymentDate)
		argIndex++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, id)

	// Thực hiện update
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	// Lấy thông tin payment sau khi update
	return r.GetByID(ctx, id)
}

// Delete xóa payment
func (r *paymentRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM payment WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// List lấy danh sách payment với phân trang
func (r *paymentRepository) List(ctx context.Context, page, pageSize int) ([]*model.Payment, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM payment`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	// Lấy danh sách payment
	offset := (page - 1) * pageSize
	query := `
		SELECT p.id, p.booking_id, p.amount, p.payment_method, p.payment_status, p.transaction_id, 
		       p.payment_date, p.created_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM payment p
		LEFT JOIN booking b ON p.booking_id = b.id
		LEFT JOIN "user" u ON b.user_id = u.id
		LEFT JOIN room r ON b.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	var payments []*model.Payment
	err = r.db.SelectContext(ctx, &payments, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list payments: %w", err)
	}

	return payments, total, nil
}

// Search tìm kiếm payment
func (r *paymentRepository) Search(ctx context.Context, req *model.PaymentSearchRequest) ([]*model.Payment, int, error) {
	// Xây dựng query tìm kiếm
	whereClauses := []string{}
	var args []interface{}
	argIndex := 1

	if len(req.BookingIds) > 0 {
		placeholders := []string{}
		for _, id := range req.BookingIds {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, id)
			argIndex++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("p.booking_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if req.PaymentStatus != nil && *req.PaymentStatus != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.payment_status = $%d", argIndex))
		args = append(args, *req.PaymentStatus)
		argIndex++
	}

	if req.PaymentMethod != nil && *req.PaymentMethod != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.payment_method = $%d", argIndex))
		args = append(args, *req.PaymentMethod)
		argIndex++
	}

	if req.StartDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("p.payment_date >= $%d", argIndex))
		args = append(args, *req.StartDate)
		argIndex++
	}

	if req.EndDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("p.payment_date <= $%d", argIndex))
		args = append(args, *req.EndDate)
		argIndex++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Đếm tổng số records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM payment p
		LEFT JOIN booking b ON p.booking_id = b.id
		%s
	`, whereClause)
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	logx.Info("Search payments query: ", countQuery, " args: ", args)

	// Lấy danh sách payment
	offset := (req.Page - 1) * req.PageSize
	query := fmt.Sprintf(`
		SELECT p.id, p.booking_id, b.booking_code, p.amount, p.payment_method, p.payment_status, p.transaction_id, 
		       p.payment_date, p.created_at
		FROM payment p
		LEFT JOIN booking b ON p.booking_id = b.id
		%s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, req.PageSize, offset)

	logx.Info("Search payments query: ", query, " args: ", args)

	var payments []*model.Payment
	err = r.db.SelectContext(ctx, &payments, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search payments: %w", err)
	}

	return payments, total, nil
}

// GetByBookingID lấy danh sách payment theo booking
func (r *paymentRepository) GetByBookingID(ctx context.Context, bookingID int, page, pageSize int) ([]*model.Payment, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM payment WHERE booking_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, bookingID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	// Lấy danh sách payment
	offset := (page - 1) * pageSize
	query := `
		SELECT p.id, p.booking_id, p.amount, p.payment_method, p.payment_status, p.transaction_id, 
		       p.payment_date, p.created_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM payment p
		LEFT JOIN booking b ON p.booking_id = b.id
		LEFT JOIN "user" u ON b.user_id = u.id
		LEFT JOIN room r ON b.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE p.booking_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var payments []*model.Payment
	err = r.db.SelectContext(ctx, &payments, query, bookingID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get payments by booking: %w", err)
	}

	return payments, total, nil
}

// GetByStatus lấy danh sách payment theo status
func (r *paymentRepository) GetByStatus(ctx context.Context, status string, page, pageSize int) ([]*model.Payment, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM payment WHERE payment_status = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	// Lấy danh sách payment
	offset := (page - 1) * pageSize
	query := `
		SELECT p.id, p.booking_id, p.amount, p.payment_method, p.payment_status, p.transaction_id, 
		       p.payment_date, p.created_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM payment p
		LEFT JOIN booking b ON p.booking_id = b.id
		LEFT JOIN "user" u ON b.user_id = u.id
		LEFT JOIN room r ON b.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE p.payment_status = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var payments []*model.Payment
	err = r.db.SelectContext(ctx, &payments, query, status, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get payments by status: %w", err)
	}

	return payments, total, nil
}

// GetByTransactionID lấy payment theo transaction ID
func (r *paymentRepository) GetByTransactionID(ctx context.Context, transactionID string) (*model.Payment, error) {
	query := `
		SELECT p.id, p.booking_id, p.amount, p.payment_method, p.payment_status, p.transaction_id, 
		       p.payment_date, p.created_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM payment p
		LEFT JOIN booking b ON p.booking_id = b.id
		LEFT JOIN "user" u ON b.user_id = u.id
		LEFT JOIN room r ON b.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE p.transaction_id = $1
	`

	var payment model.Payment
	err := r.db.GetContext(ctx, &payment, query, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &payment, nil
}
