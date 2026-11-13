<template>
  <a-result v-if="error" status="error" :title="t('error_robot')">
    <template #extra><a-button @click="initialize">{{ t('retry') }}</a-button></template>
  </a-result>
  <slot v-else-if="ready" />
</template>

<script setup>
import { computed, onBeforeUnmount, provide, ref, toRaw, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useCustomerServiceStore } from '@/stores/modules/customer-service'
import { useRobotStore } from '@/stores/modules/robot'
import { getRobotInfo, saveRobot } from '@/api/robot'
import { useI18n } from '@/hooks/web/useI18n'
import { cloneConfig, createSaveQueue, serializeRobotConfig } from '../config-session'

const emit = defineEmits(['refresh'])
const props = defineProps({ robot: { type: Object, required: true } })
const { t } = useI18n('views.customer-service.index')
const customer = useCustomerServiceStore()
const robotStore = useRobotStore()
const ready = ref(false)
const error = ref(false)
const pending = ref(0)
const savedRevision = ref(0)
const queue = createSaveQueue()
const id = String(props.robot.id)
const robotKey = props.robot.robot_key
const applicationType = Number(props.robot.application_type)
let alive = true
let submitted = 0
let syncing = false
const isCurrent = () => alive && String(customer.robot?.id) === id && customer.robot?.robot_key === robotKey &&
  Number(customer.robot?.application_type) === applicationType

function apply(data, updateForm = true) {
  if (!isCurrent()) return
  syncing = true
  try {
    if (updateForm) robotStore.setRobotInfo(cloneConfig(data))
    customer.syncRobot(data)
    emit('refresh', data)
  } finally { syncing = false }
}
function initialize() {
  try {
    robotStore.setRobotInfo(cloneConfig(props.robot))
    error.value = false
    ready.value = true
  } catch { error.value = true }
}
async function refreshRobot() {
  await queue.flush()
  const response = await getRobotInfo({ id })
  if (String(response?.data?.id) !== id || response.data.robot_key !== robotKey ||
    Number(response.data.application_type) !== applicationType) throw new Error(t('invalid_robot'))
  apply(response.data, pending.value === 0)
  return response
}
function updateRobotInfo(patch) {
  if (!isCurrent()) return Promise.resolve(false)
  // 同一配置上下文保留原应用类型，不能把智能体按聊天机器人保存。
  Object.assign(robotStore.robotInfo, cloneConfig(patch), { id, robot_key: robotKey, application_type: applicationType })
  const snapshot = serializeRobotConfig(toRaw(robotStore.robotInfo))
  return saveAction(() => saveRobot(snapshot, applicationType))
}
function saveAction(task) {
  if (!isCurrent()) return Promise.resolve(false)
  const revision = ++submitted
  pending.value++
  return queue.enqueue(async () => {
    const response = await task()
    if (response?.res !== undefined && response.res !== 0) throw new Error(response.msg || t('save_failed'))
    const refreshed = await getRobotInfo({ id })
    if (String(refreshed?.data?.id) !== id || refreshed.data.robot_key !== robotKey ||
      Number(refreshed.data.application_type) !== applicationType) throw new Error(t('invalid_robot'))
    if (isCurrent()) {
      apply(refreshed.data, revision === submitted)
      if (revision === submitted) savedRevision.value++
    }
    return true
  }).catch(err => {
    if (isCurrent()) message.error(err.message || t('save_failed'))
    return false
  }).finally(() => { pending.value-- })
}

provide('customerRobotContext', {
  robotInfo: robotStore.robotInfo,
  saving: computed(() => pending.value > 0),
  savedRevision,
  updateRobotInfo,
  saveAction,
  refreshRobot,
  waitForSave: () => queue.flush()
})
watch(() => props.robot, data => {
  if (!syncing && !pending.value && ready.value && isCurrent()) robotStore.setRobotInfo(cloneConfig(data))
})
initialize()
onBeforeUnmount(() => { alive = false })
</script>
