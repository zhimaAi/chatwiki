import router from './router'
import { useUserStoreWithOut } from '@/stores/modules/user'
import { usePermissionStoreWithOut } from '@/stores/modules/permission'
import { useCompanyStore } from '@/stores/modules/company'
import { NO_REDIRECT_WHITE_LIST } from '@/constants'
import { checkSystemPermisission } from '@/utils/permission.js'
// import { checkPermi } from '@/utils/permission'
// checkPermi(['yun/yunAdmin/ChannelCreate'])

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
  let { permissionList, getPermissionList, checkPermission } = permissionStore

  let needGetPermissionRoutes = [
    '/robot/list',
    '/library/list',
    '/user/model',
    '/user/clientDownload'
  ]
  if(to.path == '/set_token'){
    next()
    return
  }

  if (userInfo) {
    if (to.path === '/login') {
      next({ path: '/', query: to.query })
    } else {
      if (needGetPermissionRoutes.includes(to.path)) {
        await checkPermission()
      }
      // 权限显示逻辑更换了，不用以前的逻辑了
      // if (!checkRouterPermisission(to.path)) {
      //   toAuthorityPage(to, from, next)
      //   return
      // }
      // 系统管理里面，menu根据顺序，如果有权限的，显示第一个menu
      const name = await checkSystemPermisission(to)
      if (name) {
        if (name === 'AccountManage') {
          next(`/user/account`)
          return
        } else if (name === 'TokenManage') {
          next(`/user/usetoken`)
          return
        } else if (name === 'TeamManage') {
          next(`/user/manage`)
          return
        } else if (name === 'CompanyManage') {
          next(`/user/enterprise`)
          return
        } else if (name === 'ClientSideManage') {
          next(`/user/clientDownload`)
          return
        }
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
