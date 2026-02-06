<template>
  <div class="team-members-pages">
    <a-flex justify="flex-end">
      <div class="total"><div class="label total-label">{{ t('label_total') }}：</div><div class="total-num">{{ requestParams.total }} {{ t('label_item') }}</div></div>
      <div class="set-model">
        <div class="label set-model-label">{{ t('label_type') }}：</div>
        <div class="set-model-body">
          <a-select
            v-model:value="requestParams.type"
            :placeholder="t('placeholder_all_types')"
            @change="handleChangeModel"
            :style="{'width': '130px'}"
          >
            <a-select-option v-for="item in modelList" :key="item.key" :value="item.id">
              <span>{{ t(item.label) }}</span>
            </a-select-option>
          </a-select>
        </div>
      </div>

      <div class="set-date">
        <div class="label set-date-label">
          <span>{{ t('label_date') }}：</span>
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
        <a-table-column :title="t('table_column_qa')" data-index="question" width="411px">
          <template #default="{ record }">
            <div class="user-box">
              <div class="name-info">
                <div class="user-name">{{ record.question }}</div>
                <div class="user-info">{{ record.answer }}</div>
              </div>
            </div>
          </template>
        </a-table-column>
        <a-table-column :title="t('table_column_feedback')" data-index="type" width="160px">
          <template #default="{ record }">
            <div v-if="record.type == '1'" class="item-type"><svg-icon style="font-size: 24px; color: #8C8C8C;" name="like" />{{ t('label_like') }}</div>
            <div v-if="record.type == '2'" class="item-type"><svg-icon style="font-size: 24px; color: #8C8C8C;" name="dislike" />{{ t('label_dislike') }}</div>
          </template>
        </a-table-column>
        <a-table-column :title="t('table_column_feedback_content')" data-index="content" width="264px">
          <template #default="{ record }">
            <div class="item-content">
              {{ record.content }}
            </div>
          </template>
        </a-table-column>
        <a-table-column :title="t('table_column_time')" data-index="create_time" width="200px">
          <template #default="{ record }">
            <div class="item-date">{{ record.create_time }}</div>
          </template>
        </a-table-column>
        <a-table-column :title="t('table_column_action')" data-index="action" width="88px">
          <template #default="{ record }">
            <a-flex :gap="16" class="action-box">
              <a-button type="link" @click="handleOpenFeedbacksLog(record)">{{ t('btn_view_details') }}</a-button>
            </a-flex>
          </template>
        </a-table-column>
      </a-table>
    </div>
    <FeedbacksLogAlert ref="feedbacksLogAlertRef" />
  </div>
</template>
<script setup>
import { ref, reactive, onMounted  } from 'vue'
import { getFeedbackList, getFeedbackDetail } from '@/api/manage/index.js'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'
import DateSelect from './components/date.vue'
import FeedbacksLogAlert from './components/feedbacks-log-alert.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.qa-feedback.qa-feedback')

// 打开Feedbacks日志
const feedbacksLogAlertRef = ref(null)

const route = useRoute()
const datekey = ref('1')

const query = route.query
const modelList = ref([
  {
    key: 'all',
    id: 'all',
    label: 'placeholder_all_types',
    title: 'placeholder_all_types'
  },
  {
    key: '1',
    id: '1',
    label: 'label_like',
    title: 'label_like'
  },
  {
    key: '2',
    id: '2',
    label: 'label_dislike',
    title: 'label_dislike'
  }
])

const requestParams = reactive({
  robot_id: query.id, // Robot ID
  page: 1,
  size: 20,
  total: 0,
  type: 'all',
  start_date: '',
  end_date: ''
})

const onDateChange = (date) => {
  requestParams.start_date = date.start_date
  requestParams.end_date = date.end_date
  onSearch()
}

const handleOpenFeedbacksLog = async (item) => {
  const res = await getFeedbackDetail({id: item.id})
  feedbacksLogAlertRef.value.open(res.data)
}

const onReset = () => {
  // Reset
  requestParams.type = 'all'
  requestParams.start_date = ''
  requestParams.end_date = ''

  // Initialize child component
  datekey.value = 1 + '-' + Math.random()
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
const getData = () => {
  // Get list
  let parmas = {
    robot_id: requestParams.robot_id, // Robot ID
    page: requestParams.page,
    size: requestParams.size,
    start_date: requestParams.start_date,
    end_date: requestParams.end_date,
  }

  // 全部模型不传参数到后端
  if (requestParams.type != 'all') {
    parmas.type = requestParams.type
  }

  getFeedbackList(parmas).then((res) => {
    let lists = res.data.list
    lists.forEach((item) => {
      item.create_time = item.create_time > 0 ? dayjs(item.create_time * 1000).format('YYYY-MM-DD HH:mm') : '--'
      item.content = item.content ? item.content : '--'
    })
    tableData.value = lists
    requestParams.total = +res.data.total
  })
}

const handleChangeModel = (val) => {
  requestParams.type = val
  onSearch()
}

onMounted(() => {
    // Get model
    // onSearch()
})

</script>
<style lang="less" scoped>
.team-members-pages {
  position: relative;
  background: #fff;
  padding: 0 24px 24px;
  height: 100%;
  overflow-x: hidden;
  overflow-y: auto;
  
  .list-box {
    background: #fff;
    margin-top: 8px;
    .user-box {
      width: 411px;
      display: flex;
      display: -webkit-box;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 1;
      overflow: hidden;
      text-overflow: ellipsis;

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
          display: -webkit-box;
          -webkit-box-orient: vertical;
          -webkit-line-clamp: 1;
          align-self: stretch;
          overflow: hidden;
          color: #595959;
          text-overflow: ellipsis;
          font-family: "PingFang SC";
          font-size: 14px;
          font-style: normal;
          font-weight: 400;
          line-height: 22px;
        }
        .user-info {
          display: -webkit-box;
          -webkit-box-orient: vertical;
          -webkit-line-clamp: 1;
          align-self: stretch;
          overflow: hidden;
          color: #8c8c8c;
          text-overflow: ellipsis;
          font-family: "PingFang SC";
          font-size: 12px;
          font-style: normal;
          font-weight: 400;
          line-height: 20px;
        }
        .nick-name {
          color: #8c8c8c;
        }
      }
    }

    .item-type {
      color: #7a8699;
      font-family: "PingFang SC";
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
      display: flex;
      align-items: center;
    }

    .item-content {
      display: -webkit-box;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 2;
      width: 264px;
      flex: 1 0 0;
      overflow: hidden;
      color: #595959;
      text-overflow: ellipsis;
      font-family: "PingFang SC";
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
    }

    .item-date {
      color: #595959;
      font-family: "DIN";
      font-size: 14px;
      font-style: normal;
      font-weight: 500;
      line-height: 22px;
    }

    .action-box {
      color: #2475fc;
      text-align: center;
      font-family: "PingFang SC";
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
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
