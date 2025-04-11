export default {
  path: '/chat-monitor',
  name: 'ChatMonitorLayout',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/chat-monitor/index',
  children: [
    {
      path: '/chat-monitor/index',
      name: 'chat-monitor',
      component: () => import('@/views/chat-monitor/index.vue'),
      meta: {
        title: '实时会话',
        icon: 'monitor',
        hideTitle: true,
        isCustomPage: true
      }
    }
  ]
}
