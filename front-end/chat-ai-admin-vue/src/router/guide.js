export default {
  path: '/guide',
  name: 'guide',
  component: () => import('../layouts/AdminLayout/index.vue'),
  children: [
    {
      path: '/guide',
      name: 'guideIndex',
      component: () => import('../views/guide/index.vue'),
      meta: {
        title: '新手指引',
        hideTitle: true,
        activeMenu: 'guide',
        bgColor: '#F2F4F7',
        pageStyle: {
          'padding-right': 0,
          overflow: 'hidden'
        }
      }
    }
  ]
}
