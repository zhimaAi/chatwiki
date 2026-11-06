<template>
  <div class="customer-service-layout">
    <CustomerSidebar
      ref="sidebar"
      :current-robot="currentRobot"
      :robot-options="robotOptions"
      :robot-list-status="robotListStatus"
      @navigate="navigateTo"
      @robot-change="switchRobot"
      @retry-robots="customerServiceStore.loadRobots"
      @robot-removed="handleRobotRemoved"
      @robot-created="handleRobotCreated"
    />

    <main class="customer-main">
      <div v-if="localError || robotStatus === 'error'" class="context-message">
        <a-alert type="error" show-icon :message="t(localError || robotError || 'error_robot')" />
        <a-button @click="retryRobot">{{ t('retry') }}</a-button>
      </div>
      <div v-else-if="robotStatus === 'loading'" class="context-message context-info">
        <a-spin size="small" />
        <span>{{ t('loading_robot') }}</span>
      </div>
      <router-view v-if="!isAgentApplication" v-slot="{ Component }">
        <component
          v-if="route.meta.customerServiceMenuKey === 'session-record'"
          :is="Component"
          :robot="robot"
          :robot-status="robotStatus"
        />
        <BusinessContext v-else-if="isBusinessPage && robotStatus === 'ready' && robot" :key="robotId + ':' + robotKey" :robot="robot">
          <component :is="Component" :key="String(route.name) + ':' + robotId" />
        </BusinessContext>
        <a-empty v-else-if="isBusinessPage && robotStatus === 'empty'" :description="t('empty_robot')" />
        <component :is="Component" v-else-if="route.meta.customerServiceMenuKey === 'skill'" @switch-agent="sidebar?.openRobotSelector()" @robot-created="handleRobotCreated" />
        <component :is="Component" v-else-if="!isBusinessPage" />
      </router-view>
    </main>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { isNavigationFailure, NavigationFailureType, useRoute, useRouter } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import { useCustomerServiceStore } from '@/stores/modules/customer-service'
import CustomerSidebar from './components/customer-sidebar.vue'
import BusinessContext from './components/business-context.vue'
import {
  getCustomerServiceChatLocation,
  getCustomerServicePage,
  normalizeCustomerServiceQueryValue,
  pickCustomerServiceRobotQuery
} from './navigation'

const { t } = useI18n('views.customer-service.index')
const route = useRoute()
const router = useRouter()
const customerServiceStore = useCustomerServiceStore()
customerServiceStore.reset()
void customerServiceStore.loadRobots()

const { robots, robot, robotListStatus, robotStatus, robotError } = storeToRefs(customerServiceStore)
const robotId = computed(() => normalizeCustomerServiceQueryValue(route.query.id))
const robotKey = computed(() => normalizeCustomerServiceQueryValue(route.query.robot_key))
const localError = ref('')
const sidebar = ref(null)
const selectingCreatedRobot = ref(false)
let alive = true
let applicationVersion = 0
const isAgentApplication = computed(() => {
  const detail = robot.value
  if (detail && String(detail.id) === robotId.value && detail.robot_key === robotKey.value) {
    return Number(detail.application_type) === 2
  }
  const option = robots.value.find((item) => String(item.id) === robotId.value && item.robot_key === robotKey.value)
  return Number(option?.application_type) === 2
})
const isBusinessPage = computed(() => ['basic-settings', 'knowledge-base', 'function-center', 'analytics'].includes(route.meta.customerServiceMenuKey))

const currentRobot = computed(() => {
  if (
    !robot.value || normalizeCustomerServiceQueryValue(robot.value.id) !== robotId.value ||
    robot.value.robot_key !== robotKey.value
  ) return null
  return {
    id: robotId.value,
    name: robot.value.robot_name || t('app_name'),
    avatar: robot.value.robot_avatar_url || robot.value.robot_avatar || ''
  }
})

const robotOptions = computed(() => robots.value
  .map((item) => ({
    id: normalizeCustomerServiceQueryValue(item.id),
    name: item.robot_name || t('app_name'),
    avatar: item.robot_avatar_url || item.robot_avatar || '',
    intro: item.robot_intro || '',
    applicationType: Number(item.application_type),
    robotKey: normalizeCustomerServiceQueryValue(item.robot_key)
  }))
  .filter((item) => item.id)
)

const navigateTo = (menuKey) => {
  const page = getCustomerServicePage(menuKey)
  if (!page) return
  return router.push({ name: page.routeName, query: pickCustomerServiceRobotQuery(route.query) })
}

const switchRobot = async (selectedId) => {
  const id = normalizeCustomerServiceQueryValue(selectedId)
  const option = robots.value.find((item) => normalizeCustomerServiceQueryValue(item.id) === id)
  const key = normalizeCustomerServiceQueryValue(option?.robot_key)
  if (!id || !key) {
    localError.value = 'invalid_robot'
    return
  }

  localError.value = ''
  try {
    const failure = await router.push(getCustomerServiceChatLocation(option))
    if (failure && !isNavigationFailure(failure, NavigationFailureType.duplicated)) {
      localError.value = 'error_switch_robot'
    }
  } catch {
    localError.value = 'error_switch_robot'
  }
}

