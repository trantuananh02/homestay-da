import React from "react";
import { Search, MapPin, DollarSign, Users, Calendar } from "lucide-react";
import { HomestayListRequest } from "../../types";

interface SearchFiltersProps {
  filters: HomestayListRequest;
  onFiltersChange: (filters: HomestayListRequest) => void;
}

const SearchFilters: React.FC<SearchFiltersProps> = ({
  filters,
  onFiltersChange,
}) => {
  const handleFilterChange = (
    key: keyof HomestayListRequest,
    value: string | number
  ) => {
    onFiltersChange({
      ...filters,
      [key]: value,
      page: 1, // Reset to first page when filters change
    });
  };

  return (
    <div className="bg-white rounded-xl shadow-sm p-6 mb-8">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">
        Bộ lọc tìm kiếm
      </h3>

      <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-6 gap-4">
        <div className="relative">
          <Search className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
          <input
            type="text"
            placeholder="Tìm kiếm homestay..."
            className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
            value={filters.search || ""}
            onChange={(e) => handleFilterChange("search", e.target.value)}
          />
        </div>

        <div className="relative">
          <MapPin className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
          <input
            type="text"
            placeholder="Tỉnh/Thành phố"
            className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
            value={filters.city || ""}
            onChange={(e) => handleFilterChange("city", e.target.value)}
          />
        </div>

        <div className="relative">
          <MapPin className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
          <input
            type="text"
            placeholder="Quận/Huyện"
            className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
            value={filters.district || ""}
            onChange={(e) => handleFilterChange("district", e.target.value)}
          />
        </div>

        <div className="relative">
          <Calendar className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
          <input
            type="date"
            placeholder="Ngày nhận phòng"
            className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
            value={filters.checkIn || ""}
            onChange={(e) => handleFilterChange("checkIn", e.target.value)}
          />
        </div>

        <div className="relative">
          <Calendar className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
          <input
            type="date"
            placeholder="Ngày trả phòng"
            className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
            value={filters.checkOut || ""}
            onChange={(e) => handleFilterChange("checkOut", e.target.value)}
            min={filters.checkIn || undefined}
          />
        </div>

        <div className="relative">
          <Users className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
          <input
            type="number"
            placeholder="Số khách"
            min="1"
            className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
            value={filters.guests || ""}
            onChange={(e) =>
              handleFilterChange(
                "guests",
                e.target.value ? parseInt(e.target.value) : 0
              )
            }
          />
        </div>
      </div>

      <div className="mt-4 flex justify-between items-center">
        <div className="text-sm text-gray-600">
          Hiển thị {filters.pageSize || 12} kết quả mỗi trang
        </div>
        <button
          onClick={() =>
            onFiltersChange({
              page: 1,
              pageSize: 12,
              search: "",
              city: "",
              district: "",
              status: "active",
              checkIn: "",
              checkOut: "",
              guests: 0,
            })
          }
          className="bg-gray-100 text-gray-700 px-4 py-2 rounded-lg hover:bg-gray-200 transition-colors flex items-center space-x-2"
        >
          <Search className="h-4 w-4" />
          <span>Xóa bộ lọc</span>
        </button>
      </div>
    </div>
  );
};

export default SearchFilters;
