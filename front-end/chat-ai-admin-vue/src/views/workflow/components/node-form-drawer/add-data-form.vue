<style lang="less" scoped>
.add-data-form {
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
          <span
            >{{ t('title_desc_insert_data') }}
            <a target="_blank" href="/#/database/list">{{ t('btn_go_manage') }}</a></span
          >
        </template>
      </NodeFormHeader>
    </template>

    <div class="add-data-form">
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

        <div class="setting-label">
          <span>{{ t('label_insert_data') }}</span>
          <span class="tip">{{ t('tip_empty_field') }}</span>
        </div>

        <div class="setting-box">
          <FieldListSelect
            :form-id="state.formData.form_id"
            :showEmptyFieldRow="state.formData.form_id == ''"
            :list="state.formData.datas"
            :showAdd="false"
            :showDelete="false"
            @change="onChangeFields"
          />
        </div>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { onMounted, reactive, toRaw } from 'vue'
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import FieldListSelect from '../data-table/field-selector/index.vue'
import DataTableSelect from '../data-table/database-selector/index.vue'
import { useDataTableStore } from '@/stores/modules/data-table'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.add-data-form')

const emit = defineEmits(['update-node'])

const props = defineProps({
  lf: {
    type: Object,
    default: null
  },
  nodeType: {
    type: String
  },
  nodeId: {
    type: String,
    default: ''
  },
  node: {
    type: Object,
    default: () => ({})
  }
})

const dataTableStore = useDataTableStore()

let node_params = {}

const state = reactive({
  tableList: [],
  formData: {
    form_name: '',
    form_description: '',
    form_id: '',
    datas: []
  }
})

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
  }

  update()
}

const onChangeFields = (selectedRows) => {
  state.formData.datas = selectedRows
  update()
}

const init = () => {
  let dataRaw = props.node.dataRaw || props.node.node_params || '{}'

  node_params = JSON.parse(dataRaw)

  state.formData = node_params.form_insert

  update()
}

const update = () => {
  node_params.form_insert = toRaw(state.formData)

  emit('update-node', {
    ...props.node,
    node_params: JSON.stringify(node_params)
  })
}

const getTableList = async () => {
  const list = await dataTableStore.getFormList()

  state.tableList = list
}

const getFieldList = async (form_id) => {
  dataTableStore
    .getFormFieldList({ form_id: form_id })
    .then((list) => {
      list.forEach((item) => {
        item.id = item.id * 1
      })

      state.formData.datas = list

      update()
    })
    .catch(() => {})
}

onMounted(() => {
  getTableList()
  init()
})
</script>
