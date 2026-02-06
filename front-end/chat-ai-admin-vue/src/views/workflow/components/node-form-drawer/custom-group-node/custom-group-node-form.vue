<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        :desc="t('desc_custom_group')"
        @close="handleClose"
      >
        <template #runBtn>
          <a-tooltip>
            <template #title>{{ t('tooltip_run_test') }}</template>
            <div class="action-btn" @click="handleOpenTestModal">
              <CaretRightOutlined style="color: rgb(0, 173, 58)" />
            </div>
          </a-tooltip>
        </template>
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <div class="gray-block-title">
              <img src="@/assets/svg/loop-icon.svg" alt="" />
              {{ t('label_loop_settings') }}
              <a-tooltip>
                <template #title>{{ t('tip_loop_settings') }}</template>
                <QuestionCircleOutlined />
              </a-tooltip>
            </div>
            <div class="flex-form-item">
              <div class="form-label">{{ t('label_loop_type') }}</div>
              <a-select v-model:value="formState.loop_type" style="width: 170px">
                <a-select-option value="array">{{ t('label_use_array') }}</a-select-option>
                <a-select-option value="number">{{ t('label_specify_count') }}</a-select-option>
              </a-select>
              <a-input-number
                :placeholder="t('ph_input')"
                v-if="formState.loop_type == 'number'"
                style="flex: 1"
                v-model:value="formState.loop_number"
                :min="1"
                :max="100"
              />
            </div>
          </div>
          <div class="gray-block" v-if="formState.loop_type == 'array'">
            <div class="gray-block-title">
              <img src="@/assets/svg/loop-2.svg" alt="" />
              {{ t('label_loop_array') }}
              <a-tooltip>
                <template #title>{{ t('tip_loop_array') }}</template>
                <QuestionCircleOutlined />
              </a-tooltip>
            </div>
            <div class="output-box">
              <div class="output-block">
                <div class="output-item">{{ t('label_param_key') }}</div>
                <div class="output-item">{{ t('label_type') }}</div>
                <div class="output-item" style="flex: 1">{{ t('label_param_value') }}</div>
              </div>
              <div class="array-form-box">
                <div
                  class="form-item-list"
                  v-for="(item, index) in formState.loop_arrays"
                  :key="item.cu_key"
                >
                  <a-form-item :label="null">
                    <div class="flex-block-item" style="gap: 12px">
                      <a-input
                        style="width: 25%"
                        v-model:value="item.key"
                        :placeholder="t('ph_input')"
                      ></a-input>
                      <a-form-item-rest>
                        <a-input
                          style="width: 25%"
                          v-model:value="item.typ"
                          placeholder=""
                          readonly
                        ></a-input>
                        <a-cascader
                          v-model:value="item.value"
                          @change="(val, selectedOptions) => handleLoopArrChange(val, item, selectedOptions)"
                          @dropdownVisibleChange="onDropdownVisibleChange"
                          style="width: 50%"
                          :options="loopArraysOptions"
                          :allowClear="false"
                          :displayRender="({ labels }) => labels.join('/')"
                          :field-names="{ children: 'children' }"
                          :placeholder="t('ph_select')"
                        />
                      </a-form-item-rest>
                      <div class="btn-hover-wrap" v-if="false" @click="onDelLoopArrays(index)">
                        <CloseCircleOutlined />
                      </div>
                    </div>
                  </a-form-item>
                </div>
                <a-button
                  v-if="formState.loop_arrays.length == 0"
                  @click="handleAddLoopArrays"
                  :icon="h(PlusOutlined)"
                  block
                  type="dashed"
                  >{{ t('btn_add_param') }}</a-button
                >
              </div>
            </div>
          </div>

          <div class="gray-block">
            <div class="gray-block-title">
              <img src="@/assets/svg/loop-3.svg" alt="" />
              {{ t('label_intermediate_variable') }}
              <a-tooltip>
                <template #title>{{ t('tip_intermediate_variable') }}</template>
                <QuestionCircleOutlined />
              </a-tooltip>
            </div>

            <div class="output-box">
              <div class="output-block">
                <div class="output-item">{{ t('label_param_key') }}</div>
                <div class="output-item" style="margin-left: 4px;">{{ t('label_type') }}</div>
                <div class="output-item" style="flex: 1">{{ t('label_param_value') }}</div>
              </div>
              <div class="array-form-box">
                <div
                  class="form-item-list"
                  v-for="(item, index) in formState.intermediate_params"
                  :key="item.cu_key"
                >
                  <a-form-item :label="null">
                    <div class="flex-block-item" style="gap: 12px">
                      <a-input
                        style="width: 25%"
                        v-model:value="item.key"
                        :placeholder="t('ph_input')"
                      ></a-input>
                      <a-form-item-rest>
                        <a-select v-model:value="item.typ" :placeholder="t('ph_select')" style="width: 25%">
                          <a-select-option
                            v-for="op in typOptions"
                            :value="op.value"
                            :key="op.value"
                            >{{ op.value }}</a-select-option
                          >
                        </a-select>
                        <div style="width: 40%">
                          <at-input
                            inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                            :ref="(el) => setAtInputRef(el, 'intermediate_params', index)"
                            :options="variableOptions"
                            :defaultSelectedList="item.tags"
                            :defaultValue="item.value"
                            @open="getVlaueVariableList"
                            @change="
                              (text, selectedList) => changeValue(text, selectedList, item, index)
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
                      </a-form-item-rest>
                      <div class="btn-hover-wrap" @click="onDelIntermediateParams(index)">
                        <CloseCircleOutlined />
                      </div>
                    </div>
                  </a-form-item>
                </div>
                <a-button
                  @click="handleAddIntermediateParams"
                  :icon="h(PlusOutlined)"
                  block
                  type="dashed"
                  >{{ t('btn_add_param') }}</a-button
                >
              </div>
            </div>
          </div>

          <div class="gray-block">
            <div class="gray-block-title">
              <img src="@/assets/svg/output.svg" alt="" />
              {{ t('label_output') }}
              <a-tooltip>
                <template #title>{{ t('tip_output') }}</template>
                <QuestionCircleOutlined />
              </a-tooltip>
            </div>

            <div class="output-box">
              <div class="output-block">
                <div class="output-item">{{ t('label_param_key') }}</div>
                <div class="output-item" style="margin-left: 4px;">{{ t('label_type') }}</div>
                <div class="output-item" style="width: 38%">{{ t('label_param_value') }}</div>
              </div>
              <div class="array-form-box">
                <div
                  class="form-item-list"
                  v-for="(item, index) in formState.output"
                  :key="item.cu_key"
                >
                  <a-form-item :label="null">
                    <div class="flex-block-item" style="gap: 12px">
                      <a-input
                        style="width: 25%"
                        v-model:value="item.key"
                        :placeholder="t('ph_input')"
                      ></a-input>
                      <a-form-item-rest>
                        <a-input
                          style="width: 25%"
                          v-model:value="item.typ"
                          placeholder=""
                          readonly
                        ></a-input>
                        <a-cascader
                          @change="(val) => handleOutputArrChange(val, item)"
                          v-model:value="item.value"
                          @dropdownVisibleChange="onDropdownVisibleChange"
                          style="width: 40%"
                          :options="outputOptions"
                          :allowClear="false"
                          :displayRender="({ labels }) => labels.join('/')"
                          :field-names="{ children: 'children' }"
                          :placeholder="t('ph_select')"
                        />
                      </a-form-item-rest>
                      <div class="btn-hover-wrap" @click="onDelOutputArrays(index)">
                        <CloseCircleOutlined />
                      </div>
                    </div>
                  </a-form-item>
                </div>
                <a-button
                  @click.stop="handleAddOutputArrays"
                  :icon="h(PlusOutlined)"
                  block
                  type="dashed"
                  >{{ t('btn_add_param') }}</a-button
                >
              </div>
            </div>
          </div>
        </a-form>
      </div>
    </div>
    <RunTest :loop_node_key="nodeId" ref="runTestRef" />
  </NodeFormLayout>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import { ref, reactive, watch, h, onMounted, inject, computed } from 'vue'
