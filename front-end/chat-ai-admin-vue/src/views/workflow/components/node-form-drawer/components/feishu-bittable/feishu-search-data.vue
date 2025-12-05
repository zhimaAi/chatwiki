<template>
  <div class="select-data-form">
    <div class="node-box-content">
      <div class="setting-box">
        <collapse-panel title="查询条件" :count="state.filter.conditions.length">
          <QueryConditionFilter
            :disabled="!tableId"
            :where="state.filter.conditions"
            :field-options="fields"
            :conjunction="state.filter.conjunction"
            @change="onConditionChagne"
          />
        </collapse-panel>
      </div>

      <div class="setting-box" @wheel.stop @touchmove.stop>
        <collapse-panel title="查询字段" :count="state.field_names.length">
          <FieldListSelect
            :showAdd="!!tableId"
            :showInput="false"
            :showEmptyFieldRow="!tableId"
            :list="state.field_names"
            :fields="fields"
            @change="onChangeFields"
          />
        </collapse-panel>
      </div>

      <div class="setting-box" @wheel.stop @touchmove.stop>
        <collapse-panel title="排序" :count="state.sort.length">
          <SortSelector
            :disabled="!tableId"
            :list="state.sort"
            :fields="fields"
            @change="changeSoreField"
          />
        </collapse-panel>
      </div>

      <div class="setting-label">
        <span>查询数量</span>
      </div>

      <div class="setting-box" @wheel.stop @touchmove.stop>
        <a-input-number
           v-model:value="state.page_size"
           :disabled="!tableId"
           :min="20"
           :max="500"
           style="width: 205px"
           @change="update"/>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref, reactive, watch} from 'vue';
import FieldListSelect from "@/views/workflow/components/feishu-table/field-selector/index.vue";
import QueryConditionFilter from "@/views/workflow/components/feishu-table/query-condition-filter.vue";
import SortSelector from "@/views/workflow/components/feishu-table/sort-selector/index.vue";
import CollapsePanel from "@/views/workflow/components/data-table/collapse-panel.vue";
import OutputFields from "@/views/workflow/components/data-table/output-fields.vue";

const emit = defineEmits(['update'])
const props = defineProps({
  tableId: {
    type: [String, Number]
  },
  fields: {
    type: Array,
    default: () => ([])
  },
})

const stateStruct = {
  field_names: [],
  filter: {
    conjunction:'and',
    conditions: [],
  },
  sort: [],
  page_size: 100,
}
const state = reactive(JSON.parse(JSON.stringify(stateStruct)))

function init(nodeParams= null) {
  if (!nodeParams) {
    Object.assign(state, JSON.parse(JSON.stringify(stateStruct)))
  } else {
    let {field_names, filter, sort, page_size} = nodeParams?.plugin?.params?.arguments || {}
    if (Array.isArray(filter.conditions) && filter.conditions.length) {
      for (let item of filter.conditions) {
        if ('DateTime' === item.field_type) {
          // 后端要求DateTime添加了特殊标记 输入时移除标记
          item.value = item.value.replace(/\D/g, '')
        }
      }
    }
    if (Array.isArray(field_names)) {
      field_names.map(name => {
        let find = props.fields.find(i => i.field_name == name)
        find && state.field_names.push(find)
      })
    }
    if (Array.isArray(sort)) {
      state.sort = sort.map(item => {
        let find = props.fields.find(i => i.field_name == item.field_name)
        return {
          ...find,
          is_asc: item.desc ? 0 : 1
        }
      })
    }
    if (filter) state.filter = filter
    state.page_size = page_size || 100
  }
  //update()
}

const onConditionChagne = ({conjunction, list}) => {
  state.filter.conjunction = conjunction
  state.filter.conditions = list
  update()
}

const changeSoreField = (list) => {
  state.sort = [...list]
  update()
}

const onChangeFields = (selectedRows) => {
  state.field_names = selectedRows
  update()
}

function update() {
  emit('update', JSON.parse(JSON.stringify(state)))
}

defineExpose({
  init,
})
</script>

<style scoped lang="less">
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
