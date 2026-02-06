<template>
  <div class="official-draft-wrapper">
        <div class="qa-list-wrapper">
          <div class="group-content-box" ref="groupContentRef"
            :class="[{ collapsed: isHiddenGroup }, { 'no-transition': isDragging }]" :style="isHiddenGroup
              ? {}
              : {
                width: groupBoxWidth + 'px',
                minWidth: minGroupBoxWidth + 'px',
                maxWidth: maxGroupBoxWidth + 'px',
                height: groupBoxHeight ? groupBoxHeight + 'px' : undefined
              }
              ">
            <div class="group-head-box">
              <div class="head-title">
                <a-tooltip :title="t('tooltip_collapse_group')">
                  <div class="hover-btn-wrap" @click="handleChangeHideStatus">
                    <svg-icon name="put-away"></svg-icon>
                  </div>
                </a-tooltip>
                <div class="flex-between-box">
                  <div>{{ t('draft_group_title') }}</div>
                  <a-tooltip :title="t('tooltip_new_group')">
                    <div class="hover-btn-wrap" @click="openGroupModal({})">
                      <PlusOutlined />
                    </div>
                  </a-tooltip>
                </div>
              </div>
              <div class="search-box">
                <a-input
                  v-model:value="groupSearchKey"
                  allowClear
                  :placeholder="t('search_group_placeholder')"
                  style="width: 100%"
                >
                  <template #suffix>
                    <SearchOutlined @click.stop="" />
                  </template>
                </a-input>
              </div>
            </div>
            <div class="classify-box">
              <cu-scroll style="padding-right: 24px">
                <div class="classify-item" @click="handleChangeGroup(item)" :class="{ active: item.id == groupId }"
                  v-for="item in filterGroupLists" :key="item.id">
                  <div class="classify-title">{{ item.group_name }}</div>
                  <div class="right-content">
                    <div class="num" :class="{ 'num-block': item.id <= 0 }">{{ item.total }}</div>
                    <div class="btn-box" v-if="item.id > 0">
                      <a-dropdown placement="bottomRight">
                        <div class="hover-btn-wrap">
                          <EllipsisOutlined />
                        </div>
                        <template #overlay>
                          <a-menu>
                            <a-menu-item>
                              <div @click.stop="openGroupModal(item)">重命名</div>
                            </a-menu-item>
                            <a-menu-item>
                              <div @click.stop="handleDelGroup(item)">删 除</div>
                            </a-menu-item>
                          </a-menu>
                        </template>
                      </a-dropdown>
                    </div>
                  </div>
                </div>
              </cu-scroll>
            </div>
            <div v-if="!isHiddenGroup" class="resize-bar" @mousedown="handleResizeMouseDown"></div>
          </div>

          <div class="main-content-box" ref="mainContentRef">
            <cu-scroll :scrollbar="{ minSize: 80, fade: false, interactive: true, scrollbarTrackClickable: true }">
              <div class="head-content">
                <div class="toolbar">
                  <div class="title" v-if="isHiddenGroup">
                    <a-popover :title="null" placement="bottom">
                      <template #content>
                        <div class="pover-group-content-box">
                          <div class="group-head-box">
                            <div class="head-title">
                              <div class="flex-between-box">
                                <div>{{ t('draft_group_title') }}</div>
                                <a-tooltip :title="t('tooltip_new_group')">
                                  <div class="hover-btn-wrap" @click="openGroupModal({})">
                                    <PlusOutlined />
                                  </div>
                                </a-tooltip>
                              </div>
                            </div>
                            <div class="search-box">
                              <a-input
                                allowClear
                                v-model:value="groupSearchKey"
                                :placeholder="t('search_group_placeholder')"
                                style="width: 100%"
                              >
                                <template #suffix>
                                  <SearchOutlined @click.stop="" />
                                </template>
                              </a-input>
                            </div>
                          </div>
                          <div class="classify-box classify-scroll-box">
                            <div class="classify-item" @click="handleChangeGroup(item)"
                              :class="{ active: item.id == groupId }" v-for="item in filterGroupLists" :key="item.id">
                              <div class="classify-title">{{ item.group_name }}</div>
                              <div class="right-content">
                                <div class="num" :class="{ 'num-block': item.id <= 0 }">{{ item.total }}</div>
                                <div class="btn-box" v-if="item.id > 0">
                                  <a-dropdown placement="bottomRight">
                                    <div class="hover-btn-wrap">
                                      <EllipsisOutlined />
                                    </div>
                                    <template #overlay>
                                      <a-menu>
                                        <a-menu-item>
                                          <div @click.stop="openGroupModal(item)">重命名</div>
                                        </a-menu-item>
                                        <a-menu-item>
                                          <div @click.stop="handleDelGroup(item)">删 除</div>
                                        </a-menu-item>
                                      </a-menu>
                                    </template>
                                  </a-dropdown>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </template>
                      <div class="hover-btn-wrap" @click="handleChangeHideStatus"><svg-icon name="expand"></svg-icon>
                      </div>
                    </a-popover>
                  </div>
                  <a-button type="primary" @click="openSyncModal" :loading="syncLoading">
                    <SyncOutlined v-if="!syncLoading" />{{ t('btn_sync_draft') }}
                  </a-button>
                  <a-button :disabled="selectedCount === 0" @click="handleBatchGroup">{{ t('btn_batch_group') }}</a-button>
                </div>
                <div class="header-row">
                  <a-checkbox v-model:checked="allChecked" />
                  <span class="header-text">{{ t('header_content') }}</span>
                </div>
              </div>

              <div class="content-block">
                <div v-if="isLoadingDrafts" class="loading-box"><a-spin /></div>
                <template v-else>
                  <div v-if="draftLists.length === 0" class="empty-box">
                    <img src="@/assets/empty.png" />
                    <div class="title">{{ t('empty_no_draft') }}</div>
                  </div>
                  <div v-else class="draft-list">
                    <div class="draft-item" v-for="item in draftLists" :key="item.id">
                      <a-checkbox v-model:checked="selectedMap[item.id]" />
                      <img v-if="item.thumb_url" class="thumb" :src="item.thumb_url" />
                      <img v-else class="thumb" src="@/assets/img/default-cover.png" />
                      <div class="info">
                        <div class="title">{{ item.title }}</div>
                        <div class="meta">{{ t('draft_group_label') }}{{ getGroupName(item.group_id) }}</div>
                        <div class="digest">{{ item.digest }}</div>
                      </div>
                      <div class="extra">
                        <svg-icon name="group-send" class="icon" @click="handleGroupSend(item)"
                          style="font-size: 32px; cursor: pointer; color: transparent;" />
                      </div>
                    </div>
                  </div>
                  <div class="pagination-box">
                    <a-pagination v-model:current="paginations.page" v-model:page-size="paginations.size" :total="total"
                      :pageSizeOptions="['10', '20', '50', '100']" show-size-changer @change="onShowSizeChange" />
                  </div>
                </template>
              </div>
            </cu-scroll>
          </div>
        </div>

        <a-modal v-model:open="groupModalOpen" :title="t('modal_add_group_title')" @ok="handleSaveGroup">
          <div style="margin: 32px 0 4px;">{{ t('modal_group_name_label') }}</div>
          <a-input
            v-model:value="groupModalState.group_name"
            :placeholder="t('modal_group_name_placeholder')"
            style="margin-bottom: 20px;"
          />
        </a-modal>

        <a-modal v-model:open="moveOpen" :title="t('modal_batch_move_group_title')" @ok="handleMove">
          <a-flex align="center" style="margin: 24px 0">
            <div>{{ t('modal_move_to_group_label') }}</div>
            <a-select v-model:value="moveState.group_id" style="flex: 1" :placeholder="t('modal_select_group_placeholder')">
              <a-select-option v-for="item in groupLists.filter((g) => g.id > 0)" :key="item.id" :value="item.id">{{
                item.group_name }}</a-select-option>
            </a-select>
          </a-flex>
        </a-modal>

        <CreateSendModal ref="createSendModalRef" :app-id="appId" :access-key="accessKey"
          :get-group-name="getGroupName" @created="onCreatedSend" />


        <a-modal v-model:open="syncModalOpen" :title="t('modal_sync_draft_title')" :width="478" @ok="onConfirmSyncDraft">
          <div class="sync-modal">
            <div class="sync-desc-box">
              {{ t('sync_desc') }}
            </div>
            <div class="sync-select-title">{{ t('sync_select_title') }}</div>
            <div class="sync-select-box">
              <a-slider :min="0" :max="3" :step="null" v-model:value="syncLimitIndex" :marks="syncMarks"
                :tooltipOpen="false" />
            </div>
          </div>
        </a-modal>
      
  </div>
