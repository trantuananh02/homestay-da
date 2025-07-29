import { useEffect, useState } from 'react';
import {
  Calendar,
  Filter,
  Eye,
  Edit,
  ChevronLeft,
  ChevronRight,
  Phone,
  MapPin,
  Clock,
  User,
  MoreVertical,
  Check,
  X,
  Plus,
  Book,
} from 'lucide-react';
import { Booking } from '../types';
import { bookingService } from '../services/bookingService';
import { useConfirm } from '../components/ConfirmDialog';
import BookingDetailModal from '../components/Booking/DetailBookingModal';

// const mockRooms: Room[] = [];

// const mockBookings: Booking[] = [];

function BookingList() {
  const confirm = useConfirm();
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [activeDropdown, setActiveDropdown] = useState<number | null>(null);
  const [filters, setFilters] = useState({
    customerName: '',
    customerPhone: '',
    dateFrom: '',
    dateTo: '',
    status: ''
  });
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;
  const [showDetailBooking, setShowDetailBooking] = useState(false);
  const [selectedBooking, setSelectedBooking] = useState<Booking | null>(null);

  useEffect(() => {
    const fetchBookings = async () => {
      try {
        const bookingList = await bookingService.filterBookings({
          customerName: filters.customerName,
          customerPhone: filters.customerPhone,
          dateFrom: filters.dateFrom,
          dateTo: filters.dateTo,
          status: filters.status,
          page: currentPage,
          pageSize: itemsPerPage
        });

        setBookings(bookingList.bookings || []);
      } catch (error) {
        console.error('Error fetching bookings:', error);
      }
    };

    fetchBookings();
  }, [filters, currentPage]);

  const statusColors = {
    pending: 'bg-yellow-100 text-yellow-800',
    confirmed: 'bg-blue-100 text-blue-800',
    cancelled: 'bg-red-100 text-red-800',
    completed: 'bg-green-100 text-green-800'
  };

  const statusLabels = {
    pending: 'Chờ xác nhận',
    confirmed: 'Đã xác nhận',
    cancelled: 'Đã hủy',
    completed: 'Hoàn thành'
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND'
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('vi-VN', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric'
    });
  };

  const filteredBookings = bookings.filter(booking => {
    return (
      booking.customerName.toLowerCase().includes(filters.customerName.toLowerCase()) &&
      booking.customerPhone.includes(filters.customerPhone) &&
      (!filters.dateFrom || booking.bookingDate >= filters.dateFrom) &&
      (!filters.dateTo || booking.bookingDate <= filters.dateTo) &&
      (!filters.status || booking.status === filters.status)
    );
  });

  const totalPages = Math.ceil(filteredBookings.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const currentBookings = filteredBookings.slice(startIndex, startIndex + itemsPerPage);

  const handleFilterChange = (key: string, value: string) => {
    setFilters(prev => ({ ...prev, [key]: value }));
    setCurrentPage(1);
  };

  const clearFilters = () => {
    setFilters({
      customerName: '',
      customerPhone: '',
      dateFrom: '',
      dateTo: '',
      status: ''
    });
    setCurrentPage(1);
  };

  const handleEditBooking = (booking: Booking) => {
    console.log('Edit booking', booking.id);

    setActiveDropdown(null);
  };

  const getAvailableActions = (booking: Booking) => {
    const actions = [];

    // Luôn có thao tác xem chi tiết
    actions.push({
      label: 'Xem chi tiết',
      icon: Eye,
      color: 'text-blue-600 hover:text-blue-800',
      action: () => {
        console.log('View booking', booking.id);
        
        setSelectedBooking(booking);
        setShowDetailBooking(true);
      }
    });

    // Thao tác theo trạng thái
    switch (booking.status) {
      case 'pending':
        actions.push(
          {
            label: 'Xác nhận đặt phòng',
            icon: Check,
            color: 'text-green-600 hover:text-green-800',
            action: async () => {
              var result = await confirm({
                title: 'Xác nhận đặt phòng',
                description: `Bạn có chắc chắn muốn đặt phòng này?`,
                confirmText: 'Xác nhận',
                cancelText: 'Không'
              });
              if (result) {
                await bookingService.updateBookingStatus(booking.id, 'confirmed');

                const bookingList = await bookingService.filterBookings({
                  customerName: filters.customerName,
                  customerPhone: filters.customerPhone,
                  dateFrom: filters.dateFrom,
                  dateTo: filters.dateTo,
                  status: filters.status,
                  page: currentPage,
                  pageSize: itemsPerPage
                });

                setBookings(bookingList.bookings || []);
              }
              setActiveDropdown(null);
            }
          },
          {
            label: 'Hủy đặt phòng',
            icon: X,
            color: 'text-red-600 hover:text-red-800',
            action: async () => {
              var result = await confirm({
                title: 'Xác nhận hủy đặt phòng',
                description: `Bạn có chắc chắn muốn hủy đặt phòng này?`,
                confirmText: 'Hủy',
                cancelText: 'Không'
              });
              if (result) {
                await bookingService.updateBookingStatus(booking.id, 'cancelled');

                const bookingList = await bookingService.filterBookings({
                  customerName: filters.customerName,
                  customerPhone: filters.customerPhone,
                  dateFrom: filters.dateFrom,
                  dateTo: filters.dateTo,
                  status: filters.status,
                  page: currentPage,
                  pageSize: itemsPerPage
                });

                setBookings(bookingList.bookings || []);
              }
              setActiveDropdown(null);
            }
          }
        );
        break;

      case 'confirmed':
        if (booking.paidAmount < booking.totalAmount) {
          actions.push({
            label: 'Xác nhận thanh toán',
            icon: Plus,
            color: 'text-green-600 hover:text-green-800',
            action: async () => {
              var result = await confirm({
                title: 'Xác nhận thanh toán',
                description: `Bạn có chắc chắn muốn xác nhận thanh toán cho đặt phòng này?`,
                confirmText: 'Xác nhận',
                cancelText: 'Không'
              });
              if (result) {
                await bookingService.updateBookingStatus(booking.id, 'completed');

                const bookingList = await bookingService.filterBookings({
                  customerName: filters.customerName,
                  customerPhone: filters.customerPhone,
                  dateFrom: filters.dateFrom,
                  dateTo: filters.dateTo,
                  status: filters.status,
                  page: currentPage,
                  pageSize: itemsPerPage
                });

                setBookings(bookingList.bookings || []);
              }
              setActiveDropdown(null);
            }
          });
        }
        actions.push({
          label: 'Hủy đặt phòng',
          icon: X,
          color: 'text-red-600 hover:text-red-800',
          action: async () => {
            var result = await confirm({
              title: 'Xác nhận hủy đặt phòng',
              description: `Bạn có chắc chắn muốn hủy đặt phòng này?`,
              confirmText: 'Hủy',
              cancelText: 'Không'
            });
            if (result) {
              await bookingService.updateBookingStatus(booking.id, 'cancelled');

              const bookingList = await bookingService.filterBookings({
                customerName: filters.customerName,
                customerPhone: filters.customerPhone,
                dateFrom: filters.dateFrom,
                dateTo: filters.dateTo,
                status: filters.status,
                page: currentPage,
                pageSize: itemsPerPage
              });

              setBookings(bookingList.bookings || []);
            }
            setActiveDropdown(null);
          }
        });
        break;

      case 'completed':
        actions.push({
          label: 'Chỉnh sửa',
          icon: Edit,
          color: 'text-blue-600 hover:text-blue-800',
          action: () => handleEditBooking(booking)
        });
        break;

      case 'cancelled':
        break;
    }

    return actions;
  };

  return (
    <div className="min-h-screen bg-gray-50" onClick={() => setActiveDropdown(null)}>
      <div className="max-w-7xl mx-auto">
        {/* Filters */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
          <div className="flex items-center gap-2 mb-4">
            <Filter className="w-5 h-5 text-gray-600" />
            <h2 className="text-lg font-semibold text-gray-900">Bộ lọc</h2>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Tên người đặt
              </label>
              <div className="relative">
                <User className="absolute left-3 top-3 w-4 h-4 text-gray-400" />
                <input
                  type="text"
                  placeholder="Nhập tên..."
                  value={filters.customerName}
                  onChange={(e) => handleFilterChange('customerName', e.target.value)}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Số điện thoại
              </label>
              <div className="relative">
                <Phone className="absolute left-3 top-3 w-4 h-4 text-gray-400" />
                <input
                  type="text"
                  placeholder="Nhập số điện thoại..."
                  value={filters.customerPhone}
                  onChange={(e) => handleFilterChange('customerPhone', e.target.value)}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Ngày đặt từ
              </label>
              <div className="relative">
                <Calendar className="absolute left-3 top-3 w-4 h-4 text-gray-400" />
                <input
                  type="date"
                  value={filters.dateFrom}
                  onChange={(e) => handleFilterChange('dateFrom', e.target.value)}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Ngày đặt đến
              </label>
              <div className="relative">
                <Calendar className="absolute left-3 top-3 w-4 h-4 text-gray-400" />
                <input
                  type="date"
                  value={filters.dateTo}
                  onChange={(e) => handleFilterChange('dateTo', e.target.value)}
                  className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Trạng thái
              </label>
              <select
                value={filters.status}
                onChange={(e) => handleFilterChange('status', e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">Tất cả trạng thái</option>
                <option value="pending">Chờ xác nhận</option>
                <option value="confirmed">Đã xác nhận</option>
                <option value="cancelled">Đã hủy</option>
                <option value="completed">Hoàn thành</option>
              </select>
            </div>
          </div>

          <div className="flex justify-end mt-4">
            <button
              onClick={clearFilters}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
            >
              Xóa bộ lọc
            </button>
          </div>
        </div>

        {/* Results Summary */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4 mb-6">
          <div className="flex items-center justify-between">
            <div className="text-sm text-gray-600">
              Hiển thị {startIndex + 1} - {Math.min(startIndex + itemsPerPage, filteredBookings.length)} của {filteredBookings.length} kết quả
            </div>
            <div className="text-sm font-medium text-gray-900">
              Tổng doanh thu: {formatCurrency(
                filteredBookings
                  .filter(booking => booking.status === 'completed' || booking.status === 'confirmed')
                  .reduce((sum, booking) => sum + booking.totalAmount, 0)
              )}
            </div>
          </div>
        </div>

        {/* Bookings Table */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <div className="overflow-x-auto pb-32 min-h-[400px]">
            <table className="w-full">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    STT
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Mã đặt phòng
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Thông tin người đặt
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Thông tin phòng
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Ngày nhận/trả
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Tổng giá trị
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Trạng thái
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Thao tác
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {currentBookings?.map((booking, index) => (
                  <tr key={booking.id} className="hover:bg-gray-50 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {startIndex + index + 1}
                    </td>
                    <td className="px-2 py-4 whitespace-nowrap">
                      <div className="text-sm font-medium text-blue-600 bg-blue-50 py-1 rounded">
                        {booking.bookingCode}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <div className="flex-shrink-0 h-10 w-10">
                          <div className="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center">
                            <User className="h-5 w-5 text-blue-600" />
                          </div>
                        </div>
                        <div className="ml-4">
                          <div className="text-sm font-medium text-gray-900">{booking.customerName}</div>
                          <div className="text-sm text-gray-500 flex items-center gap-1">
                            <Phone className="w-3 h-3" />
                            {booking.customerPhone}
                          </div>
                          <div className="text-sm text-gray-500">{booking.customerEmail}</div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="space-y-1">
                        {booking.rooms?.map((room, index) => (
                          <div key={index}>
                            <div className="text-sm font-medium text-gray-900">{room.name}</div>
                            <div className="text-sm text-gray-500 flex items-center gap-1">
                              <MapPin className="w-3 h-3" />
                              <span className={`px-2 py-1 rounded text-xs font-medium ${room.type === 'Standard' ? 'bg-gray-100 text-gray-800' :
                                room.type === 'Deluxe' ? 'bg-blue-100 text-blue-800' :
                                  room.type === 'Premium' ? 'bg-purple-100 text-purple-800' :
                                    'bg-yellow-100 text-yellow-800'
                                }`}>
                                {room.type}
                              </span>
                            </div>
                          </div>
                        ))}
                      </div>
                      <div className="text-sm text-gray-500 mt-1">
                        {booking.rooms.length} phòng • {booking.nights} đêm
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        <div className="flex items-center gap-1 mb-1">
                          <Calendar className="w-3 h-3 text-green-600" />
                          <span className="text-green-600">Nhận:</span> {formatDate(booking.checkIn)}
                        </div>
                        <div className="flex items-center gap-1">
                          <Calendar className="w-3 h-3 text-red-600" />
                          <span className="text-red-600">Trả:</span> {formatDate(booking.checkOut)}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        <div className="font-medium text-lg">{formatCurrency(booking.totalAmount)}</div>
                        <div className="text-xs text-gray-500 flex items-center gap-1 mt-1">
                          <span className={`px-2 py-1 rounded text-xs font-medium ${booking.paidAmount >= booking.totalAmount
                            ? 'bg-green-100 text-green-800'
                            : booking.paidAmount > 0
                              ? 'bg-yellow-100 text-yellow-800'
                              : 'bg-red-100 text-red-800'
                            }`}>
                            {booking.paidAmount >= booking.totalAmount
                              ? 'Đã thanh toán'
                              : booking.paidAmount > 0
                                ? `Còn lại ${formatCurrency(booking.totalAmount - booking.paidAmount)}`
                                : 'Chưa thanh toán'
                            }
                          </span>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${statusColors[booking.status]}`}>
                        {statusLabels[booking.status]}
                      </span>
                      <div className="text-xs text-gray-500 mt-1 flex items-center gap-1">
                        <Clock className="w-3 h-3" />
                        {formatDate(booking.bookingDate)}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <div className="relative">
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setActiveDropdown(activeDropdown === booking.id ? null : booking.id);
                          }}
                          className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full transition-colors"
                        >
                          <MoreVertical className="w-4 h-4" />
                        </button>
                        {activeDropdown === booking.id && (
                          <div className="absolute right-0 top-full mt-1 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-1 z-10">
                            {getAvailableActions(booking)?.map((action, actionIndex) => (
                              <button
                                key={actionIndex}
                                onClick={(e) => {
                                  e.stopPropagation();
                                  action.action();
                                }}
                                className={`w-full px-4 py-2 text-left text-sm flex items-center gap-2 hover:bg-gray-50 transition-colors ${action.color}`}
                              >
                                <action.icon className="w-4 h-4" />
                                {action.label}
                              </button>
                            ))}
                          </div>
                        )}

                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6 rounded-lg shadow-sm mt-4">
            <div className="flex flex-1 justify-between sm:hidden">
              <button
                onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                disabled={currentPage === 1}
                className="relative inline-flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
              >
                Trước
              </button>
              <button
                onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                disabled={currentPage === totalPages}
                className="relative ml-3 inline-flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
              >
                Sau
              </button>
            </div>
            <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
              <div>
                <p className="text-sm text-gray-700">
                  Trang <span className="font-medium">{currentPage}</span> của{' '}
                  <span className="font-medium">{totalPages}</span>
                </p>
              </div>
              <div>
                <nav className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                  <button
                    onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                    disabled={currentPage === 1}
                    className="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                  >
                    <ChevronLeft className="w-4 h-4" />
                  </button>

                  {Array.from({ length: totalPages }, (_, i) => i + 1)?.map(page => (
                    <button
                      key={page}
                      onClick={() => setCurrentPage(page)}
                      className={`relative inline-flex items-center px-4 py-2 border text-sm font-medium ${page === currentPage
                        ? 'z-10 bg-blue-50 border-blue-500 text-blue-600'
                        : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
                        }`}
                    >
                      {page}
                    </button>
                  ))}

                  <button
                    onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                    disabled={currentPage === totalPages}
                    className="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                  >
                    <ChevronRight className="w-4 h-4" />
                  </button>
                </nav>
              </div>
            </div>
          </div>
        )}
      </div>

      <BookingDetailModal
        isOpen={showDetailBooking}
        onClose={() => setShowDetailBooking(false)}
        booking={selectedBooking}
      />
    </div>
  );
}

export default BookingList;