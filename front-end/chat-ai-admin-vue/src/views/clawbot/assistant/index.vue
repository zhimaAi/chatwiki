<template>
  <div class="assistant-page">
    <AddAgentAlert ref="addAgentAlertRef" @ok="onAgentCreated" />

    <div class="page-header">
      <div class="page-title">{{ t('page_title') }}</div>
      <a-button type="primary" class="create-btn" @click="handleCreateAssistant">
        <template #icon>
          <PlusOutlined />
        </template>
        {{ t('btn_create_agent') }}
      </a-button>
    </div>

    <div class="search-row">
      <a-input v-model:value="searchKeyword" allow-clear class="search-input" :placeholder="t('ph_search')">
        <template #suffix>
          <SearchOutlined />
        </template>
      </a-input>
    </div>

    <div class="loading-box" v-if="isInitializing">
      <a-spin />
      <div class="loading-text">{{ t('msg_loading_list') }}</div>
    </div>
    <div class="empty-box" v-else-if="!cards.length">
      <a-empty :description="searchKeyword ? t('msg_no_matching_agents') : undefined" />
    </div>
    <div class="card-grid" v-else>
      <div
        v-for="item in cards"
        :key="item.id"
        class="agent-card"
        @click="handleSelectAssistant(item.id)"
      >
        <div class="card-top">
          <img class="card-avatar" :src="item.robot_avatar" alt="" />
          <div class="card-info">
            <div class="card-title">{{ item.robot_name }}</div>
            <div class="card-tag">{{ t('tag_chatbot') }}</div>
          </div>
        </div>
        <div class="card-desc">{{ item.robot_intro }}</div>
        <div class="card-footer">
          <a-tooltip :title="t('tooltip_session_history')">
            <div class="footer-action" @click.stop="handleOpenSession(item)">
              <svg-icon name="session" />
            </div>
          </a-tooltip>
          <a-tooltip :title="t('tooltip_analytics')">
            <div class="footer-action" @click.stop="handleOpenStats(item)">
              <svg-icon name="analysis" />
            </div>
          </a-tooltip>
          <a-tooltip :title="t('tooltip_duplicate_agent')">
            <div class="footer-action" @click.stop="handleCopy(item)">
              <CopyOutlined />
            </div>
          </a-tooltip>
          <a-dropdown placement="bottomRight">
            <div class="footer-action more" @click.stop>
              <svg-icon name="point-h" />
            </div>
            <template #overlay>
              <a-menu>
                <a-menu-item>
                  <a class="delete-text-color" href="javascript:;" @click.stop="handleDelete(item)">{{ t('action_delete') }}</a>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, createVNode, ref } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { Modal, message } from 'ant-design-vue'