const handleRobotRemoved = async (removedId) => {
  // 先令旧列表失效，避免清空当前标识时默认选择逻辑再次选中已删除应用。
  const reload = customerServiceStore.loadRobots()
  if (normalizeCustomerServiceQueryValue(removedId) === robotId.value) {
    try {
      const failure = await router.replace({ name: 'customerServiceChat', query: {} })
      if (failure && !isNavigationFailure(failure, NavigationFailureType.duplicated)) {
        localError.value = 'error_switch_robot'
      }
    } catch {
      localError.value = 'error_switch_robot'
    }
  }
  await reload
}

const handleRobotCreated = async (created) => {
  localError.value = ''
  const sourcePath = route.fullPath
  const version = applicationVersion
  const id = normalizeCustomerServiceQueryValue(created?.id)
  selectingCreatedRobot.value = true
  try {
    // 等待新应用标识补齐期间，不能让无参入口的默认选择抢先跳转。
    const list = await customerServiceStore.loadRobots()
    if (!alive || version !== applicationVersion || route.fullPath !== sourcePath) return
    const option = list?.find((item) => String(item.id) === id) || created
    const key = normalizeCustomerServiceQueryValue(option?.robot_key)
    if (!id || !key || ![0, 2].includes(Number(option?.application_type))) {
      localError.value = list ? 'invalid_robot' : 'error_robots'
      return
    }
    if (id !== robotId.value || key !== robotKey.value) {
      const failure = await router.push(getCustomerServiceChatLocation(option))
      if (alive && route.fullPath === sourcePath && failure &&
        !isNavigationFailure(failure, NavigationFailureType.duplicated)) {
        localError.value = 'error_switch_robot'
      }
    }
  } catch {
    if (alive && version === applicationVersion && route.fullPath === sourcePath) localError.value = 'error_switch_robot'
  } finally {
    selectingCreatedRobot.value = false
  }
}

const selectDefaultRobot = () => {
  localError.value = ''
  if (selectingCreatedRobot.value || robotId.value || robotKey.value || robotListStatus.value !== 'ready') return

  // 只有完全无参才按接口顺序补首项（机器人或智能体）；首项标识坏时不跳到第二项。
  const firstRobot = robots.value[0]
  if (!firstRobot) return
  const id = normalizeCustomerServiceQueryValue(firstRobot.id)
  const key = normalizeCustomerServiceQueryValue(firstRobot.robot_key)
  if (!id || !key) {
    localError.value = 'invalid_robot'
    return
  }

  const page = getCustomerServicePage(route.meta.customerServiceMenuKey)
  if (!page) return
  if (Number(firstRobot.application_type) === 2) return router.replace(getCustomerServiceChatLocation(firstRobot))
  return router.replace({ name: page.routeName, query: { id, robot_key: key } })
}

// 客服地址也可能携带 Agent 标识；详情校验通过后才规范到新版 Agent 入口。
watch([robotStatus, robot], async ([status, info]) => {
  if (!alive || status !== 'ready' || Number(info?.application_type) !== 2 ||
    String(info.id) !== robotId.value || info.robot_key !== robotKey.value) return
  const sourcePath = route.fullPath
  try {
    const failure = await router.replace(getCustomerServiceChatLocation(info))
    if (alive && route.fullPath === sourcePath && failure &&
      !isNavigationFailure(failure, NavigationFailureType.duplicated)) localError.value = 'error_switch_robot'
  } catch {
    if (alive && route.fullPath === sourcePath) localError.value = 'error_switch_robot'
  }
})

const retryRobot = () => {
  localError.value = ''
  if (!robotId.value && !robotKey.value) return customerServiceStore.loadRobots()
  return customerServiceStore.loadRobot({ id: robotId.value, robotKey: robotKey.value })
}

// 分别比较两个归一化值，同一机器人仅换菜单时不重新请求详情。
watch(
  [() => robotId.value, () => robotKey.value],
  ([id, key]) => {
    applicationVersion++
    localError.value = ''
    void customerServiceStore.loadRobot({ id, robotKey: key })
  },
  { immediate: true }
)
watch(
  [() => robotListStatus.value, () => robotId.value, () => robotKey.value, () => route.name],
  selectDefaultRobot,
  { immediate: true }
)
onBeforeUnmount(() => {
  alive = false
  customerServiceStore.reset()
})
</script>

<style lang="less" scoped>
.customer-service-layout {
  display: flex;
  width: 100%;
  height: 100vh;
  overflow: hidden;
  color: #262626;
  background: #EEEEED;
}
.customer-main {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
  overflow: hidden;
}
.context-message {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  gap: 12px;
  padding: 12px 20px;
  background: #fff;

  .ant-alert {
    flex: 1;
    min-width: 0;
  }
}
.context-info {
  color: #8c8c8c;
}
</style>
