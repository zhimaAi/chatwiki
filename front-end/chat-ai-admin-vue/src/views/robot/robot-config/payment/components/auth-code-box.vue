<template>
  <div class="content">
    <div class="overview-box">
      <div class="overview-block">
        <div class="overview-item">
          <div class="tit">{{ t('title_today_generate') }}
            <QuestionCircleOutlined/>
          </div>
          <div class="count">{{ statData.today_generate_count }}</div>
          <div class="val">
            <span class="label">{{ t('label_price') }}：</span>{{ statData.today_generate_price || 0}}
          </div>
        </div>
        <div class="overview-item">
          <div class="tit">{{ t('title_yesterday_generate') }}
            <QuestionCircleOutlined/>
          </div>
          <div class="count">{{ statData.yesterday_generate_count }}</div>
          <div class="val">
            <span class="label">{{ t('label_price') }}：</span>{{ statData.yesterday_generate_price  || 0 }}
          </div>
        </div>
        <div class="overview-item">
          <div class="tit">{{ t('title_total_generate') }}
            <QuestionCircleOutlined/>
          </div>
          <div class="count">{{ statData.total_generate_count }}</div>
          <div class="val">
            <span class="label">{{ t('label_price') }}：</span>{{ statData.total_generate_price || 0}}
          </div>
        </div>
      </div>
      <div class="overview-block">
        <div class="overview-item">
          <div class="tit">{{ t('title_today_exchange') }}
            <QuestionCircleOutlined/>
          </div>
          <div class="count">{{ statData.today_exchange_count }}</div>
          <div class="val">
            <span class="label">{{ t('label_price') }}：</span>{{ statData.today_exchange_price || 0}}
          </div>
        </div>
        <div class="overview-item">
          <div class="tit">{{ t('title_yesterday_exchange') }}
            <QuestionCircleOutlined/>
          </div>
          <div class="count">{{ statData.yesterday_exchange_count }}</div>
          <div class="val">
            <span class="label">{{ t('label_price') }}：</span>{{ statData.yesterday_exchange_price || 0}}
          </div>
        </div>
        <div class="overview-item">
          <div class="tit">{{ t('title_total_exchange') }}
            <QuestionCircleOutlined/>
          </div>
          <div class="count">{{ statData.total_exchange_count }}</div>
          <div class="val">
            <span class="label">{{ t('label_price') }}：</span>{{ statData.total_exchange_price || 0}}
          </div>
        </div>
      </div>
    </div>
    <div class="main-box">
      <div class="main-tit">
        {{ t('title_auth_code_manager') }}
        <a-divider type="vertical"/>
        <span class="tip-info">{{ t('msg_manager_tip') }}</span>
      </div>
      <div>
        <a-button type="primary" ghost @click="showManagerModal"><PlusOutlined/> {{ t('btn_add_manager') }}</a-button>
        <div class="user-list">
          <a-tag v-for="item in managers" :key="item.id">
            <img class="avatar" :src="item.manager_avatar"/>
            {{ item.manager_nickname }}
            <CloseOutlined class="close" @click="delManager(item)"/>
          </a-tag>
        </div>
      </div>
    </div>
    <div class="main-box">
      <div class="main-tit">{{ t('title_auth_code_detail') }}</div>
      <div>
        <div class="filter-box">
          <div class="left">
            <a-input v-model:value.trim="filterData.content" @change="search" allow-clear :placeholder="t('ph_search_code')" style="width: 200px">
              <template #suffix>
                <SearchOutlined/>
              </template>
            </a-input>
            <a-input v-model:value.trim="filterData.openid" @change="search" allow-clear :placeholder="t('ph_search_openid')" style="width: 200px">
              <template #suffix>
                <SearchOutlined/>
              </template>
            </a-input>
            <a-select v-model:value="filterData.usage_status" @change="search" allow-clear :placeholder="t('label_usage_status')" style="width: 200px">
              <a-select-option v-for="(text, key) in usageStatusMap" :key="key" :value="key">{{ text }}</a-select-option>
            </a-select>
          </div>
          <div class="right">
            <a-button @click="exportData">{{ t('btn_export') }}</a-button>
            <a-button type="primary" @click="showCodeModal"><PlusOutlined/> {{ t('btn_add') }}</a-button>
          </div>
        </div>
        <a-table
          :scroll="{x: 1200}"
          :loading="loading"
          :data-source="list"
          :columns="columns"
          :pagination="pagination"
          @change="tableChange"
        >
          <template #bodyCell="{column, record}">
            <template v-if="'content' === column.dataIndex">
              <div class="code-box">
                <span>{{record.content}}</span>
                <CopyOutlined class="icon" @click="copy(record.content)"/>
              </div>
            </template>
            <template v-if="'remark' === column.dataIndex">
              <a-tooltip :title="record.remark">
                <div class="zm-line2">{{record.remark}}</div>
              </a-tooltip>
            </template>
            <template v-if="'action' === column.dataIndex">
              <a @click="delCode(record)">{{ t('btn_delete') }}</a>
            </template>
          </template>
        </a-table>
      </div>
    </div>

    <AuthCodeModal ref="codeRef" :packages="packages" :robotId="route.query.id" @add-manager="showManagerModal" @ok="loadData"/>
    <ManagerModal ref="managerRef" :robotId="route.query.id"/>
  </div>
