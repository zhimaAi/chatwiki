<style lang="less" scoped>
.message-input-box {
  width: 100%;
  margin-bottom: 0px;

  .message-input-body {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .message-input {
    position: relative;
    z-index: 2;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex: 1;
    min-height: 60px;
    max-width: 738px;
    min-width: 350px;
    border-radius: 16px;
    border: 1px solid #ddd;
    background: #fff;
    box-shadow: 0 2px 6px 0 rgba(0, 0, 0, 0.08);
    transition: all 0.2s;
    padding: 10px 0;
    margin: 0 12px;

    .send-btn {
      position: absolute;
      bottom: 13px;
      right: 10px;
      font-size: 32px;
      color: #b3b3b3;
      transition: all 0.2s;
    }
    
    .send-btn-active {
      color: #2475fc;
      cursor: pointer;
    }

    &.is-focus {
      border: 1px solid #2475fc;
    }
  }
}
</style>

<template>
  <div class="message-input-box">
    <div class="message-input-body">
      <div class="message-input" :class="{ 'is-focus':  isFocus }">
        <AutoSizeTextarea :value="value" @change="onChange" @focus="onFocus" @blur="onBlur" @enter="sendMessage" />
        <svg-icon :name="isSendBtn ? 'paper-airplane-new-active' : 'paper-airplane-new'" :class="{'send-btn-active': isSendBtn}" class="send-btn" @click="sendMessage" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AutoSizeTextarea from './auto-size-textarea.vue'

const props = defineProps({
  value: {
    type: String,
    default: ''
  }
})

const isFocus = ref(false)

const isSendBtn = ref(false)

const emit = defineEmits(['update:value', 'send'])

const onChange = ((val: string) => {
  // 会有延迟
  if (val.trim()) {
    isSendBtn.value = true
  } else {
    isSendBtn.value = false
  }
  emit('update:value', val)
})

const sendMessage = () => {
  if (props.value.trim()) {
    isSendBtn.value = false
    emit('send', props.value)
  }
}

const onFocus = (event) => {
  isFocus.value = true
  event.target.parentNode.style.borderColor = "#2475FC"
}

const onBlur = (event) => {
  isFocus.value = false
  event.target.parentNode.style.borderColor = "#DDD"
}

</script>
