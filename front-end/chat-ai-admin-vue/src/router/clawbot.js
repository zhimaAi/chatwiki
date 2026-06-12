export default {
  path: '/clawbot',
  name: 'Clawbot',
  redirect: '/clawbot/chat',
  component: () => import('@/views/clawbot/index.vue'),
  children: [
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
