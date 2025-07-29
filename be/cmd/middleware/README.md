# Middleware Package

Package này chứa các middleware functions cho ứng dụng homestay backend.

## Các Middleware Functions

### 1. AuthMiddleware
Xác thực JWT token từ header Authorization.

**Signature:**
```go
func AuthMiddleware(svc *svc.ServiceContext) gin.HandlerFunc
```

**Cách sử dụng:**
```go
router.Use(middleware.AuthMiddleware(serverCtx))
```

**Chức năng:**
- Kiểm tra header Authorization có format "Bearer <token>"
- Xác thực JWT token
- Lưu thông tin user vào context để sử dụng sau này
- Trả về lỗi 401 nếu token không hợp lệ

### 2. GetCurrentUser
Lấy thông tin user hiện tại từ context.

**Signature:**
```go
func GetCurrentUser(c *gin.Context) (*types.UserInfo, bool)
```

**Cách sử dụng:**
```go
user, exists := middleware.GetCurrentUser(c)
if !exists {
    // Xử lý khi không có user
    return
}
// Sử dụng user.ID, user.Name, user.Email, user.Role
```

**Chức năng:**
- Lấy thông tin user đã được lưu trong context bởi AuthMiddleware
- Trả về pointer đến UserInfo và boolean exists
- Sử dụng trong các handler để lấy thông tin user hiện tại

### 3. RoleMiddleware
Kiểm tra role của user có quyền truy cập hay không.

**Signature:**
```go
func RoleMiddleware(roles ...string) gin.HandlerFunc
```

**Cách sử dụng:**
```go
// Chỉ cho phép admin
router.Use(middleware.RoleMiddleware("admin"))

// Cho phép host hoặc admin
router.Use(middleware.RoleMiddleware("host", "admin"))
```

**Chức năng:**
- Kiểm tra role của user có trong danh sách roles được cho phép
- Trả về lỗi 403 nếu user không có quyền
- Phải được sử dụng sau AuthMiddleware

### 4. CORSMiddleware
Xử lý CORS (Cross-Origin Resource Sharing).

**Signature:**
```go
func CORSMiddleware() gin.HandlerFunc
```

**Cách sử dụng:**
```go
router.Use(middleware.CORSMiddleware())
```

**Chức năng:**
- Cho phép cross-origin requests
- Thêm các header CORS cần thiết
- Hỗ trợ preflight requests

## Ví dụ Sử Dụng

### Trong Handler
```go
func (h *AuthHandler) GetProfile(c *gin.Context) {
    // Lấy thông tin user từ middleware
    user, exists := middleware.GetCurrentUser(c)
    if !exists {
        response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
        return
    }

    // Sử dụng user.ID để lấy thông tin chi tiết
    authLogic := logic.NewAuthLogic(h.svc)
    result, err := authLogic.GetProfile(c.Request.Context(), user.ID)
    if err != nil {
        response.ResponseError(c, response.NotFound, response.MsgUserNotFound)
        return
    }

    response.ResponseSuccess(c, result)
}
```

### Trong Router
```go
// Protected routes
protected := api.Group("")
protected.Use(middleware.AuthMiddleware(serverCtx))
{
    protected.GET("/auth/profile", authHandler.GetProfile)
}

// Host routes
host := api.Group("/host")
host.Use(middleware.AuthMiddleware(serverCtx))
host.Use(middleware.RoleMiddleware("host", "admin"))
{
    host.GET("/homestays", homestayHandler.GetHomestayList)
}
```

## Lưu Ý

1. **Thứ tự middleware quan trọng**: AuthMiddleware phải được sử dụng trước RoleMiddleware
2. **Service Context**: AuthMiddleware cần service context để xác thực token
3. **Error Handling**: Các middleware sẽ tự động trả về lỗi và abort request nếu có vấn đề
4. **Performance**: GetCurrentUser chỉ đọc từ context, không gọi database 