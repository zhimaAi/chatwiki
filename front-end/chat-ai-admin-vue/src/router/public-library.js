import { message } from 'ant-design-vue'
import { usePublicLibraryStore } from '@/stores/modules/public-library'

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

export default {
  path: '/public-library',
  name: 'PublicLibrary',
  component: () => import('../layouts/AdminLayout/index.vue'),
  meta: {
    title: '对外文档',
    activeMenu: 'PublicLibrary'
  },
  redirect: '/public-library/list',
  children: [
    {
      path: '/public-library/list',
      name: 'PublicLibraryList',
      component: () => import('@/views/public-library/list/index.vue'),
      meta: {
        title: '对外文档',
        bgColor: '#fff',
        hideTitle: true
      }
    },
    {
      path: '/public-library/add',
      name: 'AddPublicLibrary',
      component: () => import('../views/public-library/add/index.vue'),
      meta: {
        title: '新建文档',
        activeMenu: 'PublicLibrary',
        breadcrumb: [
          {
            title: '对外文档',
            path: '/public-library/list'
          },
          {
            title: '新建文档',
            path: '/public-library/add'
          }
        ]
      }
    },
    {
      path: '/public-library/layout',
      name: 'PublicLibrary',
      component: () => import('@/views/public-library/index.vue'),
      meta: {
        title: '知识库配置',
        isCustomPage: true,
        activeMenu: 'PublicLibrary'
      },
      async beforeEnter(to, from, next) {
        const store = usePublicLibraryStore()
        try {
          let res = await store.getLibDocPartner(to.query)
          if (!res.data.operate_rights) {
            next('/public-library/list')
            message.warning('您没有权限访问此文档')
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
            title: '知识库配置',
            bgColor: '#F5F9FF',
            activeMenu: 'config'
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
            title: '访问权限',
            bgColor: '#F5F9FF',
            activeMenu: 'config'
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
            title: '首页',
            bgColor: '#F5F9FF',
            activeMenu: 'home'
          }
        },
        {
          path: '/public-library/ai',
          name: 'PublicLibraryAi',
          component: () => import('@/views/public-library/ai/index.vue'),
          meta: {
            title: '文档AI',
            bgColor: '#F5F9FF',
            activeMenu: 'config'
          }
        },
        {
          path: '/public-library/web-statistics',
          name: 'PublicLibraryWebStatistics',
          component: () => import('@/views/public-library/web-statistics/index.vue'),
          meta: {
            title: '统计设置',
            bgColor: '#F5F9FF',
            activeMenu: 'config'
          }
        },
        {
          path: '/public-library/editor',
          name: 'PublicLibraryEditor',
          component: () => import('@/views/public-library/editor/index.vue'),
          meta: {
            title: '编辑文档',
            bgColor: '#F5F9FF',
            activeMenu: 'doc'
          }
        }
      ]
    }
  ]
}
