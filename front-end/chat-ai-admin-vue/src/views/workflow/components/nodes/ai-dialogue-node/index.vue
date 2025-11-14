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
          <div class="field-item-label">输出</div>
          <div class="field-item-content">
            <div class="field-value is-required">
              <span class="field-key">AI回复内容</span>
              <span class="field-type">string</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { ref, reactive, watch, onMounted, inject, nextTick, onBeforeUnmount } from 'vue'
import { storeToRefs } from 'pinia'
import NodeCommon from '../base-node.vue'
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
  libs_node_key: void 0
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
  let llm = {}
  try {
    llm = JSON.parse(dataRaw).llm || {}
  } catch (e) {
    llm = {}
  }

  getVlaueVariableList()

  const {
    model_config_id,
    use_model,
    context_pair,
    temperature,
    max_token,
    prompt,
    prompt_tags,
    question_value,
    enable_thinking,
    libs_node_key
  } = JSON.parse(JSON.stringify(llm))

  formState.model_config_id = model_config_id
  formState.use_model = use_model
  formState.context_pair = context_pair || 0
  formState.temperature = temperature
  formState.max_token = max_token
  formState.prompt = prompt
  formState.enable_thinking = enable_thinking
  formState.prompt_tags = prompt_tags || []
  formState.libs_node_key = libs_node_key
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
  const data = JSON.stringify({
    llm: {
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

  variableOptions.value = list}

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
  resetSize()
  const mode = graphModel()

  mode.eventCenter.on('custom:setNodeName', onUpatateNodeName)
})

onBeforeUnmount(() => {
  const mode = graphModel()

  mode.eventCenter.off('custom:setNodeName', onUpatateNodeName)
})
</script>
