<template>
  <div>
    <a-modal
      v-model:open="show"
      title="运行测试"
      :footer="null"
      :width="820"
      wrapClassName="no-padding-modal"
    >
      <div class="flex-content-box">
        <div class="test-model-box">
          <div class="top-title" @click="handleOpenRecognizeModal">运行参数</div>
          <a-form
            :model="formState"
            ref="formRef"
            layout="vertical"
            :wrapper-col="{ span: 24 }"
            autocomplete="off"
          >
            <a-form-item
              v-for="(item, index) in formState.test_params"
              :key="item.key"
              :name="['test_params', index, 'value']"
              :rules="[{ required: true, message: `请输入${item.field}` }]"
            >
              <template #label>
                <a-flex :gap="8"
                  >{{ item.field }}<a-tag style="margin: 0">{{ item.typ }}</a-tag>
                  <a v-if="item.typ.includes('array')" @click="handleOpenRecognizeModal(index)">
                    <ImportOutlined />
                    一键导入
                  </a>
                </a-flex>
              </template>
              <template v-if="item.typ != 'number' && !item.typ.includes('array')">
                <a-input placeholder="请输入" v-model:value="item.value" />
              </template>
              <template v-if="item.typ == 'number'">
                <a-input-number
                  style="width: 100%"
                  placeholder="请输入"
                  v-model:value="item.value"
                />
              </template>
              <template v-if="item.typ.includes('array')">
                <div class="input-list-box">
                  <div class="input-list-item" v-for="(input, i) in item.value" :key="i">
                    <a-form-item-rest
                      ><a-input placeholder="请输入" v-model:value="input.value"
                    /></a-form-item-rest>

                    <CloseCircleOutlined
                      v-if="item.value.length > 1"
                      @click="handleDelItem(item.value, i)"
                    />
                  </div>
                  <div class="add-btn-box">
                    <a-button @click="handleAddItem(item.value)" block type="dashed">添加</a-button>
                  </div>
                </div>
              </template>
            </a-form-item>
          </a-form>
          <div class="save-btn-box">
            <a-button
              :loading="loading"
              @click="handleSubmit"
              style="background-color: #00ad3a"
              type="primary"
              ><CaretRightOutlined />运行测试</a-button
            >
          </div>
        </div>
        <div class="preview-box">
          <template v-if="resultStr">
            <div class="preview-title">
              <div class="title-text">日志详情</div>
            </div>
            <div class="preview-content-block">
              <div class="title-block">运行日志<CopyOutlined @click="handleCopy" /></div>
              <div class="preview-code-box">
                <vue-json-pretty :data="resultStr" />
              </div>
            </div>
          </template>
          <template v-if="isErrorStatus">
            <div class="preview-title">
              <div class="title-text">日志详情</div>
            </div>
            <div class="preview-content-block">
              <div class="title-block">错误日志</div>
              <div class="error-box" style="margin-top: 12px;">
                 <a-alert :message="errorMsg" type="error" />
              </div>
            </div>
          </template>
        </div>
      </div>
    </a-modal>

    <a-modal v-model:open="recognizeOpen" :width="820" title="批量导入字段值" @ok="handleSaveValue">
      <a-textarea
        v-model:value.trim="jsonData"
        placeholder="请输入JSON字符串"
        style="min-height: 200px"
      />

      <div v-if="showJsonError" style="color: #f50">
        json格式错误 请输入正确的json字符串 列如
        [{"xxx":"11","id":48985},{"brand":"aaa","id":48986}]
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import {
  CaretRightOutlined,
  CloseCircleOutlined,
  CopyOutlined,
  ImportOutlined
} from '@ant-design/icons-vue'
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'
import { reactive, ref } from 'vue'
import { testCodeRun } from '@/api/robot/index'
import { message } from 'ant-design-vue'
import { copyText } from '@/utils/index'

const emit = defineEmits(['save'])

const props = defineProps({
  variableOptions: {
    type: Array,
    default: () => []
  }
})

const show = ref(false)

const loading = ref(false)

let is_in_loop_node = false // 该节点是否在循环节点中

const formState = reactive({
  language: 'javaScript',
  main_func: '',
  test_params: []
})

