<template>
  <div class="main-box">
    <div v-if="loadingF || loading" class="loading-box">
      <a-spin/>
    </div>
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>{{ t('title_input') }}</div>
        <a-select
          v-if="Object.keys(configData).length"
          v-model:value="currentConfigName"
          @change="configChange">
          <a-select-option v-for="item in configData" :key="item.name" :value="item.name">
            {{ item.name }}
          </a-select-option>
        </a-select>
        <a-button v-else @click="showConfigModal">{{ t('btn_unauthorized') }} <DownOutlined/></a-button>
      </div>
      <template v-if="'create_bitable' != actionName">
        <div class="options-item is-required">
          <div class="options-item-tit">
            <div class="option-label">{{ t('label_multidimensional_table') }}</div>
          </div>
          <div class="min-input">
            <AtInput
              type="textarea"
              inputStyle="height: 33px;"
              :options="variableOptions"
              :defaultSelectedList="formState.tag_map?.app_token || []"
              :defaultValue="formState.app_token"
              ref="atInputRef"
              @open="emit('updateVar')"
              @change="(val, tags) => changeValue('app_token', val, tags)"
              :placeholder="t('ph_input_document_url')"
            >
              <template #option="{ label, payload }">
                <div class="field-list-item">
                  <div class="field-label">{{ label }}</div>
                  <div class="field-type">{{ payload.typ }}</div>
                </div>
              </template>
            </AtInput>
          </div>
          <div class="desc">{{ t('desc_multidimensional_table_id') }}</div>
        </div>
        <div v-if="!['table_add_members', 'update_advanced', 'roles_list', 'create_tables'].includes(actionName)" class="options-item is-required">
          <div class="options-item-tit">
            <div class="option-label">{{ t('label_select_table') }}</div>
          </div>
          <div class="flex-center min-input">
            <a-select v-model:value="formState.table_set_type" @change="tableSetChange" style="width: 112px;flex-shrink: 0;">
              <a-select-option value="1">{{ t('label_select_table') }}</a-select-option>
              <a-select-option value="2">{{ t('label_insert_variable') }}</a-select-option>
            </a-select>
            <a-select
              v-if="formState.table_set_type == 1"
              v-model:value="formState.table_id"
              @change="tableChange"
              :placeholder="formState.app_token ? t('ph_select_table') : t('ph_input_multidimensional_url_first')"
              style="width: 100%;">
              <a-select-option
                v-for="item in tables"
                :key="item.table_id"
                :value="item.table_id">
                {{ item.name }}
              </a-select-option>
            </a-select>
            <AtInput
              v-else
              type="textarea"
              inputStyle="height: 33px;"
              :options="variableOptions"
              :defaultSelectedList="formState.tag_map?.table_id || []"
              :defaultValue="formState.table_id"
              ref="atInputRef"
              @open="emit('updateVar')"
              @change="(val, tags) => changeValue('table_id', val, tags)"
              :placeholder="t('ph_input_table_id')"
            >
              <template #option="{ label, payload }">
                <div class="field-list-item">
                  <div class="field-label">{{ label }}</div>
                  <div class="field-type">{{ payload.typ }}</div>
                </div>
              </template>
            </AtInput>
          </div>
          <div class="desc">{{ t('desc_table_id_example') }}</div>
        </div>
      </template>
      <FeishuBatchBox
        v-if="BatchActions.includes(actionName)"
        ref="childRef"
        :variableOptions="variableOptions"
        :actionName="actionName"
        :action="action"
        @updateVar="emit('updateVar')"
        @update="update"
      />
      <component
        v-else
        ref="childRef"
        :is="actionComponentMap[actionName]"
        :tableId="formState.table_id"
        :fields="fields"
        :variableOptions="variableOptions"
        :config="{
           ...currentConfig,
           ...formState
        }"
        @update="update"
        @updateVar="emit('updateVar')"
      />
    </div>

    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/output.svg" class="title-icon"/>{{ t('title_output') }}</div>
      </div>
      <div class="options-item">
        <OutputFields :tree-data="outputData"/>
      </div>
    </div>

    <FeishuConfigModal ref="configModalRef" @change="loadConfig(true)"/>
  </div>
