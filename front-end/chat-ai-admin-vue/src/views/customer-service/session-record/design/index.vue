<template>
  <div class="session-record-page">
    <section class="session-list-panel">
      <header class="session-list-head">
        <a-popover
          v-model:open="panels.channel"
          trigger="click"
          placement="bottomLeft"
          overlay-class-name="session-channel-popover"
        >
          <template #content>
            <div class="channel-menu">
              <button
                v-for="option in channelOptions"
                :key="option.value"
                type="button"
                class="channel-option"
                :class="{ active: option.value === filters.appType }"
                @click="selectChannel(option)"
              >
                <span class="channel-option-icon">
                  <svg v-if="option.icon === 'all'" viewBox="0 0 16 16" fill="none" aria-hidden="true">
                    <rect x="1.7" y="1.7" width="9" height="9" rx="2.2" stroke="currentColor" stroke-width="1.3" />
                    <rect x="5.3" y="5.3" width="9" height="9" rx="2.2" stroke="currentColor" stroke-width="1.3" />
                    <circle cx="5.3" cy="5.3" r="1.35" stroke="currentColor" stroke-width="1.1" />
                    <circle cx="8.9" cy="8.9" r="1.35" stroke="currentColor" stroke-width="1.1" />
                  </svg>
                  <svg v-else-if="option.icon === 'mini'" viewBox="0 0 16 16" fill="none" aria-hidden="true">
                    <circle cx="8" cy="8" r="6.3" stroke="currentColor" stroke-width="1.3" />
                    <path
                      d="M10.2 5.7c-.5-.7-2.3-.9-3 .1-.7 1 .1 1.8 1.3 2.1 1.2.3 2 1.1 1.3 2.1-.7 1-2.5.8-3 .1"
                      stroke="currentColor"
                      stroke-width="1.3"
                      stroke-linecap="round"
                    />
                  </svg>
                  <MobileOutlined v-else-if="option.icon === 'mobile'" />
                  <ApiOutlined v-else-if="option.icon === 'api'" />
                  <GlobalOutlined v-else />
                </span>
                <span class="channel-option-name">{{ option.label }}</span>
                <span v-if="option.count !== null" class="channel-option-count" :title="t('channel_count_hint')">{{ option.count }}</span>
              </button>
            </div>
          </template>
          <button
            type="button"
            class="channel-trigger"
            :class="{ active: panels.channel || filters.appType !== 'all' }"
          >
            <svg class="channel-trigger-icon" viewBox="0 0 16 16" fill="none" aria-hidden="true">
              <rect x="1.7" y="1.7" width="9" height="9" rx="2.2" stroke="currentColor" stroke-width="1.3" />
              <rect x="5.3" y="5.3" width="9" height="9" rx="2.2" stroke="currentColor" stroke-width="1.3" />
              <circle cx="5.3" cy="5.3" r="1.35" stroke="currentColor" stroke-width="1.1" />
              <circle cx="8.9" cy="8.9" r="1.35" stroke="currentColor" stroke-width="1.1" />
            </svg>
            <span class="channel-trigger-label">{{ selectedChannelOption.label }}</span>
            <span class="channel-trigger-count">{{ selectedChannelOption.count }}</span>
          </button>
        </a-popover>

        <div class="session-head-actions">
          <a-popover
            v-model:open="panels.date"
            trigger="click"
            placement="bottomLeft"
            :align="{ offset: dateAlign }"
            :force-render="true"
            overlay-class-name="session-date-popover"
          >
            <template #content>
              <DateRangePanel
                :active="panels.date"
                :value="dateRange"
                @close="setPanel('date', false)"
                @confirm="applyDateRange"
                @clear="clearDateRange"
              />
            </template>
            <button
              ref="dateTriggerRef"
              type="button"
              class="icon-button"
              :class="{ active: panels.date || Boolean(dateRange) }"
              :aria-label="t('date_range')"
            >
              <CalendarOutlined />
            </button>
          </a-popover>

          <a-popover
            v-model:open="panels.search"
            trigger="click"
            placement="bottomLeft"
            overlay-class-name="session-search-popover"
          >
            <template #content>
              <div class="simple-panel search-panel">
                <div class="panel-head">
                  <span class="panel-title">{{ t('search_user') }}</span>
                  <button type="button" class="panel-close" @click="setPanel('search', false)">
                    <CloseOutlined />
                  </button>
                </div>
                <a-input
                  v-model:value="filters.name"
                  class="search-input"
                  allow-clear
                  :placeholder="t('search_user')"
                  @press-enter="applySearch"
                  @change="onSearchInput"
                >
                  <template #prefix><SearchOutlined /></template>
                </a-input>
              </div>
            </template>
            <button
              type="button"
              class="icon-button"
              :class="{ active: panels.search || Boolean(filters.name) }"
              :aria-label="t('search_user')"
            >
              <SearchOutlined />
            </button>
          </a-popover>

          <a-popover
            v-model:open="panels.transfer"
            trigger="click"
            placement="bottomLeft"
            overlay-class-name="session-transfer-popover"
          >
            <template #content>
              <div class="simple-panel transfer-panel">
                <div class="panel-head">
                  <span class="panel-title">{{ t('handoff_filter') }}</span>
                  <button type="button" class="panel-close" @click="setPanel('transfer', false)">
                    <CloseOutlined />
                  </button>
                </div>
                <div class="transfer-options" aria-disabled="true">
                  <div class="transfer-option active">
                    <span class="radio checked"></span>
                    <span>{{ t('all') }}</span>
                  </div>
                  <div class="transfer-option">
                    <span class="radio"></span>
                    <span>{{ t('transferred') }}</span>
                  </div>
                  <div class="transfer-option">
                    <span class="radio"></span>
                    <span>{{ t('not_transferred') }}</span>
                  </div>
                </div>
                <div class="transfer-hint">{{ t('handoff_unavailable') }}</div>
              </div>
            </template>
            <a-tooltip :title="t('handoff_unavailable')">
              <button type="button" class="icon-button" :class="{ active: panels.transfer }" :aria-label="t('handoff_filter')">
                <CustomerServiceOutlined />
              </button>
            </a-tooltip>
          </a-popover>
        </div>
      </header>

      <div class="session-list-body" :class="{ 'is-scrolling': listScrolling }" @scroll="onListScroll">
        <div v-if="loading && !records.length" class="list-state">
          <a-spin />
        </div>
        <a-empty
          v-else-if="!records.length"
          class="list-state"
          :image="Empty.PRESENTED_IMAGE_SIMPLE"
          :description="t('no_sessions')"
        />
        <div
          v-for="item in records"
          :key="item.session_id"
          role="button"
          tabindex="0"
          class="session-item"
          :class="{ active: activeRecord?.session_id === item.session_id }"
          @click="selectRecord(item)"
          @keydown.enter.prevent="selectRecord(item)"
          @keydown.space.prevent="selectRecord(item)"
        >
          <div class="avatar-wrap">
            <a-avatar :src="item.avatar" :size="48">{{ avatarText(item) }}</a-avatar>
            <span v-if="Number(item.unread_count) > 0" class="unread-badge">{{ item.unread_count }}</span>
          </div>
          <div class="session-item-body">
            <div class="row1">
              <span class="name">{{ item.name || t('unknown_user') }}</span>
              <span class="time">{{ formatListTime(item.last_chat_time) }}</span>
            </div>
            <div class="preview">{{ previewText(item.last_chat_message) }}</div>
            <div class="meta">
              <span class="chip"><span class="dot"></span><span class="chip-text">{{ previewText(channelName(item)) }}</span></span>
            </div>
          </div>
        </div>
        <div v-if="loading && records.length" class="list-loading"><a-spin size="small" /></div>
      </div>
    </section>

    <section class="message-panel">
      <div v-if="activeRecord" class="message-heading">
        <a-avatar :src="activeRecord.avatar" :size="44">{{ avatarText(activeRecord) }}</a-avatar>
        <div class="heading-text">
          <div class="heading-title">{{ activeRecord.name || t('unknown_user') }}</div>
          <div class="heading-sub">
            <span class="status-dot"></span>
            <span class="heading-sub-text">{{ channelName(activeRecord) }} · {{ robot?.robot_name || '' }}</span>
          </div>
        </div>
        <div class="heading-actions">
          <a-tooltip :title="t('handoff_unavailable')">
            <span class="customer-system" role="button" aria-disabled="true">
              <CustomerServiceOutlined />
              <span>{{ t('customer_system') }}</span>
            </span>
          </a-tooltip>
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

    <PromptLogAlert ref="promptLogAlertRef" />
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import dayjs from 'dayjs'
import { Empty } from 'ant-design-vue'
import {
  ApiOutlined,
  CalendarOutlined,
  CloseOutlined,
  CustomerServiceOutlined,
  GlobalOutlined,
  MobileOutlined,
  SearchOutlined
} from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useChatStore } from '@/stores/modules/chat'
import DesignMessageList from './components/message-list.vue'
import DateRangePanel from './components/date-range-panel.vue'
import PromptLogAlert from '@/views/robot/robot-test/components/prompt-log-alert.vue'

