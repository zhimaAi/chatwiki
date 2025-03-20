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
                      <div class="desc">上传本地 text/pdf/doc 等格式文件</div>
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

          <div class="tool-item">
            <span>嵌入模型：</span>
            <model-select
              modelType="TEXT EMBEDDING"
              :isOffline="false"
              :modeName="modelForm.use_model"
              :modeId="modelForm.model_config_id"
              style="width: 300px"
              @change="onChangeModel"
              @loaded="onVectorModelLoaded"
            />
          </div>
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
              <span class="status-tag status-queuing" v-if="record.status == 0"
                ><a-spin size="small" /> 转换中</span
              >
              <span class="status-tag status-learning" v-if="record.status == 1"
                ><a-spin size="small" /> 学习中</span
              >
              <span class="status-tag status-complete" v-if="record.status == 2"
                ><CheckCircleFilled /> 学习完成</span
              >

              <a-tooltip placement="top" v-if="record.status == 3">
                <template #title>
                  <span>{{ record.errmsg }}</span>
                </template>
                <span class="status-tag status-error"><CloseCircleFilled /> 学习失败</span>
              </a-tooltip>
              <template v-if="record.status == 4">
                <span class="status-tag status-complete"><ClockCircleFilled /> 待学习</span>
                <a class="ml8" @click="handlePreview(record)">学习</a>
              </template>
              <template v-if="record.status == 5">
                <span class="status-tag status-complete"><ClockCircleFilled /> 待获取</span>
              </template>
              <span class="status-tag status-learning" v-if="record.status == 6"
                ><a-spin size="small" /> 获取中</span
              >
              <a-tooltip placement="top" v-if="record.status == 7">
                <template #title>
                  <span>{{ record.errmsg }}</span>
                </template>
                <span class="status-tag status-error"><CloseCircleFilled /> 获取失败</span>
              </a-tooltip>
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
    >
      <div class="upload-file-box">
        <UploadFilesInput
          :type="libraryInfo.type"
          v-model:value="addFileState.fileList"
          @change="onFilesChange"
        />
      </div>
    </a-modal>
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
  </div>
</template>

<script setup>
import { reactive, ref, toRaw, onUnmounted, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useRoute, useRouter } from 'vue-router'
import {
  PlusOutlined,
  SearchOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  ClockCircleFilled,
  MoreOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import { getLibraryFileList, delLibraryFile, addLibraryFile, editLibrary } from '@/api/library'
import { formatFileSize } from '@/utils/index'
import UploadFilesInput from '../add-library/components/upload-input.vue'
import { transformUrlData } from '@/utils/validate.js'
import AddCustomDocument from './components/add-custom-document.vue'
import ModelSelect from '@/components/model-select/model-select.vue'

const rotue = useRoute()
const router = useRouter()
const query = rotue.query

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

const saveLibraryConfig = (showSuccessTip = true) => {
  editLibrary({
    ...toRaw(libraryInfo.value),
    use_model: modelForm.use_model,
    model_config_id: modelForm.model_config_id
  }).then(() => {
    libraryInfo.value.use_model = modelForm.use_model
    libraryInfo.value.model_config_id = modelForm.model_config_id

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

const columns = [
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
    width: '160px'
  },
  {
    title: '文档格式',
    dataIndex: 'file_ext',
    key: 'file_ext',
    width: '100px'
  },
  {
    title: '文档大小',
    dataIndex: 'file_size_str',
    key: 'file_size',
    width: '100px'
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
    width: '120px'
  },
  {
    title: '更新时间',
    dataIndex: 'update_time',
    key: 'update_time',
    width: '150px'
  },
  {
    title: '操作',
    dataIndex: 'action',
    key: 'action',
    width: '60px'
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

const handlePreview = (record) => {
  if (record.status == '4') {
    return router.push({
      path: '/library/document-segmentation',
      query: { document_id: record.id }
    })
  }
  if (record.status == '3') {
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

  router.push({ name: 'libraryPreview', query: { id: record.id } })
}

const getData = () => {
  getLibraryFileList(toRaw(queryParams)).then((res) => {
    let info = res.data.info

    if (!modelForm.use_model && modelForm.use_model != info.use_model) {
      modelForm.use_model = info.use_model
      modelForm.model_config_id = info.model_config_id
    }

    libraryInfo.value = { ...info }

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
  confirmLoading: false
})

const handleOpenFileUploadModal = () => {
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

    &.status-queuing {
      color: #2475fc;
      background-color: #e8effc;
    }

    &.status-learning {
      color: #2475fc;
      background-color: #e8effc;
    }

    &.status-complete {
      color: #3a4559;
      background-color: #edeff2;
    }

    &.status-error {
      cursor: pointer;
      color: #fb363f;
      background-color: #f5c6c8;
    }
    &.status-split {
      cursor: pointer;
      background: #faebe6;
      color: #ed744a;
    }
  }
}
.upload-file-box {
  padding: 30px 0;
}
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
</style>
