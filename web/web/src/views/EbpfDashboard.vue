<template>
  <div class="ebpf-dashboard-container">
    <h1 class="dashboard-title">eBPF 面板</h1>
    <div class="dashboard-frame-container">
      <iframe 
        :src="dashboardUrl" 
        frameborder="0" 
        allowfullscreen
        class="dashboard-frame"
      ></iframe>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

// Grafana dashboard URL with embedded mode parameters
const dashboardUrl = ref('http://localhost:53000/d/beihhay2q3bpcd/os-event?kiosk&theme=light');

// Ensure iframe has proper access
onMounted(() => {
  // Check if iframe loaded correctly
  const iframe = document.querySelector('.dashboard-frame');
  iframe.onload = () => {
    console.log('Grafana dashboard iframe loaded');
  };
  
  iframe.onerror = (error) => {
    console.error('Error loading Grafana dashboard:', error);
  };
});
</script>

<style scoped>
.ebpf-dashboard-container {
  height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
  padding: 0;
  margin: 0;
}

.dashboard-title {
  margin-bottom: 16px;
  font-size: 24px;
  font-weight: 500;
  color: #1f2f3d;
}

.dashboard-frame-container {
  flex: 1;
  height: calc(100vh - 180px);
  min-height: 600px;
  overflow: hidden;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  background-color: #fff;
  position: relative;
}

.dashboard-frame {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border: none;
}
</style>
