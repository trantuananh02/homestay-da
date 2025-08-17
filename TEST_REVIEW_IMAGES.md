# Hướng dẫn Test Tính năng Review với Ảnh

## Mục tiêu
Kiểm tra xem tính năng upload ảnh cho review có hoạt động giống như upload ảnh cho phòng không.

## Các bước test

### 1. Khởi động hệ thống
```bash
# Database
docker-compose up -d

# Backend
cd be
go run main.go

# Frontend
cd fe
npm run dev
```

### 2. Test Upload Ảnh Review

#### Bước 1: Đăng nhập
- Mở trình duyệt và truy cập `http://localhost:5173`
- Đăng nhập với tài khoản guest

#### Bước 2: Tạo Review với Ảnh
- Vào trang `Booking History`
- Tìm một booking đã hoàn thành
- Click "Đánh giá" để mở ReviewModal
- Chọn ảnh qua `CusFormUpload` component
- Điền rating và comment
- Submit review

#### Bước 3: Kiểm tra kết quả
- Kiểm tra xem ảnh có được upload không
- Kiểm tra xem review có được tạo với ảnh không
- Kiểm tra xem ảnh có được hiển thị trong danh sách review không

### 3. So sánh với Upload Ảnh Phòng

#### Giống nhau:
- ✅ Sử dụng cùng `CusFormUpload` component
- ✅ Sử dụng cùng `homestayService.uploadRoomImage()` method
- ✅ Sử dụng cùng `/upload` endpoint
- ✅ Cùng cách xử lý authentication
- ✅ Cùng giao diện và style

#### Khác nhau:
- ❌ Không có gì khác biệt

### 4. Kiểm tra Database

#### Kết nối database:
```bash
psql -h localhost -U postgres -d homestay_db
```

#### Kiểm tra bảng review:
```sql
SELECT id, user_id, homestay_id, rating, comment, image_urls, created_at 
FROM review 
ORDER BY created_at DESC 
LIMIT 5;
```

#### Kiểm tra trường image_urls:
```sql
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'review' 
AND column_name = 'image_urls';
```

### 5. Kiểm tra API

#### Test upload endpoint:
```bash
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@test_image.jpg"
```

#### Test create review endpoint:
```bash
curl -X POST http://localhost:8080/api/review \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bookingId": 1,
    "rating": 5,
    "comment": "Test review with images",
    "imageUrls": ["url1", "url2"]
  }'
```

## Kết quả mong đợi

### ✅ Thành công:
- Ảnh được upload thành công
- Review được tạo với ảnh
- Ảnh được hiển thị trong danh sách review
- Hoạt động giống hệt như upload ảnh phòng

### ❌ Lỗi có thể gặp:
- Lỗi authentication khi upload
- Lỗi database khi lưu review
- Ảnh không được hiển thị
- Giao diện khác với upload ảnh phòng

## Debug

### Nếu ảnh không upload:
1. Kiểm tra console browser
2. Kiểm tra network tab
3. Kiểm tra backend logs
4. Kiểm tra database connection

### Nếu review không có ảnh:
1. Kiểm tra request payload
2. Kiểm tra database schema
3. Kiểm tra backend logic
4. Kiểm tra frontend state

## Kết luận

Tính năng review với ảnh phải hoạt động **giống hệt** như upload ảnh phòng:
- Cùng component
- Cùng service
- Cùng endpoint
- Cùng giao diện
- Cùng logic xử lý

Nếu có sự khác biệt nào, cần sửa lại để đảm bảo tính nhất quán.
