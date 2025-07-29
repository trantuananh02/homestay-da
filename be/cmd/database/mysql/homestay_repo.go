package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"homestay-be/cmd/database/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

type HomestayRepo struct {
	db *sqlx.DB
}

func NewHomestayRepo(db *sqlx.DB) *HomestayRepo {
	return &HomestayRepo{db: db}
}

// Create tạo homestay mới
func (r *HomestayRepo) Create(ctx context.Context, req *model.HomestayCreateRequest) (*model.Homestay, error) {
	query := `
		INSERT INTO homestay (name, description, address, city, district, ward, latitude, longitude, owner_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'inactive')
		RETURNING id, name, description, address, city, district, ward, latitude, longitude, owner_id, status, created_at, updated_at
	`

	var homestay model.Homestay
	err := r.db.QueryRowContext(ctx, query,
		req.Name, req.Description, req.Address, req.City, req.District, req.Ward,
		req.Latitude, req.Longitude, req.OwnerID,
	).Scan(
		&homestay.ID, &homestay.Name, &homestay.Description, &homestay.Address,
		&homestay.City, &homestay.District, &homestay.Ward, &homestay.Latitude,
		&homestay.Longitude, &homestay.OwnerID, &homestay.Status,
		&homestay.CreatedAt, &homestay.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating homestay: %w", err)
	}

	return &homestay, nil
}

// GetByID lấy homestay theo ID
func (r *HomestayRepo) GetByID(ctx context.Context, id int) (*model.Homestay, error) {
	query := `
		SELECT h.id, h.name, h.description, h.address, h.city, h.district, h.ward,
		       h.latitude, h.longitude, h.owner_id, h.status, h.created_at, h.updated_at,
		       u.name as owner_name
		FROM homestay h
		LEFT JOIN "user" u ON h.owner_id = u.id
		WHERE h.id = $1
	`

	var homestay model.Homestay
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&homestay.ID, &homestay.Name, &homestay.Description, &homestay.Address,
		&homestay.City, &homestay.District, &homestay.Ward, &homestay.Latitude,
		&homestay.Longitude, &homestay.OwnerID, &homestay.Status,
		&homestay.CreatedAt, &homestay.UpdatedAt, &homestay.OwnerName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("homestay not found")
		}
		return nil, fmt.Errorf("error getting homestay: %w", err)
	}

	return &homestay, nil
}

// Update cập nhật thông tin homestay
func (r *HomestayRepo) Update(ctx context.Context, id int, req *model.HomestayUpdateRequest) (*model.Homestay, error) {
	// Xây dựng query động
	query := "UPDATE homestay SET "
	var args []interface{}
	var sets []string
	paramCount := 1

	if req.Name != nil {
		sets = append(sets, fmt.Sprintf("name = $%d", paramCount))
		args = append(args, *req.Name)
		paramCount++
	}
	if req.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", paramCount))
		args = append(args, *req.Description)
		paramCount++
	}
	if req.Address != nil {
		sets = append(sets, fmt.Sprintf("address = $%d", paramCount))
		args = append(args, *req.Address)
		paramCount++
	}
	if req.City != nil {
		sets = append(sets, fmt.Sprintf("city = $%d", paramCount))
		args = append(args, *req.City)
		paramCount++
	}
	if req.District != nil {
		sets = append(sets, fmt.Sprintf("district = $%d", paramCount))
		args = append(args, *req.District)
		paramCount++
	}
	if req.Ward != nil {
		sets = append(sets, fmt.Sprintf("ward = $%d", paramCount))
		args = append(args, *req.Ward)
		paramCount++
	}
	if req.Latitude != nil {
		sets = append(sets, fmt.Sprintf("latitude = $%d", paramCount))
		args = append(args, *req.Latitude)
		paramCount++
	}
	if req.Longitude != nil {
		sets = append(sets, fmt.Sprintf("longitude = $%d", paramCount))
		args = append(args, *req.Longitude)
		paramCount++
	}
	if req.Status != nil {
		sets = append(sets, fmt.Sprintf("status = $%d", paramCount))
		args = append(args, *req.Status)
		paramCount++
	}

	if len(sets) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	sets = append(sets, "updated_at = CURRENT_TIMESTAMP")
	query += strings.Join(sets, ", ") + fmt.Sprintf(" WHERE id = $%d", paramCount)
	args = append(args, id)

	// Thực hiện update
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error updating homestay: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("homestay not found")
	}

	// Lấy homestay đã cập nhật
	return r.GetByID(ctx, id)
}

