<style lang="less" scoped>
.robot-main-container {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.robot-page-layout {
  display: flex;
  height: 100%;
  flex: 1;
  width: 100%;
  border: 2px;
  overflow: hidden;
  background-color: #fff;

  .scroll-box {
    height: 100%;
    overflow-y: auto;
  }

  .layout-left {
    width: 255px;
    height: 100%;
    border-right: 1px solid #f2f4f7;
    border-top: 1px solid #f2f4f7;
    border-radius: 6px;
  }

  .layout-body {
    flex: 1;
    height: 100%;
    overflow: hidden;
  }
}
</style>

<template>
  <div class="robot-main-container">
    <div class="robot-page-layout">
      <div class="layout-left" v-if="isShowLeft">
        <TopHeader />
        <leftMenu @changeMenu="changeMenu" :robotInfo="robotInfo" />
      </div>
      <div class="layout-body">
        <router-view></router-view>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter, useRoute } from 'vue-router'
import { useRobotStore } from '@/stores/modules/robot'
import leftMenu from './components/left-menu.vue'
import { getRobotPermission } from '@/utils/permission'

import TopHeader from './components/top-header.vue'

export default defineComponent({
  name: 'robotPage',
  components: {
    leftMenu,
    TopHeader
  },
  async beforeRouteEnter(to, from, next) {
    const robotStore = useRobotStore()
    const { getRobot, robotInfo, getRobotLists, getGroupList } = robotStore
    await getRobot(to.query.id)
    getRobotLists()
    getGroupList()
    let key = getRobotPermission(to.query.id)
    if (key == 0) {
      next(`/no-permission`)
      return
    }
    if (key == 1) {
      // 只用查看权限
      window.location.href = `${robotInfo.h5_domain}/#/chat/pc?robot_key=${robotInfo.robot_key}`
      next(`/robot/list`)
      return
    }
    let workFlowFilterRouter = ['basicConfig', 'libraryConfig', 'skillConfig']
    if (robotInfo.application_type == 1 && workFlowFilterRouter.includes(to.name)) {
      next(`/robot/config/workflow?id=${robotInfo.id}&robot_key=${robotInfo.robot_key}`)
      return
    }
    if (robotInfo.application_type == 0 && to.name == 'robotWorkflow') {
      next(`/robot/config/basic-config?id=${robotInfo.id}&robot_key=${robotInfo.robot_key}`)
      return
    }
    next()
  },

  setup() {
    const router = useRouter()
    const robotStore = useRobotStore()
    const { robotInfo } = storeToRefs(robotStore)
    // 基本配置
    const isShowLeft = computed(() => {
      return useRoute().name != 'robotWorkflow'
    })
    const changeMenu = (item) => {
      router.push({
        path: item.path,
        query: {
          robot_key: robotInfo.value.robot_key,
          id: robotInfo.value.id
        }
      })
    }

    return {
      robotInfo,
      changeMenu,
      isShowLeft
    }
  }
})
</script>
