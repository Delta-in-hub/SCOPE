// src/store/auth.js
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import * as authApi from '@/services/authApi';
import router from '@/router'; // 引入 router

export const useAuthStore = defineStore('auth', () => {
  // State (使用 ref)
  const accessToken = ref(localStorage.getItem('accessToken') || null);
  const refreshToken = ref(localStorage.getItem('refreshToken') || null);
  const user = ref(JSON.parse(localStorage.getItem('user') || '{}') || null); // 可以存储用户信息，如 display_name, email
  const expiresIn = ref(parseInt(localStorage.getItem('expiresIn') || '0', 10) || null);
  const tokenExpiryTime = ref(parseInt(localStorage.getItem('tokenExpiryTime') || '0', 10) || null); // Token 过期的精确时间戳

  // Getters (使用 computed)
  const isAuthenticated = computed(() => !!accessToken.value);

  // Actions (普通函数)
  function setTokens(authResult) {
    accessToken.value = authResult.access_token;
    refreshToken.value = authResult.refresh_token;
    expiresIn.value = authResult.expires_in; // 单位：秒

    // 计算精确的过期时间戳 (当前时间 + 有效期 - 少量缓冲时间，例如 60 秒)
    const now = Date.now();
    tokenExpiryTime.value = now + (authResult.expires_in * 1000) - 60000; // 提前 60 秒认为过期

    localStorage.setItem('accessToken', accessToken.value);
    localStorage.setItem('refreshToken', refreshToken.value);
    localStorage.setItem('expiresIn', expiresIn.value.toString());
    localStorage.setItem('tokenExpiryTime', tokenExpiryTime.value.toString());

    // 可以在登录成功后获取用户信息并存储 (如果 API 返回或有单独接口)
    // setUser({ email: 'user@example.com' }); // 示例
  }

  function setUser(userInfo) {
    user.value = userInfo;
    localStorage.setItem('user', JSON.stringify(userInfo));
  }

  function clearTokens() {
    accessToken.value = null;
    refreshToken.value = null;
    user.value = null;
    expiresIn.value = null;
    tokenExpiryTime.value = null;

    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('user');
    localStorage.removeItem('expiresIn');
    localStorage.removeItem('tokenExpiryTime');
  }

  async function login(credentials) {
    try {
      const response = await authApi.login(credentials);
      setTokens(response.data); // 假设 API 返回 { access_token, refresh_token, expires_in }
      // 可以在此获取用户信息
      // await fetchUserInfo(); // 假设有此方法
      return true; // 指示登录成功
    } catch (error) {
      console.error('Login failed:', error.response?.data || error.message);
      // 这里可以抛出错误或返回 false，并在组件中处理（如显示错误消息）
      throw error; // 让调用者知道失败了
    }
  }

  async function register(userInfo) {
      try {
          const response = await authApi.register(userInfo);
          console.log('Registration successful:', response.data);
          // 注册成功后，可以根据业务逻辑决定是自动登录还是跳转到登录页
          // setUser(response.data); // 假设 API 返回 { user_id, display_name, email }
          return response.data;
      } catch (error) {
          console.error('Registration failed:', error.response?.data || error.message);
          throw error;
      }
  }

  async function logout() {
      if (refreshToken.value) {
          try {
              // 调用后端 logout 接口（如果需要）
              await authApi.logout(refreshToken.value); // 注意：API 可能需要的是 access token 认证
              console.log('Logout successful on server.');
          } catch (error) {
              console.error('Server logout failed, clearing tokens locally anyway:', error.response?.data || error.message);
              // 即使后端调用失败，也应该清除本地 token
          }
      }
      clearTokens();
      // 跳转到登录页
      router.push({ name: 'Login' }); // 确保路由名称正确
  }

  async function refreshAccessToken() {
    if (!refreshToken.value) {
      throw new Error('No refresh token available.');
    }
    try {
      const response = await authApi.refreshToken(refreshToken.value);
      // 更新 access token 和过期时间
      accessToken.value = response.data.access_token;
      expiresIn.value = response.data.expires_in;
      const now = Date.now();
      tokenExpiryTime.value = now + (expiresIn.value * 1000) - 60000; // 重新计算过期时间戳

      localStorage.setItem('accessToken', accessToken.value);
      localStorage.setItem('expiresIn', expiresIn.value.toString());
      localStorage.setItem('tokenExpiryTime', tokenExpiryTime.value.toString());

      console.log('Access token refreshed.');
    } catch (error) {
      console.error('Failed to refresh access token:', error.response?.data || error.message);
      // 刷新失败，需要登出用户
      logout(); // 清理并重定向
      throw error; // 将错误传递给拦截器或调用者
    }
  }

  // 检查 token 是否即将过期 (可以在需要的地方调用，或在路由守卫中检查)
  function isTokenExpiringSoon() {
      if (!tokenExpiryTime.value) return true; // 没有过期时间信息，认为已过期
      return Date.now() >= tokenExpiryTime.value;
  }


  return {
    // state
    accessToken,
    refreshToken,
    user,
    expiresIn,
    tokenExpiryTime,
    // getters
    isAuthenticated,
    // actions
    login,
    register,
    logout,
    refreshAccessToken,
    isTokenExpiringSoon,
    setUser, // 如果需要单独设置用户信息
    clearTokens // 可能在某些情况下需要外部调用
  };
});