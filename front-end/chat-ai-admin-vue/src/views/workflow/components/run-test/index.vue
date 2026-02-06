<template>
  <div>
    <!-- <a-button @click="handleOpenTestModal" style="background-color: #00ad3a" type="primary"
      ><CaretRightOutlined />{{ t('btn_run_test') }}</a-button
    > -->
    <a-modal v-model:open="show" :footer="null" :width="820" wrapClassName="no-padding-modal">
      <template #title>
        <div class="modal-title-block">
          {{ t('btn_run_test') }}{{ formState.dialogue_id }}
          <div class="run-detail" v-if="resultList.length">
            <span>{{ t('label_total_time') }}：{{ formatTime(use_mills) }}</span>
            <span>{{ t('label_token_usage') }}：{{ use_token }} {{ t('label_tokens') }}</span>
          </div>
        </div>
      </template>
      <div class="flex-content-box">
        <div class="test-model-box">
          <div class="top-title">{{ t('title_start_node_params') }}</div>
          <a-form
            :model="formState"
            ref="formRef"
            layout="vertical"
            :wrapper-col="{ span: 24 }"
            autocomplete="off"
          >
            <div v-show="!isShowQuestionForm">
              <a-form-item
                :name="['globalState', item.key]"
                v-for="item in diy_global"
                :key="item.key"
                :rules="[{ required: item.required, message: t('msg_input_key', { key: item.key }) }]"
              >
                <template #label>
                  <a-flex :gap="4"
                    >{{ item.key }} <a-tag style="margin: 0">{{ item.typ }}</a-tag>
                    <a-tooltip :title="item.desc">
                      <div class="option-desc">{{ item.desc }}</div>
                    </a-tooltip>
                  </a-flex>
                </template>
                <template v-if="item.typ == 'string'">
                  <a-input
                    :placeholder="getDefaultPlaceholder(item)"
                    v-model:value="formState.globalState[item.key]"
                  />
                </template>
                <template v-if="item.typ == 'number'">
                  <a-input-number
                    style="width: 100%"
                    :placeholder="getDefaultPlaceholder(item)"
                    v-model:value="formState.globalState[item.key]"
                  />
                </template>
                <template v-if="item.typ.includes('array')">
                  <div class="input-list-box">
                    <div
                      class="input-list-item"
                      v-for="(input, i) in formState.globalState[item.key]"
                      :key="i"
                    >
                      <a-form-item-rest
                        ><a-input
                          :placeholder="getDefaultPlaceholder(item)"
                          v-model:value="input.value"
                      /></a-form-item-rest>

                      <CloseCircleOutlined
                        v-if="formState.globalState[item.key].length > 1"
                        @click="handleDelItem(item.key, i)"
                      />
                    </div>
                    <div class="add-btn-box">
                      <a-button @click="handleAddItem(item.key)" block type="dashed">{{ t('btn_add') }}</a-button>
                    </div>
                  </div>
                </template>
              </a-form-item>
            </div>

            <template v-if="isShowQuestionForm">
              <a-form-item>
                <template #label>
                  <a-flex :gap="4">question <a-tag style="margin: 0">string</a-tag> </a-flex>
                </template>
                <a-input :placeholder="t('msg_input_key', { key: 'question' })" v-model:value="formState.question" />
              </a-form-item>
              <a-form-item v-if="question_multiple_switch">
                <template #label>
                  <a-flex :gap="4"
                    >question_multiple <a-tag style="margin: 0">string</a-tag>
                  </a-flex>
                </template>
                <div class="input-list-box">
                  <div
                    class="input-list-item"
                    v-for="(input, i) in formState.question_multiple"
                    :key="i"
                  >
                    <a-form-item-rest
                      ><a-input :placeholder="t('ph_input')" v-model:value="input.value"
                    /></a-form-item-rest>

                    <CloseCircleOutlined
                      v-if="formState.question_multiple.length > 1"
                      @click="handleDelQuestionItem(i)"
                    />
                  </div>
                  <div class="add-btn-box">
                    <a-button @click="handleAddQuetionItem()" block type="dashed">{{ t('btn_add') }}</a-button>
                  </div>
                </div>
              </a-form-item>
            </template>
          </a-form>

          <div class="result-list-box loading-box" v-if="loading">
            <a-spin v-if="loading" :tip="t('msg_generating_test_result')" />
          </div>

          <div class="result-list-box" v-if="resultList.length > 0">
            <div
              class="list-item-block"
              :class="{ active: currentNodeKey == item.node_key }"
              v-for="(item, index) in resultList"
              @click="handleChangeNodeKey(item)"
              :key="index"
            >
              <div class="status-block">
                <CheckCircleFilled v-if="item.is_success" style="color: #138b1b" />
                <CloseCircleFilled v-else style="color: #d81e06" />
              </div>
              <div class="icon-name-box">
                <img :src="item.node_icon" alt="" />
                <div class="node-name">{{ item.node_name }}</div>
              </div>
              <div class="time-tag" v-if="item.is_success">{{ item.use_time }}ms</div>
              <div class="right-active-icon"><RightCircleOutlined /></div>
              <!-- <div class="out-put-box" v-if="item.is_success">
                <a-tooltip>
                  <template #title>{{ item.output }}</template>
                  <div class="out-text-box">{{ item.output }}</div>
                </a-tooltip>
              </div> -->
            </div>
          </div>

          <div class="save-btn-box">
            <a-button
              :loading="loading"
              @click="handleSubmit"
              style="background-color: #00ad3a"
              type="primary"
              ><CaretRightOutlined />{{ isShowQuestionForm ? t('btn_continue_test') : t('btn_run_test') }} </a-button
            >
          </div>
        </div>
        <div class="preview-box">
          <template v-if="cuttentItem">
            <div class="preview-title">
              <div class="title-text">{{ t('title_log_details') }}</div>
              <div class="icon-name-box">
                <img :src="cuttentItem.node_icon" alt="" />
                <div class="node-name">{{ cuttentItem.node_name }}</div>
              </div>
              <div class="time-tag" v-if="cuttentItem.is_success">{{ cuttentItem.use_time }}ms</div>
            </div>
            <div class="preview-content-block" v-if="currentImageList.length > 0">
              <div class="title-block">{{ t('title_generated_image_log') }}</div>
              <div class="preview-img-box">
                <ImageLogs :currentImageList="currentImageList" />
              </div>
            </div>
            <div class="preview-content-block">
              <div class="title-block">{{ t('label_input') }}<CopyOutlined @click="handleCopy('input')" /></div>
              <div class="preview-code-box">
                <vue-json-pretty :data="cuttentItem.input" />
              </div>
            </div>
            <div class="preview-content-block">
              <div class="title-block">{{ t('label_output') }}<CopyOutlined @click="handleCopy('node_output')" /></div>
              <div class="preview-code-box">
                <vue-json-pretty :data="cuttentItem.node_output" />
              </div>
            </div>
            <div class="preview-content-block">
              <div class="title-block">{{ t('label_run_log') }}<CopyOutlined @click="handleCopy('output')" /></div>
              <div class="preview-code-box">
                <vue-json-pretty :data="cuttentItem.output" />
              </div>
            </div>
          </template>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { isJsonString } from '@/utils/index'
