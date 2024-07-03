export default {
  path: '/library',
  name: 'Library',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/library/list',
  children: [
    {
      path: '/library/list',
      name: 'libraryList',
      component: () => import('../views/library/library-list/library-list.vue'),
      meta: {
        title: '知识库管理',
        activeMenu: '/library',
        bgColor: '#F5F9FF'
      }
    },
    {
      path: '/library/add',
      name: 'addLibrary',
      component: () => import('../views/library/add-library/add-library.vue'),
      meta: {
        title: '新建知识库',
        activeMenu: '/library',
        breadcrumb: [
          {
            title: '知识库管理',
            path: '/library/list'
          },
          {
            title: '新建知识库',
            path: '/library/add'
          }
        ]
      }
    },
    {
      path: '/library/document-segmentation',
      name: 'documentSegmentation',
      component: () => import('../views/library/document-segmentation/document-segmentation.vue'),
      meta: {
        title: '新建知识库',
        activeMenu: '/library',
        breadcrumb: [
          {
            title: '知识库管理',
            path: '/library/list'
          },
          {
            title: '新建知识库',
            path: '/library/add'
          }
        ]
      }
    },
    {
      path: '/library/details',
      name: 'libraryDetails',
      component: () => import('../views/library/library-details/index.vue'),
      meta: {
        title: '编辑知识库',
        activeMenu: '/library',
        isCustomPage: true,
      },
      redirect: '/library/details/knowledge-document',
      children: [
        {
          path: '/library/details/knowledge-document',
          name: 'knowledgeDocument',
          component: () => import('../views/library/library-details/knowledge-document.vue'),
          meta: {
            title: '知识库文档',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/recall-testing',
          name: 'recallTesting',
          component: () => import('@/views/library/library-details/recall-testing.vue'),
          meta: {
            title: '召回测试',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/knowledge-config',
          name: 'knowledgeConfig',
          component: () => import('@/views/library/library-details/knowledge-config.vue'),
          meta: {
            title: '知识库配置',
            activeMenu: '/library'
          }
        },
      ],
    },
    {
      path: '/library/preview',
      name: 'libraryPreview',
      component: () => import('../views/library/library-preview/library-preview.vue'),
      meta: {
        title: '知识库管理',
        activeMenu: '/library',
        breadcrumb: [
          {
            title: '知识库管理',
            path: '/library/list'
          },
          {
            title: '知识库详情',
            path: '/library/preview'
          }
        ]
      }
    }
  ]
}
