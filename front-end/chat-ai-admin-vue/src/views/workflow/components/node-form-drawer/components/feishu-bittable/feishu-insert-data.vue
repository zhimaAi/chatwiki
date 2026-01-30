<style lang="less" scoped>
.add-data-form {
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

  .setting-label {
    line-height: 22px;
    margin-bottom: 8px;
    font-size: 14px;
    color: #262626;
    &.is-required::before {
      content: '*';
      color: #FB363F;
      display: inline-block;
      margin-right: 2px;
    }
    .tip {
      color: #8c8c8c;
    }
  }
}

.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.mt8 {
  margin-top: 8px;
}
</style>

<template>
  <div class="add-data-form">
    <div class="node-box-content">
      <div class="flex-between">
        <div class="setting-label is-required">
          <span>插入数据</span>
          <span class="tip">（不填写字段值则为空）</span>
        </div>
        <div class="flex-between">
          <FullscreenOutlined v-if="state.input_type_map.fields == 2" class="zm-pointer" @click="showFull('fields')"/>
          <a-select v-model:value="state.input_type_map.fields" style="width: 130px;" @change="update">
            <a-select-option :value="1">选择添加数据</a-select-option>
            <a-select-option :value="2">输入变量</a-select-option>
          </a-select>
        </div>
      </div>
      <div class="setting-box mt8">
        <FieldListSelect
          v-if="state.input_type_map.fields == 1"
          :showEmptyFieldRow="!tableId"
          :list="state.fields"
          :fields="fields"
          :showAdd="false"
          :showDelete="false"
          @change="onChangeFields"
        />
        <template v-else>
          <AtFullInput
            type="textarea"
            inputStyle="height: 64px;"
            :options="variableOptions"
            :defaultSelectedList="state.tag_map?.fields_json || []"
            :defaultValue="state.fields_json"
            :ref="el => atInputRef['fields'] = el"
            @open="emit('updateVar')"
            @change="(val, tags) => changeValue('fields_json', val, tags)"
            placeholder="请输入内容，键入“/”可以插入变量"
          />
          <div class="desc">内容示例：[{"field_name":"文本","ui_type":"Text","value":"001"},{"field_name":"日期","ui_type":"DateTime","value":1769077135214},{"field_name":"数字","ui_type":"Number","value":99},{"field_name":"单选","ui_type":"SingleSelect","value":"选项1"},{"field_name":"多选","ui_type":"MultiSelect","value":["类别A","类别B"]},{"field_name":"复选框","ui_type":"Checkbox","value":true}]</div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import {FullscreenOutlined} from '@ant-design/icons-vue'
import FieldListSelect from "@/views/workflow/components/feishu-table/field-selector/index.vue";
import AtFullInput from "@/views/workflow/components/at-input/at-full-input.vue";

const emit = defineEmits(['update'])
const props = defineProps({
  tableId: {
    type: [String, Number]
  },
  fields: {
    type: Array,
    default: () => ([])
  },
  variableOptions: {
    type: Array,
  },
})

const atInputRef = ref({})
const state = reactive({
  fields: [],
  fields_json: "",
  input_type_map: {
    fields: 1
  },
  tag_map: {
    fields: []
  },
})

function init(nodeParams=null) {
  state.fields = []
  if (nodeParams) {
    let { fields, input_type_map, tag_map} = nodeParams?.plugin?.params?.arguments || {}
    fields.map(item => {
      item.value = item.value ? item.value.toString() : ""
    })
    state.fields = fields
    state.tag_map = tag_map || {}
    if (input_type_map){
      state.input_type_map = input_type_map
      for (let f in input_type_map) {
        if (input_type_map[f] == 2 && !state[`${f}_json`]) {
          let d = JSON.stringify(state[f])
          state[`${f}_json`] = d
        }
      }
    }
  }
  if (!state.fields.length) {
    props.fields.map(item => {
      state.fields.push({
        value: '',
        atTags: [],
        ...item
      })
    })
  }
  update()
}
const onChangeFields = (selectedRows) => {
  state.fields = selectedRows
  update()
}

function showFull(field) {
  atInputRef.value[field]?.showFull()
}

function changeValue(field, val, tags) {
  state[field] = val
  state.tag_map[field] = tags
  update()
}

const update = () => {
  emit('update', {
    ...state
  })
}

defineExpose({
  init,
  clear() {
    update()
  }
})
</script>

