<template>
  <div class="library-preview-page">
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
                id: detailsInfo.library_id
              }
            }"
          >
            <div class="library-line1">{{ detailsInfo.library_name }}知识库</div>
          </router-link >
        </a-breadcrumb-item>
        <a-breadcrumb-item>知识库详情</a-breadcrumb-item>
      </a-breadcrumb>
      <div class="selct-star-box">
        <SelectStarBox :startLists="startLists" @change="handleCategoryChange" />
      </div>
    </div>
    <div class="document-title-block" style="height: 22px;">
      <LeftOutlined @click="goBack" />
      <span class="title">{{ detailsInfo.file_name }}</span>
      <SegmentationMode :detailsInfo="detailsInfo"></SegmentationMode>
    </div>
    <div class="search-content-block">
      <div class="preview-box">
        <a-button v-if="!isCollapsed" @click="toggleCollapse('put')"><svg-icon class="preview-box-icon" name="put-away"></svg-icon> 收起预览</a-button>
        <a-button v-else @click="toggleCollapse('expand')"><svg-icon class="preview-box-icon" name="expand"></svg-icon> 展开预览</a-button>
      </div>
      <a-space>
        <a-select
          v-model:value="filterData.status"
          placeholder="请选择"
          style="width: 160px"
          @change="search"
          v-if="is_qa_doc"
        >
          <a-select-option :value="-1">全部嵌入状态</a-select-option>
          <a-select-option v-for="(i, key) in listStatusMap" :key="i" :value="key">{{ i }}</a-select-option>
        </a-select>
        <a-select
          v-if="detailsInfo.graph_switch"
          v-model:value="filterData.graph_status"
          placeholder="请选择"
          @change="search"
          style="width: 160px"
        >
          <a-select-option :value="-1">全部知识图谱状态</a-select-option>
          <a-select-option v-for="(i, key) in graphStatusMap" :key="i" :value="key">{{ i }}</a-select-option>
        </a-select>
        <a-dropdown>
          <template #overlay>
            <a-menu>
              <a-menu-item @click="reEmbeddingVectors">重新嵌入向量</a-menu-item>
              <a-menu-item v-if="detailsInfo.graph_switch" @click="reExtractingGraph">重新抽取知识图谱</a-menu-item>
              <a-menu-item @click="handleRenameModal">重命名</a-menu-item>
              <a-menu-item v-if="doc_type == 2" @click="showUpdateFrequen"
                >设置更新频率</a-menu-item
              >
            </a-menu>
          </template>
          <a-button>其他操作 <DownOutlined /></a-button>
        </a-dropdown>
        <a-button v-if="doc_type != 3" @click="reSegment">重新分段</a-button>
        <a-button @click="openEditSubscription({})" type="primary">
          <template #icon>
            <PlusOutlined />
          </template>
          <span>添加分段</span>
        </a-button>
      </a-space>
    </div>
    <template v-if="pageLoading">
      <div class="page-loading-box">
        <a-spin />
      </div>
    </template>
    <template v-else>
      <template v-if="is_qa_doc">
        <!-- qa 类型文档  -->
        <Empty @openEditSubscription="openEditSubscription" v-if="isEmpty"></Empty>
        <cu-scroll :scrollbar="{ minSize: 0 }" v-else @onScrollEnd="onScrollEnd">
          <div class="content-block">
            <QaSubsectionBox
              ref="subsectionBoxRef"
              :total="total"
              :paragraphLists="paragraphLists"
              :detailsInfo="detailsInfo"
              @openEditSubscription="openEditSubscription"
              @handleDelParagraph="handleDelParagraph"
              @handleConvert="handleConvert"
              @getStatrList="getStatrList"
            ></QaSubsectionBox>
          </div>
        </cu-scroll>
      </template>
      <template v-else>
        <template v-if="isTableType || doc_type == 3">
          <!-- 表格类型 不支持预览 -->
          <Empty @openEditSubscription="openEditSubscription" v-if="isEmpty"></Empty>
          <cu-scroll :scrollbar="{ minSize: 0 }" v-else @onScrollEnd="onScrollEnd" :disableMouse="true">
            <div class="content-block">
              <SubsectionBox
                ref="subsectionBoxRef"
                :total="total"
                :detailsInfo="detailsInfo"
                :paragraphLists="paragraphLists"
                @openEditSubscription="openEditSubscription"
                @handleDelParagraph="handleDelParagraph"
                @handleConvert="handleConvert"
                @handleScrollTargetPage="paragraphPosition"
                @getStatrList="getStatrList"
                @handleSplit="handleSplit"
                @handleSplitNext="handleSplitNext"
                @handleSplitUp="handleSplitUp"
                @handleSplitDelete="handleSplitDelete"
                @handleSegmentation="handleSegmentation"
              ></SubsectionBox>
            </div>
          </cu-scroll>
        </template>
        <template v-else>
          <div class="pdf-view-box">
            <div class="view-content-wrap" :class="{ 'collapsed': isCollapsed }">
              <div v-if="detailsInfo?.file_ext == 'pdf'" class="pdf-mode-switch">
                <a-radio-group v-model:value="pdfPreviewMode">
                  <a-radio-button :value="1">
                    <a-tooltip
                      placement="right"
                      title="原文预览模式下，会展示pdf原始文件内容，手动编辑分段的不会展示。您可以通过原文预览模式，校验分段准确性，然后手动编辑分段。"
                    >
                      原文预览
                      <span v-if="pdfPreviewMode == 1 && totalPages > 0"
                        >({{ currentPdfPage }} / {{ totalPages }})</span
                      >
                    </a-tooltip>
                  </a-radio-button>
                  <a-radio-button :value="2">纯文本预览</a-radio-button>
                </a-radio-group>
              </div>
              <div class="segmentation-btn" v-if="detailsInfo?.file_ext == 'pdf'">
                <a-dropdown>
                  <template #overlay>
                    <a-menu @click="handleSegmentationMenuClick">
                      <a-menu-item :key="1" v-if="pdfPreviewMode == 1 && totalPages > 0">
                        <a-tooltip placement="right" title="将当前页重新解析并分段">
                          <div>将当页重新分段</div>
                        </a-tooltip>
                      </a-menu-item>
                      <a-menu-item :key="2">
                        <a-tooltip
                          placement="right"
                          title="将整个文档重新分段，注意，重新分段并不会重新提取文档的内容，只是将内容重新分段。"
                        >
                          <div>将文档重新分段</div>
                        </a-tooltip>
                      </a-menu-item>
                      <a-menu-item :key="3">
                        <a-tooltip
                          placement="right"
                          title="将整个文档重新学习，包含重新解析并提取文档内容，重新分段。"
                        >
                          <div>重新学习文档</div>
                        </a-tooltip>
                      </a-menu-item>
                    </a-menu>
                  </template>
                  <a-button>
                    重新分段
                    <DownOutlined />
                  </a-button>
                </a-dropdown>
              </div>
              <PdfPreview
                v-if="detailsInfo?.file_ext == 'pdf' && pdfPreviewMode == 1"
                ref="pdfRef"
                class="scroll-box pdf-render-box"
                :source="getFileUrl"
                @select="pdfPageSelect"
                @scroll="pdfPageOnScroll"
                @rendered="pdfPageRendered"
                @getTotalPage="getTotalPage"
              />
              <cu-scroll
                v-else
                :scrollbar="scrollbar"
                ref="leftScrollRef"
                class="scroll-box"
                @onScrollEnd="onScrollEnd"
              >
                <div class="view-content-box" style="padding-top: 60px">
                  <div
                    class="list-item"
                    @click="handleClickItem(index)"
                    v-for="(item, index) in paragraphLists"
                    :key="index"
                  >
                    <div class="content-box" v-if="item.question">{{ item.question }}</div>
                    <div class="content-box" v-if="item.answer">{{ item.answer }}</div>
                    <div class="content-box" v-html="item.content"></div>
                    <div class="fragment-img" v-viewer>
                      <img v-for="(item, index) in item.images" :key="index" :src="item" alt="" />
                    </div>
                  </div>
                </div>
              </cu-scroll>
            </div>
            <div class="right-contnt-box">
              <div class="right-tabs-box" v-if="!is_qa_doc">
                <a-tabs
                  v-model:activeKey="filterData.status"
                  style="background-color: #f2f4f7; padding: 0px 16px; border-radius: 6px;"
                  @change="search"
                >
                  <a-tab-pane v-for="i in generalListStatusMap" :key="i.value" :tab="`${i.label} (${i.num})`"></a-tab-pane>
                </a-tabs>
              </div>
              <Empty @openEditSubscription="openEditSubscription" v-if="isEmpty"></Empty>
              <cu-scroll
                :scrollbar="scrollbar"
                :disableMouse="true"
                ref="rightScrollRef"
                v-else
                @onScrollEnd="onScrollEnd"
              >
                <div class="content-block">
                  <SubsectionBox
                    ref="subsectionBoxRef"
                    :total="total"
                    :detailsInfo="detailsInfo"
                    :paragraphLists="paragraphLists"
                    @openEditSubscription="openEditSubscription"
                    @handleDelParagraph="handleDelParagraph"
                    @handleConvert="handleConvert"
                    @handleScrollTargetPage="paragraphPosition"
                    @getStatrList="getStatrList"
                    @handleSplit="handleSplit"
                    @handleSplitNext="handleSplitNext"
                    @handleSplitUp="handleSplitUp"
                    @handleSplitDelete="handleSplitDelete"
                    @handleSegmentation="handleSegmentation"
                  ></SubsectionBox>
                </div>
              </cu-scroll>
            </div>
          </div>
        </template>
      </template>
    </template>
    <EditSubscription
      :detailsInfo="detailsInfo"
      @handleEdit="handleEditParagraph"
      ref="editSubscriptionRef"
    ></EditSubscription>

    <UpdateFrequency @ok="saveFrequency" ref="updateFrequencyRef"></UpdateFrequency>
    <RenameModal @ok="handleSaveNameOk" ref="renameModalRef" />
    <ReSegmentationPage
      @ok="handleReSegmentationPageOk"
      @enable="handleeEableScroll"
      :detailsInfo="detailsInfo"
      ref="reSegmentationPageRef"
    />
    <SegmentationSettingModal ref="segmentationSettingModalRef" @ok="onSaveModal" />
  </div>
