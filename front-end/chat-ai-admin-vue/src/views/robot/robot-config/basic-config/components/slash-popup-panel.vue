<template>
  <div class="slash-popup-panel" :class="{ 'workspace-slash-panel': workspace }" @keydown.esc.stop="handleClose">
    <!-- 头部：Tab + Close -->
    <div class="popup-header">
      <div class="popup-tabs">
        <component :is="workspace ? 'button' : 'div'" :type="workspace ? 'button' : undefined" :aria-pressed="workspace ? activeTab === 'tag' : undefined"
          class="popup-tab"
          :class="{ active: activeTab === 'tag' }"
          @click="activeTab = 'tag'"
        >
          {{ workspace ? workspaceT('prompt_insert_variable') : t('tab_insert_tag') }}
        </component>
        <component :is="workspace ? 'button' : 'div'" :type="workspace ? 'button' : undefined" :aria-pressed="workspace ? activeTab === 'card' : undefined"
          class="popup-tab"
          :class="{ active: activeTab === 'card' }"
          @click="activeTab = 'card'"
        >
          {{ t('tab_insert_card') }}
        </component>
      </div>
      <div class="popup-close" role="button" tabindex="0" :aria-label="workspaceT('close')" @click="handleClose" @keydown.enter="handleClose">
        <img v-if="workspace" src="@/assets/img/customer-service/basic-settings/dialogs/close.svg" alt="" /><CloseOutlined v-else />
      </div>
    </div>

    <!-- 内容区 -->
    <div class="popup-body">
      <div class="popup-toolbar">
        <div v-if="workspace && activeTab === 'tag'" class="variable-search"><a-input v-model:value="search" :placeholder="workspaceT('prompt_search_variable')" :aria-label="workspaceT('prompt_search_variable')" /><button type="button" @click="handleAddTag"><img src="@/assets/img/customer-service/basic-settings/prompt/add.svg" alt="" />{{ workspaceT('prompt_new_variable') }}</button></div>
        <!-- 提示条 -->
        <div class="popup-alert">
          {{ activeTab === 'tag' ? t('tip_tag') : t('tip_card') }}
        </div>
        <div v-if="activeTab === 'tag' && !workspace" class="add-btn" @click="handleAddTag">
          <PlusOutlined />
          <span>{{ t('btn_add_tag') }}</span>
        </div>
      </div>

      <!-- 插入标签 Tab -->
      <div class="popup-content" v-if="activeTab === 'tag'">
        <div class="tag-list">
          <component :is="workspace ? 'button' : 'div'" :type="workspace ? 'button' : undefined"
            class="tag-item"
            v-for="opt in filteredOptions"
            :key="opt.value || opt.id"
            @click="handleSelectTag(opt)"
          >
            {{ workspace ? opt.name || opt.label : opt.label }}<small v-if="workspace && opt.description">{{ opt.description }}</small>
          </component>
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
import { computed, ref } from 'vue'
import { PlusOutlined, CloseOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useMiniCardStore } from '@/stores/modules/mini-card'
import MiniCardPreview from '@/components/mini-card/mini-card-preview.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.slash-popup-panel')
const { t: workspaceT } = useI18n('views.customer-service.index')

const props = defineProps({
  workspace: { type: Boolean, default: false },
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
const search = ref('')
const filteredOptions = computed(() => props.options.filter(option => !props.workspace || `${option.label} ${option.description || ''}`.toLowerCase().includes(search.value.trim().toLowerCase())))

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
.slash-popup-panel.workspace-slash-panel {
  width: 640px;
  max-width: calc(100vw - 32px);
  border: 1px solid #e8e8e8;
  border-radius: 12px;
  background: #fff;
  overflow: hidden;
  box-shadow: 0 8px 24px #08080814;
  .popup-header { min-height: 56px; padding: 0 16px; border: 0; }
  .popup-tabs { gap: 24px; }
  .popup-tab { padding: 16px 0; border: 0; background: #fff; color: #5a5a5a; font-size: 14px; line-height: 22px; &.active { color: #146ef5; font-weight: 500; } }
  .popup-close { display: flex; align-items: center; justify-content: center; width: 24px; height: 24px; padding: 0; }
  .popup-body { padding: 16px; }
  .popup-toolbar { display: flex; flex-direction: column; align-items: stretch; gap: 12px; margin-bottom: 12px; }
  .variable-search { display: flex; gap: 12px; :deep(.ant-input) { min-width: 0; height: 32px; border-radius: 6px; border-color: #d8d8d8; } button { display: flex; align-items: center; justify-content: center; gap: 6px; flex-shrink: 0; min-height: 32px; padding: 0 12px; border: 1px solid #d8d8d8; border-radius: 6px; background: #fff; color: #222; cursor: pointer; } }
  .popup-alert { width: 100%; padding: 12px 16px; border: 0; border-radius: 6px; background: #f1f6fe; color: #363636; line-height: 22px; }
  .tag-list { display: flex; flex-direction: row; flex-wrap: wrap; gap: 8px; }
  .tag-item { display: inline-flex; align-items: center; gap: 8px; min-height: 32px; margin: 0; padding: 0 12px; border: 0; border-radius: 6px; background: #f1f6fe; color: #146ef5; font-size: 14px; cursor: pointer; small { font-size: 12px; } }
}
</style>
<style lang="less" scoped>
.slash-popup-panel {
  width: 500px;
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
  padding: 16px 16px;
}

.popup-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.popup-alert {
  flex: 1;
  min-width: 0;
  padding: 8px 16px;
  background: #e9f1fe;
  border: 1px solid #99bffd;
  border-radius: 6px;
  font-size: 14px;
  line-height: 22px;
  color: #3a4559;
}

.popup-content {

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
  margin-top: 8px;
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