import {
  CloseCircleOutlined,
  PlusOutlined,
  QuestionCircleOutlined,
  CaretRightOutlined
} from '@ant-design/icons-vue'
import AtInput from '../../at-input/at-input.vue'
import RunTest from '../../nodes/custom-group-node/components/run-test.vue'

const { t } = useI18n('views.workflow.components.node-form-drawer.custom-group-node.custom-group-node-form')

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

const getNode = inject('getNode')

const variableOptions = ref([])

const loopArraysOptions = ref([])

const outputArrarysOptions = ref([])
const runTestRef = ref(null)
const atInputRefs = reactive({})

const handleOpenTestModal = () => {
  runTestRef.value.open()
}

const setAtInputRef = (el, name, index) => {
  if (el) {
    let key = `at_input_${name}_${index}`
    atInputRefs[key] = el
  }
}

const changeValue = (text, selectedList, item) => {
  item.tags = selectedList
  item.value = text
}

// 递归处理Options
function handleOptions(options) {
  options.forEach((item) => {
    if (item.typ == 'node') {
      if (item.node_type == 1) {
        item.value = 'global'
      } else {
        item.value = item.node_id
      }
    } else {
      item.value = item.key
    }

    if (item.children && item.children.length > 0) {
      item.children = handleOptions(item.children)
    }
  })

  return options
}
const getVlaueVariableList = () => {
  let list = getNode().getAllParentVariable()
  list.forEach((item) => {
    item.tags = item.tags || []
  })
  // 需要当前分组外面的变量
  list = list
    .filter((item) => item.group_node_key != item.loop_parent_key)
    .filter((item) => item.node_id != item.group_node_key)
  variableOptions.value = list
}

