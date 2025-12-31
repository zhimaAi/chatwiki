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
        v-for="item in paramList"
        :key="item.key"
        class="options-item"
        :class="{'is-required': item.meta?.required}"
      >
        <div class="options-item-tit">
          <div class="option-label">{{ item.meta?.name || item.key }}
            <a-tooltip title="同步最新的标签" v-if="item.meta?.tag_component">
              <a @click="syncFieldTags(item.key, item.meta)">同步 <a-spin v-if="syncingTags[item.key]" size="small" /></a>
            </a-tooltip>
            <a-tooltip v-if="item.meta?.tip" :title="item.meta?.tip">
              <QuestionCircleOutlined />
            </a-tooltip>
          </div>
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
        <div v-else-if="item.meta?.tag_component">
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
import {ref, reactive, onMounted, computed, inject} from 'vue'
import OutputFields from "@/views/workflow/components/feishu-table/output-fields.vue";
import {pluginOutputToTree} from "@/constants/plugin.js";
import {getWechatAppList} from "@/api/robot/index.js";
import { useEventBus } from '@/hooks/event/useEventBus.js'
import AtInput from '@/views/workflow/components/at-input/at-input.vue'
import { isValidURL } from '@/utils/validate.js'
import { runPlugin } from '@/api/plugins/index.js'
import { message } from 'ant-design-vue'

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
  return paramList.value.map((it) => {
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
    const tags = Array.isArray(formState[it.key + '_tags']) ? formState[it.key + '_tags'] : []
    if (tags.length) {
      tagMap[it.key] = tags
    }
  })
  return tagMap
}

function buildArguments() {
  const args = {}
  paramList.value.forEach(it => {
    args[it.key] = String(formState[it.key] || '')
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
</style>
