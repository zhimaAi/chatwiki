<template>
  <div class="goods-library-page">
    <PageTabs class="mb-16" :tabs="pageTabs" :active="route.path" />

    <div class="goods-library-body">
      <aside class="sidebar-panel">
        <a-spin :spinning="groupLoading">
          <GoodsGroupTree
            :groups="groupTree"
            :selected-id="selectedGroupId"
            :total-count="allGoodsCount"
            :ungrouped-count="ungroupedCount"
            :expanded-keys="expandedKeys"
            @select="handleSelectGroup"
            @create="handleCreateGroup"
            @edit="handleEditGroup"
            @delete="handleDeleteGroup"
            @sort="handleSortGroup"
            @move="handleMoveGroup"
            @toggle-expand="handleToggleExpand"
          />
        </a-spin>
      </aside>

      <section class="main-panel">
        <div class="main-header">
          <div class="main-title">
            <div class="title-text">{{ currentGroupName }}</div>
          </div>

          <GoodsToolbar
            v-model:keyword="keyword"
            @search="handleSearch"
            @import="handleShowImport"
            @export="handleShowExport"
            @create="handleShowGoodsCreate"
          />
        </div>

        <div class="main-content">
          <GoodsTable
            :rows="goodsList"
            :loading="goodsLoading"
            :selected-row-id="selectedGoodsId"
            :active-cell="activeCell"
            :pagination="pagination"
            @select-row="handleSelectRow"
            @hover-cell="handleHoverCell"
            @edit-field="handleOpenFieldEditor"
            @edit-basic-info="handleShowBasicInfoEdit"
            @toggle-status="handleToggleStatus"
            @delete-row="handleDeleteGoods"
            @edit-row="handleShowGoodsEdit"
            @change="handlePageChange"
          />
        </div>
      </section>
    </div>

    <GroupFormModal ref="groupFormRef" @ok="handleGroupSubmit" />
    <GoodsEditorModal
      ref="goodsEditorRef"
      :group-tree-options="groupTreeOptions"
      @submit-goods="handleGoodsSubmit"
      @submit-field="handleFieldSubmit"
    />
    <GoodsImportModal ref="importModalRef" :group-id="selectedGroupId" :group-options="groupOptions" @ok="handleImportDone" />
    <GoodsExportModal ref="exportModalRef" :group-id="selectedGroupId" :keyword="keyword" />
  </div>
</template>

<script setup>
import { computed, createVNode, onMounted, reactive, ref, watch } from 'vue'
import { Modal, message } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import { generateRandomId } from '@/utils/index'
import PageTabs from '@/components/cu-tabs/page-tabs.vue'
import GoodsToolbar from './components/goods-toolbar.vue'
import GoodsGroupTree from './components/goods-group-tree.vue'
import GoodsTable from './components/goods-table.vue'
import GroupFormModal from './components/group-form-modal.vue'
import GoodsEditorModal from './components/goods-editor-modal.vue'
import GoodsImportModal from './components/goods-import-modal.vue'
import GoodsExportModal from './components/goods-export-modal.vue'
import {
  deleteGoods,
  deleteGoodsGroup,
  getGoodsGroupList,
  getGoodsList,
  saveGoods,
  saveGoodsGroup,
  sortGoodsGroup,
  updateGoodsEnabled
} from '@/api/goods-library'

const { t } = useI18n('views.library.goods-library.index')
const route = useRoute()

const pageTabs = computed(() => [
  { title: t('page_tabs.library'), path: '/library/list' },
  { title: t('page_tabs.goods_library'), path: '/library/goods-library' },
  { title: t('page_tabs.database'), path: '/database/list' },
  { title: t('page_tabs.document_extract_faq'), path: '/ai-extract-faq/list' },
  { title: t('page_tabs.trigger_statistics'), path: '/trigger-statics/list' }
])

const groupTree = ref([])
const goodsList = ref([])
const selectedGroupId = ref('all')
const selectedGoodsId = ref('')
const expandedKeys = ref([])
const hasInitializedExpandState = ref(false)
const keyword = ref('')
const activeCell = ref({})
const groupLoading = ref(false)
const goodsLoading = ref(false)
const pagination = reactive({ page: 1, size: 20, total: 0 })
const allGoodsCount = ref(0)
const ungroupedCount = ref(0)

const groupFormRef = ref(null)
const goodsEditorRef = ref(null)
const importModalRef = ref(null)
const exportModalRef = ref(null)

