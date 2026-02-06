<template>
  <div class="avatar-uploader">
    <img class="avatar" v-if="imageUrl" :src="imageUrl" alt="avatar" />
    <a-upload
      v-model:file-list="fileList"
      list-type="picture-card"
      :show-upload-list="false"
      :before-upload="beforeUpload"
      accept=".jpg,.png,.jpeg"
      @preview="handlePreview"
    >
      <div>
        <LoadingOutlined v-if="loading" />
        <PlusOutlined v-else />
        <div class="ant-upload-text">{{ t('label_upload_photo') }}</div>
      </div>
    </a-upload>
    <a-modal :open="previewVisible" :title="previewTitle" :footer="null" @cancel="handleCancel">
      <img alt="example" style="width: 100%" :src="previewImage" />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { message, Form } from 'ant-design-vue'
import { LoadingOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { uploadFile } from '@/api/app'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-list.components.http-tool-avatar-input')

const emit = defineEmits(['update:value'])
const props = defineProps({
  value: {
    type: String,
    default: ''
  }
})

const fileList = ref([])
const loading = ref(false)
const imageUrl = ref('')

watch(
  () => props.value,
  (val) => {
    imageUrl.value = val || ''
  },
  { immediate: true }
)

const formItemContext = Form.useInjectFormItemContext()

const triggerChange = (link) => {
  imageUrl.value = link
  emit('update:value', imageUrl.value)
  formItemContext.onFieldChange()
}

const beforeUpload = (file) => {
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
  if (!isJpgOrPng) {
    message.error(t('msg_only_jpg_png'))
    return false
  }
  const isLt2M = file.size / 1024 < 1024 * 2
  if (!isLt2M) {
    message.error(t('msg_image_size_limit'))
    return false
  }
  fileList.value = [file]
  loading.value = true
  uploadFile({ file, category: 'http_tool_avatar' }).then((res) => {
    const link = res?.data?.link || ''
    if (!link) {
      message.error(t('msg_upload_failed'))
      loading.value = false
      return
    }
    loading.value = false
    triggerChange(link)
  }).catch(() => {
    loading.value = false
    message.error(t('msg_upload_failed'))
  })
  return false
}

const previewVisible = ref(false)
const previewImage = ref('')
const previewTitle = ref('')

const handleCancel = () => {
  previewVisible.value = false
  previewTitle.value = ''
}
const handlePreview = async (file) => {
  previewImage.value = file.url || imageUrl.value
  previewVisible.value = true
  previewTitle.value = file.name || (previewImage.value ? previewImage.value.substring(previewImage.value.lastIndexOf('/') + 1) : '')
}
</script>

<style lang="less" scoped>
.avatar-uploader::v-deep(.ant-upload) {
  margin: 0 !important;
}
.ant-upload-select-picture-card i {
  font-size: 32px;
  color: #999;
}
.ant-upload-select-picture-card .ant-upload-text {
  margin-top: 8px;
  color: #666;
}

.avatar-uploader {
  display: flex;
}

.avatar-uploader .avatar {
  margin-right: 5px;
  padding: 9px;
  border: 1px solid #D9D9D9;
  border-radius: 7px;
  width: 102px;
  height: 102px;
}
</style>

