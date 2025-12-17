<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        desc="通过设定批量运行次数和逻辑，运行批处理体内的任务"
        @close="handleClose"
      >
      </NodeFormHeader>
    </template>
    <div class="problem-optimization-form">
      <div class="node-form-content">
        <a-form ref="formRef" layout="vertical" :model="formState">
          <div class="gray-block">
            <div class="gray-block-title">
              <img src="@/assets/svg/execute.svg" alt="" />
              <span>执行设置</span>
            </div>
            <div class="row-form-item">
              <div class="form-label">
                <span>并行运行数量</span>
                <a-tooltip>
                  <template #title>并行运行数量，最大一次执行10个</template>
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
                  placeholder="请输入"
                  v-model:value="formState.chan_number"
                ></a-input-number>
              </div>
            </div>

            <div class="row-form-item">
              <div class="form-label">
                <span>最大执行次数</span>
                <a-tooltip>
                  <template #title>批量执行运行总次数不超过该上限，超过时会直接进入下一个节点</template>
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
                  placeholder="请输入"
                  v-model:value="formState.max_run_number"
                ></a-input-number>
              </div>
            </div>
          </div>
          <div class="gray-block">
            <div class="gray-block-title">
              <img src="@/assets/svg/execute_array.svg" alt="" />
              执行数组
              <a-tooltip>
                <template #title>批处理体中节点要使用的变量，仅支持引用数组，循环次数为数组的长度，执行时会按顺序输出单个数组</template>
                <QuestionCircleOutlined />
              </a-tooltip>
            </div>
            <div class="output-box">
              <div class="output-block">
                <div class="output-item">参数Key</div>
                <div class="output-item">类型</div>
                <div class="output-item" style="flex: 1">参数值</div>
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
                        placeholder="请输入"
                      ></a-input>
                      <a-form-item-rest>
                        <a-input
                          style="width: 30%"
                          v-model:value="item.typ"
                          placeholder=""
                          readonly
                        ></a-input>
                        <a-cascader
                          @change="(val) => handleBatchArrChange(val, item)"
                          v-model:value="item.value"
                          @dropdownVisibleChange="onBatchArrVisibleChange"
                          style="width: 45%"
                          :options="batchArraysOptions"
                          :allowClear="false"
                          :displayRender="({ labels }) => labels.join('/')"
                          :field-names="{ children: 'children' }"
                          placeholder="请选择"
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
              输出
              <a-tooltip>
                <template #title>循环完成后输出的内容，仅支持引用循环体中节点的输出变量，输出的内容会自动组装为数组类型</template>
                <QuestionCircleOutlined />
              </a-tooltip>
            </div>

            <div class="output-box">
              <div class="output-block">
                <div class="output-item">参数Key</div>
                <div class="output-item" style="margin-left: 4px">类型</div>
                <div class="output-item" style="width: 38%">参数值</div>
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
                        placeholder="请输入"
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
                          placeholder="请选择"
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
                  >添加参数</a-button
                >
              </div>
            </div>
          </div>
        </a-form>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { generateRandomId } from '@/utils/index'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import { ref, reactive, watch, h, onMounted, inject, computed } from 'vue'
import {
  CloseCircleOutlined,
  PlusOutlined,
  QuestionCircleOutlined
} from '@ant-design/icons-vue'

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

  batchArraysOptions.value = filterArrayFields(list)
}

function filterArrayFields(list) {
  // 只要数组变量 且 分组以外的变量
  let result = []
  list.forEach((item) => {
    let children = []
    if (item.children && item.children.length > 0) {
      children = item.children.filter((it) => it.typ.toLowerCase().includes('array'))
    }
    if (children.length > 0) {
      result.push({
        ...item,
        children
      })
    }
  })
  return result
    .filter((item) => item.group_node_key != item.loop_parent_key)
    .filter((item) => item.node_id != item.group_node_key)
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

const handleBatchArrChange = (val, item) => {
  if (val && val.length) {
    let filterItem1 = batchArraysOptions.value.filter((it) => it.value == val[0])[0].children
    if (filterItem1 && filterItem1.length) {
      item.typ = filterItem1.filter((it) => it.value == val[1])[0].typ
    }
  } else {
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
