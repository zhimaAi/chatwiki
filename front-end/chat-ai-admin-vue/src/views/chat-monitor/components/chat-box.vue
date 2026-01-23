<style lang="less" scoped>
.chat-box-warpper {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #F2F4F7;
}
.chat-box-header {
  padding: 12px 16px;
  .nickname {
    font-size: 14px;
    line-height: 22px;
    color: rgb(36, 41, 51);
    text-align: center;
  }
  .chat-source {
    line-height: 20px;
    margin-top: 4px;
    font-size: 12px;
    color: rgb(122, 134, 153);
    text-align: center;
  }
}
.chat-box-content {
  flex: 1;
  overflow: hidden;
}
</style>

<template>
  <div class="chat-box-warpper" v-if="activeChat">
    <div class="chat-box-header">
      <div class="nickname">{{ activeChat.name || activeChat.nickname }}</div>
      <div class="chat-source">
        {{ t('from_source', { robot: activeChat.come_from.robot_name, app: activeChat.come_from.app_name }) }}
      </div>
    </div>
    <div class="chat-box-content">
      <ChatMessage ref="chatMessageRef" @openLibrary="handleOpenLibraryInfo" />
    </div>
  </div>
  <div v-else></div>
  <LibraryInfoAlert ref="libraryInfoAlertRef" />
</template>

<script setup>
import { ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useChatMonitorStore } from '@/stores/modules/chat-monitor.js'
import { useI18n } from '@/hooks/web/useI18n'
import ChatMessage from './chat-message.vue'
import LibraryInfoAlert from './library-info-alert.vue'

const { t } = useI18n('views.chat-monitor.components.chat-box')

const chatMonitorStore = useChatMonitorStore()

const { activeChat } = storeToRefs(chatMonitorStore)

const chatMessageRef = ref(null)

const scrollToBottom = () => {
  chatMessageRef.value.scrollToBottom()
}

const scrollToTop = () => {
  chatMessageRef.value.scrollToTop()
}

const libraryInfoAlertRef = ref(null)
const handleOpenLibraryInfo = (files, file) => {
  libraryInfoAlertRef.value.open(files, file)
}

defineExpose({
  scrollToTop,
  scrollToBottom
})
</script>