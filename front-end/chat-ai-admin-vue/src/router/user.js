export default {
  path: '/user',
  name: 'User',
  component: () => import('@/layouts/AdminLayout/index.vue'),
  redirect: '/user/account',
  meta: {
    title: 'routes.basic.account_settings',
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
        title: 'routes.basic.account_settings',
        activeMenu: '/user'
      },
      children: [
        {
          path: '/user/account',
          name: 'userAccount',
          component: () => import('../views/user/account/index.vue'),
          meta: {
            title: 'routes.basic.account_settings',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/model',
          name: 'userModel',
          component: () => import('../views/user/model/index.vue'),
          meta: {
            title: 'routes.basic.model_management',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/chatclaw-token-log',
          name: 'userChatclawTokenLog',
          component: () => import('../views/user/chatclaw-token-log/index.vue'),
          meta: {
            title: 'routes.basic.chatclaw_token_log',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/usetoken',
          name: 'userUsetoken',
          component: () => import('../views/user/usetoken/index.vue'),
          meta: {
            title: 'routes.basic.token_usage',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/manage',
          name: 'userManage',
          component: () => import('../views/user/manage/index.vue'),
          meta: {
            title: 'routes.basic.team_management',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/aliocr',
          name: 'userAliocr',
          component: () => import('../views/user/aliocr/index.vue'),
          meta: {
            title: 'routes.basic.aliyun_ocr',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/enterprise',
          name: 'userEnterprise',
          component: () => import('../views/user/enterprise/index.vue'),
          meta: {
            title: 'routes.basic.enterprise_settings',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/domain',
          name: 'userDomain',
          component: () => import('../views/user/domain/index.vue'),
          meta: {
            title: 'routes.basic.custom_domain',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/clientDownload',
          name: 'clientDownload',
          component: () => import('../views/user/client-download/index.vue'),
          meta: {
            title: 'routes.basic.client_download',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/sensitive-words',
          name: 'SensitiveWords',
          component: () => import('../views/user/sensitive-words/index.vue'),
          meta: {
            title: 'routes.basic.sensitive_word_management',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/prompt-library',
          name: 'promptLibrary',
          component: () => import('../views/user/prompt-library/index.vue'),
          meta: {
            title: 'routes.basic.prompt_template_library',
            activeMenu: '/user'
          }
        },
        {
          path: '/user/official-account',
          name: 'userOfficialAccount',
          component: () => import('../views/user/official-account/index.vue'),
          meta: {
            title: 'routes.basic.official_account_management',
            activeMenu: '/user'
          }
        },
      ]
    }
  ]
}
