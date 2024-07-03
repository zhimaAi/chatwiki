import { defineStore } from 'pinia'
import { checkPermission } from '@/api/manage/index'
import { store } from '../index'

function generatePermissions(data) {
  var keys = []
  function getKey(list) {
    for (let i = 0; i < list.length; i++) {
      keys.push(list[i].label)

      if (list[i].sub_button_list && list[i].sub_button_list.length > 0) {
        list[i].sub_button_list.forEach((item) => {
          keys.push(...item.route_list)
        })
      }

      if (list[i].navs && list[i].navs.length > 0) {
        getKey(list[i].navs)
      }
    }
  }

  getKey(data)
  return keys
}

function generateMenus(data) {
  return data
}

export const usePermissionStore = defineStore('permission', {
  state: () => {
    return {
      permissionList: [],
      menuList: [],
      role_permission: [],
      user_roles: null,
      menus: [],
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
    },

  },
  actions: {
    async getPermissionList() {
      this.permissionList = ['admin']
      return Promise.resolve({ res: 0, data: [] })
    },
    setRolePermission(data) {
      this.role_permission = data.role_permission || [];
      this.user_roles = data.user_roles;
      this.menus = data.menu || [];
    },
    async checkPermission() {
      const res = await checkPermission()
      this.setRolePermission(res.data)
      return res.data;
    },
  },
  persist: true
})

export const usePermissionStoreWithOut = () => {
  return usePermissionStore(store)
}
