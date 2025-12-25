<template>
  <div class="main-box">
    <div v-if="loadingF || loading" class="loading-box">
      <a-spin/>
    </div>
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>输入</div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">公众号</div>
        </div>
        <div>
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
        </div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">
            接收人
            <a-tooltip title="用户的openid">
              <QuestionCircleOutlined />
            </a-tooltip>
          </div>
        </div>
        <div>
          <AtInput
            type="textarea"
            inputStyle="height: 64px;"
            :options="variableOptions"
            :defaultSelectedList="formState.receiver_tags"
            :defaultValue="formState.receiver"
            @open="emit('updateVar')"
            @change="(text, selectedList) => onChangeReceiver(text, selectedList)"
            placeholder="请输入接收人，键入“/”可以插入变量"
          >
            <template #option="{ label, payload }">
              <div class="field-list-item">
                <div class="field-label">{{ label }}</div>
                <div class="field-type">{{ payload.typ }}</div>
              </div>
            </template>
          </AtInput>
        </div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">消息类型</div>
        </div>
        <div>
          <a-select
            v-model:value="formState.msgtype"
            @change="msgTypeChange"
            :placeholder="formState.app_id ? '请选择消息类型' : '请先选择公众号'"
            style="width: 100%;">
            <a-select-option
              v-for="item in tables"
              :key="item.msgtype"
              :value="item.msgtype">
              {{ item.name }}
            </a-select-option>
          </a-select>
        </div>
      </div>
      <ReplyData
        ref="childRef"
        :tableId="formState.msgtype"
        :fields="fields"
        :variableOptions="variableOptions"
        @update="update"
        @updateVar="emit('updateVar')"
      />
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
import {ref, reactive, onMounted, watch, computed, inject, nextTick} from 'vue'
import {runPlugin} from "@/api/plugins/index.js";
import ReplyData from "./reply-data.vue";
import OutputFields from "@/views/workflow/components/feishu-table/output-fields.vue";
import {getPluginActionDefaultArguments, pluginOutputToTree} from "@/constants/plugin.js";
import {getWechatAppList} from "@/api/robot/index.js";
import { useEventBus } from '@/hooks/event/useEventBus.js'
import AtInput from '@/views/workflow/components/at-input/at-input.vue'

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
const plginName = 'official_send_message'


const childRef = ref(null)
const formState = reactive({
  app_id: null,
  msgtype: undefined,
  receiver: '',
  receiver_tags: [],
  app_secret: '',
  app_name: ''
})
const loading = ref(false)
const loadingF = ref(false)
const tableData = ref({})
const apps = ref([])
const tables = ref([])
const fields = ref([])

const outputData = ref([])


onMounted(() => {
  init()
  const bus = useEventBus({
    name: 'workflow:validate',
    callback: onWorkflowValidate
  })
})

async function init() {
  await loadWxApps()
  nodeParamsAssign()
  outputFormat()
}

function loadWxApps() {
  return getWechatAppList({app_type: 'official_account'}).then(res => {
    apps.value = res?.data || []
  })
}

function loadTables() {
  loading.value = true
  return runPlugin({
    name: plginName,
    action: "default/exec",
    params: JSON.stringify({
      business: 'getMsgType',
      arguments: {
        app_id: formState.app_id
      }
    })
  }).then(res => {
    tableData.value = res?.data || {}
    tables.value = tableData.value.items || []
    if (!formState.msgtype && tables.value.length == 1) {
      formState.msgtype = tables.value[0].msgtype
      msgTypeChange()
    }
    return res
  }).finally(() => {
    loading.value = false
  })
}

function loadFields() {
  loadingF.value = true
  return runPlugin({
    name: plginName,
    action: "default/exec",
    params: JSON.stringify({
      business: 'getMsgInput',
      arguments: {
        app_id: formState.app_id,
        app_secret: formState.app_secret,
        msgtype: formState.msgtype
      }
    })
  }).then(res => {
    let items = res?.data?.items || []
    fields.value = items.map(it => ({
      field_name: it.field || it.name,
      ui_type: mapMsgFieldType(it.type),
      name: it.name,
      description: it.desc,
      required: !!it.required,
      media_component: !!it.media_component,
      properties: it.properties || null
    }))
    return res
  }).finally(() => {
    loadingF.value = false
  })
}

function mapMsgFieldType(t) {
  switch ((t || '').toLowerCase()) {
    case 'string': return 'Text'
    case 'number': return 'Number'
    case 'integer': return 'Number'
    case 'boolean': return 'Checkbox'
    case 'array<string>': return 'MultiSelect'
    default: return t || 'Text'
  }
}

