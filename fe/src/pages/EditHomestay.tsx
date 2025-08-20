import React, { useState, useEffect } from "react";
import { MapContainer, TileLayer, Marker, useMapEvents } from "react-leaflet";
import "leaflet/dist/leaflet.css";
import L, { Icon } from "leaflet";
import { useParams, useNavigate } from "react-router-dom";
import { ArrowLeft, Save, MapPin, Building } from "lucide-react";
import { useAuth } from "../contexts/AuthContext";
import { homestayService } from "../services/homestayService";
import { Homestay, UpdateHomestayRequest } from "../types";

const EditHomestay: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();

  const [homestay, setHomestay] = useState<Homestay | null>(null);
  const [formData, setFormData] = useState<UpdateHomestayRequest>({});

  // Icon đỏ cho marker
  const redIcon: Icon = new L.Icon({
    iconUrl: "https://maps.google.com/mapfiles/ms/icons/red-dot.png",
    iconSize: [32, 32],
    iconAnchor: [16, 32],
  });

  // Ref cho input tìm kiếm
  const searchRef = React.useRef<HTMLInputElement>(null);

  // Hàm tìm kiếm địa chỉ sử dụng Nominatim API
  const handleSearch = async () => {
    const query = searchRef.current?.value;
    if (!query) return;
    const url = `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(
      query
    )}`;
    const res = await fetch(url);
    const data: Array<{ lat: string; lon: string }> = await res.json();
    if (data && data.length > 0) {
      setFormData((prev) => ({
        ...prev,
        latitude: parseFloat(data[0].lat),
        longitude: parseFloat(data[0].lon),
      }));
    } else {
      alert("Không tìm thấy địa chỉ!");
    }
  };

  // Component chọn vị trí trên bản đồ
  const LocationPicker = () => {
    useMapEvents({
      click(e: { latlng: { lat: number; lng: number } }) {
        setFormData((prev) => ({
          ...prev,
          latitude: e.latlng.lat,
          longitude: e.latlng.lng,
        }));
      },
    });
    return null;
  };
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const homestayId = parseInt(id || "0");

  const loadHomestay = async () => {
    if (!homestayId) return;

    try {
      setLoading(true);
      const response = await homestayService.getHomestayById(homestayId);
      const homestayData = response.homestay;
      setHomestay(homestayData);
      setFormData({
        name: homestayData.name,
        description: homestayData.description,
        address: homestayData.address,
        city: homestayData.city,
        district: homestayData.district,
        ward: homestayData.ward,
        latitude: homestayData.latitude,
        longitude: homestayData.longitude,
        status: homestayData.status,
      });
    } catch (error) {
      console.error("Error loading homestay:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (user?.role === "host") {
      loadHomestay();
    }
  }, [user, homestayId]);

  const handleInputChange = (
    field: keyof UpdateHomestayRequest,
    value: string | number
  ) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  const handleBack = () => {
    navigate(`/management/homestay/${homestayId}`);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (
      !formData.name ||
      !formData.description ||
      !formData.address ||
      !formData.city ||
      !formData.district ||
      !formData.ward
    ) {
      alert("Vui lòng điền đầy đủ thông tin bắt buộc!");
      return;
    }

    if (formData.latitude === 0 || formData.longitude === 0) {
      alert("Vui lòng nhập tọa độ địa lý!");
      return;
    }

    try {
      setSaving(true);
      await homestayService.updateHomestay(homestayId, formData);
      navigate(`/management/homestay/${homestayId}`);
    } catch (error) {
      console.error("Error updating homestay:", error);
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Đang tải thông tin homestay...</p>
        </div>
      </div>
    );
  }

  if (!homestay) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600">Không tìm thấy homestay</p>
          <button
            onClick={() => navigate("/management")}
            className="mt-4 px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700"
          >
            Quay lại
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <button
          onClick={handleBack}
          className="flex items-center space-x-2 text-emerald-600 hover:text-emerald-700 mb-6"
        >
          <ArrowLeft className="h-5 w-5" />
          <span>Quay lại</span>
        </button>

        <div className="bg-white rounded-xl shadow-lg p-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-8">
            Chỉnh sửa homestay
          </h1>

          <form onSubmit={handleSubmit} className="space-y-8">
            {/* Basic Information */}
            <div className="space-y-6">
              <h2 className="text-xl font-semibold text-gray-900 flex items-center">
                <Building className="h-5 w-5 mr-2" />
                Thông tin cơ bản
              </h2>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Tên homestay *
                  </label>
                  <input
                    type="text"
                    value={formData.name || ""}
                    onChange={(e) => handleInputChange("name", e.target.value)}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    placeholder="Nhập tên homestay"
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Trạng thái
                  </label>
                  <select
                    value={formData.status || "active"}
                    onChange={(e) =>
                      handleInputChange("status", e.target.value)
                    }
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                  >
                    <option value="active">Hoạt động</option>
                    <option value="inactive">Không hoạt động</option>
                    <option value="pending">Chờ duyệt</option>
                  </select>
                </div>

                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Mô tả *
                  </label>
                  <textarea
                    value={formData.description || ""}
                    onChange={(e) =>
                      handleInputChange("description", e.target.value)
                    }
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    placeholder="Mô tả chi tiết về homestay"
                    rows={4}
                    required
                  />
                </div>
              </div>
            </div>

            {/* Address Information */}
            <div className="space-y-6">
              <h2 className="text-xl font-semibold text-gray-900 flex items-center">
                <MapPin className="h-5 w-5 mr-2" />
                Thông tin địa chỉ
              </h2>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Địa chỉ chi tiết *
                  </label>
                  <input
                    type="text"
                    value={formData.address || ""}
                    onChange={(e) =>
                      handleInputChange("address", e.target.value)
                    }
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    placeholder="Số nhà, tên đường"
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Phường/Xã *
                  </label>
                  <input
                    type="text"
                    value={formData.ward || ""}
                    onChange={(e) => handleInputChange("ward", e.target.value)}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    placeholder="Phường/Xã"
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Quận/Huyện *
                  </label>
                  <input
                    type="text"
                    value={formData.district || ""}
                    onChange={(e) =>
                      handleInputChange("district", e.target.value)
                    }
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    placeholder="Quận/Huyện"
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Tỉnh/Thành phố *
                  </label>
                  <input
                    type="text"
                    value={formData.city || ""}
                    onChange={(e) => handleInputChange("city", e.target.value)}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    placeholder="Tỉnh/Thành phố"
                    required
                  />
                </div>
              </div>
            </div>

            {/* Coordinates + OpenStreetMap */}
            <div className="space-y-6">
              <h2 className="text-xl font-semibold text-gray-900">
                Tọa độ địa lý
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Vĩ độ (Latitude) *
                  </label>
                  <input
                    type="number"
                    step="any"
                    value={formData.latitude || 0}
                    onChange={(e) =>
                      handleInputChange(
                        "latitude",
                        parseFloat(e.target.value) || 0
                      )
                    }
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    placeholder="10.762622"
                    required
                  />
                  <p className="text-xs text-gray-500 mt-1">
                    Ví dụ: 10.762622 (Hà Nội), 10.823099 (TP.HCM)
                  </p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Kinh độ (Longitude) *
                  </label>
                  <input
                    type="number"
                    step="any"
                    value={formData.longitude || 0}
                    onChange={(e) =>
                      handleInputChange(
                        "longitude",
                        parseFloat(e.target.value) || 0
                      )
                    }
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    placeholder="106.660172"
                    required
                  />
                  <p className="text-xs text-gray-500 mt-1">
                    Ví dụ: 106.660172 (Hà Nội), 106.629664 (TP.HCM)
                  </p>
                </div>
              </div>
              {/* OpenStreetMap + Tìm kiếm địa chỉ */}
              <div className="mt-6">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Chọn vị trí trên bản đồ:
                </label>
                <div className="flex gap-2 mb-2">
                  <input
                    ref={searchRef}
                    type="text"
                    placeholder="Nhập địa chỉ để tìm kiếm..."
                    className="p-2 border rounded w-full"
                  />
                  <button
                    type="button"
                    onClick={handleSearch}
                    className="px-4 py-2 bg-emerald-600 text-white rounded"
                  >
                    Tìm kiếm
                  </button>
                </div>
                <div className="w-full h-[400px] rounded-lg overflow-hidden border border-gray-300">
                  <MapContainer
                    center={[
                      typeof formData.latitude === "number"
                        ? formData.latitude
                        : 21.028511,
                      typeof formData.longitude === "number"
                        ? formData.longitude
                        : 105.804817,
                    ]}
                    zoom={
                      typeof formData.latitude === "number" &&
                      typeof formData.longitude === "number"
                        ? 15
                        : 6
                    }
                    style={{ width: "100%", height: "100%" }}
                  >
                    <TileLayer
                      url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                      attribution="© OpenStreetMap contributors"
                    />
                    <LocationPicker />
                    {formData.latitude !== 0 && formData.longitude !== 0 && (
                      <Marker
                        position={[
                          typeof formData.latitude === "number"
                            ? formData.latitude
                            : 0,
                          typeof formData.longitude === "number"
                            ? formData.longitude
                            : 0,
                        ]}
                        icon={redIcon}
                      />
                    )}
                  </MapContainer>
                </div>
                <p className="text-xs text-gray-500 mt-2">
                  Click vào bản đồ để chọn vị trí, tọa độ sẽ tự động điền vào ô
                  bên trên.
                </p>
              </div>
            </div>

            {/* Submit Button */}
            <div className="flex justify-end space-x-4 pt-6 border-t border-gray-200">
              <button
                type="button"
                onClick={handleBack}
                className="px-6 py-3 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
              >
                Hủy
              </button>
              <button
                type="submit"
                disabled={saving}
                className="px-6 py-3 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center"
              >
                {saving ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                    Đang lưu...
                  </>
                ) : (
                  <>
                    <Save className="h-4 w-4 mr-2" />
                    Lưu thay đổi
                  </>
                )}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default EditHomestay;
