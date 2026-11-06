<template>
  <a-style-provider hash-priority="high" :transformers="[legacyLogicalPropertiesTransformer]">
    <GlobalConfigProvider>
      <a-app>
        <ThemeVariables>
          <RouterView />
        </ThemeVariables>
      </a-app>
    </GlobalConfigProvider>
  </a-style-provider>
</template>

<script setup>
import GlobalConfigProvider from '@/components/global-config-provider/index.vue'
import ThemeVariables from '@/components/global-config-provider/theme-variables.vue'
import { RouterView } from 'vue-router'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import { legacyLogicalPropertiesTransformer } from 'ant-design-vue';

dayjs.locale('zh-cn')
</script>

<style scoped></style>
