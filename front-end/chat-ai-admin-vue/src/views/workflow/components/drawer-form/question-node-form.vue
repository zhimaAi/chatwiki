<template>
  <div class="form-box">
    <a-form ref="formRef" :model="formState" :labelCol="labelCol" :wrapperCol="wrapperCol">
      <div class="title-block">基础信息</div>
      <a-form-item
        label="节点名称"
        name="node_name"
        :rules="[{ required: true, message: '请输入节点名称' }]"
      >
        <a-input
          v-model:value="formState.node_name"
          :maxLength="15"
          placeholder="请输入节点名称，最多15个字"
        ></a-input>
      </a-form-item>
      <a-form-item
        label="询问信息"
        name="field_key"
        v-if="node_sub_type == 32"
        :rules="[{ required: true, message: '请选择' }]"
      >
        <a-select v-model:value="formState.field_key" style="width: 100%" placeholder="请选择">
          <a-select-option v-for="item in filedLists" :value="item.field_key" :key="item.field_key">{{
            item.field_name
          }}</a-select-option>
        </a-select>
        <div class="form-tip">
          请选择添加的自动留资字段 <a @click="toAddFieldPage">去添加</a>
          <a @click="getFields(true)" class="ml8"><SyncOutlined />刷新</a>
        </div>
      </a-form-item>
      <a-form-item
        label="提问延时"
        name="delay_question"
        :rules="{ required: true, validator: (rule, value) => checkedDelay(rule, value) }"
      >
        <a-input-number
          style="width: 130px"
          v-model:value="formState.delay_question"
          :precision="1"
          :step="1.5"
          :max="600"
          :min="0"
          addon-after="秒"
        ></a-input-number>
        <div class="form-tip">进入当前节点后延时多久提问，最小为0，最大为600，只支持0.5的倍数</div>
      </a-form-item>
      <div class="title-block">提问参考话术</div>
      <div class="message-cards-box">
        <div
          class="message-card-item"
          v-for="(item, index) in formState.question_reference"
          :key="item.key"
        >
          <div class="message-title">
            <div>参考话术{{ index + 1 }}</div>
            <div class="minbtn-hover-wrap">
              <CloseCircleOutlined
                v-if="formState.question_reference.length > 1"
                @click="onDelQuestion(index, 'question_reference')"
              />
            </div>
          </div>
          <div class="message-body">
            <a-form-item
              :colon="false"
              class="hide-label"
              :labelCol="{ span: 0 }"
              :wrapperCol="{ span: 24 }"
              :name="['question_reference', index, 'value']"
              :rules="{ required: true, message: '请输入内容', trigger: 'change' }"
            >
              <template #label></template>
              <a-textarea
                v-model:value="item.value"
                placeholder="请输入"
                :maxLength="200"
                :auto-size="{ minRows: 2, maxRows: 5 }"
              />
            </a-form-item>
          </div>
        </div>
        <div>
          <a-button
            :disabled="formState.question_reference.length >= 3"
            @click="addMessageItem('question_reference')"
            type="dashed"
            block
          >
            <template #icon>
              <PlusOutlined />
            </template>
            添加参考话术
          </a-button>
        </div>
      </div>
      <div class="title-block mt24">问答匹配</div>
      <a-form-item label="问答不匹配时">
        <a-radio-group v-model:value="formState.not_match_mode">
          <a-radio :value="1">跳过</a-radio>
          <a-radio :value="2">反问</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item
        label="调用知识库回复"
        name="not_match_call_library_reply"
        required
      >
        <a-switch
          style="margin-top: 5px"
          v-model:checked="formState.not_match_call_library_reply"
          :checkedValue="1"
          :unCheckedValue="0"
          checked-children="开"
          un-checked-children="关"
          @change="onNotMatchCallLibraryReplyChange"
        />
        <div class="form-tip">开启后，跳过或者反问前，会调用知识库对用户消息进行回复</div>
      </a-form-item>
      <a-form-item
        label="反问次数"
        v-if="formState.not_match_mode == 2"
        name="trace_num"
        :rules="[{ required: true, message: '请输入反问次数' }]"
      >
        <a-input-number
          style="width: 130px"
          v-model:value="formState.trace_num"
          :precision="0"
          :step="1"
          :max="10"
          :min="1"
          addon-after="次"
        ></a-input-number>
        <div class="form-tip">最多反问10次，达到指定次数后还是不匹配，会自动进入下一流程</div>
      </a-form-item>
      <a-form-item
        label="反问延时"
        class="mb16"
        v-if="formState.not_match_mode == 2"
        name="trace_delay_max"
        :rules="[{ required: true, validator: (rule, value) => checkedTraceNum(rule, value) }]"
      >
        <a-flex align="center" :gap="4">
          <a-form-item-rest>
            <a-input-number
              style="width: 80px"
              v-model:value="formState.trace_delay_min"
              :precision="0"
              :step="1"
              :max="600"
              :min="0"
            ></a-input-number>
          </a-form-item-rest>

          <span>秒至</span>
          <a-input-number
            style="width: 80px"
            v-model:value="formState.trace_delay_max"
            :precision="0"
            :step="1"
            :max="600"
            :min="0"
          ></a-input-number>
          秒
        </a-flex>

        <div class="form-tip">
          每次反问延时发送时间，系统在指定时间区间内随机延时,最低0秒，最高600秒
        </div>
      </a-form-item>
      <template v-if="formState.not_match_mode == 2">
        <div class="message-cards-box" style="padding-left: 110px">
          <div
            class="message-card-item"
            v-for="(item, index) in formState.trace_reference"
            :key="item.key"
          >
            <div class="message-title">
              <div>反问参考话术{{ index + 1 }}</div>
              <div class="minbtn-hover-wrap">
                <CloseCircleOutlined
                  v-if="formState.trace_reference.length > 1"
                  @click="onDelQuestion(index, 'trace_reference')"
                />
              </div>
            </div>
            <div class="message-body">
              <a-form-item
                :colon="false"
                class="hide-label"
                :labelCol="{ span: 0 }"
                :wrapperCol="{ span: 24 }"
                :name="['trace_reference', index, 'value']"
                :rules="{ required: true, message: '请输入内容', trigger: 'change' }"
              >
                <template #label></template>
                <a-textarea
                  v-model:value="item.value"
                  placeholder="请输入"
                  :maxLength="200"
                  :auto-size="{ minRows: 2, maxRows: 5 }"
                />
              </a-form-item>
            </div>
          </div>
          <div>
            <a-button
              :disabled="formState.trace_reference.length >= 3"
              @click="addMessageItem('trace_reference')"
              type="dashed"
              block
            >
              <template #icon>
                <PlusOutlined />
              </template>
              添加反问参考话术
            </a-button>
          </div>
        </div>
      </template>

      <div class="title-block mt24">对话逻辑</div>
      <a-form-item class="mb8" label="沉默用户追问" name="silence_trace" required>
        <a-switch
          style="margin-top: 2px"
          v-model:checked="formState.silence_trace"
          checked-children="开"
          un-checked-children="关"
        />
      </a-form-item>
      <a-form-item
        :colon="false"
        class="hide-label"
        v-if="formState.silence_trace"
        name="silence_delay"
        :rules="[{ required: true, message: '请输入' }]"
      >
        <template #label></template>
        <div class="gray-block">
          <a-flex align="center" :gap="4">
            <span>提问后，访客超过</span>
            <a-input-number
              style="width: 80px"
              v-model:value="formState.silence_delay"
              :precision="0"
              :step="1"
              :max="300"
              :min="10"
            ></a-input-number>
            <span>秒未回复，自动追问</span>
          </a-flex>
        </div>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { ref, reactive, inject, watch, h, nextTick } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  CloseCircleFilled,
  CloseCircleOutlined,
  LoadingOutlined,
  PlusOutlined,
  EditOutlined,
  SyncOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons-vue'
