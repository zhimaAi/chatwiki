<template>
  <ConfigProvider v-bind="$props" :locale="currentLocale.antvLocale" :theme="themeStore.antTheme">
    <slot />
  </ConfigProvider>
</template>

<script setup>
import { provide, computed } from 'vue'
import { ConfigProvider } from 'ant-design-vue'
import { useLocaleStore } from '@/stores/modules/locale'
import { useThemeStore } from '@/stores/modules/theme'

const props = defineProps({
  // eslint-disable-next-line vue/require-prop-type-constructor
  size: 'default'
})

provide('configGlobal', props)

// 多语言相关
const localeStore = useLocaleStore()
const themeStore = useThemeStore()

const currentLocale = computed(() => localeStore.currentLocale)
</script>
