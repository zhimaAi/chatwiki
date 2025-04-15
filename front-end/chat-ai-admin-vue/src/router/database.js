export default {
  path: '/database',
  name: 'Database',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/database/list',
  children: [
    {
      path: '/database/list',
      name: 'databaseList',
      component: () => import('../views/database/database-list/index.vue'),
      meta: {
        title: '数据库',
        activeMenu: 'library',
        bgColor: '#fff',
        hideTitle: true
      }
    },
    {
      path: '/database/details',
      name: 'databaseDetails',
      component: () => import('../views/database/database-detail/index.vue'),
      meta: {
        title: '数据库管理',
        activeMenu: 'library',
        isCustomPage: true,
      },
      redirect: '/database/details/field-manage',
      children: [
        {
          path: '/database/details/field-manage',
          name: 'fieldManage',
          component: () => import('../views/database/database-detail/field-manage/index.vue'),
          meta: {
            title: '字段管理',
            activeMenu: 'library'
          }
        },
        {
          path: '/database/details/database-manage',
          name: 'databaseManage',
          component: () => import('../views/database/database-detail/database-manage/index.vue'),
          meta: {
            title: '数据管理',
            activeMenu: 'library'
          }
        },
      ],
    },
  ]
}
