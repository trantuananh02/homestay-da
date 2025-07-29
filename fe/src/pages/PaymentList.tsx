import { useEffect, useState } from 'react';
import { Calendar, CreditCard, Banknote, Filter, X } from 'lucide-react';
import { Payment } from '../types';
import { bookingService } from '../services/bookingService';

function formatCurrency(amount: number) {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
  }).format(amount);
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  });
}

const methodLabels: Record<Payment['paymentMethod'], string> = {
  "Tiền mặt": 'Tiền mặt',
  "Thẻ tín dụng": 'Thẻ tín dụng',
  "Chuyển khoản": 'Chuyển khoản',
  "Momo": 'Momo',
};

function PaymentList() {
  const [payments, setPayments] = useState<Payment[]>([]);
  const [filters, setFilters] = useState({
    bookingCode: '',
    method: '',
    dateFrom: '',
    dateTo: '',
  });
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

  useEffect(() => {
    const loadPayments = async () => {
      const response = await bookingService.getPayments({
        ...filters,
        page: currentPage,
        pageSize: itemsPerPage,
      });

      console.log('Payments loaded:', response);

      setPayments(response.payments || []);
    };
    loadPayments();
  }, [filters, currentPage]);

  const totalPages = Math.ceil(payments?.length  / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const currentPayments = payments?.slice(startIndex, startIndex + itemsPerPage) || [];

  const handleFilterChange = (key: string, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    setCurrentPage(1);
  };

  const clearFilters = () => {
    setFilters({
      bookingCode: '',
      method: '',
      dateFrom: '',
      dateTo: '',
    });
    setCurrentPage(1);
  };

  return (
    <div className="min-h-screen bg-gray-50 py-6">
      <div className="max-w-7xl mx-auto px-4">
        {/* Filters */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
          <div className="flex items-center gap-2 mb-4">
            <Filter className="w-5 h-5 text-gray-600" />
            <h2 className="text-lg font-semibold text-gray-900">Bộ lọc</h2>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
            <input
              type="text"
              placeholder="Mã đặt phòng..."
              value={filters.bookingCode}
              onChange={(e) => handleFilterChange('bookingCode', e.target.value)}
              className="border-gray-300 rounded-lg px-4 py-2 border"
            />
            <select
              value={filters.method}
              onChange={(e) => handleFilterChange('method', e.target.value)}
              className="border-gray-300 rounded-lg px-4 py-2 border"
            >
              <option value="">Phương thức</option>
              {Object.keys(methodLabels).map((m) => (
                <option key={m} value={m}>{methodLabels[m as keyof typeof methodLabels]}</option>
              ))}
            </select>
            <input
              type="date"
              value={filters.dateFrom}
              onChange={(e) => handleFilterChange('dateFrom', e.target.value)}
              className="border-gray-300 rounded-lg px-4 py-2 border"
            />
            <input
              type="date"
              value={filters.dateTo}
              onChange={(e) => handleFilterChange('dateTo', e.target.value)}
              className="border-gray-300 rounded-lg px-4 py-2 border"
            />
            <button
              onClick={clearFilters}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 flex items-center justify-center"
            >
              Xóa bộ lọc
            </button>
          </div>
        </div>

        {/* Table */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-x-auto pb-32 min-h-[400px]">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">STT</th>
                {/* <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Mã thanh toán</th> */}
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Mã đặt phòng</th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Số tiền</th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Phương thức</th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Ngày thanh toán</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-100 ">
              {currentPayments.map((payment, idx) => (
                <tr key={payment.id} className="hover:bg-gray-50 transition">
                  <td className="px-4 py-3 text-sm text-gray-900">{startIndex + idx + 1}</td>
                  {/* <td className="px-4 py-3 text-sm font-medium text-blue-600 bg-blue-50 rounded">{payment.id}</td> */}
                  <td className="px-4 py-3 text-sm text-gray-700">{payment.bookingCode}</td>
                  <td className="px-4 py-3 text-sm font-semibold text-gray-800">
                    <div className="flex items-center gap-1">
                      <Banknote className="w-4 h-4 text-gray-400" />
                      {formatCurrency(payment.amount)}
                    </div>
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-700 ">
                    <div className="flex items-center gap-1">
                      <CreditCard className="w-4 h-4 text-gray-400" />
                      {payment.paymentMethod}
                    </div>
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-500">
                    <div className="flex items-center gap-1">
                      <Calendar className="w-4 h-4" />
                      {formatDate(payment.paymentDate)}
                    </div>
                  </td>
                </tr>
              ))}
              {!currentPayments || currentPayments?.length === 0 && (
                <tr>
                  <td colSpan={6} className="text-center text-gray-500 py-6 text-sm">
                    Không có dữ liệu thanh toán.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex justify-center mt-6">
            <div className="flex gap-2">
              <button
                disabled={currentPage === 1}
                onClick={() => setCurrentPage((p) => Math.max(p - 1, 1))}
                className="px-3 py-1 border rounded disabled:opacity-50"
              >
                Trước
              </button>
              {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
                <button
                  key={page}
                  onClick={() => setCurrentPage(page)}
                  className={`px-3 py-1 border rounded ${page === currentPage ? 'bg-blue-100 text-blue-700 font-semibold' : ''}`}
                >
                  {page}
                </button>
              ))}
              <button
                disabled={currentPage === totalPages}
                onClick={() => setCurrentPage((p) => Math.min(p + 1, totalPages))}
                className="px-3 py-1 border rounded disabled:opacity-50"
              >
                Sau
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default PaymentList;
