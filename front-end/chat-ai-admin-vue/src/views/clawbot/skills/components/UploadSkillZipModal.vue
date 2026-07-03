<template>
  <a-modal
    class="upload-skill-zip-modal"
    :open="visible"
    :title="modalTitle"
    :width="472"
    :maskClosable="false"
    :destroyOnClose="false"
    :mask="false"
    :z-index="2222"
    @cancel="handleCancel"
  >
    <div class="upload-skill-wrap">
      <a-spin v-if="detailLoading" class="detail-loading" />

      <a-upload-dragger
        v-if="!isEdit"
        class="skill-upload-dragger"
        :file-list="fileList"
        :show-upload-list="false"
        :multiple="false"
        :max-count="1"
        accept=".zip"
        :disabled="uploading || submitLoading"
        :beforeUpload="handleBeforeUpload"
      >
        <div class="upload-icon">
          <InboxOutlined />
        </div>
        <div class="upload-title">{{ t('upload_title') }}</div>
        <div class="upload-hint">{{ t('upload_hint') }}</div>
      </a-upload-dragger>

      <div v-if="!isEdit && selectedFile" class="file-row">
        <div class="file-main">
          <a-progress
            v-if="uploading"
            type="circle"
            :percent="uploadPercent"
            :width="16"
            :show-info="false"
            :stroke-width="12"
          />
          <InboxOutlined v-else class="file-icon" />
          <div class="file-info">
            <div class="file-name">{{ selectedFile.name }}</div>
            <a-progress
              v-if="uploading"
              class="file-progress"
              :percent="uploadPercent"
              :show-info="false"
              size="small"
            />
          </div>
        </div>
        <DeleteOutlined class="delete-icon" @click="handleRemoveFile" />
      </div>

      <div v-if="isEdit" class="edit-file-tip">{{ t('edit_file_tip') }}</div>

      <a-form v-if="!detailLoading" class="skill-form" layout="vertical">
        <a-form-item required :label="t('label_skill_name')">
          <a-input
            v-model:value="formState.skill_name"
            :maxlength="50"
            :placeholder="t('ph_skill_name')"
          />
          <div class="field-tip">{{ t('tip_skill_name') }}</div>
        </a-form-item>

        <a-form-item required :label="t('label_remark_name')">
          <a-input
            v-model:value="formState.remark_name"
            :maxlength="50"
            :placeholder="t('ph_remark_name')"
          />
        </a-form-item>

        <a-form-item required :label="t('label_intro')">
          <a-textarea
            v-model:value="formState.intro"
            :maxlength="500"
            :rows="3"
            :placeholder="t('ph_intro')"
          />
        </a-form-item>
      </a-form>
    </div>

    <template #footer>
      <a-button @click="handleCancel">{{ t('btn_cancel') }}</a-button>
      <a-button type="primary" :loading="submitLoading" :disabled="uploading" @click="handleConfirm">
        {{ t('btn_confirm') }}
      </a-button>
    </template>
  </a-modal>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { DeleteOutlined, InboxOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getClawbotSkillInfo, saveClawbotSkill, uploadClawbotSkillZip } from '@/api/clawbot'

const { t } = useI18n('views.clawbot.skills.components.UploadSkillZipModal')

const MAX_SIZE = 10 * 1024 * 1024
const SKILL_NAME_REG = /^[A-Za-z0-9_-]{1,50}$/

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  robotId: {
    type: [String, Number],
    default: ''
  },
  skillId: {
    type: [String, Number],
    default: 0
  }
})

const emit = defineEmits(['update:visible', 'confirm'])

const fileList = ref([])
const uploadResult = ref(null)
const uploading = ref(false)
const submitLoading = ref(false)
const detailLoading = ref(false)
const uploadPercent = ref(0)

const formState = reactive({
  skill_name: '',
  remark_name: '',
  intro: ''
})

const selectedFile = computed(() => fileList.value[0] || null)
const isEdit = computed(() => Number(props.skillId) > 0)
const modalTitle = computed(() => isEdit.value ? t('title_edit') : t('title'))

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      resetForm()
      if (isEdit.value) {
        loadSkillInfo()
      }
    } else {
      resetForm()
    }
  }
)

const resetForm = () => {
  fileList.value = []
  uploadResult.value = null
  uploading.value = false
  submitLoading.value = false
  detailLoading.value = false
  uploadPercent.value = 0
  formState.skill_name = ''
  formState.remark_name = ''
  formState.intro = ''
}

const loadSkillInfo = async () => {
  if (!props.robotId || !props.skillId) {
    return
  }

  detailLoading.value = true
  try {
    const res = await getClawbotSkillInfo({
      id: props.robotId,
      skill_id: props.skillId
    })
    const data = res?.data || {}
    formState.skill_name = data.skill_name || ''
    formState.remark_name = data.remark_name || data.skill_name || ''
    formState.intro = data.intro || data.description || ''
  } catch (err) {
    console.error('获取 Skill 详情失败', err)
  } finally {
    detailLoading.value = false
  }
}

