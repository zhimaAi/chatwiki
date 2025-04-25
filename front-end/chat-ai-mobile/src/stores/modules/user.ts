import { defineStore } from 'pinia'
import { store } from '../index'
import { loginApi, getUserInfo } from '@/api/user'

// 定义明确的类型接口
type Token = {
  token: string
  exp: number
  ttl: number
}

type UserInfo = {
  // 补充完整的用户信息类型
  avater: string
  [key: string]: any
}

type LoginInfo = any // 根据实际情况定义登录信息类型

interface UserState {
  loginStatus: boolean
  userInfo?: UserInfo
  token?: Token
  loginInfo?: LoginInfo
}

export type User = {
  avater: string
}

export const useUserStore = defineStore('user', {
  state: (): UserState => {
    return {
      loginStatus: false, // 登录状态 用于渲染退出登录按钮
      userInfo: undefined,
      token: undefined,
      loginInfo: undefined
    }
  },
  getters: {
    getToken(): string | null {
      return this.token ? this.token.token : null
    },
    getLoginStatus(): boolean {
      return this.loginStatus
    }
  },
  actions: {
    setToken(token) {
      this.token = token
    },
    setUserInfo(userInfo) {
      this.userInfo = userInfo
    },
    setLoginStatus(loginStatus) {
      this.loginStatus = loginStatus
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
        this.setUserInfo(res.data)
        this.setLoginStatus(true)

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
    reset() {
      this.setToken(undefined)
      this.setUserInfo(undefined)
      this.setLoginStatus(false)
    },
    logout() {
      this.reset()
    },
    setLoginInfo(loginInfo) {
      this.loginInfo = loginInfo
    }
  },
  persist: true
})

export const useUserStoreWithOut = () => {
  return useUserStore(store)
}
