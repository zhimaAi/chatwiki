<template><SystemPromptWords ref="editor" dialog-mode :robot-list="robotList" /></template>
<script setup>
import { provide, reactive, ref } from 'vue'
import SystemPromptWords from '@/views/robot/robot-config/basic-config/components/system-prompt-words.vue'
import { cloneConfig } from '../../config-session'

const props = defineProps({ robotInfo: { type: Object, required: true }, robotList: { type: Array, default: () => [] } })
const editor = ref(null)
const draft = reactive(cloneConfig(props.robotInfo))
provide('robotInfo', { robotInfo: draft, updateRobotInfo: () => Promise.resolve(false) })
defineExpose({ getDraft: () => editor.value?.getDraft() })
</script>
