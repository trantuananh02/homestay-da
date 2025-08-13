import React, { useEffect, useState } from "react";
import { CheckCircle, XCircle } from "lucide-react";
import { Link, useSearchParams } from "react-router-dom";
import { authService } from "../../services/authService";

const VerifyAccount: React.FC = () => {
  const [status, setStatus] = useState<"success" | "error" | "loading">("loading");
  const [searchParams] = useSearchParams();

  useEffect(() => {
    const token = searchParams.get("token");

    if (!token) {
      setStatus("error");
      return;
    }

    // Gọi API verify
    const verifyEmail = async () => {
      const res = await authService.verifyEmail(token);
        
      console.log("Verify response:", res);
      if (res.message === "Email đã được xác nhận thành công") {
        setStatus("success");
      } else {
        setStatus("error");
      }
    };

    verifyEmail();
  }, [searchParams]);

  if (status === "loading") {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <p className="text-gray-600 text-lg">Đang xác thực tài khoản...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center px-6">
      {status === "success" ? (
        <div className="bg-white p-8 rounded-lg shadow-lg text-center max-w-md w-full">
          <CheckCircle className="mx-auto text-emerald-600 w-16 h-16 mb-4" />
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Xác thực thành công!</h1>
          <p className="text-gray-600 mb-6">
            Tài khoản của bạn đã được kích hoạt. Bạn có thể đăng nhập và bắt đầu sử dụng.
          </p>
          <Link
            to="/login"
            className="bg-emerald-600 hover:bg-emerald-700 text-white px-6 py-3 rounded-lg font-medium"
          >
            Đăng nhập ngay
          </Link>
        </div>
      ) : (
        <div className="bg-white p-8 rounded-lg shadow-lg text-center max-w-md w-full">
          <XCircle className="mx-auto text-red-500 w-16 h-16 mb-4" />
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Xác thực thất bại</h1>
          <p className="text-gray-600 mb-6">
            Liên kết xác thực không hợp lệ hoặc đã hết hạn. Vui lòng yêu cầu gửi lại email.
          </p>
          <Link
            to="/resend-verification"
            className="bg-emerald-600 hover:bg-emerald-700 text-white px-6 py-3 rounded-lg font-medium"
          >
            Gửi lại email
          </Link>
        </div>
      )}
    </div>
  );
};

export default VerifyAccount;
