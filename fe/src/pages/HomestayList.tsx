import React, { useState, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import HomestayCard from "../components/Homestay/HomestayCard";
import SearchFilters from "../components/Homestay/SearchFilters";
import { homestayService } from "../services/homestayService";
import { Homestay, HomestayListRequest } from "../types";
import { SearchX } from "lucide-react";

const DEFAULT_PAGE_SIZE = 12;

const HomestayList: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [homestays, setHomestays] = useState<Homestay[]>([]);
  const [loading, setLoading] = useState(true);
  const [hasMore, setHasMore] = useState(true);
  const [filters, setFilters] = useState<HomestayListRequest>({
    page: 1,
    pageSize: DEFAULT_PAGE_SIZE,
    search: "",
    city: "",
    district: "",
    status: "active",
    checkIn: "",
    checkOut: "",
    guests: 0,
  });

  useEffect(() => {
    if (location.state?.filters) {
      setFilters((prev) => ({ ...prev, ...location.state.filters, page: 1 }));
    }
  }, [location.state]);

  const loadHomestays = async () => {
    try {
      setLoading(true);
      const response = await homestayService.getPublicHomestayList(filters);
      const newList =
        filters.page === 1
          ? response.homestays
          : [...homestays, ...response.homestays];

      setHomestays(newList);
      setHasMore(
        response.homestays.length === (filters.pageSize || DEFAULT_PAGE_SIZE)
      );
    } catch (error) {
      console.error("Error loading homestays:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadHomestays();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filters]);

  const handleHomestayClick = (id: number) => {
    navigate(`/homestay/${id}`);
  };

  const handleFiltersChange = (newFilters: HomestayListRequest) => {
    setFilters((prev) => ({ ...prev, ...newFilters, page: 1 }));
  };

  const handleLoadMore = () => {
    setFilters((prev) => ({ ...prev, page: (prev.page || 1) + 1 }));
  };

  return (
    <div className="min-h-screen bg-white font-sans">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        {/* Section Header */}
        <div className="text-center mb-14">
          <h1 className="text-4xl font-light text-indigo-700 tracking-wide mb-4">
            Tìm kiếm Homestay
          </h1>
          <p className="text-gray-500 text-lg">
            {homestays.length > 0
              ? `${homestays.length} nơi trú chân được tìm thấy`
              : "Đang tìm những homestay đẹp nhất cho bạn..."}
          </p>
        </div>

        {/* Bộ lọc */}
        <SearchFilters
          filters={filters}
          onFiltersChange={handleFiltersChange}
        />

        {/* Danh sách */}
        {loading && homestays.length === 0 ? (
          <div className="flex justify-center items-center py-24">
            <div className="text-center">
              <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-indigo-600 mx-auto mb-4"></div>
              <p className="text-gray-500">Đang tải homestay...</p>
            </div>
          </div>
        ) : homestays.length === 0 ? (
          <div className="text-center py-24 text-gray-500">
            <SearchX className="mx-auto mb-4 h-12 w-12 text-gray-400" />
            <h3 className="text-xl font-light mb-2">
              Không tìm thấy homestay phù hợp
            </h3>
            <p className="text-sm">
              Hãy thử điều chỉnh lại bộ lọc của bạn nhé 🌿
            </p>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-3 gap-8 mt-12">
              {homestays.map((homestay) => (
                <HomestayCard
                  key={homestay.id}
                  homestay={homestay}
                  onClick={() => handleHomestayClick(homestay.id)}
                />
              ))}
            </div>

            {hasMore && (
              <div className="text-center mt-16">
                <button
                  onClick={handleLoadMore}
                  disabled={loading}
                  className="px-6 py-3 bg-indigo-600 text-white rounded-full hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition"
                >
                  {loading ? "Đang tải thêm..." : "Xem thêm homestay"}
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default HomestayList;
