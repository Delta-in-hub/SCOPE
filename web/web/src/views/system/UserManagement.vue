<template>
  <div class="user-management">
    <el-card class="user-card">
      <template #header>
        <div class="card-header">
          <span class="header-title">用户管理</span>
          <el-button type="primary" @click="openAddUserDialog" class="add-button" round>
            <el-icon><Plus /></el-icon>添加用户
          </el-button>
        </div>
      </template>
      
      <el-table :data="userList" style="width: 100%" v-loading="loading" border stripe>
        <el-table-column prop="user_id" label="用户ID" width="220" />
        <el-table-column prop="display_name" label="用户名" width="150" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editUser(row)">
              <el-icon><Edit /></el-icon>
            </el-button>
            <el-button type="danger" link @click="confirmDeleteUser(row)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页组件 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 添加/编辑用户对话框 -->
    <el-dialog
      v-model="userDialogVisible"
      :title="isEdit ? '编辑用户' : '添加用户'"
      width="500px"
    >
      <el-form
        ref="userFormRef"
        :model="userForm"
        :rules="userRules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="display_name">
          <el-input v-model="userForm.display_name" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="userForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="userDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitUserForm" :loading="submitting">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';

// 模拟用户数据
const userList = ref([
  {
    user_id: 'usr_001',
    display_name: '管理员',
    email: 'admin@example.com',
    created_at: '2025-04-01T10:00:00Z'
  },
  {
    user_id: 'usr_002',
    display_name: '测试用户',
    email: 'test@example.com',
    created_at: '2025-04-05T14:30:00Z'
  }
]);

const loading = ref(false);
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(2);
const userDialogVisible = ref(false);
const isEdit = ref(false);
const submitting = ref(false);
const userFormRef = ref(null);

// 用户表单数据
const userForm = reactive({
  user_id: '',
  display_name: '',
  email: '',
  password: ''
});

// 表单验证规则
const userRules = {
  display_name: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度应在2-20个字符之间', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: ['blur', 'change'] }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
  ]
};

// 格式化日期时间
const formatDateTime = (dateTimeString) => {
  if (!dateTimeString) return 'N/A';
  try {
    const date = new Date(dateTimeString);
    if (isNaN(date.getTime())) {
      return 'Invalid Date';
    }
    return date.toLocaleString();
  } catch (e) {
    return 'Invalid Date Format';
  }
};

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size;
  fetchUsers();
};

const handleCurrentChange = (page) => {
  currentPage.value = page;
  fetchUsers();
};

// 获取用户列表（模拟）
const fetchUsers = () => {
  loading.value = true;
  // 这里应该调用实际的API
  setTimeout(() => {
    loading.value = false;
  }, 500);
};

// 打开添加用户对话框
const openAddUserDialog = () => {
  isEdit.value = false;
  resetUserForm();
  userDialogVisible.value = true;
};

// 编辑用户
const editUser = (row) => {
  isEdit.value = true;
  resetUserForm();
  Object.assign(userForm, {
    user_id: row.user_id,
    display_name: row.display_name,
    email: row.email
  });
  userDialogVisible.value = true;
};

// 确认删除用户
const confirmDeleteUser = (row) => {
  ElMessageBox.confirm(
    `确定要删除用户 "${row.display_name}" 吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    deleteUser(row.user_id);
  }).catch(() => {
    // 用户取消操作
  });
};

// 删除用户（模拟）
const deleteUser = (userId) => {
  loading.value = true;
  // 这里应该调用实际的API
  setTimeout(() => {
    userList.value = userList.value.filter(user => user.user_id !== userId);
    ElMessage.success('删除成功');
    loading.value = false;
  }, 500);
};

// 重置用户表单
const resetUserForm = () => {
  if (userFormRef.value) {
    userFormRef.value.resetFields();
  }
  userForm.user_id = '';
  userForm.display_name = '';
  userForm.email = '';
  userForm.password = '';
};

// 提交用户表单
const submitUserForm = async () => {
  if (!userFormRef.value) return;
  
  await userFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true;
      try {
        if (isEdit.value) {
          // 编辑用户（模拟）
          setTimeout(() => {
            const index = userList.value.findIndex(user => user.user_id === userForm.user_id);
            if (index !== -1) {
              userList.value[index].display_name = userForm.display_name;
              userList.value[index].email = userForm.email;
            }
            ElMessage.success('更新成功');
            userDialogVisible.value = false;
            submitting.value = false;
          }, 500);
        } else {
          // 添加用户（模拟）
          setTimeout(() => {
            const newUser = {
              user_id: 'usr_' + Math.floor(Math.random() * 1000).toString().padStart(3, '0'),
              display_name: userForm.display_name,
              email: userForm.email,
              created_at: new Date().toISOString()
            };
            userList.value.push(newUser);
            total.value++;
            ElMessage.success('添加成功');
            userDialogVisible.value = false;
            submitting.value = false;
          }, 500);
        }
      } catch (error) {
        ElMessage.error('操作失败: ' + (error.message || '未知错误'));
        submitting.value = false;
      }
    }
  });
};

onMounted(() => {
  fetchUsers();
});
</script>

<style scoped>
.user-management {
  padding: 20px;
  height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
}

.user-card {
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

.add-button {
  margin-left: 10px;
  font-size: 14px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

:deep(.el-table) {
  background-color: transparent;
  border-radius: 8px;
  overflow: hidden;
  flex: 1;
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

:deep(.el-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.el-dialog__header) {
  background-color: rgba(245, 247, 250, 0.9);
  padding: 15px 20px;
  margin-right: 0;
}

:deep(.el-dialog__body) {
  padding: 20px;
}

:deep(.el-dialog__footer) {
  padding: 10px 20px 20px;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}
</style>
