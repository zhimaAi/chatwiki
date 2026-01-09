<template>
  <div class="main-box">
    <div v-if="loadingF || loading" class="loading-box">
      <a-spin/>
    </div>
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>输入</div>
      </div>
      <div
        v-for="item in displayParamList"
        :key="item.key"
        class="options-item"
        :class="{'is-required': item.meta?.required}"
      >
      <div class="options-item-tit">
          <div class="option-label">{{ item.meta?.name || item.key }}
            <a-tooltip title="同步最新的标签" v-if="item.meta?.tag_component" :overlayStyle="{ maxWidth: '800px' }">
              <a @click="syncFieldTags(item.key, item.meta)">同步 <a-spin v-if="syncingTags[item.key]" size="small" /></a>
            </a-tooltip>
            <a-tooltip v-if="item.meta?.tip" :overlayStyle="{ maxWidth: '800px' }">
              <template #title>
                <div class="tip-content">{{ item.meta?.tip }}</div>
              </template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </div>
          <div class="option-type">{{ item.meta?.type || 'string' }}</div>
        </div>
        <div v-if="item.meta?.select_official_component">
          <a-select
            v-model:value="formState.app_id"
            placeholder="请选择公众号"
            style="width: 100%;"
            @change="(value, option) => appChange(value, option)"
          >
            <a-select-option
              v-for="app in apps"
              :key="app.app_id"
              :value="app.app_id"
              :name="app.app_name"
              :app_secret="app.app_secret"
            >
              {{ app.app_name }}
            </a-select-option>
          </a-select>
          <div class="desc">{{ item.meta?.desc }}</div>
          <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
        </div>
        <div v-else-if="item.meta?.enum_component">
          <a-select
            v-model:value="formState[item.key]"
            :placeholder="item.meta?.required ? '请选择' : '请选择（非必填）'"
            style="width: 100%;"
            @change="update"
          >
            <a-select-option
              v-for="opt in (sortedEnum(item.meta?.enum) || [])"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.name }}
            </a-select-option>
          </a-select>
          <div class="desc">{{ item.meta?.desc }}</div>
          <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
        </div>
        <div v-else-if="item.meta?.radio_component">
          <a-radio-group
            v-model:value="formState[item.key]"
            @change="update"
          >
            <a-radio
              v-for="opt in (sortedEnum(item.meta?.enum) || [])"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.name }}
            </a-radio>
          </a-radio-group>
          <div class="desc">{{ item.meta?.desc }}</div>
          <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
        </div>
        <div v-else-if="item.meta?.date_range_begin_component">
          <div class="tag-box">
            <a-select
              v-model:value="dateRangeTypes[item.meta?.begin_date_key]"
              @change="(val) => dateRangeTypeChange(item.meta, val)"
              style="width: 120px; padding-right: 4px;"
            >
              <a-select-option :value="1">选择日期</a-select-option>
              <a-select-option :value="2">插入变量</a-select-option>
            </a-select>
            <a-range-picker
              v-if="dateRangeTypes[item.meta?.begin_date_key] == 1"
              :value="getRangeValue(item.meta)"
              format="YYYY-MM-DD"
              style="width: 100%;"
              @calendarChange="(dates) => onDateRangeCalendarChange(item.meta, dates)"
              @change="(dates, dateStrings) => onDateRangeChange(item.meta, dates, dateStrings)"
              :disabled-date="(current) => getDisabledDateForRange(item.meta, current)"
            />
            <template v-else>
              <AtInput
                inputStyle="height: 33px; width: 214px;"
                :options="variableOptions"
                :defaultSelectedList="formState[(item.meta?.begin_date_key || '') + '_tags'] || []"
                :defaultValue="formState[item.meta?.begin_date_key] || ''"
                @open="emit('updateVar')"
                @change="(text, selectedList) => onChangeDynamic(item.meta?.begin_date_key, text, selectedList)"
                placeholder="请输入开始日期，键入“/”可以插入变量"
              >
                <template #option="{ label, payload }">
                  <div class="field-list-item">
                    <div class="field-label">{{ label }}</div>
                    <div class="field-type">{{ payload.typ }}</div>
                  </div>
                </template>
              </AtInput>
              <div class="date-range-separator">-</div>
              <AtInput
                inputStyle="height: 33px; width: 214px;"
                :options="variableOptions"
                :defaultSelectedList="formState[(item.meta?.end_date_key || '') + '_tags'] || []"
                :defaultValue="formState[item.meta?.end_date_key] || ''"
                @open="emit('updateVar')"
                @change="(text, selectedList) => onChangeDynamic(item.meta?.end_date_key, text, selectedList)"
                placeholder="请输入结束日期，键入“/”可以插入变量"
              >
                <template #option="{ label, payload }">
                  <div class="field-list-item">
                    <div class="field-label">{{ label }}</div>
                    <div class="field-type">{{ payload.typ }}</div>
                  </div>
                </template>
              </AtInput>
            </template>
          </div>
          <div class="desc">{{ item.meta?.desc }}</div>
          <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
        </div>
        <div v-else-if="item.meta?.tag_component && showTagComponents">
          <div class="tag-box">
            <a-select
              v-model:value="tagTypes[item.key]"
              @change="(val) => tagTypeChange(item.key, val, item.meta)"
              style="width: 120px;"
            >
              <a-select-option :value="1">选择标签</a-select-option>
              <a-select-option :value="2">插入变量</a-select-option>
            </a-select>
            <a-select
              v-if="tagTypes[item.key] == 1"
              v-model:value="formState[item.key]"
              placeholder="请选择标签"
              style="width: 100%;"
              @change="(val, opt) => tagChange(item.key, val, opt)"
              show-search
              :filter-option="filterOption"
            >
              <a-select-option
                v-for="t in (tagLists[item.key] || [])"
                :key="t.id"
                :value="t.id"
                :name="t.name"
              >
                {{ t.name }}
              </a-select-option>
            </a-select>
            <AtInput
              v-else
              type="textarea"
              inputStyle="height: 33px;"
              :options="variableOptions"
              :defaultSelectedList="formState[item.key + '_tags'] || []"
              :defaultValue="formState[item.key] || ''"
              @open="emit('updateVar')"
              @change="(text, selectedList) => onChangeDynamic(item.key, text, selectedList)"
              placeholder="请输入内容，键入“/”可以插入变量"
            >
              <template #option="{ label, payload }">
                <div class="field-list-item">
                  <div class="field-label">{{ label }}</div>
                  <div class="field-type">{{ payload.typ }}</div>
                </div>
              </template>
            </AtInput>
          </div>
          <div class="desc">{{ item.meta?.desc }}</div>
          <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
        </div>
        <div v-else>
          <AtInput
            :type="item.meta?.in === 'body' ? 'textarea' : undefined"
            :inputStyle="item.meta?.in === 'body' ? 'height: 64px;' : undefined"
            :options="variableOptions"
            :defaultSelectedList="formState[item.key + '_tags'] || []"
            :defaultValue="formState[item.key] || ''"
            @open="emit('updateVar')"
            @change="(text, selectedList) => onChangeDynamic(item.key, text, selectedList)"
            placeholder="请输入内容，键入“/”可以插入变量"
          >
            <template #option="{ label, payload }">
              <div class="field-list-item">
                <div class="field-label">{{ label }}</div>
                <div class="field-type">{{ payload.typ }}</div>
              </div>
            </template>
          </AtInput>
          <div class="desc">{{ item.meta?.desc }}</div>
          <div v-if="errors[item.key]" class="desc" style="color:#FB363F;">{{ errors[item.key] }}</div>
        </div>
      </div>
    </div>

    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/output.svg" class="title-icon"/>输出</div>
      </div>
      <div class="options-item">
        <OutputFields :tree-data="outputData"/>
      </div>
    </div>
  </div>
