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
      gap: 8px;
    }
    .menu-list-box {
      display: flex;
      flex-direction: column;
      gap: 8px;
    }
    .field-list-item {
      display: flex;
      gap: 4px;
      align-items: center;
      line-height: 16px;
      padding: 3px 4px;
      border-radius: 4px;
      border: 1px solid #d9d9d9;
      color: #595959;
      background: #fff;
      .right-arrow {
        width: 24px;
        height: 100%;
        border-radius: 4px;
        background: #e4e6eb;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 13px;
      }
      .field-text {
        max-width: 250px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        font-size: 12px;
      }
    }
    .menu-list-item {
      width: 100%;
      gap: 4px;
      align-items: center;
      line-height: 16px;
      padding: 3px 4px;
      border-radius: 4px;
      border: 1px solid #d9d9d9;
      color: #595959;
      background: #fff;
      .field-text {
        width: 100%;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        font-size: 12px;
      }
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
      .field-key {
        max-width: 200px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
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
    style="width: 420px"
  >
    <div class="ai-dialogue-node">
      <div class="field-list">
        <div class="field-item">
          <div class="field-item-label" @click="handleDelEdg">提问内容</div>
          <div class="field-item-content">
            <div class="field-list-item">
              <div class="field-text">
                <at-text
                  :options="valueOptions"
                  :defaultValue="formState.answer_text"
                  ref="atInputRef"
                  v-if="formState.answer_text.length > 0"
                />
                <span v-else>--</span>
              </div>
            </div>
          </div>
        </div>

        <div class="field-item">
          <div class="field-item-label">回答方式</div>
          <div class="field-item-content">
            <div class="field-value">
              <span class="field-key">
                {{ formState.answer_type == 'text' ? '直接回答' : '智能菜单回答' }}</span
              >
            </div>
          </div>
        </div>
        <div class="field-item menu-list-box" v-if="formState.answer_type == 'menu'">
          <div class="menu-list-item" v-for="(item, index) in formState.menu_content" :key="index">
            <div class="field-text">
              <at-text
                :options="valueOptions"
                :defaultValue="item.content"
                ref="atInputRef"
                v-if="item.content.length > 0"
              />
              <span v-else>--</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import { ref, reactive, watch, onMounted, inject, nextTick, onBeforeUnmount } from 'vue'
import NodeCommon from '../base-node.vue'
import AtText from '../../at-input/at-text.vue'
import { haveOutKeyNode } from '@/views/workflow/components/util.js'

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

// --- State ---
const menus = ref([])
const valueOptions = ref([])
const formState = reactive({
  answer_text: '',
  answer_type: 'text',
  reply_content_list: [],
  menu_content: [
    {
      menu_type: '1',
      serial_no: '1',
      content: ''
    },
    {
      menu_type: '1',
      serial_no: '2',
      content: ''
    }
  ],

  outputs: [
    {
      key: 'question',
      typ: 'string',
      subs: []
    },
    {
      key: 'question_multiple',
      typ: 'array<object>',
      subs: []
    }
  ]
})

const handleDelEdg = () => {
  console.log(graphModel(), getNode())
  let nodeSortKey = props.properties.nodeSortKey
  let edges = graphModel().edges
  if (formState.answer_type == 'text') {
    formState.menu_content.forEach((item, index) => {
      let sourceAnchorId = nodeSortKey + '-anchor_' + index
      let delEdges = edges.find((item) => item.sourceAnchorId == sourceAnchorId)
      if (delEdges) {
        graphModel().eventCenter.emit('custom:edge:delete', delEdges)
      }
    })
  } else {
    let sourceAnchorId = nodeSortKey + '-anchor_right'

    let delEdges = edges.find((item) => item.sourceAnchorId == sourceAnchorId)
    if (delEdges) {
      graphModel().eventCenter.emit('custom:edge:delete', delEdges)
    }
  }
}

watch(
  () => formState.answer_type,
  (val) => {
    handleDelEdg()
  }
)

const reset = () => {
  getValueOptions()
  const dataRaw = props.properties.dataRaw || props.properties.node_params || '{}'
  let question = JSON.parse(dataRaw).question || {}
  question = JSON.parse(JSON.stringify(question))
  formState.answer_text = question.answer_text || ''
  formState.answer_type = question.answer_type || 'text'
  formState.reply_content_list = question.reply_content_list || []
  let reply_content_list = question.reply_content_list || []
  if (reply_content_list.length > 0) {
    let menu_content = reply_content_list[0].smart_menu?.menu_content || []
    if (menu_content.length > 0) {
      formState.menu_content = menu_content
    }
  }
  try {
  } catch (e) {}

  nextTick(() => {
    resetSize()
  })
}


const update = () => {
  const data = JSON.stringify({
    question: {
      ...formState,
    }
  })

  setData({
    ...props.node,
    node_params: data
  })
}
const onUpatateNodeName = (data) => {
  update()
}

const getValueOptions = () => {
  let options = getNode().getAllParentVariable()

  valueOptions.value = options || []
}

// --- Watchers and Lifecycle Hooks ---
watch(() => props.properties, reset, { deep: true })

onMounted(() => {
  reset()
  resetSize()
  const mode = graphModel()

  mode.eventCenter.on('custom:setNodeName', onUpatateNodeName)
})

onBeforeUnmount(() => {
  const mode = graphModel()

  mode.eventCenter.off('custom:setNodeName', onUpatateNodeName)
})
</script>
