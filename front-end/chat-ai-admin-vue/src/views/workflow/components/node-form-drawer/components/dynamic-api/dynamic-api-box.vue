<template>
  <div class="main-box">
    <div v-if="loadingF || loading" class="loading-box">
      <a-spin/>
    </div>
    <div class="node-options" v-if="hasCacheSetting">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/setting-icon.svg" class="title-icon"/>{{ t('title_settings') }}</div>
        <div class="acton-box" v-if="clickUrl">
          <a :href="clickUrl" target="_blank" rel="noopener noreferrer">{{ clickTitle || clickUrl }}</a>
        </div>
      </div>
      <FieldList
        :items="cacheParamList"
        :formState="formState"
        :errors="errors"
        :apps="apps"
        :variableOptions="variableOptions"
        :tagTypes="tagTypes"
        :tagLists="tagLists"
        :showTagComponents="showTagComponents"
        :sortedEnum="sortedEnum"
        :appChange="appChange"
        :onFieldChange="onFieldChange"
        :update="update"
        :onChangeDynamic="onChangeDynamic"
        :filterOption="filterOption"
        :tagTypeChange="tagTypeChange"
        :tagChange="tagChange"
        :syncFieldTags="syncFieldTags"
        :syncingTags="syncingTags"
        :dateRangeTypes="dateRangeTypes"
        :dateRangeTypeChange="dateRangeTypeChange"
        :getRangeValue="getRangeValue"
        :onDateRangeCalendarChange="onDateRangeCalendarChange"
        :onDateRangeChange="onDateRangeChange"
        :getDisabledDateForRange="getDisabledDateForRange"
        :onToggleAdvanced="settingsShowChange"
        :onOpenVar="() => emit('updateVar')"
      />
    </div>
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>{{ t('title_input') }}</div>
        <template v-if="hasPluginConfigComponent">
          <a-select
            v-if="Object.keys(configData).length"
            v-model:value="currentConfigName"
            @change="configChange"
            style="min-width: 160px"
          >
            <a-select-option v-for="(item, name) in configData" :key="name" :value="name">
              {{ item.name || name }}
            </a-select-option>
          </a-select>
          <a-button v-else @click="showConfigModal">{{ t('btn_unauthorized') }} <DownOutlined/></a-button>
        </template>
      </div>
      <FieldList
        :items="displayParamList"
        :formState="formState"
        :errors="errors"
        :apps="apps"
        :variableOptions="variableOptions"
        :tagTypes="tagTypes"
        :tagLists="tagLists"
        :showTagComponents="showTagComponents"
        :sortedEnum="sortedEnum"
        :appChange="appChange"
        :onFieldChange="onFieldChange"
        :update="update"
        :onChangeDynamic="onChangeDynamic"
        :filterOption="filterOption"
        :tagTypeChange="tagTypeChange"
        :tagChange="tagChange"
        :syncFieldTags="syncFieldTags"
        :syncingTags="syncingTags"
        :dateRangeTypes="dateRangeTypes"
        :dateRangeTypeChange="dateRangeTypeChange"
        :getRangeValue="getRangeValue"
        :onDateRangeCalendarChange="onDateRangeCalendarChange"
        :onDateRangeChange="onDateRangeChange"
        :getDisabledDateForRange="getDisabledDateForRange"
        :onToggleAdvanced="settingsShowChange"
        :onOpenVar="() => emit('updateVar')"
      />
      <template v-if="hasAdvancedSetting && settingsOpen">
        <FieldList
          :items="advancedParamList"
          :formState="formState"
          :errors="errors"
          :apps="apps"
          :variableOptions="variableOptions"
          :tagTypes="tagTypes"
          :tagLists="tagLists"
          :showTagComponents="showTagComponents"
          :sortedEnum="sortedEnum"
          :appChange="appChange"
          :onFieldChange="onFieldChange"
          :update="update"
          :onChangeDynamic="onChangeDynamic"
          :filterOption="filterOption"
          :tagTypeChange="tagTypeChange"
          :tagChange="tagChange"
          :syncFieldTags="syncFieldTags"
          :syncingTags="syncingTags"
          :dateRangeTypes="dateRangeTypes"
          :dateRangeTypeChange="dateRangeTypeChange"
          :getRangeValue="getRangeValue"
          :onDateRangeCalendarChange="onDateRangeCalendarChange"
          :onDateRangeChange="onDateRangeChange"
          :getDisabledDateForRange="getDisabledDateForRange"
          :onToggleAdvanced="settingsShowChange"
          :onOpenVar="() => emit('updateVar')"
        />
      </template>
    </div>

    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/output.svg" class="title-icon"/>{{ t('title_output') }}</div>
      </div>
      <div class="options-item">
        <OutputFields :tree-data="outputData"/>
      </div>
    </div>

    <PluginConfigModal ref="configModalRef" @change="loadConfig(true)"/>
  </div>