function getOptions() {
  let list = getNode().getAllParentVariable()

  list = handleOptions(list)
  list = list.filter((item) => item.group_node_key != item.loop_parent_key).filter((item) => item.node_id != item.group_node_key)
  
  loopArraysOptions.value = filterArrayFields(list)
}

function filterArrayFields(list) {
  const filterNodes = (nodes) => {
    if (!nodes || nodes.length === 0) {
      return [];
    }

    return nodes.map(node => {
      // 递归过滤子节点
      const filteredChildren = filterNodes(node.children);

      // 检查当前节点是否包含 'array' 类型，或者其子节点过滤后是否还有内容
      if (node.typ.includes('array') || (filteredChildren && filteredChildren.length > 0)) {
        return {
          ...node,
          children: filteredChildren
        };
      }

      return null;
    }).filter(Boolean); // 过滤掉 null 值
  };

  return filterNodes(list);
}

const onDropdownVisibleChange = (visible) => {
  if (!visible) {
    getOptions()
  }
}

const formRef = ref()

const loop_arrays_default = {
  key: '',
  typ: '',
  value: '',
  cu_key: Math.random() * 10000
}

const intermediate_params_default = {
  key: '',
  typ: '',
  value: '',
  cu_key: Math.random() * 10000
}

const output_default = {
  key: '',
  typ: '',
  value: '',
  cu_key: Math.random() * 10000
}

const formState = reactive({
  loop_type: 'array',
  loop_number: '',
  loop_arrays: [],
  intermediate_params: [],
  output: []
})

