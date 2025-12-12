<template>
  <KnowledgeDocument v-if="type == 0 || type == 3" :type="type"/>
  <QaKnowledgeDocument v-if="type == 2" />
</template>

<script setup>
import { computed } from 'vue'
import KnowledgeDocument from '../knowledge-document.vue'
import QaKnowledgeDocument from '../qa-knowledge-document/index.vue'
import { useLibraryStore } from '@/stores/modules/library'
const libraryStore = useLibraryStore()
const type = computed(() => libraryStore.type)
</script>

<style lang="less" scoped></style>
