<style lang="less" scoped>
.message-list-wrapper {
  min-width: 350px;
  max-width: 800px;
  margin: 0 auto;
  height: calc(100vh - 298px);
  overflow-y: auto;
  padding: 12px 24px 24px;
}

.msg-img {
  width: auto;
  height: auto;
  max-width: 100%;
  max-height: 100%;
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
  .message-info{
    display: flex;
    align-items: center;
    justify-content: flex-end;
    line-height: 20px;
    font-size: 12px;
    color: #7a8699;
    gap: 8px;
    margin-bottom: 4px;
  }

  .user-message {

    .user-avatar-box {
      margin-left: 8px;
    }

    .item-body {
      padding: 5px 0;
    }

    .message-content {
      float: right;
      display: flex;
      border-radius: 16px 4px 16px 16px;
      margin-right: 0;
      background-color: #2475fc;
      color: #fff;
      padding: 12px;
      max-width: 100%;
      line-height: 22px;
      font-size: 14px;
      font-weight: 400;
      white-space: pre-wrap;
      width: auto;
    }
  }

  .robot-message {
    .message-info{
      justify-content: left;
    }
    .item-body {
      position: relative;
      display: inline-block;
      padding: 12px;
      margin-right: 50px;
      border-radius: 4px 16px 16px;
      background-color: #fff;
      white-space: pre-wrap;
      word-break: break-all;
      max-width: 100%;
      width: auto;
      overflow: hidden;
    }

    .message-content {
      line-height: 22px;
      font-size: 14px;
      font-weight: 400;
      color: #3a4559;
      white-space: pre-wrap;
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

    .file-items {
      .file-item {
        line-height: 22px;
        margin-top: 8px;
        font-size: 14px;
        font-size: 14px;
        cursor: pointer;
      }

      .file-name {
        color: #164799;
      }
    }
  }
}

.loading-box {
  position: absolute;
  left: 50%;
  top: 40%;
}

.empty-box {
  text-align: center;
  height: 100%;
  padding-top: 148px;
  img {
    width: 200px;
    height: 200px;
  }
  .title {
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    color: #262626;
  }
}

/* 滚动条样式 */
.message-list-wrapper::-webkit-scrollbar {
  width: 4px; /*  设置纵轴（y轴）轴滚动条 */
  height: 4px; /*  设置横轴（x轴）轴滚动条 */
}
/* 滚动条滑块（里面小方块） */
.message-list-wrapper::-webkit-scrollbar-thumb {
  border-radius: 0px;
  background: transparent;
}
/* 滚动条轨道 */
.message-list-wrapper::-webkit-scrollbar-track {
  border-radius: 0;
  background: transparent;
}

/* hover时显色 */
.message-list-wrapper:hover::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
}
.message-list-wrapper:hover::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
}

.message-top-box {
  display: flex;
  padding: 12px 0;
  justify-content: center;
  align-items: center;
  background: #F0F2F5;
  z-index: 2;
  border-radius: 6px;

  .message-top-title {
    color: #242933;
    font-family: "PingFang SC";
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 20px;
  }

  .message-top-source {
    margin-left: 4px;
    color: #7a8699;
    font-family: "PingFang SC";
    font-size: 12px;
    font-style: normal;
    font-weight: 400;
    line-height: 20px;
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

.reply-item .message-content {
  padding: 0;
  color: #1a1a1a;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
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
</style>

<template>
<div class="message-box">
  <div class="message-top-box">
    <div class="message-top-title">{{ props.robotInfo.robot_name }}</div>
    <div v-if="sessionSource" class="message-top-source">({{ formatSource(sessionSource) }})</div>
  </div>
  <div class="message-list-wrapper" ref="scrollBoxRef" @scroll="onScroll">
    <div class="message-list">
      <div class="empty-box" v-if="props.isEmpty || !props.messages">
        <img src="@/assets/img/library/detail/empty.png" alt="" />
        <div class="title">{{ t('msg_empty_result') }}</div>
      </div>
      <template v-else v-for="item in props.messages" :key="item.uid">
        <!-- 用户的消息 -->
        <div class="message-item user-message" :id="'msg-' + item.uid" v-if="item.is_customer == 1">
          <div class="itme-right">
            <div class="message-info">
              <span>{{ formatDisplayChatTime(item.create_time) }}</span>
              <span>{{ item.name }}</span>
            </div>
            <div class="item-body">
              <!-- 收到消息类型处理，目前只处理了image 后续有其他的在这里添加 -->
              <template v-if="item.received_message_type == 'image' && item.media_id_to_oss_url">
                <div class="message-content">
                  <img v-viewer class="msg-img" :src="item.media_id_to_oss_url" alt="" />
                </div>
              </template>
              <div v-else class="message-content">{{ item.content }}</div>
            </div>
          </div>
          <div class="user-avatar-box">
            <img class="user-avatar" :src="item.avatar" />
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
            <div class="message-info">
              <span>{{ formatDisplayChatTime(item.create_time) }}</span>
              <span>{{ item.name }}</span>
            </div>
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
                          <div class="reply-line">
                            <span class="line-text" v-if="line.kind === 'text'">{{ line.text }}</span>
                            <div v-else-if="line.kind === 'newline'" class="empty-line"></div>
                            <span v-else-if="line.kind === 'html'" v-html="line.html"></span>
                            <a v-else-if="line.kind === 'keyword'" href="javascript:;" class="link">
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
            <div class="item-body" v-if="!(item.msg_type == 1 && item.content == '')">
              <template v-if="item.msg_type == 1">
                <div class="message-content" v-viewer>
                  <cherry-markdown :content="item.content"></cherry-markdown>
                </div>
              </template>

              <template v-if="item.msg_type == 2">
                <div class="message-content" v-html="item.menu_json.content"></div>
              </template>

              <template v-if="item.msg_type == 3">
                <div class="message-content">
                  <img v-viewer class="msg-img" :src="item.content" alt="" />
                </div>
              </template>
            </div>
            <!-- 语音消息 -->
            <div v-if="item.voice_content && item.voice_content.length">
              <VoiceMessage :voice_content="item.voice_content" />
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
import { ref, nextTick } from 'vue'
import { formatDisplayChatTime } from '@/utils/index'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import VoiceMessage from '@/views/robot/robot-test/components/voice-message.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.session-record.components.message-list')

const emit = defineEmits(['scroll', 'scrollStart', 'scrollEnd'])

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
    default: () => false
  },
  sessionSource: {
    type: String,
    default: () => null
  },
  channelItem: {
    type: Array,
    default: () => []
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
  scrollDirection: '' // 滚动方向
}

let scrollEventTimer = null // 滚动条防抖
let onScrollEventLock = false // 时间触发锁

const formatSource = (val) => {
  let newVal
  for (let i = 0; i < props.channelItem.length; i++) {
    const item = props.channelItem[i];
    if (item.app_type === val) {
      newVal = item.app_name
    }
  }
  return newVal
}

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
  }, 100)
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

    if (element) {
      // 没数据时时没有element的
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

// 重置滚动条状态
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

defineExpose({
  scrollToBottom,
  scrollToMessage,
  resetScroll
})
</script>
