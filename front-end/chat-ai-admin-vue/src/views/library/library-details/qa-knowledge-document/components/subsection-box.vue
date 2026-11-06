<template>
  <div class="subsection-box" :class="{ 'fixed-table-host': fixedHeight, 'compact-list': compact, 'has-selection': state.selectedRowKeys.length > 0 }" ref="tableBox">
    <a-table
      :data-source="paragraphLists"
      row-key="id"
      :row-selection="{ selectedRowKeys: state.selectedRowKeys, onChange: onSelectChange, ...(compact ? { columnWidth: 32 } : {}) }"
      :pagination="false"
      :loading="props.isLoading"
      :scroll="{ x: compact ? tableWidth : 1200, ...(fixedHeight ? { y: scrollY } : {}) }"
      :table-layout="compact ? 'fixed' : undefined"
      @change="tableChange"
    >
      <a-table-column key="id" data-index="id" :width="compact ? getColumnWidth('id') : 1148" class-name="qa-content-cell">
        <template #title>
          {{ t('title_qa_total', { total: props.total }) }}
          <ColumnWidthHandle v-if="compact" :width="getColumnWidth('id')" :min-width="240" :label="t('resize_column', { column: t('label_question') })" @resize="columnWidths.id = $event" @reset="delete columnWidths.id" />
        </template>
        <template #default="{ record }">
          <div class="qa-list-box" @dblclick="handleOpenEditModal(record)">
            <div class="list-item">
              <div v-if="!compact" class="list-label">{{ t('label_question') }}</div>
              <div class="list-content">
                <a-tooltip placement="top" v-if="record.status == 2">
                  <template #title>
                    <span>{{ record.errmsg }}</span>
                  </template>
                  <span>
                    <span class="status-error"><ExclamationCircleFilled /> </span>
                  </span>
                </a-tooltip>
                <span @click="handleOpenEditModal(record)" @keydown.enter.prevent="handleOpenEditModal(record)" :role="compact ? 'button' : undefined" :tabindex="compact ? 0 : undefined" v-html="textToHighlight(record.question, props.search)" class="question-text"></span>
                <template v-if="record.similar_questions && record.similar_questions.length">
                  <a-popover placement="topLeft" :overlayInnerStyle="{ 'padding-right': '2px' }">
                    <template #content>
                      <div class="similar-question-list">
                        <div
                          class="similar-question-list-item"
                          v-for="(item, index) in record.similar_questions"
                          :key="index"
                        >
                          {{ item }}
                        </div>
                      </div>
                    </template>
                    <template #title>
                      <span>{{ t('msg_similar_questions', { length: record.similar_questions.length }) }}</span>
                    </template>
                    <a>（{{ record.similar_questions.length }}）</a>
                  </a-popover>
                </template>
              </div>
            </div>
            <div class="list-item list-item-answer">
              <div v-if="!compact" class="list-label">{{ t('label_answer') }}</div>
              <div class="list-content" v-html="textToHighlight(record.answer, props.search)"></div>
            </div>
            <div v-if="record.images?.length" class="fragment-img">
              <div class="media-item" v-for="(item, index) in (record.images || [])" :key="index">
                <div
                  v-if="isVideoUrl(item)"
                  class="video-item"
                  @click.stop="handlePreviewVideo(item)"
                >
                  <video :src="item" muted preload="metadata"></video>
                  <div class="video-play-mask">
                    <CaretRightOutlined />
                  </div>
                </div>
                <div v-else class="image-item" v-viewer>
                  <img :src="item" alt="" />
                </div>
              </div>
            </div>
          </div>
        </template>
      </a-table-column>
      <a-table-column v-for="hit in hitColumns" :key="hit.key" :data-index="hit.key" :width="compact ? getColumnWidth(hit.key) : 108" :sorter="true" class-name="qa-hit-cell">
        <template #title>
          <span class="hit-column-title">{{ t(hit.title) }}<span v-if="compact" class="workspace-sort-icons" :class="{ 'sort-ascending': sortState.field === hit.key && sortState.order === 'ascend', 'sort-descending': sortState.field === hit.key && sortState.order === 'descend' }"><img class="sort-up" :src="icons.sortUp" alt="" /><img class="sort-down" :src="icons.sortDown" alt="" /></span></span>
          <ColumnWidthHandle v-if="compact" :width="getColumnWidth(hit.key)" :label="t('resize_column', { column: t(hit.title) })" @resize="columnWidths[hit.key] = $event" @reset="delete columnWidths[hit.key]" />
        </template>
      </a-table-column>
      <a-table-column
        v-for="meta in metaData"
        :key="`meta_${meta.key}`"
        :data-index="`meta_${meta.key}`"
        :width="compact ? getColumnWidth(`meta_${meta.key}`) : 160"
      >
        <template #title>
          <a-tooltip :title="meta.name?.length > 8 ? meta.name : '' ">
            <div class="zm-line1">{{meta.name}}</div>
          </a-tooltip>
          <ColumnWidthHandle v-if="compact" :width="getColumnWidth(`meta_${meta.key}`)" :label="t('resize_column', { column: meta.name })" @resize="columnWidths[`meta_${meta.key}`] = $event" @reset="delete columnWidths[`meta_${meta.key}`]" />
        </template>
        <template #default="{record}">
          <a-tooltip :title="record[`meta_${meta.key}`]?.length > 8 ? record[`meta_${meta.key}`] : '' ">
            <div class="zm-line1">{{record[`meta_${meta.key}`]}}</div>
          </a-tooltip>
        </template>
      </a-table-column>
      <a-table-column :title="t('col_action')" key="action" data-index="action" :width="compact ? 56 : 120" fixed="right" class-name="qa-action-cell">
        <template #default="{ record }">
          <a-dropdown v-if="compact" :trigger="['click']" placement="bottomRight">
            <button type="button" class="qa-more-button" :aria-label="t('more_actions')" :title="t('more_actions')" @click.stop><img :src="icons.more" alt="" /></button>
            <template #overlay>
              <a-menu>
                <a-menu-item key="edit" @click="handleOpenEditModal(record)">{{ t('edit') }}</a-menu-item>
                <a-menu-item key="delete" danger @click="hanldleDelete(record.id)">{{ t('delete') }}</a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
          <div v-else class="right-opration" >
            <template v-if="!compact">
              <div class="hover-btn-box" @click="handleSetCategory(record, 0)" v-if="record.category_id > 0">
                <StarFilled style="color: #F59A23;" />
              </div>
              <div class="hover-btn-box" v-else @click="handleSetCategory(record, 4)">
                <StarOutlined />
              </div>
            </template>
            <div class="hover-btn-box" @click.stop="handleOpenEditModal(record)">
              <EditOutlined />
            </div>
            <div class="hover-btn-box" @click.stop="hanldleDelete(record.id)">
              <DeleteOutlined />
            </div>
          </div>
        </template>
      </a-table-column>
    </a-table>
  </div>
  <VideoPreviewModal ref="videoPreviewModalRef" />
