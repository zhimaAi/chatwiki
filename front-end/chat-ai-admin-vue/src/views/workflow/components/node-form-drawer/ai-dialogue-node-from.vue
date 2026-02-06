<template>
  <NodeFormLayout>
    <template #header>
       <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        :desc="t('desc_call_llm_generate_reply')"
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content" @mousedown.stop="">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <div class="gray-block-title">{{ t('title_llm_settings') }}</div>
            <a-form-item :label="t('label_llm_model')" name="use_model">
              <div class="flex-block-item">
                <ModelSelect
                  modelType="LLM"
                  v-model:modeName="formState.use_model"
                  v-model:modeId="formState.model_config_id"
                  @loaded="onVectorModelLoaded"
                  @change="handleModelChange"
                  style="width: 348px"
                />
                <!-- <DownOutlined /> -->
                <a-button @click="hanldeShowMore"
                  >{{ t('btn_advanced_settings') }}
                  <DownOutlined v-if="showMoreBtn" />
                  <UpOutlined v-else />
                </a-button>
              </div>
            </a-form-item>

            <a-form-item name="role" v-if="showMoreBtn">
              <template #label>
                <span>{{ t('label_prompt_role') }}&nbsp;</span>
                <a-tooltip>
                  <template #title>
                    <div>
                      {{ t('tip_prompt_role_desc') }}
                      <div>{{ t('label_system_role_example') }}</div>
                      <div>{"role": "system", "content": "{{ t('sample_custom_prompt') }}"}</div>
                    </div>
                  </template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </template>
              <a-radio-group v-model:value="formState.role">
                <a-radio :value="1">{{ t('label_system_role') }}</a-radio>
                <a-radio :value="2">{{ t('label_user_role') }}</a-radio>
              </a-radio-group>
            </a-form-item>

            <a-form-item name="temperature" v-if="showMoreBtn">
              <template #label>
                <span>{{ t('label_temperature') }}&nbsp;</span>
                <a-tooltip>
                  <template #title>{{ t('tip_temperature_desc') }}</template>
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
                <span>{{ t('label_max_token') }}&nbsp;</span>
                <a-tooltip>
                  <template #title>{{ t('tip_max_token_desc') }}</template>
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
                <span>{{ t('label_deep_thinking') }}&nbsp;</span>
                <a-tooltip>
                  <template #title>{{ t('tip_deep_thinking_desc') }}</template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </template>
              <div class="number-box">
                <a-switch v-model:checked="formState.enable_thinking" />
              </div>
            </a-form-item>
            <a-form-item name="context_pair">
              <template #label>
                <span>{{ t('label_context_pair') }}&nbsp;</span>
                <a-tooltip>
                  <template #title
                    >{{ t('tip_context_pair_desc') }}</template
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
                  <div>{{ t('label_prompt') }}</div>
                  <a-flex :gap="8">
                    <div class="btn-hover-wrap" @click="onShowAddPromptModal">
                      <DownloadOutlined />{{ t('btn_import_from_prompt_library') }}
                    </div>
                    <div class="btn-hover-wrap" @click="handleOpenFullAtModal">
                      <FullscreenOutlined />
                    </div>
                  </a-flex>
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
                :placeholder="t('ph_input_message_content')"
              >
                <template #option="{ label, payload }">
                  <div class="field-list-item">
                    <div class="field-label">{{ label }}</div>
                    <div class="field-type">{{ payload.typ }}</div>
                  </div>
                </template>
              </at-input>
              <div class="form-tip">{{ t('msg_input_insert_variable') }}</div>
            </a-form-item>

            <div class="diy-form-item mt12">
              <div class="form-label">{{ t('label_knowledge_base_reference') }}</div>
              <div class="form-content">
                <a-select
                  :placeholder="t('ph_select')"
                  allowClear
                  @dropdownVisibleChange="onDropdownVisibleChange"
                  v-model:value="formState.libs_node_key"
                  style="width: 220px"
                >
                  <a-select-option
                    v-for="item in knowledgeQuoteOptions"
                    :value="item.node_id"
                    :key="item.node_id"
                    >{{ item.label }}</a-select-option
                  >
                </a-select>
              </div>
            </div>
          </div>

           <div class="gray-block mt16">
            <div class="gray-block-title">{{ t('label_input') }}</div>
            <div class="diy-form-item question-value-item">
              <div class="form-label">{{ t('label_input') }}</div>
              <div class="form-content">
                <a-cascader
                  v-model:value="formState.question_value"
                  @dropdownVisibleChange="onDropdownVisibleChange"
                  style="width: 220px"
                  :options="variableOptionsSelect"
                  :allowClear="false"
                  :displayRender="({ labels }) => labels.join('/')"
                  :field-names="{ children: 'children' }"
                  :placeholder="t('ph_select')"
                />
              </div>
            </div>
          </div>

          <div class="gray-block mt16">
            <div class="gray-block-title">{{ t('label_output') }}</div>
            <div class="options-item">
              <div class="option-label">{{ t('label_ai_reply_content') }}</div>
              <div class="option-type">string</div>
            </div>
          </div>
        </a-form>
      </div>

      <ImportPrompt @ok="handleSavePrompt" ref="importPromptRef" />
      <FullAtInput
        :options="variableOptions"
        :defaultSelectedList="formState.prompt_tags"
        :defaultValue="formState.prompt"
        :placeholder="t('ph_input_message_content')"
        type="textarea"
        @open="getVlaueVariableList"
        @change="(text, selectedList) => changeValue(text, selectedList)"
        @ok="handleRefreshAtInput"
        ref="fullAtInputRef"
      />
    </div>
  </NodeFormLayout>

