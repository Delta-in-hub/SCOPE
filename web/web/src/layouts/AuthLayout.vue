<template>
  <div class="auth-layout">
    <div class="auth-header">
      <div class="logo">
        <img src="@/assets/logo.svg" alt="Logo" />
        <h1>Scope Center</h1>
      </div>
    </div>
    <div class="auth-content">
      <router-view />
    </div>
    <div class="auth-footer">
      <p>© {{ currentYear }} Scope Center. All rights reserved.</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const currentYear = computed(() => new Date().getFullYear());
</script>

<style scoped>
.auth-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f0f2f5;
  background-image: radial-gradient(circle at 50% 50%, rgba(0, 0, 0, 0.02) 0%, rgba(0, 0, 0, 0.05) 100%);
}

.auth-header {
  padding: 20px 0;
  text-align: center;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}

.logo img {
  height: 40px;
  margin-right: 10px;
}

.logo h1 {
  font-size: 24px;
  color: #1890ff;
  margin: 0;
}

.auth-content {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.auth-footer {
  text-align: center;
  padding: 20px 0;
  color: #999;
  font-size: 14px;
}
</style>
