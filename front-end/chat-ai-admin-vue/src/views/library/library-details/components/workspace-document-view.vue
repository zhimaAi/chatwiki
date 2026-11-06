<template>
  <div class="workspace-document-view" :class="{ 'groups-collapsed': collapsed }">
    <aside v-if="!collapsed" class="document-groups" :style="{ width: `${groupWidth}px` }">
      <div class="groups-heading"><button type="button" class="icon-button" :aria-label="commonT('document_collapse_groups')" @click="emit('toggle-groups')"><img :src="icons.collapse" alt="" /></button><h2>{{ t('label_knowledge_group') }}</h2><button type="button" class="icon-button" :aria-label="t('tooltip_create_group')" @click="emit('group-action', { key: 'add' })"><img :src="icons.addGroup" alt="" /></button></div>
      <a-input v-model:value="groupSearch" class="group-search" :placeholder="t('ph_search_group')" :aria-label="t('ph_search_group')" allow-clear><template #prefix><img :src="icons.search" alt="" /></template></a-input>
      <Draggable :model-value="filteredGroups" item-key="id" handle=".group-drag" :disabled="!!groupSearch" :move="checkGroupMove" class="document-group-list" @update:model-value="emit('reorder-groups', $event)">
        <template #item="{ element }"><div class="document-group-row" :class="{ active: String(element.id) === String(groupId) }"><button type="button" class="group-select" :aria-pressed="String(element.id) === String(groupId)" @click="chooseGroup(element)"><span>{{ element.group_name }}</span><span class="group-count">{{ element.total || 0 }}</span></button><span v-if="element.id > 0" class="group-drag" :title="t('btn_rename')"><svg-icon name="drag" /></span><a-dropdown v-if="element.id > 0" :trigger="['click']"><button class="group-more icon-button" type="button" :aria-label="element.group_name"><img :src="icons.more" alt="" /></button><template #overlay><a-menu @click="({ key }) => emit('group-action', { key, record: element })"><a-menu-item key="rename">{{ t('btn_rename') }}</a-menu-item><a-menu-item key="delete">{{ t('btn_delete') }}</a-menu-item></a-menu></template></a-dropdown></div></template>
      </Draggable>
      <div class="group-width-handle" @mousedown="emit('resize-groups', $event)" />
    </aside>
    <main class="document-workspace-main">
      <div class="document-toolbar">
        <div class="document-toolbar-title">
          <a-popover v-if="collapsed" placement="bottomLeft" :trigger="['hover', 'focus']" overlay-class-name="workspace-document-group-popover">
            <button class="icon-button expand-groups" type="button" :aria-label="commonT('document_expand_groups')" @click="emit('toggle-groups')"><img :src="icons.collapse" alt="" /></button>
            <template #content><div class="collapsed-group-list"><button v-for="group in groups" :key="group.id" type="button" :class="{ active: String(group.id) === String(groupId) }" @click="chooseGroup(group)"><span>{{ group.group_name }}</span><span>{{ group.total || 0 }}</span></button></div></template>
          </a-popover>
          <h2 :title="currentGroup?.group_name">{{ currentGroup?.group_name || t('label_all_groups') }}</h2>
          <template v-if="collapsed"><span class="learning-summary" v-if="counts.learned_count">{{ metricLabel('status_learned') }} {{ counts.learned_count }}</span><span class="learning-summary" v-if="counts.learned_err_count">{{ metricLabel('status_learn_failed') }} {{ counts.learned_err_count }}</span></template>
        </div>
        <div class="document-toolbar-controls">
          <a-input class="document-search" :value="fileName" :placeholder="commonT('document_search')" :aria-label="commonT('document_search')" allow-clear @update:value="emit('update:fileName', $event)" @change="search"><template #prefix><img :src="icons.search" alt="" /></template></a-input>
          <a-select class="status-filter" :value="status" :options="statusOptions" :bordered="false" :aria-label="t('col_status')" @update:value="emit('update:status', $event)" @change="search"><template #suffixIcon><img :src="icons.chevron" alt="" /></template></a-select>
          <a-tooltip v-if="!readonly" :title="t('btn_metadata')"><button type="button" class="metadata-button icon-button" :aria-label="t('btn_metadata')" @click="emit('metadata')"><img :src="icons.metadata" alt="" /><span v-if="collapsed">{{ t('btn_metadata') }}</span></button></a-tooltip>
          <div class="document-view-toggle" role="group" :aria-label="commonT('document_view')"><button type="button" :class="{ active: viewMode === 'list' }" :aria-label="commonT('document_list')" :aria-pressed="viewMode === 'list'" @click="emit('view-change', 'list')"><img :src="viewMode === 'list' ? icons.listActive : icons.list" alt="" /></button><button type="button" :class="{ active: viewMode === 'folder' }" :aria-label="commonT('document_cards')" :aria-pressed="viewMode === 'folder'" @click="emit('view-change', 'folder')"><img :src="viewMode === 'folder' ? icons.cardsActive : icons.cards" alt="" /></button></div>
          <a-popover v-if="!readonly" v-model:open="addOpen" placement="bottomRight" :trigger="['click']" overlay-class-name="workspace-document-add-popover" :overlay-inner-style="{ padding: 0 }"><button type="button" class="add-content" :aria-expanded="addOpen"><img :src="icons.add" alt="" />{{ t('btn_add_content') }}</button><template #content><div class="document-add-grid"><button v-for="entry in addEntries" :key="entry.key" type="button" @click="addContent(entry.key)"><span class="add-entry-icon"><img :src="entry.icon" alt="" /></span><span><strong>{{ t(entry.title) }}</strong><span class="add-entry-description">{{ t(entry.description) }}</span></span></button></div></template></a-popover>
        </div>
      </div>
      <div v-if="viewMode === 'list'" ref="tableBox" class="workspace-document-table">
        <a-table :columns="tableColumns" :data-source="displayFiles" :loading="loading" :pagination="false" :custom-row="() => ({ tabindex: 0 })" :scroll="{ y: scrollY }" row-key="id" table-layout="fixed">
          <template #headerCell="{ column }"><div class="document-column-header"><button v-if="column.sortField" type="button" class="document-sort" :title="column.local ? commonT('document_sort_current_page') : column.title" @click="sortBy(column)"><span>{{ column.title }}</span><span class="sort-icons"><img :src="icons.sortUp" alt="" :class="{ chosen: activeSort.field === column.key && activeSort.order === 'ascend' }" /><img :src="icons.sortDown" alt="" :class="{ chosen: activeSort.field === column.key && activeSort.order === 'descend' }" /></span></button><span v-else>{{ column.title }}</span><ColumnWidthHandle v-if="column.key !== 'action'" :icon="icons.resize" :width="column.width" :min-width="column.minWidth" :label="commonT('resize_column', { name: column.title })" @resize="resizeColumn(column, $event)" @reset="resetColumn(column)" /></div></template>
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'file_name'"><button class="document-file-name" type="button" :title="record.file_name || record.doc_url" @click="action('preview', record)"><img v-if="fileIcon(record)" :src="fileIcon(record)" alt="" /><svg-icon v-else name="file_ext" /><span>{{ record.file_name || record.doc_url }}</span></button></template>
            <template v-else-if="column.key === 'file_size'"><span>{{ Number(record.doc_type) === 3 ? '--' : record.file_size_str }}</span><span class="document-secondary">；{{ record.paragraph_count || 0 }} {{ commonT('document_segments') }}</span></template>
            <template v-else-if="column.key === 'status'"><div class="document-status-cell"><DocumentStatus :record="record" /><button v-if="[3, 7, 8].includes(Number(record.status))" type="button" class="icon-button" :aria-label="t('tooltip_relearn')" @click="action('relearn', record)"><img :src="icons.refresh" alt="" /></button><button v-if="Number(record.status) === 0 && Number(record.pdf_parse_type) >= 2" type="button" class="cancel-ocr" @click="action('cancel-ocr', record)">{{ t('btn_cancel') }}</button></div></template>
            <template v-else-if="column.key === 'update_info'"><span>{{ dateOnly(record.doc_last_renew_time_desc && record.doc_last_renew_time_desc !== '--' ? record.doc_last_renew_time_desc : record.update_time) }}</span><span class="document-secondary">；{{ formatFrequency(record) }}</span></template>
            <template v-else-if="column.key === 'graph_status'"><span :title="record.graph_err_msg">{{ t(['status_wait_generate', 'status_queueing', 'status_generate_complete', 'status_generate_failed', 'status_generating', 'status_partial_success'][Number(record.graph_status)] || 'status_wait_generate') }}</span><a-button v-if="[0, 3].includes(Number(record.graph_status))" type="link" @click="action('generate-graph', record)">{{ t('btn_generate') }}</a-button><span v-if="Number(record.graph_status) === 5">{{ t('label_fail_count') }}{{ record.graph_err_count || 0 }}<a-button type="link" @click="action('graph-details', record)">{{ t('btn_details') }}</a-button></span></template>
            <template v-else-if="column.key === 'graph_entity_count'"><a @click="action('graph', record)">{{ record.graph_entity_count }}</a></template>
            <template v-else-if="column.key === 'action'"><div class="document-row-actions"><button type="button" class="icon-button" :aria-label="t('tooltip_re_segment')" @click="action('resegment', record)"><img :src="icons.refresh" alt="" /></button><DocumentActions :record="record" @action="emit('action', $event)" /></div></template>
            <template v-else><span class="document-cell-text" :title="String(record[column.dataIndex] ?? '')">{{ record[column.dataIndex] ?? '--' }}</span></template>
          </template>
        </a-table>
      </div>
      <div v-else class="workspace-document-cards"><a-spin :spinning="loading"><div v-if="cardItems.length" class="document-card-grid"><article v-for="item in cardItems" :key="`${item.kind}-${item.record.id}`" class="document-card" :class="{ 'group-card': item.kind === 'group' }">
        <button type="button" class="document-card-open" :title="item.record.group_name || item.record.file_name || item.record.doc_url" @click="item.kind === 'group' ? chooseGroup(item.record) : action('preview', item.record)"><span class="document-card-thumbnail" :class="{ placeholder: !cover(item.record) && item.kind !== 'group' }"><img v-if="item.kind === 'group'" class="folder-image" :src="icons.folder" alt="" /><img v-else class="file-cover" :src="cover(item.record) || defaultCover" alt="" @error="useDefaultCover" /><DocumentStatus v-if="item.kind !== 'group' && Number(item.record.status) !== 2" :record="item.record" card /></span><span class="document-card-name">{{ item.record.group_name || item.record.file_name || item.record.doc_url }}</span></button>
        <div class="document-card-footer"><span v-if="item.kind === 'group'">{{ commonT('document_group_total', { total: item.record.total || 0 }) }}</span><span v-else class="document-file-extension"><img v-if="fileIcon(item.record)" :src="fileIcon(item.record)" alt="" /><svg-icon v-else name="file_ext" />{{ item.record.file_ext || '--' }}</span><span>{{ cardDate(item.record.update_time || item.record.last_update) }}</span></div>
        <DocumentActions v-if="item.kind === 'file'" :record="item.record" card @action="emit('action', $event)" /><a-dropdown v-else :trigger="['click']"><button type="button" class="document-group-card-more icon-button" :aria-label="item.record.group_name"><img :src="icons.cardMore" alt="" /></button><template #overlay><a-menu @click="({ key }) => emit('group-action', { key, record: item.record })"><a-menu-item key="rename">{{ t('btn_rename') }}</a-menu-item><a-menu-item key="delete">{{ t('btn_delete') }}</a-menu-item></a-menu></template></a-dropdown>
      </article></div><a-empty v-else :description="t('empty_no_doc')" /></a-spin></div>
      <footer class="workspace-document-pagination"><span>{{ commonT('document_total', { total: pagination.total || 0 }) }}</span><a-pagination :current="pagination.page" :page-size="pagination.size" :total="pagination.total" :page-size-options="['10', '20', '50', '100']" show-size-changer @change="changePage" /></footer>
    </main>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import Draggable from 'vuedraggable'
