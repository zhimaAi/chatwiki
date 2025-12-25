<template>
  <div class="chat-message-warpper">
    <ChatMessageScroll
      ref="scrollViewRef"
      v-model:scrollTop="messageListScrollTop"
      :is-loading="messageListLoading"
      @scroll="onScroll"
      @scroll-top="onScrollTop"
      @scroll-bottom="onScrollBottom"
    >
      <div class="chat-message-content">
        <div class="is-loaded" v-if="chatMessageLoadCompleted">没用更多聊天记录了</div>
        <div
          v-for="(message, index) in messageList"
          :key="index"
          :class="['message-item', message.is_customer == 0 ? 'right' : '']"
        >
          <div class="avatar">
            <img :src="message.avatar" :alt="message.displayName" />
          </div>
          <div class="message-content">
            <div class="message-info">
              <span class="nickname">{{ message.displayName }}</span>
              <i class="gap"></i>
              <span class="time">{{ message.dispayTime }}</span>
            </div>
            <template v-if="parseReplyList(message.reply_content_list).length">
              <div class="reply-list">
                <template v-for="(rc, idx) in parseReplyList(message.reply_content_list)" :key="idx">
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
            <div class="message-bubble" v-if="!(message.msg_type == 1 && message.content == '')">
              <!-- 文本消息 -->
              <div v-if="message.msg_type == 1" class="text-content">
                <cherry-markdown
                  :content="message.content"
                  v-if="message.is_customer == 0"
                ></cherry-markdown>
                <div v-else>{{ message.content }}</div>
              </div>
              <!-- 收到消息类型处理，目前只处理了image 后续有其他的在这里添加 -->
              <div v-else-if="message.received_message_type == 'image' && message.media_id_to_oss_url" class="image-content">
                <img :src="message.media_id_to_oss_url" alt="image" />
              </div>
              <!-- 图片消息 -->
              <div v-else-if="message.msg_type == 3" class="image-content">
                <img :src="message.content" alt="image" />
              </div>
              <!-- 菜单消息 -->
              <div v-else-if="message.msg_type == 2" class="menu-content">
                <div class="menus-label">{{ message.menu_json.content }}</div>
                <div
                  class="menu-items"
                  v-if="message.menu_json.question && message.menu_json.question.length > 0"
                >
                  <div
                    v-for="(item, menuIndex) in message.menu_json.question"
                    :key="menuIndex"
                    class="menu-item"
                    @click="handleMenuClick(item)"
                  >
                    {{ item }}
                  </div>
                </div>
              </div>
              <!-- 多模态 -->
              <div v-else-if="message.msg_type == 99" class="multiple-content">
                <MultipleMessage :message="message.content" />
              </div>
              <!-- 参考文件 -->
              <div class="answer-reference-box" v-if="message.quote_file && message.quote_file.length">
                <div class="answer-reference-label">回答参考</div>
                <div>
                  <div
                    class="list-item"
                    v-for="(file, index) in message.quote_file"
                    :key="index"
                  >
                    <svg-icon name="attachment" />
                    <span @click="openLibrary(message.quote_file, file, message)">{{ file.file_name}}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </ChatMessageScroll>
  </div>
</template>

<script setup>
import { ref, nextTick, toRaw } from 'vue'
import { storeToRefs } from 'pinia'
import { useChatMonitorStore } from '@/stores/modules/chat-monitor.js'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import ChatMessageScroll from './chat-message-scroll.vue'
import MultipleMessage from '@/views/robot/robot-test/components/messages/multiple-message.vue'

const emit = defineEmits([
  'openLibrary'
])

const scrollViewRef = ref(null)

const chatMonitorStore = useChatMonitorStore()
const { getChatMessage } = chatMonitorStore
const { messageList, messageListScrollTop, messageListLoading, chatMessageLoadCompleted, activeChat } =
  storeToRefs(chatMonitorStore)

// 处理菜单点击事件
const handleMenuClick = () => {
  // console.log('Menu item clicked:', item)
}

