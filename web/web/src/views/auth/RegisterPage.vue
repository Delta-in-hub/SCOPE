<template>
  <div class="register-container">
    <el-card class="register-card">
      <template #header>
        <div class="card-header">
          <span>用户注册</span>
        </div>
      </template>
      <el-form ref="registerFormRef" :model="registerForm" :rules="registerRules" label-width="100px" @submit.prevent="handleRegister">
        <el-form-item label="用户名" prop="display_name">
          <el-input v-model="registerForm.display_name" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="registerForm.email" placeholder="请输入邮箱" clearable />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="registerForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="registerForm.confirmPassword" type="password" placeholder="请再次输入密码" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading">注册</el-button>
          <el-button @click="goToLogin">返回登录</el-button>
        </el-form-item>
      </el-form>
      <el-alert v-if="errorMessage" :title="errorMessage" type="error" show-icon :closable="false" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/store/auth';
import { ElMessage } from 'element-plus';

// 表单引用
const registerFormRef = ref(null);
const loading = ref(false);
const errorMessage = ref('');

// 表单数据
const registerForm = reactive({
  display_name: '',
  email: '',
  password: '',
  confirmPassword: '',
});

// 密码验证函数
const validatePassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请输入密码'));
  } else if (value.length < 6) {
    callback(new Error('密码长度不能少于6个字符'));
  } else {
    if (registerForm.confirmPassword !== '') {
      // 如果确认密码已填写，则同时验证确认密码
      registerFormRef.value.validateField('confirmPassword');
    }
    callback();
  }
};

// 确认密码验证函数
const validateConfirmPassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'));
  } else if (value !== registerForm.password) {
    callback(new Error('两次输入密码不一致'));
  } else {
    callback();
  }
};

// 表单验证规则
const registerRules = {
  display_name: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度应在2-20个字符之间', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: ['blur', 'change'] },
  ],
  password: [
    { validator: validatePassword, trigger: 'blur' },
  ],
  confirmPassword: [
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
};

// 路由和 Store
const router = useRouter();
const authStore = useAuthStore();

// 注册处理
const handleRegister = async () => {
  if (!registerFormRef.value) return;
  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true;
      errorMessage.value = ''; // 清除之前的错误信息
      try {
        // 调用 store 中的注册方法
        const response = await authStore.register({
          display_name: registerForm.display_name,
          email: registerForm.email,
          password: registerForm.password,
        });
        
        ElMessage.success('注册成功，请登录');
        console.log('Registration successful:', response);
        router.push({ name: 'Login' });
      } catch (error) {
        // 从 error 对象中提取更具体的错误信息
        if (error.response?.status === 409) {
          errorMessage.value = '该邮箱已被注册';
        } else {
          errorMessage.value = error.response?.data?.message || error.response?.data || '注册失败，请稍后再试';
        }
        console.error('Registration failed in component:', error);
      } finally {
        loading.value = false;
      }
    } else {
      console.log('Form validation failed');
      return false;
    }
  });
};

// 跳转到登录页
const goToLogin = () => {
  router.push({ name: 'Login' });
};
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f0f2f5;
}

.register-card {
  width: 450px;
}

.card-header {
  text-align: center;
  font-size: 1.2em;
}
</style>
