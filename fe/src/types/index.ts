export interface Room {
  id: number;
  homestayId: number;
  name: string;
  description: string;
  type: string;
  capacity: number;
  price: number;
  priceType: "per_night" | "per_person";
  status: "available" | "occupied" | "maintenance";
  area?: number; // Diện tích phòng (m2)
  images?: string[]; // Danh sách URL hình ảnh
  amenities?: string[]; // Danh sách tiện nghi
  createdAt?: string;
  updatedAt?: string;
}

export interface RoomAvailability {
  id: number;
  roomId: number;
  date: string;
  status: "available" | "booked" | "blocked";
  price?: number;
  createdAt: string;
  updatedAt: string;
}

export interface Review {
  id?: number;
  homestayId: number;
  bookingId: number;
  guestId?: number;
  guestName?: string;
  rating: number; // Sử dụng số nguyên từ 1 đến 5
  comment: string;
  createdAt?: string;
  imageUrls?: string[]; // Danh sách URL ảnh của review
  userName?: string; // Tên người đánh giá
  homestayName?: string; // Tên homestay
}

export interface Homestay {
  id: number;
  name: string;
  description: string;
  address: string;
  city: string;
  district: string;
  ward: string;
  latitude: number;
  longitude: number;
  hostId: number;
  status: "active" | "inactive";
  createdAt: string;
  updatedAt: string;
  rooms?: Room[];
  rating?: number; // Trung bình đánh giá
  totalReviews?: number; // Tổng số lượt đánh giá
  reviews?: Review[]; // Danh sách đánh giá
}

export interface HomestayStats {
  totalHomestays: number;
  activeHomestays: number;
  totalRooms: number;
  availableRooms: number;
  totalBookings: number;
  totalRevenue: number;
  monthlyRevenue?: number; // Doanh thu tháng này
  occupancyRate?: number; // Tỷ lệ lấp đầy (%)
}

export interface RoomStats {
  totalRooms: number;
  availableRooms: number;
  occupiedRooms: number;
  maintenanceRooms: number;
  averagePrice: number;
  totalRevenue: number;
  occupancyRate: number;
}

// export interface Booking {
//   id: number;
//   homestayId: number;
//   roomId?: number;
//   guestId: string;
//   guestName: string;
//   email: string;
//   phone: string;
//   checkIn: string;
//   checkOut: string;
//   guests: number;
//   totalPrice: number;
//   status: 'pending' | 'confirmed' | 'cancelled' | 'completed';
//   createdAt: string;
//   notes?: string;
// }

// export interface Booking {
//   id: number;
//   customerId: number;
//   customerName: string;
//   customerPhone: string;
//   customerEmail: string;
//   homestayId: number;
//   homestayName: string;
//   roomId: number;
//   roomName: string;
//   roomType: string;
//   checkIn: string;
//   checkOut: string;
//   nights: number;
//   totalAmount: number;
//   paidAmount: number;
//   status: 'pending' | 'confirmed' | 'cancelled' | 'completed';
//   bookingDate: string;
//   paymentMethod: string;
// }

// interface Room {
//   id: number;
//   name: string;
//   type: string;
//   pricePerNight: number;
//   capacity: number;
//   amenities: string[];
//   isAvailable: boolean;
// }

export interface Booking {
  id: number;
  homestayId: number;
  homestayName?: string;
  bookingCode: string;
  customerName: string;
  customerPhone: string;
  customerEmail: string;
  rooms: {
    id: number;
    name: string;
    type: string;
    pricePerNight: number;
    nights: number;
    subtotal: number;
  }[];
  checkIn: string;
  checkOut: string;
  nights: number;
  totalAmount: number;
  paidAmount: number;
  status: "pending" | "confirmed" | "cancelled" | "completed";
  bookingDate: string;
  paymentMethod: string;
}

export interface Payment {
  id: string;
  bookingCode: string;
  amount: number;
  paymentMethod: "Tiền mặt" | "Thẻ tín dụng" | "Chuyển khoản" | "Momo";
  paymentStatus: "paid" | "unpaid" | "refunded";
  paymentDate: string;
}

export interface User {
  id: number;
  name: string;
  email: string;
  phone?: string;
  role: "guest" | "host";
  avatar?: string;
  createdAt: string;
}

export interface Review {
  id?: number;
  homestayId: number;
  bookingId: number;
  guestId?: number;
  guestName?: string;
  rating: number;
  comment: string;
  createdAt?: string;
  imageUrls?: string[]; // Danh sách URL ảnh của review
  userName?: string; // Tên người đánh giá
  homestayName?: string; // Tên homestay
}

// API Request/Response Types
export interface CreateHomestayRequest {
  name: string;
  description: string;
  address: string;
  city: string;
  district: string;
  ward: string;
  latitude: number;
  longitude: number;
}

export interface UpdateHomestayRequest {
  name?: string;
  description?: string;
  address?: string;
  city?: string;
  district?: string;
  ward?: string;
  latitude?: number;
  longitude?: number;
  status?: "active" | "inactive";
}

export interface CreateRoomRequest {
  homestayId: number;
  name: string;
  description: string;
  type: "Standard" | "Deluxe" | "Premium" | "Suite";
  capacity: number;
  price: number;
  priceType: "per_night" | "per_person";
  amenities: string[];
  images: string[];
}

export interface UpdateRoomRequest {
  name?: string;
  description?: string;
  type?: "Standard" | "Deluxe" | "Premium" | "Suite";
  capacity?: number;
  price?: number;
  priceType?: "per_night" | "per_person";
  status?: "available" | "occupied" | "maintenance";
  amenities?: string[];
  images?: string[];
}

export interface CreateAvailabilityRequest {
  roomId: number;
  date: string;
  status: "available" | "booked" | "blocked";
  price?: number;
}

export interface UpdateAvailabilityRequest {
  status?: "available" | "booked" | "blocked";
  price?: number;
}

export interface BulkAvailabilityRequest {
  roomId: number;
  startDate: string;
  endDate: string;
  status: "available" | "booked" | "blocked";
  price?: number;
  excludeDates?: string[];
}

export interface HomestayListRequest {
  page?: number;
  pageSize?: number;
  status?: string;
  search?: string;
  city?: string;
  district?: string;
  checkIn?: string;
  checkOut?: string;
  guests?: number;
}

export interface RoomListRequest {
  homestayId: number;
  page?: number;
  pageSize?: number;
  status?: string;
  type?: string;
  minPrice?: number;
  maxPrice?: number;
}

export interface HomestayListResponse {
  homestays: Homestay[];
  total: number;
  page: number;
  pageSize: number;
  totalPage: number;
}

export interface RoomListResponse {
  rooms: Room[];
  total: number;
  page: number;
  pageSize: number;
  totalPage: number;
}

export interface HomestayDetailResponse {
  homestay: Homestay;
  rooms: Room[];
}

export interface RoomDetailResponse {
  room: Room;
  homestay: Homestay;
  availabilities: RoomAvailability[];
}