</template>

<script setup>
import { reactive, ref, createVNode, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  getOfficialDraftGroupList,
  saveOfficialDraftGroup,
  deleteOfficialDraftGroup,
  moveOfficialDraftGroup,
  getOfficialDraftList,
  syncOfficialDraftList
} from '@/api/robot'
import CreateSendModal from './components/create-send-modal.vue'
import { PlusOutlined, EllipsisOutlined, SearchOutlined, ExclamationCircleOutlined, SyncOutlined } from '@ant-design/icons-vue'
import { addNoReferrerMeta, removeNoReferrerMeta } from '@/utils/index.js'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.explore.article-group-send.draft-box')
const props = defineProps({
  appId: { type: String, default: '' },
  accessKey: { type: String, default: '' }
})

const groupContentRef = ref(null)
const mainContentRef = ref(null)
const groupBoxHeight = ref(0)

const syncLoading = ref(false)

watch(() => props.appId, () => {
  paginations.value.page = 1
  if (props.appId) initPage()
})

// 左侧分组
const groupLists = ref([])
const groupSearchKey = ref('')
const groupId = ref(-1)

const filterGroupLists = computed(() => {
  return groupLists.value.filter((item) => (item.group_name || '').includes(groupSearchKey.value))
})

const loadGroupLists = async () => {
  const res = await getOfficialDraftGroupList()
  let lists = res.data || []
  groupLists.value = lists.map((g) => ({ ...g, total: 0 }))
  await loadGroupTotals()
  syncGroupHeight()
}

