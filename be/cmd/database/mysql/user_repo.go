package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"homestay-be/cmd/database/model"
	"homestay-be/cmd/database/repo"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository tạo instance mới của UserRepository
func NewUserRepository(db *sqlx.DB) repo.UserRepository {
	return &userRepository{db: db}
}

// Create tạo user mới
func (r *userRepository) Create(ctx context.Context, req *model.UserCreateRequest) (*model.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO "user" (name, email, phone, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, phone, email, password, role, created_at
	`

	var user model.User
	err = r.db.GetContext(ctx, &user, query, req.Name, req.Email, req.Phone, string(hashedPassword), req.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// GetByID lấy user theo ID
func (r *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `
		SELECT id, name, phone, email, password, role, created_at
		FROM "user"
		WHERE id = $1
	`

	var user model.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByEmail lấy user theo email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, name, phone, email, password, role, created_at
		FROM "user"
		WHERE email = $1
	`

	var user model.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// Update cập nhật thông tin user
func (r *userRepository) Update(ctx context.Context, id int, req *model.UserUpdateRequest) (*model.User, error) {
	// Xây dựng query động
	query := `UPDATE "user" SET `
	var args []interface{}
	var setClauses []string
	argIndex := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}

	if req.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, *req.Email)
		argIndex++
	}

	// phone
	if req.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, *req.Phone)
		argIndex++
	}

	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		setClauses = append(setClauses, fmt.Sprintf("password = $%d", argIndex))
		args = append(args, string(hashedPassword))
		argIndex++
	}

	if req.Role != nil {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, *req.Role)
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
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Lấy thông tin user sau khi update
	return r.GetByID(ctx, id)
}

// Delete xóa user
func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM "user" WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// List lấy danh sách user với phân trang
func (r *userRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int, error) {
	// Đếm tổng số records
	countQuery := `SELECT COUNT(*) FROM "user"`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Lấy danh sách user
	offset := (page - 1) * pageSize
	query := `
		SELECT id, name, email, password, role, created_at
		FROM "user"
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var users []*model.User
	err = r.db.SelectContext(ctx, &users, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// Search tìm kiếm user
func (r *userRepository) Search(ctx context.Context, name, email, role string, page, pageSize int) ([]*model.User, int, error) {
	// Xây dựng query tìm kiếm
	whereClauses := []string{}
	var args []interface{}
	argIndex := 1

	if name != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("name ILIKE $%d", argIndex))
		args = append(args, "%"+name+"%")
		argIndex++
	}

	if email != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("email ILIKE $%d", argIndex))
		args = append(args, "%"+email+"%")
		argIndex++
	}

	if role != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, role)
		argIndex++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Đếm tổng số records
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM "user" %s`, whereClause)
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Lấy danh sách user
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, name, email, password, role, created_at
		FROM "user"
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	var users []*model.User
	err = r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	return users, total, nil
} 