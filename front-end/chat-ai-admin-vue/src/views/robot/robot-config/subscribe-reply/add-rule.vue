<template>
  <div class="subManage-edit">
    <a-breadcrumb class="subManage-breadcrumb">
      <a-breadcrumb-item><a :href="autoReplyUrl">{{ ruleType === 'subscribe_reply_duration' ? '按时段设置' : '按关注来源设置' }}</a></a-breadcrumb-item>
      <a-breadcrumb-item>新增规则</a-breadcrumb-item>
    </a-breadcrumb>

    <div class="main">
      <a-form ref="formRef" :model="form" :rules="rules">
        <a-form-item name="switch_status" :rules="[{ required: true, message: '请输入规则名称' }]">
          <template #label>
            是否开启
            <a-tooltip>
              <template #title>开启开关后，触发消息类型回复指定的内容</template>
              <span class="ml4"><QuestionCircleOutlined /></span>
            </a-tooltip>
          </template>
          <a-switch v-model:checked="form.switch_status" checked-children="开" un-checked-children="关" />
        </a-form-item>

        <a-form-item v-if="ruleType === 'subscribe_reply_source'" name="subscribe_source" label="关注来源" :rules="[{ required: true, message: '请选择关注来源' }]">
          <a-select v-model:value="form.subscribe_source" mode="multiple" style="width: 480px;" placeholder="请选择关注来源" allowClear showSearch>
            <a-select-option v-for="op in subscribeSourceOptions" :key="op.value" :value="op.value">{{ op.label }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item  v-if="ruleType === 'subscribe_reply_duration'" name="duration_type" label="选择回复时间">
          <div class="item-box">
            <a-form-item name="duration_type" :rules="[{ required: true, message: '请选择回复时间类型' }]">
              <a-radio-group v-model:value="form.duration_type">
                <a-radio value="week">每星期</a-radio>
                <a-radio value="day">每天</a-radio>
                <a-radio value="time_range">自定义时间</a-radio>
              </a-radio-group>
            </a-form-item>

            <a-form-item v-if="form.duration_type === 'week'" name="week_day"
              :rules="[{ validator: validateWeekDay }]">
              <a-select v-model:value="form.week_day" mode="multiple" style="width: 200px;" allowClear>
                <a-select-option value="1">星期一</a-select-option>
                <a-select-option value="2">星期二</a-select-option>
                <a-select-option value="3">星期三</a-select-option>
                <a-select-option value="4">星期四</a-select-option>
                <a-select-option value="5">星期五</a-select-option>
                <a-select-option value="6">星期六</a-select-option>
                <a-select-option value="7">星期日</a-select-option>
              </a-select>
            </a-form-item>

            <a-form-item v-if="form.duration_type === 'time_range'" name="date_range"
              :rules="[{ validator: validateDateRange }]">
              <a-range-picker v-model:value="form.date_range" format="YYYY-MM-DD" />
            </a-form-item>
          </div>
        </a-form-item>

        <!-- 回复时段 -->
        <a-form-item  v-if="ruleType === 'subscribe_reply_duration'" name="reply_period" label="回复时段" :rules="[{ required: true, message: '请选择回复时段' }]">
          <!-- 范围选择 时间段 -->
          <a-time-range-picker v-model:value="form.reply_period" valueFormat="HH:mm" format="HH:mm" :allowClear="false" />
          <div class="tip-box">
            如需跨天，请单独重新再添加一条时段回复：第一个是到23:59结束，第二个是从0点开始
          </div>
        </a-form-item>

        <!-- 优先级 -->
        <a-form-item  v-if="ruleType === 'subscribe_reply_duration'" name="priority_num" label="优先级" :rules="[{ required: true, message: '输入优先级，1-100之间' }]">
          <a-input-number v-model:value="form.priority_num" style="width: 200px;" :min="1" :max="100" />
          <div class="tip-box" style="color: red;">
            数字越小，则优先级越高。优先级仅在时间相互冲突时有效，比如一个每天1-9点与一个每天310点，当8点钟触发时，则只会触发其中优先级最高的一个。
          </div>
        </a-form-item>

        <!-- 触发回复间隔时间：秒 -->
        <a-form-item v-if="ruleType === 'subscribe_reply_duration'" class="interval-item" name="reply_interval" label="触发回复间隔时间">
          <div class="flex-center">
            <a-input-number v-model:value="form.reply_interval" style="width: 60px;" />
            <span class="ml4">秒 在时间间隔内只会触发一次该回复，0秒为无时间间隔限制</span>
          </div>
        </a-form-item>

        <!-- 消息类型：单选 全部消息、指定消息 -->
        

        <!-- 回复内容 -->
        <div class="nav-box" style="margin-top: 24px;">
          <svg-icon name="reply-content" style="font-size: 16px;"></svg-icon>
          回复内容
        </div>
        <div class="reply-content-box">
          <MultiReply v-for="(it, idx) in replyList" :key="idx" ref="replyRefs" v-model:value="replyList[idx]" :reply_index="idx"
            @change="onContentChange" @del="onDelItem" />
          <a-button type="dashed" style="width: 694px;" :disabled="replyList.length >= 5" @click="addReplyItem">
            <template #icon>
              <PlusOutlined />
            </template>
            添加回复内容({{replyList.length}}/5)
          </a-button>
        </div>

        <!-- 回复方式 -->
        <a-form-item name="reply_num" label="回复方式" style="margin-top: 24px;">
          <div class="radio-container">
            <a-radio-group v-model:value="form.reply_num" @change="handleReplyTypeChange">
              <a-radio value="0">全部回复</a-radio>
              <a-radio value="1">随机回复一条</a-radio>
            </a-radio-group>
          </div>
        </a-form-item>

        <!-- 保存 底部固定 -->
        <div class="btn-container">
          <a-button type="primary" @click="onSubmit">保存</a-button>
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
import { getRobotSubscribeReply, saveRobotSubscribeReply } from '@/api/explore/index.js'
import { SUBSCRIBE_SOURCE_OPTIONS } from '@/constants/index.js'
import dayjs from 'dayjs'

const replyRefs = ref([])
const query = useRoute().query
const ruleId = ref(+query.rule_id || +query['rule-id'] || 0)
const ruleType = ref(query.subscribe_rule_type || 'subscribe_reply_duration')
const router = useRouter()
const autoReplyUrl = computed(() => `/#/explore/index/subscribe-reply?subscribe_rule_type=${ruleType.value}`)
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
  switch_status: true,
  subscribe_source: []
})
const rules = {
  name: [
    { required: true, message: '请输入规则名称', trigger: 'blur' }
  ]
}
// 关键词集合
const messageTypeOptions = [
  { label: '文本', value: 'text' },
  { label: '图片', value: 'image' },
  { label: '音频', value: 'audio' },
  { label: '视频', value: 'video' }
]
const subscribeSourceOptions = SUBSCRIBE_SOURCE_OPTIONS()

