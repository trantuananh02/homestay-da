package middleware

import (
	"homestay-be/cmd/logic"
	"homestay-be/cmd/svc"
	"homestay-be/cmd/types"
	"homestay-be/core/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware middleware xác thực JWT token
func AuthMiddleware(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ResponseError(c, response.Unauthorized, response.MsgTokenRequired)
			c.Abort()
			return
		}

		// Kiểm tra format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ResponseError(c, response.Unauthorized, response.MsgTokenInvalid)
			c.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			response.ResponseError(c, response.Unauthorized, response.MsgTokenRequired)
			c.Abort()
			return
		}

		// Xác thực token
		authLogic := logic.NewAuthLogic(svc)
		userInfo, err := authLogic.ValidateToken(token)
		if err != nil {
			response.ResponseError(c, response.Unauthorized, response.MsgTokenInvalid)
			c.Abort()
			return
		}

		// Lưu user ID và thông tin user vào context
		c.Set("user_id", userInfo.ID)
		c.Set("user", userInfo)
		c.Next()
	}
}

// GetCurrentUser lấy thông tin user hiện tại từ context
func GetCurrentUser(c *gin.Context) (*types.UserInfo, bool) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	user, ok := userInterface.(*types.UserInfo)
	if !ok {
		return nil, false
	}

	return user, true
}

// RoleMiddleware middleware kiểm tra role
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := GetCurrentUser(c)
		if !exists {
			response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
			c.Abort()
			return
		}

		// Kiểm tra role
		hasRole := false
		for _, role := range roles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.ResponseError(c, response.Forbidden, response.MsgForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
