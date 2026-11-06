<style lang="less" scoped>
.chat-conversation-bar {
  min-height: 56px;
  padding: 16px 24px;
  border-bottom: 1px solid var(--app-border-secondary);
  background: var(--app-bg-container);
  color: var(--app-text-primary);
  display: flex;
  align-items: center;
  box-sizing: border-box;
}

.chat-conversation-main {
  min-width: 0;
  max-width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
  overflow: hidden;
}

.conversation-title-box {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.conversation-icon {
  flex-shrink: 0;
  color: var(--app-text-primary);
}

.conversation-title {
  min-width: 0;
  max-width: 420px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--app-text-primary);
  font-size: 16px;
  line-height: 24px;
  font-weight: 600;
  flex-shrink: 1;
}

.knowledge-trigger {
  flex-shrink: 1;
  max-width: 360px;
  padding: 2px 6px;
  border: 0;
  border-radius: var(--app-radius-control);
  background: var(--app-bg-hover);
  display: inline-flex;
  align-items: center;
  min-width: 0;
  cursor: pointer;

  &:hover {
    background: var(--app-fill-secondary);
  }

  &:active,
  &[aria-expanded='true'] {
    background: var(--app-selected);
  }
}

.knowledge-trigger-prefix,
.knowledge-trigger-name,
.knowledge-trigger-count {
  font-size: 12px;
  line-height: 16px;
}

.knowledge-trigger-prefix,
.knowledge-trigger-name {
  color: var(--app-text-secondary);
}

.knowledge-trigger-name {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.knowledge-trigger-count {
  flex-shrink: 0;
  color: var(--app-primary);
}

:deep(.chat-knowledge-dropdown-overlay) {
  .ant-dropdown-menu,
  .ant-dropdown-content {
    padding: 0;
    box-shadow: none;
    background: transparent;
  }
}

.knowledge-dropdown-menu {
  min-width: 312px;
  padding: 2px;
  border-radius: var(--app-radius-popover);
  background: var(--app-bg-elevated);
  box-shadow: var(--app-shadow-popover);
  box-sizing: border-box;
}

.knowledge-dropdown-item {
  width: 100%;
  padding: 5px 16px;
  border: 0;
  border-radius: var(--app-radius-control);
  background: transparent;
  display: flex;
  align-items: center;
  gap: 8px;
  text-align: left;
  cursor: pointer;
  transition: background-color 0.2s ease;

  &:hover {
    background: var(--app-bg-hover);
  }

  &:active {
    background: var(--app-selected);
  }
}

.knowledge-trigger,
.knowledge-dropdown-item {
  &:focus-visible {
    outline: 2px solid var(--app-primary);
    outline-offset: -2px;
  }
}

.knowledge-dropdown-item-icon {
  flex-shrink: 0;
  color: var(--app-text-primary);
  font-size: 14px;
}

.knowledge-dropdown-item-name {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--app-text-primary);
  font-size: 14px;
  line-height: 22px;
}
</style>

<template>
  <div class="chat-conversation-bar">
    <div class="chat-conversation-main">
      <div class="conversation-title-box">
        <svg-icon class="conversation-icon" name="message2" size="18" />
        <div class="conversation-title">{{ title }}</div>
      </div>

      <a-dropdown
        v-model:open="dropdownOpen"
        :trigger="['click']"
        placement="bottomLeft"
        overlay-class-name="chat-knowledge-dropdown-overlay"
      >
        <button type="button" class="knowledge-trigger" :aria-expanded="dropdownOpen">
          <span class="knowledge-trigger-prefix">{{ t('label_knowledge_base') }}</span>
          <span class="knowledge-trigger-name">{{ primaryLibrary.library_name }}</span>
          <span v-if="extraLibraryCount > 0" class="knowledge-trigger-count">（+{{ extraLibraryCount }}）</span>
        </button>

        <template #overlay>
          <div class="knowledge-dropdown-menu">
            <button
              v-for="item in libraries"
              :key="item.id"
              type="button"
              class="knowledge-dropdown-item"
              @click="handleOpenLibrary(item)"
            >
              <ExportOutlined class="knowledge-dropdown-item-icon" />
              <span class="knowledge-dropdown-item-name">{{ item.library_name }}</span>
            </button>
          </div>
        </template>
      </a-dropdown>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ExportOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.clawbot.chat.components.conversation-bar')
defineProps({
  title: {
    type: String,
    default: ''
  },
  primaryLibrary: {
    type: Object,
    default: () => ({})
  },
  extraLibraryCount: {
    type: Number,
    default: 0
  },
  libraries: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['openLibrary'])

const dropdownOpen = ref(false)

const handleOpenLibrary = (item) => {
  dropdownOpen.value = false
  emit('openLibrary', item)
}
</script>
