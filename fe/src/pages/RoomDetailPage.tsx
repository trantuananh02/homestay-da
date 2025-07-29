import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { homestayService } from '../services/homestayService';
import { Room, RoomDetailResponse } from '../types';
import { X, MapPin, Star } from 'lucide-react';

const RoomDetailPage: React.FC = () => {
  const { roomId } = useParams<{ roomId: string }>();
  const navigate = useNavigate();
  const [room, setRoom] = useState<Room | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'overview' | 'bookings' | 'reviews'>('overview');

  useEffect(() => {
    const fetchRoom = async () => {
      if (!roomId) return;
      try {
        setLoading(true);
        const res: RoomDetailResponse = await homestayService.getRoomById(Number(roomId));
        setRoom(res.room);
      } catch (error) {
        // Xử lý lỗi nếu cần
      } finally {
        setLoading(false);
      }
    };
    fetchRoom();
  }, [roomId]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Đang tải thông tin phòng...</p>
        </div>
      </div>
    );
  }

  if (!room) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600">Không tìm thấy phòng</p>
          <button
            onClick={() => navigate(-1)}
            className="mt-4 px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700"
          >
            Quay lại
          </button>
        </div>
      </div>
    );
  }

  const fakeBookings = [
    {
      id: 'B001',
      guestName: 'Nguyễn Văn A',
      checkIn: '2024-06-01',
      checkOut: '2024-06-03',
      guests: 2,
      totalPrice: 1200000,
      status: 'completed',
    },
    {
      id: 'B002',
      guestName: 'Trần Thị B',
      checkIn: '2024-06-05',
      checkOut: '2024-06-07',
      guests: 3,
      totalPrice: 1800000,
      status: 'confirmed',
    },
  ];

  const fakeReviews = [
    {
      id: 'R001',
      guestName: 'Nguyễn Văn A',
      rating: 5,
      comment: 'Phòng sạch sẽ, tiện nghi. Sẽ quay lại!',
      createdAt: '2024-06-04',
    },
    {
      id: 'R002',
      guestName: 'Trần Thị B',
      rating: 4,
      comment: 'Dịch vụ tốt, phòng đẹp nhưng hơi ồn.',
      createdAt: '2024-06-08',
    },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <button
            onClick={() => navigate(-1)}
            className="flex items-center space-x-2 text-emerald-600 hover:text-emerald-700 mb-4"
          >
            <X className="h-5 w-5" />
            <span>Quay lại</span>
          </button>
          <div className="flex flex-col md:flex-row md:justify-between md:items-center">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">{room.name}</h1>
              <div className="flex items-center text-gray-600 mb-2">
                <MapPin className="h-4 w-4 mr-1" />
                <span>{room.type} • Sức chứa: {room.capacity} người</span>
              </div>
              <div className="flex items-center space-x-4 text-sm text-gray-600">
                <span className="font-medium">{room.price.toLocaleString()} VNĐ / đêm</span>
                <span className={`px-2 py-1 rounded-full text-xs ${room.status === 'available' ? 'bg-primary-100 text-primary-700' : room.status === 'occupied' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-700'}`}>{room.status === 'available' ? 'Còn trống' : room.status === 'occupied' ? 'Đã đặt' : 'Bảo trì'}</span>
              </div>
            </div>
          </div>
        </div>
        {/* Tab Navigation */}
        <div className="bg-white rounded-xl shadow-sm mb-8">
          <div className="border-b border-gray-200">
            <nav className="flex space-x-8 px-6">
              <button
                onClick={() => setActiveTab('overview')}
                className={`py-4 px-1 border-b-2 font-medium text-sm ${activeTab === 'overview' ? 'border-emerald-500 text-emerald-600' : 'border-transparent text-gray-500 hover:text-gray-700'}`}
              >
                Tổng quan
              </button>
              <button
                onClick={() => setActiveTab('bookings')}
                className={`py-4 px-1 border-b-2 font-medium text-sm ${activeTab === 'bookings' ? 'border-emerald-500 text-emerald-600' : 'border-transparent text-gray-500 hover:text-gray-700'}`}
              >
                Đặt phòng
              </button>
              <button
                onClick={() => setActiveTab('reviews')}
                className={`py-4 px-1 border-b-2 font-medium text-sm ${activeTab === 'reviews' ? 'border-emerald-500 text-emerald-600' : 'border-transparent text-gray-500 hover:text-gray-700'}`}
              >
                Đánh giá
              </button>
            </nav>
          </div>
          <div className="p-6">
            {activeTab === 'overview' && (
              <>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">Tên phòng</label>
                    <div className="p-3 border border-gray-200 rounded-lg bg-gray-50">{room.name}</div>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">Loại phòng</label>
                    <div className="p-3 border border-gray-200 rounded-lg bg-gray-50">{room.type}</div>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">Sức chứa</label>
                    <div className="p-3 border border-gray-200 rounded-lg bg-gray-50">{room.capacity} người</div>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">Giá phòng/đêm</label>
                    <div className="p-3 border border-gray-200 rounded-lg bg-gray-50">{room.price.toLocaleString()} VNĐ</div>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">Trạng thái</label>
                    <div className="p-3 border border-gray-200 rounded-lg bg-gray-50">{room.status}</div>
                  </div>
                </div>
                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 mb-2">Mô tả</label>
                  <div className="p-3 border border-gray-200 rounded-lg bg-gray-50">{room.description}</div>
                </div>
                {/* Tiện ích */}
                {Array.isArray(room.amenities) && room.amenities.length > 0 && (
                  <div className="mb-6">
                    <label className="block text-sm font-medium text-gray-700 mb-2">Tiện ích</label>
                    <div className="flex flex-wrap gap-2">
                      {room.amenities.map((amenity, idx) => (
                        <span
                          key={idx}
                          className="px-3 py-1 bg-emerald-100 text-emerald-700 rounded-full text-xs font-medium"
                        >
                          {amenity}
                        </span>
                      ))}
                    </div>
                  </div>
                )}
                {/* Nếu có images thì hiển thị */}
                {Array.isArray(room.images) && room.images.length > 0 && (
                  <div className="mb-6">
                    <label className="block text-sm font-medium text-gray-700 mb-2">Hình ảnh</label>
                    <div className="flex flex-wrap gap-4">
                      {room.images.map((img, idx) => (
                        <img
                          key={idx}
                          src={img}
                          alt={`Hình phòng ${idx + 1}`}
                          className="w-32 h-24 object-cover rounded-lg border"
                        />
                      ))}
                    </div>
                  </div>
                )}
              </>
            )}
            {activeTab === 'bookings' && (
              <div>
                <h3 className="text-lg font-semibold mb-4 text-emerald-700">Danh sách đặt phòng</h3>
                <div className="space-y-4">
                  {fakeBookings.map((booking) => (
                    <div key={booking.id} className="border rounded-lg p-4 bg-gray-50 flex flex-col md:flex-row md:items-center md:justify-between">
                      <div>
                        <div className="font-medium text-gray-900">Khách: {booking.guestName}</div>
                        <div className="text-sm text-gray-600">Từ {booking.checkIn} đến {booking.checkOut} ({booking.guests} khách)</div>
                      </div>
                      <div className="mt-2 md:mt-0 flex flex-col md:items-end">
                        <div className="text-emerald-600 font-bold">{booking.totalPrice.toLocaleString()} VNĐ</div>
                        <span className={`text-xs px-2 py-1 rounded-full mt-1 ${booking.status === 'completed' ? 'bg-primary-100 text-primary-700' : booking.status === 'confirmed' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-700'}`}>{booking.status === 'completed' ? 'Hoàn thành' : booking.status === 'confirmed' ? 'Đã xác nhận' : 'Đã hủy'}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
            {activeTab === 'reviews' && (
              <div>
                <h3 className="text-lg font-semibold mb-4 text-emerald-700">Đánh giá của khách thuê</h3>
                <div className="space-y-4">
                  {fakeReviews.map((review) => (
                    <div key={review.id} className="border rounded-lg p-4 bg-white shadow-sm">
                      <div className="flex items-center mb-2">
                        <span className="font-medium text-gray-900 mr-2">{review.guestName}</span>
                        <span className="text-yellow-500 font-bold mr-2">{'★'.repeat(review.rating)}{'☆'.repeat(5 - review.rating)}</span>
                        <span className="text-xs text-gray-400">{review.createdAt}</span>
                      </div>
                      <div className="text-gray-700">{review.comment}</div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default RoomDetailPage;
