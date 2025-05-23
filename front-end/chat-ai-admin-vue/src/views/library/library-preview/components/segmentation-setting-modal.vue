<template>
  <div>
    <a-modal
      v-model:open="settingModal"
      @cancel="onClose"
      :title="modalTitle"
      :maskClosable="false"
      :ok-text="isLoading ? '文档解析中,请稍候...' : '确定'"
      :confirmLoading="confirmLoading || isLoading"
      @ok="handleSaveSegmentation"
      :width="1084"
    >
      <div class="main-content-box">
        <DocumentSegmentationModal
          v-if="settingModal"
          @ok="handleSaveOk"
          @finish="handleFinish"
          @loading="handleLoading"
          :paragraphType="'paragraphsSegmented'"
          :pdfState="formState"
          :currentData="currentData"
          ref="documentSegmentationModalRef"
        />
      </div>

      <!-- <template #footer></template> -->
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import DocumentSegmentationModal from '@/views/library/document-segmentation/document-segmentation-model.vue'
import { message } from 'ant-design-vue'

const settingModal = ref(false)
const modalTitle = ref('重新分段')

const emit = defineEmits(['ok', 'enable'])
let parmas = {}

const formState = reactive({
  data_ids: 0,
  file_id: 0
})

const listStatusMap = {
  0: '未转换',
  1: '已转换',
  2: '转换异常',
  3: '转换中',
  4: '分段失败',
  10: '分段中'
}

const currentData = reactive({
  status: '',
  errmsg: '',
  status_text: ''
})

const onClose = () => {
  emit('enable')
}
const isLoading = ref(false)
const handleLoading = (data) => {
  isLoading.value = data
}
const show = (data) => {
  settingModal.value = true
  formState.data_ids = data.id
  formState.file_id = data.file_id
  // 点击当前段落的重新分段，需要填充的数据
  currentData.list = [{
    content: data.content,
    number: data.number,
    title: data.title,
    word_total: data.word_total,
    question: data.question,
    answer: data.answer,
    images: data.images
  }]
  currentData.status = data.status
  currentData.errmsg = data.errmsg
  currentData.status_text = listStatusMap[data.status]
}

const confirmLoading = ref(false)

const documentSegmentationModalRef = ref(null)

const handleSaveSegmentation = () => {
  confirmLoading.value = true
  documentSegmentationModalRef.value.handleSaveLibFileSplit()
}

const handleSaveOk = () => {
  message.success('保存成功')
  settingModal.value = false
  emit('ok', parmas.type)
}
const handleFinish = () => {
  confirmLoading.value = false
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.main-content-box {
  max-height: 650px;
  height: 600px;
  overflow: hidden;
}
</style>