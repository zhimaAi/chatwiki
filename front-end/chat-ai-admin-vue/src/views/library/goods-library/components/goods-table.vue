<template>
  <div ref="tableRef" class="goods-table">
    <a-table
      :columns="columns"
      :data-source="rows"
      :loading="loading"
      :pagination="tablePagination"
      @change="handleTableChange"
      :row-key="(record) => record.id"
      :scroll="{ x: 1600}"
      :custom-row="getCustomRow"
      :row-class-name="getRowClassName"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'image'">
          <div
            class="cell-inner image-inner"
            :class="isActiveCell(record.id, 'image') ? 'cell-active' : ''"
            @click.stop="emit('edit-field', { row: record, fieldKey: 'images', fieldLabel: t('table.image'), mode: 'images' })"
          >
            <div class="image-box">
              <img v-if="getImageUrl(record)" :src="getImageUrl(record)" alt="" />
              <div v-else class="image-empty">{{ t('table.no_image') }}</div>
              <div v-if="getImageCount(record) > 1" class="image-count">
                {{ t('table.image_count', { count: getImageCount(record) }) }}
              </div>
            </div>
          </div>
        </template>

        <template v-else-if="column.key === 'basic_info'">
          <div
            class="cell-inner basic-inner"
            :class="isActiveCell(record.id, 'basic_info') ? 'cell-active' : ''"
            @click.stop="emit('edit-basic-info', record)"
          >
            <div class="info-grid">
              <div class="info-col" style="width: 280px;">
                <div class="info-row">
                  <span class="info-label">ID：</span>
                  <span class="info-value">{{ record.goods_id || record.id || '-' }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">{{ t('table.name_label') }}</span>
                  <span class="info-value">{{ record.goods_name || record.name || '-' }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">{{ t('table.category_label') }}</span>
                  <span class="info-value">{{ record.category || '-' }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">{{ t('table.brand_label') }}</span>
                  <span class="info-value">{{ record.brand || '-' }}</span>
                </div>
              </div>
              <div class="info-col" style="flex: 1;">
                <div class="info-row">
                  <span class="info-label">{{ t('table.price_label') }}</span>
                  <span class="info-value">{{ record.price || '-' }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">{{ t('table.stock_label') }}</span>
                  <span class="info-value">{{ record.stock || '-' }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">{{ t('table.link_label') }}</span>
                  <span class="info-value">
                    <a
                      v-if="record.link"
                      class="goods-link"
                      :href="record.link"
                      target="_blank"
                      :title="record.link"
                      @click.stop
                    >
                      {{ t('table.link_text') }}
                    </a>
                    <span v-else>-</span>
                  </span>
                </div>
              </div>
            </div>
          </div>
        </template>

        <template v-else-if="column.key === 'description'">
          <div
            class="cell-inner text-inner"
            :class="isActiveCell(record.id, 'description') ? 'cell-active' : ''"
            @click.stop="emit('edit-field', { row: record, fieldKey: 'description', fieldLabel: t('table.description'), mode: 'textarea' })"
          >
            <a-tooltip
              v-if="record.description"
              :mouse-enter-delay="0.3"
              :overlay-style="{ maxWidth: '400px' }"
              overlay-class-name="goods-cell-tooltip"
            >
              <template #title>
                <div class="tooltip-pre-line">{{ record.description }}</div>
              </template>
              <div class="line-clamp">
                {{ record.description }}
              </div>
            </a-tooltip>
            <div v-else class="line-clamp placeholder">
              {{ t('table.description_placeholder') }}
            </div>
          </div>
        </template>

        <template v-else-if="column.key === 'qa'">
          <div
            class="cell-inner text-inner"
            :class="isActiveCell(record.id, 'qa') ? 'cell-active' : ''"
            @click.stop="emit('edit-field', { row: record, fieldKey: 'qa', fieldLabel: t('table.qa'), mode: 'textarea' })"
          >
            <a-tooltip
              v-if="record.qa"
              :mouse-enter-delay="0.3"
              :overlay-style="{ maxWidth: '400px' }"
              overlay-class-name="goods-cell-tooltip"
            >
              <template #title>
                <div class="tooltip-pre-line">{{ record.qa }}</div>
              </template>
              <div class="line-clamp">
                {{ record.qa }}
              </div>
            </a-tooltip>
            <div v-else class="line-clamp placeholder">
              {{ t('table.qa_placeholder') }}
            </div>
          </div>
        </template>

        <template v-else-if="column.key === 'custom_info'">
          <div
            class="cell-inner text-inner"
            :class="isActiveCell(record.id, 'custom_info') ? 'cell-active' : ''"
            @click.stop="emit('edit-field', { row: record, fieldKey: 'custom_info', fieldLabel: t('custom_info.title'), mode: 'custom_info' })"
          >
            <a-tooltip
              v-if="record.custom_info"
              :mouse-enter-delay="0.3"
              :overlay-style="{ maxWidth: '400px' }"
              overlay-class-name="goods-cell-tooltip"
            >
              <template #title>
                <div class="tooltip-pre-line">{{ record.custom_info }}</div>
              </template>
              <div class="line-clamp">
                {{ record.custom_info }}
              </div>
            </a-tooltip>
            <div v-else class="line-clamp placeholder">
              {{ t('custom_info.empty') }}
            </div>
          </div>
        </template>

        <template v-else-if="column.key === 'actions'">
          <div class="actions-cell" @click.stop>
            <a-switch
              :checked="record.switch_status === 1"
              :checked-children="t('table.enabled')"
              :un-checked-children="t('table.disabled')"
              @change="(checked) => emit('toggle-status', { row: record, checked })"
            />
            <a-button type="text" class="delete-btn" @click.stop="emit('delete-row', record)">
              <DeleteOutlined />
            </a-button>
          </div>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { DeleteOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.goods-library.index')

const props = defineProps({
  rows: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  selectedRowId: {
    type: [String, Number],
    default: ''
  },
  activeCell: {
    type: Object,
    default: () => ({})
  },
  pagination: {
    type: Object,
    default: () => ({ page: 1, size: 20, total: 0 })
  }
})

const emit = defineEmits(['hover-cell', 'edit-field', 'toggle-status', 'delete-row', 'select-row', 'edit-row', 'edit-basic-info', 'change'])

const tableRef = ref(null)
const scrollY = ref(400)
let resizeObserver = null

const HEADER_HEIGHT = 54
const PAGINATION_HEIGHT = 56
const SAFE_GAP = 8

onMounted(() => {
  const container = tableRef.value?.closest('.main-content')
  if (!container) return

  resizeObserver = new ResizeObserver((entries) => {
    for (const entry of entries) {
      const height = entry.contentRect.height
      scrollY.value = Math.max(height - HEADER_HEIGHT - PAGINATION_HEIGHT - SAFE_GAP, 200)
    }
  })
  resizeObserver.observe(container)
})

onBeforeUnmount(() => {
  resizeObserver?.disconnect()
})

const tablePagination = computed(() => ({
  current: props.pagination.page,
  pageSize: props.pagination.size,
  total: props.pagination.total,
  showSizeChanger: true,
  showQuickJumper: true,
  pageSizeOptions: ['10', '20', '50'],
  size: 'default'
}))

const handleTableChange = (pagination) => {
  emit('change', {
    page: pagination.current,
    size: pagination.pageSize
  })
}

const columns = computed(() => [
  {
    title: t('table.image'),
    key: 'image',
    width: 167
  },
  {
    title: t('table.basic_info'),
    key: 'basic_info',
    width: 480
  },
  {
    title: t('table.description'),
    key: 'description',
    width: 334
  },
  {
    title: t('table.qa'),
    key: 'qa',
    width: 263
  },
  {
    title: t('custom_info.title'),
    key: 'custom_info',
    width: 383
  },
  {
    title: t('table.actions'),
    key: 'actions',
    width: 115,
    fixed: 'right',
    className: 'actions-fixed-column'
  }
])


const normalizeId = (value) => {
  if (value === undefined || value === null || value === '') {
    return ''
  }

  return String(value)
}

const getCustomRow = (record) => ({
  onClick: () => {
    emit('select-row', record)
  },
  onDblclick: () => {
    emit('edit-row', record)
  }
})

const getRowClassName = (record) => {
  return normalizeId(record.id) === normalizeId(props.selectedRowId) ? 'row-selected' : ''
}

const isActiveCell = (rowId, field) => {
  return props.activeCell?.rowId === rowId && props.activeCell?.field === field
}

const getImageList = (record) => {
  if (Array.isArray(record.images) && record.images.length) {
    return record.images
  }

  return []
}

const getImageUrl = (record) => {
  const images = getImageList(record)
  return images[0] || ''
}

const getImageCount = (record) => {
  return getImageList(record).length
}
</script>

<style lang="less" scoped>
.goods-table {
  width: 100%;
  height: 100%;
  overflow: hidden;
  overflow-y: auto;

  :deep(.ant-table-cell) {
    padding: 4px;
    background: #fff;
    border-bottom: 1px solid #e8e8e8;
  }

  :deep(.ant-table-thead > tr > th) {
    height: 54px;
    padding: 16px;
    background: #f5f5f5;
    color: #262626;
    font-size: 14px;
    line-height: 20px;
    font-weight: 400;
  }

  :deep(.ant-table-tbody > tr) {
    cursor: pointer;
  }

  :deep(.ant-table-tbody > tr > td) {
    padding: 4px;
  }

  :deep(.ant-table-tbody > tr:hover > td) {
    background: #fff !important;
  }

  :deep(.ant-table-placeholder .ant-table-cell) {
    padding: 80px 0;
    border-bottom: none;
  }

  :deep(.ant-table-cell-fix-right),
  :deep(.ant-table-cell-fix-right-first) {
    z-index: 3;
    background: #fff;
  }

  :deep(.ant-table-thead > tr > .ant-table-cell-fix-right),
  :deep(.ant-table-thead > tr > .ant-table-cell-fix-right-first) {
    background: #f5f5f5;
  }

  :deep(.ant-table-thead > tr > .actions-fixed-column),
  :deep(.ant-table-tbody > tr > .actions-fixed-column) {
    box-shadow: -8px 0 12px rgba(0, 0, 0, 0.06);
  }

  :deep(.ant-table-pagination.ant-pagination) {
    margin: 12px 0 0;
    padding: 0 4px;
  }

  .cell-inner {
    position: relative;
    width: 100%;
    height: 100%;
    padding: 12px;
    border-radius: 6px;
    background: #fff;
    transition: background-color 0.2s;
    &:hover {
      background: #E5EFFF;
    }
    &.cell-active {
      
    }
  }

  .image-inner {
    width: 158px;
    height: 116px;
    overflow: hidden;
  }

  .image-box {
    position: relative;
    width: 135px;
    height: 92px;
    overflow: hidden;
    border-radius: 8px;
    background: #fafafa;

    img {
      display: block;
      width: 100%;
      max-width: 135px;
      height: 100%;
      max-height: 92px;
    }
  }

  .image-empty {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    color: #bfbfbf;
    font-size: 12px;
  }

  .image-count {
    position: absolute;
    left: 8px;
    bottom: 8px;
    height: 20px;
    line-height: 20px;
    padding: 0 8px;
    border-radius: 10px;
    background: rgba(0, 0, 0, 0.55);
    color: #fff;
    font-size: 12px;
  }

  .basic-inner {
    height: 116px;
    overflow: hidden;
    color: #262626;
    font-size: 14px;
    line-height: 20px;
  }

  .info-grid {
    display: flex;
    gap: 0 12px;
  }

  .info-col {
    display: flex;
    flex-direction: column;
    gap: 4px;
    min-width: 0;
  }

  .info-row {
    display: flex;
    gap: 4px;
    align-items: flex-start;
    min-width: 0;
  }

  .info-label {
    flex-shrink: 0;
    width: 42px;
    color: #8c8c8c;
    white-space: nowrap;
  }

  .info-value {
    flex: 1;
    min-width: 0;
    color: #262626;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .goods-link {
    color: #2475fc;
    text-decoration: none;

    &:hover {
      text-decoration: underline;
    }
  }

  .text-inner {
    height: 116px;
    overflow: hidden;
  }

  .line-clamp {
    color: #595959;
    font-size: 14px;
    line-height: 20px;
    word-break: break-word;
    display: -webkit-box;
    -webkit-line-clamp: 4;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;

    &.placeholder {
      color: #bfbfbf;
    }
  }

  .actions-cell {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    min-height: 122px;
    padding: 12px;
  }

  .delete-btn {
    width: 24px;
    height: 24px;
    padding: 0;
    color: #595959;
  }
}
</style>

<style lang="less">
.goods-cell-tooltip {
  .tooltip-pre-line {
    max-width: 400px;
    max-height: 300px;
    overflow-y: auto;
    word-break: break-word;
    white-space: pre-line;

    &::-webkit-scrollbar {
      width: 5px;
    }

    &::-webkit-scrollbar-track {
      background: transparent;
    }

    &::-webkit-scrollbar-thumb {
      background: rgba(255, 255, 255, 0.3);
      border-radius: 3px;

      &:hover {
        background: rgba(255, 255, 255, 0.5);
      }
    }
  }
}
</style>
