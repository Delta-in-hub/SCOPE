<template>
  <div>
    <el-card>
      <template #header>
        <span>节点列表</span>
         <el-button style="float: right; padding: 3px 0" type="primary" link @click="fetchNodes" :loading="loading">
            <el-icon><Refresh /></el-icon>刷新
         </el-button>
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
        <!-- 可以添加更多列，例如 latency (需要格式化 time.Duration) -->
        <!--
        <el-table-column prop="latency" label="延迟" width="100">
            <template #default="{ row }">
                {{ formatDuration(row.latency) }}
            </template>
        </el-table-column>
        -->
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

// 格式化 Go 的 time.Duration (纳秒) 为更易读的格式 (可选)
// const formatDuration = (ns) => {
//   if (ns === undefined || ns === null) return 'N/A';
//   if (ns < 1000) return `${ns} ns`;
//   if (ns < 1000000) return `${(ns / 1000).toFixed(1)} µs`;
//   if (ns < 1000000000) return `${(ns / 1000000).toFixed(1)} ms`;
//   return `${(ns / 1000000000).toFixed(2)} s`;
// };

</script>

<style scoped>
/* 可以添加一些页面特定样式 */
.el-card {
    margin: 20px;
}
</style>