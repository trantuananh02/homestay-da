-- Tạo bảng user
CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(20) DEFAULT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'host', 'guest')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tạo bảng homestay
CREATE TABLE homestay (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    ward VARCHAR(100) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    owner_id INTEGER NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    rate DECIMAL(3,2) DEFAULT 0 CHECK (rate >= 0 AND rate <= 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tạo bảng room
CREATE TABLE room (
    id SERIAL PRIMARY KEY,
    homestay_id INTEGER NOT NULL REFERENCES homestay(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL,
    capacity INTEGER NOT NULL,
    price DECIMAL(12,2) NOT NULL,
    price_type VARCHAR(20) NOT NULL DEFAULT 'per_night' CHECK (price_type IN ('per_night', 'per_person')),
    status VARCHAR(20) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'occupied', 'maintenance')),
    image_urls TEXT, -- Mảng chứa các URL ảnh của phòng
    amenities TEXT, -- Mảng chứa các tiện nghi của phòng
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tạo bảng booking (sau khi được xác nhận)
CREATE TABLE booking (
    id SERIAL PRIMARY KEY,
    booking_code VARCHAR(50) UNIQUE NOT NULL, -- Mã đặt phòng duy nhất
    homestay_id INTEGER NOT NULL REFERENCES homestay(id) ON DELETE CASCADE,
    email VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    check_in DATE NOT NULL,
    check_out DATE NOT NULL,
    num_guests INTEGER NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payment_method VARCHAR(50) NOT NULL,
    paid_amount DECIMAL(12,2) NOT NULL DEFAULT 0
);

-- Tạo bảng booking_room để quản lý mối quan hệ giữa booking và room
CREATE TABLE booking_room (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL REFERENCES booking(id) ON DELETE CASCADE,
    room_id INTEGER NOT NULL REFERENCES room(id) ON DELETE CASCADE,
    room_name VARCHAR(100) NOT NULL,
    room_type VARCHAR(20) NOT NULL,
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    price DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tạo bảng payment để quản lý thanh toán
CREATE TABLE payment (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL REFERENCES booking(id) ON DELETE CASCADE,
    amount DECIMAL(12,2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL, -- 'cash', 'bank_transfer', 'credit_card', etc.
    payment_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    transaction_id VARCHAR(100), -- Mã giao dịch từ cổng thanh toán
    payment_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tạo bảng review
CREATE TABLE review (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    homestay_id INTEGER NOT NULL REFERENCES homestay(id) ON DELETE CASCADE,
    booking_id INTEGER REFERENCES booking(id) ON DELETE SET NULL, -- Liên kết với booking cụ thể
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