</template>

<script setup>
import { ref, reactive, computed, nextTick, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { LeftOutlined, PlusOutlined, DownOutlined } from '@ant-design/icons-vue'
import Empty from './components/empty.vue'
import SubsectionBox from './components/subsection-box.vue'
import {
  getParagraphCount,
  getLibFileInfo,
  saveLibFileSplit,
  getParagraphList,
  reconstructGraph,
  reconstructVector
} from '@/api/library'
import CuScroll from '@/components/cu-scroll/cu-scroll.vue'
import EditSubscription from './components/edit-subsection.vue'
import SegmentationMode from './components/segmentation-mode.vue'
import UpdateFrequency from './components/update-frequency.vue'
import RenameModal from '../library-details/components/rename-modal.vue'
import PdfPreview from '@/components/pdf-preview/pdf-preview.vue'
import { useUserStore } from '@/stores/modules/user.js'
import ReSegmentationPage from './components/re-segmentation-page.vue'
import QaSubsectionBox from './components/qa-subsection-box.vue'
import SelectStarBox from './components/selct-star-box.vue'
import { useCompanyStore } from '@/stores/modules/company'
import SegmentationSettingModal from './components/segmentation-setting-modal.vue'

const companyStore = useCompanyStore()
const neo4j_status = computed(()=>{
  return companyStore.companyInfo?.neo4j_status == 'true'
})

const isCollapsed = ref(false)
const pdfRef = ref(null)
const pdfLoadedPage = ref(0)
const pdfTotalPage = ref(0)
const graphStatusMap = {
  0: '待生成',
  4: '生成中',
  2: '生成成功',
  3: '生成失败'
}

const startLists = ref([])
const getStatrList = (data) => {
  startLists.value = data
}

const scrollbar = ref({
  fade: false,
  scrollbarTrackClickable: true,
  interactive: true,
  minSize: 40
})
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const isTableType = computed(() => {
  return detailsInfo.value.is_table_file == '1'
})
const doc_type = computed(() => {
  return detailsInfo.value.doc_type
})
const isEmpty = computed(() => {
  // 1学习中,2学习完成,3文件异常 4待用户切分 待学习
  let status = detailsInfo.value.status
  return !(status == '1' || status == '2') || paragraphLists.value.length == 0
})
const is_qa_doc = computed(() => {
  return detailsInfo.value.is_qa_doc == '1'
})
const detailsInfo = ref({})
const filterData = reactive({
  status: -1,
  graph_status: -1,
  category_id: -1
})
const pdfPreviewMode = ref(1)

const getFileUrl = computed(() => {
  return '/manage/getLibRawFile?id=' + route.query.id + '&token=' + userStore.getToken
})

const goBack = () => {
  router.back()
}

const leftScrollRef = ref(null)
const rightScrollRef = ref(null)
let lastIndex = -1
const handleClickItem = (index, type = 1) => {
  let elClass, scrollRef
  if (type == 1) {
    elClass = '.subsection-box .list-item'
    scrollRef = rightScrollRef.value
  } else {
    elClass = '.view-content-box .list-item'
    scrollRef = leftScrollRef.value
  }
  const el = document.querySelectorAll(elClass)[index]
  el.classList.add('flash-border')
  if (lastIndex >= 0) {
    document.querySelectorAll(elClass)[lastIndex].classList.remove('flash-border')
  }
  setTimeout(() => {
    lastIndex = -1
    el.classList.remove('flash-border')
  }, 2500)
  lastIndex = index
  scrollRef.scrollToElement({
    el,
    time: 1000
  })
}

const pageLoading = ref(true)
const library_type = ref(0)

const paginations = ref({
  page: 1,
  size: 10
})
const paragraphInfo = ref({})
const paragraphLists = ref([])
const total = ref(0)
const listStatusMap = {
  0: '未转换',
  1: '已转换',
  2: '转换异常',
  3: '转换中...'
}

const totalData = reactive({
  "split_status_exception": 0, // 分段失败
  "total": 0,// 全部
  "vector_status_converted": 0, // 已转换
  "vector_status_converting": 0, // 转换中
  "vector_status_exception": 0,// 转换失败
  "vector_status_initial": 0// 待转换
})

const generalListStatusMap = reactive([
  {
    value: -1,
    label: '全部',
    num: computed(() => totalData.total)
  },
  {
    value: 4,
    label: '分段失败',
    num: computed(() => totalData.split_status_exception)
  },
  {
    value: 0,
    label: '未转换',
    num: computed(() => totalData.vector_status_initial)
  },
  {
    value: 3,
    label: '转换中',
    num: computed(() => totalData.vector_status_converting)
  },
  {
    value: 2,
    label: '转换异常',
    num: computed(() => totalData.vector_status_exception)
  },
  {
    value: 1,
    label: '已转换',
    num: computed(() => totalData.vector_status_converted)
  }
])

const defaultAiChunkPrumpt = '你是一位文章分段助手，根据文章内容的语义进行合理分段，确保每个分段表述一个完整的语义，每个分段字数控制在500字左右，最大不超过1000字。请严格按照文章内容进行分段，不要对文章内容进行加工，分段完成后输出分段后的内容。'
const settingMode = ref(1) // 1 表格，0 非表格
let formData = {
  id: route.query.id,
  separators_no: '', // 自定义分段-分隔符序号集
  chunk_size: 512, // 自定义分段-分段最大长度 默认512，最大值不得超过2000
  chunk_overlap: 50, // 自定义分段-分段重叠长度 默认为50，最小不得低于10，最大不得超过最大分段长度的50%
  is_qa_doc: 0, // 0 普通文档 1 QA文档
  question_lable: '', // QA文档-问题开始标识符
  similar_label: '', // QA文档-相似度标识符
  answer_lable: '', // QA文档-答案开始标识符
  enable_extract_image: true,
  ai_chunk_size: 5000, // ai大模型分段最大字符数
  ai_chunk_model:'', // ai大模型分段模型名称
  ai_chunk_model_config_id: '', // ai大模型分段模型配置id
  ai_chunk_prumpt: defaultAiChunkPrumpt,
  ai_chunk_task_id: ''  //  ai分段数据id，如果有ai分段数据就有值
}
const saveLoading = ref(false)

let isLoading = true

const getFileInfo = async () => {
  getLibFileInfo({
    id: route.query.id
  }).then((res) => {
    library_type.value = res.data.library_type
    detailsInfo.value = {
      ...res.data,
      graph_switch: res.data.graph_switch == '1' && neo4j_status.value
    }
    let status = res.data.status
    if (status == 1 || status == 2) {
      getParagraphLists()
    } else {
      pageLoading.value = false
    }

    formData = {
      ...formData,
      separators_no: res.data.separators_no || '11,12',
      chunk_size: +res.data.chunk_size || 512,
      ai_chunk_size: +res.data.ai_chunk_size || 5000, // ai模型的默认值是5000
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
      ai_chunk_prumpt: res.data.ai_chunk_prumpt || '',
      ai_chunk_model: res.data.ai_chunk_model || '',
      semantic_chunk_model_config_id:
        res.data.semantic_chunk_model_config_id > 0 ? res.data.semantic_chunk_model_config_id : '',
      ai_chunk_model_config_id:
        res.data.ai_chunk_model_config_id > 0 ? res.data.ai_chunk_model_config_id : '',
      ai_chunk_task_id: res.data.ai_chunk_task_id || ''
    }

    if (res.data.chunk_type == 0) {
      formData = {
        ...formData,
        separators_no: res.data.normal_chunk_default_separators_no,
        chunk_size: +res.data.normal_chunk_default_chunk_size,
        chunk_overlap: +res.data.normal_chunk_default_chunk_overlap,
        chunk_type: +res.data.default_chunk_type,
        semantic_chunk_size: +res.data.semantic_chunk_default_chunk_size,
        semantic_chunk_overlap: +res.data.semantic_chunk_default_chunk_overlap,
        semantic_chunk_threshold: +res.data.semantic_chunk_default_threshold,
        semantic_chunk_use_model: res.data.default_use_model || '',
        semantic_chunk_model_config_id: res.data.default_model_config_id > 0 ? res.data.default_model_config_id : ''
      }
    }

    if (status == 4 || status == 2) {
      settingMode.value = parseInt(res.data.is_table_file)
    }

  })
}
getFileInfo()

const getParagraphCountFn = () => {
  getParagraphCount({
    file_id: route.query.id
  }).then((res) => {
    Object.assign(totalData, res.data)
  })
}

const onSaveModal = () => {
  // 先使用最简单的方式，刷新页面
  window.location.reload();
  // getParagraphLists()
  // getParagraphCountFn()
}

const search = () => {
  paginations.value.page = 1
  total.value = 0
  paragraphLists.value = []
  getParagraphLists()
}

const handleCategoryChange = (id)=>{
  filterData.category_id = id
  search()
}

const loadMore = (callback = null) => {
  paginations.value.page++
  getParagraphLists(callback)
}

const getParagraphLists = (callback = null) => {
  isLoading = true
  getParagraphList({
    file_id: route.query.id,
    ...paginations.value,
    ...filterData
  }).then((res) => {
    pageLoading.value = false
    isLoading = false
    let data = res.data
    let list = data.list || []
    let isExcelQa = data.info?.is_qa_doc == '1' && data.info?.is_table_file == '1'
    list.forEach((item) => {
      item.status_text = listStatusMap[item.status]
      item.graph_status_text = graphStatusMap[item.graph_status]
      item.isExcelQa = isExcelQa
      if(item.similar_questions){
        item.similar_questions = JSON.parse(item.similar_questions)
      }
    })

    paragraphInfo.value = data.info
    if (paginations.value.page == 1) {
      paragraphLists.value = []
    }
    paragraphLists.value = [...paragraphLists.value, ...list]
    total.value = data.total
    typeof callback === 'function' && callback()
  })
}
const onScrollEnd = () => {
  if (isLoading) {
    return
  }
  if (paragraphLists.value.length >= total.value) {
    return
  }
  loadMore()
  // 同时加载pdf
  pdfRef.value && pdfRef.value.loadMore()
}

const handleConvert = () => {
  updataList()
}

let timer = null

function setIntervalUpStatus() {
  clearInterval(timer)
  // if()
  let list = paragraphLists.value.filter((item) => item.status == 3)
  if (list.length) {
    timer = setInterval(() => {
      updataList()
    }, 5000)
  } else {
    clearInterval(timer)
  }
}

const updataList = () => {
  getParagraphList({
    file_id: route.query.id,
    page: 1,
    size: paragraphLists.value.length,
    ...filterData
  }).then((res) => {
    let data = res.data
    let list = data.list || []
    let isExcelQa = data.info?.is_qa_doc == '1' && data.info?.is_table_file == '1'
    list.forEach((item) => {
      item.status_text = listStatusMap[item.status]
      item.graph_status_text = graphStatusMap[item.graph_status]
      item.isExcelQa = isExcelQa
      if(item.similar_questions){
        item.similar_questions = JSON.parse(item.similar_questions)
      }
    })
    paragraphLists.value = list
    setIntervalUpStatus()
  })
}
const handleEditParagraph = (data) => {
  if (!data.id) {
    paginations.value.page = 1
    getParagraphLists()
    return
  }
  // 更新分段内容 无刷更新
  let lists = paragraphLists.value
  let index = lists.findIndex((item) => item.id == data.id)
  if (index > -1) {
    let lastItem = lists[index]
    lastItem.title = data.title
    lastItem.content = data.content
    lastItem.question = data.question
    lastItem.answer = data.answer
    lastItem.images = data.images
    lastItem.category_id = data.category_id
    lastItem.word_total = data.question.length + data.answer.length + data.content.length
    lastItem.similar_questions = data.similar_questions.length ? JSON.parse(data.similar_questions) : []
    lists.splice(index, 1, lastItem)
    paragraphLists.value = lists
  }
}
const handleDelParagraph = (id) => {
  // 无刷新 删除列表
  let lists = paragraphLists.value
  let index = lists.findIndex((item) => item.id == id)
  if (index > -1) {
    lists.splice(index, 1)
    paragraphLists.value = lists
    total.value = total.value--
  }
}

const editSubscriptionRef = ref(null)
const openEditSubscription = (data) => {
  editSubscriptionRef.value.showModal(JSON.parse(JSON.stringify(data)))
}

const updateFrequencyRef = ref(null)

const reEmbeddingVectors = () => {
  reconstructVector({ id: route.query.id }).then(() => message.success('操作完成'))
}

const reExtractingGraph = () => {
  reconstructGraph({ id: route.query.id }).then(() => message.success('操作完成'))
}
const renameModalRef = ref(null)
const handleRenameModal = () => {
  renameModalRef.value.show(detailsInfo.value)
}

const reSegment = () => {
  router.push('/library/document-segmentation?document_id=' + route.query.id + '&source=preview')
}

const showUpdateFrequen = () => {
  updateFrequencyRef.value.open({
    id: route.query.id,
    doc_auto_renew_frequency: detailsInfo.value.doc_auto_renew_frequency
  })
}

const saveFrequency = (data) => {
  detailsInfo.value.doc_auto_renew_frequency = data
}

const handleSaveNameOk = (file_name) => {
  detailsInfo.value.file_name = file_name
}

const pdfPageSelect = (page) => {
  const link = () => {
    for (let i in paragraphLists.value) {
      // 查询该页的第一个分段
      if (paragraphLists.value[i].page_num == page) {
        handleClickItem(i)
        return
      }
    }
  }
  if (page > paragraphLists.value[paragraphLists.value.length - 1].page_num) {
    message.warning('分段数据正在加载，请稍后...')
    const checkLoadFinished = () => {
      if (page <= paragraphLists.value[paragraphLists.value.length - 1].page_num) {
        message.success('加载完成')
        nextTick(() => {
          link()
        })
      } else {
        loadMore(checkLoadFinished)
      }
    }
    loadMore(checkLoadFinished)
  } else {
    link()
  }
}

const pdfPageRendered = (page, total) => {
  pdfTotalPage.value = total
  pdfLoadedPage.value = page
  const maxPage = paragraphLists.value[paragraphLists.value.length - 1].page_num
  // 当pdf最大页分段未加载时，加载分段
  if (maxPage < page) {
    loadMore()
  }
}

const paragraphPosition = (data) => {
  if (detailsInfo.value?.file_ext == 'pdf' && pdfPreviewMode.value == 1) {
    const index = data.page_num - 1
    const getPageEl = () =>
      document.querySelectorAll('.pdf-render-box > .scroll-content .vue-pdf-embed')[index]
    const link = () => {
      const el = getPageEl()
      el.classList.add('flash-border')
      setTimeout(() => el.classList.remove('flash-border'), 2500)
      pdfRef.value.getScrollInstance().scrollToElement({ el, time: 1000 })
    }
    const el = getPageEl()
    if (!el) {
      // 未加载到pdf时
      message.warning('PDF正在加载，请稍后...')
      const checkLoadFinished = () => {
        if (!getPageEl()) {
          pdfRef.value && pdfRef.value.loadMore()
        } else {
          message.success('加载完成')
          link()
        }
      }
      pdfRef.value && pdfRef.value.loadMore(checkLoadFinished)
    } else {
      link()
    }
  } else {
    handleClickItem(data.index, 2)
  }
}
const currentPdfPage = ref(1)
let pdfHeight = null

const totalPages = ref(0)
const getTotalPage = (data) => {
  totalPages.value = data
}
const pdfPageOnScroll = ({ y }) => {
  if (!pdfHeight) {
    // 一页pdf的高度
    const el = document.querySelectorAll('.pdf-render-box > .scroll-content .vue-pdf-embed')[0]
    pdfHeight = el && el.clientHeight
  }
  // 计算一下 现在滚到pdf页码
  let page = Math.floor(Math.abs(y) / pdfHeight) + 1
  currentPdfPage.value = Math.max(page, 1)
}

const reSegmentationPageRef = ref(null)
const segmentationSettingModalRef = ref(null)

const handleSegmentationMenuClick = (e) => {
  let { key } = e
  // console.log(pdfRef.value.getScrollInstance(),'==')
  pdfRef.value && pdfRef.value.getScrollInstance().disable()
  if (key == 1) {
    reSegmentationPageRef.value.show({
      type: key,
      pdf_page_num: currentPdfPage.value
    })
  }
  if (key == 2) {
    reSegment()
  }
  if (key == 3) {
    reSegmentationPageRef.value.show({
      type: key,
      pdf_page_num: 0
    })
  }
}

const handleSegmentation = (item) => {
  if (segmentationSettingModalRef.value) {
    segmentationSettingModalRef.value.show(item)
  }
}

const handleeEableScroll = () => {
  pdfRef.value && pdfRef.value.getScrollInstance().enable()
}
const handleReSegmentationPageOk = () => {
  handleeEableScroll()
  search()
}

const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
}