import { CopyOutlined, ExclamationCircleOutlined, PlusOutlined, SearchOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import AddAgentAlert from '@/views/clawbot/components/add-agent-alert.vue'

const { t } = useI18n('views.clawbot.assistant.index')
const router = useRouter()
const clawbotStore = useClawbotStore()
const { assistantList, isInitializing } = storeToRefs(clawbotStore)

const addAgentAlertRef = ref(null)
const searchKeyword = ref('')

const cards = computed(() => {
  const keyword = searchKeyword.value.trim().toLowerCase()

  if (!keyword) {
    return assistantList.value
  }

  return assistantList.value.filter((item) => {
    const name = String(item.robot_name || '').toLowerCase()
    const intro = String(item.robot_intro || '').toLowerCase()
    return name.includes(keyword) || intro.includes(keyword)
  })
})

const handleCreateAssistant = () => {
  addAgentAlertRef.value.open()
}

const onAgentCreated = (data) => {
  clawbotStore.fetchAssistants().then(() => {
    if (data && data.id) {
      clawbotStore.selectAssistant(String(data.id))
    }
  })
}

const buildAssistantQuery = (item, query = {}) => {
  return {
    ...query,
    id: item.id,
    robot_key: item.robot_key
  }
}

const openClawbotPage = (path, query = {}) => {
  const searchParams = new URLSearchParams(query)
  const queryString = searchParams.toString()
  window.open(`#${path}${queryString ? `?${queryString}` : ''}`, '_blank', 'noopener')
}

const handleSelectAssistant = (id) => {
  const assistant = clawbotStore.selectAssistant(id)
  if (!assistant?.id) {
    return
  }

  router.push({
    path: '/clawbot/chat',
    query: buildAssistantQuery(assistant)
  })
}

const handleOpenSession = (item) => {
  openClawbotPage('/clawbot/settings', buildAssistantQuery(item, { menu: 'session-record' }))
}

const handleOpenStats = (item) => {
  openClawbotPage('/clawbot/stats', buildAssistantQuery(item))
}

const handleCopy = (item) => {
  clawbotStore.copyAssistant(item.id).then(() => {
    message.success(t('msg_copy_success'))
  })
}

const handleDelete = (item) => {
  Modal.confirm({
    title: t('title_delete_agent', { name: item.robot_name }),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_confirm_delete_agent'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    okType: 'danger',
    onOk() {
      return clawbotStore.deleteAssistant(item.id).then(() => {
        message.success(t('msg_delete_success'))
      })
    }
  })
}
</script>

<style lang="less" scoped>
.assistant-page {
  --assistant-primary: #2475fc;
  --assistant-border: #e4e6eb;
  --assistant-text-main: #262626;
  --assistant-text-secondary: #595959;

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

.empty-box {
  position: relative;
  z-index: 1;
  min-height: 420px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-box {
  position: relative;
  z-index: 1;
  min-height: 420px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;

  .loading-text {
    font-size: 14px;
    color: #98a2b3;
  }
}

.page-header,
.search-row,
.card-grid {
  position: relative;
  z-index: 1;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;

  .page-title {
    color: var(--assistant-text-main);
    font-size: 20px;
    font-weight: 600;
    line-height: 28px;
  }

  .create-btn {
    height: 32px;
    padding: 0 16px;
    border: none;
    border-radius: 6px;
    background: var(--assistant-primary);
    box-shadow: none;
    font-size: 14px;
    line-height: 22px;

    &:hover,
    &:focus {
      background: #4a8dff;
    }
  }
}

.search-row {
  margin-top: 12px;
}

.search-input {
  width: 400px;
  max-width: 100%;
}

.search-input :deep(.ant-input-affix-wrapper) {
  height: 32px;
  padding: 4px 12px;
  border-color: #d9d9d9;
  border-radius: 8px;
  box-shadow: none;
}

.search-input :deep(.ant-input-affix-wrapper:hover),
.search-input :deep(.ant-input-affix-wrapper-focused) {
  border-color: var(--assistant-primary);
}

.search-input :deep(.ant-input) {
  color: var(--assistant-text-main);
  font-size: 14px;
  line-height: 22px;
}

.search-input :deep(.ant-input::placeholder) {
  color: rgba(0, 0, 0, 0.25);
}

.search-input :deep(.ant-input-suffix) {
  color: rgba(0, 0, 0, 0.25);
}

.card-grid {
  margin-top: 20px;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px 14px;
}

.agent-card {
  min-height: 188px;
  border-radius: 12px;
  border: 1px solid var(--assistant-border);
  background: #fff;
  padding: 24px;
  cursor: pointer;
  transition: box-shadow 0.2s ease, transform 0.2s ease, border-color 0.2s ease;
  display: flex;
  flex-direction: column;

  &:hover {
    border-color: #d6dce5;
    box-shadow: 0 8px 18px rgba(15, 23, 42, 0.08);
    transform: translateY(-1px);
  }

  &.active {
    border-color: #d6dce5;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  }

  .card-top {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .card-info {
    flex: 1;
    min-width: 0;
  }

  .card-avatar {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    background: #4876ff;
    object-fit: cover;
    flex-shrink: 0;
  }

  .card-title {
    color: var(--assistant-text-main);
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .card-tag {
    margin-top: 4px;
    display: inline-flex;
    align-items: center;
    height: 22px;
    padding: 0 8px;
    border: 1px solid #cde0ff;
    border-radius: 6px;
    color: var(--assistant-primary);
    font-size: 12px;
    line-height: 20px;
  }

  .card-desc {
    margin-top: 12px;
    min-height: 44px;
    color: var(--assistant-text-secondary);
    font-size: 14px;
    line-height: 22px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .card-footer {
    display: flex;
    gap: 12px;
    align-items: center;
    margin-top: auto;
    padding-top: 12px;
    color: var(--assistant-text-secondary);

    .footer-action {
      width: 24px;
      height: 24px;
      border-radius: 6px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: all 0.2s ease;

      &:hover {
        background: #f2f4f7;
        color: var(--assistant-primary);
      }

      :deep(svg),
      :deep(.anticon) {
        width: 16px;
        height: 16px;
      }
    }

    .more {
      margin-left: auto;
    }
  }
}

@media (max-width: 1400px) {
  .card-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 1120px) {
  .card-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .assistant-page {
    padding: 20px 16px 24px;
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
  }

  .card-grid {
    grid-template-columns: 1fr;
  }

  .page-header .create-btn,
  .search-input {
    width: 100%;
  }
}
</style>
