package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"homestay-be/cmd/database/model"
	"homestay-be/cmd/database/repo"
	"strings"

	"github.com/jmoiron/sqlx"
)

type bookingRequestRepository struct {
	db *sqlx.DB
}

// NewBookingRequestRepository tạo instance mới của BookingRequestRepository
func NewBookingRequestRepository(db *sqlx.DB) repo.BookingRequestRepository {
	return &bookingRequestRepository{db: db}
}

// Create tạo booking request mới
func (r *bookingRequestRepository) Create(ctx context.Context, req *model.BookingRequestCreateRequest) (*model.BookingRequest, error) {
	query := `
		INSERT INTO booking_request (user_id, room_id, check_in, check_out, num_guests, total_amount, guest_note)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, room_id, check_in, check_out, num_guests, total_amount, status, host_note, guest_note, created_at, updated_at
	`

	var bookingRequest model.BookingRequest
	err := r.db.GetContext(ctx, &bookingRequest, query, req.UserID, req.RoomID, req.CheckIn, req.CheckOut, req.NumGuests, req.TotalAmount, req.GuestNote)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking request: %w", err)
	}

	return &bookingRequest, nil
}

// GetByID lấy booking request theo ID
func (r *bookingRequestRepository) GetByID(ctx context.Context, id int) (*model.BookingRequest, error) {
	query := `
		SELECT br.id, br.user_id, br.room_id, br.check_in, br.check_out, br.num_guests, br.total_amount, 
		       br.status, br.host_note, br.guest_note, br.created_at, br.updated_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM booking_request br
		LEFT JOIN "user" u ON br.user_id = u.id
		LEFT JOIN room r ON br.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE br.id = $1
	`

	var bookingRequest model.BookingRequest
	err := r.db.GetContext(ctx, &bookingRequest, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking request not found")
		}
		return nil, fmt.Errorf("failed to get booking request: %w", err)
	}

	return &bookingRequest, nil
}

// Update cập nhật thông tin booking request
func (r *bookingRequestRepository) Update(ctx context.Context, id int, req *model.BookingRequestUpdateRequest) (*model.BookingRequest, error) {
	// Xây dựng query động
	query := `UPDATE booking_request SET `
	var args []interface{}
	var setClauses []string
	argIndex := 1

	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if req.HostNote != nil {
		setClauses = append(setClauses, fmt.Sprintf("host_note = $%d", argIndex))
		args = append(args, *req.HostNote)
		argIndex++
	}

	if req.GuestNote != nil {
		setClauses = append(setClauses, fmt.Sprintf("guest_note = $%d", argIndex))
		args = append(args, *req.GuestNote)
		argIndex++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	// Thêm updated_at
	setClauses = append(setClauses, fmt.Sprintf("updated_at = CURRENT_TIMESTAMP"))

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, id)

	// Thực hiện update
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update booking request: %w", err)
	}

	// Lấy thông tin booking request sau khi update
	return r.GetByID(ctx, id)
}

// Delete xóa booking request
func (r *bookingRequestRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM booking_request WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete booking request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking request not found")
	}

	return nil
}

