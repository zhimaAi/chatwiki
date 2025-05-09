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
  <edit-box class="setting-box" title="用户问题建议" icon-name="suggested-issues">
    <template #extra>
      <span></span>
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt">
        开启后，机器人回答问题后，自动根据对话内容提供
        <a-input-number
          v-model:value="robotInfo.question_guide_num"
          style="width: 80px"
          :precision="0"
          :min="1"
          :max="10"
          @blur="handleEdit"
        />
        条用户提问建议，最大不超过10条。
      </div>
    </div>

    <a-switch
      @change="handleEdit"
      class="switch-item"
      checkedValue="true"
      unCheckedValue="false"
      v-model:checked="robotInfo.enable_question_guide"
      checked-children="开"
      un-checked-children="关"
    />
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, toRaw } from 'vue'
import EditBox from './edit-box.vue'

const isEdit = ref(false)
const { robotInfo, updateRobotInfo } = inject('robotInfo')

const handleEdit = (val) => {
  updateRobotInfo({});
}
</script>
