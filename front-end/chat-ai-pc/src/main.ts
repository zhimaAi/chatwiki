import './assets/main.css'
import 'virtual:svg-icons-register'
import { createApp } from 'vue'
import { setupVant } from './vant'

// 初始化多语言
import { setupI18n } from '@/locales'

// 引入状态管理
import { setupStore } from '@/stores'

import App from './App.vue'
import { setupRouter } from './router/index'

import registGlobaDirective from './directives'
import { registGlobalComponent } from './components'

import './event/index'
// 创建实例
const setupAll = async () => {
  const app = createApp(App)

  await setupI18n(app)

  setupStore(app)

  setupRouter(app)

  setupVant(app)

  registGlobaDirective(app)
  registGlobalComponent(app)

  app.mount('#app')
}

setupAll()
