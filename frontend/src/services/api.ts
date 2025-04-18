import axios from 'axios';
import { jwtDecode } from 'jwt-decode';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const register = async (email: string, password: string) => {
  const response = await api.post('/register', { email, password });
  return response.data;
};

export const login = async (email: string, password: string) => {
  const response = await api.post('/login', { email, password });
  return response.data;
};

export const getBalance = async () => {
  const response = await api.get('/balance');
  return response.data;
};

export const earnPoints = async (points: number) => {
  const response = await api.post('/earn', { points });
  return response.data;
};

export const redeemPoints = async (points: number) => {
  const response = await api.post('/redeem', { points });
  return response.data;
};

export const getHistory = async () => {
  const response = await api.get('/history');
  return response.data;
};

export const getUserInfo = () => {
  const token = localStorage.getItem('token');
  if (!token) return null;
  
  try {
    const decoded = jwtDecode(token) as any;
    return {
      userId: decoded.user_id,
      accountId: decoded.account_id,
      email: localStorage.getItem('user_email') || 'User'
    };
  } catch (error) {
    console.error('Error decoding token:', error);
    return null;
  }
};