import BlankLayout from '@/layouts/BlankLayout.vue'

const routes = [
  {
    path: '/',
    name: 'Root',
    component: BlankLayout,
    redirect: '/chat'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: '登录',
      noCache: true,
      hidden: true
    }
  },
  {
    path: '/chat',
    name: 'chat',
    component: () => import('@/views/chat/index.vue'),
    meta: {
      title: '登录',
      noCache: true,
      hidden: true
    }
  }
]

if (import.meta.env.DEV) {
  routes.push({
    path: '/icons',
    name: 'icons',
    component: () => import('@/views/icons/index.vue'),
    meta: {
      title: 'icons',
      noCache: true,
      hidden: true
    }
  })
}

export default routes
