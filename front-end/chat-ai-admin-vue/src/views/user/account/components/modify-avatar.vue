<template>
  <div class="upload-box">
    <a-upload
      v-model:file-list="fileList"
      list-type="picture-card"
      :show-upload-list="false"
      :before-upload="beforeUpload"
      accept=".jpg,.png,.jpeg"
    >
      <div ref="uploadBtnRef">修改头像</div>
    </a-upload>
  </div>
</template>

<script setup>
import { ref, reactive, toRaw } from 'vue'
import { Form, message } from 'ant-design-vue'
import { saveProfile } from '@/api/manage/index.js'
const emit = defineEmits(['ok'])
function getBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = () => resolve(reader.result)
    reader.onerror = (error) => reject(error)
  })
}
const fileList = ref([])
const imageUrl = ref('')
const beforeUpload = (file) => {
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'

  if (!isJpgOrPng) {
    message.error('机器人头像只支持JPG、PNG格式的图片')
    return false
  }

  const isLt2M = file.size / 1024 < 1024

  if (!isLt2M) {
    message.error('成员头像图片大小不能超过1M')
    return false
  }

  fileList.value = [file]

  getBase64(file).then((base64Url) => {
    imageUrl.value = base64Url

    triggerChange()
  })

  return false
}
const triggerChange = () => {
  let data = {
    imageUrl: imageUrl.value,
    file: fileList.value[0].originFileObj
  }
  handleOk(data)
}
const modalTitle = ref('修改登录密码')
const formState = reactive({
  id: ''
})
const uploadBtnRef = ref(null)
const setAvatar = (record) => {
  formState.id = record.user_id
  uploadBtnRef.value.click()
}

const handleOk = (data) => {
  saveProfile({
    ...toRaw(formState),
    avatar: data.file
  }).then((res) => {
    message.success(`修改成功`)
    emit('ok', data.imageUrl)
  })
}

defineExpose({
  setAvatar
})
</script>

<style lang="less" scoped>
.upload-box {
  position: absolute;
  left: -9999999px;
}
</style>
