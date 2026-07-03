<template>
  <div class="goods-group-tree">
    <div class="tree-head">
      <div class="tree-title">{{ t('group_tree.title') }}</div>
      <a-tooltip :title="t('group_tree.new_root')">
        <a-button type="text" class="head-action" @click="handleCreateRoot">
          <template #icon>
            <PlusOutlined />
          </template>
        </a-button>
      </a-tooltip>
    </div>

    <div class="fixed-group-list">
      <div
        class="group-node fixed-node"
        :class="{ active: isSelected('all') }"
        @click.stop="handleSelectFixed('all')"
      >
        <div class="node-main">
          <span class="group-name">{{ t('group_tree.all') }}</span>
        </div>
        <span class="group-count">{{ totalCount }}</span>
      </div>
      <div
        class="group-node fixed-node"
        :class="{ active: isSelected('0') }"
        @click.stop="handleSelectFixed('0')"
      >
        <div class="node-main">
          <span class="group-name">{{ t('group_tree.ungrouped') }}</span>
        </div>
        <span class="group-count">{{ ungroupedCount }}</span>
      </div>
    </div>

    <a-tree
      v-if="treeData.length"
      :tree-data="treeData"
      :expanded-keys="expandedKeys"
      :selected-keys="selectedKeys"
      :draggable="true"
      :allow-drop="handleAllowDrop"
      block-node
      @expand="handleExpand"
      @select="handleTreeSelect"
      @drop="handleDrop"
      @dragstart="handleDragStart"
      @dragend="handleDragEnd"
    >
      <template #switcherIcon="{ expanded, isLeaf }">
        <span v-if="isLeaf" class="switcher-placeholder"></span>
        <CaretDownOutlined v-else-if="expanded" />
        <CaretRightOutlined v-else />
      </template>
      <template #title="node">
        <div
          class="tree-node-content"
          :class="{ 'action-open': isActionOpen(node.key) }"
        >
          <span class="node-name">{{ node.group_name }}</span>
          <div class="node-actions">
            <span class="node-count">{{ node.goods_count ?? 0 }}</span>
            <a-dropdown
              placement="bottomRight"
              trigger="click"
              overlayClassName="goods-group-tree-dropdown"
              @openChange="handleActionOpenChange(node.key, $event)"
            >
              <button type="button" class="icon-btn action-btn" @click.stop>
                <EllipsisOutlined />
              </button>
              <template #overlay>
                <a-menu>
                  <a-menu-item v-if="node.level < 4" @click.stop="handleCreateChild(node)">
                    <span>{{ t('group_tree.new_child') }}</span>
                  </a-menu-item>
                  <a-menu-item @click.stop="handleEdit(node)">
                    <span>{{ t('group_tree.edit') }}</span>
                  </a-menu-item>
                  <a-menu-divider />
                  <a-menu-item class="danger-item" @click.stop="handleDelete(node)">
                    <span>{{ t('group_tree.delete') }}</span>
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>
      </template>
    </a-tree>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import {
  CaretDownOutlined,
  CaretRightOutlined,
  EllipsisOutlined,
  PlusOutlined
} from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

defineOptions({
  name: 'GoodsGroupTree'
})

const { t } = useI18n('views.library.goods-library.index')

