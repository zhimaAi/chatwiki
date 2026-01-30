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
      gap: 4px;
    }
    .library-item {
      display: flex;
      align-items: center;
      line-height: 16px;
      padding: 3px 4px;
      border-radius: 4px;
      border: 1px solid #d9d9d9;
      color: #595959;
      background: #fff;
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
    :properties="properties"
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
          <div class="field-item-label">输出字段</div>
          <div class="field-item-content">
            <div class="field-value is-required">
              <span class="field-key">知识库引用</span>
              <span class="field-type">string</span>
            </div>
          </div>
        </div>

        <div class="field-item">
          <div class="field-item-label">知识库</div>
          <div class="field-item-content">
            <div class="library-item" v-for="item in selectedLibraryRows" :key="item.id">{{ item.library_name }}</div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { ref, reactive, watch, onMounted, inject, nextTick, onBeforeUnmount, computed } from 'vue'
import NodeCommon from '../base-node.vue'
import UserQuestionText from '../user-question-text.vue'
import { getLibraryList } from '@/api/library/index'
import { useRobotStore } from '@/stores/modules/robot'
const robotStore = useRobotStore()

const rrf_weight = computed(()=>{
  return robotStore.robotInfo.rrf_weight
})

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

// --- State ---
const questionTextRef = ref(null)
const menus = ref([])
const formState = reactive({
  library_ids: [],
  rerank_status: 0,
  rerank_use_model: undefined,
  rerank_model_config_id: void 0,
  top_k: 5,
  similarity: 0.5,
  search_type: 1,
  question_value: [],
  rrf_weight: {},
  recall_neighbor_switch: false,
  recall_neighbor_before_num: 1,
  recall_neighbor_after_num: 1,
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
  let libs = {}
  try {
    libs = JSON.parse(dataRaw).libs || {}
  } catch (e) {
    libs = {}
  }

  getVlaueVariableList()

  libs = JSON.parse(JSON.stringify(libs))

  for (let key in libs) {
    if (key == 'library_ids') {
      formState[key] = libs[key] ? libs[key].split(',') : []
    } else if (key == 'question_value') {
      formState.question_value = formatQuestionValue(libs['question_value'])
    }else if(key == 'rrf_weight') {
      formState.rrf_weight = libs[key] ? JSON.parse(libs[key]) : libs[key]
    } else {
      formState[key] = libs[key]
    }
  }
  
  if (!formState.rrf_weight || Object.keys(formState.rrf_weight).length == 0) {
    //  没有值 则去默认值
    formState.rrf_weight = rrf_weight.value
  }

  nextTick(() => {
    resetSize()
  })
}

const update = () => {
  const data = JSON.stringify({
    libs: {
      ...formState,
      rerank_model_config_id: formState.rerank_model_config_id
        ? +formState.rerank_model_config_id
        : void 0,
      question_value: formState.question_value.join('.'),
      library_ids: formState.library_ids.join(','),
      rrf_weight: JSON.stringify(formState.rrf_weight),
      recall_neighbor_switch: formState.recall_neighbor_switch,
      recall_neighbor_before_num: formState.recall_neighbor_before_num,
      recall_neighbor_after_num: formState.recall_neighbor_before_num,
    }
  })

  setData({
    ...props.node,
    node_params: data
  })
}

const libraryList = ref([])
const selectedLibraryRows = computed(() => {
  return libraryList.value.filter((item) => {
    return formState.library_ids.includes(item.id)
  })
})
// 获取知识库
const getList = async () => {
  const res = await getLibraryList({ type: '' })
  if (res) {
    libraryList.value = res.data || []

    nextTick(() => {
      resetSize()
    })
  }
}

const getVlaueVariableList = () => {
  let list = getNode().getAllParentVariable()

  variableOptions.value = list
}

const onUpatateNodeName = () => {
  getVlaueVariableList()

  update()

  questionTextRef.value.refresh()
}

watch(() => props.properties, reset, { deep: true })

onMounted(() => {
  getList()

  reset()
  const mode = graphModel()

  mode.eventCenter.on('custom:setNodeName', onUpatateNodeName)
})

onBeforeUnmount(() => {
  const mode = graphModel()

  mode.eventCenter.off('custom:setNodeName', onUpatateNodeName)
})
</script>
