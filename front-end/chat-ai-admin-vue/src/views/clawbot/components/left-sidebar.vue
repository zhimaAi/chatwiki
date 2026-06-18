<template>
  <aside class="clawbot-sidebar">
    <div class="sidebar-main">
      <div class="sidebar-header">
        <a-dropdown :trigger="['click']" overlay-class-name="clawbot-assistant-dropdown">
          <button class="brand-card" type="button">
            <div class="brand">
              <div class="brand-avatar">
                <img v-if="currentAssistantAvatar" :src="currentAssistantAvatar" alt="" />
                <span v-else class="brand-fallback">🔥</span>
              </div>
              <div class="brand-name">{{ currentAssistantName }}</div>
            </div>
            <DownOutlined class="brand-arrow" />
          </button>
          <template #overlay>
            <div class="assistant-dropdown-menu">
              <button
                v-for="item in assistantList"
                :key="item.id"
                type="button"
                class="assistant-option"
                :class="{ active: item.id === currentAssistantId }"
                @click="handleSelectAssistant(item.id)"
              >
                <div class="assistant-option-avatar">
                  <img v-if="getAssistantAvatar(item)" :src="getAssistantAvatar(item)" alt="" />
                  <span v-else class="brand-fallback">🔥</span>
                </div>
                <span class="assistant-option-name">{{ item.robot_name || t('label_agent_fallback') }}</span>
                <CheckOutlined v-if="item.id === currentAssistantId" class="assistant-option-check" />
              </button>
            </div>
          </template>
        </a-dropdown>
      </div>

      <nav class="sidebar-nav">
        <button
          v-for="item in navItems"
          :key="item.path"
          type="button"
          class="nav-item"
          :class="{ active: route.path === item.path }"
          @click="handleNavigate(item)"
        >
          <svg-icon :name="item.icon" class="nav-icon" />
          <span class="nav-label">{{ item.label }}</span>
        </button>
      </nav>
    </div>

    <div class="sidebar-footer">
      <button class="footer-btn primary" type="button" @click="handleCreateAssistant">
        <PlusOutlined class="footer-icon" />
        <span>{{ t('btn_new_agent') }}</span>
      </button>
      <button class="footer-btn" type="button" @click="handleBackChatwiki">
        <LeftOutlined class="footer-icon" />
        <span>{{ t('btn_back_chatwiki') }}</span>
      </button>
    </div>

    <AddAgentAlert ref="addAgentAlertRef" @ok="onAgentCreated" />
  </aside>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  DownOutlined,
  PlusOutlined,
  LeftOutlined,
  CheckOutlined,
} from '@ant-design/icons-vue'
import { storeToRefs } from 'pinia'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import AddAgentAlert from '@/views/clawbot/components/add-agent-alert.vue'

const { t } = useI18n('views.clawbot.components.left-sidebar')
const route = useRoute()
const router = useRouter()
const clawbotStore = useClawbotStore()
const { currentAssistant, currentAssistantId, robotInfo, assistantList } = storeToRefs(clawbotStore)
const currentAssistantData = computed(() => {
  const currentData = currentAssistant.value || {}
  const robotDetail =
    robotInfo.value && String(robotInfo.value.id || '') === String(currentAssistantId.value || '') ? robotInfo.value : {}

  return {
    ...currentData,
    ...robotDetail,
  }
})
const currentAssistantName = computed(() => currentAssistantData.value.robot_name || t('label_agent_fallback'))
const currentAssistantAvatar = computed(() => {
  return currentAssistantData.value.robot_avatar_url || currentAssistantData.value.robot_avatar || ''
})

const addAgentAlertRef = ref(null)

const navItems = computed(() => ([
  { label: t('nav_chat'), path: '/clawbot/chat', icon: 'chat-menu' },
  { label: t('nav_agent'), path: '/clawbot/assistant', icon: 'agent-menu' },
  { label: t('nav_skill'), path: '/clawbot/skills', icon: 'skill-menu' },
  { label: t('nav_knowledge_base'), path: '/clawbot/knowledge', icon: 'knowledge-menu' },
  { label: t('nav_external_services'), path: '/clawbot/services', icon: 'service-menu' },
  { label: t('nav_analytics'), path: '/clawbot/stats', icon: 'analytics-menu' },
  { label: t('nav_settings'), path: '/clawbot/settings', icon: 'settings-menu', query: { menu: 'basic' } },
]))

const handleCreateAssistant = () => {
  addAgentAlertRef.value.open()
}

const handleBackChatwiki = () => {
  router.replace('/robot/list')
}

