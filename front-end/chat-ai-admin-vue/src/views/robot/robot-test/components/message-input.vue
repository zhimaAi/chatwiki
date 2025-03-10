<style lang="less" scoped>
.message-input-wrapper {
  position: relative;

  .message-action {
    position: absolute;
    right: 16px;
    bottom: 12px;
  }

  .send-msg-btn {
    width: 60px;
    height: 32px;
    line-height: 32px;
    padding: 0;
    font-size: 14px;
    font-weight: 400;
    border-radius: 4px;
    border: none;
    cursor: pointer;
    transition: all 0.2s;
    color: #2475fc;
    background-color: #e6efff;

    &:hover {
      opacity: 0.7;
    }
  }

  .loading-action {
    display: block;
    height: 32px;
    line-height: 32px;
    padding-top: 6px;
    text-align: center;
  }
}
</style>

<template>
  <div class="message-input-wrapper">
    <a-textarea
      :value="value"
      :auto-size="{ minRows: 5, maxRows: 5 }"
      placeholder="在此输入您想了解的内容，Shift+Enter换行"
      @change="onChange"
      @keydown="handleKeydown"
    />
    <div class="message-action">
      <span class="loading-action">
        <a-spin :spinning="loading"></a-spin>
      </span>

      <button class="send-msg-btn" :disabled="disabled" @click="sendMesage" v-if="!loading">
        <span>发送</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const emit = defineEmits(['send', 'update:value'])

const props = defineProps({
  value: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const disabled = computed(() => {
  return props.loading || props.value.trim().length === 0
})

const onChange = (e) => {
  emit('update:value', e.target.value.trim())
}

const handleKeydown = (event) => {
  // 键盘 enter事件 和 enter+shift 组合键 绑定事件
  if (event.key === 'Enter' && !event.shiftKey) {
    // 只按下了 Enter 键
    if (!event.target.value.trim()) {
      return
    }
    event.preventDefault()
    sendMesage()
  } else if (event.key === 'Enter' && event.shiftKey) {
    // 同时按下了 Shift 和 Enter 键
    event.preventDefault()
    emit('update:value', props.value + '\n')
  }
}

const sendMesage = () => {
  emit('send', props.value.trim())
}
</script>
