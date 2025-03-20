import { defineStore } from 'pinia'
import { store } from '../index'
import { getLibDocPartner } from '@/api/public-library'

export const usePublicLibraryStore = defineStore('publicLibrary', {
  state: () => {
    return {
      library_key: '',
      library_id: '',
      operate_rights: '' // 4：管理权限  2：编辑权限 0:无
    }
  },
  getters: {},
  actions: {
    async getLibDocPartner(params) {
      return await getLibDocPartner(params).then((res) => {
        this.library_id = res.data.library_id
        this.operate_rights = res.data.operate_rights
        this.library_key = params.library_key
        return res
      })
    },
    checkPermission(permissions) {
      if (!permissions || permissions.includes('*') || permissions.length === 0) {
        return true
      }
      return permissions.includes(this.operate_rights)
    }
  }
})

export const usePublicLibraryStoreWithOut = () => {
  return usePublicLibraryStore(store)
}
