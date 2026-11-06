<template>
  <aside class="customer-sidebar">
    <header class="sidebar-header">
      <RobotSelector
        v-model:open="robotSelectorOpen"
        :value="currentRobot?.id || ''"
        :options="robotOptions"
        :status="robotListStatus"
        @change="emit('robot-change', $event)"
        @retry="emit('retry-robots')"
        @removed="emit('robot-removed', $event)"
        @created="handleRobotCreated"
      >
        <button
          type="button"
          class="app-brand"
          :class="{ 'is-open': robotSelectorOpen }"
          :aria-label="t('select_app')"
          :aria-expanded="robotSelectorOpen"
          aria-haspopup="dialog"
          aria-controls="customer-service-app-picker"
        >
          <img class="brand-logo" :src="currentRobot?.avatar || systemLogo" alt="" />
          <span class="brand-app-name" :title="currentRobot?.name || t('app_name')">{{ currentRobot?.name || t('app_name') }}</span>
          <span class="brand-caret" aria-hidden="true"><img :src="arrowDownIcon" alt="" /></span>
        </button>
      </RobotSelector>
    </header>

    <nav class="customer-menus" :aria-label="t('navigation')">
      <template v-for="item in menuItems" :key="item.key">
        <div v-if="item.children" class="menu-group">
          <button
            type="button"
            class="menu-item"
            :aria-expanded="openKeys.includes(item.key)"
            :aria-controls="'customer-service-menu-' + item.key"
            @click="toggleGroup(item.key)"
          >
            <img class="menu-icon" :src="menuIcons[item.key]" alt="" />
            <span class="menu-label">{{ t(item.titleKey) }}</span>
            <img class="group-arrow" :class="{ expanded: openKeys.includes(item.key) }" :src="arrowDownIcon" alt="" />
          </button>
          <div
            v-show="openKeys.includes(item.key)"
            :id="'customer-service-menu-' + item.key"
            class="group-children"
          >
            <button
              v-for="page in item.children"
              :key="page.key"
              type="button"
              :class="{ active: selectedKey === page.key }"
              :aria-current="selectedKey === page.key ? 'page' : undefined"
              @click="emit('navigate', page.key)"
            >
              <span class="menu-label">{{ t(page.titleKey) }}</span>
            </button>
          </div>
        </div>
        <button
          v-else
          type="button"
          class="menu-item"
          :class="{ active: selectedKey === item.key }"
          :aria-current="selectedKey === item.key ? 'page' : undefined"
          @click="emit('navigate', item.key)"
        >
          <svg-icon v-if="isAgent" :name="item.icon" size="16px" />
          <img v-else class="menu-icon" :src="menuIcons[item.key]" alt="" />
          <span class="menu-label">{{ isAgent ? agentT(item.titleKey) : t(item.titleKey) }}</span>
        </button>
      </template>
    </nav>

    <footer class="sidebar-footer">
      <CustomerUserMenu />
    </footer>
  </aside>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import systemLogo from '@/assets/default_logo.svg'
import { CUSTOMER_SERVICE_AGENT_MENU_ITEMS, CUSTOMER_SERVICE_MENU_ITEMS } from '../navigation'
import { getRobotPermission } from '@/utils/permission'
import RobotSelector from './robot-selector.vue'
import CustomerUserMenu from './customer-user-menu.vue'
import arrowDownIcon from '@/assets/img/customer-service/sidebar/arrow-down.svg'
import chatIcon from '@/assets/img/customer-service/sidebar/chat.svg'
import sessionRecordIcon from '@/assets/img/customer-service/sidebar/session-record.svg'
import channelsIcon from '@/assets/img/customer-service/sidebar/channels.svg'
import aiAgentIcon from '@/assets/img/customer-service/sidebar/ai-agent.svg'
import analyticsIcon from '@/assets/img/customer-service/sidebar/analytics.svg'

