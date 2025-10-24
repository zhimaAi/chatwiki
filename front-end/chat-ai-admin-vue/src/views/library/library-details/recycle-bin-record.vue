<template>
  <div class="details-library-page">
    <cu-scroll :scrollbar="false">
      <a-alert show-icon message="回收站内的文档依然占用账户资源，请及时清理无用数据"></a-alert>
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
          :loading="isLoading"
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
                <span>{{ record.file_name }}</span>
                <a style="margin-left: 2px" @click="handleSyncDownload(record)">下载</a>
              </div>
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
              <a-flex :gap="12">
                <a @click="handleRestore(record)">恢复</a>
                <a @click="handleDel(record)">删除</a>
              </a-flex>
            </template>
          </template>
        </a-table>
      </div>
    </cu-scroll>
  </div>
</template>

<script setup>
import { reactive, ref, toRaw, onUnmounted, onMounted, computed } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useRoute, useRouter } from 'vue-router'
import { SearchOutlined, LoadingOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import {
  getLibFileRecycleList,
  delRecycleLibraryFile,
  restoreRecycleLibraryFile
} from '@/api/library'
import { formatFileSize } from '@/utils/index'
import { useUserStore } from '@/stores/modules/user'
const userStore = useUserStore()

const rotue = useRoute()
const router = useRouter()
const query = rotue.query

const isLoading = ref(false)

const state = reactive({
  selectedRowKeys: []
})

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
    title: '删除时间',
    dataIndex: 'delete_time',
    key: 'delete_time',
    width: 150
  },
  {
    title: '操作',
    dataIndex: 'action',
    key: 'action',
    fixed: 'right',
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

const getData = () => {
  let params = toRaw(queryParams)
  isLoading.value = true
  getLibFileRecycleList(params)
    .then((res) => {
      let list = res.data.list || []
      list.forEach((item) => {
        item.file_size_str = formatFileSize(item.file_size)
        item.delete_time =
          item.delete_time > 0 ? dayjs(item.delete_time * 1000).format('YYYY-MM-DD HH:mm') : '--'
      })
      fileList.value = list
      queryParams.total = +res.data.total || 0
    })
    .finally(() => {
      isLoading.value = false
    })
}

const handleSyncDownload = (record) => {
  let targetUrl = `/manage/downloadLibraryFile?id=${record.id}&token=${userStore.getToken}`
  var aTag = document.createElement('a')
  aTag.href = targetUrl
  aTag.style.display = 'none'
  aTag.click()
}

const handleRestore = (record) => {
  Modal.confirm({
    title: `恢复确认？`,
    content: `确认恢复该文档【${record.file_name}】`,
    okText: '确定',
    onOk() {
      restoreRecycleLibraryFile({
        id: record.id
      }).then((res) => {
        message.success('恢复成功')
        getData()
      })
    }
  })
}

const handleDel = (record) => {
  Modal.confirm({
    title: `彻底删除确认？`,
    content: '文件删除后无法恢复，请确认是否删除?',
    okText: '确定',
    onOk() {
      delRecycleLibraryFile({
        id: record.id
      }).then((res) => {
        message.success('删除成功')
        getData()
      })
    }
  })
}

onMounted(() => {
  if (query.page) {
    queryParams.page = +query.page
  }
  getData()
})

onUnmounted(() => {})
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
