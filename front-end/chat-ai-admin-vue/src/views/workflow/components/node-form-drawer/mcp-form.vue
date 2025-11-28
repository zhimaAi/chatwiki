<style lang="less" scoped>
.node-icon{
  display: block;
  width: 20px;
  height: 20px;
  border-radius: 6px;
}
.node-options {
  background: #f2f4f7;
  border-radius: 6px;
  padding: 12px;
  margin-top: 16px;

  &:first-child {
    margin-top: 0;
  }

  .options-title {
    color: var(--wf-color-text-1);
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-weight: 600;
    height: 22px;;
    line-height: 22px;
    font-size: 14px;

    .title-icon {
      width: 16px;
      height: 16px;
      vertical-align: -3px;
      margin-right: 8px;;
    }

    .acton-box {
      font-weight: 400;
    }
  }

  .options-item {
    display: flex;
    flex-direction: column;
    margin-top: 12px;
    line-height: 22px;
    gap: 4px;

    .options-item-tit {
      display: flex;
      align-items: center;
    }

    .option-label {
      color: var(--wf-color-text-1);
      font-size: 14px;
      margin-right: 8px;
    }

    .desc {
      color: var(--wf-color-text-2);
    }


    &.is-required .option-label::before {
      content: '*';
      color: #FB363F;
      display: inline-block;
      margin-right: 2px;
    }

    .option-type {
      height: 22px;
      line-height: 18px;
      padding: 0 8px;
      border-radius: 6px;
      border: 1px solid rgba(0, 0, 0, 0.15);
      background-color: #fff;
      color: var(--wf-color-text-3);
      font-size: 12px;
    }

    .item-actions-box {
      display: flex;
      align-items: center;

      .action-btn {
        margin-left: 12px;
        font-size: 16px;
        color: #595959;
        cursor: pointer;
      }
    }
  }
}
</style>

<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader :title="node.node_name" :iconName="node.node_icon_name">
         <template #node-icon>
          <img class="node-icon" :src="state.mcpInfo.avatar" alt="">
        </template>
        <template #desc>
          <span>{{ state.toolInfo.description }}</span>
        </template>
      </NodeFormHeader>
    </template>

    <div class="mcp-form">
      <div class="node-options">
        <div class="options-title">
          <div><img src="@/assets/img/workflow/input.svg" class="title-icon" />输入</div>
        </div>
        <div
          v-for="(item, key) in formState.params"
          :key="key"
          :class="['options-item', { 'is-required': item.required }]"
        >
          <div class="options-item-tit">
            <div class="option-label">{{ key }}</div>
            <div class="option-type">{{ item.type }}</div>
          </div>
          <div>
            <!--          <a-input v-model:value="item.value" @change="update" placeholder="键入/插入变量"/>-->
            <AtInput
              type="textarea"
              inputStyle="height: 64px;"
              :options="variableOptions"
              :defaultSelectedList="item.tags"
              :defaultValue="item.value"
              ref="atInputRef"
              @open="getValueVariableList"
              @change="(text, selectedList) => changeValue(item, text, selectedList)"
              placeholder="请输入内容，键入“/”可以插入变量"
            >
              <template #option="{ label, payload }">
                <div class="field-list-item">
                  <div class="field-label">{{ label }}</div>
                  <div class="field-type">{{ payload.typ }}</div>
                </div>
              </template>
            </AtInput>
          </div>
          <div class="desc">{{ item.description }}</div>
        </div>
      </div>
      <div class="node-options">
        <div class="options-title">
          <div><img src="@/assets/img/workflow/output.svg" class="title-icon" />输出</div>
        </div>
        <div class="options-item">
          <div class="options-item-tit">
            <div class="option-label">输出</div>
          </div>
          <div class="options-item-tit">
            <div class="option-label">text</div>
            <div class="option-type">string</div>
          </div>
          <div class="desc">工具生成的内容</div>
        </div>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { jsonDecode } from '@/utils/index'
import { getTMcpProviderInfo } from '@/api/robot/thirdMcp.js'
import { ref, reactive, inject, onMounted } from 'vue'
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
const state = reactive({
  mcpInfo: {},
  toolInfo: {}
})

const formState = reactive({
  params: {},
})

const variableOptions = ref([])

let nodeParams = {}

function init() {
  getValueVariableList();

  nodeParams = JSON.parse(props.node.node_params)

  loadProvider()
}

function getValueVariableList() {
  variableOptions.value = getNode().getAllParentVariable()
}

function changeValue(item, text, selectedList){
  item.value = text;
  item.tags = selectedList

  update()
}

function update() {
  let data = {
    arguments: {},
    tag_map: {}
  }

  for (let key in formState.params) {
    let value = formState.params[key].value
    let field = state.toolInfo.inputSchema.properties[key] || {}

    if(field.type == 'string'){
      value = String(value)
    }else if(['number', 'integer'].includes(field.type)){
      value = Number(value)
    }

    data.arguments[key] = value
    data.tag_map[key] = formState.params[key].tags
  }

  nodeParams.mcp.arguments = data.arguments
  nodeParams.mcp.tag_map = data.tag_map || []

  setData({
    ...props.node,
    node_params: JSON.stringify(nodeParams)
  })
}

const loadProvider = () => {
  getTMcpProviderInfo({
    provider_id: nodeParams.mcp.provider_id
  }).then((res) => {
    state.mcpInfo = res?.data || {}

    let tools = jsonDecode(state.mcpInfo?.tools, [])

    state.toolInfo = tools.find((item) => item.name == nodeParams.mcp.tool_name)

    let inputSchema = state.toolInfo.inputSchema || {}
    let params = inputSchema?.properties || {}
    let requireds = inputSchema?.required || []

    params = JSON.parse(JSON.stringify(params))

    for (let key in params) {
      params[key].value = String(nodeParams.mcp.arguments[key] || '')
      params[key].tags = nodeParams.mcp.tag_map ? nodeParams.mcp.tag_map[key] : []
      params[key].required = requireds.includes(key)
    }

    formState.params = params
  })
}

onMounted(() => {
  init()
})
</script>
