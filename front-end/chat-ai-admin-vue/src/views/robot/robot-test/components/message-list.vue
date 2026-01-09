<style lang="less" scoped>
.message-list-wrapper {
  width: 100%;
  height: 100%;
}

.scroll-box {
  width: 100%;
  height: 100%;
  overflow-y: auto;

  &::-webkit-scrollbar {
    display: none;
  }
}

.markdown-content{
  white-space: normal;
}

.msg-img {
  width: 100%;
  height: 100%;
  max-width: 300px;
  max-height: 300px;
}

.message-list {
  .message-item {
    display: flex;
    margin-top: 24px;

    .user-avatar,
    .robot-avatar {
      display: block;
      width: 32px;
      height: 32px;
      border-radius: 50%;
    }

    .itme-left {
      margin-right: 8px;
    }

    .itme-right {
      flex: 1;
      overflow: hidden;
    }
  }

  .user-message {
    .item-body {
      padding: 5px 0;
    }

    .message-content {
      line-height: 22px;
      font-size: 14px;
      font-weight: 400;
      color: #3a4559;
      white-space: pre-wrap;
    }
  }

  .robot-message {
    .item-body {
      padding: 12px 16px;
      border-radius: 8px;
      width: auto;
      min-height: 32px;
      max-width: 100%;
      overflow: hidden;
      background-color: #fff;
    }

    .message-content {
      line-height: 22px;
      font-size: 14px;
      font-weight: 400;
      color: #3a4559;
      white-space: pre-wrap;
    }
  }

  .message-menus {
    .menu-item {
      line-height: 22px;
      padding: 8px 16px;
      margin-top: 8px;
      font-size: 14px;
      border-radius: 4px;
      color: rgb(22, 71, 153);
      background: #f2f4f7;
      cursor: pointer;
    }
  }

.message-action-wrap {
  .message-action {
    display: flex;
    padding-top: 12px;
    margin-top: 12px;
    border-top: 1px solid #edeff2;

      .action-btn {
        position: relative;
        padding: 0 8px;
        line-height: 22px;
        font-size: 14px;
        font-weight: 400;
        color: #7a8699;
        cursor: pointer;
        &::before {
          display: block;
          position: absolute;
          content: ' ';
          right: 0;
          top: 5px;
          width: 1px;
          height: 12px;
          background-color: rgba(5, 5, 5, 0.06);
        }
        &:last-child::before {
          display: none;
        }
      }
    }
  }

  .file-items {
    display: flex;
    flex-wrap: wrap;
    gap: 12px 24px;
    line-height: 22px;
    font-weight: 400;
    color: #164799;
    font-size: 14px;
    padding-bottom: 12px;
    border-bottom: 1px solid #edeff2;
    margin-bottom: 12px;
    .file-name {
      cursor: pointer;
      display: flex;
      align-items: baseline;
      gap: 4px;
      color: #164799;
      position: relative;
      .svg-action {
        position: relative;
        top: 2px;
      }
    }
  }
  .label-flex-block {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .thinking-label-wrapper {
    display: flex;
    align-items: center;
    margin-bottom: 8px;

    .thinking-label {
      display: flex;
      align-items: center;
      height: 32px;
      padding: 0 16px;
      border-radius: 8px;
      background: #e4e6eb;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        background: #d8dde6;
      }
      .think-icon,
      .loading {
        margin-right: 8px;
        font-size: 16px;
        color: #262626;
      }
      .label-text {
        font-size: 14px;
        font-weight: 400;
        color: #262626;
      }
      .arrow-down {
        margin-left: 8px;
        font-size: 16px;
        color: #262626;
        cursor: pointer;
      }
    }
    .tip {
      margin-left: 8px;
      font-size: 16px;
      color: #8c8c8c;
      cursor: pointer;
    }

    &.reasoning_open {
      .arrow-down {
        transform: rotate(180deg);
      }
    }
  }

  .thinking-content {
    position: relative;
    line-height: 22px;
    padding-bottom: 0;
    padding-left: 16px;
    margin-bottom: 12px;
    font-size: 14px;
    font-weight: 400;
    color: #8c8c8c;
    border-bottom: 1px solid #edeff2;
    // 4px竖线
    &::before {
      display: block;
      position: absolute;
      content: '';
      left: 0;
      top: 4px;
      bottom: 20px;
      width: 4px;
      background-color: #d9d9d9;
    }
  }

  .reply-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 8px;
  }

