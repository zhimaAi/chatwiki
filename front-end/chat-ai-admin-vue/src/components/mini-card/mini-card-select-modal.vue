<template>
  <a-modal
    v-model:open="open"
    :title="t('title_select')"
    :width="650"
    :footer="null"
    @cancel="handleClose"
  >
    <div class="mini-card-select">
      <!-- 顶部虚线按钮 -->
      <div class="add-btn-box">
        <div class="add-btn" @click="handleAdd">
          <PlusOutlined />
          <span>{{ t('btn_add_card') }}</span>
        </div>
      </div>

      <!-- 卡片列表 -->
      <div class="card-list-box" ref="scrollBoxRef" @scroll="onScroll">
        <a-spin :spinning="loading">
          <div class="card-list" v-if="cardList.length > 0">
            <MiniCardPreview
              v-for="card in cardList"
              :key="card.id"
              :card="card"
              :show-edit="true"
              :show-delete="true"
              @click="handleSelectCard"
              @edit="handleEditCard"
              @delete="handleDeleteCard"
            />
          </div>
          <div class="empty-box" v-else-if="!loading">
            <span>{{ t('msg_empty') }}</span>
          </div>
        </a-spin>
      </div>
    </div>
  </a-modal>
  <!-- 新增/编辑弹窗 -->
  <MiniCardFormModal ref="formModalRef" @ok="onFormModalOk" />
</template>

<script setup>
import { ref } from 'vue'
import { Modal, message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { useMiniCardStore } from '@/stores/modules/mini-card'
import { useI18n } from '@/hooks/web/useI18n'
import MiniCardPreview from './mini-card-preview.vue'
import MiniCardFormModal from './mini-card-form-modal.vue'

const { t } = useI18n('components.mini-card.mini-card-select-modal')

const emit = defineEmits(['select'])

const open = ref(false)
const loading = ref(false)
const formModalRef = ref()
const scrollBoxRef = ref()

const miniCardStore = useMiniCardStore()

const cardList = ref([])

const loadCards = async (force = false) => {
  loading.value = true
  try {
    await miniCardStore.fetchCardList(force)
    cardList.value = miniCardStore.cardList
  } finally {
    loading.value = false
  }
}

const show = async () => {
  open.value = true
  await loadCards(true)
}

const handleAdd = () => {
  formModalRef.value.show()
}

const handleEditCard = (card) => {
  formModalRef.value.edit(card)
}

const handleDeleteCard = (card) => {
  Modal.confirm({
    title: t('msg_confirm_delete'),
    content: t('msg_delete_desc'),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    onOk: async () => {
      try {
        await miniCardStore.removeCard({ id: card.id })
        message.success(t('msg_delete_success'))
        cardList.value = miniCardStore.cardList
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

const handleSelectCard = (card) => {
  emit('select', card)
  open.value = false
}

const onFormModalOk = () => {
  loadCards(true)
}

const handleClose = () => {
  open.value = false
}

const onScroll = () => {
  // 预留分页加载
}

defineExpose({ show })
</script>

<style lang="less" scoped>
.mini-card-select {
  .add-btn-box {
    margin-bottom: 16px;
  }

  .add-btn {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 5px 16px;
    border: 1px dashed #D9D9D9;
    border-radius: 6px;
    background: #fff;
    color: #595959;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      border-color: #2475FC;
      color: #2475FC;
    }
  }
}

.card-list-box {
  max-height: 400px;
  overflow-y: auto;
}

.card-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.empty-box {
  text-align: center;
  padding: 40px 0;
  color: #8C8C8C;
  font-size: 14px;
}
</style>
