<template>
  <div class="select-data-form">
    <div class="node-box-content">
      <div class="setting-box">
        <collapse-panel title="查询条件" :count="state.filter.conditions.length">
          <template #extra>
            <div class="flex-between">
              <div v-if="state.input_type_map.filter != 1" class="btn-hover-wrap" @click="handleOpenFullAtModal('filter_json')">
                <FullscreenOutlined/>
              </div>
              <a-select v-model:value="state.input_type_map.filter" style="width: 130px;" @change="update">
                <a-select-option :value="1">选择查询条件</a-select-option>
                <a-select-option :value="2">输入变量</a-select-option>
              </a-select>
            </div>
          </template>
          <template v-if="state.input_type_map.filter != 1">
            <AtInput
              type="textarea"
              inputStyle="height: 64px;"
              :options="variableOptions"
              :defaultSelectedList="state.tag_map?.filter_json || []"
              :defaultValue="state.filter_json"
              :ref="el => atInputRef['filter_json'] = el"
              @open="emit('updateVar')"
              @change="(val, tags) => changeValue('filter_json', val, tags)"
              placeholder="请输入内容，键入“/”可以插入变量"
            >
              <template #option="{ label, payload }">
                <div class="field-list-item">
                  <div class="field-label">{{ label }}</div>
                  <div class="field-type">{{ payload.typ }}</div>
                </div>
              </template>
            </AtInput>
            <div class="desc">内容示例：{"conditions":[{"field_name":"文章ID","operator":"is","value":"123"}],"conjunction":"and"}</div>
          </template>
          <QueryConditionFilter
            v-else
            :disabled="!tableId"
            :where="state.filter.conditions"
            :field-options="fields"
            :conjunction="state.filter.conjunction"
            @change="onConditionChange"
          />
        </collapse-panel>
      </div>

      <div class="setting-box" @wheel.stop @touchmove.stop>
        <collapse-panel title="查询字段" :count="state.field_names.length">
          <template #extra>
            <div class="flex-between">
              <div v-if="state.input_type_map.field_names != 1" class="btn-hover-wrap" @click="handleOpenFullAtModal('field_names_json')">
                <FullscreenOutlined/>
              </div>
              <a-select v-model:value="state.input_type_map.field_names" style="width: 130px;" @change="update">
                <a-select-option :value="1">选择查询字段</a-select-option>
                <a-select-option :value="2">输入变量</a-select-option>
              </a-select>
            </div>
          </template>
          <template  v-if="state.input_type_map.field_names != 1">
            <AtInput
              type="textarea"
              inputStyle="height: 64px;"
              :options="variableOptions"
              :defaultSelectedList="state.tag_map?.field_names_json || []"
              :defaultValue="state.field_names_json"
              :ref="el => atInputRef['field_names_json'] = el"
              @open="emit('updateVar')"
              @change="(val, tags) => changeValue('field_names_json', val, tags)"
              placeholder="请输入内容，键入“/”可以插入变量"
            >
              <template #option="{ label, payload }">
                <div class="field-list-item">
                  <div class="field-label">{{ label }}</div>
                  <div class="field-type">{{ payload.typ }}</div>
                </div>
              </template>
            </AtInput>
            <div class="desc">内容示例：["文章ID","封面图","创建时间","介绍"]</div>
          </template>
          <FieldListSelect
            v-else
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
          <template #extra>
            <div class="flex-between">
              <div v-if="state.input_type_map.sort != 1" class="btn-hover-wrap" @click="handleOpenFullAtModal('sort_json')">
                <FullscreenOutlined/>
              </div>
              <a-select v-model:value="state.input_type_map.sort" style="width: 130px;" @change="update">
                <a-select-option :value="1">选择排序</a-select-option>
                <a-select-option :value="2">输入变量</a-select-option>
              </a-select>
            </div>
          </template>
          <template v-if="state.input_type_map.sort != 1">
            <AtInput
              type="textarea"
              inputStyle="height: 64px;"
              :options="variableOptions"
              :defaultSelectedList="state.tag_map?.sort_json || []"
              :defaultValue="state.sort_json"
              :ref="el => atInputRef['sort_json'] = el"
              @open="emit('updateVar')"
              @change="(val, tags) => changeValue('sort_json', val, tags)"
              placeholder="请输入内容，键入“/”可以插入变量"
            >
              <template #option="{ label, payload }">
                <div class="field-list-item">
                  <div class="field-label">{{ label }}</div>
                  <div class="field-type">{{ payload.typ }}</div>
                </div>
              </template>
            </AtInput>
            <div class="desc">内容示例：[{"desc":false,"field_name":"文章ID"},{"desc":true,"field_name":"更新时间"}]</div>
          </template>
          <SortSelector
            v-else
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
    <FullAtInput
      :options="variableOptions"
      :defaultSelectedList="fullDefaultTags"
      :defaultValue="fullDefaultValue"
      placeholder="请输入内容，键入“/”可以插入变量"
      type="textarea"
      @open="emit('updateVar')"
      @change="(val, tags) => changeValueByFull(val, tags)"
      @ok="handleRefreshAtInput"
      ref="fullAtInputRef"
    />
  </div>
