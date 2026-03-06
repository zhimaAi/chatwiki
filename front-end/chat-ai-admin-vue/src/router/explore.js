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
        title: 'routes.basic.explore',
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
        title: 'routes.basic.custom_menu',
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
        title: 'routes.basic.subscribe_reply',
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
        title: 'routes.basic.add_rule',
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
        title: 'routes.basic.article_group_send',
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
          path: '/explore/index/article-group-send/group-send',
          name: 'exploreArticleGroupSendGroupSend',
          component: () => import('../views/explore/article-group-send/group-send.vue'),
          meta: {
            title: 'routes.basic.group_send_management',
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
            title: 'routes.basic.draft_box',
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
        title: 'routes.basic.ai_comment_management',
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
            title: 'routes.basic.default_rule',
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
            title: 'routes.basic.custom_rule',
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
            title: 'routes.basic.create_custom_rule',
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
            title: 'routes.basic.comment_processing_record',
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
