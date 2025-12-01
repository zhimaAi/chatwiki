<style lang="less" scoped>
.ai-dialogue-node {
  .field-list {
    .field-item {
      display: flex;
      margin-bottom: 8px;
      &:last-child {
        margin-bottom: 0;
      }
    }
    .field-item-label {
      width: 60px;
      line-height: 22px;
      margin-right: 8px;
      font-size: 14px;
      font-weight: 400;
      color: #262626;
      text-align: right;
    }
    .field-item-content {
      flex: 1;
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }
    .field-list-item {
      display: flex;
      gap: 4px;
      align-items: center;
      line-height: 16px;
      padding: 3px 4px;
      border-radius: 4px;
      border: 1px solid #d9d9d9;
      color: #595959;
      background: #fff;
      .right-arrow {
        width: 24px;
        height: 100%;
        border-radius: 4px;
        background: #e4e6eb;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 13px;
      }
      .field-text{
        font-size: 12px;
      }
    }
    .field-value {
      display: flex;
      align-items: center;
      line-height: 16px;
      padding: 3px 4px;
      border-radius: 4px;
      font-size: 12px;
      border: 1px solid #d9d9d9;
      color: #595959;
      background: #fff;

      .field-type {
        padding: 1px 8px;
        margin-left: 4px;
        border-radius: 4px;
        font-size: 12px;
        line-height: 16px;
        font-weight: 400;
        background: #e4e6eb;
      }
    }
  }
}
</style>

<template>
  <node-common
    :title="props.properties.node_name"
    :menus="menus"
    :icon-name="props.properties.node_icon_name"
    :isSelected="props.isSelected"
    :isHovered="props.isHovered"
    :node-key="props.properties.node_key"
    :node_type="props.properties.node_type"
    style="width: 420px;"
  >
    <div class="ai-dialogue-node">
      <div class="field-list">
        <div class="field-item">
          <div class="field-item-label">变量</div>
          <div class="field-item-content">
            <div class="field-list-item" v-for="(item, index) in formState.list" :key="index">
              <div class="field-text">
                <AtText
                  :options="valueOptions"
                  :default-value="item.value"
                  :defaultSelectedList="item.tags"
                />
              </div>
              <div class="right-arrow"><ArrowRightOutlined /></div>
              <div class="field-text">
                {{ getKeyValue(item.variable) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { ref, reactive, watch, onMounted, inject, nextTick, onBeforeUnmount } from 'vue'
import NodeCommon from '../base-node.vue'
import { ArrowRightOutlined } from '@ant-design/icons-vue'
import { haveOutKeyNode } from '@/views/workflow/components/util.js'
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
const graphModel = inject('getGraph')
const getNode = inject('getNode')
const resetSize = inject('resetSize')

// --- State ---
const menus = ref([])
const formState = reactive({
  list: []
})

const valueOptions = ref([])

const reset = () => {
  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  let fields = []
  try {
    let node_params = JSON.parse(dataRaw)
    fields = node_params.assign || []
  } catch (e) {
    console.log(e)
  }

  getVlaueVariableList()


  fields.forEach((item) => {
    item.tags = item.tags || []
  })

  formState.list = fields

  nextTick(() => {
    update()
  })
}

const update = () => {
  let node_params = JSON.parse(props.properties.node_params)

  node_params.assign = [...formState.list]

  setData({
    ...props.node,
    node_params: JSON.stringify(node_params)
  })
}

function getKeyValue(val) {
  if (val) {
    return val.split('.')[1]
  }
  return ''
}
const getVlaueVariableList = () => {
  let list = getNode().getAllParentVariable()

  list.forEach((item) => {
    item.tags = item.tags || []
  })
  valueOptions.value = list
}

const onUpatateNodeName = (data) => {
  if (!haveOutKeyNode.includes(data.node_type)) {
    return
  }

  getVlaueVariableList()
  nextTick(() => {
    formState.list.forEach((item) => {
      if (item.tags && item.tags.length > 0) {
        item.tags.forEach((tag) => {
          if (tag.node_id == data.node_id) {
            let arr = tag.label.split('/')
            arr[0] = data.node_name
            tag.label = arr.join('/')
            tag.node_name = data.node_name
          }
        })
      }
    })
    update()
  })
}

// --- Watchers and Lifecycle Hooks ---
watch(() => props.properties, reset, { deep: true })

onMounted(() => {
  reset()
  resetSize()
  const mode = graphModel()

  mode.eventCenter.on('custom:setNodeName', onUpatateNodeName)
})

onBeforeUnmount(() => {
  const mode = graphModel()

  mode.eventCenter.off('custom:setNodeName', onUpatateNodeName)
})
</script>
