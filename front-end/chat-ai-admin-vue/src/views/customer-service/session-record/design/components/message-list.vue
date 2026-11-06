<template>
  <div class="design-message-list">
    <div class="design-message-scroll" ref="scrollBoxRef" @scroll="onScroll">
      <div class="design-message-inner">
        <div v-if="props.isEmpty || !props.messages.length" class="design-empty">
          <a-empty :image="Empty.PRESENTED_IMAGE_SIMPLE" :description="t('no_chat_messages')" />
        </div>
        <template v-else v-for="(item, index) in props.messages" :key="item.uid">
          <div v-if="showDateDivider(index)" class="date-divider">{{ formatDate(item.create_time) }}</div>

          <div v-if="item.is_customer == 1" class="message-row user" :id="'msg-' + item.uid">
            <div class="message-main">
              <div class="bubble user-bubble">
                <img
                  v-if="item.received_message_type == 'image' && item.media_id_to_oss_url"
                  v-viewer
                  class="msg-img"
                  :src="item.media_id_to_oss_url"
                  alt=""
                />
                <audio
                  v-else-if="item.msg_type == 4 && getMediaUrl(item)"
                  class="msg-audio"
                  :src="getMediaUrl(item)"
                  controls
                  preload="metadata"
                ></audio>
                <video
                  v-else-if="item.msg_type == 5 && getMediaUrl(item)"
                  class="msg-video"
                  :src="getMediaUrl(item)"
                  controls
                  preload="metadata"
                ></video>
                <MultipleMessage v-else-if="item.msg_type == 99" :message="item.content" />
                <span v-else class="message-text">{{ item.content }}</span>
              </div>
              <div class="message-time">{{ formatTime(item.create_time) }}</div>
            </div>
            <div class="message-avatar user-avatar">
              <img v-if="item.avatar" :src="item.avatar" alt="" />
              <span v-else>{{ avatarText(item) }}</span>
            </div>
          </div>

          <div v-else class="message-row robot" :id="'msg-' + item.uid">
            <div class="message-avatar robot-avatar">
              <img v-if="robotAvatar(item)" :src="robotAvatar(item)" alt="" />
              <span v-else>{{ avatarText(item, true) }}</span>
            </div>
            <div class="message-main">
              <div
                class="bubble robot-bubble"
                :class="{ 'has-debug': showQuoteFileBox(item) }"
                :title="showQuoteFileBox(item) ? 'Prompt 日志' : undefined"
                @click="onRobotBubbleClick($event, item)"
              >
                <ProcessTimeline
                  v-if="item.msg_type == 1 && hasVisibleProcessSteps(item)"
                  :item="item"
                />

                <template v-if="parseReplyList(item.reply_content_list).length">
                  <div class="reply-list">
                    <template v-for="(rc, idx) in parseReplyList(item.reply_content_list)" :key="idx">
                      <div v-if="(rc.reply_type || rc.type) === 'text'" class="reply-item">
                        <div class="message-text" v-html="rc.description"></div>
                      </div>
                      <div v-else-if="(rc.reply_type || rc.type) === 'image'" class="reply-item">
                        <img v-viewer class="msg-img" :src="rc.pic || rc.thumb_url" alt="" />
                      </div>
                      <div v-else-if="(rc.reply_type || rc.type) === 'url'" class="reply-item">
                        <a class="url-link" :href="rc.url" target="_blank">{{ rc.url }}</a>
                      </div>
                      <div v-else-if="(rc.reply_type || rc.type) === 'card'" class="reply-item">
                        <a class="reply-card" :href="rc.url || 'javascript:;'" target="_blank">
                          <img v-if="rc.thumb_url" :src="rc.thumb_url" class="reply-thumb" alt="" />
                          <span>{{ rc.title }}</span>
                        </a>
                      </div>
                      <div v-else-if="(rc.reply_type || rc.type) === 'imageText'" class="reply-item">
                        <a class="reply-image-text" :href="rc.url" target="_blank">
                          <img v-if="rc.thumb_url" :src="rc.thumb_url" class="reply-thumb" alt="" />
                          <span class="reply-image-text-title">{{ rc.title }}</span>
                          <span class="reply-image-text-desc">{{ rc.description }}</span>
                        </a>
                      </div>
                      <div v-else-if="(rc.reply_type || rc.type) === 'smartMenu'" class="reply-item">
                        <div class="reply-smart-menu">
                          <div v-if="rc.smart_menu?.menu_description" class="reply-smart-title">
                            {{ rc.smart_menu.menu_description }}
                          </div>
                          <template v-for="(line, li) in buildMenuLines(rc.smart_menu?.menu_content || [])" :key="li">
                            <div v-if="line.kind === 'text'" class="reply-line">{{ line.text }}</div>
                            <div v-else-if="line.kind === 'newline'" class="reply-line-empty"></div>
                            <span v-else-if="line.kind === 'html'" v-html="line.html"></span>
                            <a v-else-if="line.kind === 'keyword'" href="javascript:;" class="reply-keyword">
                              <span v-if="line.serial_no">{{ line.serial_no }}</span>{{ line.text }}
                            </a>
                          </template>
                        </div>
                      </div>
                    </template>
                  </div>
                </template>

                <div
                  v-if="item.msg_type == 1 && item.content"
                  class="message-content"
                >
                  <cherry-markdown :content="item.content"></cherry-markdown>
                </div>
                <div v-if="item.msg_type == 2" class="message-content" v-html="item.menu_json.content"></div>
                <img
                  v-if="item.msg_type == 3"
                  v-viewer
                  class="msg-img"
                  :src="item.content"
                  alt=""
                />
                <audio
                  v-if="item.msg_type == 4 && getMediaUrl(item)"
                  class="msg-audio"
                  :src="getMediaUrl(item)"
                  controls
                  preload="metadata"
                ></audio>
                <video
                  v-if="item.msg_type == 5 && getMediaUrl(item)"
                  class="msg-video"
                  :src="getMediaUrl(item)"
                  controls
                  preload="metadata"
                ></video>
                <VoiceMessage
                  v-if="item.voice_content && item.voice_content.length"
                  :voice_content="item.voice_content"
                />
              </div>
              <div class="message-footer">
                <span class="message-time">{{ formatTime(item.create_time) }}</span>
              </div>
            </div>
          </div>
        </template>
        <div v-if="loading" class="loading-box"><a-spin /></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { nextTick, ref, toRaw } from 'vue'
