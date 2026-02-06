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

.setting-info-block {
  padding: 16px;
  padding-top: 0;
  display: flex;
  flex-wrap: wrap;
  gap: 12px 16px;
  color: #595959;
  line-height: 22px;
  .set-item {
    display: flex;
    align-items: center;
  }
}
</style>

<template>
  <edit-box
    class="setting-box"
    :title="t('title_chat_mode')"
    icon-name="chat-mode"
    v-model:isEdit="isEdit"
    :bodyStyle="{ padding: 0 }"
  >
    <template #tip>
      <a-tooltip placement="top" :overlayInnerStyle="{ width: '400px' }">
        <template #title>
          <span>{{ t('tip_knowledge_base_only') }}</span
          ><br />
          <span>{{ t('tip_direct_mode') }}</span>
        </template>
        <QuestionCircleOutlined />
      </a-tooltip>
    </template>
    <template #extra>
      <div class="actions-box">
        <a-button size="small" @click="handleEdit(true)">{{ t('btn_edit') }}</a-button>
      </div>
    </template>
    <div class="setting-info-block">
      <div class="set-item">
        {{ t('label_chat_mode') }}
        <span v-if="robotInfo.chat_type == '1'">{{ t('chat_mode_knowledge_base_only') }}</span>
        <span v-if="robotInfo.chat_type == '2'">{{ t('chat_mode_direct') }}</span>
        <span v-if="robotInfo.chat_type == '3'">{{ t('chat_mode_mixed') }}</span>
      </div>
      <template v-if="robotInfo.chat_type == '1'">
        <div class="set-item">
          {{ t('label_direct_reply_qa') }}
          <span>{{ robotInfo.library_qa_direct_reply_switch == 'true' ? t('status_on') : t('status_off') }}</span>
        </div>
        <div class="set-item" v-if="robotInfo.library_qa_direct_reply_switch == 'true'">
          {{ t('label_similarity_threshold') }}
          <span>{{ robotInfo.library_qa_direct_reply_score }}</span>
        </div>
      </template>
      <template v-if="robotInfo.chat_type == '3'">
        <div class="set-item">
          {{ t('label_direct_reply_qa') }}
          <span>{{ robotInfo.mixture_qa_direct_reply_switch == 'true' ? t('status_on') : t('status_off') }}</span>
        </div>
        <div class="set-item" v-if="robotInfo.mixture_qa_direct_reply_switch == 'true'">
          {{ t('label_similarity_threshold') }}
          <span>{{ robotInfo.mixture_qa_direct_reply_score }}</span>
        </div>
      </template>
    </div>
    <div class="form-box">
      <ChatEditAlert ref="ChatEditAlertRef" @save="onSave"></ChatEditAlert>
    </div>
  </edit-box>
</template>

<script setup>
import { ref, reactive, inject, toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.basic-config.components.chat-mode.index')
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
