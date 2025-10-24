<style lang="less" scoped>
.node-box {
  a{
    position: relative;
  }
  .node-desc {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: var(--wf-color-text-2);
  }

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
  <node-common
    :title="props.properties.node_name"
    :menus="menus"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
    @handleMenu="handleMenu"
  >
    <div class="node-box">
      <div class="node-desc">
        在指定数据表中更新符合条件的数据，可以在数据库中维护数据表 <a target="_blank" href="/#/database/list">去管理</a>
      </div>

      <div class="node-box-content">
        <div class="node-box-title">
          <img class="input-icon" src="../../../../../assets/img/workflow/input.svg" alt="" />
          <span class="text">输入</span>
        </div>

        <div class="setting-label">
          <span>数据表</span>
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

        <div class="setting-label">
          <span>更新条件</span>
        </div>

        <div class="setting-box">
          <QueryConditionFilter
            :disabled="!state.formData.form_id"
            :where="state.formData.where"
            :field-options="state.fieldList"
            :type="state.formData.typ"
            @change="onConditionChagne"
          />
        </div>

        <div class="setting-label">
          <span>更新数据</span>
        </div>

        <div class="setting-box">
          <FieldListSelect
            :form-id="state.formData.form_id"
            :showAdd="state.formData.form_id != ''"
            :showEmptyFieldRow="state.formData.form_id == ''"
            :list="state.formData.datas"
            @change="onChangeFields"
          />
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { ref, inject, onMounted, reactive, toRaw } from 'vue'
import NodeCommon from '../base-node.vue'
import FieldListSelect from '../../data-table/field-selector/index.vue'
import DataTableSelect from '../../data-table/database-selector/index.vue'
import QueryConditionFilter from '../../data-table/query-condition-filter.vue'
import { useDataTableStore } from '@/stores/modules/data-table'

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  },
  isSelected: { type: Boolean, default: false },
  isHovered: { type: Boolean, default: false }
})

const getNode = inject('getNode')
const getGraph = inject('getGraph')
const setData = inject('setData')

const dataTableStore = useDataTableStore()

const menus = ref([{ name: '删除', key: 'delete', color: '#fb363f' }])

const handleMenu = (item) => {
  if (item.key === 'delete') {
    let node = getNode()
    getGraph().deleteNode(node.id)
  }
}

let node_params = {}

const state = reactive({
  tableList: [],
  formData: {
    form_name: '',
    form_description: '',
    form_id: '',
    typ: 1,
    datas: [],
    where: []
  },
  fieldList: []
})

const onConditionChagne = ({type, list}) => {
  state.formData.typ = type
  state.formData.where = list

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
    // state.formData.form_id = ''
    state.formData.datas = []
    state.formData.where = []
    state.formData.typ = 1
  }

  update()
}

const onChangeFields = (selectedRows) => {
  state.formData.datas = selectedRows
  update()
}

const init = () => {
  let dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'

  node_params = JSON.parse(dataRaw)

  let where = whereDataConditions(node_params.form_update.where || [])

  node_params.form_update.where = where

  state.formData = node_params.form_update

  if (state.formData.form_id) {
    getFieldList(state.formData.form_id)
  }

  update()
}

const update = () => {
  let form = { ...toRaw(state.formData) }

  form.where = whereDataFormatter(form.where)

  node_params.form_update = form

  setData({
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
      list.forEach(item => {
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
