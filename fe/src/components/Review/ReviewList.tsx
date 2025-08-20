import React from "react";
import { Star, Calendar, User } from "lucide-react";
import { Review } from "../../types";
import { parseImageUrls } from "../../services/homestayService";

interface ReviewListProps {
  reviews: Review[];
  showUserInfo?: boolean;
}

const ReviewList: React.FC<ReviewListProps> = ({
  reviews,
  showUserInfo = true,
}) => {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("vi-VN", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const getRatingText = (rating: number) => {
    switch (rating) {
      case 1:
        return "Rất tệ";
      case 2:
        return "Tệ";
      case 3:
        return "Bình thường";
      case 4:
        return "Tốt";
      case 5:
        return "Tuyệt vời";
      default:
        return "";
    }
  };

  return (
    <div className="space-y-6">
      {reviews.map((review) => (
        <div
          key={review.id}
          className="bg-white rounded-lg border border-gray-200 p-6 shadow-sm"
        >
          {/* Header */}
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center space-x-3">
              {showUserInfo && (
                <div className="w-10 h-10 bg-emerald-100 rounded-full flex items-center justify-center">
                  <User className="w-5 h-5 text-emerald-600" />
                </div>
              )}
              <div>
                {showUserInfo && (
                  <h4 className="font-semibold text-gray-900">
                    {review.userName || "Khách hàng"}
                  </h4>
                )}
                <div className="flex items-center space-x-2 text-sm text-gray-500">
                  <Calendar className="w-4 h-4" />
                  <span>{formatDate(review.createdAt)}</span>
                </div>
              </div>
            </div>

            {/* Rating */}
            <div className="flex items-center space-x-2">
              <div className="flex items-center space-x-1">
                {[1, 2, 3, 4, 5].map((star) => (
                  <Star
                    key={star}
                    className={`w-4 h-4 ${
                      star <= review.rating
                        ? "text-yellow-400 fill-current"
                        : "text-gray-300"
                    }`}
                  />
                ))}
              </div>
              <span className="text-sm font-medium text-gray-900">
                {review.rating}/5
              </span>
              <span className="text-xs text-gray-500">
                {getRatingText(review.rating)}
              </span>
            </div>
          </div>

          {/* Comment */}
          {review.comment && (
            <div className="mb-4">
              <p className="text-gray-700 leading-relaxed">{review.comment}</p>
            </div>
          )}

          {/* Images */}
          {parseImageUrls(review.imageUrls).length > 0 && (
            <div className="mb-4">
              <h5 className="text-sm font-medium text-gray-700 mb-3">
                Ảnh đánh giá:
              </h5>
              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
                {parseImageUrls(review.imageUrls).map((imageUrl, index) => (
                  <div key={index} className="relative group">
                    <img
                      src={imageUrl}
                      alt={`Review image ${index + 1}`}
                      className="w-full h-24 object-cover rounded-lg border border-gray-200 hover:border-emerald-300 transition-colors"
                    />
                    {/* Image preview on hover */}
                    <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all duration-200 rounded-lg flex items-center justify-center">
                      <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                        <button
                          onClick={() => window.open(imageUrl, "_blank")}
                          className="bg-white bg-opacity-90 text-gray-800 px-3 py-1 rounded-full text-sm font-medium hover:bg-opacity-100 transition-all"
                        >
                          Xem ảnh
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Homestay info if available */}
          {review.homestayName && (
            <div className="pt-4 border-t border-gray-100">
              <p className="text-sm text-gray-500">
                Đánh giá cho:{" "}
                <span className="font-medium text-gray-700">
                  {review.homestayName}
                </span>
              </p>
            </div>
          )}
        </div>
      ))}

      {reviews.length === 0 && (
        <div className="text-center py-12">
          <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <Star className="w-8 h-8 text-gray-400" />
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            Chưa có đánh giá
          </h3>
          <p className="text-gray-500">
            Hãy là người đầu tiên đánh giá homestay này!
          </p>
        </div>
      )}
    </div>
  );
};

export default ReviewList;
