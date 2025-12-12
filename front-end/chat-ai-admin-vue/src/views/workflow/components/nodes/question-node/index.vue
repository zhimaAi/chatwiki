<template>
  <node-common
    :title="props.properties.node_name"
    :menus="menus"
    :icon-name="props.properties.node_icon_name"
    :isSelected="props.isSelected"
    :isHovered="props.isHovered"
    :node-key="props.properties.node_key"
    :node_type="props.properties.node_type"
    style="width: 420px"
  >
    <div class="ai-dialogue-node">
      <div class="field-list">
        <div class="field-item">
          <div class="field-item-label">用户问题</div>
          <div class="field-item-content">
            <div class="field-value">
              <span class="field-key">
                <user-question-text ref="questionTextRef" :value="formState.question_value" />
              </span>
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
          <div class="field-item-label">问题分类</div>
          <div class="field-item-content">
            <div class="category-list">
              <div class="field-value category-value" v-for="item in formState.categorys" :key="item.key">
                <span>{{ item.category }}</span>
              </div>

              <div class="field-value category-value">
                <span>默认分类</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { getUuid } from '@/utils/index'
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
const resetSize = inject('resetSize')

// --- Store ---
const robotStore = useRobotStore()
const { modelList } = storeToRefs(robotStore)

// --- State ---
const menus = ref([])
const formState = reactive({
  model_config_id: void 0,
  use_model: void 0,
  temperature: 0,
  max_token: 0,
  context_pair: 0,
  prompt: '',
  question_value: [],
  enable_thinking: false,
  categorys: []
})

function formatQuestionValue(val) {
  if (val && typeof val === 'string') {
    let lists = val.split('.')
    let str1 = lists[0]
    let str2 = lists.filter((item, index) => index > 0).join('.')
    return [str1, str2]
  } else if (val && typeof val === 'object' && val.length === 2) {
    return val
  }

  return ['global', 'question']
}

const reset = () => {
  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  let cate = {}
  try {
    cate = JSON.parse(dataRaw).cate || {}
  } catch (e) {
    cate = {}
  }

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

  if (!formState.model_config_id && modelList.value.length > 0) {
    formState.model_config_id = modelList.value[0].id
    formState.use_model = modelList.value[0].children[0].name
  }

  nextTick(() => {
    resetSize()
  })
}

const update = () => {
  const model_config_id = formState.model_config_id  ? +formState.model_config_id : formState.model_config_id;
  const data = JSON.stringify({
    cate: {
      ...formState,
      question_value: formState.question_value.join('.'),
      model_config_id: model_config_id
    }
  })
  
  setData({
    ...props.node,
    ...formState,
    node_params: data
  })
}

watch(() => props.properties, (newVal, oldVal) => {
  const newDataRaw = newVal.dataRaw || newVal.node_params || '{}'
  const oldDataRaw = oldVal.dataRaw || oldVal.node_params || '{}'
  
  if(newDataRaw != oldDataRaw) { 
    reset()
  }
}, { deep: true })

onMounted(() => {
  reset()
  resetSize()
  update()
})

onBeforeUnmount(() => {
})
</script>


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
      overflow: hidden;
    }
    .category-list {
      width: 100%;
      overflow: hidden;
    }
    .field-value {
      width: 100%;
      line-height: 16px;
      padding: 3px 4px;
      height: 24px;
      margin-bottom: 8px;
      border-radius: 4px;
      font-size: 12px;
      border: 1px solid #d9d9d9;
      color: #595959;
      background: #fff;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      &:last-child {
        margin-bottom: 0;
      }
      .category-value{
        width: 100%;
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