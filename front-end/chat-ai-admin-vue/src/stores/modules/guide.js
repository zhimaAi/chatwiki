import { defineStore } from 'pinia'
import { useGuideProcess } from '@/api/guide/index'

function getGuideProcessVal(list) {
  let finish_num = 0
  let total_num = 0
  list.forEach((item) => {
    item.steps.forEach((it) => {
      total_num++
      if (it.is_finish == 1) {
        finish_num++
      }
    })
  })

  return ((finish_num / total_num) * 100).toFixed()
}
export const useGuideStore = defineStore('guide', {
  state: () => {
    return {
      process_list: [],
      total_process: 0
    }
  },
  getters: {},
  actions: {
    async getUseGuideProcess() {
      try {
        const res = await useGuideProcess()
        this.process_list = res.data.process_list || []
        this.total_process = getGuideProcessVal(this.process_list)
        return res.data
      } catch (error) {
        return Promise.reject(error)
      }
    }
  }
})
