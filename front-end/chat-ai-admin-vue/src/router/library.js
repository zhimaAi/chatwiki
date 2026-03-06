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
        title: 'routes.basic.knowledge_management',
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
        title: 'routes.basic.create_knowledge',
        activeMenu: '/library',
        breadcrumb: [
          {
            title: 'routes.basic.knowledge_management',
            path: '/library/list'
          },
          {
            title: 'routes.basic.create_knowledge',
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
        title: 'routes.basic.create_knowledge',
        activeMenu: '/library',
        hideTitle: true
      }
    },
    {
      path: '/library/details',
      name: 'libraryDetails',
      component: () => import('../views/library/library-details/index.vue'),
      meta: {
        title: 'routes.basic.edit_knowledge',
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
            title: 'routes.basic.knowledge_document',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/import-record',
          name: 'qaDocImportRecord',
          component: () => import('../views/library/library-details/import-record.vue'),
          meta: {
            title: 'routes.basic.import_record',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/export-record',
          name: 'qaDocExportRecord',
          component: () => import('../views/library/library-details/export-record.vue'),
          meta: {
            title: 'routes.basic.export_record',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/recycle-bin-record',
          name: 'recycleBinRecord',
          component: () => import('../views/library/library-details/recycle-bin-record.vue'),
          meta: {
            title: 'routes.basic.recycle_bin',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/recall-testing',
          name: 'recallTesting',
          component: () => import('@/views/library/library-details/recall-testing.vue'),
          meta: {
            title: 'routes.basic.recall_testing',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/knowledge-graph',
          name: 'knowledgeGraph',
          component: () => import('@/views/library/knowledge-graph/index.vue'),
          meta: {
            title: 'routes.basic.knowledge_graph',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/knowledge-config',
          name: 'knowledgeConfig',
          component: () => import('../views/library/library-details/knowledge-config/index.vue'),
          meta: {
            title: 'routes.basic.knowledge_config',
            activeMenu: '/library',
            cuStyle: true,
          }
        },
        {
          path: '/library/details/related-robots',
          name: 'knowledgeRelatedRobots',
          component: () => import('../views/library/library-details/related-robots.vue'),
          meta: {
            title: 'routes.basic.related_robots',
            activeMenu: '/library'
          }
        },
        {
          path: '/library/details/categary-manage',
          name: 'categaryManages',
          component: () => import('../views/library/library-details/categary-manage/index.vue'),
          meta: {
            title: 'routes.basic.featured',
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
        title: 'routes.basic.knowledge_management',
        activeMenu: '/library',
        pageStyle:{
          'padding': '0 24px'
        },
        hideTitle: true
      }
    },
    {
      path: '/library/similar-question-list',
      name: 'similarQuestionList',
      component: () => import('../views/library/similar-question-list/index.vue'),
      meta: {
        title: 'routes.basic.similar_question_list',
        activeMenu: '/library',
        pageStyle:{
          'padding': '0 24px'
        },
        hideTitle: true
      }
    }
  ]
}
