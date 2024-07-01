<style lang="less" scoped>
.document-segmentation {
  position: relative;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-flow: column nowrap;
  padding-bottom: 56px;
  .page-container {
    display: flex;
    flex: 1;
    overflow: hidden;

    .page-left {
      width: 430px;
      height: 100%;
      overflow-y: auto;
      &::-webkit-scrollbar{
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
    <page-mini-title>文档分段与清洗</page-mini-title>

    <div class="page-container">
      <div class="page-left">
        <SegmentationSetting :excellQaLists="excellQaLists" :mode="settingMode" @change="onChangeSetting" @validate="onValidate" />
      </div>
      <div class="page-right">
        <div class="document-fragment-preview">
          <div class="preview-header">
            <span class="label-text">分段预览</span>
            <span class="fragment-number">共{{ documentFragmentTotal }}个分段</span>
          </div>
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
                @delete="handleDeleteFragment(index)"
                @edit="handleEditFragment(item, index)"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="footer-btn-box">
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="saveLoading" @click="handleSaveLibFileSplit"
        >保存</a-button
      >
    </div>
    <!-- 设置 -->
    <EditFragmentAlert ref="editFragmentAlertRef" @ok="saveFragment" />
  </div>
</template>

<script setup>
import { ref, createVNode } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import SegmentationSetting from './components/segmentation-setting.vue'
import DocumentFragment from './components/document-fragment.vue'
import EditFragmentAlert from './components/edit-fragment-alert.vue'
import { getLibFileSplit, getLibFileInfo, saveLibFileSplit, getLibFileExcelTitle } from '@/api/library/index'

const route = useRoute()
const router = useRouter()

const { document_id } = route.query
const spinning = ref(true)
const settingMode = ref(1) // 1 表格，0 非表格
let itWasEdited = false

let formData = {
  id: document_id,
  is_diy_split: 0, // 0 智能分段 1 自定义分段
  separators_no: '', // 自定义分段-分隔符序号集
  chunk_size: 512, // 自定义分段-分段最大长度 默认512，最大值不得超过2000
  chunk_overlap: 50, // 自定义分段-分段重叠长度 默认为50，最小不得低于10，最大不得超过最大分段长度的50%
  is_qa_doc: 0, // 0 普通文档 1 QA文档
  question_lable: '', // QA文档-问题开始标识符
  answer_lable: '' // QA文档-答案开始标识符
}

const onChangeSetting = (data) => {
  if (data.is_diy_split == 1) {
    data.is_qa_doc = 0
  }

  if (typeof data.separators_no == 'object') {
    data.separators_no = data.separators_no.join(',')
  }

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
const getDocumentStatus = () => {
  if (!spinning.value) {
    spinning.value = true
  }

  getLibFileInfo({ id: document_id }).then((res) => {
    const { status } = res.data

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
    } else if (status == 4) {
      spinning.value = false
      settingMode.value = parseInt(res.data.is_table_file)
      library_id = res.data.library_id
      getDocumentFragment()
    } else {
      router.replace('/library/details?id=' + res.data.library_id)
    }

    if(res.data.is_table_file == 1){
      getExcelQaTitle()
    }
  })
}

getDocumentStatus()

const excellQaLists = ref([]);
const getExcelQaTitle = () => {
  // 获取excel的QA  问题所在列 下拉列表
  getLibFileExcelTitle({ id: document_id }).then(res=>{
    let datas = [];
    for(let key in res.data){
      datas.push({
        lable: res.data[key],
        value: key
      })
    }
    excellQaLists.value = datas;
  })
}

// 获取文档切片
const documentFragmentList = ref([])
const documentFragmentTotal = ref(0)

const getDocumentFragment = () => {
  getLibFileSplit(formData).then((res) => {
    documentFragmentList.value = res.data.list || []
    documentFragmentTotal.value = res.data.list.length || 0
  })
}

// 编辑文档片段
const editFragmentAlertRef = ref(null)
let editFragmentIndex = null

const handleEditFragment = ({ title, content, question, answer  }, index) => {
  editFragmentIndex = index
  editFragmentAlertRef.value.open({ title, content, question, answer })
}

const saveFragment = ({ title, content, question, answer }) => {
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

  documentFragmentList.value[editFragmentIndex].word_total = answer.length + question.length + content.length
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
const validateMessage = ref('');
const onValidate = (data)=>{
  // 获取错误信息
  validateMessage.value = data
}

const saveLoading = ref(false)
const handleSaveLibFileSplit = () => {

  if(validateMessage.value) {
    return message.error(validateMessage.value)
  }

  let split_params = {
    ...formData,
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
  if(split_params.is_qa_doc == 1){
    // 表格类型 + QA文档
    parmas.qa_index_type = split_params.qa_index_type;
  }

  saveLoading.value = true

  saveLibFileSplit(parmas)
    .then(() => {
      message.success('保存成功')

      router.replace('/library/details?id=' + library_id)
    })
    .finally(() => {
      saveLoading.value = false
    }).catch(err=>{
      console.log(err,'==')
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

const previewBoxRef = ref(null);
const handleScrollToErrorDom = (index) => {
  index = index - 1
  let fragmentElements = previewBoxRef.value.querySelectorAll('.fragment-item')
  if(fragmentElements.length >= index){
    let scorllElement = fragmentElements[index];
    scorllElement.scrollIntoView({behavior: "smooth", block: "start", inline: "nearest"});
  }
}
</script>