const props = defineProps({
  robot: { type: Object, default: null },
  robotStatus: { type: String, default: 'loading' }
})
const { t } = useI18n('views.customer-service.index')
const chatStore = useChatStore()
const { messageList } = storeToRefs(chatStore)
const { createMsg, onGetChatMessage, getRecordList, getChannelList } = chatStore

const panels = reactive({ channel: false, date: false, search: false, transfer: false })
const records = ref([])
const channelList = ref([])
const channelStats = ref({ total: null, counts: {} })
const activeRecord = ref(null)
const loading = ref(false)
const hasMore = ref(false)
const dateRange = ref(null)
const dateTriggerRef = ref(null)
const dateAlign = ref([0, 4])
const messageListRef = ref(null)
const promptLogAlertRef = ref(null)
const filters = reactive({ appType: 'all', name: '' })
let listRequestSeq = 0
let channelRequestSeq = 0
let recordRequestSeq = 0
let searchTimer = null
let alive = true

const setPanel = (name, open) => {
  panels[name] = open
}

const updateDateAlignment = () => {
  const trigger = dateTriggerRef.value?.getBoundingClientRect()
  const panel = dateTriggerRef.value?.closest('.session-list-panel')?.getBoundingClientRect()
  if (!trigger || !panel) return
  // 设计稿日期浮层左边缘与列表卡片左边缘对齐（内缩 8px）。
  dateAlign.value = [panel.left + 8 - trigger.left, 4]
}

