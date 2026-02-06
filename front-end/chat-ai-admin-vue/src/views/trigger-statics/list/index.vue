<template>
  <div class="library-page">
    <PageTabs class="mb-16" :tabs="pageTabs" active="/trigger-statics/list"></PageTabs>
    <div class="statics-block">
      <div class="statics-item">
        <div class="title">{{ t('title_total_triggers') }}</div>
        <div class="num">{{ tip_total }}</div>
      </div>

      <div class="statics-item unknown-item cursor-pointer" @click="showUnknowQuestionModal">
        <a class="title">{{ t('unknown_questions_total') }} <RightOutlined /></a>
        <div class="num">{{ unknow_question_total }}</div>
      </div>
    </div>
    <div class="search-box">
      <a-range-picker
        v-model:value="dates"
        :allowClear="false"
        @change="handleDateChange"
        style="width: 256px"
        :presets="dateRangePresets"
      />
    </div>
    <div class="tab-box">
      <div class="main-title">{{ t('trigger_ranking') }}</div>
      <a-segmented @change="handleChangeType" v-model:value="currentType" :options="options" />
      <a-select
        v-if="currentType == 1 || currentType == 3"
        :placeholder="t('all_knowledge_base')"
        allowClear
        @change="onSearch"
        v-model:value="searchState.library_id"
        style="width: 220px"
      >
        <a-select-option v-for="item in libraryList" :value="item.id">{{
          item.library_name
        }}</a-select-option>
      </a-select>
    </div>
    <div class="library-page-body">
      <div class="list-box">
        <a-table
          :data-source="list"
          sticky
          :loading="loading"
          :pagination="{
            current: pager.page,
            total: pager.total,
            pageSize: pager.size,
            showQuickJumper: true,
            showSizeChanger: true,
            pageSizeOptions: ['10', '20', '50', '100']
          }"
          @change="onTableChange"
          :scroll="{ x: 800 }"
        >
          <a-table-column key="index" data-index="index" :title="t('ranking')" :width="100">
            <template #default="{ index }">
              {{ index + (pager.page - 1) * pager.size + 1 }}
            </template>
          </a-table-column>
          <a-table-column
            key="content"
            data-index="content"
            :title="t('knowledge_content')"
            :width="380"
            v-if="currentType == 1"
          >
            <template #default="{ record }">
              <div v-if="record.content">{{ record.content }}</div>
              <div class="qa-list-box" v-else>
                <div class="list-item">
                  <div class="list-label">{{ t('question') }}</div>
                  <div class="list-content">{{ record.question }}</div>
                </div>
                <div class="list-item" style="color: #8c8c8c">
                  <div class="list-label">{{ t('answer') }}</div>
                  <div class="list-content">{{ record.answer }}</div>
                </div>
              </div>
            </template>
          </a-table-column>

          <a-table-column key="group_name" :title="t('knowledge_base_group')" :width="120" v-if="currentType == 3">
            <template #default="{ record }"> {{ record.group_name }} </template>
          </a-table-column>

          <a-table-column key="library_name" :title="t('owned_knowledge_base')" :width="140">
            <template #default="{ record }">
              <span>{{ record.library_name }}
                <span v-if="currentType == 1">
                  <RightOutlined />
                  {{ record.group_name || t('unassigned') }}
                  <span v-if="record.library_file_name">
                    <RightOutlined />
                    {{ record.library_file_name }}
                  </span>
                </span>
              </span>
            </template>
          </a-table-column>
          <a-table-column key="tip" :title="t('trigger_count')" :width="100">
            <template #default="{ record }">
              <a-flex :gap="12">
                <span>{{ record.tip }}</span>
                <a @click="toDetail(record)">{{ t('details') }}<RightOutlined /></a>
              </a-flex>
            </template>
          </a-table-column>
          <a-table-column key="percentage" :title="t('proportion')" :width="100">
            <template #default="{ record }"> {{ record.percentage }}% </template>
          </a-table-column>
        </a-table>
      </div>
    </div>
    <DetailModal ref="detailModalRef" />
    <UnknowQuestionModal ref="unknowQuestionModalRef" />
  </div>
</template>
<script setup>
import dayjs from 'dayjs'
import { ref, reactive, onMounted } from 'vue'
import { RightOutlined } from '@ant-design/icons-vue'
import { getDateRangePresets } from '@/utils/index'
import { useI18n } from '@/hooks/web/useI18n'
import {
  statLibraryDataSort,
  statLibraryTotal,
  getLibraryList,
  statLibrarySort,
  statUnknowQuestionTotal,
  statLibraryGroupSort
} from '@/api/library'
import PageTabs from '@/components/cu-tabs/page-tabs.vue'
import DetailModal from './components/detail-modal.vue'
import UnknowQuestionModal from './components/unknow-question-modal.vue'

const { t } = useI18n('views.trigger-statics.list.index')

const dateRangePresets = getDateRangePresets()

const pageTabs = ref([
  {
    title: t('knowledge_base'),
    path: '/library/list'
  },
  {
    title: t('database'),
    path: '/database/list'
  },
  {
    title: t('document_extraction_faq'),
    path: '/ai-extract-faq/list'
  },
  {
    title: t('trigger_count_statistics'),
    path: '/trigger-statics/list'
  }
])

