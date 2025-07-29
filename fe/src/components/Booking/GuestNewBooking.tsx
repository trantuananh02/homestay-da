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
import { Booking, Room } from "../../types";
import { useNavigate, useParams } from "react-router-dom";
import { homestayService } from "../../services/homestayService";
import { bookingService } from "../../services/bookingService";

const GuestNewBooking = () => {
  const { id } = useParams<{ id: string }>();
  const [rooms, setRooms] = useState<Room[]>([]); // Replace with actual rooms data
  const [existingBookings, setExistingBookings] = useState<Booking[]>([]); // Replace with actual bookings data
  const navigate = useNavigate();

  // get user info from localStorage
  const userInfo = localStorage.getItem("user");

  useEffect(() => {
    const fetchData = async () => {
      const roomList = await homestayService.getPublicRoomList({
        homestayId: Number(id),
        page: 1,
        pageSize: 100, // Lấy tất cả phòng
      });
      setRooms(roomList.rooms || []);

      const bookingList = await homestayService.getGuestBookingsByHomestayId(
        Number(id)
      );

      setExistingBookings(bookingList.bookings || []);
    };
    fetchData();
  }, [id]);

  const formatDate = (dateString: string) => {
    if (!dateString) return "";
    const date = new Date(dateString);
    const day = date.getDate().toString().padStart(2, "0");
    const month = (date.getMonth() + 1).toString().padStart(2, "0");
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
  };

  const [showRoomSuggestions, setShowRoomSuggestions] = useState(false);
  const [roomSearchQuery, setRoomSearchQuery] = useState("");
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

  const isRoomAvailable = (room: Room) => {
    // Kiểm tra phòng có sẵn cơ bản
    if (room?.status !== "available") return false;

    // Nếu chưa chọn ngày thì hiển thị tất cả phòng có sẵn
    if (!newBooking.checkIn || !newBooking.checkOut) return true;

    // Kiểm tra xung đột với các booking hiện tại
    const checkIn = new Date(newBooking.checkIn);
    const checkOut = new Date(newBooking.checkOut);

    return !existingBookings.some((booking) => {
      // Bỏ qua booking đã hủy
      if (booking.status === "cancelled") return false;

      // Kiểm tra xem phòng có trong booking này không
      const hasRoom = booking.rooms.some(
        (bookingRoom) => bookingRoom.id === room.id
      );
      if (!hasRoom) return false;

      // Kiểm tra xung đột thời gian
      const bookingCheckIn = new Date(booking.checkIn);
      const bookingCheckOut = new Date(booking.checkOut);

      // Xung đột nếu:
      // - Ngày nhận mới nằm trong khoảng booking cũ
      // - Ngày trả mới nằm trong khoảng booking cũ
      // - Booking mới bao trùm booking cũ
      return (
        (checkIn >= bookingCheckIn && checkIn < bookingCheckOut) ||
        (checkOut > bookingCheckIn && checkOut <= bookingCheckOut) ||
        (checkIn <= bookingCheckIn && checkOut >= bookingCheckOut)
      );
    });
  };

  const filteredRooms = rooms.filter(
    (room) =>
      room.name.toLowerCase().includes(roomSearchQuery.toLowerCase()) &&
      isRoomAvailable(room) &&
      !selectedRooms.some((selected) => selected.id === room.id)
  );

  const handleAddRoom = (room: Room) => {
    setSelectedRooms((prev) => [
      ...prev,
      {
        id: room.id,
        name: room.name,
        type: room.type,
        pricePerNight: room.price,
        capacity: room.capacity,
      },
    ]);
    setRoomSearchQuery("");
    setShowRoomSuggestions(false);
  };

  const handleRemoveRoom = (roomId: number) => {
    setSelectedRooms((prev) => prev.filter((room) => room.id !== roomId));
  };

  const handleRoomSearchChange = (value: string) => {
    setRoomSearchQuery(value);
    setShowRoomSuggestions(value.length > 0);
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
            </div>

            {/* Room Information */}
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="flex items-center gap-2 mb-4">
                <Home className="w-5 h-5 text-green-600" />
                <h3 className="text-lg font-semibold text-gray-900">
                  Chọn phòng
                </h3>
              </div>

              {/* Room Search */}
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Tìm kiếm phòng{" "}
                  {newBooking.checkIn && newBooking.checkOut && (
                    <span className="text-sm font-normal text-green-600">
                      (có sẵn từ {formatDate(newBooking.checkIn)} đến{" "}
                      {formatDate(newBooking.checkOut)})
                    </span>
                  )}
                </label>
                <div className="relative">
                  <input
                    type="text"
                    value={roomSearchQuery}
                    onChange={(e) => handleRoomSearchChange(e.target.value)}
                    onFocus={() =>
                      setShowRoomSuggestions(roomSearchQuery.length > 0)
                    }
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder={
                      !newBooking.checkIn || !newBooking.checkOut
                        ? "Vui lòng chọn ngày nhận và trả phòng trước..."
                        : "Nhập tên phòng để tìm kiếm và thêm..."
                    }
                    disabled={!newBooking.checkIn || !newBooking.checkOut}
                    autoComplete="off"
                  />

                  {/* Room Suggestions Dropdown */}
                  {showRoomSuggestions && filteredRooms.length > 0 && (
                    <div className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-300 rounded-lg shadow-lg z-50 max-h-60 overflow-y-auto">
                      {filteredRooms.map((room) => (
                        <div
                          key={room.id}
                          onClick={() => handleAddRoom(room)}
                          className="px-4 py-3 hover:bg-gray-50 cursor-pointer border-b border-gray-100 last:border-b-0"
                        >
                          <div className="flex items-center justify-between">
                            <div className="flex-1">
                              <div className="font-medium text-gray-900">
                                {room.name}
                              </div>
                              <div className="text-sm text-gray-500 flex items-center gap-2 mt-1">
                                <span
                                  className={`px-2 py-1 rounded text-xs font-medium ${
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
                                <span>•</span>
                                <span>{room.capacity} người</span>
                                <span>•</span>
                                <span className="font-medium text-green-600">
                                  {formatCurrency(room.price)}/đêm
                                </span>
                                {calculateNights() > 0 && (
                                  <>
                                    <span>•</span>
                                    <span className="font-medium text-blue-600">
                                      {formatCurrency(
                                        room.price * calculateNights()
                                      )}{" "}
                                      ({calculateNights()} đêm)
                                    </span>
                                  </>
                                )}
                              </div>
                              <div className="text-xs text-gray-400 mt-1">
                                {room?.amenities?.join(", ") || ""}
                              </div>
                            </div>
                            <div className="ml-2">
                              <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                                Thêm
                              </span>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}

                  {/* No results message */}
                  {showRoomSuggestions &&
                    roomSearchQuery.length > 0 &&
                    filteredRooms.length === 0 && (
                      <div className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-300 rounded-lg shadow-lg z-50 p-4 text-center text-gray-500">
                        {!newBooking.checkIn || !newBooking.checkOut
                          ? "Vui lòng chọn ngày nhận và trả phòng để xem phòng có sẵn"
                          : selectedRooms.length > 0
                          ? "Không tìm thấy phòng khác có sẵn trong thời gian này"
                          : "Không tìm thấy phòng có sẵn trong thời gian này"}
                      </div>
                    )}
                </div>
              </div>

              {/* Selected Rooms */}
              {selectedRooms.length > 0 && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Phòng đã chọn ({selectedRooms.length})
                  </label>
                  <div className="space-y-2 max-h-40 overflow-y-auto">
                    {selectedRooms.map((room) => (
                      <div
                        key={room.id}
                        className="flex items-center justify-between p-3 bg-white border border-gray-200 rounded-lg"
                      >
                        <div className="flex-1">
                          <div className="font-medium text-gray-900">
                            {room.name}
                          </div>
                          <div className="text-sm text-gray-500 flex items-center gap-2">
                            <span
                              className={`px-2 py-1 rounded text-xs font-medium ${
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
                            <span>•</span>
                            <span className="font-medium text-green-600">
                              {formatCurrency(room.pricePerNight)}/đêm
                            </span>
                            {calculateNights() > 0 && (
                              <>
                                <span>•</span>
                                <span className="font-medium text-blue-600">
                                  {formatCurrency(
                                    room.pricePerNight * calculateNights()
                                  )}{" "}
                                  ({calculateNights()} đêm)
                                </span>
                              </>
                            )}
                          </div>
                        </div>
                        <button
                          type="button"
                          onClick={() => handleRemoveRoom(room.id)}
                          className="ml-2 p-1 text-red-600 hover:text-red-800 hover:bg-red-50 rounded-full transition-colors"
                        >
                          <X className="w-4 h-4" />
                        </button>
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
                    <option value="Thẻ tín dụng">Thẻ tín dụng</option>
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