// 单独分段
const createNewParagraph = (item, content) => {
  let newItem = JSON.parse(JSON.stringify(item))
  newItem.content = content
  newItem.word_total = content.length
  return newItem
}

const isFullSelectionRef = ref(false)
const splitType = ref(1) // 1单独分段 2合并到下一分段 3合并到上一分段 4删除

// 单独分段
const handleSplit = ({ index, beforeContent, selectedContent, afterContent, isFullSelection }) => {
  const newList = [...paragraphLists.value]

  isFullSelectionRef.value = isFullSelection
  splitType.value = 1
  if (isFullSelection) {
    // 插入选中内容段落
    newList.splice(index + 1, 0, createNewParagraph(newList[index], selectedContent))
    // 直接删除当前段落
    newList.splice(index, 1);
  } else {
    // 更新原段落
    newList[index].content = beforeContent || ''
    newList[index].word_total = beforeContent.length || 0

    // 插入剩余内容段落
    if (afterContent.trim()) {
      // 插入选中内容段落
      newList.splice(index + 1, 0, createNewParagraph(newList[index], afterContent))
    }

    // 插入选中内容段落
    newList.splice(index + 1, 0, createNewParagraph(newList[index], selectedContent))

    // 不是全选最终数据长度会加1
    newList[index + 1].images = []
  }

  paragraphLists.value = newList

  // 这里直接调用保存接口/manage/saveLibFileSplit
  handleSaveLibFileSplit(newList, index)
}

