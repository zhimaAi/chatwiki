<template>
  <svg class="sparkle-icon" viewBox="0 0 24 24" aria-hidden="true">
    <path d="M12 3l2.5 6.5L21 12l-6.5 2.5L12 21l-2.5-6.5L3 12l6.5-2.5Z"
      :fill="filled ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" />
    <path v-if="!filled" d="M20 2v4M18 4h4" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
  </svg>
</template>
<script setup>
defineProps({ filled: Boolean })
</script>
<style scoped>
.sparkle-icon { width: 18px; height: 18px; flex-shrink: 0; color: #2475fc; vertical-align: middle; }
</style>
