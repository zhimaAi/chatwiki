import 'virtual:svg-icons-register'
import './assets/main.css'
import { createApp } from 'vue'
// 初始化多语言
import { setupI18n } from '@/locales'
// 引入状态管理
import { setupStore } from '@/stores'
import { setupRouter } from './router/index'

import App from './App.vue'
import { message, Modal } from 'ant-design-vue'
import { registGlobalComponent } from './components'

// 创建实例
const setupAll = async () => {
  const app = createApp(App)

  app.config.globalProperties.$message = message
  app.config.globalProperties.$modal = Modal

  setupStore(app)

  await setupI18n(app)

  await setupRouter(app)

  registGlobalComponent(app)

  app.mount('#app')
}

setupAll()
