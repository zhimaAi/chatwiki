<template>
  <a-style-provider hash-priority="high" :transformers="[legacyLogicalPropertiesTransformer]">
    <a-app>
      <GlobalConfigProvider>
        <RouterView />
      </GlobalConfigProvider>
    </a-app>
  </a-style-provider>
</template>

<script setup>
import GlobalConfigProvider from '@/components/global-config-provider/index.vue'
import { RouterView } from 'vue-router'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import { legacyLogicalPropertiesTransformer } from 'ant-design-vue';

dayjs.locale('zh-cn')
</script>

<style scoped></style>
