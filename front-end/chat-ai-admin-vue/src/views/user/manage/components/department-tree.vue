<template>
  <div class="department-tree-search">
    <div class="flex-block">
      <a-input-search v-model:value="searchValue" placeholder="搜索部门" style="flex: 1" />
      <div class="add-btn-box" @click="openAddModal({})"><PlusOutlined /></div>
    </div>
    <div class="tree-list-box">
      <a-tree
        :expanded-keys="expandedKeys"
        v-model:selectedKeys="selectedKeys"
        :auto-expand-parent="autoExpandParent"
        @select="onSelect"
        :tree-data="gData"
        blockNode
        @expand="onExpand"
      >
        <template #switcherIcon="{ switcherCls }">
          <div style="margin-right: 4px">
            <span class="action-btn toggle-btn">
              <CaretDownOutlined style="font-size: 12px" :class="switcherCls"
            /></span>
          </div>
        </template>
        <template #title="treeNode">
          <div class="doc-item">
            <span class="tree-title" v-if="treeNode.title.indexOf(searchValue) > -1 && searchValue">
              <template v-if="treeNode.title.length > 3 && treeNode.level > 3">
                <a-tooltip>
                  <template #title>{{ treeNode.title }}</template>
                  {{ treeNode.title.substring(0, treeNode.title.indexOf(searchValue)) }}
                  <span style="color: #f50">{{ searchValue }}</span>
                  {{
                    treeNode.title.substring(
                      treeNode.title.indexOf(searchValue) + searchValue.length
                    )
                  }}
                </a-tooltip>
              </template>
              <template v-else>
                {{ treeNode.title.substring(0, treeNode.title.indexOf(searchValue)) }}
                <span style="color: #f50">{{ searchValue }}</span>
                {{
                  treeNode.title.substring(treeNode.title.indexOf(searchValue) + searchValue.length)
                }}
              </template>
            </span>
            <span class="tree-title" v-else>
              <template v-if="treeNode.title.length > 3 && treeNode.level > 3">
                <a-tooltip>
                  <template #title>{{ treeNode.title }}</template>
                  {{ treeNode.title }}
                </a-tooltip>
              </template>
              <template v-else>{{ treeNode.title }}</template>
            </span>
            <div class="use-num">
              {{ treeNode.children_nums }}
            </div>
            <div class="doc-action-box" v-if="!treeNode.is_user">
              <a-dropdown>
                <template #overlay>
                  <a-menu>
                    <a-menu-item
                      key="addDoc"
                      @click.stop.prevent="openAddModal(treeNode)"
                    >
                      <span class="menu-name">添加子部门</span>
                    </a-menu-item>
                    <a-menu-item key="importDoc" @click.stop.prevent="handleRename(treeNode)">
                      <span class="menu-name">重命名</span>
                    </a-menu-item>
                    <a-menu-item
                      v-if="treeNode.is_default != 1"
                      key="del"
                      @click.stop.prevent="openDel(treeNode)"
                    >
                      <span class="menu-name" style="color: #fb363f">删 除</span>
                    </a-menu-item>
                  </a-menu>
                </template>
                <span class="action-btn" @click.stop=""><MoreOutlined /></span>
              </a-dropdown>
            </div>
          </div>
        </template>
      </a-tree>
    </div>
  </div>
  <AddDepartment ref="addDepartmentRef" @ok="handleSave" />
  <DelDepartment ref="delDepartmentRef" @ok="handleDelDepartment" />
</template>

<script setup>
import { PlusOutlined, CaretDownOutlined, MoreOutlined } from '@ant-design/icons-vue'
import { getDepartmentList } from '@/api/department/index.js'
import AddDepartment from './add-department.vue'
import { nextTick, ref, watch } from 'vue'
import { formateDepartmentData } from '@/utils/index.js'
import DelDepartment from './del-department.vue'

const emit = defineEmits(['select', 'save', 'refresh'])

let genData = []

const dataList = []
const getParentKey = (key, tree) => {
  let parentKey
  for (let i = 0; i < tree.length; i++) {
    const node = tree[i]
    if (node.children) {
      if (node.children.some((item) => item.key === key)) {
        parentKey = node.key
      } else if (getParentKey(key, node.children)) {
        parentKey = getParentKey(key, node.children)
      }
    }
  }
  return parentKey
}
const expandedKeys = ref([])
const selectedKeys = ref([])
const searchValue = ref('')
const autoExpandParent = ref(true)
const gData = ref([])