</template>

<script setup>
import { useRobotStore } from '@/stores/modules/robot'
import { ref, reactive, watch, computed, onMounted, nextTick, inject } from 'vue'
import { QuestionCircleOutlined, UpOutlined, DownOutlined, FullscreenOutlined, DownloadOutlined } from '@ant-design/icons-vue'
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import AtInput from '../at-input/at-input.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import FullAtInput from '../at-input/full-at-input.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.ai-dialogue-node-from')


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

const getNode = inject('getNode')

const robotStore = useRobotStore()
const atInputRef = ref(null)
const modelList = computed(() => {
  return robotStore.modelList
})
const variableOptions = ref([])
const variableOptionsSelect = ref([])
const showMoreBtn = ref(false)
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
  enable_thinking: false,
  libs_node_key: void 0,
  role: 1
})

const handleModelChange = () => {
  if (formState.use_model && formState.use_model.toLowerCase().includes('deepseek-r1')) {
    formState.role = 2
  } else {
    formState.role = 1
  }
}
const hanldeShowMore = () => {
  showMoreBtn.value = !showMoreBtn.value
}

const changeValue = (text, selectedList) => {
  formState.prompt = text
  formState.prompt_tags = selectedList
}

function getOptions() {
  let list = getNode().getAllParentVariable()

  variableOptionsSelect.value = handleOptions(list)
}

const knowledgeQuoteOptions = computed(() => {
  let list = variableOptionsSelect.value.filter((item) => item.node_type == 5)
  return list
})

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
  let list = getNode().getAllParentVariable()
  variableOptions.value = list
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

const importPromptRef = ref(null)

const onShowAddPromptModal = () => {
  importPromptRef.value.show()
}

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

const fullAtInputRef = ref(null)
const handleOpenFullAtModal = () => {
  fullAtInputRef.value.show()
}

const handleRefreshAtInput = () => {
  atInputRef.value.refresh()
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

  emit('update-node', {
    ...props.node,
    ...formState,
    node_params: data
  })
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let llm = JSON.parse(dataRaw).llm || {}

    llm = JSON.parse(JSON.stringify(llm))
    let {
      model_config_id,
      use_model,
      context_pair,
      temperature,
      max_token,
      prompt,
      prompt_tags,
      question_value,
      enable_thinking,
      libs_node_key,
      role
    } = llm

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
    formState.libs_node_key = libs_node_key
    formState.question_value = formatQuestionValue(question_value)
    if (!formState.model_config_id && modelList.value.length > 0) {
      formState.model_config_id = modelList.value[0].id
      formState.use_model = modelList.value[0].children[0].name
    }
    formState.role = +role || 1
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
.question-value-item{
  display: flex;
  align-items: center;
  .form-label{
    margin-right: 8px;
    &::before {
      content: '*';
      color: #fb363f;
      display: inline-block;
      margin-right: 2px;
    }
  }
  .form-content{
    margin-top: 0;
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