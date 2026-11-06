<template>
  <div
    class="message-input-wrapper"
    :class="{
      'empty-mode': props.emptyMode,
      'has-content': hasContent,
    }"
  >
    <div class="file-toolbar-box" v-if="props.fileList.length > 0">
      <FileToolbar :file-list="props.fileList" @delete="deleteFile" />
    </div>
    <div class="message-input-box">
      <a-textarea
        @click="emit('inputClick')"
        class="message-input"
        :value="props.value"
        :placeholder="placeholderText"
        :auto-size="{ minRows: 2, maxRows: 5 }"
        @change="onChange"
        @keydown="handleKeydown"
      />
    </div>
    <div class="message-footer">
      <div class="message-footer-left">
        <ChatModelSelect
          :key="modelSelectVersion"
          class="chat-model-select"
          :disabled="props.modelLoading"
          :modeId="props.modeId"
          :modeName="props.modeName"
          @change="handleModelChange"
        />
        <div class="quick-entry-list">
          <button class="quick-entry-item" type="button" :disabled="props.modelLoading" @click="emit('openPrompt')">
            <UserOutlined />
            <span>{{ t('label_persona') }}</span>
          </button>
          <button class="quick-entry-item" type="button" @click="emit('openKnowledge')">
            <FileOutlined />
            <span>{{ t('label_knowledge_base') }}</span>
          </button>
        </div>
      </div>
      <div class="message-action">
        <button class="select-file-btn" type="button" :aria-label="chatT('upload_images')" @click="openFileDialog" v-if="props.showUpload">
          <PictureOutlined class="select-file-icon" />
          <span class="file-number" :class="{ big: fileList.length > 9 }" v-if="fileList.length > 0">{{ fileList.length }}</span>
        </button>

        <button
          v-if="props.loading"
          type="button"
          :aria-label="chatT('stop')"
          class="send-msg-btn stop"
          @click="stopMessage"
        >
          <span class="stop-icon"></span>
        </button>
        <button
          v-else
          type="button"
          :aria-label="chatT('send')"
          class="send-msg-btn"
          :class="{ disabled, 'can-send': !disabled }"
          :disabled="disabled"
          @click="sendMessage"
        >
          <ArrowUpOutlined class="paper-airplane" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { UserOutlined, FileOutlined, PictureOutlined, ArrowUpOutlined } from '@ant-design/icons-vue'
import ChatModelSelect from '@/views/clawbot/chat/components/chat-model-select.vue'
import FileToolbar from '@/views/clawbot/chat/components/file-toolbar.vue'
import { useUpload } from '@/hooks/web/useUpload.js'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.clawbot.chat.components.message-input')
const { t: chatT } = useI18n('views.customer-service.chat')

// 只处理受控草稿/上传和入口事件；正式配置保存及导航由客服页负责。
const emit = defineEmits(['send', 'stop', 'update:value', 'update:fileList', 'changeModel', 'openPrompt', 'openKnowledge', 'inputClick'])

const props = defineProps({
  value: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  },
  modelLoading: {
    type: Boolean,
    default: false
  },
  fileList: {
    type: Array,
    default: () => []
  },
  modeId: {
    type: [String, Number],
    default: ''
  },
  modeName: {
    type: String,
    default: ''
  },
  showUpload: {
    type: Boolean,
    default: false
  },
  emptyMode: {
    type: Boolean,
    default: false
  },
})

const fileList = computed(() => props.fileList)
const uploadFiles = ref(props.fileList.map((file) => ({ ...file })))
const modelSelectVersion = ref(0)
let syncingFiles = false

watch(() => props.modelLoading, (loading, previousLoading) => {
  // 共用选择器会先改内部选中值；保存结束后重新按正式props展示，失败时不能留下假切换。
  if (previousLoading && !loading) modelSelectVersion.value++
})

// 上传hook会push/splice数组，使用本地副本，不能直接写props或父页的附件对象。
watch(() => props.fileList, (files) => {
  syncingFiles = true
  uploadFiles.value = files.map((file) => ({ ...file }))
  syncingFiles = false
}, { flush: 'sync' })

watch(uploadFiles, (files) => {
  if (!syncingFiles) emit('update:fileList', files.map((file) => ({ ...file })))
}, { deep: true, flush: 'sync' })

const { openFileDialog } = useUpload({
  limit: 10,
  maxSize: 10,
  category: 'chat_image',
  fileList: uploadFiles,
  multiple: true,
  accept: 'image/bmp,image/jpeg,image/png,image/tiff,image/heic,image/gif,image/webp'
})

const deleteFile = (index) => {
  uploadFiles.value = uploadFiles.value.filter((_, i) => i !== index)
}

const hasContent = computed(() => {
  return props.fileList.length > 0 || props.value.trim().length > 0
})

const placeholderText = computed(() => {
  return props.emptyMode ? t('ph_input_question') : (t('ph_input_content') || t('ph_input_question'))
})

