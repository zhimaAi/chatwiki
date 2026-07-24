<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        desc="从商品库中查找对应商品，最多返回100条数据"
        @close="handleClose"
      />
    </template>

    <div class="goods-search-form">
      <div class="node-box-content">
        <div class="node-box-title">
          <img class="input-icon" src="@/assets/img/workflow/input.svg" alt="" />
          <span class="text">输入</span>
        </div>

        <div class="setting-box">
          <draggable
            style="display: flex; flex-direction: column; gap: 8px"
            handle=".drag-btn"
            v-model="formState.condition_groups"
            item-key="key"
          >
            <template #item="{ element: group, index }">
              <div class="gray-block condition-group">
                <div class="gray-block-title">
                  <a-flex :gap="8">
                    <HolderOutlined class="icon drag-btn" />
                    查询条件{{ index + 1 }}
                  </a-flex>
                  <div
                    v-if="formState.condition_groups.length > 1"
                    class="btn-hover-wrap"
                    @click="handleDelGroup(index)"
                  >
                    <CloseCircleOutlined />
                  </div>
                </div>

                <div class="condition-list-box">
                  <div class="left-select-box">
                    <a-select
                      size="small"
                      v-model:value="group.is_or"
                      :bordered="false"
                      style="width: 64px"
                    >
                      <a-select-option :value="0">and</a-select-option>
                      <a-select-option :value="1">or</a-select-option>
                    </a-select>
                  </div>
                  <div class="condition-body">
                    <div
                      class="condition-item"
                      v-for="(condition, conditionIndex) in group.conditions"
                      :key="condition.key"
                    >
                      <a-select
                        v-model:value="condition.field"
                        style="width: 124px"
                        placeholder="请选择"
                        @change="handleFieldChange(condition)"
                      >
                        <a-select-option
                          v-for="option in fieldOptions"
                          :key="option.value"
                          :value="option.value"
                        >
                          {{ option.label }}
                        </a-select-option>
                      </a-select>
                      <a-select
                        v-model:value="condition.match"
                        style="width: 124px"
                        placeholder="请选择"
                      >
                        <a-select-option
                          v-for="option in getMatchOptions(condition.field)"
                          :key="option.value"
                          :value="option.value"
                        >
                          {{ option.label }}
                        </a-select-option>
                      </a-select>
                      <AtInput
                        class="condition-value"
                        :options="variableOptions"
                        :defaultValue="condition.value"
                        placeholder="请输入值，输入/插入变量"
                        @change="(text) => handleValueChange(condition, text)"
                        @open="handleAtOpen"
                      />
                      <div
                        v-if="group.conditions.length > 1"
                        class="btn-hover-wrap"
                        @click="handleDelCondition(index, conditionIndex)"
                      >
                        <CloseCircleOutlined />
                      </div>
                    </div>
                  </div>
                </div>

                <div class="btn-wrap">
                  <a-button
                    @click="handleAddCondition(index)"
                    :icon="h(PlusOutlined)"
                    block
                    type="dashed"
                  >
                    添加条件
                  </a-button>
                </div>
              </div>
            </template>
          </draggable>

          <div class="add-btn-block">
            <a-button @click="handleAddGroup" :icon="h(PlusOutlined)" block type="dashed">
              添加查询条件
            </a-button>
          </div>
        </div>
      </div>

      <div class="node-box-content">
        <div class="node-box-title">
          <img class="input-icon" src="@/assets/img/workflow/output.svg" alt="" />
          <span class="text">输出</span>
        </div>

        <div class="setting-box output-list">
          <div class="options-item">
            <div class="option-label">output_list</div>
            <div class="option-type">array&lt;object&gt;</div>
          </div>
          <div class="options-item">
            <div class="option-label">row_num</div>
            <div class="option-type">number</div>
          </div>
          <div class="output-note">返回商品库全部信息，命中数据最多100条。</div>
        </div>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { h, inject, onMounted, reactive, watch } from 'vue'
import draggable from 'vuedraggable'
import { CloseCircleOutlined, HolderOutlined, PlusOutlined } from '@ant-design/icons-vue'
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import AtInput from '../at-input/at-input.vue'

