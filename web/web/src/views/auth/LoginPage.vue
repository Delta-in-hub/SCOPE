<template>
  <div class="login-container">
    <div class="background-elements">
      <div class="bg-circle circle-1"></div>
      <div class="bg-circle circle-2"></div>
      <div class="bg-circle circle-3"></div>
    </div>
    

    
    <div class="content-wrapper">
      <div class="login-content">
        <div class="login-header">
          <h1 class="welcome-text">欢迎使用 Scope</h1>
          <p class="subtitle">登录以继续使用所有功能</p>
        </div>
        
        <div class="login-form-container">
        <el-form 
          ref="loginFormRef" 
          :model="loginForm" 
          :rules="loginRules" 
          label-position="top" 
          @submit.prevent="handleLogin"
          class="login-form"
        >
          <el-form-item prop="email">
            <el-input 
              v-model="loginForm.email" 
              placeholder="邮箱" 
              clearable 
              :prefix-icon="User"
              class="custom-input"
            />
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input 
              v-model="loginForm.password" 
              type="password" 
              placeholder="密码" 
              show-password 
              :prefix-icon="Lock"
              class="custom-input"
            />
          </el-form-item>
          
          <div class="form-footer">
            <el-button 
              type="primary" 
              native-type="submit" 
              :loading="loading"
              class="login-button"
              round
            >
              登录
            </el-button>
            
            <div class="register-link">
              还没有账号？<a @click="goToRegister">立即注册</a>
            </div>
          </div>
        </el-form>
        
        <el-alert 
          v-if="errorMessage" 
          :title="errorMessage" 
          type="error" 
          show-icon 
          :closable="false"
          class="error-alert" 
        />
      </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/store/auth';
import { ElMessage } from 'element-plus';
import { User, Lock } from '@element-plus/icons-vue';

// 表单引用
const loginFormRef = ref(null);
const loading = ref(false);
const errorMessage = ref('');

// 表单数据
const loginForm = reactive({
  email: '',
  password: '',
});

// 表单验证规则
const loginRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: ['blur', 'change'] },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
  ],
};

// 路由和 Store
const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

// 登录处理
const handleLogin = async () => {
  if (!loginFormRef.value) return;
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true;
      errorMessage.value = ''; // 清除之前的错误信息
      try {
        await authStore.login({
          email: loginForm.email,
          password: loginForm.password,
        });
        ElMessage.success('登录成功');
        // 登录成功后，检查查询参数中是否有 redirect，有则跳转，否则跳转到默认页
        const redirectPath = route.query.redirect || '/nodes'; // 默认跳转到 NodeList
        router.push(redirectPath);
      } catch (error) {
        // 从 error 对象中提取更具体的错误信息（如果后端返回了）
        errorMessage.value = error.response?.data?.message || error.response?.data || '登录失败，请检查邮箱或密码';
        console.error('Login failed in component:', error);
      } finally {
        loading.value = false;
      }
    } else {
      console.log('Form validation failed');
      return false;
    }
  });
};

// 跳转到注册页
const goToRegister = () => {
  router.push({ name: 'Register' });
};
</script>

<style scoped>
/* 全局容器 */
.login-container {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-height: 100vh;
  width: 100vw;
  background-color: #fafafa;
  font-family: 'SF Pro Display', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
  padding: 0;
  margin: 0;
  z-index: -1;
}

/* 背景元素 */
.background-elements {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: -1;
}

.bg-circle {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
}

.circle-1 {
  width: 70vw;
  height: 70vw;
  background: rgba(0, 122, 255, 0.2);
  top: -10vw;
  right: -10vw;
}

.circle-2 {
  width: 70vw;
  height: 70vw;
  background: rgba(88, 86, 214, 0.15);
  bottom: -10vw;
  left: -10vw;
}

.circle-3 {
  width: 30vw;
  height: 30vw;
  background: rgba(52, 199, 89, 0.1);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}





/* 内容包装器 */
.content-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  z-index: 1;
  position: relative;
}

/* 内容区域 */
.login-content {
  width: 90%;
  max-width: 480px;
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: 24px;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  overflow: hidden;
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.2);
  z-index: 2;
}

/* 头部样式 */
.login-header {
  padding: 50px 40px 20px;
  text-align: center;
}



.welcome-text {
  font-size: 32px;
  font-weight: 600;
  color: #1d1d1f;
  margin: 0 0 12px;
  letter-spacing: -0.5px;
  background: linear-gradient(135deg, #1d1d1f 0%, #434343 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.subtitle {
  font-size: 16px;
  color: #86868b;
  margin: 0;
  font-weight: 400;
}

/* 表单容器 */
.login-form-container {
  padding: 0 40px 40px;
}

.login-form {
  margin-top: 20px;
}

/* 自定义输入框 */
.custom-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.1);
  padding: 12px 15px;
  transition: box-shadow 0.3s ease;
}

.custom-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px #0071e3;
}

.custom-input :deep(.el-input__inner) {
  height: 24px;
  font-size: 16px;
}

.custom-input :deep(.el-input__prefix) {
  margin-right: 10px;
  color: #86868b;
}

/* 表单底部 */
.form-footer {
  margin-top: 30px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

/* 登录按钮 */
.login-button {
  width: 100%;
  height: 52px;
  font-size: 16px;
  font-weight: 500;
  background: linear-gradient(90deg, #0071e3, #42a5f5);
  border: none;
  transition: all 0.3s ease;
  margin-bottom: 24px;
  letter-spacing: 0.5px;
}

.login-button:hover {
  background: linear-gradient(90deg, #005bb5, #3994e4);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 113, 227, 0.3);
}

.login-button:active {
  transform: translateY(0);
}

/* 注册链接 */
.register-link {
  font-size: 14px;
  color: #86868b;
}

.register-link a {
  color: #0071e3;
  text-decoration: none;
  cursor: pointer;
  font-weight: 500;
  transition: color 0.2s ease;
}

.register-link a:hover {
  color: #005bb5;
  text-decoration: underline;
}

/* 错误提示 */
.error-alert {
  margin-top: 20px;
  border-radius: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-content {
    width: 95%;
    max-width: 100%;
    border-radius: 20px;
  }
  
  .login-header {
    padding: 30px 30px 15px;
  }
  
  .login-form-container {
    padding: 0 30px 30px;
  }
  
  .welcome-text {
    font-size: 24px;
  }
  
  .circle-1 {
    width: 100vw;
    height: 100vw;
  }
  
  .circle-2 {
    width: 100vw;
    height: 100vw;
  }
}
</style>