const outputOptions = computed(() => {
  let list = formState.intermediate_params
  let childList = []
  
  list.forEach((item) => {
    if (item.key) {
      childList.push({
        children: [],
        id: props.nodeId,
        key: item.key,
        label: item.key,
        node_id: props.nodeId,
        node_name: props.node.node_name,
        node_type: '25',
        original_value: props.nodeId + '.' + item.key,
        text: item.key,
        typ: item.typ,
        value: item.key
      })
    }
  })

  if (childList.length > 0) {
    return [
      {
        label: props.node.node_name,
        node_id: props.nodeId,
        node_type: '25',
        typ: 'node',
        value: props.nodeId,
        children: childList
      }
    ]
  }

  return []
})

const update = () => {
  const data = JSON.stringify({
    loop: {
      ...formState,
      loop_arrays: formState.loop_arrays.map((item) => {
        return {
          ...item,
          value: item.value ? item.value.join('.') : ''
        }
      }),
      output: formState.output.map((item) => {
        return {
          ...item,
          value: item.value ? item.value.join('.') : ''
        }
      })
    }
  })

  emit('update-node', {
    ...props.node,
    ...formState,
    node_params: data
  })
}

function formatQuestionValue(val) {
  if (val) {
    let lists = val.split('.')
    let str1 = lists[0]
    let str2 = lists.filter((item, index) => index > 0).join('.')
    return [str1, str2]
  }
  return []
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let loop = JSON.parse(dataRaw).loop || {}

    loop = JSON.parse(JSON.stringify(loop))

    getVlaueVariableList()
    getOptions()
    formState.loop_type = loop.loop_type
    formState.loop_number = loop.loop_number
    formState.loop_arrays = loop.loop_arrays.map((item) => {
      return {
        ...item,
        value: formatQuestionValue(item.value)
      }
    })
    formState.intermediate_params = loop.intermediate_params.map((item) => {
      return {
        ...item
      }
    })
    formState.output = loop.output.map((item) => {
      return {
        ...item,
        value: formatQuestionValue(item.value)
      }
    })

    if (formState.loop_arrays.length == 0) {
      formState.loop_arrays = [loop_arrays_default]
    }

    if (formState.intermediate_params.length == 0) {
      formState.intermediate_params = [intermediate_params_default]
    }
    if (formState.output.length == 0) {
      formState.output = [output_default]
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

const handleAddLoopArrays = () => {
  formState.loop_arrays.push({
    key: '',
    typ: '',
    value: '',
    cu_key: Math.random() * 10000
  })
}
const onDelLoopArrays = (index) => {
  formState.loop_arrays.splice(index, 1)
}

const handleLoopArrChange = (val, item, selectedOptions) => {
  if(selectedOptions){
    item.typ = selectedOptions[selectedOptions.length - 1].typ 
  }else{
    item.typ = ''
  }
}

const handleAddIntermediateParams = () => {
  formState.intermediate_params.push({
    key: '',
    typ: '',
    value: '',
    cu_key: Math.random() * 10000
  })
}
const onDelIntermediateParams = (index) => {
  formState.intermediate_params.splice(index, 1)
}

const handleAddOutputArrays = () => {
  formState.output.push({
    key: '',
    typ: '',
    value: '',
    cu_key: Math.random() * 10000
  })
}
const onDelOutputArrays = (index) => {
  formState.output.splice(index, 1)
}

const handleOutputArrChange = (val, item) => {
  if (val && val.length) {
    let filterItem1 = outputOptions.value.filter((it) => it.value == val[0])[0].children
    if (filterItem1 && filterItem1.length) {
      item.typ = filterItem1.filter((it) => it.value == val[1])[0].typ
    }
  } else {
    item.typ = ''
  }
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

.flex-form-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.gray-block {
  margin-bottom: 16px;
}

.gray-block-title {
  img {
    width: 16px;
    height: 16px;
  }
}

.at-input-flex1 {
  flex: 1;
  overflow: hidden;
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
      width: 24%;
    }
  }
  .flex-block-item .btn-hover-wrap {
    width: 24px;
    height: 24px;
  }
}
</style>
