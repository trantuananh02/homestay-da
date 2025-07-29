package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"homestay-be/cmd/database/model"
	"homestay-be/cmd/database/repo"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
)

type bookingRepository struct {
	db *sqlx.DB
}

// NewBookingRepository tạo instance mới của BookingRepository
func NewBookingRepository(db *sqlx.DB) repo.BookingRepository {
	return &bookingRepository{db: db}
}

// Create tạo booking mới
func (r *bookingRepository) Create(ctx context.Context, req *model.BookingCreateRequest) (*model.Booking, error) {
	logx.Info(req)

	query := `
		INSERT INTO booking (name, email, phone, check_in, check_out, num_guests, total_amount, status, created_at, booking_code, paid_amount, payment_method, homestay_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), $9, $10, $11, $12)
		RETURNING id, name, email, phone, check_in, check_out, num_guests, total_amount, status, created_at, booking_code, paid_amount, payment_method, homestay_id
	`

	var booking model.Booking
	err := r.db.GetContext(ctx, &booking, query, req.Name, req.Email, req.Phone, req.CheckIn, req.CheckOut, req.NumGuests, req.TotalAmount, req.Status, req.BookingCode, req.PaidAmount, req.PaymentMethod, req.HomestayID)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

// GetByID lấy booking theo ID
func (r *bookingRepository) GetByID(ctx context.Context, id int) (*model.Booking, error) {
	query := `
		SELECT b.id, b.email, b.name, b.phone, b.check_in, b.check_out, b.num_guests, b.total_amount, b.status, b.created_at, b.booking_code, b.paid_amount, b.payment_method
		FROM booking b
		WHERE b.id = $1
	`

	var booking model.Booking
	err := r.db.GetContext(ctx, &booking, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &booking, nil
}

// Update cập nhật thông tin booking
func (r *bookingRepository) UpdateStatus(ctx context.Context, id int, req *model.BookingUpdateRequest) (*model.Booking, error) {
	// Xây dựng query động
	query := `UPDATE booking SET `
	var args []interface{}
	var setClauses []string
	argIndex := 1

	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if req.PaidAmount != nil {
		setClauses = append(setClauses, fmt.Sprintf("paid_amount = $%d", argIndex))
		args = append(args, *req.PaidAmount)
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
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	// Lấy thông tin booking sau khi update
	return r.GetByID(ctx, id)
}

// Update cập nhật thông tin booking
func (r *bookingRepository) Update(ctx context.Context, id int, booking *model.Booking) (*model.Booking, error) {
	query := `
		UPDATE booking
		SET name = $1, email = $2, phone = $3, check_in = $4, check_out = $5, num_guests = $6, total_amount = $7, status = $8, updated_at = NOW()
		WHERE id = $9
		RETURNING id, name, email, phone, check_in, check_out, num_guests, total_amount, status, created_at, booking_code, paid_amount, payment_method
	`
	var updatedBooking model.Booking

	err := r.db.GetContext(ctx, &updatedBooking, query, booking.Name, booking.Email, booking.Phone, booking.CheckIn, booking.CheckOut, booking.NumGuests, booking.TotalAmount, booking.Status, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}
	return &updatedBooking, nil
}

// Delete xóa booking
func (r *bookingRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM booking WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete booking: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

// List lấy danh sách booking với phân trang
func (r *bookingRepository) List(ctx context.Context, page, pageSize int) ([]*model.Booking, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM booking`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	// Lấy danh sách booking
	offset := (page - 1) * pageSize
	query := `
		SELECT b.id, b.booking_request_id, b.user_id, b.room_id, b.check_in, b.check_out, b.num_guests, 
		       b.total_amount, b.status, b.created_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM booking b
		LEFT JOIN "user" u ON b.user_id = u.id
		LEFT JOIN room r ON b.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		ORDER BY b.created_at DESC
		LIMIT $1 OFFSET $2
	`

	var bookings []*model.Booking
	err = r.db.SelectContext(ctx, &bookings, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list bookings: %w", err)
	}

	return bookings, total, nil
}

func (r *bookingRepository) Search(ctx context.Context, req *model.BookingSearchRequest) ([]*model.Booking, int, error) {
	whereClauses := []string{}
	var args []interface{}
	argIndex := 1

	if req.Status != nil && *req.Status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("b.status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if req.StartDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("b.check_in >= $%d", argIndex))
		args = append(args, *req.StartDate)
		argIndex++
	}

	if req.EndDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("b.check_out <= $%d", argIndex))
		args = append(args, *req.EndDate)
		argIndex++
	}

	if req.CustomerName != nil && *req.CustomerName != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("b.name ILIKE $%d", argIndex))
		args = append(args, "%"+*req.CustomerName+"%")
		argIndex++
	}

	if req.CustomerPhone != nil && *req.CustomerPhone != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("b.phone ILIKE $%d", argIndex))
		args = append(args, "%"+*req.CustomerPhone+"%")
		argIndex++
	}

	if req.CustomerEmail != nil && *req.CustomerEmail != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("b.email = $%d", argIndex))
		args = append(args, *req.CustomerEmail)
		argIndex++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Đếm tổng số records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM booking b
		%s
	`, whereClause)

	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	// Truy vấn danh sách bookings
	query := fmt.Sprintf(`
		SELECT b.id, b.email, b.name, b.phone, b.check_in, b.check_out, b.num_guests, 
		       b.total_amount, b.status, b.created_at, b.booking_code, 
		       b.paid_amount, b.payment_method
		FROM booking b
		%s
		ORDER BY b.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, req.PageSize, (req.Page-1)*req.PageSize)

	var bookings []*model.Booking
	err = r.db.SelectContext(ctx, &bookings, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search bookings: %w", err)
	}

	return bookings, total, nil
}

// GetByUserID lấy danh sách booking theo user
func (r *bookingRepository) GetByUserID(ctx context.Context, userID int, page, pageSize int) ([]*model.Booking, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM booking WHERE user_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	// Lấy danh sách booking
	offset := (page - 1) * pageSize
	query := `
		SELECT b.id, b.booking_request_id, b.user_id, b.room_id, b.check_in, b.check_out, b.num_guests, 
		       b.total_amount, b.status, b.created_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM booking b
		LEFT JOIN "user" u ON b.user_id = u.id
		LEFT JOIN room r ON b.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var bookings []*model.Booking
	err = r.db.SelectContext(ctx, &bookings, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get bookings by user: %w", err)
	}

	return bookings, total, nil
}

// GetByRoom lấy booking qua bảng booking_room
func (r *bookingRepository) GetByRoom(ctx context.Context, roomID int, page, pageSize int) ([]*model.Booking, int, error) {
	countQuery := `SELECT COUNT(DISTINCT b.id)
		FROM booking b
		JOIN booking_room br ON b.id = br.booking_id
		WHERE br.room_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, roomID)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := `SELECT b.*
	FROM booking b
	JOIN booking_room br ON b.id = br.booking_id
	WHERE br.room_id = $1
	ORDER BY b.created_at DESC
	LIMIT $2 OFFSET $3`
	var bookings []*model.Booking
	err = r.db.SelectContext(ctx, &bookings, query, roomID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	return bookings, total, nil
}

// GetRoomsByBookingID lấy danh sách BookingRoom theo booking_id
func (r *bookingRepository) GetRoomsByBookingID(ctx context.Context, bookingID int) ([]*model.BookingRoom, error) {
	query := `SELECT * FROM booking_room WHERE booking_id = $1`
	var rooms []*model.BookingRoom
	err := r.db.SelectContext(ctx, &rooms, query, bookingID)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetByStatus lấy danh sách booking theo status
func (r *bookingRepository) GetByStatus(ctx context.Context, status string, page, pageSize int) ([]*model.Booking, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM booking WHERE status = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	// Lấy danh sách booking
	offset := (page - 1) * pageSize
	query := `
		SELECT b.id, b.booking_request_id, b.user_id, b.room_id, b.check_in, b.check_out, b.num_guests, 
		       b.total_amount, b.status, b.created_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM booking b
		LEFT JOIN "user" u ON b.user_id = u.id
		LEFT JOIN room r ON b.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE b.status = $1
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var bookings []*model.Booking
	err = r.db.SelectContext(ctx, &bookings, query, status, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get bookings by status: %w", err)
	}

	return bookings, total, nil
}

// GetByBookingRequestID lấy booking theo booking request ID
func (r *bookingRepository) GetByBookingRequestID(ctx context.Context, bookingRequestID int) (*model.Booking, error) {
	query := `
		SELECT b.id, b.booking_request_id, b.user_id, b.room_id, b.check_in, b.check_out, b.num_guests, 
		       b.total_amount, b.status, b.created_at,
		       u.name as user_name, r.name as room_name, h.name as homestay_name
		FROM booking b
		LEFT JOIN "user" u ON b.user_id = u.id
		LEFT JOIN room r ON b.room_id = r.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE b.booking_request_id = $1
	`

	var booking model.Booking
	err := r.db.GetContext(ctx, &booking, query, bookingRequestID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &booking, nil
}

// InsertBookingRoom lưu 1 bản ghi booking_room
func (r *bookingRepository) InsertBookingRoom(ctx context.Context, bookingRoom *model.BookingRoom) (*model.BookingRoom, error) {
	query := `INSERT INTO booking_room (booking_id, room_id, room_name, room_type, capacity, price, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, booking_id, room_id, room_name, room_type, capacity, price, created_at`
	var br model.BookingRoom
	err := r.db.GetContext(ctx, &br, query, bookingRoom.BookingID, bookingRoom.RoomID, bookingRoom.RoomName, bookingRoom.RoomType, bookingRoom.Capacity, bookingRoom.Price, bookingRoom.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &br, nil
}

// GetBookingsByHomestayID lấy danh sách booking theo homestay ID
func (r *bookingRepository) GetByHomestayID(ctx context.Context, homestayID int, page, pageSize int) ([]*model.Booking, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM booking b
		JOIN booking_room br ON b.id = br.booking_id
		JOIN room r ON br.room_id = r.id
		WHERE r.homestay_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, homestayID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	// Lấy danh sách booking
	offset := (page - 1) * pageSize
	query := `
		SELECT b.id, b.booking_code, b.name, b.phone, b.email, b.check_in, b.check_out, b.num_guests,
			b.total_amount, b.paid_amount, b.status, b.created_at, b.payment_method
		FROM booking b
		JOIN booking_room br ON b.id = br.booking_id
		JOIN room r ON br.room_id = r.id
		JOIN homestay h ON r.homestay_id = h.id
		WHERE r.homestay_id = $1
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var bookings []*model.Booking
	err = r.db.SelectContext(ctx, &bookings, query, homestayID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get bookings by homestay ID: %w", err)
	}

	return bookings, total, nil
}

// CheckRoomExists kiểm tra xem phòng đã tồn tại trong booking hay chưa
func (r *bookingRepository) CheckRoomExists(ctx context.Context, roomID int, checkIn, checkOut time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM booking b
		JOIN booking_room br ON b.id = br.booking_id
		WHERE br.room_id = $1 AND (
			(b.check_in < $2 AND b.check_out > $2) OR 
			(b.check_in < $3 AND b.check_out > $3) OR 
			(b.check_in >= $2 AND b.check_out <= $3)
		)
	`

	var count int
	err := r.db.GetContext(ctx, &count, query, roomID, checkIn, checkOut)
	if err != nil {
		return false, fmt.Errorf("failed to check room existence: %w", err)
	}

	return count > 0, nil
}

// CreateReview tạo review cho booking
func (r *bookingRepository) CreateReview(ctx context.Context, review *model.ReviewCreateRequest) (*model.Review, error) {
	query := `INSERT INTO review (booking_id, user_id, homestay_id, comment, rating, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, booking_id, user_id, homestay_id, comment, rating, created_at`
	var rv model.Review
	err := r.db.GetContext(ctx, &rv, query, review.BookingID, review.UserID, review.HomestayID, review.Comment, review.Rating, time.Now())
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// GetReviewByBookingID lấy review theo booking ID
func (r *bookingRepository) GetReviewByBookingID(ctx context.Context, bookingID int) (*model.Review, error) {
	query := `SELECT * FROM review WHERE booking_id = $1 LIMIT 1`
	var review model.Review
	err := r.db.GetContext(ctx, &review, query, bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review not found for booking ID %d", bookingID)
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}
	return &review, nil
}

// FilterByBookingCode lọc booking theo booking code
func (r *bookingRepository) FilterByBookingCode(ctx context.Context, userId int, bookingCode string, page, pageSize int) ([]*model.Booking, int, error) {
	query := `
		SELECT b.id, b.booking_code, b.homestay_id, b.email, b.name, b.phone,
       		b.check_in, b.check_out, b.num_guests, b.total_amount,
       		b.status, b.created_at, b.payment_method, b.paid_amount
		FROM booking b
		JOIN homestay h ON b.homestay_id = h.id
		WHERE b.booking_code like $1
		AND (h.owner_id = $2)
		ORDER BY b.created_at DESC
		LIMIT $3 OFFSET $4
	`

	offset := (page - 1) * pageSize
	var bookings []*model.Booking
	err := r.db.SelectContext(ctx, &bookings, query, "%"+bookingCode+"%", userId, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to filter bookings by code: %w", err)
	}

	countQuery := `
		SELECT COUNT(*)
		FROM booking b
		JOIN homestay h ON b.homestay_id = h.id
		WHERE b.booking_code like $1
		AND (h.owner_id = $2)
	`
	var total int
	err = r.db.GetContext(ctx, &total, countQuery, "%"+bookingCode+"%", userId)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	return bookings, total, nil
}