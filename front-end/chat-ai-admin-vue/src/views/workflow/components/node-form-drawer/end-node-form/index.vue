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
            <a-segmented
              class="customer-segmented"
              v-model:value="formState.out_type"
              :options="typeOptions"
            />
          </div>
          <div class="gray-block" style="margin-top: 12px" v-if="formState.out_type == 'variable'">
            <div class="gray-block-title">输出变量</div>

            <div class="array-form-box">
              <div class="field-key-title">
                <div class="field-key" style="width: 200px">参数key</div>
                <div class="field-key" style="width: 114px">类型</div>
                <div class="field-key" style="flex: 1">参数值</div>
              </div>
              <div class="form-item-list">
                <a-form-item :label="null">
                  <div class="flex-block-item" style="gap: 8px">
                    <a-input
                      style="width: 200px"
                      disabled
                      :value="'content'"
                      placeholder="请输入"
                    ></a-input>
                    <a-form-item-rest>
                      <a-select
                        :value="'string'"
                        disabled
                        placeholder="请选择"
                        style="width: 114px"
                      >
                        <a-select-option
                          v-for="op in typOptions"
                          :value="op.value"
                          :key="op.value"
                          >{{ op.value }}</a-select-option
                        >
                      </a-select>
                    </a-form-item-rest>
                    <div>
                      <ExclamationCircleFilled style="color: #faad14" />
                      content的返回值请在下方消息中设置
                    </div>
                  </div>
                </a-form-item>
              </div>
              <div class="form-item-list" v-for="(item, index) in formState.outputs" :key="index">
                <a-form-item :label="null">
                  <div class="flex-block-item" style="gap: 8px">
                    <a-input
                      style="width: 200px"
                      v-model:value="item.key"
                      placeholder="请输入"
                    ></a-input>
                    <a-form-item-rest>
                      <a-select
                        @change="onTypeChange(item)"
                        v-model:value="item.typ"
                        placeholder="请选择"
                        style="width: 114px"
                      >
                        <a-select-option
                          v-for="op in typOptions"
                          :value="op.value"
                          :key="op.value"
                          >{{ op.value }}</a-select-option
                        >
                      </a-select>
                    </a-form-item-rest>
                    <div style="width: 200px">
                      <at-input
                        v-if="item.typ != 'object'"
                        :options="filterOutputOption(item)"
                        :defaultSelectedList="item.tags"
                        :defaultValue="item.desc"
                        @open="getVlaueVariableList"
                        @change="
                          (text, selectedList) => changeOutputValue(text, selectedList, item, index)
                        "
                        placeholder="请输入变量值，键入“/”插入变量"
                      >
                        <template #option="{ label, payload }">
                          <div class="field-list-item">
                            <div class="field-label">{{ label }}</div>
                            <div class="field-type">{{ payload.typ }}</div>
                          </div>
                        </template>
                      </at-input>
                    </div>

                    <div
                      class="btn-hover-wrap"
                      v-if="item.typ == 'object'"
                      @click="onAddSubs(index)"
                    >
                      <PlusCircleOutlined />
                    </div>

                    <div class="btn-hover-wrap" @click="onDelOutput(index)">
                      <CloseCircleOutlined />
                    </div>
                  </div>
                  <div class="sub-field-box" v-if="item.subs && item.subs.length > 0">
                    <a-form-item-rest>
                      <SubKey
                        :data="item.subs"
                        :level="2"
                        :typOptions="typOptions"
                        :variableOptions="variableOptions"
                      />
                    </a-form-item-rest>
                  </div>
                </a-form-item>
              </div>
              <a-button @click="handleAddOutPut" :icon="h(PlusOutlined)" block type="dashed"
                >添加参数</a-button
              >
            </div>
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
import {
  CloseCircleOutlined,
  PlusSquareOutlined,
  PlusOutlined,
  PlusCircleOutlined,
  ExclamationCircleFilled
} from '@ant-design/icons-vue'
import AtInput from '../../at-input/at-input.vue'
import SubKey from './subs-key.vue'
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
    label: '返回消息和变量',
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

const changeOutputValue = (text, selectedList, item) => {
  item.tags = selectedList
  item.desc = text
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

function filterOutputOption(item) {
  if (item.typ) {
    return filterCalueOptions(variableOptions.value, item.typ)
  }
  return variableOptions.value
}

function filterCalueOptions(list, typ) {
  function traverse(items) {
    const result = []

    for (const item of items) {
      if (item.children && item.children.length > 0) {
        // 有子节点的情况，递归处理子节点
        const filteredChildren = traverse(item.children)

        if (filteredChildren.length > 0) {
          // 如果过滤后的子节点不为空，保留当前节点并更新其子节点
          const newItem = { ...item, children: filteredChildren }
          result.push(newItem)
        }
      } else {
        // 叶子节点，检查类型是否符合要求
        if (item.typ && item.typ === typ) {
          result.push({ ...item })
        }
      }
    }

    return result
  }

  return traverse(list)
}

const formState = reactive({
  out_type: 'message', // variable
  messages: [
    {
      key: getUuid(16),
      type: 'text',
      content: ''
    }
  ],
  outputs: [
    {
      key: '',
      typ: '',
      subs: []
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
    console.log(finish, '==')
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
    formState.outputs = finish.outputs || []
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

const handleAddOutPut = () => {
  formState.outputs.push({
    key: '',
    typ: 'string',
    subs: [],
    cu_key: getUuid(16)
  })
}

const onDelOutput = (index) => {
  formState.outputs.splice(index, 1)
}

const onTypeChange = (data) => {
  data.subs = []
  data.desc = void 0
  data.tags = []
}

const onAddSubs = (index) => {
  formState.outputs[index].subs.push({
    key: '',
    value: '',
    subs: [],
    cu_key: getUuid(16)
  })
}

const typOptions = [
  {
    lable: 'string',
    value: 'string'
  },
  {
    lable: 'number',
    value: 'number'
  },
  {
    lable: 'boole',
    value: 'boole'
  },
  {
    lable: 'object',
    value: 'object'
  },
  {
    lable: 'array\<string>',
    value: 'array\<string>'
  },
  {
    lable: 'array\<number>',
    value: 'array\<number>'
  },
  {
    lable: 'array\<object>',
    value: 'array\<object>'
  }
]

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

.flex-block-item .btn-hover-wrap {
  width: 24px;
  height: 24px;
}

.field-key-title {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #262626;
  gap: 8px;
  margin-bottom: 4px;
  font-weight: 600;
}
</style>
