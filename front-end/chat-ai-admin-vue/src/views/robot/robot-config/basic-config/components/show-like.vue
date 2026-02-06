<style lang="less" scoped>
.setting-box {
  position: relative;
  .robot-info-box {
    .robot-prompt {
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
  <edit-box class="setting-box" :title="t('title_show_like_dislike')" icon-name="like-icon">
    <template #extra>
      <span></span>
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt">{{ t('msg_show_like_dislike_desc') }}</div>
    </div>

    <a-switch
      @change="handleEdit"
      class="switch-item"
      :checkedValue="1"
      :unCheckedValue="0"
      v-model:checked="robotInfo.feedback_switch"
      :checked-children="t('btn_on')"
      :un-checked-children="t('btn_off')"
    />
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import EditBox from './edit-box.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.show-like')

const { robotInfo, updateRobotInfo } = inject('robotInfo')

const handleEdit = (val) => {
  updateRobotInfo({})
}
</script>
