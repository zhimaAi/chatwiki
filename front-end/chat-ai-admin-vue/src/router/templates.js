export default {
  path: '/templates',
  name: 'templates',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/templates/index',
  children: [
    {
      path: '/templates/index',
      name: 'templatesIndex',
      component: () => import('../views/explore/templates/index.vue'),
      meta: {
        title: '模板广场',
        hideTitle: true,
        activeMenu: 'explore',
        bgColor: '#ffffff',
        pageStyle: {
          'padding-right': 0,
          overflow: 'hidden'
        }
      }
    },
  ]
}
