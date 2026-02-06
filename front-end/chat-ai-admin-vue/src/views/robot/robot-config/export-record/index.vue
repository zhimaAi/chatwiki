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
    <div class="page-title">{{ t('title_export_record') }}</div>
    <div class="list-wrapper">
      <div class="content-wrapper">
        <a-alert show-icon :message="t('msg_export_file_retention')"></a-alert>
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
                <ClockCircleFilled />{{ t('status_waiting_export') }}
              </div>
              <div v-if="record.status == 1" class="status-item status-1">
                <LoadingOutlined />{{ t('status_exporting') }}
              </div>
              <div v-if="record.status == 2" class="status-item status-2">
                <CheckCircleFilled />{{ t('status_export_success') }}
              </div>
              <div v-if="record.status == 3" class="status-item status-3">
                <CloseCircleFilled />{{ t('status_export_failed') }}
                <a-tooltip v-if="record.err_msg">
                  <template #title>{{ record.err_msg }}</template>
                  <QuestionCircleOutlined />
                </a-tooltip>
              </div>
            </template>

            <template v-if="column.key === 'action'">
              <a @click="handleDownload(record)" v-if="!record.is_over_7_days">{{ t('action_download') }}</a>
              <span v-else>--</span>
            </template>
          </template>
        </a-table>
      </div>
    </div>
  </div>
</template>

<script setup>
import DateSelect from './components/date.vue'
import { getExportTaskList, downloadExportFile } from '@/api/chat'
import {
  ClockCircleFilled,
  CheckCircleFilled,
  LoadingOutlined,
  CloseCircleFilled,
  QuestionCircleOutlined
} from '@ant-design/icons-vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/modules/user'
import { useI18n } from '@/hooks/web/useI18n'
import dayjs from 'dayjs'
import { reactive, ref } from 'vue'

const { t } = useI18n('views.robot.robot-config.export-record.index')
const userStore = useUserStore()
const query = useRoute().query

const tableData = ref([])

const requestParams = reactive({
  robot_id: query.id,
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
  // window.open(`/manage/downloadExportFile?id=${record.id}&robot_id=${query.id}`)
  let targetUrl = `/manage/downloadExportFile?id=${record.id}&robot_id=${query.id}&token=${userStore.getToken}`
  var aTag = document.createElement('a')
  aTag.href = targetUrl
  aTag.style.display = 'none'
  aTag.click()
  // document.body.removeChild(aTag)
}

const columns = [
  {
    title: t('title_export_time'),
    dataIndex: 'create_time_desc',
    key: 'create_time_desc'
  },
  {
    title: t('title_export_file_name'),
    dataIndex: 'file_name',
    key: 'file_name'
  },
  {
    title: t('title_source'),
    key: 'source_desc',
    dataIndex: 'source_desc'
  },
  {
    title: t('title_status'),
    key: 'status',
    dataIndex: 'status'
  },
  {
    title: t('title_action'),
    key: 'action'
  }
]
</script>
