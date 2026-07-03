<template>
  <div>
    <a-modal
      v-model:open="open"
      :title="t('import_modal.title')"
      :width="486"
      :destroyOnClose="true"
      :centered="true"
      @cancel="handleCancel"
    >
      <div class="import-alert">
        {{ t('import_modal.alert_text_before') }}
        <a class="template-link" @click="handleDownloadTemplate">{{ t('import_modal.alert_template') }}</a>
        {{ t('import_modal.alert_text_after') }}
      </div>

      <div class="form-item">
        <div class="form-label">{{ t('import_modal.group_label') }}</div>
        <a-select
          v-model:value="selectedGroup"
          class="group-select"
          :placeholder="t('import_modal.group_placeholder')"
          :options="groupOptions"
          allow-clear
        />
      </div>

      <a-upload-dragger
        v-model:fileList="fileList"
        :multiple="false"
        accept=".xlsx,.xls,.csv"
        :before-upload="handleBeforeUpload"
        :show-upload-list="false"
        class="custom-upload-dragger"
      >
        <div class="upload-inner">
          <div class="upload-icon">
            <InboxOutlined />
          </div>
          <div class="upload-text-wrapper">
            <div class="upload-main-text">
              <span class="required-asterisk">*</span>{{ t('import_modal.upload_text') }}
            </div>
            <div class="upload-hint">{{ t('import_modal.upload_hint') }}</div>
            <div class="upload-types">{{ t('import_modal.upload_types') }}</div>
          </div>
        </div>
      </a-upload-dragger>

      <div v-if="selectedFile" class="file-list">
        <div class="file-item">
          <PaperClipOutlined class="file-icon" />
          <span class="file-name">{{ selectedFile.name }}</span>
          <span class="file-size">{{ formatFileSize(selectedFile.size) }}</span>
          <span class="file-divider"></span>
          <a class="file-delete" @click="handleRemoveFile">{{ t('import_modal.delete') }}</a>
        </div>
      </div>

      <template #footer>
        <div class="modal-footer">
          <a-button class="footer-btn cancel-btn" @click="handleCancel">{{ t('import_modal.cancel') }}</a-button>
          <a-button
            type="primary"
            class="footer-btn confirm-btn"
            :disabled="!selectedFile"
            :loading="importing"
            @click="handleConfirm"
          >
            {{ t('import_modal.confirm') }}
          </a-button>
        </div>
      </template>
    </a-modal>

    <a-modal
      v-model:open="resultOpen"
      :title="resultTitle"
      :width="486"
      :destroyOnClose="true"
      :centered="true"
      @cancel="handleResultClose"
    >
      <div class="import-result-body">
        <div class="result-summary">
          <div class="result-summary-item">
            {{ t('import_modal.result_success_count', { count: resultInfo.created_count }) }}
          </div>
          <div class="result-summary-item">
            {{ t('import_modal.result_updated_count', { count: resultInfo.updated_count }) }}
          </div>
          <div class="result-summary-item">
            {{ t('import_modal.result_failed_count', { count: resultInfo.failed_count }) }}
          </div>
        </div>

        <div v-if="hasImportErrors" class="result-fail-section">
          <div class="result-fail-tip">{{ t('import_modal.result_failed_tip') }}</div>
          <div v-if="canDownloadFailedRows" class="result-file-row">
            <PaperClipOutlined class="result-file-icon" />
            <span class="result-file-name">{{ t('import_modal.failed_record') }}</span>
            <a class="result-file-link" @click="handleDownloadFailedRows">{{ t('import_modal.download') }}</a>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="modal-footer">
          <a-button class="footer-btn cancel-btn" @click="handleResultClose">{{ t('import_modal.cancel') }}</a-button>
          <a-button type="primary" class="footer-btn confirm-btn" @click="handleResultClose">
            {{ t('import_modal.confirm') }}
          </a-button>
        </div>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { InboxOutlined, PaperClipOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import { message } from 'ant-design-vue'
import { importGoodsLibrary, downloadImportTemplate } from '@/api/goods-library'
import { useI18n } from '@/hooks/web/useI18n'
import { exportAoAToXlsx } from '@/utils/export/xlsx'

const { t } = useI18n('views.library.goods-library.index')

const emit = defineEmits(['ok'])

const props = defineProps({
  groupId: {
    type: [String, Number],
    default: '0'
  },
  groupOptions: {
    type: Array,
    default: () => []
  }
})

const open = ref(false)
const fileList = ref([])
const importing = ref(false)
const downloadingFailedRows = ref(false)
const selectedGroup = ref(undefined)
const selectedFile = ref(null)
const resultOpen = ref(false)
const resultInfo = ref({
  created_count: 0,
  updated_count: 0,
  failed_count: 0,
  errors: [],
  headers: []
})

const hasImportErrors = computed(() => resultInfo.value.failed_count > 0)
const resultTitle = computed(() =>
  hasImportErrors.value ? t('message.import_failed') : t('message.import_success')
)
const exportHeaders = computed(() => {
  const headers = Array.isArray(resultInfo.value.headers) ? resultInfo.value.headers : []

  return headers
    .filter((item) => item?.Field && item?.Name)
    .map((item) => ({
      field: item.Field,
      name: item.Name
    }))
})
const canDownloadFailedRows = computed(() =>
  hasImportErrors.value &&
  Array.isArray(resultInfo.value.errors) &&
  resultInfo.value.errors.length > 0 &&
  exportHeaders.value.length > 0
)

const show = () => {
  fileList.value = []
  selectedFile.value = null
  const id = String(props.groupId)
  selectedGroup.value = id && id !== 'all' && id !== '-1' ? props.groupId : '0'
  importing.value = false
  resultOpen.value = false
  open.value = true
}

const handleDownloadTemplate = () => {
  downloadImportTemplate()
}

const handleBeforeUpload = (file) => {
  const isValidType = ['.xlsx', '.xls', '.csv'].some((ext) =>
    file.name.toLowerCase().endsWith(ext)
  )
  if (!isValidType) {
    message.error(t('import_modal.upload_failed', { name: file.name }))
    return false
  }

  const maxSize = 100 * 1024 * 1024
  if (file.size > maxSize) {
    message.error(t('import_modal.file_too_large'))
    return false
  }

  selectedFile.value = file
  fileList.value = [file]
  return false
}

const handleRemoveFile = () => {
  selectedFile.value = null
  fileList.value = []
}

const handleCancel = () => {
  open.value = false
}

const normalizeExportValue = (value) => {
  if (Array.isArray(value)) {
    return value.join(',')
  }

  if (value === undefined || value === null) {
    return ''
  }

  return value
}

const showImportResult = (result) => {
  const resultData = result?.data || {}

  resultInfo.value = {
    created_count: Number(resultData.created_count) || 0,
    updated_count: Number(resultData.updated_count) || 0,
    failed_count: Number(resultData.failed_count) || 0,
    errors: Array.isArray(resultData.errors) ? resultData.errors : [],
    headers: Array.isArray(resultData.headers) ? resultData.headers : []
  }

  open.value = false
  resultOpen.value = true
}

const handleResultClose = () => {
  resultOpen.value = false
}

const handleDownloadFailedRows = async () => {
  if (!canDownloadFailedRows.value || downloadingFailedRows.value) {
    return
  }

  downloadingFailedRows.value = true

  try {
    const headers = [
      ...exportHeaders.value,
      {
        field: 'message',
        name: t('import_modal.fail_reason')
      }
    ]
    const rows = [
      headers.map((item) => item.name),
      ...resultInfo.value.errors.map((item) =>
        headers.map(({ field }) => String(normalizeExportValue(item?.[field])))
      )
    ]
    const cols = headers.map(({ field, name }) => ({
      wch: ['images', 'description', 'qa', 'custom_info', 'message'].includes(field)
        ? 28
        : Math.max(String(name || '').length + 4, 14)
    }))

    await exportAoAToXlsx({
      sheetName: t('import_modal.failed_record'),
      fileName: `${t('import_modal.fail_data_name')}_${dayjs().format('YYYY-MM-DD_HH-mm-ss')}.xlsx`,
      rows,
      cols
    })
  } catch (e) {
    message.error(t('import_modal.download_failed_message'))
  } finally {
    downloadingFailedRows.value = false
  }
}

const handleConfirm = async () => {
  if (!selectedFile.value) {
    message.warning(t('import_modal.no_file_selected'))
    return
  }

  importing.value = true

  try {
    const result = await importGoodsLibrary({
      group_id: Number(selectedGroup.value) || 0,
      file: selectedFile.value
    })

    showImportResult(result)
    importing.value = false
    emit('ok')
  } catch (e) {
    importing.value = false
  }
}

const formatFileSize = (bytes) => {
  if (!bytes) return ''
  if (bytes < 1024) return bytes + 'B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + 'KB'
  return (bytes / (1024 * 1024)).toFixed(1) + 'MB'
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.import-alert {
  margin-bottom: 24px;
  background: #e9f1fe;
  border: 1px solid #99bffd;
  border-radius: 6px;
  padding: 9px 16px;
  font-size: 14px;
  line-height: 22px;
  color: #3a4559;

  .template-link {
    color: #2475fc;
    cursor: pointer;

    &:hover {
      text-decoration: underline;
    }
  }
}

.form-item {
  margin-bottom: 24px;
  .form-label {
    font-size: 14px;
    line-height: 22px;
    color: #262626;
    margin-bottom: 4px;
  }

  .group-select {
    width: 100%;
  }
}

.custom-upload-dragger {
  margin-top: 24px;

  .upload-inner {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
  }

  .upload-icon {
    font-size: 48px;
    color: #2475fc;
    line-height: 1;
  }

  .upload-text-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    width: 100%;
  }

  .upload-main-text {
    font-size: 14px;
    line-height: 22px;
    color: rgba(0, 0, 0, 0.85);

    .required-asterisk {
      color: #fb363f;
      margin-right: 2px;
    }
  }

  .upload-hint {
    font-size: 12px;
    line-height: 20px;
    color: rgba(0, 0, 0, 0.45);
  }

  .upload-types {
    font-size: 12px;
    line-height: 20px;
    color: rgba(0, 0, 0, 0.45);
  }
}

.file-list {
  margin-top: 8px;
  padding: 4px 0;

  .file-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 0 4px;
    font-size: 14px;
    line-height: 22px;

    .file-icon {
      font-size: 14px;
      color: #595959;
      flex-shrink: 0;
    }

    .file-name {
      color: #595959;
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .file-size {
      color: #8c8c8c;
      flex-shrink: 0;
    }

    .file-divider {
      width: 1px;
      height: 12px;
      background: #d9d9d9;
      flex-shrink: 0;
    }

    .file-delete {
      color: #2475fc;
      cursor: pointer;
      flex-shrink: 0;

      &:hover {
        text-decoration: underline;
      }
    }
  }
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
}

.import-result-body {
  padding-top: 6px;
}

.result-summary {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.result-summary-item {
  font-size: 16px;
  line-height: 28px;
  color: #262626;
}

.result-fail-section {
  margin-top: 20px;
  padding-top: 18px;
  border-top: 1px dashed #e8e8e8;
}

.result-fail-tip {
  font-size: 14px;
  line-height: 22px;
  color: #262626;
}

.result-file-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
  font-size: 14px;
  line-height: 22px;
}

.result-file-icon {
  color: #8c8c8c;
  font-size: 14px;
}

.result-file-name {
  color: #595959;
}

.result-file-link {
  color: #2475fc;
  cursor: pointer;

  &:hover {
    text-decoration: underline;
  }
}
</style>
