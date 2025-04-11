<template>
  <div class="app-header">
    <div class="logo">
      <img src="@/assets/logo.svg" alt="Logo" />
      <h1>Scope Center</h1>
    </div>
    <div class="spacer"></div>
    <div class="user-info">
      <el-dropdown trigger="click" @command="handleCommand">
        <span class="user-dropdown-link">
          <el-avatar :size="32" icon="UserFilled" />
          <span class="username">{{ user?.display_name || '用户' }}</span>
          <el-icon><ArrowDown /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>个人信息
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useAuthStore } from '@/store/auth';
import { ElMessageBox } from 'element-plus';

const authStore = useAuthStore();
const user = computed(() => authStore.user);

const handleCommand = (command) => {
  if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }).then(() => {
      authStore.logout();
    }).catch(() => {
      // 用户取消操作
    });
  } else if (command === 'profile') {
    // 可以添加跳转到个人信息页面的逻辑
    console.log('Navigate to profile page');
  }
};
</script>

<style scoped>
.app-header {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 0 20px;
  background-color: #fff;
}

.logo {
  display: flex;
  align-items: center;
}

.logo img {
  height: 32px;
  margin-right: 10px;
}

.logo h1 {
  font-size: 18px;
  color: #1890ff;
  margin: 0;
  font-weight: 600;
}

.spacer {
  flex: 1;
}

.user-info {
  display: flex;
  align-items: center;
}

.user-dropdown-link {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.username {
  margin: 0 8px;
  color: #333;
}
</style>
