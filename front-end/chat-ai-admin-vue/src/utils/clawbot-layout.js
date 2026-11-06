export const CUSTOMER_SERVICE_LAYOUT = 'customer-service'

/** 布局模式由 URL 决定，刷新和新标签页都能恢复；不影响独立入口的默认行为。 */
export function isCustomerServiceClawbotRoute(route = {}) {
  const path = route.path || ''
  const layout = Array.isArray(route.query?.layout) ? route.query.layout[0] : route.query?.layout
  return (path === '/clawbot' || path.startsWith('/clawbot/')) && layout === CUSTOMER_SERVICE_LAYOUT
}
