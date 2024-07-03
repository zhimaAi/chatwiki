import { reactive } from 'vue'
import { defineStore } from 'pinia'
import { store } from '../index'
import { DEFAULT_USER_AVATAR } from '@/constants/index'

export type User = {
  avater: string
}

const defaultUserStore: User = {
  avater: DEFAULT_USER_AVATAR
}

export const useUserStore = defineStore('user', () => {
  const userInfo = reactive({ ...defaultUserStore })

  return userInfo
})

export const useUserStoreWithOut = () => {
  return useUserStore(store)
}
