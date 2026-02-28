<template>
  <div class="input-area">
    <div class="input-wrapper" :class="{ focus: isFocused }">
      <!-- 文件工具栏 -->
      <FileToolbar
        v-if="fileList.length > 0"
        :fileList="fileList"
        @delete="handleDeleteFile"
      />
      <a-textarea v-model:value="inputValue" :auto-size="{ minRows: 2, maxRows: 5 }" class="chat-input"
        :placeholder="t(placeholder)" @keydown="handleKeydown" @focus="handleFocus" @blur="handleBlur" />
      <div class="input-actions-wrapper">
        <div class="left-actions actions-box">
          <span class="action-btn file-action" v-if="showUpload" @click="openFileDialog">
            <svg-icon name="circularNeedle" />
            <span class="file-number" :class="{ big: fileList.length > 9 }" v-if="fileList.length > 0">{{ fileList.length }}</span>
          </span>
        </div>
        
        <div class="right-actions actions-box">
          <IconSend class="send-action" :size="32" @click="handleSend" :isActive="sendBtnIsActive" />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Textarea as ATextarea } from 'ant-design-vue'
import { ref, watch, computed, toRefs } from 'vue'
import { useChatStore } from '@/stores/modules/chat'
import { useUpload } from '@/hooks/web/useUpload.js'
import IconSend from './Icon/IconSend.vue'
import FileToolbar from '@/views/chat/components/file-toolbar.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.pc.components.chat-input')
const chatStore = useChatStore()
const { robot } = chatStore

interface Props {
  modelValue?: string
  placeholder?: string
  showUpload?: boolean
  fileList?: any[]
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: 'ph_input_question',
  showUpload: false,
  fileList: () => []
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'update:fileList': [value: any[]]
  send: [value: string]
}>()

const inputValue = ref(props.modelValue)
const { fileList } = toRefs(props)
const isFocused = ref(false)

const sendBtnIsActive = computed(() => {
  if (fileList.value.length > 0) {
    return true
  }
  return inputValue.value.trim() !== ''
})

watch(() => props.modelValue, (newVal) => {
  inputValue.value = newVal
})

watch(inputValue, (newVal) => {
  emit('update:modelValue', newVal)
})

// 使用 useUpload Hook
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

const handleDeleteFile = (index: number) => {
  const newFileList = props.fileList.filter((_, i) => i !== index)
  emit('update:fileList', newFileList)
}

const handleSend = () => {
  if (inputValue.value.trim() || fileList.value.length > 0) {
    emit('send', inputValue.value)
  }
}

const handleKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    const target = event.target as HTMLTextAreaElement
    if (!target.value && fileList.value.length === 0) {
      return
    }

    event.preventDefault()
    event.stopPropagation()

    handleSend()
  } else if (event.key === 'Enter' && event.shiftKey) {
    // emit('update:value', event.target.value)
  }
}

const handleFocus = () => {
  isFocused.value = true
}

const handleBlur = () => {
  isFocused.value = false
}
</script>

<style lang="less" scoped>
.input-area {
  max-width: 100%;
  margin: 0 auto;
}

.input-wrapper {
  padding: 10px 12px;
  border-radius: 16px;
  border: 2px solid #E5E7EB;
  transition: border-color 0.2s, box-shadow 0.2s;
  cursor: text;
}
.input-wrapper.focus {
  border-color: #659DFC;
  box-shadow: 0 0 6px 6px rgba(101, 157, 252, 0.2);
}
.chat-input {
  width: 100%;
  line-height: 24px;
  padding: 0;
  border: none;
  outline: none;
  font-size: 14px;
  color: #333;
  background: transparent;
  margin-bottom: 24px;
  resize: none;
   line-height: 24px;
  border: none !important;
  outline: none !important;
  resize: none !important;
  box-shadow: none !important;

  &::placeholder {
    color: #bfbfbf;
  }
}

.input-actions-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;

  .actions-box {
    display: flex;
    align-items: center;
  }

  .action-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background-color: #E4E6EB;
    }
  }

  .file-action {
    position: relative;
    font-size: 16px;
    border-radius: 50%;
    border: 1px solid #D8DDE5;

    .file-number {
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

      &.big {
        width: auto;
        padding: 0 4px;
        border-radius: 12px;
      }
    }
  }

  .send-action{
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    cursor: pointer;
    font-size: 16px;
    transition: all 0.2s;
  }
}
</style>
