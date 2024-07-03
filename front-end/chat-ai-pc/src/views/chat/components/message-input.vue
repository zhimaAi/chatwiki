<style lang="less" scoped>
.message-input-box {
  display: flex;
  padding: 12px 28px 8px 28px;

  .message-input-body {
    flex: 1;
  }

  .message-input {
    position: relative;
    padding: 8px 0 8px 0;
    border-radius: 16px;
    border: 1px solid #f8f9fa;
    background: #fff;
    border: 1px solid #2475fc;
    transition: all 0.2s;

    .text-input {
      border: none;
      width: 100%;
      font-size: 16px;
      font-weight: 400;
      padding: 0 32px 0 12px;
      color: #1a1a1a;
      background: none;

      &::placeholder {
        font-size: 14px;
        font-weight: 400;
        color: #bfbfbf;
      }
    }

    .send-btn {
      position: absolute;
      right: 12px;
      bottom: 8px;
      font-size: 20px;
      color: #b3b3b3;
      transition: all 0.2s;
      color: #2475fc;
      cursor: pointer;
    }

    &.is-focus {
      box-shadow: 0 2px 6px 0 rgba(0, 0, 0, 0.12);
    }
  }
}
</style>

<template>
  <div class="message-input-box">
    <div class="message-input-body">
      <div class="message-input" :class="{ 'is-focus':  isFocus }">
        <AutoSizeTextarea :value="value" class="text-input"  @change="onChange" @focus="isFocus = true" @blur="isFocus = false" @enter="sendMessage" />
        <svg-icon name="paper-airplane" class="send-btn" @click="sendMessage" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AutoSizeTextarea from './auto-size-textarea.vue'

const emit = defineEmits(['update:value', 'send'])

const props = defineProps({
  value: {
    type: String,
    default: ''
  }
})

const isFocus = ref(false)

const onChange = ((val: string) => {
  emit('update:value', val)
})

const sendMessage = () => {
  if (props.value.trim()) {
    emit('send', props.value)
  }
}
</script>
