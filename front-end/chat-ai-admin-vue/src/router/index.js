import { createRouter, createWebHashHistory } from 'vue-router'
import BlankLayout from '../layouts/BlankLayout.vue'
import user from './user'
import robot from './robot'
import explore from './explore'
import library from './library'
import authority from './authority'
import database from './database'
import librarySearch from './library-search'
import publicLibrary from './public-library'
import chatMonitor from './chat-monitor'
import noPermission from './no-permission'
import AiExtractFaq from './ai-extract-faq'
import triggerStatics from './trigger-statics'
import plugins from './plugins'

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
  {
    path: '/set_token',
    name: 'SetToken',
    component: () => import('../views/set-token/index.vue'),
    meta: {
      title: '登录',
      noCache: true,
      hidden: true
    }
  },
  user,
  robot,
  explore,
  library,
  librarySearch,
  ...publicLibrary,
  authority,
  noPermission,
  database,
  chatMonitor,
  AiExtractFaq,
  triggerStatics,
  plugins,
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
