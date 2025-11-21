export default {
  path: '/explore',
  name: 'Explore',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/explore/index',
  children: [
    {
      path: '/explore/index',
      name: 'exploreIndex',
      component: () => import('../views/explore/explore-index/explore-index.vue'),
      meta: {
        title: '探索',
        hideTitle: true,
        activeMenu: '/explore',
        bgColor: '#ffffff',
        pageStyle: {
          'padding-right': 0,
          overflow: 'hidden'
        }
      }
    }
  ]
}
