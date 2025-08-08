<style lang="less" scoped>
.message-input-box {
  width: 100%;
  margin-bottom: 0px;

  .message-input-body {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .message-input {
    position: relative;
    z-index: 2;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex: 1;
    min-height: 60px;
    max-width: 738px;
    min-width: 350px;
    border-radius: 16px;
    border: 1px solid #ddd;
    background: #fff;
    box-shadow: 0 2px 6px 0 rgba(0, 0, 0, 0.08);
    transition: all 0.2s;
    padding: 10px 0;
    margin: 0 12px;

    .text-input {
      max-height: 10em;
      line-height: 1.2em;
      height: 1.5em;
      overflow: hidden;
      white-space: pre-wrap; /* 保持内容的换行，并允许自动换行 */
      resize: none;
      border: none;
      width: 100%;
      margin: 0 45px 0 12px;
      color: rgb(26, 26, 26);
      font-size: 16px;
      font-weight: 400;
      background: none;

      transition: height 0.1s ease-in-out;;

      &::placeholder {
        font-size: 16px;
        font-weight: 400;
        color: rgb(191, 191, 191);
      }
    }

    /* 滚动条样式 */
    .text-input::-webkit-scrollbar {
        width: 4px; /*  设置纵轴（y轴）轴滚动条 */
        height: 4px; /*  设置横轴（x轴）轴滚动条 */
    }
    /* 滚动条滑块（里面小方块） */
    .text-input::-webkit-scrollbar-thumb {
        cursor: pointer;
        border-radius: 5px;
        background: transparent;
    }
    /* 滚动条轨道 */
    .text-input::-webkit-scrollbar-track {
        border-radius: 5px;
        background: transparent;
    }

    /* hover时显色 */
    .text-input:hover::-webkit-scrollbar-thumb {
        background: rgba(0,0,0,0.2);
    }
    .text-input:hover::-webkit-scrollbar-track {
        background: rgba(0,0,0,0.1);
    }

    .send-btn {
      position: absolute;
      bottom: 13px;
      right: 10px;
      font-size: 32px;
      color: #b3b3b3;
      transition: all 0.2s;
    }

    .send-btn-active {
      color: #2475fc;
      cursor: pointer;
    }

    &.is-set {
      border: 1px solid #2475fc;
    }
  }
}
</style>

<template>
  <div class="message-input-box">
    <div class="message-input-body">
      <div class="message-input" :class="{ 'is-set': valueText }">
        <textarea
          ref="messageTextarea"
          :style="{height: inputHeight}"
          class="text-input"
          v-model="valueText"
          placeholder="在此输入您想了解的内容"
          @change="onChange"
          @input="onInput"
          @keyup.enter="handleEnter" 
          @focus="onFocus"
          @blur="onBlur"
        ></textarea>
        <svg-icon :name="isSendBtn ? 'paper-airplane-new-active' : 'paper-airplane-new'" :class="{'send-btn-active': isSendBtn}" class="send-btn" @click="sendMessage" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import calcTextareaHeight from '@/utils/calcTextareaHeight'
import { useChatStore } from '@/stores/modules/chat'
import { useUserStore } from '@/stores/modules/user'
import { showToast } from 'vant'

const chatStore = useChatStore()
const userStore = useUserStore()
const { robot } = chatStore
import { checkChatRequestPermission } from '@/api/robot/index'

const emit = defineEmits(['update:value', 'send', 'showLogin' ])

const isSendBtn = ref(false)
const valueText = ref("")
const inputHeight = ref("1.5em")
const messageTextarea = ref(null)

const onChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  emit('update:value', target.value)
}

const onInput = (event) => {
  // 会有延迟
  if (valueText.value.trim()) {
    isSendBtn.value = true
  } else {
    isSendBtn.value = false
  }
  inputHeight.value = calcTextareaHeight(messageTextarea.value).height // 调整高度
  if (parseInt(inputHeight.value) >= 160) {
    event.target.style.overflow = 'auto'
  } else {
    event.target.style.overflow = 'hidden'
  }
}

const sendMessage = async () => {
  //检查是否含有敏感词
  let result = await checkChatRequestPermission({
    robot_key: robot.robot_key,
    openid: robot.openid,
    question: valueText.value,
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
  if(result.data && result.data.words){
    return showToast(`提交的内容包含敏感词：[${result.data.words.join(';')}] 请修改后再提交`)
  }
  if (valueText.value.trim()) {
    isSendBtn.value = false
    emit('send', valueText.value)
    valueText.value = ""
    inputHeight.value = '1.5em'
  }
}

const handleEnter = (e: KeyboardEvent) => {
  e.preventDefault()
  const target = e.target as HTMLInputElement;
  if(e.shiftKey){
    inputHeight.value = target.scrollHeight + 'px';
    return 
  }
  if (!valueText.value.trim()) {
    return
  }
  sendMessage()
}

const onFocus = (event) => {
  event.target.parentNode.style.borderColor = "#2475FC"
}

const onBlur = (event) => {
  event.target.parentNode.style.borderColor = "#DDD"
}
const handleSetValue = (data: string) => {
  valueText.value = data;
}
defineExpose({
  handleSetValue,
  sendMessage
})
</script>
