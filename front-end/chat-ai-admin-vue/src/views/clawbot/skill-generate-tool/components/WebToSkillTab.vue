<template>
  <div class="generate-content">
    <div class="info-box">
      {{ t('web_description') }}
    </div>

    <a-button type="primary" class="add-btn" @click="createModalVisible = true">
      <template #icon>
        <PlusOutlined />
      </template>
      {{ t('btn_add') }}
    </a-button>
    
    <a-table
      class="skill-table"
      :columns="columns"
      :data-source="taskList"
      :loading="taskLoading"
      row-key="id"
      :pagination="taskPagination"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'skill_name'">
          <span class="skill-name">{{ record.skill_name }}</span>
        </template>
        <template v-else-if="column.dataIndex === 'web_count'">
          {{ Array.isArray(record.urls) ? record.urls.length : 0 }}
        </template>
        <template v-else-if="column.dataIndex === 'create_time'">
          {{ formatTime(record.create_time) }}
        </template>
        <template v-else-if="column.dataIndex === 'status'">
          <span class="status-tag" :class="getStatusConfig(record.status).className">
            <svg
              v-if="Number(record.status) === 1"
              class="status-tag-marquee"
              aria-hidden="true"
            >
              <rect pathLength="100" />
            </svg>
            <CheckCircleFilled v-if="Number(record.status) === 2" />
            <ExclamationCircleFilled v-else-if="Number(record.status) === 3" />
            <svg-icon
              v-else-if="Number(record.status) === 1"
              class="status-tag-icon"
              name="clock-filled"
              size="16"
            />
            <LoadingOutlined v-else-if="Number(record.status) === 4" />
            <ClockCircleFilled v-else />
            {{ getStatusConfig(record.status).text }}
          </span>
        </template>
        <template v-else-if="column.dataIndex === 'action'">
          <div class="action-list">
            <a
              v-for="action in getActions(record)"
              :key="action.key"
              href="javascript:void(0);"
              @click="handleActionClick(action.key, record)"
            >
              {{ action.label }}
            </a>
          </div>
        </template>
      </template>
    </a-table>

    <CreateWebToSkillModal v-model:visible="createModalVisible" @confirm="handleCreateConfirm" />
    <WebToSkillLogModal v-model:visible="logModalVisible" :task-id="logTaskId" />
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, reactive, ref, watch } from 'vue'
import {
  CheckCircleFilled,
  ClockCircleFilled,
  ExclamationCircleFilled,
  LoadingOutlined,
  PlusOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useUserStore } from '@/stores/modules/user'
import {
  getWebToSkillTaskList,
  installWebToSkill,
  regenerateWebToSkillTask,
  stopWebToSkillTask
} from '@/api/clawbot'
import CreateWebToSkillModal from './CreateWebToSkillModal.vue'
import WebToSkillLogModal from './WebToSkillLogModal.vue'

const props = defineProps({
  active: {
    type: Boolean,
    default: false
  }
})

const userStore = useUserStore()
const { t } = useI18n('views.clawbot.skill-generate-tool.index')

const columns = computed(() => [
  { title: t('column_skill_name'), dataIndex: 'skill_name', key: 'skill_name', width: 263 },
  { title: t('column_create_time'), dataIndex: 'create_time', key: 'create_time', width: 211 },
  { title: t('column_web_count'), dataIndex: 'web_count', key: 'web_count', width: 115 },
  { title: t('column_status'), dataIndex: 'status', key: 'status', width: 135 },
  { title: t('column_action'), dataIndex: 'action', key: 'action', width: 104 }
])

const taskList = ref([])
const taskLoading = ref(false)
const createModalVisible = ref(false)
const logModalVisible = ref(false)
const logTaskId = ref('')
let pollingTimer = null

const pager = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const taskPagination = computed(() => ({
  current: pager.page,
  pageSize: pager.pageSize,
  total: pager.total,
  showSizeChanger: true,
  pageSizeOptions: ['10', '20', '50', '100']
}))

const getStatusConfig = (status) => {
  const statusMap = {
    0: { text: t('status_pending'), className: 'pending' },
    1: { text: t('status_running'), className: 'running' },
    2: { text: t('status_success'), className: 'success' },
    3: { text: t('status_failed'), className: 'failed' },
    4: { text: t('status_stopping'), className: 'stopping' },
    5: { text: t('status_stopped'), className: 'stopped' }
  }
  return statusMap[Number(status)] || statusMap[0]
}

const formatTime = (time) => {
  const timestamp = Number(time || 0)
  if (!timestamp) {
    return '—'
  }
  const date = new Date(timestamp * 1000)
  const pad = (value) => String(value).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}

const isRequestSuccess = (res) => res && (res.res === 0 || res.code === 0)

const getActions = (record) => {
  const status = Number(record.status)
  if (status === 2) {
    return [
      { key: 'download', label: t('action_download_skill') },
      { key: 'install', label: t('action_install_skill') }
    ]
  }
  if (status === 3) {
    return [
      { key: 'retry', label: t('action_regenerate') },
      { key: 'log', label: t('action_view_log') }
    ]
  }
  if (status === 0 || status === 1) {
    return [{ key: 'stop', label: t('action_stop') }]
  }
  if (status === 5) {
    return [{ key: 'retry', label: t('action_regenerate') }]
  }
  return []
}

