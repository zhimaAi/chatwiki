<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        desc="发起一个http请求，实现与第三方系统数据交互"
        @close="handleClose"
      >
        <template #runBtn>
          <a-tooltip>
            <template #title>运行测试</template>
            <div class="run-test-btn" @click="handleOpenTestModal">
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
            <div class="gray-block-title  export-curl-title">
            <span> 输入</span>
            <a-button ghost type="primary" size="small" @click="showParseModal">
              <CodeOutlined /> 导入cURL
            </a-button>
            </div>

            <a-form-item
              label="请求地址"
              name="rawurl"
              :rules="[{ required: true, message: '请输入请求地址' }]"
            >
              <div class="flex-block-item">
                <a-form-item-rest>
                  <a-select v-model:value="formState.method" style="width: 120px">
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

            <div class="array-form-box">
              <div class="form-item-label">HEADERS</div>
              <div class="form-item-list" v-for="(item, index) in formState.headers" :key="index">
                <a-form-item :label="null" :name="['headers', index, 'value']">
                  <div class="flex-block-item">
                    <a-form-item-rest>
                      <a-input
                        style="width: 120px"
                        v-model:value="item.key"
                        placeholder="请输入参数KEY"
                      ></a-input>
                    </a-form-item-rest>

                    <div class="at-input-flex1">
                      <at-input
                        inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                        :ref="(el) => setAtInputRef(el, 'headers', index)"
                        :options="variableOptions"
                        :defaultSelectedList="item.tags"
                        :defaultValue="item.value"
                        @open="getVlaueVariableList"
                        @change="
                          (text, selectedList) =>
                            changeValue('headers', text, selectedList, item, index)
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
            <div class="array-form-box">
              <div class="form-item-label">PARAMS</div>
              <div class="form-item-list" v-for="(item, index) in formState.params" :key="index">
                <a-form-item :label="null" :name="['params', index, 'value']">
                  <div class="flex-block-item">
                    <a-form-item-rest>
                      <a-input
                        style="width: 120px"
                        v-model:value="item.key"
                        placeholder="请输入参数KEY"
                      ></a-input>
                    </a-form-item-rest>
                    <div class="at-input-flex1">
                      <at-input
                        inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                        :ref="(el) => setAtInputRef(el, 'params', index)"
                        :options="variableOptions"
                        :defaultSelectedList="item.tags"
                        :defaultValue="item.value"
                        @open="getVlaueVariableList"
                        @change="
                          (text, selectedList) =>
                            changeValue('params', text, selectedList, item, index)
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

            <div v-if="formState.method == 'POST'">
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
                    :rules="{
                      required: true
                    }"
                  >
                    <div class="flex-block-item">
                      <a-form-item-rest>
                        <a-input
                          style="width: 120px"
                          v-model:value="item.key"
                          placeholder="请输入参数KEY"
                        ></a-input>
                      </a-form-item-rest>
                      <div class="at-input-flex1">
                        <at-input
                          inputStyle="overflow-y: hidden; overflow-x: scroll; height: 22px;"
                          :ref="(el) => setAtInputRef(el, 'body', index)"
                          :options="variableOptions"
                          :defaultSelectedList="item.tags"
                          :defaultValue="item.value"
                          @open="getVlaueVariableList"
                          @change="
                            (text, selectedList) =>
                              changeValue('body', text, selectedList, item, index)
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
                <at-input
                  type="textarea"
                  :options="variableOptions"
                  :defaultSelectedList="formState.body_raw_tags"
                  :defaultValue="formState.body_raw"
                  :ref="(el) => setAtInputRef(el, 'body', 'body_raw')"
                  @open="getVlaueVariableList"
                  @change="(text, selectedList) => changeValue('body_raw', text, selectedList)"
                  placeholder="请输入变量值，键入“/”插入变量"
                >
                  <template #option="{ label, payload }">
                    <div class="field-list-item">
                      <div class="field-label">{{ label }}</div>
                      <div class="field-type">{{ payload.typ }}</div>
                    </div>
                  </template>
                </at-input>
              </a-form-item>
            </div>

            <a-form-item name="timeout">
              <template #label>
                <a-flex :gap="2">超时时长</a-flex>
              </template>
              <div class="flex-block-item">
                <a-input-number
                  placeholder="请输入请求地址"
                  style="width: 120px"
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
            <div class="gray-block-title">鉴权
              <a-tooltip>
                <template #title>鉴权参数：导出模板CSL文件，或者上架模板时，自动清空参数值</template>
                <QuestionCircleOutlined />
              </a-tooltip>
            </div>

            <div class="output-box">
              <div class="output-block">
                <div class="output-item" style="width: 124px">Key</div>
                <div class="output-item" style="width: 124px">Value</div>
                <div class="output-item" style="width: 124px">Add To</div>
              </div>
              <div class="array-form-box" @mousedown.stop="">
                <div class="form-item-list" v-for="(item, index) in formState.http_auth" :key="index">
                  <a-form-item :label="null">
                    <div class="flex-block-item" style="gap: 12px">
                      <a-input
                        style="width: 124px"
                        v-model:value="item.key"
                        placeholder="请输入"
                      ></a-input>
                      <a-form-item-rest>
                        <a-input
                          style="width: 124px"
                          v-model:value="item.value"
                          placeholder="请输入"
                        ></a-input>
                        <a-select v-model:value="item.add_to" style="width: 124px">
                          <a-select-option value="HEADERS">HEADERS</a-select-option>
                          <a-select-option value="PARAMS">PARAMS</a-select-option>
                          <a-select-option value="BODY">BODY</a-select-option>
                        </a-select>
                      </a-form-item-rest>
                      <div class="btn-hover-wrap" @click="handleDelAuthentication(index)">
                        <CloseCircleOutlined />
                      </div>
                    </div>
                  </a-form-item>
                </div>
                <a-button @click="handleOpenAuthenticationModal" :icon="h(PlusOutlined)" block type="dashed">添加参数</a-button>
              </div>
            </div>
          </div>

          <div class="gray-block mt16">
            <div class="gray-block-title output-title" @click="test">
              <span>输出 (输出字段提取)</span>
              <a-button type="primary" size="small" @click="handleOpenTestModal">自动提取输出参数</a-button>
            </div>
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
                          <a-select-option v-for="op in typOptions" :value="op.value" :key="op.value">{{
                            op.value
                          }}</a-select-option>
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
        </a-form>
      </div>
    </div>
    <ParseCurlModal 
      ref="parseCurlModalRef" 
      @parse="handleParseResult" 
    />
    <AddAuthenticationModal @ok="handleSaveAuthentication" ref="addAuthenticationModalRef" />
    <RunTest ref="runTestRef" :node_key="nodeId" @autoExtract="handleAutoExtract" />
  </NodeFormLayout>
  
