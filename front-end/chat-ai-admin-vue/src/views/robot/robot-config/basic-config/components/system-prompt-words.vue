<style lang="less" scoped>
.setting-box {
  .robot-info-box {
    .robot-prompt {
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
  <edit-box class="setting-box" title="系统提示词" icon-name="xitongtishici" v-model:isEdit="isEdit" @save="onSave"
    @edit="handleEdit">
    <div class="robot-form-box" v-show="isEdit">
      <div class="form-item">
        <a-textarea :rows="4" v-model:value="formState.prompt" placeholder="请输入提示词" allow-clear />
      </div>
    </div>
    <div class="robot-info-box" v-show="!isEdit">
      <div class="robot-prompt">{{ robotInfo.prompt }}</div>
    </div>
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, toRaw } from 'vue'
import EditBox from './edit-box.vue'

const isEdit = ref(false)
const { robotInfo, updateRobotInfo } = inject('robotInfo')

const formState = reactive({
  prompt: ''
})

const onSave = () => {
  updateRobotInfo({ ...toRaw(formState) })
  isEdit.value = false;
}

const handleEdit = () => {
  formState.prompt = robotInfo.prompt
}
</script>
