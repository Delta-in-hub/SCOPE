<template>
  <component :is="layout">
    <router-view />
  </component>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import DefaultLayout from '@/layouts/DefaultLayout.vue';
import AuthLayout from '@/layouts/AuthLayout.vue';
import { useAuthStore } from '@/store/auth';

const authStore = useAuthStore();
const route = useRoute();

// 根据路由元信息动态选择布局
const layout = computed(() => {
  // 检查路由是否需要认证
  const requiresAuth = route.matched.some(record => record.meta.requiresAuth);
  
  // 如果路由不需要认证，使用 AuthLayout
  if (!requiresAuth) {
    return AuthLayout;
  }
  
  // 否则使用 DefaultLayout
  return DefaultLayout;
});
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  height: 100%;
}

html, body {
  height: 100%;
  margin: 0;
  padding: 0;
}
</style>