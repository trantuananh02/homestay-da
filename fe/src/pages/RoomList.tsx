import React from 'react';
import { Room } from '../types';

interface RoomListProps {
  rooms: Room[];
  onViewRoom?: (roomId: number) => void;
}

const priceTypeLabel = (type: Room['priceType']) =>
  type === 'per_night' ? 'phòng/đêm' : 'người/đêm';

// const statusLabel = (status: Room['status']) => {
//   switch (status) {
//     case 'available':
//       return <span className="text-primary-600 font-semibold">Còn phòng</span>;
//     case 'occupied':
//       return <span className="text-yellow-600 font-semibold">Đã đặt</span>;
//     case 'maintenance':
//       return <span className="text-gray-500 font-semibold">Bảo trì</span>;
//     default:
//       return status;
//   }
// };

const RoomList: React.FC<RoomListProps> = ({ rooms, onViewRoom}) => {
  if (!rooms || rooms.length === 0) {
    return <div className="text-center text-gray-500 py-8">Chưa có phòng nào.</div>;
  }

  return (
    <div className="space-y-4">
      {rooms.map((room) => (
        <div
          key={room.id}
          className="bg-white rounded-lg shadow flex flex-col md:flex-row items-stretch overflow-hidden"
        >
          {/* Cột 1: Ảnh */}
          <div className="md:w-1/4 w-full flex-shrink-0 flex items-center justify-center bg-gray-100">
            {room.images && room.images.length > 0 ? (
              <img
                src={room.images[0]}
                alt={room.name}
                className="object-cover w-full h-40 md:h-32"
              />
            ) : (
              <div className="w-full h-40 md:h-32 flex items-center justify-center text-gray-400">
                Không có ảnh
              </div>
            )}
          </div>

          {/* Cột 2: Thông tin chi tiết */}
          <div className="md:w-2/4 w-full p-4 flex flex-col justify-center">
            <h3 className="text-lg font-semibold mb-1">{room.name}</h3>
            <div className="flex flex-wrap gap-2 mb-2 text-xs text-gray-600">
              <span className="px-2 py-1 bg-gray-100 rounded">Loại: {room.type}</span>
              <span className="px-2 py-1 bg-gray-100 rounded">Sức chứa: {room.capacity} người</span>
              {room.area && <span className="px-2 py-1 bg-gray-100 rounded">Diện tích: {room.area} m²</span>}
              <span className="px-2 py-1 bg-gray-100 rounded">Giá: {room.price?.toLocaleString('vi-VN')} đ/{priceTypeLabel(room.priceType)}</span>
              {/* <span className="px-2 py-1 bg-gray-100 rounded">Trạng thái: {statusLabel(room.status)}</span> */}
            </div>
            {room.description && (
              <div className="text-gray-600 text-sm mb-2 line-clamp-2">{room.description}</div>
            )}
            {room.amenities && room.amenities.length > 0 && (
              <div className="flex flex-wrap gap-2 mt-1">
                {room.amenities.slice(0, 5).map((item, idx) => (
                  <span key={idx} className="px-2 py-1 bg-primary-50 text-primary-700 rounded-full text-xs">
                    {item}
                  </span>
                ))}
                {room.amenities.length > 5 && (
                  <span className="px-2 py-1 bg-gray-100 text-gray-500 rounded-full text-xs">+{room.amenities.length - 5} tiện ích</span>
                )}
              </div>
            )}
          </div>

          {/* Cột 3: Nút đặt ngay */}
          <div className="md:w-1/4 w-full flex flex-col items-center justify-center p-4 border-t md:border-t-0 md:border-l border-gray-100">
            <button
              className="w-full px-4 py-2 bg-primary-600 text-white rounded-lg font-semibold hover:bg-primary-700 transition mb-2"
              onClick={() => onViewRoom && onViewRoom(room.id)}
            >
              Xem chi tiết
            </button>
          </div>
        </div>
      ))}
    </div>
  );
};

export default RoomList;
