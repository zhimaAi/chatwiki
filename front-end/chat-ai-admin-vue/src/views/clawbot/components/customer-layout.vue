<template>
  <div class="customer-agent-layout">
    <CustomerSidebar
      is-agent
      :current-robot="currentRobot"
      :robot-options="robotOptions"
      :robot-list-status="robotListStatus"
      :selected-key="selectedMenuKey"
      @navigate="navigateTo"
      @robot-change="switchRobot"
      @retry-robots="customerStore.loadRobots"
      @robot-created="handleRobotCreated"
      @robot-removed="handleRobotRemoved"
    />
    <main class="customer-agent-main">
      <a-result v-if="error" status="error" :title="t(error)">
        <template #extra><a-button @click="initialize">{{ t('retry') }}</a-button></template>
      </a-result>
      <a-empty v-else-if="!robotId && !robotKey" :description="t('select_app')" />
      <div v-else-if="!contextReady" class="agent-loading"><a-spin /></div>
      <router-view v-else v-slot="{ Component }">
        <component :is="Component" :key="String(route.name) + ':' + robotId + ':' + robotKey" />
      </router-view>
    </main>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, provide, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { isNavigationFailure, NavigationFailureType, useRoute, useRouter } from 'vue-router'
import { getRobotInfo } from '@/api/robot'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { useCustomerServiceStore } from '@/stores/modules/customer-service'
import { useRobotStore } from '@/stores/modules/robot'
import { getRobotPermission } from '@/utils/permission'
import { CUSTOMER_SERVICE_LAYOUT, isCustomerServiceClawbotRoute } from '@/utils/clawbot-layout'
import CustomerSidebar from '@/views/customer-service/components/customer-sidebar.vue'
import {
  CUSTOMER_SERVICE_AGENT_MENU_ITEMS,
  getCustomerServiceChatLocation,
  normalizeCustomerServiceQueryValue
} from '@/views/customer-service/navigation'
import { getClawbotLocation } from '../navigation'

const { t } = useI18n('views.customer-service.index')
const route = useRoute()
const router = useRouter()
const clawbotStore = useClawbotStore()
const customerStore = useCustomerServiceStore()
const robotStore = useRobotStore()
const { robots, robot, robotListStatus } = storeToRefs(customerStore)
const robotId = computed(() => normalizeCustomerServiceQueryValue(route.query.id))
const robotKey = computed(() => normalizeCustomerServiceQueryValue(route.query.robot_key))
const readyIdentity = ref(null)
const error = ref('')
let alive = true
// 分别淘汰旧初始化、store 选中项同步和共享 robot 详情请求；同一应用重试也会推进初始化版本。
let initializationVersion = 0
let selectionVersion = 0
let robotRequestVersion = 0

const isCurrentIdentity = (identity, version = initializationVersion) => alive &&
  isCustomerServiceClawbotRoute(route) && version === initializationVersion &&
  identity.id === robotId.value && identity.robotKey === robotKey.value

const selectedMenuKey = computed(() => {
  if (route.name === 'clawbotSkillGenerateTool') return 'skills'
  return CUSTOMER_SERVICE_AGENT_MENU_ITEMS.find((item) => item.routeName === route.name)?.key || ''
})
const pageAllowed = computed(() => {
  const rights = getRobotPermission(robotId.value)
  return rights > 0 && (rights !== 1 || ['chat', 'assistant', 'stats'].includes(selectedMenuKey.value))
})
// store 会跨布局保留，isReady 可能属于上一次应用，必须核对 URL、选中项和详情后才渲染页面。
const contextReady = computed(() => readyIdentity.value && isCurrentIdentity(readyIdentity.value) &&
  pageAllowed.value && clawbotStore.currentAssistantId === robotId.value &&
  clawbotStore.currentAssistant?.robot_key === robotKey.value &&
  String(clawbotStore.robotInfo?.id || '') === robotId.value && clawbotStore.robotInfo?.robot_key === robotKey.value)

