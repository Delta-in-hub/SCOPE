// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/store/auth';

// 导入页面组件
import LoginPage from '@/views/auth/LoginPage.vue';
import RegisterPage from '@/views/auth/RegisterPage.vue';
import NodeListPage from '@/views/node/NodeListPage.vue';
import NotFoundPage from '@/views/NotFoundPage.vue';

const routes = [
  // 无需认证的路由
  {
    path: '/login',
    name: 'Login',
    component: LoginPage,
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterPage,
    meta: { requiresAuth: false }
  },
  
  // 需要认证的路由
  {
    path: '/',
    redirect: '/nodes',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'nodes',
        name: 'NodeList',
        component: NodeListPage,
      },
      // 可以添加更多需要认证的页面
      {
        path: 'system/users',
        name: 'UserManagement',
        component: () => import('@/views/system/UserManagement.vue'),
      },
      {
        path: 'system/settings',
        name: 'SystemSettings',
        component: () => import('@/views/system/SystemSettings.vue'),
      }
    ]
  },
  
  // 404 页面
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFoundPage,
    meta: { requiresAuth: false }
  }
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL), // 使用 History 模式
  routes,
});

// 全局前置守卫 (Navigation Guard)
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore();
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);

  // 检查 Token 是否即将过期，如果需要认证且 Token 快过期了，尝试刷新
  if (requiresAuth && authStore.isAuthenticated && authStore.isTokenExpiringSoon()) {
      try {
          await authStore.refreshAccessToken();
      } catch (error) {
          console.error('Forced logout due to token refresh failure during navigation.');
          // refreshAccessToken 内部失败时会调用 logout 并重定向，这里可以不用再 next('/login')
          // 但为了保险，可以加一个 return，防止继续执行 next()
          return; // 停留在当前逻辑，等待 logout 完成重定向
      }
  }


  if (requiresAuth && !authStore.isAuthenticated) {
    // 如果目标路由需要认证，但用户未认证
    next({ name: 'Login', query: { redirect: to.fullPath } }); // 重定向到登录页，并带上原目标路径
  } else if ((to.name === 'Login' || to.name === 'Register') && authStore.isAuthenticated) {
    // 如果用户已认证，但试图访问登录或注册页，则重定向到主页（例如 NodeList）
    next({ name: 'NodeList' });
  } else {
    // 其他情况（访问公共页，或已认证访问需认证页），正常放行
    next();
  }
});

export default router;