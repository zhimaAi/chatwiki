<style lang="less" scoped>
  .actions-box {
    display: flex;
    align-items: center;
    line-height: 22px;
    font-size: 14px;
    color: #595959;

    .action-btn {
      cursor: pointer;
    }

    .save-btn {
      color: #2475fc;
    }

    .model-name {
      font-size: 14px;
      line-height: 22px;
      color: #8c8c8c;
    }
  }
</style>

<template>
  <edit-box
    class="setting-box"
    title="聊天模式"
    icon-name="chat-mode"
    v-model:isEdit="isEdit"
    :bodyStyle="{ padding: 0 }"
  >
    <template #tip>
      <a-tooltip placement="top" :overlayInnerStyle="{ width: '400px' }">
        <template #title>
          <span
            >仅知识库模式：用户提问时，从知识库检索文档，大语言模型（LLM）根据检索出来的文档分段进行回复。如果没有符合的分段，则不由LLM回复，直接回复未知问题提示语。</span
          ><br />
          <span>直连模式：用户提问时，直接由LLM生成答案，不从关联知识库中检索。</span>
        </template>
        <QuestionCircleOutlined />
      </a-tooltip>
    </template>
    <template #extra>
      <div class="actions-box">
        <span class="model-name">{{
        robotInfo.chat_type == '1' ? '仅知识库' : robotInfo.chat_type == '2' ? '直连模式' : '混合模式'
        }}</span>
        <a-divider type="vertical" />
        <a-button size="small" @click="handleEdit(true)">修改</a-button>
      </div>
    </template>
    <div class="form-box">
        <ChatEditAlert ref="ChatEditAlertRef" @save="onSave"></ChatEditAlert>
    </div>
  </edit-box>
</template>

<script setup>
import { ref, reactive, inject, toRaw } from 'vue'
const isEdit = ref(false)
import EditBox from '../edit-box.vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import ChatEditAlert from './chat-edit-alert.vue'
const ChatEditAlertRef = ref(null)

const { robotInfo, updateRobotInfo } = inject('robotInfo')

const formState = reactive({
  chat_type: robotInfo.chat_type,
  library_qa_direct_reply_score: robotInfo.library_qa_direct_reply_score,
  library_qa_direct_reply_switch: robotInfo.library_qa_direct_reply_switch,
  mixture_qa_direct_reply_score: robotInfo.mixture_qa_direct_reply_score,
  mixture_qa_direct_reply_switch: robotInfo.mixture_qa_direct_reply_switch
})

const onSave = (data) => {
    formState.chat_type = data.chat_type
    formState.library_qa_direct_reply_score = data.library_qa_direct_reply_score
    formState.library_qa_direct_reply_switch = data.library_qa_direct_reply_switch
    formState.mixture_qa_direct_reply_score = data.mixture_qa_direct_reply_score
    formState.mixture_qa_direct_reply_switch = data.mixture_qa_direct_reply_switch

    handleSave()
}

const handleSave = () => {

    updateRobotInfo({ ...toRaw(formState) })
    isEdit.value = false

}

const handleEdit = () => {
    ChatEditAlertRef.value.open(toRaw(formState))
    isEdit.value = true
}
</script>
