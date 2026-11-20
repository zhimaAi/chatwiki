<template>
  <SettingsDialog class="workspace-prompt-dialog" :open="true" :width="1000" :aria-label="t('prompt')" :ok-text="t('save')" :cancel-text="t('cancel')" :confirm-loading="saving" :closable="!saving" :mask-closable="!saving" :keyboard="!saving" :cancel-button-props="{ disabled: saving }" @ok="save" @cancel="cancel">
    <template #closeIcon><img src="@/assets/img/customer-service/basic-settings/prompt/close.svg" alt="" /></template>
    <div class="prompt-dialog" :inert="saving ? '' : null"><PromptEditor ref="editor" :robot-info="robotInfo" :robot-list="robotList" /></div>
  </SettingsDialog>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import PromptEditor from './prompt-editor.vue'
import SettingsDialog from './settings-dialog.vue'

const props = defineProps({
  robotInfo: { type: Object, required: true },
  robotList: { type: Array, default: () => [] },
  saving: { type: Boolean, default: false }
})
const emit = defineEmits(['save', 'cancel'])
const { t } = useI18n('views.customer-service.index')
const editor = ref(null)

function save() {
  if (props.saving) return
  const patch = editor.value?.getDraft()
  if (patch) emit('save', patch)
}
function cancel() { if (!props.saving) emit('cancel') }
</script>

<style lang="less" scoped>
.prompt-dialog { min-height: 448px; }
</style>
<style lang="less">
.workspace-prompt-dialog.ant-modal {
  .ant-modal-body { max-height: calc(100dvh - 154px); }
  .ant-modal-footer { margin-top: 48px; }
  .ant-modal-close { top: 28px; }
}
</style>
