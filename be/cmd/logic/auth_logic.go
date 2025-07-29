package logic

import (
	"context"
	"errors"
	"homestay-be/cmd/database/model"
	"homestay-be/cmd/svc"
	"homestay-be/cmd/types"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type AuthLogic struct {
	svc *svc.ServiceContext
}

func NewAuthLogic(svc *svc.ServiceContext) *AuthLogic {
	return &AuthLogic{svc: svc}
}

// Login xử lý đăng nhập
func (l *AuthLogic) Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error) {
	log.Println("Input Login Request:", req)

	// Tìm user theo email
	user, err := l.svc.UserRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("email hoặc mật khẩu không đúng")
	}

	// Kiểm tra password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logx.Error(err)
		return nil, errors.New("email hoặc mật khẩu không đúng")
	}

	// Tạo JWT token
	accessToken, err := l.generateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("không thể tạo token")
	}

	// Tạo response
	response := &types.LoginResponse{
		User: types.UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Phone: user.Phone,
			Email: user.Email,
			Role:  user.Role,
		},
		AccessToken: accessToken,
		ExpiresIn:   30 * 24 * 60 * 60, // 30 ngày tính bằng giây
	}

	return response, nil
}

// Register xử lý đăng ký
func (l *AuthLogic) Register(ctx context.Context, req *types.RegisterRequest) (*types.LoginResponse, error) {
	log.Println("Input Register Request:", req)

	// Kiểm tra email đã tồn tại chưa
	existingUser, err := l.svc.UserRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		logx.Error(err)
		return nil, errors.New("email đã được sử dụng")
	}

	// Tạo user mới
	userReq := &model.UserCreateRequest{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		Role:     req.Role,
	}

	user, err := l.svc.UserRepo.Create(ctx, userReq)
	if err != nil {
		logx.Error("Failed to create user:", err)
		return nil, errors.New("không thể tạo tài khoản")
	}

	// Tạo JWT token
	accessToken, err := l.generateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("không thể tạo token")
	}

	// Tạo response
	response := &types.LoginResponse{
		User: types.UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Phone: user.Phone,
			Email: user.Email,
			Role:  user.Role,
		},
		AccessToken: accessToken,
		ExpiresIn:   30 * 24 * 60 * 60, // 30 ngày tính bằng giây
	}

	return response, nil
}

// GetProfile lấy thông tin profile
func (l *AuthLogic) GetProfile(ctx context.Context, userID int) (*types.ProfileResponse, error) {
	// Lấy thông tin user
	user, err := l.svc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("không tìm thấy người dùng")
	}

	// Tạo response
	response := &types.ProfileResponse{
		User: types.UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Phone: user.Phone,
			Email: user.Email,
			Role:  user.Role,
		},
	}

	return response, nil
}

// generateJWT tạo JWT token
func (l *AuthLogic) generateJWT(userID int, email, role string) (string, error) {
	// Tạo claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 ngày
		"iat":     time.Now().Unix(),
	}

	// Tạo token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ký token với secret key (nên lưu trong config)
	secretKey := "your-secret-key" // TODO: Lấy từ config
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logx.Error(err)
		return "", err
	}

	return tokenString, nil
}

// ValidateToken xác thực JWT token
func (l *AuthLogic) ValidateToken(tokenString string) (*types.UserInfo, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("your-secret-key"), nil // TODO: Lấy từ config
	})

	if err != nil {
		logx.Error(err)
		return nil, errors.New("token không hợp lệ")
	}

	// Kiểm tra token có hợp lệ không
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userInfo := &types.UserInfo{
			ID:    int(claims["user_id"].(float64)),
			Email: claims["email"].(string),
			Role:  claims["role"].(string),
		}
		return userInfo, nil
	}

	return nil, errors.New("token không hợp lệ")
}

// Update Profile cập nhật thông tin profile
func (l *AuthLogic) UpdateProfile(ctx context.Context, userID int, req *types.UpdateProfileRequest) (*types.ProfileResponse, error) {
	log.Println("Input UpdateProfile Request:", req)

	// Lấy thông tin user hiện tại
	user, err := l.svc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("không tìm thấy người dùng")
	}

	// Cập nhật thông tin người dùng
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	// Xử lý phone - có thể rỗng để xóa số điện thoại
	if req.Phone != "" {
		user.Phone = &req.Phone
	} else {
		// Nếu phone rỗng, set về null
		user.Phone = nil
	}

	var userUpdate model.UserUpdateRequest
	// Populate userUpdate fields from user
	userUpdate.Name = &user.Name
	userUpdate.Email = &user.Email
	userUpdate.Phone = user.Phone

	// Lưu thay đổi
	if _, err := l.svc.UserRepo.Update(ctx, user.ID, &userUpdate); err != nil {
		logx.Error(err)
		return nil, errors.New("không thể cập nhật thông tin")
	}

	// Tạo response
	response := &types.ProfileResponse{
		User: types.UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Phone: user.Phone,
			Email: user.Email,
			Role:  user.Role,
		},
	}

	return response, nil
}
