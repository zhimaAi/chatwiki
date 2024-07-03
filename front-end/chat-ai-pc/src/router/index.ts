import type { App } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import { createRouterGuards } from './router-guards'
import routes from './routes/index'

const router = createRouter({
  history: createWebHashHistory(),
  routes: routes
})

export async function setupRouter(app: App) {
  // 创建路由守卫
  createRouterGuards(router)

  app.use(router)

  // 路由准备就绪后挂载APP实例
  await router.isReady()
}
export default router