import {
  CaretRightOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  CloseCircleOutlined,
  RightCircleOutlined,
  CopyOutlined
} from '@ant-design/icons-vue'
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'
import { reactive, ref, computed, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { callWorkFlow } from '@/api/robot/index'
import { getImageUrl, formatTime } from '../util'
import { message } from 'ant-design-vue'
import { copyText } from '@/utils/index'
import ImageLogs from '@/views/workflow/components/image-logs/index.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.run-test.index')
const DEFAULT_SUGGEST_MAP = computed(() => ({
  openid: '38fjp8344ru43',
  question: t('placeholder_hello'),
  question_multiple: `{"type":"text", "text":"${t('placeholder_what_is_this')}"}`,
  conversationid: '12345'
}))
const props = defineProps({
  lf: {
    type: Object
  },
  start_node_params: {
    default: () => {},
    type: Object
  },
  isLockedByOther: { type: Boolean, default: false }
})

const isShowQuestionForm = ref(false)

const diy_global = computed(() => {
  return props.start_node_params.diy_global || []
})

const question_multiple_switch = computed(() => {
  let trigger_list = props.start_node_params.trigger_list || []
  let result = false
  trigger_list.forEach((item) => {
    if (item.trigger_type == 1) {
      result = item.chat_config?.question_multiple_switch
    }
  })
  return result
})

// const golbalTips = `自定义全局变量（json格式）
// 示例：
//   {
//     "str": "字符串",
//     "num": 1,
//     "arr": [
//       "a",
//       "b"
//     ]
//   }`
const emit = defineEmits(['save', 'getGlobal'])
const query = useRoute().query

const show = ref(false)
const currentNodeKey = ref('')
const resultList = ref([])

const cuttentItem = computed(() => {
  if (!currentNodeKey.value) {
    return null
  }
  return resultList.value.filter((item) => item.node_key == currentNodeKey.value)[0]
})

const currentImageList = computed(() => {
  let list = []
  if (cuttentItem.value && cuttentItem.value.node_type == 33) {
    let output = cuttentItem.value.output
    for (let key in output) {
      if (key.includes('picture_url_')) {
        list.push(output[key])
      }
    }
  }
  return list
})

const loading = ref(false)

