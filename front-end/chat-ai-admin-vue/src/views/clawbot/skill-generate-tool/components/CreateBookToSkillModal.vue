<template>
  <a-modal
    class="create-book-skill-modal"
    :open="visible"
    :title="t('modal_upload_document')"
    :width="472"
    :maskClosable="false"
    :destroyOnClose="false"
    @cancel="handleCancel"
  >
    <a-form class="create-form" layout="vertical">
      <a-form-item required :label="t('label_skill_name')">
        <a-input
          v-model:value="formState.skill_name"
          :maxlength="20"
          :placeholder="t('placeholder_skill_name')"
          :disabled="submitLoading"
        />
      </a-form-item>

      <a-form-item required :label="t('label_generate_model')">
        <ModelSelect
          modelType="LLM"
          v-model:modeName="formState.use_model"
          v-model:modeId="formState.model_config_id"
          :placeholder="t('placeholder_select_model')"
        />
      </a-form-item>
    </a-form>

    <a-upload-dragger
      class="book-upload-dragger"
      :file-list="fileList"
      :show-upload-list="false"
      :multiple="true"
      :max-count="20"
      accept=".txt,.docx,.md"
      :disabled="submitLoading"
      :beforeUpload="handleBeforeUpload"
    >
      <div class="upload-icon">
        <InboxOutlined />
      </div>
      <div class="upload-title">{{ t('upload_title') }}</div>
      <div class="upload-hint">
        {{ t('upload_limit_hint') }}<br />
        {{ t('upload_type_hint') }}
      </div>
    </a-upload-dragger>

    <div v-if="fileList.length" class="file-list">
      <div v-for="file in fileList" :key="file.uid || file.name" class="file-row">
        <div class="file-main">
          <LoadingOutlined v-if="submitLoading" class="file-loading" />
          <FileTextOutlined v-else class="file-icon" />
          <div class="file-info">
            <div class="file-name">{{ file.name }}</div>
            <a-progress
              v-if="submitLoading"
              class="file-progress"
              :percent="80"
              :show-info="false"
              size="small"
            />
          </div>
        </div>
        <DeleteOutlined class="delete-icon" @click="handleRemoveFile(file)" />
      </div>
    </div>

    <template #footer>
      <a-button :disabled="submitLoading" @click="handleCancel">{{ t('btn_cancel') }}</a-button>
      <a-button type="primary" :loading="submitLoading" @click="handleConfirm">{{ t('btn_confirm') }}</a-button>
    </template>
  </a-modal>
</template>

<script setup>
import { reactive, ref, watch } from 'vue'
import {
  DeleteOutlined,
  FileTextOutlined,
  InboxOutlined,
  LoadingOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import ModelSelect from '@/components/model-select/model-select.vue'
import { createBookToSkillTask } from '@/api/clawbot'

const MAX_FILE_COUNT = 20
const MAX_FILE_SIZE = 100 * 1024 * 1024
const ACCEPT_EXTS = ['txt', 'docx', 'md']
const { t } = useI18n('views.clawbot.skill-generate-tool.index')

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  robotId: {
    type: [String, Number],
    default: ''
  }
})

const emit = defineEmits(['update:visible', 'confirm'])

const fileList = ref([])
const submitLoading = ref(false)
const formState = reactive({
  skill_name: '',
  model_config_id: '',
  use_model: ''
})

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      resetForm()
    }
  }
)

const resetForm = () => {
  formState.skill_name = ''
  formState.model_config_id = ''
  formState.use_model = ''
  fileList.value = []
  submitLoading.value = false
}

const getFileExt = (fileName = '') => {
  return fileName.split('.').pop()?.toLowerCase() || ''
}

const validateFile = (file) => {
  const ext = getFileExt(file.name)
  if (!ACCEPT_EXTS.includes(ext)) {
    message.error(t('msg_unsupported_book_format'))
    return false
  }

  if (file.size > MAX_FILE_SIZE) {
    message.error(t('msg_file_too_large'))
    return false
  }

  if (fileList.value.length >= MAX_FILE_COUNT) {
    message.error(t('msg_max_file_count'))
    return false
  }

  return true
}

