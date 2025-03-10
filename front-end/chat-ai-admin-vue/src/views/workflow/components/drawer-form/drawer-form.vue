<template>
  <div>
    <a-drawer
      :title="drawerTitle"
      placement="right"
      :width="600"
      :open="open"
      class="workflow-drawer-form"
      @close="onBeforeClose"
    >
      <div class="draw-content-body" v-if="open">
        <component
          :properties="properties"
          :is="typeComponentMap[currentComponent]"
          ref="currentComponentRef"
        ></component>
      </div>

      <template #footer>
        <div class="drawer-footer">
          <!-- <a-button key="back" @click="test">test数据</a-button> -->
          <a-button key="back" @click="onClose">取消</a-button>
          <a-button key="submit" type="primary" @click="onSave">保存</a-button>
        </div>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, provide, createVNode } from 'vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import StartNodeForm from './start-node-form.vue'
import MessageNodeForm from './message-node-form.vue'
import QuestionNodeForm from './question-node-form.vue'
import ActionNodeFrom from './action-node-form.vue'
import QaNodeFrom from './qa-node-form.vue'
import transferNodeForm from './transfer-node-form.vue'
import { message, Modal } from 'ant-design-vue'

onMounted(() => {})

const emit = defineEmits(['change'])
const drawerTitle = ref('开始节点')
const open = ref(false)

const currentComponent = ref('')
const typeComponentMap = {
  'start-node': StartNodeForm,
  'message-node': MessageNodeForm,
  'question-node': QuestionNodeForm,
  'action-node': ActionNodeFrom,
  'qa-node': QaNodeFrom,
  'transfer-node': transferNodeForm,
}

const titleMap = {
  'start-node': '开始节点',
  'message-node': '消息节点',
  'question-node': '问题节点',
  'action-node': '动作节点-发送通知',
  'qa-node': '问答节点',
  'transfer-node': '动作节点-转人工',
}
const currentComponentRef = ref(null)

const properties = reactive({})

let updateNum = 0
const onShow = (data) => {
  Object.assign(properties, JSON.parse(JSON.stringify(data.properties)))
  let type = data.type
  if (data.properties.node_sub_type == 52) {
    // 转人工
    type = 'transfer-node'
  }
  currentComponent.value = type
  drawerTitle.value = titleMap[data.type]
  if (data.properties.node_sub_type == 51) {
    // 转人工
    drawerTitle.value = '触发后结束对话'
  }
  updateNum = 0
  open.value = true
}

const onClose = () => {
  open.value = false
}

const onBeforeClose = () => {
  if (updateNum > 0) {
    Modal.confirm({
      title: '关闭确认?',
      icon: createVNode(ExclamationCircleOutlined),
      content: createVNode('div', { style: 'color:red;' }, '您当前有未保存的节点信息,是否关闭?'),
      onOk() {
        open.value = false
      },
    })
  } else {
    open.value = false
  }
}
const onSave = () => {
  currentComponentRef.value.onSave()
}
const updateNodeItem = (val) => {
  // 更新节点
  let newState = JSON.parse(JSON.stringify(val))
  Object.assign(properties, newState)
  // console.log(properties,'===')
  emit('change', {
    ...properties,
  })

  onClose()
}

const updateModifyNum = (data) => {
  updateNum = data
}

provide('nodeInfo', {
  updateNodeItem,
  updateModifyNum,
})

defineExpose({
  onShow,
  onClose,
})
</script>
<style lang="less">
.workflow-drawer-form .ant-drawer-header-title {
  flex-direction: row-reverse;
}
</style>
<style lang="less" scoped>
.drawer-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
}
</style>
