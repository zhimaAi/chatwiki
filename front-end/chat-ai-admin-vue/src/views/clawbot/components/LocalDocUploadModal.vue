<template>
  <a-modal
    :open="props.open"
    :title="t('title_upload_file')"
    :width="680"
    :maskClosable="false"
    :destroyOnClose="false"
    @cancel="handleCancel"
  >
    <div class="local-doc-upload-modal">
      <a-upload-dragger
        class="upload-dragger"
        :file-list="fileList"
        :show-upload-list="false"
        :multiple="false"
        :max-count="1"
        :accept="accept"
        :beforeUpload="handleBeforeUpload"
      >
        <div class="upload-dragger-icon">
          <InboxOutlined />
        </div>
        <div class="upload-dragger-text">{{ t('text_drag_or_upload') }}</div>
        <div class="upload-dragger-hint">
          {{ t('text_file_hint', { formats: supportedFormatsText, size: maxSizeText }) }}
        </div>
      </a-upload-dragger>

      <div v-if="selectedFile" class="selected-file">
        <div class="selected-file-name">{{ selectedFile.name }}</div>
        <button class="selected-file-remove" type="button" @click="handleRemoveFile">
          {{ t('btn_remove_file') }}
        </button>
      </div>

      <div class="form-section">
        <div class="field-label">{{ t('label_description') }}</div>
        <a-textarea
          v-model:value="formState.description"
          class="field-textarea"
          :rows="4"
          :maxlength="500"
          :placeholder="t('ph_description')"
        />
        <div class="field-tip">{{ t('tip_description') }}</div>
      </div>

      <div class="form-section keywords-section">
        <div class="field-label">{{ t('label_keywords') }}</div>
        <a-textarea
          v-model:value="formState.keywordsText"
          class="field-textarea"
          :rows="4"
          :maxlength="1000"
          :placeholder="t('ph_keywords')"
        />
        <div class="field-tip">{{ t('tip_keywords') }}</div>
      </div>
    </div>

    <template #footer>
      <div class="modal-footer">
        <a-button class="footer-btn" @click="handleCancel">{{ t('btn_cancel') }}</a-button>
        <a-button type="primary" class="footer-btn primary-btn" :loading="props.loading" @click="handleConfirm">
          {{ t('btn_confirm_upload') }}
        </a-button>
      </div>
    </template>
  </a-modal>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { InboxOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.clawbot.components.LocalDocUploadModal')

const ALLOWED_EXTS = ['docx', 'doc', 'xlsx', 'xls', 'md', 'txt', 'pdf', 'csv']
const MAX_SIZE = 100 * 1024 * 1024

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:open', 'confirm'])

const fileList = ref([])
const formState = reactive({
  description: '',
  keywordsText: ''
})

const accept = ALLOWED_EXTS.map((ext) => `.${ext}`).join(',')
const supportedFormatsText = ALLOWED_EXTS.map((ext) => `.${ext}`).join(' ')
const maxSizeText = `${Math.floor(MAX_SIZE / 1024 / 1024)}MB`

const selectedFile = computed(() => fileList.value[0] || null)

watch(
  () => props.open,
  (open) => {
    if (!open) {
      resetForm()
    }
  }
)

const resetForm = () => {
  fileList.value = []
  formState.description = ''
  formState.keywordsText = ''
}

const handleBeforeUpload = (file) => {
  const ext = file.name.split('.').pop()?.toLowerCase()
  if (!ALLOWED_EXTS.includes(ext)) {
    message.error(t('msg_unsupported_format', { formats: ALLOWED_EXTS.join('、') }))
    return false
  }

  if (file.size > MAX_SIZE) {
    message.error(t('msg_file_too_large', { size: maxSizeText }))
    return false
  }

  fileList.value = [file]
  return false
}

const handleRemoveFile = () => {
  fileList.value = []
}

const handleCancel = () => {
  emit('update:open', false)
}

const handleConfirm = () => {
  if (!selectedFile.value) {
    message.error(t('msg_select_file'))
    return
  }

  emit('confirm', {
    file: selectedFile.value,
    description: formState.description.trim(),
    keywords: formState.keywordsText
      .split(/\r?\n/)
      .map((item) => item.trim())
      .filter(Boolean)
  })
}
</script>

<style lang="less" scoped>
.local-doc-upload-modal {
  padding-top: 8px;
}

.upload-dragger {
  border-radius: 12px;

  :deep(.ant-upload-drag) {
    border-radius: 12px;
    border: 1px dashed #c8d8f8;
    background: #fbfdff;
    padding: 36px 24px;
  }
}

.upload-dragger-icon {
  color: #2475fc;
  font-size: 40px;
  line-height: 1;
}

.upload-dragger-text {
  margin-top: 14px;
  color: #262626;
  font-size: 16px;
  line-height: 24px;
  font-weight: 500;
}

.upload-dragger-hint {
  margin-top: 8px;
  color: #8c8c8c;
  font-size: 13px;
  line-height: 22px;
}

.selected-file {
  margin-top: 12px;
  padding: 10px 12px;
  border-radius: 10px;
  background: #f5f8ff;
  border: 1px solid #d9e6ff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.selected-file-name {
  color: #262626;
  font-size: 14px;
  line-height: 22px;
  word-break: break-all;
}

.selected-file-remove {
  flex-shrink: 0;
  border: none;
  background: transparent;
  color: #2475fc;
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;
  padding: 0;
}

.form-section {
  margin-top: 20px;
}

.keywords-section {
  margin-top: 16px;
}

.field-label {
  margin-bottom: 10px;
  color: #262626;
  font-size: 14px;
  line-height: 22px;
  font-weight: 600;
}

.field-textarea {
  :deep(.ant-input) {
    border-radius: 10px;
    padding: 12px 14px;
    resize: none;
  }
}

.field-tip {
  margin-top: 8px;
  color: #8c8c8c;
  font-size: 12px;
  line-height: 20px;
}

.modal-footer {
  padding-top: 16px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.footer-btn {
  min-width: 96px;
  height: 38px;
  border-radius: 10px;
}

.primary-btn {
  box-shadow: none;
}
</style>
