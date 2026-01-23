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
  <edit-box class="setting-box" title="答案生成中提示语" icon-name="suggested-issues">
    <template #extra>
      <span></span>
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt">
        开启后，WebAPP、嵌入网站对外服务中，答案生成前会显示提示语
        <a-input
          v-model:value="robotInfo.tips_before_answer_content"
          style="width: 280px"
          :maxLength="10"
          placeholder="请输入"
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
  updateRobotInfo({})
}

const handleBlur = () => {
  console.log(robotInfo.tips_before_answer_content)
  if (robotInfo.tips_before_answer_content == '') {
    robotInfo.tips_before_answer_content = '思考中、请稍等'
  }
  handleEdit()
}
</script>