const updateTaskListSilently = (nextList) => {
  const currentTaskMap = new Map(taskList.value.map((item) => [item.id, item]))
  taskList.value = nextList.map((nextTask) => {
    const currentTask = currentTaskMap.get(nextTask.id)
    if (currentTask && JSON.stringify(currentTask) === JSON.stringify(nextTask)) {
      return currentTask
    }
    return nextTask
  })
}

const loadTaskList = async ({ silent = false } = {}) => {
  if (!silent) {
    taskLoading.value = true
  }
  try {
    const res = await getWebToSkillTaskList({
      page: pager.page,
      size: pager.pageSize
    })
    const data = res.data || {}
    const nextTaskList = data.list || []
    if (silent) {
      updateTaskListSilently(nextTaskList)
    } else {
      taskList.value = nextTaskList
    }
    pager.total = Number(data.total || 0)
    pager.page = Number(data.page || pager.page)
    pager.pageSize = Number(data.size || pager.pageSize)
    updatePolling()
  } catch (error) {
    console.error('获取Web转Skill任务列表失败', error)
  } finally {
    if (!silent) {
      taskLoading.value = false
    }
  }
}

const stopPolling = () => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
  }
}

const updatePolling = () => {
  const hasActiveTask = taskList.value.some((item) => [0, 1, 4].includes(Number(item.status)))
  if (!props.active || !hasActiveTask) {
    stopPolling()
    return
  }
  if (!pollingTimer) {
    pollingTimer = setInterval(() => loadTaskList({ silent: true }), 5000)
  }
}

const handleTableChange = (pagination) => {
  pager.page = pagination.current
  pager.pageSize = pagination.pageSize
  loadTaskList()
}

const handleCreateConfirm = () => {
  pager.page = 1
  loadTaskList()
}

const handleActionClick = (action, record) => {
  if (action === 'stop') {
    handleStopTask(record)
  } else if (action === 'install') {
    handleInstallSkill(record)
  } else if (action === 'download') {
    handleDownloadSkill(record)
  } else if (action === 'log') {
    logTaskId.value = record.id
    logModalVisible.value = true
  } else if (action === 'retry') {
    handleRegenerateTask(record)
  }
}

const handleStopTask = async (record) => {
  if (![0, 1, 4].includes(Number(record.status))) {
    message.warning(t('msg_task_cannot_stop'))
    return
  }
  try {
    const res = await stopWebToSkillTask({ id: record.id })
    if (isRequestSuccess(res)) {
      message.success(t('msg_stop_success'))
      loadTaskList()
    } else {
      message.error(res?.msg || t('msg_stop_failed'))
    }
  } catch (error) {
    console.error('停止Web转Skill任务失败', error)
  }
}

const handleInstallSkill = async (record) => {
  if (Number(record.status) !== 2) {
    message.warning(t('msg_task_cannot_install'))
    return
  }
  try {
    const res = await installWebToSkill({ id: record.id })
    if (isRequestSuccess(res)) {
      message.success(t('msg_install_success'))
    } else {
      message.error(res?.msg || t('msg_install_failed'))
    }
  } catch (error) {
    console.error('安装Web转Skill失败', error)
  }
}

const handleDownloadSkill = (record) => {
  if (Number(record.status) !== 2) {
    message.warning(t('msg_task_cannot_download'))
    return
  }
  const token = userStore.getToken
  const tokenQuery = token ? `&token=${encodeURIComponent(token)}` : ''
  window.open(`/manage/downloadWebToSkillFile?id=${record.id}${tokenQuery}`, '_blank')
}

const handleRegenerateTask = async (record) => {
  try {
    const res = await regenerateWebToSkillTask({ id: record.id })
    if (isRequestSuccess(res)) {
      message.success(t('msg_regenerate_success'))
      pager.page = 1
      loadTaskList()
    } else {
      message.error(res?.msg || t('msg_regenerate_failed'))
    }
  } catch (error) {
    console.error('重新生成Web转Skill任务失败', error)
  }
}

watch(
  () => props.active,
  (active) => {
    if (active) {
      loadTaskList()
    } else {
      stopPolling()
    }
  },
  { immediate: true }
)

onBeforeUnmount(stopPolling)
</script>

<style scoped lang="less">
.info-box {
  white-space: pre-line;
}

.status-tag.running {
  position: relative;
}

.status-tag-icon {
  display: inline-flex;
  flex: none;
  width: 16px;
  height: 16px;
}

.status-tag-marquee {
  position: absolute;
  inset: -2px;
  width: calc(100% + 4px);
  height: calc(100% + 4px);
  overflow: visible;
  pointer-events: none;

  rect {
    x: 1px;
    y: 1px;
    width: calc(100% - 2px);
    height: calc(100% - 2px);
    rx: 7px;
    fill: none;
    stroke: #2475fc;
    stroke-width: 2px;
    stroke-linecap: round;
    stroke-dasharray: 58 42;
    animation: status-tag-marquee 1.6s linear infinite;
  }
}

@keyframes status-tag-marquee {
  to {
    stroke-dashoffset: -100;
  }
}

@media (prefers-reduced-motion: reduce) {
  .status-tag-marquee rect {
    animation: none;
  }
}
</style>
