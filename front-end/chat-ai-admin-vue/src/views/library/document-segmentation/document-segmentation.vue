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
      width: 430px;
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
</style>

<template>
  <!-- loading -->
  <div class="loading-wrap" v-if="spinning">
    <a-spin></a-spin>
  </div>

  <div class="document-segmentation">
    <div class="breadcrumb-block">
      <a-breadcrumb>
        <a-breadcrumb-item>
          <router-link to="/library/list"> 知识库管理 </router-link></a-breadcrumb-item
        >
        <a-breadcrumb-item>
          <router-link
            :to="{
              path: '/library/details/knowledge-document',
              query: {
                id: libFileInfo.library_id
              }
            }"
          >
            {{ libFileInfo.library_name }}知识库</router-link
          >
        </a-breadcrumb-item>
        <a-breadcrumb-item>分段设置</a-breadcrumb-item>
      </a-breadcrumb>
    </div>
    <div class="page-title">
      <LeftOutlined @click="goBack" />
      <span class="title">文档分段与清洗</span>
      <!-- <div class="page-title-right" v-if="library_type == 0" style="margin-left: auto">
        <a-button type="primary" :loading="saveLoading" @click="handleSaveLibFileSplit"
          >保存</a-button
        >
      </div> -->
    </div>
    <div class="page-container">
      <div class="page-left">
        <SegmentationSetting
          :excellQaLists="excellQaLists"
          :libFileInfo="libFileInfo"
          :library_type="library_type"
          :mode="settingMode"
          ref="segmentationSettingRef"
          @change="onChangeSetting"
          @save="handleSaveLibFileSplit"
          @validate="onValidate"
        />
      </div>
      <div class="page-right">
        <div class="document-fragment-preview">
          <div class="preview-header">
            <span class="label-text">分段预览</span>
            <span class="fragment-number">共{{ documentFragmentTotal }}个分段</span>
          </div>
          <Empty v-if="isEmpty"></Empty>
          <div class="preview-box" ref="previewBoxRef">
            <div
              class="fragment-item"
              v-for="(item, index) in documentFragmentList"
              :key="item.number"
            >
              <DocumentFragment
                :number="item.number"
                :title="item.title"
                :content="item.content"
                :total="item.word_total"
                :question="item.question"
                :answer="item.answer"
                :images="item.images"
                :similar_question_list="item.similar_question_list"
                @delete="handleDeleteFragment(index)"
                @edit="handleEditFragment(item, index)"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- <div class="footer-btn-box" v-if="library_type == 2">
      <a-button type="primary" :loading="saveLoading" @click="handleSaveLibFileSplit"
        >保存</a-button
      >
    </div> -->
    <!-- 设置 -->
    <EditFragmentAlert ref="editFragmentAlertRef" @ok="saveFragment" />
  </div>
</template>

