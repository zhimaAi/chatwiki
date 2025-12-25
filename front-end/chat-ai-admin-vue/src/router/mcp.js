export default {
  path: '/mcp',
  name: 'mcp',
  component: () => import('../layouts/AdminLayout/index.vue'),
  redirect: '/mcp/index',
  children: [
    {
      path: '/mcp/index',
      name: 'mcpIndex',
      component: () => import('../views/explore/mcp/index.vue'),
      meta: {
        title: 'MCP管理',
        hideTitle: true,
        activeMenu: 'explore',
        bgColor: '#ffffff',
        pageStyle: {
          'padding-right': 0,
          overflow: 'hidden'
        }
      }
    }
  ]
}