const handleBeforeUpload = (file) => {
  if (!validateFile(file)) {
    return false
  }
  fileList.value = [...fileList.value, file]
  return false
}

const handleRemoveFile = (file) => {
  if (submitLoading.value) {
    return
  }
  fileList.value = fileList.value.filter((item) => item.uid !== file.uid)
}

const validateForm = () => {
  if (!props.robotId) {
    message.error(t('msg_missing_robot_id'))
    return false
  }
  if (!formState.skill_name.trim()) {
    message.error(t('msg_enter_skill_name'))
    return false
  }
  if (formState.skill_name.trim().length > 20) {
    message.error(t('msg_skill_name_too_long'))
    return false
  }
  if (!formState.model_config_id || !formState.use_model) {
    message.error(t('msg_select_model'))
    return false
  }
  if (!fileList.value.length) {
    message.error(t('msg_upload_document'))
    return false
  }
  return true
}

const handleConfirm = async () => {
  if (!validateForm()) {
    return
  }

  submitLoading.value = true
  try {
    const formData = new FormData()
    formData.append('robot_id', props.robotId)
    formData.append('skill_name', formState.skill_name.trim())
    formData.append('model_config_id', formState.model_config_id)
    formData.append('use_model', formState.use_model)
    fileList.value.forEach((file) => {
      formData.append('files', file)
    })

    const res = await createBookToSkillTask(formData)
    if (res && (res.res === 0 || res.code === 0)) {
      message.success(t('msg_task_created'))
      emit('confirm', res.data)
      emit('update:visible', false)
    } else {
      message.error(res?.msg || t('msg_task_create_failed'))
    }
  } catch (err) {
    console.error('创建Book转Skill任务失败', err)
  } finally {
    submitLoading.value = false
  }
}

const handleCancel = () => {
  if (submitLoading.value) {
    return
  }
  emit('update:visible', false)
}
</script>

<style lang="less" scoped>
.create-form {
  padding-top: 10px;

  :deep(.ant-form-item) {
    margin-bottom: 16px;
  }

  :deep(.ant-form-item-label) {
    padding-bottom: 6px;
  }

  :deep(.ant-form-item-label > label) {
    color: #262626;
    font-size: 14px;
    line-height: 22px;
  }

  :deep(.ant-input),
  :deep(.ant-select-selector) {
    border-radius: 6px;
  }
}

.book-upload-dragger {
  display: block;

  :deep(.ant-upload-drag) {
    border: 1px dashed #d6dce6;
    border-radius: 6px;
    background: #f7f9fc;
    padding: 20px 16px 18px;
  }
}

.upload-icon {
  color: #2475fc;
  font-size: 38px;
  line-height: 1;
}

.upload-title {
  margin-top: 16px;
  color: #262626;
  font-size: 14px;
  line-height: 22px;
}

.upload-hint {
  margin-top: 4px;
  color: #8c8c8c;
  font-size: 14px;
  line-height: 22px;
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
}

.file-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 50px;
  padding: 8px 12px;
  border: 1px dashed #d6dce6;
  border-radius: 6px;
  background: #fff;
}

.file-main {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.file-icon,
.file-loading {
  color: #2475fc;
  font-size: 16px;
  flex-shrink: 0;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  color: #262626;
  font-size: 14px;
  line-height: 22px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-progress {
  margin-top: 3px;
  line-height: 1;
}

.delete-icon {
  color: #595959;
  cursor: pointer;
  font-size: 16px;
  flex-shrink: 0;

  &:hover {
    color: #ff4d4f;
  }
}

:deep(.ant-modal-content) {
  border-radius: 14px;
}

:deep(.ant-modal-header) {
  margin-bottom: 12px;
}

:deep(.ant-modal-title) {
  color: #262626;
  font-size: 16px;
  font-weight: 600;
  line-height: 24px;
}

:deep(.ant-modal-footer) {
  margin-top: 24px;
}
</style>
