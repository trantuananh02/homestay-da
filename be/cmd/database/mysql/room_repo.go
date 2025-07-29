package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"homestay-be/cmd/database/model"
	"homestay-be/cmd/database/repo"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
)

type roomRepository struct {
	db *sqlx.DB
}

// NewRoomRepository tạo instance mới của RoomRepository
func NewRoomRepository(db *sqlx.DB) repo.RoomRepository {
	return &roomRepository{db: db}
}

// Create tạo room mới
func (r *roomRepository) Create(ctx context.Context, req *model.RoomCreateRequest) (*model.Room, error) {
	// Serialize images, amenities
	imageUrls := ""
	if len(req.Images) > 0 {
		if b, err := json.Marshal(req.Images); err == nil {
			imageUrls = string(b)
		} else {
			imageUrls = strings.Join(req.Images, ",")
		}
	}
	amenities := ""
	if len(req.Amenities) > 0 {
		if b, err := json.Marshal(req.Amenities); err == nil {
			amenities = string(b)
		} else {
			amenities = strings.Join(req.Amenities, ",")
		}
	}
	query := `
		INSERT INTO room (homestay_id, name, description, type, capacity, price, price_type, status, image_urls, amenities)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	var roomID int
	err := r.db.GetContext(ctx, &roomID, query, req.HomestayID, req.Name, req.Description, req.Type, req.Capacity, req.Price, req.PriceType, "available", imageUrls, amenities)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	return r.GetByID(ctx, roomID)
}

// GetByID lấy room theo ID
func (r *roomRepository) GetByID(ctx context.Context, id int) (*model.Room, error) {
	query := `
		SELECT r.id, r.homestay_id, r.name, r.description, r.type, r.capacity, r.price, r.price_type, r.status, r.image_urls, r.amenities, h.name as homestay_name
		FROM room r
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE r.id = $1
	`
	var room model.Room

	err := r.db.GetContext(ctx, &room, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("room not found")
		}
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	logx.Info(room)

	return &room, nil
}

// Update cập nhật thông tin room
func (r *roomRepository) Update(ctx context.Context, id int, req *model.RoomUpdateRequest) (*model.Room, error) {
	// Xây dựng query động
	query := `UPDATE room SET `
	var args []interface{}
	var setClauses []string
	argIndex := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}

	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	if req.Type != nil {
		setClauses = append(setClauses, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, *req.Type)
		argIndex++
	}

	if req.Capacity != nil {
		setClauses = append(setClauses, fmt.Sprintf("capacity = $%d", argIndex))
		args = append(args, *req.Capacity)
		argIndex++
	}

	if req.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *req.Price)
		argIndex++
	}

	if req.PriceType != nil {
		setClauses = append(setClauses, fmt.Sprintf("price_type = $%d", argIndex))
		args = append(args, *req.PriceType)
		argIndex++
	}

	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	// json.Marshal(req.Images)
	if len(req.Images) > 0 {
		var images string
		if b, err := json.Marshal(req.Images); err == nil {
			images = string(b)
		} else {
			images = strings.Join(req.Images, ",")
		}
		setClauses = append(setClauses, fmt.Sprintf("image_urls = $%d", argIndex))
		args = append(args, images)
		argIndex++
	}
	if len(req.Amenities) > 0 {
		var amenities string
		if b, err := json.Marshal(req.Amenities); err == nil {
			amenities = string(b)
		} else {
			amenities = strings.Join(req.Amenities, ",")
		}
		setClauses = append(setClauses, fmt.Sprintf("amenities = $%d", argIndex))
		args = append(args, amenities)
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
		return nil, fmt.Errorf("failed to update room: %w", err)
	}

	// Lấy thông tin room sau khi update
	return r.GetByID(ctx, id)
}

// Delete xóa room
func (r *roomRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM room WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("room not found")
	}

	return nil
}

// List lấy danh sách room với phân trang
func (r *roomRepository) List(ctx context.Context, page, pageSize int) ([]*model.Room, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM room`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count rooms: %w", err)
	}

	// Lấy danh sách room
	offset := (page - 1) * pageSize
	query := `
		SELECT r.id, r.homestay_id, r.name, r.description, r.type, r.capacity, r.price, r.price_type, r.status, h.name as homestay_name
		FROM room r
		LEFT JOIN homestay h ON r.homestay_id = h.id
		ORDER BY r.id DESC
		LIMIT $1 OFFSET $2
	`

	var rooms []*model.Room
	err = r.db.SelectContext(ctx, &rooms, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list rooms: %w", err)
	}

	return rooms, total, nil
}

