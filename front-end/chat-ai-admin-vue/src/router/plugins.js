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
        title: 'routes.basic.plugin_management',
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
        title: 'routes.basic.plugin_detail',
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