const loadGroupTotals = async () => {
  // 计算各分组数量（请求一次列表接口，仅取total）
  const appId = props.appId
  const promises = groupLists.value.map(async (g) => {
    const rsp = await getOfficialDraftList({ page: 1, size: 1, app_id: appId, group_id: g.id })
    const total = rsp?.data?.total || 0
    g.total = total
  })
  await Promise.all(promises)
}

const openGroupModal = (data = {}) => {
  groupModalState.id = data.id || 0
  groupModalState.group_name = data.group_name || ''
  groupModalOpen.value = true
}

const groupModalOpen = ref(false)
const groupModalState = reactive({ id: 0, group_name: '' })
const handleSaveGroup = async () => {
  if (!groupModalState.group_name) return message.error(t('error_enter_group_name'))
  await saveOfficialDraftGroup({ id: groupModalState.id, group_name: groupModalState.group_name })
  groupModalOpen.value = false
  message.success(t('message_operation_success'))
  await loadGroupLists()
}

const handleDelGroup = (item) => {
  Modal.confirm({
    title: t('modal_delete_group_title', { name: item.group_name }),
    icon: createVNode(ExclamationCircleOutlined),
    content: '',
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
    onOk: async () => {
      await deleteOfficialDraftGroup({ id: item.id })
      await loadGroupLists()
      message.success(t('message_operation_success'))
      if (groupId.value == item.id) {
        groupId.value = -1
        await loadDrafts()
      }
    }
  })
}

