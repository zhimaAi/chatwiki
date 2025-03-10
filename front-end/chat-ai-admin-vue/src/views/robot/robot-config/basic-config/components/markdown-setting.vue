<style lang="less" scoped></style>

<template>
  <SwitchBox
    class="markdown-setting"
    title="Markdown格式输出"
    icon-name="jibenpeizhi"
    v-model:value="value"
    @change="onSave"
  >
    <template #extra>
      <span></span>
    </template>

    <div>开启后，在使用H5链接、Web、PC客户端时,可以要求机器人以Markdown格式输出回答</div>
  </SwitchBox>
</template>

<script setup>
import { ref, watch, inject } from 'vue'
import SwitchBox from './switch-box.vue'

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
