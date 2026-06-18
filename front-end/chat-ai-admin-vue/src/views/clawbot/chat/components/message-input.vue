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
          class="chat-model-select"
          :disabled="props.modelLoading"
          :modeId="props.modeId"
          :modeName="props.modeName"
          @change="handleModelChange"
        />
        <div class="quick-entry-list">
          <div class="quick-entry-item" @click="handleOpenPrompt">
            <UserOutlined />
            <span>{{ t('label_persona') }}</span>
          </div>
          <div class="quick-entry-item" @click="handleOpenSkill">
            <StarOutlined />
            <span>{{ t('label_skill') }}</span>
          </div>
          <div class="quick-entry-item" @click="handleOpenKnowledge">
            <FileOutlined />
            <span>{{ t('label_knowledge_base') }}</span>
          </div>
        </div>
      </div>
      <div class="message-action">
        <button class="select-file-btn" type="button" @click="openFileDialog" v-if="props.showUpload">
          <PictureOutlined class="select-file-icon" />
          <span class="file-number" :class="{ big: fileList.length > 9 }" v-if="fileList.length > 0">{{ fileList.length }}</span>
        </button>

        <button
          class="send-msg-btn"
          :class="{ loading: props.loading, disabled, 'can-send': !disabled }"
          :disabled="disabled"
          @click="sendMessage"
        >
          <a-spin size="small" class="loading-action" style="margin-right: 4px" v-if="props.loading" />
          <ArrowUpOutlined class="paper-airplane" v-else />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, toRefs } from 'vue'
import { UserOutlined, StarOutlined, FileOutlined, PictureOutlined, ArrowUpOutlined } from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import ChatModelSelect from './chat-model-select.vue'
import FileToolbar from './file-toolbar.vue'
import { useUpload } from '@/hooks/web/useUpload.js'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.clawbot.chat.components.message-input')
const router = useRouter()

const emit = defineEmits(['send', 'update:value', 'update:fileList', 'changeModel', 'openPrompt', 'openSkill'])

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

const hasContent = computed(() => {
  return props.fileList.length > 0 || props.value.trim().length > 0
})

const placeholderText = computed(() => {
  return props.emptyMode ? t('ph_input_question') : (t('ph_input_content') || t('ph_input_question'))
})

const disabled = computed(() => {
  if (props.modelLoading) {
    return true
  }

  if (props.fileList.length > 0) {
    return false
  }

  return props.loading || props.value.trim().length === 0
})

const handleModelChange = (model) => {
  emit('changeModel', model)
}

const handleOpenPrompt = () => {
  emit('openPrompt')
}

const handleOpenSkill = () => {
  emit('openSkill')
}

const handleOpenKnowledge = () => {
  router.push('/clawbot/knowledge')
}

const onChange = (e) => {
  emit('update:value', e.target.value.trim())
}

const handleKeydown = (event) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    if (!event.target.value.trim()) {
      return
    }
    event.preventDefault()
    event.stopPropagation()
    sendMessage()
  } else if (event.key === 'Enter' && event.shiftKey) {
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
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  position: relative;
  border: 2px solid #e5e7eb;
  border-radius: 16px;
  overflow: visible;
  padding: 10px 12px;
  background: #fff;
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
    font-size: 14px;
    line-height: 22px;
    white-space: nowrap;
    cursor: pointer;
    transition: background-color 0.2s ease, color 0.2s ease;

    &:hover {
      background: #f2f4f7;
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

    &.loading {
      background-color: #2475fc;
    }

    &.can-send {
      background: #2475fc;
    }

    &.can-send:hover {
      background: #1d68e5;
    }

    &.disabled:hover {
      background: #d8dde5;
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
    border-radius: 6px;
    border: none;
    background: #fff;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #f5f7fa;
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

  &.empty-mode {
    .message-input-box {
      min-height: 82px;
    }
  }

  &.has-content {
    border-color: #adc8ff;
    box-shadow: 0 4px 8px rgba(38, 81, 140, 0.24);
  }
}
</style>