// List lấy danh sách booking request với phân trang
func (r *bookingRequestRepository) List(ctx context.Context, page, pageSize int) ([]*model.BookingRequest, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM booking_request`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count booking requests: %w", err)
	}

	// Lấy danh sách booking request
	offset := (page - 1) * pageSize
	query := `
		SELECT br.id, br.user_id, br.room_id, br.check_in, br.check_out, br.num_guests, br.total_amount, 
		       br.status, br.host_note, br.guest_note, br.created_at, br.updated_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM booking_request br
		LEFT JOIN "user" u ON br.user_id = u.id
		LEFT JOIN room r ON br.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		ORDER BY br.created_at DESC
		LIMIT $1 OFFSET $2
	`

	var bookingRequests []*model.BookingRequest
	err = r.db.SelectContext(ctx, &bookingRequests, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list booking requests: %w", err)
	}

	return bookingRequests, total, nil
}

// Search tìm kiếm booking request
func (r *bookingRequestRepository) Search(ctx context.Context, req *model.BookingRequestSearchRequest) ([]*model.BookingRequest, int, error) {
	// Xây dựng query tìm kiếm
	whereClauses := []string{}
	var args []interface{}
	argIndex := 1

	if req.UserID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("br.user_id = $%d", argIndex))
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.RoomID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("br.room_id = $%d", argIndex))
		args = append(args, *req.RoomID)
		argIndex++
	}

	if req.Status != nil && *req.Status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("br.status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if req.StartDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("br.check_in >= $%d", argIndex))
		args = append(args, *req.StartDate)
		argIndex++
	}

	if req.EndDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("br.check_out <= $%d", argIndex))
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
		FROM booking_request br
		LEFT JOIN "user" u ON br.user_id = u.id
		LEFT JOIN room r ON br.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		%s
	`, whereClause)
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count booking requests: %w", err)
	}

	// Lấy danh sách booking request
	offset := (req.Page - 1) * req.PageSize
	query := fmt.Sprintf(`
		SELECT br.id, br.user_id, br.room_id, br.check_in, br.check_out, br.num_guests, br.total_amount, 
		       br.status, br.host_note, br.guest_note, br.created_at, br.updated_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM booking_request br
		LEFT JOIN "user" u ON br.user_id = u.id
		LEFT JOIN room r ON br.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		%s
		ORDER BY br.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, req.PageSize, offset)

	var bookingRequests []*model.BookingRequest
	err = r.db.SelectContext(ctx, &bookingRequests, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search booking requests: %w", err)
	}

	return bookingRequests, total, nil
}

// GetByUserID lấy danh sách booking request theo user
func (r *bookingRequestRepository) GetByUserID(ctx context.Context, userID int, page, pageSize int) ([]*model.BookingRequest, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM booking_request WHERE user_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count booking requests: %w", err)
	}

	// Lấy danh sách booking request
	offset := (page - 1) * pageSize
	query := `
		SELECT br.id, br.user_id, br.room_id, br.check_in, br.check_out, br.num_guests, br.total_amount, 
		       br.status, br.host_note, br.guest_note, br.created_at, br.updated_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM booking_request br
		LEFT JOIN "user" u ON br.user_id = u.id
		LEFT JOIN room r ON br.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE br.user_id = $1
		ORDER BY br.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var bookingRequests []*model.BookingRequest
	err = r.db.SelectContext(ctx, &bookingRequests, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get booking requests by user: %w", err)
	}

	return bookingRequests, total, nil
}

// GetByRoomID lấy danh sách booking request theo room
func (r *bookingRequestRepository) GetByRoomID(ctx context.Context, roomID int, page, pageSize int) ([]*model.BookingRequest, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM booking_request WHERE room_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, roomID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count booking requests: %w", err)
	}

	// Lấy danh sách booking request
	offset := (page - 1) * pageSize
	query := `
		SELECT br.id, br.user_id, br.room_id, br.check_in, br.check_out, br.num_guests, br.total_amount, 
		       br.status, br.host_note, br.guest_note, br.created_at, br.updated_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM booking_request br
		LEFT JOIN "user" u ON br.user_id = u.id
		LEFT JOIN room r ON br.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE br.room_id = $1
		ORDER BY br.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var bookingRequests []*model.BookingRequest
	err = r.db.SelectContext(ctx, &bookingRequests, query, roomID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get booking requests by room: %w", err)
	}

	return bookingRequests, total, nil
}

// GetByStatus lấy danh sách booking request theo status
func (r *bookingRequestRepository) GetByStatus(ctx context.Context, status string, page, pageSize int) ([]*model.BookingRequest, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM booking_request WHERE status = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count booking requests: %w", err)
	}

	// Lấy danh sách booking request
	offset := (page - 1) * pageSize
	query := `
		SELECT br.id, br.user_id, br.room_id, br.check_in, br.check_out, br.num_guests, br.total_amount, 
		       br.status, br.host_note, br.guest_note, br.created_at, br.updated_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM booking_request br
		LEFT JOIN "user" u ON br.user_id = u.id
		LEFT JOIN room r ON br.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE br.status = $1
		ORDER BY br.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var bookingRequests []*model.BookingRequest
	err = r.db.SelectContext(ctx, &bookingRequests, query, status, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get booking requests by status: %w", err)
	}

	return bookingRequests, total, nil
} 