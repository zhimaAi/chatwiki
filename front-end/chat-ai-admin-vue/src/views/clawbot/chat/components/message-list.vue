<template>
  <div class="message-list-wrapper">
    <div class="scroll-box" ref="scrollBoxRef" @scroll="onScroll">
      <div class="message-list">
        <template v-for="item in visibleMessages" :key="item.uid">
          <UserMessageItem v-if="item.is_customer == 1" :item="item" />
          <RobotMessageItem
            v-else
            :item="item"
            :robot-info="robotInfo"
            :tips-before-answer-content="tips_before_answer_content"
            :tips-before-answer-switch="tips_before_answer_switch"
            :is-show-quote-file-progress="isShowQuoteFileProgress"
            @clickMsgMeun="onClickMeun"
          />
        </template>
        <div class="reply-space" v-if="props.reserveReplySpace" ref="replySpaceRef"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, computed, watch } from 'vue'
import { useRobotStore } from '@/stores/modules/robot'
import UserMessageItem from './messages/user-message-item.vue'
import RobotMessageItem from './messages/robot-message-item.vue'

const robotStore = useRobotStore()


const currentLangConfig = computed(()=>{
  let currentLang = localStorage.getItem('lang') ? JSON.parse(localStorage.getItem('lang'))?.value : 'zh-CN' || 'zh-CN'
  let multi_lang_configs = robotStore.robotInfo.multi_lang_configs
  let configs = multi_lang_configs.find(item => item.lang_key === currentLang) || multi_lang_configs[0]
  return configs
})
const tips_before_answer_content = computed(()=>{
  return currentLangConfig.value?.tips_before_answer_content
})

const tips_before_answer_switch = computed(()=>{
  return currentLangConfig.value?.tips_before_answer_switch == 'true'
})

const emit = defineEmits([
  'clickMsgMeun',
  'scroll',
  'scrollStart',
  'scrollEnd'
])

const props = defineProps({
  messages: {
    type: Array,
    default: () => []
  },
  robotInfo: {
    type: Object,
    default: () => {}
  },
  reserveReplySpace: {
    type: Boolean,
    default: false
  }
})

const isTruthyHidden = (value) => {
  if (value === true || value === 1) {
    return true
  }
  if (typeof value === 'string') {
    const normalized = value.trim().toLowerCase()
    return normalized === 'true' || normalized === '1'
  }
  return false
}

const isMessageHidden = (item) => {
  if (!item) {
    return false
  }
  if (item.visible === false) {
    return true
  }
  return isTruthyHidden(item.hide) || isTruthyHidden(item.hidden) || isTruthyHidden(item.is_hidden)
}

const visibleMessages = computed(() => {
  return (props.messages || []).filter((item) => !isMessageHidden(item))
})

const isShowQuoteFileProgress = computed(() => {
  return (props.robotInfo.chat_type == 1 || props.robotInfo.chat_type == 3) && props.robotInfo.application_type == '0'
})

const onClickMeun = (item) => {
  emit('clickMsgMeun', item)
}

const scrollBoxRef = ref(null)
const replySpaceRef = ref(null)
const scrollOption = {
  scrollTop: 0,
  scrollHeight: 0,
  clientHeight: 0,
  scrollStartDiff: 60,
  scrollEndDiff: 60,
  scrollDirection: '',
  isAtBottom: true,
  isReplySpaceVisible: false
}

let scrollEventTimer = null
let onScrollEventLock = false

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

    updateScrollOption(e.target)
    emit('scroll', { ...scrollOption })

    let isAtTop = Math.abs(scrollOption.scrollTop) <= scrollOption.scrollStartDiff

    if (isAtTop && scrollOption.scrollDirection === 'up') {
      onScrollStart()
    }

    if (scrollOption.isAtBottom && scrollOption.scrollDirection === 'down') {
      onScrollEnd()
    }
  }, 50)
}

function updateScrollOption(scroller) {
  scrollOption.scrollTop = scroller.scrollTop
  scrollOption.scrollHeight = scroller.scrollHeight
  scrollOption.clientHeight = scroller.clientHeight
  scrollOption.isAtBottom =
    Math.abs(scrollOption.scrollHeight - scrollOption.scrollTop - scrollOption.clientHeight) <=
    scrollOption.scrollEndDiff

  if (!replySpaceRef.value) {
    scrollOption.isReplySpaceVisible = false
    return
  }

  let viewportTop = scroller.scrollTop
  let viewportBottom = viewportTop + scroller.clientHeight
  let replySpaceTop = replySpaceRef.value.offsetTop
  let replySpaceBottom = replySpaceTop + replySpaceRef.value.offsetHeight

  scrollOption.isReplySpaceVisible =
    viewportBottom > replySpaceTop && viewportTop < replySpaceBottom
}

function emitScrollState() {
  if (!scrollBoxRef.value) {
    return
  }
  updateScrollOption(scrollBoxRef.value)
  emit('scroll', { ...scrollOption })
}

function onScrollStart() {
  if (visibleMessages.value.length == 0) {
    return
  }
  emit('scrollStart', {
    msg: visibleMessages.value[0]
  })
}

function onScrollEnd() {
  if (visibleMessages.value.length == 0) {
    return
  }
  emit('scrollEnd', {
    msg: visibleMessages.value[visibleMessages.value.length - 1]
  })
}

const scrollToBottom = () => {
  if (!scrollBoxRef.value) {
    return
  }
  nextTick(() => {
    onScrollEventLock = true

    scrollBoxRef.value.scrollTop = scrollBoxRef.value.scrollHeight + 1
    setTimeout(() => {
      emitScrollState()
      onScrollEventLock = false
    }, 50)
  })
}

function scrollToMessage(id, direction) {
  if (!scrollBoxRef.value) {
    return
  }
  nextTick(() => {
    onScrollEventLock = true

    if (!direction) {
      direction = 'top'
    }

    let scroller = scrollBoxRef.value
    let element = document.querySelector('#msg-' + id)
    if (!element) {
      onScrollEventLock = false
      return
    }

    if (direction == 'top') {
      scroller.scrollTop = Math.max(element.offsetTop - 4, 0)
    } else {
      scroller.scrollTop = element.offsetTop - scroller.clientHeight + element.clientHeight
    }

    setTimeout(() => {
      emitScrollState()
      onScrollEventLock = false
    }, 50)
  })
}

function resetScroll() {
  scrollOption.scrollTop = 0
  scrollOption.scrollDirection = ''
  scrollOption.isAtBottom = true
  emit('scroll', { ...scrollOption })
}

const isAtBottom = () => scrollOption.isAtBottom

watch(
  () => [props.messages, props.reserveReplySpace],
  () => {
    nextTick(emitScrollState)
  },
  { deep: true }
)

defineExpose({
  scrollToBottom,
  scrollToMessage,
  resetScroll,
  isAtBottom
})
</script>

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

.message-list {
  :deep(.message-item) {
    display: flex;
    margin-top: 24px;
  }

  :deep(.itme-left) {
    display: none;
  }

  :deep(.itme-right) {
    flex: 1;
    overflow: hidden;
  }

  :deep(.user-message) {
    justify-content: flex-end;
  }

  :deep(.user-message .itme-right) {
    flex: 0 1 auto;
    max-width: min(72%, 520px);
  }

  :deep(.robot-message) {
    justify-content: flex-start;
  }
}

.reply-space {
  height: 65vh;
  min-height: 320px;
  max-height: 560px;
  pointer-events: none;
}
</style>