const currentRobot = computed(() => {
  if (String(robot.value?.id || '') !== robotId.value || robot.value?.robot_key !== robotKey.value) return null
  return {
    id: robotId.value,
    name: robot.value.robot_name || t('type_agent'),
    avatar: robot.value.robot_avatar_url || robot.value.robot_avatar || ''
  }
})
const robotOptions = computed(() => robots.value.map((item) => ({
  id: normalizeCustomerServiceQueryValue(item.id),
  name: item.robot_name || t('app_name'),
  avatar: item.robot_avatar_url || item.robot_avatar || '',
  intro: item.robot_intro || '',
  applicationType: Number(item.application_type),
  robotKey: normalizeCustomerServiceQueryValue(item.robot_key)
})).filter((item) => item.id))

async function navigate(location, replace = false) {
  const sourcePath = route.fullPath
  try {
    const failure = replace ? await router.replace(location) : await router.push(location)
    if (failure && !isNavigationFailure(failure, NavigationFailureType.duplicated)) {
      if (alive && route.fullPath === sourcePath) error.value = 'error_switch_robot'
    }
  } catch {
    if (alive && route.fullPath === sourcePath) error.value = 'error_switch_robot'
  }
}

function navigateTo(menuKey) {
  const page = CUSTOMER_SERVICE_AGENT_MENU_ITEMS.find((item) => item.key === menuKey)
  if (!page) return
  return navigate({ name: page.routeName, query: {
    id: robotId.value,
    robot_key: robotKey.value,
    ...page.query,
    layout: CUSTOMER_SERVICE_LAYOUT
  } })
}

async function initialize() {
  const version = ++initializationVersion
  const identity = { id: robotId.value, robotKey: robotKey.value }
  readyIdentity.value = null
  error.value = ''
  const isCurrent = () => isCurrentIdentity(identity, version)
  try {
    const info = await customerStore.loadRobot(identity)
    if (!isCurrent()) return
    if (!info) {
      if (!identity.id && !identity.robotKey) {
        const list = await customerStore.loadRobots()
        if (!isCurrent()) return
        const firstAgent = list?.find((item) => Number(item.application_type) === 2 && item.id && item.robot_key)
        if (firstAgent) await navigate(getCustomerServiceChatLocation(firstAgent), true)
        return
      }
      error.value = customerStore.robotError || 'invalid_robot'
      return
    }
    if (Number(info.application_type) !== 2) {
      await navigate(getCustomerServiceChatLocation(info), true)
      return
    }
    if (getRobotPermission(identity.id) === 0) {
      error.value = 'invalid_robot'
      return
    }
    const assistant = await clawbotStore.initModule({
      forceRefresh: true,
      targetId: identity.id,
      targetRobotKey: identity.robotKey,
      autoCreate: false,
      isCurrent
    })
    if (!isCurrent()) return
    if (!assistant || assistant.id !== identity.id || assistant.robot_key !== identity.robotKey) {
      error.value = 'invalid_robot'
      return
    }
    const details = clawbotStore.robotInfo
    if (String(details?.id || '') !== identity.id || details?.robot_key !== identity.robotKey || Number(details?.application_type) !== 2) {
      error.value = 'error_robot'
      return
    }
    readyIdentity.value = identity
  } catch {
    if (isCurrent()) error.value = 'error_robot'
  }
}

function switchRobot(selectedId) {
  const option = robots.value.find((item) => String(item.id) === String(selectedId))
  if (!option?.id || !option.robot_key) {
    error.value = 'invalid_robot'
    return
  }
  error.value = ''
  return navigate(getCustomerServiceChatLocation(option))
}

async function handleRobotCreated(created) {
  const sourcePath = route.fullPath
  const version = initializationVersion
  const list = await customerStore.loadRobots()
  if (!alive || version !== initializationVersion || sourcePath !== route.fullPath) return
  const option = list?.find((item) => String(item.id) === String(created?.id)) || created
  if (!option?.id || !option.robot_key || ![0, 2].includes(Number(option.application_type))) {
    error.value = list ? 'invalid_robot' : 'error_robots'
    return
  }
  await navigate(getCustomerServiceChatLocation(option))
}

