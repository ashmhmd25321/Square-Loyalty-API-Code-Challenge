import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

export const register = async (email: string, password: string) => {
  try {
    const response = await axios.post(`${API_BASE_URL}/register`, {
      email,
      password,
    });
    
    const { token } = response.data;
    if (token) {
      localStorage.setItem('token', token);
      localStorage.setItem('user_email', email);
      console.log('Token saved after registration:', token);
    }
    
    return response.data;
  } catch (error) {
    console.error('Registration error:', error);
    throw error;
  }
};

export const login = async (email: string, password: string) => {
  try {
    const response = await axios.post(`${API_BASE_URL}/login`, {
      email,
      password,
    });
    
    const { token } = response.data;
    if (token) {
      localStorage.setItem('token', token);
      localStorage.setItem('user_email', email);
      console.log('Token saved after login:', token);
    }
    
    return response.data;
  } catch (error) {
    console.error('Login error:', error);
    throw error;
  }
};

export const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user_email');
};

export const isAuthenticated = () => {
  const token = localStorage.getItem('token');
  return !!token;
}; 