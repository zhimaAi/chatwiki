<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        desc="运行一段JS代码，将代码return的数据输出到下一节点。一般用于进行数据处理。"
        @close="handleClose"
      >
        <template #runBtn>
          <a-tooltip>
            <template #title>运行测试</template>
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
            <div class="gray-block-title">输入</div>
            <div class="array-form-box">
              <div class="form-item-label">自定义输入参数</div>
              <div class="form-item-list" v-for="(item, index) in formState.params" :key="index">
                <a-form-item :label="null" :name="['params', index, 'variable']">
                  <div class="flex-block-item">
                    <a-form-item-rest>
                      <a-input
                        style="width: 190px"
                        v-model:value="item.field"
                        placeholder="请输入参数KEY"
                      ></a-input>
                    </a-form-item-rest>

                    <a-cascader
                      v-model:value="item.variable"
                      @dropdownVisibleChange="onDropdownVisibleChange"
                      style="width: 400px"
                      :options="variableOptions"
                      :allowClear="false"
                      :displayRender="({ labels }) => labels.join('/')"
                      :field-names="{ children: 'children' }"
                      placeholder="请选择"
                    />

                    <div class="btn-hover-wrap" @click="onDelParams(index)">
                      <CloseCircleOutlined />
                    </div>
                  </div>
                </a-form-item>
              </div>
              <a-button @click="handleAddParams" :icon="h(PlusOutlined)" block type="dashed"
                >添加参数</a-button
              >
            </div>

            <div class="code-edit-box">
              <div class="title-block">
                <div>JavaScript 代码</div>
                <a-flex :gap="16" align="center">
                  <a @click="resetTemp">还原为模板</a>
                  <div class="btn-hover-wrap" @click="openEditCodeBox"><FullscreenOutlined /></div>
                </a-flex>
              </div>
              <div class="code-box">
                <CodeEditBox v-model:value="formState.main_func" :width="410" :height="170" />
              </div>
            </div>

            <a-form-item name="timeout" class="mt16">
              <template #label>
                <a-flex :gap="2">超时时长</a-flex>
              </template>
              <div class="flex-block-item">
                <a-input-number
                  placeholder="请输入请求地址"
                  style="width: 138px"
                  :precision="0"
                  v-model:value="formState.timeout"
                  :min="0"
                  :max="3000"
                />
                秒
              </div>
            </a-form-item>
          </div>

          <div class="gray-block mt16">
            <div class="gray-block-title">输出 (自定义输出字段)</div>
            <div class="output-box">
              <div class="output-block">
                <div class="output-item">参数Key</div>
                <div class="output-item">类型</div>
              </div>
              <div class="array-form-box" @mousedown.stop="">
                <div class="form-item-list" v-for="(item, index) in formState.output" :key="index">
                  <a-form-item :label="null" :name="['output', index, 'key']">
                    <div class="flex-block-item" style="gap: 12px">
                      <a-input
                        style="width: 214px"
                        v-model:value="item.key"
                        placeholder="请输入"
                      ></a-input>
                      <a-form-item-rest>
                        <a-select
                          @change="onTypeChange(item)"
                          v-model:value="item.typ"
                          placeholder="请选择"
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

                      <div class="btn-hover-wrap" @click="onDelOutput(index)">
                        <CloseCircleOutlined />
                      </div>
                    </div>
                    <div class="sub-field-box" v-if="item.subs && item.subs.length > 0">
                      <a-form-item-rest>
                        <SubKey :data="item.subs" :level="2" :typOptions="typOptions" />
                      </a-form-item-rest>
                    </div>
                  </a-form-item>
                </div>
                <a-button @click="handleAddOutPut" :icon="h(PlusOutlined)" block type="dashed"
                  >添加参数</a-button
                >
              </div>
            </div>
          </div>

          <div class="gray-block mt16">
            <div class="gray-block-title">异常处理</div>
            <div>运行代码报错时执行该分支</div>
          </div>

          <a-modal v-model:open="open" title="编辑代码" :width="746">
            <template #footer>
              <a-button type="primary" @click="handleOk">确定</a-button>
            </template>
            <div style="margin: 40px 0 24px 0">
              <CodeEditBox
                ref="codeEditBoxRef"
                v-model:value="formState.main_func"
                :width="698"
                :height="472"
              />
            </div>
          </a-modal>
        </a-form>
      </div>
    </div>
    <RunTest :variableOptions="variableOptions" ref="runTestRef" />
  </NodeFormLayout>
