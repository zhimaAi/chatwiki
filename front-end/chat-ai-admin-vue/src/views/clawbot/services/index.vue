<template>
  <div class="external-services-page" v-if="ready">
    <div class="page-header">
      <TopMeun @change="changeMenu" v-model:value="activeMenuKey" />
    </div>
    <div class="page-body">
      <div class="scroll-box">
        <component :is="tabComponents[activeMenuKey]"></component>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, provide, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { useRobotStore } from '@/stores/modules/robot'

// 复用 robot 对外服务子组件
import TopMeun from '@/views/robot/robot-config/external-service/components/top-meun.vue'
import WebAPP from '@/views/robot/robot-config/external-service/components/web-app.vue'
import EmbedWebsite from '@/views/robot/robot-config/external-service/components/embed-website.vue'
import WeChatOfficialAccount from '@/views/robot/robot-config/external-service/components/wechat-official-account.vue'
import WeChatMiniProgram from '@/views/robot/robot-config/external-service/components/wechat-mini-program.vue'
import WeChatCustomerService from '@/views/robot/robot-config/external-service/components/wechat-customer-service.vue'
import WeComRobot from '@/views/robot/robot-config/external-service/components/wecom-robot.vue'
import FeishuRobot from '@/views/robot/robot-config/external-service/components/feishu-robot.vue'
import DingDingRobot from '@/views/robot/robot-config/external-service/components/dingding-robot.vue'

const clawbotStore = useClawbotStore()
const robotStore = useRobotStore()

const ready = ref(false)

// 桥接：将当前助手数据加载到 useRobotStore（子组件依赖它）
onMounted(async () => {
  await robotStore.getRobot(clawbotStore.currentAssistantId)
  ready.value = true
})

const { robotInfo } = storeToRefs(robotStore)
const { getRobot } = robotStore

const tabComponents = {
  WebAPP,
  EmbedWebsite,
  WeChatOfficialAccount,
  WeChatMiniProgram,
  WeChatCustomerService,
  WeComRobot,
  FeishuRobot,
  DingDingRobot
}

const activeLocalKey = '/clawbot/services/activeKey'
const activeMenuKey = ref(localStorage.getItem(activeLocalKey) || 'WebAPP')

const changeMenu = (item) => {
  activeMenuKey.value = item.id
  localStorage.setItem(activeLocalKey, activeMenuKey.value)
}

// 为 inject('robotInfo') 的子组件提供数据
provide('robotInfo', { robotInfo: robotInfo.value, getRobot })
</script>

<style lang="less" scoped>
.external-services-page {
  display: flex;
  flex-flow: column nowrap;
  width: 100%;
  height: 100%;
  .page-body {
    flex: 1;
    padding: 16px 24px;
    overflow: hidden;
  }
  .scroll-box {
    height: 100%;
    overflow-y: auto;
  }
}
</style>
