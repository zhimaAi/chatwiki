<style lang="less" scoped></style>

<template>
  <SwitchBox
    title="显示引文"
    icon-name="jibenpeizhi"
    v-model:value="value"
    @change="onSave"
  >
    <template #extra>
      <span></span>
    </template>

    <div>开启后，机器人回答时，在回答后面显示引用的知识库文档</div>
  </SwitchBox>
</template>

<script setup>
import { ref, watch, inject } from 'vue'
import SwitchBox from './switch-box.vue'

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
