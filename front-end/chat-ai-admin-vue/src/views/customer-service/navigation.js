import { CUSTOMER_SERVICE_LAYOUT } from '@/utils/clawbot-layout'

export const CUSTOMER_SERVICE_DEFAULT_MENU_KEY = 'session-record'

// 菜单和路由共享稳定标识；标签保持为i18n key，由展示层响应式翻译。
export const CUSTOMER_SERVICE_MENU_ITEMS = [
  {
    key: 'chat',
    titleKey: 'chat',
    routeName: 'customerServiceChat'
  },
  {
    key: 'session-record',
    titleKey: 'session_record',
    routeName: 'customerServiceSessionRecord'
  },
  {
    key: 'channels',
    titleKey: 'channels',
    children: [
      {
        key: 'web-app',
        titleKey: 'web_app',
        routeName: 'customerServiceWebApp'
      },
      {
        key: 'other-channels',
        titleKey: 'other_channels',
        routeName: 'customerServiceOtherChannels'
      }
    ]
  },
  {
    key: 'ai-agent',
    titleKey: 'ai_agent',
    children: [
      {
        key: 'basic-settings',
        titleKey: 'basic_settings',
        routeName: 'customerServiceBasicSettings'
      },
      {
        key: 'knowledge-base',
        titleKey: 'knowledge_base',
        routeName: 'customerServiceKnowledgeBase'
      },
      {
        key: 'function-center',
        titleKey: 'function_center',
        routeName: 'customerServiceFunctionCenter'
      },
      {
        key: 'skill',
        titleKey: 'skill',
        routeName: 'customerServiceSkill'
      }
    ]
  },
  {
    key: 'analytics',
    titleKey: 'analytics',
    routeName: 'customerServiceAnalytics'
  }
]

const agentRoutes = {
  chat: 'clawbotChat',
  'session-record': 'customerAgentSessionRecord',
  'web-app': 'customerAgentWebApp',
  'other-channels': 'customerAgentOtherChannels',
  'basic-settings': 'customerAgentBasicSettings',
  'knowledge-base': 'clawbotKnowledge',
  'function-center': 'customerAgentFunctionCenter',
  skills: 'clawbotSkills',
  analytics: 'customerAgentAnalytics'
}

// 新版 Agent 共用客服菜单结构，业务内容仍使用各自的应用上下文。
export const CUSTOMER_SERVICE_AGENT_MENU_ITEMS = CUSTOMER_SERVICE_MENU_ITEMS.map(item => {
  if (!item.children) return { ...item, routeName: agentRoutes[item.key] }
  return {
    ...item,
    children: item.children.map(page => {
      const key = page.key === 'skill' ? 'skills' : page.key
      return { ...page, key, routeName: agentRoutes[key] }
    })
  }
})

export function getCustomerServiceAgentPage(menuKey) {
  const key = menuKey === 'skill' ? 'skills' : menuKey
  return CUSTOMER_SERVICE_AGENT_MENU_ITEMS.flatMap(item => item.children || [item])
    .find(item => item.key === key)
}

export function getCustomerServicePage(menuKey) {
  for (const menuItem of CUSTOMER_SERVICE_MENU_ITEMS) {
    if (menuItem.children) {
      const page = menuItem.children.find((child) => child.key === menuKey)
      if (page) return page
    } else if (menuItem.key === menuKey) {
      return menuItem
    }
  }
  return null
}

/** 新版框架选择应用时，按类型使用原有业务路由。 */
export function getCustomerServiceChatLocation(application) {
  const isAgent = Number(application.application_type) === 2
  return {
    name: isAgent ? 'clawbotChat' : 'customerServiceChat',
    query: {
      id: normalizeCustomerServiceQueryValue(application.id),
      robot_key: normalizeCustomerServiceQueryValue(application.robot_key),
      ...(isAgent ? { layout: CUSTOMER_SERVICE_LAYOUT } : {})
    }
  }
}

/** URL重复参数按首项读取；标识不额外限制为整数，也不去除有效字符。 */
export function normalizeCustomerServiceQueryValue(value) {
  const firstValue = Array.isArray(value) ? value[0] : value
  return firstValue === undefined || firstValue === null ? '' : String(firstValue)
}

/**
 * @returns {{ robotId: string, robotKey: string }}
 * 保留单缺状态，由布局提示无效链接；不能在这里把残缺标识当成无参。
 */
export function readCustomerServiceRobotIdentity(query = {}) {
  return {
    robotId: normalizeCustomerServiceQueryValue(query.id),
    robotKey: normalizeCustomerServiceQueryValue(query.robot_key)
  }
}

/** 默认入口和局部路径回退只保留机器人字段，残缺标识继续由布局提示。 */
export function pickCustomerServiceRobotQuery(query = {}) {
  const { robotId, robotKey } = readCustomerServiceRobotIdentity(query)
  const robotQuery = {}
  if (robotId) robotQuery.id = robotId
  if (robotKey) robotQuery.robot_key = robotKey
  return robotQuery
}

