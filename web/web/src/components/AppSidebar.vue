<template>
  <el-menu
    :default-active="activeRoute"
    class="sidebar-menu"
    :collapse="isCollapse"
    background-color="rgba(0, 21, 41, 0.95)"
    text-color="#ffffff"
    active-text-color="#ffffff"
    router
  >
    <div class="sidebar-logo" @click="toggleCollapse">
      <el-icon v-if="isCollapse"><Expand /></el-icon>
      <el-icon v-else><Fold /></el-icon>
    </div>
    
    <el-menu-item index="/nodes">
      <el-icon><Monitor /></el-icon>
      <template #title>节点管理</template>
    </el-menu-item>
    
    <el-menu-item index="/ebpf-dashboard">
      <el-icon><DataAnalysis /></el-icon>
      <template #title>eBPF面板</template>
    </el-menu-item>
    
    <el-sub-menu index="system">
      <template #title>
        <el-icon><Setting /></el-icon>
        <span>系统管理</span>
      </template>
      <el-menu-item index="/system/users">
        <el-icon><User /></el-icon>
        <span>用户管理</span>
      </el-menu-item>
      <el-menu-item index="/system/settings">
        <el-icon><Tools /></el-icon>
        <span>系统设置</span>
      </el-menu-item>
    </el-sub-menu>
  </el-menu>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';
import { DataAnalysis } from '@element-plus/icons-vue';

const route = useRoute();
const isCollapse = ref(false);

const activeRoute = computed(() => {
  return route.path;
});

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value;
};
</script>

<style scoped>
.sidebar-menu {
  height: 100%;
  border-right: none;
}

.sidebar-menu:not(.el-menu--collapse) {
  width: 220px;
}

.sidebar-logo {
  height: 50px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  color: #fff;
  font-size: 20px;
}

.el-menu-item.is-active {
  background-color: rgba(24, 144, 255, 0.8) !important;
  font-weight: bold;
}

.el-menu-item:hover {
  background-color: rgba(38, 52, 69, 0.8) !important;
}

.el-menu-item, .el-sub-menu__title, .el-sub-menu .el-menu-item {
  color: #ffffff !important;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

.el-sub-menu .el-menu {
  background-color: rgba(0, 10, 20, 0.95) !important;
}

.el-sub-menu .el-menu-item:hover {
  background-color: rgba(38, 52, 69, 0.8) !important;
}

.el-sub-menu .el-menu-item.is-active {
  background-color: rgba(24, 144, 255, 0.8) !important;
}
</style>
