<template>
  <div class="subManage-edit">
    <a-breadcrumb class="subManage-breadcrumb">
      <a-breadcrumb-item><a :href="autoReplyUrl">{{ t('breadcrumb_received_reply') }}</a></a-breadcrumb-item>
      <a-breadcrumb-item>{{ ruleType === 'receive_reply_message_type' ? t('breadcrumb_add_reply') : t('breadcrumb_add_time_reply') }}</a-breadcrumb-item>
    </a-breadcrumb>

    <div class="main">
      <a-form ref="formRef" :model="form" :rules="rules">
        <a-form-item name="switch_status" :rules="[{ required: true, message: t('msg_input_rule_name') }]">
          <template #label>
            {{ t('label_enable') }}
            <a-tooltip>
              <template #title>{{ t('tip_enable') }}</template>
              <span class="ml4"><QuestionCircleOutlined /></span>
            </a-tooltip>
          </template>
          <a-switch v-model:checked="form.switch_status" :checked-children="t('switch_on')" :un-checked-children="t('switch_off')" />
        </a-form-item>

        <a-form-item  v-if="ruleType === 'receive_reply_duration'" name="duration_type" :label="t('label_select_reply_time')">
          <div class="item-box">
            <a-form-item name="duration_type" :rules="[{ required: true, message: t('msg_select_reply_time_type') }]">
              <a-radio-group v-model:value="form.duration_type">
                <a-radio value="week">{{ t('radio_every_week') }}</a-radio>
                <a-radio value="day">{{ t('radio_every_day') }}</a-radio>
                <a-radio value="time_range">{{ t('radio_custom_time') }}</a-radio>
              </a-radio-group>
            </a-form-item>

            <a-form-item v-if="form.duration_type === 'week'" name="week_day"
              :rules="[{ validator: validateWeekDay }]">
              <a-select v-model:value="form.week_day" mode="multiple" style="width: 200px;" allowClear>
                <a-select-option value="1">{{ t('week_monday') }}</a-select-option>
                <a-select-option value="2">{{ t('week_tuesday') }}</a-select-option>
                <a-select-option value="3">{{ t('week_wednesday') }}</a-select-option>
                <a-select-option value="4">{{ t('week_thursday') }}</a-select-option>
                <a-select-option value="5">{{ t('week_friday') }}</a-select-option>
                <a-select-option value="6">{{ t('week_saturday') }}</a-select-option>
                <a-select-option value="7">{{ t('week_sunday') }}</a-select-option>
              </a-select>
            </a-form-item>

            <a-form-item v-if="form.duration_type === 'time_range'" name="date_range"
              :rules="[{ validator: validateDateRange }]">
              <a-range-picker v-model:value="form.date_range" format="YYYY-MM-DD" />
            </a-form-item>
          </div>
        </a-form-item>

        <!-- 回复时段 -->
        <a-form-item  v-if="ruleType === 'receive_reply_duration'" name="reply_period" :label="t('label_reply_period')" :rules="[{ required: true, message: t('msg_select_reply_period') }]">
          <!-- 范围选择 时间段 -->
          <a-time-range-picker v-model:value="form.reply_period" valueFormat="HH:mm" format="HH:mm" :allowClear="false" />
          <div class="tip-box">
            {{ t('tip_cross_day') }}
          </div>
        </a-form-item>

        <!-- 优先级 -->
        <a-form-item  v-if="ruleType === 'receive_reply_duration'" name="priority_num" :label="t('label_priority')" :rules="[{ required: true, message: t('msg_input_priority') }]">
          <a-input-number v-model:value="form.priority_num" style="width: 200px;" :min="1" :max="100" />
          <div class="tip-box" style="color: red;">
            {{ t('tip_priority') }}
          </div>
        </a-form-item>

        <!-- 触发回复间隔时间：秒 -->
        <a-form-item class="interval-item" name="reply_interval" :label="t('label_reply_interval')">
          <div class="flex-center">
            <a-input-number v-model:value="form.reply_interval" style="width: 60px;" />
            <span class="ml4">{{ t('tip_reply_interval') }}</span>
          </div>
        </a-form-item>

        <!-- 消息类型：单选 全部消息、指定消息 -->
        <a-form-item class="message-type-item" v-if="ruleType === 'receive_reply_message_type'" name="message_type" :label="t('label_message_type')">
          <a-radio-group v-model:value="form.message_type">
            <a-radio value="0">{{ t('radio_all_messages') }}</a-radio>
            <a-radio value="1">{{ t('radio_specify_messages') }}</a-radio>
          </a-radio-group>
        </a-form-item>
        <div v-if="form.message_type === '1'" style="margin-left: 98px;">
          <a-checkbox-group v-model:value="form.message_type_list" :options="messageTypeOptions" />
        </div>

        <!-- 回复内容 -->
        <div class="nav-box" style="margin-top: 24px;">
          <svg-icon name="reply-content" style="font-size: 16px;"></svg-icon>
          {{ t('title_reply_content') }}
        </div>
        <div class="item-box">
          <MultiReply v-for="(it, idx) in replyList" :key="idx" ref="replyRefs" v-model:value="replyList[idx]" :reply_index="idx" 
            @change="onContentChange" @del="onDelItem" />
          <a-button type="dashed" style="width: 694px;" :disabled="replyList.length >= 5" @click="addReplyItem">
            <template #icon>
              <PlusOutlined />
            </template>
            {{ t('btn_add_reply_content') }}({{replyList.length}}/5)
          </a-button>
        </div>

        <!-- 回复方式 -->
        <a-form-item name="reply_num" :label="t('label_reply_method')" style="margin-top: 24px;">
          <div class="radio-container">
            <a-radio-group v-model:value="form.reply_num" @change="handleReplyTypeChange">
              <a-radio value="0">{{ t('reply_all') }}</a-radio>
              <a-radio value="1">{{ t('reply_random') }}</a-radio>
            </a-radio-group>
          </div>
        </a-form-item>

        <!-- 保存 底部固定 -->
        <div class="btn-container">
          <a-button type="primary" @click="onSubmit">{{ t('btn_save') }}</a-button>
        </div>
      </a-form>
    </div>
  </div>
