import api from './api';
import { toastService } from './toastService';
import {
  Homestay,
  CreateHomestayRequest,
  UpdateHomestayRequest,
  HomestayListRequest,
  HomestayListResponse,
  HomestayDetailResponse,
  HomestayStats,
  Room,
  CreateRoomRequest,
  UpdateRoomRequest,
  RoomListRequest,
  RoomListResponse,
  RoomDetailResponse,
  RoomStats,
  CreateAvailabilityRequest,
  UpdateAvailabilityRequest,
  BulkAvailabilityRequest,
  RoomAvailability,
  Booking
} from '../types';

class HomestayService {
  // Homestay Management
  async createHomestay(data: CreateHomestayRequest): Promise<Homestay> {
    try {
      const response = await api.post('/api/host/homestays', data);
      toastService.success('Tạo homestay thành công');
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi tạo homestay');
      throw error;
    }
  }

  async getHomestayList(params: HomestayListRequest = {}): Promise<HomestayListResponse> {
    try {
      const response = await api.get('/api/host/homestays', { params });

      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy danh sách homestay');
      throw error;
    }
  }


  async getHomestayById(id: number): Promise<HomestayDetailResponse> {
    try {
      const response = await api.get(`/api/host/homestays/${id}`);
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy thông tin homestay');
      throw error;
    }
  }

  async getPublicHomestayDetail(id: number): Promise<Homestay> {
    try {
      const response = await api.get(`/api/public/homestays/${id}`);
      console.log('Homestay detail response:', response.data.homestay);

      return response.data.homestay;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy thông tin homestay');
      throw error;
    }
  }

  async updateHomestay(id: number, data: UpdateHomestayRequest): Promise<Homestay> {
    try {
      const response = await api.put(`/api/host/homestays/${id}`, data);
      toastService.success('Cập nhật homestay thành công');
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi cập nhật homestay');
      throw error;
    }
  }

  async deleteHomestay(id: number): Promise<void> {
    try {
      await api.delete(`/api/host/homestays/${id}`);
      toastService.success('Xóa homestay thành công');
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi xóa homestay');
      throw error;
    }
  }

  async toggleHomestayStatus(id: number): Promise<Homestay> {
    try {
      const response = await api.put(`/api/host/homestays/${id}/toggle-status`);
      const newStatus = response.data.homestay.status;
      const statusText = newStatus === 'active' ? 'Hoạt động' : 'Không hoạt động';
      toastService.success(`Đã chuyển homestay sang trạng thái ${statusText}`);
      return response.data.homestay;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi thay đổi trạng thái homestay');
      throw error;
    }
  }

  async getHomestayStats(): Promise<HomestayStats> {
    try {
      const response = await api.get('/api/host/homestays/stats');

      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy thống kê homestay');
      throw error;
    }
  }

  async getHomestayStatsById(id: number): Promise<HomestayStats> {
    try {
      const response = await api.get(`/api/host/homestays/${id}/stats`);
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy thống kê homestay');
      throw error;
    }
  }

  // Room Management
  async createRoom(data: CreateRoomRequest): Promise<Room> {
    try {
      const response = await api.post('/api/host/rooms', data);
      toastService.success('Tạo phòng thành công');
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi tạo phòng');
      throw error;
    }
  }

  async getRoomList(params: RoomListRequest): Promise<RoomListResponse> {
    try {
      const response = await api.get('/api/host/rooms', { params });
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy danh sách phòng');
      throw error;
    }
  }

  async getRoomById(id: number): Promise<RoomDetailResponse> {
    try {
      const response = await api.get(`/api/host/rooms/${id}`);
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy thông tin phòng');
      throw error;
    }
  }

  async updateRoom(id: number, data: UpdateRoomRequest): Promise<Room> {
    try {
      console.log('Updating room:----------------------------', id, data);


      const response = await api.put(`/api/host/rooms/${id}`, data);
      toastService.success('Cập nhật phòng thành công');
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi cập nhật phòng');
      throw error;
    }
  }

  async deleteRoom(id: number): Promise<void> {
    try {
      await api.delete(`/api/host/rooms/${id}`);
      toastService.success('Xóa phòng thành công');
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi xóa phòng');
      throw error;
    }
  }

