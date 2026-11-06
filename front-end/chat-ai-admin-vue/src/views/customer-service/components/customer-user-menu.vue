<template>
  <a-popover
    v-model:open="popoverOpen"
    trigger="click"
    placement="topLeft"
    :arrow="false"
    :force-render="true"
    :overlay-inner-style="popoverInnerStyle"
  >
    <div class="sidebar-user">
      <button
        type="button"
        class="credits-trigger"
        :aria-expanded="popoverOpen"
        aria-controls="customer-user-popover"
        @click.stop="popoverOpen = !popoverOpen"
      >
        <span class="credits-label">{{ globalT('layout.navbar.system_manage') }}</span>
        <SettingOutlined class="system-manage-icon" aria-hidden="true" />
      </button>
    </div>
    <template #content>
      <div
        id="customer-user-popover"
        class="user-popover"
        @keydown.esc.stop="closeMenu"
        @scroll="submenu = ''"
      >
        <a
          v-for="item in accountMenus"
          :key="item.path"
          class="user-menu-item"
          :href="router.resolve(item.path).href"
          target="_blank"
          rel="noopener noreferrer"
        >
          <component :is="item.icon" />
          <span>{{ item.label }}</span>
        </a>
        <!-- 主题切换入口暂时隐藏，保留实现便于后续恢复。 -->
        <a-popover
          v-if="false"
          :open="submenu === 'appearance'"
          trigger="click"
          placement="rightTop"
          :arrow="false"
          :overlay-inner-style="popoverInnerStyle"
          @open-change="setSubmenu('appearance', $event)"
        >
          <button type="button" class="user-menu-item" :aria-expanded="submenu === 'appearance'">
            <BgColorsOutlined /><span>{{ t('appearance') }}</span>
            <span class="current-appearance">{{ t('appearance_' + themeStore.themeId) }}</span
            ><RightOutlined class="menu-arrow" />
          </button>
          <template #content>
            <div class="user-submenu">
              <button
                v-for="item in appearances"
                :key="item.id"
                type="button"
                class="submenu-item"
                :class="{ active: themeStore.themeId === item.id }"
                :aria-pressed="themeStore.themeId === item.id"
                @click="selectAppearance(item.id)"
              >
                <i class="appearance-dot" :style="{ background: item.antTheme.token.colorPrimary }"></i>
                <span>{{ t('appearance_' + item.id) }}</span
                ><CheckOutlined v-if="themeStore.themeId === item.id" class="submenu-check" />
              </button>
            </div>
          </template>
        </a-popover>
        <a-popover
          :open="submenu === 'language'"
          trigger="click"
          placement="rightTop"
          :arrow="false"
          :overlay-inner-style="popoverInnerStyle"
          @open-change="setSubmenu('language', $event)"
        >
          <button type="button" class="user-menu-item" :aria-expanded="submenu === 'language'">
            <GlobalOutlined /><span>{{ selectedLocaleName }}</span
            ><RightOutlined class="menu-arrow" />
          </button>
          <template #content>
            <div class="user-submenu">
              <button
                v-for="item in localeStore.getLocaleMap"
                :key="item.lang"
                type="button"
                class="submenu-item"
                :class="{ active: currentLocale === item.lang }"
                :aria-pressed="currentLocale === item.lang"
                @click="selectLocale(item.lang)"
              >
                <span>{{ item.name }}</span
                ><CheckOutlined v-if="currentLocale === item.lang" class="submenu-check" />
              </button>
            </div>
          </template>
        </a-popover>
        <button type="button" class="user-menu-item" @click="logout">
          <LogoutOutlined /><span>{{ globalT('layout.navbar.logout') }}</span>
        </button>
      </div>
    </template>
  </a-popover>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { App } from 'ant-design-vue'
import {
  AppstoreOutlined,
  BankOutlined,
  BarChartOutlined,
  BgColorsOutlined,
  CheckOutlined,
  GlobalOutlined,
  LogoutOutlined,
  RightOutlined,
  SettingOutlined,
  TeamOutlined
} from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useLocale } from '@/hooks/web/useLocale'
import { useLocaleStore } from '@/stores/modules/locale'
import { useUserStore } from '@/stores/modules/user'
import { useThemeStore } from '@/stores/modules/theme'
import { themes } from '@/config/themes'
import { checkRole } from '@/utils/permission'

const { t } = useI18n('views.customer-service.index')
const { t: accountT } = useI18n('views.user.left-menus')
const { t: globalT } = useI18n()
const router = useRouter()
const userStore = useUserStore()
const localeStore = useLocaleStore()
const themeStore = useThemeStore()
const { message } = App.useApp()
const { changeLocale } = useLocale()
const popoverOpen = ref(false)
const submenu = ref('')
const popoverInnerStyle = { padding: 0, background: 'transparent', boxShadow: 'none' }
const appearances = Object.values(themes)
const appearanceMessageKey = 'customer-appearance-change'
const accountMenus = computed(() =>
  [
    {
      path: '/user/enterprise',
      label: accountT('enterprise_settings'),
      icon: BankOutlined,
      permission: 'CompanyManage'
    },
    {
      path: '/user/model',
      label: accountT('model_management'),
      icon: AppstoreOutlined,
      permission: 'ModelManage'
    },
    {
      path: '/user/usetoken',
      label: t('usage_statistics'),
      icon: BarChartOutlined,
      permission: 'TokenManage'
    },
    {
      path: '/user/manage',
      label: accountT('team_management'),
      icon: TeamOutlined,
      permission: 'TeamManage'
    },
    { path: '/user/account', label: accountT('account_settings'), icon: SettingOutlined }
  ].filter((item) => !item.permission || checkRole([item.permission]))
)
const currentLocale = computed(() => {
  const lang = localeStore.getCurrentLocale.lang
  return lang === 'en' ? 'en-US' : lang
})
const selectedLocaleName = computed(
  () => localeStore.getLocaleMap.find((item) => item.lang === currentLocale.value)?.name
)

