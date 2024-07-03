<style lang="less" scoped>
.robot-page-layout {
  display: flex;
  height: 100%;
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

    .robot-name-box {
      display: flex;
      align-items: center;
      padding: 24px 24px 16px 24px;

      .robot-avatar {
        width: 20px;
        height: 20px;
        margin-right: 8px;
        border-radius: 2px;
      }

      .robot-name {
        flex: 1;
        line-height: 22px;
        font-size: 14px;
        font-weight: 600;
        color: #262626;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
    }
  }

  .layout-body {
    flex: 1;
    height: 100%;
    overflow: hidden;
  }
}
</style>

<template>
  <div class="robot-page-layout">
    <div class="layout-left">
      <div class="robot-name-box">
        <img class="robot-avatar" :src="robotInfo.robot_avatar_url" alt="" />
        <span class="robot-name">{{ robotInfo.robot_name }}</span>
      </div>
      <leftMenu @changeMenu="changeMenu" :robotInfo="robotInfo" />
    </div>
    <div class="layout-body">
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useRobotStore } from '@/stores/modules/robot'
import leftMenu from './components/left-menu.vue'

export default defineComponent({
  name: 'robotPage',
  components: {
    leftMenu
  },
  async beforeRouteEnter(to, from, next) {
    const robotStore = useRobotStore()
    const { getRobot } = robotStore

    await getRobot(to.query.id)

    next()
  },
  setup() {
    const router = useRouter()
    const robotStore = useRobotStore()
    const { robotInfo } = storeToRefs(robotStore)

    // 基本配置
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
      changeMenu
    }
  }
})
</script>