import dayjs from 'dayjs'
import { Empty } from 'ant-design-vue'
import ProcessTimeline from '@/views/clawbot/chat/components/messages/process-timeline.vue'
import { hasVisibleProcessSteps } from '@/utils/agent-process'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import VoiceMessage from '@/views/robot/robot-test/components/voice-message.vue'
import MultipleMessage from '@/views/robot/robot-test/components/messages/multiple-message.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.customer-service.index')
const emit = defineEmits(['scroll', 'scrollStart', 'scrollEnd', 'openPromptLog'])

const props = defineProps({
  messages: {
    type: Array,
    default: () => []
  },
  robotInfo: {
    type: Object,
    default: () => {}
  },
  isEmpty: {
    type: Boolean,
    default: false
  }
})

const loading = ref(false)
const scrollBoxRef = ref(null)
const scrollOption = {
  scrollTop: 0,
  scrollHeight: 0,
  clientHeight: 0,
  scrollStartDiff: 60,
  scrollEndDiff: 60,
  scrollDirection: ''
}

let scrollEventTimer = null
let onScrollEventLock = false

function onScroll(e) {
  if (onScrollEventLock) return
  if (scrollEventTimer) {
    clearTimeout(scrollEventTimer)
    scrollEventTimer = null
  }

  scrollEventTimer = setTimeout(() => {
    if (scrollOption.scrollTop - e.target.scrollTop > 0) scrollOption.scrollDirection = 'up'
    if (scrollOption.scrollTop - e.target.scrollTop < 0) scrollOption.scrollDirection = 'down'

    scrollOption.scrollTop = e.target.scrollTop
    scrollOption.scrollHeight = e.target.scrollHeight
    scrollOption.clientHeight = e.target.clientHeight
    emit('scroll', { ...scrollOption })

    const isAtTop = Math.abs(scrollOption.scrollTop) <= scrollOption.scrollStartDiff
    if (isAtTop && scrollOption.scrollDirection === 'up') onScrollStart()

    const isAtBottom = Math.abs(
      scrollOption.scrollHeight - scrollOption.scrollTop - scrollOption.clientHeight
    ) <= scrollOption.scrollEndDiff
    if (isAtBottom && scrollOption.scrollDirection === 'down') onScrollEnd()
  }, 100)
}

