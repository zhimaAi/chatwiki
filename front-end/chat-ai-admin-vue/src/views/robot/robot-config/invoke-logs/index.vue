<template>
  <div class="user-model-page">
    <div class="page-title">{{ t('title_invoke_logs') }}</div>
    <div class="list-wrapper">
      <div class="content-wrapper">
        <div class="search-block">
          <a-range-picker
            v-model:value="dates"
            :allowClear="false"
            @change="handleDateChange"
            style="width: 256px"
            :presets="dateRangePresets"
          />
          <a-input
            style="width: 256px"
            v-model:value="requestParams.openid"
            @change="onSearch"
            :placeholder="t('ph_input_openid')"
          >
            <template #suffix>
              <SearchOutlined />
            </template>
          </a-input>
          <a-input
            style="width: 256px"
            v-model:value="requestParams.question"
            @change="onSearch"
            :placeholder="t('ph_input_question')"
          >
            <template #suffix>
              <SearchOutlined />
            </template>
          </a-input>
        </div>

        <a-table
          class="table-list"
          :columns="columns"
          :data-source="tableData"
          :loading="loading"
          :scroll="{ x: tableScrollX }"
          sticky
          :pagination="{
            current: requestParams.page,
            total: requestParams.total,
            pageSize: requestParams.size,
            show_quick_jumper: true,
            show_size_changer: true,
            pageSizeOptions: ['10', '20', '50', '100']
          }"
          @change="onTableChange"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <a-space>
                <a-tag :color="getStatusInfo(record.status).color">
                  {{ getStatusInfo(record.status).text }}
                </a-tag>
                <a-button
                  v-if="record.status === WORKFLOW_STATUS_RUNNING"
                  type="link"
                  size="small"
                  :loading="stoppingLogIds.has(record.id)"
                  @click.stop="handleStop(record)"
                >
                  {{ t('btn_stop') }}
                </a-button>
              </a-space>
            </template>
            <template v-if="column.key === 'action'">
              <a @click="handleOpenDetailModal(record)">{{ t('btn_view_detail') }}</a>
            </template>
          </template>
        </a-table>
      </div>
      <LogsDetail ref="logsDetailRef" />
    </div>
  </div>
</template>

