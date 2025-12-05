<style lang="less" scoped>
.add-data-form {
  .node-box-content {
    margin-top: 16px;
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
    margin-top: 12px;
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
</style>

<template>
  <div class="add-data-form">
    <div class="node-box-content">
<!--      <div class="setting-label">-->
<!--        <span>更新数据</span>-->
<!--        <span class="tip"></span>-->
<!--      </div>-->
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">record_id</div>
          <div class="option-type">string</div>
        </div>
        <div>
          <AtInput
            type="textarea"
            inputStyle="height: 64px;"
            :options="variableOptions"
            :defaultSelectedList="state.tags"
            :defaultValue="state.value"
            ref="atInputRef"
            @open="emit('updateVar')"
            @change="changeValue"
            placeholder="请输入内容，键入“/”可以插入变量"
          >
            <template #option="{ label, payload }">
              <div class="field-list-item">
                <div class="field-label">{{ label }}</div>
                <div class="field-type">{{ payload.typ }}</div>
              </div>
            </template>
          </AtInput>
        </div>
        <div class="desc">待更新的记录ID；示例值："reCWNXZPQv"</div>
      </div>

      <div class="setting-box">
        <FieldListSelect
          :showEmptyFieldRow="!tableId"
          :list="dataItems"
          :fields="fields"
          :showAdd="true"
          :showDelete="true"
          @change="onChangeFields"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, toRaw } from 'vue'
import FieldListSelect from "@/views/workflow/components/feishu-table/field-selector/index.vue";
import AtInput from "@/views/workflow/components/at-input/at-input.vue";

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

const dataItems = ref([])
const state = reactive({
  tags: [],
  value: '',
})

function init(nodeParams = null) {
  dataItems.value = []
  if (Array.isArray(nodeParams?.plugin?.params?.arguments?.fields) && nodeParams.plugin.params.arguments.fields.length) {
    nodeParams.plugin.params.arguments.fields.map(item => {
      item.value = item.value.toString()
    })
    dataItems.value = nodeParams.plugin.params.arguments.fields
  }
  if (nodeParams && nodeParams?.plugin?.params?.arguments?.record_id) {
    state.value = nodeParams.plugin.params.arguments.record_id
    state.tags = nodeParams.plugin.params.arguments.record_tags || []
  }
}

const onChangeFields = (selectedRows) => {
  dataItems.value = selectedRows
  update()
}

function changeValue(val, tags) {
  state.value = val
  state.tags = tags
  update()
}


const update = () => {
  emit('update', {
    record_id: state.value,
    record_tags: state.tags,
    fields: toRaw(dataItems.value)
  })
}

defineExpose({
  init
})
</script>

