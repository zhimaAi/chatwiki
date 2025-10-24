<template>
  <a-modal v-model:open="open" title="复制链接" :footer="null" width="600px">
    <ShareForm ref="shareFormRef" :baseUrl="props.baseUrl" />
  </a-modal>
</template>

<script setup>
import { OPEN_BOC_BASE_URL } from '@/constants/index'
import { ref, nextTick } from 'vue'
import ShareForm from './share-form.vue'

const props = defineProps({
  baseUrl: {
    type: String,
    default: OPEN_BOC_BASE_URL + '/doc'
  }
})

const open = ref(false)
const shareFormRef = ref(null)

const show = async (data) => {
  try {
    open.value = true
    nextTick(() => {
      shareFormRef.value.init(data)
    })
  } catch (error) {
    console.log(error)
    // message.error('修改失败')
  }
}

const hide = () => {
  open.value = false
}

defineExpose({
  show,
  hide
})
</script>

<style lang="less" scoped>
.share-form {
  :deep(.share-form-body) {
    padding: 0;
  }
}
</style>
