<template>
  <div class="generate-content">
    <div class="info-box">
      {{ t('book_description') }}
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
        <template v-else-if="column.dataIndex === 'create_time'">
          {{ formatTime(record.create_time) }}
        </template>
        <template v-else-if="column.dataIndex === 'status'">
          <span class="status-tag" :class="getStatusConfig(record.status).className">
            <CheckCircleFilled v-if="Number(record.status) === 2" />
            <ExclamationCircleFilled v-else-if="Number(record.status) === 3" />
            <LoadingOutlined v-else-if="[1, 4].includes(Number(record.status))" />
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

    <CreateBookToSkillModal
      v-model:visible="createModalVisible"
      @confirm="handleCreateConfirm"
    />
    <BookToSkillLogModal v-model:visible="logModalVisible" :task-id="logTaskId" />
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
  getDocToSkillTaskList,
  installDocToSkill,
  regenerateDocToSkillTask,
  stopDocToSkillTask
} from '@/api/clawbot'
import CreateBookToSkillModal from './CreateBookToSkillModal.vue'
import BookToSkillLogModal from './BookToSkillLogModal.vue'

const props = defineProps({
  active: {
    type: Boolean,
    default: false
  }
})

const userStore = useUserStore()
const { t } = useI18n('views.clawbot.skill-generate-tool.index')

const columns = computed(() => [
  { title: t('column_skill_name'), dataIndex: 'skill_name', key: 'skill_name' },
  { title: t('column_upload_time'), dataIndex: 'create_time', key: 'create_time', width: 190 },
  { title: t('column_status'), dataIndex: 'status', key: 'status', width: 130 },
  { title: t('column_action'), dataIndex: 'action', key: 'action', width: 120 }
])

const taskList = ref([])
const taskLoading = ref(false)
const createModalVisible = ref(false)
const logModalVisible = ref(false)
const logTaskId = ref('')
let pollingTimer = null

const pager = reactive({
  page: 1,
  pageSize: 10,
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
      { key: 'install', label: t('action_install') },
      { key: 'download', label: t('action_download') }
    ]
  }
  if (status === 3) {
    return [
      { key: 'retry', label: t('action_retry') },
      { key: 'log', label: t('action_log') }
    ]
  }
  if (status === 0 || status === 1) {
    return [{ key: 'stop', label: t('action_stop') }]
  }
  if (status === 5) {
    return [{ key: 'retry', label: t('action_retry') }]
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
    const res = await getDocToSkillTaskList({
      status: -1,
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
    console.error('获取Book转Skill任务列表失败', error)
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
    handleRetryTask(record)
  }
}

const handleStopTask = async (record) => {
  const status = Number(record.status)
  if (status !== 0 && status !== 1) {
    message.warning(t('msg_task_cannot_stop'))
    return
  }
  try {
    const res = await stopDocToSkillTask({ id: record.id })
    if (isRequestSuccess(res)) {
      message.success(t('msg_stop_success'))
      loadTaskList()
    } else {
      message.error(res?.msg || t('msg_stop_failed'))
    }
  } catch (error) {
    console.error('停止Book转Skill任务失败', error)
  }
}

const handleInstallSkill = async (record) => {
  if (Number(record.status) !== 2) {
    message.warning(t('msg_task_cannot_install'))
    return
  }
  try {
    const res = await installDocToSkill({ id: record.id })
    if (isRequestSuccess(res)) {
      message.success(t('msg_install_success'))
    } else {
      message.error(res?.msg || t('msg_install_failed'))
    }
  } catch (error) {
    console.error('安装Book转Skill失败', error)
  }
}

const handleDownloadSkill = (record) => {
  if (Number(record.status) !== 2) {
    message.warning(t('msg_task_cannot_download'))
    return
  }
  const token = userStore.getToken
  const tokenQuery = token ? `&token=${encodeURIComponent(token)}` : ''
  window.open(`/manage/downloadDocToSkillFile?id=${record.id}${tokenQuery}`, '_blank')
}

const handleRetryTask = async (record) => {
  try {
    const res = await regenerateDocToSkillTask({ id: record.id })
    if (isRequestSuccess(res)) {
      message.success(t('msg_retry_success'))
      loadTaskList()
    } else {
      message.error(res?.msg || t('msg_retry_failed'))
    }
  } catch (error) {
    console.error('重试Book转Skill任务失败', error)
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
