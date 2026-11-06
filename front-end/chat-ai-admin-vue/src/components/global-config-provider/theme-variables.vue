<template>
  <slot />
</template>

<script setup>
import { onBeforeUnmount, onMounted, watchEffect } from 'vue'
import { theme } from 'ant-design-vue'
import { useThemeStore } from '@/stores/modules/theme'

const themeStore = useThemeStore()
// 使用当前 Provider 推导出的 Token，使业务悬停、选中等状态与组件库保持一致。
const { token } = theme.useToken()

watchEffect(() => {
  const current = themeStore.currentTheme
  const colors = token.value
  // 仅默认主题保留历史业务值，不再用原始配置覆盖 Provider 的最终 Token。
  const legacyBusiness = current.legacyBusiness || {}
  const variables = {
    '--app-primary': colors.colorPrimary,
    '--app-primary-hover': colors.colorPrimaryHover,
    '--app-primary-active': colors.colorPrimaryActive,
    '--app-primary-border': colors.colorPrimaryBorder,
    '--app-link': legacyBusiness.link ?? colors.colorLink,
    '--app-link-hover': legacyBusiness.linkHover ?? colors.colorLinkHover,
    '--app-link-active': legacyBusiness.linkActive ?? colors.colorLinkActive,
    '--app-selected': colors.colorPrimaryBg,
    '--app-bg-layout': legacyBusiness.layoutBackground ?? colors.colorBgLayout,
    '--app-bg-container': legacyBusiness.containerBackground ?? colors.colorBgContainer,
    '--app-bg-elevated': colors.colorBgElevated,
    '--app-bg-disabled': colors.colorBgContainerDisabled,
    '--app-bg-hover': colors.colorFillTertiary,
    '--app-fill-secondary': colors.colorFillSecondary,
    '--app-text-primary': legacyBusiness.textPrimary ?? colors.colorText,
    '--app-text-secondary': legacyBusiness.textSecondary ?? colors.colorTextSecondary,
    '--app-text-tertiary': legacyBusiness.textTertiary ?? colors.colorTextTertiary,
    '--app-text-disabled': colors.colorTextDisabled,
    '--app-text-on-primary': colors.colorTextLightSolid,
    '--app-border': colors.colorBorder,
    '--app-border-secondary': colors.colorBorderSecondary,
    '--app-success': colors.colorSuccess,
    '--app-warning': colors.colorWarning,
    '--app-error': colors.colorError,
    '--app-font-family': colors.fontFamily,
    '--app-radius-control': `${colors.borderRadius}px`,
    '--app-radius-lg': `${colors.borderRadiusLG}px`,
    '--app-radius-sm': `${colors.borderRadiusSM}px`,
    '--app-shadow-card': colors.boxShadowSecondary,
    '--app-shadow-focus': `0 0 0 ${colors.controlOutlineWidth}px ${colors.controlOutline}`,
    '--app-radius-sidebar': `${current.radius.sidebarItem ?? current.radius.control}px`,
    '--app-radius-sidebar-child': `${current.radius.sidebarChildItem ?? current.radius.control}px`,
    '--app-radius-popover': `${current.radius.popover ?? colors.borderRadiusLG}px`,
    '--app-shadow-popover': current.shadow.popover || current.shadow.elevated
  }

  // 放在 html 上，挂到 body 的自定义浮层也能继承；不在侧栏上再维护一份样式。
  const root = document.documentElement
  root.dataset.theme = themeStore.themeId
  Object.entries(variables).forEach(([name, value]) => root.style.setProperty(name, value))
})

const handleStorage = (event) => themeStore.syncFromStorage(event)
onMounted(() => window.addEventListener('storage', handleStorage))
onBeforeUnmount(() => window.removeEventListener('storage', handleStorage))
</script>
