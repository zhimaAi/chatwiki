<template>
  <div class="user-model-page">
    <!-- 关键词回复 开关 -->
    <div class="switch-block">
      <span class="switch-title">{{ t('auto_reply') }}</span>
      <a-switch
        @change="keyWordReplySwitchChange"
        :checked="keywordReplyStatus"
        :checked-children="t('switch_on')"
        :un-checked-children="t('switch_off')"
      />
      <span class="switch-desc">{{ t('auto_reply_desc') }}</span>
    </div>
    <a-alert show-icon>
      <template #message>
        <p>{{ t('received_reply_tip') }}</p>
        <p>{{ t('received_reply_tip2') }}</p>
      </template>
    </a-alert>
    <a-tabs v-model:activeKey="rule_type" @change="onRuleTypeChange" style="margin-top: 8px;">
      <a-tab-pane key="receive_reply_message_type" :tab="t('tab_default_reply')" />
      <a-tab-pane key="receive_reply_duration" :tab="t('tab_time_reply')" />
    </a-tabs>
    <div class="search-block">
      <div class="left-block">
        <a-button type="primary" @click="handleAddReply">
          <template #icon>
            <PlusOutlined />
          </template>
          {{ rule_type === 'receive_reply_message_type' ? t('btn_add_reply') : t('btn_add_time_reply') }}
        </a-button>
        <!-- 回复内容：下拉选择 图文  文本  图片  小程序 和链接 -->
        <div class="search-item">
          <a-select
            v-model="reply_type"
            :placeholder="t('label_reply_content')"
            allowClear
            :options="replyTypeOptions"
            style="width: 240px;"
            @change="onReplyTypeChange"
          />
        </div>

        <!-- <div class="search-item">
          <a-input-search
            v-model:value="search_keyword"
            placeholder="请输入关键词名称和规则名称搜索"
            allowClear
            style="width: 240px;"
            @search="onSearch"
          >
          </a-input-search>
        </div> -->
      </div>
    </div>
      <div class="list-box">
        <a-table
        :columns="columns"
        :data-source="tableData"
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
      >
        <template #headerCell="{ column }">
          <span v-if="column.key === 'priority_num'">
            {{ t('col_priority') }}
            <a-tooltip :title="t('col_priority_tip')">
              <QuestionCircleOutlined />
            </a-tooltip>
          </span>
          <span v-else-if="column.key === 'switch_status'">
            {{ t('col_enabled_status') }}
            <a-tooltip :title="rule_type === 'receive_reply_message_type' ? t('col_enabled_tip_default') : t('col_enabled_tip_time')">
              <QuestionCircleOutlined />
            </a-tooltip>
          </span>
          <span v-else>{{ column.title }}</span>
        </template>
        <template #bodyCell="{ column, record }">

          <!-- 消息类型 -->
          <template v-if="column.key === 'message_type'">
            <span style="color:#595959;">{{ formatMessageType(record) }}</span>
          </template>
          <!-- 回复内容 -->
          <template v-if="column.key === 'reply_content'">
            <span style="color:#595959;">{{ summarizeReplyTypes(record.reply_content) || '--' }}</span>
          </template>

          <template v-if="column.key === 'reply_num'">
            {{ record.reply_num == 0 ? t('reply_all') : t('reply_random') }}
          </template>

          <!-- 创建时间 -->
          <template v-if="column.key === 'create_time'">
            <span style="color:#595959;">{{ formatDateFn(record.create_time) || '--' }}</span>
          </template>

          <template v-if="column.key === 'duration'">
            <span style="color:#595959; white-space: pre-wrap;">{{ formatDurationLabel(record) }}</span>
          </template>

          <!-- 优先级可编辑 -->
          <template v-if="column.key === 'priority_num'">
            <a-input-number
              v-model:value="record.priority_num"
              :min="1"
              :max="100"
              :precision="0"
              style="width: 96px;"
              @blur="onPriorityChange(record)"
              @pressEnter="onPriorityChange(record)"
            />
          </template>

          <!-- 启用状态 开关-->
          <template v-if="column.key === 'switch_status'">
            <a-switch
              :checked="record.switch_status"
              :checkedValue="'1'"
              :un-checkedValue="'0'"
              :checked-children="t('switch_on')"
              :un-checked-children="t('switch_off')"
              @change="handleReplySwitchChange(record, $event)"
            />
          </template>


          <!-- 操作 -->
          <template v-if="column.key === 'action'">
            <a-flex :gap="8">
              <a @click="handleEdit(record)">{{ t('btn_edit') }}</a>
              <a @click="handleDelete(record)">{{ t('btn_delete') }}</a>
              <a @click="handleCopy(record)">{{ t('btn_copy') }}</a>
            </a-flex>
          </template>
        </template>
      </a-table>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed } from 'vue'
import { QuestionCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { saveRobotAbilitySwitchStatus, getRobotReceivedMessageReplyList, updateRobotReceivedMessageReplyPriorityNum, updateRobotReceivedMessageReplySwitchStatus, deleteRobotReceivedMessageReply } from '@/api/explore/index.js'
import { REPLY_TYPE_OPTIONS, REPLY_TYPE_LABEL_MAP, SUBSCRIBE_REPLY_TYPE_LABEL_MAP } from '@/constants/index'
import { useRobotStore } from '@/stores/modules/robot'
import { message, Modal } from 'ant-design-vue'
import dayjs from 'dayjs'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.auto-reply.received-reply')
const robotStore = useRobotStore()

// 来自左侧菜单的能力开关（关键词回复）
const keywordReplyStatus = computed(() => robotStore.keywordReplySwitchStatus === '1')

const query = useRoute().query
const router = useRouter()

const rule_type = ref(query.rule_type || 'receive_reply_message_type')

const columnsMessageType = [
  { title: t('col_message_type'), dataIndex: 'message_type', key: 'message_type', width: 120 },
  { title: t('col_reply_content'), dataIndex: 'reply_content', key: 'reply_content', width: 120 },
  { title: t('col_reply_method'), dataIndex: 'reply_num', key: 'reply_num', width: 120 },
  { title: t('col_create_time'), dataIndex: 'create_time', key: 'create_time', width: 120 },
  { title: t('col_enabled_status'), dataIndex: 'switch_status', key: 'switch_status', width: 120 },
  { title: t('col_action'), dataIndex: 'action', key: 'action', width: 120 }
]

const columnsDuration = [
  { title: t('col_duration'), dataIndex: 'duration', key: 'duration', width: 220 },
  { title: t('col_reply_content'), dataIndex: 'reply_content', key: 'reply_content', width: 120 },
  { title: t('col_reply_method'), dataIndex: 'reply_num', key: 'reply_num', width: 120 },
  { title: t('col_priority'), dataIndex: 'priority_num', key: 'priority_num', width: 120 },
  { title: t('col_create_time'), dataIndex: 'create_time', key: 'create_time', width: 120 },
  { title: t('col_enabled_status'), dataIndex: 'switch_status', key: 'switch_status', width: 120 },
  { title: t('col_action'), dataIndex: 'action', key: 'action', width: 120 }
]

const columns = computed(() => rule_type.value === 'receive_reply_duration' ? columnsDuration : columnsMessageType)

const pager = reactive({
  page: 1,
  size: 10,
  total: 0
})
const replyTypeOptions = REPLY_TYPE_OPTIONS()
const tableData = ref([])
const loading = ref(false)
const reply_type = ref('')
// const search_keyword = ref('')
const getTableData = () => {
  const parmas = {
    robot_id: query.id,
    rule_type: rule_type.value,
    reply_type: reply_type.value || '',
    page: pager.page,
    size: pager.size
  }
  loading.value = true
  getRobotReceivedMessageReplyList({
    ...parmas
  })
    .then((res) => {
      const data = res?.data || { list: [], total: 0, page: pager.page, size: pager.size }
      tableData.value = (data.list || []).map((item) => ({
        ...item,
        reply_content: Array.isArray(item.reply_content) ? item.reply_content : [],
        switch_status: String(item.switch_status ?? '0'),
        duration_type: item.duration_type || '',
        start_duration: item.start_duration || '',
        end_duration: item.end_duration || '',
        priority_num: item.priority_num ?? '',
        message_type: String(item.message_type ?? ''),
        specify_message_type: Array.isArray(item.specify_message_type) ? item.specify_message_type : []
      }))
      pager.total = +data.total || 0
    })
    .finally(() => {
      loading.value = false
    })
}
getTableData()

const onTableChange = (pagination) => {
  pager.page = pagination.current
  pager.size = pagination.pageSize
  getTableData()
}

const onSearch = () => {
  pager.page = 1
  getTableData()
}

const onReplyTypeChange = (val) => {
  reply_type.value = val
  onSearch()
}

const onRuleTypeChange = () => {
  pager.page = 1
  getTableData()
}

async function onPriorityChange (record) {
  const val = Number(record.priority_num)
  if (!Number.isInteger(val) || val < 0) {
    message.error(t('msg_priority_invalid'))
    return
  }
  try {
    await updateRobotReceivedMessageReplyPriorityNum({ id: record.id, robot_id: query.id, priority_num: val })
    message.success(t('msg_priority_updated'))
    getTableData()
  } catch (e) {
    // message.error('更新失败，请稍后重试')
  }
}

const handleAddReply = () => {
  router.push({
    path: '/robot/ability/auto-reply/add-reply',
    query: {
      id: query.id,
      robot_key: query.robot_key,
      rule_type: rule_type.value
    }
  })
}

  const handleEdit = (record) => {
    router.push({
      path: '/robot/ability/auto-reply/add-reply',
      query: {
        id: query.id,
        robot_key: query.robot_key,
        rule_id: record.id
      }
    })
  }

  const handleCopy = (record) => {
    router.push({
      path: '/robot/ability/auto-reply/add-reply',
      query: {
        id: query.id,
        robot_key: query.robot_key,
        copy_id: record.id
      }
    })
  }

const keyWordReplySwitchChange = (checked) => {
  const switch_status = checked ? '1' : '0'
  if (switch_status === '0') {
    Modal.confirm({
      title: t('title_tip'),
      content: t('msg_close_warning'),
      onOk: () => {
        saveRobotAbilitySwitchStatus({ robot_id: query.id, ability_type: 'robot_auto_reply', switch_status }).then((res) => {
          if (res && res.res == 0) {
            robotStore.setKeywordReplySwitchStatus(switch_status)
            message.success(t('msg_operation_success'))
            window.dispatchEvent(new CustomEvent('robotAbilityUpdated', { detail: { robotId: query.id } }))
          }
        })
      }
    })
    return
  }
  saveRobotAbilitySwitchStatus({ robot_id: query.id, ability_type: 'robot_auto_reply', switch_status }).then((res) => {
    if (res && res.res == 0) {
      robotStore.setKeywordReplySwitchStatus(switch_status)
      message.success(t('msg_operation_success'))
      window.dispatchEvent(new CustomEvent('robotAbilityUpdated', { detail: { robotId: query.id } }))
    }
  })
}

const handleReplySwitchChange = (record, checked) => {
  const switch_status = checked
  updateRobotReceivedMessageReplySwitchStatus({ id: record.id, robot_id: query.id, switch_status }).then((res) => {
    if (res && res.res == 0) {
      record.switch_status = switch_status
      message.success(t('msg_operation_success'))
    }
  })
}

const handleDelete = (record) => {
  // 确认删除
  Modal.confirm({
    title: t('title_confirm_delete'),
    okText: t('btn_confirm'),
    onOk: () => {
      deleteRobotReceivedMessageReply({ id: record.id, robot_id: query.id }).then((res) => {
        if (res && res.res == 0) {
          message.success(t('msg_delete_success'))
          getTableData()
        }
      })
    }
  })
}

function mapReplyTypeLabel (t) {
  return REPLY_TYPE_LABEL_MAP()[t] || ''
}

function summarizeReplyTypes (list) {
  if (!Array.isArray(list)) return ''
  const labels = list
    .map((rc) => mapReplyTypeLabel(rc?.type))
    .filter((s) => !!s)
  // 去重并使用/连接
  const uniq = [...new Set(labels)]
  return uniq.join('/')
}

function formatWeek (v) {
  const map = { '1': t('week_monday'), '2': t('week_tuesday'), '3': t('week_wednesday'), '4': t('week_thursday'), '5': t('week_friday'), '6': t('week_saturday'), '7': t('week_sunday') }
  const s = String(v || '')
  return map[s] || s
}

function formatDate (s) {
  const str = String(s || '')
  if (!str) return ''
  if (/^\d{4}-\d{2}-\d{2}/.test(str)) return str.slice(0, 10)
  return str
}

function formatDateFn (date, format = 'YYYY-MM-DD') {
  if (!date) return ''
  return dayjs(date * 1000).format(format)
}

function formatTime (s) {
  const str = String(s || '')
  if (!str) return ''
  if (/^\d{2}:\d{2}/.test(str)) return str.slice(0, 5)
  return str
}

function formatDurationLabel (record) {
  const type = record?.duration_type || ''
  const wd = record?.week_duration || ''
  const sd = record?.start_duration || ''
  const ed = record?.end_duration || ''
  if (type === 'week') {
    const weekList = Array.isArray(wd) ? wd.map(v => String(v)) : (wd ? [String(wd)] : [])
    const weekText = weekList.length ? weekList.map(v => formatWeek(v)).join('、') : ''
    const timeText = `${formatTime(sd)} ${t('label_to')} ${formatTime(ed)}`
    return weekText ? `${t('label_every_week')}${weekText}：${timeText}` : `${t('label_every_week')}：${timeText}`
  }
  if (type === 'day') {
    return `${t('label_every_day')}：${formatDate(sd)} ${t('label_to')} ${formatDate(ed)}`
  }
  if (type === 'time_range') {
    const sDay = record?.start_day || ''
    const eDay = record?.end_day || ''
    const startStr = `${formatDate(sDay)} ${formatTime(sd)}`.trim()
    const endStr = `${formatDate(eDay)} ${formatTime(ed)}`.trim()
    return `${t('label_custom_time')}：\n${startStr} ${t('label_to')} ${endStr}`
  }
  return `${type || ''}：${sd || ''} ${t('label_to')} ${ed || ''}`
}

function mapMsgLabel (t) {
  const m = SUBSCRIBE_REPLY_TYPE_LABEL_MAP()
  return m[String(t)] || String(t)
}

function formatMessageType (record) {
  const mt = String(record?.message_type ?? '')
  if (mt === '0') return t('msg_all_messages')
  const list = Array.isArray(record?.specify_message_type) ? record.specify_message_type : []
  if (!list.length) return '--'
  return list.map(mapMsgLabel).join('/')
}
</script>

<style lang="less" scoped>
.user-model-page {
  width: 100%;
  .search-block {
    display: flex;
    align-items: center;
    justify-content: space-between;
    .left-block {
      display: flex;
      align-items: center;
      gap: 16px;
    }
    .right-block {
      display: flex;
      align-items: center;
      gap: 2px;
    }
  }
  .list-box {
    margin-top: 8px;
  }
  ::v-deep(.ant-alert) {
    align-items: baseline;
  }
}

.switch-block {
  display: flex;
  align-items: center;
  margin-bottom: 16px;

  .switch-title {
    margin-right: 12px;
    color: #262626;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }
}
.switch-desc {
  margin-left: 4px;
  color: #8c8c8c;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}

.flex {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.popover-cont {
  max-width: 560px;
}
</style>