</template>

<script setup>
import {ref, reactive, onMounted, computed, inject, watch} from 'vue'
import { DownOutlined } from '@ant-design/icons-vue'
import OutputFields from "./output-fields.vue";
import FieldList from "./field-list.vue";
import {pluginOutputToTree, getPluginConfigData} from "@/constants/plugin.js";
import {getWechatAppList} from "@/api/robot/index.js";
import { useEventBus } from '@/hooks/event/useEventBus.js'
import { isValidURL } from '@/utils/validate.js'
import { runPlugin } from '@/api/plugins/index.js'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'
import { useRoute } from 'vue-router'
import PluginConfigModal from "@/views/explore/plugins/components/plugin-config-modal.vue";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.components.dynamic-api.dynamic-api-box')

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
const route = useRoute()
let pluginNameForCache = ''
const hasCacheSetting = computed(() => {
  return (paramList.value || []).some((it) => {
    const v = it.meta?.setting_cache_component
    return v === true || v === 1 || v === 'true' || v === '1'
  })
})
const clickUrl = computed(() => String(formState.click_url || '').trim())
const clickTitle = computed(() => String(formState.click_title || '').trim())

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
  app_name: '',
  extras: {}
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

const configModalRef = ref(null)
const configData = ref({})
const currentConfigName = ref(null)
const settingsOpen = ref(false)

const showTagComponents = computed(() => {
  const hasTag = paramList.value.some((it) => it.meta?.tag_component)
  if (!hasTag) return false
  return paramList.value.some((it) => {
    const meta = it.meta || {}
    const enumList = Array.isArray(meta.enum) ? meta.enum : []
    if (!(meta.enum_component || meta.radio_component || meta.switch_component) || enumList.length === 0) return false
    const selected = String(formState[it.key] ?? '')
    const opt = enumList.find((e) => String(e?.value) === selected)
    const related = opt?.related_tags
    return related === true || related === 1 || related === 'true' || related === '1'
  })
})

const hasPluginConfigComponent = computed(() => {
  return (paramList.value || []).some((it) => {
    const v = it.meta?.plugin_config_component
    return v === true || v === 1 || v === 'true' || v === '1'
  })
})
const hasAdvancedSetting = computed(() => {
  const useAdv = props.action?.use_advanced_settings
  const anyAdvField = (paramList.value || []).some((it) => {
    const v = it.meta?.advanced_settings
    return v === true || v === 1 || v === 'true' || v === '1'
  })
  return (useAdv === true || useAdv === 1 || useAdv === 'true' || useAdv === '1') || anyAdvField
})