</template>

<script setup>
import {ref, reactive, onMounted, computed, inject, nextTick} from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import {runPlugin} from "@/api/plugins/index.js";
import FeishuInsertData from "./feishu-insert-data.vue";
import FeishuUpdateData from "./feishu-update-data.vue";
import FeishuDelData from "./feishu-del-data.vue";
import FeishuSearchData from "./feishu-search-data.vue";
import OutputFields from "@/views/workflow/components/feishu-table/output-fields.vue";
import {BatchActions, ShowFieldTypes} from "@/constants/feishu-table.js";

const { t } = useI18n('views.workflow.components.node-form-drawer.components.feishu-bittable.feishu-bittable-box')
import FeishuBatchBox from "./feishu-batch-box.vue";
import {getPluginActionDefaultArguments, pluginOutputToTree, getPluginConfigData} from "@/constants/plugin.js";
import {DownOutlined} from '@ant-design/icons-vue';
import FeishuConfigModal from "@/views/explore/plugins/components/feishu-config-modal.vue";
import FeishuCreateBitable from "./feishu-create-bitable.vue";
import FeishuCreateTable from "./feishu-create-table.vue";
import FeishuCreateView from "./feishu-create-view.vue";
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import {isValidURL} from "@/utils/validate.js";
import {jsonDecode} from "@/utils/index.js";
import FeishuAddMembers from "./feishu-add-members.vue";
import FeishuUpdateAdvanced from "./feishu-update-advanced.vue";

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
const actionComponentMap = {
  search_records: FeishuSearchData,
  create_record: FeishuInsertData,
  update_record: FeishuUpdateData,
  delete_record: FeishuDelData,
  create_bitable: FeishuCreateBitable,
  create_tables: FeishuCreateTable,
  create_views: FeishuCreateView,
  table_add_members: FeishuAddMembers,
  update_advanced: FeishuUpdateAdvanced,
}
const pluginName = 'feishu_bitable'

const childRef = ref(null)
const configModalRef = ref(null)
const configData = ref({})
const currentConfigName = ref(null)
const formState = reactive({
  app_token: '',
  table_id: undefined,
  tag_map: {},
  table_set_type: '1',
})
const loading = ref(false)
const loadingF = ref(false)
const tableData = ref({})
const tables = ref([])
const fields = ref([])

const currentConfig = computed(() => {
  return configData.value?.[currentConfigName.value] || {}
})
const outputData = ref([])

onMounted(() => {
  init()
})

async function init() {
  await loadConfig()
  nodeParamsAssign()
  outputFormat()
}

function showConfigModal() {
  configModalRef.value.show(configData.value)
}

async function loadConfig(refresh=false) {
  await getPluginConfigData(pluginName, refresh).then(res => {
    configData.value = res || {}
    for (let name in configData.value) {
      if (configData.value[name].is_default) {
        currentConfigName.value = name
        break
      }
    }
    if (!currentConfigName.value) currentConfigName.value = Object.keys(configData.value)[0]
  })
}

function loadTables() {
  loading.value = true
  tableData.value = {}
  tables.value = []
  if (!isValidURL(formState.app_token) || (Array.isArray(formState.tag_map?.app_token) && formState.tag_map.app_token.length)) {
    loading.value = false
    return Promise.resolve(null)
  }
  return runPlugin({
    name: 'feishu_bitable',
    action: "default/exec",
    params: JSON.stringify({
      business: 'get_tables',
      arguments: {
        app_id: currentConfig.value.appid,
        app_secret: currentConfig.value.app_secret,
        app_token: formState.app_token
      }
    })
  }).then(res => {
    tableData.value = res?.data || {}
    tables.value = tableData.value.items || []
    if (!formState.table_id && tables.value.length == 1) {
      formState.table_id = tables.value[0].table_id
      tableChange()
    }
    return res
  }).catch(() => {
    tableData.value = {}
    tables.value = []
    if (formState.table_set_type == 1) {
      formState.table_id = undefined
    }
    fields.value = []
    return null
  }).finally(() => {
    loading.value = false
  })
}

