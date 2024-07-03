export default {
  path: '/authority',
  name: 'authority',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/authority/index',
  children: [
    {
      path: '/authority/index',
      name: 'authorityIndex',
      component: () => import('../views/authority/index.vue'),
      meta: {
        title: '',
        isCustomPage: true
      }
    },
  ]
}