</template>

<script setup>
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {ref, reactive, onMounted, computed, inject, watch} from 'vue'
import OutputFields from "./output-fields.vue";
import {pluginOutputToTree} from "@/constants/plugin.js";
import {getWechatAppList} from "@/api/robot/index.js";
import { useEventBus } from '@/hooks/event/useEventBus.js'
import AtInput from '@/views/workflow/components/at-input/at-input.vue'
import { isValidURL } from '@/utils/validate.js'
import { runPlugin } from '@/api/plugins/index.js'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'

const emit = defineEmits(['updateVar'])
const props = defineProps({
  node: {
    type: Object,
    default: () => ({})
  },
  action: {
    type: Object,
  },
  actionName: {
    type: String,
  },
  variableOptions: {
    type: Array,
  }
})

const setData = inject('setData')
const paramsMap = computed(() => props.action?.params || {})
const hasOfficial = computed(() => !!(paramsMap.value.app_id && paramsMap.value.app_secret))
const paramList = computed(() => {
  const list = Object.entries(paramsMap.value || {}).map(([key, meta]) => ({ key, meta }))
  return list
    .filter(it => !it.meta?.hide_official_component)
    .sort((a, b) => (+(a.meta?.sort_num || 0)) - (+(b.meta?.sort_num || 0)))
})
function sortedEnum(list) {
  return (list || []).slice().sort((a, b) => (+(a.sort_num || 0)) - (+(b.sort_num || 0)))
}