const disabled = computed(() => {
  const filesUnavailable = props.fileList.some((file) => file.status !== 'done' || !file.url)
  return props.loading || props.modelLoading || filesUnavailable || !hasContent.value
})

const handleModelChange = (model) => {
  if (!props.modelLoading) emit('changeModel', model)
}

const onChange = (e) => {
  emit('update:value', e.target.value)
}

const handleKeydown = (event) => {
  // 输入法确认不发送；Shift+Enter沿textarea原生光标位置换行，保留未修剪草稿。
  if (event.isComposing || event.keyCode === 229) return
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    event.stopPropagation()
    sendMessage()
  }
}

const sendMessage = () => {
  if (!disabled.value) emit('send', props.value.trim())
}

const stopMessage = () => {
  emit('stop')
}
</script>

<style lang="less" scoped>
.message-input-wrapper {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  position: relative;
  border: 2px solid #e5e7eb;
  border-radius: 16px;
  overflow: visible;
  padding: 10px 12px;
  background: #fff;
  color: #262626;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  transition: border-color 0.2s ease, box-shadow 0.2s ease;

  .file-toolbar-box{
    margin-bottom: 8px;
  }

  .message-input-box {
    min-height: 48px;
  }

  .message-input {
    width: 100%;
    padding: 0;
    line-height: 24px;
    border: none !important;
    outline: none !important;
    resize: none !important;
    box-shadow: none !important;
    font-size: 16px;
    background: transparent !important;

    ::v-deep(textarea) {
      font-size: 16px;
      line-height: 24px;
      color: #262626;
      background: transparent !important;
      padding: 0 !important;

      &::placeholder {
        color: #8c8c8c;
      }
    }
  }

  .message-footer {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 12px;
  }

  .message-footer-left {
    display: flex;
    align-items: center;
    gap: 4px;
    min-width: 0;
    flex-wrap: wrap;
  }

  .chat-model-select {
    flex-shrink: 0;
  }

  .quick-entry-list {
    display: flex;
    align-items: center;
    gap: 4px;
    min-width: 0;
    color: #3a4559;
  }

  .quick-entry-item {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    height: 32px;
    padding: 0 10px;
    border-radius: 8px;
    background: transparent;
    border: 0;
    color: inherit;
    font-size: 14px;
    line-height: 22px;
    white-space: nowrap;
    cursor: pointer;
    transition: background-color 0.2s ease, color 0.2s ease;

    &:hover:not(:disabled) {
      background: #f2f4f7;
    }

    &:active:not(:disabled) {
      background: var(--app-fill-secondary);
    }

    &:disabled {
      cursor: not-allowed;
      color: var(--app-text-disabled);
    }

    :deep(.anticon) {
      font-size: 14px;
    }
  }

  .message-action {
    display: flex;
    gap: 8px;
    align-items: center;
    justify-content: flex-end;
    margin-left: auto;
  }

  .send-msg-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    padding: 0;
    font-weight: 400;
    border-radius: 50%;
    border: none;
    cursor: pointer;
    transition: all 0.2s;
    color: #fff;
    background: #d8dde5;

    &:hover {
      background: #c8d0db;
    }

    &:disabled {
      cursor: not-allowed;
    }
    .paper-airplane {
      font-size: 18px;
    }

    &.stop {
      color: #fff;
      background: #2475fc;

      &:hover {
        background: #1d68e5;
      }
    }

    .stop-icon {
      width: 12px;
      height: 12px;
      border-radius: var(--app-radius-sm);
      background: currentColor;
      display: inline-block;
    }

    &.can-send {
      color: #fff;
      background: #2475fc;
    }

    &.can-send:hover {
      background: #1d68e5;
    }

    &.disabled:hover {
      background: #d8dde5;
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
    border-radius: var(--app-radius-control);
    border: none;
    background: var(--app-bg-container);
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: var(--app-bg-hover);
    }

    &:active {
      background: var(--app-fill-secondary);
    }

    .select-file-icon {
      font-size: 16px;
      color: var(--app-text-secondary);
    }

    .file-number{
      position: absolute;
      right: -8px;
      top: -8px;
      width: 16px;
      height: 16px;
      border-radius: 50%;
      background: var(--app-error);
      color: var(--app-text-on-primary);
      font-size: 12px;
      font-weight: 400;
      display: flex;
      align-items: center;
      justify-content: center;

      &.big{
        width: auto;
        padding: 0 4px;
        border-radius: var(--app-radius-lg);
      }
    }
  }

  &.empty-mode {
    .message-input-box {
      min-height: 82px;
    }
  }

  &.has-content {
    border-color: #adc8ff;
    box-shadow: 0 4px 8px rgba(38, 81, 140, 0.24);
  }

  &:focus-within {
    border-color: var(--app-primary);
    box-shadow: var(--app-shadow-focus);
  }

  .quick-entry-item,
  .send-msg-btn,
  .select-file-btn {
    &:focus-visible {
      outline: 2px solid var(--app-primary);
      outline-offset: 2px;
    }
  }
}
</style>