const normalizeId = (value) => {
  if (value === undefined || value === null || value === '') return ''
  return String(value)
}

const toPayloadId = (value) => {
  if (value === undefined || value === null || value === '') return 0
  const numericValue = Number(value)
  return Number.isNaN(numericValue) ? value : numericValue
}

const getApiGroupId = () => {
  const id = normalizeId(selectedGroupId.value)
  if (id === 'all') return -1
  return Number(id) || 0
}

const normalizeGoodsRecord = (item = {}, fallbackGroupId = '0') => ({
  id: normalizeId(item.id ?? item.goods_id ?? generateRandomId(10)),
  goods_id: item.goods_id ?? '',
  goods_name: item.goods_name ?? '',
  category: item.category ?? '',
  brand: item.brand ?? '',
  group_id: normalizeId(item.group_id ?? fallbackGroupId ?? '0'),
  group_name: item.group_name ?? '',
  price: item.price ?? '',
  stock: item.stock ?? '',
  link: item.link ?? '',
  images: Array.isArray(item.images) ? item.images : [],
  description: item.description ?? '',
  qa: item.qa ?? '',
  custom_info: item.custom_info ?? '',
  switch_status: item.switch_status !== undefined ? Number(item.switch_status) : 1
})

const normalizeGroupNode = (item = {}) => ({
  id: normalizeId(item.id ?? item.group_id ?? generateRandomId(8)),
  parent_id: normalizeId(item.parent_id ?? 0),
  group_name: item.group_name ?? item.name ?? '',
  goods_count: Number(item.goods_count ?? 0),
  total_goods_count: Number(item.total_goods_count ?? item.goods_count ?? item.total ?? 0),
  sort: Number(item.sort ?? item.group_sort ?? 0),
  children: []
})

const sortGroupNodes = (nodes = []) => {
  nodes.sort((a, b) => {
    const sortDiff = (a.sort || 0) - (b.sort || 0)
    if (sortDiff !== 0) return sortDiff
    return (a.group_name || '').localeCompare(b.group_name || '', 'zh-Hans-CN')
  })
  nodes.forEach((node) => {
    if (node.children?.length) sortGroupNodes(node.children)
  })
  return nodes
}

const normalizeGroupTree = (list = []) => {
  if (!Array.isArray(list)) return []

  const hasNestedChildren = list.some((item) => Array.isArray(item.children))
  if (hasNestedChildren) {
    const normalized = list.map((item) => ({
      ...normalizeGroupNode(item),
      children: normalizeGroupTree(item.children || [])
    }))
    return sortGroupNodes(normalized)
  }

  const map = new Map()
  list.forEach((item) => {
    const node = normalizeGroupNode(item)
    map.set(node.id, node)
  })

  const roots = []
  map.forEach((node) => {
    const parentId = normalizeId(node.parent_id)
    if (parentId && parentId !== '0' && map.has(parentId)) {
      map.get(parentId).children.push(node)
    } else {
      roots.push(node)
    }
  })
  return sortGroupNodes(roots)
}

const collectExpandableIds = (nodes = [], list = []) => {
  nodes.forEach((node) => {
    if (node.children?.length) {
      list.push(node.id)
      collectExpandableIds(node.children, list)
    }
  })
  return list
}

const resolveExpandedKeys = (nodes = [], currentKeys = [], forceExpandKeys = []) => {
  const expandableKeys = collectExpandableIds(nodes, []).map((id) => normalizeId(id))

  if (!hasInitializedExpandState.value) {
    hasInitializedExpandState.value = true
    return expandableKeys
  }

  const expandableKeySet = new Set(expandableKeys)
  const nextKeys = currentKeys
    .map((id) => normalizeId(id))
    .filter((id, index, list) => expandableKeySet.has(id) && list.indexOf(id) === index)

  forceExpandKeys
    .map((id) => normalizeId(id))
    .forEach((id) => {
      if (expandableKeySet.has(id) && !nextKeys.includes(id)) {
        nextKeys.push(id)
      }
    })

  return nextKeys
}

const findGroupNode = (nodes = [], id) => {
  for (const node of nodes) {
    if (normalizeId(node.id) === normalizeId(id)) return node
    if (node.children?.length) {
      const found = findGroupNode(node.children, id)
      if (found) return found
    }
  }
  return null
}