function loadFields() {
  loadingF.value = true
  if (!isValidURL(formState.app_token) || !formState.table_id || formState.tag_map?.table_id?.length) {
    loadingF.value = false
    return Promise.resolve(null)
  }
  return runPlugin({
    name: 'feishu_bitable',
    action: "default/exec",
    params: JSON.stringify({
      business: 'get_table_fields',
      arguments: {
        app_id: currentConfig.value.appid,
        app_secret: currentConfig.value.app_secret,
        app_token: formState.app_token,
        table_id: formState.table_id
      }
    })
  }).then(res => {
    let items = res?.data?.items || []
    fields.value = items.filter(i => ShowFieldTypes.includes(i.ui_type))
    return res
  }).finally(() => {
    loadingF.value = false
  })
}

function paramsFormat(val) {
  if (!val) {
    return getPluginActionDefaultArguments(pluginName, props.actionName)
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
          if (item?.atTags?.length) continue
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
            if (Array.isArray(item.value)) {
              // 已经是数组，无需处理
            } else if (typeof item.value === 'string') {
              const v = item.value.trim()
              item.value = v ? v.replace(/，/g, ',').split(',') : []
            } else {
              item.value = []
            }
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
      case 'input_type_map':
        // 查询条件的json输入
        for (let f in val[key]) {
          if (val[key][f] == 2) {
            val[f] = jsonDecode(val[`${f}_json`], f == 'filter' ? {} : [])
          }
        }
        break
    }
  }
  if (val.tag_map) {
    formState.tag_map = {
      ...val.tag_map,
      app_token: formState?.tag_map?.app_token || [],
      table_id: formState?.tag_map?.table_id || [],
    }
  }

  delete val.field_names_json
  delete val.filter_json
  delete val.sort_json
  delete val.fields_json
  return val
}

function nodeParamsAssign() {
  let nodeParams = JSON.parse(props.node.node_params)
  let arg = nodeParams?.plugin?.params?.arguments || {}
  if (arg.table_set_type) formState.table_set_type = arg.table_set_type
  if (arg.tag_map) formState.tag_map = arg.tag_map
  if (arg.app_token) formState.app_token = arg.app_token
  if (arg.table_id) formState.table_id = arg.table_id
  if (arg.app_id && arg.app_secret) {
    for (let name in configData.value) {
      let i = configData.value[name]
      if (i.appid == arg.app_id && i.app_secret == arg.app_secret) {
        currentConfigName.value = name
      }
    }
  }
  compInit(nodeParams)
  if (formState.app_token
    && !['table_add_members', 'update_advanced'].includes(props.actionName)) {
    loadTables()
  }
  if (formState.table_id) {
    loadFields().then(() => compInit(nodeParams))
  }
}

function tableSetChange () {
  if (formState.table_set_type == 1) {
    formState.table_id = tables.value?.[0]?.table_id || undefined
    formState.tag_map.table_id = []
    fields.value = []
    loadFields().then(() => compInit())
  } else {
    formState.table_id = ''
    fields.value = []
    if (props.actionName === 'create_record') {
      compInit({ plugin: { params: { arguments: { fields: [] } } } })
    } else {
      compInit()
    }
  }
  update({ fields: [] })
}

function changeValue(field, val, tags) {
  formState[field] = val
  formState.tag_map[field] = tags
  if (field === 'app_token') {
    tokenChange()
  } else if (field === 'table_id') {
    tableChange()
  } else {
    update()
  }
}

function configChange() {
  if (formState.app_token) {
    loadTables().then(() => compInit())
  } else {
    compInit()
  }
  if (formState.table_id) {
    loadFields().then(() => compInit())
  }
  update()
}

