import React from "react";
import { AiOutlineMail } from "react-icons/ai";

interface VerifyEmailPopupProps {
  email: string;
  onClose: () => void;
  onResend?: () => void;
}

const VerifyEmailPopup: React.FC<VerifyEmailPopupProps> = ({ email, onClose, onResend }) => {
  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-40 z-50">
      <div className="bg-white rounded-lg shadow-lg p-6 max-w-sm w-full text-center animate-fadeIn">
        {/* Icon */}
        <AiOutlineMail className="mx-auto text-emerald-600 text-5xl mb-4" />

        {/* Title */}
        <h2 className="text-xl font-bold text-gray-900">Xác thực tài khoản</h2>

        {/* Description */}
        <p className="text-gray-700 mt-2">
          Chúng tôi đã gửi email xác thực đến:
          <br />
          <span className="font-semibold text-gray-900">{email}</span>
        </p>
        <p className="text-gray-500 text-sm mt-2">
          Vui lòng kiểm tra hộp thư của bạn và nhấn vào liên kết để kích hoạt tài khoản.
        </p>

        {/* Buttons */}
        <div className="mt-6 flex flex-col gap-3">
          <button
            onClick={onClose}
            className="font-semibold text-white transition-colors bg-primary-600 hover:bg-primary-70 py-2 px-4 rounded-lg font-medium"
          >
            Đã hiểu
          </button>
          {/* {onResend && (
            <button
              onClick={onResend}
              className="text-emerald-600 text-sm hover:underline"
            >
              Gửi lại email
            </button>
          )} */}
        </div>
      </div>
    </div>
  );
};

export default VerifyEmailPopup;
