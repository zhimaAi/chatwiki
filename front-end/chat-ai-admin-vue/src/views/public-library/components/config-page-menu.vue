<template>
  <div class="page-menus-box">
    <a-menu :selectedKeys="current" mode="horizontal" :items="menus" @click="handleMenu" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()

const current = computed(() => [route.name])

// const current = ref(['permissions'])
const menus = ref([
  {
    key: 'PublicLibraryConfig',
    label: '知识库设置',
    title: '知识库设置',
    path: '/public-library/config'
  },
  {
    key: 'PublicLibraryPermissions',
    label: '访问权限',
    title: '访问权限',
    path: '/public-library/permissions'
  },
  {
    key: 'PublicLibraryAi',
    label: '文档AI',
    title: '文档AI',
    path: '/public-library/ai'
  },
  {
    key: 'PublicLibraryWebStatistics',
    label: '统计设置',
    title: '统计设置',
    path: '/public-library/web-statistics'
  }
])

const handleMenu = ({ item }) => {
  router.push({
    path: item.path,
    query: {
      library_id: route.query.library_id,
      library_key: route.query.library_key
    }
  })
}
</script>

<style lang="less" scoped>
.page-menus-box {
  margin-bottom: 16px;
  padding-left: 24px;
}
</style>
