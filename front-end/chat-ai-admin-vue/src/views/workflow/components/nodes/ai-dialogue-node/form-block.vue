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
              <div class="btn-hover-wrap" @click="onShowAddPromptModal">
                <DownloadOutlined />从提示词库导入
              </div>
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
          <div class="form-content">流程开始>用户问题</div>
        </div>
        <div class="diy-form-item mt12">
          <div class="form-label">知识库引用</div>
          <div class="form-content">知识库搜索>知识库引用</div>
        </div>
      </div>
      <div class="gray-block mt16">
        <div class="gray-block-title">输出</div>
        <div class="options-item">
          <div class="option-label">AI回复内容</div>
          <div class="option-type">string</div>
        </div>
      </div>
    </a-form>
    <ImportPrompt @ok="handleSavePrompt" ref="importPromptRef" />
  </div>
</template>

<script setup>
import { ref, reactive, watch, inject, computed, nextTick, onMounted, onBeforeUnmount } from 'vue'
import {
  QuestionCircleOutlined,
  UpOutlined,
  DownOutlined,
  DownloadOutlined
} from '@ant-design/icons-vue'
import AtInput from '../at-input/at-input.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import ImportPrompt from '@/components/import-prompt/index.vue'

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

const atInputRef = ref(null)

const modelList = computed(() => {
  return robotStore.modelList
})

const isFocus = ref(false)
const variableOptions = ref([])

const showMoreBtn = ref(false)

const hanldeShowMore = () => {
  showMoreBtn.value = !showMoreBtn.value
  emit('setData', {
    height: showMoreBtn.value ? 890 : 674
  })
}

const changeValue = (text, selectedList) => {
  formState.prompt = text
  formState.prompt_tags = selectedList
}

const getVlaueVariableList = () => {
  let list = getNode().getAllParentVariable()

  variableOptions.value = list
}

const emit = defineEmits(['setData'])
const formRef = ref()

const formState = reactive({
  model_config_id: void 0,
  use_model: void 0,
  temperature: 0,
  max_token: 0,
  context_pair: 0,
  prompt: '',
  prompt_tags: [],
  enable_thinking: false
})
let lock = false
watch(
  () => props.properties,
  (val) => {
    try {
      if (lock) {
        return
      }

      let dataRaw = val.dataRaw || val.node_params || '{}'
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
        enable_thinking,
      } = llm

      getVlaueVariableList()

      formState.model_config_id = model_config_id
      formState.use_model = use_model
      formState.context_pair = context_pair || 0
      formState.temperature = temperature
      formState.max_token = max_token
      formState.prompt = prompt
      formState.enable_thinking = enable_thinking
      formState.prompt_tags = prompt_tags || []
      if (!formState.model_config_id && modelList.value.length > 0) {
        formState.model_config_id = modelList.value[0].id
        formState.use_model = modelList.value[0].children[0].name
      }
      lock = true
      setTimeout(() => {
        emit('setData', {
          ...formState,
          node_params: JSON.stringify({
            llm: {
              ...formState,
              model_config_id: formState.model_config_id
                ? +formState.model_config_id
                : formState.model_config_id
            }
          })
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
        llm: {
          ...formState,
          model_config_id: formState.model_config_id
            ? +formState.model_config_id
            : formState.model_config_id
        }
      })
    })
  },
  { deep: true }
)

const checkedHeader = (rule, value) => {
  // if (value == null) {
  //   return Promise.reject('请输入延迟发送时间')
  // }
  // if (!Number.isInteger(value / 0.5)) {
  //   return Promise.reject('必须为0.5秒的倍数')
  // }
  return Promise.resolve()
}

const onTextChange = () => {
  let regex = / +【/g
  formState.prompt = formState.prompt.replace(/\//g, '').replace(regex, '【')
}

function transformArray(arr, parentLabel = '') {
  let result = []

  arr.forEach((item) => {
    let newLabel = parentLabel ? `${parentLabel}.${item.key}` : String(item.key)
    let newValue = parentLabel ? `${parentLabel}.${item.key}` : String(item.key)

    result.push({
      label: newLabel,
      value: newValue,
      payload: { typ: item.typ },
      hasSub: item.subs && item.subs.length > 0
    })

    if (item.subs && Array.isArray(item.subs)) {
      result = result.concat(transformArray(item.subs, newLabel))
    }
  })
  return result
}

const onUpatateNodeName = (data) => {
  if (data.node_type !== 'http-node') {
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
