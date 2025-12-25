<style lang="less" scoped>

</style>

<template>
  <div class="text-message">
    {{ message }}
  </div>
</template>

<script setup>
  const props = defineProps({
    message: String,
  })
</script>