</template>
<script setup>
import { reactive, computed, createVNode, ref, watch, nextTick, onBeforeUnmount } from 'vue'
import { message } from 'ant-design-vue'
import {
  ExclamationCircleOutlined,
  ExclamationCircleFilled,
  EditOutlined,
  DeleteOutlined,
  StarOutlined,
  StarFilled,
  CaretRightOutlined
} from '@ant-design/icons-vue'
import { Modal } from 'ant-design-vue'
import ColumnWidthHandle from './column-width-handle.vue'
import icons from '../knowledge-icons'
import VideoPreviewModal from '@/components/upload-img/video-preview-modal.vue'
import { useTableScroll } from '../../components/use-table-scroll'
import {
  deleteParagraph,
  updateParagraphCategory
} from '@/api/library'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-details.qa-knowledge-document.components.subsection-box')

const emit = defineEmits([
  'handleDelParagraph',
  'handleScrollTargetPage',
  'openEditSubscription',
  'handleSort',
  'getList'
])
const props = defineProps({
  compact: { type: Boolean, default: false },
  fixedHeight: { type: Boolean, default: false },
  paragraphLists: {
    type: Array,
    default: () => []
  },
  total: {
    type: [Number, String],
    default: 0
  },
  isLoading: {
    type: Boolean,
    default: false
  },
  search: {
    type: String,
    default: ''
  }
})
const tableBox = ref(null)
const { scrollY } = useTableScroll(tableBox, () => props.fixedHeight)

const metaData = computed(() => {
  return props.paragraphLists?.[0]?.meta_list || []
})