// 打开知识库
function openLibrary(files, file, message) {
  let newfiles = toRaw(files)
  
  file.robot_key = activeChat.value.robot_key
  file.robot_id = activeChat.value.robot_id

  newfiles.forEach((item) => {
    item.message_id = message.id
    item.openid = message.openid
    item.robot_key = activeChat.value.robot_key
    item.robot_id = activeChat.value.robot_id
  })

  emit('openLibrary', newfiles, toRaw(file))
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

const onScroll = (e) => {
  const { scrollTop } = e

  messageListScrollTop.value = scrollTop
}

const onScrollTop = async () => {
  if (chatMessageLoadCompleted.value) {
    return
  }

  let oldState = scrollViewRef.value.getState()

  await getChatMessage()

  nextTick(() => {
    let newState = scrollViewRef.value.getState()

    // 如果是向上滚动，则将滚动条位置设置为新的高度减去旧的高度
    if (newState.scrollDirection == 'up') {
      messageListScrollTop.value = newState.scrollHeight - oldState.scrollHeight
    }
  })
}

const onScrollBottom = async () => {
  // console.log('onScrollBottom')
}

const scrollToBottom = () => {
  scrollViewRef.value.scrollToBottom()
}

const scrollToTop = () => {
  scrollViewRef.value.scrollToTop()
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

defineExpose({
  scrollToTop,
  scrollToBottom
})
</script>

<style lang="less" scoped>
.chat-message-warpper {
  height: 100%;
  width: 100%;
  overflow: hidden;
}
.chat-message-content {
  padding: 8px 24px 8px 24px;

  .is-loaded {
    text-align: center;
    line-height: 20px;
    font-size: 12px;
    padding-bottom: 8px;
    color: rgb(122, 134, 153);
  }
}
.message-item {
  display: flex;
  margin-bottom: 16px;
  &.right {
    flex-direction: row-reverse;
    .avatar {
      margin-right: 0;
      margin-left: 12px;
    }
    .message-info {
      flex-direction: row-reverse;
    }
  }
}
.avatar {
  width: 48px;
  height: 48px;
  margin-right: 12px;
  border-radius: 12px;
  overflow: hidden;
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}
.message-content {
  max-width: 70%;
}
.message-info {
  display: flex;
  align-items: center;
  margin-bottom: 4px;
  line-height: 20px;
  justify-content: flex-start;
  .nickname,
  .time {
    line-height: 20px;
    font-size: 12px;
    color: rgb(122, 134, 153);
  }
  .gap {
    width: 8px;
  }
}
.message-bubble {
  display: inline-block;
  border-radius: 8px;
  background-color: #fff;
  .text-content {
    padding: 12px 16px;
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    word-break: break-all;
    &::v-deep(.cherry-markdown) {
      p:last-child {
        margin-bottom: 0;
      }
    }
  }
  .image-content {
    padding: 12px 16px;
    max-width: 200px;
    img {
      width: 100%;
      border-radius: 4px;
    }
  }
  .menu-content {
    padding: 12px 16px;

    .menus-label {
      line-height: 22px;
      font-size: 14px;
      color: rgb(38, 38, 38);
      text-align: right;
    }
    .menu-items {
      margin-top: 8px;
    }
    .menu-item {
      line-height: 20px;
      padding: 6px 12px;
      margin-bottom: 8px;
      border-radius: 6px;
      font-size: 14px;
      font-weight: 400;
      cursor: pointer;
      transition: all 0.3s;
      color: rgb(22, 71, 153);
      background-color: rgba(230, 239, 255, 1);
      &:hover {
        opacity: 0.8;
      }
      &:last-child {
        margin-bottom: 0;
      }
    }
  }

  .multiple-content{
    padding: 12px 16px;
  }
}

.answer-reference-box{
  padding: 12px 16px;
  border-top: 1px solid #EDEFF2;
  .answer-reference-label{
    color: #7a8699;
    font-size: 14px;
    line-height: 22px;
    font-weight: 400;
  }
  .list-item{
    cursor: pointer;
    margin-top: 8px;
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
    color: #164799;
  }
}

.reply-list {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
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
  width: 100%;
  border-radius: 16px 4px 16px 16px;
  padding: 16px 12px;
  background: #fff;
}

.reply-text .message-content {
  max-width: 100%;
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