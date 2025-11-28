export default {
  path: '/plugins',
  name: 'plugins',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/plugins/index',
  children: [
    {
      path: '/plugins/index',
      name: 'pluginsIndex',
      component: () => import('../views/explore/plugins/index.vue'),
      meta: {
        title: '插件管理',
        hideTitle: true,
        activeMenu: 'explore',
        bgColor: '#ffffff',
        pageStyle: {
          'padding-right': 0,
          overflow: 'hidden'
        }
      }
    },
    {
      path: '/plugins/detail',
      name: 'pluginsDetail',
      component: () => import('../views/explore/plugins/detail.vue'),
      meta: {
        title: '插件详情',
        hideTitle: true,
        activeMenu: 'explore',
        bgColor: '#ffffff',
        pageStyle: {
          'padding': 0,
        }
      }
    },
  ]
}
