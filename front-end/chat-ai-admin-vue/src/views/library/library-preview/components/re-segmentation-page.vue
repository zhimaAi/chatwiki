<template>
  <div>
    <a-modal
      v-model:open="open"
      :title="modalTitle"
      @ok="handleSaveType"
      @cancel="onClose"
      :width="720"
      :confirmLoading="confirmLoading"
    >
      <div class="select-box">
        <div
          class="select-list-item"
          :class="{ active: item.pdf_parse_type == formState.pdf_parse_type }"
          v-for="item in analysisLists"
          @click="handleChangeAnalysis(item)"
          :key="item.pdf_parse_type"
        >
          <div class="title">{{ item.title }}</div>
          <div class="desc">{{ item.desc }}</div>
        </div>
      </div>
    </a-modal>
    <a-modal
      v-model:open="settingModal"
      @cancel="onClose"
      :title="modalTitle"
      :maskClosable="false"
      :ok-text="isLoading ? '文档解析中,请稍候...' : '确定'"
      :confirmLoading="confirmLoading || isLoading"
      @ok="handleSaveSegmentation"
      :width="1000"
    >
      <div class="main-content-box">
        <DocumentSegmentationModal
          v-if="settingModal"
          @ok="handleSaveOk"
          @finish="handleFinish"
          @loading="handleLoading"
          :pdfState="formState"
          ref="documentSegmentationModalRef"
        />
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DocumentSegmentationModal from '@/views/library/document-segmentation/document-segmentation-model.vue'
import { message } from 'ant-design-vue'
import { restudyLibraryFile } from '@/api/library'

const query = useRoute().query
const router = useRouter()
const props = defineProps({
  detailsInfo: {
    type: Object,
    default: () => {}
  }
})

const open = ref(false)
const settingModal = ref(false)
const modalTitle = ref('')

const emit = defineEmits(['ok', 'enable'])
let parmas = {}

const analysisLists = [
  {
    title: '图文OCR解析',
    desc: '通过OCR文字识别提取pdf文件内容，可以兼容扫描件，但是解析速度较慢。',
    pdf_parse_type: 2
  },
  {
    title: '纯文本解析',
    desc: '只提取Pdf中的文字内容，如果文档为扫描件可能提取不到内容。',
    pdf_parse_type: 1
  },
  {
    title: '图文解析',
    desc: '提取PDF文档中的图片和文字',
    pdf_parse_type: 3
  }
]

const formState = reactive({
  pdf_parse_type: 2,
  pdf_page_num: 0
})

const onClose = () => {
  emit('enable')
}
const isLoading = ref(false)
const handleLoading = (data) => {
  isLoading.value = data
}
const show = (data) => {
  parmas = {
    ...data
  }
  formState.pdf_parse_type = 2
  formState.pdf_page_num = data.pdf_page_num
  if (data.type == 1) {
    modalTitle.value = `当页(${data.pdf_page_num})重新分段`
  }
  if (data.type == 2) {
    modalTitle.value = `将文档重新分段`
  }
  if (data.type == 3) {
    modalTitle.value = `重新学习文档`
  }
  open.value = true
}

const handleChangeAnalysis = (item) => {
  formState.pdf_parse_type = item.pdf_parse_type
}
const confirmLoading = ref(false)

const handleSaveType = () => {
  if (parmas.type == 3) {
    if (formState.pdf_parse_type == 1) {
      // 跳转到文档分段页
      router.push('/library/document-segmentation?document_id=' + query.id + '&source=preview')
      return
    }
    confirmLoading.value = true
    restudyLibraryFile({
      id: query.id,
      pdf_parse_type: formState.pdf_parse_type
    })
      .then((res) => {
        message.success('重新学习成功')
        if (formState.pdf_parse_type == 2 || formState.pdf_parse_type == 3) {
          //跳转回知识库文档列表，当前文档状态进入转换中
          router.push({
            path: '/library/details/knowledge-document',
            query: {
              id: props.detailsInfo.library_id
            }
          })
        }
        if (formState.pdf_parse_type == 1) {
          // 跳转到文档分段页
          router.push('/library/document-segmentation?document_id=' + query.id + '&source=preview')
        }
        open.value = false
      })
      .finally(() => {
        confirmLoading.value = false
      })
    return
  }
  open.value = false
  settingModal.value = true
}

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
.select-box {
  .select-list-item {
    padding: 16px;
    margin-top: 16px;
    border: 1px solid #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    &.active {
      border: 1px solid #2475fc;
    }
    .title {
      color: #262626;
      font-weight: 600;
      font-size: 16px;
    }
    .desc {
      color: #595959;
    }
  }
}
.main-content-box {
  max-height: 650px;
  height: 600px;
  overflow: hidden;
}
</style>
