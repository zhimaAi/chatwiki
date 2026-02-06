<style lang="less" scoped></style>

<template>
  <SwitchBox
    class="markdown-setting"
    :title="t('title_markdown_output')"
    icon-name="jibenpeizhi"
    v-model:value="value"
    @change="onSave"
  >
    <template #extra>
      <span></span>
    </template>

    <div>{{ t('msg_markdown_description') }}</div>
  </SwitchBox>
</template>

<script setup>
import { ref, watch, inject } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import SwitchBox from './switch-box.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.markdown-setting')

const { robotInfo, updateRobotInfo } = inject('robotInfo')

const value = ref(false)

watch(
  () => robotInfo.show_type,
  (val) => {
    value.value = val == 1
  },
  { immediate: true }
)

const onSave = () => {
  updateRobotInfo({
    show_type: value.value ? 1 : 0
  })
}
</script>
