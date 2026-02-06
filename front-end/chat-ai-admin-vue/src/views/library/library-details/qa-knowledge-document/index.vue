<template>
  <div class="qa-list-wrapper">
    <div
      class="group-content-box"
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
          <a-tooltip :title="t('collapse_group')">
            <div class="hover-btn-wrap" @click="handleChangeHideStatus">
              <svg-icon name="put-away"></svg-icon>
            </div>
          </a-tooltip>
          <div class="flex-between-box">
            <div>{{ t('qa_group') }}</div>
            <a-tooltip :title="t('create_group')">
              <div class="hover-btn-wrap" @click="openGroupModal({})"><PlusOutlined /></div>
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
            <template #suffix>
              <SearchOutlined @click.stop="" />
            </template>
          </a-input>
        </div>
      </div>
      <div class="classify-box">
        <cu-scroll style="padding-right: 24px">
          <div
            class="classify-item"
            @click="handleChangeGroup(item)"
            :class="{ active: item.id == groupId }"
            v-for="item in filterGroupLists"
            :key="item.id"
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
      <cu-scroll
        :scrollbar="{ minSize: 80, fade: false, interactive: true, scrollbarTrackClickable: true }"
      >
        <div class="head-content">
          <div class="title">
            <a-popover :title="null" placement="bottom" v-if="isHiddenGroup">
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
            <div>{{ currentGroupItem.group_name }}</div>
            <div class="btn-right-box">
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
          <div class="des-content-block">
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
            :total="total"
            :isLoading="isLoading"
            :paragraphLists="paragraphLists"
            :search="filterData.search"
            @openEditSubscription="openEditSubscription"
            @handleDelParagraph="initData"
            @getList="getParagraphLists"
            @handleSort="handleSort"
          ></SubsectionBox>
          <div class="pagination-box">
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
      </cu-scroll>
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
import { reactive, ref, createVNode, computed, onMounted } from 'vue'
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
import EditSubscription from '@/views/library/library-preview/components/edit-subsection.vue'
import QaUploadModal from '../components/qa-upload-modal.vue'
import AddGroup from './components/add-group.vue'
import router from '@/router'
import MetadataManageModal from "@/views/library/library-details/components/metadata-manage-modal.vue";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-details.qa-knowledge-document.index')

const subsectionBoxRef = ref(null)
const metaRef = ref(null)
const route = useRoute()
const query = route.query
const groupSearchKey = ref('')
const groupLists = ref([])
const groupId = ref(-1)
const qaIds = ref([])

const currentGroupItem = computed(() => {
  return groupLists.value.find((item) => item.id == groupId.value) || {}
})

const filterGroupLists = computed(() => {
  return groupLists.value.filter((item) => item.group_name.includes(groupSearchKey.value))
})

const props = defineProps({
  library_id: {
    type: [Number, String],
    default: () => ''
  }
})

const libraryId = computed(() => {
  return props.library_id || query.id
})

const paragraphLists = ref([])

const isHiddenGroup = ref(localStorage.getItem('qa_document_group_hide_key') == 1)

const total = ref(0)
const exception_total = ref(0)

const sourceTypeMap = {
  1: '本地文档',
  2: '在线文档',
  3: '自定义文档',
  4: '手工新增问答',
  5: '导入问答'
}
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
  isHiddenGroup.value = !isHiddenGroup.value
  localStorage.setItem('qa_document_group_hide_key', isHiddenGroup.value ? 1 : 2)
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

const isLoading = ref(false)

const getParagraphLists = () => {
  isLoading.value = true
  getParagraphList({
    ...paginations.value,
    ...filterData,
    group_id: groupId.value
  }).then((res) => {
    isLoading.value = false
    let data = res.data
    let list = data.list || []
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
            i.value = sourceTypeMap[i.value]
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
      deleteParagraph({ id: subsectionBoxRef.value.state.selectedRowKeys.join(',') }).then(
        (res) => {
          message.success(t('delete_success'))
          getParagraphLists()
        }
      )
    },
    onCancel() {}
  })
}

const handleMove = () => {
  if (!moveState.group_id) {
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

const GROUP_BOX_WIDTH_KEY = 'qa_document_group_box_width'
const minGroupBoxWidth = 200
const maxGroupBoxWidth = 256

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


const downLoadModalOpen = ref(false)
const handleSyncDownload = () => {
  let data_ids = ''
  if (subsectionBoxRef.value.state.selectedRowKeys.length > 0) {
    data_ids = subsectionBoxRef.value.state.selectedRowKeys.join(',')
  }
  createExportLibFileTask({
    library_id: query.id,
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
    query,
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
