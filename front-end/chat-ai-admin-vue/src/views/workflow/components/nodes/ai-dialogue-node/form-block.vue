<template>
  <div class="form-block" @mousedown.stop="">
    <a-form ref="formRef" layout="vertical" :model="formState">
      <div class="gray-block">
        <div class="gray-block-title">输入</div>
        <a-form-item label="LLM模型" name="use_model">
          <div class="flex-block-item">
            <a-select
              v-model:value="formState.use_model"
              placeholder="请选择LLM模型"
              @change="handleChangeModel"
              style="width: 348px"
            >
              <a-select-opt-group v-for="item in modelList" :key="item.id">
                <template #label>
                  <a-flex align="center" :gap="6">
                    <img class="model-icon" :src="item.icon" alt="" />{{ item.name }}
                  </a-flex>
                </template>
                <a-select-option
                  :value="
                    modelDefine.indexOf(item.model_define) > -1 && val.deployment_name
                      ? val.deployment_name
                      : val.name + val.id
                  "
                  :model_config_id="item.id"
                  :current_obj="val"
                  v-for="val in item.children"
                  :key="val.name + val.id"
                >
                  <span v-if="modelDefine.indexOf(item.model_define) > -1 && val.deployment_name">{{
                    val.deployment_name
                  }}</span>
                  <span v-else>{{ val.name }}</span>
                </a-select-option>
              </a-select-opt-group>
            </a-select>
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
        <a-form-item name="context_pair">
          <template #label>
            <span>上下文数量&nbsp;</span>
            <a-tooltip>
              <template #title
                >提示词中携带的历史聊天记录轮次。设置为0则不携带聊天记录。最多设置10轮。注意，携带的历史聊天记录越多，消耗的token相应也就越多。</template
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
                  :max="10"
                />
              </a-form-item-rest>
            </div>
            <div class="number-input-box">
              <a-input-number v-model:value="formState.context_pair" :min="0" :max="10" />
            </div>
          </div>
        </a-form-item>
        <a-form-item label="提示词" name="prompt">
          <a-mentions
            @wheel.stop=""
            v-model:value="formState.prompt"
            prefix="/"
            placeholder="输入 / 插入变量"
            :options="variableOptions"
            @select="onTextChange"
            @blur="() => (isFocus = false)"
            @focus="() => (isFocus = true)"
            rows="4"
          >
            <template #option="{ value, label, payload }">
              <div class="field-list-item">
                <div class="field-label">{{ label }}</div>
                <div class="field-type">{{ payload.typ }}</div>
              </div>
            </template>
          </a-mentions>
          <div class="form-tip">输入 / 插入变量</div>
        </a-form-item>
        <div class="diy-form-item">
          <div class="form-label">用户问题</div>
          <div class="form-content">流程开始>用户问题</div>
        </div>
        <div class="diy-form-item mt12">
          <div class="form-label">知识库引用</div>
          <div class="form-content">流程开始>用户问题</div>
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
  </div>
</template>

<script setup>
import { ref, reactive, watch, h, inject, computed } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  CloseCircleFilled,
  CloseCircleOutlined,
  QuestionCircleOutlined,
  UpOutlined,
  DownOutlined,
  LoadingOutlined,
  PlusOutlined,
  EditOutlined,
  SyncOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons-vue'

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

const modelList = computed(() => {
  return robotStore.modelList
})

const isFocus = ref(false)
const variableOptions = ref([])

const showMoreBtn = ref(false)

const hanldeShowMore = () => {
  showMoreBtn.value = !showMoreBtn.value
  emit('setData', {
    height: showMoreBtn.value ? 810 : 670
  })
}

const emit = defineEmits(['setData'])
const formRef = ref()

const formState = reactive({
  model_config_id: void 0,
  use_model: void 0,
  temperature: 0,
  max_token: 0,
  context_pair: 0,
  prompt: ''
})
let lock = false
watch(
  () => props.properties,
  (val) => {
    try {
      if (!isFocus.value) {
        getVariableOptions()
      }
      if (lock) {
        return
      }
      let llm = JSON.parse(val.node_params).llm || {}
      llm = JSON.parse(JSON.stringify(llm))
      let { model_config_id, use_model, context_pair, temperature, max_token, prompt } = llm
      formState.model_config_id = model_config_id
      formState.use_model = use_model
      formState.context_pair = context_pair
      formState.temperature = temperature
      formState.max_token = max_token
      formState.prompt = prompt
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
  console.log('===')
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

function getVariableOptions() {
  let node = getNode()
  let preNodes = graphModel().getNodeIncomingNode(node.id)
  let outOptions = []

  if (preNodes && preNodes.length) {
    preNodes.forEach((item) => {
      let node_type = item.properties.node_type
      let output = item.properties.output
      if (node_type == 4 && output && output.length) {
        outOptions = transformArray(output)
        outOptions = outOptions.filter((it) => !it.hasSub)
        outOptions = outOptions.map((it) => {
          return {
            ...it,
            value: `【${it.value}】`
          }
        })
      }
    })
  }
  let lists = [
    {
      label: '用户消息',
      value: '【global.question】',
      payload: { typ: 'string' }
    },
    {
      label: 'open_id',
      value: '【global.openid】',
      payload: { typ: 'string' }
    },
    ...outOptions
  ]
  variableOptions.value = lists
}

const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']

const handleChangeModel = (val, option) => {
  const self = option.current_obj
  formState.use_model =
    modelDefine.indexOf(self.model_define) > -1 && self.deployment_name
      ? self.deployment_name
      : self.name
  formState.model_config_id = self.id || option.model_config_id
}

defineExpose({})
</script>

<style lang="less" scoped>
@import '../form-block.less';

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