const getAssistantAvatar = (item) => {
  return item?.robot_avatar_url || item?.robot_avatar || ''
}

const buildAssistantQuery = (assistant = currentAssistant.value) => {
  if (!assistant?.id) {
    return {}
  }

  return {
    id: assistant.id,
    robot_key: assistant.robot_key
  }
}

const handleNavigate = (item) => {
  router.push({
    path: item.path,
    query: {
      ...buildAssistantQuery(),
      ...(item.query || {})
    }
  })
}

const handleSelectAssistant = (id) => {
  const assistant = clawbotStore.selectAssistant(id)
  if (!assistant?.id) {
    return
  }

  router.replace({
    path: route.path,
    query: {
      ...route.query,
      ...buildAssistantQuery(assistant)
    }
  })
}

const onAgentCreated = (data) => {
  clawbotStore.fetchAssistants().then(() => {
    if (data && data.id) {
      const assistant = clawbotStore.selectAssistant(String(data.id))
      if (assistant?.id) {
        router.replace({
          path: route.path,
          query: {
            ...route.query,
            ...buildAssistantQuery(assistant)
          }
        })
      }
    }
  })
}
</script>

<style lang="less" scoped>
.clawbot-sidebar {
  width: 256px;
  height: 100vh;
  padding: 0;
  background: #f9f9f9;
  border-right: 1px solid #f0f0f0;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  flex-shrink: 0;
}

.sidebar-main {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
}

.sidebar-header {
  padding: 12px 8px 8px;
}

.brand-card,
.nav-item,
.footer-btn {
  width: 100%;
  padding: 0;
  border: 0;
  background: transparent;
  outline: none;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  font-family: inherit;
  text-align: left;
  cursor: pointer;
}

.brand-card {
  min-height: 48px;
  justify-content: space-between;
  padding: 12px;
  border: 1px solid #e4e6eb;
  border-radius: 20px;
  background: #fff;
}

.brand {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}

.brand-avatar {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  overflow: hidden;
  background: linear-gradient(180deg, #ffb347 0%, #ff4d4f 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.brand-fallback {
  font-size: 12px;
  line-height: 1;
}

.brand-name {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #262626;
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
}

.brand-arrow {
  flex-shrink: 0;
  color: #595959;
  font-size: 12px;
}

:deep(.clawbot-assistant-dropdown) {
  .ant-dropdown-menu,
  .ant-dropdown-content {
    padding: 0;
    box-shadow: none;
    background: transparent;
  }
}

.assistant-dropdown-menu {
  width: 240px;
  margin-top: 4px;
  padding: 6px;
  border: 1px solid #e4e6eb;
  border-radius: 16px;
  background: #fff;
  box-shadow: 0 10px 24px rgba(0, 0, 0, 0.08);
}

.assistant-option {
  width: 100%;
  min-height: 40px;
  padding: 8px 10px;
  border: 0;
  border-radius: 12px;
  background: transparent;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #262626;
  cursor: pointer;
  transition: background-color 0.2s ease, color 0.2s ease;

  &:hover {
    background: #f5f5f5;
  }

  &.active {
    background: #f4f8ff;
    color: #2475fc;
  }
}

.assistant-option-avatar {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  overflow: hidden;
  background: linear-gradient(180deg, #ffb347 0%, #ff4d4f 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.assistant-option-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
  line-height: 20px;
}

.assistant-option-check {
  flex-shrink: 0;
  font-size: 12px;
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 0 8px;
  padding-bottom: 12px;
}

.nav-item {
  min-height: 36px;
  gap: 9px;
  margin: 0;
  padding: 7px 17px;
  border-radius: 8px;
  color: #262626;
  transition: background-color 0.2s ease, box-shadow 0.2s ease, color 0.2s ease;

  &:hover {
    background: #e4e6eb;
  }

  &.active {
    background: #fff;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.06);
    color: #2475fc;
  }
}

.nav-icon {
  flex-shrink: 0;
  font-size: 18px;
  line-height: 1;
}

.nav-label {
  flex: 1;
  min-width: 0;
  font-size: 14px;
  font-weight: 400;
  line-height: 22px;
}

.sidebar-footer {
  padding: 0 8px 16px;
}

.footer-btn {
  min-height: 36px;
  justify-content: center;
  gap: 4px;
  padding: 5px 16px;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  background: #fff;
  color: #595959;
  font-size: 14px;
  font-weight: 400;
  line-height: 22px;

  & + .footer-btn {
    margin-top: 8px;
  }

  &.primary {
    border-color: #2475fc;
    color: #2475fc;
  }
}

.footer-icon {
  flex-shrink: 0;
  font-size: 16px;
}
</style>