const displayParamList = computed(() => {
  const isCacheField = (meta) => {
    const v = meta?.setting_cache_component
    return v === true || v === 1 || v === 'true' || v === '1'
  }
  return (paramList.value || [])
    .filter((it) => !(it.meta?.tag_component && !showTagComponents.value))
    .filter((it) => !(it.meta?.plugin_config_component))
    .filter((it) => {
      const v = it.meta?.advanced_settings
      const isAdv = v === true || v === 1 || v === 'true' || v === '1'
      return !isAdv
    })
    .filter((it) => !isCacheField(it.meta))
})
const cacheParamList = computed(() => {
  return (paramList.value || []).filter((it) => {
    const v = it.meta?.setting_cache_component
    return v === true || v === 1 || v === 'true' || v === '1'
  })
})
const advancedParamList = computed(() => {
  return (paramList.value || []).filter((it) => {
    const v = it.meta?.advanced_settings
    const isAdv = v === true || v === 1 || v === 'true' || v === '1'
    const pc = it.meta?.plugin_config_component
    const isPc = pc === true || pc === 1 || pc === 'true' || pc === '1'
    return isAdv && !isPc
  })
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
  if (hasPluginConfigComponent.value) {
    await loadConfig()
  }
  paramList.value.forEach(async it => {
    const def = String(it.meta?.enum_default || '').trim()
    if ((it.meta?.radio_component || it.meta?.enum_component || it.meta?.switch_component) && !String(formState[it.key] || '').trim() && def) {
      formState[it.key] = def
    }
    if (it.meta?.tag_component) {
      tagTypes[it.key] = typeof tagTypes[it.key] === 'number' ? tagTypes[it.key] : 1
      await loadFieldTags(it.key, it.meta?.query_tag_service)
    }
  })
  loadCacheForEmptyFields()
  outputFormat()
}

function loadWxApps() {
  return getWechatAppList({app_type: 'official_account'}).then(res => {
    apps.value = res?.data || []
  })
}

function nodeParamsAssign() {
  let nodeParams = JSON.parse(props.node.node_params)
  pluginNameForCache = nodeParams?.plugin?.name || ''
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
  if (arg.extras) formState.extras = arg.extras
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
  console.log('node', props.node)
  const bus = useEventBus()
  const res = validateAll()
  if (res && res.ok === false) {
    bus.emit && bus.emit('workflow:validate:error', {
      component: props.node?.node_name,
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
  const args = buildArguments()
  injectConfigArgs(args)
  Object.assign(nodeParams.plugin.params, {
    arguments: {
      ...args,
      tag_map: {
        ...buildTagMap()
      }
    },
    rendering: buildRendering()
  })
  if (hasPluginConfigComponent.value) {
    nodeParams.plugin.params.config_name = currentConfigName.value || ''
  }
  nodeParams.plugin.params.use_plugin_config = !!hasPluginConfigComponent.value
  setData({
    ...props.node,
    node_params: JSON.stringify(nodeParams)
  })
  storeAllToCache()
}

function onChangeDynamic (key, text, selectedList, extras = {}) {
  formState[key] = text
  formState[key + '_tags'] = selectedList
  if (typeof formState.extras === "object") {
    Object.assign(formState.extras, extras)
  } else {
    formState.extras = extras
  }
  update()
}

function onFieldChange (key, value) {
  formState[key] = value
  update()
}

function sanitizeStr(s) {
  return String(s || '').replace(/`/g, '').trim()
}
function buildCacheKey(fieldKey) {
  const wfId = String(route.query.id || '')
  return `${wfId}:${pluginNameForCache}:${fieldKey}`
}
function loadCacheForEmptyFields() {
  if (!hasCacheSetting.value) return
  let url = ''
  let title = ''
  for (const it of (paramList.value || [])) {
    const v = it.meta?.setting_cache_component
    const enabled = v === true || v === 1 || v === 'true' || v === '1'
    if (!enabled) continue
    if (!url) {
      const u = sanitizeStr(it.meta?.click_url)
      const t = sanitizeStr(it.meta?.click_title)
      if (u) {
        url = u
        title = t
      }
    }
    const k = it.key
    const cur = String(formState[k] == null ? '' : formState[k]).trim()
    if (!cur) {
      const cache = localStorage.getItem(buildCacheKey(k))
      if (cache != null) {
        formState[k] = it.key === 'app_id' ? cache : String(cache)
      }
    }
  }
  if (url) formState.click_url = url
  if (title) formState.click_title = title
}
function storeAllToCache() {
  if (!hasCacheSetting.value) return
  paramList.value.forEach((it) => {
    const cfg = it.meta?.setting_cache_component
    const enabled = cfg === true || cfg === 1 || cfg === 'true' || cfg === '1'
    if (!enabled) return
    const k = it.key
    const val = String(formState[k] == null ? '' : formState[k])
    localStorage.setItem(buildCacheKey(k), val)
  })
}

function buildRendering () {
  return paramList.value
    .filter((it) => {
      if (it.meta?.tag_component && !showTagComponents.value) return false
      const v = it.meta?.setting_cache_component
      const isCache = v === true || v === 1 || v === 'true' || v === '1'
      if (isCache) return false
      if (it.meta?.plugin_config_component) return false
      if (it.meta?.advanced_settings) return false
      return true
    })
    .map((it) => {
    let label = it.meta?.name || it.key
    let value = String(formState[it.key] ?? '')

    const tags = Array.isArray(formState[it.key + '_tags']) ? formState[it.key + '_tags'] : []
    if (it.key === 'app_id') {
      label = t('label_official_account_name')
      value = String(formState.app_name ?? '')
    }
    if (it.meta?.select_remote_component && it.key === 'sheetIdOrName' && formState?.extras?.table_name) {
      value = String(formState.extras.table_name)
    }
    if (it.meta?.enum_component || it.meta?.radio_component || it.meta?.switch_component) {
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
    if (it.meta?.plugin_config_component) return
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
    if (it.meta?.plugin_config_component) return
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
  if (formState.extras) {
    args.extras = formState.extras
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
    if (meta.plugin_config_component) {
      return
    }
    if (meta.select_official_component) {
      if (meta.required && !val) {
        const msg = meta.error || t('msg_select_official_account')
        errors[key] = msg
        result.ok = false
        result.errors.push({ field_name: meta.name || key, message: msg })
      }
      return
    }
    if (meta.enum_component || meta.radio_component || meta.switch_component) {
      const enums = (meta.enum || []).map(it => String(it.value))
      if (meta.required && !val) {
        const msg = meta.error || t('msg_please_select')
        errors[key] = msg
        result.ok = false
        result.errors.push({ field_name: meta.name || key, message: msg })
      } else if (val && enums.length && !enums.includes(val)) {
        const msg = meta.error || t('msg_invalid_selection')
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
        const msg = meta.error || t('msg_select_date_range')
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
          const msg = meta.error || t('msg_start_date_format')
          errors[key] = msg
          result.ok = false
          result.errors.push({ field_name: meta.name || key, message: msg })
          return
        }
        if (endVal && endTags.length === 0 && !isValidDateStr(endVal)) {
          const msg = meta.error || t('msg_end_date_format')
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
          const msg = meta.error || t('msg_end_after_start')
          errors[key] = msg
          result.ok = false
          result.errors.push({ field_name: meta.name || key, message: msg })
          return
        }
        const interval = e.diff(b, 'day')
        const allowedDays = Number(meta.date_range_max_interval || 0)
        if (allowedDays > 0 && (interval + 1) > allowedDays) {
          const msg = meta.error || t('msg_date_range_limit', { days: allowedDays })
          errors[key] = msg
          result.ok = false
          result.errors.push({ field_name: meta.name || key, message: msg })
          return
        }
      }
      return
    }
    if (meta.required && !val && (!Array.isArray(tags) || tags.length === 0)) {
      const msg = meta.error || t('msg_please_input_content')
      errors[key] = msg
      result.ok = false
      result.errors.push({ field_name: meta.name || key, message: msg })
    } else if (key.toLowerCase().includes('url') && val && (!Array.isArray(tags) || tags.length === 0)) {
      const ok = /^https?:\/\//.test(val) || isValidURL(val)
      if (!ok) {
        const msg = meta.error || t('msg_correct_url')
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
    return message.warning(t('msg_select_official_account_first'))
  }
  if (syncingTags[fieldKey]) return
  syncingTags[fieldKey] = true
  loadFieldTags(fieldKey, meta?.query_tag_service)
    .then(() => {
      message.success(t('msg_sync_completed'))
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

function loadConfig(refresh=false) {
  return getPluginConfigData(pluginNameForCache, refresh).then(res => {
    configData.value = res || {}
    const names = Object.keys(configData.value)
    const defName = names.find((n) => configData.value[n]?.is_default)
    let wantName = null
    try {
      const nodeParams = JSON.parse(props.node.node_params)
      wantName = nodeParams?.plugin?.params?.config_name || ''
    } catch (e) {}
    if (wantName && names.includes(wantName)) {
      currentConfigName.value = wantName
    } else {
      currentConfigName.value = defName || names[0] || null
    }
  })
}
function showConfigModal() {
  const schemaData = { [props.actionName || 'default']: props.action || {} }
  configModalRef.value && configModalRef.value.show(configData.value, pluginNameForCache, schemaData)
}
function configChange() {
  let args = []
  injectConfigArgs(args)
  Object.assign(formState, args)
  update()
}
function settingsShowChange() {
  settingsOpen.value = !settingsOpen.value
}

function injectConfigArgs(args) {
  try {
    const cfgName = currentConfigName.value || ''
    const cfg = cfgName ? (configData.value?.[cfgName] || {}) : {}
    if (!cfg || Object.keys(cfg).length === 0) return
    paramList.value.forEach((it) => {
      const meta = it.meta || {}
      const isCfgField = meta.plugin_config_component === true || meta.plugin_config_component === 1 || meta.plugin_config_component === 'true' || meta.plugin_config_component === '1'
      if (!isCfgField) return
      const k = it.key
      const v = cfg[k]
      if (v != null) {
        args[k] = sanitizeStr(v)
      }
    })
  } catch (e) {}
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
