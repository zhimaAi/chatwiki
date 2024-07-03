import { defineStore } from 'pinia'
import { getCompany } from '@/api/user'

export const useCompanyStore = defineStore('company', {
  state: () => {
    return {
      name: '',
      id: '',
      companyInfo: undefined,
    }
  },
  getters: {

  },
  actions: {
    async getCompanyInfo() {
      const res = await getCompany()

      if (!res) {
        return Promise.reject(res)
      }
      
      this.setCompanyInfo(res.data)

      return res
    },

    setCompanyInfo(data) {
      this.companyInfo = data;
      this.name = data ? data.name : '';
      this.id = data ? data.id : '';
    },
  },
  persist: true
})
