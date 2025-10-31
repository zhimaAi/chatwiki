<template>
  <div bg-color="#f5f9ff" class="library-page">
    <PageTabs class="mb-16" :tabs="pageTabs" active="/library/list"></PageTabs>

    <page-alert class="mb-16 mr-16" title="使用说明">
      <div>
        <p>
          1、知识库可以关联到聊天机器人中使用，创建机器人之前请先创建知识库，然后去机器人设置中关联。
        </p>
        <p>
          2、知识库支持普通知识库和问答知识库。普通知识库适用于非结构化数据，支持text/doc/pdf/md/html等格式文件，上传后系统会自动分段处理，也支持自定义分段。问答知识库适用于一问一答形式结构化数据，支持通过excel/word批量上传或自定义添加问题和答案。
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
              新建知识库
            </a-button>
            <template #overlay>
              <a-menu>
                <a-menu-item @click.prevent="handleAdd(0)">
                  <span class="create-action">
                    <img class="icon" :src="LIBRARY_NORMAL_AVATAR" alt="" />
                    <span>普通知识库</span>
                  </span>
                </a-menu-item>
                <a-menu-item @click.prevent="handleAdd(2)">
                  <span class="create-action">
                    <img class="icon" :src="LIBRARY_QA_AVATAR" alt="" />
                    <span>问答知识库</span>
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
                <div>知识库分组</div>
                <a-tooltip title="新建分组">
                  <div class="hover-btn-wrap" @click="openGroupModal({})"><PlusOutlined /></div>
                </a-tooltip>
              </div>
              <!-- <div class="search-box">
                <a-input
                  v-model:value="groupSearchKey"
                  allowClear
                  placeholder="搜索分组"
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
                </template>
              </draggable>
            </div>
          </cu-scroll>
          <a-tooltip placement="right" :title="isHideGroup ? '展开分组' : '收起分组'">
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
  </div>
</template>

<script setup>
import { ref, createVNode, computed } from 'vue'
import { useRouter } from 'vue-router'
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
  SearchOutlined,
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

const router = useRouter()

const pageTabs = ref([
  {
    title: '知识库',
    path: '/library/list'
  },
  {
    title: '数据库',
    path: '/database/list'
  },
  {
    title: '文档提取FAQ',
    path: '/ai-extract-faq/list'
  },
  {
    title: '触发次数统计',
    path: '/trigger-statics/list'
  },
])

const tabs = ref([
  {
    title: '全部 (0)',
    value: 'all'
  },
  {
    title: '普通知识库 (0)',
    value: '0'
  },
  {
    title: '问答知识库 (0)',
    value: '2'
  }
])

const permissionStore = usePermissionStore()
let { role_permission, role_type } = permissionStore
const libraryCreate = computed(() => role_type == 1 || role_permission.includes('LibraryCreate'))

const addLibraryModelRef = ref(null)
const addLibrayPopup = ref(null)
const activeKey = ref('all')

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
    // 计算每个分组的机器人数量
    lists.forEach((group) => {
      totalNumber += +group.total
    })
    groupLists.value = [
      {
        group_name: '全部',
        total: totalNumber,
        id: ''
      },
      ...lists
    ]
  })
}
getGroupList()
const handleAdd = (type) => {
  // addLibrayPopup.value.show(type)
  toAdd(type)
}

const toAdd = (val) => {
  addLibraryModelRef.value.show({ type: val, group_id: group_id.value })
}

const list = ref([])
const updateTabNumber = (data) => {
  let all = 0
  let normal = 0
  let qa = 0
  data.forEach((item) => {
    if (item.type == 0) {
      normal += 1
    }
    if (item.type == 2) {
      qa += 1
    }
    all += 1
  })

  tabs.value = [
    {
      title: '全部 (' + all + ')',
      value: 'all'
    },
    {
      title: '普通知识库 (' + normal + ')',
      value: '0'
    },
    {
      title: '问答知识库 (' + qa + ')',
      value: '2'
    }
  ]
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
    return message.error('您没有删除该知识库的权限')
  }
  let secondsToGo = 3

  let modal = Modal.confirm({
    title: `删除${data.library_name}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: '您确定要删除此知识库吗？',
    okText: secondsToGo + ' 确 定',
    okType: 'danger',
    okButtonProps: {
      disabled: true
    },
    cancelText: '取 消',
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
        okText: '确 定',
        okButtonProps: {
          disabled: false
        }
      })

      clearInterval(interval)
      interval = undefined
    } else {
      secondsToGo -= 1

      modal.update({
        okText: secondsToGo + ' 确 定',
        okButtonProps: {
          disabled: true
        }
      })
    }
  }, 1000)
}

const onDelete = ({ id }) => {
  deleteLibrary({ id }).then(() => {
    message.success('删除成功')
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
    title: `确认删除分组${item.group_name}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: '',
    okText: '确认',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      deleteLibraryListGroup({
        id: item.id
      }).then(() => {
        message.success('删除成功')
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
  // 只允许id>0的项目拖拽
  return e.draggedContext.element.id > 0 && e.relatedContext.element.id > 0
}
const handleDragEnd = async () => {
  try {
    // 过滤掉"全部分组"(-1)
    const sortList = groupLists.value
      .filter((item) => item.id > 0)
      .map((item, index) => ({
        id: item.id,
        sort: groupLists.value.length - index
      }))

    // 调用API保存排序
    await sortLibararyListGroup({
      sort_group: JSON.stringify(sortList)
    })
    message.success('排序已保存')
  } catch (error) {
    console.error('排序保存失败:', error)
    // 恢复原顺序
    getGroupLists()
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
