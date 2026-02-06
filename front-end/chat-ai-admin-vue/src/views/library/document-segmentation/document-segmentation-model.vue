<style lang="less" scoped>
.breadcrumb-block {
  height: 48px;
  display: flex;
  align-items: center;
}
.page-title {
  height: 38px;
  border-radius: 2px;
  display: flex;
  gap: 8px;
  align-items: center;
  color: #262626;
  margin-bottom: 16px;
  font-weight: 600;
  font-size: 14px;
  line-height: 22px;
  .title {
    font-weight: 600;
    font-size: 14px;
    line-height: 22px;
  }
}
.document-segmentation {
  position: relative;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-flow: column nowrap;
  .page-container {
    display: flex;
    flex: 1;
    overflow: hidden;

    .page-left {
      width: 400px;
      height: 100%;
      overflow-y: auto;
      &::-webkit-scrollbar {
        display: none;
      }
    }
    .page-right {
      flex: 1;
      height: 100%;
      padding-left: 24px;
      overflow: hidden;
    }
  }
}
.document-fragment-preview {
  display: flex;
  flex-flow: column nowrap;
  height: 100%;
  overflow: hidden;

  padding: 14px 16px;
  background-color: #f2f4f7;

  .preview-header {
    display: flex;
    height: 22px;
    line-height: 22px;
    font-size: 14px;

    .label-text {
      font-weight: 600;
      color: #242933;
    }
    .fragment-number {
      padding-left: 8px;
      color: #7a8699;
    }
  }

  .preview-box {
    flex: 1;
    overflow-y: auto;
    .fragment-item {
      margin-top: 8px;
    }
  }
}
.footer-btn-box {
  height: 56px;
  position: fixed;
  left: 16px;
  right: 16px;
  bottom: 16px;
  background: #fff;
  border-top: 1px solid #f0f0f0;
  box-shadow: 0 -8px 4px 0 #00000014;
  border-radius: 0 0 8px 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.loading-wrap {
  position: absolute;
  z-index: 99;
  top: 100px;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(255, 255, 255, 0.3);
}

.loading-box {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: #595959;
  margin-top: 200px;
}
</style>

<template>
  <!-- loading -->
  <div class="loading-wrap" v-if="spinning">
    <a-spin></a-spin>
  </div>

  <div class="document-segmentation">
    <div class="page-container">
      <div class="page-left">
        <SegmentationSetting
          :excellQaLists="excellQaLists"
          :libFileInfo="libFileInfo"
          :library_type="library_type"
          :mode="settingMode"
          :hideSave="true"
          :status="'paragraphsSegmented'"
          ref="segmentationSettingRef"
          @change="onChangeSetting"
          @changeChunkType="onChangeChunkType"
          @save="handleSaveLibFileSplit"
          @validate="onValidate"
        />
      </div>
      <div class="page-right">
        <div class="document-fragment-preview">
          <div class="preview-header">
            <span class="label-text">{{ t('segmentationPreview') }}</span>
            <span class="fragment-number">{{ t('totalSegments', { count: documentFragmentTotal }) }}</span>
          </div>
          <Empty v-if="isEmpty && !aiLoading"></Empty>
          <div v-if="aiLoading" class="loading-box"><a-spin />{{ t('processingData') }}</div>
          <div class="preview-box" ref="previewBoxRef">
            <div
              class="fragment-item"
              v-for="(item, index) in documentFragmentList"
              :key="index"
            >
              <DocumentFragment
                :status="'paragraphsSegmented'"
                :chunk_type="formData.chunk_type"
                :father_chunk_paragraph_number="item.father_chunk_paragraph_number"
                :currentData="currentData"
                :number="item.number"
                :title="item.title"
                :content="item.content"
                :total="item.word_total"
                :question="item.question"
                :answer="item.answer"
                :images="item.images"
                @delete="handleDeleteFragment(index)"
                @edit="handleEditFragment(item, index)"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    <EditFragmentAlert ref="editFragmentAlertRef" @ok="saveFragment" />
  </div>
</template>

<script setup>
import { ref, createVNode, computed, onUnmounted, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import { ExclamationCircleOutlined, LeftOutlined } from '@ant-design/icons-vue'
import Empty from './components/empty.vue'
import SegmentationSetting from './components/segmentation-setting.vue'
import DocumentFragment from './components/document-fragment.vue'
import EditFragmentAlert from './components/edit-fragment-alert.vue'
import {
  getLibFileSplit,
  getSplitParagraph,
  getLibFileSplitAiChunks,
  getLibFileInfo,
  saveSplitParagraph,
  getLibFileExcelTitle
} from '@/api/library/index'
import { useLibraryStore } from '@/stores/modules/library'
import { useI18n } from '@/hooks/web/useI18n'

const libraryStore = useLibraryStore()
const { t } = useI18n('views.library.document-segmentation.document-segmentation-model')

const { setInitDocumentFragmentList } = libraryStore
const route = useRoute()

const emit = defineEmits(['ok', 'finish', 'loading'])

const { id } = route.query
const spinning = ref(true)
const settingMode = ref(1) // 1 表格，0 非表格
let itWasEdited = false

const props = defineProps({
  pdfState: {
    type: Object,
    default: () => {}
  },
  currentData: {
    type: Object,
    default: () => {}
  },
  paragraphType: {
    type: String,
    default: ''
  }
})
const current_chunk_type = ref(1)
const defaultAiChunkPrumpt = 'defaultAiChunkPrumpt'
let formData = {
  id: id,
  separators_no: '', // 自定义分段-分隔符序号集
  chunk_size: 512, // 自定义分段-分段最大长度 默认512，最大值不得超过2000
  not_merged_text: false,
  chunk_overlap: 50, // 自定义分段-分段重叠长度 默认为50，最小不得低于10，最大不得超过最大分段长度的50%
  is_qa_doc: 0, // 0 普通文档 1 QA文档
  question_lable: '', // QA文档-问题开始标识符
  answer_lable: '', // QA文档-答案开始标识符
  enable_extract_image: true,
  ai_chunk_size: 5000, // ai大模型分段最大字符数
  ai_chunk_model: '', // ai大模型分段模型名称
  ai_chunk_model_config_id: '', // ai大模型分段模型配置id
  ai_chunk_prumpt: defaultAiChunkPrumpt,
  ai_chunk_task_id: '', //  ai分段数据id，如果有ai分段数据就有值
  father_chunk_paragraph_type: 2,
  father_chunk_separators_no: [],
  father_chunk_chunk_size: 1024,
  son_chunk_separators_no: [],
  son_chunk_chunk_size: 512
}
let isEdit = false

const segmentationSettingRef = ref(null)

const onChangeSetting = (data) => {
  if (typeof data.separators_no == 'object') {
    data.separators_no = JSON.stringify(data.separators_no)
  }
  if (typeof data.father_chunk_separators_no == 'object') {
    data.father_chunk_separators_no = JSON.stringify(data.father_chunk_separators_no)
  }
  if (typeof data.son_chunk_separators_no == 'object') {
    data.son_chunk_separators_no = JSON.stringify(data.son_chunk_separators_no)
  }
  isEdit = true
  formData = {
    ...formData,
    ...data
  }
  if (itWasEdited) {
    Modal.confirm({
      title: t('reminder'),
      icon: createVNode(ExclamationCircleOutlined),
      content: t('reFetchConfirm'),
      okText: t('confirm'),
      okType: 'danger',
      cancelText: t('cancel'),
      onOk() {
        itWasEdited = false
        getDocumentFragment('create')
      },
      onCancel() {}
    })
  } else {
    itWasEdited = false
    getDocumentFragment('create')
  }
}

// 获取文档状态（是否已经转换为PDF）
let maxLoopNumber = 60 * 10
let loopNumber = 0
let library_id = null
const libFileInfo = ref({})

const library_type = ref(0)
const getDocumentStatus = () => {
  if (!spinning.value) {
    spinning.value = true
  }

  getLibFileInfo({ id: id }).then(async (res) => {
    const { status } = res.data
    library_type.value = res.data.library_type
    libFileInfo.value = res.data
    formData = {
      ...formData,
      separators_no: res.data.separators_no || '[12,11]',
      chunk_size: +res.data.chunk_size || 512,
      not_merged_text: res.data.not_merged_text == 'true',
      ai_chunk_size: +res.data.ai_chunk_size || 5000, // ai模型的默认值是5000
      chunk_overlap: +res.data.chunk_overlap || 50,
      is_qa_doc: library_type.value == 2 ? 1 : 0,
      question_lable: res.data.question_lable,
      answer_lable: res.data.answer_lable,
      question_column: res.data.question_column,
      answer_column: res.data.answer_column,
      enable_extract_image: res.data.enable_extract_image == 'true',
      qa_index_type: +res.data.qa_index_type,

      chunk_type: +res.data.chunk_type,
      semantic_chunk_size: +res.data.semantic_chunk_size || 512,
      semantic_chunk_overlap: +res.data.semantic_chunk_overlap || 50,
      semantic_chunk_threshold: +res.data.semantic_chunk_threshold || 90,
      semantic_chunk_use_model: res.data.semantic_chunk_use_model || '',
      ai_chunk_prumpt: res.data.ai_chunk_prumpt || t('defaultAiChunkPrumpt'),
      ai_chunk_model: res.data.ai_chunk_model || '',
      semantic_chunk_model_config_id:
        res.data.semantic_chunk_model_config_id > 0 ? res.data.semantic_chunk_model_config_id : '',
      ai_chunk_model_config_id:
        res.data.ai_chunk_model_config_id > 0 ? res.data.ai_chunk_model_config_id : '',
      ai_chunk_task_id: res.data.ai_chunk_task_id || '',

      father_chunk_paragraph_type: +res.data.father_chunk_paragraph_type || 2,
      father_chunk_separators_no: res.data.father_chunk_separators_no || '[12,11]',
      father_chunk_chunk_size: +res.data.father_chunk_chunk_size || 1024,
      son_chunk_separators_no: res.data.son_chunk_separators_no || '[8,10]',
      son_chunk_chunk_size: +res.data.son_chunk_chunk_size || 512
    }

    if (res.data.chunk_type == 0) {
      formData = {
        ...formData,
        separators_no: res.data.normal_chunk_default_separators_no,
        chunk_size: res.data.normal_chunk_default_chunk_size,
        not_merged_text: res.data.normal_chunk_default_not_merged_text == 'true',
        chunk_overlap: res.data.normal_chunk_default_chunk_overlap,
        chunk_type: +res.data.default_chunk_type,
        semantic_chunk_size: +res.data.semantic_chunk_default_chunk_size,
        semantic_chunk_overlap: +res.data.semantic_chunk_default_chunk_overlap,
        semantic_chunk_threshold: +res.data.semantic_chunk_default_threshold,
        semantic_chunk_use_model: res.data.default_use_model || '',
        semantic_chunk_model_config_id:
          res.data.default_model_config_id > 0 ? res.data.default_model_config_id : ''
      }
    }

    // await getInfo(res.data.library_id)
    if (status == 0) {
      loopNumber++
      if (loopNumber > maxLoopNumber) {
        Modal.error({
          title: t('reminder'),
          content: t('documentParsingSlow')
        })
        return
      }
      setTimeout(() => {
        getDocumentStatus()
      }, 1000)
    } else if (status == 4 || status == 2) {
      spinning.value = false
      settingMode.value = parseInt(res.data.is_table_file)
      library_id = res.data.library_id

      if (library_type.value == 2) {
        if (formData.question_lable || formData.question_column) {
          getDocumentFragment()
        }
      } else {
        getDocumentFragment()
      }
    }

    if (res.data.is_table_file == 1) {
      getExcelQaTitle()
    }
  })
}

onMounted(async () => {
  await getDocumentStatus()
})

const excellQaLists = ref([])
const getExcelQaTitle = () => {
  // 获取excel的QA  问题所在列 下拉列表
  getLibFileExcelTitle({ id: id }).then((res) => {
    let datas = []
    for (let key in res.data) {
      datas.push({
        lable: res.data[key],
        value: key
      })
    }
    excellQaLists.value = datas
  })
}

// 获取文档切片
const initDocumentFragmentList = ref([])
const isAiSave = ref(false)
const aiLoading = ref(false)
const task_id = ref('')
const error = ref(null)
let timer = null // 轮询定时器
const documentFragmentList = ref([])
const documentFragmentTotal = ref(0)
const isEmpty = computed(() => documentFragmentTotal.value <= 0)

const getDocumentFragment = (type) => {
  emit('loading', true)
  let params = {
    ...formData,
    semantic_chunk_model_config_id: formData.semantic_chunk_model_config_id
      ? +formData.semantic_chunk_model_config_id
      : 0,
    ai_chunk_model_config_id: formData.ai_chunk_model_config_id
      ? +formData.ai_chunk_model_config_id
      : 0,
    ...props.pdfState
  }

  if (formData.chunk_type == 3) {
    // 之前有分段数据可以传taskid过去直接获取
    if (formData.ai_chunk_task_id && !type) {
      params.ai_chunk_task_id = formData.ai_chunk_task_id
    }

    if (params.ai_chunk_task_id && type && type === 'create') {
      // 如果是主动点击的生成分段预览则不传taskid
      delete params.ai_chunk_task_id
    }
  }

  let requiredUrl = getLibFileSplit

  if (props.paragraphType && props.paragraphType === 'paragraphsSegmented' && type == 'create') {
    requiredUrl = getSplitParagraph
    delete params.id
  }

  if (props.paragraphType && props.paragraphType === 'paragraphsSegmented' && !type) {
    // 不请求接口，直接将数据填充进去
    documentFragmentList.value = props.currentData.list || []
    // 这里是点击重新分段选中的数据，只有一条
    documentFragmentTotal.value = props.currentData.list.length

    if (formData.chunk_type == 3) {
      // ai分段
      task_id.value = props.currentData.split_params?.ai_chunk_task_id || ''
      aiLoading.value = true
      segmentationSettingRef.value.reLoading = true
      error.value = null

      // 之前没有ai分段数据，重新异步请求ai分段数据
      // 之前不管有没有ai分段数据，只要是点击生成分段预览则重新异步请求ai分段数据
      if ((!formData.ai_chunk_task_id && settingMode.value != 1) || (type && type === 'create')) {
        // 清空之前的页面分段数据，重新请求
        documentFragmentList.value = []
        documentFragmentTotal.value = 0
        startPolling()
      } else {
        aiLoading.value = false
        segmentationSettingRef.value.reLoading = false
      }
    }

    if (formData.chunk_type != 3) {
      segmentationSettingRef.value.reLoading = false
      segmentationSettingRef.value.saveLoading = false
    }
    emit('loading', false)

    return false
  }

  return requiredUrl(params)
    .then((res) => {
      initDocumentFragmentList.value = res.data.list || []
      setInitDocumentFragmentList(initDocumentFragmentList.value)
      documentFragmentList.value = res.data.list || []
      documentFragmentTotal.value = res.data.list.length || 0

      if (formData.chunk_type == 3) {
        // ai分段
        task_id.value = res.data.split_params?.ai_chunk_task_id || ''
        aiLoading.value = true
        segmentationSettingRef.value.reLoading = true
        segmentationSettingRef.value.saveLoading = true
        error.value = null

        // 之前没有ai分段数据，重新异步请求ai分段数据
        // 之前不管有没有ai分段数据，只要是点击生成分段预览则重新异步请求ai分段数据
        if ((!formData.ai_chunk_task_id && settingMode.value != 1) || (type && type === 'create')) {
          // 清空之前的页面分段数据，重新请求
          documentFragmentList.value = []
          documentFragmentTotal.value = 0
          startPolling()
        } else {
          aiLoading.value = false
          segmentationSettingRef.value.reLoading = false
          segmentationSettingRef.value.saveLoading = false
        }
      }
    })
    .finally(() => {
      if (formData.chunk_type != 3) {
        segmentationSettingRef.value.reLoading = false
        segmentationSettingRef.value.saveLoading = false
      }
      emit('loading', false)
    })
}

const onChangeChunkType = (chunk_type) => {
  current_chunk_type.value = chunk_type
}

// 轮询查询结果
const pollData = async () => {
  try {
    const res = await getAiDocumentFragment()

    // 条件1: 接口返回错误信息
    if (res.data.err_msg) {
      error.value = res.data.err_msg
      return true // 停止轮询
    }

    // 条件2: 接口返回有效数据
    if (res.data.list?.length > 0) {
      initDocumentFragmentList.value = res.data.list || []
      setInitDocumentFragmentList(initDocumentFragmentList.value)
      documentFragmentList.value = res.data.list || []
      documentFragmentTotal.value = res.data.list.length || 0
      return true // 停止轮询
    }

    return false // 继续轮询
  } catch (err) {
    error.value = '请求异常，停止轮询'
    return true // 停止轮询
  }
}

function formatError(errorStr) {
  return errorStr.split('message:')[1]
}

// 启动轮询控制
const startPolling = () => {
  const executePoll = async () => {
    const shouldStop = await pollData()
    if (!shouldStop) {
      timer = setTimeout(executePoll, 3000) // 3秒后再次执行
    } else {
      aiLoading.value = false
      segmentationSettingRef.value.reLoading = false
      segmentationSettingRef.value.saveLoading = false
      if (timer !== null) {
        clearTimeout(timer)
        timer = null
      }
      if (isAiSave.value) {
        // 保存
        handleSaveLibFileSplit()
      }

      if (error.value) {
        let errorText = formatError(error.value)
        Modal.error({
          title: t('segmentationFailed'),
          content: errorText ? t('modelCallFailedWithReason', { reason: errorText }) : t('modelCallFailed')
        })
      }
    }
  }
  executePoll() // 立即执行首次查询
}

const getAiDocumentFragment = () => {
  return getLibFileSplitAiChunks({
    id: formData.id,
    task_id: task_id.value
  })
}

// 编辑文档片段
const editFragmentAlertRef = ref(null)
let editFragmentIndex = null

const handleEditFragment = ({ title, content, question, answer, images }, index) => {
  editFragmentIndex = index
  editFragmentAlertRef.value.open({ title, content, question, answer, images })
}

const saveFragment = ({ title, content, question, answer, images }) => {
  if (
    documentFragmentList.value[editFragmentIndex].title != title ||
    documentFragmentList.value[editFragmentIndex].content != content ||
    documentFragmentList.value[editFragmentIndex].question != question ||
    documentFragmentList.value[editFragmentIndex].answer != answer
  ) {
    itWasEdited = true
  }

  documentFragmentList.value[editFragmentIndex].title = title
  documentFragmentList.value[editFragmentIndex].content = content
  documentFragmentList.value[editFragmentIndex].question = question
  documentFragmentList.value[editFragmentIndex].answer = answer
  documentFragmentList.value[editFragmentIndex].images = images

  documentFragmentList.value[editFragmentIndex].word_total =
    answer.length + question.length + content.length
}

// 删除文档片段
const handleDeleteFragment = (index) => {
  Modal.confirm({
    title: t('reminder'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('deleteSegmentConfirm'),
    okText: t('confirm'),
    okType: 'danger',
    cancelText: t('cancel'),
    onOk() {
      itWasEdited = true
      documentFragmentList.value.splice(index, 1)
    },
    onCancel() {}
  })
}
const validateMessage = ref('')
const onValidate = (data) => {
  // 获取错误信息
  validateMessage.value = data
}

const saveLoading = ref(false)

const updataFormData = () => {
  let data = segmentationSettingRef.value.formState
  data = JSON.parse(JSON.stringify(data))
  if (typeof data.separators_no == 'object') {
    data.separators_no = JSON.stringify(data.separators_no)
  }
  if (typeof data.father_chunk_separators_no == 'object') {

    data.father_chunk_separators_no = JSON.stringify(data.father_chunk_separators_no)
  }
  if (typeof data.son_chunk_separators_no == 'object') {
    data.son_chunk_separators_no = JSON.stringify(data.son_chunk_separators_no)

  }
  formData = {
    ...formData,
    ...data
  }
}

const formatWoedTotal = (arr) => {
  let num = 0
  arr.forEach((item) => {
    num += parseInt(item.word_total)
  })
  return num
}
const handleSaveLibFileSplit = async () => {
  // 如果右侧的数据不是当前保存选中的分段类型则清空内容重新分段
  // 如果之前已经分段成功了，保存的时候不再另外分段了
  if (documentFragmentTotal.value <= 0 || formData.chunk_type != current_chunk_type.value) {
    updataFormData()
    await getDocumentFragment()
  }

  updataFormData()
  // await getDocumentFragment()
  if (validateMessage.value) {
    return message.error(validateMessage.value)
  }

  if (formData.chunk_type == 3 && !documentFragmentList.value.length) {
    // ai分段的保存
    isAiSave.value = true
    return false
  }

  if (task_id.value) {
    formData.ai_chunk_task_id = task_id.value
  }

  let split_params = {
    ...formData,
    semantic_chunk_model_config_id: formData.semantic_chunk_model_config_id
      ? +formData.semantic_chunk_model_config_id
      : 0,
    ai_chunk_model_config_id: formData.ai_chunk_model_config_id
      ? +formData.ai_chunk_model_config_id
      : 0,
    is_table_file: settingMode.value,
    ...props.pdfState
  }

  delete split_params.id

  let parmas = {
    file_id: id,
    word_total: formatWoedTotal(documentFragmentList.value),
    split_params: JSON.stringify(split_params),
    list: JSON.stringify(documentFragmentList.value),
    ...props.pdfState
  }

  // 非表格的也需要存储qa_index_type
  if (split_params.is_qa_doc == 1) {
    // 表格类型 + QA文档
    parmas.qa_index_type = split_params.qa_index_type
  }

  saveLoading.value = true

  if (props.paragraphType && props.paragraphType === 'paragraphsSegmented') {
    // get的接口和save没有统一，多了个s，这里要特殊处理一下
    parmas.data_id = parmas.data_ids
    delete parmas.data_ids
  }

  saveSplitParagraph(parmas)
    .then(() => {
      emit('ok')
      emit('finish')
    })
    .finally(() => {
      saveLoading.value = false
    })
    .catch((err) => {
      err.data && err.data.index && handleScrollToErrorDom(err.data.index)
      emit('finish')
    })
}

const previewBoxRef = ref(null)
const handleScrollToErrorDom = (index) => {
  index = index - 1
  let fragmentElements = previewBoxRef.value.querySelectorAll('.fragment-item')
  if (fragmentElements.length >= index) {
    let scorllElement = fragmentElements[index]
    scorllElement.scrollIntoView({ behavior: 'smooth', block: 'start', inline: 'nearest' })
  }
}

// 组件卸载时清理
onUnmounted(() => {
  if (timer) {
    clearTimeout(timer)
    timer = null
  }
})

defineExpose({
  handleSaveLibFileSplit
})
</script>
