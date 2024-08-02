<template>
  <div class="team-members-pages">
    <a-flex justify="flex-start">
      <div class="set-model">
        <div class="label set-model-label">渠道：</div>
        <div class="set-model-body">
          <a-select
            v-model:value="requestParams.channel"
            placeholder="全部渠道"
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
          <span>日期：</span>
        </div>
        <div class="set-date-body">
          <DateSelect 
            @dateChange="onDateChange"
            :datekey="datekey"
          ></DateSelect>
        </div>
      </div>

      <div class="set-reset">
        <a-button @click="onReset">重置</a-button>
      </div>
    </a-flex>
    <div class="list-box">
      <div class="list-item">
        <div class="item-top">
          <div class="item-top-left">
            <div class="item-top-title">日活用户数</div>
            <a-tooltip placement="top">
              <template #title>每天进行有效沟通的用户数。有效沟通的判定标准为一次及以上的问—答。</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </div>
          <div class="item-top-right">单位（人）</div>
        </div>
        <Echarts :xDataArray="xDataArray.res1" :yDataArray="yDataArray.res1"></Echarts>
      </div>
      <div class="list-item">
        <div class="item-top">
          <div class="item-top-left">
            <div class="item-top-title">日新增用户数</div>
            <a-tooltip placement="top">
              <template #title>每天进行有效沟通的新增用户数。有效沟通的判定标准为一次及以上的—问—答。</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </div>
          <div class="item-top-right">单位（人）</div>
        </div>
        <Echarts :xDataArray="xDataArray.res2" :yDataArray="yDataArray.res2"></Echarts>
      </div>
      <div class="list-item">
        <div class="item-top">
          <div class="item-top-left">
            <div class="item-top-title">总消息数</div>
            <a-tooltip placement="top">
              <template #title>每天用户和机器人沟通过程中产生的所有消息总数。</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </div>
          <div class="item-top-right">单位（条）</div>
        </div>
        <Echarts :xDataArray="xDataArray.res3" :yDataArray="yDataArray.res3"></Echarts>
      </div>
      <div class="list-item">
        <div class="item-top">
          <div class="item-top-left">
            <div class="item-top-title">Token消耗数</div>
            <a-tooltip placement="top">
              <template #title>每天大模型消耗的token总数，包括使用大模型回答问题、优化问题、推荐问题的token消耗。</template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </div>
          <div class="item-top-right">单位（千）</div>
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

const route = useRoute()
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
    label: '全部渠道',
    title: '全部渠道'
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
    label: '嵌入网站',
    title: '嵌入网站'
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
  name: '日活用户数（人）',
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
        currentName = '日活用户数（人）'
          break;
      case '2':
        currentName = '日新增用户数（人）'
          break;
      case '3':
        currentName = '总消息数（条）'
          break;
      case '4':
        currentName = 'token消耗数（千）'
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
