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
    .manage-link {
      color: #2475fc;
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
  <edit-box class="setting-box" :compact="compact" :title="t('title_sensitive_words')" icon-name="sensitive-icon">
    <template #extra>
      <a-button v-if="compact" size="small" @click="toManagePage"><template #icon><img src="@/assets/img/customer-service/basic-settings/arrow-right.svg" alt="" /></template>{{ t('btn_manage_sensitive_words') }}</a-button>
      <span v-else></span>
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt" :title="t('msg_sensitive_words_desc')">
        {{ t('msg_sensitive_words_desc') }}<span v-if="!compact"
          class="manage-link"
          @click="toManagePage"
          style="cursor: pointer"
          >{{ t('btn_manage_sensitive_words') }}</span
        >
      </div>
    </div>

    <template #switch><a-switch
      @change="handleEdit"
      :class="{ 'switch-item': !compact }"
      :checkedValue="1"
      :unCheckedValue="0"
      v-model:checked="robotInfo.sensitive_words_switch"
      :checked-children="t('switch_on')"
      :un-checked-children="t('switch_off')"
    /></template>
  </edit-box>
</template>
<script setup>
import { inject } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import EditBox from './edit-box.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.sensitive-words')
defineProps({ compact: { type: Boolean, default: false } })

const { robotInfo, updateRobotInfo } = inject('robotInfo')

const toManagePage = () => {
  window.open(`/#/user/sensitive-words`)
}

const handleEdit = () => {
  updateRobotInfo({})
}
</script>
