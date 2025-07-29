import React from 'react';
import { Star, MapPin, Users, Bed, Bath, Square } from 'lucide-react';
import { Homestay } from '../../types';

interface HomestayCardProps {
  homestay: Homestay;
  onClick: () => void;
}

const HomestayCard: React.FC<HomestayCardProps> = ({ homestay, onClick }) => {
  const firstRoom = homestay.rooms?.[0] || null;

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND',
    }).format(price);
  };

  return (
    <div
      className="bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-xl transition-all duration-300 cursor-pointer transform hover:-translate-y-1"
      onClick={onClick}
    >
      <div className="relative h-48 overflow-hidden">
        <img
          src={firstRoom?.images?.[0] ?? ""}
          alt={homestay.name}
          className="w-full h-full object-cover transition-transform duration-300 hover:scale-110"
        />
      </div>

      <div className="p-6">
        <div className="flex items-center justify-between mb-2">
          <h3 className="text-xl font-semibold text-gray-900 truncate">{homestay.name}</h3>
          <div className="flex items-center space-x-1">
            <Star className="h-4 w-4 text-yellow-400 fill-current" />
            <span className="text-sm font-medium">{homestay.rating ?? 0}</span>
            <span className="text-sm text-gray-500">({homestay.totalReviews ?? 0})</span>
          </div>
        </div>

        <div className="flex items-center text-gray-600 mb-3">
          <MapPin className="h-4 w-4 mr-1" />
          <span className="text-sm">{`${homestay.address}, ${homestay.ward}, ${homestay.district}, ${homestay.city}`}</span>
        </div>

        <div className="flex items-center justify-between text-sm text-gray-600 mb-4">
          <div className="flex items-center space-x-4">
            <div className="flex items-center">
              <Users className="h-4 w-4 mr-1" />
              <span>{firstRoom?.capacity ?? '-'}</span>
            </div>
            <div className="flex items-center">
              <Bed className="h-4 w-4 mr-1" />
              <span>{firstRoom ? 1 : '-'}</span>
            </div>
            <div className="flex items-center">
              <Bath className="h-4 w-4 mr-1" />
              <span>-</span> {/* nếu có trường bathrooms, thay thế tại đây */}
            </div>
            <div className="flex items-center">
              <Square className="h-4 w-4 mr-1" />
              <span>{firstRoom?.area ? `${firstRoom.area} m²` : '- m²'}</span>
            </div>
          </div>
        </div>

        <p className="text-gray-600 text-sm mb-4 line-clamp-2">
          {homestay.description}
        </p>

        <div className="flex items-center justify-between">
          <div className="text-right">
            {(() => {
              const rooms = homestay.rooms || [];
              if (rooms.length === 0) {
                return (
                  <>
                    <p className="text-2xl font-bold text-emerald-600">Đang cập nhật</p>
                    <p className="text-sm text-gray-500">/đêm</p>
                  </>
                );
              }

              if (rooms.length === 1) {
                return (
                  <>
                    <p className="text-2xl font-bold text-emerald-600">
                      {formatPrice(rooms[0].price)}
                    </p>
                    <p className="text-sm text-gray-500">/đêm</p>
                  </>
                );
              }

              const prices = rooms.map(r => r.price);
              const min = Math.min(...prices);
              const max = Math.max(...prices);

              return (
                <>
                  <p className="text-2xl font-bold text-emerald-600">
                    {`${formatPrice(min)} - ${formatPrice(max)}`}
                  </p>
                  <p className="text-sm text-gray-500">/đêm</p>
                </>
              );
            })()}

          </div>
        </div>
      </div>
    </div>
  );
};

export default HomestayCard;
