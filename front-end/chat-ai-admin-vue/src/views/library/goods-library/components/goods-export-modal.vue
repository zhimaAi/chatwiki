<template>
  <a-modal
    v-model:open="open"
    :title="t('export_modal.title')"
    :ok-text="t('export_modal.export')"
    :cancel-text="t('export_modal.cancel')"
    :width="480"
    @ok="handleOk"
    :destroyOnClose="true"
  >
    <div class="export-info">
      <p>{{ t('export_modal.confirm_text') }}</p>
      <p class="info-line">
        <span class="info-label">{{ t('export_modal.scope_label') }}</span>
        <span>{{ currentScopeText }}</span>
      </p>
      <p class="info-line" v-if="keyword">
        <span class="info-label">{{ t('export_modal.keyword_label') }}</span>
        <span>{{ keyword }}</span>
      </p>
    </div>
  </a-modal>
</template>

<script setup>
import { computed, ref } from 'vue'
import { message } from 'ant-design-vue'
import { exportGoodsList } from '@/api/goods-library'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.goods-library.index')

const props = defineProps({
  groupId: {
    type: [String, Number],
    default: '0'
  },
  keyword: {
    type: String,
    default: ''
  }
})

const open = ref(false)

const currentScopeText = computed(() => {
  const id = String(props.groupId)
  if (id === 'all' || id === '-1') {
    return t('export_modal.scope_all')
  }
  if (id === '0') {
    return t('group_tree.ungrouped')
  }
  return t('export_modal.scope_group')
})

const show = () => {
  open.value = true
}

const handleOk = () => {
  const id = String(props.groupId)
  let groupId = -1
  if (id === 'all' || id === '-1') {
    groupId = -1
  } else if (id === '0') {
    groupId = 0
  } else {
    groupId = Number(id)
  }

  exportGoodsList({
    group_id: groupId,
    keyword: props.keyword,
    switch_status: -1
  })
  message.success(t('message.export_triggered'))
  open.value = false
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.export-info {
  margin-top: 4px;

  p {
    margin: 0 0 12px;
    color: #262626;
    font-size: 14px;
    line-height: 22px;
  }
}

.info-line {
  display: flex;
  gap: 8px;
}

.info-label {
  flex-shrink: 0;
  color: #8c8c8c;
}
</style>