const flattenGroupOptions = (nodes = [], list = []) => {
  nodes.forEach((node) => {
    list.push({
      label: node.group_name,
      value: node.id
    })
    if (node.children?.length) {
      flattenGroupOptions(node.children, list)
    }
  })
  return list
}

const buildGroupTreeOptions = (nodes = []) => {
  return nodes.map((node) => ({
    title: node.group_name,
    value: node.id,
    key: node.id,
    children: node.children?.length ? buildGroupTreeOptions(node.children) : undefined
  }))
}

const removeGroupNode = (nodes = [], id) => {
  const nextNodes = []
  const removed = []
  nodes.forEach((node) => {
    if (normalizeId(node.id) === normalizeId(id)) {
      const collectIds = (treeNode) => {
        removed.push(normalizeId(treeNode.id))
        if (treeNode.children?.length) {
          treeNode.children.forEach(collectIds)
        }
      }
      collectIds(node)
      return
    }
    let nextNode = node
    if (node.children?.length) {
      const childResult = removeGroupNode(node.children, id)
      nextNode = { ...node, children: childResult.tree }
      removed.push(...childResult.removed)
    }
    nextNodes.push(nextNode)
  })
  return { tree: nextNodes, removed }
}

const loadGroupTree = async (options = {}) => {
  const { forceExpandKeys = [] } = options
  groupLoading.value = true
  try {
    const res = await getGoodsGroupList()
    const data = res?.data || {}
    groupTree.value = normalizeGroupTree(data.list || [])
    allGoodsCount.value = Number(data.total ?? 0)
    ungroupedCount.value = Number(data.ungrouped_count ?? 0)
    expandedKeys.value = resolveExpandedKeys(groupTree.value, expandedKeys.value, forceExpandKeys)
  } catch (error) {
    console.warn(error)
  }
  groupLoading.value = false
}

const loadGoodsList = async () => {
  goodsLoading.value = true
  try {
    const res = await getGoodsList({
      group_id: getApiGroupId(),
      keyword: keyword.value,
      switch_status: -1,
      page: pagination.page,
      size: pagination.size
    })
    const data = res?.data || {}
    goodsList.value = (data.list || []).map((item) => normalizeGoodsRecord(item))
    pagination.total = Number(data.total ?? 0)
  } catch (error) {
    console.warn(error)
    goodsList.value = []
    pagination.total = 0
  }
  goodsLoading.value = false
}

const currentGroupName = computed(() => {
  if (normalizeId(selectedGroupId.value) === 'all') return t('group_tree.all')
  if (normalizeId(selectedGroupId.value) === '0') return t('group_tree.ungrouped')
  const node = findGroupNode(groupTree.value, selectedGroupId.value)
  return node?.group_name || t('group_tree.ungrouped')
})

const groupOptions = computed(() => ([
  {
    label: t('group_tree.ungrouped'),
    value: '0'
  },
  ...flattenGroupOptions(groupTree.value).filter((item) => normalizeId(item.value) !== '0')
]))

const groupTreeOptions = computed(() => ([
  {
    title: t('group_tree.ungrouped'),
    value: '0',
    key: '0'
  },
  ...buildGroupTreeOptions(groupTree.value).filter((item) => normalizeId(item.value) !== '0')
]))

watch(goodsList, (list) => {
  if (!list.length) {
    selectedGoodsId.value = ''
    activeCell.value = {}
    return
  }
  const exists = list.some((item) => normalizeId(item.id) === normalizeId(selectedGoodsId.value))
  if (!exists) {
    selectedGoodsId.value = normalizeId(list[0].id)
  }
}, { immediate: true })

watch(selectedGoodsId, () => {
  activeCell.value = {}
})

watch(selectedGroupId, () => {
  activeCell.value = {}
})

const handleSelectGroup = (payload) => {
  selectedGroupId.value = normalizeId(payload.id)
  pagination.page = 1
  loadGoodsList()
}

const handleSearch = (value) => {
  keyword.value = value || ''
  pagination.page = 1
  loadGoodsList()
}

const handlePageChange = ({ page, size }) => {
  pagination.page = page
  pagination.size = size
  loadGoodsList()
}

const handleSelectRow = (record) => {
  selectedGoodsId.value = normalizeId(record.id)
}

const handleHoverCell = (payload) => {
  activeCell.value = payload
}

const handleShowGoodsCreate = () => {
  const data = {}
  if (selectedGroupId.value !== 'all') {
    data.group_id = selectedGroupId.value
  }
  goodsEditorRef.value?.show({
    name: 'goods_full',
    data
  })
}

