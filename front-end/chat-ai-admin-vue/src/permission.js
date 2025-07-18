import router from './router'
import { useUserStoreWithOut } from '@/stores/modules/user'
import { usePermissionStoreWithOut } from '@/stores/modules/permission'
import { useCompanyStore } from '@/stores/modules/company'
import { NO_REDIRECT_WHITE_LIST } from '@/constants'
import { checkSystemPermisission } from '@/utils/permission.js'

function toLogin(to, from, next) {
  if (to.path === '/login') {
    next()
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
  let str = `Chatwiki ` + `${companyInfo?.name || '开源大模型企业知识库问答系统'}`
  if (to.meta.title) {
    document.title = to.meta.title + ' - ' + str
  } else {
    document.title = str
  }
}


router.beforeEach(async (to, from, next) => {
  const companyStore = useCompanyStore()
  const { companyInfo, getCompanyInfo } = companyStore
  if (!companyInfo) {
    await getCompanyInfo()
  }
  setTitle(to, companyInfo)

  const userStore = useUserStoreWithOut()
  const permissionStore = usePermissionStoreWithOut()
  // 不是白名单的路由
  const notWhitePath = NO_REDIRECT_WHITE_LIST.indexOf(to.path) !== -1

  let { userInfo, getUserInfo } = userStore
  let { getPermissionList, checkPermission } = permissionStore

  if (to.path == '/set_token') {
    next()
    return
  }

  if (userInfo) {
    if (to.path === '/login') {
      next({ path: '/', query: to.query })
    } else {
      await checkPermission()
      
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

        next()
      } catch (e) {
        toLogin(to, from, next)
      }
    }
  }
})