// 合并到下一个分段
const handleSplitNext = ({ index, beforeContent, selectedContent, afterContent, isFullSelection }) => {
  const newList = [...paragraphLists.value]

  isFullSelectionRef.value = isFullSelection
  splitType.value = 2

  if (isFullSelection) {
    // 合并内容到下一段落
    newList[index + 1].content = selectedContent + newList[index + 1].content;
    newList[index + 1].word_total += selectedContent.length
    // 删除当前段落
    newList.splice(index, 1);
    paragraphLists.value = newList;
    handleSaveLibFileSplit(newList, index);
  } else {
    // 更新原段落
    newList[index].content = beforeContent || ''
    newList[index].word_total = beforeContent.length || 0

    // 修改下一个段落的内容
    newList[index + 1].content = selectedContent + newList[index + 1].content;

    // 插入剩余内容到原段落
    if (afterContent.trim()) {
      newList[index].content += afterContent || ''
      newList[index].word_total += afterContent.length || 0
    }

    paragraphLists.value = newList

    // 这里直接调用保存接口/manage/saveLibFileSplit
    handleSaveLibFileSplit(newList, index + 1)
  }
}

// 合并到上一个分段
const handleSplitUp = ({ index, beforeContent, selectedContent, afterContent, isFullSelection }) => {
  const newList = [...paragraphLists.value]

  isFullSelectionRef.value = isFullSelection
  splitType.value = 3

  if (isFullSelection) {
    // 合并内容到上一段落
    newList[index - 1].content += selectedContent
    newList[index - 1].word_total += selectedContent.length
    // 删除当前段落
    newList.splice(index, 1);
    paragraphLists.value = newList;
    handleSaveLibFileSplit(newList, index);
  } else {
    // 更新原段落
    newList[index].content = beforeContent || ''
    newList[index].word_total = beforeContent.length || 0

    // 修改上一个段落的内容
    newList[index - 1].content += selectedContent

    // 插入剩余内容段落
    if (afterContent.trim()) {
      newList[index].content += afterContent || ''
      newList[index].word_total += afterContent.length || 0
    }

    paragraphLists.value = newList

    // 这里直接调用保存接口/manage/saveLibFileSplit
    handleSaveLibFileSplit(newList, index - 1)
  }
}

