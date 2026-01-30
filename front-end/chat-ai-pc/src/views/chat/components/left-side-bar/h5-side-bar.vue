<template>
  <div>
    <van-popup v-model:show="showLeft" position="left" :style="{ width: '85%', height: '100%' }">
      <div class="h5-side-content">
        <SideHeader
          @openNewChat="openNewChat"
          @emptyAllChat="emptyAllChat"
          @handleClose="handleClose"
          :isMobileDevice="true"
          :sidebarHide="false"
        />
        <SessionList @handleOpenChat="handleOpenChat" />
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref } from 'vue'

import SideHeader from './side-header.vue'
import SessionList from './session-list.vue'

const emit = defineEmits(['handleOpenChat', 'openNewChat', 'emptyAllChat'])

const showLeft = ref(false)
const handleOpenChat = (item) => {
  emit('handleOpenChat', item)
  handleClose()
}

const openNewChat = () => {
  emit('openNewChat')
  handleClose()
}

const emptyAllChat = () => {
  emit('emptyAllChat')
  handleClose()
}

const handleShow = () => {
  showLeft.value = true
}

const handleClose = () => {
  showLeft.value = false
}

defineExpose({
  handleShow
})
</script>

<style lang="less" scoped>
.h5-side-content {
  width: 100%;
  height: 100%;
  overflow: hidden;
  padding: 16px;
  display: flex;
  flex-direction: column;
}
</style>