<script setup>
import { workflowLogs, stopWorkFlow } from '@/api/chat'
import { SearchOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { getDateRangePresets } from '@/utils/index'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'
import LogsDetail from './components/logs-detail.vue'
import { reactive, ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.invoke-logs.index')
const query = useRoute().query

const WORKFLOW_STATUS_RUNNING = 0
const WORKFLOW_STATUS_COMPLETED = 1
const WORKFLOW_STATUS_STOPPED = 2

const dateRangePresets = getDateRangePresets()

const tableData = ref([])

const dates = ref([dayjs().subtract(6, 'day'), dayjs()])

const loading = ref(false)
const stoppingLogIds = ref(new Set())
let statusPollingTimer = null
let fastPollingTimer = null

const requestParams = reactive({
  robot_id: query.id,
  page: 1,
  size: 20,
  total: 0,
  openid: '',
  question: '',
  start_date: dates.value[0].format('YYYY-MM-DD'),
  end_date: dates.value[1].format('YYYY-MM-DD')
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

const handleDateChange = () => {
  requestParams.start_date = ''
  requestParams.end_date = ''
  if (dates.value && dates.value.length > 0) {
    requestParams.start_date = dates.value[0].format('YYYY-MM-DD')
    requestParams.end_date = dates.value[1].format('YYYY-MM-DD')
  }
  onSearch()
}

const hasRunningRecord = () => {
  return tableData.value.some((item) => item.status === WORKFLOW_STATUS_RUNNING)
}

const clearStatusPolling = () => {
  if (statusPollingTimer) {
    clearInterval(statusPollingTimer)
    statusPollingTimer = null
  }
}

const refreshPollingStatus = () => {
  if (hasRunningRecord()) {
    if (!statusPollingTimer) {
      statusPollingTimer = setInterval(() => {
        getData(false)
      }, 5000)
    }
  } else {
    clearStatusPolling()
  }
}

const clearFastPolling = () => {
  if (fastPollingTimer) {
    clearInterval(fastPollingTimer)
    fastPollingTimer = null
  }
}

const startFastPollingAfterStop = () => {
  clearFastPolling()
  let times = 0
  fastPollingTimer = setInterval(() => {
    times += 1
    getData(false).finally(() => {
      if (times >= 5 || !hasRunningRecord()) {
        clearFastPolling()
      }
    })
  }, 1000)
}

const getStatusInfo = (status) => {
  const statusMap = {
    [WORKFLOW_STATUS_RUNNING]: { text: t('status_running'), color: 'blue' },
    [WORKFLOW_STATUS_COMPLETED]: { text: t('status_completed'), color: 'green' },
    [WORKFLOW_STATUS_STOPPED]: { text: t('status_stopped'), color: 'orange' }
  }
  return statusMap[Number(status)] || { text: t('status_unknown'), color: 'default' }
}

const getData = (showLoading = true) => {
  if (showLoading) {
    loading.value = true
  }
  return workflowLogs({
    ...requestParams
  })
    .then((res) => {
      let lists = res.data.list || []

      tableData.value = lists.map((item) => {
        item.duration_mills_desc = `${item.duration_mills}ms`
        item.total_token_desc = `${(item.total_token / 1000).toFixed(3)}`
        item.create_time_desc = dayjs(item.create_time * 1000).format('YYYY-MM-DD HH:mm:ss')
        item.node_logs = item.node_logs ? JSON.parse(item.node_logs) : []
        item.status = Number(item.status ?? WORKFLOW_STATUS_COMPLETED)
        return item
      })

      requestParams.total = +res.data.total || 0
      refreshPollingStatus()
    })
    .finally(() => {
      if (showLoading) {
        loading.value = false
      }
    })
}

const logsDetailRef = ref(null)
const handleOpenDetailModal = (record) => {
  logsDetailRef.value.show({ ...record })
}

const handleStop = async (record) => {
  stoppingLogIds.value.add(record.id)
  try {
    const res = await stopWorkFlow({ log_id: record.id })
    if (res?.data?.stopped) {
      record.status = WORKFLOW_STATUS_STOPPED
      message.success(t('msg_stop_success'))
    } else {
      message.warning(t('msg_stop_unavailable'))
    }
    await getData()
    startFastPollingAfterStop()
  } catch (err) {
    message.error(err?.message || t('msg_stop_failed'))
  } finally {
    stoppingLogIds.value.delete(record.id)
  }
}

const columns = [
  {
    title: t('label_openid'),
    dataIndex: 'openid',
    key: 'openid',
    width: 150
  },
  {
    title: t('label_question'),
    dataIndex: 'question',
    key: 'question',
    width: 250
  },
  {
    title: t('label_status'),
    key: 'status',
    dataIndex: 'status',
    width: 160
  },
  {
    title: t('label_workflow_version'),
    key: 'version',
    dataIndex: 'version',
    width: 130
  },
  {
    title: t('label_duration'),
    key: 'duration_mills_desc',
    dataIndex: 'duration_mills_desc',
    width: 130
  },
  {
    title: t('label_token_consumption'),
    key: 'total_token_desc',
    dataIndex: 'total_token_desc',
    width: 130
  },
  {
    title: t('label_invoke_time'),
    key: 'create_time_desc',
    dataIndex: 'create_time_desc',
    width: 150
  },
  {
    title: t('label_action'),
    key: 'action',
    width: 100
  }
]

const tableScrollX = columns.reduce((total, column) => total + Number(column.width || 0), 0)

onMounted(() => {
  onSearch()
})

onUnmounted(() => {
  clearStatusPolling()
  clearFastPolling()
})
</script>

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
    display: flex;
    align-items: center;
    gap: 16px;
  }
  .table-list {
    margin-top: 16px;
  }
}
</style>