function onScrollStart() {
  if (!props.messages.length) return
  emit('scrollStart', { msg: props.messages[0] })
}

function onScrollEnd() {
  if (!props.messages.length) return
  emit('scrollEnd', { msg: props.messages[props.messages.length - 1] })
}

const scrollToBottom = () => {
  if (!scrollBoxRef.value) return
  nextTick(() => {
    onScrollEventLock = true
    scrollBoxRef.value.scrollTop = scrollBoxRef.value.scrollHeight + 1
    setTimeout(() => {
      scrollOption.scrollTop = scrollBoxRef.value.scrollTop
      onScrollEventLock = false
    }, 50)
  })
}

function scrollToMessage(id, direction) {
  nextTick(() => {
    onScrollEventLock = true
    if (!direction) direction = 'top'

    const scroller = scrollBoxRef.value
    const element = document.querySelector('#msg-' + id)
    if (element) {
      if (direction == 'top') {
        scroller.scrollTop = element.offsetTop
      } else {
        scroller.scrollTop = element.offsetTop - scroller.clientHeight + element.clientHeight
      }
    }

    setTimeout(() => {
      scrollOption.scrollTop = scrollBoxRef.value.scrollTop
      onScrollEventLock = false
    }, 50)
  })
}

function resetScroll() {
  scrollOption.scrollTop = 0
  scrollOption.scrollDirection = ''
}

function parseReplyList(val) {
  try {
    if (!val) return []
    if (Array.isArray(val)) return val
    if (typeof val === 'string') return JSON.parse(val || '[]')
    return []
  } catch {
    return []
  }
}

