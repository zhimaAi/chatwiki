<template>
  <div class="skill-generate-tool-page">
    <div class="generate-layout">
      <div class="left-panel">
        <div class="page-header">
          <div class="page-title">
            <a-segmented :value="route.path" :options="titleOptios" @change="handleTitleChange" />
          </div>
        </div>

        <div class="template-list">
          <div
            v-for="item in templateList"
            :key="item.key"
            class="template-card"
            :class="{ active: activeTemplate === item.key }"
            @click="handleTemplateClick(item.key)"
          >
            <div class="template-card-head">
              <div class="template-icon">
                <ThunderboltOutlined />
              </div>
              <a-tooltip>
                <template #title>{{ item.title }}</template>
                <div class="template-title">{{ item.title }}</div>
              </a-tooltip>
            </div>
            <a-tooltip>
              <template #title>{{ item.desc }}</template>
              <div class="template-desc">{{ item.desc }}</div>
            </a-tooltip>
          </div>
        </div>
      </div>

      <BookToSkillTab
        v-show="activeTemplate === 'book'"
        :active="activeTemplate === 'book'"
        :robot-id="currentAssistant?.id"
      />
      <WebToSkillTab v-show="activeTemplate === 'web'" :active="activeTemplate === 'web'" />
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { ThunderboltOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import BookToSkillTab from './components/BookToSkillTab.vue'
import WebToSkillTab from './components/WebToSkillTab.vue'

const route = useRoute()
const router = useRouter()
const clawbotStore = useClawbotStore()
const { currentAssistant } = storeToRefs(clawbotStore)
const { t } = useI18n('views.clawbot.skill-generate-tool.index')

const titleOptios = computed(() => [
  {
    label: t('nav_my_skill'),
    value: '/clawbot/skills'
  },
  {
    label: t('nav_skill_generator'),
    value: '/clawbot/skill-generate-tool'
  }
])

const handleTitleChange = (path) => {
  if (path !== route.path) {
    router.push(path)
  }
}

const templateList = computed(() => [
  {
    key: 'book',
    title: t('book_template_title'),
    desc: t('book_description')
  },
  {
    key: 'web',
    title: t('web_template_title'),
    desc: t('web_description')
  }
])

const activeTemplate = ref('book')

const handleTemplateClick = (key) => {
  activeTemplate.value = key
}
</script>

<style lang="less" scoped>
.skill-generate-tool-page {
  --skills-primary: #2475fc;
  --skills-border: #d9d9d9;
  --skills-title: #262626;
  --skills-text: #595959;
  --skills-text-light: #8c8c8c;

  position: relative;
  min-height: 100vh;
  padding: 0;
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
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;

  .page-title {
    color: var(--skills-title);
    font-size: 20px;
    line-height: 28px;

    .ant-segmented {
      background: #edeff2;
    }

    &::v-deep(.ant-segmented-item-selected) {
      color: var(--skills-primary);
    }
  }
}

.generate-layout {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: 308px minmax(0, 1fr);
  min-height: 100vh;
}

.left-panel {
  padding: 24px 16px 32px;
  border-right: 1px solid #f0f0f0;
}

.template-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.template-card {
  min-height: 114px;
  padding: 24px 22px;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  background: #fff;
  cursor: pointer;
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease;

  &.active {
    border-color: var(--skills-primary);
    box-shadow: 0 0 0 3px rgba(36, 117, 252, 0.16);
  }
}

.template-card-head {
  display: flex;
  align-items: center;
  gap: 12px;
}

.template-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  flex: 0 0 28px;
  border-radius: 8px;
  color: #fff;
  background: #1fb6ff;
  font-size: 16px;
}

.template-title {
  min-width: 0;
  color: var(--skills-title);
  font-size: 14px;
  font-weight: 600;
  line-height: 22px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.template-desc {
  display: -webkit-box;
  margin-top: 14px;
  color: var(--skills-text);
  font-size: 13px;
  line-height: 22px;
  overflow: hidden;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.generate-content {
  min-width: 0;
  padding: 28px 20px 32px;
}

.skill-generate-tool-page :deep(.info-box) {
  padding: 10px 14px;
  border: 1px solid #9bc3ff;
  border-radius: 6px;
  color: #345176;
  background: #eef5ff;
  font-size: 14px;
  line-height: 22px;
}

.skill-generate-tool-page :deep(.add-btn) {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 32px;
  margin-top: 14px;
  padding: 0 14px;
  border-radius: 4px;
  background: var(--skills-primary);
  box-shadow: none;
}

.skill-generate-tool-page :deep(.skill-table) {
  margin-top: 14px;
  color: var(--skills-text);
  font-size: 13px;

  .ant-table {
    color: var(--skills-text);
    font-size: 13px;
  }

  .ant-table-thead > tr > th {
    height: 46px;
    color: #262626;
    background: #f5f5f5;
    font-weight: 400;
  }

  .ant-table-tbody > tr > td {
    height: 45px;
    border-bottom: 1px solid #edf0f2;
  }
}

.skill-generate-tool-page :deep(.skill-name) {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.skill-generate-tool-page :deep(.status-tag) {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 20px;
  padding: 0 7px;
  border-radius: 5px;
  font-size: 12px;
  line-height: 20px;

  &.success {
    color: #12b36f;
    background: #c9f8df;
  }

  &.failed {
    color: #ff7a45;
    background: #ffe6dc;
  }

  &.running,
  &.stopping {
    color: #2475fc;
    background: #dbeaff;
  }

  &.pending,
  &.stopped {
    color: #8c8c8c;
    background: #f0f0f0;
  }
}

.skill-generate-tool-page :deep(.action-list) {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 14px;

  a {
    color: var(--skills-primary);
    white-space: nowrap;
  }
}

@media (max-width: 768px) {
  .skill-generate-tool-page {
    overflow: auto;
  }

  .generate-layout {
    grid-template-columns: 1fr;
  }

  .left-panel {
    padding: 20px 16px 18px;
    border-right: none;
    border-bottom: 1px solid #f0f0f0;
  }

  .generate-content {
    padding: 20px 16px 24px;
  }

  .skill-generate-tool-page :deep(.skill-table .ant-table) {
    min-width: 828px;
  }
}
</style>
