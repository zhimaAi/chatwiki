<template>
  <div class="pc-side-content" :class="{ 'is-hide': sidebarHide }">
    <div class="sidebar-handle-wrapper">
      <span class="sidebar-handle" @click="onHandleClick">
        <span class="handle-line handle-line01"></span>
        <span class="handle-line handle-line02"></span>
      </span>
    </div>
    <div class="sidebar-container">
      <SideHeader
        @openNewChat="openNewChat"
        @emptyAllChat="emptyAllChat"
        :isMobileDevice="false"
        :sidebarHide="sidebarHide"
      />
      <SessionList @handleOpenChat="handleOpenChat" v-if="!sidebarHide" />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

import SideHeader from './side-header.vue'
import SessionList from './session-list.vue'
let defaultData = localStorage.getItem('sidebar_hide_mobile') === null || localStorage.getItem('sidebar_hide_mobile') == 1
const sidebarHide = ref(defaultData)
const emit = defineEmits(['handleOpenChat', 'openNewChat', 'emptyAllChat'])

const handleOpenChat = (item) => {
  emit('handleOpenChat', item)
}

const openNewChat = () => {
  emit('openNewChat')
}

const emptyAllChat = () => {
  emit('emptyAllChat')
}
function onHandleClick() {
  sidebarHide.value = !sidebarHide.value
  localStorage.setItem('sidebar_hide_mobile', sidebarHide.value ? 1 : 0)
}
</script>

<style lang="less" scoped>
.pc-side-content {
  width: 298px;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition: width 0.2s ease;
  position: relative;
  .sidebar-container {
    width: 280px;
    height: 100%;
    overflow: hidden;
    padding: 16px;
    display: flex;
    flex-direction: column;
    background: #fff;
  }

  &:hover {
    .sidebar-handle {
      opacity: 1;
    }
  }

  .sidebar-handle-wrapper {
    position: absolute;
    top: 0;
    right: 0;
    width: 18px;
    height: 100%;
    z-index: 100;
  }
  .sidebar-handle {
    position: absolute;
    right: 0;
    top: 50%;
    width: 12px;
    height: 26px;
    transform: translateY(-50%);
    cursor: pointer;
    transition: all 0.2s ease;
    opacity: 0;
    // opacity: 1;

    .handle-line {
      position: absolute;
      width: 4px;
      height: 13px;
      left: 4px;
      position: absolute;
      transition: all 0.2s ease;
      background-color: #bfbfbf;
    }

    .handle-line01 {
      top: 0;
      border-top-left-radius: 4px;
      border-top-right-radius: 4px;
      transform-origin: 50% 0;
    }

    .handle-line02 {
      bottom: 0;
      border-bottom-left-radius: 4px;
      border-bottom-right-radius: 4px;
      transform-origin: 50% 100%;
    }
  }

  .sidebar-handle:hover {
    .handle-line01 {
      background-color: #595959;
      transform: rotate(18deg) translateY(0);
      border-top-left-radius: 4px;
      border-top-right-radius: 4px;
      border-bottom-left-radius: 10px;
      height: 16px;
    }

    .handle-line02 {
      background-color: #595959;
      transform: rotate(-18deg) translateY(0);
      border-bottom-left-radius: 4px;
      border-bottom-right-radius: 4px;
      border-top-left-radius: 10px;
      height: 16px;
    }
  }
}

.pc-side-content.is-hide {
  width: 82px;
  .sidebar-container {
    width: 64px;
  }

  .sidebar-handle {
    opacity: 1 !important;
  }

  .sidebar-handle:hover {
    .handle-line01 {
      background-color: #595959;
      transform: rotate(-18deg) translateY(0);
      border-bottom-left-radius: 4px;
      border-bottom-right-radius: 4px;
      border-top-left-radius: 10px;
      height: 16px;
    }

    .handle-line02 {
      background-color: #595959;
      transform: rotate(18deg) translateY(0);
      border-top-left-radius: 4px;
      border-top-right-radius: 4px;
      border-bottom-left-radius: 10px;
      height: 16px;
    }
  }
}
</style>
