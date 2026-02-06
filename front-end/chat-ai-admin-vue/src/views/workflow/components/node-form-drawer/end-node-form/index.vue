<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        :desc="t('desc_end_node')"
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <div class="gray-block-title">{{ t('title_output_type') }}</div>
            <a-segmented
              class="customer-segmented"
              v-model:value="formState.out_type"
              :options="typeOptions"
            />
          </div>
          <div class="gray-block" style="margin-top: 12px" v-if="formState.out_type == 'variable'">
            <div class="gray-block-title">{{ t('title_output_variable') }}</div>

            <div class="array-form-box">
              <div class="field-key-title">
                <div class="field-key" style="width: 200px">{{ t('label_param_key') }}</div>
                <div class="field-key" style="width: 114px">{{ t('label_type') }}</div>
                <div class="field-key" style="flex: 1">{{ t('label_param_value') }}</div>
              </div>
              <div class="form-item-list">
                <a-form-item :label="null">
                  <div class="flex-block-item" style="gap: 8px">
                    <a-input
                      style="width: 200px"
                      disabled
                      :value="'content'"
                      :placeholder="t('ph_input_value')"
                    ></a-input>
                    <a-form-item-rest>
                      <a-select
                        :value="'string'"
                        disabled
                        :placeholder="t('ph_select_value')"
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
                    <div class="msg_content_return_setting">
                      <ExclamationCircleFilled style="color: #faad14" />
                      {{ t('msg_content_return_setting') }}
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
                      :placeholder="t('ph_input_value')"
                    ></a-input>
                    <a-form-item-rest>
                      <a-select
                        @change="onTypeChange(item)"
                        v-model:value="item.typ"
                        :placeholder="t('ph_select_value')"
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
                        :placeholder="t('ph_input_variable_value')"
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
                >{{ t('btn_add_param') }}</a-button
              >
            </div>
          </div>
          <div class="gray-block" style="margin-top: 12px">
            <div class="gray-block-title">{{ t('title_custom_message') }}</div>
            <div class="array-form-box">
              <div class="form-item-label" style="margin-bottom: 8px">
                {{ t('msg_custom_message_desc') }}
              </div>
              <div class="form-item-list" v-for="(item, index) in formState.messages" :key="index">
                <a-form-item :label="null">
                  <div class="input-block-item">
                    <div class="input-header">
                      <div>
                        {{ index + 1 }}、
                        <span v-if="item.type == 'text'">{{ t('label_text_message') }}</span>
                        <span v-if="item.type == 'image'">{{ t('label_image_message') }}</span>
                        <span v-if="item.type == 'voice'">{{ t('label_voice_message') }}</span>
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
                    ><PlusOutlined />{{ t('btn_add_message', { count: formState.messages.length }) }}</a-button
                  >
                </div>

                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="handleAddItem('text')">
                      <a>{{ t('label_text') }}</a>
                    </a-menu-item>
                    <a-menu-item @click="handleAddItem('image')">
                      <a>{{ t('label_image') }}</a>
                    </a-menu-item>
                    <a-menu-item @click="handleAddItem('voice')">
                      <a>{{ t('label_voice') }}</a>
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
import { useI18n } from '@/hooks/web/useI18n'
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

const { t } = useI18n('views.workflow.components.node-form-drawer.end-node-form.index')
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
  text: t('ph_text_message'),
  image: t('ph_image_message'),
  voice: t('ph_voice_message')
}

const typeOptions = [
  {
    label: t('label_return_message'),
    value: 'message'
  },
  {
    label: t('label_return_message_and_variable'),
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
    return message.warning(t('msg_max_5_messages'))
  }
  formState.messages.push({
    key: getUuid(16),
    type,
    content: ''
  })
}

const onDelItem = (index) => {
  if (formState.messages.length <= 1) {
    return message.warning(t('msg_min_1_message'))
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

.msg_content_return_setting{
  line-height: 16px;
  font-size: 13px;
}
</style>
