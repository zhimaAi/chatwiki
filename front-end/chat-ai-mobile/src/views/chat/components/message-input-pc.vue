<style lang="less" scoped>
.message-input-box {
  display: flex;
  align-items: center;
  justify-content: center;
}
.message-input-wrapper {
  position: relative;
  width: 760px;
  .message-action {
    position: absolute;
    right: 16px;
    bottom: 12px;
  }

  .send-msg-btn {
    width: 60px;
    height: 32px;
    line-height: 32px;
    padding: 0;
    font-size: 14px;
    font-weight: 400;
    border-radius: 4px;
    border: none;
    cursor: pointer;
    transition: all 0.2s;
    color: #2475fc;
    background-color: #e6efff;

    &:hover {
      opacity: 0.7;
    }
  }

  .loading-action {
    display: block;
    height: 32px;
    line-height: 32px;
    padding-top: 6px;
    text-align: center;
  }
}
</style>

<template>
  <div class="message-input-box">
    <div class="message-input-wrapper">
      <a-textarea
        v-model:value="valueText"
        :auto-size="{ minRows: 5, maxRows: 5 }"
        placeholder="在此输入您想了解的内容，Shift+Enter换行"
        @change="onChange"
        @keydown="handleKeydown"
      />
      <div class="message-action">
        <span class="loading-action">
          <a-spin :spinning="props.loading"></a-spin>
        </span>

        <button class="send-msg-btn" @click="sendMessage" v-if="!props.loading">
          <span>发送</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useChatStore } from '@/stores/modules/chat'
import { useUserStore } from '@/stores/modules/user'
import { Textarea as ATextarea, Spin as ASpin } from 'ant-design-vue'
import { showToast } from 'vant'

const chatStore = useChatStore()
const userStore = useUserStore()
const { robot } = chatStore
import { checkChatRequestPermission } from '@/api/robot/index'

const emit = defineEmits(['update:value', 'send', 'showLogin'])

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  }
})

const isSendBtn = ref(false)
const valueText = ref('')

const onChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('update:value', target.value)
}

const sendMessage = async () => {
  if (props.loading) {
    return
  }
  //检查是否含有敏感词
  let result = await checkChatRequestPermission({
    robot_key: robot.robot_key,
    openid: robot.openid,
    question: valueText.value
  })
  // 未登录
  if (result.data?.code == 10002) {
    // 弹出登录
    showToast('请登录账号')
    emit('showLogin')
    userStore.setLoginStatus(false)
    return false
  }
  // 有权限的账号登录后才可访问
  if (result.data?.code == 10003) {
    // 弹出登录
    showToast('当前账号无访问权限')
    emit('showLogin')
    userStore.setLoginStatus(false)
    return false
  }
  userStore.setLoginStatus(true)
  if (result.data && result.data.words) {
    return showToast(`提交的内容包含敏感词：[${result.data.words.join(';')}] 请修改后再提交`)
  }
  if (valueText.value.trim()) {
    isSendBtn.value = false
    emit('send', valueText.value)
    valueText.value = ''
  }
}

const handleKeydown = (e: KeyboardEvent) => {
  // const target = e.target as HTMLInputElement
  if (e.key === 'Enter' && !e.shiftKey) {
    if (!valueText.value.trim()) {
      return
    }
    e.preventDefault()
    sendMessage()
  } else if (e.key === 'Enter' && e.shiftKey) {
    e.preventDefault()
    valueText.value += '\n'
  }
}

const handleSetValue = (data: string) => {
  valueText.value = data
}
defineExpose({
  handleSetValue,
  sendMessage
})
</script>
