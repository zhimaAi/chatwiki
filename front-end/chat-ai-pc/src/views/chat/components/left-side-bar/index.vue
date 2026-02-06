<template>
  <div class="side-bar-box">
    <H5SizeBar
      @handleOpenChat="handleOpenChat"
      @openNewChat="openNewChat"
      @emptyAllChat="emptyAllChat"
      v-if="props.isMobileDevice"
      ref="h5SizeBarRef"
    />
    <PcSideBar
      @handleOpenChat="handleOpenChat"
      @openNewChat="openNewChat"
      @emptyAllChat="emptyAllChat"
      v-else
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import PcSideBar from './pc-side-bar.vue'
import H5SizeBar from './h5-side-bar.vue'
import { showConfirmDialog } from 'vant'
import { useChatStore } from '@/stores/modules/chat'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.chat.components.left-side-bar.index')
const chatStore = useChatStore()

const emit = defineEmits(['openChat', 'openNewChat'])

const props = defineProps({
  isMobileDevice: {
    type: Boolean,
    default: false
  }
})

const handleOpenChat = (item) => {
  emit('openChat', item)
}

const openNewChat = () => {
  emit('openNewChat')
}

const emptyAllChat = () => {
  showConfirmDialog({
    title: t('title_confirm_clear'),
    message: t('msg_confirm_clear_history')
  })
    .then(() => {
      chatStore.delDialogue({
        id: -1
      })
    })
    .catch(() => {})
}

const h5SizeBarRef = ref(null)
const handleShowH5Chat = () => {
  h5SizeBarRef.value.handleShow()
}

defineExpose({
  handleShowH5Chat
})
</script>

<style lang="less" scoped>
.side-bar-box {
  height: 100%;
  overflow: hidden;
  background: #F0F2F5;
}
</style>
