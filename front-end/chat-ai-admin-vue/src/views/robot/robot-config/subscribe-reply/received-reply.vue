<template>
  <div class="user-model-page">
    <div class="page-title">关注后自动回复</div>
    <div class="switch-block">
      <a-switch
        @change="keyWordReplySwitchChange"
        :checked="keywordReplyStatus"
        checked-children="开"
        un-checked-children="关"
      />
      <span class="switch-desc">
        开启后，用户关注公众号后，回复指定的内容，<span style="color: #FF4D4F;">该功能仅支持公众号内回复</span>
      </span>
    </div>
    <!-- 公众号列表 -->
    <div class="mp-list-block">
      <div class="mp-list" :class="{ expanded }" ref="mpListRef">
        <div class="mp-card" v-for="mp in (expanded ? mpAccounts : mpAccounts.slice(0, visibleCount))" :key="mp.id" :class="{ selected: mp.appid === selectedAppid }" @click="selectMp(mp)">
          <img :src="mp.logo" class="mp-logo" />
          <span class="mp-name">{{ mp.name }}</span>
        </div>
        <a-button v-if="!expanded && mpAccounts.length > visibleCount" type="dashed" class="more-btn" @click="expanded = true">
          更多 +{{ mpAccounts.length - visibleCount }}
        </a-button>
      </div>
    </div>
    <a-tabs v-model:activeKey="rule_type" @change="onRuleTypeChange" style="margin-top: 8px;">
      <a-tab-pane key="subscribe_reply_default" tab="默认回复" />
      <a-tab-pane key="subscribe_reply_duration" tab="按时段设置" />
      <a-tab-pane key="subscribe_reply_source" tab="按关注来源设置" />
    </a-tabs>
    <a-alert show-icon>
      <template #message>
        <p v-if="rule_type === 'subscribe_reply_default'">关注后1分钟内只能给粉丝发送3条消息。建议关注后自动回复数量不要超过2条。可以根据关注来源或关注时间段单独设置回复</p>
        <p v-else-if="rule_type === 'subscribe_reply_duration'">优先级设置的区间为1-100的自然整数，数字越小，优先级越高。触发优先级为：按时间段回复>按关注来源回复>默认回复。</p>
        <p v-else>触发优先级为：按时间段回复>按关注来源回复>默认回复。</p>
      </template>
    </a-alert>
    <div class="search-block" v-if="rule_type !== 'subscribe_reply_default'">
      <div class="left-block">
        <a-button type="primary" @click="handleAddReply">
          <template #icon>
            <PlusOutlined />
          </template>
          {{ rule_type === 'subscribe_reply_default' ? '新增回复' : '增加时段回复' }}
        </a-button>
        <!-- 回复内容：下拉选择 图文  文本  图片  小程序 和链接 -->
        <div class="search-item">
          <a-select
            v-model="reply_type"
            placeholder="回复内容"
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
      <div class="list-box" v-if="rule_type !== 'subscribe_reply_default'">
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
            优先级
            <a-tooltip title="优先级设置的数值不能重复，数值越小，优先级越高">
              <QuestionCircleOutlined />
            </a-tooltip>
          </span>
          <span v-else-if="column.key === 'switch_status'">
            启用状态
            <a-tooltip :title="rule_type === 'receive_reply_message_type' ? '开启开关后，触发消息类型回复指定的内容' : '开启后，按照设置的时间段回复指定的内容'">
              <QuestionCircleOutlined />
            </a-tooltip>
          </span>
          <span v-else>{{ column.title }}</span>
        </template>
        <template #bodyCell="{ column, record }">

          <!-- 关注来源 -->
          <template v-if="column.key === 'subscribe_source'">
            <span style="color:#595959;">{{ formatMessageType(record) }}</span>
          </template>
          <!-- 回复内容 -->
          <template v-if="column.key === 'reply_content'">
            <span style="color:#595959;">{{ summarizeReplyTypes(record.reply_content) || '--' }}</span>
          </template>

          <template v-if="column.key === 'reply_num'">
            {{ record.reply_num == 0 ? '全部回复' : '随机回复一条' }}
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
              checked-children="开"
              un-checked-children="关"
              @change="handleReplySwitchChange(record, $event)"
            />
          </template>


          <!-- 操作 -->
          <template v-if="column.key === 'action'">
            <a-flex :gap="8">
              <a @click="handleEdit(record)">编辑</a>
              <a @click="handleDelete(record)">删除</a>
              <a @click="handleCopy(record)">复制</a>
            </a-flex>
          </template>
        </template>
      </a-table>
    </div>

    <div v-if="rule_type === 'subscribe_reply_default'" class="reply-editor">
      <div class="reply-editor-item">
        <div class="switch-box">
          <div class="nav-box">关注后自动回复：</div>
          <a-switch
            :checked="switch_status"
            :checkedValue="'1'"
            :un-checkedValue="'0'"
            checked-children="开"
            un-checked-children="关"
            @change="handleSwitchChange"
          />
        </div>
        <div class="content-box">
          <div class="nav-box">回复内容：</div>
          <div class="item-box">
            <MultiReply v-for="(it, idx) in replyList" :key="idx" v-model:value="replyList[idx]" :reply_index="idx"
              @change="onContentChange" @del="onDelItem" />
            <a-button type="dashed" style="width: 694px;" :disabled="replyList.length >= 5" @click="addReplyItem">
              <template #icon>
                <PlusOutlined />
              </template>
              添加回复内容({{replyList.length}}/5)
            </a-button>
          </div>
        </div>
        <div class="method-box">
          <div class="nav-box">回复方式：</div>
          <a-radio-group v-model:value="reply_num">
            <a-radio value="0">全部回复</a-radio>
            <a-radio value="1">随机回复一条</a-radio>
          </a-radio-group>
        </div>
        <div style="margin-top: 16px;">
          <a-button type="primary" @click="onSaveDefault">保存</a-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import MultiReply from '@/components/replay-card/multi-reply.vue'
import { QuestionCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { saveRobotAbilitySwitchStatus, getRobotSubscribeReplyList, updateRobotSubscribeReplyPriorityNum, updateRobotSubscribeReplySwitchStatus, deleteRobotSubscribeReply, saveRobotSubscribeReply } from '@/api/explore/index.js'
import { REPLY_TYPE_OPTIONS, REPLY_TYPE_LABEL_MAP } from '@/constants/index'
import { useRobotStore } from '@/stores/modules/robot'
import { message, Modal } from 'ant-design-vue'
import dayjs from 'dayjs'

const robotStore = useRobotStore()

// 来自左侧菜单的能力开关（关键词回复）
const keywordReplyStatus = computed(() => robotStore.subscribeReplySwitchStatus === '1')

const query = useRoute().query
const router = useRouter()

const rule_type = ref(query.rule_type || 'subscribe_reply_default')
const replyList = ref([{ type: 'text', description: '' }])
const reply_num = ref('0')
const currentRuleId = ref(0)
const switch_status = ref('0')

const mpAccounts = ref([])
const selectedAppid = ref('')
const expanded = ref(false)
const mpListRef = ref(null)
const visibleCount = ref(0)
const CARD_WIDTH = 160
const GAP = 8
const MORE_BTN_WIDTH = 96

function calcVisibleCount () {
  const el = mpListRef.value
  if (!el) { visibleCount.value = 0; return }
  const w = el.clientWidth || 0
  const per = CARD_WIDTH + GAP
  const count = Math.floor((w - MORE_BTN_WIDTH) / per)
  visibleCount.value = Math.max(count, 0)
}

onMounted(() => {
  mpAccounts.value = [
    { id: 'mp-1', appid: 'wx_appid_A', name: '企业公众号A', logo: 'https://dummyimage.com/48x48/2475fc/ffffff&text=A' },
    { id: 'mp-2', appid: 'wx_appid_B', name: '企业公众号B', logo: 'https://dummyimage.com/48x48/52c41a/ffffff&text=B' },
    { id: 'mp-3', appid: 'wx_appid_C', name: '企业公众号C', logo: 'https://dummyimage.com/48x48/f5222d/ffffff&text=C' },
    { id: 'mp-4', appid: 'wx_appid_D', name: '企业公众号D', logo: 'https://dummyimage.com/48x48/13c2c2/ffffff&text=D' },
    { id: 'mp-5', appid: 'wx_appid_E', name: '企业公众号E', logo: 'https://dummyimage.com/48x48/722ed1/ffffff&text=E' },
    { id: 'mp-6', appid: 'wx_appid_F', name: '企业公众号F', logo: 'https://dummyimage.com/48x48/fa8c16/ffffff&text=F' },
  ]
  selectedAppid.value = mpAccounts.value[0]?.appid || ''
  nextTick(calcVisibleCount)
  window.addEventListener('resize', calcVisibleCount)
  if (rule_type.value === 'subscribe_reply_default') {
    loadDefaultRule()
  }
})
onUnmounted(() => { window.removeEventListener('resize', calcVisibleCount) })

const columnsMessageType = [
  { title: '关注来源', dataIndex: 'subscribe_source', key: 'subscribe_source', width: 120 },
  { title: '回复内容', dataIndex: 'reply_content', key: 'reply_content', width: 120 },
  { title: '回复方式', dataIndex: 'reply_num', key: 'reply_num', width: 120 },
  { title: '创建时间', dataIndex: 'create_time', key: 'create_time', width: 120 },
  { title: '启用状态', dataIndex: 'switch_status', key: 'switch_status', width: 120 },
  { title: '操作', dataIndex: 'action', key: 'action', width: 120 }
]

const columnsDuration = [
  { title: '时间段', dataIndex: 'duration', key: 'duration', width: 220 },
  { title: '回复内容', dataIndex: 'reply_content', key: 'reply_content', width: 120 },
  { title: '回复方式', dataIndex: 'reply_num', key: 'reply_num', width: 120 },
  { title: '优先级', dataIndex: 'priority_num', key: 'priority_num', width: 120 },
  { title: '创建时间', dataIndex: 'create_time', key: 'create_time', width: 120 },
  { title: '启用状态', dataIndex: 'switch_status', key: 'switch_status', width: 120 },
  { title: '操作', dataIndex: 'action', key: 'action', width: 120 }
]

const columns = computed(() => rule_type.value === 'subscribe_reply_duration' ? columnsDuration : columnsMessageType)

const pager = reactive({
  page: 1,
  size: 10,
  total: 0
})
const replyTypeOptions = REPLY_TYPE_OPTIONS
const tableData = ref([])
const loading = ref(false)
const reply_type = ref('')
// const search_keyword = ref('')
const getTableData = () => {
  const parmas = {
    robot_id: query.id,
    rule_type: rule_type.value,
    reply_type: reply_type.value || '',
    appid: selectedAppid.value || '',
    page: pager.page,
    size: pager.size
  }
  loading.value = true
  getRobotSubscribeReplyList({
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
        subscribe_source: Array.isArray(item.subscribe_source) ? item.subscribe_source : []
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
  if (rule_type.value === 'subscribe_reply_default') {
    loadDefaultRule()
  }
}

function selectMp (mp) {
  selectedAppid.value = mp.appid
  expanded.value = true
  pager.page = 1
  getTableData()
  if (rule_type.value === 'subscribe_reply_default') {
    loadDefaultRule()
  }
}

async function onPriorityChange (record) {
  const val = Number(record.priority_num)
  if (!Number.isInteger(val) || val < 0) {
    message.error('请输入有效的优先级')
    return
  }
  try {
    await updateRobotSubscribeReplyPriorityNum({ id: record.id, robot_id: query.id, priority_num: val, appid: selectedAppid.value || '' })
    message.success('优先级已更新')
    getTableData()
  } catch (e) {
    message.error('更新失败，请稍后重试')
  }
}

const handleAddReply = () => {
  router.push({
    path: '/robot/ability/subscribe-reply/add-rule',
    query: {
      id: query.id,
      robot_key: query.robot_key,
      rule_type: rule_type.value,
      appid: selectedAppid.value || ''
    }
  })
}

  const handleEdit = (record) => {
    router.push({
      path: '/robot/ability/subscribe-reply/add-rule',
      query: {
        id: query.id,
        robot_key: query.robot_key,
        rule_id: record.id,
        appid: selectedAppid.value || ''
      }
    })
  }

  const handleCopy = (record) => {
    router.push({
      path: '/robot/ability/subscribe-reply/add-rule',
      query: {
        id: query.id,
        robot_key: query.robot_key,
        copy_id: record.id,
        appid: selectedAppid.value || ''
      }
    })
  }

const keyWordReplySwitchChange = (checked) => {
  const switch_status = checked ? '1' : '0'
  saveRobotAbilitySwitchStatus({ robot_id: query.id, ability_type: 'robot_subscribe_reply', switch_status }).then((res) => {
    if (res && res.res == 0) {
      robotStore.setSubscribeReplySwitchStatus(switch_status)
      message.success('操作成功')
      window.dispatchEvent(new CustomEvent('robotAbilityUpdated', { detail: { robotId: query.id } }))
    }
  })
}

const handleReplySwitchChange = (record, checked) => {
  const switch_status = checked
  updateRobotSubscribeReplySwitchStatus({ id: record.id, robot_id: query.id, switch_status, appid: selectedAppid.value || '' }).then((res) => {
    if (res && res.res == 0) {
      record.switch_status = switch_status
      message.success('操作成功')
    }
  })
}

const handleDelete = (record) => {
  // 确认删除
  Modal.confirm({
    title: '确认删除吗？',
    okText: '确认',
    onOk: () => {
      deleteRobotSubscribeReply({ id: record.id, robot_id: query.id }).then((res) => {
        if (res && res.res == 0) {
          message.success('删除成功')
          getTableData()
        }
      })
    }
  })
}

function serializeReplyContent (list) {
  return (list || []).map((it) => ({ ...it, status: '1' }))
}

function serializeReplyTypeCodes (list) {
  const map = { text: '2', image: '4', card: '3', imageText: '1', url: '5' }
  return list.map((it) => map[it.type] || '').filter(Boolean)
}

async function loadDefaultRule () {
  try {
    const res = await getRobotSubscribeReplyList({ robot_id: query.id, rule_type: 'subscribe_reply_default', appid: selectedAppid.value || '', page: 1, size: 1 })
    const first = Array.isArray(res?.data?.list) ? res.data.list[0] : null
    if (!first) {
      currentRuleId.value = 0
      replyList.value = [{ type: 'text', description: '' }]
      reply_num.value = '0'
      switch_status.value = '0'
      return
    }
    currentRuleId.value = Number(first.id || 0)
    const list = Array.isArray(first.reply_content) ? first.reply_content : []
    replyList.value = list.map((rc) => ({
      type: rc?.type || rc?.reply_type || 'text',
      description: rc?.description || '',
      thumb_url: rc?.thumb_url || rc?.pic || '',
      title: rc?.title || '',
      url: rc?.url || '',
      appid: rc?.appid || '',
      page_path: rc?.page_path || ''
    }))
    reply_num.value = String(first.reply_num ?? '0')
    switch_status.value = String(first.switch_status ?? '0')
  } catch (e) {
    message.error('加载默认回复失败')
  }
}

async function onSaveDefault () {
  if (!selectedAppid.value) {
    message.warning('请选择公众号')
    return
  }
  if (!replyList.value.length) {
    message.warning('请至少添加一条回复内容')
    return
  }
  const payload = {
    robot_id: query.id,
    appid: selectedAppid.value,
    rule_type: 'subscribe_reply_default',
    reply_content: JSON.stringify(serializeReplyContent(replyList.value)),
    reply_type: serializeReplyTypeCodes(replyList.value),
    reply_num: Number(reply_num.value) || 0,
    switch_status: Number(switch_status.value) || 0
  }
  if (currentRuleId.value) payload.id = currentRuleId.value
  try {
    const res = await saveRobotSubscribeReply(payload)
    if (res && res.res == 0) {
      message.success('保存成功')
      loadDefaultRule()
    }
  } catch (e) {
    message.error('保存失败，请稍后重试')
  }
}

function handleSwitchChange (val) {
  const next = String(val)
  if (!currentRuleId.value) {
    switch_status.value = next
    message.info('请先完善回复内容并保存')
    return
  }
  updateRobotSubscribeReplySwitchStatus({ id: currentRuleId.value, robot_id: query.id, switch_status: next, appid: selectedAppid.value || '' })
    .then((res) => {
      if (res && res.res == 0) {
        switch_status.value = next
        message.success('操作成功')
      }
    })
}

function mapReplyTypeLabel (t) {
  return REPLY_TYPE_LABEL_MAP[t] || ''
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
  const map = { '1': '周一', '2': '周二', '3': '周三', '4': '周四', '5': '周五', '6': '周六', '7': '周日' }
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
    const timeText = `${formatTime(sd)} 至 ${formatTime(ed)}`
    return weekText ? `每星期${weekText}：${timeText}` : `每星期：${timeText}`
  }
  if (type === 'day') {
    return `每天：${formatDate(sd)} 至 ${formatDate(ed)}`
  }
  if (type === 'time_range') {
    const sDay = record?.start_day || ''
    const eDay = record?.end_day || ''
    const startStr = `${formatDate(sDay)} ${formatTime(sd)}`.trim()
    const endStr = `${formatDate(eDay)} ${formatTime(ed)}`.trim()
    return `自定义时间：\n${startStr} 至 ${endStr}`
  }
  return `${type || ''}：${sd || ''} 至 ${ed || ''}`
}

function formatMessageType (record) {
  const list = Array.isArray(record?.subscribe_source) ? record.subscribe_source : []
  if (!list.length) return '--'
  return list.join('、')
}

function addReplyItem () {
  if (replyList.value.length >= 5) return
  replyList.value.push({ type: 'text', description: '' })
}
function onContentChange ({ reply_index, ...rest }) {
  if (reply_index >= 0 && reply_index < replyList.value.length) {
    replyList.value[reply_index] = rest
  }
}
function onDelItem (index) { replyList.value.splice(index, 1) }
</script>

<style lang="less" scoped>
.user-model-page {
  width: 100%;
  .page-title {
    display: flex;
    align-items: center;
    gap: 24px;
    padding-bottom: 16px;
    background-color: #fff;
    color: #000000;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }
  .search-block {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 16px;
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

.mp-list-block {
  margin: 8px 0 4px 0;
}
.mp-list {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: nowrap;
}
.mp-list.expanded {
  flex-wrap: wrap;
}
.mp-card {
  width: 160px;
  padding: 8px 12px;
  border-radius: 8px;
  background: #fff;
  border: 1px solid #edeff2;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
.mp-logo {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}
.mp-name {
  color: #262626;
  font-size: 14px;
  line-height: 22px;
}
.more-btn {
  flex: 0 0 auto;
}

.reply-editor {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 16px;

  .reply-editor-item {
    border: 1px solid #edeff2;
    padding: 12px;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    gap: 8px;

    .switch-box {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .content-box {
      margin-top: 8px;
      display: flex;
      flex-direction: column;
      gap: 8px;
    }

    .method-box {
      margin-top: 8px;
      display: flex;
      align-items: center;
      gap: 8px;
    }
  }
}
</style>