const formState = reactive({
  app_id: '',
  app_secret: '',
  app_name: ''
})
const loading = ref(false)
const loadingF = ref(false)
const apps = ref([])
const errors = reactive({})

const outputData = ref([])

const tagLists = reactive({})
const tagTypes = reactive({})
const syncingTags = reactive({})
const dateRangeTypes = reactive({})
const dateAnchors = reactive({})

const showTagComponents = computed(() => {
  const hasTag = paramList.value.some((it) => it.meta?.tag_component)
  if (!hasTag) return false
  return paramList.value.some((it) => {
    const meta = it.meta || {}
    const enumList = Array.isArray(meta.enum) ? meta.enum : []
    if (!(meta.enum_component || meta.radio_component) || enumList.length === 0) return false
    const selected = String(formState[it.key] ?? '')
    const opt = enumList.find((e) => String(e?.value) === selected)
    const related = opt?.related_tags
    return related === true || related === 1 || related === 'true' || related === '1'
  })
})

const displayParamList = computed(() => {
  return (paramList.value || [])
    .filter((it) => !(it.meta?.tag_component && !showTagComponents.value))
})

watch(showTagComponents, (val) => {
  if (!val) {
    let changed = false
    paramList.value.forEach((it) => {
      if (it.meta?.tag_component) {
        if (formState[it.key]) {
          formState[it.key] = ''
          changed = true
        }
        if (Array.isArray(formState[it.key + '_tags']) && formState[it.key + '_tags'].length) {
          formState[it.key + '_tags'] = []
          changed = true
        }
      }
    })
    if (changed) update()
  }
})

onMounted(() => {
  init()
  const bus = useEventBus({
    name: 'workflow:validate',
    callback: onWorkflowValidate
  })
})

async function init() {
  if (hasOfficial.value) {
    await loadWxApps()
  }
  nodeParamsAssign()
  paramList.value.forEach(async it => {
    const def = String(it.meta?.enum_default || '').trim()
    if ((it.meta?.radio_component || it.meta?.enum_component) && !String(formState[it.key] || '').trim() && def) {
      formState[it.key] = def
    }
    if (it.meta?.tag_component) {
      tagTypes[it.key] = typeof tagTypes[it.key] === 'number' ? tagTypes[it.key] : 1
      await loadFieldTags(it.key, it.meta?.query_tag_service)
    }
  })
  outputFormat()
}

