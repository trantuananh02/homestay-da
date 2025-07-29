import React from 'react';
import AuthModal from './AuthModal';

interface LoginModalProps {
  isOpen: boolean;
  onClose: () => void;
  initialMode?: 'login' | 'register';
}

const LoginModal: React.FC<LoginModalProps> = ({ isOpen, onClose, initialMode = 'login' }) => {

  return (
    <AuthModal isOpen={isOpen} onClose={onClose} initialMode={initialMode} />
  );
};

export default LoginModal;