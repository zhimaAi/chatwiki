<template>
  <a-modal
    v-model:open="open"
    :title="t('title')"
    :width="480"
    :mask-closable="false"
    centered
  >
    <div class="scope-modal-body">
      <a-radio-group v-model:value="scopeMode" class="scope-radio-group">
        <a-radio value="all">{{ t('radio_scope_all') }}</a-radio>
        <a-radio value="partial">{{ t('radio_scope_partial') }}</a-radio>
      </a-radio-group>

      <div v-if="scopeMode === 'partial'" class="scope-tree-wrapper">
        <a-spin :spinning="groupLoading">
          <a-tree
            :checked-keys="checkedGroupKeys"
            :tree-data="groupTreeData"
            checkable
            :check-strictly="true"
            :expanded-keys="expandedGroupKeys"
            :selectable="false"
            class="scope-tree"
            @check="handleTreeCheck"
            @expand="handleTreeExpand"
          >
            <template #title="node">
              <div class="scope-tree-node">
                <span class="scope-node-name">{{ node.title }}</span>
                <span v-if="node.count != null" class="scope-node-count">{{ node.count }}</span>
              </div>
            </template>
          </a-tree>
        </a-spin>
      </div>
    </div>

    <template #footer>
      <a-button @click="handleClose">{{ t('btn_cancel') }}</a-button>
      <a-button type="primary" :loading="scopeSaving" @click="handleSave">{{ t('btn_confirm') }}</a-button>
    </template>
  </a-modal>
</template>

<script setup>
import { ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { storeToRefs } from 'pinia'
import { getGoodsGroupList } from '@/api/goods-library'

const { t } = useI18n('views.clawbot.skills.components.GoodsRecommendScopeModal')

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'confirm'])

const clawbotStore = useClawbotStore()
const { robotInfo } = storeToRefs(clawbotStore)
const { updateClawbotConf } = clawbotStore

const open = ref(false)
const scopeMode = ref('all')
const groupTreeData = ref([])
const groupLoading = ref(false)
const checkedGroupKeys = ref([])
const expandedGroupKeys = ref([])
const scopeSaving = ref(false)

watch(
  () => props.visible,
  (val) => {
    open.value = val
    if (val) {
      initForm()
    }
  }
)

watch(open, (val) => {
  if (!val) {
    emit('update:visible', false)
  }
})

const convertGroupNode = (item) => ({
  key: String(item.id ?? item.group_id ?? ''),
  title: item.group_name ?? item.name ?? '',
  count: Number(item.goods_count ?? 0),
  children: (item.children || []).map(convertGroupNode)
})

const buildGroupTreeData = (list, ungroupedCount) => {
  const tree = list.map(convertGroupNode)
  tree.unshift({
    key: '0',
    title: t('label_ungrouped'),
    count: ungroupedCount,
    children: []
  })
  return tree
}

const collectAllKeys = (nodes) => {
  const keys = []
  const walk = (list) => {
    list.forEach((node) => {
      if (node.children && node.children.length) {
        keys.push(node.key)
        walk(node.children)
      }
    })
  }
  walk(nodes)
  return keys
}

const loadGoodsGroupList = async () => {
  groupLoading.value = true
  try {
    const res = await getGoodsGroupList()
    const data = res?.data || {}
    groupTreeData.value = buildGroupTreeData(data.list || [], Number(data.ungrouped_count ?? 0))
    expandedGroupKeys.value = collectAllKeys(groupTreeData.value)
  } catch (err) {
    console.error(t('msg_load_groups_failed'), err)
    message.error(t('msg_load_groups_failed'))
  } finally {
    groupLoading.value = false
  }
}

const initForm = () => {
  const savedIds = (robotInfo.value?.goods_lib_recommend_group_ids || '').split(',').filter(Boolean)
  if (savedIds.length === 0) {
    scopeMode.value = 'all'
    checkedGroupKeys.value = []
  } else {
    scopeMode.value = 'partial'
    checkedGroupKeys.value = savedIds
  }
  if (!groupTreeData.value.length) {
    loadGoodsGroupList()
  }
}

const handleTreeCheck = (checkedKeys) => {
  checkedGroupKeys.value = Array.isArray(checkedKeys) ? checkedKeys : (checkedKeys?.checked || [])
}

const handleTreeExpand = (keys) => {
  expandedGroupKeys.value = keys
}

const buildConfPayload = (overrides = {}) => {
  return {
    id: robotInfo.value?.id,
    search_knowledge_close: !Number(robotInfo.value?.search_knowledge_close || 0) ? 0 : 1,
    query_local_docs_close: !Number(robotInfo.value?.query_local_docs_close || 0) ? 0 : 1,
    open_agent_write_file_tool: Number(robotInfo.value?.open_agent_write_file_tool || 0) === 1 ? 1 : 0,
    open_agent_execute_tool: Number(robotInfo.value?.open_agent_execute_tool || 0) === 1 ? 1 : 0,
    open_agent_edit_file_tool: Number(robotInfo.value?.open_agent_edit_file_tool || 0) === 1 ? 1 : 0,
    goods_lib_recommend_switch: Number(robotInfo.value?.goods_lib_recommend_switch || 0) === 1 ? 1 : 0,
    goods_lib_recommend_group_ids: robotInfo.value?.goods_lib_recommend_group_ids || '',
    ...overrides
  }
}

const handleSave = async () => {
  const ids = scopeMode.value === 'all' ? '' : checkedGroupKeys.value.join(',')

  scopeSaving.value = true
  try {
    await updateClawbotConf(buildConfPayload({ goods_lib_recommend_group_ids: ids }))
    message.success(t('msg_saved'))
    emit('confirm', { groupIds: ids })
    handleClose()
  } catch (err) {
    message.error(t('msg_save_failed'))
  } finally {
    scopeSaving.value = false
  }
}

const handleClose = () => {
  open.value = false
}
</script>

<style lang="less" scoped>
.scope-modal-body {
  .scope-radio-group {
    display: flex;
    gap: 24px;
    margin-bottom: 16px;
  }

  .scope-tree-wrapper {
    max-height: 360px;
    overflow-y: auto;
    border: 1px solid #f0f0f0;
    border-radius: 8px;
    padding: 8px;
  }

  .scope-tree {
    :deep(.ant-tree-node-content-wrapper) {
      flex: 1;
    }

    .scope-tree-node {
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
      padding-right: 8px;

      .scope-node-name {
        font-size: 14px;
        color: #262626;
      }

      .scope-node-count {
        font-size: 12px;
        color: #8c8c8c;
        margin-left: 8px;
      }
    }
  }
}
</style>