</template>

<script setup>
import {ref, reactive, onMounted, computed, h} from 'vue';
import {useRoute} from 'vue-router';
import dayjs from 'dayjs';
import {Modal, message} from 'ant-design-vue';
import {QuestionCircleOutlined, PlusOutlined, SearchOutlined, CopyOutlined, CloseOutlined, ExclamationCircleOutlined} from '@ant-design/icons-vue'
import AuthCodeModal from "./auth-code-modal.vue";
import ManagerModal from "./manager-modal.vue";
import {
  delAuthCode,
  delAuthCodeManager,
  getAuthCodeList,
  getAuthCodeManager,
  getAuthCodeStats
} from "@/api/robot/payment.js";
import {tableToExcel, copyText} from "@/utils/index.js";
import {useI18n} from '@/hooks/web/useI18n';

const {t} = useI18n('views.robot.robot-config.payment.components.auth-code-box')

const props = defineProps({
  config: {
    type: Object
  }
})
const route = useRoute()
const codeRef = ref(null)
const managerRef = ref(null)
const loading = ref(false)
const statData = reactive({
  "today_exchange_count": "0",
  "today_exchange_price": "",
  "today_generate_count": "0",
  "today_generate_price": "",
  "total_exchange_count": "0",
  "total_exchange_price": "",
  "total_generate_count": "0",
  "total_generate_price": "",
  "yesterday_exchange_count": "0",
  "yesterday_exchange_price": "",
  "yesterday_generate_count": "0",
  "yesterday_generate_price": "",
})
const managers = ref([])
const list = ref([])
const columns = ref([
  {
    title: t('label_auth_code'),
    dataIndex: 'content',
    width: 200,
  },
  {
    title: t('label_package_name'),
    dataIndex: 'package_name',
    width: 100,
  },
  {
    title: t('label_duration'),
    dataIndex: 'package_duration',
    width: 120,
  },
  {
    title: t('label_count'),
    dataIndex: 'package_count',
    width: 120,
  },
  {
    title: t('label_cost'),
    dataIndex: 'package_price',
    width: 100,
  },
  {
    title: t('label_usage_status'),
    dataIndex: 'usage_status_text',
    width: 120,
  },
  {
    title: t('label_exchanger'),
    dataIndex: 'exchanger_openid',
    width: 160,
  },
  {
    title: t('label_exchange_time'),
    dataIndex: 'exchange_date',
    width: 160,
  },
  {
    title: t('label_use_time'),
    dataIndex: 'use_date',
    width: 160,
  },
  {
    title: t('label_create_time'),
    dataIndex: 'create_date',
    width: 160,
  },
  {
    title: t('label_creator'),
    dataIndex: 'creator_name',
    width: 100,
  },
  {
    title: t('label_remark'),
    dataIndex: 'remark',
    width: 160,
  },
  {
    title: t('label_action'),
    dataIndex: 'action',
    fixed: 'right',
    width: 100,
  },
])
const filterData = reactive({
  content: "",
  openid: "",
  usage_status: undefined,
})
const pagination = reactive({
  current: 1,
  pageSize: 50,
  total: 0
})
const usageStatusMap = {
  1: t('usage_status_unused'),
  2: t('usage_status_exchanged'),
  3: t('usage_status_used'),
}
const packages = computed(() => {
  let res = props.config?.package_type == 1 ? props.config.count_package : props.config.duration_package
  return Array.isArray(res) ? res : []
})

