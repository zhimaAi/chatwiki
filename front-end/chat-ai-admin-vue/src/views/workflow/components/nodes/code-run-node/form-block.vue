<template>
  <div class="form-block">
    <a-form ref="formRef" layout="vertical" :model="formState">
      <div class="gray-block">
        <div class="gray-block-title">输入</div>

        <div class="array-form-box" @mousedown.stop="">
          <div class="form-item-label">自定义输入参数</div>
          <div class="form-item-list" v-for="(item, index) in formState.params" :key="index">
            <a-form-item
              :label="null"
              :name="['params', index, 'variable']"
              :rules="{ required: true, validator: (rule, value) => checkedHeader(rule, value) }"
            >
              <div class="flex-block-item">
                <a-form-item-rest>
                  <a-input
                    style="width: 160px"
                    v-model:value="item.field"
                    placeholder="请输入参数KEY"
                  ></a-input>
                </a-form-item-rest>

                <a-cascader
                  v-model:value="item.variable"
                  @dropdownVisibleChange="onDropdownVisibleChange"
                  style="width: 300px"
                  :options="variableOptions"
                  :allowClear="false"
                  :displayRender="({ labels }) => labels.join('/')"
                  :field-names="{ children: 'children' }"
                  placeholder="请选择"
                />

                <div
                  class="btn-hover-wrap"
                  @click="onDelParams(index)"
                  @mousedown.stop=""
                  @wheel.stop=""
                >
                  <CloseCircleOutlined />
                </div>
              </div>
            </a-form-item>
          </div>
          <a-button @click="handleAddParams" :icon="h(PlusOutlined)" block type="dashed"
            >添加参数</a-button
          >
        </div>
        <div class="code-edit-box"  @mousedown.stop="">
          <div class="title-block">
            <div>JavaScript 代码</div>
            <a-flex :gap="16" align="center">
              <a @click="resetTemp">还原为模板</a>
              <div class="btn-hover-wrap" @click="openEditCodeBox"><FullscreenOutlined /></div>
            </a-flex>
          </div>
          <div class="code-box">
            <CodeEditBox v-model:value="formState.main_func" :width="512" :height="150" />
          </div>
        </div>
      </div>
      <a-form-item name="timeout">
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
      <div class="gray-block mt16">
        <div class="gray-block-title" @click="test">输出 (自定义输出字段)</div>
        <div class="output-box">
          <div class="output-block">
            <div class="output-item">参数Key</div>
            <div class="output-item">类型</div>
          </div>
          <div class="array-form-box" @mousedown.stop="">
            <div class="form-item-list" v-for="(item, index) in formState.output" :key="index">
              <a-form-item
                :label="null"
                :name="['output', index, 'key']"
                :rules="{ required: true, validator: (rule, value) => checkedHeader(rule, value) }"
              >
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
                      <a-select-option v-for="op in typOptions" :value="op.value">{{
                        op.value
                      }}</a-select-option>
                    </a-select>
                  </a-form-item-rest>

                  <div class="btn-hover-wrap" v-if="item.typ == 'object'" @click="onAddSubs(index)">
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
    </a-form>
    <a-modal
      v-model:open="open"
      title="编辑代码"
      :width="746"
    >
      <template #footer>
        <a-button type="primary"  @click="handleOk">确定</a-button>
      </template>
      <div style="margin: 40px 0 24px 0">
        <CodeEditBox ref="codeEditBoxRef" v-model:value="formState.main_func" :width="698" :height="472" />
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, watch, h, inject, toRaw, nextTick, onMounted, onBeforeUnmount } from 'vue'
import {
  CloseCircleOutlined,
  PlusOutlined,
  PlusCircleOutlined,
  FullscreenOutlined
} from '@ant-design/icons-vue'
import SubKey from './subs-key.vue'
import AtInput from '../at-input/at-input.vue'
import CodeEditBox from './code-edit-box.vue'
import { haveOutKeyNode } from '@/views/workflow/components/util.js'

const graphModel = inject('getGraph')
const getNode = inject('getNode')

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['setData'])
const atInputRefs = reactive({})

const setAtInputRef = (el, name, index) => {
  if (el) {
    let key = `at_input_${name}_${index}`
    atInputRefs[key] = el
  }
}

const changeValue = (type, text, selectedList, item) => {
  if (type == 'body_raw') {
    formState.body_raw = text
    formState.body_raw_tags = selectedList
  } else {
    item.tags = selectedList
    item.value = text
  }
}