const props = defineProps({
  groups: {
    type: Array,
    default: () => []
  },
  selectedId: {
    type: [String, Number],
    default: ''
  },
  totalCount: {
    type: Number,
    default: 0
  },
  ungroupedCount: {
    type: Number,
    default: 0
  },
  expandedKeys: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['select', 'create', 'edit', 'delete', 'sort', 'move', 'toggle-expand'])

const openedActionId = ref('')
const dragNodeRef = ref(null)
const isDragging = ref(false)

/* ---------- computed ---------- */

const treeData = computed(() => buildTreeData(props.groups))

const selectedKeys = computed(() => {
  const id = String(props.selectedId)
  if (id === 'all' || id === '0') return []
  return [props.selectedId]
})

/* ---------- helpers ---------- */

function buildTreeData(nodes, level = 0) {
  return nodes.map((node) => {
    const result = {
      key: node.id,
      title: node.group_name,
      level,
      ...node
    }
    if (node.children?.length) {
      result.children = buildTreeData(node.children, level + 1)
    }
    return result
  })
}


function findNodeAndParent(nodes, id, parent = null, siblings = null) {
  for (let i = 0; i < nodes.length; i++) {
    if (String(nodes[i].id) === String(id)) {
      return { node: nodes[i], parent, siblings: siblings || nodes, index: i }
    }
    if (nodes[i].children?.length) {
      const result = findNodeAndParent(nodes[i].children, id, nodes[i], nodes[i].children)
      if (result) return result
    }
  }
  return null
}

function getNodeDepth(node) {
  if (!node.children || !node.children.length) return 0
  return 1 + Math.max(...node.children.map(getNodeDepth))
}

function getNodeLevel(nodes, id, level = 1) {
  for (const node of nodes) {
    if (String(node.id) === String(id)) return level
    if (node.children?.length) {
      const found = getNodeLevel(node.children, id, level + 1)
      if (found) return found
    }
  }
  return 0
}

function isIdInSubtree(node, id) {
  const nodeId = node.key !== undefined ? node.key : node.id
  if (String(nodeId) === String(id)) return true
  if (node.children?.length) {
    return node.children.some((child) => isIdInSubtree(child, id))
  }
  return false
}

/* ---------- selection & expand ---------- */

const isSelected = (id) => String(props.selectedId) === String(id)
const isActionOpen = (id) => openedActionId.value === String(id)

const handleSelectFixed = (id) => {
  emit('select', { id, fixed: true })
}

const handleTreeSelect = (keys, { selected, node }) => {
  if (selected) {
    emit('select', node)
  }
}

const handleExpand = (keys, { node }) => {
  if (isDragging.value) return
  emit('toggle-expand', node.key)
}

/* ---------- action menu ---------- */

const handleActionOpenChange = (id, open) => {
  openedActionId.value = open ? String(id) : ''
}

const resetActionOpen = () => {
  openedActionId.value = ''
}

const handleCreateRoot = () => {
  emit('create', { parentId: 0, parentName: '', level: 0 })
}

const handleCreateChild = (node) => {
  resetActionOpen()
  emit('create', {
    parentId: node.id,
    parentName: node.group_name,
    level: node.level + 1
  })
}

const handleEdit = (node) => {
  resetActionOpen()
  emit('edit', node)
}

const handleDelete = (node) => {
  resetActionOpen()
  emit('delete', node)
}

/* ---------- drag & drop ---------- */

const handleDragStart = ({ node }) => {
  dragNodeRef.value = node
  isDragging.value = true
}

const handleDragEnd = () => {
  isDragging.value = false
  dragNodeRef.value = null
}

const handleAllowDrop = ({ dropNode, dropPosition }) => {
  const dragNode = dragNodeRef.value
  if (!dragNode) return true

  // 无论是拖到节点内部还是节点间隙，都不能落到自己的子孙节点上
  if (String(dragNode.key) !== String(dropNode.key) && isIdInSubtree(dragNode, dropNode.key)) {
    return false
  }

  // 只有拖到节点内容上时，才视为移动到该节点内部并校验层级
  if (dropPosition === 0) {
    // 层级深度校验：目标 level + 1 + 被拖拽子树深度 <= 5（最多5级）
    const dropLevel = getNodeLevel(props.groups, String(dropNode.key))
    const targetLevel = dropLevel + 1
    const dragDepth = getNodeDepth(dragNode)
    if (targetLevel + dragDepth > 5) return false
  }

  return true
}

const handleDrop = (info) => {
  const { dragNode, node, dropToGap } = info
  const dragKey = String(dragNode.key)
  const dropKey = String(node.key)

  // 标准化 dropPosition（官方 demo 方式）
  const dropPos = node.pos.split('-')
  const dropPosition = info.dropPosition - Number(dropPos[dropPos.length - 1])

  // 从原始数据获取被拖拽元素的原始父 id
  const dragResult = findNodeAndParent(props.groups, dragKey)
  if (!dragResult) return
  const oldParentId = dragResult.parent ? String(dragResult.parent.id) : '0'

  // 深拷贝数据（官方 demo 方式：先移除被拖拽元素，再插入）
  const data = JSON.parse(JSON.stringify(props.groups))

  const loop = (list, key, callback) => {
    list.forEach((item, index) => {
      if (String(item.id) === String(key)) {
        return callback(item, index, list)
      }
      if (item.children) {
        return loop(item.children, key, callback)
      }
    })
  }

  // 1. 找到被拖拽元素并从原位置移除
  let dragObj
  loop(data, dragKey, (item, index, arr) => {
    arr.splice(index, 1)
    dragObj = item
  })
  if (!dragObj) return

  let newParentId = '0'
  let targetList = null

  if (!dropToGap) {
    // 情况1：拖到节点内容上 → 成为子节点
    loop(data, dropKey, (item) => {
      item.children = item.children || []
      item.children.unshift(dragObj)
      newParentId = String(item.id)
      targetList = item.children
    })
  } else {
    // 情况2：节点间隙排序，保持与目标节点同级
    let ar = []
    let i = 0
    loop(data, dropKey, (_item, index, arr) => {
      ar = arr
      i = index
    })
    if (dropPosition === -1) {
      ar.splice(i, 0, dragObj)
    } else {
      ar.splice(i + 1, 0, dragObj)
    }
    const dropResult = findNodeAndParent(props.groups, dropKey)
    newParentId = dropResult?.parent ? String(dropResult.parent.id) : '0'
    targetList = ar
  }

  // 生成排序数据
  const sortData = targetList.map((item, index) => ({
    id: item.id,
    sort: index + 1
  }))

  if (newParentId !== oldParentId) {
    emit('move', {
      id: dragKey,
      parentId: newParentId,
      data: sortData,
      expandId: dropToGap ? '' : dropKey
    })
  } else {
    emit('sort', { parentId: newParentId, data: sortData })
  }

  dragNodeRef.value = null
}
</script>

<style lang="less" scoped>
.goods-group-tree {
  display: flex;
  flex-direction: column;
  min-height: 0;

  .tree-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 34px;
    margin-bottom: 4px;
  }

  .tree-title {
    font-size: 16px;
    font-weight: 600;
    color: #262626;
    line-height: 24px;
  }

  .head-action {
    color: #595959;
    width: 24px;
    height: 24px;
    padding: 0;
    border-radius: 6px;
  }

  .fixed-group-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-bottom: 4px;
  }

  .group-node {
    display: flex;
    align-items: center;
    gap: 8px;
    min-height: 34px;
    padding: 5px 8px;
    border-radius: 2px;
    cursor: pointer;
    transition: background-color 0.2s, border-radius 0.2s;

    &:hover {
      background: #f2f4f7;
    }

    &.active {
      background: #e5efff;
      border-radius: 6px;

      .group-name,
      .group-count {
        color: #2475fc;
      }
    }
  }

  .fixed-node {
    padding-left: 36px;

    .group-count {
      display: inline-flex;
    }
  }

  .node-main {
    display: flex;
    align-items: center;
    flex: 1;
    min-width: 0;
    gap: 4px;
  }

  .group-name {
    flex: 1;
    min-width: 0;
    color: #262626;
    font-size: 14px;
    line-height: 22px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .group-count {
    display: inline-flex;
    align-items: center;
    justify-content: flex-end;
    flex-shrink: 0;
    min-width: 24px;
    height: 20px;
    line-height: 20px;
    color: #8c8c8c;
    font-size: 12px;
    text-align: center;
  }

  /* ---------- a-tree overrides ---------- */
  :deep(.ant-tree) {
    .ant-tree-list-holder-inner {
      display: flex;
      flex-direction: column;
      gap: 4px;
    }

    .ant-tree-treenode {
      min-height: 34px;
      padding: 0;
      margin: 0;
      border-radius: 2px;
      transition: background-color 0.2s, border-radius 0.2s;
      display: flex;
      align-items: center;

      &:hover {
        background: #f2f4f7;
        border-radius: 6px;
      }

      &.ant-tree-treenode-selected {
        background: #e5efff;
        border-radius: 6px;

        .node-name,
        .node-count {
          color: #2475fc;
        }
      }

      &.drop-target {
        background: #f0f7ff;
        box-shadow: inset 0 0 0 1px #2475fc;
        border-radius: 6px;
      }

      &.dragging {
        opacity: 0.5;
      }
    }

    .ant-tree-node-content-wrapper {
      flex: 1;
      min-width: 0;
      padding: 5px 8px 5px 4px;
      line-height: 22px;

      &:hover {
        background: transparent;
      }

      &.ant-tree-node-selected {
        background: transparent;
      }
    }

    .ant-tree-switcher {
      width: 24px;
      height: 34px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #8c8c8c;
      flex-shrink: 0;
      line-height: 34px;

      .anticon {
        font-size: 12px;
      }
    }

    .switcher-placeholder {
      display: inline-block;
      width: 24px;
      height: 24px;
    }

    .ant-tree-draggable-icon {
      display: none;
    }

    .ant-tree-indent-unit {
      width: 24px;
    }
  }
}

