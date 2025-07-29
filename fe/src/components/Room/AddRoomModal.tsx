import React, { useEffect, useState } from 'react';
import { X, Plus } from 'lucide-react';
import { TbXboxX } from 'react-icons/tb';
import { Room } from '../../types';
import CusFormUpload from '../UploadFile';
import { homestayService } from '../../services/homestayService';

interface AddRoomModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (room: Omit<Room, 'id' | 'createdAt'>) => void;
  homestayId: number;
  room?: Room | null;
  action?: 'add' | 'edit' | 'view';
}

const AddRoomModal: React.FC<AddRoomModalProps> = ({
  isOpen,
  onClose,
  onSubmit,
  homestayId,
  room,
  action = 'add',
}) => {
  const isView = action === 'view';
  const isEdit = action === 'edit';
  const isAdd = action === 'add';

  const [isUploading, setIsUploading] = useState(false);
  const [formData, setFormData] = useState<Omit<Room, 'id' | 'createdAt'>>({
    name: '',
    type: 'Standard',
    capacity: 2,
    price: 0,
    priceType: 'per_night', // or whatever default value is appropriate
    description: '',
    amenities: [],
    images: [],
    status: 'available',
    homestayId,
  });

  const [newAmenity, setNewAmenity] = useState('');

  useEffect(() => {
    if ((isEdit || isView) && room) {
      const { id, createdAt, ...rest } = room;
      setFormData({ ...rest, price: rest.price });
    } else if (isAdd) {
      setFormData({
        name: '',
        type: 'Standard',
        capacity: 2,
        price: 0,
        priceType: 'per_night', // or whatever default value is appropriate
        description: '',
        amenities: [],
        images: [],
        status: 'available',
        homestayId,
      });
    }
  }, [room, isEdit, isView, isAdd, homestayId]);

  const roomTypes = [
    { value: 'Standard', label: 'Phòng Standard' },
    { value: 'Deluxe', label: 'Phòng Deluxe' },
    { value: 'Premium', label: 'Phòng Premium' },
    { value: 'Suite', label: 'Phòng Suite' },
  ];

  const commonAmenities = [
    'Wi-Fi', 'Điều hòa', 'TV', 'Tủ lạnh', 'Máy sấy tóc', 'Két an toàn',
    'Ban công', 'Tầm nhìn ra biển', 'Tầm nhìn ra núi', 'Phòng tắm riêng',
    'Bồn tắm', 'Vòi sen', 'Đồ vệ sinh cá nhân', 'Khăn tắm',
  ];

  const addAmenityToList = (amenity: string) => {
    if (amenity && !formData.amenities?.includes(amenity)) {
      setFormData((prev) => ({
        ...prev,
        amenities: [...(prev?.amenities ?? []), amenity],
      }));
    }
  };

  const removeAmenity = (amenity: string) => {
    setFormData((prev) => ({
      ...prev,
      amenities: prev.amenities?.filter((a) => a !== amenity),
    }));
  };

  const addCustomAmenity = () => {
    if (newAmenity.trim()) {
      addAmenityToList(newAmenity.trim());
      setNewAmenity('');
    }
  };

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

    setFormData((prev) => ({
      ...prev,
      images: [...(prev?.images ?? []), ...(uploaded.filter(Boolean) as string[])],
    }));
    setIsUploading(false);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (isView) return;

    if (!formData.name || !formData.price || !formData.description) {
      alert('Vui lòng điền đầy đủ thông tin bắt buộc!');
      return;
    }

    if (formData.images?.length === 0) {
      alert('Vui lòng thêm ít nhất một hình ảnh!');
      return;
    }

    const roomData = {
      ...formData,
      price: formData.price,
    };

    onSubmit(roomData);
    onClose();
  };

  if (!isOpen) return null;

  const title = isAdd ? 'Thêm phòng mới' : isEdit ? 'Chỉnh sửa phòng' : 'Xem chi tiết phòng';
  const submitLabel = isAdd ? 'Thêm phòng' : 'Lưu thay đổi';

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-gray-900">{title}</h2>
            <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
              <X className="h-6 w-6" />
            </button>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <span className='text-sm font-medium text-gray-700 mb-3 block'>Tên phòng</span>
                <input
                  disabled={isView}
                  required
                  type="text"
                  placeholder="VD: Phòng Deluxe 101"
                  className="w-full p-3 border rounded-lg"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                />
              </div>

              <div>
                <span className='text-sm font-medium text-gray-700 mb-3 block'>Loại phòng</span>
                <select
                  disabled={isView}
                  required
                  value={formData.type}
                  onChange={(e) => setFormData({ ...formData, type: e.target.value })}
                  className="w-full p-3 border rounded-lg"
                >
                  {roomTypes.map((type) => (
                    <option key={type.value} value={type.value}>
                      {type.label}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <span className='text-sm font-medium text-gray-700 mb-3 block'>Số người</span>
                <input
                  disabled={isView}
                  required
                  type="number"
                  min="1"
                  placeholder="2"
                  className="w-full p-3 border rounded-lg"
                  value={formData.capacity}
                  onChange={(e) => setFormData({ ...formData, capacity: parseInt(e.target.value) })}
                />
              </div>

              <div>
                <span className='text-sm font-medium text-gray-700 mb-3 block'>Giá phòng</span>
                <input
                  disabled={isView}
                  required
                  type="number"
                  placeholder="500000"
                  className="w-full p-3 border rounded-lg"
                  value={formData.price}
                  onChange={(e) => setFormData({ ...formData, price: parseInt(e.target.value) })}
                />

              </div>
            </div>

            <textarea
              disabled={isView}
              required
              rows={3}
              className="w-full p-3 border rounded-lg"
              placeholder="Mô tả chi tiết về phòng..."
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            />

            <div>
              <label className="text-sm font-medium text-gray-700 mb-3 block">Tiện nghi phòng</label>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-2 mb-4">
                {commonAmenities.map((amenity) => (
                  <button
                    key={amenity}
                    disabled={isView}
                    type="button"
                    onClick={() =>
                      formData.amenities?.includes(amenity)
                        ? removeAmenity(amenity)
                        : addAmenityToList(amenity)
                    }
                    className={`p-2 text-sm rounded-lg border ${formData.amenities?.includes(amenity)
                        ? 'bg-emerald-100 border-emerald-300'
                        : 'bg-white border-gray-300'
                      }`}
                  >
                    {amenity}
                  </button>
                ))}
              </div>

              {!isView && (
                <div className="flex space-x-2 mb-4">
                  <input
                    type="text"
                    value={newAmenity}
                    onChange={(e) => setNewAmenity(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addCustomAmenity())}
                    className="flex-1 p-3 border rounded-lg"
                    placeholder="Thêm tiện nghi khác..."
                  />
                  <button
                    type="button"
                    onClick={addCustomAmenity}
                    className="bg-emerald-600 text-white px-4 py-2 rounded-lg"
                  >
                    <Plus className="h-5 w-5" />
                  </button>
                </div>
              )}

              <div className="flex flex-wrap gap-2">
                {formData.amenities?.map((amenity) => (
                  <span
                    key={amenity}
                    className="bg-emerald-100 text-emerald-800 px-3 py-1 rounded-full text-sm flex items-center space-x-1"
                  >
                    <span>{amenity}</span>
                    {!isView && (
                      <button
                        type="button"
                        onClick={() => removeAmenity(amenity)}
                        className="text-emerald-600 hover:text-emerald-800"
                      >
                        <X className="h-3 w-3" />
                      </button>
                    )}
                  </span>
                ))}
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-3">Hình ảnh phòng *</label>
              <div className="flex flex-wrap gap-4">
                {formData.images?.map((image, index) => (
                  <div key={index} className="relative">
                    <img
                      src={image}
                      alt={`Ảnh phòng ${index + 1}`}
                      className="w-40 h-40 object-cover rounded-lg"
                    />
                    {!isView && (
                      <TbXboxX
                        className="text-red-500 text-2xl absolute top-2 right-2 cursor-pointer"
                        onClick={() =>
                          setFormData((prev) => ({
                            ...prev,
                            images: prev?.images?.filter((img) => img !== image),
                          }))
                        }
                      />
                    )}
                  </div>
                ))}
                {!isView && (
                  <CusFormUpload
                    disabled={false}
                    handleUpload={handleImageUpload}
                    isUploading={isUploading}
                  />
                )}
              </div>
            </div>

            {!isView && (
              <div className="flex space-x-4 pt-6 border-t">
                <button
                  type="button"
                  onClick={onClose}
                  className="flex-1 bg-gray-300 text-gray-700 py-3 rounded-lg"
                >
                  Hủy
                </button>
                <button
                  type="submit"
                  className="flex-1 bg-emerald-600 text-white py-3 rounded-lg"
                >
                  {submitLabel}
                </button>
              </div>
            )}
          </form>
        </div>
      </div>
    </div>
  );
};

export default AddRoomModal;
