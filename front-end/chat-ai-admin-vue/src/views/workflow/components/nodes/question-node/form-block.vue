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
              style="width: 348px"
            />
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
        <a-form-item label="提示词" name="prompt" v-if="false">
          <a-mentions
            @wheel.stop=""
            v-model:value="formState.prompt"
            @blur="() => (isFocus = false)"
            @focus="() => (isFocus = true)"
            prefix="/"
            placeholder="输入 / 插入变量"
            :options="variableOptions"
            @select="onTextChange"
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
      </div>
      <div class="gray-block mt16">
        <div class="gray-block-title">问题分类设置</div>
        <div class="array-form-box">
          <div class="form-item-list" v-for="(item, index) in formState.categorys" :key="item.key">
            <a-form-item
              :label="null"
              :name="['categorys', index, 'category']"
              :rules="{ required: true, validator: (rule, value) => checkedHeader(rule, value) }"
            >
              <div class="flex-block-item">
                <a-input class="flex1" v-model:value="item.category" placeholder="请输入"></a-input>
                <div class="btn-hover-wrap" @click="onDelcategory(index)">
                  <CloseCircleOutlined />
                </div>
              </div>
            </a-form-item>
            </div>

            <div class="form-item-list">
              <a-form-item
              :label="null"
            >
              <div class="flex-block-item">
                <a-input class="flex1" value="默认分类" readonly></a-input>
              </div>
            </a-form-item>
          </div>
          <a-button @click="handleAddcategory" :icon="h(PlusOutlined)" block type="dashed"
            >添加问题分类</a-button
          >
        </div>
      </div>
    </a-form>
  </div>
</template>

<script setup>
import { getSystemVariable } from '../../util'
import { ref, reactive, watch, h, inject, computed } from 'vue'
import {
  CloseCircleOutlined,
  QuestionCircleOutlined,
  UpOutlined,
  DownOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue'
import ModelSelect from '@/components/model-select/model-select.vue'
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

const emit = defineEmits(['setData'])

const showMoreBtn = ref(false)

const hanldeShowMore = () => {
  showMoreBtn.value = !showMoreBtn.value
  emit('setData', {
    ...formState,
    showMoreBtn: showMoreBtn.value,
    height: getNodeHeight()
  })
}

const formRef = ref()

const formState = reactive({
  model_config_id: void 0,
  use_model: void 0,
  temperature: 0,
  max_token: 0,
  context_pair: 0,
  prompt: '',
  categorys: [
    {
      category: '',
      next_node_key: '',
      key: Math.random() * 10000
    }
  ]
})
let lock = false
const isFocus = ref(false)
const variableOptions = ref([])
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
      let cate = JSON.parse(val.node_params).cate || {}
      cate = JSON.parse(JSON.stringify(cate))
      for (let key in cate) {
        if (key == 'categorys') {
          if (cate.categorys && cate.categorys.length > 0) {
            let items = cate.categorys.map((item) => {
              return {
                ...item,
                key: Math.random() * 10000
              }
            })
            formState[key] = items
          } else {
            formState[key] = [
              {
                category: '',
                next_node_key: '',
                key: Math.random() * 10000
              }
            ]
          }
        } else {
          formState[key] = cate[key]
        }
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
            cate: {
              ...formState,
              model_config_id: formState.model_config_id
                ? +formState.model_config_id
                : formState.model_config_id
            }
          }),
          showMoreBtn: showMoreBtn.value,
          height: getNodeHeight()
        })
      }, 100)
    } catch (error) {}
  },
  { immediate: true, deep: true }
)

watch(
  () => formState,
  (val) => {
    emit('setData', {
      ...formState,
      node_params: JSON.stringify({
        cate: {
          ...formState,
          model_config_id: formState.model_config_id
            ? +formState.model_config_id
            : formState.model_config_id
        }
      }),
      showMoreBtn: showMoreBtn.value,
      height: getNodeHeight()
    })
  },
  { deep: true }
)

function getNodeHeight() {
  let baseHeight = 72 + 260 + 132 + (formState.categorys.length + 1) * 38
  return showMoreBtn.value ? baseHeight + 144 : baseHeight
}

const handleAddcategory = () => {
  formState.categorys.push({
    category: '',
    next_node_key: '',
    key: Math.random() * 10000
  })
}

const onDelcategory = (index) => {
  formState.categorys.splice(index, 1)
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
    ...getSystemVariable(),
    ...outOptions
  ]
  variableOptions.value = lists
}

const checkedHeader = (rule, value) => {
  // if (value == null) {
  //   return Promise.reject('请输入延迟发送时间')
  // }
  // if (!Number.isInteger(value / 0.5)) {
  //   return Promise.reject('必须为0.5秒的倍数')
  // }
  return Promise.resolve()
}


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
