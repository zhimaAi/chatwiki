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
      gap: 4px;
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

      &.is-required .field-key::before{
        content: '*';
        color: #FB363F;
        display: inline-block;
        margin-right: 2px;
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
          <div class="field-item-label">用户问题</div>
          <div class="field-item-content">
            <div class="field-value">
              <user-question-text ref="questionTextRef" :value="formState.question_value" />
            </div>
          </div>
        </div>

        <div class="field-item">
          <div class="field-item-label">LLM模型</div>
          <div class="field-item-content">
            <div class="field-value">
              <span class="field-key">
                <model-name-text :useModel="formState.use_model" :modelConfigId="formState.model_config_id" />
              </span>
            </div>
          </div>
        </div>

        <div class="field-item">
          <div class="field-item-label">输出字段</div>
          <div class="field-item-content">
            <div class="field-value" :class="{ 'is-required': item.required }" v-for="item in formState.output" :key="item.cu_key">
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
import {
  ref,
  reactive,
  watch,
  inject,
  nextTick,
  onMounted,
  onBeforeUnmount
} from 'vue'
import { storeToRefs } from 'pinia'
import NodeCommon from '../base-node.vue'
import { haveOutKeyNode } from '@/views/workflow/components/util.js'
import { useRobotStore } from '@/stores/modules/robot'
import ModelNameText from '../model-name-text.vue'
import UserQuestionText from '../user-question-text.vue'

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

// --- Store ---
const robotStore = useRobotStore()
const { modelList } = storeToRefs(robotStore)

// --- State ---
const questionTextRef = ref(null)
const menus = ref([])
const formState = reactive({
  model_config_id: undefined,
  use_model: undefined,
  temperature: 0,
  max_token: 0,
  context_pair: 0,
  prompt: '',
  prompt_tags: [],
  question_value: [],
  enable_thinking: false,
  output: []
})

const variableOptions = ref([])
const variableOptionsSelect = ref([])

function formatQuestionValue(val) {
  if (val) {
    let lists = val.split('.')
    let str1 = lists[0]
    let str2 = lists.filter((item, index) => index > 0).join('.')
    return [str1, str2]
  }
  return ['global', 'question']
}

function recursionData(data) {
  data.forEach((item) => {
    item.cu_key = Math.random() * 10000
    if (item.subs && item.subs.length) {
      recursionData(item.subs)
    } else {
      item.subs = []
    }
  })
  return data
}

const reset = () => {
  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  let params_extractor = {}
  try {
    params_extractor = JSON.parse(dataRaw).params_extractor || {}
  } catch (e) {
    params_extractor = {}
  }
  params_extractor = JSON.parse(JSON.stringify(params_extractor))

  getVlaueVariableList()
  getOptions()

  for (let key in params_extractor) {
    if (key == 'output') {
      formState['output'] = recursionData(params_extractor[key])
      continue
    }
    if (key == 'question_value') {
      formState.question_value = formatQuestionValue(params_extractor[key])
      continue
    }
    formState[key] = params_extractor[key]
  }

  if (!formState.model_config_id && modelList.value.length > 0) {
    formState.model_config_id = modelList.value[0].id
    formState.use_model = modelList.value[0].children[0].name
  }

  nextTick(() => {
    resetSize()
  })
}

const update = () => {
  const data = JSON.stringify({
    params_extractor: {
      ...formState,
      question_value: formState.question_value.join('.'),
      model_config_id: formState.model_config_id
        ? +formState.model_config_id
        : formState.model_config_id
    }
  })

  setData({
    ...props.node,
    ...formState,
    node_params: data
  })
}

const getVlaueVariableList = () => {
  let list = getNode().getAllParentVariable()
  list.forEach((item) => {
    item.tags = item.tags || []
  })

  variableOptions.value = list
}

function getOptions() {
  let list = getNode().getAllParentVariable()

  variableOptionsSelect.value = handleOptions(list)
}

// 递归处理Options
function handleOptions(options) {
  options.forEach((item) => {
    if (item.typ == 'node') {
      if (item.node_type == 1) {
        item.value = 'global'
      } else {
        item.value = item.node_id
      }
    } else {
      item.value = item.key
    }

    if (item.children && item.children.length > 0) {
      item.children = handleOptions(item.children)
    }
  })

  return options
}

const onUpatateNodeName = (data) => {
  if (!haveOutKeyNode.includes(data.node_type)) {
    return
  }
  getVlaueVariableList()

  if (formState.prompt_tags && formState.prompt_tags.length > 0) {
    formState.prompt_tags.forEach((tag) => {
      if (tag.node_id == data.node_id) {
        let arr = tag.label.split('/')
        arr[0] = data.node_name
        tag.label = arr.join('/')
        tag.node_name = data.node_name
      }
    })
  }

  update()

  questionTextRef.value.refresh()
}

// --- Watchers and Lifecycle Hooks ---
watch(() => props.properties.dataRaw, (newVal, oldVal) => {
  if(newVal != oldVal) {
    reset()
  }
}, { deep: true })

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
