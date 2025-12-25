import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import { getTriggerConfigList, getTriggerOfficialMessage } from '@/api/trigger/index'

export const useWorkflowStore = defineStore('workflow', () => {
  const triggerList = ref([])
  const officialList = ref([])
  const triggerOfficialList = ref([])

  const getTriggerList = async (isRefresh) => {
    try {
      if (triggerList.value.length === 0 || isRefresh) {
        const res = await getTriggerConfigList()
        triggerList.value = res.data.map((item) => {
          item.subMenus = []
          item.expend = false
          if (item.trigger_type == 4) {
            item.subMenus = [
              {
                title: '私信消息',
                value: 'message'
              },
              {
                title: '关注/取消关注事件',
                value: 'subscribe_unsubscribe'
              },
              {
                title: '扫描带参数二维码事件',
                value: 'qrcode_scan'
              },
              {
                title: '自定义菜单事件',
                value: 'menu_click'
              }
            ]
          }
          return item
        })
      }

      return triggerList.value
    } catch (error) {
      console.log(error)
    }
  }


  const getTriggerOfficialMsg = (robot_key) => {
    getTriggerOfficialMessage({ robot_key }).then((res) => {
      triggerOfficialList.value = res.data.messages || []
      officialList.value = res.data.apps || []
    })
  }

  return {
    triggerList,
    getTriggerList,
    officialList,
    getTriggerOfficialMsg,
    triggerOfficialList,
  }
})
