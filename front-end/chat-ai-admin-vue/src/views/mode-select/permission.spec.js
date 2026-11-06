import { beforeEach, describe, expect, it, vi } from 'vitest'

const state = vi.hoisted(() => ({ guard: null, user: null, company: null, mode: null }))
const api = vi.hoisted(() => ({ company: vi.fn(), user: vi.fn(), permissions: vi.fn(), listPermissions: vi.fn(), loadMode: vi.fn(), canManage: vi.fn() }))
vi.mock('@/router', () => ({ default: { beforeEach: guard => { state.guard = guard } } }))
vi.mock('@/stores/modules/user', () => ({ useUserStoreWithOut: () => state.user }))
vi.mock('@/stores/modules/company', () => ({ useCompanyStore: () => state.company }))
vi.mock('@/stores/modules/system-mode', () => ({ useSystemModeStore: () => state.mode }))
vi.mock('@/stores/modules/permission', () => ({ usePermissionStoreWithOut: () => ({ checkPermission: api.permissions, getPermissionList: api.listPermissions }) }))
vi.mock('@/constants', () => ({ NO_REDIRECT_WHITE_LIST: ['/login', '/privacy_policy'] }))
vi.mock('@/utils/permission.js', () => ({ checkRole: api.canManage, checkSystemPermisission: () => undefined }))
vi.mock('@/utils/clawbot-domain', () => ({ getCrossDomainTarget: () => '', buildSetTokenRedirectUrl: vi.fn() }))
vi.mock('@/stores/modules/locale', () => ({ useLocaleStoreWithOut: () => ({ getLocaleMap: [] }) }))
vi.mock('@/hooks/web/useLocale', () => ({ useLocale: () => ({}) }))
vi.mock('@/locales', () => ({ i18n: { global: { t: key => key } } }))

import '@/permission'

async function navigate(path = '/', query = {}, from = '/') {
  const next = vi.fn()
  await state.guard({ path, query, fullPath: path, meta: {} }, { path: from }, next)
  return next
}

describe('first-mode navigation', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    state.user = { getToken: 'session', userInfo: { user_id: 7, admin_user_id: 7 }, user_is_admin: true, getUserInfo: api.user }
    state.company = { id: 1, companyInfo: { id: 1 }, getCompanyInfo: api.company }
    state.mode = { needsModeSelection: true, loadForScope: api.loadMode }
    api.permissions.mockResolvedValue({})
    api.listPermissions.mockResolvedValue({ res: 0 })
    api.loadMode.mockResolvedValue(true)
    api.canManage.mockReturnValue(false)
    api.company.mockReset().mockResolvedValue({})
    api.user.mockReset().mockImplementation(async () => {
      state.user.userInfo = { user_id: 7, admin_user_id: 7 }
      return { data: state.user.userInfo }
    })
  })

  it.each(['/', '/robot/list', '/customer-service/chat', '/clawbot/chat'])('shows mode selection before the default landing page %s', async path => {
    const next = await navigate(path)
    expect(next).toHaveBeenCalledWith('/mode-select')
  })

  it('does not redirect old enterprises, invited members or explicit business deep links', async () => {
    state.mode.needsModeSelection = false
    expect(await navigate()).toHaveBeenCalledWith()
    state.mode.needsModeSelection = true
    state.user.user_is_admin = false
    expect(await navigate()).toHaveBeenCalledWith()
    state.user.user_is_admin = true
    for (const [path, query] of [['/robot/config/workflow', { id: 900 }], ['/user/enterprise', {}], ['/customer-service/chat', { id: 7, robot_key: 'customer' }], ['/mode-select', {}]]) {
      expect(await navigate(path, query)).toHaveBeenCalledWith()
    }
  })

  it('refreshes the authenticated enterprise before checking onboarding after login', async () => {
    api.company.mockImplementation(async () => { state.company.id = 2; state.company.companyInfo = { id: 2 } })
    const next = await navigate('/', {}, '/login')
    expect(api.company).toHaveBeenCalledTimes(1)
    expect(api.loadMode).toHaveBeenCalledWith({ companyId: 2 })
    expect(next).toHaveBeenCalledWith('/mode-select')
  })

  it('checks mode selection after restoring an authenticated user on refresh', async () => {
    state.user.userInfo = undefined
    const next = await navigate()
    expect(api.user).toHaveBeenCalledTimes(1)
    expect(api.permissions).toHaveBeenCalledTimes(1)
    expect(next).toHaveBeenCalledWith('/mode-select')
  })

  it('lets company managers complete onboarding while leaving unauthenticated routes under login control', async () => {
    state.user.user_is_admin = false
    api.canManage.mockReturnValue(true)
    expect(await navigate()).toHaveBeenCalledWith('/mode-select')
    state.user.userInfo = undefined
    state.user.getToken = null
    api.user.mockResolvedValueOnce(null)
    expect(await navigate()).toHaveBeenCalledWith('/login?redirect=/')
  })
})
