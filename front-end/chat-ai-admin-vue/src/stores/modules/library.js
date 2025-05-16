import { defineStore } from 'pinia'
import { getLibraryInfo } from '@/api/library'

export const useLibraryStore = defineStore('library', {
  state: () => {
    return {
      library_name: '',
      avatar: '',
      library_intro: '',
      robot_nums: 0,
      graph_switch: 0
    }
  },
  getters: {},
  actions: {
    async getLibraryInfo(id){
      const res = await getLibraryInfo({ id: id })

      this.library_name = res.data.library_name
      this.avatar = res.data.avatar
      this.library_intro = res.data.library_intro
      this.robot_nums = res.data.robot_nums
      this.graph_switch = res.data.graph_switch * 1
    },
    changeGraphSwitch(val){
      this.graph_switch = val
    }
  },
})

