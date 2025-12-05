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
      <NodeFormHeader :title="node.node_name" :desc="showDesc">
        <template #node-icon>
          <img v-if="info.icon" class="node-icon" :src="info.icon"/>
          <a-spin v-else size="small"/>
        </template>
      </NodeFormHeader>
    </template>

    <FeishuBittableBox
      v-if="nodeParams?.plugin?.name === 'feishu_bitable'"
      :node="node"
      :action="actionInfo"
      :actionName="nodeParams.plugin.params.business"
      :params="formState"
      :variableOptions="variableOptions"
      @updateVar="getValueVariableList"
    />
    <div v-else class="node-options">
        <div class="options-title">
          <div><img src="@/assets/img/workflow/input.svg" class="title-icon" />输入</div>
        </div>
        <div v-for="(item, key) in formState"
             :key="key"
             :class="['options-item', {'is-required': item.required}]">
          <div class="options-item-tit">
            <div class="option-label">{{ key }}</div>
            <div class="option-type">{{ item.type }}</div>
          </div>
          <div>
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
  </NodeFormLayout>
</template>

<script setup>
import { ref, reactive, inject, onMounted, computed} from 'vue'
import NodeFormLayout from './node-form-layout.vue'
import NodeFormHeader from './node-form-header.vue'
import AtInput from '../at-input/at-input.vue'
import {getPluginInfo, runPlugin} from "@/api/plugins/index.js";
import defaultIcon from '@/assets/svg/plugin-node-type.svg'
import FeishuBittableBox from "./components/feishu-bittable/feishu-bittable-box.vue";

const getNode = inject('getNode')
const setData = inject('setData')

const props = defineProps({
  node: {
    type: Object,
    default: () => ({})
  }
})

const formState = reactive({})
const info = ref({})
const actionInfo = ref(null)
const variableOptions = ref([])
const showDesc = computed(() => {
  if (actionInfo.value) {
    return actionInfo.value.desc
  }
  return info.value.description
})

let nodeParams = {}

function init() {
  getValueVariableList();
  nodeParams = JSON.parse(props.node.node_params)
  loadPluginParams()
  loadPluginInfo()
}

function loadPluginInfo() {
  getPluginInfo({name: nodeParams?.plugin?.name}).then(res => {
    let data = res?.data || {}
    info.value = data.remote
  })
}

function loadPluginParams() {
  runPlugin({
    name: nodeParams?.plugin?.name || '',
    action: "default/get-schema",
    params: {}
  }).then(res => {
    let data = res?.data || {}
    if (nodeParams?.plugin?.name === 'feishu_bitable') {
      actionInfo.value = data[nodeParams.plugin.params.business]
      //data = actionInfo.value.params || {}
    } else {
      Object.assign(formState, data)
      for (let key in formState) {
        formState[key].value = String(nodeParams?.plugin.params[key] || '')
        formState[key].tags = nodeParams.tag_map ? nodeParams?.plugin?.tag_map[key] : []
      }
    }
  })
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
    params: {},
    tag_map: {}
  }

  for (let key in formState) {
    let value = formState[key].value
    if(formState[key].type == 'string'){
      value = String(value)
    } else if(['number', 'integer'].includes(formState[key].type)){
      value = Number(value)
    }
    data.params[key] = value
    data.tag_map[key] = formState[key].tags
  }

  nodeParams.plugin.params = data.params
  nodeParams.plugin.tag_map = data.tag_map

  setData({
    ...props.node,
    node_params: JSON.stringify(nodeParams)
  })
}

onMounted(() => {
  init()
})
</script>
