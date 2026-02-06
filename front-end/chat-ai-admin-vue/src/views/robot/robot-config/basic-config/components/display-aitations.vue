<style lang="less" scoped></style>

<template>
  <SwitchBox
    :title="t('title_show_citations')"
    icon-name="jibenpeizhi"
    v-model:value="value"
    @change="onSave"
  >
    <template #extra>
      <span></span>
    </template>

    <div>{{ t('desc_show_citations') }}</div>
  </SwitchBox>
</template>

<script setup>
import { ref, watch, inject } from 'vue'
import SwitchBox from './switch-box.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.basic-config.components.display-aitations')

const { robotInfo, updateRobotInfo } = inject('robotInfo')

const value = ref(false)

watch(
  () => robotInfo.answer_source_switch,
  (val) => {
    value.value = val == 'true'
  },
  { immediate: true }
)

const onSave = () => {
  updateRobotInfo({
    answer_source_switch: value.value ? 'true' : 'false'
  })
}
</script>
