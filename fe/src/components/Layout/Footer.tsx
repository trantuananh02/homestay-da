import React from 'react';
import { Home, Facebook, Instagram, Twitter, Mail, Phone, MapPin } from 'lucide-react';

const Footer: React.FC = () => {
  return (
    <footer className="bg-gray-900 text-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* Company Info */}
          <div className="col-span-1 md:col-span-2">
            <div className="flex items-center space-x-2 mb-4">
              <Home className="h-8 w-8 text-primary-400" />
              <span className="text-xl font-bold">Mây Lang Thang</span>
            </div>
            <p className="text-gray-300 mb-4 max-w-md">
              Nền tảng homestay hàng đầu Việt Nam, kết nối du khách với những trải nghiệm 
              lưu trú đích thực và đáng nhớ khắp đất nước.
            </p>
            <div className="flex space-x-4">
              <a href="#" className="text-gray-300 hover:text-primary-400 transition-colors">
                <Facebook className="h-6 w-6" />
              </a>
              <a href="#" className="text-gray-300 hover:text-primary-400 transition-colors">
                <Instagram className="h-6 w-6" />
              </a>
              <a href="#" className="text-gray-300 hover:text-primary-400 transition-colors">
                <Twitter className="h-6 w-6" />
              </a>
            </div>
          </div>

          {/* Quick Links */}
          <div>
            <h3 className="text-lg font-semibold mb-4">Liên kết nhanh</h3>
            <ul className="space-y-2">
              <li><a href="#" className="text-gray-300 hover:text-primary-400 transition-colors">Trang chủ</a></li>
              <li><a href="#" className="text-gray-300 hover:text-primary-400 transition-colors">Homestay</a></li>
              <li><a href="#" className="text-gray-300 hover:text-primary-400 transition-colors">Giới thiệu</a></li>
              <li><a href="#" className="text-gray-300 hover:text-primary-400 transition-colors">Liên hệ</a></li>
              <li><a href="#" className="text-gray-300 hover:text-primary-400 transition-colors">Hỗ trợ</a></li>
            </ul>
          </div>

          {/* Contact Info */}
          <div>
            <h3 className="text-lg font-semibold mb-4">Liên hệ</h3>
            <div className="space-y-3">
              <div className="flex items-center space-x-2">
                <MapPin className="h-5 w-5 text-primary-400" />
                <span className="text-gray-300 text-sm">123 Đường Lê Lợi, Q1, TP.HCM</span>
              </div>
              <div className="flex items-center space-x-2">
                <Phone className="h-5 w-5 text-primary-400" />
                <span className="text-gray-300 text-sm">+84 28 1234 5678</span>
              </div>
              <div className="flex items-center space-x-2">
                <Mail className="h-5 w-5 text-primary-400" />
                <span className="text-gray-300 text-sm">info@homestayvietnam.com</span>
              </div>
            </div>
          </div>
        </div>

        <div className="border-t border-gray-800 mt-8 pt-8">
          <div className="flex flex-col md:flex-row justify-between items-center">
            <div className="text-gray-300 text-sm">
              © 2025 Mây Lang Thang. Tất cả quyền được bảo lưu.
            </div>
            <div className="flex space-x-6 mt-4 md:mt-0">
              <a href="#" className="text-gray-300 hover:text-primary-400 text-sm transition-colors">
                Điều khoản sử dụng
              </a>
              <a href="#" className="text-gray-300 hover:text-primary-400 text-sm transition-colors">
                Chính sách bảo mật
              </a>
              <a href="#" className="text-gray-300 hover:text-primary-400 text-sm transition-colors">
                Chính sách hoàn tiền
              </a>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;