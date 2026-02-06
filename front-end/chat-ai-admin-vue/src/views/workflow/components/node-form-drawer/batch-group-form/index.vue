<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        :desc="t('desc_batch_group')"
        @close="handleClose"
      >
        <template #runBtn>
          <a-tooltip>
            <template #title>{{ t('btn_run_test') }}</template>
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
              <img src="@/assets/svg/execute.svg" alt="" />
              <span>{{ t('title_execution_settings') }}</span>
            </div>
            <div class="row-form-item">
              <div class="form-label">
                <span>{{ t('label_parallel_count') }}</span>
                <a-tooltip>
                  <template #title>{{ t('tip_parallel_count') }}</template>
                  <QuestionCircleOutlined class="tip-icon" />
                </a-tooltip>
              </div>
              <div>
                <a-input-number
                  style="width: 300px"
                  :max="10"
                  :min="1"
                  :precision="0"
                  :step="1"
                  :placeholder="t('ph_input_value')"
                  v-model:value="formState.chan_number"
                ></a-input-number>
              </div>
            </div>

            <div class="row-form-item">
              <div class="form-label">
                <span>{{ t('label_max_execution_count') }}</span>
                <a-tooltip>
                  <template #title>{{ t('tip_max_execution_count') }}</template>
                  <QuestionCircleOutlined class="tip-icon" />
                </a-tooltip>
              </div>
              <div>
                <a-input-number
                  style="width: 300px"
                  :max="500"
                  :min="1"
                  :precision="0"
                  :step="1"
                  :placeholder="t('ph_input_value')"
                  v-model:value="formState.max_run_number"
                ></a-input-number>
              </div>
            </div>
          </div>
          <div class="gray-block">
            <div class="gray-block-title">
              <img src="@/assets/svg/execute_array.svg" alt="" />
              {{ t('title_execution_array') }}
              <a-tooltip>
                <template #title>{{ t('tip_execution_array') }}</template>
                <QuestionCircleOutlined />
              </a-tooltip>
            </div>
            <div class="output-box">
              <div class="output-block">
                <div class="output-item" style="width: 25%">{{ t('label_param_key') }}</div>
                <div class="output-item" style="width: 30%">{{ t('label_type') }}</div>
                <div class="output-item" style="width: 45%">{{ t('label_param_value') }}</div>
              </div>
              <div class="array-form-box">
                <div
                  class="form-item-list"
                  v-for="(item, index) in formState.batch_arrays"
                  :key="item._id"
                >
                  <a-form-item :label="null">
                    <div class="flex-block-item" style="gap: 12px">
                      <a-input
                        style="width: 25%"
                        v-model:value="item.key"
                        :placeholder="t('ph_input_value')"
                      ></a-input>
                      <a-form-item-rest>
                        <a-input
                          style="width: 30%"
                          v-model:value="item.typ"
                          placeholder=""
                          readonly
                        ></a-input>
                        <a-cascader
                          @change="(val, selectedOptions) => handleBatchArrChange(val, item, selectedOptions)"
                          v-model:value="item.value"
                          @dropdownVisibleChange="onBatchArrVisibleChange"
                          style="width: 45%"
                          :options="batchArraysOptions"
                          :allowClear="false"
                          :displayRender="({ labels }) => labels.join('/')"
                          :field-names="{ children: 'children' }"
                          :placeholder="t('ph_select_value')"
                        />
                      </a-form-item-rest>
                    </div>
                  </a-form-item>
                </div>
              </div>
            </div>
          </div>

          <div class="gray-block">
            <div class="gray-block-title">
              <img src="@/assets/svg/output.svg" alt="" />
              {{ t('title_output') }}
              <a-tooltip>
                <template #title>{{ t('tip_output') }}</template>
                <QuestionCircleOutlined />
              </a-tooltip>
            </div>

            <div class="output-box">
              <div class="output-block">
               <div class="output-item" style="width: 23%">{{ t('label_param_key') }}</div>
                <div class="output-item" style="width: 27%">{{ t('label_type') }}</div>
                <div class="output-item" style="width: 45%">{{ t('label_param_value') }}</div>
              </div>
              <div class="array-form-box">
                <div
                  class="form-item-list"
                  v-for="(item, index) in formState.output"
                  :key="item._id"
                >
                  <a-form-item :label="null">
                    <div class="flex-block-item" style="gap: 12px">
                      <a-input
                        style="width: 25%"
                        v-model:value="item.key"
                        :placeholder="t('ph_input_value')"
                      ></a-input>
                      <a-form-item-rest>
                        <a-input
                          style="width: 30%"
                          v-model:value="item.typ"
                          placeholder=""
                          readonly
                        ></a-input>
                        <a-cascader
                          @change="(val) => handleOutputArrChange(val, item)"
                          v-model:value="item.value"
                          @dropdownVisibleChange="onDropdownVisibleChange"
                          style="width: 45%"
                          :options="outputArrarysOptions"
                          :allowClear="false"
                          :displayRender="({ labels }) => labels.join('/')"
                          :field-names="{ children: 'children' }"
                          :placeholder="t('ph_select_value')"
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
    <RunTest :batch_node_key="nodeId" ref="runTestRef" />
  </NodeFormLayout>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { generateRandomId } from '@/utils/index'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import RunTest from '../../nodes/batch-group-node/components/run-test.vue'
