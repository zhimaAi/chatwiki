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
          :class="['select-list-item', { active: item.pdf_parse_type == formState.pdf_parse_type }, { disabled: item.pdf_parse_type == 4 && !ali_ocr_switch }]"
          v-for="item in analysisLists"
          @click="handleChangeAnalysis(item)"
          :key="item.pdf_parse_type"
        >
          <div class="title">
            <svg-icon :name="item.icon" style="font-size: 16px"></svg-icon>
            {{ item.title }}
          </div>
          <div class="desc">{{ item.desc }}</div>
          <div class="card-switch" v-if="item.pdf_parse_type == 4 && !ali_ocr_switch">
            {{ t('ali_ocr_not_enabled') }}
            <div class="card-switch-btn" @click.stop="onGoSwitch">{{ t('go_to_enable') }}</div>
          </div>
        </div>
      </div>
    </a-modal>
    <a-modal
      v-model:open="settingModal"
      @cancel="onClose"
      :title="modalTitle"
      :maskClosable="false"
      :ok-text="isLoading ? t('document_parsing_please_wait') : t('confirm')"
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
import { ref, reactive, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DocumentSegmentationModal from '@/views/library/document-segmentation/document-segmentation-model.vue'
import { message } from 'ant-design-vue'
import { restudyLibraryFile } from '@/api/library'
import { useCompanyStore } from '@/stores/modules/company'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-preview.components.re-segmentation-page')

const companyStore = useCompanyStore()
const ali_ocr_switch = computed(() => {
  return companyStore.companyInfo?.ali_ocr_switch == '1'
})
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

const analysisLists = ref([
  {
    title: t('image_text_ocr'),
    desc: t('image_text_ocr_desc'),
    icon: 'pdf-ocr',
    pdf_parse_type: 2
  },
  {
    title: t('text_only'),
    desc: t('text_only_desc'),
    icon: 'pdf-text',
    pdf_parse_type: 1
  },
  {
    title: t('image_text'),
    desc: t('image_text_desc'),
    icon: 'pdf-img',
    pdf_parse_type: 3
  }
])

const formState = reactive({
  pdf_parse_type: 2,
  pdf_page_num: 0
})

const onClose = () => {
  // analysisLists.value的数据还原，不要阿里云OCR解析了
  analysisLists.value = analysisLists.value.filter(item => item.pdf_parse_type != 4)

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
    modalTitle.value = t('resegment_current_page', { page: data.pdf_page_num })
  }
  if (data.type == 2) {
    modalTitle.value = t('resegment_document')
  }
  if (data.type == 3) {
    modalTitle.value = t('restudy_document')
    analysisLists.value.push({
      title: t('ali_cloud_ocr'),
      desc: t('ali_cloud_ocr_desc'),
      icon: 'ali-ocr',
      pdf_parse_type: 4
    })
  }
  open.value = true
}

const onGoSwitch = () => {
  window.open('#/user/aliocr')
}

const handleChangeAnalysis = (item) => {
  if (item.pdf_parse_type == 4 && !ali_ocr_switch.value) {
    return false
  }
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
        message.success(t('restudy_success'))
        if (formState.pdf_parse_type == 2 || formState.pdf_parse_type == 3 || formState.pdf_parse_type == 4) {
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
  message.success(t('save_success'))
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
      .title {
        color: #2475fc;
      }
    }
    .title {
      color: #262626;
      font-weight: 600;
      font-size: 16px;
    }
    .desc {
      color: #595959;
    }

    .card-switch {
      display: flex;
      gap: 4px;
      color: #595959;

      .card-switch-btn {
        cursor: pointer;
        color: #2475fc;

        &:hover {
          opacity: 0.8;
        }
      }
    }

    &.disabled {
      cursor: no-drop;

      .title,.desc,.card-switch {
        color: #999;
      }
    }
  }
}
.main-content-box {
  max-height: 650px;
  height: 600px;
  overflow: hidden;
}
</style>
