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
          >
            <div class="itme-left">
              <img class="user-avatar" :src="item.avatar" />
            </div>
            <div class="itme-right">
              <div class="item-body">
                <div class="message-content">{{ item.content }}</div>
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
              <div class="item-body">
                <template v-if="item.msg_type == 1">
                  <div class="message-content">
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
                  <div class="message-content" v-html="item.menu_json.content"></div>
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
                    <img class="msg-img" :src="item.content" alt="" />
                  </div>
                </template>

                <div
                  class="message-action-wrap"
                  v-if="showQuoteFileBox(item)"
                  v-show="!item.question"
                >
                  <div class="message-action">
                    <div class="action-btn" v-if="item.quote_file && item.quote_file.length > 0">
                      <span>答案来源&nbsp;</span>
                      <a-tooltip placement="top" :overlayInnerStyle="{ width: '400px' }">
                        <template #title>
                          <span
                            >您可以修改相应文档并重新上传，以调整机器人的回答效果。答案来源仅后台测试时可见。</span
                          >
                        </template>
                        <QuestionCircleOutlined />
                      </a-tooltip>
                    </div>
                    <div class="action-btn" v-if="item.debug && item.debug.length > 0">
                      <span><a @click="openPromptLog(item)">Prompt 日志</a></span>
                    </div>
                  </div>

                  <div class="file-items" v-if="item.quote_file && item.quote_file.length > 0">
                    <div
                      class="file-item"
                      v-for="file in item.quote_file"
                      :key="file.id"
                      @click="openLibrary(item.quote_file, file, item.id)"
                    >
                      <a class="file-name">{{ file.file_name }}</a>
                    </div>
                  </div>
                </div>
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
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import GuessYouWant from './guess-you-want.vue'

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

const common_question_list = computed(() => {
  if (props.robotInfo.common_question_list.length) {
    return JSON.parse(props.robotInfo.common_question_list)
  }
  return []
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

  newfiles.forEach((item) => {
    item.message_id = message_id
  })

  emit('openLibrary', newfiles, toRaw(file))
}

defineExpose({
  scrollToBottom,
  scrollToMessage,
  resetScroll
})
</script>