const handleChangeGroup = (item) => {
  groupId.value = item.id
  searchDrafts()
}

// 右侧列表
const draftLists = ref([])
const isLoadingDrafts = ref(false)
const total = ref(0)
const paginations = ref({ page: 1, size: 20 })
const filterData = reactive({ title: '' })

const getSelectedIds = () => Object.keys(selectedMap).filter((k) => selectedMap[k])
const selectedMap = reactive({})
const selectedCount = computed(() => getSelectedIds().length)
const allChecked = computed({
  get () {
    return draftLists.value.length > 0 && draftLists.value.every((it) => !!selectedMap[it.id])
  },
  set (val) {
    draftLists.value.forEach((it) => {
      selectedMap[it.id] = val
    })
  }
})

const resetSelection = () => {
  Object.keys(selectedMap).forEach((k) => delete selectedMap[k])
}

watch(draftLists, () => { nextTick(syncGroupHeight) }, { deep: true })
watch(groupLists, () => { nextTick(syncGroupHeight) }, { deep: true })

const loadDrafts = async () => {
  isLoadingDrafts.value = true
  const rsp = await getOfficialDraftList({
    page: paginations.value.page,
    size: paginations.value.size,
    title: filterData.title,
    app_id: props.appId,
    group_id: groupId.value
  })
  const data = rsp?.data || {}
  draftLists.value = data.list || []
  total.value = data.total || 0
  isLoadingDrafts.value = false
  resetSelection()
  syncGroupHeight()
}

const searchDrafts = () => {
  paginations.value.page = 1
  loadDrafts()
}

const onShowSizeChange = (current, pageSize) => {
  paginations.value.page = current
  paginations.value.size = pageSize
  loadDrafts()
}

const handleOpenMoveModal = () => {
  const ids = getSelectedIds()
  if (ids.length === 0) return message.error(t('error_select_draft_to_move'))
  moveState.ids = ids.join(',')
  moveState.group_id = void 0
  moveOpen.value = true
}

const handleBatchGroup = () => {
  handleOpenMoveModal()
}

const createSendModalRef = ref(null)
const handleGroupSend = (item) => { createSendModalRef.value && createSendModalRef.value.show({ draft: item }) }

const moveOpen = ref(false)
const moveState = reactive({ group_id: void 0, ids: '' })
const handleMove = async () => {
  if (!moveState.group_id) return message.error(t('error_select_group'))
  await moveOfficialDraftGroup({ draft_id: moveState.ids, group_id: moveState.group_id, app_id: props.appId })
  moveOpen.value = false
  message.success(t('message_move_success'))
  await loadDrafts()
  await loadGroupTotals()
}

// 折叠与宽度拖拽
const GROUP_BOX_WIDTH_KEY = 'official_draft_group_box_width'
const HIDE_KEY = 'official_draft_group_hide_key'
const minGroupBoxWidth = 200
const maxGroupBoxWidth = 256
const groupBoxWidth = ref(minGroupBoxWidth)
const isHiddenGroup = ref(localStorage.getItem(HIDE_KEY) == 1)
const isDragging = ref(false)

onMounted(() => {
  addNoReferrerMeta()
  const width = parseInt(localStorage.getItem(GROUP_BOX_WIDTH_KEY))
  if (width && width >= minGroupBoxWidth && width <= maxGroupBoxWidth) {
    groupBoxWidth.value = width
  }
  if (props.appId) initPage()
  window.addEventListener('resize', syncGroupHeight)
})

onUnmounted(() => { window.removeEventListener('resize', syncGroupHeight); removeNoReferrerMeta() })

const handleChangeHideStatus = () => {
  isHiddenGroup.value = !isHiddenGroup.value
  localStorage.setItem(HIDE_KEY, isHiddenGroup.value ? 1 : 2)
  syncGroupHeight()
}

