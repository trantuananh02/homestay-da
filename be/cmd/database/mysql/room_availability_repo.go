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

type roomAvailabilityRepository struct {
	db *sqlx.DB
}

// NewRoomAvailabilityRepository tạo instance mới của RoomAvailabilityRepository
func NewRoomAvailabilityRepository(db *sqlx.DB) repo.RoomAvailabilityRepository {
	return &roomAvailabilityRepository{db: db}
}

// Create tạo room availability mới
func (r *roomAvailabilityRepository) Create(ctx context.Context, req *model.RoomAvailabilityCreateRequest) (*model.RoomAvailability, error) {
	query := `
		INSERT INTO room_availability (room_id, date, status, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id, room_id, date, status, price
	`

	var availability model.RoomAvailability
	err := r.db.GetContext(ctx, &availability, query, req.RoomID, req.Date, req.Status, req.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to create room availability: %w", err)
	}

	return &availability, nil
}

// GetByID lấy room availability theo ID
func (r *roomAvailabilityRepository) GetByID(ctx context.Context, id int) (*model.RoomAvailability, error) {
	query := `
		SELECT id, room_id, date, status, price
		FROM room_availability
		WHERE id = $1
	`

	var availability model.RoomAvailability
	err := r.db.GetContext(ctx, &availability, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("room availability not found")
		}
		return nil, fmt.Errorf("failed to get room availability: %w", err)
	}

	return &availability, nil
}

// Update cập nhật thông tin room availability
func (r *roomAvailabilityRepository) Update(ctx context.Context, id int, req *model.RoomAvailabilityUpdateRequest) (*model.RoomAvailability, error) {
	// Xây dựng query động
	query := `UPDATE room_availability SET `
	var args []interface{}
	var setClauses []string
	argIndex := 1

	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if req.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *req.Price)
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
		return nil, fmt.Errorf("failed to update room availability: %w", err)
	}

	// Lấy thông tin room availability sau khi update
	return r.GetByID(ctx, id)
}

// Delete xóa room availability
func (r *roomAvailabilityRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM room_availability WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete room availability: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("room availability not found")
	}

	return nil
}

// List lấy danh sách room availability với phân trang
func (r *roomAvailabilityRepository) List(ctx context.Context, page, pageSize int) ([]*model.RoomAvailability, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM room_availability`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count room availabilities: %w", err)
	}

	// Lấy danh sách room availability
	offset := (page - 1) * pageSize
	query := `
		SELECT id, room_id, date, status, price
		FROM room_availability
		ORDER BY date DESC
		LIMIT $1 OFFSET $2
	`

	var availabilities []*model.RoomAvailability
	err = r.db.SelectContext(ctx, &availabilities, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list room availabilities: %w", err)
	}

	return availabilities, total, nil
}

// Search tìm kiếm room availability
func (r *roomAvailabilityRepository) Search(ctx context.Context, req *model.RoomAvailabilitySearchRequest) ([]*model.RoomAvailability, int, error) {
	// Xây dựng query tìm kiếm
	whereClauses := []string{}
	var args []interface{}
	argIndex := 1

	if req.RoomID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("room_id = $%d", argIndex))
		args = append(args, *req.RoomID)
		argIndex++
	}

	if req.StartDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("date >= $%d", argIndex))
		args = append(args, *req.StartDate)
		argIndex++
	}

	if req.EndDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("date <= $%d", argIndex))
		args = append(args, *req.EndDate)
		argIndex++
	}

	if req.Status != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Đếm tổng số records
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM room_availability %s`, whereClause)
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count room availabilities: %w", err)
	}

	// Lấy danh sách room availability
	query := fmt.Sprintf(`
		SELECT id, room_id, date, status, price
		FROM room_availability
		%s
		ORDER BY date DESC
	`, whereClause)

	var availabilities []*model.RoomAvailability
	err = r.db.SelectContext(ctx, &availabilities, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search room availabilities: %w", err)
	}

	return availabilities, total, nil
}

// GetByRoomID lấy danh sách room availability theo room
func (r *roomAvailabilityRepository) GetByRoomID(ctx context.Context, roomID int, page, pageSize int) ([]*model.RoomAvailability, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM room_availability WHERE room_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, roomID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count room availabilities: %w", err)
	}

	// Lấy danh sách room availability
	offset := (page - 1) * pageSize
	query := `
		SELECT id, room_id, date, status, price
		FROM room_availability
		WHERE room_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3
	`

	var availabilities []*model.RoomAvailability
	err = r.db.SelectContext(ctx, &availabilities, query, roomID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get room availabilities by room: %w", err)
	}

	return availabilities, total, nil
}

// GetByDateRange lấy room availability trong khoảng thời gian
func (r *roomAvailabilityRepository) GetByDateRange(ctx context.Context, roomID int, startDate, endDate string) ([]*model.RoomAvailability, error) {
	query := `
		SELECT id, room_id, date, status, price
		FROM room_availability
		WHERE room_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date ASC
	`

	var availabilities []*model.RoomAvailability
	err := r.db.SelectContext(ctx, &availabilities, query, roomID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get room availabilities by date range: %w", err)
	}

	return availabilities, nil
}

// CreateBatch tạo nhiều room availability cùng lúc
func (r *roomAvailabilityRepository) CreateBatch(ctx context.Context, req *model.RoomAvailabilityBatchRequest) error {
	// Bắt đầu transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Tạo danh sách các ngày từ startDate đến endDate
	currentDate := req.StartDate
	for currentDate.Before(req.EndDate) || currentDate.Equal(req.EndDate) {
		query := `
			INSERT INTO room_availability (room_id, date, status, price)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (room_id, date) DO UPDATE SET
			status = EXCLUDED.status,
			price = EXCLUDED.price
		`

		_, err := tx.ExecContext(ctx, query, req.RoomID, currentDate, req.Status, req.Price)
		if err != nil {
			return fmt.Errorf("failed to create room availability for date %s: %w", currentDate.Format("2006-01-02"), err)
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CheckAvailability kiểm tra tính khả dụng của room trong khoảng thời gian
func (r *roomAvailabilityRepository) CheckAvailability(ctx context.Context, roomID int, checkIn, checkOut string) (bool, error) {
	query := `
		SELECT COUNT(*) as unavailable_count
		FROM room_availability
		WHERE room_id = $1 
		AND date >= $2 
		AND date < $3 
		AND status = false
	`

	var unavailableCount int
	err := r.db.GetContext(ctx, &unavailableCount, query, roomID, checkIn, checkOut)
	if err != nil {
		return false, fmt.Errorf("failed to check room availability: %w", err)
	}

	// Nếu có ít nhất một ngày không khả dụng, thì room không khả dụng
	return unavailableCount == 0, nil
}
