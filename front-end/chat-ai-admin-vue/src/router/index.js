import { createRouter, createWebHashHistory } from 'vue-router'
import BlankLayout from '../layouts/BlankLayout.vue'
import user from './user'
import robot from './robot'
import library from './library'
import authority from './authority'

const routes = [
  {
    path: '/',
    name: 'Root',
    component: BlankLayout,
    redirect: '/robot/list'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/login.vue'),
    meta: {
      title: '登录',
      noCache: true,
      hidden: true
    }
  },
  user,
  robot,
  library,
  authority
]

if (import.meta.env.DEV) {
  routes.push({
    path: '/icons',
    name: 'icons',
    component: () => import('../views/icons/index.vue'),
    meta: {
      title: 'icons',
      noCache: true,
      hidden: true
    }
  })
}

const router = createRouter({
  history: createWebHashHistory(),
  routes: routes
})

export default router
