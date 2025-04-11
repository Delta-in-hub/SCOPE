<template>
  <div class="node-list-container">
    <el-card class="node-card">
      <template #header>
        <div class="card-header">
          <span class="header-title">节点列表</span>
          <el-button class="refresh-button" type="primary" @click="fetchNodes" :loading="loading" round>
            <el-icon><Refresh /></el-icon>刷新
          </el-button>
        </div>
      </template>

      <el-table v-loading="loading" :data="nodes" style="width: 100%" border stripe>
        <el-table-column prop="id" label="节点 ID" min-width="180" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'online' ? 'success' : 'danger'">
              {{ row.status === 'online' ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
         <el-table-column label="IP 地址" min-width="200">
            <template #default="{ row }">
                <div v-if="row.ips && Object.keys(row.ips).length > 0">
                    <div v-for="(ip, iface) in row.ips" :key="iface">
                        {{ iface }}: {{ ip }}
                    </div>
                </div>
                <span v-else>N/A</span>
            </template>
        </el-table-column>
        <el-table-column prop="last_seen" label="最后在线时间" width="180">
           <template #default="{ row }">
               {{ formatDateTime(row.last_seen) }}
           </template>
        </el-table-column>
        <el-table-column prop="latency" label="通信延迟" width="120">
            <template #default="{ row }">
                <el-tag :type="getLatencyTagType(row.latency)" size="small" effect="light">
                    {{ formatDuration(row.latency) }}
                </el-tag>
            </template>
        </el-table-column>
      </el-table>

       <el-alert v-if="error" :title="error" type="error" show-icon :closable="false" style="margin-top: 15px;" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useNodeStore } from '@/store/node'; // 假设你创建了 node store
import { ElMessage } from 'element-plus';
// import { Refresh } from '@element-plus/icons-vue'; // 如果在 main.js 全局注册了，这里不用单独引入

const nodeStore = useNodeStore();
const loading = ref(false);
const error = ref(null);

// 从 store 获取节点数据
const nodes = computed(() => nodeStore.nodes);

// 获取节点列表的方法
const fetchNodes = async () => {
  loading.value = true;
  error.value = null;
  try {
    await nodeStore.fetchNodes();
  } catch (err) {
    error.value = err.response?.data?.message || err.response?.data || '获取节点列表失败';
    ElMessage.error(error.value);
    console.error('Failed to fetch nodes:', err);
  } finally {
    loading.value = false;
  }
};

// 组件挂载时获取数据
onMounted(() => {
  fetchNodes();
});

// 格式化日期时间 (可以放到 utils 文件中)
const formatDateTime = (dateTimeString) => {
  if (!dateTimeString) return 'N/A';
  try {
    const date = new Date(dateTimeString);
    // 检查日期是否有效
    if (isNaN(date.getTime())) {
        return 'Invalid Date';
    }
    return date.toLocaleString(); // 或者使用更专业的日期格式化库如 date-fns 或 dayjs
  } catch (e) {
    return 'Invalid Date Format';
  }
};

// 格式化Go的time.Duration类型
const formatDuration = (duration) => {
  if (!duration) return 'N/A';
  
  // Go的time.Duration是纳秒为单位的字符串
  // 例如："1.234s"、"4.56ms"、"789µs"、"12ns"
  
  // 如果已经是格式化的字符串，直接返回
  if (typeof duration === 'string') {
    // 将µs替换为ms，使其更易读
    if (duration.includes('µs')) {
      const microseconds = parseFloat(duration.replace('µs', ''));
      return (microseconds / 1000).toFixed(2) + 'ms';
    }
    return duration;
  }
  
  // 如果是数字（纳秒），则转换为合适的单位
  const ns = Number(duration);
  if (isNaN(ns)) return 'N/A';
  
  if (ns < 1000) {
    return ns + 'ns';
  } else if (ns < 1000000) {
    return (ns / 1000).toFixed(2) + 'µs';
  } else if (ns < 1000000000) {
    return (ns / 1000000).toFixed(2) + 'ms';
  } else {
    return (ns / 1000000000).toFixed(2) + 's';
  }
};

// 根据延迟时间获取标签类型
const getLatencyTagType = (duration) => {
  if (!duration) return 'info';
  
  let ms = 0;
  
  // 处理字符串格式
  if (typeof duration === 'string') {
    if (duration.includes('ns')) {
      ms = parseFloat(duration) / 1000000;
    } else if (duration.includes('µs')) {
      ms = parseFloat(duration) / 1000;
    } else if (duration.includes('ms')) {
      ms = parseFloat(duration);
    } else if (duration.includes('s') && !duration.includes('ms') && !duration.includes('µs') && !duration.includes('ns')) {
      ms = parseFloat(duration) * 1000;
    }
  } else {
    // 处理数字格式（纳秒）
    ms = Number(duration) / 1000000;
  }
  
  // 根据延迟时间返回不同的标签类型
  if (ms < 10) return 'success'; // 小于10ms，非常好
  if (ms < 50) return ''; // 小于50ms，正常
  if (ms < 100) return 'warning'; // 小于100ms，警告
  return 'danger'; // 大于等于100ms，危险
};

</script>

<style scoped>
.node-list-container {
  padding: 20px;
  height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
}

.node-card {
  flex: 1;
  margin-bottom: 20px;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  background-color: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: #1d1d1f;
}

.refresh-button {
  margin-left: 10px;
  font-size: 14px;
}

:deep(.el-table) {
  background-color: transparent;
  border-radius: 8px;
  overflow: hidden;
}

:deep(.el-table__header) {
  background-color: rgba(245, 247, 250, 0.8);
}

:deep(.el-table__row) {
  background-color: rgba(255, 255, 255, 0.6);
}

:deep(.el-table__row:hover) {
  background-color: rgba(245, 247, 250, 0.9) !important;
}

:deep(.el-table--striped .el-table__row.striped) {
  background-color: rgba(250, 250, 252, 0.8);
}

:deep(.el-tag) {
  border-radius: 12px;
  padding: 2px 10px;
}
</style>