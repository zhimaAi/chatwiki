<template>
  <a-dropdown v-model:open="open" :trigger="['click']" placement="bottomRight">
    <button type="button" class="document-more" :class="{ 'card-more': card, 'is-open': open }" :aria-label="commonT('document_more', { name: record.file_name || record.doc_url || '' })" :aria-expanded="open" @click.stop><img :src="card ? icons.cardMore : icons.more" alt="" /></button>
    <template #overlay>
      <a-menu class="workspace-document-actions" @click="({ key }) => choose(key)">
        <a-menu-item v-for="entry in entries" :key="entry.key"><span class="document-menu-entry"><img :src="entry.icon" alt="" />{{ t(entry.label) }}</span></a-menu-item>
      </a-menu>
    </template>
  </a-dropdown>
</template>
<script setup>
import { computed, ref } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import icons from './workspace-document-icons'
const props = defineProps({ record: { type: Object, required: true }, card: Boolean })
const emit = defineEmits(['action'])
const { t } = useI18n('views.library.library-details.knowledge-document')
const { t: commonT } = useI18n('views.customer-service.index')
const open = ref(false)
const entries = computed(() => [
  { key: 'base', icon: icons.base, label: 'btn_base_data' },
  { key: 'preview', icon: icons.preview, label: 'btn_preview' },
  ...(Number(props.record.doc_type) === 2 ? [{ key: 'edit', icon: icons.rename, label: 'btn_edit' }] : []),
  { key: 'rename', icon: icons.rename, label: 'btn_rename' },
  { key: 'group', icon: icons.group, label: 'btn_edit_group' },
  { key: 'download', icon: icons.download, label: 'btn_download_doc' },
  { key: 'delete', icon: icons.remove, label: 'btn_delete' }
])
function choose(key) { open.value = false; emit('action', { key, record: props.record }) }
</script>
<style lang="less" scoped>
.document-more { display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0; width: 32px; height: 32px; padding: 0; border: 0; border-radius: 6px; background: transparent; cursor: pointer; }
.document-more:hover, .document-more.is-open { background: #e8f0fe; }
.document-more:focus-visible { outline: 2px solid #146ef5; outline-offset: -2px; }
.document-more.card-more { position: absolute; top: 12px; right: 12px; width: 24px; height: 24px; opacity: 0; }
.document-more.card-more:hover, .document-more.card-more.is-open { background: #e8f0fe; opacity: 1; }
</style>
<style lang="less">
.workspace-document-actions.ant-dropdown-menu { width: 144px; padding: 4px; border: 1px solid #e8e8e8; border-radius: 12px; box-shadow: 0 8px 24px 8px #0808080a, 0 6px 12px #0808080a, 0 4px 8px -8px #0808080f; }
.workspace-document-actions .ant-dropdown-menu-item { min-height: 32px; padding: 4px 8px; border-radius: 4px; color: #222; font-size: 14px; line-height: 22px; }
.document-menu-entry { display: flex; align-items: center; gap: 8px; white-space: nowrap; }
</style>