function buildMenuLines(menuContent) {
  const out = []
  ;(Array.isArray(menuContent) ? menuContent : []).forEach((mc) => {
    const type = String(mc?.menu_type || '')
    const text = String(mc?.content || '')
    if (type === '0') {
      if (text === '') {
        out.push({ kind: 'newline' })
      } else if (/<a[\s\S]*?<\/a>/.test(text)) {
        const sanitized = /href\s*=\s*['"]\s*#\s*['"]/i.test(text)
          ? text.replace(/href\s*=\s*['"]\s*#\s*['"]/ig, 'href="javascript:;"')
          : text.replace(/href=/ig, 'target="_blank" href=')
        out.push({ kind: 'html', html: sanitized })
      } else {
        out.push({ kind: 'text', text })
      }
    } else if (type === '1') {
      out.push({ kind: 'keyword', text, serial_no: mc?.serial_no || '' })
    }
  })
  return out.slice(0, 20)
}

function getMediaUrl(item) {
  return item.media_id_to_oss_url || item.content || ''
}

function showQuoteFileBox(item) {
  const newItem = JSON.parse(JSON.stringify(item))
  newItem.debug = newItem.debug_log ? JSON.parse(newItem.debug_log) : []
  return newItem.debug && newItem.debug.length > 0
}

function openPromptLog(item) {
  const newItem = JSON.parse(JSON.stringify(item))
  newItem.debug = newItem.debug_log ? JSON.parse(newItem.debug_log) : []
  emit('openPromptLog', toRaw(newItem))
}

function onRobotBubbleClick(event, item) {
  if (!showQuoteFileBox(item)) return
  if (event.target?.closest?.('a, button')) return
  openPromptLog(item)
}

function parseTime(time) {
  const value = Number(time)
  if (!Number.isFinite(value) || value <= 0) return null
  const date = dayjs(value * 1000)
  return date.isValid() ? date : null
}

function formatDate(time) {
  const date = parseTime(time)
  return date ? date.format('YYYY-MM-DD') : ''
}

function formatTime(time) {
  const date = parseTime(time)
  return date ? date.format('MM-DD HH:mm:ss') : ''
}

function dateKey(time) {
  return formatDate(time)
}

function showDateDivider(index) {
  const current = dateKey(props.messages[index]?.create_time)
  if (!current) return false
  if (index === 0) return true
  return current !== dateKey(props.messages[index - 1]?.create_time)
}

function avatarText(item, robot = false) {
  const source = item?.name || (robot ? props.robotInfo?.robot_name : '') || (robot ? 'AI' : '客')
  return source.slice(0, 1)
}

function robotAvatar(item) {
  return item?.robot_avatar || props.robotInfo?.robot_avatar_url || props.robotInfo?.robot_avatar || ''
}

defineExpose({
  scrollToBottom,
  scrollToMessage,
  resetScroll
})
</script>

<style lang="less" scoped>
.design-message-list {
  display: flex;
  flex: 1;
  min-height: 0;
  background: #fff;
  font-family: 'PingFang SC', 'Noto Sans SC', sans-serif;
}
.design-message-scroll {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}
.design-message-inner {
  display: flex;
  flex-direction: column;
  min-height: 100%;
  padding: 0 28px 32px;
}
.design-empty {
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  min-height: 240px;
}
.date-divider {
  align-self: center;
  margin: 28px 0 36px;
  padding: 0 12px;
  border-radius: 999px;
  color: #888;
  background: #f5f6f7;
  font-size: 12px;
  line-height: 24px;
}
.message-row {
  display: flex;
  gap: 10px;
  max-width: 88%;
  margin: 0 0 32px;

  &.user {
    align-self: flex-end;
    flex-direction: row-reverse;
  }

  &.robot {
    align-self: flex-start;
  }
}
.message-avatar {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  color: #fff;
  font-size: 12px;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}
.user-avatar {
  background: #86868b;
}
.robot-avatar {
  background: var(--app-primary);
}
.message-main {
  display: flex;
  flex-direction: column;
  min-width: 0;

  .user & {
    align-items: flex-end;
  }
}
.bubble {
  max-width: 480px;
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 14px;
  line-height: 24px;
  word-break: break-word;
}
.user-bubble {
  color: #242424;
  background: #e8efff;
}
.robot-bubble {
  color: #242424;
  background: #f2f3f5;

  &.has-debug {
    cursor: pointer;
  }

  &.has-debug:hover {
    box-shadow: 0 0 0 1px rgba(0, 82, 217, 0.2);
  }
}
.message-text {
  white-space: pre-wrap;
}
.message-time {
  margin-top: 6px;
  color: #888;
  font-size: 12px;
  line-height: 16px;
}
.message-footer {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 6px;
}
.message-content {
  white-space: pre-wrap;
  font-family: inherit;
  line-height: 24px;

  :deep(.cherry-markdown) {
    --font-family-sans: 'PingFang SC', 'Noto Sans SC', sans-serif;
    font-family: 'PingFang SC', 'Noto Sans SC', sans-serif;
    font-size: 14px;
    line-height: 24px;
    color: inherit;
  }

  :deep(p) {
    margin: 0 0 8px;
    line-height: 24px;

    &:last-child {
      margin-bottom: 0;
    }
  }
}
.msg-img {
  display: block;
  max-width: 300px;
  max-height: 300px;
  border-radius: 8px;
}
.msg-audio,
.msg-video {
  max-width: 100%;
}
.reply-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.reply-item {
  min-width: 0;
}
.url-link {
  color: var(--app-primary);
  word-break: break-all;
}
.reply-card {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #1d1d1f;
  text-decoration: none;
}
.reply-thumb {
  width: 48px;
  height: 48px;
  border-radius: 6px;
  object-fit: cover;
}
.reply-image-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
  color: #1d1d1f;
  text-decoration: none;
}
.reply-image-text-title {
  font-weight: 500;
}
.reply-image-text-desc {
  color: #8c8c8c;
  font-size: 12px;
}
.reply-smart-menu {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.reply-smart-title {
  font-weight: 500;
}
.reply-line {
  color: #1d1d1f;
}
.reply-line-empty {
  height: 8px;
}
.reply-keyword {
  color: var(--app-primary);
}
.loading-box {
  display: flex;
  justify-content: center;
  padding: 16px 0;
}
</style>
