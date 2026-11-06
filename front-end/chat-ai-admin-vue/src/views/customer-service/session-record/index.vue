<template>
  <div class="session-record-page">
    <header class="page-header">
      <h1>{{ t('session_record') }}</h1>
    </header>

    <div class="session-toolbar">
      <a-select
        v-model:value="filters.appType"
        class="filter-control filter-channel"
        :options="channelOptions"
        :aria-label="t('all_channels')"
        @change="reload"
      />
      <a-range-picker
        v-model:value="dateRange"
        class="filter-control filter-date"
        :placeholder="[t('start_date'), t('end_date')]"
        @change="reload"
      />
      <a-input
        v-model:value="filters.name"
        class="filter-control filter-search"
        allow-clear
        :placeholder="t('search_user')"
        @pressEnter="reload"
        @change="onSearchChange"
      >
        <template #suffix><SearchOutlined /></template>
      </a-input>
      <a-tooltip :title="t('handoff_unavailable')">
        <a-select
          class="filter-control filter-handoff"
          :placeholder="t('handoff_filter')"
          :options="[]"
          disabled
          :aria-label="t('handoff_filter')"
        />
      </a-tooltip>
      <a-tooltip :title="t('handoff_unavailable')">
        <a-button class="customer-system" type="primary" disabled>
          <CommentOutlined />{{ t('customer_system') }}
        </a-button>
      </a-tooltip>
    </div>

    <div class="session-panels">
      <section class="session-list-panel">
        <div class="session-list-head">
          <h2>{{ t('session_list') }}</h2>
          <span class="session-count">{{ t('sessions_count', { count: records.length }) }}</span>
        </div>
        <div class="session-list-body" @scroll="onListScroll">
          <div v-if="loading && !records.length" class="list-state">
            <a-spin />
          </div>
          <a-empty
            v-else-if="!records.length"
            class="list-state"
            :image="Empty.PRESENTED_IMAGE_SIMPLE"
            :description="t('no_sessions')"
          />
          <button
            v-for="item in records"
            :key="item.session_id"
            type="button"
            class="session-item"
            :class="{ active: activeRecord?.session_id === item.session_id }"
            @click="selectRecord(item)"
          >
            <a-avatar :src="item.avatar" :size="38">{{ avatarText(item) }}</a-avatar>
            <div class="session-item-body">
              <div class="row1">
                <span class="name">{{ item.name || t('unknown_user') }}</span>
                <span class="time">{{ formatListTime(item.last_chat_time) }}</span>
              </div>
              <div class="preview">{{ item.last_chat_message || '' }}</div>
              <div class="meta">
                <span class="chip"><span class="dot"></span>{{ channelName(item) }}</span>
              </div>
            </div>
          </button>
          <div v-if="loading && records.length" class="list-loading"><a-spin size="small" /></div>
        </div>
      </section>

      <section class="message-panel">
        <div v-if="activeRecord" class="message-heading">
          <a-avatar :src="activeRecord.avatar" :size="34">{{ avatarText(activeRecord) }}</a-avatar>
          <div class="heading-text">
            <div class="heading-title">{{ activeRecord.name || t('unknown_user') }}</div>
            <div class="heading-sub">
              <span class="status-dot"></span>
              <span>{{ channelName(activeRecord) }} - {{ robot?.robot_name || '' }}</span>
            </div>
          </div>
        </div>
        <div class="message-body">
          <DesignMessageList
            v-if="activeRecord"
            ref="messageListRef"
            :messages="messageList"
            :robotInfo="robot || {}"
            :isEmpty="!messageList.length"
            @openPromptLog="handleOpenPromptLog"
            @scrollStart="onMessageScrollStart"
          />
          <div v-else class="message-empty">
            <a-empty
              :image="Empty.PRESENTED_IMAGE_SIMPLE"
              :description="t('select_session')"
            />
          </div>
        </div>
      </section>
    </div>

    <PromptLogAlert ref="promptLogAlertRef" />
  </div>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import dayjs from 'dayjs'
