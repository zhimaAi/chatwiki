<template>
  <div>
    <a-modal
      v-model:open="settingModal"
      @cancel="onClose"
      :title="modalTitle"
      :maskClosable="false"
      :ok-text="isLoading ? t('btn_ok_parsing') : t('btn_ok')"
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
import { useI18n } from '@/hooks/web/useI18n'
import DocumentSegmentationModal from '@/views/library/document-segmentation/document-segmentation-model.vue'
import { message } from 'ant-design-vue'

const { t } = useI18n('views.library.library-preview.components.segmentation-setting-modal')

const settingModal = ref(false)
const modalTitle = ref(t('modal_title'))

const emit = defineEmits(['ok', 'enable'])
let parmas = {}

const formState = reactive({
  data_ids: 0,
  file_id: 0
})

const getListStatusMap = () => ({
  0: t('status_not_converted'),
  1: t('status_converted'),
  2: t('status_conversion_exception'),
  3: t('status_converting'),
  4: t('tab_segmentation_failed'),
  10: t('status_segmenting')
})

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
    father_chunk_paragraph_number: +data.father_chunk_paragraph_number,
    content: data.content,
    number: +data.number,
    title: data.title,
    word_total: +data.word_total,
    question: data.question,
    answer: data.answer,
    images: data.images
  }]
  currentData.status = data.status
  currentData.errmsg = data.errmsg
  currentData.status_text = getListStatusMap()[data.status]
}

const confirmLoading = ref(false)

const documentSegmentationModalRef = ref(null)

const handleSaveSegmentation = () => {
  confirmLoading.value = true
  documentSegmentationModalRef.value.handleSaveLibFileSplit()
}

const handleSaveOk = () => {
  message.success(t('msg_save_success'))
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