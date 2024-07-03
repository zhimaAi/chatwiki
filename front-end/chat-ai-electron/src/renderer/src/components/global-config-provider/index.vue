<template>
  <ConfigProvider v-bind="$props" :locale="currentLocale.elLocale" :theme="theme">
    <slot />
  </ConfigProvider>
</template>

<script setup>
import { provide, computed } from 'vue'
import { ConfigProvider } from 'ant-design-vue'
import { useLocaleStore } from '@/stores/modules/locale'
// 主题
import theme from '@/config/theme'

const props = defineProps({
  // eslint-disable-next-line vue/require-prop-type-constructor
  size: 'default'
})

provide('configGlobal', props)

// 多语言相关
const localeStore = useLocaleStore()

const currentLocale = computed(() => localeStore.currentLocale)
</script>