onMounted(() => {
  updateDateAlignment()
  window.addEventListener('resize', updateDateAlignment)
})

const channelOptions = computed(() => {
  const all = {
    value: 'all',
    label: t('all_channels'),
    count: channelStats.value.total,
    icon: 'all'
  }
  const channels = channelList.value.map((channel) => {
    const value = channel.app_id ? `${channel.app_type}::${channel.app_id}` : channel.app_type
    return {
      value,
      label: channel.app_name,
      count: channelStats.value.counts[value] ?? null,
      icon: channelIcon(channel)
    }
  })
  return [all, ...channels]
})

const selectedChannelOption = computed(() => (
  channelOptions.value.find((option) => option.value === filters.appType) || channelOptions.value[0]
))

const channelIcon = (channel) => {
  const name = String(channel?.app_name || '').toLowerCase()
  if (name.includes('小程序') || name.includes('wechat') || name.includes('weixin')) return 'mini'
  if (name.includes('h5') || name.includes('手机')) return 'mobile'
  if (name.includes('开放') || name.includes('api')) return 'api'
  return 'globe'
}

const findChannel = (item) => {
  if (!item) return null
  return channelList.value.find((channel) => (
    channel.app_type === item.app_type && String(channel.app_id || '') === String(item.app_id || '')
  )) || null
}

