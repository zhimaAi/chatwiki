<template>
  <div class="mini-card-tab-content">
    <!-- 虚线按钮 -->
    <div class="add-btn" @click="handleOpenSelect">
      <PlusOutlined />
      <span>{{ t('btn_add_mini_card') }}</span>
    </div>
    <span class="tip-text">{{ t('tip_channel_available') }}</span>

    <!-- 已添加卡片列表 -->
    <div class="card-list" v-if="modelValue.length > 0">
      <div
        class="card-item-wrapper"
        v-for="(card, index) in modelValue"
        :key="card.id || index"
      >
        <MiniCardPreview
          :card="card"
          :show-edit="false"
          :show-delete="true"
          @delete="handleDeleteCard(index)"
        />
      </div>
    </div>

    <!-- 选择弹窗 -->
    <MiniCardSelectModal ref="selectModalRef" @select="handleSelectCard" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import MiniCardPreview from './mini-card-preview.vue'
import MiniCardSelectModal from './mini-card-select-modal.vue'

const { t } = useI18n('components.mini-card.mini-card-tab-content')

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const selectModalRef = ref()

const handleOpenSelect = () => {
  selectModalRef.value.show()
}

const handleSelectCard = (card) => {
  const existsIndex = props.modelValue.findIndex((item) => item.id === card.id)
  if (existsIndex !== -1) {
    const list = [...props.modelValue]
    list.splice(existsIndex, 1, card)
    emit('update:modelValue', list)
    return
  }
  emit('update:modelValue', [...props.modelValue, card])
}

const handleDeleteCard = (index) => {
  const list = [...props.modelValue]
  list.splice(index, 1)
  emit('update:modelValue', list)
}
</script>

<style lang="less" scoped>
.mini-card-tab-content {
  padding: 16px 32px 16px 16px;

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

  .tip-text {
    margin-left: 16px;
    font-size: 14px;
    color: #8C8C8C;
    line-height: 22px;
  }
}

.card-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}
</style>