const closeMenu = () => {
  popoverOpen.value = false
  submenu.value = ''
}
const setSubmenu = (name, open) => {
  if (open) submenu.value = name
  else if (submenu.value === name) submenu.value = ''
}
const selectAppearance = async (themeId) => {
  if (themeId === themeStore.themeId && !themeStore.storageError) {
    submenu.value = ''
    return
  }
  const saved = themeStore.setTheme(themeId)
  submenu.value = ''
  // 消息在调用时读取样式 hash，需等待 Provider 与消息 Holder 完成主题更新。
  await nextTick()
  // 复用同一条提示，更新内容与当前主题的样式 hash，避免旧提示叠加或丢失样式。
  if (!saved) {
    message.warning({ key: appearanceMessageKey, content: t('appearance_save_failed') })
    return
  }
  message.success({
    key: appearanceMessageKey,
    content: t('appearance_changed', { name: t('appearance_' + themeId) })
  })
}
const selectLocale = async (lang) => {
  if (lang === currentLocale.value) {
    submenu.value = ''
    return
  }
  await changeLocale(lang)
  window.location.reload()
}
const logout = () => {
  closeMenu()
  userStore.logoutConfirm()
}
const closeSubmenu = () => {
  submenu.value = ''
}
const handleEscape = (event) => {
  if (event.key === 'Escape') closeMenu()
}
watch(popoverOpen, (open) => {
  if (!open) submenu.value = ''
})
onMounted(() => {
  window.addEventListener('scroll', closeSubmenu, true)
  document.addEventListener('keydown', handleEscape)
})
onBeforeUnmount(() => {
  window.removeEventListener('scroll', closeSubmenu, true)
  document.removeEventListener('keydown', handleEscape)
})
</script>

<style lang="less" scoped>
.sidebar-user {
  display: flex;
  align-items: center;
  gap: 6px;
  height: 36px;
  padding: 0 10px;
  border-radius: 8px;
  background: var(--app-bg-hover);
  transition: background 150ms ease-out;

  &:hover {
    background: var(--app-fill-secondary);
  }
}
.credits-trigger {
  display: flex;
  flex: 1;
  align-items: center;
  gap: 6px;
  min-width: 0;
  height: 100%;
  padding: 0;
  border: 0;
  background: transparent;
  cursor: pointer;
}
.credits-label {
  color: var(--app-text-secondary);
  font-size: 12px;
  white-space: nowrap;
}
.system-manage-icon {
  flex-shrink: 0;
  margin-left: auto;
  color: var(--app-text-secondary);
  font-size: 16px;
}
.user-popover,
.user-submenu {
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  padding: 6px;
  border: 1px solid var(--app-border);
  border-radius: var(--app-radius-popover, 12px);
  color: var(--app-text-primary);
  background: var(--app-bg-container);
  box-shadow: var(--app-shadow-popover, 0 8px 24px rgba(0, 0, 0, 0.12));
}
.user-popover {
  width: 280px;
  max-width: calc(100vw - 24px);
  max-height: calc(100vh - 64px);
  overflow-y: auto;
}
.user-menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  box-sizing: border-box;
  padding: 10px 12px;
  border: 0;
  border-radius: 8px;
  color: var(--app-text-primary);
  background: transparent;
  font: inherit;
  font-size: 14px;
  text-align: left;
  text-decoration: none;
  cursor: pointer;

  &:hover {
    color: var(--app-text-primary);
    background: var(--app-bg-hover);
  }
  .anticon {
    flex-shrink: 0;
    width: 18px;
    height: 18px;
    color: var(--app-text-secondary);
    font-size: 18px;
  }
  .menu-arrow {
    margin-left: auto;
    font-size: 16px;
  }
}
.current-appearance {
  margin-left: auto;
  color: var(--app-text-secondary);
  font-size: 12px;

  + .menu-arrow {
    margin-left: 0;
  }
}
.user-submenu {
  gap: 2px;
  min-width: 184px;
}
.submenu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  border: 0;
  border-radius: 8px;
  color: var(--app-text-primary);
  background: transparent;
  font: inherit;
  font-size: 14px;
  text-align: left;
  cursor: pointer;
  transition: background 150ms ease-out;

  &:hover {
    background: var(--app-bg-hover);
  }
  &.active {
    color: var(--app-primary);
    background: var(--app-selected);
    font-weight: 500;
  }
}
.appearance-dot {
  width: 12px;
  height: 12px;
  border: 1px solid var(--app-border);
  border-radius: 50%;
}
.submenu-check {
  margin-left: auto;
  font-size: 16px;
}
.credits-trigger:focus-visible,
.user-menu-item:focus-visible,
.submenu-item:focus-visible {
  outline: 2px solid var(--app-primary);
  outline-offset: 2px;
}
@media (prefers-reduced-motion: reduce) {
  .sidebar-user,
  .submenu-item {
    transition: none;
  }
}
</style>
