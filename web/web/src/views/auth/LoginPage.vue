<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <span>用户登录</span>
        </div>
      </template>
      <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" label-width="80px" @submit.prevent="handleLogin">
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="loginForm.email" placeholder="请输入邮箱" clearable />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading">登录</el-button>
          <el-button @click="goToRegister">注册</el-button>
        </el-form-item>
      </el-form>
      <el-alert v-if="errorMessage" :title="errorMessage" type="error" show-icon :closable="false" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/store/auth';
import { ElMessage } from 'element-plus'; // 引入 ElMessage

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
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh; /* 使容器至少占满整个视口高度 */
  background-color: #f0f2f5; /* 添加一个浅灰色背景 */
}

.login-card {
  width: 400px;
}

.card-header {
  text-align: center;
  font-size: 1.2em;
}

/* 可以添加更多样式美化 */
</style>