// 删除
const handleSplitDelete = ({ index, beforeContent, afterContent, isFullSelection }) => {
  const newList = [...paragraphLists.value]

  isFullSelectionRef.value = isFullSelection
  splitType.value = 4

  if (isFullSelection) {
    // 删除当前段落
    newList.splice(index, 1);
  } else {
    // 更新原段落
    newList[index].content = beforeContent || ''
    newList[index].word_total = beforeContent.length || 0

    // 插入剩余内容段落
    if (afterContent.trim()) {
      newList[index].content += afterContent || ''
      newList[index].word_total += afterContent.length || 0
    }
  }

  paragraphLists.value = newList

  // 这里直接调用保存接口/manage/saveLibFileSplit
  handleSaveLibFileSplit(newList, index)
}

function formatList (array, keys, renameMap = {}, index) {
  let sourceNumber = index + 1  // 索引从0开始
  if (!Array.isArray(array) || !Array.isArray(keys)) {
    return JSON.stringify([]);
  }
  const newArray = array.map((originalObj, originalIndex) => {
      const newObj = {};
      for (const key of keys) {
        const sourceKey = renameMap[key] || key;
        const sourceValue = originalObj[sourceKey];

        if (!Object.prototype.hasOwnProperty.call(originalObj, sourceKey)) continue;

        if (key === 'number') {
          newObj[key] = originalIndex + 1;
        } else if (key === 'page_num' || key === 'word_total') {
          newObj[key] = +sourceValue;
        } else if (key === 'content') {
          if (!sourceValue) return null;
          newObj[key] = sourceValue;
        } else if (key === 'images') {
          if (splitType.value == 1) {
            // 单独分段
            // 是否全选
            if (isFullSelectionRef.value) {
              // 全选最终的数据长度不会改变，因为全选的源数据被删除，然后新插入数据体
              if (sourceNumber == newObj.number) {
                // 如果之前的数据有图片，还是携带
                newObj[key] = sourceValue
              } else {
                // 其他的数据该有图片就有，没有就没有保持原样
                newObj[key] = sourceValue
              }
            } else {
              // 不是全选最终数据长度会加1
              if (sourceNumber + 1 == newObj.number) {
                newObj[key] = []
              } else {
                newObj[key] = sourceValue
              }
            }
          } else if (splitType.value == 2) {
            // 2合并到下一分段
            // 是否全选
            if (isFullSelectionRef.value) {
              // 全选最终的数据长度-1，因为全选的源数据被删除，然后修改了下一分段的内容
              if (sourceNumber == newObj.number) {
                newObj[key] = sourceValue
              } else {
                newObj[key] = sourceValue
              }
            } else {
              // 不是全选最终数据不会发生改变
              if (sourceNumber == newObj.number) {
                newObj[key] = sourceValue
              } else {
                newObj[key] = sourceValue
              }
            }
          } else if (splitType.value == 3) {
            // 3合并到上一分段
            // 是否全选
            if (isFullSelectionRef.value) {
              // 全选最终的数据长度-1，因为全选的源数据被删除，然后修改了下一分段的内容
              if (sourceNumber - 1 == newObj.number) {
                newObj[key] = sourceValue
              } else {
                newObj[key] = sourceValue
              }
            } else {
              // 不是全选最终数据不会发生改变
              if (sourceNumber == newObj.number) {
                newObj[key] = sourceValue
              } else {
                newObj[key] = sourceValue
              }
            }
          } else if (splitType.value == 4) {
            // 4删除
            // 是否全选
            if (isFullSelectionRef.value) {
              // 全选最终的数据长度-1，因为全选的源数据被删除
              newObj[key] = sourceValue
            } else {
              // 不是全选最终数据不会发生改变
              newObj[key] = sourceValue
            }
          }
        } else {
          newObj[key] = sourceValue;
        }
      }
      return newObj;
    })
    .filter(obj => obj !== null && (obj.content !== undefined && obj.content !== ''));

  return JSON.stringify(newArray);
}

