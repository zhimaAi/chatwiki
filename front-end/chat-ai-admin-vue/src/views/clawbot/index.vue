<template>
  <CustomerLayout v-if="isCustomerLayout" />
  <div class="clawbot-layout" v-else-if="clawbotStore.isReady && !initializingStandalone">
    <LeftSidebar />
    <div class="clawbot-layout-container">
      <router-view v-if="clawbotStore.currentAssistant"></router-view>
    </div>
  </div>
  <ClawbotLoading v-else />
</template>

<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import LeftSidebar from './components/left-sidebar.vue'
import ClawbotLoading from './components/clawbot-loading.vue'
import CustomerLayout from './components/customer-layout.vue'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { getRobotPermission } from '@/utils/permission'
import { isCustomerServiceClawbotRoute } from '@/utils/clawbot-layout'

const route = useRoute()
const router = useRouter()
const clawbotStore = useClawbotStore()
const isSyncingRouteQuery = ref(false)
// 切回独立布局时 store 可能仍是新版应用的数据，初始化结束前暂停渲染和旧入口的 URL 同步。
const initializingStandalone = ref(false)
const isCustomerLayout = computed(() => isCustomerServiceClawbotRoute(route))
let alive = true
let initializationVersion = 0
const isStandaloneRoute = () => alive && !isCustomerLayout.value &&
  (route.path === '/clawbot' || route.path.startsWith('/clawbot/'))

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
  if (!isStandaloneRoute()) return
  const assistant = clawbotStore.currentAssistant
  if (!assistant?.id || isSyncingRouteQuery.value) {
    return
  }

  const nextQuery = buildAssistantQuery(assistant)
  const rights = getRobotPermission(assistant.id)
  const targetPath = rights == 0 ? '/no-permission'
    : rights == 1 && !['/clawbot/chat', '/clawbot/assistant', '/clawbot/stats'].includes(route.path) ? '/clawbot/chat' : route.path
  if (targetPath === route.path && isSameQuery(nextQuery)) {
    return
  }

  isSyncingRouteQuery.value = true
  try {
    await router.replace({
      path: targetPath,
      query: nextQuery
    })
  } finally {
    isSyncingRouteQuery.value = false
  }
}

const initClawbotModule = async () => {
  const version = ++initializationVersion
  initializingStandalone.value = isStandaloneRoute()
  if (!isStandaloneRoute()) return
  const { id, robotKey } = getRouteAssistantQuery()
  try {
    await clawbotStore.initModule({
      forceRefresh: clawbotStore.isReady,
      targetId: id,
      targetRobotKey: robotKey,
      isCurrent: () => isStandaloneRoute() && version === initializationVersion
    })
    if (!isStandaloneRoute() || version !== initializationVersion) return
    await syncAssistantQuery()
  } finally {
    if (version === initializationVersion) initializingStandalone.value = false
  }
}

watch(isCustomerLayout, initClawbotModule, { immediate: true })

watch(
  () => [route.query.id, route.query.robot_key],
  async () => {
    if (!isStandaloneRoute() || initializingStandalone.value || isSyncingRouteQuery.value || !clawbotStore.isReady) {
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
  () => [route.path, clawbotStore.isReady, clawbotStore.currentAssistant?.id, clawbotStore.currentAssistant?.robot_key, getRobotPermission(clawbotStore.currentAssistant?.id), isSyncingRouteQuery.value, initializingStandalone.value],
  async () => {
    if (!isStandaloneRoute() || initializingStandalone.value || !clawbotStore.isReady) {
      return
    }
    if (!clawbotStore.currentAssistant?.id) {
      await router.replace('/robot/list')
      return
    }
    await syncAssistantQuery()
  },
  { immediate: true }
)

onBeforeUnmount(() => {
  alive = false
  initializationVersion++
})
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
