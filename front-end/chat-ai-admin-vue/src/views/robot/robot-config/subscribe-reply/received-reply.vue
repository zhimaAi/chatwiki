<template>
  <div class="user-model-page">
    <!-- <div class="page-title">{{ t('title_subscribe_reply') }}</div> -->
    <div class="breadcrumb-wrap">
      <svg-icon @click="goBack" name="back" style="font-size: 20px;" />
      <div @click="goBack" class="breadcrumb-title">{{ t('title_subscribe_reply') }}</div>
       <a-switch v-model:checked="enabled_status" :checkedValue="'1'" :un-checkedValue="'0'" :checkedChildren="t('switch_on')" :unCheckedChildren="t('switch_off')" @change="handleSwitchChange" />
       <span class="switch-tip">{{ t('switch_tip') }}</span>
    </div>
    <div class="mp-list-block">
      <div class="mp-list" :class="{ expanded }" ref="mpListRef">
        <div class="mp-card" v-for="mp in (expanded ? mpAccounts : mpAccounts.slice(0, visibleCount))" :key="mp.id" :class="{ selected: mp.appid === selectedAppid }" @click="selectMp(mp)">
          <img :src="mp.logo" class="mp-logo" />
          <span class="mp-name">{{ mp.name }}</span>
        </div>
        <a-button v-if="!expanded && mpAccounts.length > visibleCount" type="dashed" class="more-btn" @click="expanded = true">
          {{ t('btn_more') }} +{{ mpAccounts.length - visibleCount }}
        </a-button>
      </div>
    </div>
    <a-tabs v-model:activeKey="rule_type" @change="onRuleTypeChange" style="margin-top: 8px;">
      <a-tab-pane key="subscribe_reply_default" :tab="t('tab_default_reply')" />
      <a-tab-pane key="subscribe_reply_duration" :tab="t('tab_duration_reply')" />
      <a-tab-pane key="subscribe_reply_source" :tab="t('tab_source_reply')" />
    </a-tabs>
    <a-alert show-icon>
      <template #message>
        <p v-if="rule_type === 'subscribe_reply_default'">
          {{ t('msg_default_reply_limit') }}<span class="text_color_ed744a">{{ t('label_verified_only') }}</span>
        </p>
        <p v-else-if="rule_type === 'subscribe_reply_duration'">
          {{ t('msg_duration_priority') }} <span class="text_color_ed744a">{{ t('label_verified_only') }}</span>
        </p>
        <p v-else>{{ t('msg_source_priority') }}<span class="text_color_ed744a"> {{ t('label_verified_only') }}</span></p>
      </template>
    </a-alert>
    <div class="search-block" v-if="rule_type !== 'subscribe_reply_default'">
      <div class="left-block">
        <a-button type="primary" @click="handleAddReply">
          <template #icon>
            <PlusOutlined />
          </template>
          {{ rule_type === 'subscribe_reply_duration' ? t('btn_add_duration_reply') : t('btn_add_reply') }}
        </a-button>
        <div class="search-item">
          <a-select
            v-model="reply_type"
            :placeholder="t('ph_reply_content')"
            allowClear
            :options="replyTypeOptions"
            style="width: 240px;"
            @change="onReplyTypeChange"
          />
        </div>
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
            {{ t('label_priority') }}
            <a-tooltip :title="t('tooltip_priority')">
              <QuestionCircleOutlined />
            </a-tooltip>
          </span>
          <span v-else-if="column.key === 'switch_status'">
            {{ t('label_switch_status') }}
            <a-tooltip :title="rule_type === 'receive_reply_message_type' ? t('tooltip_switch_status_message') : t('tooltip_switch_status_duration')">
              <QuestionCircleOutlined />
            </a-tooltip>
          </span>
          <span v-else>{{ column.title }}</span>
        </template>
        <template #bodyCell="{ column, record }">

          <template v-if="column.key === 'subscribe_source'">
            <span style="color:#595959;">{{ formatMessageType(record) }}</span>
          </template>
          <template v-if="column.key === 'reply_content'">
            <span style="color:#595959;">{{ summarizeReplyTypes(record.reply_content) || '--' }}</span>
          </template>

          <template v-if="column.key === 'reply_num'">
            {{ record.reply_num == 0 ? t('label_all_reply') : t('label_random_reply') }}
          </template>

          <template v-if="column.key === 'create_time'">
            <span style="color:#595959;">{{ formatDateFn(record.create_time) || '--' }}</span>
          </template>

          <template v-if="column.key === 'duration'">
            <span style="color:#595959; white-space: pre-wrap;">{{ formatDurationLabel(record) }}</span>
          </template>

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

          <template v-if="column.key === 'switch_status'">
            <a-switch
              :checked="record.switch_status"
              :checkedValue="'1'"
              :un-checkedValue="'0'"
              :checkedChildren="t('switch_on')"
              :unCheckedChildren="t('switch_off')"
              @change="handleReplySwitchChange(record, $event)"
            />
          </template>

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

    <div v-if="rule_type === 'subscribe_reply_default'" class="reply-editor">
      <div class="reply-editor-item" v-for="(rule, ri) in defaultRules" :key="rule.id">
        <div class="switch-box">
          <div class="nav-box">{{ t('label_auto_reply') }}（{{ t('label_priority') }}：{{ rule.priority_num }}）</div>
          <a-switch
            :checked="rule.switch_status"
            :checkedValue="'1'"
            :un-checkedValue="'0'"
            :checkedChildren="t('switch_on')"
            :unCheckedChildren="t('switch_off')"
            @change="(val) => handleRuleSwitchChange(ri, val)"
          />
        </div>
        <transition name="collapse">
        <div class="reply-main" v-show="rule.switch_status === '1'">
          <div class="content-box">
            <div class="nav-box">{{ t('label_reply_content') }}：</div>
            <div class="item-box">
              <MultiReply v-for="(it, idx) in rule.replyList" :key="idx" :ref="el => setRuleReplyRef(ri, idx, el)" v-model:value="rule.replyList[idx]" :reply_index="idx"
                @change="(payload) => onRuleContentChange(ri, payload)" @del="(index) => onRuleDelItem(ri, index)" />
              <a-button type="dashed" style="width: 694px;" :disabled="rule.replyList.length >= 5" @click="() => addRuleReplyItem(ri)">
                <template #icon>
                  <PlusOutlined />
                </template>
                {{ t('btn_add_reply_content') }}({{rule.replyList.length}}/5)
              </a-button>
            </div>
          </div>
          <div class="method-box">
            <div class="nav-box">{{ t('label_reply_method') }}：</div>
            <a-radio-group v-model:value="rule.reply_num">
              <a-radio :value="0">{{ t('label_all_reply') }}</a-radio>
              <a-radio :value="1">{{ t('label_random_reply') }}</a-radio>
            </a-radio-group>
          </div>
          <div style="margin-top: 8px;">
            <a-button type="primary" @click="onSaveRule(ri)">{{ t('btn_save') }}</a-button>
          </div>
        </div>
        </transition>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import MultiReply from '@/components/replay-card/multi-reply.vue'
import { QuestionCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { getRobotSubscribeReplyList, updateRobotSubscribeReplyPriorityNum, updateRobotSubscribeReplySwitchStatus, deleteRobotSubscribeReply, saveRobotSubscribeReply, getSpecifyAbilityConfig, saveUserAbility } from '@/api/explore/index.js'
import { getWechatAppList } from '@/api/robot'
import { REPLY_TYPE_OPTIONS, REPLY_TYPE_LABEL_MAP, SUBSCRIBE_SOURCE_OPTIONS } from '@/constants/index'
import { message, Modal } from 'ant-design-vue'
import dayjs from 'dayjs'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.subscribe-reply.received-reply')

const replyRefs = ref([])
const ruleReplyRefs = ref({})
const query = useRoute().query
const route = useRoute()
const router = useRouter()

const rule_type = ref(query.subscribe_rule_type || 'subscribe_reply_default')
const defaultRules = ref([])
const enabled_status = ref('0')

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

const getWechatAppListFn = async() => {
  try {
    const res = await getWechatAppList({ app_type: 'official_account', app_name: '' })
    const list = Array.isArray(res?.data) ? res.data : []
    mpAccounts.value = list.filter((it) => it.account_is_verify == 'true').map((it) => ({ id: it.id, appid: it.app_id, name: it.app_name, logo: it.app_avatar }))
    selectedAppid.value = mpAccounts.value[0]?.appid || ''
  } catch (_e) {
    mpAccounts.value = []
    selectedAppid.value = ''
  }
}

onMounted(async () => {
  try {
    const res = await getSpecifyAbilityConfig({ ability_type: 'robot_subscribe_reply' })
    const item = res?.data
    const status = String(item?.user_config?.switch_status ?? '0')
    enabled_status.value = status
  } catch (_) { enabled_status.value = '0' }
  await getWechatAppListFn()
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
const replyTypeOptions = REPLY_TYPE_OPTIONS()
const tableData = ref([])
const loading = ref(false)
const reply_type = ref('')
const getTableData = () => {
  const parmas = {
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
    message.error(t('msg_invalid_priority'))
    return
  }
  try {
    await updateRobotSubscribeReplyPriorityNum({ id: record.id, priority_num: val, appid: selectedAppid.value || '' })
    message.success(t('msg_priority_updated'))
    getTableData()
  } catch (e) {
  }
}

const handleAddReply = () => {
  router.push({
    path: '/explore/index/subscribe-reply/add-rule',
    query: {
      subscribe_rule_type: rule_type.value,
      appid: selectedAppid.value || ''
    }
  })
}

  const handleEdit = (record) => {
    router.push({
      path: '/explore/index/subscribe-reply/add-rule',
      query: {
        rule_id: record.id,
        appid: selectedAppid.value || ''
      }
    })
  }

  const handleCopy = (record) => {
    router.push({
      path: '/explore/index/subscribe-reply/add-rule',
      query: {
        copy_id: record.id,
        appid: selectedAppid.value || ''
      }
    })
  }

const handleReplySwitchChange = (record, checked) => {
  const switch_status = checked
  updateRobotSubscribeReplySwitchStatus({ id: record.id, switch_status, appid: selectedAppid.value || '' }).then((res) => {
    if (res && res.res == 0) {
      record.switch_status = switch_status
      message.success(t('msg_operation_success'))
    }
  })
}

const handleDelete = (record) => {
  Modal.confirm({
    title: t('modal_confirm_delete'),
    okText: t('modal_confirm_ok'),
    onOk: () => {
      deleteRobotSubscribeReply({ id: record.id }).then((res) => {
        if (res && res.res == 0) {
          message.success(t('msg_delete_success'))
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
  const map = { text: '2', image: '4', card: '3', imageText: '1', url: '5', smartMenu: '6' }
  return list.map((it) => map[it.type] || '').filter(Boolean)
}

async function loadDefaultRule () {
  try {
    const res = await getRobotSubscribeReplyList({ rule_type: 'subscribe_reply_default', appid: selectedAppid.value || '', page: 1, size: 3 })
    const list = Array.isArray(res?.data?.list) ? res.data.list : []
    defaultRules.value = list.map((item) => ({
      id: Number(item.id || 0),
      priority_num: Number(item.priority_num || 0),
      switch_status: String(item.switch_status ?? '0'),
      reply_num: Number(item.reply_num || 0),
      replyList: (Array.isArray(item.reply_content) ? item.reply_content : []).map((rc) => ({
        type: rc?.type || rc?.reply_type || 'text',
        description: rc?.description || '',
        thumb_url: rc?.thumb_url || rc?.pic || '',
        title: rc?.title || '',
        url: rc?.url || '',
        appid: rc?.appid || '',
        page_path: rc?.page_path || '',
        smart_menu_id: rc?.smart_menu_id || '',
        smart_menu: rc?.smart_menu || {},
      }))
    }))
    const need = 3 - defaultRules.value.length
    if (need > 0) {
      const start = defaultRules.value.length + 1
      for (let i = 0; i < need; i++) {
        defaultRules.value.push({ id: 0, priority_num: start + i, switch_status: '0', reply_num: 0, replyList: [{ type: 'text', description: '' }] })
      }
    }
  } catch (e) {
    defaultRules.value = [
      { id: 0, priority_num: 1, switch_status: '0', reply_num: 0, replyList: [{ type: 'text', description: '' }] },
      { id: 0, priority_num: 2, switch_status: '0', reply_num: 0, replyList: [{ type: 'text', description: '' }] },
      { id: 0, priority_num: 3, switch_status: '0', reply_num: 0, replyList: [{ type: 'text', description: '' }] }
    ]
  }
}

async function onSaveRule (ri) {
  if (!selectedAppid.value) {
    message.warning(t('msg_select_mp'))
    return
  }
  const rule = defaultRules.value[ri]
  if (!rule || !Array.isArray(rule.replyList) || rule.replyList.length === 0) {
    message.warning(t('msg_complete_reply_content'))
    return
  }
  const arr = Array.isArray(ruleReplyRefs.value[ri]) ? ruleReplyRefs.value[ri].filter(Boolean) : []
  for (const comp of arr) {
    if (comp && comp.validate) {
      const ok = await comp.validate()
      if (!ok) { return }
    }
  }
  try {
    const payload = {
      appid: selectedAppid.value,
      rule_type: 'subscribe_reply_default',
      reply_content: JSON.stringify(serializeReplyContent(rule.replyList)),
      reply_type: serializeReplyTypeCodes(rule.replyList),
      priority_num: Number(rule.priority_num) || 0,
      reply_num: Number(rule.reply_num) || 0,
      switch_status: Number(rule.switch_status) || 0,
      id: Number(rule.id) || undefined
    }
    const res = await saveRobotSubscribeReply(payload)
    if (res && res.res == 0) {
      message.success(t('msg_save_success'))
      loadDefaultRule()
    }
  } catch (e) {
  }
}

function handleRuleSwitchChange (ri, val) {
  const next = String(val)
  const rule = defaultRules.value[ri]
  if (!rule?.id) {
    defaultRules.value[ri].switch_status = next
    message.info(t('msg_please_save_first'))
    return
  }
  updateRobotSubscribeReplySwitchStatus({ id: rule.id, switch_status: next, appid: selectedAppid.value || '' })
    .then((res) => {
      if (res && res.res == 0) {
        defaultRules.value[ri].switch_status = next
        message.success(t('msg_operation_success'))
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
  const uniq = [...new Set(labels)]
  return uniq.join('/')
}

function formatWeek (v) {
  const map = { '1': t('weekday_1'), '2': t('weekday_2'), '3': t('weekday_3'), '4': t('weekday_4'), '5': t('weekday_5'), '6': t('weekday_6'), '7': t('weekday_7') }
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
    const timeText = t('format_time_range', { start: formatTime(sd), end: formatTime(ed) })
    return weekText ? t('format_weekly', { week: weekText, time: timeText }) : t('format_weekly', { week: '', time: timeText })
  }
  if (type === 'day') {
    return t('format_daily', { time: t('format_time_range', { start: formatDate(sd), end: formatDate(ed) }) })
  }
  if (type === 'time_range') {
    const sDay = record?.start_day || ''
    const eDay = record?.end_day || ''
    const startStr = `${formatDate(sDay)} ${formatTime(sd)}`.trim()
    const endStr = `${formatDate(eDay)} ${formatTime(ed)}`.trim()
    return t('format_custom_time', { time: t('format_time_range', { start: startStr, end: endStr }) })
  }
  return `${type || ''}：${sd || ''} 至 ${ed || ''}`
}

const SUBSCRIBE_SOURCE_LABEL_MAP = Object.fromEntries((SUBSCRIBE_SOURCE_OPTIONS() || []).map((op) => [String(op.value), op.label]))

function formatMessageType (record) {
  const list = Array.isArray(record?.subscribe_source) ? record.subscribe_source : []
  if (!list.length) return '--'
  const labels = list.map((v) => SUBSCRIBE_SOURCE_LABEL_MAP[String(v)] || String(v)).filter(Boolean)
  return labels.length ? labels.join('、') : '--'
}

function addRuleReplyItem (ri) {
  const list = defaultRules.value[ri]?.replyList || []
  if (list.length >= 5) return
  list.push({ type: 'text', description: '' })
}
function onRuleContentChange (ri, payload) {
  const { reply_index, ...rest } = payload
  const list = defaultRules.value[ri]?.replyList || []
  if (reply_index >= 0 && reply_index < list.length) {
    list[reply_index] = rest
  }
}
function onRuleDelItem (ri, index) { (defaultRules.value[ri]?.replyList || []).splice(index, 1) }

function setRuleReplyRef (ri, idx, el) {
  const map = ruleReplyRefs.value
  const arr = Array.isArray(map[ri]) ? map[ri] : []
  arr[idx] = el
  map[ri] = arr
}

const goBack = () => {
  if (route.query.id && route.query.robot_key) {
    router.push({ path: '/robot/config/function-center', query: { id: route.query.id, robot_key: route.query.robot_key } })
  } else {
    router.push({ path: '/explore/index' })
  }
}

const handleSwitchChange = (checked) => {
  const prev = enabled_status.value
  const next = checked
  if (next === '0') {
    Modal.confirm({
      title: t('modal_tip'),
      content: t('modal_close_confirm'),
      onOk: () => {
        saveUserAbility({ ability_type: 'robot_subscribe_reply', switch_status: next }).then((res) => {
          if (res && res.res == 0) {
            enabled_status.value = next
            message.success(t('msg_operation_success'))
          } else {
            enabled_status.value = prev
          }
        }).catch(() => { enabled_status.value = prev })
      },
      onCancel: () => { enabled_status.value = '1' }
    })
    return
  }
  saveUserAbility({ ability_type: 'robot_subscribe_reply', switch_status: next }).then((res) => {
    if (res && res.res == 0) {
      enabled_status.value = next
      message.success(t('msg_operation_success'))
    } else {
      enabled_status.value = prev
    }
  }).catch(() => { enabled_status.value = prev })
}
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
  cursor: pointer;
  min-width: 160px;
  padding: 8px 12px;
  border-radius: 8px;
  background: #fff;
  border: 1px solid #edeff2;
  display: inline-flex;
  align-items: center;
  gap: 8px;

  &:hover {
    box-shadow: 0px 2px 4px 0px rgba(0, 0, 0, 0.08);
  }
}

.selected {
  border-color: #1890ff;
  background-color: rgba(24, 144, 255, 0.04);
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

.collapse-enter-active,
.collapse-leave-active {
  transition: max-height .2s ease, opacity .2s ease, padding .2s ease;
}
.collapse-enter-from,
.collapse-leave-to {
  max-height: 0;
  opacity: 0;
  padding-top: 0;
  padding-bottom: 0;
}
.collapse-enter-to,
.collapse-leave-from {
  max-height: 1000px;
  opacity: 1;
}
.reply-main {
  overflow: hidden;
}
.subManage-breadcrumb {
  display: flex;
  align-items: center;
  color: #000000;
  font-family: "PingFang SC";
  font-size: 14px;
  font-style: normal;
  line-height: 22px;
  padding-bottom: 16px;
}

.breadcrumb-wrap {
  width: fit-content;
  display: flex;
  align-items: center;
  cursor: pointer;
  margin-bottom: 16px;
}
.breadcrumb-title {
  margin: 0 12px 0 2px;
  color: #262626;
  font-size: 16px;
  font-style: normal;
  font-weight: 600;
  line-height: 24px;
}

.switch-tip {
  margin-left: 4px;
  color: #8c8c8c;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}
</style>
