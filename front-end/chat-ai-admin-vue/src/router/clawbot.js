import customerServiceRoutes from './customer-service'
import { getCustomerServiceAgentPage } from '@/views/customer-service/navigation'

// 仅新版 Agent 的菜单会进入这些路由；直接复用客服业务组件，不复制页面。
const customerAgentRoutes = customerServiceRoutes.children
  .filter(page => ['session-record', 'web-app', 'other-channels', 'basic-settings', 'function-center', 'analytics'].includes(page.path))
  .map(page => ({ ...page, name: getCustomerServiceAgentPage(page.path).routeName }))

export default {
  path: '/clawbot',
  name: 'Clawbot',
  redirect: '/clawbot/chat',
  component: () => import('@/views/clawbot/index.vue'),
  children: [
    ...customerAgentRoutes,
    {
      path: 'chat',
      name: 'clawbotChat',
      meta: {
        title: 'Clawbot Chat'
      },
      component: () => import('@/views/clawbot/chat/index.vue')
    },
    {
      path: 'assistant',
      name: 'clawbotAssistant',
      meta: {
        title: 'Clawbot Assistant'
      },
      component: () => import('@/views/clawbot/assistant/index.vue')
    },
    {
      path: 'skills',
      name: 'clawbotSkills',
      meta: {
        title: 'Clawbot Skills'
      },
      component: () => import('@/views/clawbot/skills/index.vue')
    },
    {
      path: 'skill-generate-tool',
      name: 'clawbotSkillGenerateTool',
      meta: {
        title: 'skill生成工具'
      },
      component: () => import('@/views/clawbot/skill-generate-tool/index.vue')
    },
    {
      path: 'knowledge',
      name: 'clawbotKnowledge',
      meta: {
        title: 'Clawbot Knowledge'
      },
      component: () => import('@/views/clawbot/knowledge/index.vue')
    },
    {
      path: 'services',
      name: 'clawbotServices',
      meta: {
        title: 'Clawbot Services'
      },
      component: () => import('@/views/clawbot/services/index.vue')
    },
    {
      path: 'stats',
      name: 'clawbotStats',
      meta: {
        title: 'Clawbot Stats'
      },
      component: () => import('@/views/clawbot/stats/index.vue')
    },
    {
      path: 'settings',
      name: 'clawbotSettings',
      meta: {
        title: 'Clawbot Settings'
      },
      component: () => import('@/views/clawbot/settings/index.vue')
    }
  ]
}
