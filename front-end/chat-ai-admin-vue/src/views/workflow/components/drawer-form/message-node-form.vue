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
      <div class="title-block">消息内容</div>
      <div class="message-cards-box">
        <div
          class="message-card-item"
          v-for="(item, index) in formState.message_list"
          :key="item.key"
        >
          <div class="message-title">
            <div>消息内容{{ index + 1 }}</div>
            <div class="minbtn-hover-wrap">
              <CloseCircleOutlined @click="onDelMsgCard(index)" />
            </div>
          </div>
          <div class="message-body">
            <a-form-item
              label="延迟发送"
              :name="['message_list', index, 'delay']"
              :rules="{ required: true, validator: (rule, value) => checkedDelay(rule, value) }"
            >
              <a-input-number
                style="width: 130px"
                v-model:value="item.delay"
                :precision="1"
                :step="1.5"
                :max="600"
                :min="0"
                addon-after="秒"
              ></a-input-number>
              <div class="form-tip">可设置延迟0-600s发送，必须为0.5秒的倍数</div>
            </a-form-item>
            <a-form-item label="消息内容" class="mb8" required>
              <a-radio-group v-model:value="item.ratio_type">
                <a-radio value="text">文本</a-radio>
                <a-radio value="image">图片</a-radio>
              </a-radio-group>
            </a-form-item>
            <a-form-item
              v-if="item.ratio_type == 'text'"
              :colon="false"
              class="hide-label"
              :name="['message_list', index, 'content']"
              :rules="{ required: true, message: '请输入内容', trigger: 'change' }"
            >
              <template #label></template>
              <a-textarea
                v-model:value="item.content"
                placeholder="请输入"
                :maxLength="200"
                :auto-size="{ minRows: 2, maxRows: 5 }"
              />
            </a-form-item>
            <a-form-item
              v-if="item.ratio_type == 'text'"
              :colon="false"
              class="hide-label"
              :name="['message_list', index, 'question']"
              :rules="{
                required: true,
                trigger: 'change',
                validator: (rule, value) => checkedQuestion(rule, value),
              }"
            >
              <template #label></template>
              <div class="question-guide-box">
                <draggable
                  style="display: flex; flex-direction: column; gap: 8px"
                  handle=".drag-btn"
                  v-model="item.question"
                  item-key="key"
                >
                  <template #item="{ element, index: i }">
                    <div class="question-guide-item" :key="element.key">
                      <HolderOutlined class="drag-btn" />
                      <div class="input-box">
                        <a-form-item-rest>
                          <a-input
                            :maxLength="14"
                            placeholder="请输入问题引导"
                            v-model:value="element.value"
                          ></a-input>
                        </a-form-item-rest>
                      </div>
                      <div class="btn-hover-wrap" @click="onDelGuideItem(index, i)">
                        <CloseCircleOutlined />
                      </div>
                    </div>
                  </template>
                </draggable>

                <div class="btn-box">
                  <a-button
                    :disabled="item.question.length >= 6"
                    @click="addGuideItem(index)"
                    type="dashed"
                    block
                  >
                    <template #icon>
                      <PlusOutlined />
                    </template>
                    添加问题引导
                  </a-button>
                </div>
              </div>
            </a-form-item>
            <a-form-item
              v-if="item.ratio_type == 'image'"
              :colon="false"
              class="hide-label"
              :name="['message_list', index, 'image_url']"
              :rules="{ required: true, message: '请上传图片', trigger: 'change' }"
            >
              <template #label></template>
              <UploadImage v-model:value="item.image_url" />
            </a-form-item>
          </div>
        </div>
        <div>
          <a-button
            :disabled="formState.message_list.length >= 3"
            @click="addMessageItem(index)"
            type="dashed"
            block
          >
            <template #icon>
              <PlusOutlined />
            </template>
            添加消息内容
          </a-button>
        </div>
      </div>
      <div class="title-block">对话逻辑</div>
      <a-form-item
        label="访客开口后不再推送"
        name="open_after_not_push"
        required
        :labelCol="{ span: 7 }"
      >
        <a-switch
          style="margin-top: 5px"
          v-model:checked="formState.open_after_not_push"
          checked-children="开"
          un-checked-children="关"
        />
        <div class="form-tip">开启后，如果访客发送了消息，讲不再发送后面的消息，进入下一节点</div>
      </a-form-item>
      <a-form-item
        label="跳出等待时长"
        name="click_wait_time"
        :labelCol="{ span: 7 }"
        :rules="{ required: true, message: '请输入跳出等待时长' }"
      >
        <a-input-number
          style="width: 130px"
          v-model:value="formState.click_wait_time"
          :precision="0"
          :step="1"
          :max="10"
          :min="3"
          addon-after="分钟"
        ></a-input-number>
        <div class="form-tip">
          发送的消息中包含问题引导时，系统等待访客响应达到指定时长才会跳出当前节点,进入下一节点。最低3分钟,最高10分钟
        </div>
      </a-form-item>
      <a-form-item
        label="调用知识库回复"
        name="receive_call_library_reply"
        required
        :labelCol="{ span: 7 }"
      >
        <a-switch
          style="margin-top: 5px"
          v-model:checked="formState.receive_call_library_reply"
          :checkedValue="1"
          :unCheckedValue="0"
          checked-children="开"
          un-checked-children="关"
        />
        <div class="form-tip">开启后，如果处于消息节点，用户发了消息，先调用大模型回复用户，再继续发送下一条消息或者跳出到下一个节点。</div>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { ref, reactive, inject, watch } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  CloseCircleFilled,
  CloseCircleOutlined,
  LoadingOutlined,
  PlusOutlined,
  EditOutlined,
  HolderOutlined,
} from '@ant-design/icons-vue'
import draggable from 'vuedraggable'
import UploadImage from './components/upload-image.vue'
import { useRoute } from 'vue-router'
const query = useRoute().query

