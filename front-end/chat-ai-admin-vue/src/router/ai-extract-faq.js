export default {
  path: '/ai-extract-faq',
  name: 'aiExtractFAQ',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/ai-extract-faq/list',
  children: [
    {
      path: '/ai-extract-faq/list',
      name: 'aiExtractFAQList',
      component: () => import('../views/ai-extract-faq/list/index.vue'),
      meta: {
        title: 'AI提取FAQ',
        activeMenu: 'library',
        bgColor: '#fff',
        hideTitle: true
      }
    },
    {
      path: '/ai-extract-faq/details',
      name: 'aiExtractFAQDetails',
      component: () => import('../views/ai-extract-faq/detail/index.vue'),
      meta: {
        title: 'AI提取FAQ',
        activeMenu: 'library',
        isCustomPage: true
      }
    }
  ]
}
