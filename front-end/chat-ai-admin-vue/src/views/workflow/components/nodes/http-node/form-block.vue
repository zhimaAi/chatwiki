<template>
  <div class="form-block">
    <a-form ref="formRef" layout="vertical" :model="formState">
      <div class="gray-block">
        <div class="gray-block-title">输入</div>
        <a-form-item
          label="请求地址"
          name="rawurl"
          :rules="[{ required: true, message: '请输入请求地址' }]"
        >
          <div class="flex-block-item" @mousedown.stop="">
            <a-form-item-rest>
              <a-select v-model:value="formState.method" style="width: 160px">
                <a-select-option value="POST">POST</a-select-option>
                <a-select-option value="GET">GET</a-select-option>
              </a-select>
            </a-form-item-rest>
            <a-input
              class="flex1"
              v-model:value="formState.rawurl"
              placeholder="请输入请求地址"
            ></a-input>
          </div>
        </a-form-item>
        <div class="array-form-box" @mousedown.stop="">
          <div class="form-item-label">HEADERS</div>
          <div class="form-item-list" v-for="(item, index) in formState.headers" :key="index">
            <a-form-item
              :label="null"
              :name="['headers', index, 'value']"
              :rules="{ required: true, validator: (rule, value) => checkedHeader(rule, value) }"
            >
              <div class="flex-block-item">
                <a-form-item-rest>
                  <a-input
                    style="width: 160px"
                    v-model:value="item.key"
                    placeholder="请输入参数KEY"
                  ></a-input>
                </a-form-item-rest>
                <a-mentions
                  class="flex1"
                  v-model:value="item.value"
                  prefix="/"
                  placeholder="请输入参数值，键入“/”插入变量"
                  :options="variableOptions"
                  @blur="() => (isFocus = false)"
                  @focus="() => (isFocus = true)"
                  @select="onTextChange('headers', index, item.value)"
                >
                  <template #option="{ value, label, payload }">
                    <div class="field-list-item">
                      <div class="field-label">{{ label }}</div>
                      <div class="field-type">{{ payload.typ }}</div>
                    </div>
                  </template>
                </a-mentions>

                <div class="btn-hover-wrap" @click="onDelHeader(index)">
                  <CloseCircleOutlined />
                </div>
              </div>
            </a-form-item>
          </div>
          <a-button @click="handleAddHeader" :icon="h(PlusOutlined)" block type="dashed"
            >添加参数</a-button
          >
        </div>
        <div class="array-form-box" @mousedown.stop="">
          <div class="form-item-label">PARAMS</div>
          <div class="form-item-list" v-for="(item, index) in formState.params" :key="index">
            <a-form-item
              :label="null"
              :name="['params', index, 'value']"
              :rules="{ required: true, validator: (rule, value) => checkedHeader(rule, value) }"
            >
              <div class="flex-block-item">
                <a-form-item-rest>
                  <a-input
                    style="width: 160px"
                    v-model:value="item.key"
                    placeholder="请输入参数KEY"
                  ></a-input>
                </a-form-item-rest>
                <a-mentions
                  class="flex1"
                  v-model:value="item.value"
                  prefix="/"
                  @blur="() => (isFocus = false)"
                  @focus="() => (isFocus = true)"
                  placeholder="请输入参数值，键入“/”插入变量"
                  :options="variableOptions"
                  @select="onTextChange('params', index, item.value)"
                >
                  <template #option="{ value, label, payload }">
                    <div class="field-list-item">
                      <div class="field-label">{{ label }}</div>
                      <div class="field-type">{{ payload.typ }}</div>
                    </div>
                  </template>
                </a-mentions>
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
        <div @mousedown.stop="" v-if="formState.method == 'POST'">
          <a-form-item label="BODY" name="type" class="line-height-22">
            <a-radio-group v-model:value="formState.type">
              <a-radio :value="0">none</a-radio>
              <a-radio :value="1">x-www-form-urlencoded</a-radio>
              <a-radio :value="2">JSON</a-radio>
            </a-radio-group>
          </a-form-item>
          <div class="array-form-box" v-if="formState.type == 1">
            <div class="form-item-list" v-for="(item, index) in formState.body" :key="index">
              <a-form-item
                :label="null"
                :name="['body', index, 'value']"
                :rules="{ required: true, validator: (rule, value) => checkedHeader(rule, value) }"
              >
                <div class="flex-block-item">
                  <a-form-item-rest>
                    <a-input
                      style="width: 160px"
                      v-model:value="item.key"
                      placeholder="请输入参数KEY"
                    ></a-input>
                  </a-form-item-rest>
                  <!-- <a-input
                    class="flex1"
                    v-model:value="item.value"
                    placeholder="请输入参数值，键入“/”插入变量"
                  ></a-input> -->
                  <a-mentions
                    class="flex1"
                    @blur="() => (isFocus = false)"
                    @focus="() => (isFocus = true)"
                    v-model:value="item.value"
                    prefix="/"
                    placeholder="请输入参数值，键入“/”插入变量"
                    :options="variableOptions"
                    @select="onTextChange('body', index, item.value)"
                  >
                    <template #option="{ value, label, payload }">
                      <div class="field-list-item">
                        <div class="field-label">{{ label }}</div>
                        <div class="field-type">{{ payload.typ }}</div>
                      </div>
                    </template>
                  </a-mentions>
                  <div class="btn-hover-wrap" @click="onDelBody(index)">
                    <CloseCircleOutlined />
                  </div>
                </div>
              </a-form-item>
            </div>
            <a-button @click="handleAddBody" :icon="h(PlusOutlined)" block type="dashed"
              >添加参数</a-button
            >
          </div>
          <a-form-item :label="null" name="body_raw" v-if="formState.type == 2">
            <a-mentions
              @wheel.stop=""
              v-model:value="formState.body_raw"
              @blur="() => (isFocus = false)"
              @focus="() => (isFocus = true)"
              prefix="/"
              :rows="4"
              :options="variableOptions"
              @select="onBodyRawChange()"
              placeholder="请输入json, 输入 / 插入变量"
            ></a-mentions>
          </a-form-item>
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
      </div>
      <div class="gray-block mt16">
        <div class="gray-block-title" @click="test">输出 (输出字段提取)</div>
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
    </a-form>
  </div>