import dayjs from 'dayjs'
import { useI18n } from '@/hooks/web/useI18n'
import icons, { fileIcons } from './workspace-document-icons'
import ColumnWidthHandle from '../qa-knowledge-document/components/column-width-handle.vue'
import DocumentActions from './document-actions.vue'
import DocumentStatus from './document-status.vue'
import { useTableScroll } from './use-table-scroll'
import defaultCover from '@/assets/img/default-bg.png'
const props = defineProps({
  groups: { type: Array, default: () => [] }, groupId: { type: [Number, String], default: -1 }, groupWidth: { type: Number, default: 200 }, collapsed: Boolean,
  files: { type: Array, default: () => [] }, counts: { type: Object, default: () => ({}) }, pagination: { type: Object, required: true }, extraColumns: { type: Array, default: () => [] },
  fileName: String, status: { type: String, default: '' }, statusOptions: { type: Array, default: () => [] },
  viewMode: { type: String, default: 'list' }, libraryId: { type: [Number, String], required: true }, loading: Boolean, readonly: Boolean, updateFrequency: { type: Function, required: true }
})
const emit = defineEmits(['toggle-groups', 'resize-groups', 'reorder-groups', 'group-action', 'select-group', 'update:fileName', 'update:status', 'search', 'view-change', 'metadata', 'add', 'action', 'sort', 'page-change'])
const { t } = useI18n('views.library.library-details.knowledge-document')
const { t: commonT } = useI18n('views.customer-service.index')
const groupSearch = ref('')
const addOpen = ref(false)
const tableBox = ref(null)
const hostWidth = ref(0)
const { scrollY } = useTableScroll(tableBox, () => props.viewMode === 'list')
const currentGroup = computed(() => props.groups.find(group => String(group.id) === String(props.groupId)))
const filteredGroups = computed(() => props.groups.filter(group => String(group.group_name).toLowerCase().includes(groupSearch.value.trim().toLowerCase())))
const addEntries = computed(() => [
  { key: 1, icon: icons.upload, title: 'menu_local_doc', description: 'menu_local_doc_desc' },
  { key: 2, icon: icons.online, title: 'menu_online_data', description: 'menu_online_data_desc' },
  { key: 3, icon: icons.custom, title: 'menu_custom_doc', description: 'menu_custom_doc_desc' },
  { key: 4, icon: icons.feishu, title: 'menu_feishu', description: 'menu_feishu_desc' }
])
const cardItems = computed(() => [
  ...(Number(props.groupId) === -1 && !props.fileName && !props.status ? props.groups.filter(group => Number(group.id) > 0).map(record => ({ kind: 'group', record })) : []),
  ...displayFiles.value.map(record => ({ kind: 'file', record }))
])
const widthKey = computed(() => `customer_service_document_columns_${props.libraryId}`)
const widths = ref({})
watch(widthKey, key => { try { const value = JSON.parse(localStorage.getItem(key) || '{}'); widths.value = value && typeof value === 'object' ? value : {} } catch { widths.value = {} } }, { immediate: true })
const columnDefinitions = computed(() => [
  { key: 'file_name', dataIndex: 'file_name', title: t('col_doc_name'), width: 148, minWidth: 120 },
  { key: 'file_size', dataIndex: 'file_size', title: t('col_size_segment'), width: 160, minWidth: 140, sortField: 'file_size', local: true },
  ...['total_hits', 'yesterday_hits', 'today_hits'].map((key, i) => ({ key, dataIndex: key, title: t(['col_total', 'col_yesterday', 'col_today'][i]), width: 72, minWidth: 64, sortField: key })),
  { key: 'status', title: t('col_status'), width: 136, minWidth: 120, sortField: 'status', local: true },
  { key: 'update_info', title: t('col_update_info'), width: 160, minWidth: 150, sortField: 'update_time', local: true },
  ...props.extraColumns.map(column => ({ ...column, width: 160, minWidth: 100 })),
  { key: 'action', title: t('col_action'), width: 96 }
])
const tableColumns = computed(() => {
  const columns = columnDefinitions.value.map(column => ({ ...column, ellipsis: true, width: Math.max(column.minWidth || 0, Number(widths.value[column.key]) || column.width) }))
  const total = columns.reduce((sum, column) => sum + column.width, 0)
  const available = hostWidth.value || total
  let flexible = columns.filter(column => !widths.value[column.key] && column.key !== 'action')
  let fixed = total - flexible.reduce((sum, column) => sum + column.width, 0)
  // 已保存的过大列宽和新增元数据也必须适配表体可用宽度，给操作按钮保留空间。
  if (!flexible.length || fixed >= available) {
    flexible = columns.filter(column => column.key !== 'action')
    fixed = columns.at(-1).width
  }
  if (fixed >= available) { flexible = columns; fixed = 0 }
  if (total > available) {
    const scale = (available - fixed) / flexible.reduce((sum, column) => sum + column.width, 0)
    flexible.forEach(column => { column.width *= scale })
  }
  columns.forEach(column => { column.width = Math.floor(column.width); column.minWidth = Math.min(column.minWidth || column.width, column.width) })
  const remaining = available - columns.reduce((sum, column) => sum + column.width, 0)
  const expandedColumn = flexible[0] || columns[0]
  expandedColumn.width += remaining
  return columns
})
const localSort = ref({ field: '', order: '' })
const activeSort = computed(() => localSort.value.field ? localSort.value : { field: props.pagination.sort_field, order: props.pagination.sort_type === 'ASC' ? 'ascend' : props.pagination.sort_type === 'DESC' ? 'descend' : '' })
const displayFiles = computed(() => {
  const { field, order } = localSort.value
  if (!field || !order) return props.files
  return [...props.files].sort((a, b) => {
    const key = field === 'update_info' ? 'update_time' : field
    const result = key === 'update_time' ? String(a[key]).localeCompare(String(b[key])) : Number(a[key] || 0) - Number(b[key] || 0)
    return order === 'ascend' ? result : -result
  })
})
function resizeColumn(column, width) { widths.value = { ...widths.value, [column.key]: Math.max(column.minWidth || 64, width) }; localStorage.setItem(widthKey.value, JSON.stringify(widths.value)) }
function resetColumn(column) { const next = { ...widths.value }; delete next[column.key]; widths.value = next; localStorage.setItem(widthKey.value, JSON.stringify(next)) }
function sortBy(column) {
  const previous = activeSort.value
  const order = previous.field !== column.key || !previous.order ? 'ascend' : previous.order === 'ascend' ? 'descend' : ''
  if (column.local) localSort.value = { field: column.key, order }
  else { localSort.value = { field: '', order: '' }; emit('sort', { field: column.sortField, order }) }
}
function chooseGroup(group) { localSort.value = { field: '', order: '' }; emit('select-group', group) }
function checkGroupMove(event) { return Number(event.draggedContext.element.id) > 0 && Number(event.relatedContext?.element?.id) > 0 }
function search() { localSort.value = { field: '', order: '' }; emit('search') }
function addContent(key) { addOpen.value = false; emit('add', { key }) }
function action(key, record) { emit('action', { key, record }) }
function changePage(page, size) { localSort.value = { field: '', order: '' }; emit('page-change', page, size) }
function fileIcon(record) { return fileIcons[String(record.file_ext || '').toLowerCase().replace(/^\./, '')] }
function cover(record) { return record.thumb_path || record.cover_url || record.cover || record.article_cover_url || '' }
function useDefaultCover(event) { if (!event.target.src.endsWith(defaultCover)) event.target.src = defaultCover }
function dateOnly(value) { const date = typeof value === 'number' || /^\d{9,13}$/.test(String(value)) ? dayjs(Number(value) < 1e12 ? Number(value) * 1000 : Number(value)) : dayjs(value); return value && date.isValid() ? date.format('YYYY-MM-DD') : '--' }
function cardDate(value) { return dateOnly(value).replaceAll('-', '/') }
function formatFrequency(record) { const value = props.updateFrequency(record); return !value || value === '--' ? commonT('document_manual_update') : value }
function metricLabel(key) { return t(key).replace(/^[，,]\s*/, '').replace(/[：:]\s*$/, '') }
let stopWidthObserver
watch(tableBox, element => {
  stopWidthObserver?.()
  if (!element) return
  const body = element.querySelector('.ant-table-body')
  const measure = () => { hostWidth.value = body?.clientWidth || element.clientWidth }
  measure()
  if (typeof ResizeObserver === 'undefined') { window.addEventListener('resize', measure); stopWidthObserver = () => window.removeEventListener('resize', measure) }
  else { const observer = new ResizeObserver(measure); observer.observe(element); if (body) observer.observe(body); stopWidthObserver = () => observer.disconnect() }
}, { flush: 'post' })
onBeforeUnmount(() => stopWidthObserver?.())
</script>

