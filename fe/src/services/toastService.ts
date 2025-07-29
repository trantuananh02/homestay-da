import { toast, ToastOptions } from 'react-toastify';

// Cấu hình mặc định cho toast
const defaultOptions: ToastOptions = {
  position: "top-right",
  autoClose: 5000,
  hideProgressBar: false,
  closeOnClick: true,
  pauseOnHover: true,
  draggable: true,
};

// Toast service để quản lý thông báo
export const toastService = {
  // Hiển thị thông báo thành công
  success: (message: string, options?: ToastOptions) => {
    toast.success(message, { ...defaultOptions, ...options });
  },

  // Hiển thị thông báo lỗi
  error: (message: string, options?: ToastOptions) => {
    toast.error(message, { ...defaultOptions, ...options });
  },

  // Hiển thị thông báo cảnh báo
  warning: (message: string, options?: ToastOptions) => {
    toast.warning(message, { ...defaultOptions, ...options });
  },

  // Hiển thị thông báo thông tin
  info: (message: string, options?: ToastOptions) => {
    toast.info(message, { ...defaultOptions, ...options });
  },

  // Hiển thị thông báo từ response API
  showApiError: (error: any) => {
    let message = 'Đã xảy ra lỗi không xác định';
    
    if (error?.response?.data?.result?.message) {
      message = error.response.data.result.message;
    } else if (error?.response?.data?.message) {
      message = error.response.data.message;
    } else if (error?.message) {
      message = error.message;
    }

    toast.error(message, { ...defaultOptions });
  },

  // Hiển thị thông báo thành công từ response API
  showApiSuccess: (response: any) => {
    let message = 'Thao tác thành công';
    
    if (response?.data?.result?.message) {
      message = response.data.result.message;
    } else if (response?.data?.message) {
      message = response.data.message;
    }

    toast.success(message, { ...defaultOptions });
  },

  // Xóa tất cả toast
  dismissAll: () => {
    toast.dismiss();
  },
};

export default toastService; 