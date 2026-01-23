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
  width: 352px;
  height: 100%;
  overflow: hidden;
  background-color: #fff;
  .app-list-box {
    padding: 16px 16px 8px;
  }

  .search-box {
    padding: 0 16px 16px;
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
        <a-input
          v-model:value="keyword"
          :placeholder="t('search_placeholder')"
          enter-button
          style="width: 100%;"
          @search="onSearch"
          allowClear
          @pressEnter="onSearch"
        />
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

const { t } = useI18n('views.chat-monitor')
const emitter = useEventBus()

const chatMonitorStore = useChatMonitorStore()
const { init, changeRobot, switchChat, closeIM } = chatMonitorStore
const { robotList, selectedRobotId, activeChat } = storeToRefs(chatMonitorStore)

const chatListRef = ref(null)
const chatBoxRef = ref(null)
const keyword = ref('')
const onChangeRobot = () => {
  const params = {
    keyword: keyword.value,
    page: 1
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
    page: 1
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