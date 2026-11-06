import router from './router'
import { useUserStoreWithOut } from '@/stores/modules/user'
import { usePermissionStoreWithOut } from '@/stores/modules/permission'
import { useCompanyStore } from '@/stores/modules/company'
import { NO_REDIRECT_WHITE_LIST } from '@/constants'
import { checkRole, checkSystemPermisission } from '@/utils/permission.js'
import { useSystemModeStore } from '@/stores/modules/system-mode'
import { buildSetTokenRedirectUrl, getCrossDomainTarget } from '@/utils/clawbot-domain'
import { isCustomerServiceClawbotRoute } from '@/utils/clawbot-layout'
import { useLocaleStoreWithOut } from '@/stores/modules/locale'
import { useLocale } from '@/hooks/web/useLocale'
import { i18n } from '@/locales'

function toLogin(to, from, next) {
  if (to.path === '/login') {
    next()
  } else if (isCustomerServiceClawbotRoute(to)) {
    // 新版深链需要恢复应用、菜单及布局模式；登录页按现有约定解码 redirect。
    next({ path: '/login', query: { redirect: encodeURIComponent(to.fullPath) } })
  } else {
    next(`/login?redirect=${to.path}`)
  }
}

// function toAuthorityPage(to, from, next) {
//   if (to.path.includes('/user') && to.path != '/user/account') {
//     next(`/user/account`)
//   }
//   next(`/authority/index`)
// }

function setTitle(to, companyInfo) {
  const appName = companyInfo?.name || (i18n && i18n.global ? i18n.global.t('utils.index.app_name_default') : '开源大模型企业知识库问答系统')
  let str = `Chatwiki ` + appName
  let title = to.meta.title
  if (typeof title === 'string' && title.startsWith('routes.basic.')) {
    title = i18n && i18n.global ? i18n.global.t(title) : title
  }
  if (title) {
    document.title = title + ' - ' + str
  } else {
    document.title = str
  }
}

async function needsModeSelection(to, userStore, companyStore) {
  // 新版 Agent 无应用参数时显示应用选择入口，不能被独立入口的模式选择截走。
  if (isCustomerServiceClawbotRoute(to)) return false
  const landingPaths = ['/', '/robot/list', '/customer-service/chat', '/customer-service/session-record', '/clawbot/chat']
  // 只拦截默认入口，保留设置页、工作流和指定应用的业务深链。
  if (!landingPaths.includes(to.path) || to.query.id || to.query.robot_key) return false
  if (!userStore.user_is_admin && !checkRole(['CompanyManage'])) return false
  const modeStore = useSystemModeStore()
  return await modeStore.loadForScope({ companyId: companyStore.id }) && modeStore.needsModeSelection
}


router.beforeEach(async (to, from, next) => {
  // 处理 URL 中的 lang 参数
  if (to.query.lang) {
    const localeStore = useLocaleStoreWithOut()
    const { changeLocale } = useLocale()
    
    const validLocales = localeStore.getLocaleMap.map(v => v.lang)
    
    if (validLocales.includes(to.query.lang)) {
      await changeLocale(to.query.lang)
      
      // 移除 URL 中的 lang 参数
      const query = { ...to.query }
      delete query.lang
      next({ path: to.path, query })
      return
    }
  }

  const companyStore = useCompanyStore()
  const { companyInfo, getCompanyInfo } = companyStore
  if (!companyInfo) {
    await getCompanyInfo()
  }
  const userStore = useUserStoreWithOut()
  if (userStore.getToken && ['/login', '/set_token', '/chatclaw/login'].includes(from.path)) {
    // 登录或跨域落地后重读当前企业，不能使用登录页缓存的上一企业配置。
    await getCompanyInfo()
  }
  setTitle(to, companyStore.companyInfo)
  const permissionStore = usePermissionStoreWithOut()
  // 不是白名单的路由
  const notWhitePath = NO_REDIRECT_WHITE_LIST.indexOf(to.path) !== -1

  let { userInfo, getUserInfo } = userStore
  let { getPermissionList, checkPermission } = permissionStore

  const crossDomainTarget = getCrossDomainTarget({
    toPath: to.path,
    fromPath: from.path,
    toQuery: to.query,
    fromQuery: from.query,
    adminDomain: companyStore.admin_domain,
    agentDomain: companyStore.agent_domain,
    currentOrigin: window.location.origin
  })

  if (crossDomainTarget && userStore.getToken && userStore.userInfo?.user_id) {
    // admin/clawbot 分域部署时，通过 set_token 把当前登录态带到目标域，再回跳原始页面。
    const targetUrl = buildSetTokenRedirectUrl({
      domain: crossDomainTarget,
      redirectUrl: to.fullPath,
      token: userStore.getToken,
      exp: userStore.token?.exp,
      ttl: userStore.token?.ttl,
      userId: userStore.userInfo.user_id,
      userName: userStore.userInfo.user_name,
      refreshUserInfo: true
    })

    window.location.replace(targetUrl)
    return
  }

  if (to.path == '/set_token' || to.path == '/chatclaw/login') {
    next()
    return
  }

  if (userInfo) {
    if (to.path === '/login') {
      next({ path: '/', query: to.query })
    } else {
      await checkPermission()
      if (await needsModeSelection(to, userStore, companyStore)) { next('/mode-select'); return }
      
      const name = checkSystemPermisission(to)
      if (name === 'AccountManage') {
          next(`/user/account`)
          return
      }
      next()
    }
  } else {
    if (notWhitePath) {
      next()
    } else {
      try {
        let res1 = await getUserInfo()

        if (!res1) {
          toLogin(to, from, next)
          return
        }

        let res2 = await getPermissionList()

        if (!res2) {
          toLogin(to, from, next)
          return
        }

        await checkPermission()
        if (await needsModeSelection(to, userStore, companyStore)) { next('/mode-select'); return }
        next()
      } catch (e) {
        toLogin(to, from, next)
      }
    }
  }
})