  async getRoomStats(homestayId: number): Promise<RoomStats> {
    try {
      const response = await api.get(`/api/host/homestays/${homestayId}/rooms/stats`);
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy thống kê phòng');
      throw error;
    }
  }

  // Room Availability Management
  async createAvailability(data: CreateAvailabilityRequest): Promise<RoomAvailability> {
    try {
      const response = await api.post('/api/host/rooms/availability', data);
      toastService.success('Tạo availability thành công');
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi tạo availability');
      throw error;
    }
  }

  async updateAvailability(id: number, data: UpdateAvailabilityRequest): Promise<RoomAvailability> {
    try {
      const response = await api.put(`/api/host/rooms/availability/${id}`, data);
      toastService.success('Cập nhật availability thành công');
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi cập nhật availability');
      throw error;
    }
  }

  async bulkUpdateAvailability(data: BulkAvailabilityRequest): Promise<void> {
    try {
      await api.post('/api/host/rooms/availability/bulk', data);
      toastService.success('Cập nhật availability hàng loạt thành công');
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi cập nhật availability hàng loạt');
      throw error;
    }
  }

  async uploadRoomImage(
    file: File
  ): Promise<string> {
    const formData = new FormData();

    formData.append('image', file);

    try {
      const response = await api.post(`/upload`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      });

      toastService.success('Tải lên hình ảnh phòng thành công');
      return response.data.url;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi tải lên hình ảnh phòng');
      throw error;
    }
  }
  // Helper methods for data transformation
  formatPrice(price: number): string {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND'
    }).format(price);
  }

  formatPriceType(priceType: string): string {
    return priceType === 'per_night' ? 'theo đêm' : 'theo người';
  }

  formatRoomType(type: string): string {
    const typeMap: Record<string, string> = {
      single: 'Phòng đơn',
      double: 'Phòng đôi',
      family: 'Phòng gia đình',
      dormitory: 'Phòng tập thể'
    };
    return typeMap[type] || type;
  }

  formatStatus(status: string): string {
    const statusMap: Record<string, string> = {
      active: 'Hoạt động',
      inactive: 'Không hoạt động',
      available: 'Có thể đặt',
      occupied: 'Đã được đặt',
      maintenance: 'Bảo trì'
    };
    return statusMap[status] || status;
  }

  getStatusColor(status: string): string {
    const colorMap: Record<string, string> = {
      active: 'bg-green-100 text-green-800',
      inactive: 'bg-red-100 text-red-800',
      available: 'bg-green-100 text-green-800',
      occupied: 'bg-red-100 text-red-800',
      maintenance: 'bg-orange-100 text-orange-800'
    };
    return colorMap[status] || 'bg-gray-100 text-gray-800';
  }

  async getPublicHomestayList(params: HomestayListRequest = {}): Promise<HomestayListResponse> {
    try {
      const response = await api.get('/api/public/homestays', { params });

      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy danh sách homestay');
      throw error;
    }
  }

  async getTopHomestays(limit: number = 8): Promise<Homestay[]> {
    try {
      const response = await api.get('/api/public/homestays/top', { params: { limit } });
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy danh sách homestay nổi bật');
      throw error;
    }
  }

  async getPublicHomestayById(id: number): Promise<HomestayDetailResponse> {
    try {
      const response = await api.get(`/api/public/homestays/${id}`);
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy thông tin homestay');
      throw error;
    }
  }

  async getPublicRoomList(params: RoomListRequest): Promise<RoomListResponse> {
    try {
      const response = await api.get('/api/guest/rooms', { params });

      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy danh sách phòng');
      throw error;
    }
  }

  async getBookingsByHomestayId(homestayId: number) : Promise<{ bookings: Booking[] }> {
    try {
      const response = await api.get(`/api/host/homestays/${homestayId}/bookings`);
      
      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy danh sách đặt phòng');
      throw error;
    }
  }

  async getGuestBookingsByHomestayId(homestayId: number) : Promise<{ bookings: Booking[] }> {
    try {
      const response = await api.get(`/api/guest/homestays/${homestayId}/bookings`);

      return response.data;
    } catch (error: any) {
      toastService.error(error.response?.data?.result?.message || 'Lỗi khi lấy danh sách đặt phòng');
      throw error;
    }
  }

}

const homestayService = new HomestayService();
export { homestayService };