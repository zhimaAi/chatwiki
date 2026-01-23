<template>
  <div bg-color="#f5f9ff" class="library-page">
    <PageTabs class="mb-16" :tabs="pageTabs" active="/library/list"></PageTabs>

    <page-alert class="mb-16 mr-16" :title="t('instruction_title')">
      <div>
        <p>
          {{ t('instruction_content_1') }}
        </p>
        <p>
          {{ t('instruction_content_2') }}
        </p>
      </div>
    </page-alert>

    <div class="library-page-body">
      <div class="list-toolbar">
        <div class="toolbar-box">
          <ListTabs :tabs="tabs" v-model:value="activeKey" @change="onChangeTab" />
        </div>

        <div class="toolbar-box">
          <a-dropdown v-if="libraryCreate">
            <a-button type="primary" @click.prevent="">
              <template #icon>
                <PlusOutlined />
              </template>
              {{ t('create_library_btn') }}
            </a-button>
            <template #overlay>
              <a-menu>
                <a-menu-item @click.prevent="handleAdd(0)">
                  <span class="create-action">
                    <img class="icon" :src="LIBRARY_NORMAL_AVATAR" alt="" />
                    <span>{{ t('normal_library') }}</span>
                  </span>
                </a-menu-item>
                <a-menu-item @click.prevent="handleAdd(2)">
                  <span class="create-action">
                    <img class="icon" :src="LIBRARY_QA_AVATAR" alt="" />
                    <span>{{ t('qa_library') }}</span>
                  </span>
                </a-menu-item>
                <a-menu-item v-if="wxAppLibary" @click.prevent="handleAdd(3)">
                  <span class="create-action">
                    <img class="icon" src="@/assets/svg/library_ability_official_account.svg" alt="" />
                    <span>{{ t('official_account_library') }}</span>
                  </span>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>

      <div class="list-group-box">
        <div class="group-list-box" :class="{ 'hide-group': isHideGroup }">
          <cu-scroll style="padding-right: 16px">
            <div class="group-head-box">
              <div class="head-title" v-if="!isHideGroup">
                <div>{{ t('library_groups') }}</div>
                <a-tooltip :title="t('create_group')">
                  <div class="hover-btn-wrap" @click="openGroupModal({})"><PlusOutlined /></div>
                </a-tooltip>
              </div>
              <!-- <div class="search-box">
                <a-input
                  v-model:value="groupSearchKey"
                  allowClear
                  :placeholder="t('search_groups')"
                  style="width: 100%"
                >
                  <template #suffix>
                    <SearchOutlined @click.stop="" />
                  </template>
                </a-input>
              </div> -->
            </div>
            <div class="classify-box" style="margin-top: 16px">
              <draggable
                v-model="groupLists"
                item-key="id"
                @end="handleDragEnd"
                :move="checkMove"
                handle=".drag-btn"
              >
                <template #item="{ element: item }">
                  <div
                    class="classify-item"
                    @click="handleChangeGroup(item)"
                    :class="{ active: item.id == group_id }"
                  >
                    <div class="classify-title">
                      <span v-if="item.id > 0" class="drag-btn">
                        <svg-icon name="drag" />
                      </span>
                      <span v-else style="width: 20px; display: inline-block"></span>
                      {{ item.group_name }}
                    </div>
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
                </template>
              </draggable>
            </div>
          </cu-scroll>
          <a-tooltip placement="right" :title="isHideGroup ? t('expand_groups') : t('collapse_groups')">
            <div class="hide-group-box" @click="handleChangeHideGroup">
              <LeftOutlined v-if="!isHideGroup" />
              <RightOutlined v-else />
            </div>
          </a-tooltip>
        </div>
        <cu-scroll style="padding-right: 16px; flex: 1">
          <LibraryList
            :show-create="libraryCreate"
            :list="list"
            @edit="toEdit"
            @delete="handleDelete"
            @openEditGroupModal="openEditGroupModal"
          />
        </cu-scroll>
      </div>
    </div>
    <AddLibrayPopup ref="addLibrayPopup" @ok="toAdd" />
    <AddLibraryModel ref="addLibraryModelRef" />
    <AddGroup ref="addGroupRef" @ok="onChangeTab" />
    <EditGroup ref="editGroupRef" @ok="onChangeTab" />
    <SelectWechatApp
      ref="wxAppRef"
      :title="t('create_library_btn')"
      :okText="t('next_step')"
      :disabled-app-ids="wxAppids"
      @ok="showWxLibraryModel"
    />
  </div>
