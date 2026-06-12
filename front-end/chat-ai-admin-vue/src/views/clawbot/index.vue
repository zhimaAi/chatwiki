<template>
  <div class="clawbot-layout" v-if="clawbotStore.isReady">
    <LeftSidebar />
    <div class="clawbot-layout-container">
      <router-view></router-view>
    </div>
  </div>
  <ClawbotLoading v-else />
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import LeftSidebar from './components/left-sidebar.vue'
import ClawbotLoading from './components/clawbot-loading.vue'
import { useClawbotStore } from '@/stores/modules/clawbot'

const route = useRoute()
const router = useRouter()
const clawbotStore = useClawbotStore()
const isSyncingRouteQuery = ref(false)

const getQueryValue = (value) => {
  if (Array.isArray(value)) {
    return value[0] || ''
  }
  return value === undefined || value === null ? '' : String(value)
}

const getRouteAssistantQuery = () => {
  return {
    id: getQueryValue(route.query.id),
    robotKey: getQueryValue(route.query.robot_key)
  }
}

const isSameQuery = (nextQuery) => {
  const keys = new Set([...Object.keys(route.query), ...Object.keys(nextQuery)])
  for (const key of keys) {
    if (getQueryValue(route.query[key]) !== getQueryValue(nextQuery[key])) {
      return false
    }
  }
  return true
}

const buildAssistantQuery = (assistant) => {
  return {
    ...route.query,
    id: assistant.id,
    robot_key: assistant.robot_key
  }
}

const syncAssistantQuery = async () => {
  const assistant = clawbotStore.currentAssistant
  if (!assistant?.id || isSyncingRouteQuery.value) {
    return
  }

  const nextQuery = buildAssistantQuery(assistant)
  if (isSameQuery(nextQuery)) {
    return
  }

  isSyncingRouteQuery.value = true
  try {
    await router.replace({
      path: route.path,
      query: nextQuery
    })
  } finally {
    isSyncingRouteQuery.value = false
  }
}

const initClawbotModule = async () => {
  const { id, robotKey } = getRouteAssistantQuery()
  await clawbotStore.initModule({
    forceRefresh: clawbotStore.isReady,
    targetId: id,
    targetRobotKey: robotKey
  })
  await syncAssistantQuery()
}

initClawbotModule()

watch(
  () => [route.query.id, route.query.robot_key],
  async () => {
    if (isSyncingRouteQuery.value || !clawbotStore.isReady) {
      return
    }

    const { id, robotKey } = getRouteAssistantQuery()
    if (id && clawbotStore.selectAssistantByQuery(id, robotKey)) {
      await syncAssistantQuery()
      return
    }

    await syncAssistantQuery()
  }
)

watch(
  () => [clawbotStore.currentAssistant?.id, clawbotStore.currentAssistant?.robot_key],
  async () => {
    if (!clawbotStore.isReady) {
      return
    }
    await syncAssistantQuery()
  }
)
</script>

<style lang="less" scoped>
.clawbot-layout {
  display: flex;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: #f7f8fb;

  .clawbot-layout-container {
    flex: 1;
    min-width: 0;
    height: 100%;
    overflow: auto;
  }
}
</style>
