<template>
  <div class="team-members-pages">
    <a-flex justify="flex-start">
      <div class="set-model">
        <div class="label set-model-label">{{ t('label_channel') }}</div>
        <div class="set-model-body">
          <a-select
            v-model:value="requestParams.channel"
            :placeholder="t('ph_all_channels')"
            @change="handleChangeModel"
            :style="{'width': '200px'}"
          >
            <a-select-option v-for="item in channelItem" :key="item.key" :value="item.id">
              <span>{{ item.label }}</span>
            </a-select-option>
          </a-select>
        </div>
      </div>

      <div class="set-date">
        <div class="label set-date-label">
          <span>{{ t('label_date') }}</span>
        </div>
        <div class="set-date-body">
          <DateSelect 
            @dateChange="onDateChange"
            :datekey="datekey"
          ></DateSelect>
        </div>
      </div>

      <div class="set-reset">
        <a-button @click="onReset">{{ t('btn_reset') }}</a-button>
      </div>
    </a-flex>
    <div class="list-box">
      <div class="list-item">
        <div class="item-top">
          <div class="item-top-left">
            <div class="item-top-title">{{ t('title_daily_active_users') }}</div>
            <a-tooltip placement="top">
              <template #title>{{ t('tooltip_daily_active_users') }}</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </div>
          <div class="item-top-right">{{ t('unit_people') }}</div>
        </div>
        <Echarts :xDataArray="xDataArray.res1" :yDataArray="yDataArray.res1"></Echarts>
      </div>
      <div class="list-item">
        <div class="item-top">
          <div class="item-top-left">
            <div class="item-top-title">{{ t('title_daily_new_users') }}</div>
            <a-tooltip placement="top">
              <template #title>{{ t('tooltip_daily_new_users') }}</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </div>
          <div class="item-top-right">{{ t('unit_people') }}</div>
        </div>
        <Echarts :xDataArray="xDataArray.res2" :yDataArray="yDataArray.res2"></Echarts>
      </div>
      <div class="list-item">
        <div class="item-top">
          <div class="item-top-left">
            <div class="item-top-title">{{ t('title_total_messages') }}</div>
            <a-tooltip placement="top">
              <template #title>{{ t('tooltip_total_messages') }}</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </div>
          <div class="item-top-right">{{ t('unit_messages') }}</div>
        </div>
        <Echarts :xDataArray="xDataArray.res3" :yDataArray="yDataArray.res3"></Echarts>
      </div>
      <div class="list-item">
        <div class="item-top">
          <div class="item-top-left">
            <div class="item-top-title">{{ t('title_token_consumption') }}</div>
            <a-tooltip placement="top">
              <template #title>{{ t('tooltip_token_consumption') }}</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </div>
          <div class="item-top-right">{{ t('unit_thousands') }}</div>
        </div>
        <Echarts :xDataArray="xDataArray.res4" :yDataArray="yDataArray.res4"></Echarts>
      </div>
    </div>
  </div>
</template>
<script setup>
import * as echarts from 'echarts'
import { ref, reactive, onMounted } from 'vue'
import { getAnalyse } from '@/api/manage/index.js'
import DateSelect from './components/date.vue'
import Echarts from './components/echarts.vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'

const route = useRoute()
const { t } = useI18n('views.robot.robot-config.statistical-analysis.statistical-analysis')
// type = 1日活用户数 2日新增用户数 3总消息数 4token消耗数
const echartItem = [1, 2, 3, 4]

const xDataArray = reactive({
  'res1': null, // 日活用户数返回的数据
  'res2': null, // 日新增用户数返回的数据
  'res3': null, // 总消息数返回的数据
  'res4': null // token消耗数返回的数据
})

const yDataArray = reactive({
  'res1': null, // 日活用户数返回的数据
  'res2': null, // 日新增用户数返回的数据
  'res3': null, // 总消息数返回的数据
  'res4': null // token消耗数返回的数据
})

// master环境只有2个 yun_master环境有5个
const channelItem = ref([
  {
    key: 'all',
    id: 'all',
    label: t('ph_all_channels'),
    title: t('ph_all_channels')
  },
  {
    key: 'yun_h5',
    id: 'yun_h5',
    label: 'WebAPP',
    title: 'WebAPP'
  },
  {
    key: 'yun_pc',
    id: 'yun_pc',
    label: t('channel_embedded_website'),
    title: t('channel_embedded_website')
  }
])

const datekey = ref('2')

const requestParams = reactive({
  channel: 'all',
  start_date: '',
  end_date: ''
})

const onDateChange = (date) => {
  requestParams.start_date = date.start_date
  requestParams.end_date = date.end_date
  onSearch()
}

const onReset = () => {
  // 重置
  requestParams.channel = 'all'

  // 初始化子组件
  datekey.value = 2 + '-' + Math.random()
}

