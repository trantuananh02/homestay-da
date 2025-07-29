import React, { useEffect, useState } from 'react';
import {
    X,
    User,
    Home,
    Calendar,
    DollarSign
} from 'lucide-react';
import { Booking } from '../../types';

interface BookingDetailModalProps {
    isOpen: boolean;
    onClose: () => void;
    booking: Booking | null;
}

const BookingDetailModal: React.FC<BookingDetailModalProps> = ({
    isOpen,
    onClose,
    booking
}) => {
    console.log('BookingDetailModal props:', { isOpen, booking });

    const [selectedRooms, setSelectedRooms] = useState<{
        id: number;
        name: string;
        type: string;
        pricePerNight: number;
        capacity: number;
    }[]>(
        (booking?.rooms
            ? booking?.rooms?.map(room => ({
                id: room.id,
                name: room.name,
                type: room.type,
                pricePerNight: room.pricePerNight,
                capacity: (room as any).capacity ?? 1 // default to 1 if missing
            }))
            : [])
    );
    const [newBooking, setNewBooking] = useState({
        customerName: booking?.customerName || '',
        customerPhone: booking?.customerPhone || '',
        customerEmail: booking?.customerEmail || '',
        checkIn: booking?.checkIn || '',
        checkOut: booking?.checkOut || '',
        paidAmount: booking?.paidAmount || 0,
        paymentMethod: booking?.paymentMethod || 'Tiền mặt'
    });

    useEffect(() => {
        if (booking) {
            setNewBooking({
                customerName: booking.customerName,
                customerPhone: booking.customerPhone,
                customerEmail: booking.customerEmail,
                checkIn: booking.checkIn,
                checkOut: booking.checkOut,
                paidAmount: booking.paidAmount,
                paymentMethod: booking.paymentMethod
            });
            setSelectedRooms(
                booking?.rooms?.map(room => ({
                    id: room.id,
                    name: room.name,
                    type: room.type,
                    pricePerNight: room.pricePerNight,
                    capacity: (room as any).capacity ?? 1 // default to 1 if missing
                }))
            );
        }
    }, [booking]);

    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat('vi-VN', {
            style: 'currency',
            currency: 'VND'
        }).format(amount);
    };

    const calculateNights = () => {
        if (newBooking.checkIn && newBooking.checkOut) {
            const checkIn = new Date(newBooking.checkIn);
            const checkOut = new Date(newBooking.checkOut);
            const diffTime = checkOut.getTime() - checkIn.getTime();
            const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
            return diffDays > 0 ? diffDays : 0;
        }
        return 0;
    };

    const calculateTotalAmount = () => {
        const nights = calculateNights();
        return selectedRooms?.reduce((total, room) => total + (room.pricePerNight * nights), 0);
    };

    const handleDateChange = (field: 'checkIn' | 'checkOut', value: string) => {
        setNewBooking(prev => ({ ...prev, [field]: value }));
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
    };

    const handleClose = () => {
        onClose();
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
                {/* Modal Header */}
                <div className="flex items-center justify-between p-6 border-b border-gray-200">
                    <h2 className="text-2xl font-bold text-gray-900">Chi tiết đặt phòng</h2>
                    <button
                        onClick={handleClose}
                        className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full transition-colors"
                    >
                        <X className="w-6 h-6" />
                    </button>
                </div>

                {/* Modal Content */}
                <div className="p-6">
                    <form onSubmit={handleSubmit} className="space-y-6">
                        {/* Customer Information */}
                        <div className="bg-gray-50 rounded-lg p-4">
                            <div className="flex items-center gap-2 mb-4">
                                <User className="w-5 h-5 text-blue-600" />
                                <h3 className="text-lg font-semibold text-gray-900">Thông tin khách hàng</h3>
                            </div>
                            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-2">
                                        Tên khách hàng *
                                    </label>
                                    <input
                                        type="text"
                                        required
                                        value={newBooking.customerName}
                                        onChange={(e) => setNewBooking(prev => ({ ...prev, customerName: e.target.value }))}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                        placeholder="Nhập tên khách hàng"
                                        disabled={true}
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-2">
                                        Số điện thoại *
                                    </label>
                                    <input
                                        type="tel"
                                        required
                                        value={newBooking.customerPhone}
                                        onChange={(e) => setNewBooking(prev => ({ ...prev, customerPhone: e.target.value }))}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                        placeholder="Nhập số điện thoại"
                                        disabled={true}
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-2">
                                        Email
                                    </label>
                                    <input
                                        type="email"
                                        value={newBooking.customerEmail}
                                        onChange={(e) => setNewBooking(prev => ({ ...prev, customerEmail: e.target.value }))}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                        placeholder="Nhập email"
                                        disabled={true}
                                    />
                                </div>
                            </div>
                        </div>

                        {/* Stay Duration */}
                        <div className="bg-gray-50 rounded-lg p-4">
                            <div className="flex items-center gap-2 mb-4">
                                <Calendar className="w-5 h-5 text-purple-600" />
                                <h3 className="text-lg font-semibold text-gray-900">
                                    Thời gian lưu trú
                                    <span className="text-sm font-normal text-gray-600 ml-2">
                                        (Chọn trước để xem phòng có sẵn)
                                    </span>
                                </h3>
                            </div>
                            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-2">
                                        Ngày nhận phòng *
                                    </label>
                                    <input
                                        type="date"
                                        required
                                        value={newBooking.checkIn}
                                        onChange={(e) => handleDateChange('checkIn', e.target.value)}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                        disabled={true}
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-2">
                                        Ngày trả phòng *
                                    </label>
                                    <input
                                        type="date"
                                        required
                                        value={newBooking.checkOut}
                                        onChange={(e) => handleDateChange('checkOut', e.target.value)}
                                        min={newBooking.checkIn}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                        disabled={true}
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-2">
                                        Số đêm
                                    </label>
                                    <input
                                        type="number"
                                        value={calculateNights()}
                                        readOnly
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-100 text-gray-600"
                                    />
                                </div>
                            </div>
                        </div>

                        {/* Room Information */}
                        <div className="bg-gray-50 rounded-lg p-4">
                            <div className="flex items-center gap-2 mb-4">
                                <Home className="w-5 h-5 text-green-600" />
                                <h3 className="text-lg font-semibold text-gray-900">Danh sách phòng</h3>
                            </div>

                            {/* Selected Rooms */}
                            {selectedRooms?.length > 0 && (
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-2">
                                        Phòng đã chọn ({selectedRooms.length})
                                    </label>
                                    <div className="space-y-2 max-h-40 overflow-y-auto">
                                        {selectedRooms?.map((room) => (
                                            <div key={room.id} className="flex items-center justify-between p-3 bg-white border border-gray-200 rounded-lg">
                                                <div className="flex-1">
                                                    <div className="font-medium text-gray-900">{room.name}</div>
                                                    <div className="text-sm text-gray-500 flex items-center gap-2">
                                                        <span className={`px-2 py-1 rounded text-xs font-medium ${room.type === 'Standard' ? 'bg-gray-100 text-gray-800' :
                                                                room.type === 'Deluxe' ? 'bg-blue-100 text-blue-800' :
                                                                    room.type === 'Premium' ? 'bg-purple-100 text-purple-800' :
                                                                        'bg-yellow-100 text-yellow-800'
                                                            }`}>
                                                            {room.type}
                                                        </span>
                                                        <span>•</span>
                                                        <span className="font-medium text-green-600">
                                                            {formatCurrency(room.pricePerNight)}/đêm
                                                        </span>
                                                        {calculateNights() > 0 && (
                                                            <>
                                                                <span>•</span>
                                                                <span className="font-medium text-blue-600">
                                                                    {formatCurrency(room.pricePerNight * calculateNights())} ({calculateNights()} đêm)
                                                                </span>
                                                            </>
                                                        )}
                                                    </div>
                                                </div>
                                                {/* <button
                                                    type="button"
                                                    onClick={() => handleRemoveRoom(room.id)}
                                                    className="ml-2 p-1 text-red-600 hover:text-red-800 hover:bg-red-50 rounded-full transition-colors"
                                                >
                                                    <X className="w-4 h-4" />
                                                </button> */}
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            )}
                        </div>

                        {/* Payment Information */}
                        <div className="bg-gray-50 rounded-lg p-4">
                            <div className="flex items-center gap-2 mb-4">
                                <DollarSign className="w-5 h-5 text-orange-600" />
                                <h3 className="text-lg font-semibold text-gray-900">Thông tin thanh toán</h3>
                            </div>
                            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-2">
                                        Tổng tiền *
                                    </label>
                                    <input
                                        type="number"
                                        min="0"
                                        value={calculateTotalAmount()}
                                        readOnly
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 text-gray-600 cursor-not-allowed"
                                        placeholder="Sẽ tự động tính khi chọn phòng"
                                    />
                                    <div className="text-xs text-gray-500 mt-1">
                                        {formatCurrency(calculateTotalAmount())}
                                    </div>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-2">
                                        Đã thanh toán
                                    </label>
                                    <input
                                        type="number"
                                        min="0"
                                        max={calculateTotalAmount()}
                                        value={newBooking.paidAmount}
                                        onChange={(e) => setNewBooking(prev => ({ ...prev, paidAmount: Number(e.target.value) }))}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                        placeholder="Nhập số tiền đã thanh toán"
                                        disabled={true}
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-2">
                                        Phương thức thanh toán
                                    </label>
                                    <select
                                        value={newBooking.paymentMethod}
                                        onChange={(e) => setNewBooking(prev => ({ ...prev, paymentMethod: e.target.value }))}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                        disabled={true}
                                    >
                                        <option value="Tiền mặt">Tiền mặt</option>
                                        <option value="Chuyển khoản">Chuyển khoản</option>
                                        <option value="Thẻ tín dụng">Thẻ tín dụng</option>
                                    </select>
                                </div>
                            </div>

                            {calculateTotalAmount() > 0 && newBooking.paidAmount < calculateTotalAmount() && (
                                <div className="mt-3 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                                    <p className="text-sm text-yellow-800">
                                        <strong>Còn lại:</strong> {formatCurrency(calculateTotalAmount() - newBooking.paidAmount)}
                                    </p>
                                </div>
                            )}
                        </div>

                        {/* Form Actions */}
                        <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
                            <button
                                type="button"
                                onClick={handleClose}
                                className="px-6 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors font-medium"
                            >
                                Đóng
                            </button>
                            {/* <button
                                type="submit"
                                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium flex items-center gap-2"
                            >
                                <Save className="w-4 h-4" />
                                Tạo đặt phòng
                            </button> */}
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default BookingDetailModal;