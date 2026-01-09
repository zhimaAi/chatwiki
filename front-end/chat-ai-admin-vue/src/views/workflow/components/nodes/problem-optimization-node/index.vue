<style lang="less" scoped>
@import '../form-block.less';
.ai-dialogue-node {
  
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
    style="width: 420px;"
  >
    <div class="ai-dialogue-node">
      <div class="static-field-list">
        <div class="static-field-item">
          <div class="static-field-item-label">用户问题</div>
          <div class="static-field-item-content">
            <div class="static-field-value">
              <user-question-text ref="questionTextRef" :value="formState.question_value" />
            </div>
          </div>
        </div>

        <div class="static-field-item">
          <div class="static-field-item-label">LLM模型</div>
          <div class="static-field-item-content">
            <div class="static-field-value">
              <model-name-text :useModel="formState.use_model" :modelConfigId="formState.model_config_id" />
            </div>
          </div>
        </div>

        <div class="static-field-item">
          <div class="static-field-item-label">输出</div>
          <div class="static-field-item-content">
            <div class="static-field-value">
              <span class="static-field-key">问题优化结果</span>
              <span class="static-field-type">string</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { useRobotStore } from '@/stores/modules/robot'
import { ref, reactive, watch, nextTick, onMounted, inject, onBeforeUnmount } from 'vue'
import { storeToRefs } from 'pinia'
import NodeCommon from '../base-node.vue'
import ModelNameText from '../model-name-text.vue'
import UserQuestionText from '../user-question-text.vue'

// --- Props and Injections ---
const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  },
  isSelected: { type: Boolean, default: false },
  isHovered: { type: Boolean, default: false }
})
const questionTextRef = ref(null)
const setData = inject('setData')
const graphModel = inject('getGraph')
const getNode = inject('getNode')
const resetSize = inject('resetSize')

// --- Store ---
const robotStore = useRobotStore()
const { modelList } = storeToRefs(robotStore)

// --- State ---
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
  enable_thinking: false
})

const variableOptions = ref([])

function formatQuestionValue(val) {
  if (val) {
    let lists = val.split('.')
    let str1 = lists[0]
    let str2 = lists.filter((item, index) => index > 0).join('.')
    return [str1, str2]
  }
  return ['global', 'question']
}

const reset = () => {
  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  
  let question_optimize = {}
  try {
    question_optimize = JSON.parse(dataRaw).question_optimize || {}
  } catch (e) {
    question_optimize = {}
  }

  getVlaueVariableList()

  // Use a deep copy to avoid side effects
  const {
    model_config_id,
    use_model,
    context_pair,
    temperature,
    max_token,
    prompt,
    prompt_tags,
    question_value,
    enable_thinking
  } = JSON.parse(JSON.stringify(question_optimize))

  formState.model_config_id = model_config_id
  formState.use_model = use_model
  formState.context_pair = context_pair || 0
  formState.temperature = temperature
  formState.max_token = max_token
  formState.prompt = prompt
  formState.enable_thinking = enable_thinking
  formState.prompt_tags = prompt_tags || []
  formState.question_value = formatQuestionValue(question_value)

  if (!formState.model_config_id && modelList.value.length > 0) {
    formState.model_config_id = modelList.value[0].id
    formState.use_model = modelList.value[0].children[0].name
  }

  nextTick(() => {
    resetSize()
  })
}

const update = () => {
  const data = {
    ...formState,
    node_params: JSON.stringify({
      question_optimize: {
        ...formState,
        question_value: formState.question_value.join('.'),
        model_config_id: formState.model_config_id
          ? +formState.model_config_id
          : formState.model_config_id
      }
    })
  }

  setData(data)
}

const getVlaueVariableList = () => {
  let list = getNode().getAllParentVariable()

  variableOptions.value = list
}

const onUpatateNodeName = (data) => {
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
watch(() => props.properties, reset, { deep: true })

onMounted(() => {
  reset()

  nextTick(() => {
    resetSize()
  })

  const mode = graphModel()

  mode.eventCenter.on('custom:setNodeName', onUpatateNodeName)
})

onBeforeUnmount(() => {
  const mode = graphModel()

  mode.eventCenter.off('custom:setNodeName', onUpatateNodeName)
})
</script>
