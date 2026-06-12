<template>
  <div class="knowledge-page">
    <div class="page-header">
      <div class="page-title">{{ t('page_title') }}</div>
    </div>

    <div class="knowledge-content">
      <div class="top-tabs">
        <div
          class="top-tab"
          :class="{ active: activeTab === 'local' }"
          @click="activeTab = 'local'"
        >
          {{ t('tab_local_docs') }}
        </div>
        <div
          class="top-tab"
          :class="{ active: activeTab === 'related' }"
          @click="activeTab = 'related'"
        >
          {{ t('tab_related_knowledge') }}
        </div>
      </div>

      <div class="tab-panel">
        <LocalDocs v-if="activeTab === 'local'" />
        <RelatedKnowledge v-if="activeTab === 'related'" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import LocalDocs from './LocalDocs.vue'
import RelatedKnowledge from './RelatedKnowledge.vue'

const { t } = useI18n('views.clawbot.knowledge.index')
const activeTab = ref('local')
</script>

<style lang="less" scoped>
.knowledge-page {
  --knowledge-primary: #2475fc;
  --knowledge-text-main: #262626;
  --knowledge-text-secondary: #595959;

  position: relative;
  min-height: 100vh;
  padding: 20px 24px 32px;
  background: #fff;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: -88px;
    right: -46px;
    width: 284px;
    height: 274px;
    border-radius: 999px;
    background: rgba(198, 210, 255, 0.3);
    filter: blur(64px);
    pointer-events: none;
  }
}

.page-header {
  position: relative;
  z-index: 1;

  .page-title {
    color: var(--knowledge-text-main);
    font-size: 20px;
    font-weight: 600;
    line-height: 28px;
  }
}

.knowledge-content {
  position: relative;
  z-index: 1;
  margin-top: 20px;
}

.top-tabs {
  display: flex;
  gap: 8px;
}

.top-tab {
  height: 32px;
  padding: 0 16px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  color: var(--knowledge-text-secondary);
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;
  transition: all 0.2s ease;

  &.active {
    color: var(--knowledge-primary);
    background: #e5efff;
  }
}

.tab-panel {
  margin-top: 10px;
}
</style>