// Search tìm kiếm room
func (r *roomRepository) Search(ctx context.Context, req *model.RoomSearchRequest) ([]*model.Room, int, error) {
	// Xây dựng query tìm kiếm
	whereClauses := []string{}
	var args []interface{}
	argIndex := 1

	if req.HomestayID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("r.homestay_id = $%d", argIndex))
		args = append(args, *req.HomestayID)
		argIndex++
	}

	if req.Name != nil && *req.Name != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("r.name ILIKE $%d", argIndex))
		args = append(args, "%"+*req.Name+"%")
		argIndex++
	}

	if req.Type != nil && *req.Type != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("r.type = $%d", argIndex))
		args = append(args, *req.Type)
		argIndex++
	}

	if req.MinPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("r.price >= $%d", argIndex))
		args = append(args, *req.MinPrice)
		argIndex++
	}

	if req.MaxPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("r.price <= $%d", argIndex))
		args = append(args, *req.MaxPrice)
		argIndex++
	}

	if req.Capacity != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("r.capacity >= $%d", argIndex))
		args = append(args, *req.Capacity)
		argIndex++
	}

	if req.Status != nil && *req.Status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("r.status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Đếm tổng số records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM room r
		LEFT JOIN homestay h ON r.homestay_id = h.id
		%s
	`, whereClause)
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count rooms: %w", err)
	}

	// Lấy danh sách room
	offset := (req.Page - 1) * req.PageSize
	query := fmt.Sprintf(`
		SELECT r.id, r.homestay_id, r.name, r.description, r.type, r.capacity, r.price, r.price_type, r.status, r.image_urls, h.name as homestay_name
		FROM room r
		LEFT JOIN homestay h ON r.homestay_id = h.id
		%s
		ORDER BY r.id DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, req.PageSize, offset)

	var rooms []*model.Room
	err = r.db.SelectContext(ctx, &rooms, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search rooms: %w", err)
	}

	return rooms, total, nil
}

// GetByHomestayID lấy danh sách room theo homestay
func (r *roomRepository) GetByHomestayID(ctx context.Context, homestayID int, page, pageSize int) ([]*model.Room, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM room WHERE homestay_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, homestayID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count rooms: %w", err)
	}

	// Lấy danh sách room
	offset := (page - 1) * pageSize
	query := `
		SELECT r.id, r.homestay_id, r.name, r.description, r.type, r.capacity, r.price, r.price_type, r.status, h.name as homestay_name, r.image_urls
		FROM room r
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE r.homestay_id = $1
		ORDER BY r.id DESC
		LIMIT $2 OFFSET $3
	`

	var rooms []*model.Room
	err = r.db.SelectContext(ctx, &rooms, query, homestayID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get rooms by homestay: %w", err)
	}

	return rooms, total, nil
}

// GetAvailableRooms lấy danh sách room có sẵn trong khoảng thời gian
func (r *roomRepository) GetAvailableRooms(ctx context.Context, homestayID int, checkIn, checkOut string, numGuests int) ([]*model.Room, error) {
	query := `
		SELECT DISTINCT r.id, r.homestay_id, r.name, r.description, r.price, r.max_guests, r.is_active, h.name as homestay_name
		FROM room r
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE r.homestay_id = $1 
		AND r.is_active = true 
		AND r.max_guests >= $2
		AND r.id NOT IN (
			SELECT DISTINCT br.room_id
			FROM booking_request br
			WHERE br.status IN ('pending', 'approved')
			AND (
				(br.check_in <= $3 AND br.check_out >= $3) OR
				(br.check_in <= $4 AND br.check_out >= $4) OR
				(br.check_in >= $3 AND br.check_out <= $4)
			)
		)
		ORDER BY r.price ASC
	`

	var rooms []*model.Room
	err := r.db.SelectContext(ctx, &rooms, query, homestayID, numGuests, checkIn, checkOut)
	if err != nil {
		return nil, fmt.Errorf("failed to get available rooms: %w", err)
	}

	return rooms, nil
}