const handleShowGoodsEdit = (record) => {
  goodsEditorRef.value?.show({ name: 'goods_full', data: record })
}

const handleShowBasicInfoEdit = (record) => {
  goodsEditorRef.value?.show({ name: 'goods_basic', data: record })
}

const fieldMetaMap = {
  images: { name: 'goods_images' },
  description: { name: 'goods_description' },
  qa: { name: 'goods_qa' },
  custom_info: { name: 'goods_custom_info' }
}

const getFieldValue = (row, fieldKey) => {
  switch (fieldKey) {
    case 'images':
      return row.images || []
    case 'description':
      return row.description || ''
    case 'qa':
      return row.qa || ''
    case 'custom_info':
      return row.custom_info || ''
    default:
      return row[fieldKey] || ''
  }
}

const handleOpenFieldEditor = ({ row, fieldKey }) => {
  if (fieldKey === 'basic_info') {
    handleShowBasicInfoEdit(row)
    return
  }
  const fieldConfig = fieldMetaMap[fieldKey] || {}
  goodsEditorRef.value?.show({
    name: fieldConfig.name,
    row,
    value: getFieldValue(row, fieldKey)
  })
}

const handleFieldSubmit = async ({ fieldKey, value, row }, actions) => {
  const rowId = normalizeId(row?.id)
  if (!rowId) {
    actions?.setSubmitting(false)
    return
  }

  const currentRecord = goodsList.value.find((item) => normalizeId(item.id) === rowId)
  if (!currentRecord) {
    actions?.setSubmitting(false)
    return
  }

  const record = {
    ...currentRecord,
    id: toPayloadId(rowId),
    group_id: toPayloadId(currentRecord.group_id),
    [fieldKey]: value
  }

  try {
    await saveGoods(record)
    await loadGoodsList()
    message.success(t('message.save_success'))
    actions?.close()
  } catch (error) {
    console.warn(error)
  } finally {
    actions?.setSubmitting(false)
  }
}

const handleToggleStatus = async ({ row, checked }) => {
  const rowId = normalizeId(row?.id)
  if (!rowId) return

  try {
    await updateGoodsEnabled({
      id: toPayloadId(rowId),
      switch_status: checked ? 1 : 0
    })
    await loadGoodsList()
    await loadGroupTree()
    message.success(t('message.status_update_success'))
  } catch (error) {
    console.warn(error)
    await loadGoodsList()
    await loadGroupTree()
  }
}

