<template>
  <a-modal
    :title="t('title_version_history')"
    v-model:open="open"
    width="500px"
    :footer="null"
  >
    <a-table :columns="columns" :data-source="list"></a-table>
  </a-modal>
</template>

<script setup>
import {ref} from 'vue';
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.explore.plugins.components.version-modal');

const open = ref(false)
const list = ref([])
const columns = [
  {
    title: t('label_version'),
    dataIndex: 'version',
    width: 100,
  },
  {
    title: t('label_update_time'),
    dataIndex: 'update_time_text',
    width: 160,
  },
  {
    title: t('label_publisher'),
    dataIndex: 'author',
    width: 100,
  }
]

function show(versions= []) {
  list.value = versions
  open.value = true
}

defineExpose({
  show,
})
</script>

<style scoped>

</style>
