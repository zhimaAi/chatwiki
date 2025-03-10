<style lang="less" scoped>
.prompt-log-alert {
  .close-btn {
    font-size: 16px;
    color: rgba(0, 0, 0, 0.45);
    cursor: pointer;
  }

  .prompt-log-content {
    .prompt-log-items {
      margin-bottom: 24px;
    }
    .prompt-log-label {
      margin-bottom: 4px;
      font-size: 16px;
      font-weight: 600;
      line-height: 24px;
      color: rgb(0, 0, 0);
    }

    .prompt-log-item {
      line-height: 22px;
      padding: 12px 16px;
      margin-bottom: 8px;
      font-size: 14px;
      color: #3a4559;
      background-color: #f2f4f7;
      border-radius: 4px;
      white-space: pre-wrap;

      &:last-child {
        margin-bottom: 0;
      }
    }
  }
}
</style>

<template>
  <a-drawer
    class="prompt-log-alert"
    v-model:open="show"
    title="Prompt 日志"
    placement="right"
    width="746px"
    :closable="false"
  >
    <template #extra>
      <span class="close-btn" @click="onClose"><CloseOutlined /></span>
    </template>

    <div class="prompt-log-content">
      <div class="prompt-log-items">
        <div class="prompt-log-label">
          <span>SYSTEM </span>
          <a-tooltip>
            <template #title>系统提示词和文档分段。</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>

        <div class="prompt-log-item">
          <p>{{ promptLog.prompt }}</p>
        </div>

        <div class="prompt-log-item" v-for="(item, index) in promptLog.library" :key="index">
          <p>{{ item }}</p>
        </div>
      </div>

      <div class="prompt-log-items">
        <div class="prompt-log-label">
          <span>上下文 </span>
          <a-tooltip>
            <template #title>传递的历史提问消息</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="prompt-log-item" v-for="(item, index) in promptLog.context_qa" :key="index">
          <p>Q：{{ item.question }}</p>
          <p>A：{{ item.answer }}</p>
        </div>
      </div>

      <div class="prompt-log-items">
        <div class="prompt-log-label">
          <span>USER </span>
          <a-tooltip>
            <template #title>本次用户的提问</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="prompt-log-item">
          <p>{{ promptLog.cur_question }}</p>
        </div>
      </div>

      <div class="prompt-log-items">
        <div class="prompt-log-label">
          <span>ASSISTANT </span>
          <a-tooltip>
            <template #title>语言模型输出的答案</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="prompt-log-item">
          <p>{{ promptLog.cur_answer }}</p>
        </div>
      </div>

      <div class="prompt-log-items" v-if="promptLog.recall_time > 0">
        <div class="prompt-log-label">
          <span>Recall time </span>
          <a-tooltip>
            <template #title>从知识库中检索到有效分段所需要的时间，包括对分段进行排序时间。</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="prompt-log-item">
          <p>{{ promptLog.recall_time }}s</p>
        </div>
      </div>

      <div class="prompt-log-items" v-if="promptLog.request_time > 0">
        <div class="prompt-log-label">
          <span>Request time </span>
          <a-tooltip>
            <template #title>从发送问题和上下文信息给大模型到大模型开始返回答案的时间。</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="prompt-log-item">
          <p>{{ promptLog.request_time }}s</p>
        </div>
      </div>

      <div class="prompt-log-items">
        <div class="prompt-log-label">
          <span>Error </span>
          <a-tooltip>
            <template #title>报错信息，调用聊天接口报错时会显示。</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="prompt-log-item">
          <p>{{ promptLog.error }}</p>
        </div>
      </div>
    </div>
  </a-drawer>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { CloseOutlined, QuestionCircleOutlined } from '@ant-design/icons-vue'

const show = ref(false)
const promptLog = reactive({
  prompt: '',
  library: [],
  context_qa: [],
  cur_question: '',
  cur_answer: '',
  error: ''
})

const onClose = () => {
  show.value = false
}

const reset = () => {
  promptLog.prompt = ''
  promptLog.library = []
  promptLog.context_qa = []
  promptLog.cur_question = ''
  promptLog.cur_answer = ''
  promptLog.error = ''
  promptLog.recall_time = ''
  promptLog.request_time = ''
}

const open = (msg) => {
  reset()
  promptLog.error = msg.error
  promptLog.recall_time = msg.recall_time ?  (msg.recall_time / 1000).toFixed(2) : '';
  promptLog.request_time = msg.request_time ? (msg.request_time / 1000).toFixed(2) : '';

  let items = msg.debug || []

  for (let i = 0; i < items.length; i++) {
    let item = items[i]

    if (item.type == 'library') {
      promptLog.library.push(item.content)
    } else if (item.type == 'context_qa') {
      promptLog.context_qa.push(item)
    } else {
      promptLog[item.type] = item.content
    }
  }

  show.value = true
}

defineExpose({
  open
})
</script>
