<template>
  <div class="upload-input" :class="{ disabled: fileList.length >= props.maxCount }">
    <a-upload-dragger :file-list="fileList" name="file" :multiple="true" :max-count="props.maxCount" :accept="getAccept"
      :beforeUpload="beforeUpload" @remove="handleRemove" class="ant-upload-drag">
      <div class="ant-upload-drag-icon">
        <inbox-outlined></inbox-outlined>
      </div>
      <div class="ant-upload-text">点击或将文件拖拽到这里上传</div>
      <div class="ant-upload-hint">
        <p>一次只能上传一个文档，单个文件不超过10M</p>
        <p>
          <span>支持文件类型：</span><span class="ant-upload-hint-ext" v-for="ext in fileTypes" :key="ext">.{{ ext }}</span>
        </p>
      </div>
    </a-upload-dragger>

    <div class="disabled-mask" @click="showMaxCountErrorTip"></div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import { InboxOutlined } from '@ant-design/icons-vue'

const emit = defineEmits(['change', 'update:value'])

const props = defineProps({
  value: {
    type: Array,
    default: () => []
  },
  maxCount: {
    type: Number,
    default: 1
  }
})

const fileTypes = ['pdf', 'docx', 'txt', 'md', 'xlsx', 'csv', 'html']
const getAccept = computed(() => {
  return fileTypes.map((item) => `.${item}`).join(',')
})

const fileList = ref([])

watch(
  () => props.value,
  (val) => {
    fileList.value = val
  },
  {
    immediate: true
  }
)

let timer = null

const showMaxCountErrorTip = () => {
  if (fileList.value.length >= props.maxCount) {
    return message.error('一次最多上传' + props.maxCount + '个文件')
  }
}

const onChange = () => {
  if (fileList.value.length > props.maxCount) {
    message.error('一次最多上传' + props.maxCount + '个文件')
  }

  let files = [...fileList.value]

  emit('change', files)

  emit('update:value', files)
}

const handleRemove = (file) => {
  removeFile(file)

  onChange()
}

const removeFile = (file) => {
  for (let i = 0; i < fileList.value.length; i++) {
    if (fileList.value[i].uid === file.uid) {
      fileList.value.splice(i, 1)
      break
    }
  }
}

const beforeUpload = (file) => {
  const ext = file.name.split('.').pop()
  const isWhitelistFile = fileTypes.includes(ext.toLowerCase())

  if (!isWhitelistFile) {
    message.error('不支持的文件格式')
    return false
  }

  const isLt100M = file.size / 1024 / 1024 < 10

  if (!isLt100M) {
    message.error('文件大小不能超过10M')
    return false
  }

  fileList.value = [...(fileList.value || []), file]

  if (fileList.value.length > props.maxCount) {
    fileList.value = fileList.value.splice(0, props.maxCount)
  }

  if (timer) {
    clearTimeout(timer)
  }

  timer = setTimeout(() => {
    onChange()
    clearTimeout(timer)
  }, 50)

  return false
}
</script>

<style lang="less" scoped>
.ant-upload-drag {
  .ant-upload-drag-icon {
    font-size: 48px;
    color: #2475fc;
  }

  .ant-upload-text {
    line-height: 22px;
    font-size: 14px;
    text-align: center;
    color: #000000;
  }

  .ant-upload-hint {
    line-height: 20px;
    margin-top: 4px;
    font-size: 12px;
    color: rgba(0, 0, 0, 0.45);
    text-align: center;
  }
}

.upload-input {
  .ant-upload-hint-ext {
    padding-right: 4px;
  }

  .disabled-mask {
    display: none;
    cursor: pointer;
    height: 176px;
  }

  &.disabled {
    position: relative;

    .disabled-mask {
      display: block;
      position: absolute;
      left: 0;
      right: 0;
      top: 0;
      bottom: 0;
      opacity: 0;
    }
  }
}
</style>
