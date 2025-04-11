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
  <edit-box class="setting-box" title="敏感词" icon-name="sensitive-icon">
    <template #extra>
      <span></span>
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt">
        开启后，访客提交的问题中如果包含敏感词，将不支持提交，<span
          @click="toManagePage"
          style="cursor: pointer"
          >设置敏词 ></span
        >
      </div>
    </div>

    <a-switch
      @change="handleEdit"
      class="switch-item"
      :checkedValue="1"
      :unCheckedValue="0"
      v-model:checked="robotInfo.sensitive_words_switch"
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

const toManagePage = () => {
  window.open(`/#/user/sensitive-words`)
}

const handleEdit = (val) => {
  updateRobotInfo({})
}
</script>
