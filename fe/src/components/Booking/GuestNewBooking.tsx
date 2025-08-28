import React, { useEffect, useState } from "react";
import {
  X,
  User,
  Home,
  Calendar,
  DollarSign,
  Save,
  ArrowLeft,
} from "lucide-react";
import { Booking } from "../../types";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import { homestayService } from "../../services/homestayService";
import { bookingService } from "../../services/bookingService";

const GuestNewBooking = () => {
  const { id } = useParams<{ id: string }>();
  const location = useLocation();
  const navigate = useNavigate();
  const [existingBookings, setExistingBookings] = useState<Booking[]>([]);
  const [homestayName, setHomestayName] = useState<string>("");

  // get user info from localStorage
  const userInfo = localStorage.getItem("user");

  useEffect(() => {
    const fetchData = async () => {
      // Lấy tên homestay
      try {
        const homestay = await homestayService.getPublicHomestayDetail(
          Number(id)
        );
        setHomestayName(homestay.name || "");
      } catch {
        /* empty */
      }

      const roomList = await homestayService.getPublicRoomList({
        homestayId: Number(id),
        page: 1,
        pageSize: 100, // Lấy tất cả phòng
      });
      const fetchedRooms = roomList.rooms || [];

      // Fetch existing bookings to check availability
      try {
        const bookingList = await homestayService.getGuestBookingsByHomestayId(
          Number(id)
        );
        setExistingBookings(bookingList.bookings || []);
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
      } catch (err) {
        // ignore; toast already handled in service
      }

      // Auto-select room if roomId is provided in query string
      const searchParams = new URLSearchParams(location.search);
      const selectedRoomIdParam = searchParams.get("roomId");
      if (selectedRoomIdParam) {
        const targetRoom = fetchedRooms.find(
          (r) => r.id === Number(selectedRoomIdParam)
        );
        if (targetRoom) {
          setSelectedRooms([
            {
              id: targetRoom.id,
              name: targetRoom.name,
              type: targetRoom.type,
              pricePerNight: targetRoom.price,
              capacity: targetRoom.capacity,
            },
          ]);
        }
      }
    };
    fetchData();
  }, [id, location.search]);

  const [selectedRooms, setSelectedRooms] = useState<
    {
      id: number;
      name: string;
      type: string;
      pricePerNight: number;
      capacity: number;
    }[]
  >([]);
  const [newBooking, setNewBooking] = useState({
    customerName: userInfo ? JSON.parse(userInfo).name : "",
    customerPhone: userInfo ? JSON.parse(userInfo).phone || "" : "",
    customerEmail: userInfo ? JSON.parse(userInfo).email : "",
    checkIn: "",
    checkOut: "",
    paidAmount: 0,
    paymentMethod: "Tiền mặt",
  });

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
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
    return selectedRooms.reduce(
      (total, room) => total + room.pricePerNight * nights,
      0
    );
  };

  // Check availability for selected rooms against existing bookings
  const getUnavailableRooms = (): number[] => {
    if (!newBooking.checkIn || !newBooking.checkOut) return [];
    const checkIn = new Date(newBooking.checkIn);
    const checkOut = new Date(newBooking.checkOut);

    const isOverlap = (aStart: Date, aEnd: Date, bStart: Date, bEnd: Date) => {
      return (
        (aStart >= bStart && aStart < bEnd) ||
        (aEnd > bStart && aEnd <= bEnd) ||
        (aStart <= bStart && aEnd >= bEnd)
      );
    };

    const selectedIds = selectedRooms.map((r) => r.id);

    const unavailable = new Set<number>();
    for (const booking of existingBookings) {
      if (booking.status === "cancelled") continue;
      const bStart = new Date(booking.checkIn);
      const bEnd = new Date(booking.checkOut);
      if (!isOverlap(checkIn, checkOut, bStart, bEnd)) continue;
      for (const r of booking.rooms) {
        if (selectedIds.includes(r.id)) {
          unavailable.add(r.id);
        }
      }
    }
    return Array.from(unavailable);
  };

  const handleRemoveRoom = (roomId: number) => {
    setSelectedRooms((prev) => prev.filter((room) => room.id !== roomId));
  };

  const handleDateChange = (field: "checkIn" | "checkOut", value: string) => {
    setNewBooking((prev) => ({ ...prev, [field]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (selectedRooms.length === 0) {
      alert("Vui lòng chọn ít nhất một phòng");
      return;
    }

    const nights = calculateNights();
    if (nights <= 0) {
      alert("Ngày trả phòng phải sau ngày nhận phòng");
      return;
    }

    // Validate availability against existing bookings
    const unavailable = getUnavailableRooms();
    if (unavailable.length > 0) {
      alert(
        "Một hoặc nhiều phòng đã hết trong khoảng thời gian bạn chọn. Vui lòng đổi ngày khác hoặc chọn phòng khác."
      );
      return;
    }

    const totalAmount = calculateTotalAmount();
    const bookingRooms = selectedRooms.map((room) => ({
      id: room.id,
      name: room.name,
      type: room.type,
      pricePerNight: room.pricePerNight,
      nights: nights,
      subtotal: room.pricePerNight * nights,
      capacity: room.capacity,
    }));

    await bookingService.createGuestBooking({
      homestayId: Number(id),
      customerName: newBooking.customerName,
      customerPhone: newBooking.customerPhone,
      customerEmail: newBooking.customerEmail,
      rooms: bookingRooms,
      checkIn: newBooking.checkIn,
      checkOut: newBooking.checkOut,
      nights: nights,
      totalAmount: totalAmount,
      paidAmount: newBooking.paidAmount,
      paymentMethod: newBooking.paymentMethod,
    });

    navigate(`/bookings`);
  };

  const handleBack = () => {
    // Logic to navigate back, e.g., using useNavigate from react-router
    window.history.back();
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <button
          onClick={handleBack}
          className="flex items-center space-x-2 text-emerald-600 hover:text-emerald-700 mb-4"
        >
          <ArrowLeft className="h-5 w-5" />
          <span>Quay lại</span>
        </button>

        {/* Modal Header */}
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold text-gray-900">
            Tạo đặt phòng mới
          </h2>
        </div>

        {/* Modal Content */}
        <div className="p-6">
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Customer Information */}
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="flex items-center gap-2 mb-4">
                <User className="w-5 h-5 text-blue-600" />
                <h3 className="text-lg font-semibold text-gray-900">
                  Thông tin khách hàng
                </h3>
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
                    onChange={(e) =>
                      setNewBooking((prev) => ({
                        ...prev,
                        customerName: e.target.value,
                      }))
                    }
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Nhập tên khách hàng"
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
                    onChange={(e) =>
                      setNewBooking((prev) => ({
                        ...prev,
                        customerPhone: e.target.value,
                      }))
                    }
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Nhập số điện thoại"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Email
                  </label>
                  <input
                    type="email"
                    value={newBooking.customerEmail}
                    onChange={(e) =>
                      setNewBooking((prev) => ({
                        ...prev,
                        customerEmail: e.target.value,
                      }))
                    }
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Nhập email"
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
                    onChange={(e) =>
                      handleDateChange("checkIn", e.target.value)
                    }
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
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
                    onChange={(e) =>
                      handleDateChange("checkOut", e.target.value)
                    }
                    min={newBooking.checkIn}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
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

              {(() => {
                const unavailable = getUnavailableRooms();
                if (unavailable.length > 0) {
                  return (
                    <div className="mt-3 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700">
                      Một hoặc nhiều phòng bạn chọn đã được đặt trong khoảng
                      thời gian này. Vui lòng đổi ngày khác hoặc chọn phòng
                      khác.
                    </div>
                  );
                }
                return null;
              })()}
            </div>

            {/* Homestay Name + Room Information */}
            <div className="bg-gray-50 rounded-lg p-4">
              {homestayName && (
                <div className="flex items-center gap-2 mb-2">
                  <span className="text-base font-semibold text-purple-700">
                    Homestay:
                  </span>
                  <span className="text-base font-bold text-gray-900">
                    {homestayName}
                  </span>
                </div>
              )}
              <div className="flex items-center gap-2 mb-4">
                <Home className="w-5 h-5 text-green-600" />
                <h3 className="text-lg font-semibold text-gray-900">
                  Phòng đã chọn
                </h3>
              </div>

              {/* Selected Rooms Display */}
              {selectedRooms.length > 0 ? (
                <div className="space-y-4">
                  {selectedRooms.map((room) => (
                    <div
                      key={room.id}
                      className="bg-white border border-gray-200 rounded-lg p-4"
                    >
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="flex items-center gap-3 mb-3">
                            <h4 className="text-lg font-semibold text-gray-900">
                              {room.name}
                            </h4>
                            <span
                              className={`px-3 py-1 rounded-full text-sm font-medium ${
                                room.type === "Standard"
                                  ? "bg-gray-100 text-gray-800"
                                  : room.type === "Deluxe"
                                  ? "bg-blue-100 text-blue-800"
                                  : room.type === "Premium"
                                  ? "bg-purple-100 text-purple-800"
                                  : "bg-yellow-100 text-yellow-800"
                              }`}
                            >
                              {room.type}
                            </span>
                          </div>

                          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-3">
                            <div className="space-y-2">
                              <div className="flex items-center gap-2">
                                <span className="text-sm text-gray-600">
                                  Sức chứa:
                                </span>
                                <span className="font-medium text-gray-900">
                                  {room.capacity} người
                                </span>
                              </div>
                              <div className="flex items-center gap-2">
                                <span className="text-sm text-gray-600">
                                  Giá/đêm:
                                </span>
                                <span className="font-medium text-emerald-600">
                                  {formatCurrency(room.pricePerNight)}
                                </span>
                              </div>
                            </div>

                            {calculateNights() > 0 && (
                              <div className="space-y-2">
                                <div className="flex items-center gap-2">
                                  <span className="text-sm text-gray-600">
                                    Số đêm:
                                  </span>
                                  <span className="font-medium text-gray-900">
                                    {calculateNights()}
                                  </span>
                                </div>
                                <div className="flex items-center gap-2">
                                  <span className="text-sm text-gray-600">
                                    Tổng tiền:
                                  </span>
                                  <span className="font-medium text-blue-600">
                                    {formatCurrency(
                                      room.pricePerNight * calculateNights()
                                    )}
                                  </span>
                                </div>
                              </div>
                            )}
                          </div>
                        </div>

                        <button
                          type="button"
                          onClick={() => handleRemoveRoom(room.id)}
                          className="ml-4 p-2 text-red-600 hover:text-red-800 hover:bg-red-50 rounded-full transition-colors"
                        >
                          <X className="w-5 h-5" />
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8 text-gray-500">
                  <Home className="w-12 h-12 mx-auto mb-3 text-gray-300" />
                  <p>Chưa có phòng nào được chọn</p>
                  <p className="text-sm">
                    Phòng sẽ được tự động chọn nếu bạn bấm "Đặt phòng" từ trang
                    chi tiết
                  </p>
                </div>
              )}
            </div>

            {/* Payment Information */}
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="flex items-center gap-2 mb-4">
                <DollarSign className="w-5 h-5 text-orange-600" />
                <h3 className="text-lg font-semibold text-gray-900">
                  Thông tin thanh toán
                </h3>
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
                    onChange={(e) =>
                      setNewBooking((prev) => ({
                        ...prev,
                        paidAmount: Number(e.target.value),
                      }))
                    }
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Nhập số tiền đã thanh toán"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Phương thức thanh toán
                  </label>
                  <select
                    value={newBooking.paymentMethod}
                    onChange={(e) =>
                      setNewBooking((prev) => ({
                        ...prev,
                        paymentMethod: e.target.value,
                      }))
                    }
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="Tiền mặt">Tiền mặt</option>
                    <option value="Chuyển khoản">Chuyển khoản</option>
                  </select>
                </div>
              </div>

              {calculateTotalAmount() > 0 &&
                newBooking.paidAmount < calculateTotalAmount() && (
                  <div className="mt-3 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                    <p className="text-sm text-yellow-800">
                      <strong>Còn lại:</strong>{" "}
                      {formatCurrency(
                        calculateTotalAmount() - newBooking.paidAmount
                      )}
                    </p>
                  </div>
                )}
            </div>

            {/* Form Actions */}
            <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
              <button
                type="submit"
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium flex items-center gap-2"
              >
                <Save className="w-4 h-4" />
                Tạo đặt phòng
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default GuestNewBooking;