const formState = reactive({
  is_draft: true,
  robot_key: query.robot_key,
  global: '',
  globalState: {},
  question: '',
  question_multiple: [
    {
      value: ''
    }
  ],
  dialogue_id: '',
  session_id: ''
})

const handleOpenTestModal = () => {
  isShowQuestionForm.value = false
  getQuestionMultipleSwitchStatus()
  let localData = localStorage.getItem('workflow_run_test_data') || '{}'
  localData = JSON.parse(localData)

  // formState.question = localData.question || ''
  // formState.openid = localData.openid || ''
  formState.global = localData.global || ''
  formState.question = ''
  formState.question_multiple = [
    {
      value: ''
    }
  ]
  formState.dialogue_id = ''
  formState.session_id = ''

  resultList.value = []
  currentNodeKey.value = ''
  emit('getGlobal')
  nextTick(() => {
    diy_global.value.forEach((item) => {
      formState.globalState[item.key] = setGlobalDefaultVal(item)
    })
    try {
      let global = localData.global ? JSON.parse(localData.global) : {}
      for (let key in formState.globalState) {
        if (global[key]) {
          if (Array.isArray(global[key])) {
            formState.globalState[key] = global[key].map((item) => {
              // 处理数组元素，将对象转换为字符串
              let value = typeof item === 'object' ? JSON.stringify(item) : item
              return {
                value: value,
                key: Math.random() * 1000
              }
            })
          } else {
            formState.globalState[key] = global[key]
          }
        }
      }
    } catch (error) {
      console.log(error)
    }
  })
  show.value = true
}

function setGlobalDefaultVal(item) {
  if (item.typ == 'string' || item.typ == 'number') {
    return ''
  }
  return [
    {
      value: '',
      key: Math.random() * 10000
    }
  ]
}

function getGlobalDefaultVal() {
  let result = {}
  diy_global.value.forEach((item) => {
    if (item.typ == 'string') {
      const v = (formState.globalState[item.key] || '').trim()
      result[item.key] = v ? v : (DEFAULT_SUGGEST_MAP.value[item.key] ?? '')
    }
    if (item.typ == 'number') {
      const raw = formState.globalState[item.key]
      if (raw !== '' && raw != null) {
        result[item.key] = +raw
      } else {
        const def = DEFAULT_SUGGEST_MAP.value[item.key]
        result[item.key] = def != null && def !== '' && !isNaN(Number(def)) ? Number(def) : ''
      }
    }

    if (item.typ.includes('array')) {
      let list = formState.globalState[item.key]
        .map((it) => {
          it.value = typeof it.value == 'string' ? it.value.trim() : it.value

          if (isJsonString(it.value)) {
            return JSON.parse(it.value)
          } else {
            return it.value
          }
        })
        .filter((it) => it)

      if ((!Array.isArray(list) || list.length === 0) && DEFAULT_SUGGEST_MAP.value[item.key]) {
        const def = DEFAULT_SUGGEST_MAP.value[item.key]
        if (isJsonString(def)) {
          list = [JSON.parse(def)]
        } else {
          list = [def]
        }
      }

      result[item.key] = list
    }

    if (item.typ === 'object') {
      if (isJsonString(item.value)) {
        result[item.key] = JSON.parse(item.value)
      }
    }
  })
  return JSON.stringify(result)
}

function getQuestionMultipleSwitchStatus() {
  const graphData = props.lf.getGraphData()
  const sessionTriggerNode = graphData.nodes.find((node) => node.type === 'session-trigger-node')

  if (sessionTriggerNode) {
    const node_params = JSON.parse(sessionTriggerNode.properties.node_params || '{}')
    const data = node_params.trigger || {}

    return data.chat_config ? data.chat_config.question_multiple_switch : false
  }
}

const handleDelItem = (key, index) => {
  formState.globalState[key].splice(index, 1)
}
const handleAddItem = (key) => {
  formState.globalState[key].push({
    value: '',
    key: Math.random() * 10000
  })
}

const handleDelQuestionItem = (index) => {
  formState.question_multiple.splice(index, 1)
}
const handleAddQuetionItem = (key) => {
  formState.question_multiple.push({
    value: ''
  })
}

const formRef = ref(null)

const use_token = ref(0)
const use_mills = ref(0)