</template>

<script setup>
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import { ref, reactive, watch, h, onMounted } from 'vue'
import {
  CloseCircleOutlined,
  PlusOutlined,
  PlusCircleOutlined,
  FullscreenOutlined,
  CaretRightOutlined
} from '@ant-design/icons-vue'

import SubKey from './subs-key.vue'
import CodeEditBox from './code-edit-box.vue'
import RunTest from './run-test.vue'

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

const variableOptions = ref([])

const defaultCode = `function main({data1, data2}){
	return {
		data1,
    	data2
	}
}`

function getOptions() {
  // const node = props.lf.getNodeDataById(props.nodeId)
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let list = nodeModel.getAllParentVariable()

    variableOptions.value = handleOptions(list)
  }
  console.log(variableOptions.value, '==')
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

const formRef = ref()

const formState = reactive({
  main_func: '',
  params: [
    {
      field: '',
      variable: ''
    }
  ],
  timeout: 30,
  output: [
    {
      key: '',
      typ: '',
      subs: []
    }
  ],
  exception: ''
})

const resetTemp = () => {
  formState.main_func = defaultCode
}

let specialNodeList = [
  'special.lib_paragraph_list',
  'special.llm_reply_content',
  'specify-reply-node'
]

function recursionData(data) {
  data.forEach((item) => {
    item.cu_key = Math.random() * 10000
    if (item.subs && item.subs.length) {
      recursionData(item.subs)
    } else {
      item.subs = []
    }
  })
  return data
}

const update = () => {
  const data = JSON.stringify({
    code_run: {
      ...formState,
      params: formatParams()
    }
  })

  emit('update-node', {
    ...props.node,
    ...formState,
    node_params: data
  })
}

function formatParams() {
  let list = []
  formState.params.forEach((item) => {
    list.push({
      field: item.field,
      variable: item.variable && item.variable.length > 0 ? item.variable.join('.') : ''
    })
  })
  return list
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let code_run = JSON.parse(dataRaw).code_run || {}

    code_run = JSON.parse(JSON.stringify(code_run))

    getOptions()

    for (let key in code_run) {
      if (key == 'params') {
        if (code_run[key] && code_run[key].length > 0) {
          formState[key] = code_run[key].map((item) => {
            return {
              ...item,
              cu_key: Math.random() * 10000
            }
          })
        } else {
          formState[key] = []
        }
        continue
      }

      if (key == 'output') {
        formState['output'] = recursionData(code_run[key])
        continue
      }
      formState[key] = code_run[key]
    }
    formState.main_func = formState.main_func || defaultCode

    formState.params = formState.params.map((it) => {
      let specialKey = ''
      for (let i = 0; i < specialNodeList.length; i++) {
        if (it.variable.indexOf(specialNodeList[i]) > -1) {
          specialKey = specialNodeList[i]
          break
        }
      }
      if (specialKey != '') {
        let arr = it.variable.split('.')
        it.variable = [arr[0], specialKey]
      } else {
        it.variable = it.variable.split('.')
      }
      return {
        ...it,
        key: Math.random() * 10000
      }
    })
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

const handleAddParams = () => {
  formState.params.push({
    key: '',
    value: '',
    cu_key: Math.random() * 10000
  })
}

const onDelParams = (index) => {
  formState.params.splice(index, 1)
}

const handleAddOutPut = () => {
  formState.output.push({
    key: '',
    typ: 'string',
    subs: [],
    cu_key: Math.random() * 10000
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
    cu_key: Math.random() * 10000
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

const onDropdownVisibleChange = (visible) => {
  if (!visible) {
    getOptions()
  }
}

const open = ref(false)
const codeEditBoxRef = ref(null)
const openEditCodeBox = () => {
  open.value = true
  codeEditBoxRef.value && codeEditBoxRef.value.handleRefresh()
}
const handleOk = () => {
  open.value = false
}
const handleClose = () => {
  emit('close')
}

const runTestRef = ref(null)
const handleOpenTestModal = () => {
  runTestRef.value.open(JSON.parse(JSON.stringify(formState)), props.node.loop_parent_key != '')
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
