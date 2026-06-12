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
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, computed } from 'vue'
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

    let isAtTop = Math.abs(scrollOption.scrollTop) <= scrollOption.scrollStartDiff

    if (isAtTop && scrollOption.scrollDirection === 'up') {
      onScrollStart()
    }

    let isAtBottom =
      Math.abs(scrollOption.scrollHeight - scrollOption.scrollTop - scrollOption.clientHeight) <=
      scrollOption.scrollEndDiff

    if (isAtBottom && scrollOption.scrollDirection === 'down') {
      onScrollEnd()
    }
  }, 50)
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
      scrollOption.scrollTop = scrollBoxRef.value.scrollTop
      onScrollEventLock = false
    }, 50)
  })
}

function scrollToMessage(id, direction) {
  nextTick(() => {
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

function resetScroll() {
  scrollOption.scrollTop = 0
  scrollOption.scrollDirection = ''
}

defineExpose({
  scrollToBottom,
  scrollToMessage,
  resetScroll
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

  :deep(.user-avatar),
  :deep(.robot-avatar) {
    display: block;
    width: 32px;
    height: 32px;
    border-radius: 50%;
  }

  :deep(.itme-left) {
    margin-right: 8px;
  }

  :deep(.itme-right) {
    flex: 1;
    overflow: hidden;
  }
}
</style>
