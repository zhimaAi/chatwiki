import { message } from 'ant-design-vue'
import { usePublicLibraryStore } from '@/stores/modules/public-library'
import { i18n } from '@/locales'

const checkRouteManagePermission = (to, from, next) => {
  const { operate_rights } = usePublicLibraryStore()

  if (operate_rights == 4) {
    next()
  } else {
    next({
      path: '/public-library/home',
      query: to.query
    })
  }
}

export default [
  {
    path: '/public-library',
    name: 'PublicLibrary',
    component: () => import('../layouts/AdminLayout/index.vue'),
    meta: {
      title: 'routes.basic.public_document',
      activeMenu: 'PublicLibrary'
    },
    redirect: '/public-library/list',
    children: [
      {
        path: '/public-library/list',
        name: 'PublicLibraryList',
        component: () => import('@/views/public-library/list/index.vue'),
        meta: {
          title: 'routes.basic.public_document',
          bgColor: '#fff',
          hideTitle: true
        }
      },
      {
        path: '/public-library/add',
        name: 'AddPublicLibrary',
        component: () => import('@/views/public-library/add/index.vue'),
        meta: {
          title: 'routes.basic.create_document',
          activeMenu: 'PublicLibrary',
          breadcrumb: [
            {
              title: 'routes.basic.public_document',
              path: '/public-library/list'
            },
            {
              title: 'routes.basic.create_document',
              path: '/public-library/add'
            }
          ]
        }
      }
    ]
  },
  {
    path: '/public-library/layout',
    name: 'PublicLibraryLayout',
    component: () => import('@/views/public-library/index.vue'),
    meta: {
      title: 'routes.basic.knowledge_config',
      isCustomPage: true,
      activeMenu: 'PublicLibrary'
    },
    async beforeEnter(to, from, next) {
      const store = usePublicLibraryStore()
      try {
        await store.getInfo(to.query)

        if (!store.operate_rights) {
          next('/public-library/list')
          message.warning(i18n && i18n.global ? i18n.global.t('common.noPermissionAccess') : '您没有权限访问此文档')
        } else {
          next()
        }
      } catch (error) {
        next(false)
      }
    },
    redirect: '/public-library/home',
    children: [
      {
        path: '/public-library/config',
        name: 'PublicLibraryConfig',
        component: () => import('@/views/public-library/config/index.vue'),
        meta: {
          title: 'routes.basic.knowledge_config',
          bgColor: '#F5F9FF',
          subActiveMenu: 'config',
          activeMenu: 'PublicLibrary'
        },
        beforeEnter: (to, from, next) => {
          checkRouteManagePermission(to, from, next)
        }
      },
      {
        path: '/public-library/permissions',
        name: 'PublicLibraryPermissions',
        component: () => import('@/views/public-library/permissions/index.vue'),
        meta: {
          title: 'routes.basic.access_permission',
          bgColor: '#F5F9FF',
          subActiveMenu: 'config',
          activeMenu: 'PublicLibrary'
        },
        beforeEnter: (to, from, next) => {
          checkRouteManagePermission(to, from, next)
        }
      },
      {
        path: '/public-library/home',
        name: 'PublicLibraryHome',
        component: () => import('@/views/public-library/home/index.vue'),
        meta: {
          title: 'routes.basic.home',
          bgColor: '#F5F9FF',
          subActiveMenu: 'home',
          activeMenu: 'PublicLibrary'
        }
      },
      {
        path: '/public-library/ai',
        name: 'PublicLibraryAi',
        component: () => import('@/views/public-library/ai/index.vue'),
        meta: {
          title: 'routes.basic.document_ai',
          bgColor: '#F5F9FF',
          subActiveMenu: 'config',
          activeMenu: 'PublicLibrary'
        }
      },
      {
        path: '/public-library/web-statistics',
        name: 'PublicLibraryWebStatistics',
        component: () => import('@/views/public-library/web-statistics/index.vue'),
        meta: {
          title: 'routes.basic.statistics_settings',
          bgColor: '#F5F9FF',
          subActiveMenu: 'config',
          activeMenu: 'PublicLibrary'
        }
      },
      {
        path: '/public-library/editor',
        name: 'PublicLibraryEditor',
        component: () => import('@/views/public-library/editor/index.vue'),
        meta: {
          title: 'routes.basic.edit_document',
          bgColor: '#F5F9FF',
          subActiveMenu: 'doc',
          activeMenu: 'PublicLibrary'
        }
      }
    ]
  }
]
