<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        desc="工作流完结的节点，在此节点可以定义工作流返回的信息"
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <div class="gray-block-title">输出类型</div>
            <a-segmented :value="formState.out_type" :options="typeOptions" />
          </div>
          <div class="gray-block" style="margin-top: 12px">
            <div class="gray-block-title">自定义消息</div>
            <div class="array-form-box">
              <div class="form-item-label" style="margin-bottom: 8px">
                未自定义消息内容，系统默认将结束节点的上一级大模型节点的输出或者指
                定回复节点的输出作为消息返回。最多添加5条消息。
              </div>
              <div class="form-item-list" v-for="(item, index) in formState.messages" :key="index">
                <a-form-item :label="null">
                  <div class="input-block-item">
                    <div class="input-header">
                      <div>
                        {{ index + 1 }}、
                        <span v-if="item.type == 'text'">文本消息</span>
                        <span v-if="item.type == 'image'">图片消息</span>
                        <span v-if="item.type == 'voice'">语音消息</span>
                      </div>
                      <div class="btn-hover-wrap" @click="onDelItem(index)">
                        <CloseCircleOutlined />
                      </div>
                    </div>
                    <at-input
                      inputStyle="height: 100px;"
                      :ref="(el) => setAtInputRef(el, index)"
                      :options="variableOptions"
                      :defaultSelectedList="item.tags"
                      :defaultValue="item.content"
                      @open="getVlaueVariableList"
                      type="textarea"
                      @change="(text, selectedList) => changeValue(text, selectedList, item, index)"
                      :placeholder="placeholderMap[item.type]"
                    >
                      <template #option="{ label, payload }">
                        <div class="field-list-item">
                          <div class="field-label">{{ label }}</div>
                          <div class="field-type">{{ payload.typ }}</div>
                        </div>
                      </template>
                    </at-input>
                  </div>
                </a-form-item>
              </div>
              <a-dropdown>
                <div>
                  <a-button type="dashed" block
                    ><PlusOutlined />添加消息 （{{ formState.messages.length }} / 5）</a-button
                  >
                </div>

                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="handleAddItem('text')">
                      <a>文本</a>
                    </a-menu-item>
                    <a-menu-item @click="handleAddItem('image')">
                      <a>图片</a>
                    </a-menu-item>
                    <a-menu-item @click="handleAddItem('voice')">
                      <a>语音</a>
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
          </div>
        </a-form>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { getUuid } from '@/utils/index'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import { ref, reactive, watch, h, onMounted } from 'vue'
import { CloseCircleOutlined, PlusSquareOutlined, PlusOutlined } from '@ant-design/icons-vue'
import AtInput from '../../at-input/at-input.vue'
import { message } from 'ant-design-vue'

const emit = defineEmits(['update-node'])
const props = defineProps({
  lf: {
    type: Object,
    default: null
  },
  nodeId: {
    type: String,
    default: ''
  },
  node: {
    type: Object,
    default: () => ({})
  }
})

const placeholderMap = {
  text: '请输入文本消息内容，键入/可以插入变量',
  image: '请输入图片消息的url（系统发送时会自动转成图片发送），键入/可插入变量',
  voice: '请输入语音的url（系统发送时会自动转成语音发送），键入/可插入变量'
}

const typeOptions = [
  {
    label: '返回消息',
    value: 'message'
  },
  {
    label: '返回变量',
    value: 'variable'
  }
]

const variableOptions = ref([])

const atInputRefs = reactive({})
const setAtInputRef = (el, index) => {
  if (el) {
    let key = `at_input_${index}`
    atInputRefs[key] = el
  }
}

const changeValue = (text, selectedList, item) => {
  item.tags = selectedList
  item.content = text
}

const getVlaueVariableList = () => {
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let list = nodeModel.getAllParentVariable()
    list.forEach((item) => {
      item.tags = item.tags || []
    })

    variableOptions.value = list
  }
}

const formState = reactive({
  out_type: 'message', // variable
  messages: [
    {
      key: getUuid(16),
      type: 'text',
      content: ''
    }
  ]
})

const update = () => {
  const data = JSON.stringify({
    finish: {
      ...formState
    }
  })

  emit('update-node', {
    ...props.node,
    node_params: data
  })
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let finish = JSON.parse(dataRaw).finish || {}
    finish = JSON.parse(JSON.stringify(finish))
    getVlaueVariableList()

    formState.out_type = finish.out_type || 'message'
    formState.messages = finish.messages || []

    if (formState.messages.length == 0) {
      formState.messages = [
        {
          key: getUuid(16),
          type: 'text',
          content: ''
        }
      ]
    }
  } catch (error) {
    console.log(error)
  }
}

watch(
  () => formState,
  () => {
    update()
  },
  { deep: true }
)

const handleAddItem = (type) => {
  if (formState.messages.length >= 5) {
    return message.warning('最多添加5条消息')
  }
  formState.messages.push({
    key: getUuid(16),
    type,
    content: ''
  })
}

const onDelItem = (index) => {
  if (formState.messages.length <= 1) {
    return message.warning('最少添加1条消息')
  }
  formState.messages.splice(index, 1)
}

const handleClose = () => {
  emit('close')
}

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import '../form-block.less';
.gray-block-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.input-block-item {
  background: #e4e6eb;
  border-radius: 6px;
  padding: 8px;
  .input-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
}
</style>