const renameMap = {
  similar_question_list: 'similar_questions' // 新属性名: 旧属性名
};

const handleSaveLibFileSplit = async (documentFragmentList, index) => {
  // 如果右侧的数据不是当前保存选中的分段类型则清空内容重新分段
  // 如果之前已经分段成功了，保存的时候不再另外分段了

  let split_params = {
    ...formData,
    semantic_chunk_model_config_id: formData.semantic_chunk_model_config_id
      ? +formData.semantic_chunk_model_config_id
      : 0,
    ai_chunk_model_config_id: formData.ai_chunk_model_config_id
      ? +formData.ai_chunk_model_config_id
      : 0,
    is_table_file: settingMode.value
  }

  delete split_params.id

  let parmas = {
    id: route.query.id,
    word_total: total.value,
    split_params: JSON.stringify(split_params),
    list: formatList(documentFragmentList, ['number', 'page_num', 'title', 'content', 'question', 'similar_question_list', 'answer', 'word_total', 'images'], renameMap, index)
  }

  // 非表格的也需要存储qa_index_type
  if (split_params.is_qa_doc == 1) {
    // 表格类型 + QA文档
    parmas.qa_index_type = split_params.qa_index_type
  }

  saveLoading.value = true

  saveLibFileSplit(parmas)
    .then((res) => {
      message.success('操作成功')
      // 后端返回的id替换成新的id
      const newIds = res.data
      for (let index = 0; index < newIds.length; index++) {
        const item = newIds[index];
        const pItem = paragraphLists.value[index]
        if (index < paragraphLists.value.length) {
          pItem.id = item.toString()
          pItem.status = '0'
          pItem.status_text = listStatusMap[pItem.status]
        } else {
          break;
        }
      }
    })
    .finally(() => {
      saveLoading.value = false
    })
}

