<template>
  <div class="app-container">
    <el-container>
      <el-header height="60px">
        <app-header />
      </el-header>
      <el-container>
        <el-aside width="220px">
          <app-sidebar />
        </el-aside>
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import AppHeader from '@/components/AppHeader.vue';
import AppSidebar from '@/components/AppSidebar.vue';
</script>

<style scoped>
.app-container {
  height: 100vh;
  width: 100%;
}

.el-header {
  padding: 0;
  background-color: #ffffff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.el-aside {
  background-color: #001529;
  transition: width 0.3s;
  overflow-x: hidden;
}

.el-main {
  background-color: #f5f7fa;
  padding: 20px;
  overflow-y: auto;
}
</style>
