export default {
    path: '/workbench',
    name: 'Workbench',
    meta: {
        title: 'routes.basic.workbench'
    },
    component: () => import('@/views/workbench/index.vue'),
    children: [
        {
            path: 'chat',
            name: 'workbenchChat',
            meta: {
                title: 'routes.basic.home_page',
            },
            component: () => import('@/views/workbench/chat/index.vue')
        },
        {
            path: 'team-tools',
            name: 'workbenchTeamTools',
            meta: {
                title: 'routes.basic.team_apps'
            },
            component: () => import('@/views/workbench/team-tools/index.vue')
        },
        {
            path: 'home-config',
            name: 'workbenchHomeConfig',
            meta: {
                title: 'routes.basic.home_config'
            },
            component: () => import('@/views/workbench/home-config/index.vue')
        }
    ]
}