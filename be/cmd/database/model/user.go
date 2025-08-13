package model

import (
	"time"
)

// User đại diện cho bảng user
type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Phone     *string   `db:"phone" json:"phone,omitempty"` // Phone có thể là null
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"` // Không trả về password trong JSON
	Role      string    `db:"role" json:"role"`
	Status    string    `db:"status" json:"status"` // Trạng thái người dùng (active, inactive)
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// UserCreateRequest request tạo user mới
type UserCreateRequest struct {
	Name     string  `json:"name" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Phone    *string `json:"phone,omitempty"`
	Password string  `json:"password" binding:"required,min=6"`
	Role     string  `json:"role" binding:"required,eq=guest"`
	Status   string  `json:"status" binding:"required,oneof=active inactive"`
}

// UserUpdateRequest request cập nhật user
type UserUpdateRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Phone    *string `json:"phone,omitempty"`
	Password *string `json:"password" binding:"omitempty,min=6"`
	Role     *string `json:"role" binding:"omitempty,oneof=host guest"`
}

// UserLoginRequest request đăng nhập
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