<style lang="less" scoped>
.workspace-document-view { display: flex; flex: 1; width: 100%; height: 100%; min-width: 0; min-height: 0; color: #222; font-size: 14px; line-height: 22px; }
.icon-button { display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0; width: 24px; height: 24px; padding: 0; border: 0; border-radius: 6px; background: transparent; cursor: pointer; }
.icon-button:hover { background: #f5f5f5; }
button:focus-visible { outline: 2px solid #146ef5; outline-offset: -2px; }
h2 { min-width: 0; margin: 0; overflow: hidden; color: #080808; font-size: 16px; font-weight: 600; line-height: 24px; text-overflow: ellipsis; white-space: nowrap; }
.document-groups { position: relative; display: flex; flex-direction: column; flex-shrink: 0; min-height: 0; padding: 24px 16px 24px 0; border-right: 1px solid #e8e8e8; }
.groups-heading { display: flex; align-items: center; gap: 8px; margin-bottom: 12px; h2 { flex: 1; } }
.group-search { flex-shrink: 0; height: 32px; margin-bottom: 6px; }
.document-group-list { flex: 1; min-height: 0; overflow-y: auto; }
.document-group-row { position: relative; height: 32px; margin-top: 2px; border-radius: 4px; &:first-child { margin-top: 0; } &:hover { background: #f5f5f5; } &.active { color: #146ef5; background: #e8f0fe; } }
.group-select { display: flex; align-items: center; justify-content: space-between; gap: 8px; width: 100%; height: 100%; padding: 0 12px; border: 0; border-radius: inherit; color: inherit; background: transparent; cursor: pointer; span:first-child { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; } }
.group-count { color: #ababab; }
.group-more { position: absolute; top: 4px; right: 4px; display: none; background: #f5f5f5; }
.group-drag { position: absolute; top: 5px; left: 0; display: none; width: 16px; cursor: grab; }
.document-group-row:hover, .document-group-row:focus-within { .group-more, .group-drag { display: inline-flex; } &:has(.group-more) .group-count { opacity: 0; } }
.group-width-handle { position: absolute; top: 0; right: -3px; z-index: 1; width: 6px; height: 100%; cursor: col-resize; }
.document-workspace-main { display: flex; flex-direction: column; flex: 1; min-width: 0; min-height: 0; padding: 24px 0 24px 24px; }
.groups-collapsed .document-workspace-main { padding-left: 0; }
.document-toolbar, .document-toolbar-title, .document-toolbar-controls { display: flex; align-items: center; gap: 8px; }
.document-toolbar { flex-shrink: 0; gap: 16px; min-height: 32px; margin-bottom: 16px; }
.document-toolbar-title { flex: 1; min-width: 0; }
.document-toolbar-controls { justify-content: flex-end; flex-shrink: 0; }
.learning-summary { flex-shrink: 0; color: #5a5a5a; font-size: 12px; }
.expand-groups img { transform: rotate(180deg); }
.document-search { width: 200px; height: 32px; }
.status-filter { width: 104px; }
.metadata-button { gap: 6px; width: auto; min-width: 32px; height: 32px; padding: 0 8px; color: #5a5a5a; white-space: nowrap; }
.document-view-toggle { display: flex; gap: 2px; padding: 2px; border-radius: 6px; background: #f5f5f5; button { display: flex; align-items: center; justify-content: center; width: 28px; height: 28px; padding: 0; border: 0; border-radius: 6px; background: transparent; cursor: pointer; &.active { background: #fff; box-shadow: 0 2px 4px #0808080a; } } }
.add-content { display: inline-flex; align-items: center; justify-content: center; gap: 6px; height: 32px; padding: 0 12px; border: 1px solid #146ef5; border-radius: 6px; color: #fff; background: #146ef5; white-space: nowrap; cursor: pointer; }
.workspace-document-table { flex: 1; min-width: 0; min-height: 0; overflow: hidden; }
.workspace-document-table :deep(.ant-table-thead > tr > th) { height: 48px; padding: 0 8px; border-bottom: 1px solid #d8d8d8; color: #5a5a5a; background: #fff; font-size: 14px; font-weight: 500; &::before { display: none !important; } }
.workspace-document-table :deep(.ant-table-tbody > tr > td) { padding: 12px 8px; border-bottom: 1px solid #d8d8d8; color: #5a5a5a; background: #fff; font-size: 14px; line-height: 22px; }
.workspace-document-table :deep(.ant-table-tbody > tr:hover > td), .workspace-document-table :deep(.ant-table-tbody > tr:focus-within > td) { background: #f5f5f5; border-bottom-color: transparent; }
.workspace-document-table :deep(.ant-table-tbody > tr.ant-table-row-selected > td) { background: #f1f6fe; border-bottom-color: transparent; }
.workspace-document-table :deep(.ant-table-tbody > tr:hover > td:first-child), .workspace-document-table :deep(.ant-table-tbody > tr.ant-table-row-selected > td:first-child) { border-radius: 8px 0 0 8px; }
.workspace-document-table :deep(.ant-table-tbody > tr:hover > td:last-child), .workspace-document-table :deep(.ant-table-tbody > tr.ant-table-row-selected > td:last-child) { border-radius: 0 8px 8px 0; }
.document-column-header { position: relative; display: flex; align-items: center; height: 48px; white-space: nowrap; }
.document-sort { display: inline-flex; align-items: center; gap: 6px; max-width: 100%; padding: 0; border: 0; color: inherit; background: transparent; cursor: pointer; }
.document-sort > span:first-child { overflow: hidden; text-overflow: ellipsis; }
.sort-icons { display: inline-flex; flex-direction: column; flex-shrink: 0; img { opacity: .5; &:first-child { transform: rotate(180deg); } &.chosen { opacity: 1; filter: brightness(0) saturate(100%) invert(35%) sepia(97%) saturate(2775%) hue-rotate(206deg) brightness(100%) contrast(94%); } } }
.document-file-name { display: flex; align-items: center; gap: 8px; width: 100%; padding: 0; border: 0; color: #222; background: transparent; text-align: left; cursor: pointer; img, .svg-action { flex-shrink: 0; width: 24px; height: 24px; } span { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; } }
.document-secondary { color: #ababab; font-size: 12px; }
.document-cell-text { display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.document-status-cell, .document-row-actions { display: flex; align-items: center; gap: 4px; white-space: nowrap; }
.document-row-actions { justify-content: center; gap: 8px; .icon-button { width: 32px; height: 32px; } }
.cancel-ocr { border: 0; color: #146ef5; background: transparent; cursor: pointer; }
.workspace-document-cards { flex: 1; min-height: 0; overflow-y: auto; container-type: inline-size; :deep(.ant-spin-nested-loading) { min-height: 100%; } }
.document-card-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); align-content: start; gap: 12px; padding: 2px; }
.document-card { position: relative; display: flex; flex-direction: column; gap: 8px; height: 192px; min-width: 0; padding: 12px; border: 1px solid #e8e8e8; border-radius: 12px; background: #fff; }
.document-card:hover, .document-card:focus-within { border-color: transparent; background: #f5f5f5; }
.document-card:hover :deep(.card-more), .document-card:focus-within :deep(.card-more), .document-card:hover .document-group-card-more, .document-card:focus-within .document-group-card-more { opacity: 1; }
.document-card-open { display: flex; flex-direction: column; gap: 8px; width: 100%; padding: 0; border: 0; border-radius: 8px; color: #222; background: transparent; cursor: pointer; }
.document-card-thumbnail { position: relative; display: flex; align-items: center; justify-content: center; width: 100%; height: 86px; overflow: hidden; border-radius: 8px; }
.file-cover { width: 100%; height: 100%; object-fit: cover; }
.folder-image { width: 120px; height: auto; }
.document-card-name { display: -webkit-box; width: 100%; height: 44px; overflow: hidden; -webkit-line-clamp: 2; -webkit-box-orient: vertical; line-height: 22px; text-align: center; word-break: break-word; }
.document-card-footer { display: flex; align-items: center; justify-content: space-between; gap: 8px; color: #ababab; font-size: 12px; line-height: 20px; white-space: nowrap; }
.document-file-extension { display: inline-flex; align-items: center; gap: 6px; text-transform: lowercase; img, .svg-action { width: 16px; height: 16px; } }
.document-group-card-more { position: absolute; top: 12px; right: 12px; opacity: 0; }
.workspace-document-pagination { display: flex; align-items: center; justify-content: space-between; flex-shrink: 0; gap: 16px; min-height: 32px; padding-top: 16px; color: #5a5a5a; font-size: 12px; }
.workspace-document-view :deep(.ant-input-affix-wrapper), .workspace-document-view :deep(.ant-select-selector) { border-color: #d8d8d8; border-radius: 6px; }
.workspace-document-view :deep(.ant-input-prefix) { margin-right: 8px; }
.workspace-document-view :deep(.ant-table-body) { overflow-x: hidden !important; scrollbar-width: thin; @supports selector(::-webkit-scrollbar) { scrollbar-width: auto; } &::-webkit-scrollbar { width: 6px; height: 6px; } &::-webkit-scrollbar-thumb { border-radius: 6px; background: #e5e7eb; } }
@container (min-width: 900px) { .document-card-grid { grid-template-columns: repeat(4, minmax(0, 1fr)); } }
@container (min-width: 1100px) { .document-card-grid { grid-template-columns: repeat(6, minmax(0, 1fr)); } }
@media (max-width: 1100px) { .document-toolbar { flex-wrap: wrap; } .document-toolbar-controls { flex-wrap: wrap; } .document-toolbar-title { min-width: 160px; } .learning-summary { display: none; } }
@media (hover: none) { .document-card :deep(.card-more), .document-group-card-more { opacity: 1; } }
</style>

<style lang="less">
.workspace-document-add-popover .ant-popover-inner { padding: 0; border-radius: 12px; overflow: hidden; }
.document-add-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 8px; width: 520px; max-width: calc(100vw - 32px); padding: 12px; button { display: grid; grid-template-columns: 32px minmax(0, 1fr); align-items: start; gap: 12px; min-height: 88px; padding: 12px; border: 0; border-radius: 6px; color: #080808; background: #fff; text-align: left; cursor: pointer; &:hover, &:focus-visible { outline: 0; background: #f5f5f5; } } strong { font-size: 14px; font-weight: 500; line-height: 22px; } }
.add-entry-icon { display: flex; align-items: center; justify-content: center; width: 32px; height: 32px; border-radius: 6px; background: #f5f5f5; }
.add-entry-description { display: block; margin-top: 2px; color: #5a5a5a; font-size: 12px; line-height: 20px; }
.workspace-document-group-popover .ant-popover-inner { width: 200px; padding: 8px; }
.collapsed-group-list { display: flex; flex-direction: column; gap: 2px; max-height: 360px; overflow-y: auto; button { display: flex; align-items: center; justify-content: space-between; gap: 8px; min-height: 32px; padding: 0 12px; border: 0; border-radius: 4px; color: #5a5a5a; background: #fff; cursor: pointer; &.active, &:hover, &:focus-visible { color: #146ef5; background: #e8f0fe; outline: 0; } } }
@media (max-width: 560px) { .document-add-grid { grid-template-columns: minmax(0, 1fr); } }
</style>