function tokenChange() {
  const prevTableId = formState.table_id
  fields.value = []
  loadTables().then(() => {
    const exists = tables.value.find(i => i.table_id == prevTableId)
    if (formState.table_set_type == 2) {
      formState.table_id = prevTableId
      if (props.actionName === 'create_record') {
        typeof childRef.value?.clear === 'function' && childRef.value.clear()
      }
      if (props.actionName === 'search_records') {
        typeof childRef.value?.clearSelections === 'function' && childRef.value.clearSelections()
      }
      compInit()
    } else {
      if (exists) {
        formState.table_id = prevTableId
        loadFields().then(() => {
          if (props.actionName === 'create_record') {
            typeof childRef.value?.clear === 'function' && childRef.value.clear()
          } else if (props.actionName !== 'create_tables' && props.actionName !== 'create_views' && props.actionName !== 'create_bitable') {
            compInit({ plugin: { params: { arguments: { fields: [] } } } })
          } else {
            compInit()
          }
          if (props.actionName === 'search_records') {
            typeof childRef.value?.clearSelections === 'function' && childRef.value.clearSelections()
          }
        })
      } else {
        formState.table_id = undefined
        if (props.actionName === 'create_record') {
          typeof childRef.value?.clear === 'function' && childRef.value.clear()
        } else if (props.actionName !== 'create_tables' && props.actionName !== 'create_views' && props.actionName !== 'create_bitable') {
          compInit({ plugin: { params: { arguments: { fields: [] } } } })
        } else {
          compInit()
        }
        if (props.actionName === 'search_records') {
          typeof childRef.value?.clearSelections === 'function' && childRef.value.clearSelections()
        }
      }
    }
  }).catch(() => {
    tables.value = []
    fields.value = []
    if (formState.table_set_type == 1) {
      formState.table_id = undefined
      if (props.actionName === 'create_record') {
        typeof childRef.value?.clear === 'function' && childRef.value.clear()
      } else if (props.actionName !== 'create_tables' && props.actionName !== 'create_views' && props.actionName !== 'create_bitable') {
        compInit({ plugin: { params: { arguments: { fields: [] } } } })
      } else {
        compInit()
      }
    } else {
      // 插入变量模式保持用户输入
      if (props.actionName === 'create_record') {
        typeof childRef.value?.clear === 'function' && childRef.value.clear()
      } else {
        compInit()
      }
    }
    if (props.actionName === 'search_records') {
      typeof childRef.value?.clearSelections === 'function' && childRef.value.clearSelections()
    }
  })
  update()
}

function tableChange() {
  fields.value = []
  loadFields().then(() => {
    update({ fields: [] })
    compInit()
  })
}

function compInit(nodeParams=null) {
  nextTick(() => {
    const np = nodeParams ?? JSON.parse(props.node.node_params)
    typeof childRef.value?.init === "function" && childRef.value.init(np)
  })
}

function update(val = null) {
  let nodeParams = JSON.parse(props.node.node_params)
  let table = tables.value.find(i => i.table_id == formState.table_id)
  const prevArgs = nodeParams?.plugin?.params?.arguments || {}
  const nextPart = val ? paramsFormat(val) : {}
  const mergedTagMap = {
    ...(prevArgs?.tag_map || {}),
    ...(nextPart?.tag_map || {}),
    ...(formState?.tag_map || {})
  }
  const mergedArgs = {
    ...prevArgs,
    ...nextPart,
    ...formState,
    tag_map: mergedTagMap,
    app_id: currentConfig.value.appid,
    app_secret: currentConfig.value.app_secret,
  }
  nodeParams.plugin.output_obj = outputData.value
  nodeParams.plugin.params = nodeParams.plugin.params || {}
  nodeParams.plugin.params.config_name = currentConfigName.value
  nodeParams.plugin.params.table_name = table?.name || ''
  nodeParams.plugin.params.arguments = mergedArgs
  setData({
    ...props.node,
    node_params: JSON.stringify(nodeParams)
  })
}

function outputFormat() {
  outputData.value = pluginOutputToTree(props.action.output)
}
</script>

<style scoped lang="less">
@import "common";
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
}

.flex-center {
  display: flex;
  align-items: center;
}
</style>
