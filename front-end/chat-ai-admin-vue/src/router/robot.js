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
          path: '/robot/config/external-services',
          name: 'externalServices',
          component: () => import('@/views/robot/robot-config/external-service/index.vue'),
          meta: {
            title: '对外服务',
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
    }
  ]
}
