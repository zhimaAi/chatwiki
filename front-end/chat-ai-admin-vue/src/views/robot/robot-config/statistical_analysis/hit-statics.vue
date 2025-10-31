<template>
  <div class="team-members-pages">
    <a-flex justify="flex-start" :gap="16">
      <div class="search-item">
        <div class="search-label">渠道：</div>
        <div class="search-content">
          <a-select
            v-model:value="requestParams.channel"
            placeholder="全部渠道"
            @change="handleChangeModel"
            style="width: 200px"
          >
            <a-select-option
              v-for="item in channelItem"
              :data="item"
              :key="item.app_id + item.app_type"
              :value="item.app_id + item.app_type"
            >
              <span>{{ item.app_name }}</span>
            </a-select-option>
          </a-select>
        </div>
      </div>

      <div class="search-item">
        <div class="search-label">
          <span>日期：</span>
        </div>
        <div class="search-content">
          <DateSelect
            @dateChange="onDateChange"
            :datekey="datekey"
            :disabledValue="'2025-01-01'"
          ></DateSelect>
        </div>
      </div>

      <div class="search-item">
        <a-button @click="onReset">重置</a-button>
      </div>
    </a-flex>
    <div class="statics-header">
      <div class="statics-item">
        <div class="title">
          <a-flex align="center" :gap="4"
            >知识命中率
            <a-tooltip>
              <template #title>命中知识库消息数 / 机器人回复消息总数</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </a-flex>
        </div>
        <div class="num">{{ staticsHeader.library_hit_rate }}%</div>
      </div>
      <div class="statics-item">
        <div class="title">
          <a-flex align="center" :gap="4"
            >消息总数
            <a-tooltip>
              <template #title>仅统计机器人回复对应消息数量，不含人工服务期间消息</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </a-flex>
        </div>
        <div class="num">{{ staticsHeader.message_total }}</div>
      </div>
      <div class="statics-item">
        <div class="title">
          <a-flex align="center" :gap="4"
            >知识库命中消息数
            <a-tooltip>
              <template #title>机器人回复时，检索到知识库文档，则认为命中消息</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </a-flex>
        </div>
        <div class="num">{{ staticsHeader.library_hit_total }}</div>
      </div>
      <div class="statics-item">
        <div class="title">
          <a-flex align="center" :gap="4">知识库未命中消息数 </a-flex>
          <a @click="handleToUnknow">详情</a>
        </div>
        <div class="num">{{ staticsHeader.library_miss_total }}</div>
      </div>
    </div>

    <div>
      <HitLineChart
        ref="hitLineChartRef"
        :xDataArray="lineChartData.xAxis"
        :yDataArray="lineChartData.series"
      />
    </div>
  </div>
</template>
<script setup>
import { ref, reactive, onMounted } from 'vue'
import { statAiTipAnalyse } from '@/api/manage/index.js'
import DateSelect from './components/date.vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import HitLineChart from './components/hit-line-charts.vue'
import { useRoute, useRouter } from 'vue-router'
import { useChatStore } from '@/stores/modules/chat'
import dayjs from 'dayjs'
const router = useRouter()

const chatStore = useChatStore()
const { getChannelList } = chatStore
const route = useRoute()

const channelItem = ref([])

const datekey = ref('2')

const requestParams = reactive({
  robot_id: route.query.id, // 机器人ID
  start_date: '',
  end_date: '',
  channel: ''
})

const staticsHeader = reactive({
  library_hit_rate: 0,
  library_hit_total: 0,
  library_miss_total: 0,
  message_total: 0
})

const lineChartData = reactive({
  xAxis: [],
  series: []
})

const getChannelLists = async () => {
  const res = await getChannelList({ robot_id: route.query.id })
  channelItem.value = [...[{ app_type: '', app_name: '全部渠道', app_id: '' }], ...res.data]
}

const onDateChange = (date) => {
  requestParams.start_date = date.start_date
  requestParams.end_date = date.end_date
  onSearch()
}

const onReset = () => {
  // 重置
  requestParams.channel = ''
  datekey.value = 2 + '-' + Math.random()
}

const onSearch = () => {
  getStatistics()
}

const getStatistics = async () => {
  statAiTipAnalyse({ ...requestParams }).then((res) => {
    Object.assign(staticsHeader, res.data.header)

    let list = res.data.chart_list || []
    let xData = []
    let yData = []

    xData = list.map((item) => dayjs(item.date).format('MM-DD'))

    yData = forMatYDataArray(list)

    lineChartData.xAxis = xData
    lineChartData.series = yData
  })
}

const forMatYDataArray = (arr) => {
  const series = []
  const library_hit_rate = arr.map((item) => item.library_hit_rate)
  const library_hit_total = arr.map((item) => item.library_hit_total)
  const library_miss_total = arr.map((item) => item.library_miss_total)
  const message_total = arr.map((item) => item.message_total)

  // 创建三个系列配置
  const createSeries = (name, data, color) => ({
    name,
    type: 'line',
    showSymbol: false,
    symbol: 'circle',
    symbolSize: 6,
    itemStyle: {
      color,
      borderColor: '#fff',
      borderWidth: 2
    },
    lineStyle: {
      width: 2,
      color
    },
    smooth: true,
    data
  })

  series.push(
    createSeries('消息总数', message_total, '#0079FE'),
    createSeries('命中消息', library_hit_total, '#00C292'),
    createSeries('未命中消息', library_miss_total, '#FF6B6B'),
    createSeries('命中率', library_hit_rate, '#333')
  )

  return series
}

const handleChangeModel = (val, options) => {
  requestParams.app_type = options.data.app_type
  requestParams.app_id = options.data.app_id
  onSearch()
}

const handleToUnknow = () => {
  localStorage.setItem('/robot/config/unknown_issue/activeKey', 1)
  router.push({
    path: '/robot/config/unknown_issue',
    query: {
      id: route.query.id,
      robot_key: route.query.robot_key,
      start_date: requestParams.start_date,
      end_date: requestParams.end_date,
    }
  })
}

onMounted(() => {
  getChannelLists()
})
</script>
<style lang="less" scoped>
.team-members-pages {
  background: #fff;
  padding: 0 24px 24px;
  height: 100%;
}
.search-item {
  display: flex;
  align-items: center;
  .search-label {
    color: #262626;
    font-family: 'PingFang SC';
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
  }
}

.statics-header {
  display: flex;
  gap: 24px;
  margin-top: 24px;
  width: 95%;
  .statics-item {
    flex: 1;
    height: 90px;
    border-radius: 6px;
    background: #f2f4f7;
    padding: 16px 24px;
    .title {
      color: #7a8699;
      display: flex;
      align-items: center;
      justify-content: space-between;
    }
    .num {
      line-height: 32px;
      font-weight: 600;
      font-size: 24px;
      color: #242933;
      margin-top: 4px;
    }
  }
}
</style>
