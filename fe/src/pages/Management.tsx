import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, Edit, Trash2, Eye, DollarSign, Users, Building, RefreshCw, Power } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';
import { homestayService } from '../services/homestayService';
import { bookingService } from '../services/bookingService';
import { Homestay, HomestayStats, Booking } from '../types';
import BookingList from './BookingList';
import { useConfirm } from '../components/ConfirmDialog';
import PaymentList from './PaymentList';

const Management: React.FC = () => {
  const confirm = useConfirm();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [activeTab, setActiveTab] = useState('overview');
  const [homestays, setHomestays] = useState<Homestay[]>([]);
  const [stats, setStats] = useState<HomestayStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const [bookings, setBookings] = useState<Booking[]>([]);
  const [bookingLoading, setBookingLoading] = useState(false);

  // Load data
  const loadData = async () => {
    try {
      setLoading(true);
      const [homestayList, homestayStats] = await Promise.all([
        homestayService.getHomestayList(),
        homestayService.getHomestayStats()
      ]);
      setHomestays(homestayList.homestays);
      setStats(homestayStats);
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
    }
  };

  const refreshData = async () => {
    setRefreshing(true);
    await loadData();
    setRefreshing(false);
  };

  useEffect(() => {
    if (user?.role === 'host' || user?.role === 'admin') {
      loadData();
    }
  }, [user]);

  const loadBookings = async () => {
    setBookingLoading(true);
    try {
      let allBookings: Booking[] = [];
      for (const homestay of homestays) {
        const res = await bookingService.filterBookings(homestay.id.toString());
        allBookings = allBookings.concat(res);
      }

      setBookings(allBookings);
    } catch (error) {
      setBookings([]);
    } finally {
      setBookingLoading(false);
    }
  };

  useEffect(() => {
    if (activeTab === 'bookings') {
      loadBookings();
    }
    // eslint-disable-next-line
  }, [activeTab, homestays]);

  const handleAddHomestay = () => {
    navigate('/add-homestay');
  };

  const handleViewHomestay = (homestay: Homestay) => {
    navigate(`/management/homestay/${homestay.id}`);
  };

  const handleEditHomestay = (homestay: Homestay) => {
    navigate(`/management/homestay/${homestay.id}/edit`);
  };

  const handleDeleteHomestay = async (id: number) => {
    var result = await confirm({
      title: 'Xác nhận xóa homestay',
      description: `Bạn có chắc chắn muốn xóa homestay này?`,
      confirmText: 'Xóa',
      cancelText: 'Không'
    });
    if (result) {
      try {
        await homestayService.deleteHomestay(id);
        await loadData(); // Reload data after deletion
      } catch (error) {
        console.error('Error deleting homestay:', error);
      }
    }
  };

  const handleToggleStatus = async (homestay: Homestay) => {
    const action = homestay.status === 'active' ? 'tắt' : 'bật';
    
    var result = await confirm({
      title: `Xác nhận ${action} homestay`,
      description: `Bạn có chắc chắn muốn ${action} homestay "${homestay.name}"?`,
      confirmText: action === 'bật' ? 'Bật' : 'Tắt',
      cancelText: 'Không'
    });
    if (result) {
      try {
        await homestayService.toggleHomestayStatus(homestay.id);
        await loadData(); // Reload data after status change
      } catch (error) {
        console.error('Error toggling homestay status:', error);
      }
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Đang tải dữ liệu...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8 flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">Quản lý Homestay</h1>
            <p className="text-gray-600">Quản lý homestay và đặt phòng của bạn</p>
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
              onClick={handleAddHomestay}
              className="flex items-center px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-lg hover:bg-primary-700"
            >
              <Plus className="h-4 w-4 mr-2" />
              Thêm Homestay
            </button>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div className="bg-white p-6 rounded-xl shadow-sm">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">Homestay hoạt động</p>
                <p className="text-2xl font-bold text-primary-600">
                  {stats?.activeHomestays || 0}
                </p>
              </div>
              <div className="bg-primary-100 p-3 rounded-full">
                <Building className="h-6 w-6 text-primary-600" />
              </div>
            </div>
          </div>

          <div className="bg-white p-6 rounded-xl shadow-sm">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">Tổng homestay</p>
                <p className="text-2xl font-bold text-blue-600">
                  {stats?.totalHomestays || 0}
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
                  {stats?.availableRooms || 0}
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
                <p className="text-sm text-gray-600">Tổng doanh thu</p>
                <p className="text-2xl font-bold text-primary-600">
                  {homestayService.formatPrice(stats?.totalRevenue || 0)}
                </p>
              </div>
              <div className="bg-primary-100 p-3 rounded-full">
                <DollarSign className="h-6 w-6 text-primary-600" />
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
                className={`py-4 px-1 border-b-2 font-medium text-sm ${
                  activeTab === 'overview'
                    ? 'border-emerald-500 text-emerald-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                }`}
              >
                Tổng quan
              </button>
              <button
                onClick={() => setActiveTab('homestays')}
                className={`py-4 px-1 border-b-2 font-medium text-sm ${
                  activeTab === 'homestays'
                    ? 'border-emerald-500 text-emerald-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                }`}
              >
                Homestay của tôi ({homestays.length})
              </button>
              <button
                onClick={() => setActiveTab('bookings')}
                className={`py-4 px-1 border-b-2 font-medium text-sm ${
                  activeTab === 'bookings'
                    ? 'border-emerald-500 text-emerald-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                }`}
              >
                Đặt phòng ({stats?.totalBookings || 0})
              </button>
              <button
                onClick={() => setActiveTab('payments')}
                className={`py-4 px-1 border-b-2 font-medium text-sm ${
                  activeTab === 'payments'
                    ? 'border-emerald-500 text-emerald-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                }`}
              >
                Thanh toán ({homestayService.formatPrice(stats?.totalRevenue || 0)})
              </button>
            </nav>
          </div>

          <div className="p-6">
            {activeTab === 'overview' && (
              <div>
                <h2 className="text-xl font-semibold mb-4">Thống kê chi tiết</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="bg-gray-50 p-6 rounded-lg">
                    <h3 className="text-lg font-medium mb-4">Thông tin Homestay</h3>
                    <div className="space-y-3">
                      <div className="flex justify-between">
                        <span className="text-gray-600">Tổng homestay:</span>
                        <span className="font-medium">{stats?.totalHomestays || 0}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Homestay hoạt động:</span>
                        <span className="font-medium text-primary-600">{stats?.activeHomestays || 0}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Homestay không hoạt động:</span>
                        <span className="font-medium text-red-600">
                          {(stats?.totalHomestays || 0) - (stats?.activeHomestays || 0)}
                        </span>
                      </div>
                    </div>
                  </div>
                  <div className="bg-gray-50 p-6 rounded-lg">
                    <h3 className="text-lg font-medium mb-4">Thông tin Phòng</h3>
                    <div className="space-y-3">
                      <div className="flex justify-between">
                        <span className="text-gray-600">Tổng phòng:</span>
                        <span className="font-medium">{stats?.totalRooms || 0}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Phòng có sẵn:</span>
                        <span className="font-medium text-primary-600">{stats?.availableRooms || 0}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Tổng đặt phòng:</span>
                        <span className="font-medium">{stats?.totalBookings || 0}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'homestays' && (
              <div>
                <div className="flex justify-between items-center mb-6">
                  <h2 className="text-xl font-semibold">Danh sách Homestay</h2>
                  <button
                    onClick={handleAddHomestay}
                    className="flex items-center px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-lg hover:bg-primary-700"
                  >
                    <Plus className="h-4 w-4 mr-2" />
                    Thêm Homestay
                  </button>
                </div>

                {homestays.length === 0 ? (
                  <div className="text-center py-12">
                    <Building className="h-16 w-16 text-gray-400 mx-auto mb-4" />
                    <h3 className="text-lg font-medium text-gray-900 mb-2">Chưa có homestay nào</h3>
                    <p className="text-gray-600 mb-6">Bắt đầu bằng cách tạo homestay đầu tiên của bạn</p>
                    <button
                      onClick={handleAddHomestay}
                      className="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-lg hover:bg-primary-700"
                    >
                      <Plus className="h-4 w-4 mr-2" />
                      Tạo Homestay
                    </button>
                  </div>
                ) : (
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {homestays.map((homestay) => (
                      <div key={homestay.id} className="bg-white border border-gray-200 rounded-lg overflow-hidden shadow-sm hover:shadow-md transition-shadow">
                        <div className="p-6">
                          <div className="flex justify-between items-start mb-4">
                            <h3 className="text-lg font-semibold text-gray-900 truncate">
                              {homestay.name}
                            </h3>
                            <span className={`px-2 py-1 text-xs font-medium rounded-full ${homestayService.getStatusColor(homestay.status)}`}>
                              {homestayService.formatStatus(homestay.status)}
                            </span>
                          </div>
                          
                          <p className="text-gray-600 text-sm mb-4 line-clamp-2">
                            {homestay.description}
                          </p>
                          
                          <div className="space-y-2 mb-4">
                            <div className="flex items-center text-sm text-gray-600">
                              <Building className="h-4 w-4 mr-2" />
                              {homestay.address}, {homestay.ward}, {homestay.district}, {homestay.city}
                            </div>
                            <div className="flex items-center text-sm text-gray-600">
                              <Users className="h-4 w-4 mr-2" />
                              {homestay.rooms?.length || 0} phòng
                            </div>
                          </div>

                          <div className="grid grid-cols-2 gap-2">
                            <button
                              onClick={() => handleViewHomestay(homestay)}
                              className="flex items-center justify-center px-3 py-2 text-sm font-medium text-primary-600 bg-primary-50 rounded-lg hover:bg-primary-100"
                            >
                              <Eye className="h-4 w-4 mr-1" />
                              Xem
                            </button>
                            <button
                              onClick={() => handleEditHomestay(homestay)}
                              className="flex items-center justify-center px-3 py-2 text-sm font-medium text-blue-600 bg-blue-50 rounded-lg hover:bg-blue-100"
                            >
                              <Edit className="h-4 w-4 mr-1" />
                              Sửa
                            </button>
                            <button
                              onClick={() => handleToggleStatus(homestay)}
                              className="flex items-center justify-center px-3 py-2 text-sm font-medium text-yellow-600 bg-yellow-50 rounded-lg hover:bg-yellow-100"
                            >
                              <Power className="h-4 w-4 mr-1" />
                              {homestay.status === 'active' ? 'Tắt' : 'Bật'}
                            </button>
                            <button
                              onClick={() => handleDeleteHomestay(homestay.id)}
                              className="flex items-center justify-center px-3 py-2 text-sm font-medium text-red-600 bg-red-50 rounded-lg hover:bg-red-100"
                            >
                              <Trash2 className="h-4 w-4 mr-1" />
                              Xóa
                            </button>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}

            {activeTab === 'bookings' && (
              <BookingList />
            )}

            {activeTab === 'payments' && (
              <PaymentList />
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Management;