const getVlaueVariableList = () => {
  // let list = getNode().getAllParentVariable()
  // list.forEach((item) => {
  //   item.tags = item.tags || []
  // })
  // variableOptions.value = list
}

const formRef = ref()

const defaultCode = `function main({data1, data2}){
	return {
		data1,
    	data2
	}
}`

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

const resetTemp = ()=>{
  formState.main_func = defaultCode
}
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
let lock = false
// 特殊节点列表
let specialNodeList = [
  'special.lib_paragraph_list',
  'special.llm_reply_content',
  'specify-reply-node'
]
const variableOptions = ref([])
watch(
  () => props.properties,
  (val) => {
    try {
      if (lock) {
        return
      }
      getVlaueVariableList()

      let dataRaw = val.dataRaw || val.node_params || '{}'
      let code_run = JSON.parse(dataRaw).code_run || {}

      code_run = JSON.parse(JSON.stringify(code_run))

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

      lock = true
      setTimeout(() => {
        emit('setData', {
          ...formState,
          node_params: JSON.stringify({
            code_run: {
              ...formState,
              params: formatParams()
            }
          }),
          height: getNodeHeight()
        })
      }, 100)
    } catch (error) {
      console.log(error)
    }
  },
  { immediate: true, deep: true }
)

watch(
  () => formState,
  (val) => {
    emit('setData', {
      ...formState,
      node_params: JSON.stringify({
        code_run: {
          ...formState,
          params: formatParams()
        }
      }),
      height: getNodeHeight()
    })
  },
  { deep: true }
)

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
const test = () => {
  // console.log(formState.output, '==')
}

function getNodeHeight() {
  let topDefaultHeight = 562
  let height = 0

  // 下面输出字段高度

  let outLens = calculateTotalLength(formState.output)

  return height + 180 + (outLens - 1) * 36
}

function calculateTotalLength(array) {
  let totalLength = array.length // 先加上主数组的长度

  // 遍历数组中的每个对象
  for (const item of array) {
    if (item.subs && Array.isArray(item.subs)) {
      // 如果对象有 subs 属性且 subs 是数组，则递归计算 subs 的长度
      totalLength += calculateTotalLength(item.subs)
    }
  }

  return totalLength
}

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

const checkedHeader = (rule, value) => {
  // if (value == null) {
  //   return Promise.reject('请输入延迟发送时间')
  // }
  // if (!Number.isInteger(value / 0.5)) {
  //   return Promise.reject('必须为0.5秒的倍数')
  // }
  return Promise.resolve()
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

const onUpatateNodeName = (data) => {
  if(!haveOutKeyNode.includes(data.node_type)){
    return
  }

  getVlaueVariableList()

  nextTick(() => {
    if (formState.body_raw_tags && formState.body_raw_tags.length > 0) {
      formState.body_raw_tags.forEach((tag) => {
        if (tag.node_id == data.node_id) {
          let arr = tag.label.split('/')
          arr[0] = data.node_name
          tag.label = arr.join('/')
          tag.node_name = data.node_name
        }
      })
    }

    let keys = ['headers', 'params', 'body']

    keys.forEach((key) => {
      let items = formState[key]

      items.forEach((item) => {
        if (item.tags && item.tags.length > 0) {
          item.tags.forEach((tag) => {
            if (tag.node_id == data.node_id) {
              let arr = tag.label.split('/')
              arr[0] = data.node_name
              tag.label = arr.join('/')
              tag.node_name = data.node_name
            }
          })
        }
      })
    })

    Object.keys(toRaw(atInputRefs)).forEach((key) => {
      if (atInputRefs[key] && atInputRefs[key].$refs.JMention) {
        atInputRefs[key].refresh()
      }
    })
  })
}

const open = ref(false)
const codeEditBoxRef = ref(null)
const openEditCodeBox = () => {
  open.value = true
  codeEditBoxRef.value && codeEditBoxRef.value.handleRefresh()
}

function getOptions() {
  let list = getNode().getAllParentVariable()

  variableOptions.value = handleOptions(list)
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

const onDropdownVisibleChange = (visible) => {
  if (!visible) {
    getOptions()
  }
}
const handleOk = () => {
  open.value = false
}

onMounted(() => {
  const mode = graphModel()
  getOptions()

  mode.eventCenter.on('custom:setNodeName', onUpatateNodeName)
})

onBeforeUnmount(() => {
  const mode = graphModel()

  mode.eventCenter.off('custom:setNodeName', onUpatateNodeName)
})
defineExpose({})
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
