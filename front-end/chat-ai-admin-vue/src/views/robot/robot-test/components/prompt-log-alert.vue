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
    :title="t('title_prompt_log')"
    placement="right"
    width="746px"
    :closable="false"
  >
    <template #extra>
      <span class="close-btn" @click="onClose"><CloseOutlined /></span>
    </template>

    <div class="prompt-log-content" ref="contRef">
      <div class="prompt-log-items">
        <div class="prompt-log-label">
          <span>{{ t('label_prompt') }} </span>
          <a-tooltip>
            <template #title>{{ t('msg_prompt_tooltip') }}</template>
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
          <span>{{ t('label_context') }} </span>
          <a-tooltip>
            <template #title>{{ t('msg_context_tooltip') }}</template>
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
          <span>{{ t('label_user') }} </span>
          <a-tooltip>
            <template #title>{{ t('msg_user_tooltip') }}</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="prompt-log-item">
          <p>{{ promptLog.cur_question }}</p>
        </div>
      </div>

      <div class="prompt-log-items">
        <div class="prompt-log-label">
          <span>{{ t('label_assistant') }} </span>
          <a-tooltip>
            <template #title>{{ t('msg_assistant_tooltip') }}</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="prompt-log-item">
          <p>{{ promptLog.cur_answer }}</p>
        </div>
      </div>

      <div class="prompt-log-items" v-if="promptLog.recall_time > 0">
        <div class="prompt-log-label">
          <span>{{ t('label_recall_time') }} </span>
          <a-tooltip>
            <template #title>{{ t('msg_recall_time_tooltip') }}</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="prompt-log-item">
          <p>{{ promptLog.recall_time }}s</p>
        </div>
      </div>

      <div class="prompt-log-items" v-if="promptLog.request_time > 0">
        <div class="prompt-log-label">
          <span>{{ t('label_request_time') }} </span>
          <a-tooltip>
            <template #title>{{ t('msg_request_time_tooltip') }}</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="prompt-log-item">
          <p>{{ promptLog.request_time }}s</p>
        </div>
      </div>

      <div class="prompt-log-items">
        <div class="prompt-log-label">
          <span>{{ t('label_error') }} </span>
          <a-tooltip>
            <template #title>{{ t('msg_error_tooltip') }}</template>
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
import { useMathJax } from "@/composables/useMathJax.js";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-test.components.prompt-log-alert')

const {renderMath} = useMathJax()

const contRef = ref(null)
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
  renderMath(contRef.value)
}

defineExpose({
  open
})
</script>