const validateFile = (file) => {
  const ext = file.name.split('.').pop()?.toLowerCase()
  if (ext !== 'zip') {
    message.error(t('msg_zip_only'))
    return false
  }

  if (file.size > MAX_SIZE) {
    message.error(t('msg_file_too_large'))
    return false
  }

  return true
}

const handleBeforeUpload = (file) => {
  if (!validateFile(file)) {
    return false
  }

  fileList.value = [file]
  uploadResult.value = null
  uploadPercent.value = 0
  uploadSkillZip(file)
  return false
}

const uploadSkillZip = async (file) => {
  if (!props.robotId) {
    message.error(t('msg_missing_robot'))
    return
  }

  uploading.value = true
  try {
    const res = await uploadClawbotSkillZip(
      {
        id: props.robotId,
        file
      },
      (event) => {
        if (!event.total) {
          return
        }
        uploadPercent.value = Math.min(99, Math.round((event.loaded / event.total) * 100))
      }
    )

    uploadResult.value = res?.data || null
    uploadPercent.value = 100
    if (uploadResult.value) {
      formState.skill_name = uploadResult.value.skill_name || formState.skill_name
      formState.remark_name = uploadResult.value.skill_name || formState.remark_name
      formState.intro = uploadResult.value.description || formState.intro
    }
  } catch (err) {
    console.error('上传 skill zip 失败', err)
    uploadResult.value = null
    uploadPercent.value = 0
  } finally {
    uploading.value = false
  }
}

const handleRemoveFile = () => {
  fileList.value = []
  uploadResult.value = null
  uploadPercent.value = 0
}

const validateForm = () => {
  if (!props.robotId) {
    message.error(t('msg_missing_robot'))
    return false
  }

  if (!isEdit.value && !selectedFile.value) {
    message.error(t('msg_select_file'))
    return false
  }

  if (!isEdit.value && !uploadResult.value?.upload_key) {
    message.error(t('msg_wait_upload'))
    return false
  }

  if (!SKILL_NAME_REG.test(formState.skill_name.trim())) {
    message.error(t('msg_skill_name_invalid'))
    return false
  }

  if (!formState.remark_name.trim()) {
    message.error(t('msg_remark_required'))
    return false
  }

  if (formState.remark_name.trim().length > 50) {
    message.error(t('msg_remark_too_long'))
    return false
  }

  if (formState.intro.trim().length > 500) {
    message.error(t('msg_intro_too_long'))
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
    const res = await saveClawbotSkill({
      id: props.robotId,
      skill_id: isEdit.value ? props.skillId : 0,
      skill_name: formState.skill_name.trim(),
      remark_name: formState.remark_name.trim(),
      intro: formState.intro.trim(),
      ...(!isEdit.value ? { upload_key: uploadResult.value.upload_key } : {})
    })

    if (res && res.res === 0) {
      message.success(t('msg_save_success'))
      emit('confirm', res.data)
      emit('update:visible', false)
    } else {
      message.error(res?.msg || t('msg_save_failed'))
    }
  } catch (err) {
    console.error('保存 skill 失败', err)
  } finally {
    submitLoading.value = false
  }
}

const handleCancel = () => {
  if (uploading.value || submitLoading.value) {
    return
  }
  emit('update:visible', false)
}
</script>

<style lang="less" scoped>
.upload-skill-wrap {
  padding-top: 10px;
}

.skill-upload-dragger {
  display: block;

  :deep(.ant-upload-drag) {
    border: 1px dashed #d6dce6;
    border-radius: 6px;
    background: #f7f9fc;
    padding: 20px 16px 18px;
  }
}

.detail-loading {
  display: flex;
  justify-content: center;
  padding: 48px 0;
}

.edit-file-tip {
  padding: 10px 12px;
  border-radius: 6px;
  background: #f7f9fc;
  color: #8c8c8c;
  font-size: 14px;
  line-height: 22px;
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

.file-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 50px;
  margin-top: 8px;
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

.file-icon {
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

.skill-form {
  margin-top: 16px;

  :deep(.ant-form-item) {
    margin-bottom: 18px;
  }

  :deep(.ant-form-item-label) {
    padding-bottom: 4px;
  }

  :deep(.ant-form-item-label > label) {
    color: #262626;
    font-size: 14px;
    line-height: 22px;
  }

  :deep(.ant-input) {
    border-radius: 6px;
  }

  :deep(textarea.ant-input) {
    resize: vertical;
  }
}

.field-tip {
  margin-top: 4px;
  color: #8c8c8c;
  font-size: 14px;
  line-height: 22px;
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
  margin-top: 0;
}
</style>