const onSearch = () => {
  // shellphy说调四个接口还快一些
  for (let i = 0; i < echartItem.length; i++) {
    const item = echartItem[i];
    getData(item)
  }
  
}

const yDataArrayDefault = reactive({
  name: t('label_daily_active_users'),
  type: 'line',
  showSymbol: false,
  symbol: 'circle', // 标记的图形
  symbolSize: 10,
  itemStyle: {
    borderColor: '#fff', // 边框颜色
    borderWidth: 2,      // 边框宽度
    borderType: 'solid',  // 边框类型
    borderDashOffset: 0, // 控制边框的虚线样式
    borderDashArray: [10, 5], // 控制边框的虚线样式
    gap: 2,              // 间隔透明度的大小
    borderRadius: 50     // 圆角半径，可以使边框变成圆形
  },
  color: '#2475FC',
  lineStyle: {
    width: 1,
    color: '#2475FC'
  },
  label: {
    show: true, // 在折线拐点上显示数据
    fontSize: 8,
    color: "#fff",
    fontWeight: 10,
  },
  smooth: true, // 开启平滑过渡
  areaStyle: {//区域样式
    origin: "start",//向最小值方向渐变，y轴有负值要写
    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
      {
        offset: 0,
        color: "rgba(36,117,252, 0.3)",
      },
      {
        offset: 1,
        color: "rgba(36,117,252, 0)",
      },
    ]),
  },
  data: [542, 1222, 147, 5]
})

const forMatXDataArray = (arr) => {
  const newArr = []
  arr.map((item) => newArr.push(item.date.slice(5, item.date.length)))
  return newArr
}

const formatAmount = (num) => {
  // 处理计算后的精度问题
  return (num / 1000).toFixed(2)
}

const forMatYDataArray = (arr, type) => {
  const newYDataArrayDefault = JSON.parse(JSON.stringify(yDataArrayDefault))  // 深拷贝，不能改变原对象
  const newArr = []
  let currentName = '日活用户数（人）'
  switch (type) {
      case '1':
        currentName = t('label_daily_active_users')
          break;
      case '2':
        currentName = t('label_daily_new_users')
          break;
      case '3':
        currentName = t('label_total_messages')
          break;
      case '4':
        currentName = t('label_token_consumption')
          break;
  }
  arr.map((item) => newArr.push(type == '4' ? formatAmount(item.amount) : item.amount))
  newYDataArrayDefault.name = currentName
  newYDataArrayDefault.data = newArr
  return [newYDataArrayDefault]
}

const getData = (type) => {
  // 获取用户列表
  let parmas = {
    type: type,
    start_date: requestParams.start_date,
    end_date: requestParams.end_date,
    robot_id: route.query.id
  }

  // 全部渠道不传参数到后端
  if (requestParams.channel != 'all') {
    parmas.channel = requestParams.channel
  }

  getAnalyse(parmas).then((res) => {
    if (res.data.length > 0) {
      xDataArray['res' + res.data[0].type] = forMatXDataArray(res.data)
      yDataArray['res' + res.data[0].type] = forMatYDataArray(res.data, res.data[0].type)
    }
  })
}

const handleChangeModel = (val) => {
  requestParams.channel = val
  onSearch()
}

onMounted(() => {
})

</script>
<style lang="less" scoped>
.team-members-pages {
  background: #fff;
  padding: 0 24px 24px;
  height: 100%;
  .list-box {
    margin-top: 16px;
    display: flex;
    flex-wrap: wrap;
    gap: 24px;
    width: 100%;

    .list-item {
      width: calc((100% / 2) - 12px);
      height: 310px;
      flex-shrink: 0;
      border-radius: 6px;
      border: 1px solid var(--07, #F0F0F0);
      background: #FFF;

      .item-top {
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 16px 16px 0;

        .item-top-left {
          display: flex;

          .item-top-title {
            color: #262626;
            text-align: center;
            font-family: "PingFang SC";
            font-size: 14px;
            font-style: normal;
            font-weight: 600;
            line-height: 22px;
            margin-right: 6px;
          }
        }

        .item-top-right {
          color: #595959;
          text-align: center;
          font-family: "PingFang SC";
          font-size: 12px;
          font-style: normal;
          font-weight: 400;
          line-height: 20px;
        }
      }
    }
  }
}

.label {
  color: #262626;
  font-family: "PingFang SC";
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
}

.set-model {
  display: flex;
  align-items: center;

  .set-model-body {

    .set-model-select {
      display: flex;
      padding: 4px 12px;
      align-items: flex-start;
      align-self: stretch;
      border-radius: 2px;
      border: 1px solid var(--06, #D9D9D9);
      background: var(--Neutral-1, #FFF);
    }
  }
}

.set-date {
  display: flex;
  align-items: center;
  margin: 0 24px;
}
</style>
