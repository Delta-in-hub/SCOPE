<template>
  <div class="system-settings">
    <el-card class="settings-card">
      <template #header>
        <div class="card-header">
          <span class="header-title">系统设置</span>
        </div>
      </template>
      
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本设置" name="basic">
          <el-form
            ref="basicFormRef"
            :model="basicForm"
            :rules="basicRules"
            label-width="120px"
            class="settings-form"
          >
            <el-form-item label="系统名称" prop="systemName">
              <el-input v-model="basicForm.systemName" />
            </el-form-item>
            
            <el-form-item label="API 地址" prop="apiUrl">
              <el-input v-model="basicForm.apiUrl" />
            </el-form-item>
            
            <el-form-item label="日志级别" prop="logLevel">
              <el-select v-model="basicForm.logLevel" style="width: 100%">
                <el-option label="调试" value="debug" />
                <el-option label="信息" value="info" />
                <el-option label="警告" value="warn" />
                <el-option label="错误" value="error" />
              </el-select>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="saveBasicSettings" :loading="saving">保存设置</el-button>
              <el-button @click="resetBasicForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="安全设置" name="security">
          <el-form
            ref="securityFormRef"
            :model="securityForm"
            :rules="securityRules"
            label-width="120px"
            class="settings-form"
          >
            <el-form-item label="会话超时(分钟)" prop="sessionTimeout">
              <el-input-number v-model="securityForm.sessionTimeout" :min="5" :max="1440" />
            </el-form-item>
            
            <el-form-item label="密码最小长度" prop="minPasswordLength">
              <el-input-number v-model="securityForm.minPasswordLength" :min="6" :max="32" />
            </el-form-item>
            
            <el-form-item label="密码复杂度" prop="passwordComplexity">
              <el-select v-model="securityForm.passwordComplexity" style="width: 100%">
                <el-option label="低 (仅字母和数字)" value="low" />
                <el-option label="中 (字母、数字和特殊字符)" value="medium" />
                <el-option label="高 (大小写字母、数字和特殊字符)" value="high" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="启用双因素认证" prop="enableTwoFactor">
              <el-switch v-model="securityForm.enableTwoFactor" />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="saveSecuritySettings" :loading="saving">保存设置</el-button>
              <el-button @click="resetSecurityForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="节点设置" name="node">
          <el-form
            ref="nodeFormRef"
            :model="nodeForm"
            :rules="nodeRules"
            label-width="120px"
            class="settings-form"
          >
            <el-form-item label="节点心跳间隔(秒)" prop="heartbeatInterval">
              <el-input-number v-model="nodeForm.heartbeatInterval" :min="5" :max="300" />
            </el-form-item>
            
            <el-form-item label="节点超时时间(秒)" prop="nodeTimeout">
              <el-input-number v-model="nodeForm.nodeTimeout" :min="10" :max="600" />
            </el-form-item>
            
            <el-form-item label="自动清理离线节点" prop="autoCleanOfflineNodes">
              <el-switch v-model="nodeForm.autoCleanOfflineNodes" />
            </el-form-item>
            
            <el-form-item label="离线清理时间(天)" prop="offlineCleanupDays" v-if="nodeForm.autoCleanOfflineNodes">
              <el-input-number v-model="nodeForm.offlineCleanupDays" :min="1" :max="365" />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="saveNodeSettings" :loading="saving">保存设置</el-button>
              <el-button @click="resetNodeForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { ElMessage } from 'element-plus';

const activeTab = ref('basic');
const saving = ref(false);

// 表单引用
const basicFormRef = ref(null);
const securityFormRef = ref(null);
const nodeFormRef = ref(null);

// 基本设置表单
const basicForm = reactive({
  systemName: 'Scope Center',
  apiUrl: 'http://127.0.0.1:18080/api/v1',
  logLevel: 'info',
});

// 安全设置表单
const securityForm = reactive({
  sessionTimeout: 30,
  minPasswordLength: 8,
  passwordComplexity: 'medium',
  enableTwoFactor: false,
});

// 节点设置表单
const nodeForm = reactive({
  heartbeatInterval: 30,
  nodeTimeout: 90,
  autoCleanOfflineNodes: true,
  offlineCleanupDays: 7,
});

// 表单验证规则
const basicRules = {
  systemName: [
    { required: true, message: '请输入系统名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度应在2-50个字符之间', trigger: 'blur' }
  ],
  apiUrl: [
    { required: true, message: '请输入API地址', trigger: 'blur' },
    { pattern: /^https?:\/\/.+/i, message: '请输入有效的URL地址', trigger: 'blur' }
  ],
  logLevel: [
    { required: true, message: '请选择日志级别', trigger: 'change' }
  ]
};