function matchAngleBrackets(str) {
  if (typeof str !== 'string') {
    return null;
  }
  
  // 匹配第一个 < > 内的内容
  const regex = /<([^<>]+)>/;
  const match = str.match(regex);
  
  return match ? match[1] : null;
}

function foundFieldTyp(value) {
  // 匹配字段的类型  如果节点在循环里面 则需要将数组类型 转为非数组类型
  if (value && value.length == 2) {
    let filterItem = props.variableOptions.filter((item) => item.value == value[0])
    if (filterItem.length > 0) {
      let type = filterItem[0].children.filter((item) => item.value == value[1])[0]?.typ || 'string'
      if(type.includes('array') && is_in_loop_node){
        return matchAngleBrackets(type)
      }
      return type
    }
  }

  return 'string'
}
let local_test_params = []
const handleOpenTestModal = async (data) => {
  resultStr.value = ''
  formState.main_func = data.main_func
  let localFormData = localStorage.getItem('code_run_node_test_modal_data') || '{}'
  local_test_params = JSON.parse(localFormData).test_params || []

  let params = data.params.filter((item) => item.field && item.variable && item.variable.length > 0)
  formState.test_params = params.map((item) => {
    let typ = foundFieldTyp(item.variable)
    return {
      key: Math.random() * 1000,
      field: item.field,
      typ,
      value: setGlobalDefaultVal(typ, item)
    }
  })

  // console.log(json)

  show.value = true
}

function setGlobalDefaultVal(typ, item) {
  // 设置默认值  先去缓存里找
  let filetrLocalData = local_test_params.filter((it) => it.field == item.field)
  if (!typ.includes('array')) {
    if (filetrLocalData.length && typeof filetrLocalData[0].value == 'string') {
      return filetrLocalData[0].value
    }
    return ''
  }

  if (filetrLocalData.length && typeof filetrLocalData[0].value == 'object') {
    return filetrLocalData[0].value
  }

  return [
    {
      value: '',
      key: Math.random() * 10000
    }
  ]
}

const handleDelItem = (item, index) => {
  item.splice(index, 1)
}
const handleAddItem = (item) => {
  item.push({
    value: '',
    key: Math.random() * 10000
  })
}

const formRef = ref(null)

const resultStr = ref('')

const isErrorStatus = ref(false)
const errorMsg = ref('')
const handleSubmit = () => {
  formRef.value.validate().then(() => {
    let postData = { ...formState }
    let test_params_data = JSON.parse(JSON.stringify(formState.test_params))
    let params = {}

    test_params_data.forEach((item) => {
      params[item.field] = item.value
      if (item.typ.includes('array')) {
        params[item.field] = item.value.map((it) => {
          it.value = it.value.trim()
          if (isJson(it.value)) {
            return JSON.parse(it.value)
          } else {
            return it.value
          }
        })
        params[item.field] = params[item.field].filter((it) => it)
      }
      if(item.typ === 'object'){
        if(isJson(item.value)){
          params[item.field] = JSON.parse(item.value)
        }
      }
    })

    postData.params = JSON.stringify(params)
    delete postData.test_params

    localStorage.setItem('code_run_node_test_modal_data', JSON.stringify(formState))

    loading.value = true
    isErrorStatus.value = false
    testCodeRun({
      ...postData
    })
      .then((res) => {
        message.success('测试结果生成完成')
        resultStr.value = JSON.parse(res.data)
      })
      .catch((res) => {
        // resultStr.value = JSON.parse(res.data)
        isErrorStatus.value = true
        errorMsg.value = res.msg
      })
      .finally(() => {
        loading.value = false
      })
  })
}

function isJson(str) {
  // 检查是否为字符串类型
  if (typeof str !== 'string') {
    return false
  }

  // 去除首尾空格
  const trimmedStr = str.trim()

  // 检查是否以 { 开头、} 结尾 或者 [ 开头、] 结尾
  const isValidFormat =
    (trimmedStr.startsWith('{') && trimmedStr.endsWith('}')) ||
    (trimmedStr.startsWith('[') && trimmedStr.endsWith(']'))

  if (!isValidFormat) {
    return false
  }

  try {
    // 尝试解析 JSON
    const parsed = JSON.parse(trimmedStr)
    // 确保解析结果是对象或数组
    return typeof parsed === 'object' && parsed !== null
  } catch (e) {
    // 如果解析失败，则不是有效 JSON
    return false
  }
}

