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
  <edit-box class="setting-box" title="问题优化" icon-name="problem-optimization">
    <template #extra>
      <span></span>
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt">
        开启后，进行知识库搜索时，会根据对话记录，利用AI补全问题缺失的信息
      </div>
    </div>

    <a-switch
      @change="handleEdit"
      class="switch-item"
      checkedValue="true"
      unCheckedValue="false"
      v-model:checked="robotInfo.enable_question_optimize"
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

const formState = reactive({
  prompt: ''
})

const onSave = () => {
  updateRobotInfo({ ...toRaw(formState) })
  isEdit.value = false
}

const handleEdit = (val) => {
  updateRobotInfo({});
}
</script>