.reply-url .url-link {
  color: #2475fc;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  text-decoration: none;
}

.reply-url .url-row {
  display: inline-flex;
  align-items: center;
  gap: 12px;
}

  .reply-card .card-row {
    display: inline-flex;
    align-items: center;
    gap: 12px;
  }

  .reply-card .card-title-box {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .reply-card .card-thumb {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
  }

  .reply-card .card-title {
    color: #595959;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }

  .reply-imageText .imageText-row {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .reply-imageText .imageText-thumb {
    width: 40px;
    height: 40px;
    border-radius: 6px;
    object-fit: cover;
  }

  .reply-imageText .imageText-text {
    display: flex;
    flex-direction: column;
  }

  .reply-imageText .imageText-title {
    color: #595959;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }

  .reply-imageText .imageText-desc {
    color: #8c8c8c;
    font-size: 12px;
    font-style: normal;
    font-weight: 400;
    line-height: 16px;
  }

  .reply-item {
    border-radius: 4px 16px 16px 16px;
    padding: 16px 12px;
    background: #fff;
  }

  .reply-image .msg-img {
    max-width: 500px;
  }

  .reply-smartMenu {
    color: #1a1a1a;
    .smart-menu-box {
      .card-title {
        white-space: pre-wrap;
        align-self: stretch;
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 20px;
        margin-bottom: 14px;
      }
      .card-text {
        display: flex;
        flex-direction: column;
        gap: 10px;
        .reply-line {
          line-height: 22px;
          .line-text {
            color: #3a4559;
          }
          .empty-line {
            height: 22px;
          }
          .link {
            display: flex;
            align-items: center;
            gap: 4px;
          }
        }
      }
    }
  }
}
</style>

<template>
  <div class="message-list-wrapper">
    <div class="scroll-box" ref="scrollBoxRef" @scroll="onScroll">
      <div class="message-list">
        <template v-for="item in props.messages" :key="item.uid">
          <!-- 用户的消息 -->
          <div
            class="message-item user-message"
            :id="'msg-' + item.uid"
            v-if="item.is_customer == 1"
            :data-msg_type="item.msg_type"
          >
            <div class="itme-left">
              <img class="user-avatar" :src="item.avatar" />
            </div>
            <div class="itme-right">
              <div class="item-body">
                <!-- 收到消息类型处理，目前只处理了image 后续有其他的在这里添加 -->
                <template v-if="item.msg_type == 3">
                  <div class="message-content">
                    <img v-viewer class="msg-img" :src="item.media_id_to_oss_url" alt="" />
                  </div>
                </template>
                <template v-else-if="item.msg_type == 99">
                  <div class="message-content">
                    <MultipleMessage :message="item.content" />
                  </div>
                </template>
                <template v-else>
                  <TextMessage class="message-content" :message="item.content" />
                </template>
              </div>
            </div>
          </div>
          <!-- 机器人的消息 -->
          <div class="message-item robot-message" :id="'msg-' + item.uid" v-else>
            <div class="itme-left">
              <a-spin size="small" :spinning="item.loading">
                <img class="robot-avatar" :src="item.robot_avatar" />
              </a-spin>
            </div>

            <div class="itme-right">
              <template v-if="parseReplyList(item.reply_content_list).length">
                <div class="reply-list">
                  <template v-for="(rc, idx) in parseReplyList(item.reply_content_list)" :key="idx">
                    <div v-if="(rc.reply_type || rc.type) === 'text'" class="reply-item reply-text">
                      <div class="message-content" v-html="rc.description"></div>
                    </div>
                    <div v-else-if="(rc.reply_type || rc.type) === 'image'" class="reply-item reply-image">
                      <div class="message-content">
                        <img v-viewer class="msg-img" :src="rc.pic || rc.thumb_url" />
                      </div>
                    </div>
                    <div v-else-if="(rc.reply_type || rc.type) === 'url'" class="reply-item reply-url">
                      <div class="url-row">
                        <a class="url-link" :href="rc.url" target="_blank">{{ rc.url }}</a>
                      </div>
                    </div>
                    <div v-else-if="(rc.reply_type || rc.type) === 'card'" class="reply-item reply-card">
                      <div class="card-row">
                        <img v-if="rc.thumb_url" :src="rc.thumb_url" class="card-thumb" />
                        <div class="card-title-box">
                          <svg-icon class="think-icon" name="applet"></svg-icon>
                          <span class="card-title">{{ rc.title }}</span>
                        </div>
                      </div>
                    </div>
                  <div v-else-if="(rc.reply_type || rc.type) === 'imageText'" class="reply-item reply-imageText">
                    <a class="imageText-row" :href="rc.url" target="_blank">
                      <img v-if="rc.thumb_url" :src="rc.thumb_url" class="imageText-thumb" />
                      <div class="imageText-text">
                        <div class="imageText-title">{{ rc.title }}</div>
                        <div class="imageText-desc">{{ rc.description }}</div>
                      </div>
                    </a>
                  </div>
                  <div v-else-if="(rc.reply_type || rc.type) === 'smartMenu'" class="reply-item reply-smartMenu">
                    <div class="smart-menu-box">
                      <div class="card-title" v-if="rc.smart_menu && rc.smart_menu.menu_description">{{ rc.smart_menu.menu_description }}</div>
                      <div class="card-text">
                        <template v-for="(line, li) in buildMenuLines(rc.smart_menu?.menu_content || [])" :key="li">
                          <div class="reply-line" @click="onSmartReplyLineClick($event)">
                            <span class="line-text" v-if="line.kind === 'text'">{{ line.text }}</span>
                            <div v-else-if="line.kind === 'newline'" class="empty-line"></div>
                            <span v-else-if="line.kind === 'html'" v-html="line.html"></span>
                            <a v-else-if="line.kind === 'keyword'" href="javascript:;" class="link" @click.prevent="onClickSmartMenuKeyword(line.text)">
                              <div v-if="line.serial_no">{{ line.serial_no }}</div>
                              {{ line.text }}
                            </a>
                          </div>
                        </template>
                      </div>
                    </div>
                  </div>
                  </template>
                </div>
              </template>
              <!-- 检索知识库 -->
              <div class="label-flex-block">
                <div
                  class="thinking-label-wrapper"
                  :class="{ reasoning_open: item.show_quote_file }"
                  v-if="item.msg_type == 1 && isShowQuoteFileProgress"
                >
                  <div class="thinking-label" @click="toggleQuoteFiel(item)">
                    <template v-if="item.quote_loading">
                      <LoadingOutlined class="loading" />
                      <span class="label-text">正在检索知识库...</span>
                    </template>
                    <template v-else>
                      <svg-icon class="think-icon" name="quote-file"></svg-icon>
                      <span class="label-text">检索到{{ item.quote_file.length }}个知识库文档</span>
                    </template>
                    <svg-icon
                      name="arrow-down"
                      class="arrow-down"
                      v-if="item.quote_file.length"
                    ></svg-icon>
                  </div>
                </div>

                <!-- 思考过程label -->
                <div
                  class="thinking-label-wrapper"
                  :class="{ reasoning_open: item.show_reasoning }"
                  v-if="item.reasoning_content"
                >
                  <div class="thinking-label" @click="toggleReasonProcess(item)">
                    <LoadingOutlined class="loading" v-if="item.reasoning_status" />
                    <svg-icon class="think-icon" name="think" v-else></svg-icon>
                    <span class="label-text">{{
                      item.reasoning_status ? '深度思考中...' : '已完成深度思考'
                    }}</span>

                    <svg-icon
                      name="arrow-down"
                      class="arrow-down"
                      v-if="!item.reasoning_status"
                    ></svg-icon>
                  </div>
                  <a-tooltip>
                    <template #title>在应用设置中关闭推理过程开关，将不再显示推理过程</template>
                    <InfoCircleOutlined class="tip" />
                  </a-tooltip>
                </div>
              </div>

              <div class="item-body">
                <div
                  class="file-items"
                  v-if="item.show_quote_file && item.quote_file && item.quote_file.length > 0"
                >
                  <div
                    class="file-item"
                    v-for="file in item.quote_file"
                    :key="file.id"
                    @click="openLibrary(item.quote_file, file, item.id)"
                  >
                    <a class="file-name">
                      <svg-icon class="think-icon" name="quote-file"></svg-icon>
                      <span v-if="file.file_name">{{ file.file_name }}</span>
                      <span v-else>{{ file.library_name }}-精选</span>
                    </a>
                  </div>
                </div>
                <div class="thinking-content" v-if="item.show_reasoning && item.reasoning_content">
                  <cherry-markdown :content="item.reasoning_content"></cherry-markdown>
                </div>
                <template v-if="item.msg_type == 1">
                  <div class="message-content" v-viewer>
                    <cherry-markdown :content="item.content"></cherry-markdown>
                  </div>
                  <div class="message-menus">
                    <div
                      class="menu-item"
                      @click="onClickMeun(item)"
                      v-for="(item, index) in item.question"
                      :key="index"
                    >
                      {{ item }}
                    </div>
                  </div>
                </template>

                <template v-if="item.msg_type == 2">
                  <cherry-markdown class="markdown-content" :content="item.menu_json.content" v-if="item.isWelcome"></cherry-markdown>
                  <div class="message-content" v-html="item.menu_json.content" v-else></div>
                  <div class="message-menus">
                    <div
                      class="menu-item"
                      @click="onClickMeun(item)"
                      v-for="(item, index) in item.menu_json.question"
                      :key="index"
                    >
                      {{ item }}
                    </div>
                  </div>
                </template>

                <template v-if="item.msg_type == 3">
                  <div class="message-content">
                    <img v-viewer class="msg-img" :src="item.content" alt="" />
                  </div>
                </template>

                <div
                  class="message-action-wrap"
                  v-if="showQuoteFileBox(item)"
                  v-show="!item.question"
                >
                  <div class="message-action" v-if="item.debug && item.debug.length > 0">
                    <div class="action-btn">
                      <span><a @click="openPromptLog(item)">Prompt 日志</a></span>
                    </div>
                  </div>
                </div>
              </div>
              <!-- 语音消息 -->
              <div v-if="item.voice_content && item.voice_content.length">
                <VoiceMessage :voice_content="item.voice_content" />
              </div>
              <GuessYouWant
                v-if="
                  ((item.guess_you_want && item.guess_you_want.length) ||
                    (common_question_list.length &&
                      props.robotInfo.enable_common_question == 'true')) &&
                  item.question_tabkey > 0
                "
                :item="item"
                :enable_common_question="props.robotInfo.enable_common_question == 'true'"
                :common_question_list="common_question_list"
                @clickMeun="onClickMeun"
              ></GuessYouWant>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, toRaw, computed } from 'vue'
import { LoadingOutlined, InfoCircleOutlined } from '@ant-design/icons-vue'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import GuessYouWant from './guess-you-want.vue'
import VoiceMessage from './voice-message.vue'
import TextMessage from './messages/text-message.vue'
import MultipleMessage from './messages/multiple-message.vue'


const emit = defineEmits([
  'clickMsgMeun',
  'scroll',
  'scrollStart',
  'scrollEnd',
  'openPromptLog',
  'openLibrary'
])

const props = defineProps({
  messages: {
    type: Array,
    default: () => []
  },
  robotInfo: {
    type: Object,
    default: () => {}
  }
})

const toggleReasonProcess = (item) => {
  item.show_reasoning = !item.show_reasoning
}

const toggleQuoteFiel = (item) => {
  item.show_quote_file = !item.show_quote_file
}

const common_question_list = computed(() => {
  if (props.robotInfo.common_question_list.length) {
    return JSON.parse(props.robotInfo.common_question_list)
  }
  return []
})

const isShowQuoteFileProgress = computed(() => {
  // 普通机器人才显示
  return (props.robotInfo.chat_type == 1 || props.robotInfo.chat_type == 3) && props.robotInfo.application_type == '0'
})

const onClickMeun = (item) => {
  emit('clickMsgMeun', item)
}

const scrollBoxRef = ref(null)
const scrollOption = {
  scrollTop: 0,
  scrollHeight: 0,
  clientHeight: 0,
  scrollStartDiff: 60,
  scrollEndDiff: 60,
  scrollDirection: '' // 滚动方向
}

let scrollEventTimer = null // 滚动条防抖
let onScrollEventLock = false // 时间触发锁

// 监听滚动条滚动
function onScroll(e) {
  if (onScrollEventLock) {
    return
  }

  if (scrollEventTimer) {
    clearTimeout(scrollEventTimer)
    scrollEventTimer = null
  }

  scrollEventTimer = setTimeout(() => {
    if (scrollOption.scrollTop - e.target.scrollTop > 0) {
      scrollOption.scrollDirection = 'up'
    }

    if (scrollOption.scrollTop - e.target.scrollTop < 0) {
      scrollOption.scrollDirection = 'down'
    }

    scrollOption.scrollTop = e.target.scrollTop
    scrollOption.scrollHeight = e.target.scrollHeight
    scrollOption.clientHeight = e.target.clientHeight

    emit('scroll', { ...scrollOption })
    // 触顶
    let isAtTop = Math.abs(scrollOption.scrollTop) <= scrollOption.scrollStartDiff

    if (isAtTop && scrollOption.scrollDirection === 'up') {
      onScrollStart()
    }
    // 触底
    let isAtBottom =
      Math.abs(scrollOption.scrollHeight - scrollOption.scrollTop - scrollOption.clientHeight) <=
      scrollOption.scrollEndDiff

    if (isAtBottom && scrollOption.scrollDirection === 'down') {
      onScrollEnd()
    }
  }, 50)
}

function onScrollStart() {
  // 如果消息列表为空可能是断线重连等逻辑手动清空了消息列表造成的抖动，此时不触发事件
  if (props.messages.length == 0) {
    return
  }
  emit('scrollStart', {
    msg: props.messages[0]
  })
}

function onScrollEnd() {
  // 如果消息列表为空可能是断线重连等逻辑手动清空了消息列表造成的抖动，此时不触发事件
  if (props.messages.length == 0) {
    return
  }
  emit('scrollEnd', {
    msg: props.messages[props.messages.length - 1]
  })
}

const scrollToBottom = () => {
  if (!scrollBoxRef.value) {
    return
  }
  nextTick(() => {
    // 手动控制滚动到底部不触发触底事件
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
    // 手动控制滚动到底部不触发触底事件
    onScrollEventLock = true

    if (!direction) {
      direction = 'top'
    }

    let scroller = scrollBoxRef.value
    let element = document.querySelector('#msg-' + id)

    if (direction == 'top') {
      scroller.scrollTop = element.offsetTop
    } else {
      scroller.scrollTop = element.offsetTop - scroller.clientHeight + element.clientHeight
    }

    setTimeout(() => {
      scrollOption.scrollTop = scrollBoxRef.value.scrollTop
      onScrollEventLock = false
    }, 50)
  })
}

// 重置滚动条状态
function resetScroll() {
  scrollOption.scrollTop = 0
  scrollOption.scrollDirection = ''
}

function showQuoteFileBox(item) {
  if (item.quote_file && item.quote_file.length > 0) {
    return true
  }

  if (item.debug && item.debug.length > 0) {
    return true
  }

  return false
}

// 打开Prompt日志
function openPromptLog(item) {
  emit('openPromptLog', toRaw(item))
}

// 打开知识库
function openLibrary(files, file, message_id) {
  let newfiles = toRaw(files)
  file.message_id = file.message_id || message_id
  newfiles.forEach((item) => {
    item.message_id = item.message_id || message_id
  })

  emit('openLibrary', newfiles, toRaw(file))
}

function parseReplyList(val) {
  try {
    if (!val) return []
    if (Array.isArray(val)) return val
    if (typeof val === 'string') return JSON.parse(val || '[]')
    return []
  } catch (_e) {
    return []
  }
}

function buildMenuLines(menu_content) {
  const out = []
  ;(Array.isArray(menu_content) ? menu_content : []).forEach((mc) => {
    const t = String(mc?.menu_type || '')
    const txt = String(mc?.content || '')
    if (t === '0') {
      if (txt === '') { out.push({ kind: 'newline' }) }
      else if (/<a[\s\S]*?<\/a>/.test(txt)) {
        const sanitized = /href\s*=\s*['"]\s*#\s*['"]/i.test(txt)
          ? txt.replace(/href\s*=\s*['"]\s*#\s*['"]/ig, 'href="javascript:;"')
          : txt.replace(/href=/ig, 'target="_blank" href=')
        out.push({ kind: 'html', html: sanitized })
      } else { out.push({ kind: 'text', text: txt }) }
    } else if (t === '1') { out.push({ kind: 'keyword', text: txt, serial_no: mc?.serial_no || '' }) }
  })
  return out.slice(0, 20)
}

function onSmartReplyLineClick(e) {
  const a = e.target?.closest?.('a')
  if (!a) return
  const href = String(a.getAttribute('href') || '')
  if (href === '#' || href === 'javascript:;') { e.preventDefault(); e.stopPropagation() }
}

function onClickSmartMenuKeyword(text) {
  const t = String(text || '').trim()
  if (!t) return
  emit('clickMsgMeun', t)
}

defineExpose({
  scrollToBottom,
  scrollToMessage,
  resetScroll
})
</script>
