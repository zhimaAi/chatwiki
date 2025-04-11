<style lang="less" scoped>
.chat-monitor-page {
  display: flex;
  overflow: hidden;
  width: 100%;
  height: 100%;
  border-radius: 6px;
  border: 1px solid #fff;
}
.page-left {
  display: flex;
  flex-direction: column;
  width: 352px;
  height: 100%;
  overflow: hidden;
  background-color: #fff;
  .app-list-box {
    padding: 16px;
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
          <a-select-option value="">全部应用</a-select-option>
          <a-select-option :value="item.id" v-for="item in robotList" :key="item.id">
            <span>{{ item.robot_name }}</span>
          </a-select-option>
        </a-select>
      </div>

      <div class="chat-list-wrapper">
        <ChatList @switchChat="handleSwitchChat" />
      </div>
    </div>
    <div class="page-body">
      <ChatBox v-if="activeChat" ref="chatBoxRef" />
      <list-empty size="250" v-else>
        <div>
          <p>请在左侧列表先选择会话</p>
          <p>通过本功能，可以实时查看机器人接待中的会话消息，监控机器人回复效果</p>
        </div>
      </list-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { storeToRefs } from 'pinia'
import { useEventBus } from '@/hooks/event/useEventBus'
import { useChatMonitorStore } from '@/stores/modules/chat-monitor.js'
import ChatList from './components/chat-list .vue'
import ChatBox from './components/chat-box.vue'
import ListEmpty from './components/list-empty.vue'

const emitter = useEventBus()

const chatMonitorStore = useChatMonitorStore()
const { init, changeRobot, switchChat, closeIM } = chatMonitorStore
const { robotList, selectedRobotId, activeChat } = storeToRefs(chatMonitorStore)

const chatBoxRef = ref(null)

const onChangeRobot = () => {
  changeRobot()
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