const props = defineProps({
  node: {
    type: Object,
    default: () => ({})
  },
  lf: {
    type: Object,
    default: null
  },
  nodeId: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close'])
const setData = inject('setData')

const fieldOptions = [
  { label: 'ID', value: 'goods_id', type: 'text' },
  { label: '名称', value: 'goods_name', type: 'text' },
  { label: '类目', value: 'category', type: 'text' },
  { label: '品牌', value: 'brand', type: 'text' },
  { label: '价格', value: 'price', type: 'number' },
  { label: '库存', value: 'stock', type: 'number' }
]

const textMatchOptions = [
  { label: '精准检索', value: 1 },
  { label: '模糊检索', value: 2 }
]

const numberMatchOptions = [
  { label: '大于', value: 3 },
  { label: '小于', value: 4 },
  { label: '等于', value: 5 },
  { label: '大于等于', value: 6 },
  { label: '小于等于', value: 7 }
]

const formState = reactive({
  condition_groups: []
})

const variableOptions = reactive([])

const createCondition = () => ({
  key: Math.random() * 10000,
  field: undefined,
  match: undefined,
  value: ''
})

const createGroup = () => ({
  key: Math.random() * 10000,
  is_or: 0,
  conditions: [createCondition()]
})

const getFieldType = (field) => {
  return fieldOptions.find((item) => item.value === field)?.type || ''
}

const getMatchOptions = (field) => {
  return getFieldType(field) === 'number' ? numberMatchOptions : textMatchOptions
}

const ensureGroups = (groups) => {
  const result = Array.isArray(groups) && groups.length ? groups : [createGroup()]

  return result.map((group) => {
    const conditions =
      Array.isArray(group.conditions) && group.conditions.length
        ? group.conditions
        : [createCondition()]

    return {
      key: group.key || Math.random() * 10000,
      is_or: group.is_or ? 1 : 0,
      conditions: conditions.map((condition) => ({
        key: condition.key || Math.random() * 10000,
        field: condition.field || undefined,
        match: condition.match || undefined,
        value: condition.value || ''
      }))
    }
  })
}

const getOptions = () => {
  const nodeModel = props.lf?.getNodeModelById(props.nodeId)
  const list = nodeModel ? nodeModel.getAllParentVariable() : []

  variableOptions.splice(0, variableOptions.length, ...(list || []))
}

const init = () => {
  const dataRaw = props.node.dataRaw || props.node.node_params || '{}'
  let params = {}

  try {
    params = JSON.parse(dataRaw)
  } catch (error) {
    params = {}
  }

  formState.condition_groups = ensureGroups(params.goods_search?.condition_groups)

  getOptions()
  update()
}

const formatGroups = () => {
  return formState.condition_groups.map((group) => ({
    is_or: !!group.is_or,
    conditions: group.conditions.map((condition) => ({
      field: condition.field || '',
      match: condition.match || '',
      value: condition.value || ''
    }))
  }))
}

const update = () => {
  setData({
    ...props.node,
    node_params: JSON.stringify({
      goods_search: {
        condition_groups: formatGroups()
      }
    })
  })
}

const handleFieldChange = (condition) => {
  const matchOptions = getMatchOptions(condition.field)
  if (!matchOptions.some((item) => item.value === condition.match)) {
    condition.match = matchOptions[0]?.value
  }
}

const handleValueChange = (condition, text) => {
  condition.value = text
}

const handleAddCondition = (groupIndex) => {
  formState.condition_groups[groupIndex].conditions.push(createCondition())
}

const handleDelCondition = (groupIndex, conditionIndex) => {
  const conditions = formState.condition_groups[groupIndex].conditions
  if (conditions.length <= 1) return

  conditions.splice(conditionIndex, 1)
}

const handleAddGroup = () => {
  formState.condition_groups.push(createGroup())
}

const handleDelGroup = (index) => {
  if (formState.condition_groups.length <= 1) return

  formState.condition_groups.splice(index, 1)
}

const handleAtOpen = (open) => {
  if (!open) {
    getOptions()
  }
}

const handleClose = () => {
  emit('close')
}

watch(
  () => formState,
  () => {
    update()
  },
  { deep: true }
)

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import './form-block.less';

.goods-search-form {
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
  .setting-box {
    padding: 0 12px 12px;
  }
  .condition-group {
    margin-bottom: 12px;
    &:last-child {
      margin-bottom: 0;
    }
  }
  .condition-list-box {
    display: flex;
    align-items: center;
  }
  .left-select-box {
    width: 72px;
    flex-shrink: 0;
    ::v-deep(.ant-select-selector) {
      color: #2475fc;
    }
    ::v-deep(.ant-select-arrow) {
      color: #2475fc;
    }
  }
  .condition-body {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 0;
  }
  .condition-item {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .condition-value {
    width: 190px;
  }
  .btn-wrap {
    margin-top: 8px;
    padding-left: 72px;
    padding-right: 40px;
  }
  .add-btn-block {
    margin-top: 8px;
  }
  .output-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .output-note {
    color: #8c8c8c;
    font-size: 13px;
    line-height: 20px;
  }
}
</style>