onMounted(() => {
  init()
})

function init() {
  loadStat()
  loadManagers()
  loadData()
}

function loadStat() {
  getAuthCodeStats({
    robot_id: route.query.id
  }).then(res => {
    Object.assign(statData, res?.data || {})
  })
}

function loadManagers() {
  getAuthCodeManager({robot_id: route.query.id}).then(res => {
    managers.value = res?.data || []
  })
}

function search() {
  list.value = []
  pagination.total = 0
  pagination.current = 1
  loadData()
}

function loadData() {
  let params = {
    robot_id: route.query.id,
    page: pagination.current,
    size: pagination.pageSize,
    ...filterData
  }
  getAuthCodeList(params).then(res => {
    let _list = res?.data?.list || []
    _list.forEach(item => {
      item.usage_status_text = usageStatusMap[item?.usage_status]
      item.exchange_date = time2date(item.exchange_time)
      item.create_date = time2date(item.create_time)
      item.use_date = time2date(item.use_time)
    })
    list.value = _list
    pagination.total = Number(res?.data?.total || 0)
  })
}

function time2date(time) {
  if (time > 0) {
    return dayjs(time * 1000).format('YYYY-MM-DD HH:mm')
  }
  return ''
}

function tableChange(p) {
  Object.assign(pagination, p)
  loadData()
}

function showManagerModal() {
  managerRef.value.show()
}

function showCodeModal() {
  codeRef.value.show()
}

function delCode(record) {
  Modal.confirm({
    title: t('title_confirm'),
    icon: h(ExclamationCircleOutlined),
    content: t('msg_confirm_delete_code'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    onOk() {
      delAuthCode({id: record.id}).then(() => {
        message.success(t('msg_deleted'))
        loadData()
      })
    }
  })
}

function delManager(record) {
  Modal.confirm({
    title: t('title_confirm'),
    icon: h(ExclamationCircleOutlined),
    content: t('msg_confirm_delete_manager'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    onOk() {
      delAuthCodeManager({id: record.id, robot_id: route.query.id}).then(() => {
        message.success(t('msg_deleted'))
        loadManagers()
      })
    }
  })
}

function exportData() {
  let str = [], fields = []
  columns.value.forEach(item => {
    if (item.dataIndex === 'action') return
    str.push(item.title)
    fields.push(item.dataIndex)
  })
  tableToExcel(str.toString()+ '\n', list.value, fields, `${route.query.robot_key}${t('title_auth_code_detail')}.csv`)
}

function copy(text) {
  copyText(text)
  message.success(t('msg_copied'))
}
</script>

<style scoped lang="less">
.overview-box {
  display: flex;
  align-items: center;
  gap: 16px;
  .overview-block {
    display: flex;
    align-items: center;
    color: #595959;
    font-size: 14px;
    font-weight: 400;
    border-radius: 12px;
    background: #F2F4F7;
    width: calc(50% - 8px);

    .overview-item {
      flex-shrink: 0;
      width: 33.33%;
      display: flex;
      flex-direction: column;
      gap: 4px;
      padding: 16px;

      .count {
        color: #262626;
        font-size: 24px;
        font-weight: 500;
      }

      .label {
        color: #8c8c8c;
      }
    }
  }
}

.main-box {
  margin-top: 24px;
  .main-tit {
    color: #262626;
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 16px;
  }
  .filter-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;
    > div {
      display: flex;
      align-items: center;
      gap: 16px;
    }
  }
}

.user-list {
  margin-top: 8px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  :deep(.ant-tag) {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 5px 16px;
    border-radius: 6px;
    border: 1px solid #D9D9D9;
    background: #F5F5F5;

    .close {
      display: none;
      margin-left: 4px;
      cursor: pointer;
    }

    &:hover {
      .close {
        display: inline-block;
      }
    }
  }

  .avatar {
    flex-shrink: 0;
    width: 20px;
    height: 20px;
    border-radius: 4px;
  }
}

.tip-info {
  color: #8c8c8c;
  font-size: 14px;
  font-weight: 400;
}

.code-box {
  display: flex;
  align-items: center;
  gap: 4px;

  > span {
    overflow: hidden;
    white-space: nowrap;
    max-width: 88%;
    text-overflow: ellipsis;
  }

  .icon {
    color: #BFBFBF;
    cursor: pointer;
  }
}

.mt8 {
  margin-top: 8px;
}

</style>
