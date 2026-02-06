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
  <edit-box class="setting-box" :title="t('title_suggested_questions')" icon-name="suggested-issues">
    <template #extra>
      <span></span>
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt">
        {{ t('msg_guide_description', { num: robotInfo.question_guide_num }) }}
        <a-input-number
          v-model:value="robotInfo.question_guide_num"
          style="width: 80px"
          :precision="0"
          :min="1"
          :max="10"
          @blur="handleEdit"
        />
      </div>
    </div>

    <a-switch
      @change="handleEdit"
      class="switch-item"
      checkedValue="true"
      unCheckedValue="false"
      v-model:checked="robotInfo.enable_question_guide"
      :checked-children="t('btn_on')"
      :un-checked-children="t('btn_off')"
    />
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import EditBox from './edit-box.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.suggested-issues')

const isEdit = ref(false)
const { robotInfo, updateRobotInfo } = inject('robotInfo')

const handleEdit = (val) => {
  updateRobotInfo({});
}
</script>