function loadWxApps() {
  return getWechatAppList({app_type: 'official_account'}).then(res => {
    apps.value = res?.data || []
  })
}

function nodeParamsAssign() {
  let nodeParams = JSON.parse(props.node.node_params)
  let arg = nodeParams?.plugin?.params?.arguments || {}
  paramList.value.forEach(it => {
    formState[it.key] = it.key !== 'app_id' ? String(arg[it.key] || '').trim() : arg[it.key]
    const tagMap = arg?.tag_map || {}
    formState[it.key + '_tags'] = Array.isArray(tagMap[it.key]) ? tagMap[it.key] : []
    if (it.meta?.date_range_begin_component) {
      const endKey = it.meta?.end_date_key
      const typeKey = it.meta?.date_type_key
      formState[endKey] = String(arg[endKey] || '').trim()
      formState[endKey + '_tags'] = Array.isArray(tagMap[endKey]) ? tagMap[endKey] : []
      const typeVal = arg[typeKey]
      dateRangeTypes[it.meta?.begin_date_key] =
        typeof dateRangeTypes[it.meta?.begin_date_key] === 'number'
          ? dateRangeTypes[it.meta?.begin_date_key]
          : (typeVal == null ? 1 : Number(typeVal) || 1)
    }
  })
  if (arg.app_secret) formState.app_secret = String(arg.app_secret)
  if (arg.app_name) formState.app_name = String(arg.app_name)
}

function appChange (value, option) {
  if (option && option.app_secret) {
    formState.app_secret = option.app_secret
  }
  if (option && option.name) {
    formState.app_name = option.name
  }
  update()
  paramList.value.forEach(async it => {
    if (it.meta?.tag_component) {
      await loadFieldTags(it.key, it.meta?.query_tag_service)
    }
  })
}

function onWorkflowValidate() {
  const bus = useEventBus()
  const res = validateAll()
  if (res && res.ok === false) {
    bus.emit && bus.emit('workflow:validate:error', {
      component: '公众号智能API',
      type: '',
      typeDisplay: '',
      message: '请完善字段后再继续',
      errors: res.errors || []
    })
  }
}

function update(val = null) {
  let nodeParams = JSON.parse(props.node.node_params)
  nodeParams.plugin.output_obj = outputData.value
  Object.assign(nodeParams.plugin.params, {
    arguments: {
      ...buildArguments(),
      tag_map: {
        ...buildTagMap()
      }
    },
    rendering: buildRendering()
  })
  setData({
    ...props.node,
    node_params: JSON.stringify(nodeParams)
  })
}

function onChangeDynamic (key, text, selectedList) {
  formState[key] = text
  formState[key + '_tags'] = selectedList
  update()
}

function buildRendering () {
  return paramList.value
    .filter((it) => !(it.meta?.tag_component && !showTagComponents.value))
    .map((it) => {
    let label = it.meta?.name || it.key
    let value = String(formState[it.key] ?? '')

    const tags = Array.isArray(formState[it.key + '_tags']) ? formState[it.key + '_tags'] : []
    if (it.key === 'app_id') {
      label = '公众号名称'
      value = String(formState.app_name ?? '')
    }
    if (it.meta?.enum_component || it.meta?.radio_component) {
      const opt = (it.meta?.enum || []).find((e) => String(e?.value) === value)
      if (opt && opt.name) {
        value = String(opt.name)
      }
    }
    if (it.meta?.tag_component && tagTypes[it.key] == 1) {
      const opt = (tagLists[it.key] || []).find((t) => String(t.id) === String(value))
      if (opt && opt.name) {
        value = String(opt.name)
      }
    }
    if (it.meta?.date_range_begin_component) {
      const beginKey = it.meta?.begin_date_key
      const endKey = it.meta?.end_date_key
      const beginVal = String(formState[beginKey] || '')
      const endVal = String(formState[endKey] || '')
      const beginTags = Array.isArray(formState[beginKey + '_tags']) ? formState[beginKey + '_tags'] : []
      const endTags = Array.isArray(formState[endKey + '_tags']) ? formState[endKey + '_tags'] : []
      value = [beginVal, endVal].filter(Boolean).join(' ~ ')
      const mergedTags = [...beginTags, ...endTags]
      return {
        label,
        value,
        key: it.key,
        tags: mergedTags
      }
    }
    return {
      label,
      value,
      key: it.key,
      tags
    }
  })
}