async function handleRobotRemoved(removedId) {
  const identity = { id: robotId.value, robotKey: robotKey.value }
  const wasCurrent = String(removedId) === identity.id
  if (wasCurrent) {
    readyIdentity.value = null
    initializationVersion++
  }
  const version = initializationVersion
  const list = await customerStore.loadRobots()
  if (!isCurrentIdentity(identity, version)) return
  if (!list) {
    if (wasCurrent) error.value = 'error_robots'
    return
  }
  if (wasCurrent) {
    const location = list.length ? getCustomerServiceChatLocation(list[0])
      : { name: 'clawbotChat', query: { layout: CUSTOMER_SERVICE_LAYOUT } }
    await navigate(location, true)
  } else {
    await clawbotStore.fetchAssistants()
  }
}

// 设置和对外服务复用 robot store；身份变化后，旧请求不能回填这个共享对象。
async function getRobot(targetId = robotId.value) {
  const identity = readyIdentity.value
  if (!identity || !isCurrentIdentity(identity) || String(targetId) !== identity.id) return false
  const version = ++robotRequestVersion
  const initialization = initializationVersion
  const canApply = () => isCurrentIdentity(identity, initialization) &&
    version === robotRequestVersion && clawbotStore.currentAssistantId === identity.id
  try {
    const response = await getRobotInfo({ id: identity.id })
    if (!canApply()) return false
    const info = response?.data
    if (String(info?.id || '') !== identity.id || info?.robot_key !== identity.robotKey || Number(info?.application_type) !== 2) {
      error.value = 'invalid_robot'
      return false
    }
    if (info.library_ids) info.library_ids = info.library_ids.split(',').filter(Boolean).join(',')
    robotStore.setRobotInfo(info)
    return response
  } catch {
    if (canApply()) error.value = 'error_robot'
    return false
  }
}
provide('customerAgentContext', { getRobot })

// Agent 页面内部的创建和删除会先修改 store，初始化完成后才允许反向同步 URL。
watch(() => clawbotStore.currentAssistantId, async () => {
  const version = ++selectionVersion
  const initialization = initializationVersion
  const identity = readyIdentity.value
  if (!identity || !isCurrentIdentity(identity) || clawbotStore.isInitializing) return
  const assistant = clawbotStore.currentAssistant
  if (assistant?.id === identity.id) return
  const list = await customerStore.loadRobots()
  if (!isCurrentIdentity(identity, initialization) || version !== selectionVersion) return
  if (!list) {
    error.value = 'error_robots'
    return
  }
  const option = assistant || list[0]
  await navigate(option ? getCustomerServiceChatLocation(option)
    : { name: 'clawbotChat', query: { layout: CUSTOMER_SERVICE_LAYOUT } }, true)
})

watch(() => clawbotStore.assistantList, () => {
  if (readyIdentity.value && isCurrentIdentity(readyIdentity.value) && !clawbotStore.isInitializing &&
    clawbotStore.currentAssistantId === robotId.value) void customerStore.loadRobots()
})
watch(() => clawbotStore.robotInfo, (info) => {
  if (contextReady.value) customerStore.syncRobot(info)
})
watch([readyIdentity, pageAllowed], ([identity, allowed]) => {
  if (!identity || !isCurrentIdentity(identity) || allowed) return
  if (getRobotPermission(identity.id) === 0) error.value = 'invalid_robot'
  else void navigate(getClawbotLocation(route, 'chat'), true)
})

customerStore.reset()
void customerStore.loadRobots()
watch([robotId, robotKey], initialize, { immediate: true })
onBeforeUnmount(() => {
  alive = false
  initializationVersion++
  robotRequestVersion++
  customerStore.reset()
})
</script>

<style scoped>
.customer-agent-layout {
  display: flex;
  width: 100%;
  height: 100vh;
  overflow: hidden;
  color: #262626;
  background: #f5f6f8;
}
.customer-agent-main {
  flex: 1;
  min-width: 0;
  min-height: 0;
  height: 100%;
  overflow: auto;
}
.agent-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}
</style>
