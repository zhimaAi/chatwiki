<style lang="less" scoped>
.chat-monitor-page {
  display: flex;
  overflow: hidden;
  width: 100%;
  height: 100%;
  border-radius: 6px;
}
.page-left {
  display: flex;
  flex-direction: column;
  width: 368px;
  height: 100%;
  overflow: hidden;
  background-color: #fff;
  .app-list-box {
    padding: 16px 16px 8px;
  }

  .search-box {
    padding: 0 16px 16px;
    .search-input-row {
      display: flex;
      align-items: center;
      gap: 8px;
    }
    .time-filter-btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 32px;
      height: 32px;
      border-radius: 6px;
      cursor: pointer;
      background-color: #E4E6EB;
      transition: background-color 0.2s ease, color 0.2s ease;
      &:hover {
        background-color: rgba(242, 244, 247, 1);
        color: #262626;
      }
      .time-icon {
        font-size: 20px;
        color: #595959;
      }
      &.active {
        background-color: #E5EFFF;
        .time-icon {
          color: #2475FC;
        }
      }
    }
    .date-filter-wrapper {
      margin-top: 8px;
      padding: 8px;
      border-radius: 6px;
      background-color: #fff;
    }
  }

  .chat-list-wrapper {
    flex: 1;
    overflow: hidden;
  }
}
.page-body {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  width: 100%;
  height: 100%;

  
}
</style>

<template>
  <div class="chat-monitor-page">
    <div class="page-left">
      <div class="app-list-box">
        <a-select v-model:value="selectedRobotId" style="width: 100%" @change="onChangeRobot">
          <a-select-option value="">{{ t('all_apps') }}</a-select-option>
          <a-select-option :value="item.id" v-for="item in robotList" :key="item.id">
            <span>{{ item.robot_name }}</span>
          </a-select-option>
        </a-select>
      </div>

      <div class="search-box">
        <div class="search-input-row">
          <a-input
            v-model:value="keyword"
            :placeholder="t('search_placeholder')"
            enter-button
            style="flex: 1;"
            @search="onSearch"
            allowClear
            @pressEnter="onSearch"
          />
          <div class="time-filter-btn" :class="{ active: showDateFilter }" @click="toggleDateFilter">
            <svg-icon name="time-line-icon" class="time-icon"></svg-icon>
          </div>
        </div>
        <div class="date-filter-wrapper" v-if="showDateFilter">
          <DateFilter :datekey="defaultDateKey" @dateChange="onDateChange" />
        </div>
      </div>

      <div class="chat-list-wrapper">
        <ChatList ref="chatListRef" @switchChat="handleSwitchChat" />
      </div>
    </div>
    <div class="page-body">
      <ChatBox v-if="activeChat" ref="chatBoxRef" />
      <list-empty style="background: #F2F4F7;" size="250" v-else>
        <div>
          <p>{{ t('select_conversation_tip') }}</p>
          <p>{{ t('feature_description') }}</p>
        </div>
      </list-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { storeToRefs } from 'pinia'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useChatMonitorStore } from '@/stores/modules/chat-monitor.js'
import ChatList from './components/chat-list .vue'
import ChatBox from './components/chat-box.vue'
import ListEmpty from './components/list-empty.vue'
import DateFilter from './components/date-filter.vue'
import dayjs from 'dayjs'

const { t } = useI18n('views.chat-monitor')
const emitter = useEventBus()

const chatMonitorStore = useChatMonitorStore()
const { init, changeRobot, switchChat, closeIM } = chatMonitorStore
const { robotList, selectedRobotId, activeChat } = storeToRefs(chatMonitorStore)

const chatListRef = ref(null)
const chatBoxRef = ref(null)
const keyword = ref('')
const showDateFilter = ref(false)
const defaultDateKey = ref('2')
const start_time = ref('')
const end_time = ref('')

const toggleDateFilter = () => {
  showDateFilter.value = !showDateFilter.value
}

const onDateChange = ({ start_time: s, end_time: e }) => {
  start_time.value = dayjs(String(s), 'YYYYMMDD').startOf('day').unix()
  end_time.value = dayjs(String(e), 'YYYYMMDD').endOf('day').unix()
  const params = {
    keyword: keyword.value,
    page: 1,
    start_time: start_time.value,
    end_time: end_time.value
  }
  if (chatListRef.value) {
    chatListRef.value.getData(params)
  }
}
const onChangeRobot = () => {
  const params = {
    keyword: keyword.value,
    page: 1,
    start_time: start_time.value || undefined,
    end_time: end_time.value || undefined
  }
  changeRobot(params)
}


const handleSwitchChat = async (item) => {
  await switchChat(item)

  chatBoxRef.value?.scrollToBottom()
}

const onAddMessage = () => {
  nextTick(() => {
    chatBoxRef.value?.scrollToBottom()
  })
}

const onSearch = () => {
  const params = {
    keyword: keyword.value,
    page: 1,
    start_time: start_time.value || undefined,
    end_time: end_time.value || undefined
  }
  if (chatListRef.value) {
    chatListRef.value.getData(params)
  }
}

onMounted(async () => {
  await init()

  nextTick(() => {
    chatBoxRef.value?.scrollToBottom()
  })

  // 监听 updateAiMessage 触发消息列表滚动
  emitter.on('onAddMessage', onAddMessage)
})

onUnmounted(() => {
  emitter.off('onAddMessage', onAddMessage)

  closeIM()
})
</script>
