import { toast, ToastOptions } from "react-toastify";

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
  showApiError: (error: unknown) => {
    let message = "Đã xảy ra lỗi không xác định";
    if (
      typeof error === "object" &&
      error !== null &&
      "response" in error &&
      typeof (error as Record<string, unknown>).response === "object" &&
      (error as Record<string, unknown>).response !== null
    ) {
      const response = (error as Record<string, unknown>).response;
      const data = (response as Record<string, unknown>).data as
        | Record<string, unknown>
        | undefined;
      if (data) {
        if (
          typeof data.result === "object" &&
          data.result !== null &&
          "message" in (data.result as object)
        ) {
          const msg = (data.result as Record<string, unknown>).message;
          if (typeof msg === "string") message = msg;
        }
        if ("message" in data && typeof data.message === "string") {
          message = data.message;
        }
      }
      // Nếu không có message trong data, kiểm tra message trực tiếp trên error
      if (
        "message" in error &&
        typeof (error as Record<string, unknown>).message === "string"
      ) {
        message = (error as Record<string, unknown>).message as string;
      }
    }
    toast.error(message, { ...defaultOptions });
  },

  // Hiển thị thông báo thành công từ response API
  showApiSuccess: (response: unknown) => {
    let message = "Thao tác thành công";
    if (
      typeof response === "object" &&
      response !== null &&
      "data" in response &&
      typeof (response as Record<string, unknown>).data === "object" &&
      (response as Record<string, unknown>).data !== null
    ) {
      const data = (response as Record<string, unknown>).data as
        | Record<string, unknown>
        | undefined;
      if (data) {
        if (
          typeof data.result === "object" &&
          data.result !== null &&
          "message" in (data.result as object)
        ) {
          const msg = (data.result as Record<string, unknown>).message;
          if (typeof msg === "string") message = msg;
        }
        if ("message" in data && typeof data.message === "string") {
          message = data.message;
        }
      }
    }
    toast.success(message, { ...defaultOptions });
  },

  // Xóa tất cả toast
  dismissAll: () => {
    toast.dismiss();
  },
};

export default toastService;
