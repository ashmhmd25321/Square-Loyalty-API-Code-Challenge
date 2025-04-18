import React, { useEffect } from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const ProtectedRoute: React.FC = () => {
  const { isLoggedIn } = useAuth();
  
  useEffect(() => {
    // Debug log to check token and authentication state
    const token = localStorage.getItem('token');
    console.log('Protected route check - Token exists:', !!token);
    console.log('Protected route check - isLoggedIn state:', isLoggedIn);
  }, [isLoggedIn]);

  if (!isLoggedIn) {
    console.log('Not authenticated, redirecting to login');
    return <Navigate to="/login" replace />;
  }

  console.log('Authenticated, rendering protected content');
  return <Outlet />;
};

export default ProtectedRoute; 