import { ref, reactive, watch, h, onMounted, inject } from 'vue'
import {
  CloseCircleOutlined,
  PlusOutlined,
  QuestionCircleOutlined,
  CaretRightOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'

const { t } = useI18n('views.workflow.components.node-form-drawer.batch-group-form.index')
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
const batchArraysOptions = ref([])
const outputArrarysOptions = ref([])
const runTestRef = ref(null)


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


function getOptions() {
  let list = getNode().getAllParentVariable()

  list = handleOptions(list)
  list = list.filter((item) => item.group_node_key != item.loop_parent_key).filter((item) => item.node_id != item.group_node_key)

  batchArraysOptions.value = filterArrayFields(list)
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

const getOutputOptions = () => {
  let list = getNode().getAllParentVariable()
  list.forEach((item) => {
    item.tags = item.tags || []
  })
 
  // 需要当前分组内的变量
  list = list.filter((item) => item.group_node_key == item.loop_parent_key)

  // 递归处理成后端需要的格式
  const loop = (tree) => {
    tree.forEach((item) => {
      if(item.original_value){
        item.value = item.original_value
      }
      if(item.children && item.children.length) {
        loop(item.children)
      }
    })
  }

  // loop(list)
  outputArrarysOptions.value = handleOptions(list)
}

const onDropdownVisibleChange = (visible) => {
  if (!visible) {
    getOutputOptions()
  }
}

const formRef = ref()

const batch_arrays_default = {
  key: '',
  typ: '',
  value: '',
  _id: generateRandomId(16)
}

const output_default = {
  key: '',
  typ: '',
  value: '',
  _id: generateRandomId(16)
}

const formState = reactive({
  chan_number: 10,
  max_run_number: 500,
  batch_arrays: [],
  output: []
})



function formatQuestionValue(val) {
  if (val) {
    let lists = val.split('.')
    let str1 = lists[0]
    let str2 = lists.filter((item, index) => index > 0).join('.')
    return [str1, str2]
  }
  return []
}

const onBatchArrVisibleChange = (visible) => {
  if(visible){
    getOptions()
  }
} 

const handleBatchArrChange = (val, item, selectedOptions) => {
  if(selectedOptions){
    item.typ = selectedOptions[selectedOptions.length - 1].typ 
  }else{
    item.typ = ''
  }
}

const handleAddOutputArrays = () => {
  formState.output.push({
    key: '',
    typ: '',
    value: '',
    _id: generateRandomId(16)
  })
}
const onDelOutputArrays = (index) => {
  formState.output.splice(index, 1)
}

const handleOutputArrChange = (val, item) => {
  if (val && val.length) {
    let filterItem1 = outputArrarysOptions.value.filter((it) => it.value == val[0])[0].children
    if (filterItem1 && filterItem1.length) {
      let typ = filterItem1.filter((it) => it.value == val[1])[0].typ
      item.typ = item.typ = typ.startsWith('array') ? typ : `array\<${typ}>`
    }
  } else {
    item.typ = ''
  }
}

const handleClose = () => {
  emit('close')
}

const handleOpenTestModal = () => {
  let data = formState.batch_arrays[0]

  if (data && data.key && data.typ && data.value) {
    runTestRef.value?.open()
  } else {
    message.error(t('msg_fill_execution_array'))
  }
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let batch = JSON.parse(dataRaw).batch || {}

    getOptions()
    getOutputOptions()

    formState.chan_number = batch.chan_number
    formState.max_run_number = batch.max_run_number
    formState.batch_arrays = batch.batch_arrays.map((item) => {
      return {
        ...item,
        value: formatQuestionValue(item.value)
      }
    })

    formState.output = batch.output.map((item) => {
      return {
        ...item,
        value: formatQuestionValue(item.value)
      }
    })

    if (formState.batch_arrays.length == 0) {
      formState.batch_arrays = [batch_arrays_default]
    }

    if (formState.output.length == 0) {
      formState.output = [output_default]
    }
  } catch (error) {
    console.log(error)
  }
}

const update = () => {
  const data = JSON.stringify({
    batch: {
      ...formState,
      batch_arrays: formState.batch_arrays.map((item) => {
        return {
          ...item,
          value: item.value ? item.value.join('.') : ''
        }
      }),
      chan_number: formState.chan_number,
      max_run_number: formState.max_run_number,
      output: formState.output.map((item) => {
        let value = item.value ? item.value.join('.') : ''

        return {
          ...item,
          value: value
        }
      })
    }
  })

  emit('update-node', {
    ...props.node,
    node_params: data
  })
}

watch(
  () => formState,
  () => {
    update()
  },
  { deep: true }
)

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import '../form-block.less';

.row-form-item {
  margin-bottom: 12px;

  &:last-child {
    margin-bottom: 0;
  }

  .form-label {
    display: flex;
    line-height: 22px;
    font-size: 14px;
    margin-bottom: 4px;
    color: #262626;

    .tip-icon {
      margin-left: 2px;
      color: #8c8c8c;
    }
  }
}

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
