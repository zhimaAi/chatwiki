import { defineStore } from 'pinia'
import { checkPermission } from '@/api/manage/index'
import { store } from '../index'

export const usePermissionStore = defineStore('permission', {
  state: () => {
    return {
      permissionList: [],
      menuList: [],
      role_permission: [],
      user_roles: null,
      role_type: null,
      menus: []
    }
  },
  getters: {
    permissionMap: (state) => {
      return state.permissionList.reduce((acc, cur) => {
        // 把/替换成_
        let key = cur.replace(/\//g, '_')

        acc[key] = true

        return acc
      }, {})
    }
  },
  actions: {
    async getPermissionList() {
      this.permissionList = ['admin']
      return Promise.resolve({ res: 0, data: [] })
    },
    setRolePermission(data) {
      this.role_permission = data.role_permission || []
      this.user_roles = data.user_roles
      this.role_type = data.role_type
      this.menus = data.menu || []
    },
    async checkPermission() {
      const res = await checkPermission()
      this.setRolePermission(res.data)
      return res.data
    }
  },
  persist: true
})

export const usePermissionStoreWithOut = () => {
  return usePermissionStore(store)
}