// Delete xóa homestay
func (r *HomestayRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM homestay WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting homestay: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("homestay not found")
	}

	return nil
}

// List lấy danh sách homestay với phân trang
func (r *HomestayRepo) List(ctx context.Context, page, pageSize int) ([]*model.Homestay, int, error) {
	// Đếm tổng số
	countQuery := "SELECT COUNT(*) FROM homestay"
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting homestays: %w", err)
	}

	// Lấy danh sách
	offset := (page - 1) * pageSize
	query := `
		SELECT h.id, h.name, h.description, h.address, h.city, h.district, h.ward,
		       h.latitude, h.longitude, h.owner_id, h.status, h.created_at, h.updated_at,
		       u.name as owner_name
		FROM homestay h
		LEFT JOIN "user" u ON h.owner_id = u.id
		ORDER BY h.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying homestays: %w", err)
	}
	defer rows.Close()

	var homestays []*model.Homestay
	for rows.Next() {
		var homestay model.Homestay
		err := rows.Scan(
			&homestay.ID, &homestay.Name, &homestay.Description, &homestay.Address,
			&homestay.City, &homestay.District, &homestay.Ward, &homestay.Latitude,
			&homestay.Longitude, &homestay.OwnerID, &homestay.Status,
			&homestay.CreatedAt, &homestay.UpdatedAt, &homestay.OwnerName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning homestay: %w", err)
		}
		homestays = append(homestays, &homestay)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating homestays: %w", err)
	}

	return homestays, total, nil
}

// Search tìm kiếm homestay
func (r *HomestayRepo) Search(ctx context.Context, req *model.HomestaySearchRequest) ([]*model.Homestay, int, error) {
	// Xây dựng query đếm
	countQuery := "SELECT COUNT(*) FROM homestay h WHERE 1=1"
	searchQuery := `
		SELECT h.id, h.name, h.description, h.address, h.city, h.district, h.ward,
		       h.latitude, h.longitude, h.owner_id, h.status, h.created_at, h.updated_at,
		       u.name as owner_name
		FROM homestay h
		LEFT JOIN "user" u ON h.owner_id = u.id
		WHERE 1=1
	`

	var args []interface{}
	paramCount := 1

	// Thêm điều kiện tìm kiếm
	if req.Name != nil && *req.Name != "" {
		countQuery += fmt.Sprintf(" AND h.name ILIKE $%d", paramCount)
		searchQuery += fmt.Sprintf(" AND h.name ILIKE $%d", paramCount)
		args = append(args, "%"+*req.Name+"%")
		paramCount++
	}

	if req.City != nil && *req.City != "" {
		countQuery += fmt.Sprintf(" AND h.city ILIKE $%d", paramCount)
		searchQuery += fmt.Sprintf(" AND h.city ILIKE $%d", paramCount)
		args = append(args, "%"+*req.City+"%")
		paramCount++
	}

	if req.District != nil && *req.District != "" {
		countQuery += fmt.Sprintf(" AND h.district ILIKE $%d", paramCount)
		searchQuery += fmt.Sprintf(" AND h.district ILIKE $%d", paramCount)
		args = append(args, "%"+*req.District+"%")
		paramCount++
	}

	if req.Status != nil && *req.Status != "" {
		countQuery += fmt.Sprintf(" AND h.status = $%d", paramCount)
		searchQuery += fmt.Sprintf(" AND h.status = $%d", paramCount)
		args = append(args, *req.Status)
		paramCount++
	}

	if req.OwnerID != nil && *req.OwnerID != 0 {
		countQuery += fmt.Sprintf(" AND h.owner_id = $%d", paramCount)
		searchQuery += fmt.Sprintf(" AND h.owner_id = $%d", paramCount)
		args = append(args, *req.OwnerID)
		paramCount++
	}

	// Đếm tổng số
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting homestays: %w", err)
	}

	// Thêm phân trang
	offset := (req.Page - 1) * req.PageSize
	searchQuery += fmt.Sprintf(" ORDER BY h.created_at DESC LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	args = append(args, req.PageSize, offset)

	// Thực hiện tìm kiếm
	rows, err := r.db.QueryContext(ctx, searchQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error searching homestays: %w", err)
	}
	defer rows.Close()

	var homestays []*model.Homestay
	for rows.Next() {
		var homestay model.Homestay
		err := rows.Scan(
			&homestay.ID, &homestay.Name, &homestay.Description, &homestay.Address,
			&homestay.City, &homestay.District, &homestay.Ward, &homestay.Latitude,
			&homestay.Longitude, &homestay.OwnerID, &homestay.Status,
			&homestay.CreatedAt, &homestay.UpdatedAt, &homestay.OwnerName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning homestay: %w", err)
		}
		homestays = append(homestays, &homestay)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating homestays: %w", err)
	}

	return homestays, total, nil
}

func (r *HomestayRepo) SearchAvailable(ctx context.Context, req *model.HomestaySearchRequest) ([]*model.Homestay, int, error) {
	countQuery := `
		SELECT COUNT(*) FROM homestay h
		WHERE h.status = 'active'
		AND EXISTS (
			SELECT 1 FROM room r
			WHERE r.homestay_id = h.id
			  AND r.status = 'available'
			  AND r.id NOT IN (
				SELECT br.room_id
				FROM booking_room br
				JOIN booking b ON br.booking_id = b.id
				WHERE b.status = 'confirmed'
				  AND (b.check_in, b.check_out) OVERLAPS ($1::date, $2::date)
			  )
			GROUP BY r.homestay_id
			HAVING SUM(r.capacity) >= $3
		)
	`

	searchQuery := `
		SELECT h.id, h.name, h.description, h.address, h.city, h.district, h.ward,
			   h.latitude, h.longitude, h.owner_id, h.status, h.created_at, h.updated_at,
			   u.name as owner_name
		FROM homestay h
		LEFT JOIN "user" u ON h.owner_id = u.id
		WHERE h.status = 'active'
		AND EXISTS (
			SELECT 1 FROM room r
			WHERE r.homestay_id = h.id
			  AND r.status = 'available'
			  AND r.id NOT IN (
				SELECT br.room_id
				FROM booking_room br
				JOIN booking b ON br.booking_id = b.id
				WHERE b.status = 'confirmed'
				  AND (b.check_in, b.check_out) OVERLAPS ($1::date, $2::date)
			  )
			GROUP BY r.homestay_id
			HAVING SUM(r.capacity) >= $3
		)
	`

	var args []interface{}
	args = append(args, req.CheckIn, req.CheckOut, req.GuestCount)
	paramCount := 4

	// Các điều kiện lọc bổ sung
	if req.Name != nil && *req.Name != "" {
		countQuery += fmt.Sprintf(" AND h.name ILIKE $%d", paramCount)
		searchQuery += fmt.Sprintf(" AND h.name ILIKE $%d", paramCount)
		args = append(args, "%"+*req.Name+"%")
		paramCount++
	}

	if req.City != nil && *req.City != "" {
		countQuery += fmt.Sprintf(" AND h.city ILIKE $%d", paramCount)
		searchQuery += fmt.Sprintf(" AND h.city ILIKE $%d", paramCount)
		args = append(args, "%"+*req.City+"%")
		paramCount++
	}

	if req.District != nil && *req.District != "" {
		countQuery += fmt.Sprintf(" AND h.district ILIKE $%d", paramCount)
		searchQuery += fmt.Sprintf(" AND h.district ILIKE $%d", paramCount)
		args = append(args, "%"+*req.District+"%")
		paramCount++
	}

	if req.OwnerID != nil && *req.OwnerID != 0 {
		countQuery += fmt.Sprintf(" AND h.owner_id = $%d", paramCount)
		searchQuery += fmt.Sprintf(" AND h.owner_id = $%d", paramCount)
		args = append(args, *req.OwnerID)
		paramCount++
	}

	// Đếm tổng số
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting homestays: %w", err)
	}

	// Phân trang
	offset := (req.Page - 1) * req.PageSize
	searchQuery += fmt.Sprintf(" ORDER BY h.created_at DESC LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	args = append(args, req.PageSize, offset)

	// Truy vấn dữ liệu
	rows, err := r.db.QueryContext(ctx, searchQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error searching homestays: %w", err)
	}
	defer rows.Close()

	var homestays []*model.Homestay
	for rows.Next() {
		var homestay model.Homestay
		err := rows.Scan(
			&homestay.ID, &homestay.Name, &homestay.Description, &homestay.Address,
			&homestay.City, &homestay.District, &homestay.Ward, &homestay.Latitude,
			&homestay.Longitude, &homestay.OwnerID, &homestay.Status,
			&homestay.CreatedAt, &homestay.UpdatedAt, &homestay.OwnerName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning homestay: %w", err)
		}
		homestays = append(homestays, &homestay)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating homestays: %w", err)
	}

	return homestays, total, nil
}


// GetStats lấy thống kê homestay
func (r *HomestayRepo) GetStats(ctx context.Context) (*model.HomestayStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_homestays,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_homestays,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_homestays,
			COUNT(CASE WHEN status = 'inactive' THEN 1 END) as inactive_homestays
		FROM homestay
	`

	var stats model.HomestayStats
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalHomestays, &stats.ActiveHomestays,
		&stats.PendingHomestays, &stats.InactiveHomestays,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting homestay stats: %w", err)
	}

	return &stats, nil
}

