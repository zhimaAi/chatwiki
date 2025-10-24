<template>
  <div class="form-block" @mousedown.stop="">
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
              <a-input-number v-model:value="formState.temperature" :min="0" :max="2" :step="0.1" />
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
              <template #title
                >开启时，调用大模型时会指定走深度思考模式</template
              >
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
              <div>提示词</div>
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
            placeholder="请输入消息内容，键入“/”可以插入变量"
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
        <div class="gray-block-title" @click="test">输出 (输出字段提取)</div>
        <div class="output-box">
          <div class="output-block">
            <div class="output-item">参数Key</div>
            <div class="output-item">类型</div>
          </div>
          <div class="array-form-box" @mousedown.stop="">
            <div class="form-item-list" v-for="(item, index) in formState.output" :key="index">
              <a-form-item :label="null" :name="['output', index, 'key']">
                <div class="flex-block-item" style="gap: 12px">
                  <a-input
                    style="width: 214px"
                    v-model:value="item.key"
                    placeholder="请输入"
                  ></a-input>
                  <a-form-item-rest>
                    <a-select
                      @change="onTypeChange(item)"
                      v-model:value="item.typ"
                      placeholder="请选择"
                      style="width: 214px"
                    >
                      <a-select-option v-for="op in typOptions" :value="op.value">{{
                        op.value
                      }}</a-select-option>
                    </a-select>
                  </a-form-item-rest>

                  <div class="btn-hover-wrap" v-if="item.typ == 'object'" @click="onAddSubs(index)">
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
              >添加参数</a-button
            >
          </div>
        </div>
      </div>
    </a-form>
    <AddParamsModal
      @add="onOutputAdd"
      @edit="onOutputEdit"
      @addSub="onOutputAddSub"
      ref="addParamsModalRef"
    />
  </div>
</template>

<script setup>
import {
  ref,
  reactive,
  watch,
  h,
  inject,
  computed,
  nextTick,
  onMounted,
  onBeforeUnmount
} from 'vue'
import {
  PlusOutlined,
  PlusCircleOutlined,
  DownOutlined,
  QuestionCircleOutlined,
  UpOutlined,
  CloseCircleOutlined,
  EditOutlined,
  DownloadOutlined
} from '@ant-design/icons-vue'
import SubKey from './subs-key.vue'
import AtInput from '../at-input/at-input.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import AddParamsModal from './add-params-modal.vue'
import { haveOutKeyNode } from '@/views/workflow/components/util.js'
import { useRobotStore } from '@/stores/modules/robot'
const robotStore = useRobotStore()

const graphModel = inject('getGraph')
const getNode = inject('getNode')

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['setData'])

const atInputRef = ref(null)

const modelList = computed(() => {
  return robotStore.modelList
})

const changeValue = (text, selectedList) => {
  formState.prompt = text
  formState.prompt_tags = selectedList
}
const getVlaueVariableList = () => {
  let list = getNode().getAllParentVariable()
  list.forEach((item) => {
    item.tags = item.tags || []
  })

  variableOptions.value = list
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
  output: [
    // {
    //   key: '',
    //   typ: 'string',
    //   required: false,
    //   default: '',
    //   enum: '',
    //   subs: []
    // }
  ]
})
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

function formatQuestionValue(val) {
  if (val) {
    let lists = val.split('.')
    let str1 = lists[0]
    let str2 = lists.filter((item, index) => index > 0).join('.')
    return [str1, str2]
  }
  return ['global', 'question']
}

let lock = false
const variableOptions = ref([])
const variableOptionsSelect = ref([])

const showMoreBtn = ref(false)

const hanldeShowMore = () => {
  showMoreBtn.value = !showMoreBtn.value
  emit('setData', {
    height: getNodeHeight()
  })
}

watch(
  () => props.properties,
  (val) => {
    try {
      if (lock) {
        return
      }
      getVlaueVariableList()
      getOptions()
      let dataRaw = val.dataRaw || val.node_params || '{}'
      let params_extractor = JSON.parse(dataRaw).params_extractor || {}

      params_extractor = JSON.parse(JSON.stringify(params_extractor))

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
      lock = true
      setTimeout(() => {
        emit('setData', {
          ...formState,
          node_params: JSON.stringify({
            params_extractor: {
              ...formState,
              question_value: formState.question_value.join('.'),
              model_config_id: formState.model_config_id
                ? +formState.model_config_id
                : formState.model_config_id
            }
          }),
          height: getNodeHeight()
        })
      }, 100)
    } catch (error) {
      console.log(error)
    }
  },
  { immediate: true, deep: true }
)

watch(
  () => formState,
  (val) => {
    emit('setData', {
      ...formState,
      node_params: JSON.stringify({
        params_extractor: {
          ...formState,
          question_value: formState.question_value.join('.'),
          model_config_id: formState.model_config_id
                ? +formState.model_config_id
                : formState.model_config_id
        }
      }),
      height: getNodeHeight()
    })
  },
  { deep: true }
)
const test = () => {
  // console.log(formState.output, '==')
}

function getNodeHeight() {
  let height = showMoreBtn.value ? 670 : 520

  // 下面输出字段高度

  let outLens = calculateTotalLength(formState.output)

  return height + 180 + (outLens - 1) * 36
}

function calculateTotalLength(array) {
  let totalLength = array.length // 先加上主数组的长度

  // 遍历数组中的每个对象
  for (const item of array) {
    if (item.subs && Array.isArray(item.subs)) {
      // 如果对象有 subs 属性且 subs 是数组，则递归计算 subs 的长度
      totalLength += calculateTotalLength(item.subs)
    }
  }

  return totalLength
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
  // formState.output[index].subs.push({
  //   key: '',
  //   value: '',
  //   subs: [],
  //   cu_key: Math.random() * 10000
  // })
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


const onUpatateNodeName = (data) => {
  if(!haveOutKeyNode.includes(data.node_type)){
    return
  }

  getVlaueVariableList()
  nextTick(() => {
    if (formState.prompt_tags && formState.prompt_tags.length > 0) {
      formState.prompt_tags.forEach((tag) => {
        if (tag.node_id == data.node_id) {
          let arr = tag.label.split('/')
          arr[0] = data.node_name
          tag.label = arr.join('/')
          tag.node_name = data.node_name
        }
      })

      atInputRef.value.refresh()
    }
  })
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
const onDropdownVisibleChange = (visible) => {
  if (!visible) {
    getOptions()
  }
}

const choosable_thinking = ref({})
const onVectorModelLoaded = (list, choosable_thinking_map) => {
  choosable_thinking.value = choosable_thinking_map
}

const show_enable_thinking = computed(()=>{
  if(!formState.model_config_id){
    return false
  }
  let key = formState.model_config_id + '#' + formState.use_model
  return choosable_thinking.value[key]
})

onMounted(() => {
  const mode = graphModel()

  mode.eventCenter.on('custom:setNodeName', onUpatateNodeName)
})

onBeforeUnmount(() => {
  const mode = graphModel()

  mode.eventCenter.off('custom:setNodeName', onUpatateNodeName)
})
defineExpose({})
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