const props = defineProps({
  currentRobot: { type: Object, default: null },
  robotOptions: { type: Array, default: () => [] },
  robotListStatus: { type: String, default: 'idle' },
  isAgent: { type: Boolean, default: false },
  selectedKey: { type: String, default: '' }
})
const emit = defineEmits(['navigate', 'robot-change', 'retry-robots', 'robot-removed', 'robot-created'])
const { t } = useI18n('views.customer-service.index')
const { t: agentT } = useI18n('views.clawbot.components.left-sidebar')
const route = useRoute()
const menuItems = computed(() => {
  if (!props.isAgent) return CUSTOMER_SERVICE_MENU_ITEMS
  return CUSTOMER_SERVICE_AGENT_MENU_ITEMS.filter((item) =>
    getRobotPermission(props.currentRobot?.id) !== 1 || ['chat', 'assistant', 'stats'].includes(item.key)
  )
})
const selectedKey = computed(() => props.selectedKey || route.meta.customerServiceMenuKey)
const robotSelectorOpen = ref(false)
const openKeys = ref(CUSTOMER_SERVICE_MENU_ITEMS.filter((item) => item.children).map((item) => item.key))
const menuIcons = {
  chat: chatIcon,
  'session-record': sessionRecordIcon,
  channels: channelsIcon,
  'ai-agent': aiAgentIcon,
  analytics: analyticsIcon
}

const toggleGroup = (key) => {
  if (openKeys.value.includes(key)) {
    openKeys.value = openKeys.value.filter((openKey) => openKey !== key)
  } else {
    openKeys.value.push(key)
  }
}

const handleRobotCreated = (robot) => {
  robotSelectorOpen.value = false
  emit('robot-created', robot)
}

defineExpose({ openRobotSelector: () => { robotSelectorOpen.value = true } })

// 只在初次进入或所属组变化时展开，保留用户在同组内手动收起的状态。
watch(() => route.meta.customerServiceGroupKey, (groupKey) => {
  if (groupKey && !openKeys.value.includes(groupKey)) openKeys.value.push(groupKey)
}, { immediate: true })
</script>

<style lang="less" scoped>
.customer-sidebar {
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  box-sizing: border-box;
  width: 212px;
  height: 100%;
  min-height: 0;
  padding: 20px 12px 16px;
  color: #222;
  background: #f5f5f5;
  overflow: hidden;
}
.sidebar-header {
  flex-shrink: 0;
}
.app-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  height: 40px;
  padding: 0 4px;
  border: 0;
  border-radius: 8px;
  color: #222;
  background: transparent;
  font: inherit;
  font-size: 16px;
  line-height: 24px;
  font-weight: 500;
  text-align: left;
  cursor: pointer;

  &:hover {
    background: #ebebeb;
  }
  &.is-open {
    background: #e8f0fe;
    .brand-caret { filter: brightness(0) saturate(100%) invert(35%) sepia(97%) saturate(2775%) hue-rotate(206deg) brightness(100%) contrast(94%); }
  }
  &:focus-visible {
    outline: 2px solid #146ef5;
    outline-offset: 2px;
  }
}
.brand-logo {
  display: block;
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  object-fit: contain;
}
.brand-app-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.brand-caret {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 12px;
  height: 12px;
}
.customer-menus {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-height: 0;
  padding-top: 24px;
  overflow-y: auto;
  scrollbar-width: thin;
}
.menu-group {
  flex-shrink: 0;
}
.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  flex-shrink: 0;
  height: 36px;
  padding: 0 14px;
  border: 0;
  border-radius: 8px;
  color: #5a5a5a;
  background: transparent;
  font: inherit;
  font-size: 14px;
  line-height: 22px;
  text-align: left;
  cursor: pointer;

  .group-arrow {
    flex-shrink: 0;
    margin-left: auto;
    transform: rotate(-90deg);

    &.expanded {
      transform: rotate(0);
    }
  }
}
.menu-icon {
  display: block;
  flex-shrink: 0;
}
.menu-label {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.group-children {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding-top: 2px;

  button {
    display: flex;
    align-items: center;
    height: 32px;
    padding: 0 14px 0 42px;
    border: 0;
    border-radius: 8px;
    color: #5a5a5a;
    background: transparent;
    font: inherit;
    font-size: 14px;
    line-height: 22px;
    text-align: left;
    cursor: pointer;
  }
}
.menu-item:hover,
.group-children button:hover {
  background: #ebebeb;
}
.menu-item.active,
.group-children button.active {
  color: #146ef5;
  background: #fff;
  font-weight: 500;
}
.menu-item:focus-visible,
.group-children button:focus-visible {
  outline: 2px solid #146ef5;
  outline-offset: -2px;
}
.sidebar-footer {
  flex-shrink: 0;
  padding-top: 12px;
  :deep(.sidebar-user) {
    height: 40px;
    gap: 8px;
    background: transparent;
    &:hover { background: #ebebeb; }
  }
  :deep(.credits-trigger) { gap: 8px; }
  :deep(.credits-label) { color: #5a5a5a; line-height: 20px; }
  :deep(.system-manage-icon) { order: -1; margin-left: 0; color: #5a5a5a; }
}
</style>