</template>

<script setup>
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import { ref, reactive, watch, h, onMounted, inject } from 'vue'
import { CloseCircleOutlined, PlusOutlined, PlusCircleOutlined, CodeOutlined, QuestionCircleOutlined, CaretRightOutlined, } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import AtInput from '../../at-input/at-input.vue'
import SubKey from './subs-key.vue'
import ParseCurlModal from './parse-curl-modal.vue'
import { getUuid } from '@/utils/index'
import AddAuthenticationModal from './add-authentication-modal.vue'
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

const handleAutoSaveDraft = inject('handleAutoSaveDraft')

const variableOptions = ref([])
const atInputRefs = reactive({})
const runTestRef = ref(null)
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
  body_raw_tags: [],
  timeout: 30,
  output: [],
  http_auth:[],
})

const parseCurlModalRef = ref()

const showParseModal = () => {
  parseCurlModalRef.value.show()
}

const handleParseResult = (parsedData) => {
  // 将解析结果填充到表单
  let headers = [{
    key: '',
    value: ''
  }]
  let params = [{
    key: '',
    value: ''
  }]
  let body = []
  let body_raw = ''
  let isJson = false

  formState.method = parsedData.method
  formState.rawurl = parsedData.url
  // 处理value
  function handleValue(value) {
    if (typeof value == 'string') {
      return value
    }
    return JSON.stringify(value)
  }

  if(parsedData.header && typeof parsedData.header == 'object') {
    headers = Object.keys(parsedData.header).map((key) => {
      if(key.toLowerCase() == 'content-type' && parsedData.header[key].indexOf('application/json') !== -1) {
        isJson = true
      }

      return {
        key: key,
        value: handleValue(parsedData.header[key]),
        cu_key: getUuid(16)
      }
    })

    headers = headers.filter(item => item.key.toLowerCase() !== 'content-type')
  }
  
  if(parsedData.data && typeof parsedData.data == 'object') {
    body = Object.keys(parsedData.data).map((key) => {
      return {
        key: key,
        value: handleValue(parsedData.data[key]),
        cu_key: getUuid(16)
      }
    })

    body_raw = JSON.stringify(parsedData.data)
  }

  if(parsedData.params && typeof parsedData.params == 'object') {
    params = Object.keys(parsedData.params).map((key) => {
      return {
        key: key,
        value: handleValue(parsedData.params[key]),
        cu_key: getUuid(16)
      }
    })
  }

  if(isJson) {
    formState.type = 2
    formState.body_raw = body_raw
  }else {
    formState.type = 1
    formState.body = body
  }

  formState.headers = headers
  formState.params = params
}

