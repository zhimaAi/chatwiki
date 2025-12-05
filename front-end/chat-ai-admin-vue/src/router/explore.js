export default {
  path: '/explore',
  name: 'Explore',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/explore/index',
  children: [
    {
      path: '/explore/index',
      name: 'exploreIndex',
      component: () => import('../views/explore/explore-index/explore-index.vue'),
      meta: {
        title: '探索',
        hideTitle: true,
        activeMenu: '/explore',
        bgColor: '#ffffff',
        pageStyle: {
          'padding-right': 0,
          overflow: 'hidden'
        }
      }
    },
    {
      path: '/explore/index/subscribe-reply',
      name: 'exploreSubscribeReply',
      component: () => import('@/views/robot/robot-config/subscribe-reply/index.vue'),
      meta: {
        title: '关注后回复',
        isCustomPage: true,
        activeMenu: '/explore',
        bgColor: '#ffffff',
        pageStyle: {
          'padding-right': 0,
          overflow: 'hidden'
        }
      }
    },
    {
      path: '/explore/index/subscribe-reply/add-rule',
      name: 'exploreSubscribeReplyAddRule',
      component: () => import('@/views/robot/robot-config/subscribe-reply/add-rule.vue'),
      meta: {
        title: '新增规则',
        isCustomPage: true,
        activeMenu: '/explore',
        bgColor: '#ffffff',
        pageStyle: {
          'padding-right': 0,
          overflow: 'hidden'
        }
      }
    }
  ]
}
