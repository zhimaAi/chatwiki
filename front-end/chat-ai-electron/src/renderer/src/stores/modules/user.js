import { defineStore } from 'pinia'
import { store } from '../index'
import router from '@/router'
import { login, getCompany } from '@/api/user'
import { appTitle } from '@/constants/index'

export const useUserStore = defineStore('user', {
  persist: true,
  state: () => {
    return {
      companyInfo: null,
      userInfo: undefined,
      token: undefined
    }
  },
  getters: {
    getToken() {
      return this.token ? this.token.token : null
    }
  },
  actions: {
    getCompany() {
      return getCompany().then((res) => {
        this.companyInfo = res.data
        // 设置title
        if (this.companyInfo && this.companyInfo.name) {
          document.title = appTitle + ' | ' + this.companyInfo.name
        }

        return res
      })
    },
    setToken(token) {
      this.token = token
    },
    // set
    setUserInfo(userInfo) {
      this.userInfo = userInfo
    },
    async login(params) {
      try {
        const { username, password } = params

        const res = await login(username, password)

        this.setToken({
          token: res.data.token,
          exp: res.data.exp,
          ttl: res.data.ttl
        })

        delete res.data.token
        delete res.data.exp
        delete res.data.ttl

        this.setUserInfo(res.data)

        return res.data
      } catch (error) {
        return Promise.reject(error)
      }
    },
    reset(goLogin) {
      this.setToken(undefined)
      this.setUserInfo(undefined)

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
    }
  }
})

export const useUserStoreWithOut = () => {
  return useUserStore(store)
}
