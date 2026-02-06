<template>
  <div class="header-box">
    <template v-if="props.isMobileDevice || !props.sidebarHide">
      <div class="logo-box">
        <img class="logo" src="@/assets/logo.svg" alt="" />
        <div class="left-btn" @click="handleClose" v-if="isMobileDevice">
          <img src="@/assets/icons/left-arrow.svg" alt="" />
        </div>
      </div>
      <div class="btn-box">
        <div style="flex: 1">
          <van-button block @click="openNewChat" type="primary" :size="size">{{ t('btn_new_chat') }}</van-button>
        </div>
        <van-button @click="emptyAllChat" :size="size">{{ t('btn_clear_history') }}</van-button>
      </div>
    </template>
    <template v-else>
      <div class="short-header">
        <div class="logo-box">
          <img src="@/assets/logo.png" alt="" />
        </div>
        <div class="new-chat-box" @click="openNewChat" v-tooltip="t('btn_new_chat')">
          <img src="@/assets/icons/add-chat.svg" alt="" />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Button } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.chat.components.left-side-bar.side-header')

const emit = defineEmits(['openNewChat', 'emptyAllChat', 'handleClose'])
const props = defineProps({
  isMobileDevice: {
    default: false,
    type: Boolean
  },
  sidebarHide: {
    default: true,
    type: Boolean
  }
})

const size = computed(() => {
  if (props.isMobileDevice) {
    return 'normal'
  }
  return 'small'
})

const openNewChat = () => {
  emit('openNewChat')
}

const emptyAllChat = () => {
  emit('emptyAllChat')
}

const handleClose = () => {
  emit('handleClose')
}
</script>

<style lang="less" scoped>
.header-box {
  .logo-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    .logo {
      width: 116px;
    }
    .left-btn {
      width: 40px;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: var(--07, #f0f2f5);
      border-radius: 8px;
      cursor: pointer;
      img {
        width: 24px;
        height: 24px;
      }
    }
  }
  .btn-box {
    margin-top: 16px;
    display: flex;
    align-items: center;
    gap: 8px;
  }
}
.short-header {
  .logo-box {
    margin-bottom: 16px;
    img {
      width: 32px;
    }
  }
  .new-chat-box {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    img {
      width: 20px;
    }
  }
}
</style>
