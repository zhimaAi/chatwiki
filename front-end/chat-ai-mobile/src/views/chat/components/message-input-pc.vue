<template>
  <div class="message-input-wrapper"  :class="{ 'is-set': props.value }">
    <FileToolbar :file-list="fileList" @delete="deleteFile" v-if="fileList.length > 0" />
    <div class="message-input-box">
      <ATextarea
        class="message-input"
        :value="props.value"
        :auto-size="{ minRows: 2, maxRows: 5 }"
        :placeholder="t('ph_input_message_with_shift')"
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
        <ASpin size="small" class="loading-action" style="margin-right: 4px" v-if="props.loading" />
        <svg-icon class="paper-airplane" name="paper-airplane-new-active" v-else />
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, toRefs, computed } from 'vue'
import { useChatStore } from '@/stores/modules/chat'
import { useUserStore } from '@/stores/modules/user'
import { Textarea as ATextarea, Spin as ASpin } from 'ant-design-vue'
import { showToast } from 'vant'
import { useI18n } from '@/hooks/web/useI18n'
import { useUpload } from '@/hooks/web/useUpload.js'
import { checkChatRequestPermission } from '@/api/robot/index'
import FileToolbar from './file-toolbar.vue'

const chatStore = useChatStore()
const userStore = useUserStore()
const { robot } = chatStore


const emit = defineEmits(['update:value', 'send', 'showLogin', 'update:fileList'])

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

const { t } = useI18n('views.chat.components.message-input-pc')

const { openFileDialog } = useUpload({
  limit: 10,
  maxSize: 10,
  category: 'chat_image',
  fileList: fileList,
  multiple: true,
  accept: 'image/bmp,image/jpeg,image/png,image/tiff,image/heic,image/gif,image/webp',
  extraData: {
    robot_key: robot.robot_key,
    openid: robot.openid
  }
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

const onChange = (event) => {
  emit('update:value', event.target.value)
}

const sendMessage = async () => {
  emit('send', props.value)
}

const handleKeydown = (event) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    if (!event.target.value) {
      return
    }
    event.preventDefault()
    event.stopPropagation()
    sendMessage()
  } else if (event.key === 'Enter' && event.shiftKey) {
    emit('update:value', event.target.value)
  }
}

const handleSetValue = (data) => {
  emit('update:value', data)
}

defineExpose({
  handleSetValue,
  sendMessage
})
</script>


<style lang="less" scoped>
.message-input-wrapper {
  position: relative;
  max-width: 736px;
  margin: 0 auto;
  padding: 12px;
  border-radius: 20px;
  border: 1px solid #d9d9d9;
  overflow: hidden;
  background-color: #fff;
  transition: all 0.2s;

  &.is-set {
    border: 1px solid #2475fc;
  }

  .message-input-box {
    padding: 0 0 12px 0;
  }

  .message-input {
    width: 100%;
    padding: 0 12px;
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
    padding: 0 12px 12px 12px;
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