const handleCopy = () => {
  copyText(JSON.stringify(resultStr.value))
  message.success('复制成功')
}

const open = (data, isInLoopNode) => {
  isErrorStatus.value = false
  is_in_loop_node = isInLoopNode
  handleOpenTestModal(data)
}

const recognizeOpen = ref(false)
const showJsonError = ref(false)
const editIndex = ref(null)
const jsonData = ref('')

const handleOpenRecognizeModal = (index) => {
  editIndex.value = index
  showJsonError.value = false
  recognizeOpen.value = true
}

const handleSaveValue = () => {
  if (isJson(jsonData.value)) {
    let lists = JSON.parse(jsonData.value).map((it) => {
      return {
        value: typeof it == 'string' ? it : JSON.stringify(it)
      }
    })
    let data = formState.test_params[editIndex.value]
    let newData = {
      ...data,
      value: [...data.value.filter(it => it.value != ''), ...lists]
    }
    formState.test_params.splice(editIndex.value, 1, newData)
    recognizeOpen.value = false
  } else {
    showJsonError.value = true
  }
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.flex-content-box {
  display: flex;
  max-height: 600px;
  overflow: hidden;
}
.test-model-box {
  flex: 1;
  margin: 24px 0 0 0;
  padding-right: 16px;
  overflow-y: auto;
  .top-title {
    font-weight: 600;
    margin-bottom: 16px;
  }
  .save-btn-box {
    margin: 32px 0;
    margin-top: 50px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
.tooltip-content {
  white-space: pre-wrap;
}
.loading-box {
  height: 100px;
  justify-content: center;
}
.result-list-box {
  margin: 24px 0;
  width: 100%;
  border: 1px solid #ebebeb;
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  padding: 8px;
  .list-item-block {
    display: flex;
    align-items: center;
    overflow: hidden;
    gap: 8px;
    padding: 8px;
    color: #333;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    .right-active-icon {
      margin-left: auto;
      color: #2475fc;
      opacity: 0;
    }
    &:hover {
      background: #f2f4f7;
      .right-active-icon {
        opacity: 1;
      }
    }
    &.active {
      color: #2475fc;
      background: #e6efff;
      .right-active-icon {
        opacity: 0;
      }
    }
    .status-block {
      font-size: 20px;
    }
    .icon-name-box {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      font-weight: 600;
      img {
        width: 24px;
        height: 24px;
      }
    }
    .time-tag {
      width: fit-content;
      border-radius: 4px;
      height: 22px;
      background: #d2f1dc;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0 4px;
      font-size: 12px;
    }
    .out-put-box {
      flex: 1;
      margin-left: 24px;
      overflow: hidden;
      .out-text-box {
        background: #f2f2f2;
        border-radius: 6px;
        padding: 8px;
        width: 100%;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}

.input-list-box {
  display: flex;
  flex-direction: column;
  gap: 8px;
  .input-list-item {
    display: flex;
    gap: 8px;
  }
}

.preview-box {
  flex: 1;
  border-left: 1px solid #d9d9d9;
  padding: 16px;
  overflow-y: auto;
  .preview-title {
    display: flex;
    align-items: center;
    gap: 8px;
    .title-text {
      font-size: 15px;
      font-weight: 600;
    }
    .icon-name-box {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      margin-left: 12px;
      img {
        width: 16px;
        height: 16px;
      }
    }
    .time-tag {
      width: fit-content;
      border-radius: 4px;
      height: 22px;
      background: #d2f1dc;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0 4px;
      font-size: 12px;
    }
  }
  .preview-content-block {
    margin-top: 16px;
    .title-block {
      font-size: 15px;
      color: #262626;
      display: flex;
      align-items: center;
      gap: 4px;
      .anticon-copy {
        cursor: pointer;
        &:hover {
          color: #2475fc;
        }
      }
    }
    .preview-code-box {
      margin-top: 16px;
      padding: 8px;
      border-radius: 8px;
      border: 1px solid #d9d9d9;
    }
  }
}
</style>
