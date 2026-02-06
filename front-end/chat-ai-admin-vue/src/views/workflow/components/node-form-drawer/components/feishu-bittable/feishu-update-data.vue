<style lang="less" scoped>
.add-data-form {
  .node-box-content {
    overflow: hidden;
    border-radius: 6px;
    background: #f2f4f7;
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

  }
}

.options-item {
  display: flex;
  flex-direction: column;
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
    word-break: break-all;
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
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_record_id') }}</div>
          <div class="option-type">string</div>
        </div>
        <div>
          <AtInput
            type="textarea"
            inputStyle="height: 64px;"
            :options="variableOptions"
            :defaultSelectedList="state.record_tags"
            :defaultValue="state.record_id"
            :ref="el => atInputRef['record_id'] = el"
            @open="emit('updateVar')"
            @change="(val, tags) => changeValue('record_id', val, tags)"
            :placeholder="t('ph_input_content')"
          >
            <template #option="{ label, payload }">
              <div class="field-list-item">
                <div class="field-label">{{ label }}</div>
                <div class="field-type">{{ payload.typ }}</div>
              </div>
            </template>
          </AtInput>
        </div>
        <div class="desc">{{ t('desc_update_record_id') }}</div>
      </div>
      <div class="options-item is-required">
        <div class="flex-between">
          <div class="options-item-tit">
            <div class="option-label">{{ t('label_update_data') }}</div>
          </div>
          <div class="flex-between">
            <FullscreenOutlined v-if="state.input_type_map.fields == 2" class="zm-pointer" @click="showFull('fields')"/>
            <a-select v-model:value="state.input_type_map.fields" style="width: 130px;" @change="update">
              <a-select-option :value="1">{{ t('opt_select_update_data') }}</a-select-option>
              <a-select-option :value="2">{{ t('opt_input_variable') }}</a-select-option>
            </a-select>
          </div>
        </div>
        <div class="setting-box">
          <FieldListSelect
            v-if="state.input_type_map.fields == 1"
            :showEmptyFieldRow="!tableId"
            :list="state.fields"
            :fields="fields"
            :showAdd="true"
            :showDelete="true"
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
              :placeholder="t('ph_input_content')"
            />
            <div class="desc">{{ t('desc_update_data_example') }}</div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, toRaw } from 'vue'
import {FullscreenOutlined} from '@ant-design/icons-vue'
import FieldListSelect from "@/views/workflow/components/feishu-table/field-selector/index.vue";
import AtFullInput from "@/views/workflow/components/at-input/at-full-input.vue";
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.components.feishu-bittable.feishu-update-data')

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

const atInputRef = ref({})
const state = reactive({
  record_id: '',
  record_tags: [],
  fields: [],
  fields_json: "",
  input_type_map: {
    fields: 1
  },
  tag_map: {
    fields: []
  },
})

function init(nodeParams = null) {
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
  if (nodeParams && nodeParams?.plugin?.params?.arguments?.record_id) {
    state.record_id = nodeParams.plugin.params.arguments.record_id
    state.record_tags = nodeParams.plugin.params.arguments.record_tags || []
  }
  update()
}

const onChangeFields = (selectedRows) => {
  state.fields = selectedRows
  update()
}

function changeValue(field, val, tags) {
  state[field] = val
  if (field == 'record_id') {
    state.record_tags = tags
  } else {
    state.tag_map[field] = tags
  }
  update()
}

function showFull(field) {
  atInputRef.value[field]?.showFull()
}

const update = () => {
  emit('update', {
    ...state
  })
}

defineExpose({
  init
})
</script>

