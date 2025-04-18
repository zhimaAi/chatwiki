export default {
  path: '/library-search',
  name: 'LibrarySearchLayout',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/library-search/index',
  children: [
    {
      path: '/library-search/index',
      name: 'library-search',
      component: () => import('@/views/library-search/index.vue'),
      meta: {
        title: '搜索',
        icon: 'monitor',
        hideTitle: true,
        isCustomPage: true
      }
    }
  ]
}
