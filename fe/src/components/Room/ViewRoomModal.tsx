import React, { useState, useEffect } from 'react';
import { X, Edit, Plus } from 'lucide-react';
import { Room } from '../../types';
import { homestayService } from '../../services/homestayService';

interface ViewRoomModalProps {
  isOpen: boolean;
  onClose: () => void;
  room: Room | null;
  isEdit: boolean;
  readonly?: boolean; // New prop to indicate read-only mode
}

const roomTypes = [
  { value: 'Standard', label: 'Phòng Standard' },
  { value: 'Deluxe', label: 'Phòng Deluxe' },
  { value: 'Premium', label: 'Phòng Premium' },
  { value: 'Suite', label: 'Phòng Suite' }
];

const commonAmenities = [
  'Wi-Fi', 'Điều hòa', 'TV', 'Tủ lạnh', 'Máy sấy tóc', 'Két an toàn',
  'Ban công', 'Tầm nhìn ra biển', 'Tầm nhìn ra núi', 'Phòng tắm riêng',
  'Bồn tắm', 'Vòi sen', 'Đồ vệ sinh cá nhân', 'Khăn tắm'
];

const ViewRoomModal: React.FC<ViewRoomModalProps> = ({ isOpen, onClose, room, isEdit, readonly }) => {
  const [isEditing, setIsEditing] = useState(isEdit);
  const [formData, setFormData] = useState<Room | null>(room);
  const [newAmenity, setNewAmenity] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    setFormData(room);
    setIsEditing(false);
  }, [room, isOpen]);

  const handleFormChange = (field: keyof Room, value: any) => {
    if (!formData) return;
    setFormData({ ...formData, [field]: value });
  };

  const addAmenityToList = (amenity: string) => {
    if (!formData) return;
    if (amenity && !formData.amenities?.includes(amenity)) {
      setFormData(prev => prev ? { ...prev, amenities: [...(prev.amenities || []), amenity] } : prev);
    }
  };
  const removeAmenity = (amenity: string) => {
    if (!formData) return;
    setFormData(prev => prev ? { ...prev, amenities: (prev.amenities || []).filter(a => a !== amenity) } : prev);
  };
  const addCustomAmenity = () => {
    if (newAmenity.trim()) {
      addAmenityToList(newAmenity.trim());
      setNewAmenity('');
    }
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData) return;
    setIsSaving(true);
    try {
      await homestayService.updateRoom(formData.id, {
        name: formData.name,
        description: formData.description,
        type: formData.type as 'Standard' | 'Deluxe' | 'Premium' | 'Suite',
        capacity: formData.capacity,
        price: formData.price,
        status: formData.status,
        amenities: formData.amenities || [],
        images: formData.images || []
      });
      setIsEditing(false);
    } catch (error) {
      alert('Có lỗi khi cập nhật phòng!');
    } finally {
      setIsSaving(false);
    }
  };

  console.log(formData);
  console.log(isEdit);

  if (!isOpen || !formData) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-gray-900">Chi tiết phòng</h2>
            <div className="flex items-center gap-2">
              {!readonly && !isEditing && (
                <button
                  onClick={() => setIsEditing(true)}
                  className="text-blue-600 hover:text-blue-800 p-1 rounded-full border border-blue-200 hover:bg-blue-50 flex items-center gap-1 px-4"
                  title="Sửa phòng"
                >
                  <Edit className="h-5 w-5" />
                  Sửa
                </button>
              )}
              <button
                onClick={onClose}
                className="text-gray-400 hover:text-gray-600 ml-2"
              >
                <X className="h-6 w-6" />
              </button>
            </div>
          </div>

          {isEditing ? (
            <form onSubmit={handleSave} className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Tên phòng *</label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={e => handleFormChange('name', e.target.value)}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Loại phòng *</label>
                  <select
                    value={formData.type}
                    onChange={e => handleFormChange('type', e.target.value)}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    required
                  >
                    {roomTypes.map(type => (
                      <option key={type.value} value={type.value}>{type.label}</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Sức chứa *</label>
                  <input
                    type="number"
                    value={formData.capacity}
                    onChange={e => handleFormChange('capacity', parseInt(e.target.value))}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Giá phòng/đêm *</label>
                  <input
                    type="number"
                    value={formData.price}
                    onChange={e => handleFormChange('price', parseInt(e.target.value))}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    required
                  />
                </div>
                {!readonly && (<div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Trạng thái</label>
                  <select
                    value={formData.status}
                    onChange={e => handleFormChange('status', e.target.value)}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                  >
                    <option value="available">Sẵn sàng</option>
                    <option value="occupied">Đang sử dụng</option>
                    <option value="maintenance">Bảo trì</option>
                  </select>
                </div>)}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Mô tả</label>
                <textarea
                  value={formData.description}
                  onChange={e => handleFormChange('description', e.target.value)}
                  rows={3}
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-3">
                  Tiện nghi phòng
                </label>

                {/* Common Amenities */}
                <div className="grid grid-cols-2 md:grid-cols-3 gap-2 mb-4">
                  {commonAmenities.map((amenity) => (
                    <button
                      key={amenity}
                      type="button"
                      onClick={() =>
                        formData.amenities?.includes(amenity)
                          ? removeAmenity(amenity)
                          : addAmenityToList(amenity)
                      }
                      className={`p-2 text-sm rounded-lg border transition-colors ${formData.amenities?.includes(amenity)
                        ? 'bg-emerald-100 border-emerald-300 text-emerald-800'
                        : 'bg-white border-gray-300 text-gray-700 hover:bg-gray-50'
                        }`}
                    >
                      {amenity}
                    </button>
                  ))}
                </div>

                {/* Custom Amenity */}
                <div className="flex space-x-2 mb-4">
                  <input
                    type="text"
                    value={newAmenity}
                    onChange={(e) => setNewAmenity(e.target.value)}
                    className="flex-1 p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                    placeholder="Thêm tiện nghi khác..."
                    onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addCustomAmenity())}
                  />
                  <button
                    type="button"
                    onClick={addCustomAmenity}
                    className="bg-emerald-600 text-white px-4 py-3 rounded-lg hover:bg-emerald-700 transition-colors"
                  >
                    <Plus className="h-5 w-5" />
                  </button>
                </div>

                {/* Selected Amenities */}
                {Array.isArray(formData.amenities) && formData.amenities.length > 0 && (
                  <div className="flex flex-wrap gap-2">
                    {formData.amenities.map((amenity) => (
                      <span
                        key={amenity}
                        className="bg-emerald-100 text-emerald-800 px-3 py-1 rounded-full text-sm flex items-center space-x-1"
                      >
                        <span>{amenity}</span>
                        <button
                          type="button"
                          onClick={() => removeAmenity(amenity)}
                          className="text-emerald-600 hover:text-emerald-800"
                        >
                          <X className="h-3 w-3" />
                        </button>
                      </span>
                    ))}
                  </div>
                )}
              </div>
              {/* Ảnh phòng */}
              {formData.images && formData.images.length > 0 && (
                <div className="mb-6 flex flex-wrap gap-3">
                  {formData.images.map((img, idx) => (
                    <img
                      key={idx}
                      src={img}
                      alt={`Ảnh phòng ${idx + 1}`}
                      className="w-32 h-24 object-cover rounded-lg border"
                    />
                  ))}
                </div>
              )}
              {/* TODO: Thêm upload ảnh nếu cần */}
              <div className="flex gap-2 mt-4">
                <button type="submit" disabled={isSaving} className="bg-emerald-600 text-white px-4 py-2 rounded-lg">{isSaving ? 'Đang lưu...' : 'Lưu'}</button>
                <button type="button" onClick={() => setIsEditing(false)} className="ml-2 px-4 py-2 rounded-lg border">Hủy</button>
              </div>
            </form>
          ) : (
            <>
              {/* Ảnh phòng */}
              {formData.images && formData.images.length > 0 && (
                <div className="mb-6 flex flex-wrap gap-3">
                  {formData.images.map((img, idx) => (
                    <img
                      key={idx}
                      src={img}
                      alt={`Ảnh phòng ${idx + 1}`}
                      className="w-32 h-24 object-cover rounded-lg border"
                    />
                  ))}
                </div>
              )}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Tên phòng</label>
                  <div className="p-3 border border-gray-200 rounded-lg bg-gray-50">{formData.name}</div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Loại phòng</label>
                  <div className="p-3 border border-gray-200 rounded-lg bg-gray-50">{formData.type}</div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Sức chứa</label>
                  <div className="p-3 border border-gray-200 rounded-lg bg-gray-50">{formData.capacity} người</div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Giá phòng/đêm</label>
                  <div className="p-3 border border-gray-200 rounded-lg bg-gray-50">{formData.price.toLocaleString()} VNĐ</div>
                </div>
                {!readonly && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">Trạng thái</label>
                    <div className="p-3 border border-gray-200 rounded-lg bg-gray-50">{formData.status}</div>
                  </div>
                )}
              </div>
              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-700 mb-2">Mô tả</label>
                <div className="p-3 border border-gray-200 rounded-lg bg-gray-50">{formData.description}</div>
              </div>
              {/* Tiện ích */}

                {console.log(formData)}

              {Array.isArray(formData.amenities) && formData.amenities.length > 0 && (
                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 mb-2">Tiện ích</label>
                  <div className="flex flex-wrap gap-2">
                    {formData.amenities.map((amenity, idx) => (
                      <span
                        key={idx}
                        className="px-3 py-1 bg-emerald-100 text-emerald-700 rounded-full text-xs font-medium"
                      >
                        {amenity}
                      </span>
                    ))}
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default ViewRoomModal;
