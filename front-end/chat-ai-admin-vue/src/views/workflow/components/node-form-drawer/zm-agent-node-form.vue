

<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader :title="node.node_name" :desc="robotInfo?.robot_intro">
        <template #node-icon>
          <img v-if="robotInfo.robot_avatar" class="node-icon" :src="robotInfo.robot_avatar"/>
          <a-spin v-else size="small"/>
        </template>
      </NodeFormHeader>
    </template>
    <div class="agent-form">
      <div class="node-options">
        <div class="options-title">
          <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>提问内容</div>
        </div>
        <div class="options-item">
          <AtInput
            type="textarea"
            inputStyle="height: 140px;"
            :options="variableOptions"
            :defaultSelectedList="questionTags"
            :defaultValue="questionValue"
            placeholder="请输入消息内容，键入”/”插入变量"
            @open="getValueVariableList"
            @change="changeValue"
          >
            <template #option="{ label, payload }">
              <div class="field-list-item">
                <div class="field-label">{{ label }}</div>
                <div class="field-type">{{ payload.typ }}</div>
              </div>
            </template>
          </AtInput>
        </div>
      </div>
      <div class="node-options">
        <div class="options-title">
          <div><img src="@/assets/img/workflow/output.svg" class="title-icon"/>输出</div>
        </div>
        <div class="agent-output-list">
          <div
            v-for="item in outputData"
            :key="item.key"
            class="agent-output-item"
          >
            <span class="output-name">{{ item.desc || item.key }}</span>
            <span class="output-type">{{ item.typ }}</span>
          </div>
        </div>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { inject, onMounted, ref } from 'vue'
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import AtInput from '../at-input/at-input.vue'

const getNode = inject('getNode')
const setData = inject('setData')

const props = defineProps({
  node: {
    type: Object,
    default: () => ({})
  }
})

const variableOptions = ref([])
const robotInfo = ref({})
const outputData = ref([])
const questionValue = ref('')
const questionTags = ref([])

let nodeParams = {}

function init() {
  getValueVariableList()
  nodeParams = JSON.parse(props.node.node_params || '{}')
  const agent = nodeParams.agent || {}
  robotInfo.value = agent.robot_info || {}
  outputData.value = agent.output || []
  questionValue.value = agent.question_value || ''
  questionTags.value = Array.isArray(agent.question_tags) ? agent.question_tags : []
}

function getValueVariableList() {
  variableOptions.value = getNode().getAllParentVariable()
}

function changeValue(text, selectedList) {
  questionValue.value = text
  questionTags.value = selectedList

  update()
}

function update() {
  if (!nodeParams.agent) {
    nodeParams.agent = {}
  }

  nodeParams.agent.question_value = questionValue.value
  nodeParams.agent.question_tags = questionTags.value
  const node_params = JSON.stringify(nodeParams)

  setData({
    ...props.node,
    node_params,
    dataRaw: node_params
  })
}

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import "./components/node-options";

.node-icon {
  display: block;
  width: 20px;
  height: 20px;
  border-radius: 6px;
}

.agent-form {
  :deep(.mention-input-warpper) {
    height: 140px;
  }

  .agent-output-list {
    padding: 10px 14px 0;
  }

  .agent-output-item {
    display: flex;
    align-items: center;
    min-height: 24px;
    line-height: 22px;
    color: #262626;
    font-size: 14px;
  }

  .output-name {
    &::before {
      content: '*';
      color: #fb363f;
      margin-right: 2px;
    }
  }

  .output-type {
    margin-left: 6px;
    padding: 1px 8px;
    line-height: 20px;
    border-radius: 6px;
    border: 1px solid #d9d9d9;
    background: #fff;
    color: #595959;
    font-size: 12px;
  }
}
</style>
