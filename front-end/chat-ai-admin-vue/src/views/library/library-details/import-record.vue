<template>
  <div class="details-library-page">
    <cu-scroll :scrollbar="false">
      <a-alert
        show-icon
        message="通过文档导入记录，可以查看文档学习状态，学习失败的文档可以重新学习"
      ></a-alert>
      <div class="list-tools">
        <div class="tools-items">
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
      <div class="list-content">
        <a-table
          :columns="columns"
          :data-source="fileList"
          :scroll="{ x: 1000 }"
          row-key="id"
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
                <a-popover :title="null" v-if="record.doc_type == 2">
                  <template #content>
                    原链接：<a :href="record.doc_url" target="_blank">{{ record.doc_url }} </a>
                    <CopyOutlined
                      v-copy="`${record.doc_url}`"
                      style="margin-left: 4px; cursor: pointer"
                    />
                  </template>
                  <span>
                    <span v-if="['5', '6', '7'].includes(record.status)">{{ record.doc_url }}</span>
                    <span v-else>{{ record.file_name }}</span>
                  </span>
                </a-popover>
                <span v-else>
                  <span v-if="['5', '6', '7'].includes(record.status)">{{ record.doc_url }}</span>
                  <span v-else>{{ record.file_name }}</span>
                  <a style="margin-left: 2px;" @click="handleSyncDownload(record)">下载</a>
                </span>
                <div v-if="record.doc_type == 2 && record.remark" class="url-remark">
                  备注：{{ record.remark }}
                </div>
              </div>
            </template>
            <template v-if="column.key === 'status'">
              <template v-if="record.file_ext == 'pdf' && record.pdf_parse_type >= 2">
                <div class="pdf-progress-box" v-if="record.status == 0">
                  <div class="progress-title">
                    <span class="status-box"><LoadingOutlined />文档解析中</span>
                    <a @click="handleCancelOcrPdf(record)">取消</a>
                  </div>
                  <div class="progress-bar">
                    <a-progress
                      size="small"
                      class="progress-bar-box"
                      :percent="parseInt((record.ocr_pdf_index / record.ocr_pdf_total) * 100)"
                      :show-info="false"
                    />
                    <div class="num-box">
                      {{ record.ocr_pdf_index }} / {{ record.ocr_pdf_total }}
                    </div>
                  </div>
                </div>
              </template>
              <template v-else>
                <span class="status-tag running" v-if="record.status == 0"
                  ><a-spin size="small" /> 转换中</span
                >
              </template>

              <span class="status-tag running" v-if="record.status == 1"
                ><a-spin size="small" /> 学习中</span
              >

              <span class="status-tag complete" v-if="record.status == 2"
                ><CheckCircleFilled /> 学习完成</span
              >

              <a-tooltip placement="top" v-if="record.status == 3">
                <template #title>
                  <span>{{ record.errmsg }}</span>
                </template>
                <span>
                  <span class="status-tag status-error"><CloseCircleFilled /> 转换失败</span>
                  <a class="ml8" v-if="libraryInfo.type == 2" @click="handlePreview(record)"
                    >学习</a
                  >
                </span>
              </a-tooltip>
              <a-tooltip placement="top" v-if="record.status == 8">
                <template #title>
                  <span>{{ record.errmsg }}</span>
                </template>
                <span>
                  <span class="status-tag status-error"><CloseCircleFilled /> 转化异常</span>
                </span>
              </a-tooltip>
              <template v-if="record.status == 4">
                <span class="status-tag"><ClockCircleFilled /> 待学习</span>
                <a class="ml8" @click="handlePreview(record)">学习</a>
              </template>
              <template v-if="record.status == 5">
                <span class="status-tag"><ClockCircleFilled /> 待获取</span>
              </template>
              <span class="status-tag running" v-if="record.status == 6"
                ><a-spin size="small" /> 获取中</span
              >
              <a-tooltip placement="top" v-if="record.status == 7">
                <template #title>
                  <span>{{ record.errmsg }}</span>
                </template>
                <span class="status-tag error"><CloseCircleFilled /> 获取失败</span>
              </a-tooltip>
              <template v-if="record.status == 9">
                <span class="status-tag cancel"><ExclamationCircleOutlined /> 取消解析</span>
              </template>

              <span class="status-tag subning" v-if="record.status == 10">
                <a-spin size="small" />
                正在分段
              </span>
            </template>
            <template v-if="column.key === 'file_size'">
              <span v-if="record.doc_type == 3">-</span>
              <span v-else>{{ record.file_size_str }}</span>
            </template>
            <template v-if="column.key === 'paragraph_count'">
              <span v-if="record.status == 0 || record.status == 1">-</span>
              <span v-else>{{ record.paragraph_count }}</span>
            </template>
          </template>
        </a-table>
      </div>
    </cu-scroll>
    <a-modal v-model:open="downLoadModalOpen" :title="null" :footer="null" :width="640">
      <a-result
        status="success"
        title="导出任务创建成功"
        sub-title="系统会在后台导出。导出数据量越大，耗时越久。您可以稍后点击导出记录查看并下载导出的文件。"
      >
        <template #extra>
          <a-button style="margin-right: 16px;" @click="downLoadModalOpen = false">知道了</a-button>
          <a-button @click="toDownloadPage" type="primary">去下载</a-button>
        </template>
      </a-result>
    </a-modal>
  </div>
