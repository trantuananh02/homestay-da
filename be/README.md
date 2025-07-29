# Homestay Backend API

Backend API cho ứng dụng quản lý homestay được xây dựng bằng Go và Gin framework.

## Cấu trúc dự án

```
be/
├── cmd/
│   ├── config/     # Cấu hình ứng dụng
│   ├── server/     # Logic server và routes
│   ├── handler/    # HTTP handlers
│   ├── logic/      # Business logic
│   ├── svc/        # Services
│   ├── database/   # Database models và repositories
│   ├── types/      # Type definitions
│   └── utils/      # Utility functions
├── configs/        # File cấu hình
├── core/           # Core components
└── main.go         # Entry point
```

## Cài đặt và chạy

### Yêu cầu
- Go 1.24.5 trở lên
- MySQL database

### Cài đặt dependencies
```bash
go mod tidy
```

### Cấu hình
Chỉnh sửa file `configs/config.yaml`:
```yaml
http:
  path: localhost
  port: 8080
database:
  driver: mysql
  source: root:password@tcp(localhost:3306)/homestay
```

### Chạy server
```bash
go run main.go
```

Server sẽ chạy tại `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - Kiểm tra trạng thái server

### Homestay API
- `GET /api/v1/homestays` - Lấy danh sách homestay
- `GET /api/v1/homestays/:id` - Lấy chi tiết homestay
- `POST /api/v1/homestays` - Tạo homestay mới
- `PUT /api/v1/homestays/:id` - Cập nhật homestay
- `DELETE /api/v1/homestays/:id` - Xóa homestay

### Booking API
- `GET /api/v1/bookings` - Lấy danh sách booking
- `POST /api/v1/bookings` - Tạo booking mới
- `PUT /api/v1/bookings/:id` - Cập nhật booking
- `DELETE /api/v1/bookings/:id` - Xóa booking

### Authentication API
- `POST /api/v1/auth/login` - Đăng nhập
- `POST /api/v1/auth/register` - Đăng ký

## Tính năng

- ✅ Đọc cấu hình từ file YAML
- ✅ Khởi động server với Gin framework
- ✅ CORS middleware
- ✅ Health check endpoint
- ✅ API routes cơ bản
- ✅ Error handling
- ✅ Logging

## Phát triển

### Thêm route mới
1. Thêm handler method trong `cmd/server/server.go`
2. Đăng ký route trong `SetupRoutes()`

### Thêm cấu hình mới
1. Thêm field vào struct `Config` trong `cmd/config/config.go`
2. Cập nhật file `configs/config.yaml`

## License

MIT
