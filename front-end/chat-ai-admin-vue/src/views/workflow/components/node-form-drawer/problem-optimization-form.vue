<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        desc="用大模型优化用户问题，补全潜在信息缺口，提升知识召回率"
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content" @mousedown.stop="">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <div class="gray-block-title">输入</div>
            <a-form-item label="LLM模型" name="use_model">
              <div class="flex-block-item">
                <ModelSelect
                  modelType="LLM"
                  v-model:modeName="formState.use_model"
                  v-model:modeId="formState.model_config_id"
                  @loaded="onVectorModelLoaded"
                  style="width: 348px"
                />
                <!-- <DownOutlined /> -->
                <a-button @click="hanldeShowMore"
                  >高级设置
                  <DownOutlined v-if="showMoreBtn" />
                  <UpOutlined v-else />
                </a-button>
              </div>
            </a-form-item>
            <a-form-item name="temperature" v-if="showMoreBtn">
              <template #label>
                <span>温度&nbsp;</span>
                <a-tooltip>
                  <template #title>温度越低，回答越严谨。温度越高，回答越发散。</template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </template>
              <div class="number-box">
                <div class="number-slider-box">
                  <a-form-item-rest>
                    <a-slider
                      class="custom-slider"
                      v-model:value="formState.temperature"
                      :min="0"
                      :max="2"
                      :step="0.1"
                    />
                  </a-form-item-rest>
                </div>
                <div class="number-input-box">
                  <a-input-number
                    v-model:value="formState.temperature"
                    :min="0"
                    :max="2"
                    :step="0.1"
                  />
                </div>
              </div>
            </a-form-item>
            <a-form-item name="max_token" v-if="showMoreBtn">
              <template #label>
                <span>最大token&nbsp;</span>
                <a-tooltip>
                  <template #title>问题+答案的最大token数，如果出现回答被截断，可调高此值</template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </template>
              <div class="number-box">
                <div class="number-slider-box">
                  <a-form-item-rest>
                    <a-slider
                      class="custom-slider"
                      v-model:value="formState.max_token"
                      :min="0"
                      :max="100 * 1024"
                    />
                  </a-form-item-rest>
                </div>
                <div class="number-input-box">
                  <a-input-number v-model:value="formState.max_token" :min="0" :max="100 * 1024" />
                </div>
              </div>
            </a-form-item>
            <a-form-item name="enable_thinking" v-if="showMoreBtn && show_enable_thinking">
              <template #label>
                <span>深度思考&nbsp;</span>
                <a-tooltip>
                  <template #title>开启时，调用大模型时会指定走深度思考模式</template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </template>
              <div class="number-box">
                <a-switch v-model:checked="formState.enable_thinking" />
              </div>
            </a-form-item>
            <a-form-item name="context_pair">
              <template #label>
                <span>上下文数量&nbsp;</span>
                <a-tooltip>
                  <template #title
                    >提示词中携带的历史聊天记录轮次。设置为0则不携带聊天记录。最多设置50轮。注意，携带的历史聊天记录越多，消耗的token相应也就越多。</template
                  >
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </template>
              <div class="number-box">
                <div class="number-slider-box">
                  <a-form-item-rest>
                    <a-slider
                      class="custom-slider"
                      v-model:value="formState.context_pair"
                      :min="0"
                      :max="50"
                    />
                  </a-form-item-rest>
                </div>
                <div class="number-input-box">
                  <a-input-number v-model:value="formState.context_pair" :min="0" :max="50" />
                </div>
              </div>
            </a-form-item>
            <a-form-item name="prompt" class="width-100">
              <template #label>
                <div class="space-between-box">
                  <div>对话背景</div>
                  <!-- <div class="btn-hover-wrap" @click="onShowAddPromptModal">
                  <DownloadOutlined />从提示词库导入
                </div> -->
                </div>
              </template>
              <at-input
                type="textarea"
                inputStyle="height: 100px;"
                :options="variableOptions"
                :defaultSelectedList="formState.prompt_tags"
                :defaultValue="formState.prompt"
                ref="atInputRef"
                @open="getVlaueVariableList"
                @change="(text, selectedList) => changeValue(text, selectedList)"
                placeholder="描述当前对话的背景，便于大模型补全用户问题，比如：当前对话是关于chatwiki的使用问题和功能介绍"
              >
                <template #option="{ label, payload }">
                  <div class="field-list-item">
                    <div class="field-label">{{ label }}</div>
                    <div class="field-type">{{ payload.typ }}</div>
                  </div>
                </template>
              </at-input>
              <div class="form-tip">输入 / 插入变量</div>
            </a-form-item>
            <div class="diy-form-item">
              <div class="form-label">用户问题</div>
              <div class="form-content">
                <a-cascader
                  v-model:value="formState.question_value"
                  @dropdownVisibleChange="onDropdownVisibleChange"
                  style="width: 220px"
                  :options="variableOptionsSelect"
                  :allowClear="false"
                  :displayRender="({ labels }) => labels.join('/')"
                  :field-names="{ children: 'children' }"
                  placeholder="请选择"
                />
              </div>
            </div>
          </div>
          <div class="gray-block mt16">
            <div class="gray-block-title">输出</div>
            <div class="options-item">
              <div class="option-label">问题优化结果</div>
              <div class="option-type">string</div>
            </div>
          </div>
        </a-form>
        <ImportPrompt @ok="handleSavePrompt" ref="importPromptRef" />
      </div>
    </div>
  </NodeFormLayout>
  
