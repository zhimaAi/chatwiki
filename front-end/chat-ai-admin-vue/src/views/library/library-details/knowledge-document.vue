<template>
  <div class="details-library-page">
    <cu-scroll :scrollbar="false">
      <div class="list-tools">
        <div class="tools-items">
          <div class="tool-item">
            <a-dropdown :trigger="['click']" overlayClassName="add-dropdown-btn">
              <template #overlay>
                <a-menu @click="handleMenuClick">
                  <a-menu-item :key="1">
                    <div class="dropdown-btn-menu">
                      <a-flex class="title-block" :gap="4">
                        <svg-icon name="doc-icon"></svg-icon>
                        <div class="title">本地文档</div>
                      </a-flex>
                      <div class="desc" v-if="libraryInfo.type == 2">上传本地 docx/csv/xlsx 等格式文件</div>
                      <div class="desc" v-else>上传本地 pdf/docx/ofd/txt/md/xlsx/csv/html 等格式文件</div>
                    </div>
                  </a-menu-item>
                  <a-menu-item :key="2" v-if="libraryInfo.type != 2">
                    <div class="dropdown-btn-menu">
                      <a-flex class="title-block" :gap="4">
                        <svg-icon name="link-icon"></svg-icon>
                        <div class="title">在线数据</div>
                      </a-flex>
                      <div class="desc">获取在线网页内容</div>
                    </div>
                  </a-menu-item>
                  <a-menu-item :key="3">
                    <div class="dropdown-btn-menu">
                      <a-flex class="title-block" :gap="4">
                        <svg-icon name="cu-doc-icon"></svg-icon>
                        <div class="title">自定义文档</div>
                      </a-flex>
                      <div class="desc">自定义一个空文档，手动添加内容</div>
                    </div>
                  </a-menu-item>
                </a-menu>
              </template>
              <a-button type="primary">
                <template #icon>
                  <PlusOutlined />
                </template>
                <span>添加内容</span>
              </a-button>
            </a-dropdown>
          </div>

          <a-flex align="center" class="tool-item custom-select-box">
            <span>嵌入模型：</span>
            <model-select
              modelType="TEXT EMBEDDING"
              :isOffline="false"
              :modeName="modelForm.use_model"
              :modeId="modelForm.model_config_id"
              style="width: 240px"
              @change="onChangeModel"
              @loaded="onVectorModelLoaded"
            />
          </a-flex>
          <a-flex align="center" class="tool-item custom-select-box pd-5-8">
            <span>生成知识图谱：</span>
            <a-switch
              v-model:checked="createGraphSwitch"
              @change="createGraphSwitchChange"
              checked-children="开"
              un-checked-children="关"
            />
          </a-flex>
        </div>
        <div>
          <div class="tool-item">
            <a-input
              style="width: 282px"
              v-model:value="queryParams.file_name"
              placeholder="请输入文档名称搜索"
              @change="onSearch"
            >
              <template #suffix>
                <SearchOutlined @click="onSearch" style="color: rgba(0, 0, 0, 0.25)" />
              </template>
            </a-input>
          </div>
        </div>
      </div>
      <div class="list-content">
        <a-table
          :columns="columns"
          :data-source="fileList"
          :pagination="{
            current: queryParams.page,
            total: queryParams.total,
            pageSize: queryParams.size,
            showQuickJumper: true,
            showSizeChanger: true,
            pageSizeOptions: ['10', '20', '50', '100']
          }"
          @change="onTableChange"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'file_name'">
              <div class="doc-name-td">
                <a @click="handlePreview(record)">
                  <span v-if="['5', '6', '7'].includes(record.status)">{{ record.doc_url }}</span>
                  <span v-else>{{ record.file_name }}</span>
                </a>
                <div v-if="record.doc_type == 2 && record.remark" class="url-remark">
                  备注：{{ record.remark }}
                </div>
              </div>
            </template>
            <template v-if="column.key === 'status'">
              <span class="status-tag running" v-if="record.status == 0"><a-spin size="small" /> 转换中</span>
              <span class="status-tag running" v-if="record.status == 1"><a-spin size="small" /> 学习中</span>
              <span class="status-tag complete" v-if="record.status == 2"><CheckCircleFilled /> 学习完成</span>

              <a-tooltip placement="top" v-if="record.status == 3">
                <template #title>
                  <span>{{ record.errmsg }}</span>
                </template>
                <span>
                  <span class="status-tag status-error"><CloseCircleFilled /> 学习失败</span>
                  <a class="ml8" v-if="libraryInfo.type == 2" @click="handlePreview(record)">学习</a>
                </span>
              </a-tooltip>
              <template v-if="record.status == 4">
                <span class="status-tag"><ClockCircleFilled /> 待学习</span>
                <a class="ml8" @click="handlePreview(record)">学习</a>
              </template>
              <template v-if="record.status == 5">
                <span class="status-tag"><ClockCircleFilled /> 待获取</span>
              </template>
              <span class="status-tag running" v-if="record.status == 6"><a-spin size="small" /> 获取中</span>
              <a-tooltip placement="top" v-if="record.status == 7">
                <template #title>
                  <span>{{ record.errmsg }}</span>
                </template>
                <span class="status-tag error"><CloseCircleFilled /> 获取失败</span>
              </a-tooltip>
            </template>
            <template v-if="column.key === 'graph_status'">
              <!--0待生成 1排队中 2生成完成 3生成失败 4生成中 5部分成功-->
              <template v-if="record.graph_status == 0">
                <span class="status-tag"><ClockCircleFilled /> 待生成</span>
                <a class="ml8" @click="createGraphTask(record)">生成</a>
              </template>
              <span v-else-if="record.graph_status == 1" class="status-tag running"><HourglassFilled /> 排队中</span>
              <span v-else-if="record.graph_status == 2" class="status-tag complete"><CheckCircleFilled /> 生成完成</span>
              <template v-else-if="record.graph_status == 3">
                <span class="status-tag error"><CloseCircleFilled /> 生成失败</span>
                <a class="ml8" @click="createGraphTask(record)">生成</a>
                <a-tooltip v-if="record.graph_err_msg" :title="record.graph_err_msg">
                  <div class="zm-line1 reason-text">原因：{{record.graph_err_msg}}</div>
                </a-tooltip>
              </template>
              <span v-else-if="record.graph_status == 4" class="status-tag running"><a-spin size="small" /> 生成中</span>
              <template v-else-if="record.graph_status == 5">
                <span  class="status-tag warning"><CheckCircleFilled /> 部分成功</span>
                <div class="reason-text">
                  失败数：{{ record.graph_err_count || 0 }}
                  <a class="ml8" @click="handlePreview(record, {graph_status: 3})">详情</a>
                </div>
              </template>
            </template>
            <template v-if="column.key === 'file_size'">
              <span v-if="record.doc_type == 3">-</span>
              <span v-else>{{ record.file_size_str }}</span>
            </template>
            <template v-if="column.key === 'paragraph_count'">
              <span v-if="record.status == 0 || record.status == 1">-</span>
              <span v-else>{{ record.paragraph_count }}</span>
            </template>
            <template v-if="column.key === 'action'">
              <a-dropdown>
                <div class="table-btn" @click.prevent>
                  <MoreOutlined />
                </div>
                <template #overlay>
                  <a-menu>
                    <a-menu-item
                      :disabled="record.status == 6 || record.status == 7 || record.status == 0"
                    >
                      <div @click="handlePreview(record)">预览</div>
                    </a-menu-item>
                    <a-menu-item>
                      <a-popconfirm title="确定要删除吗?" @confirm="onDelete(record)">
                        <span>删除</span>
                      </a-popconfirm>
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </template>
          </template>
        </a-table>
      </div>
    </cu-scroll>
    <a-modal
      v-model:open="addFileState.open"
      :confirm-loading="addFileState.confirmLoading"
      :maskClosable="false"
      title="上传文档"
      @ok="handleSaveFiles"
      @cancel="handleCloseFileUploadModal"
      width="740px"
    >
      <a-form class="mt24" :label-col="{ span: 4 }" :wrapper-col="{ span: 20 }">
        <a-form-item required label="上传文档">
          <div class="upload-file-box">
            <UploadFilesInput
              :type="libraryInfo.type"
              v-model:value="addFileState.fileList"
              @change="onFilesChange"
            />
          </div>
        </a-form-item>
        <a-form-item v-if="existPdfFile" required label="PDF解析模式">
          <div class="select-card-box">
            <div
              v-for="item in PDF_PARSE_MODE"
              :key="item.key"
              :class="['select-card-item', {active: addFileState.pdf_parse_type == item.key}]"
              @click="pdfParseTypeChange(item)"
            >
              <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
              <div class="card-title">{{ item.title }}</div>
              <div class="card-desc">{{ item.desc }}</div>
            </div>
          </div>
        </a-form-item>
      </a-form>
    </a-modal>
    <QaUploadModal @ok="getData" ref="qaUploadModalRef" />
    <a-modal
      v-model:open="addUrlState.open"
      :confirm-loading="addUrlState.confirmLoading"
      :maskClosable="false"
      title="添加在线数据"
      width="746px"
      @ok="handleSaveUrl"
      @cancel="handleCloseUrlModal"
    >
      <a-form
        class="url-add-form"
        layout="vertical"
        ref="urlFormRef"
        :model="addUrlState"
        :rules="addUrlState.rules"
      >
        <a-form-item name="urls" label="网页链接">
          <a-textarea
            style="height: 120px"
            v-model:value="addUrlState.urls"
            placeholder="请输入网页链接,形式：一行标题一行网页链接"
          />
        </a-form-item>
        <a-form-item name="doc_auto_renew_frequency" label="更新频率" required>
          <a-select v-model:value="addUrlState.doc_auto_renew_frequency" style="width: 100%">
            <a-select-option :value="1">不自动更新</a-select-option>
            <a-select-option :value="2">每天</a-select-option>
            <a-select-option :value="3">每3天</a-select-option>
            <a-select-option :value="4">每7天</a-select-option>
            <a-select-option :value="5">每30天</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
    <AddCustomDocument
      :libraryInfo="libraryInfo"
      @ok="onSearch"
      ref="addCustomDocumentRef"
    ></AddCustomDocument>
    <OpenGrapgModal @ok="handleOpenGrapgOk" ref="openGrapgModalRef" />
  </div>
