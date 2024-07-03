import './assets/main.css'
import 'virtual:svg-icons-register'

import { createApp } from 'vue'
import { message, Modal } from 'ant-design-vue'

// 初始化多语言
import { setupI18n } from '@/locales'

// 引入状态管理
import { setupStore } from '@/stores'

import router from './router'

import App from './App.vue'

import './permission'

import registGlobaDirective from './directives'
import { registGlobalComponent } from './components'

// 创建实例
const setupAll = async () => {
  const app = createApp(App)

  app.config.globalProperties.$message = message
  app.config.globalProperties.$modal = Modal

  await setupI18n(app)

  setupStore(app)

  app.use(router)

  registGlobaDirective(app)
  registGlobalComponent(app)
  app.mount('#app')
}

setupAll()