</template>

<script setup>
import { useStorage } from '@/hooks/web/useStorage'
import { reactive, ref, toRaw, onUnmounted, onMounted, computed } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useRoute, useRouter } from 'vue-router'
import {
  SearchOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  ClockCircleFilled,
  ExclamationCircleOutlined,
  CopyOutlined,
  LoadingOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import { getLibraryFileList, cancelOcrPdf, createExportLibFileTask } from '@/api/library'
import { formatFileSize } from '@/utils/index'
import { useLibraryStore } from '@/stores/modules/library'
import { useUserStore } from '@/stores/modules/user'
const { setStorage } = useStorage('localStorage')
const userStore = useUserStore()
const libraryStore = useLibraryStore()
const { changeGraphSwitch } = libraryStore

const rotue = useRoute()
const router = useRouter()
const query = rotue.query

const isLoading = ref(false)
const libraryInfo = ref({
  library_intro: '',
  library_name: '',
  use_model: '',
  is_offline: null,
  type: 0
})
const state = reactive({
  selectedRowKeys: []
})

const modelForm = reactive({
  use_model: '',
  model_config_id: ''
})

const fileList = ref([])

const count_data = reactive({
  learned_err_count: 0, // 失败数
  learned_count: 0, // 成功数
  learned_wait_count: 0 // 待学习
})

const queryParams = reactive({
  library_id: query.id,
  file_name: undefined,
  status: '', // 2:成功 3:全部失败 8:部分失败 4:待学习
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
    width: 300
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
  {
    title: '分段',
    dataIndex: 'paragraph_count',
    key: 'paragraph_count',
    width: 120
  },
  {
    title: '文档状态',
    dataIndex: 'status',
    key: 'status',
    width: 200
  },
  {
    title: '操作时间',
    dataIndex: 'update_time',
    key: 'update_time',
    width: 150
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

const handlePreview = (record, params = {}) => {
  if (record.status == '4' || record.status == '3' ) {
    return router.push({
      path: '/library/document-segmentation',
      query: { document_id: record.id, page: queryParams.page }
    })
  }
  return
  if (record.status == '3' && libraryInfo.value.type != 2) {
    return message.error('学习失败,不可预览')
  }
  if (record.status == '0') {
    return message.error('转换中,稍候可预览')
  }
  if (record.status == '1') {
    return message.error('学习中,不可预览')
  }
  if (record.status == '6') {
    return message.error('获取中,不可预览')
  }
  if (record.status == '7') {
    return message.error('获取失败,不可预览')
  }
  if (record.status == '10') {
    return message.error('正在分段,不可预览')
  }

  router.push({ name: 'libraryPreview', query: { id: record.id, ...params } })
}

const getData = () => {
  let params = toRaw(queryParams)
  if (params.status == 0) {
    params.status = ''
  }
  isLoading.value = true
  getLibraryFileList(params)
    .then((res) => {
      let info = res.data.info

      if (!modelForm.use_model && modelForm.use_model != info.use_model) {
        modelForm.use_model = info.use_model
        modelForm.model_config_id = info.model_config_id
      }

      libraryInfo.value = { ...info }

      columns.value = columnsDefault

      let list = res.data.list || []
      let countData = res.data.count_data || {}

      queryParams.total = res.data.total

      count_data.learned_count = countData.learned_count
      count_data.learned_err_count = countData.learned_err_count
      count_data.learned_wait_count = countData.learned_wait_count

      let needRefresh = false
      fileList.value = list.map((item) => {
        // , '4' 是待学习，如果加进去会一直刷新状态不会改变
        if (['1', '6', '0', '5'].includes(item.status)) {
          needRefresh = true
        }
        item.file_size_str = formatFileSize(item.file_size)
        item.update_time = dayjs(item.update_time * 1000).format('YYYY-MM-DD HH:mm')
        item.doc_last_renew_time_desc =
          item.doc_last_renew_time > 0
            ? dayjs(item.doc_last_renew_time * 1000).format('YYYY-MM-DD HH:mm')
            : '--'
        return item
      })
      needRefresh && timingRefreshStatus()
      !needRefresh && clearInterval(timingRefreshStatusTimer.value)
    })
    .finally(() => {
      isLoading.value = false
    })
}

const timingRefreshStatusTimer = ref(null)
const timingRefreshStatus = () => {
  clearInterval(timingRefreshStatusTimer.value)
  timingRefreshStatusTimer.value = setInterval(() => {
    getData()
  }, 1000 * 5)
}

const handleCancelOcrPdf = (record) => {
  Modal.confirm({
    title: `取消确认？`,
    content: '确认取消该文档解析',
    okText: '确定',
    onOk() {
      cancelOcrPdf({
        id: record.id
      }).then((res) => {
        message.success('取消成功')
        getData()
      })
    }
  })
}

const downLoadModalOpen = ref(false)
const handleSyncDownload = (record) => {
  let targetUrl = `/manage/downloadLibraryFile?id=${record.id}&token=${userStore.getToken}`
  var aTag = document.createElement('a')
  aTag.href = targetUrl
  aTag.style.display = 'none'
  aTag.click()
}

const toDownloadPage = ()=>{
  router.push({
    path: '/library/details/export-record',
    query,
  })
}

onMounted(() => {
  if (query.page) {
    queryParams.page = +query.page
  }
  getData()
})

onUnmounted(() => {
  timingRefreshStatusTimer.value && clearInterval(timingRefreshStatusTimer.value)
})
</script>

<style lang="less" scoped>
.details-library-page {
  height: 100%;
  padding-top: 24px;
  padding-left: 24px;
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
    overflow: hidden;
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
  margin-top: 16px;
}

.list-tools {
  margin-bottom: 16px;
  .tools-items {
    display: flex;
    align-items: center;
    .tool-item {
      margin-right: 16px;
    }
  }
}

.list-content {
  .text-block {
    color: #595959;
  }
  .c-gray {
    color: #8c8c8c;
  }
  .time-content-box {
    display: flex;
    color: #8c8c8c;
  }
  .btn-hover-block {
    height: 24px;
    display: flex;
    align-items: center;
    padding: 0 8px;
    cursor: pointer;
    width: fit-content;
    border-radius: 6px;
    transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
    &:hover {
      background: #e4e6eb;
    }
  }
  .status-tag {
    display: inline-block;
    height: 24px;
    line-height: 24px;
    padding: 0 6px;
    border-radius: 2px;
    font-size: 14px;
    font-weight: 500;
    text-align: center;
    color: #595959;
    background-color: #edeff2;

    &.running {
      color: #2475fc;
      background-color: #e8effc;
    }

    &.subning {
      color: #6524fc;
      background-color: #eae0ff;
    }

    &.complete {
      color: #21a665;
      background: #e8fcf3;
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

    &.status-learning {
      color: #2475fc;
      // background-color: #e8effc;
    }

    &.status-complete {
      color: #3a4559;
      background-color: #edeff2;
    }

    &.status-error {
      cursor: pointer;
      color: #fb363f;
      // background-color: #f5c6c8;
    }
    &.warning {
      cursor: pointer;
      // background: #faebe6;
      color: #ed744a;
    }
    &.cancel {
      background: #fae4dc;
      color: #ed744a;
    }
  }
}
.pdf-progress-box {
  .progress-title {
    display: flex;
    align-items: center;
    gap: 8px;
    white-space: nowrap;
    .status-box {
      width: fit-content;
      padding: 0 6px;
      display: flex;
      align-items: center;
      gap: 2px;
      height: 22px;
      border-radius: 6px;
      background: #e8effc;
      color: #2475fc;
      font-weight: 500;
    }
  }
  .progress-bar {
    display: flex;
    align-items: center;
    gap: 8px;
    line-height: 20px;
    .ant-progress-line {
      margin: 0;
    }
    .progress-bar-box {
      flex: 1;
    }
    .num-box {
      font-size: 12px;
      color: #8c8c8c;
    }
  }
}
//.upload-file-box {
//  padding: 30px 0;
//}
.ml8 {
  margin-left: 8px;
}
.ml4 {
  margin-left: 4px;
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
  color: #8c8c8c;
  font-size: 12px;
  line-height: 24px;
}

.select-card-box {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
  margin-top: 8px;
  .select-card-item {
    min-height: 105px;
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
      gap: 4px;
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
      line-height: 22px;
      font-size: 14px;
      color: #595959;
    }

    .card-switch {
      display: flex;
      gap: 4px;
      color: #8c8c8c;
      font-size: 14px;
      line-height: 22px;

      .card-switch-btn {
        cursor: pointer;
        color: #2475fc;

        &:hover {
          opacity: 0.8;
        }
      }
    }

    .card-switch-box {
      width: 52px;
      height: 22px;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 10px;
      border-radius: 6px;
      border: 1px solid #edb8a6;
      background: #ffece5;
      color: #ed744a;
      font-size: 12px;
      font-style: normal;
      font-weight: 400;
      line-height: 20px;
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

.right-box {
  display: flex;
  gap: 16px;
}

.status-box {
  display: flex;
  align-items: center;

  .status-item {
    display: flex;
    align-items: center;

    .status-label {
      color: #595959;
    }

    .content-tip {
      color: red;
    }
  }
}

.select-card-main {
  :deep(.ant-row) {
    display: block;
  }
}

.subning {
  :deep(.ant-spin-dot-item) {
    background-color: #6524fc !important;
  }
}
</style>
