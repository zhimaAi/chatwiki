<style lang="less" scoped>
@import '../form-block.less';

.goods-search-node {
  .condition-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 0;
    width: 100%;
  }
  .if-box {
    display: flex;
    align-items: center;
    min-width: 260px;
  }
  .if-box-label {
    width: 72px;
    line-height: 22px;
    color: #262626;
    flex-shrink: 0;
  }
  .if-box-content {
    flex: 1;
    padding: 8px;
    border-radius: 6px;
    overflow: hidden;
    background: #f2f4f7;
  }
  .condition-box {
    position: relative;
    display: flex;
    align-items: center;
  }
  .condition-left-box {
    margin-right: 4px;
  }
  .connection-text {
    display: inline-block;
    height: 18px;
    line-height: 18px;
    padding: 0 8px;
    font-size: 12px;
    font-weight: 400;
    border-radius: 4px;
    color: #595959;
    background: #e4e6eb;
  }
  .condition-line {
    width: 24px;
    height: 100%;
    margin-right: 4px;

    &::after {
      content: '';
      position: absolute;
      width: 24px;
      top: 12px;
      bottom: 12px;
      border: 1px solid #d9d9d9;
      border-right: 0;
      border-radius: 2px;
      border-top-right-radius: 0;
      border-bottom-right-radius: 0;
    }
  }
  .condition-body {
    flex: 1;
    min-width: 0;
    overflow: hidden;
  }
  .field-item {
    display: flex;
    align-items: center;
    min-height: 24px;
    line-height: 16px;
    padding: 2px 4px;
    margin-bottom: 2px;
    border-radius: 4px;
    border: 1px solid #d9d9d9;
    background: #fff;

    &:last-child {
      margin-bottom: 0;
    }
  }
  .field-name,
  .field-value {
    flex: 1;
    min-width: 0;
    line-height: 16px;
    font-size: 12px;
    font-weight: 400;
    color: #595959;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
  .field-rule {
    height: 18px;
    line-height: 18px;
    padding: 0 8px;
    margin: 0 4px;
    font-size: 12px;
    font-weight: 400;
    border-radius: 4px;
    color: #595959;
    background: #e4e6eb;
    flex-shrink: 0;
  }
}
</style>

<template>
  <node-common
    :properties="properties"
    :title="props.properties.node_name"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
  >
    <div class="goods-search-node">
      <div class="static-field-list">
        <div class="static-field-item">
          <!-- <div class="static-field-item-label">输入</div> -->
          <div class="static-field-item-content">
            <div class="condition-list">
              <template v-if="conditionGroups.length">
                <div
                  v-for="(group, groupIndex) in conditionGroups"
                  :key="group.key || groupIndex"
                  class="if-box"
                >
                  <div class="if-box-label">查询条件{{ groupIndex + 1 }}</div>
                  <div class="if-box-content">
                    <div class="condition-box">
                      <div class="condition-left-box" v-if="getConditionList(group).length > 1">
                        <span class="connection-text">{{ group.is_or ? 'or' : 'and' }}</span>
                      </div>
                      <div class="condition-line" v-if="getConditionList(group).length > 1"></div>
                      <div class="condition-body">
                        <div class="field-items" v-if="getConditionList(group).length">
                          <div
                            class="field-item"
                            v-for="(condition, conditionIndex) in getConditionList(group)"
                            :key="condition.key || conditionIndex"
                          >
                            <span style="max-width: 40px;" class="field-name">{{ getLabel(fieldOptions, condition.field) }}</span>
                            <span class="field-rule">{{ getLabel(matchOptions, condition.match) }}</span>
                            <span class="field-value">
                              <AtText
                                ref="atTextRefs"
                                :options="variableOptions"
                                :default-value="condition.value || '--'"
                              />
                            </span>
                          </div>
                        </div>
                        <span v-else class="static-field-value">--</span>
                      </div>
                    </div>
                  </div>
                </div>
              </template>
              <span v-else class="static-field-value">--</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { computed, inject, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import NodeCommon from '../base-node.vue'
import AtText from '../../at-input/at-text.vue'

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  },
  isSelected: { type: Boolean, default: false },
  isHovered: { type: Boolean, default: false }
})

const setData = inject('setData')
const resetSize = inject('resetSize')
const getNode = inject('getNode')
const graphModel = inject('getGraph')

const fieldOptions = [
  { label: 'ID', value: 'goods_id' },
  { label: '名称', value: 'goods_name' },
  { label: '类目', value: 'category' },
  { label: '品牌', value: 'brand' },
  { label: '价格', value: 'price' },
  { label: '库存', value: 'stock' }
]

const matchOptions = [
  { label: '精准检索', value: 1 },
  { label: '模糊检索', value: 2 },
  { label: '大于', value: 3 },
  { label: '小于', value: 4 },
  { label: '等于', value: 5 },
  { label: '大于等于', value: 6 },
  { label: '小于等于', value: 7 }
]

const formState = ref({
  condition_groups: []
})

const atTextRefs = ref([])
const variableOptions = ref([])

const conditionGroups = computed(() => formState.value.condition_groups || [])

const getLabel = (list, value) => {
  const item = list.find((option) => option.value === value)
  return item ? item.label : '--'
}

const getConditionList = (group) => Array.isArray(group.conditions) ? group.conditions : []

const getVariableOptions = () => {
  variableOptions.value = getNode().getAllParentVariable() || []
}

const init = () => {
  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  let params = {}

  try {
    params = JSON.parse(dataRaw)
  } catch (error) {
    params = {}
  }

  getVariableOptions()
  formState.value.condition_groups = params.goods_search?.condition_groups || []

  update()

  nextTick(() => {
    resetSize()
  })
}

const update = () => {
  setData({
    ...props.properties,
    node_params: JSON.stringify({
      goods_search: {
        condition_groups: formState.value.condition_groups
      }
    })
  })
}

const handleNodeNameUpdate = () => {
  getVariableOptions()

  nextTick(() => {
    atTextRefs.value.forEach((atTextRef) => atTextRef?.refresh())
  })
}

watch(
  () => props.properties.dataRaw,
  (newVal, oldVal) => {
    if (newVal !== oldVal) {
      init()
    }
  }
)

onMounted(() => {
  init()

  const mode = graphModel()
  mode.eventCenter.on('custom:setNodeName', handleNodeNameUpdate)
})

onBeforeUnmount(() => {
  const mode = graphModel()
  mode.eventCenter.off('custom:setNodeName', handleNodeNameUpdate)
})
</script>
