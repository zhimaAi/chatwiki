<style lang="less" scoped>

</style>

<template>
  <div class="text-message">
    {{ message }}
  </div>
</template>

<script setup>
  defineProps({
    message: String,
  })
</script>
