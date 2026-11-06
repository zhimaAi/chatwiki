import { isCustomerServiceClawbotRoute } from './clawbot-layout'

function normalizeDomain(domain = '') {
  return String(domain || '').trim().replace(/\/+$/, '')
}

export function isClawbotPath(path = '') {
  const targetPath = String(path || '')
  return targetPath === '/clawbot' || targetPath.startsWith('/clawbot/')
}

export function getCrossDomainTarget({
  toPath = '',
  fromPath = '',
  toQuery = {},
  fromQuery = {},
  adminDomain = '',
  agentDomain = '',
  currentOrigin = ''
} = {}) {
  // set_token 只是跨域登录态落地页，继续参与切域会造成重定向循环。
  if (toPath === '/set_token') {
    return ''
  }

  // 新版框架中的 Agent 留在当前域；返回客服时也按同一框架分类，避免再次切域。
  if (isCustomerServiceClawbotRoute({ path: toPath, query: toQuery })) {
    return ''
  }
  const isToClawbot = isClawbotPath(toPath)
  const isFromClawbot = isClawbotPath(fromPath) &&
    !isCustomerServiceClawbotRoute({ path: fromPath, query: fromQuery })

  if (isToClawbot === isFromClawbot) {
    return ''
  }

  const targetDomain = normalizeDomain(isToClawbot ? agentDomain : adminDomain)
  const normalizedOrigin = normalizeDomain(currentOrigin)

  // 域名为空时保持当前 origin 内路由跳转，兼容后台未配置 admin/agent 域名的场景。
  if (!targetDomain || targetDomain === normalizedOrigin) {
    return ''
  }

  return targetDomain
}

export function buildSetTokenRedirectUrl({
  domain = '',
  redirectUrl = '/',
  token = '',
  exp = '',
  ttl = '',
  userId = '',
  userName = '',
  refreshUserInfo = true
} = {}) {
  const normalizedDomain = normalizeDomain(domain)
  const searchParams = new URLSearchParams({
    token: String(token || ''),
    exp: String(exp || ''),
    ttl: String(ttl || ''),
    user_id: String(userId || ''),
    user_name: String(userName || ''),
    redirect_url: String(redirectUrl || '/')
  })

  if (refreshUserInfo) {
    searchParams.set('refresh_user_info', '1')
  }

  return `${normalizedDomain}/#/set_token?${searchParams.toString()}`
}
