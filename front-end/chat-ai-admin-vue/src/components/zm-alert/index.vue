<template>
  <div class="zm-alert zm-alert-info" :class="[`zm-alert-${props.type}`]">
    <span class="zm-alert-icon">
      <slot name="icon"><ExclamationCircleFilled /></slot>
    </span>

    <div class="zm-alert-content">
      <div class="zm-alert-title">
        <slot name="title">{{ props.title }}</slot>
      </div>

      <div class="zm-alert-message">
        <slot name="message">{{ props.message }}</slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ExclamationCircleFilled } from '@ant-design/icons-vue'

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'info'
  }
})
</script>

<style lang="less" scoped>
.zm-alert {
  display: flex;
  align-items: start;
  padding: 9px 16px;

  .zm-alert-icon {
    height: 22px;
    line-height: 22px;
    margin-right: 12px;
    font-size: 16px;
  }
  .zm-alert-title {
    line-height: 22px;
    font-size: 14px;
    font-weight: 600;
    color: #242933;
  }
  .zm-alert-message {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: #3a4559;
  }
}

.zm-alert.zm-alert-info {
  border: 1px solid #91caff;
  border-radius: 2px;
  background-color: #e6f4ff;

  .zm-alert-icon {
    color: #2475fc;
  }
}
</style>
