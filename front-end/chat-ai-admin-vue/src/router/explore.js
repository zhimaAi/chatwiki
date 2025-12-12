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
      path: '/explore/index/custom-menu',
      name: 'exploreCustomMenu',
      component: () => import('@/views/robot/robot-config/custom-menu/index.vue'),
      meta: {
        title: '自定义菜单',
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
    },
    {
      path: '/explore/index/article-group-send',
      name: 'exploreArticleGroupSend',
      component: () => import('../views/explore/article-group-send/index.vue'),
      meta: {
        title: '文章群发',
        isCustomPage: true,
        activeMenu: '/explore',
        bgColor: '#ffffff',
        pageStyle: {
          'padding-right': 0,
          overflow: 'hidden'
        }
      },
      // 二级tab：群发管理 草稿箱
      children: [
        {
          path: '/explore/index/article-group-send/group-send',
          name: 'exploreArticleGroupSendGroupSend',
          component: () => import('../views/explore/article-group-send/group-send.vue'),
          meta: {
            title: '群发管理',
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
          path: '/explore/index/article-group-send/draft-box',
          name: 'exploreArticleGroupSendDraftBox',
          component: () => import('../views/explore/article-group-send/draft-box.vue'),
          meta: {
            title: '草稿箱',
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
    },
    {
      path: '/explore/index/ai-comment-management',
      name: 'exploreAiCommentManagement',
      component: () => import('../views/explore/ai-comment-management/index.vue'),
      meta: {
        title: 'AI评论管理',
        isCustomPage: true,
        activeMenu: '/explore',
        bgColor: '#ffffff',
        pageStyle: {
          'padding-right': 0,
          overflow: 'hidden'
        }
      },
      children: [
        {
          path: '/explore/index/ai-comment-management/default-rule',
          name: 'exploreAiCommentManagementDefaultRule',
          component: () => import('../views/explore/ai-comment-management/default-rule.vue'),
          meta: {
            title: '默认规则',
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
          path: '/explore/index/ai-comment-management/custom-rule',
          name: 'exploreAiCommentManagementCustomRule',
          component: () => import('../views/explore/ai-comment-management/custom-rule.vue'),
          meta: {
            title: '自定义规则',
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
          path: '/explore/index/ai-comment-management/create-custom-rule',
          name: 'exploreAiCommentManagementCreateCustomRule',
          component: () => import('../views/explore/ai-comment-management/create-custom-rule.vue'),
          meta: {
            title: '新建自定义规则',
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
          path: '/explore/index/ai-comment-management/comment-processing-record',
          name: 'exploreAiCommentManagementCommentProcessingRecord',
          component: () => import('../views/explore/ai-comment-management/comment-processing-record.vue'),
          meta: {
            title: '评论处理记录',
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
  ]
}