const handleDeleteGoods = (row) => {
  const rowId = normalizeId(row?.id)
  if (!rowId) return

  Modal.confirm({
    title: t('confirm.delete_goods_title'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirm.delete_goods_content', { name: row.goods_name || row.goods_id || row.id }),
    okText: t('confirm.confirm_btn'),
    okType: 'danger',
    cancelText: t('confirm.cancel_btn'),
    onOk: async () => {
      try {
        await deleteGoods({ id: toPayloadId(rowId) })
        await loadGoodsList()
        await loadGroupTree()
        message.success(t('message.delete_success'))
      } catch (error) {
        console.warn(error)
      }
    }
  })
}

const handleGoodsSubmit = async (payload, actions) => {
  const id = normalizeId(payload.id)
  const currentRecord = goodsList.value.find((item) => normalizeId(item.id) === id)

  const record = {
    ...(currentRecord || {}),
    ...payload,
    id: id ? toPayloadId(id) : 0,
    group_id: toPayloadId(payload.group_id || 0),
    switch_status: currentRecord?.switch_status ?? 1
  }

  try {
    await saveGoods(record)
    await loadGoodsList()
    await loadGroupTree()
    message.success(t('message.save_success'))
    actions?.close()
  } catch (error) {
    console.warn(error)
  } finally {
    actions?.setSubmitting(false)
  }
}

const buildParentTreeOptions = (excludeId = '') => {
  let tree = groupTree.value
  if (excludeId) {
    const result = removeGroupNode(groupTree.value, excludeId)
    tree = result.tree
  }
  return [
    { title: t('group_modal.root_group'), value: '0', key: '0' },
    ...buildGroupTreeOptions(tree)
  ]
}

const handleCreateGroup = (payload) => {
  groupFormRef.value?.show({ ...payload, parentTreeOptions: buildParentTreeOptions() })
}

const handleEditGroup = (payload) => {
  const excludeId = normalizeId(payload.id)
  groupFormRef.value?.show({ ...payload, parentTreeOptions: buildParentTreeOptions(excludeId) })
}

const handleGroupSubmit = async (payload, actions) => {
  const groupId = normalizeId(payload.id)
  const parentId = normalizeId(payload.parent_id ?? '0')

  const nextPayload = {
    ...payload,
    id: groupId ? toPayloadId(groupId) : 0,
    parent_id: toPayloadId(parentId),
    group_name: payload.group_name
  }

  try {
    await saveGoodsGroup(nextPayload)
    await loadGroupTree()
    message.success(groupId ? t('group_modal.success_edit') : t('group_modal.success_add'))
    actions?.close()
  } catch (error) {
    console.warn(error)
  } finally {
    actions?.setSubmitting(false)
  }
}

const handleDeleteGroup = (payload) => {
  const groupId = normalizeId(payload?.id)
  if (!groupId) return

  Modal.confirm({
    title: t('confirm.delete_group_title'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirm.delete_group_content', { name: payload.group_name }),
    okText: t('confirm.confirm_btn'),
    okType: 'danger',
    cancelText: t('confirm.cancel_btn'),
    onOk: async () => {
      const removalResult = removeGroupNode(groupTree.value, groupId)
      groupTree.value = removalResult.tree

      if (removalResult.removed.includes(normalizeId(selectedGroupId.value))) {
        selectedGroupId.value = 'all'
      }

      expandedKeys.value = resolveExpandedKeys(groupTree.value, expandedKeys.value)

      try {
        await deleteGoodsGroup({ id: toPayloadId(groupId) })
        await loadGroupTree()
        await loadGoodsList()
        message.success(t('message.delete_success'))
      } catch (error) {
        console.warn(error)
        await loadGroupTree()
        await loadGoodsList()
      }
    }
  })
}

const handleSortGroup = async ({ data }) => {
  try {
    await sortGoodsGroup({
      data: data.map((item) => ({
        id: toPayloadId(item.id),
        sort: item.sort
      }))
    })
  } catch (error) {
    console.warn(error)
  }
  await loadGroupTree()
}

const handleMoveGroup = async ({ id, parentId, data, expandId }) => {
  const node = findGroupNode(groupTree.value, id)
  if (!node) {
    await loadGroupTree()
    return
  }

  try {
    await saveGoodsGroup({
      id: toPayloadId(id),
      parent_id: toPayloadId(parentId),
      group_name: node.group_name
    })
    await sortGoodsGroup({
      data: data.map((item) => ({
        id: toPayloadId(item.id),
        sort: item.sort
      }))
    })
    message.success(t('message.move_success'))
  } catch (error) {
    console.warn(error)
  }
  await loadGroupTree({ forceExpandKeys: expandId ? [expandId] : [] })
}

const handleToggleExpand = (id) => {
  const nextId = normalizeId(id)
  if (expandedKeys.value.includes(nextId)) {
    expandedKeys.value = expandedKeys.value.filter((item) => item !== nextId)
    return
  }
  expandedKeys.value = expandedKeys.value.concat(nextId)
}

const handleShowImport = () => {
  importModalRef.value?.show()
}

const handleShowExport = () => {
  exportModalRef.value?.show()
}

const handleImportDone = async () => {
  await loadGroupTree()
  await loadGoodsList()
}

onMounted(() => {
  loadGroupTree()
  loadGoodsList()
})
</script>

<style lang="less" scoped>
.goods-library-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding-bottom: 24px;
  overflow: hidden;
  background: #fff;

  .goods-library-body {
    display: flex;
    flex: 1;
    height: 100%;
    background: #fff;
    overflow: hidden;
  }

  .sidebar-panel {
    width: 247px;
    height: 100%;
    padding: 0 16px 16px 0;
    background: #fff;
    border-right: 1px solid #f0f0f0;
    overflow: hidden;
    overflow-y: auto;
  }

  .main-panel {
    display: flex;
    flex: 1;
    height: 100%;
    padding: 0 24px;
    flex-direction: column;
    overflow: hidden;
  }

  .main-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    margin-bottom: 12px;
  }

  .main-title {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .title-text {
    font-size: 16px;
    font-weight: 600;
    color: #262626;
    line-height: 24px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .main-content {
    flex: 1;
    background: #fff;
    overflow: hidden;
  }
}
</style>