function buildTagMap() {
  const tagMap = {}
  paramList.value.forEach((it) => {
    if (it.meta?.tag_component && !showTagComponents.value) return
    const tags = Array.isArray(formState[it.key + '_tags']) ? formState[it.key + '_tags'] : []
    if (tags.length) {
      tagMap[it.key] = tags
    }
    if (it.meta?.date_range_begin_component) {
      const endKey = it.meta?.end_date_key
      const endTags = Array.isArray(formState[endKey + '_tags']) ? formState[endKey + '_tags'] : []
      if (endTags.length) {
        tagMap[endKey] = endTags
      }
    }
  })
  return tagMap
}

function buildArguments() {
  const args = {}
  paramList.value.forEach(it => {
    if (it.meta?.tag_component && !showTagComponents.value) return
    args[it.key] = String(formState[it.key] || '')
    if (it.meta?.date_range_begin_component) {
      const endKey = it.meta?.end_date_key
      const typeKey = it.meta?.date_type_key
      args[endKey] = String(formState[endKey] || '')
      const typeVal = dateRangeTypes[it.meta?.begin_date_key] || 1
      args[typeKey] = String(typeVal)
    }
  })
  if (formState.app_secret) {
    args.app_secret = String(formState.app_secret || '')
  }
  if (formState.app_name) {
    args.app_name = String(formState.app_name || '')
  }
  return args
}

function outputFormat() {
  outputData.value = pluginOutputToTree(props.action.output)
  update()
}