const { updateNodeItem, updateModifyNum } = inject('nodeInfo')

const props = defineProps({
  properties: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['ok'])
const formRef = ref()

const labelCol = {
  span: 4,
}
const wrapperCol = {
  span: 20,
}

const formState = reactive({
  node_name: '',
  open_after_not_push: false,
  message_list: [
    {
      delay: 0,
      ratio_type: 'text',
      content: '',
      image_url: '',
      key: Math.random() * 10000,
      question: [],
    },
  ],
  click_wait_time: 3,
  receive_call_library_reply: 0,
})
let updateNum = 0
watch(
  () => props.properties,
  (val) => {
    try {
      formState.node_name = val.node_name
      formState.open_after_not_push = val.open_after_not_push == 1
      formState.receive_call_library_reply = val.receive_call_library_reply * 1

      formState.click_wait_time =
        (val.click_wait_time / 60).toFixed() >= 3 ? (val.click_wait_time / 60).toFixed() : 3
      let message_list = val.message_list || []
      formState.message_list = []
      if (message_list.length == 0) {
        formState.message_list = [
          {
            delay: 0,
            ratio_type: 'text',
            content: '',
            image_url: '',
            key: Math.random() * 10000,
            question: [],
          },
        ]
      }
      message_list.forEach((item) => {
        let data = {
          delay: item.delay,
          ratio_type: 'text',
          content: '',
          image_url: '',
          key: Math.random() * 10000,
          question: [],
        }
        if (item.msg_type == 'image') {
          data.ratio_type = 'image'
          data.image_url = item.content
        } else {
          data.ratio_type = 'text'
          data.content = item.content
        }
        if (item.msg_type == 'menu') {
          if (item.question && item.question.length) {
            data.question = item.question.map((it) => {
              return {
                value: it.content,
                key: new Date().getTime(),
                node_key: it.node_key,
              }
            })
          }
        }
        formState.message_list.push(data)
        updateNum = 0
      })
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

const checkedDelay = (rule, value) => {
  if (value == null) {
    return Promise.reject('请输入延迟发送时间')
  }
  if (!Number.isInteger(value / 0.5)) {
    return Promise.reject('必须为0.5秒的倍数')
  }
  return Promise.resolve()
}

const checkedQuestion = (rule, value) => {
  if (value.length > 0) {
    if (value.map((item) => item.value).join('') == '') {
      return Promise.reject('请输入至少一个问题引导选项')
    }
  }
  return Promise.resolve()
}

const addGuideItem = (index) => {
  formState.message_list[index].question.push({
    value: '',
    key: new Date().getTime(),
    node_key: '',
  })
}

const addMessageItem = () => {
  formState.message_list.push({
    delay: 0,
    ratio_type: 'text',
    content: '',
    image_url: '',
    key: Math.random() * 10000,
    question: [],
  })
}

const onDelMsgCard = (index) => {
  formState.message_list.splice(index, 1)
}

const onDelGuideItem = (index, i) => {
  formState.message_list[index].question.splice(i, 1)
}

const saveForm = () => {
  // 保存更新表单数据
  let updateInfo = {
    node_name: formState.node_name,
    open_after_not_push: formState.open_after_not_push ? 1 : 0,
    receive_call_library_reply: formState.receive_call_library_reply,
    message_list: [],
    click_wait_time: formState.click_wait_time * 60,
  }
  formState.message_list.forEach((item) => {
    let data = {
      delay: item.delay,
      msg_type: item.ratio_type == 'image' ? 'image' : item.question.length ? 'menu' : 'text',
      content: item.ratio_type == 'image' ? item.image_url : item.content,
      question: [],
    }
    data.question = item.question.map((it) => {
      return {
        content: it.value,
        node_key: it.node_key,
      }
    })
    updateInfo.message_list.push(data)
  })
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
.drag-btn{
  cursor: pointer;
  color: #595959;
}
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
      padding-right: 12px;
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
</style>
