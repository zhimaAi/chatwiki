<template>
  <div class="list-wrapper">
    <div class="content-wrapper">
      <a-alert show-icon class="alert-box">
        <template #icon>
          <div class="tip-icon">
            <ExclamationCircleFilled />
          </div>
        </template>
        <template #message>
          <div class="tip-title">统计用户在与机器人对话中，没有从知识库中检索出知识的用户问题</div>
          <div class="tip-content">
            <p>
              1.
              如果您认为知识库里面有对应的知识，只是没有检索出来，可能是您设置的召回阈值过高或者top
              K较少。您调整机器人的召回设置来提升召回率。
            </p>
          </div>
        </template>
      </a-alert>
      <div class="search-block">
        <DateSelect datekey="4" @dateChange="onDateChange" />
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
          <template v-if="column.key === 'question'">
            <div v-if="record.question.length > 100">
              <a-tooltip overlayClassName="search-content-tip" placement="top">
                <template #title>{{ record.question }}</template>
                <div>{{ record.question.slice(0, 100) + '...' }}</div>
              </a-tooltip>
            </div>
            <div v-else>{{ record.question }}</div>
          </template>
        </template>
      </a-table>
    </div>
  </div>
  <AiGenerate ref="aiGenerateRef" @handleDownload="handleDownload" />
</template>

<script setup>
import DateSelect from './components/date.vue'
import AiGenerate from './components/ai-generate-modal.vue'
import { unknownIssueStats } from '@/api/chat'
import { ExclamationCircleFilled } from '@ant-design/icons-vue'
import { useRoute } from 'vue-router'
import { reactive, ref, onMounted } from 'vue'


const query = useRoute().query

const activeKey = ref(1)
const tableData = ref([])

const requestParams = reactive({
  robot_id: query.id,
  start_day: '',
  end_day: '',
  page: 1,
  size: 10,
  total: 0
})

const aiGenerateRef = ref(null)


const onTableChange = (pagination) => {
  requestParams.page = pagination.current
  requestParams.size = pagination.pageSize
  getData()
}

const onSearch = () => {
  requestParams.page = 1
  getData()
}

const handleDownload = () => {
  activeKey.value = 2
  if (aiGenerateRef.value) {
    aiGenerateRef.value.hideModal()
  }
}

const onDateChange = (data) => {
  requestParams.start_day = data.start_time
  requestParams.end_day = data.end_time

  onSearch()
}

const getData = () => {
  unknownIssueStats({
    ...requestParams
  }).then((res) => {
    let lists = res.data.list
    tableData.value = lists
    requestParams.total = +res.data.total
  })
}

const columns = [
  {
    title: '未知问题',
    dataIndex: 'question',
    key: 'question',
    width: '60%'
  },
  {
    title: '触发日期',
    dataIndex: 'show_date',
    key: 'show_date',
    width: '20%'
  },
  {
    title: '当日出现次数',
    key: 'trigger_total',
    dataIndex: 'trigger_total',
    width: '20%'
  }
]

onMounted(() => {
  // getData()
})
</script>

<style lang="less" scoped>
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
  padding-top: 16px;
}
.content-wrapper {
  padding: 0 24px 16px 24px;
}
.search-block {
  margin-top: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
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

.alert-box {
  position: relative;

  ::v-deep(.ant-alert-icon) {
    position: absolute;
    top: 13px;
    left: 20px;
  }

  ::v-deep(.ant-alert-content) {
    padding-left: 32px;
  }

  .tip-title {
    color: #242933;
    font-size: 14px;
    font-style: normal;
    font-weight: 600;
    line-height: 22px;
  }

  .tip-content p {
    color: #3a4559;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }
}
</style>

<style lang="less">
.search-content-tip {
  max-width: 800px !important;

  .ant-tooltip-inner {
    max-width: 800px !important;
    width: max-content !important;
    max-height: 300px !important;
    min-height: 30px !important;
    overflow: hidden !important;
    overflow-y: auto !important;

    /* 滚动条样式 */
    &::-webkit-scrollbar {
      width: 4px; /*  设置纵轴（y轴）轴滚动条 */
      height: 4px; /*  设置横轴（x轴）轴滚动条 */
    }
    /* 滚动条滑块（里面小方块） */
    &::-webkit-scrollbar-thumb {
      border-radius: 0px;
      background: transparent;
    }
    /* 滚动条轨道 */
    &::-webkit-scrollbar-track {
      border-radius: 0;
      background: transparent;
    }

    /* hover时显色 */
    &:hover::-webkit-scrollbar-thumb {
      background: rgba(0, 0, 0, 0.2);
    }
    &:hover::-webkit-scrollbar-track {
      background: rgba(0, 0, 0, 0.1);
    }
  }
}
</style>
