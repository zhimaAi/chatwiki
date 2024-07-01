import { NO_REDIRECT_WHITE_LIST, appTitle } from '@/constants/index'
import { useUserStoreWithOut } from '@/stores/modules/user'

// 辅助函数：检查是否需要重定向
function shouldRedirect(to) {
  return NO_REDIRECT_WHITE_LIST.some((n) => n === to.name)
}

function setTitle(title) {
  if (title) {
    document.title = appTitle + ' | ' + title
  }
}

export async function createRouterGuards(router) {
  const useUserStore = useUserStoreWithOut()
  // 导航守卫
  router.beforeEach(async (to, from, next) => {
    if (shouldRedirect(to)) {
      next()
    } else {
      if (!useUserStore.companyInfo) {
        const res = await useUserStore.getCompany()

        if (res.data && res.data.name) {
          setTitle(res.data.name)
        }
        next()
      } else {
        setTitle(useUserStore.companyInfo.name)
        next()
      }
    }
  })
}
