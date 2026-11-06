<template>
  <section class="knowledge-page">
    <CustomerPageHeader v-model:active-key="tab" :title="t('knowledge_base')" :tabs="headerTabs" @change="selectLibrary" />
      <!-- <a-button v-if="tab !== 'other'" class="library-settings-button" :disabled="loading || !settingsTargetId || uploading" :title="tab === 'normal' && !normalId ? t('library_settings_empty_tip') : t('library_settings')" @click="openSettings"><SettingOutlined />{{ t('library_settings') }}</a-button> -->
    <div v-if="error" class="feedback"><a-alert type="error" :message="error" show-icon /><a-button @click="retry">{{ t('retry') }}</a-button></div>
    <div v-if="loading" class="feedback"><a-spin /></div>
    <div v-else class="knowledge-content" :class="{ 'qa-document-content': tab === 'qa', 'normal-document-content': tab === 'normal', 'document-content': tab !== 'other', 'fixed-document-content': detailId && (tab === 'qa' || (tab === 'normal' && normalId && !files.length && !uploading)), 'other-library-content': tab === 'other' }">
      <template v-if="tab === 'other'">
        <div class="library-filters" role="group" :aria-label="t('other_libraries')"><button v-for="filter in libraryFilters" :key="filter.key" type="button" :class="{ active: libraryFilter === filter.key }" :aria-pressed="libraryFilter === filter.key" @click="libraryFilter = filter.key">{{ filter.label }}</button></div>
        <div class="other-libraries-list">
        <div v-if="filteredOtherLibraries.length" class="library-grid">
          <article v-for="item in filteredOtherLibraries" :key="item.id" class="library-card" :class="{ linked: isAssociated(item.id), 'actions-open': openLibraryMenuId === String(item.id), 'is-busy': relatingId === String(item.id) }">
            <button type="button" class="library-info" :aria-label="item.library_name" @click="openLibrary(item)">
              <a-avatar class="library-avatar" :class="{ faq: Number(item.type) === 2 }" :size="40" shape="square" :src="item.avatar || libraryTypeIcons[Number(item.type)] || otherLibraryIcon" />
              <span class="library-copy"><span class="library-heading"><span class="library-name" :title="item.library_name">{{ item.library_name }}</span><span class="library-badge" :class="{ 'linked-badge': isAssociated(item.id) }">{{ t(isAssociated(item.id) ? 'library_linked' : 'library_unlinked') }}</span></span>
              <span class="library-metrics">
                <a-tooltip v-for="metric in libraryMetrics(item)" :key="metric.key" :arrow="false" placement="top" overlay-class-name="other-library-tooltip">
                  <template #title>{{ metric.title }}<img class="library-tooltip-beak" :src="tooltipBeak" alt="" /></template>
                  <span :aria-label="metric.title"><img :src="metric.icon" alt="" />{{ metric.value }}</span>
                </a-tooltip>
              </span>
              </span>
            </button>
              <div class="library-actions">
                <a-tooltip :title="t(isAssociated(item.id) ? 'remove_relation' : 'link_library')">
                  <a-button size="small" class="relation-button" :class="{ 'is-linked': isAssociated(item.id) }" :type="isAssociated(item.id) ? 'default' : 'primary'" :loading="relatingId === String(item.id)" :disabled="!!relatingId" :aria-label="t(isAssociated(item.id) ? 'remove_relation' : 'link_library')" :aria-pressed="isAssociated(item.id)" @click="toggleLibrary(item)">
                    <template #icon><img :src="isAssociated(item.id) ? unlinkLibraryIcon : linkLibraryIcon" alt="" /></template>
                  </a-button>
                </a-tooltip>
                <a-dropdown :trigger="['click']" :open="openLibraryMenuId === String(item.id)" @openChange="openLibraryMenuId = $event ? String(item.id) : ''">
                  <a-button type="text" class="library-more" :aria-label="t('library_more')" :aria-expanded="openLibraryMenuId === String(item.id)"><img :src="libraryMoreIcon" alt="" /></a-button>
                  <template #overlay><a-menu @click="openLibrary(item)"><a-menu-item key="view">{{ t('view') }}</a-menu-item></a-menu></template>
                </a-dropdown>
              </div>
          </article>
        </div>
        <a-empty v-else :description="t('empty_libraries')" />
        </div>
        <footer class="other-libraries-total">{{ t('other_libraries_total', { total: filteredOtherLibraries.length }) }}</footer>
      </template>
      <template v-else-if="tab === 'normal' && (!normalId || files.length || uploading)">
        <div class="upload-empty"><UploadInput v-model:value="files" :max-count="20" :type="0" />
          <a-button type="primary" :loading="uploading" :disabled="!files.length" @click="upload">{{ t('upload_documents') }}</a-button>
          <a-button v-if="files.length" :disabled="uploading" @click="files = []">{{ t('cancel') }}</a-button>
          <p>{{ t('normal_library_lazy_tip') }}</p>
          <a v-if="uploadError" @click="openModelSettings">{{ t('model_management') }}</a>
        </div>
      </template>
      <QaDocument v-else-if="tab === 'qa' && detailId" :key="detailId" :library_id="detailId" compact />
      <KnowledgeDocument v-else-if="tab === 'normal' && detailId" :key="detailId" :library_id="detailId" compact />
      <a-empty v-else :description="t('empty_libraries')" />
    </div>
    <a-modal :open="settingsOpen" :title="t(settingsTab === 'qa' ? 'qa_library_settings' : 'normal_library_settings')" :width="960" :footer="null" :closable="!settingsClosing" :keyboard="!settingsClosing" :mask-closable="false" destroy-on-close @cancel="closeSettings">
      <div class="library-settings-body"><KnowledgeConfig v-if="settingsOpen" :key="settingsLibraryId" ref="settingsForm" :library-id="settingsLibraryId" embedded @view-documents="closeSettings" /></div>
    </a-modal>
  </section>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, ref, watch } from 'vue'
