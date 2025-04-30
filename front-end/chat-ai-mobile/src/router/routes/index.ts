
import type { RouteRecordRaw } from 'vue-router';

const routes:Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/home/index.vue')
  },
  {
    path: '/chat',
    name: 'chat',
    component: () => import('@/views/chat/index.vue'),
    children: [
      {
        path: '/chat/pc',
        name: 'chatPc',
        component: () => import('@/views/chat/index.vue')
      },
      {
        path: '/chat/h5',
        name: 'chatH5',
        component: () => import('@/views/chat/index.vue')
      },
    ]
  }
]

if (import.meta.env.DEV) {
  routes.push({
    path: '/icons',
    name: 'icons',
    component: () => import('@/views/icons/index.vue')
  })
}

export default routes
