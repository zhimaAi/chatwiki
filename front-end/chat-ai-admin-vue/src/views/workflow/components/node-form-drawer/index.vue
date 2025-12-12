<template>
  <transition name="slide-right" appear>
    <div
      v-if="show"
      :width="width"
      :style="{ width: width + 'px' }"
      class="node-form-drawer"
      aria-label="节点表单编辑器"
      role="dialog"
      @keydown.esc="handleClose"
    >
      <div class="node-form-drawer-body">
        <div class="drawer-form" :node-key="props.nodeType">
          <template v-if="nodeFormComponent">
            <component
              :is="nodeFormComponent"
              :lf="props.lf"
              :node-id="props.nodeId"
              :node-type="props.nodeType"
              :node="props.node"
              :key="props.nodeId"
              @close="handleClose"
              @update-node="handleChange"
            />
          </template>
          <a-result
            v-else
            status="warning"
            title="组件加载失败"
            sub-title="无法找到对应的节点表单组件"
          >
            <template #extra>
              <a-button type="primary" @click="handleClose">关闭</a-button>
            </template>
          </a-result>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
/**
 * 节点表单抽屉组件
 * @component NodeFormDrawer
 * @description 用于显示和编辑不同类型节点的表单配置
 * @example
 * <node-form-drawer ref="drawerRef" node-key="start-node" />
 */

import { ref, computed, provide } from 'vue'
import StartNodeForm from './start-node-form/index.vue'
import ProblemOptimizationForm from './problem-optimization-form.vue'
import AiDialogueNodeFrom from './ai-dialogue-node-from.vue'
import AddDataForm from './add-data-form.vue'
import KnowledgeBaseNodeForm from './knowledge-base-node-form.vue'
import CodeRunNodeForm from './code-run-node/code-run-node-form.vue'
import HttpNodeForm from './http-node/http-node-form.vue'
import DeleteDataForm from './delete-data-form.vue'
import VariableAssignmentNodeForm from './variable-assignment-node-form.vue'
import UpdateDataForm from './update-data-form.vue'
import ParameterExtractionNodeForm from './parameter-extraction-form/index.vue'
import SelectDataForm from './select-data-form.vue'
import QuestionNodeForm from './question-node-form.vue'
import JudgeNodeForm from './judge-node-form.vue'
import SpecifyReplyNodeForm from './specify-reply-form.vue'
import McpForm from './mcp-form.vue'
import CustomGroupNodeForm from './custom-group-node/custom-group-node-form.vue'
import ZmPluginsNodeForm from "./zm-plugins-node-form.vue";
import SessionTriggerForm from './session-trigger-form.vue'
import TimingTriggerForm from './timing-trigger-node/timing-trigger-node-form.vue'
import {jsonDecode} from "@/utils/index.js";

// 预定义所有可能的表单组件
const formComponents = {
  'start-node': StartNodeForm,
  'problem-optimization-node': ProblemOptimizationForm,
  'ai-dialogue-node': AiDialogueNodeFrom,
  'add-data-node': AddDataForm,
  'knowledge-base-node': KnowledgeBaseNodeForm,
  'code-run-node': CodeRunNodeForm,
  'http-node': HttpNodeForm,
  'delete-data-node': DeleteDataForm,
  'variable-assignment-node': VariableAssignmentNodeForm,
  'update-data-node': UpdateDataForm,
  'parameter-extraction-node': ParameterExtractionNodeForm,
  'select-data-node': SelectDataForm,
  'question-node': QuestionNodeForm,
  'judge-node': JudgeNodeForm,
  'specify-reply-node': SpecifyReplyNodeForm,
  'mcp-node': McpForm,
  'custom-group': CustomGroupNodeForm,
  'zm-plugins-node': ZmPluginsNodeForm,
  'session-trigger-node': SessionTriggerForm,
  'timing-trigger-node': TimingTriggerForm
  // 其他表单组件可以在这里添加
  // 'problem-optimization-node': defineAsyncComponent(() => import('./problem-optimization-form.vue')),
}

const emit = defineEmits(['update:open', 'update-node', 'change-title'])

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  },
  lf: {
    type: Object,
    default: null
  },
  nodeId: {
    type: String,
    default: ''
  },
  nodeType: {
    type: String,
    default: 'start-node'
  },
  node: {
    type: Object,
    default: () => ({})
  }
})

// 使用计算属性获取对应的表单组件，并添加 getter 和 setter
const show = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value)
})

const close = () => {
  show.value = false
}

const width = computed(() => {
  let node_params = jsonDecode(props.node?.node_params)
  if (node_params?.plugin?.name == 'feishu_bitable') return 600
  let widthMap = {
    'delete-data-node': 600,
    'update-data-node': 600,
    'select-data-node': 600,
    'judge-node': 600,
    'custom-group': 568,
    'timing-trigger-node': 420,
  }
  return widthMap[props.nodeType] || 480
})
const getNode = () => {
  return props.lf.getNodeModelById(props.nodeId)
}

const getGraph = () => {
  return props.lf.graphModel
}

const setData = (data) => {
  handleChange(data)
}

const changeTitle = (title) => {
  emit('change-title', title)
}

const deleteNode = () => {
  show.value = false
  let node = getNode()

  getGraph().eventCenter.emit('custom:node:delete', node)
}

// 使用provide提供lf对象给子组件
provide('lf', props.lf)
provide('getNode', getNode)
provide('getGraph', getGraph)
provide('setData', setData)
provide('changeTitle', changeTitle)
provide('deleteNode', deleteNode)
provide('close', close)

// 使用计算属性获取对应的表单组件
const nodeFormComponent = computed(() => {
  const componentName = `${props.nodeType}`
  return formComponents[componentName] || null
})

const handleChange = (data) => {
  data.dataRaw = data.node_params
  emit('update-node', data)
}

const handleClose = () => {
  show.value = false
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.node-form-drawer {
  position: fixed;
  z-index: 100;
  top: 56px;
  right: 0;
  bottom: 0;
  background-color: #fff;
  box-shadow: 0 0 12px rgba(0, 0, 0, 0.1);
  .node-form-drawer-body {
    position: relative;
    height: 100%;
  }
  .drawer-form {
    height: 100%;
    overflow: hidden;
  }
}
// 定义抽屉滑入滑出动画
.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 0.3s ease;
}

.slide-right-enter-from,
.slide-right-leave-to {
  transform: translateX(100%);
}
</style>