const channelName = (item) => findChannel(item)?.app_name || item?.app_type || t('all_channels')
const avatarText = (item) => (item?.name || '').slice(0, 1) || '客'
const formatListTime = (time) => (time ? dayjs(Number(time) * 1000).format('MM-DD HH:mm') : '')
const previewText = (message) => String(message || '').replace(/[\s\u2028\u2029]+/g, ' ').trim()

const applyDateRange = (range) => {
  if (!range?.[0] || !range?.[1]) return
  dateRange.value = [range[0], range[1]]
  setPanel('date', false)
  reload()
}

const clearDateRange = () => {
  dateRange.value = null
  reload()
}

const selectChannel = (option) => {
  filters.appType = option.value
  setPanel('channel', false)
  reload()
}

const applySearch = () => {
  if (searchTimer) {
    clearTimeout(searchTimer)
    searchTimer = null
  }
  reload()
}

const onSearchInput = () => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    searchTimer = null
    if (alive) reload()
  }, 400)
}

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
  const requestSeq = ++channelRequestSeq
  try {
    const res = await getChannelList({ ...buildParams(1), with_stats: 1 })
    if (!alive || requestSeq !== channelRequestSeq) return
    const list = Array.isArray(res?.data) ? res.data : res?.data?.list || []
    channelList.value = list
    const counts = {}
    list.forEach(channel => {
      const key = channel.app_id ? `${channel.app_type}::${channel.app_id}` : channel.app_type
      if (channel.count !== undefined) counts[key] = Number(channel.count)
    })
    channelStats.value = { total: res?.data?.total ?? null, counts }
  } catch {
    if (!alive || requestSeq !== channelRequestSeq) return
    channelList.value = []
    channelStats.value = { total: null, counts: {} }
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
  loadChannels()
  loadRecords({ reset: true })
}

const listScrolling = ref(false)
let listScrollTimer
const onListScroll = (event) => {
  listScrolling.value = true
  clearTimeout(listScrollTimer)
  listScrollTimer = setTimeout(() => { listScrolling.value = false }, 800)
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
    channelStats.value = { total: null, counts: {} }
    channelRequestSeq += 1
    hasMore.value = false
    loading.value = false
    resetConversation()
    listRequestSeq += 1
    recordRequestSeq += 1
    if (!id) return
    const requestSeq = listRequestSeq
    await loadChannels()
    if (!alive || requestSeq !== listRequestSeq) return
    await loadRecords({ reset: true })
  },
  { immediate: true }
)

onBeforeUnmount(() => {
  alive = false
  clearTimeout(listScrollTimer)
  if (searchTimer) clearTimeout(searchTimer)
  window.removeEventListener('resize', updateDateAlignment)
  listRequestSeq += 1
  channelRequestSeq += 1
  recordRequestSeq += 1
})
</script>

