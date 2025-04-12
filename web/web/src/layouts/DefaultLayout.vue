<template>
  <div class="app-container">
    <div class="background-elements">
      <div class="bg-circle circle-1"></div>
      <div class="bg-circle circle-2"></div>
      <div class="bg-circle circle-3"></div>
    </div>
    
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
  position: relative;
  font-family: 'SF Pro Display', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  overflow: hidden;
}

/* 背景元素 */
.background-elements {
  position: fixed;
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
  filter: blur(100px);
}

.circle-1 {
  width: 70vw;
  height: 70vw;
  background: rgba(0, 122, 255, 0.05);
  top: -10vw;
  right: -10vw;
}

.circle-2 {
  width: 70vw;
  height: 70vw;
  background: rgba(88, 86, 214, 0.05);
  bottom: -10vw;
  left: -10vw;
}

.circle-3 {
  width: 30vw;
  height: 30vw;
  background: rgba(52, 199, 89, 0.05);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.el-header {
  padding: 0;
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  z-index: 1000;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.el-aside {
  background-color: rgba(0, 21, 41, 0.9);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  transition: width 0.3s;
  overflow-x: hidden;
  border-right: 1px solid rgba(255, 255, 255, 0.05);
}

.el-main {
  background-color: rgba(245, 247, 250, 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  padding: 20px;
  overflow-y: auto;
  border-radius: 10px;
  margin: 10px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
}
</style>