function validateAll () {
  Object.keys(errors).forEach(k => { errors[k] = '' })
  const result = { ok: true, errors: [] }
  paramList.value.forEach(it => {
    const key = it.key
    const meta = it.meta || {}
    const raw = formState[key]
    const val = String(raw == null ? '' : raw).trim()
    const tags = Array.isArray(formState[key + '_tags']) ? formState[key + '_tags'] : []
    if (meta.select_official_component) {
      if (meta.required && !val) {
        const msg = meta.error || '请选择公众号'
        errors[key] = msg
        result.ok = false
        result.errors.push({ field_name: meta.name || key, message: msg })
      }
      return
    }
    if (meta.enum_component || meta.radio_component) {
      const enums = (meta.enum || []).map(it => String(it.value))
      if (meta.required && !val) {
        const msg = meta.error || '请选择'
        errors[key] = msg
        result.ok = false
        result.errors.push({ field_name: meta.name || key, message: msg })
      } else if (val && enums.length && !enums.includes(val)) {
        const msg = meta.error || '选择不合法'
        errors[key] = msg
        result.ok = false
        result.errors.push({ field_name: meta.name || key, message: msg })
      }
      return
    }
    if (meta.date_range_begin_component) {
      const beginKey = meta.begin_date_key
      const endKey = meta.end_date_key
      const beginVal = String(formState[beginKey] || '').trim()
      const endVal = String(formState[endKey] || '').trim()
      const beginTags = Array.isArray(formState[beginKey + '_tags']) ? formState[beginKey + '_tags'] : []
      const endTags = Array.isArray(formState[endKey + '_tags']) ? formState[endKey + '_tags'] : []
      if (meta.required && (!beginVal && beginTags.length === 0 || !endVal && endTags.length === 0)) {
        const msg = meta.error || '请选择起止时间'
        errors[key] = msg
        result.ok = false
        result.errors.push({ field_name: meta.name || key, message: msg })
        return
      }
      const typeVal = dateRangeTypes[beginKey]
      if (typeVal === 2) {
        // 插入变量模式：仅校验纯文本输入的格式为 YYYY-MM-DD；引用变量不校验格式
        const isValidDateStr = (str) => {
          if (!/^\d{4}-\d{2}-\d{2}$/.test(str)) return false
          const d = dayjs(str)
          return d.isValid()
        }
        if (beginVal && beginTags.length === 0 && !isValidDateStr(beginVal)) {
          const msg = meta.error || '开始日期格式需为YYYY-MM-DD'
          errors[key] = msg
          result.ok = false
          result.errors.push({ field_name: meta.name || key, message: msg })
          return
        }
        if (endVal && endTags.length === 0 && !isValidDateStr(endVal)) {
          const msg = meta.error || '结束日期格式需为YYYY-MM-DD'
          errors[key] = msg
          result.ok = false
          result.errors.push({ field_name: meta.name || key, message: msg })
          return
        }
      }
      if (beginVal && endVal) {
        const b = dayjs(beginVal)
        const e = dayjs(endVal)
        if (e.isBefore(b)) {
          const msg = meta.error || '结束时间不能早于开始时间'
          errors[key] = msg
          result.ok = false
          result.errors.push({ field_name: meta.name || key, message: msg })
          return
        }
        const interval = e.diff(b, 'day')
        const allowedDays = Number(meta.date_range_max_interval || 0)
        if (allowedDays > 0 && (interval + 1) > allowedDays) {
          const msg = meta.error || `时间跨度不能超过${allowedDays}天`
          errors[key] = msg
          result.ok = false
          result.errors.push({ field_name: meta.name || key, message: msg })
          return
        }
      }
      return
    }
    if (meta.required && !val && (!Array.isArray(tags) || tags.length === 0)) {
      const msg = meta.error || '请输入内容'
      errors[key] = msg
      result.ok = false
      result.errors.push({ field_name: meta.name || key, message: msg })
    } else if (key.toLowerCase().includes('url') && val && (!Array.isArray(tags) || tags.length === 0)) {
      const ok = /^https?:\/\//.test(val) || isValidURL(val)
      if (!ok) {
        const msg = meta.error || '请输入正确的URL'
        errors[key] = msg
        result.ok = false
        result.errors.push({ field_name: meta.name || key, message: msg })
      }
    }
  })
  return result
}

function tagTypeChange(fieldKey, val, meta) {
  tagTypes[fieldKey] = val
  if (val == 1) {
    loadFieldTags(fieldKey, meta?.query_tag_service)
  }
}

function tagChange(fieldKey, val, opt) {
  formState[fieldKey] = val
  formState[fieldKey + '_name'] = opt?.name || ''
  update()
}

function loadFieldTags(fieldKey, service) {
  const nodeParams = JSON.parse(props.node.node_params)
  const pluginName = nodeParams?.plugin?.name
  if (!pluginName || !service) {
    tagLists[fieldKey] = []
    return Promise.resolve([])
  }
  if (!formState.app_id || !formState.app_secret) {
    tagLists[fieldKey] = []
    return Promise.resolve([])
  }
  return runPlugin({
    name: pluginName,
    action: 'default/exec',
    params: JSON.stringify({
      business: service,
      arguments: {
        app_id: formState.app_id,
        app_secret: formState.app_secret
      }
    })
  }).then((res) => {
    tagLists[fieldKey] = res?.data?.tags || []
    return res
  })
}

function syncFieldTags(fieldKey, meta) {
  if (!formState.app_id || !formState.app_secret) {
    return message.warning('请先选择公众号')
  }
  if (syncingTags[fieldKey]) return
  syncingTags[fieldKey] = true
  loadFieldTags(fieldKey, meta?.query_tag_service)
    .then(() => {
      message.success('同步完成')
    })
    .finally(() => {
      syncingTags[fieldKey] = false
    })
}

const filterOption = (input, option) => {
  return option.name.toLowerCase().indexOf(input.toLowerCase()) >= 0
}

