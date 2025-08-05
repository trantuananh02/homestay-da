package types

// RegisterRequest request đăng ký
type RegisterRequest struct {
	Name     string  `json:"name" binding:"required,min=2,max=100" example:"Nguyễn Văn A"`
	Email    string  `json:"email" binding:"required,email" example:"user@example.com"`
	Phone    *string `json:"phone,omitempty" example:"0123456789"`
	Password string  `json:"password" binding:"required,min=6" example:"password123"`
	Role     string  `json:"role" binding:"required,eq=guest" example:"guest"`
}

// LoginRequest request đăng nhập
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse response đăng nhập
type LoginResponse struct {
	User        UserInfo `json:"user"`
	AccessToken string   `json:"access_token"`
	ExpiresIn   int64    `json:"expires_in"`
}

// UserInfo thông tin user
type UserInfo struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Phone *string `json:"phone,omitempty"`
	Email string  `json:"email"`
	Role  string  `json:"role"`
}

// ProfileResponse response profile
type ProfileResponse struct {
	User UserInfo `json:"user"`
}

// UpdateProfileRequest request cập nhật profile
type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone,omitempty"`
}