</template>

<script setup>
import { ref, reactive, watch, h, inject } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  CloseCircleFilled,
  CloseCircleOutlined,
  LoadingOutlined,
  PlusOutlined,
  PlusCircleOutlined,
  EditOutlined,
  SyncOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons-vue'
import SubKey from './subs-key.vue'

const graphModel = inject('getGraph')
const getNode = inject('getNode')

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['setData'])
const formRef = ref()

const formState = reactive({
  method: 'POST',
  rawurl: '',
  headers: [
    {
      key: '',
      value: ''
    }
  ],
  params: [
    {
      key: '',
      value: ''
    }
  ],
  type: 1,
  body: [],
  body_raw: '',
  timeout: 30,
  output: [
    {
      key: '',
      typ: '',
      subs: []
    }
  ]
})
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
const variableOptions = ref([])
const isFocus = ref(false)
watch(
  () => props.properties,
  (val) => {
    try {
      if (!isFocus.value) {
        getVariableOptions()
      }
      if (lock) {
        return
      }
      let curl = JSON.parse(val.node_params).curl || {}
      curl = JSON.parse(JSON.stringify(curl))
      for (let key in curl) {
        if (key == 'headers' || key == 'params' || key == 'body') {
          if (curl[key] && curl[key].length > 0) {
            formState[key] = curl[key].map((item) => {
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
          formState['output'] = recursionData(curl[key])
          continue
        }
        formState[key] = curl[key]
      }

      lock = true
      setTimeout(() => {
        emit('setData', {
          ...formState,
          node_params: JSON.stringify({
            curl: {
              ...formState
            }
          }),
          height: getNodeHeight()
        })
      }, 100)
    } catch (error) {}
  },
  { immediate: true, deep: true }
)

watch(
  () => formState,
  (val) => {
    emit('setData', {
      ...formState,
      node_params: JSON.stringify({
        curl: {
          ...formState
        }
      }),
      height: getNodeHeight()
    })
  },
  { deep: true }
)
const test = () => {
  console.log(formState.output, '==')
}

function getNodeHeight() {
  let topDefaultHeight = 562
  let height = 0
  if (formState.type == 1) {
    topDefaultHeight = 640 + (formState.body.length - 1) * 36
  }
  if (formState.type == 2) {
    topDefaultHeight = 670
  }
  if (formState.method == 'GET') {
    topDefaultHeight = 505
  }

  height =
    topDefaultHeight + (formState.headers.length - 1) * 36 + (formState.params.length - 1) * 36

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

const handleAddHeader = () => {
  formState.headers.push({
    key: '',
    value: '',
    cu_key: Math.random() * 10000
  })
}

const onDelHeader = (index) => {
  formState.headers.splice(index, 1)
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

const handleAddBody = () => {
  formState.body.push({
    key: '',
    value: '',
    cu_key: Math.random() * 10000
  })
}

const onDelBody = (index) => {
  formState.body.splice(index, 1)
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
  }
  // {
  //   lable: 'array\<object>',
  //   value: 'array\<object>'
  // }
]
const onTextChange = (key, index, data) => {
  let regex = / +【/g
  formState[key][index]['value'] = data.replace(/\//g, '').replace(regex, '【')
}
const onBodyRawChange = () => {
  let regex = / +【/g
  formState.body_raw = formState.body_raw.replace(/\//g, '').replace(regex, '【')
}

function transformArray(arr, parentLabel = '') {
  let result = []

  arr.forEach((item) => {
    let newLabel = parentLabel ? `${parentLabel}.${item.key}` : String(item.key)
    let newValue = parentLabel ? `${parentLabel}.${item.key}` : String(item.key)

    result.push({
      label: newLabel,
      value: newValue,
      payload: { typ: item.typ },
      hasSub: item.subs && item.subs.length > 0
    })

    if (item.subs && Array.isArray(item.subs)) {
      result = result.concat(transformArray(item.subs, newLabel))
    }
  })
  return result
}

function getVariableOptions() {
  let node = getNode()
  let preNodes = graphModel().getNodeIncomingNode(node.id)
  let outOptions = []

  if (preNodes && preNodes.length) {
    preNodes.forEach((item) => {
      let node_type = item.properties.node_type
      let output = item.properties.output
      if (node_type == 4 && output && output.length) {
        outOptions = transformArray(output)
        outOptions = outOptions.filter((it) => !it.hasSub)
        outOptions = outOptions.map((it) => {
          return {
            ...it,
            value: `【${it.value}】`
          }
        })
      }
    })
  }
  let lists = [
    {
      label: '用户消息',
      value: '【global.question】',
      payload: { typ: 'string' }
    },
    {
      label: 'open_id',
      value: '【global.openid】',
      payload: { typ: 'string' }
    },
    ...outOptions
  ]
  variableOptions.value = lists
}

defineExpose({})
</script>

<style lang="less" scoped>
@import '../form-block.less';

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
