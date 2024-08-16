import { defineStore } from 'pinia'
import { getFormInfo } from '@/api/database'

export const useDatabaseStore = defineStore('database', {
  state: () => {
    return {
      databaseInfo: {},
    }
  },
  getters: {

  },
  actions: {
    async getDatabaseInfo(data) {
      const res = await getFormInfo(data)

      if (!res) {
        return Promise.reject(res)
      }

      this.setDatabaseInfo(res.data)

      return res
    },

    setDatabaseInfo(data) {
      this.databaseInfo = data;
    },
  },
})
