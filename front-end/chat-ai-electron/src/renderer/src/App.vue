<template>
  <GlobalConfigProvider>
    <RouterView />
  </GlobalConfigProvider>
</template>

<script setup>
import GlobalConfigProvider from '@/components/global-config-provider/index.vue'
import { RouterView } from 'vue-router'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'

dayjs.locale('zh-cn')
</script>

<style scoped></style>