</template>

<script setup>
import { ref, createVNode, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import {
  getLibraryList,
  deleteLibrary,
  getLibraryListGroup,
  deleteLibraryListGroup,
  sortLibararyListGroup
} from '@/api/library'
import { formatFileSize } from '@/utils/index'
import { usePermissionStore } from '@/stores/modules/permission'
import { LIBRARY_NORMAL_AVATAR, LIBRARY_QA_AVATAR } from '@/constants/index'
import {
  ExclamationCircleOutlined,
  PlusOutlined,
  EllipsisOutlined,
  LeftOutlined,
  RightOutlined
} from '@ant-design/icons-vue'
import LibraryList from './components/libray-list/index.vue'
import AddLibrayPopup from './components/add-libray-popup.vue'
import AddLibraryModel from '@/views/library/add-library/add-library-model.vue'
import PageTabs from '@/components/cu-tabs/page-tabs.vue'
import PageAlert from '@/components/page-alert/page-alert.vue'
import ListTabs from '@/components/cu-tabs/list-tabs.vue'
import { getLibraryPermission } from '@/utils/permission'
import AddGroup from '@/views/library/library-list/components/add-group.vue'
import EditGroup from '@/views/library/library-list/components/edit-group.vue'
import Draggable from 'vuedraggable'
import {getSpecifyAbilityConfig, getAbilityList} from "@/api/explore/index.js";
import SelectWechatApp from "@/components/common/select-wechat-app.vue";
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.library.library-list.library-list');

const router = useRouter()
const route = useRoute()

const pageTabs = ref([
  {
    title: t('library_tab'),
    path: '/library/list'
  },
  {
    title: t('database_tab'),
    path: '/database/list'
  },
  {
    title: t('faq_extraction_tab'),
    path: '/ai-extract-faq/list'
  },
  {
    title: t('trigger_statistics_tab'),
    path: '/trigger-statics/list'
  },
])

const tabs = ref([
  {
    title: t('all_count', { count: 0 }),
    value: 'all'
  },
  {
    title: t('normal_count', { count: 0 }),
    value: '0'
  },
  {
    title: t('qa_count', { count: 0 }),
    value: '2'
  }
])

const permissionStore = usePermissionStore()
let { role_permission, role_type } = permissionStore
const libraryCreate = computed(() => role_type == 1 || role_permission.includes('LibraryCreate'))

const addLibraryModelRef = ref(null)
const addLibrayPopup = ref(null)
const wxAppRef = ref(null)
const wxAppLibary = ref(null)
const wxAppids = ref([])
const activeKey = ref(String(route.query?.active || 'all'))

const group_id = ref('')

let hideGroupLocalKey = 'library-list-hide-group-key'
const isHideGroup = ref(localStorage.getItem(hideGroupLocalKey) == 1)

const handleChangeHideGroup = () => {
  isHideGroup.value = !isHideGroup.value
  localStorage.setItem(hideGroupLocalKey, isHideGroup.value ? 1 : 0)
}

const groupSearchKey = ref('')
const groupLists = ref([])

const filterGroupLists = computed(() => {
  return groupLists.value.filter((item) => item.group_name.includes(groupSearchKey.value))
})

const getGroupList = () => {
  let type = activeKey.value === 'all' ? '' : activeKey.value
  getLibraryListGroup({
    type
  }).then((res) => {
    let lists = res.data || []
    let totalNumber = 0
    // Calculate total number for each group
    lists.forEach((group) => {
      totalNumber += +group.total
    })
    groupLists.value = [
      {
        group_name: t('all_tab'),
        total: totalNumber,
        id: ''
      },
      ...lists
    ]
  })
}
getGroupList()
const handleAdd = (type) => {
  if (type == 3) {
    wxAppRef.value.open()
  } else {
    toAdd(type)
  }
}

const toAdd = (val) => {
  addLibraryModelRef.value.show({ type: val, group_id: group_id.value })
}

const showWxLibraryModel = (ids, rows) => {
  if (!rows.length) return message.warning(t('select_official_account'))
  wxAppRef.value.close()
  addLibraryModelRef.value.show({ type: 3, group_id: group_id.value, wx_app_ids: ids})
}

const list = ref([])
const updateTabNumber = async (data) => {
  // Check if official account library is enabled
  await getSpecifyAbilityConfig({ability_type: 'library_ability_official_account'}).then((res) => {
    let _data = res?.data || {}
    if (_data?.user_config?.switch_status == 1) {
      wxAppLibary.value = _data
    }
  })

  let all = 0
  let normal = 0
  let qa = 0
  let wx = 0
  wxAppids.value = []
  data.forEach((item) => {
    if (item.type == 0) {
      normal += 1
    }
    if (item.type == 2) {
      qa += 1
    }
    if (item.type == 3) {
      wx += 1
      wxAppids.value.push(item.official_app_id)
    }
    all += 1
  })

  tabs.value = [
    {
      title: t('all_count', { count: all }),
      value: 'all'
    },
    {
      title: t('normal_count', { count: normal }),
      value: '0'
    },
    {
      title: t('qa_count', { count: qa }),
      value: '2'
    }
  ]
  if (wxAppLibary.value) {
    tabs.value.push({
      title: t('official_account_count', { count: wx }),
      value: '3'
    })
  }
}

let allList = ref([])
const getList = () => {
  let type = activeKey.value === 'all' ? '' : activeKey.value

  getLibraryList({ type }).then((res) => {
    let data = res.data || []

    data.forEach((item) => {
      item.file_size_str = formatFileSize(item.file_size)

      if (!item.avatar) {
        item.avatar = item.type == 0 ? LIBRARY_NORMAL_AVATAR : LIBRARY_QA_AVATAR
      }
    })
    if (group_id.value != '') {
      data = data.filter((item) => item.group_id == group_id.value)
    }

    list.value = data
    if (activeKey.value === 'all') {
      allList.value = res.data
      updateTabNumber(allList.value)
    } else {
      getLibraryList().then((res) => {
        allList.value = res.data
        updateTabNumber(allList.value)
      })
    }
  })
}

getList()

const onChangeTab = () => {
  getList()
  getGroupList()
}

const toEdit = (data) => {
  if (data.type == '1') {
    // router.push({
    //   path: '/public-library/config',
    //   query: {
    //     library_id: data.id
    //   }
    // })
    window.open(`/#/public-library/config?library_id=${data.id}`, '_blank', 'noopener') // 建议添加 noopener 防止安全漏洞
  } else {
    // router.push({
    //   name: 'libraryDetails',
    //   query: {
    //     id: data.id
    //   }
    // })
    window.open(`/#/library/details?id=${data.id}`, '_blank', 'noopener') // 建议添加 noopener 防止安全漏洞
  }
}

const handleDelete = (data) => {
  let key = getLibraryPermission(data.id)
  if (key != 4) {
    return message.error(t('no_permission_delete'))
  }
  let secondsToGo = 3

  let modal = Modal.confirm({
    title: t('confirm_delete_library_title', { library_name: data.library_name }),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirm_delete_content'),
    okText: secondsToGo + ' ' + t('confirm_btn'),
    okType: 'danger',
    okButtonProps: {
      disabled: true
    },
    cancelText: t('cancel_btn'),
    onOk() {
      onDelete(data)
    },
    onCancel() {
      // console.log('Cancel')
    }
  })

  let interval = setInterval(() => {
    if (secondsToGo == 1) {
      modal.update({
        okText: t('confirm_btn'),
        okButtonProps: {
          disabled: false
        }
      })

      clearInterval(interval)
      interval = undefined
    } else {
      secondsToGo -= 1

      modal.update({
        okText: secondsToGo + ' ' + t('confirm_btn'),
        okButtonProps: {
          disabled: true
        }
      })
    }
  }, 1000)
}

const onDelete = ({ id }) => {
  deleteLibrary({ id }).then(() => {
    message.success(t('delete_success'))
    getList()
    getGroupList()
  })
}

const handleChangeGroup = (item) => {
  group_id.value = item.id
  getList()
}

const addGroupRef = ref(null)
const openGroupModal = (data) => {
  addGroupRef.value.show({
    ...data
  })
}

const editGroupRef = ref(null)
const openEditGroupModal = (item) => {
  editGroupRef.value.show({
    ...item
  })
}

const handleDelGroup = (item) => {
  Modal.confirm({
    title: t('confirm_delete_group_title', { group_name: item.group_name }),
    icon: createVNode(ExclamationCircleOutlined),
    content: '',
    okText: t('confirm_btn'),
    okType: 'danger',
    cancelText: t('cancel_btn'),
    onOk() {
      deleteLibraryListGroup({
        id: item.id
      }).then(() => {
        message.success(t('delete_success'))
        getGroupList()
        if (group_id.value == item.id) {
          group_id.value = ''
          getList()
        } else {
          getList()
        }
      })
    }
  })
}

const isDragEnabled = ref(false)
const checkMove = (e) => {
  // Only allow dragging items with id > 0
  return e.draggedContext.element.id > 0 && e.relatedContext.element.id > 0
}
const handleDragEnd = async () => {
  try {
    // Filter out "All Groups" (id <= 0)
    const sortList = groupLists.value
      .filter((item) => item.id > 0)
      .map((item, index) => ({
        id: item.id,
        sort: groupLists.value.length - index
      }))

    // Call API to save sorting
    await sortLibararyListGroup({
      sort_group: JSON.stringify(sortList)
    })
    message.success(t('sort_saved'))
  } catch (error) {
    console.error(t('sort_failed'), error)
    // Restore original order
    getGroupList()
  }
}
</script>

<style lang="less" scoped>
.library-page {
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  .list-toolbar {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
  }
}
.create-action {
  display: flex;
  align-items: center;
  .icon {
    width: 20px;
    height: 20px;
    margin-right: 8px;
  }
}

.toolbar-box {
  padding-right: 16px;
}

.library-page-body {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.list-group-box {
  display: flex;
  gap: 16px;
  flex: 1;
  overflow: hidden;
}
.group-list-box {
  width: 256px;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  margin-top: 8px;
  padding: 16px;
  padding-right: 0;
  position: relative;
  &.hide-group {
    width: 0;
    padding: 0;
    border-left: 0;
    border-top: 0;
    border-bottom: 0;
  }
  .hide-group-box {
    position: absolute;
    right: -8px;
    top: 40%;
    height: 50px;
    width: 13px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #e5e5ea;
    color: #8c8c8c;
    cursor: pointer;
    opacity: 0.78;
    &:hover {
      opacity: 1;
    }
  }
  .head-title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 4px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
  }

  .search-box {
    margin-top: 16px;
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

      .drag-btn {
        margin-right: 8px;
        cursor: grab;
        opacity: 0;
        transition: opacity 0.2s;
      }
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
        .drag-btn {
          opacity: 1;
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
}

.mr-16 {
  margin-right: 16px;
}

// 大于1920px
@media screen and (min-width: 1920px) {
  .library-page {
  }
}
</style>