const securityRules = {
  sessionTimeout: [
    { required: true, message: '请输入会话超时时间', trigger: 'blur' },
    { type: 'number', message: '会话超时必须为数字', trigger: 'blur' }
  ],
  minPasswordLength: [
    { required: true, message: '请输入密码最小长度', trigger: 'blur' },
    { type: 'number', message: '密码长度必须为数字', trigger: 'blur' }
  ],
  passwordComplexity: [
    { required: true, message: '请选择密码复杂度', trigger: 'change' }
  ]
};

const nodeRules = {
  heartbeatInterval: [
    { required: true, message: '请输入心跳间隔', trigger: 'blur' },
    { type: 'number', message: '心跳间隔必须为数字', trigger: 'blur' }
  ],
  nodeTimeout: [
    { required: true, message: '请输入节点超时时间', trigger: 'blur' },
    { type: 'number', message: '节点超时时间必须为数字', trigger: 'blur' }
  ],
  offlineCleanupDays: [
    { required: true, message: '请输入离线清理时间', trigger: 'blur' },
    { type: 'number', message: '离线清理时间必须为数字', trigger: 'blur' }
  ]
};

// 保存基本设置
const saveBasicSettings = async () => {
  if (!basicFormRef.value) return;
  
  await basicFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true;
      try {
        // 这里应该调用实际的API保存设置
        await simulateApiCall();
        ElMessage.success('基本设置保存成功');
      } catch (error) {
        ElMessage.error('保存失败: ' + (error.message || '未知错误'));
      } finally {
        saving.value = false;
      }
    }
  });
};

// 保存安全设置
const saveSecuritySettings = async () => {
  if (!securityFormRef.value) return;
  
  await securityFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true;
      try {
        // 这里应该调用实际的API保存设置
        await simulateApiCall();
        ElMessage.success('安全设置保存成功');
      } catch (error) {
        ElMessage.error('保存失败: ' + (error.message || '未知错误'));
      } finally {
        saving.value = false;
      }
    }
  });
};

// 保存节点设置
const saveNodeSettings = async () => {
  if (!nodeFormRef.value) return;
  
  await nodeFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true;
      try {
        // 这里应该调用实际的API保存设置
        await simulateApiCall();
        ElMessage.success('节点设置保存成功');
      } catch (error) {
        ElMessage.error('保存失败: ' + (error.message || '未知错误'));
      } finally {
        saving.value = false;
      }
    }
  });
};

// 模拟API调用
const simulateApiCall = () => {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve();
    }, 800);
  });
};

// 重置表单
const resetBasicForm = () => {
  if (basicFormRef.value) {
    basicFormRef.value.resetFields();
  }
};

const resetSecurityForm = () => {
  if (securityFormRef.value) {
    securityFormRef.value.resetFields();
  }
};

const resetNodeForm = () => {
  if (nodeFormRef.value) {
    nodeFormRef.value.resetFields();
  }
};

// 加载设置（模拟）
const loadSettings = async () => {
  try {
    // 这里应该调用实际的API加载设置
    // 模拟加载完成后的操作
    console.log('Settings loaded');
  } catch (error) {
    ElMessage.error('加载设置失败: ' + (error.message || '未知错误'));
  }
};

onMounted(() => {
  loadSettings();
});
</script>

<style scoped>
.system-settings {
  padding: 20px;
  height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
}

.settings-card {
  flex: 1;
  margin-bottom: 20px;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  background-color: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  display: flex;
  flex-direction: column;
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

.settings-form {
  max-width: 600px;
  margin-top: 20px;
  background-color: rgba(255, 255, 255, 0.7);
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.el-tabs {
  margin-top: 10px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

:deep(.el-tabs__content) {
  flex: 1;
  overflow: auto;
  padding: 10px 0;
}

:deep(.el-tabs__header) {
  background-color: rgba(245, 247, 250, 0.8);
  border-radius: 8px 8px 0 0;
  margin-bottom: 15px;
  padding: 5px 15px 0;
}

:deep(.el-tabs__nav) {
  border: none !important;
}

:deep(.el-tabs__item) {
  height: 40px;
  line-height: 40px;
  font-size: 15px;
  color: #606266;
}

:deep(.el-tabs__item.is-active) {
  color: #0071e3;
  font-weight: 500;
}

:deep(.el-tabs__active-bar) {
  background-color: #0071e3;
  height: 3px;
  border-radius: 3px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #303133;
}

:deep(.el-button--primary) {
  background: linear-gradient(90deg, #0071e3, #42a5f5);
  border: none;
  transition: all 0.3s ease;
}

:deep(.el-button--primary:hover) {
  background: linear-gradient(90deg, #005bb5, #3994e4);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 113, 227, 0.3);
}
</style>
