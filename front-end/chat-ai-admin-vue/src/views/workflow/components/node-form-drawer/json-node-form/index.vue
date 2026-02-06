<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        :desc="t('desc_json_encode')"
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <div class="gray-block-title">{{ t('label_input') }}</div>
            <div class="array-form-box">
              <at-input
                :options="valueOptions"
                :defaultSelectedList="formState.input_tags"
                :defaultValue="formState.input_variable"
                :checkAnyLevel="true"
                ref="atInputRef"
                :placeholder="t('ph_input_message')"
                @open="showAtList"
                @change="(text, selectedList) => changeValue(text, selectedList)"
              />
            </div>
          </div>

          <div class="gray-block mt16">
            <div class="gray-block-title">{{ t('label_output_custom') }}</div>
            <div class="output-box">
              <div class="output-block">
                <div class="output-item">{{ t('label_param_key') }}</div>
                <div class="output-item">{{ t('label_type') }}</div>
              </div>
              <div class="array-form-box" @mousedown.stop="">
                <div class="form-item-list" v-for="(item, index) in formState.output" :key="index">
                  <a-form-item :label="null" :name="['output', index, 'key']">
                    <div class="flex-block-item" style="gap: 12px">
                      <a-input
                        style="width: 214px"
                        v-model:value="item.key"
                        :placeholder="t('ph_input')"
                      ></a-input>
                      <a-form-item-rest>
                        <a-select
                          @change="onTypeChange(item)"
                          v-model:value="item.typ"
                          :placeholder="t('ph_select')"
                          :showArrow="false"
                          disabled
                          style="width: 214px"
                        >
                          <a-select-option
                            v-for="op in typOptions"
                            :value="op.value"
                            :key="op.value"
                            >{{ op.value }}</a-select-option
                          >
                        </a-select>
                      </a-form-item-rest>

                      <div
                        class="btn-hover-wrap"
                        v-if="item.typ == 'object'"
                        @click="onAddSubs(index)"
                      >
                        <PlusCircleOutlined />
                      </div>

                      <!-- <div class="btn-hover-wrap" @click="onDelOutput(index)">
                        <CloseCircleOutlined />
                      </div> -->
                    </div>
                    <div class="sub-field-box" v-if="item.subs && item.subs.length > 0">
                      <a-form-item-rest>
                        <SubKey :data="item.subs" :level="2" :typOptions="typOptions" />
                      </a-form-item-rest>
                    </div>
                  </a-form-item>
                </div>
                <!-- <a-button @click="handleAddOutPut" :icon="h(PlusOutlined)" block type="dashed"
                  >添加参数</a-button
                > -->
              </div>
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
import { CloseCircleOutlined, PlusOutlined, PlusCircleOutlined } from '@ant-design/icons-vue'
import AtInput from '../../at-input/at-input.vue'
import SubKey from './subs-key.vue'

const { t } = useI18n('views.workflow.components.node-form-drawer.json-node-form.index')

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


const valueOptions = ref([])

function getOptions() {
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let list = nodeModel.getAllParentVariable()

    valueOptions.value = filterCalueOptions(list)
  }
}

function filterCalueOptions(list) {
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
        if (item.typ && (item.typ.startsWith('array<') || item.typ === 'object')) {
          result.push({ ...item })
        }
      }
    }
    
    return result
  }
  
  return traverse(list)
}

const showAtList = (val) => {
  if (val) {
    getOptions()
  }
}


const formState = reactive({
  input_variable: '',
  input_tags: [],
  output: [
    {
      key: '',
      typ: 'string',
      subs: []
    }
  ]
})

const update = () => {
  const data = JSON.stringify({
    json_encode: {
      ...formState
    }
  })

  emit('update-node', {
    ...props.node,
    node_params: data
  })
}

const changeValue = (text, selectedList) => {
  formState.input_tags = selectedList
  formState.input_variable = text
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let json_encode = JSON.parse(dataRaw).json_encode || {}

    json_encode = JSON.parse(JSON.stringify(json_encode))
    getOptions()
    formState.input_variable = json_encode.input_variable
    if (json_encode.output.length > 0) {
      formState.output = json_encode.output
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

const handleAddOutPut = () => {
  formState.output.push({
    key: '',
    typ: 'string',
    subs: [],
    cu_key: getUuid(16)
  })
}

const onDelOutput = (index) => {
  formState.output.splice(index, 1)
}

const onTypeChange = (data) => {
  data.subs = []
}

const onAddSubs = (index) => {
  formState.output[index].subs.push({
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
    lable: 'float',
    value: 'float'
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
    lable: 'array\<boole>',
    value: 'array\<boole>'
  },
  {
    lable: 'array\<float>',
    value: 'array\<float>'
  },
  {
    lable: 'array\<object>',
    value: 'array\<object>'
  }
]

const handleClose = () => {
  emit('close')
}

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import '../form-block.less';
.code-edit-box {
  margin-top: 12px;
  .title-block {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
}
.action-btn {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s ease-in;
    margin-left: 8px;
    &:hover {
      background: #e4e6eb;
    }
  }

.output-box {
  .output-block {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 4px;
    color: #262626;
    .output-item {
      width: 214px;
    }
  }
  .flex-block-item .btn-hover-wrap {
    width: 24px;
    height: 24px;
  }
}
</style>