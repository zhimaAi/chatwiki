<style lang="less" scoped>
.feedbacks-log-alert {
  .feedbacks-log-content {
    flex: 0 0 300px;
    overflow-y: auto;
    padding: 24px;
    border-right: 4px solid #f2f4f7;

    .feedbacks-log-items {
      margin-bottom: 24px;
    }
    .feedbacks-log-label {
      margin-bottom: 8px;
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
      white-space: pre-wrap;
      overflow-wrap: anywhere;

      p { margin: 0; }
      .item-type + p { margin-top: 8px; }

      &:last-child {
        margin-bottom: 0;
      }

      .item-type {
        color: #3a4559;
        font-family: "PingFang SC";
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 22px;
        display: flex;
        gap: 8px;
        align-items: center;
      }

      .quotes-title {
        display: flex;
        margin-top: 8px;
        align-items: center;
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
        color: #3a4559;
        white-space: pre-wrap;
      }
    }
  }
}
.detail-state { display: flex; flex: 1; align-items: center; justify-content: center; flex-direction: column; gap: 16px; padding: 24px; }
.feedback-history { display: flex; flex: 1; flex-direction: column; min-width: 0; min-height: 0; }
.load-more { display: flex; justify-content: center; }
@media (max-width: 640px) {
  :global(.feedbacks-log-alert .ant-drawer-body) { flex-direction: column; }
  .feedbacks-log-alert .feedbacks-log-content { flex-basis: auto; max-height: 35%; border-right: 0; border-bottom: 4px solid #f2f4f7; }
}
</style>

<template>
  <a-drawer
    class="feedbacks-log-alert"
    v-model:open="show"
    :title="t('title_feedback_details')"
    width="min(1100px, 100vw)"
    :body-style="{ padding: 0, display: 'flex', minHeight: 0 }"
    destroy-on-close
    @close="onClose"
  >
    <div v-if="loading" class="detail-state"><a-spin /></div>
    <div v-else-if="error" class="detail-state">
      <a-alert type="error" :message="error" show-icon />
      <a-button @click="loadDetail">{{ t('btn_retry') }}</a-button>
    </div>
    <div v-else class="feedbacks-log-content">
      <div class="feedbacks-log-items">
        <div class="feedbacks-log-label">
          <span>{{ t('label_feedback') }} </span>
        </div>
        <div class="feedbacks-log-item">
          <div class="item-type"><svg-icon style="color: #2475FC;" :name="feedbacksLog.type == '1' ? 'like-active' : 'dislike-active'" />{{ feedbacksLog.type == '1' ? t('label_like') : t('label_dislike') }}</div>
          <p v-if="feedbacksLog.content">{{ feedbacksLog.content }}</p>
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
          <p v-else>--</p>
        </div>
      </div>

      <div class="feedbacks-log-items">
        <div class="feedbacks-log-label">
          <span>{{ t('label_used_model') }} </span>
        </div>
        <div class="feedbacks-log-item">
          <p>{{ modelName }}</p>
        </div>
      </div>
      <div class="feedbacks-log-items">
        <div class="feedbacks-log-label"><span>{{ t('label_customer_name') }}</span></div>
        <div class="feedbacks-log-item">{{ feedbacksLog.customer?.name || feedbacksLog.customer?.nickname || '--' }}</div>
      </div>
      <div class="feedbacks-log-items">
        <div class="feedbacks-log-label"><span>{{ t('label_channel') }}</span></div>
        <div class="feedbacks-log-item">{{ feedbacksLog.channel_name || '--' }}</div>
      </div>

      <div class="feedbacks-log-items" v-for="(item, index) in feedbacksLog.quotes" :key="`${item.FileId}-${index}`">
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
            <router-link class="quotes-title" :to="{ path: '/library/preview', query: { id: item.FileId } }">{{ item.FileName }}</router-link>
          </div>
        </div>
      </div>

    </div>
    <div v-if="!loading && !error" class="feedback-history" :aria-label="t('label_chat_history')" :aria-busy="pageLoading">
      <a-alert v-if="messageError" type="error" :message="messageError" show-icon>
        <template #action><a-button type="link" @click="loadMessages(failedDirection, true)">{{ t('btn_retry') }}</a-button></template>
      </a-alert>
      <MessageList
        ref="messageListRef"
        :messages="messages"
        :robot-info="feedbacksLog.robot"
        :highlight-message-id="feedbacksLog.ai_message_id"
        read-only
        @scroll-start="loadMessages('before')"
        @scroll-end="loadMessages('after')"
      >
        <template #before>
          <div v-if="hasBefore" class="load-more">
            <a-button type="link" :loading="pageLoading" @click="loadMessages('before', true)">{{ t('btn_earlier_messages') }}</a-button>
          </div>
        </template>
        <template #after>
          <div v-if="hasAfter" class="load-more">
            <a-button type="link" :loading="pageLoading" @click="loadMessages('after', true)">{{ t('btn_later_messages') }}</a-button>
          </div>
        </template>
      </MessageList>
    </div>
  </a-drawer>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { getFeedbackDetail, getFeedbackMessages } from '@/api/manage/index.js'
import { useI18n } from '@/hooks/web/useI18n'
import { DEFAULT_USER_AVATAR } from '@/constants'
import { extractVoiceInfo, removeVoiceFormat } from '@/utils'
import { restoreAgentProcess } from '@/utils/agent-process'
import MessageList from '../../session-record/components/message-list.vue'

const props = defineProps({ robotId: { type: [String, Number], default: '' } })
const { t } = useI18n('views.robot.robot-config.qa-feedback.components.feedbacks-log-alert')
const show = ref(false)
const loading = ref(false)
const error = ref('')
const feedbacksLog = ref({ robot: {}, quotes: [] })
const messages = ref([])
const messageListRef = ref(null)
const hasBefore = ref(false)
const hasAfter = ref(false)
const pageLoading = ref(false)
const messageError = ref('')
const failedDirection = ref('before')
let feedbackId = ''
let requestVersion = 0

const modelName = computed(() => {
  const { use_model, corp_name } = feedbacksLog.value.robot
  return use_model ? `${use_model}${corp_name ? `(${corp_name})` : ''}` : '--'
})

function parseJson(value, fallback) {
  if (typeof value !== 'string') return value ?? fallback
  try { return JSON.parse(value) ?? fallback } catch { return fallback }
}

function formatMessages(list) {
  const { robot, customer } = feedbacksLog.value
  return list.map((raw) => {
    const item = { ...raw, uid: `feedback-${raw.id}`, loading: false }
    if (item.is_customer == 1) {
      item.name = item.name || customer?.name || customer?.nickname
      item.avatar = item.avatar || customer?.avatar || DEFAULT_USER_AVATAR
    } else {
      item.name = item.name || robot.robot_name
      item.robot_avatar = item.avatar || robot.robot_avatar || '/upload/default/robot_avatar.svg'
    }
    item.menu_json = parseJson(item.menu_json, {})
    item.reply_content_list = parseJson(item.reply_content_list, [])
    item.voice_content = extractVoiceInfo(item.content)
    item.content = removeVoiceFormat(item.content)
    restoreAgentProcess(item, { historical: true })
    return item
  })
}

async function loadDetail() {
  const version = ++requestVersion
  loading.value = true
  pageLoading.value = false
  error.value = ''
  messageError.value = ''
  messages.value = []
  feedbacksLog.value = { robot: {}, quotes: [] }
  try {
    const [info, history] = await Promise.all([
      getFeedbackDetail({ id: feedbackId }),
      getFeedbackMessages({ id: feedbackId })
    ])
    if (version !== requestVersion) return
    feedbacksLog.value = {
      ...info.data,
      robot: parseJson(info.data.robot, {}),
      quotes: parseJson(info.data.quotes, [])
    }
    messages.value = formatMessages(history.data.list)
    hasBefore.value = history.data.has_more_before
    hasAfter.value = history.data.has_more_after
    loading.value = false
    await nextTick()
    if (version === requestVersion) messageListRef.value?.scrollToMessage(`feedback-${feedbacksLog.value.ai_message_id}`)
  } catch (err) {
    if (version === requestVersion) error.value = err?.message || err?.msg || t('msg_load_failed')
  } finally {
    if (version === requestVersion) loading.value = false
  }
}

async function loadMessages(direction, retry = false) {
  if (loading.value || pageLoading.value || (messageError.value && !retry)) return
  if (direction === 'before' ? !hasBefore.value : !hasAfter.value) return
  const version = requestVersion
  const anchor = direction === 'before' ? messages.value[0] : messages.value[messages.value.length - 1]
  if (!anchor) return
  pageLoading.value = true
  messageError.value = ''
  try {
    const res = await getFeedbackMessages({
      id: feedbackId,
      [direction === 'before' ? 'min_id' : 'max_id']: anchor.id
    })
    if (version !== requestVersion) return
    const page = formatMessages(res.data.list)
    if (direction === 'before') {
      messages.value = [...page, ...messages.value]
      hasBefore.value = res.data.has_more_before
      await nextTick()
      if (version === requestVersion) messageListRef.value?.scrollToMessage(anchor.uid)
    } else {
      messages.value.push(...page)
      hasAfter.value = res.data.has_more_after
    }
  } catch (err) {
    if (version === requestVersion) {
      messageError.value = err?.message || err?.msg || t('msg_load_failed')
      failedDirection.value = direction
    }
  } finally {
    if (version === requestVersion) pageLoading.value = false
  }
}

function onClose() {
  show.value = false
  requestVersion++
  loading.value = false
  pageLoading.value = false
}

function open(record) {
  feedbackId = record.id
  show.value = true
  loadDetail()
}

watch(() => props.robotId, onClose)
onBeforeUnmount(onClose)
defineExpose({ open })
</script>
