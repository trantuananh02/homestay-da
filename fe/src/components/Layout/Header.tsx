import React, { useState, useEffect, useRef } from 'react';
import { Link, useLocation } from 'react-router-dom';
import {
  Home,
  User,
  Menu,
  X,
  Calendar,
  Building,
  LogOut
} from 'lucide-react';
import { useAuth } from '../../contexts/AuthContext';
import UserUpdateModal from '../Auth/UserUpdateModal';

const Header: React.FC = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const { user, logout, isAuthenticated } = useAuth();
  const location = useLocation();
  const [isUserModalOpen, setIsUserModalOpen] = useState(false);


  const isActive = (path: string) => location.pathname === path;

  const handleLogout = async () => {
    try {
      await logout();
      setIsDropdownOpen(false);
      setIsMenuOpen(false);
    } catch (error) {
      console.error('Logout error:', error);
    }
  };

  const getRoleDisplayName = (role: string) => {
    switch (role) {
      case 'host':
        return 'Chủ nhà';
      case 'admin':
        return 'Quản trị viên';
      case 'guest':
      default:
        return 'Khách hàng';
    }
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsDropdownOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <header className="bg-white shadow-sm border-b border-gray-100 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <Link to="/" className="flex items-center space-x-2">
            <Home className="h-8 w-8 text-primary-600" />
            <span className="text-xl font-bold text-gray-900">Mây Lang Thang</span>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center space-x-8">
            <Link
              to="/"
              className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${isActive('/') ? 'text-primary-600 bg-primary-50' : 'text-gray-700 hover:text-primary-600'
                }`}
            >
              Trang chủ
            </Link>
            <Link
              to="/homestays"
              className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${isActive('/homestays') ? 'text-primary-600 bg-primary-50' : 'text-gray-700 hover:text-primary-600'
                }`}
            >
              Homestay
            </Link>
            <Link
              to="/about"
              className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${isActive('/about') ? 'text-primary-600 bg-primary-50' : 'text-gray-700 hover:text-primary-600'
                }`}
            >
              Giới thiệu
            </Link>

            {isAuthenticated && user ? (
              <div className="flex items-center space-x-6">
                {(user.role === 'host' || user.role === 'admin') && (
                  <Link
                    to="/management"
                    className={`px-3 py-2 rounded-md text-sm font-medium transition-colors flex items-center space-x-1 ${isActive('/management') ? 'text-primary-600 bg-primary-50' : 'text-gray-700 hover:text-primary-600'
                      }`}
                  >
                    <Building className="h-4 w-4" />
                    <span>Quản lý</span>
                  </Link>
                )}

                {user.role === 'guest' && (
                  <Link
                    to="/bookings"
                    className={`px-3 py-2 rounded-md text-sm font-medium transition-colors flex items-center space-x-1 ${isActive('/bookings') ? 'text-primary-600 bg-primary-50' : 'text-gray-700 hover:text-primary-600'
                      }`}
                  >
                    <Calendar className="h-4 w-4" />
                    <span>Đặt phòng của tôi</span>
                  </Link>
                )}

                {/* User Dropdown */}
                <div className="relative" ref={dropdownRef}>
                  <button
                    onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                    className="flex items-center space-x-2 focus:outline-none"
                  >
                    <div className="w-8 h-8 bg-primary-100 rounded-full flex items-center justify-center">
                      <User className="h-4 w-4 text-primary-600" />
                    </div>
                    <div className="text-sm text-left">
                      <div className="text-gray-700 font-medium">{user.name}</div>
                      <div className="text-gray-500 text-xs">{getRoleDisplayName(user.role)}</div>
                    </div>
                  </button>

                  {isDropdownOpen && (
                    <div className="absolute right-0 mt-2 w-48 bg-white border border-gray-200 rounded-md shadow-lg z-10">
                      <button
                        onClick={() => {
                          setIsUserModalOpen(true);
                          setIsDropdownOpen(false);
                        }}
                        className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 w-full text-left"
                      >
                        Cập nhật thông tin
                      </button>
                      <button
                        onClick={handleLogout}
                        className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-100"
                      >
                        Đăng xuất
                      </button>
                    </div>
                  )}
                </div>
              </div>
            ) : (
              <div className="flex items-center space-x-3">
                <Link
                  to="/login"
                  className="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
                >
                  Đăng nhập
                </Link>
                <Link
                  to="/register"
                  className="bg-primary-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-primary-700 transition-colors"
                >
                  Đăng ký
                </Link>
              </div>
            )}
          </nav>

          {/* Mobile menu toggle */}
          <button className="md:hidden p-2" onClick={() => setIsMenuOpen(!isMenuOpen)}>
            {isMenuOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
          </button>
        </div>

        {/* Mobile menu */}
        {isMenuOpen && (
          <div className="md:hidden border-t border-gray-200 py-4">
            <div className="flex flex-col space-y-2">
              <Link
                to="/"
                onClick={() => setIsMenuOpen(false)}
                className="text-left px-3 py-2 text-gray-700 hover:text-primary-600"
              >
                Trang chủ
              </Link>
              <Link
                to="/homestays"
                onClick={() => setIsMenuOpen(false)}
                className="text-left px-3 py-2 text-gray-700 hover:text-primary-600"
              >
                Homestay
              </Link>
              <Link
                to="/about"
                onClick={() => setIsMenuOpen(false)}
                className="text-left px-3 py-2 text-gray-700 hover:text-primary-600"
              >
                Giới thiệu
              </Link>
              {isAuthenticated && user ? (
                <>
                  {(user.role === 'host' || user.role === 'admin') && (
                    <Link
                      to="/management"
                      onClick={() => setIsMenuOpen(false)}
                      className="text-left px-3 py-2 text-gray-700 hover:text-primary-600"
                    >
                      Quản lý
                    </Link>
                  )}
                  {user.role === 'guest' && (
                    <Link
                      to="/bookings"
                      onClick={() => setIsMenuOpen(false)}
                      className="text-left px-3 py-2 text-gray-700 hover:text-primary-600"
                    >
                      Đặt phòng của tôi
                    </Link>
                  )}
                  <button
                    onClick={() => setIsUserModalOpen(true)}
                    className="text-left px-3 py-2 text-gray-700 hover:text-primary-600"
                  >
                    Cập nhật thông tin
                  </button>
                  <button
                    onClick={handleLogout}
                    className="text-left px-3 py-2 text-red-600 hover:text-red-700 flex items-center space-x-2"
                  >
                    <LogOut className="h-4 w-4" />
                    <span>Đăng xuất</span>
                  </button>
                </>
              ) : (
                <div className="px-3 space-y-2">
                  <Link
                    to="/login"
                    onClick={() => setIsMenuOpen(false)}
                    className="w-full text-left px-3 py-2 text-gray-700 hover:text-primary-600"
                  >
                    Đăng nhập
                  </Link>
                  <Link
                    to="/register"
                    onClick={() => setIsMenuOpen(false)}
                    className="w-full text-left px-3 py-2 bg-primary-600 text-white rounded-md"
                  >
                    Đăng ký
                  </Link>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
      {user && user.id && (
        <UserUpdateModal
          isOpen={isUserModalOpen}
          onClose={() => setIsUserModalOpen(false)}
          userId={user?.id || 0} // Assuming user.id is available
        />
      )}

    </header>
  );
};

export default Header;
