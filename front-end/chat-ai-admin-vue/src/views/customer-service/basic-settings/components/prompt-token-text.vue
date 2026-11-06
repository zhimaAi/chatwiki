<template>
  <template v-for="(part, index) in parts" :key="index">
    <span v-if="part.token" class="prompt-token" :class="{ 'card-token': part.token.type === 'card' }"><svg-icon v-if="part.token.type === 'card'" name="applet" />{{ part.token.label }}</span><template v-else>{{ part.text }}</template>
  </template>
</template>
<script setup>
import { computed } from 'vue'
import { tokenizePrompt } from '../prompt-tokens'
const props = defineProps({ value: { type: String, default: '' }, options: { type: Array, default: () => [] } })
const parts = computed(() => tokenizePrompt(props.value, props.options))
</script>
<style lang="less" scoped>
.prompt-token { display: inline-flex; align-items: center; min-height: 24px; padding: 0 8px; margin-right: 6px; border: 1px solid #d9d9d9; border-radius: 6px; background: #fff; color: #222; line-height: 22px; vertical-align: baseline; white-space: normal; }
.card-token { gap: 4px; min-height: 22px; border: 0; background: #e5efff; color: #164799; font-size: 12px; line-height: 20px; .svg-icon { font-size: 12px; } }
</style>
