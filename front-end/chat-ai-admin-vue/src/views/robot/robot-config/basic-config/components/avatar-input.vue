<template>
  <div class="avatar-uploader">
    <a-upload v-model:file-list="fileList" list-type="picture-card" :show-upload-list="false"
      :before-upload="beforeUpload" accept=".jpg,.png,.jpeg" @preview="handlePreview">
      <img class="avatar" v-if="imageUrl" :src="imageUrl" alt="avatar" />
      <div v-else>
        <loading-outlined v-if="loading"></loading-outlined>
        <plus-outlined v-else></plus-outlined>
        <div class="ant-upload-text">上传照片</div>
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

function getBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = () => resolve(reader.result)
    reader.onerror = (error) => reject(error)
  })
}

const emit = defineEmits(['update:value', 'change'])
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
    imageUrl.value = val
  },
  {
    immediate: true
  }
)

const formItemContext = Form.useInjectFormItemContext()

const triggerChange = () => {
  let data = {
    imageUrl: imageUrl.value,
    file: fileList.value[0].originFileObj
  }

  emit('change', data)
  emit('update:value', data.imageUrl)

  formItemContext.onFieldChange()
}

const beforeUpload = (file) => {
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'

  if (!isJpgOrPng) {
    message.error('机器人头像只支持JPG、PNG格式的图片')
    return false
  }

  const isLt2M = file.size / 1024 < 100

  if (!isLt2M) {
    message.error('机器人头像图片大小不能超过100kb')
    return false
  }

  fileList.value = [file]

  getBase64(file).then((base64Url) => {
    imageUrl.value = base64Url
    loading.value = false

    triggerChange()
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
  if (!file.url && !file.preview) {
    file.preview = await getBase64(file.originFileObj)
  }
  previewImage.value = file.url || file.preview
  previewVisible.value = true
  previewTitle.value = file.name || file.url.substring(file.url.lastIndexOf('/') + 1)
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

.avatar-uploader .avatar {
  width: 100%;
  height: 100%;
}
</style>
