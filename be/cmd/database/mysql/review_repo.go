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

type reviewRepository struct {
	db *sqlx.DB
}

// NewReviewRepository tạo instance mới của ReviewRepository
func NewReviewRepository(db *sqlx.DB) repo.ReviewRepository {
	return &reviewRepository{db: db}
}

// Create tạo review mới
func (r *reviewRepository) Create(ctx context.Context, req *model.ReviewCreateRequest) (*model.Review, error) {
	query := `
		INSERT INTO review (user_id, homestay_id, booking_id, rating, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, homestay_id, booking_id, rating, comment, created_at
	`

	var review model.Review
	err := r.db.GetContext(ctx, &review, query, req.UserID, req.HomestayID, req.BookingID, req.Rating, req.Comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return &review, nil
}

// GetByID lấy review theo ID
func (r *reviewRepository) GetByID(ctx context.Context, id int) (*model.Review, error) {
	query := `
		SELECT r.id, r.user_id, r.homestay_id, r.booking_id, r.rating, r.comment, r.created_at,
		       u.name as user_name, h.name as homestay_name
		FROM review r
		LEFT JOIN "user" u ON r.user_id = u.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE r.id = $1
	`

	var review model.Review
	err := r.db.GetContext(ctx, &review, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	return &review, nil
}

// Update cập nhật thông tin review
func (r *reviewRepository) Update(ctx context.Context, id int, req *model.ReviewUpdateRequest) (*model.Review, error) {
	// Xây dựng query động
	query := `UPDATE review SET `
	var args []interface{}
	var setClauses []string
	argIndex := 1

	if req.Rating != nil {
		setClauses = append(setClauses, fmt.Sprintf("rating = $%d", argIndex))
		args = append(args, *req.Rating)
		argIndex++
	}

	if req.Comment != nil {
		setClauses = append(setClauses, fmt.Sprintf("comment = $%d", argIndex))
		args = append(args, *req.Comment)
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
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	// Lấy thông tin review sau khi update
	return r.GetByID(ctx, id)
}

// Delete xóa review
func (r *reviewRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM review WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("review not found")
	}

	return nil
}

// List lấy danh sách review với phân trang
func (r *reviewRepository) List(ctx context.Context, page, pageSize int) ([]*model.Review, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM review`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	// Lấy danh sách review
	offset := (page - 1) * pageSize
	query := `
		SELECT r.id, r.user_id, r.homestay_id, r.booking_id, r.rating, r.comment, r.created_at,
		       u.name as user_name, h.name as homestay_name
		FROM review r
		LEFT JOIN "user" u ON r.user_id = u.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		ORDER BY r.created_at DESC
		LIMIT $1 OFFSET $2
	`

	var reviews []*model.Review
	err = r.db.SelectContext(ctx, &reviews, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list reviews: %w", err)
	}

	return reviews, total, nil
}

// Search tìm kiếm review
func (r *reviewRepository) Search(ctx context.Context, req *model.ReviewSearchRequest) ([]*model.Review, int, error) {
	// Xây dựng query tìm kiếm
	whereClauses := []string{}
	var args []interface{}
	argIndex := 1

	if req.UserID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("r.user_id = $%d", argIndex))
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.HomestayID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("r.homestay_id = $%d", argIndex))
		args = append(args, *req.HomestayID)
		argIndex++
	}

	if req.Rating != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("r.rating = $%d", argIndex))
		args = append(args, *req.Rating)
		argIndex++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Đếm tổng số records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM review r
		LEFT JOIN "user" u ON r.user_id = u.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		%s
	`, whereClause)
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	// Lấy danh sách review
	offset := (req.Page - 1) * req.PageSize
	query := fmt.Sprintf(`
		SELECT r.id, r.user_id, r.homestay_id, r.booking_id, r.rating, r.comment, r.created_at,
		       u.name as user_name, h.name as homestay_name
		FROM review r
		LEFT JOIN "user" u ON r.user_id = u.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		%s
		ORDER BY r.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, req.PageSize, offset)

	var reviews []*model.Review
	err = r.db.SelectContext(ctx, &reviews, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search reviews: %w", err)
	}

	return reviews, total, nil
}

// GetByUserID lấy danh sách review theo user
func (r *reviewRepository) GetByUserID(ctx context.Context, userID int, page, pageSize int) ([]*model.Review, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM review WHERE user_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	// Lấy danh sách review
	offset := (page - 1) * pageSize
	query := `
		SELECT r.id, r.user_id, r.homestay_id, r.booking_id, r.rating, r.comment, r.created_at,
		       u.name as user_name, h.name as homestay_name
		FROM review r
		LEFT JOIN "user" u ON r.user_id = u.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE r.user_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var reviews []*model.Review
	err = r.db.SelectContext(ctx, &reviews, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get reviews by user: %w", err)
	}

	return reviews, total, nil
}

// GetByHomestayID lấy danh sách review theo homestay
func (r *reviewRepository) GetByHomestayID(ctx context.Context, homestayID int, page, pageSize int) ([]*model.Review, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM review WHERE homestay_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, homestayID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	// Lấy danh sách review
	offset := (page - 1) * pageSize
	query := `
		SELECT r.id, r.user_id, r.homestay_id, r.booking_id, r.rating, r.comment, r.created_at,
		       u.name as user_name, h.name as homestay_name
		FROM review r
		LEFT JOIN "user" u ON r.user_id = u.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE r.homestay_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var reviews []*model.Review
	err = r.db.SelectContext(ctx, &reviews, query, homestayID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get reviews by homestay: %w", err)
	}

	return reviews, total, nil
}

// GetByRating lấy danh sách review theo rating
func (r *reviewRepository) GetByRating(ctx context.Context, rating int, page, pageSize int) ([]*model.Review, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM review WHERE rating = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, rating)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	// Lấy danh sách review
	offset := (page - 1) * pageSize
	query := `
		SELECT r.id, r.user_id, r.homestay_id, r.booking_id, r.rating, r.comment, r.created_at,
		       u.name as user_name, h.name as homestay_name
		FROM review r
		LEFT JOIN "user" u ON r.user_id = u.id
		LEFT JOIN homestay h ON r.homestay_id = h.id
		WHERE r.rating = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var reviews []*model.Review
	err = r.db.SelectContext(ctx, &reviews, query, rating, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get reviews by rating: %w", err)
	}

	return reviews, total, nil
}

// GetAverageRating lấy rating trung bình của homestay
func (r *reviewRepository) GetAverageRating(ctx context.Context, homestayID int) (float64, error) {
	query := `
		SELECT COALESCE(AVG(rating), 0) as average_rating
		FROM review
		WHERE homestay_id = $1
	`

	var averageRating float64
	err := r.db.GetContext(ctx, &averageRating, query, homestayID)
	if err != nil {
		return 0, fmt.Errorf("failed to get average rating: %w", err)
	}

	return averageRating, nil
} 