<style lang="less" scoped>
.session-record-page {
  display: flex;
  flex: 1;
  gap: 8px;
  min-height: 0;
  padding: 0 8px 8px 0;
  overflow: hidden;
  background: #f2f3f5;
  font-family: 'PingFang SC', 'Noto Sans SC', sans-serif;
}
.session-list-panel,
.message-panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
  border-radius: 12px;
  background: #fff;
}
.session-list-panel {
  width: 320px;
  flex-shrink: 0;
}
.message-panel {
  flex: 1;
  min-width: 0;
}
.session-list-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
  height: 72px;
  padding: 0 24px;
}
.channel-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 40px;
  margin-left: -12px;
  padding: 0 12px;
  border: 0;
  border-radius: 8px;
  color: #242424;
  background: transparent;
  font: inherit;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;

  &:hover {
    background: #f5f6f8;
  }

  &.active {
    color: #1456f0;
    background: #eef3ff;
  }
}
.channel-trigger-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}
.channel-trigger-count {
  color: #888;
  font-weight: 400;
}
.channel-trigger.active .channel-trigger-count {
  color: #242424;
}
.session-head-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}
.icon-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 0;
  border-radius: 8px;
  color: #4e5969;
  background: transparent;
  font-size: 16px;
  cursor: pointer;

  &:hover {
    background: #f5f6f8;
  }

  &.active {
    color: #1456f0;
    background: #eef3ff;
  }
}
.session-list-body {
  flex: 1;
  min-height: 0;
  padding: 0 8px 8px;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: transparent transparent;

  &::-webkit-scrollbar { width: 6px; height: 6px; }
  &::-webkit-scrollbar-track { background: transparent; }
  &::-webkit-scrollbar-thumb { background: transparent; border-radius: 4px; }
  &::-webkit-scrollbar-button { display: none; width: 0; height: 0; }
  &.is-scrolling {
    scrollbar-color: #d1d5db transparent;
    &::-webkit-scrollbar-thumb { background: #d1d5db; }
  }
  @supports selector(::-webkit-scrollbar) {
    scrollbar-width: auto;
    scrollbar-color: auto;
    &.is-scrolling { scrollbar-color: auto; }
  }
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
  height: 108px;
  padding: 12px;
  border: 0;
  border-radius: 8px;
  background: transparent;
  font: inherit;
  text-align: left;
  cursor: pointer;
  overflow: hidden;
  transition: background 150ms;

  &:hover {
    background: #f7f8fa;
  }

  &.active {
    background: #e8eefd;
  }
}
.avatar-wrap {
  position: relative;
  flex-shrink: 0;
  margin-top: -1px;
}
.unread-badge {
  position: absolute;
  top: -2px;
  right: -2px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border: 2px solid #fff;
  border-radius: 999px;
  color: #fff;
  background: #f53f3f;
  font-size: 11px;
  line-height: 14px;
  text-align: center;
}
.session-item-body {
  display: flex;
  flex: 1;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0;
  height: 84px;
  overflow: hidden;
}
.row1 {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  gap: 8px;
}
.name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  color: #242424;
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.time {
  flex-shrink: 0;
  color: #888;
  font-size: 12px;
}
.preview {
  display: block;
  flex-shrink: 0;
  max-width: 100%;
  height: 22px;
  min-width: 0;
  overflow: hidden;
  color: #666;
  font-size: 14px;
  line-height: 22px;
  text-overflow: ellipsis;
  white-space: nowrap;
  word-break: normal;
  overflow-wrap: normal;
}
.meta {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  gap: 6px;
  height: 20px;
  min-width: 0;
  overflow: hidden;
}
.chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  max-width: 100%;
  min-width: 0;
  height: 20px;
  padding: 0 8px;
  overflow: hidden;
  border-radius: 4px;
  color: #666;
  background: #f0f1f3;
  font-size: 12px;
  line-height: 20px;
  white-space: nowrap;
}
.chip-text {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.chip .dot {
  flex-shrink: 0;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #00b42a;
}
.chip.transfer {
  flex-shrink: 0;
  color: #a65b00;
  background: #fff4df;
}
.message-heading {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
  height: 78px;
  padding: 0 24px;
  border-bottom: 1px solid #e6e6e6;
}
.heading-text {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}
.heading-title {
  overflow: hidden;
  color: #242424;
  font-size: 18px;
  font-weight: 500;
  line-height: 26px;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.heading-sub {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
  overflow: hidden;
  color: #888;
  font-size: 12px;
  line-height: 18px;
}
.heading-sub-text {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.status-dot {
  flex-shrink: 0;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #00b42a;
}
.heading-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
  flex-shrink: 0;
}
.transfer-tag {
  display: inline-flex;
  align-items: center;
  height: 20px;
  padding: 0 8px;
  border-radius: 4px;
  color: #a65b00;
  background: #fff4df;
  font-size: 12px;
  line-height: 20px;
}
.customer-system {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 32px;
  padding: 0 12px;
  border-radius: 6px;
  color: #1456f0;
  background: #f0f4ff;
  font-size: 14px;
  font-weight: 500;
  cursor: not-allowed;
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

/* popovers */
:global(.session-channel-popover .ant-popover-inner),
:global(.session-date-popover .ant-popover-inner),
:global(.session-search-popover .ant-popover-inner),
:global(.session-transfer-popover .ant-popover-inner) {
  padding: 0;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}
:global(.session-channel-popover .ant-popover-inner-content),
:global(.session-date-popover .ant-popover-inner-content),
:global(.session-search-popover .ant-popover-inner-content),
:global(.session-transfer-popover .ant-popover-inner-content) {
  padding: 0;
}
:global(.session-channel-popover .ant-popover-arrow),
:global(.session-date-popover .ant-popover-arrow),
:global(.session-search-popover .ant-popover-arrow),
:global(.session-transfer-popover .ant-popover-arrow) {
  display: none;
}
.channel-menu {
  width: 256px;
  padding: 4px;
}
.channel-option {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  height: 50px;
  padding: 0 12px;
  border: 0;
  border-radius: 6px;
  color: #242424;
  background: transparent;
  font: inherit;
  font-size: 14px;
  text-align: left;
  cursor: pointer;

  &:hover {
    background: #f5f6f8;
  }

  &.active {
    color: #1456f0;
    background: #eef3ff;
  }
}
.channel-option-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  flex-shrink: 0;
  color: #6b7280;
  font-size: 16px;

  svg {
    width: 16px;
    height: 16px;
  }
}
.channel-option.active .channel-option-icon {
  color: #1456f0;
}
.channel-option-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.channel-option-count {
  color: #888;
  font-size: 14px;
}
.channel-option.active .channel-option-count {
  color: #242424;
}
.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}
.panel-title {
  color: #242424;
  font-size: 16px;
  font-weight: 500;
}
.panel-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  padding: 0;
  border: 0;
  color: #4e5969;
  background: transparent;
  font-size: 14px;
  cursor: pointer;
}
.simple-panel {
  padding: 20px 24px 24px;
}
.search-panel {
  width: 346px;
}
.search-panel .panel-head {
  margin-bottom: 16px;
}
.search-input {
  height: 40px;
  border-color: #e5e6eb;
  border-radius: 8px;
  font-size: 14px;
}
.search-input :deep(.ant-input) {
  font-size: 14px;
}
.search-input :deep(.ant-input-prefix) {
  margin-inline-end: 8px;
  color: #888;
}
.transfer-panel {
  width: 345px;
}
.transfer-panel .panel-head {
  margin-bottom: 12px;
}
.transfer-options {
  display: flex;
  flex-direction: column;
  cursor: not-allowed;
}
.transfer-option {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 44px;
  padding: 0 12px;
  border-bottom: 1px solid #f2f3f5;
  color: #242424;
  font-size: 16px;

  &:last-child {
    border-bottom: 0;
  }

  &.active {
    border-radius: 6px;
    color: #1456f0;
    background: #eef3ff;
  }
}
.radio {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border: 1.5px solid #c9cdd4;
  border-radius: 50%;
}
.radio.checked {
  border-color: #1456f0;

  &::after {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #1456f0;
    content: '';
  }
}
.transfer-hint {
  margin-top: 12px;
  color: #888;
  font-size: 12px;
  line-height: 18px;
}

@media (max-width: 1280px) {
  .session-list-panel {
    width: 300px;
  }
}
</style>