/* ---------- tree node content ---------- */
.tree-node-content {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
  gap: 4px;

  .node-name {
    flex: 1;
    min-width: 0;
    color: #262626;
    font-size: 14px;
    line-height: 22px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .node-actions {
    display: flex;
    align-items: center;
    gap: 4px;
    flex-shrink: 0;
    min-width: 24px;
    justify-content: flex-end;

    .node-count {
      display: inline-flex;
      align-items: center;
      color: #8c8c8c;
      font-size: 12px;
      min-width: 24px;
      text-align: center;
    }

    .action-btn {
      display: none;
      cursor: pointer;
      width: 24px;
      height: 24px;
      padding: 0;
      border: 0;
      border-radius: 6px;
      color: #595959;
      background: transparent;
      align-items: center;
      justify-content: center;

      &:hover {
        background: #e4e6eb;
      }
    }
  }

  &:hover,
  &.action-open {
    .node-count {
      display: none;
    }

    .action-btn {
      display: inline-flex;
    }
  }
}

/* ---------- dropdown menu ---------- */
:deep(.goods-group-tree-dropdown .ant-dropdown-menu) {
  padding: 2px;
  border-radius: 6px;
  box-shadow:
    0 6px 15px rgba(0, 0, 0, 0.05),
    0 16px 12px rgba(0, 0, 0, 0.04),
    0 8px 5px rgba(0, 0, 0, 0.08);
}

:deep(.goods-group-tree-dropdown .ant-dropdown-menu-item) {
  min-height: auto;
  margin: 0;
  padding: 5px 16px;
  border-radius: 6px;
  color: #595959;
  font-size: 14px;
  line-height: 22px;
}

:deep(.goods-group-tree-dropdown .ant-dropdown-menu-item:hover) {
  background: #f2f4f7;
}

:deep(.goods-group-tree-dropdown .ant-dropdown-menu-item-divider) {
  margin: 2px 6px;
  background: #f0f0f0;
}

:deep(.goods-group-tree-dropdown .ant-dropdown-menu-item.danger-item) {
  color: #fb363f;
}

:deep(.goods-group-tree-dropdown .ant-dropdown-menu-item.danger-item:hover) {
  color: #fb363f;
}
</style>
