import React, { useState } from "react";
import { Search, MapPin, Calendar, Users } from "lucide-react";
import home from "../../asset/home.webp";
import { HomestayListRequest } from "../../types";

interface HeroProps {
  onSearch: (filters: HomestayListRequest) => void;
}

const Hero: React.FC<HeroProps> = ({ onSearch }) => {
  const [searchData, setSearchData] = useState({
    location: "",
    checkIn: "",
    checkOut: "",
    guests: 1,
  });

  const handleSearch = () => {
    console.log("Hero search data:", searchData);

    // Kiểm tra xem có đủ thông tin không
    if (!searchData.checkIn || !searchData.checkOut) {
      alert("Vui lòng chọn ngày nhận phòng và ngày trả phòng");
      return;
    }

    if (searchData.checkOut <= searchData.checkIn) {
      alert("Ngày trả phòng phải sau ngày nhận phòng");
      return;
    }

    onSearch({
      city: searchData.location,
      checkIn: searchData.checkIn,
      checkOut: searchData.checkOut,
      guests: searchData.guests,
    });
  };

  return (
    <div className="relative bg-gradient-to-br from-primary-600 to-primary-800 text-white">
      <div className="absolute inset-0 bg-black opacity-20"></div>
      <div
        className="relative min-h-[600px] bg-cover bg-center flex items-center"
        style={{
          backgroundImage: `url(${home})`,
        }}
      >
        <div className="absolute inset-0 bg-primary-900 opacity-60"></div>
        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h1 className="text-4xl md:text-6xl font-bold mb-6">
            Khám Phá Homestay
            <span className="block text-orange-300">Tuyệt Vời Nhất</span>
          </h1>
          <p className="text-xl mb-12 max-w-2xl mx-auto opacity-90">
            Trải nghiệm không gian sống ấm cúng như nhà, với dịch vụ chuyên
            nghiệp và vị trí đắc địa khắp Việt Nam
          </p>

          {/* Search Form */}
          <div className="bg-white rounded-xl shadow-2xl p-6 max-w-4xl mx-auto">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div className="relative">
                <MapPin className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
                <input
                  type="text"
                  placeholder="Bạn muốn đi đâu?"
                  className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg text-gray-900 focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                  value={searchData.location}
                  onChange={(e) =>
                    setSearchData({ ...searchData, location: e.target.value })
                  }
                />
              </div>

              <div className="relative">
                <Calendar className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
                <input
                  type="date"
                  className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg text-gray-900 focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                  value={searchData.checkIn}
                  onChange={(e) =>
                    setSearchData({ ...searchData, checkIn: e.target.value })
                  }
                />
              </div>

              <div className="relative">
                <Calendar className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
                <input
                  type="date"
                  className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg text-gray-900 focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                  value={searchData.checkOut}
                  onChange={(e) =>
                    setSearchData({ ...searchData, checkOut: e.target.value })
                  }
                  min={searchData.checkIn || undefined}
                />
              </div>

              <div className="relative">
                <Users className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
                <select
                  className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg text-gray-900 focus:ring-2 focus:ring-primary-500 focus:border-transparent appearance-none"
                  value={searchData.guests}
                  onChange={(e) =>
                    setSearchData({
                      ...searchData,
                      guests: parseInt(e.target.value),
                    })
                  }
                >
                  {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((num) => (
                    <option key={num} value={num}>
                      {num} khách
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <button
              onClick={handleSearch}
              className="w-full md:w-auto mt-4 bg-primary-600 text-white px-8 py-3 rounded-lg font-semibold hover:bg-primary-700 transition-colors flex items-center justify-center space-x-2"
            >
              <Search className="h-5 w-5" />
              <span>Tìm kiếm</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Hero;
