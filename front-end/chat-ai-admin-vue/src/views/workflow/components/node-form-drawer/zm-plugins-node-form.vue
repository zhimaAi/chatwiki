<style lang="less" scoped>
@import "./components/node-options";
.node-icon{
  display: block;
  width: 20px;
  height: 20px;
  border-radius: 6px;
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

    <!-- 根据params动态渲染的组件 -->
    <template v-if="actionInfo?.common_template">
      <!--自定义渲染-->
      <DynamicApiBox
        :node="node"
        :action="actionInfo"
        :actionName="nodeParams.plugin.params.business"
        :variableOptions="variableOptions"
        @updateVar="getValueVariableList"
      />
    </template>

    <!--方法插件-->
    <template v-else-if="pluginHasAction(nodeParams?.plugin?.name)">
      <!--自定义渲染-->
      <component
        v-if="pluginCompMap[nodeParams.plugin.name]"
        :is="pluginCompMap[nodeParams.plugin.name]"
        :node="node"
        :action="actionInfo"
        :actionName="nodeParams.plugin.params.business"
        :variableOptions="variableOptions"
        @updateVar="getValueVariableList"
      />
      <!--默认渲染-->
      <PluginFormRender
        v-else
        :node="node"
        :params="actionInfo.params"
        :output="actionInfo.output"
        :variableOptions="variableOptions"
        @updateVar="getValueVariableList"
      />
    </template>
    <!--独立插件-->
    <div v-else class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>输入</div>
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
import FeishuBittableBox from "./components/feishu-bittable/feishu-bittable-box.vue";
import OfficialAccountBox from "./components/official-account/official-account-box.vue";
import OfficialTemplateBox from "./components/official-account/official-template-box.vue";
import OfficialTagBox from "./components/official-account/official-tag-box.vue";
import OfficialSendMessageBox from "./components/official-send-message/official-send-message-box.vue";
import {pluginHasAction} from "@/constants/plugin.js";
import PluginFormRender from "./components/pluginFormRender.vue";
import {sortObjectKeys} from "@/utils/index.js";
import OfficialDraftBox from "./components/official-draft/official-draft-box.vue";
import DynamicApiBox from "./components/dynamic-api/dynamic-api-box.vue";
import OfficialArticleBox from "./components/official-account/official-article-box.vue";

const getNode = inject('getNode')
const setData = inject('setData')

const props = defineProps({
  node: {
    type: Object,
    default: () => ({})
  }
})

const pluginCompMap = {
  feishu_bitable: FeishuBittableBox,
  official_account_profile: OfficialAccountBox,
  official_batch_tag: OfficialTagBox,
  official_send_template_message: OfficialTemplateBox,
  official_send_message: OfficialSendMessageBox,
  official_draft: OfficialDraftBox,
  official_article: OfficialArticleBox,
}

const customFieldsSortMap = {
  'official_article:get_official_article': ['url', 'number']
}

const formState = reactive({})
const info = ref({})
const actionInfo = ref(null)
const variableOptions = ref([])
const isMultiNode = ref(false)
const remoteInfo = ref({})  // 用于临时存储info需要的数据
const showDesc = computed(() => {
  if (actionInfo.value) {
    return actionInfo.value.desc
  }
  return info.value.description
})

let nodeParams = {}

async function init() {
  getValueVariableList();
  nodeParams = JSON.parse(props.node.node_params)
  await loadPluginInfo()
  await loadPluginParams()
}

function loadPluginInfo() {
  return getPluginInfo({name: nodeParams?.plugin?.name}).then(res => {
    let data = res?.data || {}
    remoteInfo.value = data.remote // 用于临时存储info需要的数据
    isMultiNode.value = data?.local?.multiNode || false
  })
}

function loadPluginParams() {
  return runPlugin({
    name: nodeParams?.plugin?.name || '',
    action: "default/get-schema",
    params: {}
  }).then(res => {
    let data = res?.data || {}
    info.value = remoteInfo.value // info会导致页面开始渲染，等loadPluginParams有数据再渲染
    const pName = nodeParams?.plugin?.name
    if (pluginHasAction(pName) || isMultiNode.value) {
      const business = nodeParams.plugin.params.business
      actionInfo.value = data[business]
      let sortFields = customFieldsSortMap[`${pName}:${business}`]
      if (Array.isArray(sortFields)) {
        actionInfo.value.params = sortObjectKeys(actionInfo.value.params, sortFields)
      }
    } else {
      Object.assign(formState, data)
      for (let key in formState) {
        formState[key].value = String(nodeParams?.plugin.params[key] || '')
        formState[key].tags = nodeParams?.plugin?.tag_map?.[key] || []
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
