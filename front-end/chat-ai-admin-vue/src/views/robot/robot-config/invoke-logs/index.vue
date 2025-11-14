<template>
  <div class="user-model-page">
    <div class="page-title">调用日志</div>
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
            placeholder="请输入openid"
          >
            <template #suffix>
              <SearchOutlined />
            </template>
          </a-input>
          <a-input
            style="width: 256px"
            v-model:value="requestParams.question"
            @change="onSearch"
            placeholder="请输入问题"
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
          sticky
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
            <template v-if="column.key === 'action'">
              <a @click="handleOpenDetailModal(record)">查看明细</a>
            </template>
          </template>
        </a-table>
      </div>
      <LogsDetail ref="logsDetailRef" />
    </div>
  </div>
</template>

<script setup>
import { workflowLogs } from '@/api/chat'
import { SearchOutlined } from '@ant-design/icons-vue'
import { getDateRangePresets } from '@/utils/index'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'
import LogsDetail from './components/logs-detail.vue'
import { reactive, ref, onMounted } from 'vue'
const query = useRoute().query

const dateRangePresets = getDateRangePresets()

const tableData = ref([])

const dates = ref([dayjs().subtract(6, 'day'), dayjs()])

const loading = ref(false)

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

const handleDateChange = (data) => {
  requestParams.start_date = ''
  requestParams.end_date = ''
  if (dates.value && dates.value.length > 0) {
    requestParams.start_date = dates.value[0].format('YYYY-MM-DD')
    requestParams.end_date = dates.value[1].format('YYYY-MM-DD')
  }
  onSearch()
}

const getData = () => {
  loading.value = true
  workflowLogs({
    ...requestParams
  })
    .then((res) => {
      let lists = res.data.list || []

      tableData.value = lists.map((item) => {
        item.duration_mills_desc = `${item.duration_mills}ms`
        item.total_token_desc = `${(item.total_token / 1000).toFixed(3)}`
        item.create_time_desc = dayjs(item.create_time * 1000).format('YYYY-MM-DD HH:mm:ss')
        item.node_logs = item.node_logs ? JSON.parse(item.node_logs) : []
        return item
      })

      requestParams.total = +res.data.total || 0
    })
    .finally(() => {
      loading.value = false
    })
}

const logsDetailRef = ref(null)
const handleOpenDetailModal = (record) => {
  logsDetailRef.value.show({ ...record })
}

const columns = [
  {
    title: 'openid',
    dataIndex: 'openid',
    key: 'openid',
    width: 150
  },
  {
    title: 'question',
    dataIndex: 'question',
    key: 'question',
    width: 250
  },
  {
    title: '工作流版本',
    key: 'version',
    dataIndex: 'version',
    width: 130
  },
  {
    title: '耗时',
    key: 'duration_mills_desc',
    dataIndex: 'duration_mills_desc',
    width: 130
  },
  {
    title: 'token消耗（K）',
    key: 'total_token_desc',
    dataIndex: 'total_token_desc',
    width: 130
  },
  {
    title: '调用时间',
    key: 'create_time_desc',
    dataIndex: 'create_time_desc',
    width: 150
  },
  {
    title: '操作',
    key: 'action',
    width: 100
  }
]

onMounted(() => {
  onSearch()
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
