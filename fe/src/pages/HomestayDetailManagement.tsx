import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import { ArrowLeft, Edit, Trash2, Plus, Eye, Calendar, DollarSign, Users, Building, RefreshCw, MapPin, Star } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';
import { homestayService } from '../services/homestayService';
import { Booking, Homestay, HomestayStats, Room, RoomStats } from '../types';
import { useConfirm } from '../components/ConfirmDialog';
import NewBookingModal from '../components/Booking/NewBookingModal';
import DEFAULT_ROOM_IMAGE from '../asset/default-image.jpg';
import AddRoomModal from '../components/Room/AddRoomModal';
import ViewRoomModal from '../components/Room/ViewRoomModal';
import { bookingService } from '../services/bookingService';

const HomestayDetailManagement: React.FC = () => {
  const confirm = useConfirm();
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useAuth();

  const [showNewBookingForm, setShowNewBookingForm] = useState(false);
  const [showAddRoomModal, setShowAddRoomModal] = useState(false);
  const [showEditRoomModal, setShowEditRoomModal] = useState(false);
  const [roomData, setRoomData] = useState<Room | null>(null);
  const [isEdit, setIsEdit] = useState(false);
  const [action, setAction] = useState<'add' | 'edit' | 'view'>('add');

  const [bookings, setBookings] = useState<Booking[]>([]); // Danh sách đặt phòng
  const [homestay, setHomestay] = useState<Homestay | null>(null);
  const [homestayStats, setHomestayStats] = useState<HomestayStats | null>(null);
  const [roomStats, setRoomStats] = useState<RoomStats | null>(null);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [roomsLoading, setRoomsLoading] = useState(false);
  const [activeTab, setActiveTab] = useState(location.state?.activeTab || 'stats');

  const homestayId = parseInt(id || '0');

  const loadData = async () => {
    if (!homestayId) return;

    try {
      setLoading(true);
      const [homestayDetail, stats, roomStatsData] = await Promise.all([
        homestayService.getHomestayById(homestayId),
        homestayService.getHomestayStatsById(homestayId),
        homestayService.getRoomStats(homestayId)
      ]);

      setHomestay(homestayDetail.homestay);
      setHomestayStats(stats);
      setRoomStats(roomStatsData);
    } catch (error) {
      console.error('Error loading homestay data:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadRooms = async () => {
    if (!homestayId) return;

    try {
      setRoomsLoading(true);
      const roomList = await homestayService.getRoomList({
        homestayId: homestayId,
        page: 1,
        pageSize: 100 // Lấy tất cả phòng
      });
      setRooms(roomList.rooms || []);

      const bookingList = await homestayService.getBookingsByHomestayId(homestayId);

      setBookings(bookingList.bookings || []);
    } catch (error) {
      console.error('Error loading rooms:', error);
    } finally {
      setRoomsLoading(false);
    }
  };

  const refreshData = async () => {
    setRefreshing(true);
    await Promise.all([loadData(), loadRooms()]);
    setRefreshing(false);
  };

  useEffect(() => {
    if (user?.role === 'host' || user?.role === 'admin') {
      loadData();
      loadRooms();
    }
  }, [user, homestayId]);

  // Reload rooms when returning from add room page
  useEffect(() => {
    if (location.state?.refreshRooms) {
      loadRooms();
      // Clear the state to prevent infinite reload
      navigate(location.pathname, { replace: true, state: { activeTab: location.state.activeTab } });
    }
  }, [location.state?.refreshRooms]);

  const handleBack = () => {
    navigate('/management');
  };

  const handleEditHomestay = () => {
    navigate(`/management/homestay/${id}/edit`);
  };

  const handleAddRoomClose = () => {
    setShowAddRoomModal(false);
  };

  const handleAddRoomSubmit = async (room: any) => {
    if (action === 'edit' && !roomData) return;

    if (action === 'edit') {
      // Cập nhật phòng
      console.log("Updating room:", roomData?.id, room);
      try {
        console.log("Submitting room update:", roomData?.id, room);


        await homestayService.updateRoom( roomData?.id || 0, {
          name: room.name,
          description: room.description,
          type: room.type,
          capacity: room.capacity,
          price: room.price,
          priceType: room.priceType || 'per_night',
          amenities: room.amenities,
          images: room.images,
        });
        // Chuyển về trang quản lý homestay với state để reload danh sách phòng
        navigate(`/management/homestay/${id}`, {
          state: {
            activeTab: 'rooms',
            refreshRooms: true
          }
        });
      } catch (error) {
        // Có thể show toast lỗi ở đây nếu muốn
      }
      return;
    }


    if (!id) return;

    console.log("Submitting room:", room);

    try {
      await homestayService.createRoom({
        homestayId: Number(id),
        name: room.name,
        description: room.description,
        type: room.type,
        capacity: room.capacity,
        price: room.price,
        priceType: room.priceType || 'per_night',
        amenities: room.amenities,
        images: room.images,
      });
      // Chuyển về trang quản lý homestay với state để reload danh sách phòng
      navigate(`/management/homestay/${id}`, {
        state: {
          activeTab: 'rooms',
          refreshRooms: true
        }
      });
    } catch (error) {
      // Có thể show toast lỗi ở đây nếu muốn
    }
  };

  const handleDeleteHomestay = async () => {
    if (!homestay) return;

    var result = await confirm({
      title: 'Xác nhận xóa homestay',
      description: `Bạn có chắc chắn muốn xóa homestay "${homestay.name}"?`,
      confirmText: 'Xóa',
      cancelText: 'Không'
    });
    if (result) {
      try {
        await homestayService.deleteHomestay(homestayId);
        navigate('/management');
      } catch (error) {
        console.error('Error deleting homestay:', error);
      }
    }
  };

  const handleAddRoom = () => {
    setShowAddRoomModal(true);
  };

  const handleViewRoom = async (roomId: number) => {
    console.log("Viewing room:", roomId);

    const roomDetail = await homestayService.getRoomById(roomId);
    setRoomData({
      id: roomDetail.room.id,
      name: roomDetail.room.name,
      description: roomDetail.room.description,
      type: roomDetail.room.type,
      capacity: roomDetail.room.capacity,
      price: roomDetail.room.price,
      priceType: roomDetail.room.priceType,
      amenities: roomDetail.room.amenities,
      images: roomDetail.room.images,
    } as Room);

    setAction('view');
    setShowAddRoomModal(true);
  };

  const handleRefreshRooms = async () => {
    await loadRooms();
  };

  const handleCreateBooking = (bookingData: Omit<Booking, 'id' | 'bookingCode' | 'status' | 'bookingDate' | 'nights'>) => {
    console.log("Creating booking:", bookingData);

    bookingService.createBooking({
      ...bookingData,
      homestayId: homestayId,
    })
  };

  const handleEditRoom = async (roomId: number) => {
    console.log("Viewing room:", roomId);

    const roomDetail = await homestayService.getRoomById(roomId);
    setRoomData({
      id: roomDetail.room.id,
      name: roomDetail.room.name,
      description: roomDetail.room.description,
      type: roomDetail.room.type,
      capacity: roomDetail.room.capacity,
      price: roomDetail.room.price,
      priceType: roomDetail.room.priceType,
      amenities: roomDetail.room.amenities,
      images: roomDetail.room.images,
    } as Room);

    setAction('edit');
    setShowAddRoomModal(true);
  }

  const handleDeleteRoom = (roomId: number) => {
    confirm({
      title: 'Xác nhận xóa phòng',
      description: 'Bạn có chắc chắn muốn xóa phòng này?',
      confirmText: 'Xóa',
      cancelText: 'Không'
    }).then(async (result) => {
      if (result) {
        try {
          await homestayService.deleteRoom(roomId);
          setRooms(prev => prev.filter(room => room.id !== roomId));
        } catch (error) {
          console.error('Error deleting room:', error);
        }
      }
    });
  }


  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Đang tải thông tin homestay...</p>
        </div>
      </div>
    );
  }

  if (!homestay) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600">Không tìm thấy homestay</p>
          <button
            onClick={handleBack}
            className="mt-4 px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700"
          >
            Quay lại
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <button
            onClick={handleBack}
            className="flex items-center space-x-2 text-emerald-600 hover:text-emerald-700 mb-4"
          >
            <ArrowLeft className="h-5 w-5" />
            <span>Quay lại</span>
          </button>

          <div className="flex justify-between items-start">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">{homestay.name}</h1>
              <div className="flex items-center text-gray-600 mb-4">
                <MapPin className="h-4 w-4 mr-1" />
                <span>{homestay.address}, {homestay.ward}, {homestay.district}, {homestay.city}</span>
              </div>
              <div className="flex items-center space-x-4 text-sm text-gray-600">
                <div className="flex items-center">
                  <Star className="h-4 w-4 text-yellow-400 fill-current mr-1" />
                  <span>{homestay.rating || 0} ({homestay.totalReviews || 0} đánh giá)</span>
                </div>
                <span className={`px-2 py-1 rounded-full text-xs ${homestayService.getStatusColor(homestay.status)}`}>
                  {homestayService.formatStatus(homestay.status)}
                </span>
              </div>
            </div>
            <div className="flex space-x-3">
              <button
                onClick={refreshData}
                disabled={refreshing}
                className="flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50"
              >
                <RefreshCw className={`h-4 w-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
                Làm mới
              </button>
              <button
                onClick={handleEditHomestay}
                className="flex items-center px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700"
              >
                <Edit className="h-4 w-4 mr-2" />
                Chỉnh sửa
              </button>
              <button
                onClick={() => setShowNewBookingForm(true)}
                className="px-6 py-3 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors font-medium flex items-center gap-2"
              >
                <Plus className="w-4 h-4" />
                Tạo đặt phòng mới
              </button>
              <button
                onClick={handleDeleteHomestay}
                className="flex items-center px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700"
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Xóa
              </button>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div className="bg-white p-6 rounded-xl shadow-sm">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">Tổng phòng</p>
                <p className="text-2xl font-bold text-blue-600">
                  {homestayStats?.totalRooms || 0}
                </p>
              </div>
              <div className="bg-blue-100 p-3 rounded-full">
                <Building className="h-6 w-6 text-blue-600" />
              </div>
            </div>
          </div>

          <div className="bg-white p-6 rounded-xl shadow-sm">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">Phòng có sẵn</p>
                <p className="text-2xl font-bold text-primary-600">
                  {homestayStats?.availableRooms || 0}
                </p>
              </div>
              <div className="bg-primary-100 p-3 rounded-full">
                <Users className="h-6 w-6 text-primary-600" />
              </div>
            </div>
          </div>

          <div className="bg-white p-6 rounded-xl shadow-sm">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">Tổng đặt phòng</p>
                <p className="text-2xl font-bold text-emerald-600">
                  {homestayStats?.totalBookings || 0}
                </p>
              </div>
              <div className="bg-emerald-100 p-3 rounded-full">
                <Calendar className="h-6 w-6 text-emerald-600" />
              </div>
            </div>
          </div>

          <div className="bg-white p-6 rounded-xl shadow-sm">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">Doanh thu</p>
                <p className="text-2xl font-bold text-emerald-600">
                  {homestayService.formatPrice(homestayStats?.totalRevenue || 0)}
                </p>
              </div>
              <div className="bg-emerald-100 p-3 rounded-full">
                <DollarSign className="h-6 w-6 text-emerald-600" />
              </div>
            </div>
          </div>
        </div>

        {/* Tab Navigation */}
        <div className="bg-white rounded-xl shadow-sm mb-8">
          <div className="border-b border-gray-200">
            <nav className="flex space-x-8 px-6">
              <button
                onClick={() => setActiveTab('stats')}
                className={`py-4 px-1 border-b-2 font-medium text-sm ${activeTab === 'stats'
                    ? 'border-emerald-500 text-emerald-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                  }`}
              >
                Tổng quan
              </button>
              <button
                onClick={() => setActiveTab('rooms')}
                className={`py-4 px-1 border-b-2 font-medium text-sm ${activeTab === 'rooms'
                    ? 'border-emerald-500 text-emerald-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                  }`}
              >
                Phòng ({rooms.length})
              </button>

              <button
                onClick={() => setActiveTab('reviews')}
                className={`py-4 px-1 border-b-2 font-medium text-sm ${activeTab === 'reviews'
                    ? 'border-emerald-500 text-emerald-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                  }`}
              >
                Đánh giá ({homestay.totalReviews || 0})
              </button>
            </nav>
          </div>

          <div className="p-6">
            {activeTab === 'rooms' && (
              <div>
                <div className="flex justify-between items-center mb-6">
                  <h2 className="text-xl font-semibold">Danh sách Phòng</h2>
                  <div className="flex space-x-3">
                    <button
                      onClick={handleRefreshRooms}
                      disabled={roomsLoading}
                      className="flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50"
                    >
                      <RefreshCw className={`h-4 w-4 mr-2 ${roomsLoading ? 'animate-spin' : ''}`} />
                      Làm mới
                    </button>
                    <button
                      onClick={handleAddRoom}
                      className="flex items-center px-4 py-2 text-sm font-medium text-white bg-emerald-600 rounded-lg hover:bg-emerald-700"
                    >
                      <Plus className="h-4 w-4 mr-2" />
                      Thêm Phòng
                    </button>
                  </div>
                </div>

                {roomsLoading ? (
                  <div className="text-center py-12">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600 mx-auto mb-4"></div>
                    <p className="text-gray-600">Đang tải danh sách phòng...</p>
                  </div>
                ) : rooms.length === 0 ? (
                  <div className="text-center py-12">
                    <Building className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                    <p className="text-gray-600 mb-4">Chưa có phòng nào</p>
                    <button
                      onClick={handleAddRoom}
                      className="px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700"
                    >
                      Thêm phòng đầu tiên
                    </button>
                  </div>
                ) : (
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {rooms.map((room) => (
                      <div key={room.id} className="bg-gray-50 rounded-lg p-6">
                        {/* Ảnh đại diện phòng */}
                        {room.images && room.images.length > 0 ? (
                          <img
                            src={room.images[0]}
                            alt={room.name}
                            className="w-full h-40 object-cover rounded-md mb-4"
                          />
                        ) : (
                          <img
                            src={DEFAULT_ROOM_IMAGE}
                            alt="Ảnh mặc định"
                            className="w-full h-40 object-cover rounded-md mb-4"
                          />
                        )}
                        <div className="flex justify-between items-start mb-4">
                          <h3 className="text-lg font-medium">{room.name}</h3>
                          <span className={`px-2 py-1 rounded-full text-xs ${homestayService.getStatusColor(room.status)}`}>
                            {homestayService.formatStatus(room.status)}
                          </span>
                        </div>
                        <div className="space-y-2 mb-4">
                          <p className="text-sm text-gray-600">{room.description}</p>
                          <p className="text-lg font-bold text-emerald-600">
                            {homestayService.formatPrice(room.price)}
                          </p>
                          <p className="text-sm text-gray-500">
                            {homestayService.formatPriceType(room.priceType)}
                          </p>
                        </div>
                        <div className="flex space-x-2">
                          <button
                            onClick={() => handleViewRoom(room.id)}
                            className="flex-1 flex items-center justify-center px-3 py-2 text-sm text-emerald-600 border border-emerald-600 rounded-lg hover:bg-emerald-50"
                          >
                            <Eye className="h-4 w-4 mr-1" />
                            Xem
                          </button>
                          <button
                            onClick={() => handleEditRoom(room.id)}
                            className="flex-1 flex items-center justify-center px-3 py-2 text-sm text-blue-600 border border-blue-600 rounded-lg hover:bg-blue-50"
                          >
                            <Edit className="h-4 w-4 mr-1" />
                            Sửa
                          </button>
                          <button
                            onClick={() => handleDeleteRoom && handleDeleteRoom(room.id)}
                            className="flex-1 flex items-center justify-center px-3 py-2 text-sm text-red-600 border border-red-600 rounded-lg hover:bg-red-50"
                          >
                            <Trash2 className="h-4 w-4 mr-1" />
                            Xóa
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}

            {activeTab === 'stats' && (
              <div className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="bg-gray-50 p-6 rounded-lg">
                    <h3 className="text-lg font-medium mb-4">Thống kê Homestay</h3>
                    <div className="space-y-3">
                      <div className="flex justify-between">
                        <span className="text-gray-600">Tổng phòng:</span>
                        <span className="font-medium">{homestayStats?.totalRooms || 0}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Phòng có sẵn:</span>
                        <span className="font-medium text-primary-600">{homestayStats?.availableRooms || 0}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Phòng đã đặt:</span>
                        <span className="font-medium text-red-600">
                          {(homestayStats?.totalRooms || 0) - (homestayStats?.availableRooms || 0)}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Tổng đặt phòng:</span>
                        <span className="font-medium">{homestayStats?.totalBookings || 0}</span>
                      </div>
                    </div>
                  </div>
                  <div className="bg-gray-50 p-6 rounded-lg">
                    <h3 className="text-lg font-medium mb-4">Thống kê Doanh thu</h3>
                    <div className="space-y-3">
                      <div className="flex justify-between">
                        <span className="text-gray-600">Tổng doanh thu:</span>
                        <span className="font-medium text-emerald-600">
                          {homestayService.formatPrice(homestayStats?.totalRevenue || 0)}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Doanh thu tháng này:</span>
                        <span className="font-medium text-emerald-600">
                          {homestayService.formatPrice(homestayStats?.monthlyRevenue || 0)}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Tỷ lệ lấp đầy:</span>
                        <span className="font-medium">
                          {homestayStats?.occupancyRate ? `${homestayStats.occupancyRate}%` : '0%'}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>

                {roomStats && (
                  <div className="bg-gray-50 p-6 rounded-lg">
                    <h3 className="text-lg font-medium mb-4">Thống kê Phòng</h3>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      <div className="text-center">
                        <p className="text-2xl font-bold text-blue-600">{roomStats.totalRooms}</p>
                        <p className="text-sm text-gray-600">Tổng phòng</p>
                      </div>
                      <div className="text-center">
                        <p className="text-2xl font-bold text-primary-600">{roomStats.availableRooms}</p>
                        <p className="text-sm text-gray-600">Phòng có sẵn</p>
                      </div>
                      <div className="text-center">
                        <p className="text-2xl font-bold text-emerald-600">
                          {homestayService.formatPrice(roomStats.averagePrice || 0)}
                        </p>
                        <p className="text-sm text-gray-600">Giá trung bình</p>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            )}

            {activeTab === 'reviews' && (
              <div className="space-y-6">
                <h2 className="text-xl font-semibold mb-4">Đánh giá Homestay</h2>
                {homestay.totalReviews && homestay.totalReviews > 0 ? (
                  <div className="space-y-4">
                    {homestay.reviews?.map((review) => (
                      <div key={review.id} className="bg-white p-4 rounded-lg shadow-sm">
                        <div className="flex items-center space-x-2 mb-2">
                          <Star className="h-5 w-5 text-yellow-400" />
                          <span className="font-medium">{review.rating}</span>
                        </div>
                        <p className="text-gray-700">{review.comment}</p>
                        <p className="text-sm text-gray-500 mt-2">Bởi {review.guestName} vào {new Date(review.createdAt ?? '').toLocaleDateString()}</p>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-gray-600">Chưa có đánh giá nào cho homestay này.</p>
                )}
              </div>
            )}


          </div>
        </div>
      </div>

      <NewBookingModal
        isOpen={showNewBookingForm}
        onClose={() => setShowNewBookingForm(false)}
        onCreateBooking={handleCreateBooking}
        rooms={rooms}
        existingBookings={bookings}
      />

      <AddRoomModal
        isOpen={showAddRoomModal}
        onClose={handleAddRoomClose}
        onSubmit={handleAddRoomSubmit}
        homestayId={id ? Number(id) : 0}
        room={roomData}
        action={action}
      />

      {/* <ViewRoomModal 
        isOpen={showEditRoomModal}
        onClose={() => setShowEditRoomModal(false)}
        room={roomData}
        isEdit={isEdit}
      /> */}
    </div>
  );
};

export default HomestayDetailManagement;