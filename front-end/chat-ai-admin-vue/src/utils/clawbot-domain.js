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
  adminDomain = '',
  agentDomain = '',
  currentOrigin = ''
} = {}) {
  // set_token 只是跨域登录态落地页，继续参与切域会造成重定向循环。
  if (toPath === '/set_token') {
    return ''
  }

  const isToClawbot = isClawbotPath(toPath)
  const isFromClawbot = isClawbotPath(fromPath)

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
