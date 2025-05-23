export default {
  path: '/no-permission',
  name: 'noPermission',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/no-permission/index',
  children: [
    {
      path: '/no-permission/index',
      component: () => import('../views/no-permission/index.vue'),
      meta: {
        title: '',
        isCustomPage: true
      }
    },
  ]
}