let isResizing = false
let startX = 0
let startWidth = 0
const handleResizeMouseDown = (e) => {
  if (isHiddenGroup.value) return
  isResizing = true
  isDragging.value = true
  startX = e.clientX
  startWidth = groupBoxWidth.value
  document.body.style.cursor = 'col-resize'
  document.addEventListener('mousemove', handleResizing)
  document.addEventListener('mouseup', handleResizeMouseUp)
}
const handleResizing = (e) => {
  if (!isResizing) return
  let newWidth = startWidth + (e.clientX - startX)
  newWidth = Math.max(minGroupBoxWidth, Math.min(maxGroupBoxWidth, newWidth))
  groupBoxWidth.value = newWidth
}
const handleResizeMouseUp = () => {
  if (!isResizing) return
  isResizing = false
  isDragging.value = false
  localStorage.setItem(GROUP_BOX_WIDTH_KEY, groupBoxWidth.value)
  document.body.style.cursor = ''
  document.removeEventListener('mousemove', handleResizing)
  document.removeEventListener('mouseup', handleResizeMouseUp)
  syncGroupHeight()
}

const initPage = async () => {
  await loadGroupLists()
  await loadDrafts()
  syncGroupHeight()
}

const syncGroupHeight = () => {
  nextTick(() => {
    const gc = groupContentRef.value
    const mc = mainContentRef.value
    if (!gc || !mc) return
    const prevH = gc.style.height
    gc.style.height = 'auto'
    const gh = gc.scrollHeight || gc.offsetHeight || 0
    gc.style.height = prevH
    const mh = mc.scrollHeight || mc.offsetHeight || 0
    groupBoxHeight.value = Math.max(gh, mh)
  })
}

const getGroupName = (gid) => {
  const item = groupLists.value.find((g) => g.id == gid)
  if (item) return item.group_name
  if (gid == 0) return t('group_unassigned')
  return ''
}

const onCreatedSend = () => { message.success(t('message_create_send_success')) }


const syncModalOpen = ref(false)
const syncLimitIndex = ref(0)
const syncMarks = {
  0: t('sync_mark_100'),
  1: t('sync_mark_200'),
  2: t('sync_mark_1000'),
  3: t('sync_mark_all')
}
const syncOptions = [100, 200, 1000, 0]
const openSyncModal = () => {
  syncModalOpen.value = true
}
const onConfirmSyncDraft = async () => {
  const limit = syncOptions[syncLimitIndex.value] || 0
  syncLoading.value = true
  await syncOfficialDraftList({ access_key: props.accessKey, limit })
  syncLoading.value = false
  message.success(t('message_sync_submitted'))
  syncModalOpen.value = false
  setTimeout(() => { loadDrafts(); loadGroupTotals() }, 800)
  setTimeout(() => { syncGroupHeight() }, 1000)
}
</script>

<style lang="less" scoped>
.official-draft-wrapper {
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;

  .head-app-select {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 16px 24px 0 24px;

    .label {
      color: #595959;
    }
  }
}

.line-box {
  height: 1px;
  background: #F0F0F0;
  margin-top: 16px;
}

