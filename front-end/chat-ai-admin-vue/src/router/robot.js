export default {
  path: '/robot',
  name: 'Robot',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/robot/list',
  children: [
    {
      path: '/robot/list',
      name: 'robotList',
      component: () => import('../views/robot/robot-list/robot-list.vue'),
      meta: {
        title: '机器人管理',
        hideTitle: true,
        activeMenu: '/robot',
        bgColor: '#ffffff',
        pageStyle: {
          'padding-right': 0,
          overflow: 'hidden'
        }
      }
    },
    {
      path: '/robot/config',
      name: 'robotDetails',
      component: () => import('../views/robot/robot-config/index.vue'),
      meta: {
        title: '编辑机器人',
        activeMenu: '/robot',
        isCustomPage: true
      },
      redirect: '/robot/config/basic-config',
      children: [
        {
          path: '/robot/config/basic-config',
          name: 'basicConfig',
          component: () => import('@/views/robot/robot-config/basic-config/index.vue'),
          meta: {
            title: '编辑机器人',
            isCustomPage: true,
            activeMenu: '/robot'
          }
        },
        {
          path: '/robot/config/workflow',
          name: 'robotWorkflow',
          component: () => import('../views/workflow/index.vue'),
          meta: {
            title: '工作流编排',
            isCustomPage: true,
          }
        },
        {
          path: '/robot/ability/smart-menu',
          name: 'smartMenu',
          component: () => import('@/views/robot/robot-config/smart-menu/index.vue'),
          meta: {
            title: '智能菜单',
            isCustomPage: true
          }
        },
        {
          path: '/robot/ability/smart-menu/add-rule',
          name: 'smartMenuAddRule',
          component: () => import('@/views/robot/robot-config/smart-menu/add-rule.vue'),
          meta: {
            title: '新增规则',
            isCustomPage: true
          }
        },
        {
          path: '/robot/config/function-center',
          name: 'functionCenter',
          component: () => import('@/views/robot/robot-config/function-center/index.vue'),
          meta: {
            title: '功能中心',
            isCustomPage: true
          }
        },
        {
          path: '/robot/config/external-services',
          name: 'externalServices',
          component: () => import('@/views/robot/robot-config/external-service/index.vue'),
          meta: {
            title: '对外服务',
            isCustomPage: true
          }
        },
        {
          path: '/robot/ability/auto-reply',
          name: 'autoReply',
          component: () => import('@/views/robot/robot-config/auto-reply/index.vue'),
          meta: {
            title: '关键词回复',
            isCustomPage: true
          }
        },
        {
          path: '/robot/ability/auto-reply/add-rule',
          name: 'addRule',
          component: () => import('@/views/robot/robot-config/auto-reply/add-rule.vue'),
          meta: {
            title: '新增规则',
            isCustomPage: true
          }
        },
        {
          path: '/robot/ability/auto-reply/add-reply',
          name: 'addReply',
          component: () => import('@/views/robot/robot-config/auto-reply/add-reply.vue'),
          meta: {
            title: '新增回复',
            isCustomPage: true
          }
        },
        {
          path: '/robot/config/library-config',
          name: 'libraryConfig',
          component: () => import('@/views/robot/robot-config/library-config/index.vue'),
          meta: {
            title: '知识库',
            isCustomPage: true
          }
        },
        {
          path: '/robot/config/skill-config',
          name: 'skillConfig',
          component: () => import('@/views/robot/robot-config/skill-config/index.vue'),
          meta: {
            title: '工作流',
            isCustomPage: true
          }
        },
        {
          path: '/robot/config/api-key-manage',
          name: 'apiKeyManage',
          component: () => import('@/views/robot/api-key-manage/index.vue'),
          meta: {
            title: 'API key管理',
            isCustomPage: true
          }
        },
        {
          path: '/robot/config/qa-feedbacks',
          name: 'qaFeedbacks',
          component: () => import('@/views/robot/robot-config/qa-feedback/index.vue'),
          meta: {
            title: '问答反馈',
            isCustomPage: true
          }
        },
        {
          path: '/robot/config/session-record',
          name: 'sessionRecord',
          component: () => import('@/views/robot/robot-config/session-record/index.vue'),
          meta: {
            title: '会话记录',
            isCustomPage: true
          }
        },
        {
          path: '/robot/config/export-record',
          name: 'exportRecord',
          component: () => import('@/views/robot/robot-config/export-record/index.vue'),
          meta: {
            title: '导出记录',
            isCustomPage: true
          }
        },
        {
          path: '/robot/config/invoke-logs',
          name: 'invokeLogs',
          component: () => import('@/views/robot/robot-config/invoke-logs/index.vue'),
          meta: {
            title: '调用日志',
            isCustomPage: true
          }
        },
        {
          path: '/robot/config/unknown_issue',
          name: 'unknownIssue',
          component: () => import('@/views/robot/robot-config/unknown_issue/unknow-index.vue'),
          meta: {
            title: '未知问题',
            isCustomPage: true
          }
        },
        {
          path: '/robot/config/statistical_analysis',
          name: 'statisticalAnalysis',
          component: () => import('@/views/robot/robot-config/statistical_analysis/index.vue'),
          meta: {
            title: '统计分析',
            isCustomPage: true
          }
        },
      ]
    },
    {
      path: '/robot/test',
      name: 'robotTest',
      component: () => import('../views/robot/robot-test/index.vue'),
      meta: {
        title: '机器人管理',
        hideTitle: true,
        activeMenu: '/robot',
        bgColor: '#F5F9FF',
        breadcrumb: false,
        isCustomPage: true
      }
    },

  ]
}