function paramsFormat(val) {
  if (!val) {
    return getPluginActionDefaultArguments(props.actionName)
  }
  val = JSON.parse(JSON.stringify(val))
  for (let key in val) {
    switch (key) {
      case 'filter':
        if (Array.isArray(val.filter.conditions) && val.filter.conditions.length) {
          for (let item of val.filter.conditions) {
            if ('DateTime' === item.field_type && item.value) {
              if (item.value < 1e11) item.value *= 1000
              // 后端要求DateTime添加了特殊标记
              item.value = ['ExactDate', item.value].join(',')
            }
          }
        }
        break
      case 'fields':
        for (let item of val.fields) {
          if (item.value != '') {
            if (['DateTime', 'Number', 'Rating', 'Progress', 'Currency'].includes(item.ui_type)) {
              item.value = Number(item.value)
            }
            if ('Currency' === item.ui_type) item.value = item.value.toFixed(2)
          }
          if ('Checkbox' === item.ui_type) {
            item.value = (item.value == 'true')
          }
          if ('MultiSelect' === item.ui_type) {
            item.value = item.value.replace(/，/g, ',').split(',')
          }
          if ('DateTime' === item.ui_type) {
            if (item.value < 1e11) item.value *= 1000
          }
        }
        break
      case 'field_names':
        val.field_names = val.field_names.map(i => i.field_name)
        break
      case 'sort':
        val.sort = val.sort.map(i => {
          return {
            field_name: i.field_name,
            desc: (i.is_asc != 1)
          }
        })
        break
    }
  }
  return val
}

function nodeParamsAssign() {
  let nodeParams = JSON.parse(props.node.node_params)
  let arg = nodeParams?.plugin?.params?.arguments || {}
  if (arg.app_id) formState.app_id = arg.app_id
  if (arg.msgtype) formState.msgtype = arg.msgtype
  if (arg.receiver) formState.receiver = String(arg.receiver)
  if (arg.openid) formState.receiver = String(arg.openid)
  if (arg.touser) formState.receiver = String(arg.touser)
  if (arg.tag_map && arg.tag_map.receiver) formState.receiver_tags = arg.tag_map.receiver
  if (arg.app_secret) formState.app_secret = String(arg.app_secret)
  if (arg.app_name) formState.app_name = String(arg.app_name)
  if (formState.app_id) {
    loadTables()
  }
  if (formState.msgtype) {
    loadFields().then(() => compInit(nodeParams))
  }
}

function appChange(value, option) {
  tableData.value = {}
  tables.value = []
  fields.value = []
  formState.msgtype = undefined
  if (option && option.app_secret) {
    formState.app_secret = option.app_secret
  }
  if (option && option.name) {
    formState.app_name = option.name
  }
  if (formState.app_id) {
    loadTables().then(() => compInit())
  } else {
    compInit()
  }
  update()
}

function msgTypeChange() {
  fields.value = []
  loadFields().then(() => compInit())
  update()
}

function compInit(nodeParams=null) {
  nextTick(() => {
    if (childRef.value && typeof childRef.value.init === 'function') {
      childRef.value.init(nodeParams)
    }
  })
}

function onWorkflowValidate() {
  const bus = useEventBus()
  if (childRef.value && typeof childRef.value.validateAll === 'function') {
    const res = childRef.value.validateAll()
    const typeItem = tables.value.find(i => i.msgtype == formState.msgtype)
    if (res && res.ok === false) {
      bus.emit && bus.emit('workflow:validate:error', {
        component: '公众号发送消息',
        type: formState.msgtype,
        typeDisplay: typeItem?.name || formState.msgtype,
        message: '请完善发送消息的必填项',
        errors: res.errors || []
      })
    }
  }
}

function update(val = null) {
  let nodeParams = JSON.parse(props.node.node_params)
  let typeItem = tables.value.find(i => i.msgtype == formState.msgtype)
  //let params = nodeParams?.plugin?.params || {}
  nodeParams.plugin.output_obj = outputData.value
  Object.assign(nodeParams.plugin.params, {
    arguments: {
      app_id: String(formState.app_id || ''),
      app_secret: String(formState.app_secret || ''),
      msgtype: String(formState.msgtype || ''),
      touser: String(formState.receiver || ''),
      fields: (paramsFormat(val) || {}).fields || [],
      app_name: String(formState.app_name || ''),
      msg_type_name: String(typeItem?.name || ''),
    }
  })
  setData({
    ...props.node,
    node_params: JSON.stringify(nodeParams)
  })
}

function onChangeReceiver (text, selectedList) {
  formState.receiver = text
  formState.receiver_tags = selectedList
  update()
}


function outputFormat() {
  outputData.value = pluginOutputToTree(props.action.output)
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
</style>
