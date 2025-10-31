export default {
  path: '/trigger-statics',
  name: 'triggerStatics',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/trigger-statics/list',
  children: [
    {
      path: '/trigger-statics/list',
      name: 'triggerStaticsList',
      component: () => import('../views/trigger-statics/list/index.vue'),
      meta: {
        title: '触发次数统计',
        activeMenu: 'library',
        bgColor: '#fff',
        hideTitle: true
      }
    },
  ]
}