import { SettingOutlined } from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useLibraryStore } from '@/stores/modules/library'
import { formatFileSize } from '@/utils/index'
import { getLibraryList, addLibraryFile } from '@/api/library'
import { ensureRobotNormalLibrary, relationLibrary } from '@/api/robot'
import QaDocument from '@/views/library/library-details/qa-knowledge-document/index.vue'
import KnowledgeDocument from '@/views/library/library-details/knowledge-document.vue'
import UploadInput from '@/views/library/add-library/components/upload-input.vue'
import KnowledgeConfig from '@/views/library/library-details/knowledge-config.vue'
import CustomerPageHeader from '../components/page-header.vue'
import qaLibraryIcon from '@/assets/img/customer-service/normal-knowledge/qa-library.png'
import normalLibraryIcon from '@/assets/img/customer-service/normal-knowledge/normal-library.png'
import otherLibraryIcon from '@/assets/img/customer-service/normal-knowledge/other-library.png'
import documentCountIcon from '@/assets/img/customer-service/normal-knowledge/document-count.svg'
import storageIcon from '@/assets/img/customer-service/normal-knowledge/storage.svg'
import relationsIcon from '@/assets/img/customer-service/normal-knowledge/relations.svg'
import libraryMoreIcon from '@/assets/img/customer-service/normal-knowledge/more.svg'
import linkLibraryIcon from '@/assets/img/customer-service/other-knowledge/link.svg'
import unlinkLibraryIcon from '@/assets/img/customer-service/other-knowledge/unlink.svg'
import tooltipBeak from '@/assets/img/customer-service/other-knowledge/tooltip-beak.svg'
const { t } = useI18n('views.customer-service.index')
const router = useRouter()
const context = inject('customerRobotContext')
const robotId = context.robotInfo.id
const libraryStore = useLibraryStore()
const tab = ref('qa')
const headerTabs = computed(() => [
  { key: 'qa', label: t('qa_library') },
  { key: 'normal', label: t('normal_library') },
  { key: 'other', label: t('other_libraries') }
])
const libraries = ref([])
const loading = ref(false)
const error = ref('')
const uploadError = ref(false)
const uploading = ref(false)
const files = ref([])
const detailId = ref('')
const relatingId = ref('')
const settingsOpen = ref(false)
const settingsClosing = ref(false)
const settingsLibraryId = ref('')
const settingsTab = ref('qa')
const settingsForm = ref(null)
const associatedIds = computed(() => String(context.robotInfo.library_ids || '').split(',').filter(Boolean))
const associatedLibraries = computed(() => libraries.value.filter(item => associatedIds.value.includes(String(item.id))))
const normalId = computed(() => Number(context.robotInfo.normal_library_id) || 0)
const settingsTargetId = computed(() => tab.value === 'normal' ? (normalId.value ? String(normalId.value) : '') : detailId.value)
let sequence = 0
let loadSequence = 0
let alive = true
const isProtected = id => [String(context.robotInfo.default_library_id), String(normalId.value)].includes(String(id))
const otherLibraries = computed(() => libraries.value.filter(item => !isProtected(item.id)))
const isAssociated = id => associatedIds.value.includes(String(id))
const libraryTypeIcons = { 0: normalLibraryIcon, 2: qaLibraryIcon }
const libraryFilter = ref('all')
const openLibraryMenuId = ref('')
watch([tab, libraryFilter], () => { openLibraryMenuId.value = '' })
const libraryFilters = computed(() => [
  { key: 'all', label: t('all') }, { key: 'normal', label: t('library_filter_normal') },
  { key: 'qa', label: t('library_filter_qa') }, { key: 'other', label: t('category_other') },
  { key: 'linked', label: t('library_linked') }
])
const filteredOtherLibraries = computed(() => otherLibraries.value.filter(item => {
  if (libraryFilter.value === 'linked') return isAssociated(item.id)
  if (libraryFilter.value === 'normal') return Number(item.type) === 0
  if (libraryFilter.value === 'qa') return Number(item.type) === 2
  if (libraryFilter.value === 'other') return ![0, 2].includes(Number(item.type))
  return true
}))
function libraryMetrics(item) {
  const size = formatFileSize(item.file_size)
  return [
    { key: 'documents', icon: documentCountIcon, value: item.file_total || 0, title: t('library_document_count_tip', { count: item.file_total || 0 }) },
    { key: 'storage', icon: storageIcon, value: size, title: t('library_storage_size_tip', { size: size.replace(/(\d)([A-Z]+)$/, '$1 $2') }) },
    { key: 'relations', icon: relationsIcon, value: item.robot_nums || 0, title: t('library_related_apps_tip', { count: item.robot_nums || 0 }) }
  ]
}
function openSettings() {
  if (tab.value === 'other' || loading.value || !settingsTargetId.value || uploading.value) return
  settingsLibraryId.value = settingsTargetId.value
  settingsTab.value = tab.value
  settingsOpen.value = true
}
async function closeSettings() {
  if (settingsClosing.value) return
  settingsClosing.value = true
  try { await settingsForm.value?.waitForSave() }
  catch (err) { if (alive) message.error(err.message || t('save_failed')) }
  finally {
    if (alive) {
      settingsOpen.value = false
      settingsClosing.value = false
      await load(false)
    }
  }
}
async function selectLibrary() {
  const version = ++sequence
  detailId.value = ''
  const item = tab.value === 'normal' ? libraries.value.find(item => String(item.id) === String(normalId.value) && Number(item.type) === 0) :
    associatedLibraries.value.find(item => Number(item.type) === 2 && String(item.id) === String(context.robotInfo.default_library_id)) || associatedLibraries.value.find(item => Number(item.type) === 2)
  if (tab.value === 'other' || !item) return
  try {
    await libraryStore.getLibraryInfo(item.id, () => alive && version === sequence)
    if (alive && version === sequence) detailId.value = String(item.id)
  } catch (err) { if (alive && version === sequence) error.value = err.message || t('load_failed') }
}
async function load(showLoading = true) {
  const version = ++loadSequence
  if (showLoading) loading.value = true
  error.value = ''
  try {
    const response = await getLibraryList({ type: '', show_open_docs: 1 })
    if (!alive || version !== loadSequence) return
    libraries.value = response.data || []
    await selectLibrary()
  } catch (err) { if (alive && version === loadSequence) error.value = err.message || t('load_failed') }
  finally { if (alive && version === loadSequence) loading.value = false }
}
async function retry() {
  try { await context.refreshRobot(); if (alive) await load() }
  catch (err) { if (alive) error.value = err.message || t('load_failed') }
}
async function upload() {
  if (uploading.value || !files.value.length) return
  if (files.value.some(file => !file.size)) { message.error(t('empty_file')); return }
  uploading.value = true
  uploadError.value = false
  const id = context.robotInfo.id
  const selectedFiles = [...files.value]
  try {
    await context.waitForSave()
    const ensured = await ensureRobotNormalLibrary(id)
    if (!alive) return
    const libraryId = ensured.data.library_id
    const data = new FormData()
    data.append('library_id', libraryId)
    data.append('group_id', '0')
    selectedFiles.forEach(file => data.append('library_files', file))
    await addLibraryFile(data)
    if (!alive) return
    files.value = []
    await context.refreshRobot()
    await load()
    message.success(t('operation_success'))
  } catch (err) {
    if (alive) {
      uploadError.value = true
      error.value = err.message || t('upload_failed')
      // 创建成功而上传失败时保留服务端引用，重试只会取回同一个库。
      await context.refreshRobot().catch(() => {})
    }
  } finally { if (alive) uploading.value = false }
}
async function setRelation(item, linked) {
  if (!alive || relatingId.value || isProtected(item.id)) return
  relatingId.value = String(item.id)
  try {
    await context.waitForSave()
    if (!alive) return
    // 关联接口提交完整 ID 列表，写入前读取最新状态，避免覆盖其他页面刚保存的关联。
    await context.refreshRobot()
    if (!alive || isProtected(item.id) || isAssociated(item.id) === linked) return
    const ids = linked ? [...associatedIds.value, String(item.id)] : associatedIds.value.filter(id => id !== String(item.id))
    const preserved = [context.robotInfo.default_library_id, normalId.value].filter(id => Number(id) > 0).map(String)
    await relationLibrary({ id: robotId, library_ids: [...new Set([...ids, ...preserved])].join(',') })
    if (!alive) return
    await context.refreshRobot()
    if (!alive) return
    await load(false)
    if (alive) message.success(t('operation_success'))
  } catch (err) { if (alive) message.error(err.message || t('save_failed')) }
  finally { if (alive) relatingId.value = '' }
}
function toggleLibrary(item) {
  if (relatingId.value) return
  if (isAssociated(item.id)) {
    Modal.confirm({ title: t('confirm_remove_relation'), onOk: () => setRelation(item, false) })
  } else {
    void setRelation(item, true)
  }
}
function openLibrary(item) {
  openLibraryMenuId.value = ''
  const target = Number(item.type) === 1 ? { path: '/public-library/config', query: { library_id: item.id } } : { path: '/library/details', query: { id: item.id } }
  window.open(router.resolve(target).href, '_blank', 'noopener,noreferrer')
}
function openModelSettings() { window.open(router.resolve('/user/model').href, '_blank', 'noopener,noreferrer') }
void load()
onBeforeUnmount(() => { alive = false; sequence++ })
</script>
<style lang="less" scoped>
.knowledge-page { flex: 1; display: flex; flex-direction: column; min-height: 0; margin: 8px 8px 8px 0; border-radius: 12px; overflow: hidden; background: #fff; }
.library-settings-button { margin-left: auto; border-radius: 8px; }
.library-settings-body { height: min(70vh, 720px); min-height: 240px; overflow: hidden; }
.knowledge-content { flex: 1; min-height: 0; overflow: auto; padding: 16px 24px; }
.knowledge-content.fixed-document-content {
  display: flex; flex-direction: column; overflow: hidden;
  :deep(.workspace-document) { flex: 1; min-height: 0; }
}
.knowledge-content.document-content {
  padding: 12px 16px;
  :deep(.workspace-document) { position: relative; }
  :deep(.group-content-box:not(.collapsed)) { padding-top: 12px; padding-left: 12px; }
  :deep(.group-content-box) {
    flex-shrink: 0;
    .group-head-box { padding-right: 12px; }
    .search-box { margin-top: 12px; margin-bottom: 8px; }
    .classify-box { min-height: 0; }
    .classify-item { height: 36px; margin-top: 2px; padding: 0 12px; border-radius: 4px; }
  }
  :deep(.main-content-box) { min-width: 0; padding: 12px 0 0 24px; }
  :deep(.group-header-title) { gap: 8px; margin-bottom: 16px; }
  :deep(.head-content .title), :deep(.head-content .btn-right-box) { flex-wrap: wrap; gap: 8px; }
  :deep(.document-main-layout) { display: flex; flex-direction: column; height: 100%; min-height: 0; overflow: hidden; }
  :deep(.document-main-layout > .head-content),
  :deep(.document-main-layout > .group-header-title),
  :deep(.document-main-layout > .list-tools) { flex-shrink: 0; }
  :deep(.document-main-layout > .list-tools) {
    flex-wrap: wrap; gap: 8px;
    .tools-items { flex-wrap: wrap; gap: 8px; }
    .tool-item { margin-right: 0; }
  }
  :deep(.document-main-layout > .content-block) { display: flex; flex-direction: column; flex: 1; min-height: 0; overflow: hidden; }
  :deep(.document-main-layout .fixed-table-host) { flex: 1; min-height: 0; overflow: hidden; }
  :deep(.document-main-layout .document-pagination) { display: flex; flex-shrink: 0; align-items: center; justify-content: flex-end; margin-top: 0; padding-top: 12px; }
  :deep(.document-main-layout > .document-pagination) { margin-top: auto; }
  :deep(.document-main-layout > .folder-view) { flex: 1; min-height: 0; overflow-y: auto; }
  :deep(.document-main-layout .ant-table-body),
  :deep(.document-main-layout > .folder-view) {
    scrollbar-width: thin;
    scrollbar-color: #e5e7eb transparent;
    @supports selector(::-webkit-scrollbar) {
      // 表格按 WebKit 滚动条尺寸预留表头占位，避免 thin 覆盖实际宽度造成列错位。
      scrollbar-width: auto;
      scrollbar-color: auto;
    }
    &::-webkit-scrollbar { width: 6px; height: 6px; }
    &::-webkit-scrollbar-track { background: transparent; }
    &::-webkit-scrollbar-thumb { background: #e5e7eb; border-radius: 6px; }
    &::-webkit-scrollbar-thumb:hover { background: #d1d5db; }
    &::-webkit-scrollbar-corner { background: transparent; }
  }
}
.knowledge-content.normal-document-content { padding: 0 24px; }
.other-library-content { display: flex; flex-direction: column; gap: 16px; padding: 24px; overflow: hidden; background: #fff; }
.library-filters { display: flex; flex-wrap: wrap; gap: 8px; flex-shrink: 0; button { min-height: 32px; padding: 0 12px; border: 1px solid transparent; border-radius: 6px; color: #5a5a5a; background: #f5f5f5; font-size: 14px; line-height: 22px; cursor: pointer; &.active { color: #146ef5; background: #e8f0fe; font-weight: 500; } &:focus-visible { outline: 2px solid #146ef5; outline-offset: -2px; } } }
.other-libraries-list { flex: 1; min-height: 0; overflow-y: auto; padding: 2px; }
.other-libraries-total { flex-shrink: 0; color: #5a5a5a; font-size: 12px; line-height: 20px; }
.knowledge-content.qa-document-content {
  padding: 0 32px;
  :deep(.group-content-box:not(.collapsed)) { padding: 32px 16px 32px 0; border-color: #e8e8e8; }
  :deep(.group-content-box) {
    .group-head-box { padding-right: 0; }
    .search-box { margin-top: 12px; margin-bottom: 6px; }
    .classify-item { height: 32px; margin-top: 2px; }
    .classify-item:first-child { margin-top: 0; }
  }
  :deep(.main-content-box) { padding: 32px 0 32px 24px; }
  :deep(.group-content-box.collapsed + .main-content-box) { padding-left: 0; }
  :deep(.head-content .title), :deep(.head-content .btn-right-box) { flex-wrap: nowrap; }
  :deep(.document-main-layout .document-pagination) { justify-content: space-between; padding-top: 16px; }
}
.library-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 16px; align-content: start; }
.library-card {
  position: relative; min-width: 0; min-height: 128px; padding: 16px; border: 1px solid #e8e8e8; border-radius: 12px;
  display: flex; flex-direction: column; background: #fff; transition: border-color .2s;
  &:hover, &:focus-within, &.actions-open, &.is-busy {
    border-color: #2475fc;
    .library-actions { opacity: 1; pointer-events: auto; }
    .library-heading { padding-right: 72px; }
  }
}
.library-badge { display: inline-block; flex-shrink: 0; min-height: 24px; border-radius: 4px; padding: 2px 6px; font-size: 12px; line-height: 20px; color: #5a5a5a; background: #f5f5f5; }
.linked-badge { color: #007012; background: #f0fdf2; }
.relation-button {
  display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0; width: 24px; height: 24px; padding: 0; border-radius: 6px; border-color: #2475fc; color: #fff; background: #2475fc; box-shadow: none;
  &.is-linked { color: #2475fc; background: #fff; }
  &:not(:disabled):hover { border-color: #2475fc; color: #fff; background: #2475fc; }
  &.is-linked:not(:disabled):hover { color: #2475fc; background: #fff; }
  &:disabled { opacity: .6; }
  :deep(.ant-btn-loading-icon) { margin: 0; }
}
.library-info {
  display: grid; grid-template-columns: 40px minmax(0, 1fr); align-items: center; flex: 1; gap: 12px; width: 100%; padding: 0; border: 0; background: transparent;
  text-align: left; cursor: pointer; color: #262626;
  &:focus-visible { outline: 2px solid var(--app-primary); outline-offset: 4px; border-radius: 12px; }
}
.library-avatar { flex-shrink: 0; border-radius: 8px; background: transparent; }
.library-copy { display: flex; flex-direction: column; gap: 8px; min-width: 0; }
.library-heading { display: flex; align-items: center; gap: 8px; min-width: 0; }
.library-name { min-width: 0; max-width: 100%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 16px; line-height: 24px; font-weight: 500; color: #080808; }
.library-actions { position: absolute; top: 34px; right: 16px; display: flex; align-items: center; gap: 8px; width: 64px; height: 28px; opacity: 0; pointer-events: none; transition: opacity .15s; }
.library-metrics {
  display: flex; flex-wrap: wrap; align-items: center; gap: 16px; font-size: 14px; line-height: 22px; color: #5a5a5a;
  span { display: inline-flex; align-items: center; gap: 4px; white-space: nowrap; }
  .svg-action { font-size: 16px; }
}
.library-more { display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0; width: 24px; height: 24px; padding: 0; border-radius: 6px; }
.library-tooltip-beak { position: absolute; top: 100%; left: 50%; transform: translateX(-50%) scaleY(-1); }
@media (max-width: 1280px) { .library-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
@media (max-width: 900px) { .library-grid { grid-template-columns: minmax(0, 1fr); } }
@media (hover: none) { .library-actions { opacity: 1; pointer-events: auto; } .library-heading { padding-right: 72px; } }
.upload-empty { max-width: 720px; margin: 60px auto; text-align: center; }
.upload-empty .ant-btn { margin: 20px 8px 12px; }
.upload-empty p { color: #8c8c8c; }
.feedback { padding: 16px 24px; display: flex; gap: 12px; }
</style>
<style lang="less">
.other-library-tooltip {
  .ant-tooltip-content { filter: drop-shadow(0 2px 4px #00000026); }
  .ant-tooltip-inner { position: relative; min-height: 34px; padding: 6px 8px; border-radius: 6px; color: #fff; background: rgba(0, 0, 0, .75); box-shadow: none; font-size: 14px; line-height: 22px; }
}
</style>