</template>

<script setup>
import {reactive, ref} from 'vue';
import FieldListSelect from "@/views/workflow/components/feishu-table/field-selector/index.vue";
import QueryConditionFilter from "@/views/workflow/components/feishu-table/query-condition-filter.vue";
import SortSelector from "@/views/workflow/components/feishu-table/sort-selector/index.vue";
import CollapsePanel from "@/views/workflow/components/data-table/collapse-panel.vue";
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import FullAtInput from "@/views/workflow/components/at-input/full-at-input.vue";
import { FullscreenOutlined } from '@ant-design/icons-vue'
import {jsonDecode} from "@/utils/index.js";

const emit = defineEmits(['update', 'updateVar'])
const props = defineProps({
  variableOptions: {
    type: Array,
  },
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

  field_names_json: "",
  filter_json: "",
  sort_json: "",

  input_type_map: { // 1选择 2 输入
    field_names: 1,
    filter: 1,
    sort: 1,
  },
  tag_map: {},
}
const state = reactive(JSON.parse(JSON.stringify(stateStruct)))
const fullAtInputRef = ref(null)
const atInputRef = ref({})
const activeField = ref('')
const fullDefaultValue = ref('')
const fullDefaultTags = ref([])

function init(nodeParams= null) {
  if (!nodeParams) {
    Object.assign(state, JSON.parse(JSON.stringify(stateStruct)))
  } else {
    let {field_names, filter, sort, page_size, input_type_map, tag_map} = nodeParams?.plugin?.params?.arguments || {}
    if (Array.isArray(filter.conditions) && filter.conditions.length) {
      for (let item of filter.conditions) {
        if ('DateTime' === item.field_type) {
          // 后端要求DateTime添加了特殊标记 输入时移除标记
          item.value = item.value.replace(/\D/g, '')
        }
      }
      state.filter = filter
    } else {
      state.filter = {
        conjunction:'and',
        conditions: [],
      }
    }

    if (Array.isArray(field_names) && input_type_map?.field_names != 2) {
      field_names.map(name => {
        let find = props.fields.find(i => i.field_name == name)
        find && state.field_names.push(find)
      })
    } else {
      state.field_names = Array.isArray(field_names) ? field_names : []
    }
    if (Array.isArray(sort) && input_type_map?.sort != 2) {
      state.sort = sort.map(item => {
        let find = props.fields.find(i => i.field_name == item.field_name)
        return {
          ...find,
          is_asc: item.desc ? 0 : 1
        }
      })
    } else {
      state.sort = Array.isArray(sort) ? sort : []
    }
    if (input_type_map) state.input_type_map = input_type_map
    for (let f in input_type_map) {
      if (input_type_map[f] == 2) {
        let d = JSON.stringify(state[f])
        if (Array.isArray(state[f])) {
          if (!state[f].length) d = ''
        } else if (!state[f]?.conditions?.length) {
          d = ''
        }
        state[`${f}_json`] = d
      }
    }

    state.tag_map = tag_map || {}
    state.page_size = page_size || 100
  }
  //update()
}

const onConditionChange = ({conjunction, list}) => {
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

function changeValue(field, val, tags) {
  state[field] = val
  state.tag_map[field] = tags
  update()
}

function changeValueByFull(val, tags) {
  const field = activeField.value
  if (!field) return
  state[field] = val
  state.tag_map[field] = tags
  update()
}

function handleOpenFullAtModal(field) {
  activeField.value = field
  fullDefaultValue.value = state[field] || ''
  fullDefaultTags.value = state.tag_map?.[field] || []
  fullAtInputRef.value.show()
}

function update() {
  emit('update', JSON.parse(JSON.stringify(state)))
}

function handleRefreshAtInput() {
  const field = activeField.value
  const target = atInputRef.value?.[field]
  target && target.refresh && target.refresh()
}

defineExpose({
  init,
  clearSelections() {
    if (state.input_type_map?.filter == 1) {
      state.filter.conditions = []
    }
    if (state.input_type_map?.field_names == 1) {
      state.field_names = []
    }
    if (state.input_type_map?.sort == 1) {
      state.sort = []
    }
    update()
  },
})
</script>

<style scoped lang="less">
.select-data-form {
  :deep(.type-textarea){
    word-break: break-all;
  }

  .desc {
    color: var(--wf-color-text-2);
    margin-top: 4px;
    word-break: break-all;
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
    //padding: 0 12px;
    margin-bottom: 12px;
  }
}
.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.btn-hover-wrap {
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
}
</style>