const handleSubmit = () => {
  formRef.value.validate().then(() => {
    let postData = { ...formState }
    const question_multiple_switch = getQuestionMultipleSwitchStatus()

    postData.global = getGlobalDefaultVal()

    delete postData.globalState
    let question_multiple = postData.question_multiple.filter(item => item.value)
    postData.question_multiple = JSON.stringify(question_multiple.map(item => item.value))

    loading.value = true
    resultList.value = []

    const overrides = buildStorageOverrides()
    if (Object.keys(overrides).length) {
      localStorage.setItem(
        'workflow_run_test_data',
        JSON.stringify({ global: JSON.stringify(overrides) })
      )
    }

    callWorkFlow({
      ...postData,
      question_multiple_switch
    })
      .then((res) => {
        formState.dialogue_id = res.data.dialog_id
        formState.session_id = res.data.session_id
        let node_logs = res.data.node_logs || []
        use_token.value = res.data.use_token
        use_mills.value = res.data.use_mills
        formateData(node_logs)
      })
      .catch((res) => {
        resultList.value = []
        isShowQuestionForm.value = false

        let node_logs = res.data.node_logs || []
        if (node_logs && node_logs.length) {
          formateData(node_logs)
        }
      })
      .finally(() => {
        loading.value = false
      })
  })
}

const formateData = (data) => {
  let lastItem = data[data.length - 1]
  if(lastItem.node_type == 43){
    // 问答节点
    isShowQuestionForm.value = true
    message.success(t('msg_continue_qa_params'))
  }else{
    message.success(t('msg_test_result_generated'))
    formState.question = ''
    formState.question_multiple = [{value: ''}]
    isShowQuestionForm.value = false
  }
  resultList.value = data.map((item) => {
    let nodeIcon = getImageUrl(item.node_type)
    // 如果node_type == 45 就是HTTP工具，node_icon用接口返回的动态数据
    if (item.node_type == 45) {
      nodeIcon = item.node_icon || getImageUrl(item.node_type)
    }
    return {
      ...item,
      is_success: item.error_msg === '<nil>',
      node_icon: nodeIcon
    }
  })
  currentNodeKey.value = resultList.value[0]?.node_key
}

const handleChangeNodeKey = (item) => {
  currentNodeKey.value = item.node_key
}

const handleCopy = (key) => {
  copyText(JSON.stringify(cuttentItem.value[key]))
  message.success(t('msg_copy_success'))
}

const open = () => {
  handleOpenTestModal()
}

function getDefaultPlaceholder(item) {
  const def = DEFAULT_SUGGEST_MAP.value[item.key]
  return def ? `${def}` : t('ph_input')
}

function parseJSONMaybe(str) {
  try {
    return JSON.parse(str)
  } catch {
    return null
  }
}

function deepEqual(a, b) {
  return JSON.stringify(a) === JSON.stringify(b)
}

function buildStorageOverrides() {
  const overrides = {}
  diy_global.value.forEach((item) => {
    const key = item.key
    if (item.typ === 'string') {
      const v = String(formState.globalState[key] || '').trim()
      if (!v) return
      const def = DEFAULT_SUGGEST_MAP.value[key]
      if (def != null && v === String(def)) return
      overrides[key] = v
      return
    }
    if (item.typ === 'number') {
      const raw = formState.globalState[key]
      if (raw === '' || raw == null) return
      const val = +raw
      const def = DEFAULT_SUGGEST_MAP.value[key]
      if (def != null && !isNaN(Number(def)) && val === Number(def)) return
      overrides[key] = val
      return
    }
    if (item.typ.includes('array')) {
      const list = Array.isArray(formState.globalState[key]) ? formState.globalState[key] : []
      const typed = list
        .map((it) => String(it?.value || '').trim())
        .filter((s) => s)
        .map((s) => {
          const o = parseJSONMaybe(s)
          return o != null ? o : s
        })
      if (!typed.length) return
      let defList = []
      const def = DEFAULT_SUGGEST_MAP.value[key]
      if (def != null) {
        const o = parseJSONMaybe(String(def))
        defList = o != null ? [o] : [String(def)]
      }
      if (defList.length && deepEqual(typed, defList)) return
      overrides[key] = typed
      return
    }
  })
  return overrides
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.flex-content-box {
  display: flex;
  max-height: 70vh;
  overflow: hidden;
}
.test-model-box {
  flex: 1;
  margin: 24px 0 0 0;
  overflow-y: auto;
  padding-right: 16px;
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
  overflow-y: auto;
  border-left: 1px solid #d9d9d9;
  padding: 16px;
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
      width: fit-content;
      min-width: 100%;
      margin-top: 16px;
      padding: 8px;
      border-radius: 8px;
      border: 1px solid #d9d9d9;

      &::v-deep(.vjs-tree) {
        width: fit-content;
      }

      &::v-deep(.vjs-tree-node) {
        width: calc(100% + 16px);
        padding-right: 16px;
      }
    }
  }
}

.modal-title-block {
  display: flex;
  align-items: center;
  gap: 12px;
  .run-detail {
    display: flex;
    align-items: center;
    gap: 16px;
    background: #bffbd7;
    padding: 4px 16px;
    font-size: 14px;
    color: #595959;
    border-radius: 8px;
  }
}

.option-desc {
  max-width: 90px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
}

</style>
