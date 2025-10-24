<template>
  <QaKnowledgeDocument v-if="type == 2" :library_id="default_library_id" />
  <KnowledgeDocument v-if="type == 0" :library_id="default_library_id" />
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import QaKnowledgeDocument from '@/views/library/library-details/qa-knowledge-document/index.vue'
import KnowledgeDocument from '@/views/library/library-details/knowledge-document.vue'
import { getLibraryInfo } from '@/api/library'
import { useRobotStore } from '@/stores/modules/robot'
const robotStore = useRobotStore()
const default_library_id = computed(() => {
  return robotStore.robotInfo.default_library_id
})

const type = ref(null) // 0 普通知识库 2qa知识库

onMounted(() => {
  getLibraryInfo({
    id: default_library_id.value
  }).then((res) => {
    type.value = +res.data.type || 0
  })
})
</script>

<style lang="less" scoped></style>
