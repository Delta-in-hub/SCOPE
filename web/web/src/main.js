// src/main.js
import { createApp } from 'vue'
import App from './App.vue'
import router from './router' // 稍后创建
import { createPinia } from 'pinia' // 稍后创建 store
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css' // 引入 Element Plus 样式
import * as ElementPlusIconsVue from '@element-plus/icons-vue' // 引入图标
import './assets/styles/index.css' // 引入你的全局样式

const app = createApp(App)

// 注册所有 Element Plus 图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(createPinia()) // 使用 Pinia
app.use(router)       // 使用 Vue Router
app.use(ElementPlus)  // 使用 Element Plus

app.mount('#app')