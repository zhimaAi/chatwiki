import { defineStore } from 'pinia'
import { store } from '../index'
import { Modal } from 'ant-design-vue'
import { loginApi, getUserInfo, getUnReadMessageTotal, refreshUserToken } from '@/api/user'
// import { getTokenCache, setTokenCache } from '@/storage/user'
import router from '@/router'
import { DEFAULT_USER_AVATAR } from '@/constants/index'

export const useUserStore = defineStore('user', {
  state: () => {
    return {
      userInfo: undefined,
      token: undefined,
      roleRouters: undefined,
      // 记住我
      rememberMe: true,
      loginInfo: undefined,
      myAppList: [],
      unReadMessageTotal: 0,
      isLayoutScroll: true,
      isShowResetPassModal: false,
      robot_nums: 0,
      guideLearningTips: false  // 是否显示引导弹窗
    }
  },
  getters: {
    getToken() {
      return this.token ? this.token.token : null
    },
    isAdmin() {
      return this.userInfo && this.userInfo.isAdministator
    },
    user_is_admin() {
      return this.userInfo && this.userInfo.user_id === this.userInfo.admin_user_id
    },
    user_id() {
      return this.userInfo ? this.userInfo.user_id : null
    },
    user_name() {
      return this.userInfo ? this.userInfo.user_name : null
    },
    avatar() {
      if (this.userInfo && this.userInfo.avatar) {
        return this.userInfo.avatar
      }

      return DEFAULT_USER_AVATAR
    },
    getRoleRouters() {
      return this.roleRouters
    },
    getRememberMe() {
      return this.rememberMe
    },
    getLoginInfo() {
      return this.loginInfo
    },
    getCurrentApp() {
      return (appId) => {
          return this.myAppList.find((item) => item.wechatapp_id === appId)
      }
    },
    getAppType() {
      return (appId) => {
          let app = this.myAppList.find((item) => item.wechatapp_id === appId)

          return app ? app.app_type : null
      }
    },
    getRobotNums () {
      return this.robot_nums
    },
    getGuideLearningTips() {
      return this.guideLearningTips
    }
  },
  actions: {
    setToken(token) {
      this.token = token
    },
    // set
    setUserInfo(userInfo) {
      this.userInfo = userInfo
    },
    setAvatar(url) {
      this.userInfo.avatar = url
    },
    setResetPassModal() {
      this.isShowResetPassModal = false
    },
    async login(params) {
      try {
        const { username, password } = params

        const res = await loginApi(username, password)

        this.setToken({
          token: res.data.token,
          exp: res.data.exp,
          ttl: res.data.ttl
        })
        this.isShowResetPassModal = true
        this.setUserInfo(res.data)

        return res.data
      } catch (error) {
        return Promise.reject(error)
      }
    },
    async refreshToken() {
      try {
        const res = await refreshUserToken()

        this.setToken({
          token: res.data.token,
          exp: res.data.exp,
          ttl: res.data.ttl
        })

        return res.data
      } catch (error) {
        return Promise.reject(error)
      }
    },
    async getUserInfo() {
      const res = await getUserInfo()

      if (!res) {
        return Promise.reject(res)
      }

      this.setUserInfo(res.data)

      return res
    },
    async getUnReadMessageTotal() {
      let query = router.currentRoute.value.query

      const res = await getUnReadMessageTotal({
          wechatapp_id: query.wechatapp_id
      })

      if (res) {
          this.unReadMessageTotal = res.data.totalCount
      }

      return res
    },
    setRoleRouters(roleRouters) {
      this.roleRouters = roleRouters
    },
    logoutConfirm() {
      Modal.confirm({
        type: 'warning',
        title: '温馨提示',
        content: '是否退出本系统？',
        onOk: () => {
          this.reset(true)
          // const res = loginOutApi()
          //   .then(() => {
          //     if (res) {
          //       this.reset(true)
          //     }
          //   })
          //   .catch(() => {})
        },
        onCancel() {}
      })
    },
    reset(goLogin) {
      this.setToken(undefined)
      this.setUserInfo(undefined)
      this.setRoleRouters([])

      if (goLogin) {
        // 直接回登陆页
        router.replace('/login')
      } else {
        // 回登陆页带上当前路由地址
        router.replace({
          path: '/login',
          query: {
            redirect: encodeURIComponent(router.currentRoute.value.fullPath)
          }
        })
      }
    },
    logout() {
      this.reset()
    },
    setRememberMe(rememberMe) {
      this.rememberMe = rememberMe
    },
    setLoginInfo(loginInfo) {
      this.loginInfo = loginInfo
    },
    setRobotNums (robot_nums) {
      this.robot_nums = robot_nums
    },
    setGuideLearningTips(guideLearningTips) {
      this.guideLearningTips = guideLearningTips
    }
  },
  persist: true
})

export const useUserStoreWithOut = () => {
  return useUserStore(store)
}
