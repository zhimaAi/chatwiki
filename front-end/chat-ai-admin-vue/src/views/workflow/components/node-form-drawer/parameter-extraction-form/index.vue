<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        :desc="t('desc_param_extraction')"
        @close="handleClose"
      ></NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <div class="gray-block-title">{{ t('label_input') }}</div>
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
                  <template #title>{{ t('tip_temperature') }}</template>
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
                  <template #title>{{ t('tip_max_token') }}</template>
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
                  <template #title>{{ t('tip_deep_thinking') }}</template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </template>
              <div class="number-box">
                <a-switch v-model:checked="formState.enable_thinking" />
              </div>
            </a-form-item>
            <a-form-item name="context_pair">
              <template #label>
                <span>{{ t('label_context_count') }}&nbsp;</span>
                <a-tooltip>
                  <template #title>{{ t('tip_context_count') }}</template>
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
                </div>
              </template>
              <!-- {{ formState.prompt }} -->
              <at-input
                type="textarea"
                inputStyle="height: 100px;"
                :options="variableOptions"
                :defaultSelectedList="formState.prompt_tags"
                :defaultValue="formState.prompt"
                ref="atInputRef"
                @open="getVlaueVariableList"
                @change="(text, selectedList) => changeValue(text, selectedList)"
                :placeholder="t('ph_input_message')"
              >
                <template #option="{ label, payload }">
                  <div class="field-list-item">
                    <div class="field-label">{{ label }}</div>
                    <div class="field-type">{{ payload.typ }}</div>
                  </div>
                </template>
              </at-input>
              <div class="form-tip">{{ t('tip_input_variable') }}</div>
            </a-form-item>
            <div class="diy-form-item">
              <div class="form-label">{{ t('label_user_question') }}</div>
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
            <div class="gray-block-title" @click="test">{{ t('label_output_extract') }}</div>
            <div class="output-box">
              <div class="output-block">
                <div class="output-item">{{ t('label_param_key') }}</div>
                <div class="output-item">{{ t('label_type') }}</div>
              </div>
              <div class="array-form-box">
                <div class="form-item-list" v-for="(item, index) in formState.output" :key="index">
                  <a-form-item :label="null" :name="['output', index, 'key']">
                    <div class="flex-block-item" style="gap: 12px">
                      <a-input
                        style="width: 214px"
                        v-model:value="item.key"
                        :placeholder="t('ph_input')"
                      ></a-input>
                      <a-form-item-rest>
                        <a-select
                          @change="onTypeChange(item)"
                          v-model:value="item.typ"
                          :placeholder="t('ph_select')"
                          style="width: 214px"
                        >
                          <a-select-option
                            v-for="op in typOptions"
                            :value="op.value"
                            :key="op.value"
                            >{{ op.value }}</a-select-option
                          >
                        </a-select>
                      </a-form-item-rest>

                      <div
                        class="btn-hover-wrap"
                        v-if="item.typ == 'object'"
                        @click="onAddSubs(index)"
                      >
                        <PlusCircleOutlined />
                      </div>
                      <div class="btn-hover-wrap" @click="handleEditOutput(item, index)">
                        <EditOutlined />
                      </div>
                      <div class="btn-hover-wrap" @click="onDelOutput(index)">
                        <CloseCircleOutlined />
                      </div>
                    </div>
                    <div class="sub-field-box" v-if="item.subs && item.subs.length > 0">
                      <a-form-item-rest>
                        <SubKey :data="item.subs" :level="2" :typOptions="typOptions" />
                      </a-form-item-rest>
                    </div>
                  </a-form-item>
                </div>
                <a-button @click="handleAddOutPut" :icon="h(PlusOutlined)" block type="dashed"
                  >{{ t('btn_add_param') }}</a-button
                >
              </div>
            </div>
          </div>
        </a-form>
      </div>
      <AddParamsModal
        @add="onOutputAdd"
        @edit="onOutputEdit"
        @addSub="onOutputAddSub"
        ref="addParamsModalRef"
      />
    </div>

    
  </NodeFormLayout>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { useRobotStore } from '@/stores/modules/robot'
import { ref, reactive, watch, h, computed, onMounted } from 'vue'
import {
  PlusOutlined,
  PlusCircleOutlined,
  DownOutlined,
  QuestionCircleOutlined,
  UpOutlined,
  CloseCircleOutlined,
  EditOutlined
} from '@ant-design/icons-vue'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import SubKey from './subs-key.vue'
import AtInput from '../../at-input/at-input.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import AddParamsModal from './add-params-modal.vue'

const { t } = useI18n('views.workflow.components.node-form-drawer.parameter-extraction-form.index')

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
  question_value: '',
  enable_thinking: false,
  output: [],
  role: 1
})

const handleModelChange = () => {
  if (formState.use_model && formState.use_model.toLowerCase().includes('deepseek-r1')) {
    formState.role = 2
  } else {
    formState.role = 1
  }
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

  emit('update-node', {
    ...props.node,
    ...formState,
    node_params: data
  })
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let params_extractor = JSON.parse(dataRaw).params_extractor || {}

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
    formState.role = +params_extractor.role || 1
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

const addParamsModalRef = ref(null)
const handleAddOutPut = () => {
  addParamsModalRef.value.add()
}

const handleEditOutput = (data, index) => {
  addParamsModalRef.value.edit(data, index)
}

const onOutputAdd = (data) => {
  formState.output.push(data)
}

const onOutputEdit = (data, index) => {
  formState.output.splice(index, 1, data)
}

const onDelOutput = (index) => {
  formState.output.splice(index, 1)
}

const onTypeChange = (data) => {
  data.subs = []
}

const onAddSubs = (index) => {
  addParamsModalRef.value.addSub(index)
}
const onOutputAddSub = (data, index) => {
  formState.output[index].subs.push(data)
}

const typOptions = [
  {
    lable: 'string',
    value: 'string'
  },
  {
    lable: 'number',
    value: 'number'
  },
  {
    lable: 'boole',
    value: 'boole'
  },
  {
    lable: 'float',
    value: 'float'
  },
  {
    lable: 'object',
    value: 'object'
  },
  {
    lable: 'array\<string>',
    value: 'array\<string>'
  },
  {
    lable: 'array\<number>',
    value: 'array\<number>'
  },
  {
    lable: 'array\<boole>',
    value: 'array\<boole>'
  },
  {
    lable: 'array\<float>',
    value: 'array\<float>'
  },
  {
    lable: 'array\<object>',
    value: 'array\<object>'
  }
]

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
@import '../form-block.less';
.output-box {
  .output-block {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 4px;
    color: #262626;
    .output-item {
      width: 214px;
    }
  }
  .flex-block-item .btn-hover-wrap {
    width: 24px;
    height: 24px;
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
