import {
  CUSTOMER_SERVICE_DEFAULT_MENU_KEY,
  getCustomerServicePage,
  pickCustomerServiceRobotQuery
} from '@/views/customer-service/navigation'

const redirectToSessionRecord = (to) => ({
  name: getCustomerServicePage(CUSTOMER_SERVICE_DEFAULT_MENU_KEY).routeName,
  query: pickCustomerServiceRobotQuery(to.query)
})

export default {
  path: '/customer-service',
  name: 'CustomerService',
  component: () => import('@/views/customer-service/index.vue'),
  redirect: redirectToSessionRecord,
  children: [
    {
      path: 'chat',
      name: 'customerServiceChat',
      component: () => import('@/views/customer-service/chat/index.vue'),
      meta: {
        title: 'routes.basic.chat',
        customerServiceMenuKey: 'chat',
        customerServiceGroupKey: null
      }
    },
    {
      path: 'session-record',
      name: 'customerServiceSessionRecord',
      // 设计稿版本；旧实现保留在 session-record/index.vue，便于回退。
      component: () => import('@/views/customer-service/session-record/design/index.vue'),
      meta: {
        title: 'routes.basic.session_record',
        customerServiceMenuKey: 'session-record',
        customerServiceGroupKey: null
      }
    },
    {
      path: 'web-app',
      name: 'customerServiceWebApp',
      component: () => import('@/views/customer-service/web-app/index.vue'),
      meta: {
        title: 'routes.basic.web_app',
        customerServiceMenuKey: 'web-app',
        customerServiceGroupKey: 'channels'
      }
    },
    {
      path: 'other-channels',
      name: 'customerServiceOtherChannels',
      component: () => import('@/views/customer-service/other-channels/index.vue'),
      meta: {
        title: 'routes.basic.other_channels',
        customerServiceMenuKey: 'other-channels',
        customerServiceGroupKey: 'channels'
      }
    },
    {
      path: 'basic-settings',
      name: 'customerServiceBasicSettings',
      component: () => import('@/views/customer-service/basic-settings/index.vue'),
      meta: {
        title: 'routes.basic.basic_settings',
        customerServiceMenuKey: 'basic-settings',
        customerServiceGroupKey: 'ai-agent'
      }
    },
    {
      path: 'knowledge-base',
      name: 'customerServiceKnowledgeBase',
      component: () => import('@/views/customer-service/knowledge-base/index.vue'),
      meta: {
        title: 'routes.basic.knowledge_base',
        customerServiceMenuKey: 'knowledge-base',
        customerServiceGroupKey: 'ai-agent'
      }
    },
    {
      path: 'function-center',
      name: 'customerServiceFunctionCenter',
      component: () => import('@/views/customer-service/function-center/index.vue'),
      meta: {
        title: 'routes.basic.function_center',
        customerServiceMenuKey: 'function-center',
        customerServiceGroupKey: 'ai-agent'
      }
    },
    {
      path: 'skill',
      name: 'customerServiceSkill',
      component: () => import('@/views/customer-service/skill/index.vue'),
      meta: {
        title: 'routes.basic.skill',
        customerServiceMenuKey: 'skill',
        customerServiceGroupKey: 'ai-agent'
      }
    },
    {
      path: 'analytics',
      name: 'customerServiceAnalytics',
      component: () => import('@/views/customer-service/analytics/index.vue'),
      meta: {
        title: 'routes.basic.statistical_analysis',
        customerServiceMenuKey: 'analytics',
        customerServiceGroupKey: null
      }
    },
    {
      path: ':pathMatch(.*)*',
      redirect: redirectToSessionRecord
    }
  ]
}
