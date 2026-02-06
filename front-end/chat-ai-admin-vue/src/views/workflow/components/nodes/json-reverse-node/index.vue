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
      width: auto;
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
      .field-text {
        max-width: 150px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
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
      .field-key {
        max-width: 200px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
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
    :properties="properties"
    :title="props.properties.node_name"
    :menus="menus"
    :icon-name="props.properties.node_icon_name"
    :isSelected="props.isSelected"
    :isHovered="props.isHovered"
    :node-key="props.properties.node_key"
    :node_type="props.properties.node_type"
    style="width: 420px"
  >
    <div class="ai-dialogue-node">
      <div class="field-list">
        <div class="field-item">
          <div class="field-item-label">{{ t('label_input_params') }}</div>
          <div class="field-item-content">
            <div class="field-list-item">
              <div class="field-text">
                <UserQuestionText :value="formState.input_variable" />
              </div>
            </div>
          </div>
        </div>

        <div class="field-item">
          <div class="field-item-label">{{ t('label_output_fields') }}</div>
          <div class="field-item-content">
            <div class="field-value" v-for="(item, index) in formState.output" :key="index">
              <span class="field-key"> {{ item.key }}</span>
              <span class="field-type">{{ item.typ }}</span>
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
import UserQuestionText from '../user-question-text.vue'
import { haveOutKeyNode, formatSpacialKey } from '@/views/workflow/components/util.js'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.nodes.json-reverse-node.index')

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
  input_variable: [],
  output: [
    {
      key: '',
      typ: 'object',
      subs: []
    }
  ]
})

const reset = () => {
  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  let json_decode = {}
  try {
    json_decode = JSON.parse(dataRaw).json_decode || {}
  } catch (e) {
    json_decode = {}
  }

  json_decode = JSON.parse(JSON.stringify(json_decode))
  formState.input_variable = formatSpacialKey(json_decode.input_variable)
  formState.output = json_decode.output

  nextTick(() => {
    resetSize()
  })
}

const update = () => {
  const data = JSON.stringify({
    json_decode: {
      input_variable: formState.input_variable ? formState.input_variable.join('.') : '',
      ...formState
    }
  })

  setData({
    ...props.node,
    node_params: data
  })
}

const onUpatateNodeName = (data) => {
  if (!haveOutKeyNode.includes(data.node_type)) {
    return
  }

  nextTick(() => {
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
