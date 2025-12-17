<template>
  <div class="main-box">
    <div v-if="loadingF || loading" class="loading-box">
      <a-spin/>
    </div>
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>输入</div>
        <a-select v-model:value="currentConfigName" @change="configChange">
          <a-select-option v-for="item in configData" :key="item.name" :value="item.name">
            {{ item.name }}
          </a-select-option>
        </a-select>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">多维表</div>
        </div>
        <div>
          <a-input v-model:value="formState.app_token" placeholder="请输入文档url" @blur="tokenChange"/>
        </div>
        <div class="desc">多维表格的唯一标识符，支持输入文档 url</div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">选择数据表</div>
        </div>
        <div>
          <a-select
            v-model:value="formState.table_id"
            @change="tableChange"
            :placeholder="formState.app_token ? '请选择数据表' : '请先输入多维表文档url'"
            style="width: 100%;">
            <a-select-option
              v-for="item in tables"
              :key="item.table_id"
              :value="item.table_id">
              {{ item.name }}
            </a-select-option>
          </a-select>
        </div>
      </div>
      <component
        ref="childRef"
        :is="actionComponentMap[actionName]"
        :tableId="formState.table_id"
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
import {ref, reactive, onMounted, computed, inject, nextTick} from 'vue'
import {runPlugin} from "@/api/plugins/index.js";
import FeishuInsertData from "./feishu-insert-data.vue";
import FeishuUpdateData from "./feishu-update-data.vue";
import FeishuDelData from "./feishu-del-data.vue";
import FeishuSearchData from "./feishu-search-data.vue";
import OutputFields from "@/views/workflow/components/feishu-table/output-fields.vue";
import {ShowFieldTypes} from "@/constants/feishu-table.js";
import {getPluginActionDefaultArguments, pluginOutputToTree, getPluginConfigData} from "@/constants/plugin.js";

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
}
const pluginName = 'feishu_bitable'

const childRef = ref(null)
const configData = ref({})
const currentConfigName = ref(null)
const formState = reactive({
  app_token: '',
  table_id: undefined,
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

async function loadConfig() {
  await getPluginConfigData(pluginName).then(res => {
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
  }).finally(() => {
    loading.value = false
  })
}

function loadFields() {
  loadingF.value = true
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
  if (formState.app_token) {
    loadTables()
  }
  if (formState.table_id) {
    loadFields().then(() => compInit(nodeParams))
  }
}

function configChange() {
  formState.app_token = ''
  formState.table_id = undefined
  tableData.value = {}
  tables.value = []
  fields.value = []
  compInit()
  update()
}

function tokenChange() {
  tableData.value = {}
  tables.value = []
  fields.value = []
  formState.table_id = undefined
  if (formState.app_token) {
    loadTables().then(res => compInit())
  } else {
    compInit()
  }
  update()
}

function tableChange() {
  fields.value = []
  loadFields().then(() => compInit())
  update()
}

function compInit(nodeParams=null) {
  nextTick(() => {
    childRef.value && childRef.value.init(nodeParams)
  })
}

function update(val = null) {
  let nodeParams = JSON.parse(props.node.node_params)
  let table = tables.value.find(i => i.table_id == formState.table_id)
  //let params = nodeParams?.plugin?.params || {}
  nodeParams.plugin.output_obj = outputData.value
  Object.assign(nodeParams.plugin.params, {
    config_name: currentConfigName.value,
    table_name: table?.name || '',
    arguments: {
      // ...params?.arguments,
      ...paramsFormat(val),
      ...formState,
      app_id: currentConfig.value.appid,
      app_secret: currentConfig.value.app_secret,
    }
  })
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
