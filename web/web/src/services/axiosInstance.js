// src/services/axiosInstance.js
import axios from 'axios';
import { useAuthStore } from '@/store/auth'; // 假设你的 auth store 路径
import router from '@/router'; // 引入 router

// 根据当前环境选择不同的 API 基础 URL
let baseURL;

// 判断是否是开发环境
// import.meta.env.MODE 在 Vite 中会自动设置为 'development' 或 'production'
if (import.meta.env.MODE === 'development') {
  // 开发环境使用环境变量或默认值
  baseURL = import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:18080/api/v1';
  console.log('开发环境使用 API 基础 URL:', baseURL);
} else {
  // 生产环境使用相对路径，让 Nginx 反向代理处理
  baseURL = '/api/v1';
  console.log('生产环境使用相对路径:', baseURL);
}

const axiosInstance = axios.create({
  baseURL: baseURL,
  timeout: 10000, // 请求超时时间
  headers: {
    'Content-Type': 'application/json',
  }
});

// 请求拦截器：添加 Authorization Header
axiosInstance.interceptors.request.use(
  config => {
    const authStore = useAuthStore(); // Pinia store 必须在 setup 或 action 中使用，但在拦截器外层定义引用
    const token = authStore.accessToken; // 从 Pinia store 获取 token

    // 检查是否是需要认证的接口 (根据你的 API 设计决定，这里简单判断 config.url 是否包含 'auth/login' 或 'auth/register')
    const requiresAuth = !config.url.includes('/auth/login') && !config.url.includes('/auth/register') && !config.url.includes('/node/up') && !config.url.includes('/node/down') && !config.url.includes('/auth/refreshToken');

    if (token && requiresAuth) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    console.error('Request interceptor error:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器：处理 401 未授权等错误，可以尝试刷新 token
axiosInstance.interceptors.response.use(
  response => {
    // 对响应数据做点什么
    return response;
  },
  async error => {
    const originalRequest = error.config;
    const authStore = useAuthStore();

    // 处理 401 错误 (未授权)
    if (error.response && error.response.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true; // 标记以防止无限重试循环

        // 检查是否是刷新 token 的请求本身失败了，或者是没有 refresh token
        if (originalRequest.url.includes('/auth/refreshToken') || !authStore.refreshToken) {
            console.error('Refresh token failed or not available. Logging out.');
            authStore.logout(); // 清除 token 并重定向到登录页
            router.push({ name: 'Login' }); // 确保 router 已正确配置
            return Promise.reject(error);
        }

        try {
            console.log('Access token expired, attempting to refresh...');
            await authStore.refreshAccessToken(); // 调用 Pinia store 中的刷新 token action
            console.log('Token refreshed successfully.');
            // 更新原始请求的 Authorization header
            originalRequest.headers.Authorization = `Bearer ${authStore.accessToken}`;
            // 重新发送原始请求
            return axiosInstance(originalRequest);
        } catch (refreshError) {
            console.error('Unable to refresh token:', refreshError);
            authStore.logout(); // 刷新失败，登出
            router.push({ name: 'Login' });
            return Promise.reject(refreshError);
        }
    }

    // 处理其他错误 (可以根据需要添加，例如 400, 403, 500 等)
    if (error.response) {
        console.error(`Error ${error.response.status}:`, error.response.data);
        // 可以使用 ElMessage 显示错误信息给用户
        // import { ElMessage } from 'element-plus'; ElMessage.error(...)
    } else if (error.request) {
        console.error('Network error or no response:', error.request);
    } else {
        console.error('Axios setup error:', error.message);
    }

    return Promise.reject(error); // 将错误继续传递下去
  }
);

export default axiosInstance;