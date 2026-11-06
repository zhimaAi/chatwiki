<template>
  <div />
</template>

<script setup>
import { onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getRobotList } from '@/api/robot'
import { useCompanyStore } from '@/stores/modules/company'
import { SYSTEM_MODES, useSystemModeStore } from '@/stores/modules/system-mode'
import { CUSTOMER_SERVICE_LAYOUT } from '@/utils/clawbot-layout'

const route = useRoute()
const router = useRouter()
const companyStore = useCompanyStore()
const modeStore = useSystemModeStore()
let isActive = true

const resolveHome = async () => {
  // 首页只负责分流；未保存模式或模式读取失败时，沿用高级模式入口及原始参数。
  let target = { path: '/robot/list', query: route.query, hash: route.hash }

  try {
    const isModeReady = await modeStore.loadForScope({ companyId: companyStore.id })
    if (isModeReady && modeStore.currentMode !== SYSTEM_MODES.ADVANCED) {
      const isCustomerService = modeStore.currentMode === SYSTEM_MODES.AI_SERVICE
      const applicationType = isCustomerService ? 0 : 2
      target = {
        path: isCustomerService ? '/customer-service/chat' : '/clawbot/chat',
        // Agent 走新版客服布局；独立域靠 system_mode 识别，不能只看本地模式。
        query: isCustomerService ? {} : {
          system_mode: SYSTEM_MODES.AGENT,
          layout: CUSTOMER_SERVICE_LAYOUT
        }
      }

      const response = await getRobotList(
        isCustomerService ? {} : { filter_claw: false, application_type: applicationType }
      )
      if (!isActive) return

      // 保留首个应用的入口参数；列表或应用信息异常由目标页面处理。
      const firstApplication = Array.isArray(response?.data)
        ? response.data.find((item) => Number(item?.application_type) === applicationType)
        : null
      if (firstApplication?.id && firstApplication?.robot_key) {
        target.query.id = firstApplication.id
        target.query.robot_key = firstApplication.robot_key
      }
    }
  } catch (error) {
    console.error('读取首页入口信息失败', error)
  }

  if (!isActive) return
  await router.replace(target)
}

resolveHome().catch((error) => {
  console.error('首页跳转失败', error)
})
onBeforeUnmount(() => {
  // 用户离开首页入口后，未完成的列表请求不能再覆盖他的新导航。
  isActive = false
})
</script>
