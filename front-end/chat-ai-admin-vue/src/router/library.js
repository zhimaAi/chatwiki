import { useLibraryStore } from '@/stores/modules/library'

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
        hideTitle: true,
        pageStyle:{
          'padding-right': 0,
          overflow: 'hidden'
        }
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
        hideTitle: true
      }
    },
    {
      path: '/library/details',
      name: 'libraryDetails',
      component: () => import('../views/library/library-details/index.vue'),
      meta: {
        title: '编辑知识库',
        activeMenu: '/library',
        isCustomPage: true
      },
      redirect: '/library/details/knowledge-document',
      async beforeEnter (to, from, next){
        try {
          const { getLibraryInfo } = useLibraryStore()

          await getLibraryInfo(to.query.id)
    
          next()
        } catch (error) {
          next(false)
        }
      },
      children: [
        {
          path: '/library/details/knowledge-document',
          name: 'knowledgeDocument',
          component: () => import('../views/library/library-details/knowledge-document/index.vue'),
          meta: {
            title: '知识库文档',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/import-record',
          name: 'qaDocImportRecord',
          component: () => import('../views/library/library-details/import-record.vue'),
          meta: {
            title: '导入记录',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/export-record',
          name: 'qaDocExportRecord',
          component: () => import('../views/library/library-details/export-record.vue'),
          meta: {
            title: '导出记录',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/recycle-bin-record',
          name: 'recycleBinRecord',
          component: () => import('../views/library/library-details/recycle-bin-record.vue'),
          meta: {
            title: '回收站',
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
          path: '/library/details/knowledge-graph',
          name: 'knowledgeGraph',
          component: () => import('@/views/library/knowledge-graph/index.vue'),
          meta: {
            title: '知识图谱',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/knowledge-config',
          name: 'knowledgeConfig',
          component: () => import('@/views/library/library-details/knowledge-config/index.vue'),
          meta: {
            title: '知识库配置',
            activeMenu: '/library',
            cuStyle: true,
          }
        },
        {
          path: '/library/details/related-robots',
          name: 'knowledgeRelatedRobots',
          component: () => import('@/views/library/library-details/related-robots.vue'),
          meta: {
            title: '关联机器人',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/categary-manage',
          name: 'categaryManages',
          component: () => import('@/views/library/library-details/categary-manage/index.vue'),
          meta: {
            title: '精选',
            activeMenu: '/library'
          }
        }
      ]
    },
    {
      path: '/library/preview',
      name: 'libraryPreview',
      component: () => import('../views/library/library-preview/library-preview.vue'),
      meta: {
        title: '知识库管理',
        activeMenu: '/library',
        pageStyle:{
          'padding': '0 24px'
        },
        hideTitle: true
      }
    }
  ]
}
