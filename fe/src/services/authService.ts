import api from "./api";
import { toastService } from "./toastService";

// Types for authentication
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  phone?: string;
  password: string;
  role: "guest"; // Only allow guest registration
}

export interface LoginResponse {
  user: {
    id: number;
    name: string;
    phone?: string;
    email: string;
    role: string;
  };
  access_token: string;
  expires_in: number;
}

export interface ProfileResponse {
  user: {
    id: number;
    name: string;
    phone?: string;
    email: string;
    role: string;
  };
}

export interface VerifyEmailResponse {
  message: string;
}

// Authentication service
export const authService = {
  // Login user
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    try {
      const response = await api.post("/api/auth/login", credentials);
      toastService.showApiSuccess(response);
      return response.data;
    } catch (error: unknown) {
      toastService.showApiError(error);
      let message = "Đăng nhập thất bại";
      if (typeof error === "object" && error !== null && "response" in error) {
        const errObj = error as {
          response?: {
            data?: { result?: { message?: string }; message?: string };
          };
        };
        message =
          errObj.response?.data?.result?.message ||
          errObj.response?.data?.message ||
          message;
      }
      throw new Error(message);
    }
  },

  // Register new user
  async register(userData: RegisterRequest): Promise<LoginResponse> {
    try {
      const response = await api.post("/api/auth/register", userData);
      toastService.showApiSuccess(response);
      return response.data;
    } catch (error: unknown) {
      toastService.showApiError(error);
      let message = "Đăng ký thất bại";
      if (typeof error === "object" && error !== null && "response" in error) {
        const errObj = error as {
          response?: {
            data?: { result?: { message?: string }; message?: string };
          };
        };
        message =
          errObj.response?.data?.result?.message ||
          errObj.response?.data?.message ||
          message;
      }
      throw new Error(message);
    }
  },

  // Get user profile
  async getProfile(): Promise<ProfileResponse> {
    try {
      const response = await api.get("/api/auth/profile");
      return response.data;
    } catch (error: unknown) {
      toastService.showApiError(error);
      let message = "Không thể lấy thông tin profile";
      if (typeof error === "object" && error !== null && "response" in error) {
        const errObj = error as {
          response?: {
            data?: { result?: { message?: string }; message?: string };
          };
        };
        message =
          errObj.response?.data?.result?.message ||
          errObj.response?.data?.message ||
          message;
      }
      throw new Error(message);
    }
  },

  // Update user profile
  async updateUser(data: {
    name: string;
    email: string;
    phone?: string;
  }): Promise<void> {
    try {
      const response = await api.put("/api/auth/profile", data);
      toastService.showApiSuccess(response);

      return response.data;
    } catch (error: unknown) {
      let message = "Cập nhật thông tin thất bại";
      if (typeof error === "object" && error !== null && "response" in error) {
        const errObj = error as {
          response?: {
            data?: { result?: { message?: string }; message?: string };
          };
        };
        message =
          errObj.response?.data?.result?.message ||
          errObj.response?.data?.message ||
          message;
      }
      throw new Error(message);
    }
  },

  // Logout user
  async logout(): Promise<void> {
    try {
      await api.post("/api/auth/logout");
      toastService.success("Đăng xuất thành công");
    } catch (error: unknown) {
      console.error("Logout error:", error);
      toastService.showApiError(error);
    } finally {
      // Always clear local storage
      localStorage.removeItem("accessToken");
      localStorage.removeItem("user");
      localStorage.removeItem("tokenExpiry");
    }
  },

  // Verify email
  async verifyEmail(token: string): Promise<VerifyEmailResponse> {
    try {
      const response = await api.post("/api/auth/verify", { token });
      toastService.showApiSuccess(response);

      return response.data;
    } catch (error: unknown) {
      toastService.showApiError(error);
      let message = "Xác thực email thất bại";
      if (typeof error === "object" && error !== null && "response" in error) {
        const errObj = error as {
          response?: {
            data?: { result?: { message?: string }; message?: string };
          };
        };
        message =
          errObj.response?.data?.result?.message ||
          errObj.response?.data?.message ||
          message;
      }
      throw new Error(message);
    }
  },

  // Check if token is expired
  isTokenExpired(): boolean {
    const expiry = localStorage.getItem("tokenExpiry");
    if (!expiry) return true;

    const expiryTime = parseInt(expiry);
    const currentTime = Date.now();

    return currentTime >= expiryTime;
  },

  // Save auth data to localStorage
  saveAuthData(data: LoginResponse): void {
    const expiryTime = Date.now() + data.expires_in * 1000;

    localStorage.setItem("accessToken", data.access_token);
    localStorage.setItem("user", JSON.stringify(data.user));
    localStorage.setItem("tokenExpiry", expiryTime.toString());
  },

  // Get user from localStorage
  getUser(): import("../types").User | null {
    const userStr = localStorage.getItem("user");
    return userStr ? JSON.parse(userStr) : null;
  },

  // Get token from localStorage
  getToken(): string | null {
    return localStorage.getItem("accessToken");
  },

  // Check if user is authenticated
  isAuthenticated(): boolean {
    const token = this.getToken();
    const user = this.getUser();
    const isExpired = this.isTokenExpired();

    return !!(token && user && !isExpired);
  },
};
