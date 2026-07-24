<template>
  <div class="run-test-panel">
    <div class="test-model-box">
      <div class="test-model-scroll customize-scroll-style">
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
                <a-flex :gap="4">{{ questionLabel }} <a-tag style="margin: 0">string</a-tag> </a-flex>
              </template>
              <a-input :placeholder="t('msg_input_key', { key: 'question' })" v-model:value="formState.question" />
            </a-form-item>
            <a-form-item v-if="questionMultipleSwitch">
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

        <div class="loading-box" v-if="loading">
          <a-spin :tip="t('msg_generating_test_result')" />
        </div>
      </div>

      <div class="save-btn-box">
        <a-button
          :loading="loading"
          @click="handleSubmit"
          style="background-color: #00ad3a"
          type="primary"
          block
          ><CaretRightOutlined />{{ isShowQuestionForm ? t('btn_continue_test') : t('btn_run_test') }} </a-button
        >
      </div>
    </div>

    <NodeRunLogs
      v-model:currentNodeKey="currentNodeKey"
      :resultList="resultList"
      :loading="loading"
      :hasRunTested="hasRunTested"
    />
  </div>
</template>

<script setup>
import { isJsonString } from '@/utils/index'
import {
  CaretRightOutlined,
  CloseCircleOutlined
} from '@ant-design/icons-vue'
import { reactive, ref, computed, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { callWorkFlow } from '@/api/robot/index'
import { getImageUrl } from '../util'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import NodeRunLogs from './NodeRunLogs.vue'

const { t } = useI18n('views.workflow.components.run-test.index')
const DEFAULT_SUGGEST_MAP = computed(() => ({
  openid: '38fjp8344ru43',
  question: t('placeholder_hello'),
  question_multiple: `{"type":"text", "text":"${t('placeholder_what_is_this')}"}`,
  conversationid: '12345'
}))

const props = defineProps({
  start_node_params: {
    default: () => {},
    type: Object
  }
})

const emit = defineEmits(['stateChange'])
const query = useRoute().query

const isShowQuestionForm = ref(false)
const questionLabel = ref('question')
const diy_global = computed(() => props.start_node_params.diy_global || [])
const questionMultipleSwitch = computed(() => {
  let trigger_list = props.start_node_params.trigger_list || []
  let result = false
  trigger_list.forEach((item) => {
    if (item.trigger_type == 1) {
      result = item.chat_config?.question_multiple_switch
    }
  })
  return result
})

const currentNodeKey = ref('')
const resultList = ref([])
const hasRunTested = ref(false)
const loading = ref(false)
const formRef = ref(null)
const use_token = ref(0)
const use_mills = ref(0)

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

const syncState = () => {
  emit('stateChange', {
    hasRunTested: hasRunTested.value,
    hasRunResult: resultList.value.length > 0,
    use_token: use_token.value,
    use_mills: use_mills.value
  })
}

const open = () => {
  isShowQuestionForm.value = false
  questionLabel.value = 'question'
  let localData = localStorage.getItem('workflow_run_test_data') || '{}'
  localData = JSON.parse(localData)

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
  hasRunTested.value = false
  use_token.value = 0
  use_mills.value = 0
  syncState()
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
const handleAddQuetionItem = () => {
  formState.question_multiple.push({
    value: ''
  })
}

const handleSubmit = () => {
  formRef.value.validate().then(() => {
    let postData = { ...formState }

    postData.global = getGlobalDefaultVal()

    delete postData.globalState
    let question_multiple = postData.question_multiple.filter(item => item.value)
    postData.question_multiple = JSON.stringify(question_multiple.map(item => item.value))

    loading.value = true
    resultList.value = []
    hasRunTested.value = true
    syncState()

    const overrides = buildStorageOverrides()
    if (Object.keys(overrides).length) {
      localStorage.setItem(
        'workflow_run_test_data',
        JSON.stringify({ global: JSON.stringify(overrides) })
      )
    }

    callWorkFlow({
      ...postData,
      question_multiple_switch: questionMultipleSwitch.value
    })
      .then((res) => {
        formState.dialogue_id = res.data.dialog_id
        formState.session_id = res.data.session_id
        let node_logs = res.data.node_logs || []
        use_token.value = res.data.use_token
        use_mills.value = res.data.use_mills
        formatData(node_logs)
      })
      .catch((res) => {
        resultList.value = []
        isShowQuestionForm.value = false

        let node_logs = res.data.node_logs || []
        if (node_logs && node_logs.length) {
          formatData(node_logs)
        }
      })
      .finally(() => {
        loading.value = false
        syncState()
      })
  })
}

const formatData = (data) => {
  let lastItem = data[data.length - 1]
  if (lastItem?.node_type == 43) {
    questionLabel.value = getQuestionLabel(lastItem)
    isShowQuestionForm.value = true
    message.success(t('msg_continue_qa_params'))
  } else {
    message.success(t('msg_test_result_generated'))
    formState.question = ''
    formState.question_multiple = [{ value: '' }]
    questionLabel.value = 'question'
    isShowQuestionForm.value = false
  }
  resultList.value = data.map((item, index) => formatNodeLog(item, index))
  currentNodeKey.value = resultList.value[0]?.log_key
  syncState()
}

function getQuestionLabel(item) {
  const special = item?.output?.special || item?.node_output?.special || {}
  const replyContentList = parseReplyContentList(special.reply_content_list)
  const smartMenu = replyContentList.find((reply) => (reply.reply_type || reply.type) === 'smartMenu')
  return smartMenu?.smart_menu?.menu_description || special.llm_reply_content || 'question'
}

function parseReplyContentList(value) {
  if (Array.isArray(value)) {
    return value
  }
  if (typeof value !== 'string' || !value) {
    return []
  }
  const parsed = parseJSONMaybe(value)
  return Array.isArray(parsed) ? parsed : []
}

function formatNodeLog(item, index) {
  let nodeIcon = getImageUrl(item.node_type)
  if (item.node_type == 45) {
    nodeIcon = item.node_icon || getImageUrl(item.node_type)
  }
  const nodeKey = item.node_key
  return {
    ...item,
    log_key: item.log_key || `${nodeKey}_${item.start_time ?? index}`,
    is_success: item.error_msg === '<nil>',
    node_icon: nodeIcon
  }
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
    }
  })
  return overrides
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.run-test-panel {
  display: flex;
  height: 70vh;
  overflow: hidden;
  align-items: stretch;
}
.test-model-box {
  width: var(--chat-test-chat-width, 360px);
  min-width: var(--chat-test-chat-width, 360px);
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #f0f0f0;
  background: #fff;
  .test-model-scroll {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
  }
  .top-title {
    font-weight: 600;
    margin-bottom: 16px;
  }
  .save-btn-box {
    padding: 14px 20px 12px;
    border-top: 1px solid #f0f0f0;
    background: #fff;
    box-shadow: 0 -6px 20px rgba(15, 23, 42, 0.06);
  }
}
.loading-box {
  min-height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
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
.option-desc {
  max-width: 90px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
}
</style>
