<style lang="less" scoped>
.user-model-page {
  width: 100%;
  height: 100%;
  background-color: #f2f4f7;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  .page-title {
    display: flex;
    align-items: center;
    gap: 24px;
    padding: 24px 24px 16px;
    background-color: #fff;
    color: #000000;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }

  .list-wrapper {
    background: #fff;
    flex: 1;
    overflow-x: hidden;
    overflow-y: auto;
  }
  .content-wrapper {
    padding: 0 24px 16px 24px;
  }
  .search-block {
    margin-top: 24px;
  }
  .table-list {
    margin-top: 16px;
  }
  .status-item {
    display: flex;
    height: 22px;
    align-items: center;
    width: fit-content;
    border-radius: 6px;
    gap: 4px;
    font-size: 14px;
    line-height: 22px;
    padding: 0 6px;
    font-weight: 500;
    white-space: nowrap;
    .anticon {
      font-size: 15px;
    }
    &.status-0 {
      background: #edeff2;
      color: #3a4559;
    }
    &.status-1 {
      background: #e8effc;
      color: #2475fc;
    }
    &.status-2 {
      background: #e8fcf3;
      color: #21a665;
    }
    &.status-3 {
      color: #fb363f;
      background-color: #f5c6c8;
    }
  }
}
</style>

<template>
  <div class="user-model-page">
    <div class="page-title">导出记录</div>
    <div class="list-wrapper">
      <div class="content-wrapper">
        <a-alert show-icon message="导出文件仅保留7天，7天后自动删除，请及时下载到本地"></a-alert>
        <div class="search-block">
          <DateSelect datekey="2" @dateChange="onDateChange" />
        </div>

        <a-table
          class="table-list"
          :columns="columns"
          :data-source="tableData"
          :pagination="{
            current: requestParams.page,
            total: requestParams.total,
            pageSize: requestParams.size,
            showQuickJumper: true,
            showSizeChanger: true,
            pageSizeOptions: ['10', '20', '50', '100']
          }"
          @change="onTableChange"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <div v-if="record.status == 0" class="status-item status-0">
                <ClockCircleFilled />等待导出
              </div>
              <div v-if="record.status == 1" class="status-item status-1">
                <LoadingOutlined />导出中
              </div>
              <div v-if="record.status == 2" class="status-item status-2">
                <CheckCircleFilled />导出成功
              </div>
              <div v-if="record.status == 3" class="status-item status-3">
                <CloseCircleFilled />导出失败
                <a-tooltip v-if="record.err_msg">
                  <template #title>{{ record.err_msg }}</template>
                  <QuestionCircleOutlined />
                </a-tooltip>
              </div>
            </template>

            <template v-if="column.key === 'action'">
              <a @click="handleDownload(record)" v-if="!record.is_over_7_days">下载</a>
              <span v-else>--</span>
            </template>
          </template>
        </a-table>
      </div>
    </div>
  </div>
</template>

<script setup>
import DateSelect from '@/views/robot/robot-config/export-record/components/date.vue'
import { getExportTaskList } from '@/api/chat'
import {
  ClockCircleFilled,
  CheckCircleFilled,
  LoadingOutlined,
  CloseCircleFilled,
  QuestionCircleOutlined
} from '@ant-design/icons-vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/modules/user'
import dayjs from 'dayjs'
import { reactive, ref } from 'vue'
const userStore = useUserStore()
const query = useRoute().query

const tableData = ref([])

const requestParams = reactive({
  library_id: query.id,
  source: 2,
  page: 1,
  size: 10,
  total: 0,
  start_time: '',
  end_time: ''
})

const onTableChange = (pagination) => {
  requestParams.page = pagination.current
  requestParams.size = pagination.pageSize
  getData()
}

const onSearch = () => {
  requestParams.page = 1
  getData()
}

const onDateChange = (data) => {
  requestParams.start_time = data.start_time
  requestParams.end_time = data.end_time
  onSearch()
}

const getData = () => {
  getExportTaskList({
    ...requestParams
  }).then((res) => {
    let lists = res.data.list
    let maps = res.data.map
    lists.forEach((item) => {
      item.create_time_desc = dayjs(item.create_time * 1000).format('YYYY/MM/DD HH:mm')
      item.source_desc = maps.filter((it) => it.source == item.source)[0]?.source_name
      item.is_over_7_days = item.create_time < Date.now() / 1000 - 7 * 24 * 60 * 60
    })
    tableData.value = lists
    requestParams.total = +res.data.total
  })
}

const handleDownload = (record) => {
  let targetUrl = `/manage/downloadExportFile?id=${record.id}&library_id=${query.id}&token=${userStore.getToken}`
  var aTag = document.createElement('a')
  aTag.href = targetUrl
  aTag.style.display = 'none'
  aTag.click()
}

const columns = [
  {
    title: '导出时间',
    dataIndex: 'create_time_desc',
    key: 'create_time_desc'
  },
  {
    title: '导出文件名称',
    dataIndex: 'file_name',
    key: 'file_name'
  },
  {
    title: '来源',
    key: 'source_desc',
    dataIndex: 'source_desc'
  },
  {
    title: '状态',
    key: 'status',
    dataIndex: 'status'
  },
  {
    title: '操作',
    key: 'action'
  }
]
</script>