const tip_total = ref(0)
const unknow_question_total = ref(0)
const unknowQuestionModalRef = ref(null)
const dates = ref([dayjs().startOf('month'), dayjs()])
const currentType = ref(1)

const libraryList = ref([])

const searchState = reactive({
  begin_date_ymd: dates.value[0].format('YYYYMMDD'),
  end_date_ymd: dates.value[1].format('YYYYMMDD'),
  library_id: void 0
})

const pager = reactive({
  page: 1,
  size: 100,
  total: 0
})

const options = [
  {
    label: t('by_content'),
    value: 1
  },
  {
    label: t('by_knowledge_base_group'),
    value: 3
  },
  {
    label: t('by_knowledge_base'),
    value: 2
  }
]

const list = ref([])

const loading = ref(false)
const getContentList = () => {
  loading.value = true
  statLibraryDataSort({
    ...searchState,
    ...pager
  })
    .then((res) => {
      let datas = res.data.list || []
      datas = datas.map((item) => {
        return {
          ...item
        }
      })
      list.value = datas
      pager.total = +res.data.total
    })
    .finally(() => {
      loading.value = false
    })
}

const getListByLibrary = () => {
  loading.value = true
  statLibrarySort({
    ...searchState,
    ...pager
  })
    .then((res) => {
      let datas = res.data.list || []
      datas = datas.map((item) => {
        return {
          ...item
        }
      })
      list.value = datas
      pager.total = +res.data.total
    })
    .finally(() => {
      loading.value = false
    })
}

const getListByGroup = () => {
  loading.value = true
  statLibraryGroupSort({
    ...searchState,
    ...pager
  })
    .then((res) => {
      let datas = res.data.list || []
      datas = datas.map((item) => {
        return {
          ...item
        }
      })
      list.value = datas
      pager.total = +res.data.total
    })
    .finally(() => {
      loading.value = false
    })
}

const handleChangeType = () => {
  pager.total = 0
  if (currentType.value == 1) {
    pager.size = 100
  } else {
    pager.size = 20
  }
  onSearch()
}

const onTableChange = (pagination) => {
  pager.page = pagination.current
  pager.size = pagination.pageSize
  if (currentType.value == 1) {
    getContentList()
  }
  if (currentType.value == 2) {
    getListByLibrary()
  }

  if (currentType.value == 3) {
    getListByGroup()
  }
}

const onSearch = () => {
  pager.page = 1
  if (currentType.value == 1) {
    getContentList()
  }
  if (currentType.value == 2) {
    getListByLibrary()
  }

  if (currentType.value == 3) {
    getListByGroup()
  }
}

const handleDateChange = () => {
  searchState.begin_date_ymd = ''
  searchState.end_date_ymd = ''
  if (dates.value && dates.value.length > 0) {
    searchState.begin_date_ymd = dates.value[0].format('YYYYMMDD')
    searchState.end_date_ymd = dates.value[1].format('YYYYMMDD')
  }
  onSearch()
}

const detailModalRef = ref(null)
const toDetail = (record) => {
  detailModalRef.value.show(
    {
      begin_date_ymd: searchState.begin_date_ymd,
      end_date_ymd: searchState.end_date_ymd,
      library_id: record.library_id,
      data_id: record.data_id,
      group_id: record.group_id
    },
    currentType.value
  )
}

const getStatics = () => {
  statLibraryTotal().then((res) => {
    tip_total.value = res.data.tip_total
  })
}

const getUnknowTotal = () => {
  statUnknowQuestionTotal({
    begin_date_ymd: searchState.begin_date_ymd,

    end_date_ymd: searchState.end_date_ymd,
  }).then(res => {
    unknow_question_total.value = res.data.unknow_question_total
  })
}

const showUnknowQuestionModal = () => {
  unknowQuestionModalRef.value.show({
    begin_date_ymd: searchState.begin_date_ymd,
    end_date_ymd: searchState.end_date_ymd,
  })
}

onMounted(() => {
  onSearch()
  getStatics()
  getUnknowTotal()
  getLibraryList().then((res) => {
    libraryList.value = res.data
  })
})
</script>

<style lang="less" scoped>
.library-page {
  .list-box {
    margin-top: 8px;
  }
}

.statics-block {
  display: flex;
  align-items: center;
  gap: 16px;
  .statics-item {
    width: 208px;
    height: 90px;
    border-radius: 6px;
    background: #f2f4f7;
    padding: 16px 24px;
    .title {
      color: #7a8699;
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

.search-box {
  margin-top: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.tab-box {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 16px;
  .main-title {
    color: #333;
    font-weight: 600;
    font-size: 15px;
    display: flex;
    align-items: center;
    gap: 8px;
    &::before {
      content: '';
      display: block;
      width: 4px;
      height: 20px;
      background: #3a84ff;
      border-radius: 8px;
    }
  }
}

.qa-list-box {
  .list-item {
    display: flex;
    flex-wrap: wrap;
    line-height: 22px;
    font-size: 14px;
    color: #262626;
    margin-bottom: 6px;
    .list-label {
      margin-right: 12px;
    }
    .list-content {
      flex: 1;
      display: flex;
      flex-wrap: wrap;
      gap: 4px;
    }
  }
}
</style>
