<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        :desc="t('desc_question_classification')"
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content" @mousedown.stop="">
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
                  <template #title>{{ t('tip_context_pair_desc') }}</template>
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
            <div class="gray-block-title">{{ t('label_question_classification_settings') }}</div>
            <div class="array-form-box">
              <div
                class="form-item-list"
                v-for="(item, index) in formState.categorys"
                :key="item.key"
              >
                <a-form-item
                  :label="null"
                  :name="['categorys', index, 'category']"
                  :rules="{ required: true, validator: (rule, value) => checkedHeader(rule, value) }"
                >
                  <div class="flex-block-item">
                    <a-input class="flex1" v-model:value="item.category" :placeholder="t('ph_input')"></a-input>
                    <div class="btn-hover-wrap" @click="onDelcategory(index)">
                      <CloseCircleOutlined />
                    </div>
                  </div>
                </a-form-item>
              </div>

              <div class="form-item-list">
                <a-form-item :label="null">
                  <div class="flex-block-item">
                    <a-input class="flex1" :value="t('label_default_category')" readonly></a-input>
                  </div>
                </a-form-item>
              </div>
              <a-button @click="handleAddcategory" :icon="h(PlusOutlined)" block type="dashed"
                >{{ t('btn_add_question_category') }}</a-button
              >
            </div>
          </div>
        </a-form>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { getUuid } from '@/utils/index'
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import { ref, reactive, watch, computed, onMounted, h, inject  } from 'vue'
import {
  CloseCircleOutlined,
  QuestionCircleOutlined,
  UpOutlined,
  DownOutlined,
  PlusOutlined
} from '@ant-design/icons-vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { useRobotStore } from '@/stores/modules/robot'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.question-node-form')

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
const modelList = computed(() => {
  return robotStore.modelList
})
const showMoreBtn = ref(false)

const hanldeShowMore = () => {
  showMoreBtn.value = !showMoreBtn.value
}

// const variableOptions = ref([])

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
  categorys: [],
  role: 1
})


const handleModelChange = () => {
  if (formState.use_model && formState.use_model.toLowerCase().includes('deepseek-r1')) {
    formState.role = 2
  } else {
    formState.role = 1
  }
}

const variableOptionsSelect = ref([])

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

const onDropdownVisibleChange = (visible) => {
  if (!visible) {
    getOptions()
  }
}

function formatQuestionValue(val) {
  if (val ) {
    let lists = val.split('.')
    let str1 = lists[0]
    let str2 = lists.filter((item, index) => index > 0).join('.')
    return [str1, str2]
  }
  return ['global', 'question']
}

const update = () => {
  const model_config_id = formState.model_config_id ? +formState.model_config_id : formState.model_config_id;
  const data = JSON.stringify({
    cate: {
      ...formState,
      question_value: formState.question_value.join('.'),
      model_config_id: model_config_id
    }
  })

  emit('update-node', {
    ...props.node,
    ...formState,
    node_params: data
  })
}

const init = () => {
  getOptions()

  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let cate = JSON.parse(dataRaw).cate || {}

    cate = JSON.parse(JSON.stringify(cate))
    
    for (let key in cate) {
      if(key === 'question_value'){
      formState.question_value = formatQuestionValue(cate['question_value'])
    }else if (key == 'categorys') {
        if (cate.categorys && cate.categorys.length > 0) {
          let items = cate.categorys.map((item) => {
            return {
              ...item,
            }
          })
          formState[key] = items
        } else {
          formState[key] = [
            {
              category: '',
              next_node_key: '',
              key: getUuid(16)
            }
          ]
        }
      } else {
        formState[key] = cate[key]
      }
    }

    formState.role = +cate.role || 1

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

const handleAddcategory = () => {
  formState.categorys.push({
    category: '',
    next_node_key: '',
    key: getUuid(16)
  })
}

const onDelcategory = (index) => {
  formState.categorys.splice(index, 1)
}

const checkedHeader = () => {
  // if (value == null) {
  //   return Promise.reject('请输入延迟发送时间')
  // }
  // if (!Number.isInteger(value / 0.5)) {
  //   return Promise.reject('必须为0.5秒的倍数')
  // }
  return Promise.resolve()
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
