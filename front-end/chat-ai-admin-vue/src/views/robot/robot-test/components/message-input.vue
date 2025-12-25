<template>
  <div class="message-input-wrapper">
    <div class="file-toolbar-box" v-if="props.fileList.length > 0">
      <FileToolbar :file-list="props.fileList" @delete="deleteFile" />
    </div>
    <div class="message-input-box">
      <a-textarea
        class="message-input"
        :value="props.value"
        placeholder="在此输入您想了解的内容，Shift+Enter换行"
        :auto-size="{ minRows: 2, maxRows: 5 }"
        @change="onChange"
        @keydown="handleKeydown"
      />
    </div>
    <div class="message-action">
      <div class="select-file-btn" @click="openFileDialog" v-if="props.showUpload">
        <svg-icon class="select-file-icon" name="circularNeedle"></svg-icon>
        <span class="file-number" :class="{ big: fileList.length > 9 }" v-if="fileList.length > 0">{{ fileList.length }}</span>
      </div>

      <button
        class="send-msg-btn"
        :class="{ loading: props.loading }"
        :disabled="disabled"
        @click="sendMessage"
      >
        <a-spin size="small" class="loading-action" style="margin-right: 4px" v-if="props.loading" />
        <svg-icon class="paper-airplane" name="paper-airplane" v-else></svg-icon>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed, toRefs } from 'vue'
import FileToolbar from './file-toolbar.vue'
import { useUpload } from '@/hooks/web/useUpload.js'

const emit = defineEmits(['send', 'update:value', 'update:fileList'])

const props = defineProps({
  value: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  },
  fileList: {
    type: Array,
    default: () => []
  },
  showUpload: {
    type: Boolean,
    default: false
  },
})

const { fileList } = toRefs(props)

const { openFileDialog } = useUpload({
  limit: 10,
  maxSize: 10,
  category: 'chat_image',
  fileList: fileList,
  multiple: true,
  accept: 'image/bmp,image/jpeg,image/png,image/tiff,image/heic,image/gif,image/webp'
})

const deleteFile = (index) => {
  const newFileList = props.fileList.filter((_, i) => i !== index);
  emit('update:fileList', newFileList);
}

const disabled = computed(() => {
  if(props.fileList.length > 0){
    return false
  }

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
    event.stopPropagation()
    sendMessage()
  } else if (event.key === 'Enter' && event.shiftKey) {
    // 同时按下了 Shift 和 Enter 键
    event.preventDefault()
    event.stopPropagation()
    emit('update:value', props.value + '\n')
  }
}

const sendMessage = () => {
  emit('send', props.value.trim())
}
</script>

<style lang="less" scoped>
.message-input-wrapper {
  position: relative;
  border-radius: 20px;
  border: 1px solid #d9d9d9;
  overflow: hidden;
  padding: 12px;
  background-color: #fff;

  .file-toolbar-box{
    margin-bottom: 8px;
  }

  .message-input-box {
    padding: 0 0 12px 0;
  }

  .message-input {
    width: 100%;
    padding: 0;
    line-height: 24px;
    border: none !important;
    outline: none !important;
    resize: none !important;
    box-shadow: none !important;
  }

  .message-action {
    display: flex;
    align-items: center;
    justify-content: flex-end;
  }

  .send-msg-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    padding: 0;
    font-size: 14px;
    font-weight: 400;
    border-radius: 4px;
    border: none;
    cursor: pointer;
    transition: all 0.2s;
    color: #2475fc;
    border-radius: 50%;
    background: none;

    &:hover {
      opacity: 0.8;
    }
    &:disabled {
      opacity: 0.5;
    }
    .paper-airplane {
      font-size: 32px;
    }

    &.loading {
      background-color: #2475fc;
    }
    .loading-action{
      display: flex;
      align-items: center;
      justify-content: center;
      margin-left: 3px;
      ::v-deep(.ant-spin-dot-item) {
        background-color: #fff;
      }
    }
    
  }

  .select-file-btn {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    padding: 0;
    margin-right: 8px;
    border-radius: 50%;
    border: none;
    background: #fff;
    cursor: pointer;
    transition: all 0.2s;
    border: 1px solid #f0f0f0;

    &:hover {
      background: #e4e6eb;
    }

    .select-file-icon {
      font-size: 16px;
      color: #595959;
    }

    .file-number{
      position: absolute;
      right: -8px;
      top: -8px;
      width: 16px;
      height: 16px;
      border-radius: 50%;
      background: #f00;
      color: #fff;
      font-size: 12px;
      font-weight: 400;
      display: flex;
      align-items: center;
      justify-content: center;

      &.big{
        width: auto;
        padding: 0 4px;
        border-radius: 12px;
      }
    }
  }
}
</style>