<template>
  <a-drawer
    :open="props.open"
    placement="right"
    :width="500"
    :closable="false"
    :maskClosable="!props.loading"
    :bodyStyle="{ padding: '24px', background: '#fff' }"
    :headerStyle="{ padding: '10px 16px', borderBottom: '1px solid #f0f0f0' }"
    @close="handleClose"
  >
    <template #title>
      <div class="drawer-header">
        <span class="drawer-title">{{ t('title_persona') }}</span>
        <button class="drawer-close" type="button" :disabled="props.loading" @click="handleClose">
          <CloseOutlined />
        </button>
      </div>
    </template>

    <div class="prompt-card">
      <div class="prompt-card-header">
        <span class="prompt-label">{{ t('label_prompt') }}</span>
        <div class="prompt-actions">
          <template v-if="editing">
            <button class="action-btn action-btn-cancel" type="button" :disabled="props.loading" @click="handleCancelEdit">
              {{ t('btn_cancel') }}
            </button>
            <button class="action-btn" type="button" :disabled="props.loading" @click="handleSave">
              <a-spin v-if="props.loading" size="small" />
              <template v-else>{{ t('btn_save') }}</template>
            </button>
          </template>
          <button v-else class="action-btn" type="button" @click="handleStartEdit">
            <EditOutlined />
            <span>{{ t('btn_edit') }}</span>
          </button>
        </div>
      </div>

      <div class="prompt-content-box" :class="{ editing }">
        <a-textarea
          v-if="editing"
          v-model:value="draftPrompt"
          class="prompt-textarea"
          :disabled="props.loading"
        />
        <div v-else class="prompt-content">{{ displayPrompt }}</div>
      </div>
    </div>
  </a-drawer>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { CloseOutlined, EditOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.clawbot.chat.components.prompt-drawer')
const emit = defineEmits(['close', 'save'])

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  },
  prompt: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const editing = ref(false)
const draftPrompt = ref('')

const displayPrompt = computed(() => {
  return props.prompt || t('msg_no_prompt')
})

watch(
  () => props.open,
  (open) => {
    if (open) {
      editing.value = false
      draftPrompt.value = props.prompt || ''
    }
  }
)

watch(
  () => props.prompt,
  (prompt) => {
    if (!editing.value) {
      draftPrompt.value = prompt || ''
    }
  }
)

const handleClose = () => {
  if (props.loading) {
    return
  }
  editing.value = false
  draftPrompt.value = props.prompt || ''
  emit('close')
}

const handleStartEdit = () => {
  draftPrompt.value = props.prompt || ''
  editing.value = true
}

const handleCancelEdit = () => {
  if (props.loading) {
    return
  }
  editing.value = false
  draftPrompt.value = props.prompt || ''
}

const handleSave = () => {
  emit('save', draftPrompt.value)
}
</script>

<style lang="less" scoped>
.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.drawer-title {
  font-size: 16px;
  line-height: 24px;
  font-weight: 600;
  color: #262626;
}

.drawer-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  padding: 0;
  border: 0;
  border-radius: 6px;
  background: transparent;
  color: #595959;
  cursor: pointer;

  &:hover {
    background: #f5f5f5;
  }

  &:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }
}

.prompt-card {
  display: flex;
  flex-direction: column;
  height: 100%;
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  background: #fff;
}

.prompt-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
}

.prompt-label {
  font-size: 16px;
  line-height: 24px;
  color: #262626;
}

.prompt-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 0;
  border: 0;
  background: transparent;
  color: #2475fc;
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;

  &:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }
}

.action-btn-cancel {
  color: #595959;
}

.prompt-content-box {
  flex: 1;
  margin: 0 16px 16px;
  
  border-radius: 8px;
  background: #f2f4f7;
  overflow: hidden;
  &.editing {
    padding: 0;
  }
}

.prompt-content {
  height: 100% !important;
  padding: 12px 16px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 16px;
  line-height: 24px;
  color: #262626;
}

.prompt-textarea {
  width: 100%;
  height: 100% !important;
  border: none !important;
  box-shadow: none !important;
  resize: none !important;
  background: #f2f4f7 !important;

  :deep(textarea) {
    min-height: 216px !important;
    padding: 12px 16px !important;
    border: none !important;
    box-shadow: none !important;
    resize: none !important;
    background: #f2f4f7 !important;
    font-size: 16px;
    line-height: 24px;
    color: #262626;
  }
}
</style>
