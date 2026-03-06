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
        title: 'routes.basic.template_square',
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
      path: '/templates/detail',
      name: 'templatesDetail',
      component: () => import('../views/explore/templates/detail.vue'),
      meta: {
        title: 'routes.basic.template_detail',
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
