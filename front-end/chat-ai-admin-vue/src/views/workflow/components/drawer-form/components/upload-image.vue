<template>
  <div>
    <a-upload
      v-model:file-list="fileList"
      name="file"
      accept=".jpeg,.png,.jpg"
      list-type="picture-card"
      class="avatar-uploader"
      :show-upload-list="false"
      :action="uploadUrl"
      :data="{
        upload_type: 'image',
        business_type: 'retain_data_robot',
      }"
      :before-upload="beforeUpload"
      @change="handleChange"
    >
      <img class="robot-avatar" v-if="imageUrl" :src="imageUrl" alt="avatar" />
      <div v-else>
        <loading-outlined v-if="loading"></loading-outlined>
        <plus-outlined v-else></plus-outlined>
        <div class="ant-upload-text">上 传</div>
      </div>
    </a-upload>
  </div>
</template>

<script setup>
import { ref, reactive, h, watch } from 'vue'
import { Form, message, Modal } from 'ant-design-vue'
import { CloseCircleFilled, LoadingOutlined, PlusOutlined } from '@ant-design/icons-vue'
const emit = defineEmits(['update:value'])
const props = defineProps({
  value: {
    type: [String, Array],
    default: '',
  },
})
const fileList = ref([])
const imageUrl = ref('')
watch(
  () => props.value,
  (val) => {
    imageUrl.value = val
  },
  {
    immediate: true,
  },
)
const uploadUrl = ref('/Upload/ManagerPUpload')
if (import.meta.env.DEV) {
  uploadUrl.value = '/api/Upload/ManagerPUpload'
}

const loading = ref(false)
const beforeUpload = (file) => {
  fileList.value = []
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
  if (!isJpgOrPng) {
    message.error('请上传图片格式为JPG/PNG的图片!')
  }
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isLt2M) {
    message.error('图片尺寸不能大于2M!')
  }
  return isJpgOrPng && isLt2M
}

const handleChange = (info) => {
  if (info.file.status === 'uploading') {
    loading.value = true
    return
  }
  if (info.file.status === 'done') {
    loading.value = false
    imageUrl.value = info.fileList[0].response.data.filePath
    emit('update:value', imageUrl.value)
  }
  if (info.file.status === 'error') {
    loading.value = false
    message.error('上传失败')
  }
}
</script>

<style lang="less" scoped>
.robot-avatar {
  width: 100%;
  height: 100%;
}
</style>