import { getFieldsList } from '@/api/retention/index.js'
import { useRoute } from 'vue-router'
import { useRobotStore } from '@/stores/modules/robot.js'

const RobotStore = useRobotStore()

const query = useRoute().query

const { updateNodeItem, updateModifyNum } = inject('nodeInfo')

const node_sub_type = ref(31)
const props = defineProps({
  properties: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['ok'])
const formRef = ref()

const labelCol = {
  span: 6,
}
const wrapperCol = {
  span: 19,
}

const formState = reactive({
  node_name: '',
  delay_question: 0,
  field_key: void 0,
  question_reference: [
    {
      value: '',
      key: Math.random() * 10000,
    },
  ],
  trace_reference: [
    {
      value: '',
      key: Math.random() * 10000,
    },
  ],
  not_match_mode: 1,
  trace_num: 1,
  trace_delay_min: 1,
  trace_delay_max: 5,
  silence_trace: false,
  silence_delay: 0,
  not_match_call_library_reply: 0,
})
let updateNum = 0
watch(
  () => props.properties,
  (val) => {
    try {
      node_sub_type.value = val.node_sub_type
      formState.node_name = val.node_name
      formState.field_key = val.field_key ? val.field_key : void 0
      formState.delay_question = val.delay_question || 0
      formState.not_match_call_library_reply = val.not_match_call_library_reply * 1 || 0
      formState.question_reference = val.question_reference.map((item) => {
        return {
          value: item,
          key: Math.random() * 10000,
        }
      })

      formState.trace_reference = val.trace_reference.map((item) => {
        return {
          value: item,
          key: Math.random() * 10000,
        }
      })

      if (formState.question_reference.length == 0) {
        formState.question_reference = [
          {
            value: '',
            key: Math.random() * 10000,
          },
        ]
      }

      if (formState.trace_reference.length == 0) {
        formState.trace_reference = [
          {
            value: '',
            key: Math.random() * 10000,
          },
        ]
      }

      formState.not_match_mode = +val.not_match_mode || 2
      formState.trace_num = +val.trace_num || 3
      formState.trace_delay_min = +val.trace_delay_min || 1
      formState.trace_delay_max = +val.trace_delay_max || 5
      formState.silence_trace = val.silence_trace == 1
      formState.silence_delay = +val.silence_delay || 10
      if (val.node_sub_type == 32) {
        // 自定义字段
        getFields()
      }
      updateNum = 0
    } catch (error) {
      console.log(error)
    }
  },
  { immediate: true, deep: true },
)
watch(
  () => formState,
  () => {
    updateNum++
    updateModifyNum(updateNum)
  },
  {
    deep: true,
  },
)
const filedLists = ref([])
function getFields(showMessage) {
  getFieldsList({
    page: 1,
    limit: 9999,
    robot_id: query.robotId,
  }).then((res) => {
    if (!res.data) {
      return
    }
    filedLists.value = res.data.data.filter((item) => item.is_default == 0)
    showMessage && message.success('刷新成功')
  })
}

const checkedDelay = (rule, value) => {
  if (value == null) {
    return Promise.reject('请输入延迟发送时间')
  }
  if (!Number.isInteger(value / 0.5)) {
    return Promise.reject('必须为0.5秒的倍数')
  }
  return Promise.resolve()
}

const checkedTraceNum = (rule, value) => {
  if (value == null || formState.trace_delay_min == null) {
    return Promise.reject('请输入反问延时间')
  }
  if (+formState.trace_delay_min > +value) {
    return Promise.reject('反问延时间范围不正确')
  }
  return Promise.resolve()
}

const addMessageItem = (key) => {
  formState[key].push({
    value: '',
    key: Math.random() * 10000,
  })
}

const onDelQuestion = (index, key) => {
  formState[key].splice(index, 1)
}

const toAddFieldPage = () => {
  if (import.meta.env.DEV) {
    window.open(`/#/retention-information/management?robotId=${query.robotId}`)
  } else {
    window.open(`/kf/retainDataRobot/#/retention-information/management?robotId=${query.robotId}`)
  }
}

const onNotMatchCallLibraryReplyChange = (val) => {
  if(val == 0){
    return
  }


  if(!RobotStore.robotInfo.library_ids){
    formState.not_match_call_library_reply = 0
    Modal.confirm({
      title: '无法开启调用知识库回复',
      icon: h(CloseCircleFilled),
      content: '当前机器人还未关联知识库，请先关联后再试。',
      okText: '去关联',
      onOk() {
        let path = `/#/base-setting?robotId=${query.robotId}&name=${query.name}`;

        if (import.meta.env.DEV) {
          window.open(path)
        } else {
          window.open(`/kf/retainDataRobot`+path)
        }
      },
      onCancel() {
      },
    });
  }
}

const saveForm = () => {
  // 保存更新表单数据
  let updateInfo = {
    node_name: formState.node_name,
    field_key: formState.field_key,
    delay_question: formState.delay_question,
    question_reference: formState.question_reference.map((item) => item.value),
    trace_reference: formState.trace_reference.map((item) => item.value),
    not_match_mode: formState.not_match_mode,
    trace_num: formState.trace_num,
    trace_delay_min: formState.trace_delay_min,
    trace_delay_max: formState.trace_delay_max,
    silence_trace: formState.silence_trace ? 1 : 0,
    silence_delay: formState.silence_delay,
    not_match_call_library_reply: formState.not_match_call_library_reply
  }
  updateNodeItem({ ...updateInfo })
}

const onSave = () => {
  formRef.value
    .validate()
    .then(() => {
      saveForm()
    })
    .catch((err) => {
      console.log('error', err)
    })
}

defineExpose({
  onSave,
})
</script>

<style lang="less" scoped>
@import './common.less';

.message-cards-box {
  display: flex;
  flex-direction: column;
  gap: 16px;
  .message-card-item {
    background: #f2f4f7;
    border-radius: 6px;
    .message-title {
      display: flex;
      align-items: center;
      justify-content: space-between;
      height: 32px;
      padding: 0 16px;
      color: #000000;
      border-bottom: 1px solid #e4e6eb;
    }
    .message-body {
      padding: 16px;
      padding-bottom: 0;
      ::v-deep(.ant-form-item) {
        margin-bottom: 16px;
      }
      ::v-deep(.mb8.ant-form-item) {
        margin-bottom: 8px;
      }
    }
  }
  .question-guide-box {
    display: flex;
    flex-direction: column;
    gap: 8px;
    .question-guide-item {
      display: flex;
      align-items: center;
      gap: 8px;
      .input-box {
        flex: 1;
      }
    }
  }
}
.gray-block {
  background: #f2f4f7;
  padding: 12px 16px;
  border-radius: 6px;
}
</style>
