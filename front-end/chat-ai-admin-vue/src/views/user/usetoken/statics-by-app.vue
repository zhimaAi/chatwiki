<template>
  <div class="statics-content">
    <div class="search-box">
      <a-range-picker
        v-model:value="dates"
        :allowClear="false"
        @change="handleDateChange"
        style="width: 256px"
        :presets="dateRangePresets"
      />
      <a-cascader
        style="width: 256px"
        changeOnSelect
        @change="handleCascaderChange"
        v-model:value="cascaderValue"
        :options="options"
        placeholder="全部应用"
      />
    </div>
    <div class="list-box">
      <div class="main-title">使用趋势（单位：千）</div>
      <LineChart :options="lineChartData" />
      <div class="main-title">
        明细数据（单位：千）
        <div style="margin-left: auto"><a-button @click="handleExport">导出</a-button></div>
      </div>
      <a-table
        style="margin-top: 16px"
        :data-source="tableData"
        :pagination="{
          current: pager.page,
          total: pager.total,
          pageSize: pager.size,
          showQuickJumper: true,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100']
        }"
        @change="onTableChange"
      >
        <a-table-column title="类型" data-index="token_app_type_desc" :width="100">
          <template #default="{ record }">{{ record.token_app_type_desc }}</template>
        </a-table-column>
        <a-table-column title="应用名称" data-index="robot_name" :width="140"> </a-table-column>
        <a-table-column title="合计消耗(k)" data-index="total_token_desc" :width="140">
          <template #default="{ record }">{{ record.total_token_desc }}</template>
        </a-table-column>
        <a-table-column title="输入(k)" data-index="prompt_token_desc" :width="140">
          <template #default="{ record }">{{ record.prompt_token_desc }}</template>
        </a-table-column>
        <a-table-column title="输出(k)" data-index="completion_token_desc" :width="140">
          <template #default="{ record }">{{ record.completion_token_desc }}</template>
        </a-table-column>
        <a-table-column title="日期" data-index="date" :width="140">
          <template #default="{ record }">{{ record.date }}</template>
        </a-table-column>
      </a-table>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import dayjs from 'dayjs'
import { getRobotList } from '@/api/robot/index.js'
import { tokenApp, tokenAppChart } from '@/api/manage/index.js'
import LineChart from './components/line-chart.vue'
import { getDateRangePresets, tableToExcel } from '@/utils/index'
const dateRangePresets = getDateRangePresets()

const dates = ref([dayjs().subtract(7, 'day'), dayjs()])
const cascaderValue = ref([])

const searchState = reactive({
  start_date: dates.value[0].format('YYYY-MM-DD'),
  end_date: dates.value[1].format('YYYY-MM-DD'),
  token_app_type: '',
  robot_id: ''
})

const pager = reactive({
  page: 1,
  size: 20,
  total: 0
})

const tableData = ref([])

const options = ref([])
const onSearch = () => {
  pager.page = 1
  getAppCharts()
  getAppList()
}

const onTableChange = (pagination) => {
  pager.page = pagination.current
  pager.size = pagination.pageSize
  getAppList()
}

let token_app_type_map = {
  chatwiki_robot: '机器人',
  workflow: '工作流',
  other: '其他'
}

const getAppList = () => {
  tokenApp({
    ...searchState,
    page: pager.page,
    size: pager.size
  }).then((res) => {
    let lsit = res.data.list || []

    tableData.value = lsit.map((item) => {
      return {
        ...item,
        token_app_type_desc: token_app_type_map[item.token_app_type],
        completion_token_desc: formatNum(item.completion_token),
        prompt_token_desc: formatNum(item.prompt_token),
        total_token_desc: formatNum(item.total_token)
      }
    })

    pager.total = +res.data.total
  })
}

const handleExport = () => {
  tokenApp({
    ...searchState,
    page: pager.page,
    size: 10000
  }).then((res) => {
    let lsit = res.data.list || []

    lsit = lsit.map((item) => {
      return {
        ...item,
        robot_name: item.robot_name || '',
        token_app_type_desc: token_app_type_map[item.token_app_type],
        completion_token_desc: formatNum(item.completion_token),
        prompt_token_desc: formatNum(item.prompt_token),
        total_token_desc: formatNum(item.total_token)
      }
    })

    let headers = '类型,应用名称,合计消耗(k),输入(k),输出(k),日期\n'
    let fieds = [
      'token_app_type_desc',
      'robot_name',
      'total_token_desc',
      'prompt_token_desc',
      'completion_token_desc',
      'date'
    ]
    tableToExcel(headers, lsit, fieds, '明细数据-' + dayjs().format('YYYY-MM-DD HH:mm:ss'))
  })
}

const lineChartData = reactive({
  xAxis: [],
  series: []
})
const getAppCharts = () => {
  tokenAppChart({
    ...searchState
  }).then((res) => {
    let list = res.data.list
    let xData = []
    let yData = []
    list.forEach((item) => {
      xData.push(dayjs(item.date).format('YY/MM/DD'))
      yData.push({
        completion_token: formatNum(item.completion_token),
        prompt_token: formatNum(item.prompt_token),
        total_token: formatNum(item.total_token)
      })
    })
    Object.assign(lineChartData, {
      xAxis: xData,
      series: yData
    })
  })
}

function formatNum(num) {
  if (num <= 0) {
    return 0
  }
  return (num / 1000).toFixed(2)
}

const handleDateChange = () => {
  searchState.start_date = ''
  searchState.end_date = ''
  if (dates.value && dates.value.length > 0) {
    searchState.start_date = dates.value[0].format('YYYY-MM-DD')
    searchState.end_date = dates.value[1].format('YYYY-MM-DD')
  }
  onSearch()
}

const handleCascaderChange = () => {
  searchState.token_app_type = ''
  searchState.robot_id = ''
  if (!cascaderValue.value) {
    onSearch()
    return
  }
  searchState.token_app_type = cascaderValue.value[0]
  if (cascaderValue.value.length >= 2) {
    searchState.robot_id = cascaderValue.value[1]
  }
  onSearch()
}

onMounted(() => {
  getRobotList().then((res) => {
    let robotLists = res.data || []
    let chatRobot = robotLists.filter((item) => item.application_type == 0)
    let workflowRobot = robotLists.filter((item) => item.application_type == 1)
    options.value = [
      {
        label: '机器人',
        value: 'chatwiki_robot',
        children: chatRobot.map((item) => {
          return {
            label: item.robot_name,
            value: item.id
          }
        })
      },
      {
        label: '工作流',
        value: 'workflow',
        children: workflowRobot.map((item) => {
          return {
            label: item.robot_name,
            value: item.id
          }
        })
      },
      {
        label: '其他',
        value: 'other'
      }
    ]
  })
  onSearch()
})
</script>

<style lang="less" scoped>
.statics-content {
  padding: 0 16px;
  .search-box {
    display: flex;
    align-items: center;
    gap: 16px;
  }
}

.list-box {
  .main-title {
    margin-top: 16px;
    font-size: 16px;
    font-weight: 600;
    display: flex;
    align-items: center;
    &::before {
      content: '';
      display: inline-block;
      width: 4px;
      height: 20px;
      background: #2475fc;
      margin-right: 8px;
      border-radius: 6px;
    }
  }
}
</style>
