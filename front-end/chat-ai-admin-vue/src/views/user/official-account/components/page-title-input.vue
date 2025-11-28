<style lang="less" scoped>
.page-title-input {
  display: flex;
  align-items: center;

  .avatar {
    display: block;
    width: 40px;
    height: 40px;
    margin-right: 8px;
  }
}
</style>

<template>
  <div class="page-title-input">
    <div class="upload-logo">
      <cu-upload @change="onChangeAvatar">
        <img class="avatar" :src="props.avatar" />
      </cu-upload>
    </div>

    <a-input size="large" :value="props.value" :placeholder="placeholder" @input="onInput" />
  </div>
</template>

<script setup>
import { getBase64 } from '@/utils/index'
import { uploadFile } from '@/api/app'
import CuUpload from '@/components/cu-upload/cu-upload.vue'
const emit = defineEmits(['update:value', 'update:avatar', 'changeAvatar'])

const props = defineProps({
  avatar: {
    type: String,
    default: ''
  },
  value: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: '请输入标题'
  },
  autoUpload: {
    type: Boolean,
    default: true
  }
})

const onInput = (e) => {
  emit('update:value', e.srcElement.value)
}

const onChangeAvatar = (fileList) => {
  let data = {
    file: fileList[0],
    category: 'icon'
  }

  if (props.autoUpload) {
    uploadFile(data).then((res) => {
      emit('update:avatar', res.data.link)
    })
  }

  getBase64(fileList[0]).then((base64Url) => {
    emit('changeAvatar', {
      file: fileList[0],
      url: base64Url
    })
  })
}
</script>
