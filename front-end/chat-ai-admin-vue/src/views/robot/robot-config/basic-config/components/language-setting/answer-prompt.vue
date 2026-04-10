<style lang="less" scoped>
.setting-box {
  position: relative;
  .robot-info-box {
    display: flex;
    flex-direction: column;
    gap: 8px;
    .robot-prompt {
      display: flex;
      align-items: center;
      line-height: 22px;
      font-size: 14px;
      white-space: pre-wrap;
      word-break: break-all;
      color: #595959;
    }
  }
}
</style>

<template>
  <edit-box
    class="setting-box"
    :title="t('title_answer_generation_prompt')"
    icon-name="suggested-issues"
  >
    <div class="robot-info-box">
      <div class="robot-prompt">
        {{ t('msg_answer_generation_description') }}
      </div>
      <div>
        <a-switch
          :disabled="!isEdit"
          class="switch-item"
          checkedValue="true"
          unCheckedValue="false"
          v-model:checked="formState.tips_before_answer_switch"
          :checked-children="t('btn_on')"
          :un-checked-children="t('btn_off')"
        />
      </div>

      <div>
        <a-input
          v-model:value="formState.tips_before_answer_content"
          :readOnly="!isEdit"
          style="width: 480px"
          :maxLength="10"
          :placeholder="t('ph_input')"
          @blur="handleBlur"
        />
      </div>
    </div>
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, toRaw, computed, watch } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import EditBox from './edit-box.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.answer-prompt')

const emit = defineEmits(['save'])

const props = defineProps({
  isEdit: {
    type: Boolean,
    default: false
  },
  currentMultiLangConfig: {
    type: Object,
    default: () => {}
  },
  backCurrentMultiLangConfig: {
    type: Object,
    default: () => {}
  }
})

const formState = reactive({
  tips_before_answer_switch: 'false',
  tips_before_answer_content: '思考中，请稍后...'
})

watch(
  formState,
  () => {
    emit('save', {
      tips_before_answer_switch: formState.tips_before_answer_switch,
      tips_before_answer_content: formState.tips_before_answer_content
    })
  },
  { deep: true }
)
const handleEdit = () => {
  formState.tips_before_answer_switch = props.backCurrentMultiLangConfig.tips_before_answer_switch
  formState.tips_before_answer_content = props.backCurrentMultiLangConfig.tips_before_answer_content
}

const handleBlur = () => {
  if (formState.tips_before_answer_content == '') {
    formState.tips_before_answer_content = t('msg_thinking_please_wait')
  }
}

watch(
  () => props.currentMultiLangConfig,
  (newVal) => {
    handleEdit()
  },
  {
    deep: true,
    immediate: true
  }
)
</script>
