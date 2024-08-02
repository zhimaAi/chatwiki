<template>
  <div class="team-members-pages">
    <a-flex justify="flex-end">
      <div class="total"><div class="label total-label">共计：</div><div class="total-num">{{ requestParams.total }}条</div></div>
      <div class="set-model">
        <div class="label set-model-label">选择模型：</div>
        <div class="set-model-body">
          <a-select
            v-model:value="requestParams.use_model"
            placeholder="全部模型"
            @change="handleChangeModel"
            :style="{'width': '200px'}"
          >
            <a-select-option v-for="item in modelList" :key="item" :value="item">
              <span>{{ item }}</span>
            </a-select-option>
          </a-select>
        </div>
      </div>

      <div class="set-date">
        <div class="label set-date-label">
          <span>统计日期：</span>
        </div>
        <div class="set-date-body">
          <DateSelect 
            @dateChange="onDateChange"
            :key="datekey"
          ></DateSelect>
        </div>
      </div>

      <div class="set-reset">
        <a-button @click="onReset">重置</a-button>
      </div>
    </a-flex>
    <div class="list-box">
      <a-table
        :data-source="tableData"
        :pagination="requestParams.total > requestParams.size ? {
          current: requestParams.page,
          total: requestParams.total,
          pageSize: requestParams.size,
          showQuickJumper: true,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100']
        } : false"
        @change="onTableChange"
      >
        <a-table-column title="模型名称" data-index="model" width="300px">
          <template #default="{ record }">
            <div class="user-box">
              <div class="name-info">
                <div class="user-name">{{ record.corp }}（{{ record.model }}）</div>
              </div>
            </div>
          </template>
        </a-table-column>
        <a-table-column title="类型" data-index="IP" width="190px">
          <template #default="{ record }">{{ record.type }}</template>
        </a-table-column>
        <a-table-column title="Token消耗(k)" data-index="amount" width="190px">
          <template #default="{ record }">{{ record.amount }}</template>
        </a-table-column>
        <a-table-column title="日期" data-index="date" width="190px">
          <template #default="{ record }">{{ record.date }}</template>
        </a-table-column>
      </a-table>
    </div>
  </div>
</template>
<script setup>
import { ref, reactive, onMounted  } from 'vue'
import { getTokenStats } from '@/api/manage/index.js'
import dayjs from 'dayjs'
import { getTokenModels } from '@/api/model/index'
import DateSelect from './components/date.vue'

const datekey = ref(Date.now())

const modelList = ref(['全部模型'])
const requestParams = reactive({
  page: 1,
  size: 10,
  total: 0,
  use_model: '全部模型',
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
  requestParams.use_model = '全部模型'
  requestParams.start_date = ''
  requestParams.end_date = ''

  // 初始化子组件
  datekey.value = Date.now()
  onSearch()
}

const onTableChange = (pagination) => {
  requestParams.page = pagination.current
  requestParams.size = pagination.pageSize
  getData()
}
const onSearch = () => {
  requestParams.page = 1
  getData()
}
const tableData = ref([])

const formatAmount = (num1, num2) => {
  // 处理计算后的精度问题
  const num1Len = (num1.split('.')[1] ? num1.split('.')[1] : num1.split('.')[0] ).length;
  const num2Len = (num2.split('.')[1] ? num2.split('.')[1] : num2.split('.')[0] ).length;
  const maxLen = Math.pow(10, Math.max(num1Len, num2Len));
  return ((num1 * maxLen) + (num2 * maxLen)) / maxLen / 1000
}

const getData = () => {
  // 获取列表
  let parmas = {
    page: requestParams.page,
    size: requestParams.size,
    start_date: requestParams.start_date,
    end_date: requestParams.end_date,
  }

  // 全部模型不传参数到后端
  if (requestParams.use_model != '全部模型') {
    parmas.model = requestParams.use_model
  }

  getTokenStats(parmas).then((res) => {
    let lists = res.data.list
    lists.forEach((item) => {
      item.create_time = item.create_time > 0 ? dayjs(item.create_time * 1000).format('YYYY-MM-DD HH:mm') : '--'
      item.amount = formatAmount(item.prompt_token, item.completion_token)
    })
    tableData.value = lists
    requestParams.total = +res.data.total
  })
}

const handleChangeModel = (val) => {
  requestParams.use_model = val
  onSearch()
}

const getModelList = () => {
  getTokenModels({}).then((res) => {
    modelList.value = [...['全部模型'], ...res.data] || ['全部模型']
  })
}

onMounted(() => {
    // 获取模型
    getModelList()
    onSearch()
})

</script>
<style lang="less" scoped>
.team-members-pages {
  position: relative;
  background: #fff;
  padding: 0 24px 24px;
  height: 100%;
  .list-box {
    background: #fff;
    margin-top: 8px;
    .user-box {
      display: flex;
      img {
        width: 40px;
        height: 40px;
        border-radius: 8px;
        margin-right: 8px;
      }
      .name-info {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        font-size: 14px;
        line-height: 22px;
        font-weight: 400;
        .user-name {
          color: #595959;
        }
        .nick-name {
          color: #8c8c8c;
        }
      }
    }
  }
}

.total {
  position: absolute;
  left: 24px;
  top: 5px;
  display: flex;
  align-items: center;

  .total-num {
    color: #595959;
    font-family: "PingFang SC";
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
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
  margin-left: 24px;

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

.model-icon {
  height: 18px;
}

.set-date {
  display: flex;
  align-items: center;
  margin: 0 24px;
}
</style>
