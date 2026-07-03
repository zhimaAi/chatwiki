<template>
  <div class="slash-popup-panel">
    <!-- 头部：Tab + Close -->
    <div class="popup-header">
      <div class="popup-tabs">
        <div
          class="popup-tab"
          :class="{ active: activeTab === 'tag' }"
          @click="activeTab = 'tag'"
        >
          {{ t('tab_insert_tag') }}
        </div>
        <div
          class="popup-tab"
          :class="{ active: activeTab === 'card' }"
          @click="activeTab = 'card'"
        >
          {{ t('tab_insert_card') }}
        </div>
      </div>
      <div class="popup-close" @click="handleClose">
        <CloseOutlined />
      </div>
    </div>

    <!-- 内容区 -->
    <div class="popup-body">
      <!-- 提示条 -->
      <div class="popup-alert">
        {{ activeTab === 'tag' ? t('tip_tag') : t('tip_card') }}
      </div>

      <!-- 插入标签 Tab -->
      <div class="popup-content" v-if="activeTab === 'tag'">
        <div class="add-btn" @click="handleAddTag">
          <PlusOutlined />
          <span>{{ t('btn_add_tag') }}</span>
        </div>
        <div class="tag-list">
          <div
            class="tag-item"
            v-for="opt in options"
            :key="opt.value || opt.id"
            @click="handleSelectTag(opt)"
          >
            {{ opt.label }}
          </div>
        </div>
      </div>

      <!-- 插入小程序卡片 Tab -->
      <div class="popup-content" v-else>
        <div class="add-btn" @click="handleAddCard">
          <PlusOutlined />
          <span>{{ t('btn_add_card') }}</span>
        </div>
        <div class="card-list" v-if="cardOptions.length > 0">
          <MiniCardPreview
            v-for="card in cardOptions"
            :key="card.id"
            :card="card"
            :show-edit="true"
            :show-delete="true"
            @click="handleSelectCard"
            @edit="handleEditCard"
            @delete="handleDeleteCard"
          />
        </div>
        <div class="card-empty" v-else>
          {{ t('msg_no_cards') }}
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref } from 'vue'
import { PlusOutlined, CloseOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useMiniCardStore } from '@/stores/modules/mini-card'
import MiniCardPreview from '@/components/mini-card/mini-card-preview.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.slash-popup-panel')

const props = defineProps({
  options: {
    type: Array,
    default: () => []
  },
  cardOptions: {
    type: Array,
    default: () => []
  },
  select: {
    type: Function,
    default: () => {}
  },
  selectCard: {
    type: Function,
    default: () => {}
  },
  close: {
    type: Function,
    default: () => {}
  }
})

const emit = defineEmits(['addTag', 'addCard', 'editCard'])

const miniCardStore = useMiniCardStore()

const activeTab = ref('tag')

const handleClose = () => {
  props.close()
}

const handleSelectTag = (opt) => {
  props.select(opt)
}

const handleSelectCard = (card) => {
  props.selectCard(card)
}

const handleAddTag = () => {
  props.close()
  emit('addTag')
}

const handleAddCard = () => {
  props.close()
  emit('addCard')
}

const handleEditCard = (card) => {
  props.close()
  emit('editCard', card)
}

const handleDeleteCard = (card) => {
  props.close()
  Modal.confirm({
    title: t('msg_confirm_delete'),
    content: t('msg_delete_desc'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    onOk: async () => {
      try {
        await miniCardStore.removeCard({ id: card.id })
        message.success(t('msg_delete_success'))
      } catch (e) {
        if (e && e.data && e.data.error === 'admin_mini_card_in_use') {
          message.error(t('msg_card_in_use'))
        } else {
          message.error(t('msg_delete_failed'))
        }
      }
    }
  })
}

</script>

<style lang="less" scoped>
.slash-popup-panel {
  width: 472px;
  background: #fff;
}

.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  border-bottom: 1px solid #f0f0f0;
}

.popup-tabs {
  display: flex;
  gap: 32px;
}

.popup-tab {
  padding: 16px 0;
  font-size: 16px;
  line-height: 24px;
  color: #595959;
  cursor: pointer;
  position: relative;

  &.active {
    color: #2475fc;
    font-weight: 600;

    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      height: 2px;
      background: #2475fc;
    }
  }
}

.popup-close {
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #595959;
  font-size: 16px;

  &:hover {
    color: #262626;
  }
}

.popup-body {
  padding: 16px 32px 32px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.popup-alert {
  padding: 8px 16px;
  background: #e9f1fe;
  border: 1px solid #99bffd;
  border-radius: 6px;
  font-size: 14px;
  line-height: 22px;
  color: #3a4559;
}

.popup-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.add-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 5px 16px;
  background: #2475fc;
  border-radius: 6px;
  color: #fff;
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;
  align-self: flex-start;

  &:hover {
    opacity: 0.9;
  }
}

.tag-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tag-item {
  display: inline-block;
  padding: 1px 8px;
  border: 1px solid rgba(0, 0, 0, 0.15);
  border-radius: 6px;
  background: #fff;
  font-size: 14px;
  line-height: 22px;
  color: #595959;
  cursor: pointer;
  transition: all 0.2s;
  align-self: flex-start;

  &:hover {
    border-color: #2475fc;
    color: #2475fc;
  }
}

.card-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  max-height: 300px;
  overflow-y: auto;
}

.card-empty {
  text-align: center;
  padding: 24px 0;
  color: #8c8c8c;
  font-size: 14px;
}
</style>