</template>

<script setup>
import { reactive, ref, toRaw, onUnmounted, onMounted, computed} from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useRoute, useRouter } from 'vue-router'
import {
  PlusOutlined,
  SearchOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  ClockCircleFilled,
  MoreOutlined,
  HourglassFilled,
  ExclamationCircleFilled
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import { getLibraryFileList, delLibraryFile, addLibraryFile, editLibrary, createGraph } from '@/api/library'
import { formatFileSize } from '@/utils/index'
import UploadFilesInput from '../add-library/components/upload-input.vue'
import { transformUrlData } from '@/utils/validate.js'
import AddCustomDocument from './components/add-custom-document.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import OpenGrapgModal from "@/views/library/library-details/components/open-grapg-modal.vue";
import QaUploadModal from './components/qa-upload-modal.vue'

const PDF_PARSE_MODE = [
  {key: 2, title: '图文OCR解析', desc: '通过OCR文字识别提取pdf文件内容，可以兼容扫描件，但是解析速度较慢。'},
  {key: 1, title: '纯文本解析', desc: '只提取pdf中的文字内容，如果文档为扫描件可能提取不到内容。'},
]
const rotue = useRoute()
const router = useRouter()
const query = rotue.query

const openGrapgModalRef = ref(null)
const createGraphSwitch = ref(false)
const libraryInfo = ref({
  library_intro: '',
  library_name: '',
  use_model: '',
  is_offline: null,
  type: 0
})

const modelForm = reactive({
  use_model: '',
  model_config_id: ''
})

const onChangeModel = (val, option) => {
  let new_use_model = option.modelName
  let new_model_config_id = option.modelId

  if (fileList.value.length > 0) {
    Modal.confirm({
      title: `确定切换模型为${new_use_model}吗？`,
      content: '切换后，所有学习文档将自动重新学习。',
      onOk() {
        modelForm.use_model = new_use_model
        modelForm.model_config_id = new_model_config_id
        saveLibraryConfig()
      },
      onCancel() {}
    })
  } else {
    modelForm.use_model = new_use_model
    modelForm.model_config_id = new_model_config_id
    saveLibraryConfig()
  }
}

const saveLibraryConfig = (showSuccessTip = true, callback=null) => {
  editLibrary({
    ...toRaw(libraryInfo.value),
    use_model: modelForm.use_model,
    model_config_id: modelForm.model_config_id
  }).then(() => {
    libraryInfo.value.use_model = modelForm.use_model
    libraryInfo.value.model_config_id = modelForm.model_config_id

    typeof callback === 'function' && callback()
    if (showSuccessTip) {
      message.success('保存成功')
    }
  })
}

const vectorModelList = ref([])

const onVectorModelLoaded = (list) => {
  vectorModelList.value = list

  setDefaultModel()
}

const setDefaultModel = () => {
  // 防止没有数据时，切换模型报错
  if (!libraryInfo.value.id) {
    setTimeout(() => {
      setDefaultModel()
    }, 100)

    return
  }

  if (vectorModelList.value.length > 0 && !libraryInfo.value.use_model) {
    // 遍历查找chatwiki模型
    let modelConfig = null
    let model = null

    for (let item of vectorModelList.value) {
      if (item.model_define === 'chatwiki') {
        modelConfig = item
        for (let child of modelConfig.children) {
          if (child.name === 'text-embedding-v2') {
            model = child
            break
          }
        }
        break
      }
    }

    if (!modelConfig) {
      modelConfig = vectorModelList.value[0]
      model = modelConfig.children[0]
    }

    if (modelConfig) {
      modelForm.use_model = model.name
      modelForm.model_config_id = model.model_config_id

      saveLibraryConfig(false)
    }
  }
}

const fileList = ref([])

const queryParams = reactive({
  library_id: query.id,
  file_name: undefined,
  page: 1,
  size: 10,
  total: 0
})

const columns = ref([])
const columnsDefault = [
  {
    title: '文档名称',
    dataIndex: 'file_name',
    key: 'file_name',
    width: 450
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 180
  },
  {
    title: '知识图谱',
    dataIndex: 'graph_status',
    key: 'graph_status',
    width: 200
  },
  {
    title: '文档格式',
    dataIndex: 'file_ext',
    key: 'file_ext',
    width: 100
  },
  {
    title: '文档大小',
    dataIndex: 'file_size_str',
    key: 'file_size',
    width: 100
  },
  // {
  //   title: '更新频率',
  //   dataIndex: 'doc_auto_renew_frequency',
  //   key: 'doc_auto_renew_frequency',
  //   width: '150px'
  // },
  {
    title: '分段',
    dataIndex: 'paragraph_count',
    key: 'paragraph_count',
    width: 120
  },
  {
    title: '知识图谱实体数',
    dataIndex: 'graph_entity_count',
    key: 'graph_entity_count',
    width: 160
  },
  {
    title: '更新时间',
    dataIndex: 'update_time',
    key: 'update_time',
    width: 150
  },
  {
    title: '操作',
    dataIndex: 'action',
    key: 'action',
    width: 100
  }
]

const onTableChange = (pagination) => {
  queryParams.page = pagination.current
  queryParams.size = pagination.pageSize
  getData()
}

const onSearch = () => {
  queryParams.page = 1
  getData()
}

let isLast = false

const onDelete = ({ id }) => {
  if (fileList.value.length == 1) {
    isLast = true
  }

  delLibraryFile({ id }).then(() => {
    if (isLast && queryParams.page > 1) {
      queryParams.page--
    }

    getData()
    message.success('删除成功')
  })
}

const handlePreview = (record, params={}) => {
  if (record.status == '4') {
    return router.push({
      path: '/library/document-segmentation',
      query: { document_id: record.id }
    })
  }
  if (record.status == '3' && libraryInfo.value.type != 2) {
    return message.error('学习失败,不可预览')
  }
  if (record.status == '0') {
    return message.error('转换中,稍候可预览')
  }
  if (record.status == '6') {
    return message.error('获取中,不可预览')
  }
  if (record.status == '7') {
    return message.error('获取失败,不可预览')
  }

  router.push({ name: 'libraryPreview', query: { id: record.id, ...params} })
}

const getData = () => {
  getLibraryFileList(toRaw(queryParams)).then((res) => {
    let info = res.data.info

    if (!modelForm.use_model && modelForm.use_model != info.use_model) {
      modelForm.use_model = info.use_model
      modelForm.model_config_id = info.model_config_id
    }

    libraryInfo.value = { ...info }
    createGraphSwitch.value = (info.graph_switch == 1)
    if (info.graph_switch == '0') {
      columns.value = columnsDefault.filter((item) => !['graph_status','graph_entity_count'].includes(item.key))
    } else {
      columns.value = columnsDefault
    }

    let list = res.data.list || []

    queryParams.total = res.data.total
    let needRefresh = false
    fileList.value = list.map((item) => {
      if (['1', '6', '0', '5', '4'].includes(item.status)) {
        needRefresh = true
      }
      item.file_size_str = formatFileSize(item.file_size)
      item.update_time = dayjs(item.update_time * 1000).format('YYYY-MM-DD HH:mm')
      return item
    })
    needRefresh && timingRefreshStatus()
    !needRefresh && clearInterval(timingRefreshStatusTimer.value)
  })
}

const timingRefreshStatusTimer = ref(null)
const timingRefreshStatus = () => {
  clearInterval(timingRefreshStatusTimer.value)
  timingRefreshStatusTimer.value = setInterval(() => {
    getData()
  }, 1000 * 5)
}

const addCustomDocumentRef = ref(null)

const handleMenuClick = (e) => {
  if (vectorModelList.value.length == 0) {
    Modal.confirm({
      title: `请先到模型管理中添加嵌入模型？`,
      content:
        '知识库学习需要使用到嵌入模型，请在系统管理-模型管理中添加。推荐使用通义千问、openai或者火山引擎的嵌入模型。',
      okText: '去添加',
      onOk() {
        router.push({ path: '/user/model' })
      }
    })
    return
  }

  let { key } = e

  if (key == 1) {
    handleOpenFileUploadModal()
  }
  if (key == 2) {
    handleOpenUrlModal()
  }
  if (key == 3) {
    addCustomDocumentRef.value.add()
  }
}
const addUrlState = reactive({
  open: false,
  urls: '',
  library_id: query.id,
  doc_auto_renew_frequency: 1,
  confirmLoading: false,
  rules: {
    urls: [
      {
        message: '请输入网页地址',
        required: true
      },
      {
        validator: (_rule, value) => {
          if (transformUrlData(value) === false) {
            return Promise.reject(new Error('网页地址不合法'))
          }
          return Promise.resolve()
        }
      }
    ]
  }
})

const handleOpenUrlModal = () => {
  addUrlState.open = true
  addUrlState.confirmLoading = false
  addUrlState.urls = ''
  addUrlState.doc_auto_renew_frequency = 1
}
const urlFormRef = ref(null)
const handleSaveUrl = () => {
  // 保存本地内容
  urlFormRef.value
    .validate()
    .then(() => {
      addUrlState.confirmLoading = true
      addLibraryFile({
        library_id: addUrlState.library_id,
        urls: JSON.stringify(transformUrlData(addUrlState.urls)),
        doc_auto_renew_frequency: addUrlState.doc_auto_renew_frequency,
        doc_type: 2
      }).then(() => {
        addUrlState.open = false
        addUrlState.confirmLoading = false
        onSearch()
      })
    })
    .catch(() => {
      addUrlState.confirmLoading = false
    })
}

const handleCloseUrlModal = () => {
  addUrlState.open = false
}

const addFileState = reactive({
  open: false,
  fileList: [],
  confirmLoading: false,
  pdf_parse_type: 1, //1纯文本解析，2ocr解析
})

const existPdfFile = computed(() => addFileState.fileList.filter(i => i.type === 'application/pdf').length > 0)

const qaUploadModalRef = ref(null)
const handleOpenFileUploadModal = () => {
  if(libraryInfo.value.type == 2){
    qaUploadModalRef.value.show()
    return
  }
  addFileState.fileList = []
  addFileState.open = true
}

const handleCloseFileUploadModal = () => {
  addFileState.fileList = []
}

const onFilesChange = (files) => {
  addFileState.fileList = files
}

const handleSaveFiles = () => {
  if (addFileState.fileList.length == 0) {
    message.error('请选择文件')
    return
  }

  addFileState.confirmLoading = true

  let formData = new FormData()

  formData.append('library_id', queryParams.library_id)
  let isTableType = false
  addFileState.fileList.forEach((file) => {
    if (file.name.includes('.xlsx') || file.name.includes('.csv')) {
      isTableType = true
    }
    formData.append('library_files', file)
  })
  if (existPdfFile.value) {
    formData.append('pdf_parse_type', addFileState.pdf_parse_type)
  }
  addLibraryFile(formData)
    .then((res) => {
      getData()
      addFileState.open = false
      addFileState.fileList = []
      addFileState.confirmLoading = false
      if (isTableType) {
        router.push('/library/document-segmentation?document_id=' + res.data.file_ids[0])
      }
    })
    .catch(() => {
      addFileState.confirmLoading = false
    })
}

const createGraphTask = record => {
  createGraph({id: record.id}).then(() => {
    getData()
  })
}

const createGraphSwitchChange = () => {
  if (createGraphSwitch.value) {
    createGraphSwitch.value = false
    let data = {
      graph_model_config_id: libraryInfo.value.graph_model_config_id,
      graph_use_model: libraryInfo.value.graph_use_model
    }
    if ((!data.graph_model_config_id || !data.graph_use_model) && vectorModelList.value.length > 0) {
      let modelConfig = vectorModelList.value[0]
      let model = modelConfig.children[0]
      data.graph_use_model = model.name
      data.graph_model_config_id = model.model_config_id
    }
    openGrapgModalRef.value.show(data)
  } else {
    libraryInfo.value.graph_switch = 0
    saveLibraryConfig(false, () => {
      getData()
    })
  }
}

const handleOpenGrapgOk = data => {
  createGraphSwitch.value = true
  libraryInfo.value.graph_switch = 1
  libraryInfo.value.graph_model_config_id = data.graph_model_config_id
  libraryInfo.value.graph_use_model = data.graph_use_model
  saveLibraryConfig(false, () => {
    getData()
  })
}

const pdfParseTypeChange = item => {
  addFileState.pdf_parse_type = item.key
}

onMounted(() => {
  getData()
})

onUnmounted(() => {
  timingRefreshStatusTimer.value && clearInterval(timingRefreshStatusTimer.value)
})
</script>

<style lang="less" scoped>
.details-library-page {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.doc-name-td {
  word-break: break-all;
}
.url-remark {
  color: #8c8c8c;
  margin-top: 2px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}
.padding-0 {
  padding: 0;
}
.test-menu-icon {
  color: #fff;
}
.library-name {
  height: 38px;
  line-height: 38px;
  padding: 0 16px;
  margin-bottom: 16px;
  font-size: 14px;
  font-weight: 600;
  color: #262626;
  border-radius: 2px;
  background-color: #f2f4f7;
  display: flex;
  align-items: center;
  .anticon-edit {
    margin-left: 8px;
    color: #8c8c8c;
    cursor: pointer;
  }
}
.between-content-box {
  display: flex;
  flex: 1;
  overflow: hidden;
  .left-menu-box {
    width: 232px;
    margin-right: 24px;
  }
  .right-content-box {
    flex: 1;
  }
}

.menu-item {
  width: 232px;
  height: 56px;
  display: flex;
  align-items: center;
  padding: 0 24px;
  background: #f5f5f5;
  border-radius: 2px;
  margin-bottom: 16px;
  cursor: pointer;
  &.active {
    background: #e6efff;
    border: 1px solid #2475fc;
  }
  .title {
    color: #242933;
    font-size: 14px;
    font-weight: 600;
    line-height: 22px;
    margin-left: 8px;
  }
}

.list-tools {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.list-tools {
  margin-bottom: 8px;
  .tools-items {
    display: flex;
    align-items: center;
    .tool-item {
      margin-right: 16px;
    }
  }
}

.list-content {
  .status-tag {
    display: inline-block;
    height: 24px;
    line-height: 24px;
    padding: 0 6px;
    border-radius: 2px;
    font-size: 14px;
    font-weight: 500;
    text-align: center;
    color: #3a4559;
    background-color: #edeff2;

    &.running {
      color: #2475fc;
      background-color: #e8effc;
    }
    &.complete {
      color: #21A665;
      background: #E8FCF3;
    }

    &.error {
      cursor: pointer;
      color: #fb363f;
      background-color: #f5c6c8;
    }
    &.warning {
      cursor: pointer;
      background: #faebe6;
      color: #ed744a;
    }
  }
}
//.upload-file-box {
//  padding: 30px 0;
//}
.ml8 {
  margin-left: 8px;
}
.url-add-form {
  margin-top: 24px;
}

.add-dropdown-btn.ant-dropdown {
  .ant-dropdown-menu {
    padding: 0;
    border-radius: 0;
    ::v-deep(.ant-dropdown-menu-item) {
      padding: 12px 16px;
    }
  }
}

.dropdown-btn-menu {
  .title-block {
    color: #262626;
    font-size: 14px;
    font-weight: 600;
    line-height: 22px;
  }
  .desc {
    color: #8c8c8c;
    font-size: 14px;
    line-height: 22px;
  }
}
.table-btn {
  cursor: pointer;
  &:hover {
    color: #2475fc;
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
.pd-5-8 {
  padding: 5px 8px;
}
.reason-text {
  color: #8C8C8C;
  font-size: 12px;
  line-height: 24px;
}

.select-card-box {
  display: flex;
  align-items: center;
  gap: 16px;
  .select-card-item {
    width: calc(50% - 8px);
    position: relative;
    padding: 16px;
    border-radius: 6px;
    border: 1px solid #d9d9d9;
    cursor: pointer;
    .check-arrow {
      position: absolute;
      display: block;
      right: -1px;
      bottom: -1px;
      width: 24px;
      height: 24px;
      font-size: 24px;
      color: #fff;
      opacity: 0;
      transition: all 0.2s cubic-bezier(0.8, 0, 1, 1);
    }
    .card-title {
      display: flex;
      align-items: center;
      line-height: 22px;
      margin-bottom: 4px;
      color: #262626;
      font-weight: 600;
      font-size: 14px;
    }
    .title-icon {
      margin-right: 4px;
      font-size: 16px;
    }
    .card-desc {
      min-height: 44px;
      line-height: 22px;
      font-size: 14px;
      color: #595959;
    }

    &.active {
      background: var(--01-, #e5efff);
      border: 2px solid #2475fc;
      .check-arrow {
        opacity: 1;
      }
      .card-title {
        color: #2475fc;
      }
    }
  }
}
.mt24 {
  margin-top: 24px;
}
</style>