onMounted(() => {
  getParagraphCountFn()
  if (route.query.graph_status) {
    filterData.graph_status = route.query.graph_status
    search()
  }
})
</script>

<style lang="less" scoped>
.breadcrumb-block {
  width: 100%;
  position: relative;
  height: 56px;
  display: flex;
  align-items: center;

  :deep(.ant-breadcrumb-link) {
    display: inline-flex;
  }

  .library-line1 {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 180px;
  }
}

.selct-star-box {
  position: absolute;
  top: 10px;
  left: 50%;
  transform: translateX(-50%);
}
.library-preview-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  .page-loading-box {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .document-title-block {
    margin-bottom: 16px;
    margin-top: 12px;
    width: 100%;
    height: 38px;
    border-radius: 2px;
    display: flex;
    align-items: center;
    color: #262626;
    .anticon-left {
      cursor: pointer;
    }
    .title {
      font-weight: 600;
      font-size: 14px;
      line-height: 22px;
      margin-left: 8px;
    }
  }

  .search-content-block{
    display: flex;
    align-items: center;
    margin-bottom: 16px;
    justify-content: space-between;

    .preview-box {
      display: flex;
      gap: 10px;
      color: #595959;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;

      .preview-box-icon {
        margin-right: 8px;
      }
    }
  }

  .pdf-view-box {
    flex: 1;
    display: flex;
    overflow: hidden;
    .view-content-wrap {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;
      position: relative;
      transition: all 0.3s ease;

      .pdf-mode-switch {
        position: absolute;
        top: 12px;
        left: 12px;
        z-index: 999;
      }
      .segmentation-btn {
        position: absolute;
        right: 36px;
        top: 12px;
        z-index: 999;
      }

      .scroll-box {
        flex: 1;
        margin-right: 24px;
        border: 1px solid #ccc;
        border-radius: 4px;
      }
      .view-content-box {
        padding: 12px 32px 24px;
        .list-item {
          cursor: pointer;
          padding: 6px 8px;
          border-radius: 4px;
          border: 1px solid #fff;
          transition: all 0.3s;
          &:hover {
            background: #eee;
            border: 1px solid #e0e0e0;
          }
        }
        .content-box {
          color: #333;
          font-size: 14px;
          font-weight: 400;
          line-height: 22px;
          margin-top: 4px;
          white-space: pre-wrap;
          word-wrap: break-word;
        }
        .fragment-img {
          display: flex;
          gap: 8px;
          align-items: baseline;
          flex-wrap: wrap;
          img {
            max-width: 100%;
            max-height: 100%;
            width: auto;
            height: fit-content;
          }
        }
      }
    }

    .view-content-wrap.collapsed {
      width: 0;
      opacity: 0;
      flex: 0;
      transform: translateX(-100%);
    }

    .right-contnt-box {
      flex: 1;
      overflow: hidden;
      transition: all 0.3s ease;
    }
  }

  ::v-deep(.bscroll-vertical-scrollbar) {
    z-index: 999 !important;
  }
}
.custom-select-box {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  > span {
    white-space: nowrap;
  }
  :deep(.ant-select-selector) {
    border: none !important;
    padding-left: 0 !important;
    height: unset !important;
  }
  padding-left: 8px;
}
@keyframes flash-border {
  0%,
  100% {
    background: transparent;
  }
  50% {
    background: #c8d9f4;
  }
}

.flash-border {
  background: #c8d9f4;
  animation: flash-border 1s infinite; /* 持续时间1秒，无限次重复 */
}
.mb8 {
  margin-bottom: 8px;
}
</style>
