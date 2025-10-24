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

<template>
  <div class="external-services-page">
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
import { ref, provide } from 'vue'
import { storeToRefs } from 'pinia'
import { useRobotStore } from '@/stores/modules/robot'
import TopMeun from './components/top-meun.vue'
import WebAPP from './components/web-app.vue'
import EmbedWebsite from './components/embed-website.vue'
import WeChatOfficialAccount from './components/wechat-official-account.vue'
import WeChatMiniProgram from './components/wechat-mini-program.vue'
import WeChatCustomerService from './components/wechat-customer-service.vue'

const robotStore = useRobotStore()

const { robotInfo } = storeToRefs(robotStore)
const { getRobot } = robotStore

const tabComponents = {
  WebAPP,
  EmbedWebsite,
  WeChatOfficialAccount,
  WeChatMiniProgram,
  WeChatCustomerService
}

const activeLocalKey = '/robot/config/external-services/activeKey'

const activeMenuKey = ref(localStorage.getItem(activeLocalKey) || 'WebAPP')

const changeMenu = (item) => {
  activeMenuKey.value = item.id
  localStorage.setItem(activeLocalKey, activeMenuKey.value)
}

provide('robotInfo', {
  robotInfo: robotInfo.value,
  getRobot
})
</script>
