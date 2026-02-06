<style lang="less" scoped>
.feedbacks-log-alert {
  .close-btn {
    font-size: 16px;
    color: rgba(0, 0, 0, 0.45);
    cursor: pointer;
  }

  .feedbacks-log-content {
    .feedbacks-log-items {
      margin-bottom: 24px;
    }
    .feedbacks-log-label {
      margin-bottom: 4px;
      font-size: 16px;
      font-weight: 600;
      line-height: 24px;
      color: rgb(0, 0, 0);
    }

    .feedbacks-log-item {
      line-height: 22px;
      padding: 12px 16px;
      margin-bottom: 8px;
      font-size: 14px;
      color: #3a4559;
      background-color: #f2f4f7;
      border-radius: 4px;
      word-break: break-word;

      &:last-child {
        margin-bottom: 0;
      }

      .item-type {
        color: #595959;
        font-family: "PingFang SC";
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 22px;
        display: flex;
        align-items: center;
      }

      .quotes-title {
        display: flex;
        margin-top: 8px;
        align-items: center;
        width: 430px;
        color: #164799;
        font-family: "PingFang SC";
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 22px;
        cursor: pointer;
      }

      .document-content {
        line-height: 22px;
        margin-top: 8px;
        font-size: 14px;
        color: #595959;
        white-space: pre-wrap;
      }
    }
  }
}
</style>

<template>
  <a-drawer
    class="feedbacks-log-alert"
    v-model:open="show"
    :title="t('title_feedback_details')"
    placement="right"
    width="746px"
    :closable="false"
  >
    <template #extra>
      <span class="close-btn" @click="onClose"><CloseOutlined /></span>
    </template>

    <div class="feedbacks-log-content">
      <div class="feedbacks-log-items">
        <div class="feedbacks-log-label">
          <span>{{ t('label_question') }} </span>
        </div>

        <div class="feedbacks-log-item">
          <p>{{ feedbacksLog.question }}</p>
        </div>
      </div>

      <div class="feedbacks-log-items">
        <div class="feedbacks-log-label">
          <span>{{ t('label_answer') }} </span>
        </div>
        
        <div class="feedbacks-log-item">
          <p>{{ feedbacksLog.answer }}</p>
        </div>
      </div>

      <div class="feedbacks-log-items">
        <div class="feedbacks-log-label">
          <span>{{ t('label_feedback') }} </span>
        </div>
        <div class="feedbacks-log-item">
          <div v-if="feedbacksLog.type == '1'" class="item-type"><svg-icon style="font-size: 16px; color: #2475FC; margin-right: 4px;" name="like-active" />{{ t('label_like') }}</div>
          <div v-if="feedbacksLog.type == '2'" class="item-type"><svg-icon style="font-size: 16px; color: #2475FC; margin-right: 4px;" name="dislike-active" />{{ t('label_dislike') }}</div>
          <p>{{ feedbacksLog.content }}</p>
        </div>
      </div>

      <div class="feedbacks-log-items" v-for="(item, index) in feedbacksLog.quotes" :key="item.FileId">
        <div class="feedbacks-log-label">
          <span>{{ t('label_reference_content') }}{{ index + 1 }} </span>
        </div>
        <div class="feedbacks-log-item">
          <div class="quotes-box" >
            <div class="document-content" v-if="item.Type == 1">{{ item.Content }}</div>
            <div v-else-if="item.Type == 2">
              <div class="document-content">Q：{{ item.Question }}</div>
              <div class="document-content">A：{{ item.Answer }}</div>
            </div>
            <div class="quotes-title" @click="onToQuotes(item)"><svg-icon style="font-size: 16px; color: #2475FC; margin-right: 4px;" name="quotes" />{{ item.FileName }}</div>
          </div>
        </div>
      </div>

      <div class="feedbacks-log-items">
        <div class="feedbacks-log-label">
          <span>{{ t('label_chat_mode') }} </span>
        </div>
        <div class="feedbacks-log-item">
          <p v-if="feedbacksLog.robot.chat_type == '1'">{{ t('label_only_knowledge_base') }}</p>
          <p v-else-if="feedbacksLog.robot.chat_type == '2'">{{ t('label_direct_connection') }}</p>
          <p v-else-if="feedbacksLog.robot.chat_type == '3'">{{ t('label_hybrid') }}</p>
        </div>
      </div>

      <div class="feedbacks-log-items">
        <div class="feedbacks-log-label">
          <span>{{ t('label_used_model') }} </span>
        </div>
        <div class="feedbacks-log-item">
          <p>{{ feedbacksLog.robot.use_model }}({{ feedbacksLog.robot.corp_name }})</p>
        </div>
      </div>
    </div>
  </a-drawer>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { CloseOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.qa-feedback.components.feedbacks-log-alert')
const router = useRouter()
const show = ref(false)
const feedbacksLog = reactive({
  question: '',
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
  feedbacksLog.question = ''
  feedbacksLog.library = []
  feedbacksLog.context_qa = []
  feedbacksLog.cur_question = ''
  feedbacksLog.cur_answer = ''
  feedbacksLog.error = ''
}

const open = (msg) => {
  reset()
  for(const key in msg) {
    feedbacksLog[key] = msg[key]
    if (key == 'robot') {
      feedbacksLog[key] = JSON.parse(msg[key])
    }
    if (key == 'quotes') {
      feedbacksLog[key] = JSON.parse(msg[key])
    }
  }
  show.value = true
}

const onToQuotes = (item) => {
  router.push('/library/preview?id=' + item.FileId)
}

defineExpose({
  open
})
</script>
