// src/services/authApi.js
import axiosInstance from './axiosInstance';

export const login = (credentials) => {
  // credentials 应为 { email: '...', password: '...' }
  return axiosInstance.post('/auth/login', credentials);
};

export const register = (userInfo) => {
  // userInfo 应为 { display_name: '...', email: '...', password: '...' }
  return axiosInstance.post('/auth/register', userInfo);
};

export const logout = (refreshToken) => {
    // refreshToken 应为 { refresh_token: '...' }
    // 注意：Logout 通常也需要认证（带 Access Token），取决于后端实现
    // 如果需要 Access Token，拦截器会自动添加
  return axiosInstance.post('/auth/logout', { refresh_token: refreshToken });
};

export const refreshToken = (refreshToken) => {
  // refreshToken 应为 { refresh_token: '...' }
  return axiosInstance.post('/auth/refreshToken', { refresh_token: refreshToken });
};