function validateDateRange (_rule, value) {
  if (form.duration_type !== 'time_range') return Promise.resolve()
  const ok = Array.isArray(value) && value.length === 2 && value[0] && value[1]
  return ok ? Promise.resolve() : Promise.reject('请选择起止日期')
}

function validateMessageTypeList (_rule, value) {
  if (form.message_type !== '1') return Promise.resolve()
  const ok = Array.isArray(value) && value.length > 0
  return ok ? Promise.resolve() : Promise.reject('请至少选择一种指定消息类型')
}

function validateWeekDay (_rule, value) {
  if (form.duration_type !== 'week') return Promise.resolve()
  const ok = Array.isArray(value) && value.length > 0
  return ok ? Promise.resolve() : Promise.reject('请选择至少一个星期')
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
    message.warning('最多添加5条回复内容')
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
      message.warning('请至少添加一条回复内容')
      return
    }

    if (ruleType.value === 'subscribe_reply_duration') {
      if (!form.reply_period || form.reply_period.length !== 2) {
        message.warning('请选择回复时段')
        return
      }
      const p = Number(form.priority_num)
      if (!Number.isInteger(p) || p < 1 || p > 100) {
        message.warning('优先级需为1-100之间的整数')
        return
      }
      if (form.duration_type === 'time_range') {
        const ok = Array.isArray(form.date_range) && form.date_range.length === 2 && form.date_range[0] && form.date_range[1]
        if (!ok) {
          message.warning('请选择自定义日期范围')
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
      appid: query.appid || '',
      switch_status: form.switch_status ? '1' : '0',
      rule_type: ruleType.value,
      reply_content: JSON.stringify(serializeReplyContent(replyList.value)),
      reply_type: serializeReplyTypeCodes(replyList.value),
      reply_num: form.reply_num
    }

    let extra = {}
    if (ruleType.value === 'subscribe_reply_duration') {
      extra = {
        duration_type: form.duration_type,
        start_day: form.start_day,
        end_day: form.end_day,
        reply_interval: Number(form.reply_interval) || 0,
        week_duration: form.week_duration,
        priority_num: Number(form.priority_num) || 0,
        start_duration: Array.isArray(form.reply_period) ? (form.reply_period[0].toString() || '') : '',
        end_duration: Array.isArray(form.reply_period) ? (form.reply_period[1].toString() || '') : ''
      }
    } else {
      extra = {
        subscribe_source: Array.isArray(form.subscribe_source) ? form.subscribe_source.join(',') : (form.subscribe_source || '')
      }
    }

    const payload = { ...base, ...extra }
    if (ruleId.value) payload.id = ruleId.value
    try {
      const res = await saveRobotSubscribeReply(payload)
      if (res && res.res == 0) {
        message.success('保存成功')
        router.push({ path: '/explore/index/subscribe-reply', query: { subscribe_rule_type: ruleType.value } })
      }
    } catch (e) {
    }
  })
}


onMounted(async () => {
  const copyId = +(query.copy_id || 0)
  const fetchOne = async (rid) => {
    const res = await getRobotSubscribeReply({ id: rid })
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
    const srcRaw = data?.subscribe_source
    if (Array.isArray(srcRaw)) {
      form.subscribe_source = srcRaw
    } else if (typeof srcRaw === 'string') {
      try {
        const arr = JSON.parse(srcRaw)
        form.subscribe_source = Array.isArray(arr) ? arr : (srcRaw ? srcRaw.split(',').filter(Boolean) : [])
      } catch (_e) {
        form.subscribe_source = srcRaw ? srcRaw.split(',').filter(Boolean) : []
      }
    }

    // message type
    
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
  right: 0;
  display: flex;
  width: 100%;
  padding: 16px 32px;
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
  width: 98px;
  min-width: 98px;
  flex: 0 0 98px;
}

.interval-item ::v-deep(.ant-form-item-label) {
  width: 130px;
  min-width: 130px;
  flex: 0 0 130px;
}

.message-type-item {
  margin-bottom: 0;
}
.item-box ::v-deep(.ant-form-item) {
  margin-bottom: 0;
}
</style>