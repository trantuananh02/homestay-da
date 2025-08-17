import React, { useState } from 'react';
import { Star, X, MapPin, Trash2 } from 'lucide-react';
import { Review, Booking } from '../../types';
import { useAuth } from '../../contexts/AuthContext';
import CusFormUpload from '../UploadFile';
import { homestayService } from '../../services/homestayService';

interface ReviewModalProps {
  isOpen: boolean;
  onClose: () => void;
  booking: Booking;
  onSubmit: (review: Review) => void;
}

const ReviewModal: React.FC<ReviewModalProps> = ({
  isOpen,
  onClose,
  booking,
  onSubmit
}) => {
  const { user } = useAuth();
  const [rating, setRating] = useState(5);
  const [comment, setComment] = useState('');
  const [hoveredRating, setHoveredRating] = useState(0);
  const [imageUrls, setImageUrls] = useState<string[]>([]);
  const [isUploading, setIsUploading] = useState(false);

  const handleImageUpload = async (
    e: React.ChangeEvent<HTMLInputElement>
  ): Promise<void> => {
    const files: File[] = Array.from(e.target.files ?? []);
    setIsUploading(true);

    const uploaded = await Promise.all(
      files.map(async (file) => {
        try {
          const url = await homestayService.uploadRoomImage(file);
          return url;
        } catch {
          return null;
        }
      })
    );

    setImageUrls(prev => [
      ...prev,
      ...(uploaded.filter(Boolean) as string[])
    ]);
    setIsUploading(false);
  };

  const removeImage = (index: number) => {
    setImageUrls(prev => prev.filter((_, i) => i !== index));
  };

    const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!comment.trim()) {
      alert('Vui lòng nhập nhận xét!');
      return;
    }

    const review: Review = {
      homestayId: 0, // This should be set to the actual homestay ID
      bookingId: booking.id,
      guestId: user?.id || 0,
      rating,
      comment,
      createdAt: new Date().toISOString(),
      imageUrls: imageUrls,
    };

    onSubmit(review);
    onClose();

    // Reset form
    setRating(5);
    setComment('');
    setImageUrls([]);
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-gray-900">Đánh giá homestay</h2>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600"
            >
              <X className="h-6 w-6" />
            </button>
          </div>

          <div className="mb-6">
            <div className="flex items-center space-x-4 p-4 bg-gray-50 rounded-lg">
              <div>
                <div className="space-y-1">
                  {booking.rooms.map((room, index) => (
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
                <p className="text-sm text-gray-500 mt-2">
                 Từ {booking.checkIn} đến {booking.checkOut}
                </p>
              </div>
            </div>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Rating */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-3">
                Đánh giá tổng thể
              </label>
              <div className="flex items-center space-x-2">
                {[1, 2, 3, 4, 5].map((star) => (
                  <button
                    key={star}
                    type="button"
                    onClick={() => setRating(star)}
                    onMouseEnter={() => setHoveredRating(star)}
                    onMouseLeave={() => setHoveredRating(0)}
                    className="focus:outline-none"
                  >
                    <Star
                      className={`h-8 w-8 ${star <= (hoveredRating || rating)
                          ? 'text-yellow-400 fill-current'
                          : 'text-gray-300'
                        }`}
                    />
                  </button>
                ))}
                <span className="ml-2 text-sm text-gray-600">
                  {rating === 1 && 'Rất tệ'}
                  {rating === 2 && 'Tệ'}
                  {rating === 3 && 'Bình thường'}
                  {rating === 4 && 'Tốt'}
                  {rating === 5 && 'Tuyệt vời'}
                </span>
              </div>
            </div>

            {/* Comment */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Nhận xét chi tiết *
              </label>
              <textarea
                value={comment}
                onChange={(e) => setComment(e.target.value)}
                rows={4}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                placeholder="Chia sẻ trải nghiệm của bạn về homestay này..."
                required
              />
            </div>

            {/* Image Upload */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Thêm ảnh (tùy chọn)
              </label>
              
              <div className="flex flex-wrap gap-4">
                {imageUrls.map((imageUrl, index) => (
                  <div key={index} className="relative">
                    <img
                      src={imageUrl}
                      alt={`Ảnh review ${index + 1}`}
                      className="w-40 h-40 object-cover rounded-lg"
                    />
                    <button
                      type="button"
                      onClick={() => removeImage(index)}
                      className="absolute top-2 right-2 text-red-500 hover:text-red-700"
                    >
                      <Trash2 className="w-5 h-5" />
                    </button>
                  </div>
                ))}
                <CusFormUpload
                  disabled={false}
                  handleUpload={handleImageUpload}
                  isUploading={isUploading}
                />
              </div>
            </div>

            {/* Submit Buttons */}
            <div className="flex space-x-4 pt-6 border-t">
              <button
                type="button"
                onClick={onClose}
                className="flex-1 bg-gray-300 text-gray-700 py-3 rounded-lg font-semibold hover:bg-gray-400 transition-colors"
                disabled={isUploading}
              >
                Hủy
              </button>
              <button
                type="submit"
                className="flex-1 bg-emerald-600 text-white py-3 rounded-lg font-semibold hover:bg-emerald-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                disabled={isUploading}
              >
                {isUploading ? 'Đang xử lý...' : 'Gửi đánh giá'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default ReviewModal;