.qa-list-wrapper {
  flex: 1;
  overflow: hidden;
  display: flex;

  .group-content-box {
    position: relative;
    width: 200px;
    border-right: 1px solid #F0F0F0;
    padding-top: 24px;
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    transition:
      width 0.3s cubic-bezier(0.65, 0, 0.35, 1),
      transform 0.3s cubic-bezier(0.65, 0, 0.35, 1),
      padding 0.3s cubic-bezier(0.65, 0, 0.35, 1),
      border-right-width 0.3s cubic-bezier(0.65, 0, 0.35, 1);

    &.no-transition {
      transition: none !important;
    }

    &.collapsed {
      width: 0 !important;
      min-width: 0 !important;
      max-width: 0 !important;
      padding-left: 0;
      padding-right: 0;
      border-right-width: 0;
      transform: translateX(-100%);
    }

    .group-head-box {
      padding-right: 24px;
    }

    .head-title {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 16px;
      font-weight: 600;
      color: #262626;
    }

    .search-box {
      margin-top: 16px;
    }
  }

  .main-content-box {
    flex: 1;
    height: 100%;
    overflow: hidden;
    padding: 24px 0 0 24px;
    font-size: 14px;
    line-height: 22px;
    transition: all 0.3s ease;

    .title {
      display: flex;
      align-items: center;
      color: #262626;
      font-size: 16px;
      font-weight: 600;
      line-height: 24px;

      .btn-right-box {
        margin-left: auto;
        display: flex;
        align-items: center;
        gap: 8px;
      }
    }

    .toolbar {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .header-row {
      margin-top: 16px;
      display: flex;
      height: 54px;
      padding: 16px;
      align-items: center;
      gap: 16px;
      align-self: stretch;
      background: #F5F5F5;
    }

    .header-text {
      color: #262626;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
    }

    .content-block {
      margin-top: 8px;
    }

    .pagination-box {
      display: flex;
      align-items: center;
      justify-content: flex-end;
      margin-top: 16px;
    }
  }
}

.flex-between-box {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.classify-box {
  flex: 1;
  overflow: hidden;
  font-size: 14px;

  .classify-item {
    height: 32px;
    padding: 0 8px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 4px;
    cursor: pointer;
    border-radius: 6px;
    color: #595959;
    transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);

    .classify-title {
      flex: 1;
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: nowrap;
    }

    .num {
      display: block;
    }

    .btn-box {
      display: none;
    }

    &:hover {
      background: #f2f4f7;

      .num {
        display: none;
      }

      .num.num-block {
        display: block;
      }

      .btn-box {
        display: block;
      }
    }

    &.active {
      color: #2475fc;
      background: #e6efff;
    }
  }
}

.hover-btn-wrap {
  width: fit-content;
  height: 24px;
  border-radius: 6px;
  padding: 0 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);

  &:hover {
    background: #e4e6eb;
  }
}

.pover-group-content-box {
  width: 232px;
  padding: 4px;
  overflow: hidden;

  .head-title {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
  }

  .search-box {
    margin-top: 16px;
  }
}

.empty-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding-top: 148px;

  img {
    width: 200px;
    height: 200px;
  }

  .title {
    color: #262626;
    text-align: center;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    margin-bottom: 16px;
  }
}

.resize-bar {
  position: absolute;
  top: 0;
  right: 0;
  width: 4px;
  height: 100%;
  cursor: col-resize;
  z-index: 10;
  background: transparent;
  transition: background 0.2s;
}

.resize-bar:hover {
  background: #2475fc;
}

.group-content-box.collapsed .resize-bar {
  display: none;
}

.draft-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.draft-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 16px;
  border-bottom: 1px solid #D9D9D9;
}

.draft-item .thumb {
  width: 146px;
  height: 96px;
  border-radius: 4px;
  object-fit: cover;
}

.draft-item .info {
  flex: 1;
  overflow: hidden;
}

.draft-item .info .title {
  align-self: stretch;
  color: #262626;
  font-size: 16px;
  font-style: normal;
  font-weight: 600;
  line-height: 24px;
}

.draft-item .info .meta {
  color: #8c8c8c;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
  margin: 4px 0 6px;
}

.draft-item .info .digest {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  flex: 1 0 0;
  overflow: hidden;
  color: #595959;
  text-overflow: ellipsis;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}

.draft-item .icon:hover {
  opacity: 0.7;
}

.sync-modal {
  .sync-desc-box {
    background: #E6F4FF;
    border: 1px solid #99BFFD;
    color: #3A4559;
    padding: 9px 12px;
    border-radius: 6px;
    margin-bottom: 24px;
    line-height: 20px;
  }

  .sync-select-title {
    margin-bottom: 4px;
    color: #262626;
  }

  .sync-select-box {
    width: 388px;
    margin: 0 auto;
    padding-bottom: 20px;
  }
}

.loading-box {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 0;
}
</style>
