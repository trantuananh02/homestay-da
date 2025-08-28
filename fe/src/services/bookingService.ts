import api from "./api";
import { toastService } from "./toastService";

import { Booking, Review, Payment } from "../types";

export const bookingService = {
  // Lọc danh sách booking (filter)
  filterBookings: async (
    params: Record<string, unknown>
  ): Promise<{
    [x: string]: number; bookings: Booking[] 
}> => {
    const response = await api.get("/api/host/booking", { params });
    return response.data;
  },

  // Lấy danh sách booking của khách
  getGuestBookings: async (
    params: Record<string, unknown>
  ): Promise<{ bookings: Booking[] }> => {
    const response = await api.get(`/api/guest/booking`, { params });
    return response.data;
  },

  // Tạo booking mới
  createBooking: async (
    bookingData: Omit<Booking, "id" | "bookingCode" | "status" | "bookingDate">
  ): Promise<Booking> => {
    const response = await api.post("/api/host/booking", bookingData);
    toastService.success("Tạo booking thành công");
    return response.data;
  },

  // Tạo booking cho khách
  createGuestBooking: async (
    bookingData: Omit<Booking, "id" | "bookingCode" | "status" | "bookingDate">
  ): Promise<Booking> => {
    const response = await api.post("/api/guest/booking", bookingData);
    toastService.success("Tạo booking thành công");
    return response.data;
  },

  // Lấy chi tiết booking
  getBookingDetail: async (id: number): Promise<Booking> => {
    const response = await api.get(`/api/booking/${id}`);
    return response.data;
  },

  // Cập nhật trạng thái booking
  updateBookingStatus: async (id: number, status: string): Promise<Booking> => {
    const response = await api.put(`/api/host/booking/${id}/status`, {
      status,
    });
    if (response.status === 200) {
      toastService.success("Cập nhật trạng thái booking thành công");
    } else {
      toastService.error("Cập nhật trạng thái booking thất bại");
    }
    return response.data;
  },

  // Cập nhật trạng thái booking của khách
  updateGuestBookingStatus: async (
    id: number,
    status: string
  ): Promise<Booking> => {
    const response = await api.put(`/api/guest/booking/${id}/status`, {
      status,
    });
    if (response.status === 200) {
      toastService.success("Cập nhật trạng thái booking thành công");
    } else {
      toastService.error("Cập nhật trạng thái booking thất bại");
    }
    return response.data;
  },

  // Tạo đánh giá cho booking
  createReview: async (reviewData: Review): Promise<Review> => {
    const response = await api.post("/api/guest/review", reviewData);
    toastService.success("Đánh giá đã được gửi thành công");
    return response.data;
  },

  // Tạo link thanh toán VNPAY cho khách thuê
  createVnpayPayment: async (
    amount: number,
    orderId: string,
    orderInfo: string
  ): Promise<string> => {
    const response = await api.get("/api/guest/payment/vnpay", {
      params: { amount, orderId, orderInfo },
    });
    return response.data.paymentUrl;
  },

  // lấy danh sách thanh toán
  getPayments: async (
    params: Record<string, unknown>
  ): Promise<{ payments: Payment[] }> => {
    const response = await api.get("/api/host/payments", { params });
    return response.data;
  },
};