<script setup>
import { ref, createVNode, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import { ExclamationCircleOutlined, LeftOutlined } from '@ant-design/icons-vue'
import Empty from './components/empty.vue'
import SegmentationSetting from './components/segmentation-setting.vue'
import DocumentFragment from './components/document-fragment.vue'
import EditFragmentAlert from './components/edit-fragment-alert.vue'
import {
  getLibFileSplit,
  getLibFileInfo,
  saveLibFileSplit,
  getLibraryInfo,
  getLibFileExcelTitle
} from '@/api/library/index'

const route = useRoute()
const router = useRouter()

const { document_id } = route.query
const spinning = ref(true)
const settingMode = ref(1) // 1 表格，0 非表格
let itWasEdited = false

let formData = {
  id: document_id,
  separators_no: '', // 自定义分段-分隔符序号集
  chunk_size: 512, // 自定义分段-分段最大长度 默认512，最大值不得超过2000
  chunk_overlap: 50, // 自定义分段-分段重叠长度 默认为50，最小不得低于10，最大不得超过最大分段长度的50%
  is_qa_doc: 0, // 0 普通文档 1 QA文档
  question_lable: '', // QA文档-问题开始标识符
  similar_label: '', // QA文档-相似度标识符
  answer_lable: '', // QA文档-答案开始标识符
  enable_extract_image: true
}
let isEdit = false

const segmentationSettingRef = ref(null)

const onChangeSetting = (data) => {
  if (typeof data.separators_no == 'object') {
    data.separators_no = data.separators_no.join(',')
  }
  isEdit = true
  formData = {
    ...formData,
    ...data
  }
  if (itWasEdited) {
    Modal.confirm({
      title: '提醒',
      icon: createVNode(ExclamationCircleOutlined),
      content: '文档片段已被编辑重新获取文档片段会丢失当前修改过的文档片段内容，确定要重新获取吗？',
      okText: '确定',
      okType: 'danger',
      cancelText: '取消',
      onOk() {
        itWasEdited = false
        getDocumentFragment()
      },
      onCancel() {}
    })
  } else {
    itWasEdited = false
    getDocumentFragment()
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

  getLibFileInfo({ id: document_id }).then(async (res) => {
    const { status } = res.data
    library_type.value = res.data.library_type
    libFileInfo.value = res.data
    formData = {
      ...formData,
      separators_no: res.data.separators_no || '11,12',
      chunk_size: +res.data.chunk_size || 512,
      chunk_overlap: +res.data.chunk_overlap || 50,
      is_qa_doc: library_type.value == 2 ? 1 : 0,
      question_lable: res.data.question_lable,
      similar_label: res.data.similar_label,
      answer_lable: res.data.answer_lable,
      question_column: res.data.question_column,
      similar_column: res.data.similar_column,
      answer_column: res.data.answer_column,
      enable_extract_image: res.data.enable_extract_image == 'true',
      qa_index_type: +res.data.qa_index_type,

      chunk_type: +res.data.chunk_type,
      semantic_chunk_size: +res.data.semantic_chunk_size || 512,
      semantic_chunk_overlap: +res.data.semantic_chunk_overlap || 50,
      semantic_chunk_threshold: +res.data.semantic_chunk_threshold || 90,
      semantic_chunk_use_model: res.data.semantic_chunk_use_model || '',
      semantic_chunk_model_config_id:
        res.data.semantic_chunk_model_config_id > 0 ? res.data.semantic_chunk_model_config_id : ''
    }

    if (res.data.chunk_type == 0) {
      formData = {
        ...formData,
        separators_no: res.data.normal_chunk_default_separators_no,
        chunk_size: res.data.normal_chunk_default_chunk_size,
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
          title: '提醒',
          content: '文档解析速度慢请稍后再试'
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
      // library_type.value  != 2 && getDocumentFragment()
    } else {
      router.replace('/library/details?id=' + res.data.library_id)
    }

    if (res.data.is_table_file == 1) {
      getExcelQaTitle()
    }
  })
}

getDocumentStatus()

const excellQaLists = ref([])
const getExcelQaTitle = () => {
  // 获取excel的QA  问题所在列 下拉列表
  getLibFileExcelTitle({ id: document_id }).then((res) => {
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
const documentFragmentList = ref([])
const documentFragmentTotal = ref(0)
const isEmpty = computed(() => documentFragmentTotal.value <= 0)

const getDocumentFragment = () => {
  return getLibFileSplit({
    ...formData,
    semantic_chunk_model_config_id: formData.semantic_chunk_model_config_id
      ? +formData.semantic_chunk_model_config_id
      : 0
  })
    .then((res) => {
      documentFragmentList.value = res.data.list || []
      documentFragmentTotal.value = res.data.list.length || 0
    })
    .finally(() => {
      segmentationSettingRef.value.reLoading = false
    })
}

// 编辑文档片段
const editFragmentAlertRef = ref(null)
let editFragmentIndex = null

const handleEditFragment = (item, index) => {
  let { title, content, question, answer, images, similar_question_list } = item
  editFragmentIndex = index
  editFragmentAlertRef.value.open({ title, content, question, answer, images, similar_question_list })
}

const saveFragment = ({ title, content, question, answer, images,similar_question_list }) => {
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
  documentFragmentList.value[editFragmentIndex].similar_question_list = similar_question_list ? similar_question_list.split('\n') : []

  documentFragmentList.value[editFragmentIndex].word_total =
    answer.length + question.length + content.length
}

// 删除文档片段
const handleDeleteFragment = (index) => {
  Modal.confirm({
    title: '提醒',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确定要删除这个片段吗?',
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
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
    data.separators_no = data.separators_no.join(',')
  }
  formData = {
    ...formData,
    ...data
  }
}
const handleSaveLibFileSplit = async () => {
  updataFormData()
  if (validateMessage.value) {
    return message.error(validateMessage.value)
  }

  let split_params = {
    ...formData,
    semantic_chunk_model_config_id: formData.semantic_chunk_model_config_id
      ? +formData.semantic_chunk_model_config_id
      : 0,
    is_table_file: settingMode.value
  }

  delete split_params.id

  let parmas = {
    id: document_id,
    word_total: documentFragmentTotal.value,
    split_params: JSON.stringify(split_params),
    list: JSON.stringify(documentFragmentList.value)
  }

  // 非表格的也需要存储qa_index_type
  if (split_params.is_qa_doc == 1) {
    // 表格类型 + QA文档
    parmas.qa_index_type = split_params.qa_index_type
  }

  saveLoading.value = true

  saveLibFileSplit(parmas)
    .then(() => {
      message.success('保存成功')
      if (route.query.source == 'preview' && !isEdit) {
        goBack()
        return
      }
      router.replace('/library/details?id=' + library_id)
    })
    .finally(() => {
      saveLoading.value = false
    })
    .catch((err) => {
      console.log(err, '==')
      err.data && err.data.index && handleScrollToErrorDom(err.data.index)
    })
}

// 取消和上一步
const handleCancel = () => {
  Modal.confirm({
    title: '确定要退出吗?',
    icon: createVNode(ExclamationCircleOutlined),
    content: '',
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      router.replace('/library/details?id=' + library_id)
    },
    onCancel() {}
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
const goBack = () => {
  router.back()
}
</script>