</template>
<script setup>
import { QuestionCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { onMounted, ref, reactive, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import MultiReply from '@/components/replay-card/multi-reply.vue'
import { getRobotReceivedMessageReply, saveRobotReceivedMessageReply } from '@/api/explore/index.js'
import dayjs from 'dayjs'
import { SUBSCRIBE_REPLY_TYPE_OPTIONS } from '@/constants/index'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.auto-reply.add-reply')
const replyRefs = ref([])
const query = useRoute().query
const ruleId = ref(+query.rule_id || +query['rule-id'] || 0)
const ruleType = ref(query.rule_type || 'receive_reply_duration')
const router = useRouter()
const autoReplyUrl = computed(() => `/#/robot/ability/auto-reply?id=${query.id}&robot_key=${query.robot_key}&active_key=received&rule_type=${ruleType.value}`)
const formRef = ref(null)
const form = reactive({
  duration_type: 'day',
  reply_period: [],
  priority_num: null,
  week_day: ['1'],
  week_duration: ['1'],
  date_range: [],
  start_day: '',
  end_day: '',
  reply_num: '0',
  message_type: '0',
  message_type_list: [],
  reply_interval: 0,
  switch_status: true
})
const rules = {
  name: [
    { required: true, message: t('msg_input_rule_name'), trigger: 'blur' }
  ]
}
// 关键词集合
const messageTypeOptions = SUBSCRIBE_REPLY_TYPE_OPTIONS()

function validateDateRange (_rule, value) {
  if (form.duration_type !== 'time_range') return Promise.resolve()
  const ok = Array.isArray(value) && value.length === 2 && value[0] && value[1]
  return ok ? Promise.resolve() : Promise.reject(t('msg_select_date_range'))
}

function validateWeekDay (_rule, value) {
  if (form.duration_type !== 'week') return Promise.resolve()
  const ok = Array.isArray(value) && value.length > 0
  return ok ? Promise.resolve() : Promise.reject(t('msg_select_at_least_one_week'))
}

watch(() => form.duration_type, (val) => {
  if (val === 'week') {
    if (!Array.isArray(form.week_day) || form.week_day.length === 0) {
      form.week_day = ['1']
    }
    form.week_duration = Array.isArray(form.week_day) ? form.week_day : [form.week_day]
  } else if (val === 'time_range') {
    if (Array.isArray(form.date_range) && form.date_range.length === 2 && form.date_range[0] && form.date_range[1]) {
      form.start_day = dayjs(form.date_range[0]).format('YYYY-MM-DD')
      form.end_day = dayjs(form.date_range[1]).format('YYYY-MM-DD')
    } else {
      form.start_day = ''
      form.end_day = ''
    }
  } else {
    form.start_day = ''
    form.end_day = ''
  }
})

watch(() => form.week_day, (v) => {
  if (form.duration_type === 'week') {
    form.week_duration = Array.isArray(v) ? v : [v]
  }
})

watch(() => form.date_range, (v) => {
  if (form.duration_type === 'time_range' && Array.isArray(v) && v.length === 2 && v[0] && v[1]) {
    form.start_day = dayjs(v[0]).format('YYYY-MM-DD')
    form.end_day = dayjs(v[1]).format('YYYY-MM-DD')
  }
})

// 回复内容列表
const replyList = ref([{ type: 'text', description: '' }])
const addReplyItem = () => {
  if (replyList.value.length >= 5) {
    message.warning(t('msg_max_reply_content'))
    return
  }
  replyList.value.push({ type: 'text', description: '' })
}
const onContentChange = (payload) => {
  const { reply_index, ...rest } = payload
  if (reply_index >= 0 && reply_index < replyList.value.length) {
    replyList.value[reply_index] = rest
  }
}
const onDelItem = (index) => {
  replyList.value.splice(index, 1)
}

const handleReplyTypeChange = () => {

}

function serializeReplyContent (list) {
  return (list || []).map((it) => ({ ...it, status: '1' }))
}

function serializeReplyTypeCodes (list) {
  const map = { text: '2', image: '4', card: '3', imageText: '1', url: '5', smartMenu: '6' }
  return list.map((it) => map[it.type] || '').filter(Boolean)
}

const onSubmit = () => {
  formRef.value?.validate().then(async () => {
    if (!replyList.value.length) {
      message.warning(t('msg_at_least_one_reply'))
      return
    }

    if (ruleType.value === 'receive_reply_duration') {
      if (!form.reply_period || form.reply_period.length !== 2) {
        message.warning(t('msg_select_reply_period_warning'))
        return
      }
      const p = Number(form.priority_num)
      if (!Number.isInteger(p) || p < 1 || p > 100) {
        message.warning(t('msg_priority_range'))
        return
      }
      if (form.duration_type === 'time_range') {
        const ok = Array.isArray(form.date_range) && form.date_range.length === 2 && form.date_range[0] && form.date_range[1]
        if (!ok) {
          message.warning(t('msg_select_custom_date_range'))
          return
        }
      }
    }

    if (ruleType.value === 'receive_reply_message_type') {
      if (form.message_type === '1') {
        if (!Array.isArray(form.message_type_list) || form.message_type_list.length === 0) {
          message.warning(t('msg_select_message_type'))
          return
        }
      }
    }

    for (const comp of replyRefs.value) {
      if (comp && comp.validate) {
        const ok = await comp.validate()
        if (!ok) { return }
      }
    }

    const base = {
      robot_id: query.id,
      switch_status: form.switch_status ? '1' : '0',
      rule_type: ruleType.value,
      reply_interval: Number(form.reply_interval) || 0,
      reply_content: JSON.stringify(serializeReplyContent(replyList.value)),
      reply_type: serializeReplyTypeCodes(replyList.value),
      reply_num: form.reply_num
    }

    let extra = {}
    if (ruleType.value === 'receive_reply_message_type') {
      extra = {
        message_type: form.message_type,
        specify_message_type: Array.isArray(form.message_type_list) ? form.message_type_list.join(',') : '',
      }
    } else {
      extra = {
        duration_type: form.duration_type,
        start_day: form.start_day,
        end_day: form.end_day,
        week_duration: form.week_duration,
        priority_num: Number(form.priority_num) || 0,
        start_duration: Array.isArray(form.reply_period) ? (form.reply_period[0].toString() || '') : '',
        end_duration: Array.isArray(form.reply_period) ? (form.reply_period[1].toString() || '') : ''
      }
    }

    const payload = { ...base, ...extra }
    if (ruleId.value) payload.id = ruleId.value
    try {
      const res = await saveRobotReceivedMessageReply(payload)
      if (res && res.res == 0) {
        message.success(t('msg_save_success'))
        router.push({ path: '/robot/ability/auto-reply', query: { id: query.id, robot_key: query.robot_key, active_key: 'received', rule_type: ruleType.value } })
      }
    } catch (e) {
      // message.error('保存失败，请稍后重试')
    }
  })
}


onMounted(async () => {
  const copyId = +(query.copy_id || 0)
  const fetchOne = async (rid) => {
    const res = await getRobotReceivedMessageReply({ id: rid, robot_id: query.id })
    const data = res?.data || {}

    // rule type
    ruleType.value = data?.rule_type || ruleType.value

    // switch status
    form.switch_status = String(data?.switch_status || '0') === '1'

    // duration settings
    form.duration_type = data?.duration_type || form.duration_type
    // week/day/time_range specifics
    if (form.duration_type === 'week') {
      const weeks = Array.isArray(data?.week_duration) ? data.week_duration.map(w => String(w)) : []
      form.week_day = weeks.length ? weeks : ['1']
      form.week_duration = [...form.week_day]
    }
    if (form.duration_type === 'time_range' || form.duration_type === 'day') {
      const sd = data?.start_day || ''
      const ed = data?.end_day || ''
      if (sd && ed) {
        form.date_range = [dayjs(sd), dayjs(ed)]
      }
    }
    // reply time period
    const rs = data?.start_duration || ''
    const re = data?.end_duration || ''
    form.reply_period = [rs, re].filter(Boolean).length === 2 ? [rs, re] : []
    form.start_duration = rs
    form.end_duration = re

    // priority
    form.priority_num = Number(data?.priority_num || form.priority_num || 0)

    // message type
    form.message_type = String(data?.message_type ?? form.message_type)
    form.message_type_list = Array.isArray(data?.specify_message_type) ? data.specify_message_type : []
    form.reply_interval = Number(data?.reply_interval || 0)

    // replies
    const list = Array.isArray(data?.reply_content) ? data.reply_content : []
    replyList.value = list.map((rc) => ({
      type: (rc?.type || rc?.reply_type || 'text'),
      description: rc?.description || '',
      thumb_url: rc?.thumb_url || rc?.pic || '',
      title: rc?.title || '',
      url: rc?.url || '',
      appid: rc?.appid || '',
      page_path: rc?.page_path || '',
      smart_menu_id: rc?.smart_menu_id || '',
      smart_menu: rc?.smart_menu || {},
    }))
    form.reply_num = String(data?.reply_num ?? form.reply_num)
  }

  try {
    if (!ruleId.value && copyId) {
      await fetchOne(copyId)
      return
    }
    if (!ruleId.value) return
    await fetchOne(ruleId.value)
  } catch (e) {
    message.error(t('msg_load_failed'))
  }
})

</script>
<style lang="less" scoped>
.subManage-edit {
  padding: 16px 24px;
  width: 100%;
  height: 100%;
  border-bottom: 1px solid #fff;
  border-right: 1px solid #fff;
  background-color: #fff;
  overflow-x: hidden;
  overflow-y: auto;

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

  .main {
    padding: 0 8px;
    border-radius: 6px;
    background-color: white;
    padding-bottom: 24px;

    .title {
      border-radius: 6px;
      padding: 12px 0 12px 24px;
      align-items: flex-start;
      border-bottom: 1px solid var(--07, #F0F0F0);
      background: #FFF;
      display: flex;
      align-items: center;
      color: #262626;
      font-family: "PingFang SC";
      font-size: 14px;
      font-style: normal;
      font-weight: 600;
      line-height: 22px;
      gap: 8px;
      margin-bottom: 24px;
    }
  }
}

.mr-8 {
  margin-right: 8px;
}

.mr16 {
  margin-right: 16px;
}

.nav-box {
  color: #262626;
  font-size: 14px;
  font-style: normal;
  font-weight: 600;
  line-height: 22px;
  margin-bottom: 16px;
}

.flex {
  display: flex;
}

.btn-container {
  position: fixed;
  bottom: 0;
  right: 16px;
  display: flex;
  width: calc(100% - 270px);
  padding: 16px 1055px 16px 32px;
  align-items: center;
  border-radius: 0 0 2px 2px;
  background: #FFF;
  box-shadow: 0 -8px 4px 0 #0000000a;
}

.flex-center {
  display: flex;
  align-items: center;
}

.ml4 {
  margin-left: 4px;
}

.tip-box {
  color: #999;
  font-size: 12px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  white-space: wrap;
  max-width: 624px;
}

::v-deep(.ant-form-item-label) {
  width: 150px;
}

.message-type-item {
  margin-bottom: 0;
}

.item-box ::v-deep(.ant-form-item) {
  margin-bottom: 0;
}
</style>