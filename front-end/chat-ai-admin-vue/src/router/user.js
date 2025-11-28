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
          path: '/user/usetoken',
          name: 'userUsetoken',
          component: () => import('../views/user/usetoken/index.vue'),
          meta: {
            title: 'Token使用',
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
          path: '/user/aliocr',
          name: 'userAliocr',
          component: () => import('../views/user/aliocr/index.vue'),
          meta: {
            title: '阿里云OCR',
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
          path: '/user/domain',
          name: 'userDomain',
          component: () => import('../views/user/domain/index.vue'),
          meta: {
            title: '自定义域名',
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
        },
        {
          path: '/user/sensitive-words',
          name: 'SensitiveWords',
          component: () => import('../views/user/sensitive-words/index.vue'),
          meta: {
            title: '敏感词管理',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/prompt-library',
          name: 'promptLibrary',
          component: () => import('../views/user/prompt-library/index.vue'),
          meta: {
            title: '提示词模板库',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/official-account',
          name: 'userOfficialAccount',
          component: () => import('../views/user/official-account/index.vue'),
          meta: {
            title: '公众号管理',
            activeMenu: '/user'
          }
        },
      ]
    }
  ]
}
