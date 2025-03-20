<style lang="less" scoped></style>

<template>
  <input
    type="file"
    style="opacity: 0; position: fixed; bottom: -9999px; left: -9999px"
    ref="fileInput"
    @change="handleFileUpload"
    accept=".md"
  />
</template>

<script setup>
import { ref } from 'vue'

const emit = defineEmits(['file-loaded'])
const fileInput = ref(null)
let params = {}

// 暴露给父组件的触发方法
const triggerUpload = (data) => {
  params = data || {}
  fileInput.value.click() // 模拟点击隐藏的 input
}

// 文件上传处理
const handleFileUpload = (e) => {
  const file = e.target.files[0]
  if (!file) return

  // 校验文件类型
  if (!file.name.endsWith('.md')) {
    alert('请选择.md文件')
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    emit('file-loaded', {
      content: e.target.result,
      title: file.name,
      file: file,
      ...params
    })
  }

  reader.readAsText(file)

  // 清空选择以便重复上传相同文件
  e.target.value = null
}

defineExpose({ triggerUpload })
</script>