// GetByOwnerID lấy homestay theo owner ID
func (r *HomestayRepo) GetByOwnerID(ctx context.Context, ownerID int, page, pageSize int) ([]*model.Homestay, int, error) {
	// Đếm tổng số
	countQuery := "SELECT COUNT(*) FROM homestay WHERE owner_id = $1"
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, ownerID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting homestays: %w", err)
	}

	// Lấy danh sách
	offset := (page - 1) * pageSize
	query := `
		SELECT h.id, h.name, h.description, h.address, h.city, h.district, h.ward,
		       h.latitude, h.longitude, h.owner_id, h.status, h.created_at, h.updated_at,
		       u.name as owner_name
		FROM homestay h
		LEFT JOIN "user" u ON h.owner_id = u.id
		WHERE h.owner_id = $1
		ORDER BY h.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, ownerID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying homestays: %w", err)
	}
	defer rows.Close()

	var homestays []*model.Homestay
	for rows.Next() {
		var homestay model.Homestay
		err := rows.Scan(
			&homestay.ID, &homestay.Name, &homestay.Description, &homestay.Address,
			&homestay.City, &homestay.District, &homestay.Ward, &homestay.Latitude,
			&homestay.Longitude, &homestay.OwnerID, &homestay.Status,
			&homestay.CreatedAt, &homestay.UpdatedAt, &homestay.OwnerName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning homestay: %w", err)
		}
		homestays = append(homestays, &homestay)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating homestays: %w", err)
	}

	return homestays, total, nil
}

func (r *HomestayRepo) GetTopHomestays(ctx context.Context, limit int) ([]*model.Homestay, error) {
	query := `
		SELECT h.id, h.name, h.description, h.address, h.city, h.district, h.ward,
		       h.latitude, h.longitude, h.owner_id, h.status, h.created_at, h.updated_at,
		       COALESCE(SUM(r.rating), 0) AS total_rating
		FROM homestay h
		LEFT JOIN review r ON h.id = r.homestay_id
		WHERE h.status = 'active'
		GROUP BY h.id
		ORDER BY total_rating DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying top homestays: %w", err)
	}
	defer rows.Close()

	var homestays []*model.Homestay
	for rows.Next() {
		var homestay model.Homestay
		var totalRating int

		err := rows.Scan(
			&homestay.ID, &homestay.Name, &homestay.Description, &homestay.Address,
			&homestay.City, &homestay.District, &homestay.Ward, &homestay.Latitude,
			&homestay.Longitude, &homestay.OwnerID, &homestay.Status,
			&homestay.CreatedAt, &homestay.UpdatedAt,
			&totalRating,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning top homestay: %w", err)
		}

		homestays = append(homestays, &homestay)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating top homestays: %w", err)
	}

	return homestays, nil
}