function dateRangeTypeChange(meta, val) {
  dateRangeTypes[meta.begin_date_key] = val
  const typeKey = meta?.date_type_key
  if (typeKey) {
    formState[typeKey] = String(val)
  }
  const beginKey = meta.begin_date_key
  const endKey = meta.end_date_key
  if (val === 1) {
    formState[beginKey] = ''
    formState[endKey] = ''
    formState[beginKey + '_tags'] = []
    formState[endKey + '_tags'] = []
    dateAnchors[beginKey] = null
  } else if (val === 2) {
    formState[beginKey] = ''
    formState[endKey] = ''
    dateAnchors[beginKey] = null
  }
  update()
}

function getRangeValue(meta) {
  const b = formState[meta.begin_date_key]
  const e = formState[meta.end_date_key]
  return [b ? dayjs(b) : null, e ? dayjs(e) : null]
}

function onDateRangeCalendarChange(meta, dates) {
  const key = meta.begin_date_key
  if (Array.isArray(dates) && dates[0]) {
    dateAnchors[key] = dates[0]
  } else {
    dateAnchors[key] = null
  }
}

function onDateRangeChange(meta, dates, dateStrings) {
  formState[meta.begin_date_key] = dateStrings[0] || ''
  formState[meta.end_date_key] = dateStrings[1] || ''
  if ((!dateStrings[0] && !dateStrings[1])) {
    dateAnchors[meta.begin_date_key] = null
  }
  update()
}

function getDisabledDateForRange(meta, current) {
  if (!current) return false
  const endMax = meta.end_date_max ? dayjs(meta.end_date_max) : null
  if (endMax && current.isAfter(endMax, 'day')) return true
  const anchor = dateAnchors[meta.begin_date_key]
  const maxInterval = Number(meta.date_range_max_interval || 0)
  if (anchor && maxInterval > 0) {
    const limit = dayjs(anchor).add(maxInterval - 1, 'day')
    return current.isAfter(limit, 'day')
  }
  return false
}
</script>

<style scoped lang="less">
.loading-box {
  width: 100%;
  height: 100%;
  position: absolute;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 99;
}
.node-options {
  background: #f2f4f7;
  border-radius: 6px;
  padding: 12px;
  margin-top: 16px;

  &:first-child {
    margin-top: 0;
  }

  .options-title {
    color: var(--wf-color-text-1);
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-weight: 600;
    height: 22px;;
    line-height: 22px;
    font-size: 14px;
    margin-bottom: 16px;

    .title-icon {
      width: 16px;
      height: 16px;
      vertical-align: -3px;
      margin-right: 8px;;
    }

    .acton-box {
      font-weight: 400;
    }
  }

.options-item {
  display: flex;
  flex-direction: column;
  margin-top: 12px;
  line-height: 22px;
  gap: 4px;

    .options-item-tit {
      display: flex;
      align-items: center;
    }

    .option-label {
      color: var(--wf-color-text-1);
      font-size: 14px;
      margin-right: 8px;
    }

    .desc {
      color: var(--wf-color-text-2);
    }


    &.is-required .option-label::before {
      content: '*';
      color: #FB363F;
      display: inline-block;
      margin-right: 2px;
    }

    .option-type {
      height: 22px;
      line-height: 18px;
      padding: 0 8px;
      border-radius: 6px;
      border: 1px solid rgba(0, 0, 0, 0.15);
      background-color: #fff;
      color: var(--wf-color-text-3);
      font-size: 12px;
    }

    .item-actions-box {
      display: flex;
      align-items: center;

      .action-btn {
        margin-left: 12px;
        font-size: 16px;
        color: #595959;
        cursor: pointer;
      }
    }
  }
}
.tag-box {
  display: flex;
  align-items: center;
  :deep(.mention-input-warpper) {
    height: 33px;
  }
}

.date-range-separator {
  padding: 0 6px;
}

.tip-content {
  white-space: pre-wrap;
}
</style>
