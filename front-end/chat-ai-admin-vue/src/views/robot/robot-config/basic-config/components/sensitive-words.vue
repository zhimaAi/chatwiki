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
  <edit-box class="setting-box" :title="t('title_sensitive_words')" icon-name="sensitive-icon">
    <template #extra>
      <span></span>
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt">
        {{ t('msg_sensitive_words_desc') }}<span
          @click="toManagePage"
          style="cursor: pointer"
          >{{ t('btn_manage_sensitive_words') }}</span
        >
      </div>
    </div>

    <a-switch
      @change="handleEdit"
      class="switch-item"
      :checkedValue="1"
      :unCheckedValue="0"
      v-model:checked="robotInfo.sensitive_words_switch"
      :checked-children="t('switch_on')"
      :un-checked-children="t('switch_off')"
    />
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import EditBox from './edit-box.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.sensitive-words')

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
