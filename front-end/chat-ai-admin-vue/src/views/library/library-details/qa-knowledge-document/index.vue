<template>
  <div class="qa-list-wrapper" :class="{ 'workspace-document': compact }">
    <div
      class="group-content-box"
      :inert="compact && isHiddenGroup ? '' : undefined"
      :aria-hidden="compact && isHiddenGroup ? true : undefined"
      :class="[{ collapsed: isHiddenGroup }, { 'no-transition': isDragging }]"
      :style="
        isHiddenGroup
          ? {}
          : {
              width: groupBoxWidth + 'px',
              minWidth: minGroupBoxWidth + 'px',
              maxWidth: maxGroupBoxWidth + 'px'
            }
      "
    >
      <div class="group-head-box">
        <div class="head-title">
          <button v-if="compact" type="button" class="workspace-icon-button collapse-group" :aria-label="t('collapse_group')" :title="t('collapse_group')" :aria-expanded="!isHiddenGroup" @click="handleChangeHideStatus">
            <img :src="icons.collapse" alt="" />
          </button>
          <a-tooltip v-if="!compact" :title="t('collapse_group')">
            <div class="hover-btn-wrap" @click="handleChangeHideStatus">
              <svg-icon name="put-away"></svg-icon>
            </div>
          </a-tooltip>
          <div class="flex-between-box">
            <div>{{ t('qa_group') }}</div>
            <a-tooltip :title="t('create_group')">
              <button v-if="compact" type="button" class="workspace-icon-button" :aria-label="t('create_group')" @click="openGroupModal({})"><img :src="icons.addGroup" alt="" /></button>
              <div v-else class="hover-btn-wrap" @click="openGroupModal({})"><PlusOutlined /></div>
            </a-tooltip>
          </div>
        </div>
        <div class="search-box">
          <a-input
            v-model:value="groupSearchKey"
            allowClear
            :placeholder="t('search_group')"
            style="width: 100%"
          >
            <template v-if="compact" #prefix><img class="workspace-icon" :src="icons.search" alt="" /></template>
            <template v-if="!compact" #suffix>
              <SearchOutlined @click.stop="" />
            </template>
          </a-input>
        </div>
      </div>
      <div class="classify-box">
        <cu-scroll :style="{ paddingRight: compact ? '0' : '24px' }">
          <div
            class="classify-item"
            @click="handleChangeGroup(item)"
            :class="{ active: item.id == groupId }"
            v-for="item in filterGroupLists"
            :key="item.id"
            :role="compact ? 'button' : undefined"
            :tabindex="compact ? 0 : undefined"
            @keydown.enter.prevent="compact && handleChangeGroup(item)"
            @keydown.space.prevent="compact && handleChangeGroup(item)"
          >
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
                        <div @click.stop="openGroupModal(item)">{{ t('rename') }}</div>
                      </a-menu-item>
                      <a-menu-item>
                        <div @click.stop="handleDelGroup(item)">{{ t('delete') }}</div>
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
    <div class="main-content-box">
      <component
        :is="compact ? 'div' : CuScroll"
        :class="{ 'document-main-layout': compact }"
        :scrollbar="compact ? undefined : { minSize: 80, fade: false, interactive: true, scrollbarTrackClickable: true }"
      >
        <div class="head-content">
          <div class="title">
            <a-popover v-if="compact && isHiddenGroup" v-model:open="groupPopoverOpen" :trigger="['hover', 'focus']" placement="bottomLeft" :arrow="false" :overlay-inner-style="{ padding: '6px' }" overlay-class-name="workspace-qa-group-popover">
              <template #content>
                <div class="workspace-group-options">
                  <button v-for="item in groupLists" :key="item.id" type="button" class="workspace-group-option" :class="{ active: item.id == groupId }" @click="selectPopoverGroup(item)">
                    <span>{{ item.group_name }}</span><span>{{ item.total }}</span>
                  </button>
                </div>
              </template>
              <button type="button" class="workspace-icon-button expand-group" :aria-label="t('expand_group')" :aria-expanded="!isHiddenGroup" @click="handleChangeHideStatus" @keydown.esc="groupPopoverOpen = false">
                <img :src="icons.collapse" alt="" />
              </button>
            </a-popover>
            <a-popover :title="null" placement="bottom" v-if="isHiddenGroup && !compact">
              <template #content>
                <div class="pover-group-content-box">
                  <div class="group-head-box">
                    <div class="head-title">
                      <div class="flex-between-box">
                        <div>{{ t('qa_group') }}</div>
                        <a-tooltip :title="t('create_group')">
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
                        :placeholder="t('search_group')"
                        style="width: 100%"
                      >
                        <template #suffix>
                          <SearchOutlined @click.stop="" />
                        </template>
                      </a-input>
                    </div>
                  </div>
                  <div class="classify-box classify-scroll-box">
                    <div
                      class="classify-item"
                      @click="handleChangeGroup(item)"
                      :class="{ active: item.id == groupId }"
                      v-for="item in filterGroupLists"
                      :key="item.id"
                    >
                      <div class="classify-title">{{ item.group_name }}</div>
                      <div class="right-content">
                        <div class="num" :class="{ 'num-block': item.id <= 0 }">
                          {{ item.total }}
                        </div>
                        <div class="btn-box" v-if="item.id > 0">
                          <a-dropdown placement="bottomRight">
                            <div class="hover-btn-wrap">
                              <EllipsisOutlined />
                            </div>
                            <template #overlay>
                              <a-menu>
                                <a-menu-item>
                                  <div @click.stop="openGroupModal(item)">{{ t('rename') }}</div>
                                </a-menu-item>
                                <a-menu-item>
                                  <div @click.stop="handleDelGroup(item)">{{ t('delete') }}</div>
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
              <div class="hover-btn-wrap" @click="handleChangeHideStatus">
                <svg-icon name="expand"></svg-icon>
              </div>
            </a-popover>
            <div class="current-group-name" :title="currentGroupItem.group_name">{{ currentGroupItem.group_name }}</div>
            <div v-if="compact" class="btn-right-box workspace-tools">
              <a-input class="workspace-qa-search" @change="search" v-model:value="filterData.search" :placeholder="t('search_qa')" :aria-label="t('search_qa')">
                <template #prefix><img class="workspace-icon" :src="icons.search" alt="" /></template>
              </a-input>
              <a-button type="text" class="workspace-tool" :title="t('metadata')" :aria-label="t('metadata')" @click="showMetaModal(1)"><img class="workspace-icon" :src="icons.metadata" alt="" /><span class="tool-label">{{ t('metadata') }}</span></a-button>
              <a-dropdown :trigger="['click']">
                <a-button type="text" class="workspace-tool" :title="t('batch_operation')" :aria-label="t('batch_operation')"><span class="workspace-batch-icon"><img :src="icons.batch" alt="" /></span><span class="tool-label">{{ t('batch_operation') }}</span></a-button>
                <template #overlay>
                  <a-menu>
                    <a-menu-item key="delete" @click="handleBathDel">{{ t('batch_delete') }}</a-menu-item>
                    <a-menu-item key="move" @click="handleOpenMoveModal">{{ t('batch_move') }}</a-menu-item>
                    <a-menu-item key="metadata" @click="showMetaModal(2)">{{ t('modify_metadata') }}</a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
              <a-button type="text" class="workspace-tool" :title="t('batch_import')" :aria-label="t('batch_import')" @click="handleOpenFileUploadModal()"><img class="workspace-icon" :src="icons.upload" alt="" /><span class="tool-label">{{ t('batch_import') }}</span></a-button>
              <a-button type="text" class="workspace-tool" :title="t('batch_export')" :aria-label="t('batch_export')" @click="handleSyncDownload()"><img class="workspace-icon" :src="icons.download" alt="" /><span class="tool-label">{{ t('batch_export') }}</span></a-button>
              <a-button type="primary" class="workspace-tool" :title="t('qa')" :aria-label="t('qa')" @click="openEditSubscription({})"><img class="workspace-icon" :src="icons.add" alt="" /><span class="tool-label">{{ t('qa') }}</span></a-button>
            </div>
            <div v-else class="btn-right-box">
              <div>
                <a-input @change="search" v-model:value="filterData.search" :placeholder="t('search_qa')">
                  <template #suffix>
                    <SearchOutlined />
                  </template>
                </a-input>
              </div>
              <a-button @click="showMetaModal(1)">{{ t('metadata') }} <SettingOutlined/></a-button>
              <a-dropdown>
               <a-button>{{ t('batch_operation') }} <DownOutlined/></a-button>
               <template #overlay>
                 <a-menu>
                   <a-menu-item key="1"><div @click="handleBathDel">{{ t('batch_delete') }}</div></a-menu-item>
                   <a-menu-item key="2"><div @click="handleOpenMoveModal">{{ t('batch_move') }}</div></a-menu-item>
                   <a-menu-item key="2"><div @click="showMetaModal(2)">{{ t('modify_metadata') }}</div></a-menu-item>
                 </a-menu>
               </template>
              </a-dropdown>
              <a-button @click="handleOpenFileUploadModal()">{{ t('batch_import') }}</a-button>
              <a-button @click="handleSyncDownload()">{{ t('batch_export') }}</a-button>
              <a-button
                @click="openEditSubscription({})"
                type="primary"
                :icon="createVNode(PlusOutlined)"
                >{{ t('qa') }}</a-button
              >
            </div>
          </div>
          <div v-if="!compact || exception_total > 0" class="des-content-block">
            <div class="alert-icon-box" v-if="exception_total > 0">
              <svg-icon name="alert-icon"></svg-icon>
              <div>{{ exception_total }}{{ t('learn_failed') }}</div>
              <a @click="reEmbeddingVectors">{{ t('re_learn') }}</a>
            </div>
          </div>
        </div>
        <div class="content-block">
          <SubsectionBox
            ref="subsectionBoxRef"
            :compact="compact"
            :fixed-height="compact"
            :total="total"
            :isLoading="isLoading"
            :paragraphLists="paragraphLists"
            :search="filterData.search"
            @openEditSubscription="openEditSubscription"
            @handleDelParagraph="initData"
            @getList="getParagraphLists"
            @handleSort="handleSort"
          ></SubsectionBox>
          <div class="pagination-box" :class="{ 'document-pagination': compact }">
            <span v-if="compact" class="pagination-total">{{ t('total_count', { total }) }}</span>
            <a-pagination
              v-model:current="paginations.page"
              v-model:page-size="paginations.size"
              :total="total"
              :pageSizeOptions="['100', '200', '500', '1000']"
              show-size-changer
              @change="onShowSizeChange"
            >
            </a-pagination>
          </div>
        </div>
      </component>
    </div>
  </div>
  <EditSubscription
    :detailsInfo="detailsInfo"
    @handleEdit="initData"
    ref="editSubscriptionRef"
  ></EditSubscription>
  <AddGroup ref="addGroupRef" @ok="initData" />
  <QaUploadModal :groupLists="groupLists" :library_id="library_id" @ok="initData" ref="qaUploadModalRef" />
  <a-modal v-model:open="moveOpen" :title="t('batch_move_group')" @ok="handleMove">
    <a-flex class="move-group-form" align="center" style="margin: 24px 0">
      <div class="form-label">{{ t('move_to_group') }}</div>
      <a-select v-model:value="moveState.group_id" style="flex: 1" :placeholder="t('select_group')">
        <a-select-option
          v-for="item in groupLists.filter((item) => item.id >= 0)"
          :key="item.id"
          :value="item.id"
          >{{ item.group_name }}</a-select-option
        >
      </a-select>
    </a-flex>
  </a-modal>
  <a-modal v-model:open="downLoadModalOpen" :title="null" :footer="null" :width="640">
    <a-result
      status="success"
      :title="t('export_task_created')"
      :sub-title="t('export_task_desc')"
    >
      <template #extra>
        <a-button style="margin-right: 16px;" @click="downLoadModalOpen = false">{{ t('got_it') }}</a-button>
        <a-button @click="toDownloadPage" type="primary">{{ t('go_download') }}</a-button>
      </template>
    </a-result>
  </a-modal>
  <MetadataManageModal ref="metaRef" :library-id="libraryId" :qa-ids="qaIds" @change="initData"/>