import { Empty } from 'ant-design-vue'
import { CommentOutlined, SearchOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useChatStore } from '@/stores/modules/chat'
import DesignMessageList from './components/message-list.vue'
import PromptLogAlert from '@/views/robot/robot-test/components/prompt-log-alert.vue'

const props = defineProps({
  robot: { type: Object, default: null },
  robotStatus: { type: String, default: 'loading' }
})
const { t } = useI18n('views.customer-service.index')
const chatStore = useChatStore()
const { messageList } = storeToRefs(chatStore)
const { createMsg, onGetChatMessage, getRecordList, getChannelList } = chatStore

const records = ref([])
const channelList = ref([])
const activeRecord = ref(null)
const loading = ref(false)
const hasMore = ref(false)
const dateRange = ref(null)
const messageListRef = ref(null)
const promptLogAlertRef = ref(null)
const filters = reactive({ appType: 'all', name: '' })
let listRequestSeq = 0
let recordRequestSeq = 0

const channelOptions = computed(() => [
  { value: 'all', label: t('all_channels') },
  ...channelList.value.map((item) => ({
    value: item.app_id ? `${item.app_type}::${item.app_id}` : item.app_type,
    label: item.app_name
  }))
])

const findChannel = (item) => {
  if (!item) return null
  return channelList.value.find((channel) => (
    channel.app_type === item.app_type && String(channel.app_id || '') === String(item.app_id || '')
  )) || channelList.value.find((channel) => channel.app_type === item.app_type) || null
}

const channelName = (item) => findChannel(item)?.app_name || item?.app_type || t('all_channels')
const avatarText = (item) => (item?.name || '').slice(0, 1) || '客'
const formatListTime = (time) => (time ? dayjs(Number(time) * 1000).format('MM-DD HH:mm') : '')

const buildParams = (page) => {
  const params = {
    robot_id: props.robot?.id,
    page,
    size: 20,
    name: filters.name.trim()
  }
  if (filters.appType !== 'all') {
    const [appType, appId] = String(filters.appType).split('::')
    params.app_type = appType
    if (appId) params.app_id = appId
  }
  if (dateRange.value?.[0] && dateRange.value?.[1]) {
    params.start_time = dateRange.value[0].startOf('day').unix()
    params.end_time = dateRange.value[1].endOf('day').unix()
  }
  return params
}

const loadChannels = async () => {
  if (!props.robot?.id) return
  try {
    const res = await getChannelList({ robot_id: props.robot.id })
    channelList.value = res?.data || []
  } catch {
    channelList.value = []
  }
}

const resetConversation = () => {
  activeRecord.value = null
  messageList.value = []
}

const loadRecords = async ({ reset = false } = {}) => {
  if (!props.robot?.id) return
  // 筛选变更必须能抢占在途请求；仅阻止重复的翻页加载。
  if (loading.value && !reset) return
  const requestSeq = ++listRequestSeq
  const page = reset ? 1 : Math.floor(records.value.length / 20) + 1
  loading.value = true
  try {
    const res = await getRecordList(buildParams(page))
    if (requestSeq !== listRequestSeq) return
    const list = res?.data?.list || []
    hasMore.value = Boolean(res?.data?.has_more)
    if (reset) {
      records.value = []
      resetConversation()
    }
    records.value = [...records.value, ...list]
    if (reset && records.value.length) {
      await selectRecord(records.value[0])
    }
  } catch {
    if (requestSeq === listRequestSeq) hasMore.value = false
  } finally {
    if (requestSeq === listRequestSeq) loading.value = false
  }
}

const selectRecord = async (item) => {
  if (!item || !props.robot?.robot_key) return
  const requestSeq = ++recordRequestSeq
  activeRecord.value = item
  await createMsg({
    robot_key: props.robot.robot_key,
    avatar: item.avatar,
    name: item.name,
    nickname: item.name,
    is_background: 1,
    openid: item.openid,
    rel_user_id: item.rel_user_id
  })
  const res = await onGetChatMessage()
  if (requestSeq === recordRequestSeq && res) {
    messageListRef.value?.scrollToBottom()
  }
}

