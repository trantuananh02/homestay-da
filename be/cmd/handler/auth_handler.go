package handler

import (
	"homestay-be/cmd/logic"
	"homestay-be/cmd/middleware"
	"homestay-be/cmd/svc"
	"homestay-be/cmd/types"
	"homestay-be/core/response"
	"log"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *svc.ServiceContext
}

func NewAuthHandler(svc *svc.ServiceContext) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login xử lý đăng nhập
func (h *AuthHandler) Login(c *gin.Context) {
	var req types.LoginRequest

	// Bind request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		response.ResponseError(c, response.BadRequest, "Email và mật khẩu không được để trống")
		return
	}

	// Xử lý logic
	authLogic := logic.NewAuthLogic(h.svc)
	resp, err := authLogic.Login(c.Request.Context(), &req)
	if err != nil {
		log.Print(err)
		response.ResponseError(c, response.BadRequest, response.MsgInvalidCredentials)
		return
	}

	// Trả về response thành công
	response.ResponseSuccess(c, resp)
}

// Register xử lý đăng ký
func (h *AuthHandler) Register(c *gin.Context) {
	var req types.RegisterRequest

	// Bind request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	// Validate request
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		response.ResponseError(c, response.BadRequest, response.MsgRequiredField)
		return
	}

	// Validate role
	validRoles := map[string]bool{"admin": true, "host": true, "guest": true}
	if !validRoles[req.Role] {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidRole+". Chỉ chấp nhận: admin, host, guest")
		return
	}

	// Xử lý logic
	authLogic := logic.NewAuthLogic(h.svc)
	result, err := authLogic.Register(c.Request.Context(), &req)
	if err != nil {
		log.Print(err)
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	// Trả về response thành công
	response.ResponseSuccess(c, result)
}

// GetProfile lấy thông tin profile của user hiện tại
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Lấy thông tin user từ middleware
	user, exists := middleware.GetCurrentUser(c)
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}

	// Lấy thông tin chi tiết từ database
	authLogic := logic.NewAuthLogic(h.svc)
	result, err := authLogic.GetProfile(c.Request.Context(), user.ID)
	if err != nil {
		response.ResponseError(c, response.NotFound, response.MsgUserNotFound)
		return
	}

	// Trả về response thành công
	response.ResponseSuccess(c, result)
}

// Logout đăng xuất
func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: Implement blacklist token logic
	response.ResponseSuccess(c, gin.H{
		"message": "Đăng xuất thành công",
	})
}

	// UpdateProfile cập nhật thông tin người dùng
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req types.UpdateProfileRequest

	// Bind request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	// Lấy thông tin user từ middleware
	user, exists := middleware.GetCurrentUser(c)
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}

	// Xử lý logic
	authLogic := logic.NewAuthLogic(h.svc)
	if _, err := authLogic.UpdateProfile(c.Request.Context(), user.ID, &req); err != nil {
		log.Print(err)
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	// Trả về response thành công
	response.ResponseSuccess(c, gin.H{"message": "Cập nhật thông tin thành công"})
}