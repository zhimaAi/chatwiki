<style lang="less" scoped>
.setting-box {
  position: relative;
  .robot-info-box {
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
  .switch-item {
    position: absolute;
    right: 16px;
    top: calc(50% - 8px);
  }
}
</style>

<template>
  <edit-box class="setting-box" :title="t('title_answer_generation_prompt')" icon-name="suggested-issues">
    <template #extra>
      <span></span>
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt">
        {{ t('msg_answer_generation_description') }}
        <a-input
          v-model:value="robotInfo.tips_before_answer_content"
          style="width: 280px"
          :maxLength="10"
          :placeholder="t('ph_input')"
          @blur="handleBlur"
        />
      </div>
    </div>

    <a-switch
      @change="handleEdit"
      class="switch-item"
      checkedValue="true"
      unCheckedValue="false"
      v-model:checked="robotInfo.tips_before_answer_switch"
      :checked-children="t('btn_on')"
      :un-checked-children="t('btn_off')"
    />
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import EditBox from './edit-box.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.answer-prompt')

const isEdit = ref(false)
const { robotInfo, updateRobotInfo } = inject('robotInfo')

const handleEdit = (val) => {
  updateRobotInfo({})
}

const handleBlur = () => {
  console.log(robotInfo.tips_before_answer_content)
  if (robotInfo.tips_before_answer_content == '') {
    robotInfo.tips_before_answer_content = t('msg_thinking_please_wait')
  }
  handleEdit()
}
</script>
