<script setup lang="ts">
import { RouterView } from 'vue-router'
import { reactive } from 'vue'
const themeVars = reactive({
  popupRoundRadius: '6px',
  popupCloseIconSize: '18px'
})
</script>

<template>
  <van-config-provider theme="light" :theme-vars="themeVars">
    <RouterView />
  </van-config-provider>
</template>

<style scoped></style>
