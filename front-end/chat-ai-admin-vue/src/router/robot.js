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
        bgColor: '#F5F9FF',
        pageStyle: {}
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
          path: '/robot/config/external-services',
          name: 'externalServices',
          component: () => import('@/views/robot/robot-config/external-service/index.vue'),
          meta: {
            title: '对外服务',
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
          path: '/robot/config/statistical_analysis',
          name: 'statisticalAnalysis',
          component: () => import('@/views/robot/robot-config/statistical_analysis/index.vue'),
          meta: {
            title: '统计分析',
            isCustomPage: true
          }
        }
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
