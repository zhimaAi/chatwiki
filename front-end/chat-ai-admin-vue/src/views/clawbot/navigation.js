import { CUSTOMER_SERVICE_LAYOUT, isCustomerServiceClawbotRoute } from '@/utils/clawbot-layout'

/** Agent 始终使用原路由；新版入口额外保留布局模式和当前应用标识。 */
export function getClawbotLocation(route, page, query = route.query) {
  const isCustomerLayout = isCustomerServiceClawbotRoute(route)
  return {
    path: `/clawbot/${page}`,
    query: isCustomerLayout ? {
      id: route.query.id,
      robot_key: route.query.robot_key,
      ...query,
      layout: CUSTOMER_SERVICE_LAYOUT
    } : query
  }
}
