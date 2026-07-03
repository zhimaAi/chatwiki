<template>
  <a-modal
    v-model:open="open"
    :title="isEdit ? t('title_edit') : t('title_add')"
    :width="472"
    :confirmLoading="saveLoading"
    @ok="handleOk"
    @cancel="handleCancel"
  >
    <div class="mini-card-form">
      <a-alert
        type="info"
        style="margin-bottom: 16px;"
        :message="t('tip_channel_limit')"
      />
      <a-form
        ref="formRef"
        :model="formState"
        :rules="rules"
        layout="vertical"
      >
        <a-form-item name="title" :label="t('label_title')">
          <a-input
            v-model:value="formState.title"
            :placeholder="t('ph_title')"
            :maxLength="50"
          />
        </a-form-item>
        <a-form-item name="appid" :label="t('label_appid')">
          <a-input
            v-model:value="formState.appid"
            :placeholder="t('ph_appid')"
          />
        </a-form-item>
        <a-form-item name="page_path" :label="t('label_path')">
          <a-input
            v-model:value="formState.page_path"
            :placeholder="t('ph_path')"
          />
          <div class="form-tip">{{ t('tip_path') }}</div>
        </a-form-item>
        <a-form-item name="thumb_url" :label="t('label_cover')">
          <a-input v-model:value="formState.thumb_url" style="display:none" />
          <div class="cover-upload-box">
            <a-upload
              :show-upload-list="false"
              :before-upload="handleBeforeUpload"
              :custom-request="handleUpload"
              accept=".png,.jpg,.jpeg"
            >
              <div class="cover-upload-area" v-if="!formState.thumb_url">
                <PlusOutlined style="font-size: 24px; color: rgba(0,0,0,0.45)" />
                <span class="cover-upload-text">{{ t('btn_upload') }}</span>
              </div>
            </a-upload>
            <div class="cover-preview" v-if="formState.thumb_url">
              <img v-viewer :src="formState.thumb_url" alt="" />
              <span class="cover-del-icon" @click="removeCover">
                <CloseOutlined />
              </span>
            </div>
          </div>
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, CloseOutlined } from '@ant-design/icons-vue'
import { uploadFile } from '@/api/app/index'
import { useMiniCardStore } from '@/stores/modules/mini-card'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('components.mini-card.mini-card-form-modal')

const emit = defineEmits(['ok'])

const open = ref(false)
const saveLoading = ref(false)
const isEdit = ref(false)
const formRef = ref()

const formState = reactive({
  id: '',
  title: '',
  appid: '',
  page_path: '',
  thumb_url: ''
})

const rules = {
  title: [{ required: true, message: t('msg_input_title'), trigger: 'blur' }],
  appid: [{ required: true, message: t('msg_input_appid'), trigger: 'blur' }],
  page_path: [{ required: true, message: t('msg_input_path'), trigger: 'blur' }],
  thumb_url: [{ required: true, message: t('msg_upload_cover') }]
}

const miniCardStore = useMiniCardStore()

const show = () => {
  isEdit.value = false
  Object.assign(formState, {
    id: '',
    title: '',
    appid: '',
    page_path: '',
    thumb_url: ''
  })
  open.value = true
}

const edit = (data) => {
  isEdit.value = true
  Object.assign(formState, {
    id: data.id || '',
    title: data.title || '',
    appid: data.appid || '',
    page_path: data.page_path || '',
    thumb_url: data.thumb_url || ''
  })
  open.value = true
}

const handleBeforeUpload = (file) => {
  const isValidType = file.type === 'image/jpeg' || file.type === 'image/png'
  if (!isValidType) {
    message.error(t('msg_image_format'))
    return false
  }
  const isLt2M = file.size / 1024 < 1024 * 2
  if (!isLt2M) {
    message.error(t('msg_image_size'))
    return false
  }
  return true
}

const handleUpload = async ({ file, onError, onSuccess }) => {
  try {
    const res = await uploadFile({
      category: 'received_message_images',
      file
    })
    const url = res?.data?.link || res?.data?.url || ''
    formState.thumb_url = url
    formRef.value?.validateFields?.(['thumb_url'])
    onSuccess && onSuccess(res)
  } catch (e) {
    message.error(t('msg_upload_failed'))
    onError && onError(e)
  }
}

const removeCover = () => {
  formState.thumb_url = ''
  formRef.value?.validateFields?.(['thumb_url'])
}

const handleOk = async () => {
  try {
    await formRef.value.validate()
    saveLoading.value = true
    const data = { ...formState }
    if (isEdit.value) {
      await miniCardStore.editCard(data)
    } else {
      await miniCardStore.addCard(data)
    }
    message.success(t('msg_save_success'))
    open.value = false
    emit('ok')
  } catch (e) {
    // 验证失败或保存失败
  } finally {
    saveLoading.value = false
  }
}

const handleCancel = () => {
  open.value = false
}

defineExpose({ show, edit })
</script>

<style lang="less" scoped>
.mini-card-form {
  .form-tip {
    margin-top: 4px;
    font-size: 14px;
    line-height: 22px;
    color: #8C8C8C;
  }
}

.cover-upload-box {
  .cover-upload-area {
    width: 104px;
    height: 104px;
    border: 1px dashed rgba(0, 0, 0, 0.15);
    border-radius: 6px;
    background: rgba(0, 0, 0, 0.04);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 8px;
    cursor: pointer;
    transition: border-color 0.2s;

    &:hover {
      border-color: #2475FC;
    }
  }

  .cover-upload-text {
    font-size: 14px;
    color: rgba(0, 0, 0, 0.65);
    line-height: 22px;
  }

  .cover-preview {
    position: relative;
    width: 104px;
    height: 104px;
    border-radius: 6px;
    overflow: hidden;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .cover-del-icon {
      position: absolute;
      top: 4px;
      right: 4px;
      width: 20px;
      height: 20px;
      background: rgba(0, 0, 0, 0.5);
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      color: #fff;
      font-size: 10px;
    }
  }
}
</style>
