import { defineStore } from 'pinia'

export const useGlobalStore = defineStore('global', {
  state: () => {
    return {
      hideLayoutTopAndBottom: true,
    }
  },
  getters: {},
  actions: {
    setHideLayoutTopAndBottom(val) {
      this.hideLayoutTopAndBottom = val
    }
  }
})
