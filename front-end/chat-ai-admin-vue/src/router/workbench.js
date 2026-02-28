export default {
    path: '/workbench',
    name: 'Workbench',
    meta: {
        title: '工作台'
    },
    component: () => import('@/views/workbench/index.vue'),
    children: [
        {
            path: 'chat',
            name: 'workbenchChat',
            meta: {
                title: '首页',
            },
            component: () => import('@/views/workbench/chat/index.vue')
        },
        {
            path: 'team-tools',
            name: 'workbenchTeamTools',
            meta: {
                title: '团队应用'
            },
            component: () => import('@/views/workbench/team-tools/index.vue')
        },
        {
            path: 'home-config',
            name: 'workbenchHomeConfig',
            meta: {
                title: '首页配置'
            },
            component: () => import('@/views/workbench/home-config/index.vue')
        }
    ]
}