</template>

<script setup>
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import { ref, reactive, watch, computed, nextTick, onMounted } from 'vue'
import { QuestionCircleOutlined, UpOutlined, DownOutlined } from '@ant-design/icons-vue'
import AtInput from '../at-input/at-input.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import ImportPrompt from '@/components/import-prompt/index.vue'
import { useRobotStore } from '@/stores/modules/robot'

const emit = defineEmits(['update-node'])
const props = defineProps({
  lf: {
    type: Object,
    default: null
  },
  nodeId: {
    type: String,
    default: ''
  },
  node: {
    type: Object,
    default: () => ({})
  }
})

const robotStore = useRobotStore()
const atInputRef = ref(null)
const modelList = computed(() => {
  return robotStore.modelList
})
const variableOptions = ref([])
const variableOptionsSelect = ref([])
const showMoreBtn = ref(false)

const hanldeShowMore = () => {
  showMoreBtn.value = !showMoreBtn.value
}

const changeValue = (text, selectedList) => {
  formState.prompt = text
  formState.prompt_tags = selectedList
}

function getOptions() {
  // const node = props.lf.getNodeDataById(props.nodeId)
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let list = nodeModel.getAllParentVariable()

    variableOptionsSelect.value = handleOptions(list)
  }
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

const getVlaueVariableList = () => {
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let list = nodeModel.getAllParentVariable()

    variableOptions.value = list
  }
}

function formatQuestionValue(val) {
  if (val) {
    let lists = val.split('.')
    let str1 = lists[0]
    let str2 = lists.filter((item, index) => index > 0).join('.')
    return [str1, str2]
  }
  return ['global', 'question']
}

const formRef = ref()

const formState = reactive({
  model_config_id: void 0,
  use_model: void 0,
  temperature: 0,
  max_token: 0,
  context_pair: 0,
  prompt: '',
  prompt_tags: [],
  question_value: '',
  enable_thinking: false
})

const update = () => {
  const dataRaw = props.node.dataRaw || props.node.node_params || '{}'
  const oldData = JSON.stringify({
    question_optimize: JSON.parse(dataRaw).question_optimize || {}
  }) 

  const data = JSON.stringify({
    question_optimize: {
      ...formState,
      question_value: formState.question_value.join('.'),
      model_config_id: formState.model_config_id
        ? +formState.model_config_id
        : formState.model_config_id
    }
  })

  if (oldData == data) {
    return
  }
  
  emit('update-node', {
    ...props.node,
    ...formState,
    node_params: data
  })
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let question_optimize = JSON.parse(dataRaw).question_optimize || {}

    question_optimize = JSON.parse(JSON.stringify(question_optimize))

    let {
      model_config_id,
      use_model,
      context_pair,
      temperature,
      max_token,
      prompt,
      prompt_tags,
      question_value,
      enable_thinking
    } = question_optimize
    
    getVlaueVariableList()
    getOptions()

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
  } catch (error) {
    console.log(error)
  }
}

watch(
  () => formState,
  () => {
    update()
  },
  { deep: true }
)


const importPromptRef = ref(null)

const handleSavePrompt = (item) => {
  formState.prompt = ''
  if (item.prompt_type == 1) {
    formState.prompt = item.markdown
  } else {
    formState.prompt = item.prompt
  }
  nextTick(() => {
    atInputRef.value.initData()
  })
}

const onDropdownVisibleChange = (visible) => {
  if (!visible) {
    getOptions()
  }
}

const choosable_thinking = ref({})
const onVectorModelLoaded = (list, choosable_thinking_map) => {
  choosable_thinking.value = choosable_thinking_map
}

const show_enable_thinking = computed(() => {
  if (!formState.model_config_id) {
    return false
  }
  let key = formState.model_config_id + '#' + formState.use_model
  return choosable_thinking.value[key]
})

const handleClose = () => {
  emit('close')
}

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import './form-block.less';

.width-100 {
  ::v-deep(.ant-form-item-label) {
    width: 100%;
    label {
      width: 100%;
    }
  }
}
.space-between-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  .btn-hover-wrap {
    width: fit-content;
    padding: 0 6px;
    color: #2475fc;
    gap: 4px;
  }
}

.options-item {
  margin-top: 12px;
  height: 22px;
  line-height: 22px;
  display: flex;
  align-items: center;
  gap: 8px;
  .option-label {
    color: var(--wf-color-text-1);
    font-size: 14px;
    &::before {
      content: '*';
      color: #fb363f;
      display: inline-block;
      margin-right: 2px;
    }
  }
  .option-type {
    height: 22px;
    width: fit-content;
    padding: 0 8px;
    border-radius: 6px;
    border: 1px solid rgba(0, 0, 0, 0.15);
    background-color: #fff;
    color: var(--wf-color-text-3);
    font-size: 12px;
    display: flex;
    align-items: center;
  }
}
.number-box {
  display: flex;
  align-items: center;

  .number-slider-box {
    width: 244px;
  }

  .number-input-box {
    margin-left: 24px;
  }
}
.model-icon {
  height: 18px;
}
</style>
