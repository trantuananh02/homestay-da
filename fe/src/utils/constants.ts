export const API_ENDPOINTS = {
  HOMESTAYS: '/homestays',
  BOOKINGS: '/bookings',
  AUTH: '/auth',
  USERS: '/users',
  REVIEWS: '/reviews',
} as const;

export const ROUTES = {
  HOME: '/',
  HOMESTAYS: '/homestays',
  HOMESTAY_DETAIL: '/homestay/:id',
  ABOUT: '/about',
  BOOKINGS: '/bookings',
  MANAGEMENT: '/management',
  ADD_HOMESTAY: '/add-homestay',
} as const;

export const USER_ROLES = {
  GUEST: 'guest',
  HOST: 'host',
  ADMIN: 'admin',
} as const;

export const BOOKING_STATUS = {
  PENDING: 'pending',
  CONFIRMED: 'confirmed',
  CANCELLED: 'cancelled',
  COMPLETED: 'completed',
} as const;

export const HOMESTAY_STATUS = {
  ACTIVE: 'active',
  INACTIVE: 'inactive',
} as const;

export const PRICE_RANGES = [
  { label: 'Dưới 500k', value: '0-500000' },
  { label: '500k - 1tr', value: '500000-1000000' },
  { label: '1tr - 2tr', value: '1000000-2000000' },
  { label: 'Trên 2tr', value: '2000000+' },
];

export const GUEST_OPTIONS = Array.from({ length: 10 }, (_, i) => i + 1);

export const COMMON_AMENITIES = [
  'Wi-Fi',
  'Điều hòa',
  'Bếp đầy đủ',
  'Máy giặt',
  'Hồ bơi',
  'Bãi đậu xe',
  'BBQ',
  'Sân vườn',
  'Netflix',
  'Lò sưởi',
  'Thang máy',
  'An ninh 24/7',
  'Xe đạp miễn phí',
  'Dịch vụ dọn phòng',
  'Hồ bơi riêng',
  'Bếp mini',
];