</template>

<script setup>
import { reactive, ref, createVNode, computed, onMounted, nextTick, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import dayjs from 'dayjs'
import {
  getLibraryGroup,
  reconstructCategoryVector,
  getParagraphList,
  deleteLibraryGroup,
  setParagraphGroup,
  deleteParagraph,
  createExportLibFileTask,
} from '@/api/library'
import {
  PlusOutlined,
  ExclamationCircleOutlined,
  EllipsisOutlined,
  SearchOutlined,
  DownOutlined,
  SettingOutlined,
} from '@ant-design/icons-vue'
import SubsectionBox from './components/subsection-box.vue'
import CuScroll from '@/components/cu-scroll/cu-scroll.vue'
import EditSubscription from '@/views/library/library-preview/components/edit-subsection.vue'
import QaUploadModal from '../components/qa-upload-modal.vue'
import AddGroup from './components/add-group.vue'
import icons from './knowledge-icons'
import router from '@/router'
import MetadataManageModal from "@/views/library/library-details/components/metadata-manage-modal.vue";
import { useI18n } from '@/hooks/web/useI18n'
import {KNOWLEDGE_SOURCE_TYPE_MAP} from "@/constants/index.js";

const { t } = useI18n('views.library.library-details.qa-knowledge-document.index')

const SOURCE_TYPE_MAP = KNOWLEDGE_SOURCE_TYPE_MAP()
const subsectionBoxRef = ref(null)
const metaRef = ref(null)
const route = useRoute()
const query = route.query
const groupSearchKey = ref('')
const groupLists = ref([])
const groupId = ref(-1)
const groupPopoverOpen = ref(false)
const qaIds = ref([])

const currentGroupItem = computed(() => {
  return groupLists.value.find((item) => item.id == groupId.value) || {}
})

const filterGroupLists = computed(() => {
  return groupLists.value.filter((item) => item.group_name.includes(groupSearchKey.value))
})

const props = defineProps({
  compact: { type: Boolean, default: false },
  library_id: {
    type: [Number, String],
    default: () => ''
  }
})

const libraryId = computed(() => {
  return props.library_id || query.id
})

const paragraphLists = ref([])

const GROUP_HIDE_KEY = props.compact ? 'customer_service_qa_document_group_hide_key' : 'qa_document_group_hide_key'
const isHiddenGroup = ref(localStorage.getItem(GROUP_HIDE_KEY) == 1)

const total = ref(0)
const exception_total = ref(0)
const listStatusMap = {
  0: t('not_converted'),
  1: t('converted'),
  2: t('convert_error'),
  3: t('converting')
}
const detailsInfo = ref({
  is_qa_doc: '1',
  library_id: libraryId.value
})

const paginations = ref({
  page: 1,
  size: 100
})

const filterData = reactive({
  status: -1,
  category_id: -1,
  sort_field: '',
  sort_type: '',
  library_id: libraryId.value,
  search: '',
})

const handleSort = (sort) => {
  filterData.sort_field = sort.sort_field
  filterData.sort_type = sort.sort_type
  search()
}

const search = () => {
  paginations.value.page = 1
  getParagraphLists()
}

const onShowSizeChange = (current, pageSize) => {
  paginations.value.page = current
  paginations.value.size = pageSize
  getParagraphLists()
}

const handleChangeHideStatus = () => {
  groupPopoverOpen.value = false
  isHiddenGroup.value = !isHiddenGroup.value
  localStorage.setItem(GROUP_HIDE_KEY, isHiddenGroup.value ? 1 : 2)
}

const reEmbeddingVectors = () => {
  reconstructCategoryVector({ id: libraryId.value }).then(() => {
    message.success(t('operation_complete'))
    getParagraphLists()
  })
}

const editSubscriptionRef = ref(null)
const openEditSubscription = (data) => {
  if (!data.group_id) {
    if (groupId.value >= 0) {
      data.group_id = groupId.value
    }
  }
  editSubscriptionRef.value.showModal(JSON.parse(JSON.stringify(data)))
}

// 监听跨窗口消息，处理QA合并后的刷新
const handleMessage = (event) => {
  if (event.data?.type === 'qa-merged' && event.data?.libraryId == libraryId.value) {
    // 获取当前编辑的QA ID和弹窗状态
    const currentId = editSubscriptionRef.value?.id
    const isModalOpen = editSubscriptionRef.value?.open
    
    // 如果弹窗是打开状态，等待列表刷新后重新打开弹窗
    if (currentId && isModalOpen) {
      getParagraphLists().then(() => {
        nextTick(() => {
          const currentData = paragraphLists.value.find(item => item.id == currentId)
          if (currentData) {
            openEditSubscription(currentData)
          }
        })
      })
    } else {
      // 弹窗未打开，只刷新列表
      getParagraphLists()
    }
  }
}

onMounted(() => {
  window.addEventListener('message', handleMessage)
})

onUnmounted(() => {
  window.removeEventListener('message', handleMessage)
})

const isLoading = ref(false)

const getParagraphLists = () => {
  isLoading.value = true
  return getParagraphList({
    ...paginations.value,
    ...filterData,
    group_id: groupId.value
  }).then((res) => {
    isLoading.value = false
    let data = res.data
    let list = data.list || []
    const lastPage = Math.max(1, Math.ceil(Number(data.total) / paginations.value.size))
    if (paginations.value.page > lastPage) {
      paginations.value.page = lastPage
      return getParagraphLists()
    }
    list.forEach((item) => {
      item.status_text = listStatusMap[item.status]
      if (item.similar_questions) {
        item.similar_questions = JSON.parse(item.similar_questions)
      }
      if (Array.isArray(item.meta_list)) {
        item.meta_list.forEach(i => {
          if (i.type == 1 && i.value > 0) {
            i.value = dayjs(i.value * 1000).format('YYYY-MM-DD HH:mm')
          }
          if (i.key == 'source') {
            i.value = SOURCE_TYPE_MAP[i.value]
          }
          item[`meta_${i.key}`] = i.value
        })
      }
    })
    paragraphLists.value = list
    total.value = data.total
    exception_total.value = data.exception_total
    subsectionBoxRef.value && subsectionBoxRef.value.resetSelect()
  })
}

getParagraphLists()

const getGroupLists = () => {
  getLibraryGroup({
    library_id: libraryId.value
  }).then((res) => {
    let lists = res.data || []
    let allTotal = lists.reduce((total, item) => {
      return total + +item.total
    }, 0)
    groupLists.value = [
      {
        group_name: t('all_groups'),
        id: -1,
        total: allTotal
      },
      ...res.data
    ]
  })
}

getGroupLists()

const addGroupRef = ref(null)
const openGroupModal = (data) => {
  addGroupRef.value.show({
    ...data,
    library_id: libraryId.value
  })
}
const initData = () => {
  getGroupLists()
  getParagraphLists()
}

const handleChangeGroup = (item) => {
  groupId.value = item.id
  search()
}

const selectPopoverGroup = (item) => {
  groupPopoverOpen.value = false
  handleChangeGroup(item)
}

const handleDelGroup = (item) => {
  Modal.confirm({
    title: `${t('confirm_delete_group')}${item.group_name}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: '',
    okText: t('confirm'),
    okType: 'danger',
    cancelText: t('cancel'),
    onOk() {
      deleteLibraryGroup({
        id: item.id
      }).then(() => {
        message.success(t('delete_success'))
        getGroupLists()
        if (groupId.value == item.id) {
          groupId.value = -1
          search()
        }
      })
    }
  })
}

const qaUploadModalRef = ref(null)

const handleOpenFileUploadModal = () => {
  qaUploadModalRef.value.show(groupId.value)
}

const moveOpen = ref(false)
const moveState = reactive({
  group_id: void 0,
  ids: ''
})

const handleOpenMoveModal = () => {
  if (subsectionBoxRef.value.state.selectedRowKeys.length == 0) {
    return message.error(t('select_qa_to_move'))
  }
  moveState.ids = subsectionBoxRef.value.state.selectedRowKeys.join(',')
  moveState.group_id = void 0
  moveOpen.value = true
}

const handleBathDel = () => {
  if (subsectionBoxRef.value.state.selectedRowKeys.length == 0) {
    return message.error(t('select_qa_to_delete'))
  }

  Modal.confirm({
    title: t('tip'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirm_delete_selected'),
    onOk() {
      return deleteParagraph({ id: subsectionBoxRef.value.state.selectedRowKeys.join(',') }).then(
        (res) => {
          message.success(t('delete_success'))
          initData()
        }
      )
    },
    onCancel() {}
  })
}

const handleMove = () => {
  if (moveState.group_id === undefined || moveState.group_id === null) {
    return message.error(t('select_group_tip'))
  }
  setParagraphGroup({
    ...moveState
  }).then(() => {
    moveOpen.value = false
    message.success(t('move_success'))
    initData()
  })
}

const GROUP_BOX_WIDTH_KEY = props.compact ? 'customer_service_qa_document_group_box_width' : 'qa_document_group_box_width'
const minGroupBoxWidth = 200
const maxGroupBoxWidth = props.compact ? 320 : 256

const groupBoxWidth = ref(minGroupBoxWidth)
const isDragging = ref(false) // 新增

onMounted(() => {
  const width = parseInt(localStorage.getItem(GROUP_BOX_WIDTH_KEY))
  if (width && width >= minGroupBoxWidth && width <= maxGroupBoxWidth) {
    groupBoxWidth.value = width
  }
})

let isResizing = false
let startX = 0
let startWidth = 0

const handleResizeMouseDown = (e) => {
  if (isHiddenGroup.value) return
  isResizing = true
  isDragging.value = true // 新增
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
  isDragging.value = false // 新增
  localStorage.setItem(GROUP_BOX_WIDTH_KEY, groupBoxWidth.value)
  document.body.style.cursor = ''
  document.removeEventListener('mousemove', handleResizing)
  document.removeEventListener('mouseup', handleResizeMouseUp)
}

onUnmounted(handleResizeMouseUp)


const downLoadModalOpen = ref(false)
const handleSyncDownload = () => {
  let data_ids = ''
  if (subsectionBoxRef.value.state.selectedRowKeys.length > 0) {
    data_ids = subsectionBoxRef.value.state.selectedRowKeys.join(',')
  }
  createExportLibFileTask({
    library_id: libraryId.value,
    export_type: 2,
    group_id: groupId.value,
    data_ids,
  }).then(res => {
    downLoadModalOpen.value = true;
  })
}

const toDownloadPage = ()=>{
  router.push({
    path: '/library/details/export-record',
    query:{
      id: libraryId.value
    },
  })
}

const showMetaModal = (type) => {
  if (type == 2) {
    if (subsectionBoxRef.value.state.selectedRowKeys.length == 0) {
      return message.error(t('select_qa_to_modify'))
    }
    qaIds.value = subsectionBoxRef.value.state.selectedRowKeys
  }
  metaRef.value.show(type == 2)
}
</script>

<style lang="less" scoped>
.workspace-icon { flex-shrink: 0; vertical-align: middle; }
.workspace-icon-button {
  display: inline-flex; flex-shrink: 0; align-items: center; justify-content: center; width: 24px; height: 24px; padding: 0;
  border: 0; border-radius: 6px; background: transparent; cursor: pointer;
  &:hover { background: #f0f0f0; }
  &:focus-visible { outline: 2px solid #146ef5; outline-offset: 2px; }
}
.workspace-group-options { width: 228px; max-height: 360px; overflow-y: auto; }
.workspace-group-option {
  display: flex; align-items: center; justify-content: space-between; gap: 12px; width: 100%; min-height: 36px; padding: 6px 12px;
  border: 0; border-radius: 4px; background: transparent; color: #242424; text-align: left; cursor: pointer;
  span:first-child { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  span:last-child { color: #666; }
  &:hover { background: #f5f5f5; }
  &.active { color: #1456f0; background: #e8eeff; }
  &:focus-visible { outline: 2px solid #146ef5; outline-offset: -2px; }
}
.qa-list-wrapper.workspace-document {
  .group-content-box .head-title { gap: 8px; height: 24px; }
  .head-content .title { gap: 8px; }
  .current-group-name { flex-shrink: 1; min-width: 64px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .expand-group { margin-right: 0; img { transform: scaleX(-1); } }
  .classify-item.active { color: #146ef5; background: #e8efff; }
  .classify-item:focus-visible { outline: 2px solid #146ef5; outline-offset: -2px; }
  .classify-item .num { color: #666; }
  .head-content .workspace-tools { min-width: 0; flex-shrink: 1; font-weight: 400; }
  .workspace-qa-search { flex: 0 1 208px; min-width: 120px; width: 208px; height: 32px; border-radius: 6px; font-weight: 400; }
  .workspace-tool {
    display: inline-flex; flex-shrink: 0; align-items: center; justify-content: center; gap: 6px; height: 32px; padding: 0 12px; border-radius: 6px; color: #5a5a5a; box-shadow: none;
    :deep(> span) { margin-inline-start: 0; }
    &.ant-btn-primary { color: #fff; background: #146ef5; border-color: #146ef5; }
    &.ant-btn-primary:hover { background: #3686ff; border-color: #3686ff; }
  }
  .workspace-batch-icon { display: inline-flex; align-items: center; justify-content: center; width: 16px; height: 16px; }
  .main-content-box .content-block { margin-top: 16px; }
  .pagination-total { color: #666; font-size: 12px; line-height: 20px; }
  :deep(.ant-input-affix-wrapper) { height: 32px; padding: 4px 10px; border-color: #d8d8d8; border-radius: 6px; box-shadow: none; }
  :deep(.ant-input-prefix) { margin-inline-end: 8px; }
  :deep(.ant-input::placeholder) { color: #ababab; }
  :deep(.ant-pagination) { display: flex; align-items: center; margin-bottom: 0; }
}
@media (max-width: 1439.98px) {
  .qa-list-wrapper.workspace-document .workspace-tool { width: 32px; padding: 0; .tool-label { display: none; } }
}

.qa-list-wrapper {
  height: 100%;
  overflow: hidden;
  display: flex;

  .group-content-box {
    position: relative;
    width: 200px;
    border-right: 1px solid #d9d9d9;
    padding-top: 24px;
    padding-left: 24px;
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
    .head-btn-box {
      display: flex;
      gap: 16px;
      margin-top: 16px;
      div {
        flex: 1;
      }
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
    .head-content {
      .title {
        display: flex;
        align-items: center;
        color: #262626;
        font-size: 16px;
        font-weight: 600;
        line-height: 24px;
      }
      .btn-right-box {
        margin-left: auto;
        display: flex;
        align-items: center;
        gap: 8px;
      }
      .des-content-block {
        display: flex;
        align-items: center;
        margin-top: 8px;
        color: #8c8c8c;
      }
      .alert-icon-box {
        display: flex;
        align-items: center;
        gap: 8px;
        color: #fb363f;
      }
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
  .head-btn-box {
    display: flex;
    gap: 16px;
    margin-top: 16px;
    div {
      flex: 1;
    }
  }
  .search-box {
    margin-top: 16px;
  }
}

.classify-scroll-box {
  max-height: 400px;
  min-height: 180px;
  margin-top: 4px;
  overflow: hidden;
  overflow-y: auto;
  /* 整个页面的滚动条 */
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

.empty-box {
  text-align: center;
  height: 100%;
  padding-top: 148px;
  img {
    width: 200px;
    height: 200px;
  }
  .title {
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    color: #262626;
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

.move-group-form{
  .form-label{
    padding-right: 4px;
  }
}
</style>