const hitColumns = [
  { key: 'total_hits', title: 'col_total' },
  { key: 'yesterday_hits', title: 'col_yesterday' },
  { key: 'today_hits', title: 'col_today' }
]
const columnWidths = reactive({})
const availableWidth = ref(932)
const sortState = ref({})
const getColumnWidth = key => columnWidths[key] ?? (key === 'id' ? Math.max(240, availableWidth.value - 32 - 56 - hitColumns.reduce((sum, hit) => sum + (columnWidths[hit.key] ?? 72), 0)) : key.startsWith('meta_') ? 160 : 72)
const tableWidth = computed(() => 32 + 56 + getColumnWidth('id') + hitColumns.reduce((sum, hit) => sum + getColumnWidth(hit.key), 0) + metaData.value.reduce((sum, meta) => sum + getColumnWidth(`meta_${meta.key}`), 0))
let widthObserver
const measureWidth = () => {
  if (!props.compact || !tableBox.value) return
  const body = tableBox.value.querySelector('.ant-table-body')
  const scrollbarWidth = body ? Math.max(0, body.offsetWidth - body.clientWidth) : 0
  if (tableBox.value.clientWidth) availableWidth.value = tableBox.value.clientWidth - scrollbarWidth
}
watch(() => props.compact, async (enabled, _, onCleanup) => {
  let active = true
  onCleanup(() => { active = false; widthObserver?.disconnect(); window.removeEventListener('resize', measureWidth) })
  if (!enabled) return
  await nextTick()
  if (!active) return
  measureWidth()
  if (typeof ResizeObserver === 'function') {
    widthObserver = new ResizeObserver(measureWidth)
    widthObserver.observe(tableBox.value)
  } else window.addEventListener('resize', measureWidth)
}, { immediate: true })
watch(() => props.paragraphLists, async () => { await nextTick(); measureWidth() })
onBeforeUnmount(() => widthObserver?.disconnect())

const state = reactive({
  selectedRowKeys: []
})
const videoPreviewModalRef = ref(null)
const onSelectChange = (selectedRowKeys) => {
  state.selectedRowKeys = selectedRowKeys
}

const resetSelect = () => {
  state.selectedRowKeys = []
}

