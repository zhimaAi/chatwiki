import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import { getTriggerConfigList } from '@/api/trigger/index'

export const useWorkflowStore = defineStore('workflow', () => {
  const triggerList = ref([])

  const getTriggerList = async (isRefresh) => {
    try {
      if (triggerList.value.length === 0 || isRefresh) {
        const res = await getTriggerConfigList()
        triggerList.value = res.data
      }

      return triggerList.value
    } catch (error) {
      console.log(error)
    }
  }

  return {
    triggerList,
    getTriggerList
  }
})
