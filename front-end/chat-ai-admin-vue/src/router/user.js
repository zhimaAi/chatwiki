export default {
  path: '/user',
  name: 'User',
  component: () => import('@/layouts/AdminLayout/index.vue'),
  redirect: '/user/account',
  meta: {
    title: '账号设置',
    activeMenu: '/user',
    breadcrumb: false,
    hideTitle: true,
    isCustomPage: true
  },
  children: [
    {
      path: '/user/layout',
      name: 'userLayout',
      component: () => import('../views/user/index.vue'),
      redirect: '/user/account',
      meta: {
        title: '账号设置',
        activeMenu: '/user'
      },
      children: [
        {
          path: '/user/account',
          name: 'userAccount',
          component: () => import('../views/user/account/index.vue'),
          meta: {
            title: '账号设置',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/model',
          name: 'userModel',
          component: () => import('../views/user/model/index.vue'),
          meta: {
            title: '模型管理',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/manage',
          name: 'userManage',
          component: () => import('../views/user/manage/index.vue'),
          meta: {
            title: '团队管理',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/enterprise',
          name: 'userEnterprise',
          component: () => import('../views/user/enterprise/index.vue'),
          meta: {
            title: '企业设置',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/clientDownload',
          name: 'clientDownload',
          component: () => import('../views/user/client-download/index.vue'),
          meta: {
            title: '客户端下载',
            activeMenu: '/user'
          }
        }
      ]
    }
  ]
}
