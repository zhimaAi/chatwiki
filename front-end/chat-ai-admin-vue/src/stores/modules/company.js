import { defineStore } from 'pinia'
import { getCompany } from '@/api/user'

export const useCompanyStore = defineStore('company', {
  state: () => {
    return {
      name: '',
      id: '',
      companyInfo: undefined,
      ali_ocr_switch: '2',
      ali_ocr_key: '',
      ali_ocr_secret: ''
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
      this.ali_ocr_switch = data ? data.ali_ocr_switch : '';
      this.ali_ocr_key = data ? data.ali_ocr_key : '';
      this.ali_ocr_secret = data ? data.ali_ocr_secret : '2';
    },
  },
})