const handleOpenEditModal = (item) => {
  emit('openEditSubscription', item)
}
const hanldleDelete = (id) => {
  Modal.confirm({
    title: t('title_hint'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_confirm_delete'),
    onOk() {
      return new Promise((resolve, reject) => {
        deleteParagraph({ id })
          .then((res) => {
            message.success(t('msg_delete_success'))
            emit('handleDelParagraph', id)
            resolve()
          })
          .catch(() => {
            resolve()
          })
      })
    },
    onCancel() {}
  })
}

const handleSetCategory = (item, category_id) => {
  updateParagraphCategory({
    id: item.id,
    category_id
  }).then((res) => {
    message.success(t('msg_set_success'))
    emit('getList')
  })
}

const tableChange = (a, b, sort) => {
  sortState.value = sort
  let sort_field = sort.field
  let sort_type = sort.order
  if (sort_type) {
    sort_type = sort_type === 'ascend' ? 'asc' : 'desc'
  } else {
    sort_field = ''
  }
  emit('handleSort', {
    sort_field,
    sort_type
  })
}

const isVideoUrl = (url) => {
  if (!url || typeof url !== 'string') return false
  return /\.mp4(\?|#|$)/i.test(url) || /library_video/i.test(url) || /\/video\//i.test(url)
}

const handlePreviewVideo = (url) => {
  videoPreviewModalRef.value?.show(url)
}

function textToHighlight(fullText, highlightText, options = {}) {
  if (!highlightText || !fullText) return fullText;

  const {
    highlightClass = 'highlight',
    caseSensitive = false,
    wholeWord = false
  } = options;

  const flags = caseSensitive ? 'g' : 'gi';
  let regexPattern;

  if (wholeWord) {
    // 使用单词边界匹配完整单词
    regexPattern = new RegExp(`\\b${escapeRegExp(highlightText)}\\b`, flags);
  } else {
    regexPattern = new RegExp(escapeRegExp(highlightText), flags);
  }

  return fullText.replace(regexPattern, match =>
    `<span class="${highlightClass}">${match}</span>`
  );
}

/**
 * 转义正则表达式特殊字符
 * @param {string} string
 * @returns {string}
 */
function escapeRegExp(string) {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

defineExpose({ handleOpenEditModal, state, resetSelect })
</script>
<style lang="less" scoped>
.subsection-box {
  width: 100%;
  .subsection-box-title {
    display: flex;
    align-items: center;
    font-size: 14px;
    line-height: 22px;
    font-weight: 600;
    color: #242933;
    span {
      color: #7a8699;
      font-weight: 400;
      margin-left: 8px;
    }
  }
}

.qa-list-box {
  .question-text {
    cursor: pointer;
    &:hover {
      color: #3475fc;
    }
  }
  .list-item {
    display: flex;
    flex-wrap: wrap;
    line-height: 22px;
    font-size: 14px;
    color: #262626;
    margin-bottom: 6px;
    .list-label {
      margin-right: 12px;
    }
    .list-content {
      flex: 1;
      display: flex;
      flex-wrap: wrap;
      gap: 4px;
      word-break: break-all;
      &::v-deep(.highlight) {
        background-color: #FFEB3B; /* 黄色高亮 */
        padding: 0 2px;
        border-radius: 2px;
        box-shadow: 0 1px 1px rgba(0,0,0,0.1);
      }
    }
  }
  .list-item-answer {
    color: #8c8c8c;
  }
  .fragment-img {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    padding-left: 40px;
    .media-item {
      width: 80px;
      height: 80px;
    }
    .image-item,
    .video-item {
      width: 80px;
      height: 80px;
    }
    img,
    video {
      width: 80px;
      height: 80px;
      border-radius: 6px;
      cursor: pointer;
      object-fit: cover;
    }
    .video-item {
      position: relative;
      cursor: pointer;
      .video-play-mask {
        position: absolute;
        left: 0;
        right: 0;
        top: 0;
        bottom: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
        background: rgba(0, 0, 0, 0.28);
        border-radius: 6px;
        font-size: 18px;
        transition: all 0.2s;
      }
      &:hover .video-play-mask {
        background: rgba(0, 0, 0, 0.4);
      }
    }
  }
}

.status-error {
  cursor: pointer;
  color: #fb363f;
}

.right-opration {
  display: flex;
  align-items: center;
  gap: 8px;
}
.hover-btn-box {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
  &:hover {
    background: #e4e6eb;
    border-radius: 6px;
  }
}

.similar-question-list {
  padding-right: 16px;
  width: 260px;
  max-height: 280px;
  overflow-y: auto;
  padding-top: 4px;
  .similar-question-list-item {
    margin-bottom: 8px;
    color: #000;
    font-size: 14px;
    line-height: 22px;
  }
  &::-webkit-scrollbar {
    width: 6px; /* 垂直滚动条宽度 */
    height: 6px; /* 水平滚动条高度 */
  }

  /* 滚动条轨道 */
  &::-webkit-scrollbar-track {
    background: #f1f1f1;
    border-radius: 10px;
  }

  /* 滚动条滑块 */
  &::-webkit-scrollbar-thumb {
    background: #888;
    border-radius: 10px;
    transition: background 0.3s ease;
  }

  /* 滚动条滑块悬停状态 */
  &::-webkit-scrollbar-thumb:hover {
    background: #555;
  }

  /* 滚动条角落 */
  &::-webkit-scrollbar-corner {
    background: #f1f1f1;
  }
}

.compact-list {
  .qa-list-box {
    .list-item { flex-wrap: nowrap; margin-bottom: 4px; }
    .list-content { min-width: 0; }
    .list-item:first-child .list-content { flex-wrap: nowrap; }
    .question-text { display: block; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #262626; }
    .question-text:focus-visible { outline: 2px solid #146ef5; outline-offset: -2px; }
    .list-item-answer { margin-bottom: 0; }
    .list-item-answer .list-content { display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2; min-height: 44px; max-height: 44px; overflow: hidden; color: #999; }
    .list-item-answer :deep(p) { margin: 0; }
    .fragment-img {
      padding-left: 0; margin-top: 4px;
      .media-item, .image-item, .video-item, img, video { width: 64px; height: 64px; }
    }
  }
  .hit-column-title { display: inline-flex; align-items: center; gap: 8px; white-space: nowrap; }
  .workspace-sort-icons {
    display: flex; flex-direction: column; flex-shrink: 0;
    .sort-up { transform: rotate(180deg); }
    &.sort-ascending .sort-up, &.sort-descending .sort-down { filter: brightness(0) saturate(100%) invert(35%) sepia(97%) saturate(2775%) hue-rotate(206deg) brightness(100%) contrast(94%); }
  }
  .qa-more-button {
    display: inline-flex; align-items: center; justify-content: center; width: 32px; height: 32px; padding: 0; border: 0; border-radius: 6px; background: transparent; cursor: pointer;
    &:hover { background: #e8e8e8; }
    &:focus-visible { outline: 2px solid #146ef5; }
  }
  :deep(.ant-table) { font-size: 14px; line-height: 22px; border-radius: 0; }
  :deep(.ant-table-container), :deep(.ant-table-container table > thead > tr:first-child > *:first-child), :deep(.ant-table-container table > thead > tr:first-child > *:last-child) { border-radius: 0; }
  :deep(.ant-table-thead > tr > th) { height: 48px; padding: 0 8px; background: #fff; border-bottom-color: #d8d8d8; color: #5a5a5a; font-weight: 400; }
  :deep(.ant-table-thead > tr > th::before) { display: none; }
  :deep(.ant-table-column-sorter) { display: none; }
  :deep(.ant-table-column-sorters) { position: static; justify-content: flex-start; }
  :deep(.ant-table-column-title) { position: static; }
  :deep(.ant-table-tbody > tr > td) { padding: 12px 8px; border-bottom-color: #d8d8d8; background: #fff; vertical-align: middle; }
  :deep(.ant-table-tbody > tr > td.ant-table-selection-column) { padding: 12px 4px; vertical-align: top; }
  :deep(.ant-table-thead > tr > th.ant-table-selection-column) { padding: 0 4px; }
  :deep(.ant-table-tbody > tr > td.qa-content-cell), :deep(.ant-table-thead > tr > th.qa-content-cell) { padding-left: 0; }
  :deep(.ant-table-tbody > tr > td.qa-hit-cell) { padding-left: 16px; }
  :deep(.ant-table-tbody > tr > td.qa-action-cell) { padding: 0 12px; }
  :deep(.ant-table-selection-column .ant-checkbox-wrapper) { opacity: 0; }
  :deep(.ant-table-tbody > tr:hover > td .ant-checkbox-wrapper),
  :deep(.ant-table-tbody > tr:focus-within > td .ant-checkbox-wrapper),
  :deep(.ant-table-thead:focus-within .ant-checkbox-wrapper),
  :deep(.ant-table-thead:hover .ant-checkbox-wrapper) { opacity: 1; }
  &.has-selection :deep(.ant-table-selection-column .ant-checkbox-wrapper) { opacity: 1; }
  :deep(.ant-table-tbody > tr:hover > td), :deep(.ant-table-tbody > tr:focus-within > td) { background: #f5f5f5; border-bottom-color: transparent; }
  :deep(.ant-table-tbody > tr.ant-table-row-selected > td), :deep(.ant-table-tbody > tr.ant-table-row-selected:hover > td) { background: #f0f5ff; border-bottom-color: transparent; }
  :deep(.ant-table-tbody > tr:has(+ tr:hover) > td), :deep(.ant-table-tbody > tr:has(+ tr:focus-within) > td), :deep(.ant-table-tbody > tr:has(+ tr.ant-table-row-selected) > td) { border-bottom-color: transparent; }
  :deep(.ant-table-tbody > tr:hover > td:first-child), :deep(.ant-table-tbody > tr:focus-within > td:first-child), :deep(.ant-table-tbody > tr.ant-table-row-selected > td:first-child) { border-radius: 6px 0 0 6px; }
  :deep(.ant-table-tbody > tr:hover > td:last-child), :deep(.ant-table-tbody > tr:focus-within > td:last-child), :deep(.ant-table-tbody > tr.ant-table-row-selected > td:last-child) { border-radius: 0 6px 6px 0; }
  :deep(.ant-table-cell-fix-right-first::after) { box-shadow: none; }
  @media (hover: none) { :deep(.ant-table-selection-column .ant-checkbox-wrapper) { opacity: 1; } }
}
</style>
