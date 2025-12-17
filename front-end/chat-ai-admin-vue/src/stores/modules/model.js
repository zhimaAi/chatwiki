import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getModelConfigList } from '@/api/model/index'

export const useModelStore = defineStore('model', () => {
  const allModelList = ref([])

  const getAllmodelList = async () => {
    try {
      if (allModelList.value.length == 0) {
        const res = await getModelConfigList()
        allModelList.value = res.data
      }
      return allModelList.value
    } catch (error) {
      console.log(error)
    }
  }

  return {
    allModelList,
    getAllmodelList
  }
})
