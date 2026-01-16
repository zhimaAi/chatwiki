<style lang="less" scoped>
@import "./components/node-options";

.node-icon {
  display: block;
  width: 20px;
  height: 20px;
  border-radius: 6px;
}

.workflow-form {
  :deep(.mention-input-warpper) {
    height: 32px;
  }
}
</style>

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
    <div class="workflow-form">
      <div class="node-options">
        <div class="options-title">
          <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>输入</div>
        </div>
        <template v-for="item in options" :key="item.key">
          <div :class="['options-item', {'is-required': item.required}]">
            <div class="options-item-tit">
              <div class="option-label">{{ item.name || item.key }}</div>
              <div class="option-type">{{ item.typ }}</div>
            </div>
            <div>
              <AtInput
                type="textarea"
                inputStyle="height: 64px;"
                :options="variableOptions"
                :defaultSelectedList="item.tags"
                :defaultValue="item.variable"
                ref="atInputRef"
                placeholder="请输入内容，键入“/”可以插入变量"
                @open="getValueVariableList"
                @change="(text, selectedList) => changeValue(item, text, selectedList)"
              >
                <template #option="{ label, payload }">
                  <div class="field-list-item">
                    <div class="field-label">{{ label }}</div>
                    <div class="field-type">{{ payload.typ }}</div>
                  </div>
                </template>
              </AtInput>
            </div>
            <div class="desc">{{ item.desc }}</div>
          </div>
        </template>
      </div>
      <div class="node-options">
        <div class="options-title">
          <div><img src="@/assets/img/workflow/output.svg" class="title-icon"/>输出</div>
        </div>
        <div class="options-item">
          <OutputFields :tree-data="outputData"/>
        </div>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import {ref, inject, onMounted} from 'vue'
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import AtInput from '../at-input/at-input.vue'
import {getRobotStartNode} from "@/api/robot/index.js";
import {jsonDecode} from "@/utils/index.js";
import OutputFields from "@/views/workflow/components/feishu-table/output-fields.vue";

const getNode = inject('getNode')
const setData = inject('setData')

const props = defineProps({
  node: {
    type: Object,
    default: () => ({})
  }
})

const options = ref([])
const variableOptions = ref([])
const robotInfo = ref({})
const outputData = ref([])

let nodeParams = {}

async function init() {
  getValueVariableList();
  nodeParams = JSON.parse(props.node.node_params)
  outputData.value = nodeParams?.workflow?.output || []
  robotInfo.value = nodeParams?.workflow?.robot_info || {}
  loadRobotNode()
}

function loadRobotNode() {
  getRobotStartNode({
    robot_id: nodeParams?.workflow?.robot_id
  }).then(res => {
    let _params = jsonDecode(res?.data?.node_params)
    let _data = nodeParams?.workflow?.params || []
    options.value = _params?.start?.diy_global || []
    for (let item of options.value) {
      let _it = _data.find(i => item.key == i.key)
      item.variable = _it?.variable || ''
      item.tags = Array.isArray(_it?.tags) ? _it.tags : []
    }
  })
}

function getValueVariableList() {
  variableOptions.value = getNode().getAllParentVariable()
}

function changeValue(item, text, selectedList) {
  item.variable = text;
  item.tags = selectedList

  update()
}

function update() {
  let _state = JSON.parse(JSON.stringify(options.value))
  for (let item of _state) {
    item.variable = String(item.variable)
  }
  nodeParams.workflow.params = _state
  setData({
    ...props.node,
    node_params: JSON.stringify(nodeParams)
  })
}

onMounted(() => {
  init()
})
</script>