function recursionData(data) {
  data.forEach((item) => {
    item.cu_key = getUuid(16)
    if (item.subs && item.subs.length) {
      recursionData(item.subs)
    } else {
      item.subs = []
    }
  })
  return data
}

function getFieldDefaultValue(item) {
  if(item.typ.includes('array')) {
    return []
  }

  return null
}

const handleOpenTestModal = async () => {
  let data = JSON.parse(JSON.stringify(formState))
  
  await handleAutoSaveDraft()
  
  let fields = []

  data.headers.forEach((item) => {
    if(item.tags && item.tags.length){
      fields = fields.concat(item.tags)
    }
  })

  data.params.forEach((item) => {
    if(item.tags && item.tags.length){
      fields = fields.concat(item.tags)
    }
  })

  if(formState.type == 1) {
    data.body.forEach((item) => {
    if(item.tags && item.tags.length){
        fields = fields.concat(item.tags)
      }
    })
  }else if(formState.type == 2) {
    if(data.body_raw_tags.length){
      fields = fields.concat(data.body_raw_tags)
    }
  }

  fields = fields.map(item => ({
    node_key: item.node_id,
    field: {
      key: item.original_value,
      label: item.label,
      typ: item.typ,
      Vals: getFieldDefaultValue(item),
      required: false,
      sys:false,
    }
  }))

  runTestRef.value?.open(fields)
}

const handleAutoExtract = (fields) => {
  formState.output = fields
  message.success('自动提取成功')
}

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
  const nodeModel = props.lf.getNodeModelById(props.nodeId)
  if (nodeModel) {
    let list = nodeModel.getAllParentVariable()
    list.forEach((item) => {
      item.tags = item.tags || []
    })

    variableOptions.value = list
  }
}


const update = () => {
  const data = JSON.stringify({
    curl: {
      ...formState
    }
  })

  emit('update-node', {
    ...props.node,
    ...formState,
    node_params: data
  })
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'
    let curl = JSON.parse(dataRaw).curl || {}

    curl = JSON.parse(JSON.stringify(curl))

    getVlaueVariableList()

    for (let key in curl) {
      if (key == 'headers' || key == 'params' || key == 'body') {
        if (curl[key] && curl[key].length > 0) {
          formState[key] = curl[key].map((item) => {
            return {
              ...item,
              cu_key: getUuid(16)
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
    formState.http_auth = curl.http_auth || []
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

const handleAddHeader = () => {
  formState.headers.push({
    key: '',
    value: '',
    cu_key: getUuid(16)
  })
}

const onDelHeader = (index) => {
  formState.headers.splice(index, 1)
}

const handleAddParams = () => {
  formState.params.push({
    key: '',
    value: '',
    cu_key: getUuid(16)
  })
}

const onDelParams = (index) => {
  formState.params.splice(index, 1)
}

const handleAddBody = () => {
  formState.body.push({
    key: '',
    value: '',
    cu_key: getUuid(16)
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
  if(!formState.output[index].subs || !Array.isArray(formState.output[index].subs)){
    formState.output[index].subs = []
  }

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


const addAuthenticationModalRef = ref(null)
const handleOpenAuthenticationModal = () => {
  addAuthenticationModalRef.value.show()
}

const handleSaveAuthentication = (list) => {
  let data = list.map(item => {
    return {
      key: item.auth_key,
      value: item.auth_value,
      add_to: item.auth_value_addto
    }
  })
  formState.http_auth = [...formState.http_auth, ...data]
}

const handleDelAuthentication = (index) => {
  formState.http_auth.splice(index, 1)
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
.run-test-btn {
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

.output-title{
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.export-curl-title{
  display: flex;
  align-items: center;
  justify-content: space-between;
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
      width: 214px;
    }
  }
  .flex-block-item .btn-hover-wrap {
    width: 24px;
    height: 24px;
  }
}
</style>