let defaltGroupId = '' // 默认分组的id
let defaultItem = {}
const onExpand = (keys) => {
  expandedKeys.value = keys
  autoExpandParent.value = false
}
watch(searchValue, (value) => {
  const expanded = dataList
    .map((item) => {
      if (item.title.indexOf(value) > -1) {
        return getParentKey(item.key, gData.value)
      }
      return null
    })
    .filter((item, i, self) => item && self.indexOf(item) === i)
  expandedKeys.value = expanded
  searchValue.value = value
  autoExpandParent.value = true
})

const generateList = (data) => {
  for (let i = 0; i < data.length; i++) {
    const node = data[i]
    const key = node.key
    dataList.push({
      key,
      title: node.title
    })
    if (node.children) {
      generateList(node.children)
    }
  }
}

const getLists = (isInit) => {
  getDepartmentList({}).then((res) => {
    let treeData = res.data || []

    genData = []
    genData = formateDepartmentData(treeData)
    generateList(genData)
    gData.value = genData
    if (isInit) {
      treeData.forEach((item) => {
        if (item.is_default == 1) {
          defaltGroupId = item.id
          defaultItem = item
        }
      })
      nextTick(() => {
        if (defaltGroupId) {
          selectedKeys.value = [defaltGroupId]
          emit('select', defaultItem)
        }
      })
    }
  })
}
getLists(true)

let treeSelectData = {}

const handleDelDepartment = (id) => {
  // 删除部门
  if (selectedKeys.value[0] == id) {
    nextTick(() => {
      if (defaltGroupId) {
        selectedKeys.value = [defaltGroupId]
        emit('select', defaultItem)
      }
    })
  } else {
    emit('refresh')
  }
  getLists()
}

const addDepartmentRef = ref()
const openAddModal = (data) => {
  // console.log(gData.value, genData, '==')
  // console.log(autoExpandParent.value, '===')
  // console.log(selectedKeys.value, '===')
  addDepartmentRef.value.show(data)
}

const handleRename = (data) => {
  addDepartmentRef.value.rename(data)
}

const onSelect = (_, e) => {
  let data = {
    user_id: '',
    department_id: ''
  }
  if (e.node.is_user) {
    data.user_id = e.node.id
  } else {
    data.department_id = e.node.id
  }
  if (e.selected) {
  } else {
    nextTick(() => {
      selectedKeys.value = [e.node.id]
    })
  }
  treeSelectData = {
    ...e.node
  }
  emit('select', treeSelectData)
}

const delDepartmentRef = ref(null)
const openDel = (treeNode) => {
  delDepartmentRef.value.show(treeNode)
}

const handleSave = () => {
  getLists()
  emit('refresh')
}

defineExpose({
  getLists
})
</script>

<style lang="less" scoped>
.add-btn-box {
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  width: 32px;
  height: 32px;
  font-size: 16px;
  cursor: pointer;
  &:hover {
    background: #e4e6eb;
  }
}
.flex-block {
  display: flex;
  align-items: center;
  gap: 8px;
}
.department-tree-search {
}

.tree-list-box {
  margin-top: 16px;
  ::v-deep(.ant-tree) {
    .action-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 22px;
      height: 24px;
      line-height: 24px;
      border-radius: 6px;
      cursor: pointer;
      transition: all 0.2s;
      font-size: 16px;
      color: #595959;

      &:hover {
        background-color: #e4e6eb;
      }
    }
    .ant-tree-treenode {
      min-height: 32px;
      &:hover {
        background: #f2f4f7;
        border-radius: 6px;
      }
      &.ant-tree-treenode-selected {
        background: #e5efff;
        color: #2475fc;
        border-radius: 6px;
      }
      .ant-tree-node-content-wrapper.ant-tree-node-selected {
        background: unset;
      }
      .ant-tree-switcher {
        line-height: 32px;
        display: flex;
        align-items: center;
        padding-left: 4px;
      }
      .ant-tree-node-content-wrapper {
        line-height: 32px;
        overflow: hidden;
        &:hover {
          background: unset;
        }
      }
    }

    .doc-item {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .tree-title {
        display: inline-block;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .doc-action-box {
        width: 24px;
        height: 24px;
        display: none;
        transition: all 0.2s cubic-bezier(0.075, 0.82, 0.165, 1);
      }
      &:hover {
        .doc-action-box {
          display: block;
        }
        .use-num {
          display: none;
        }
      }
    }
  }
}
.menu-name {
  text-align: center;
  display: block;
  width: 100%;
}
</style>
