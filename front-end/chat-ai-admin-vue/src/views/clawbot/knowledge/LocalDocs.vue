<template>
  <div class="local-docs">
    <div class="toolbar-row">
      <div class="section-desc">{{ t('section_desc') }}</div>

      <div class="toolbar">
        <div class="view-switch">
          <div class="view-btn" :class="{ active: viewMode === 'list' }" @click="handleChangeViewMode('list')">
            <UnorderedListOutlined />
          </div>
          <div class="view-btn" :class="{ active: viewMode === 'card' }" @click="handleChangeViewMode('card')">
            <AppstoreOutlined />
          </div>
        </div>
        <a-button type="primary" class="upload-btn" :loading="uploading" @click="handleUploadClick">
          {{ t('btn_upload_file') }}
        </a-button>
      </div>
    </div>

    <!-- 列表视图 -->
    <div v-if="viewMode === 'list'" class="file-list">
      <div v-for="item in fileList" :key="item.name" class="file-item">
        <div class="file-left">
          <div class="file-type" :class="item.ext">{{ item.ext.toUpperCase() }}</div>
          <div>
            <div class="file-name">{{ item.name }}</div>
            <div class="file-size">{{ formatSize(item.size) }}</div>
          </div>
        </div>
        <a-button danger ghost class="remove-btn" :loading="deletingMap[item.name]" @click="handleRemove(item)">
          {{ t('btn_remove') }}
        </a-button>
      </div>
      <div v-if="!fileList.length && !loading" class="empty-tip">
        {{ t('empty_tip') }}
      </div>
    </div>

    <!-- 卡片视图 -->
    <div v-else class="file-grid">
      <div v-for="item in fileList" :key="item.name" class="file-card">
        <div class="file-card-cover">
          <img :src="clawbotCardCover" alt="" />
        </div>
        <div class="file-card-name">{{ item.name }}</div>
        <div class="file-card-footer">
          <div class="file-card-ext">
            <svg-icon :name="getFileCardIconName(item.ext)" class="file-card-ext-icon" />
            <span>{{ item.ext.toLowerCase() }}</span>
          </div>
          <div v-if="formatDate(item.time)" class="file-card-date">{{ formatDate(item.time) }}</div>
        </div>
        <a-dropdown :trigger="['click']" placement="bottomRight">
          <div class="file-card-remove" @click.stop>
            <EllipsisOutlined />
          </div>
          <template #overlay>
            <a-menu>
              <a-menu-item key="delete" @click="handleRemove(item)">
                <span class="remove-menu-text">{{ t('action_delete') }}</span>
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </div>
      <div v-if="!fileList.length && !loading" class="empty-tip">
        {{ t('empty_tip') }}
      </div>
    </div>

    <LocalDocUploadModal
      v-model:open="uploadModalOpen"
      :loading="uploading"
      @confirm="handleUploadConfirm"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, createVNode } from 'vue'