const reload = () => {
  loadRecords({ reset: true })
}

const onSearchChange = () => {
  if (!filters.name) reload()
}

const onListScroll = (event) => {
  const el = event.target
  if (!hasMore.value || loading.value) return
  if (el.scrollHeight - el.scrollTop - el.clientHeight <= 60) {
    loadRecords()
  }
}

const onMessageScrollStart = async () => {
  if (!messageList.value.length) return
  const msgId = messageList.value[0].uid
  const res = await onGetChatMessage()
  if (res) {
    messageListRef.value?.scrollToMessage(msgId)
  }
}

const handleOpenPromptLog = (item) => {
  promptLogAlertRef.value?.open(item)
}

watch(
  () => props.robot?.id,
  async (id) => {
    records.value = []
    channelList.value = []
    hasMore.value = false
    loading.value = false
    resetConversation()
    listRequestSeq += 1
    recordRequestSeq += 1
    if (!id) return
    await loadChannels()
    await loadRecords({ reset: true })
  },
  { immediate: true }
)
</script>

<style lang="less" scoped>
.session-record-page {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
  background: #f5f6f8;
  font-family: 'PingFang SC', 'Noto Sans SC', sans-serif;
}
.page-header {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  height: 61px;
  padding: 0 24px;
  border-bottom: 1px solid #eee;
  background: #fff;

  h1 {
    margin: 0;
    color: #1d1d1f;
    font-size: 18px;
    font-weight: 500;
  }
}
.session-toolbar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  flex-shrink: 0;
  gap: 12px;
  padding: 16px 24px;
  border-bottom: 1px solid var(--app-border);
  background: #fff;
}
.filter-control {
  width: 180px;
}
.filter-channel {
  width: 170px;
}
.filter-date {
  width: 260px;
}
.filter-handoff {
  width: 170px;
}
.filter-search {
  width: 200px;
}
.customer-system {
  margin-left: auto;
}
.session-toolbar {
  :deep(.ant-select),
  :deep(.ant-picker),
  :deep(.ant-input-affix-wrapper),
  :deep(.ant-input),
  :deep(.ant-btn),
  :deep(.ant-empty) {
    font-family: inherit;
  }

  :deep(.ant-select) {
    height: 38px;
  }

  :deep(.ant-select-selector) {
    display: flex;
    align-items: center;
    height: 38px !important;
    padding: 0 12px !important;
    border-color: var(--app-border) !important;
    border-radius: 10px !important;
    background: #fff !important;
    box-shadow: none !important;
  }

  :deep(.ant-select-selection-item),
  :deep(.ant-select-selection-placeholder) {
    display: flex;
    align-items: center;
    margin: 0;
    padding: 0;
    color: #1d1d1f;
    font-size: 14px;
    line-height: 36px;
  }

  :deep(.ant-select-selection-search-input) {
    height: 36px;
  }

  :deep(.ant-select-arrow) {
    top: 50%;
    margin-top: 0;
    transform: translateY(-50%);
  }

  :deep(.ant-picker) {
    display: flex;
    align-items: center;
    height: 38px;
    padding: 0 12px;
    border-color: var(--app-border) !important;
    border-radius: 10px !important;
    background: #fff !important;
    box-shadow: none !important;
  }

  :deep(.ant-picker-input) {
    display: flex;
    align-items: center;
    height: 36px;
  }

  :deep(.ant-picker-input > input) {
    height: 36px;
    padding: 0;
    color: #1d1d1f;
    font-size: 14px;
    line-height: 36px;
  }

  :deep(.ant-picker-range-separator) {
    display: flex;
    align-items: center;
    height: 36px;
    padding: 0 4px;
  }

  :deep(.ant-picker-suffix) {
    display: flex;
    align-items: center;
  }

  :deep(.ant-input-affix-wrapper) {
    display: flex;
    align-items: center;
    height: 38px;
    padding: 0 12px;
    border-color: var(--app-border) !important;
    border-radius: 10px !important;
    background: #fff !important;
    box-shadow: none !important;
  }

  :deep(.ant-input-affix-wrapper > input.ant-input) {
    height: 36px;
    padding: 0;
    color: #1d1d1f;
    font-size: 14px;
    line-height: 36px;
  }

  :deep(.ant-input-suffix) {
    display: flex;
    align-items: center;
    margin-inline-start: 8px;
    color: #8c8c8c;
  }

  :deep(.ant-picker-input > input::placeholder),
  :deep(.ant-input::placeholder),
  :deep(.ant-select-selection-placeholder) {
    color: #8c8c8c;
  }

  .filter-channel {
    :deep(.ant-select-selector) {
      border-color: var(--app-primary) !important;
    }

    :deep(.ant-select-selection-item) {
      color: var(--app-primary);
      font-weight: 500;
    }

    :deep(.ant-select-arrow) {
      color: var(--app-primary);
    }
  }
}
.customer-system {
  display: inline-flex;
  align-items: center;
  height: 38px;
  line-height: 38px;
  border-radius: 10px;
}
.session-panels {
  display: flex;
  flex: 1;
  min-height: 0;
  gap: 16px;
  padding: 16px 24px 24px;
}
.session-list-panel,
.message-panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
  border: 1px solid var(--app-border);
  border-radius: 16px;
  background: #fff;
}
.session-list-panel {
  width: 340px;
  flex-shrink: 0;
}
.session-list-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
  padding: 16px 16px 10px;

  h2 {
    margin: 0;
    color: #1d1d1f;
    font-size: 17px;
    font-weight: 500;
  }

  .session-count {
    color: #8c8c8c;
    font-size: 13px;
  }
}
.session-list-body {
  flex: 1;
  min-height: 0;
  padding: 0 8px 8px;
  overflow-y: auto;
}
.list-state {
  padding: 60px 16px;
}
.list-loading {
  display: flex;
  justify-content: center;
  padding: 12px 0;
}
.session-item {
  display: flex;
  gap: 10px;
  width: 100%;
  min-height: 64px;
  padding: 12px;
  border: 0;
  border-radius: 10px;
  background: transparent;
  font: inherit;
  text-align: left;
  cursor: pointer;
  transition: background 150ms;

  &:hover {
    background: #f5f6f8;
  }

  &.active {
    background: rgba(0, 82, 217, 0.08);
  }
}
.session-item-body {
  flex: 1;
  min-width: 0;
}
.row1 {
  display: flex;
  align-items: center;
  gap: 6px;
}
.name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  color: #1d1d1f;
  font-size: 14px;
  font-weight: 500;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.time {
  flex-shrink: 0;
  color: #8c8c8c;
  font-size: 12px;
}
.preview {
  margin-top: 2px;
  overflow: hidden;
  color: #8c8c8c;
  font-size: 13px;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.meta {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 2px;
}
.chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 1px 8px;
  border-radius: 999px;
  color: #8c8c8c;
  background: #f5f6f8;
  font-size: 11px;
}
.chip .dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #30a46c;
}
.message-panel {
  flex: 1;
  min-width: 0;
}
.message-heading {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
  min-height: 60px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--app-border);
}
.heading-text {
  min-width: 0;
}
.heading-title {
  overflow: hidden;
  color: #1d1d1f;
  font-size: 15px;
  font-weight: 500;
  line-height: 1.3;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.heading-sub {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #8c8c8c;
  font-size: 12px;
}
.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #30a46c;
}
.message-body {
  display: flex;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}
.message-empty {
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}
.session-record-page :deep(.ant-avatar) {
  font-family: inherit;
}
@media (max-width: 1280px) {
  .session-list-panel {
    width: 280px;
  }
  .filter-control {
    width: 160px;
  }
  .filter-date {
    width: 230px;
  }
  .filter-search {
    width: 180px;
  }
}
</style>
