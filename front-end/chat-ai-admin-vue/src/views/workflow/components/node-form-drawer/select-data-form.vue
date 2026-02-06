<style lang="less" scoped>
.select-data-form {
  .node-box-content {
    margin-top: 16px;
    overflow: hidden;
    border-radius: 6px;
    background: #f2f4f7;
  }
  .node-box-title {
    display: flex;
    align-items: center;
    height: 48px;
    padding: 0 12px;

    .input-icon {
      width: 16px;
      height: 16px;
      margin-right: 4px;
    }

    .text {
      font-weight: 600;
      font-size: 14px;
      color: #262626;
    }
  }

  .setting-label {
    line-height: 22px;
    padding: 0 12px;
    margin-bottom: 8px;
    font-size: 14px;
    color: #262626;
    .tip {
      color: #8c8c8c;
    }
  }

  .setting-box {
    padding: 0 12px;
    margin-bottom: 12px;
  }
}
</style>

<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader :title="node.node_name" :iconName="node.node_icon_name">
        <template #desc>
          <span >{{ t('title_desc_select_data') }}
            <a target="_blank" href="/#/database/list">{{ t('btn_go_manage') }}</a>
          </span>
        </template>
      </NodeFormHeader>
    </template>

    <div class="select-data-form">
      <div class="node-box-content">
          <div class="node-box-title">
            <img class="input-icon" src="@/assets/img/workflow/input.svg" alt="" />
            <span class="text">{{ t('label_input') }}</span>
          </div>

          <div class="setting-label">
            <span>{{ t('label_data_table') }}</span>
          </div>

          <div class="setting-box">
            <DataTableSelect
              :options="state.tableList"
              :title="state.formData.form_name"
              :description="state.formData.form_description"
              v-model:value="state.formData.form_id"
              @change="onSelectTable"
            />
          </div>

          <div class="setting-box">
            <collapse-panel :title="t('label_query_condition')" :count="state.formData.where.length">
              <QueryConditionFilter
                :disabled="!state.formData.form_id"
                :where="state.formData.where"
                :field-options="state.fieldList"
                :type="state.formData.typ"
                @change="onConditionChagne"
              />
            </collapse-panel>
          </div>

          <div class="setting-box" @wheel.stop @touchmove.stop>
            <collapse-panel :title="t('label_query_fields')" :count="state.formData.fields.length">
              <FieldListSelect
                :form-id="state.formData.form_id"
                :showAdd="state.formData.form_id != ''"
                :showInput="false"
                :showEmptyFieldRow="state.formData.form_id == ''"
                :list="state.formData.fields"
                @change="onChangeFields"
              />
            </collapse-panel>
          </div>

          <div class="setting-box" @wheel.stop @touchmove.stop>
            <collapse-panel :title="t('label_sort')" :count="state.formData.order.length">
              <SortSelector
                :form-id="state.formData.form_id"
                :disabled="!state.formData.form_id"
                :list="state.formData.order"
                @change="changeSoreField"
              />
            </collapse-panel>
          </div>

          <div class="setting-label">
            <span>{{ t('label_query_count') }}</span>
          </div>

          <div class="setting-box" @wheel.stop @touchmove.stop>
            <a-input-number :min="1" :max="1000" v-model:value="state.formData.limit" :disabled="!state.formData.form_id" style="width: 205px" @change="changeLimit" />
          </div>
        </div>

        <div class="node-box-content">
          <div class="node-box-title">
            <img class="input-icon" src="@/assets/img/workflow/output.svg" alt="" />
            <span class="text">{{ t('label_output') }}</span>
          </div>

          <div class="setting-box">
            <OutputFields :treeData="outputFields" />
          </div>
        </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import {reactive, inject, onMounted, toRaw, computed } from 'vue'
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import CollapsePanel from '../data-table/collapse-panel.vue'
import FieldListSelect from '../data-table/field-selector/index.vue'
import DataTableSelect from '../data-table/database-selector/index.vue'
import QueryConditionFilter from '../data-table/query-condition-filter.vue'
import SortSelector from '../data-table/sort-selector/index.vue'
import OutputFields from '../data-table/output-fields.vue'
import { useDataTableStore } from '@/stores/modules/data-table'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.select-data-form')

const props = defineProps({
  node: {
    type: Object,
    default: () => ({})
  },
})

const setData = inject('setData')

const dataTableStore = useDataTableStore()

let node_params = {}

const state = reactive({
  tableList: [],
  formData: {
    form_name: '',
    form_description: '',
    form_id: '',
    typ: 1,
    fields: [],
    where: [],
    order: [],
    limit: 100
  },
  fieldList: []
})

const outputFields = computed(() => {
  let fields = state.formData.fields.map((item) => {
    return {
      title: item.name,
      key: item.name,
      type: item.type
    }
  })

  return [
    {
      title: 'output_list',
      key: 'output_list',
      type: 'array<object>',
      children: fields
    }, {
      title: 'row_num',
      key: 'row_num',
      type: 'integer',
    }
  ]
})

const onConditionChagne = ({type, list}) => {
  state.formData.typ = type
  state.formData.where = list

  update()
}

const changeSoreField = (list) => {
  state.formData.order = [...list]

  update()
}

const onSelectTable = (value, record) => {
  if (record) {
    state.formData.form_name = record.name
    state.formData.form_description = record.description
    // state.formData.form_id = record.id
    getFieldList(record.id)
  } else {
    state.formData.form_name = ''
    state.formData.form_description = ''
    state.formData.typ = 1
    state.formData.fields = []
    state.formData.where = []
    state.formData.order = []
    state.formData.limit = 100
    // state.formData.form_id = ''
  }

  update()
}

const onChangeFields = (selectedRows) => {
  state.formData.fields = selectedRows
  update()
}

const changeLimit = () => {
  update()
}

const init = () => {
  let dataRaw = props.node.dataRaw || props.node.node_params || '{}'

  node_params = JSON.parse(dataRaw)

  let where = whereDataConditions(node_params.form_select.where || [])

  node_params.form_select.where = where

  state.formData = node_params.form_select

  if (state.formData.form_id) {
    getFieldList(state.formData.form_id)
  }

  update()
}

const update = () => {
  let form = { ...toRaw(state.formData) }

  form.where = whereDataFormatter(form.where)

  node_params.form_select = form

  setData({
    ...props.node,
    node_params: JSON.stringify(node_params)
  })
}

const whereDataConditions = (where) => {
  let conditions = []

  where.forEach((item) => {
    let data = JSON.parse(JSON.stringify(item))
    let ruleArr = data.rule.split('_')

    data.form_field_id = data.form_field_id * 1
    data.field_type = ruleArr[0]
    // 删除数组的第一个
    ruleArr.shift()
    // 数组还原成字符串
    data.rule = ruleArr.join('_')

    conditions.push(data)
  })

  return conditions
}

const whereDataFormatter = (where) => {
  let conditions = []

  where.forEach((item) => {
    let data = JSON.parse(JSON.stringify(item))

    data.form_field_id = data.form_field_id * 1
    data.rule = item.field_type + '_' + item.rule

    conditions.push(data)
  })

  return conditions
}

const getTableList = async () => {
  const list = await dataTableStore.getFormList()
  if (list) {
    state.tableList = list
  }
}

const getFieldList = async (form_id) => {
  dataTableStore.getFormFieldList({ form_id: form_id })
    .then((list) => {
      list.forEach((item) => {
        item.id = item.id * 1
      })

      state.fieldList = list
    })
    .catch(() => {})
}

onMounted(() => {
  getTableList()
  init()
})
</script>