import { UnorderedListOutlined, AppstoreOutlined, ExclamationCircleOutlined, EllipsisOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import { storeToRefs } from 'pinia'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { uploadClawbotLocalDoc, getClawbotLocalDocList, deleteClawbotLocalDoc } from '@/api/clawbot'
import dayjs from 'dayjs'
import clawbotCardCover from '@/assets/img/clawbot/01.png'
import LocalDocUploadModal from '@/views/clawbot/components/LocalDocUploadModal.vue'

const { t } = useI18n('views.clawbot.knowledge.LocalDocs')
const clawbotStore = useClawbotStore()
const { currentAssistant } = storeToRefs(clawbotStore)

const fileList = ref([])
const loading = ref(false)
const uploading = ref(false)
const deletingMap = ref({})
const uploadModalOpen = ref(false)
const LOCAL_DOCS_VIEW_MODE_KEY = 'clawbot-local-docs-view-mode'
const viewMode = ref(getCachedViewMode())

function getCachedViewMode() {
  try {
    const cached = localStorage.getItem(LOCAL_DOCS_VIEW_MODE_KEY)
    return cached === 'card' ? 'card' : 'list'
  } catch {
    return 'list'
  }
}

function handleChangeViewMode(mode) {
  viewMode.value = mode === 'card' ? 'card' : 'list'
  try {
    localStorage.setItem(LOCAL_DOCS_VIEW_MODE_KEY, viewMode.value)
  } catch {
    // localStorage 不可用时仅回退为当前会话内状态
  }
}

function formatSize(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatDate(time) {
  if (!time) return ''
  const value = dayjs(time)
  return value.isValid() ? value.format('YYYY/MM/DD') : ''
}

function getFileCardIconName(ext) {
  const normalizedExt = String(ext || '').toLowerCase()

  if (normalizedExt === 'csv') return 'csv'
  if (normalizedExt === 'md' || normalizedExt === 'txt') return 'document'
  if (normalizedExt === 'pdf') return 'pdf'
  if (normalizedExt === 'xlsx' || normalizedExt === 'xls') return 'excel'
  if (normalizedExt === 'docx' || normalizedExt === 'doc') return 'docx'

  return 'document'
}

async function fetchFileList() {
  const id = currentAssistant.value?.id
  if (!id) return
  loading.value = true
  try {
    const res = await getClawbotLocalDocList({ id })
    if (res && res.res === 0) {
      fileList.value = (res.data || []).map((item) => ({
        name: item.name,
        size: item.size,
        time: item.time,
        ext: String(item.ext || '').toLowerCase()
      }))
    } else {
      message.error(res?.msg || t('msg_fetch_failed'))
    }
  } catch (err) {
    console.error('获取文档列表失败', err)
    message.error(t('msg_fetch_failed'))
  } finally {
    loading.value = false
  }
}

function handleUploadClick() {
  uploadModalOpen.value = true
}

async function handleUploadConfirm({ file, description, keywords }) {
  // 文件名不能包含路径
  const pureName = file.name.replace(/[\\/]/g, '')
  if (!pureName || pureName === '.' || pureName === '..') {
    message.error(t('msg_invalid_name'))
    return
  }

  const id = currentAssistant.value?.id
  if (!id) {
    message.error(t('msg_no_agent_selected'))
    return
  }

  uploading.value = true
  try {
    const res = await uploadClawbotLocalDoc({
      id,
      file,
      description,
      'keywords[]': keywords
    })
    if (res && res.res === 0) {
      message.success(t('msg_upload_success'))
      uploadModalOpen.value = false
      await fetchFileList()
    } else {
      message.error(res?.msg || t('msg_upload_failed'))
    }
  } catch (err) {
    console.error('上传失败', err)
    message.error(t('msg_upload_failed'))
  } finally {
    uploading.value = false
  }
}

function handleRemove(item) {
  Modal.confirm({
    title: t('title_delete_document'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_confirm_delete_document', { name: item.name }),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    okType: 'danger',
    onOk: async () => {
      const id = currentAssistant.value?.id
      if (!id) return
      deletingMap.value[item.name] = true
      try {
        const res = await deleteClawbotLocalDoc({ id, name: item.name })
        if (res && res.res === 0) {
          message.success(t('msg_delete_success'))
          await fetchFileList()
        } else {
          message.error(res?.msg || t('msg_delete_failed'))
        }
      } catch (err) {
        console.error('删除失败', err)
        message.error(t('msg_delete_failed'))
      } finally {
        deletingMap.value[item.name] = false
      }
    }
  })
}

onMounted(() => {
  fetchFileList()
})
</script>

<style lang="less" scoped>
.local-docs {
  min-height: calc(100vh - 146px);
}

.toolbar-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 8px 0 16px;
}

.section-desc {
  color: #8c8c8c;
  font-size: 14px;
  line-height: 22px;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
}

.view-switch {
  display: flex;
  overflow: hidden;
  border-radius: 6px;
  border: 1px solid #d9d9d9;
  background: #fff;
}

.view-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #8c8c8c;
  background: #fff;
  cursor: pointer;
  transition: all 0.2s ease;

  & + .view-btn {
    border-left: 1px solid #d9d9d9;
  }

  &.active {
    color: #2475fc;
    background: #f5f8ff;
  }
}

.upload-btn {
  height: 32px;
  padding: 0 16px;
  border-radius: 6px;
  border: none;
  box-shadow: none;
  font-size: 14px;
  line-height: 22px;
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.file-item {
  height: 60px;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
  padding: 0 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
}

.file-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.file-type {
  width: 38px;
  height: 38px;
  border-radius: 8px;
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;

  &.pdf {
    background: #ff6b6b;
  }

  &.docx {
    background: #4c7cf7;
  }

  &.md {
    background: #7a4ce0;
  }

  &.xlsx {
    background: #2ecc71;
  }

  &.txt {
    background: #95a5a6;
  }
}

.file-name {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
}

.file-size {
  margin-top: 2px;
  color: #999;
  font-size: 12px;
}

.remove-btn {
  border-radius: 6px;
  height: 28px;
  font-size: 13px;
}

.remove-menu-text {
  color: #ff4d4f;
}

.empty-tip {
  text-align: center;
  padding: 40px 0;
  color: #999;
  font-size: 14px;
}

// 卡片视图
.file-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(164px, 164px));
  gap: 28px 24px;
  align-items: start;
  justify-content: start;
}

.file-card {
  position: relative;
  width: 164px;
  min-height: 178px;
  padding: 7px;
  border-radius: 12px;
  border: 1px solid #f0f0f0;
  background: #fff;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;

  &:hover {
    border-color: #d9d9d9;
    box-shadow: 0 6px 16px rgba(15, 23, 42, 0.08);
  }
}

.file-card-cover {
  width: 100%;
  height: auto;
  aspect-ratio: 150 / 86;
  overflow: hidden;
  border-radius: 6px;
  border: 1px solid #f0f0f0;
  background: #fff;

  img {
    width: 100%;
    height: 100%;
    display: block;
    object-fit: cover;
  }
}

.file-card-name {
  margin-top: 10px;
  min-height: 44px;
  color: #242933;
  font-size: 14px;
  line-height: 22px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  word-break: break-all;
}

.file-card-footer {
  margin-top: auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding-top: 8px;
}

.file-card-ext,
.file-card-date {
  color: #8c8c8c;
  font-size: 12px;
  line-height: 20px;
}

.file-card-ext {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.file-card-ext-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.file-card-remove {
  position: absolute;
  top: 9px;
  right: 9px;
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(0, 0, 0, 0.45);
  background: rgba(0, 0, 0, 0.25);
  cursor: pointer;
  transition: background-color 0.2s ease, color 0.2s ease;

  &:hover {
    background: rgba(0, 0, 0, 0.38);
    color: #fff;
  }

  :deep(svg) {
    width: 16px;
    height: 16px;
  }
}

@media (max-width: 1200px) {
  .file-grid {
    grid-template-columns: repeat(auto-fill, minmax(148px, 1fr));
    gap: 20px 16px;
  }

  .file-card {
    width: 100%;
    min-height: 170px;
  }
}

@media (max-width: 768px) {
  .toolbar-row {
    align-items: stretch;
    flex-direction: column;
    gap: 12px;
  }

  .toolbar {
    justify-content: flex-end;
  }

  .file-grid {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 16px 12px;
  }

  .file-card {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .file-grid {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
